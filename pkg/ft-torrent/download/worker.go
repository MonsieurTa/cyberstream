package download

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/client"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/message"
)

type worker struct {
	parent *Downloader
	client *client.Client
	state  pieceProgress
}

type pieceProgress struct {
	index      int
	downloaded int
	requested  int
	backlog    int
	output     []byte
}

// MaxBlockSize is the largest number of bytes a request can ask for
const MaxBlockSize = 16384

// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
const MaxBacklog = 5

var (
	err_wrong_hash = func(expected, got [20]byte) error { return fmt.Errorf("expected %x, got %x", expected, got) }

	log_piece_request_failed   = func(index int) { log.Printf("piece#%d: request failed", index) }
	log_failed_integrity_check = func(index int) { log.Printf("piece#%d: failed integrity check\n", index) }
)

func (w *worker) start(workQueue chan *pieceWork, results chan *pieceResult) {
	c := w.client

	if err := c.Ask(w.parent.t.InfoHash()); err != nil {
		log.Printf("%v\n", err)
		return
	}
	defer c.Close()

	c.SendUnchoke()
	c.SendInterested()

	for pw := range workQueue {
		if !c.Bitfield().HasPiece(pw.index) {
			workQueue <- pw
			return
		}

		data, err := w.download(pw)
		if err != nil {
			log.Println(err)
			workQueue <- pw
			return
		}

		err = checkIntegrity(data, pw.hash)
		if err != nil {
			log_failed_integrity_check(pw.index)
			workQueue <- pw
			continue
		}

		c.SendHave(pw.index)
		results <- &pieceResult{pw.index, data}
	}
}

func (w *worker) download(pw *pieceWork) ([]byte, error) {
	w.state = pieceProgress{
		index:  pw.index,
		output: make([]byte, pw.length),
	}

	c := w.client
	s := &w.state

	c.Conn().SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn().SetDeadline(time.Now())

	for s.downloaded < pw.length {
		if !c.State.AmChoking {
			for s.backlog < MaxBacklog {
				err := requestBlock(c, pw, s)
				if err != nil {
					log_piece_request_failed(pw.index)
				}
			}
		}
		err := w.readMessage()
		if err != nil {
			return nil, err
		}
	}
	return s.output, nil
}

func requestBlock(c *client.Client, pw *pieceWork, state *pieceProgress) error {
	blockSize := MaxBlockSize
	// Last block might be shorter than the typical block
	if pw.length-state.requested < blockSize {
		blockSize = pw.length - state.requested
	}

	err := c.SendRequest(pw.index, state.requested, blockSize)
	if err != nil {
		return err
	}
	state.backlog++
	state.requested += blockSize
	return nil
}

func (w *worker) readMessage() error {
	msg, err := w.client.Read()
	if err != nil {
		return err
	}

	if msg == nil {
		return nil
	}

	s := &w.state
	switch msg.ID() {
	case message.UNCHOKE:
		w.client.State.AmChoking = false
	case message.CHOKE:
		w.client.State.AmChoking = true
	case message.HAVE:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}

		w.client.Bitfield().SetPiece(index)
	case message.PIECE:
		n, err := message.ParsePiece(s.index, s.output, msg)
		if err != nil {
			return err
		}
		s.downloaded += n
		s.backlog--
	}
	return nil
}

func checkIntegrity(data []byte, expectedHash [20]byte) error {
	got := sha1.Sum(data)

	if !bytes.Equal(got[:], expectedHash[:]) {
		return err_wrong_hash(expectedHash, got)
	}
	return nil
}
