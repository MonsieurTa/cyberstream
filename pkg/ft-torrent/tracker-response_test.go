package torrent

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackerResponse(t *testing.T) {
	data := map[string]interface{}{
		"interval": int64(900),
		"peers":    []byte{},
	}

	trResp := NewTrackerResponse(data)
	fmt.Printf("%v\n", trResp)
}

func TestPeersRequest(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	tfiles, err := ParseFromFile(wdPath + `/debian.torrent`)
	// tfiles, err := ParseFromFile(wdPath + `/test.mkv.torrent`)
	assert.Nil(t, err)

	for _, tfile := range tfiles {
		trackers, err := tfile.Trackers()
		assert.Nil(t, err)

		for _, tr := range trackers {
			trResp, err := tr.RequestPeers(tfile.InfoHash)
			if err != nil {
				fmt.Printf("could not request peers from %v\n", tr.announce)
				continue
			}
			_, err = trResp.Peers()
			assert.Nil(t, err)
		}
	}

}
