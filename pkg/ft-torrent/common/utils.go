package common

import "crypto/rand"

type Bitfield []byte

func GeneratePeerID() ([20]byte, error) {
	peerID := [20]byte{}
	_, err := rand.Read(peerID[:])
	if err != nil {
		return [20]byte{}, err
	}
	return peerID, nil
}

func (b Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	return b[byteIndex]>>(7-offset)&1 != 0
}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	bf[byteIndex] |= 1 << (7 - offset)
}
