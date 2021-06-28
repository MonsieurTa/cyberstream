package handler

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/MonsieurTa/hypertube/common/tcp"
	"github.com/MonsieurTa/hypertube/pkg/transcoder/internal/hls"
)

func TCPHandler(conn net.Conn) error {
	header, err := readHeader(conn)
	if err != nil {
		return err
	}
	log.Printf("transcoder: received header\n%v\n", header)
	c := hls.NewHLSConverter(&hls.Config{
		Reader:      conn,
		Length:      int64(header.FileSize),
		PieceLength: int64(header.PieceLength),
		OutputDir:   string(header.DirName),
	})

	go c.Convert()

	go func() {
		c.WaitUntilDone()
		c.Close()
		conn.Close()
	}()
	return nil
}

func readHeader(conn net.Conn) (*tcp.Header, error) {
	buf := make([]byte, 24)

	at := 0
	for at < 24 {
		n, err := conn.Read(buf[at:])
		if err != nil {
			return nil, err
		}
		at += n
	}

	fileSize, pieceLength, dirNameSize := toUInt64(buf[0:8]), toUInt64(buf[8:16]), toUInt64(buf[16:24])

	buf = make([]byte, dirNameSize)
	at = 0
	for at < int(dirNameSize) {
		n, err := conn.Read(buf[at:])
		if err != nil {
			return nil, err
		}
		at += n
	}

	return &tcp.Header{
		FileSize:    fileSize,
		PieceLength: pieceLength,
		DirNameSize: dirNameSize,
		DirName:     buf,
	}, nil
}

func toUInt64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
