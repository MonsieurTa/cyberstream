package torrent

import (
	"io"
	"os"
)

func Parse(r io.Reader) ([]Torrent, error) {
	bto, err := Read(r)
	if err != nil {
		return nil, err
	}
	return bto.toTorrent()
}

func ParseFromFile(filepath string) ([]Torrent, error) {
	r, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return Parse(r)
}
