package torrent

import (
	"errors"
	"io"
)

type Handshake struct {
	pstr     string
	infoHash [20]byte
	peerID   [20]byte
}

const BITTORENT_PROTOCOL = "BitTorrent protocol"

var (
	err_invalid_pstrlen = errors.New("invalid pstr len")
)

func (h Handshake) Serialize() []byte {
	offset := 0

	totalBufSize := len(h.pstr) + 49
	buf := make([]byte, totalBufSize)

	pstrLen := []byte{byte(len(h.pstr))}
	offset += copy(buf[offset:], pstrLen)
	offset += copy(buf[offset:], make([]byte, 8))
	offset += copy(buf[offset:], h.infoHash[:])
	offset += copy(buf[offset:], h.peerID[:])
	return buf
}

func ReadHandshake(r io.Reader) (Handshake, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return Handshake{}, err
	}

	pstrLen := int(buf[0])
	if pstrLen == 0 {
		return Handshake{}, err_invalid_pstrlen
	}

	hsBuf := make([]byte, pstrLen+48)
	_, err = io.ReadFull(r, hsBuf)
	if err != nil {
		return Handshake{}, nil
	}

	var infoHash, peerID [20]byte

	offset := pstrLen + 8
	copy(infoHash[:], hsBuf[offset:offset+20])
	copy(peerID[:], hsBuf[offset+20:])

	rv := Handshake{
		pstr:     string(hsBuf[:pstrLen]),
		infoHash: infoHash,
		peerID:   peerID,
	}
	return rv, nil
}
