package transcoder

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/MonsieurTa/hypertube/common/tcp"
	"github.com/anacrolix/torrent"
)

const (
	TCP_START_PORT      = 1024
	TCP_MAX_BUFFER_SIZE = 8096000
)

type transcodeClient struct {
	cfg          *Config
	ongoing      safeCounter
	tcpConnPorts map[int]bool
}

type safeCounter struct {
	mu sync.Mutex
	v  int
}

func NewSafeCounter(start int) safeCounter {
	return safeCounter{v: start}
}

func (sc *safeCounter) Value() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.v
}

func (sc *safeCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.v++
}

func (sc *safeCounter) Decrement() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.v--
}

type Config struct {
	CoreNb int
	Url    string
	Port   string
}

func NewService(cfg *Config) UseCase {
	return &transcodeClient{
		cfg:          cfg,
		ongoing:      NewSafeCounter(0),
		tcpConnPorts: make(map[int]bool),
	}
}

type TranscoderParams struct {
	Reader      torrent.Reader
	FileSize    int64
	PieceLength int64
	DirName     string
}

type TranscodeResponse struct {
	PlaylistURL string
}

func (tp *TranscoderParams) header() *tcp.Header {
	return &tcp.Header{
		FileSize:    uint64(tp.FileSize),
		PieceLength: uint64(tp.PieceLength),
		DirNameSize: uint64(len(tp.DirName)),
		DirName:     []byte(tp.DirName),
	}
}

func (tc *transcodeClient) Transcode(tp *TranscoderParams) error {
	if tc.ongoing.Value() == tc.cfg.CoreNb {
		return errors.New("busy")
	}

	url := tc.cfg.Url + ":" + tc.cfg.Port
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return err
	}

	log.Printf("tcp connection established with %s\n", url)
	go func() {
		tc.ongoing.Increment()
		defer tc.ongoing.Decrement()
		defer conn.Close()

		header := tp.header()

		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, header.FileSize)
		conn.Write(b)
		binary.BigEndian.PutUint64(b, header.PieceLength)
		conn.Write(b)

		binary.BigEndian.PutUint64(b, header.DirNameSize)
		conn.Write(b)
		conn.Write(header.DirName)

		at := 0
		end := int(tp.FileSize)
		buf := make([]byte, TCP_MAX_BUFFER_SIZE)
		r := tp.Reader
		for at < end {
			n, _ := r.Read(buf)

			conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			_, err := conn.Write(buf)
			if err != nil {
				break
			}
			at += n
		}
		log.Printf("media: %d bytes sent\n", tp.FileSize)
	}()
	return nil
}
