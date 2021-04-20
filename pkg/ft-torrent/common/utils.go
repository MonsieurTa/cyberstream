package common

import "crypto/rand"

func GeneratePeerID() ([20]byte, error) {
	peerID := [20]byte{}
	_, err := rand.Read(peerID[:])
	if err != nil {
		return [20]byte{}, err
	}
	return peerID, nil
}
