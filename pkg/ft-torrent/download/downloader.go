package download

import (
	"fmt"

	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/client"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/common"
)

type Downloader struct {
	peerID [20]byte
	t      Torrent

	workQueue chan *pieceWork
	results   chan *pieceResult
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

func NewDownloader(t Torrent) (*Downloader, error) {
	peerID, err := common.GeneratePeerID()
	if err != nil {
		return nil, err
	}
	return &Downloader{
		peerID:    peerID,
		t:         t,
		workQueue: make(chan *pieceWork, t.Pieces()),
		results:   make(chan *pieceResult),
	}, nil
}

func (d *Downloader) Download() error {
	peers, err := d.t.Peers(d.peerID)
	if err != nil {
		return err
	}

	d.loadWorkQueue()

	for _, peer := range peers {
		worker := d.newWorker(peer)
		go worker.start(d.workQueue, d.results)
	}
	
	buf := make([]byte, d.t.Length())
	nbPieces := len(d.t.PieceHashes())
	donePieces := 0
	for donePieces < nbPieces {
		res := <-d.results
		begin, end := d.t.PieceBounds(res.index)
		copy(buf[begin:end], res.buf)
		donePieces++
		fmt.Printf("Progress(%d%%): piece#%d downloaded\n", donePieces*100/nbPieces, res.index)
	}
	close(d.workQueue)
	return nil
}

func (d *Downloader) newWorker(peer common.Peer) *worker {
	client := client.NewClient(peer, d.peerID)
	return &worker{
		parent: d,
		client: client,
	}
}

func (d *Downloader) loadWorkQueue() {
	hashes := d.t.PieceHashes()
	for index, hash := range hashes {
		length := d.t.PieceSize(index)
		d.workQueue <- &pieceWork{index, hash, length}
	}
}
