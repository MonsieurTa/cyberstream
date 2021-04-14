package torrent

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeersRequest(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	// tfile, err := ParseFromFile(wdPath + `/debian.torrent`)
	tfile, err := ParseFromFile(wdPath + `/test.mkv.torrent`)
	assert.Nil(t, err)

	peerID := [20]byte{}
	_, err = rand.Read(peerID[:])
	assert.Nil(t, err)

	trackers, err := tfile.Trackers()
	assert.Nil(t, err)

	for _, tr := range trackers {
		trResp, err := tr.RequestPeers(tfile.InfoHash)
		if err != nil {
			fmt.Printf("could not request peers from %v\n", tr.announce)
			continue
		}

		peers, err := trResp.Peers()
		assert.Nil(t, err)
		for _, peer := range peers {
			fmt.Printf("%s\n", peer.String())
		}
		fmt.Println()
	}
}
