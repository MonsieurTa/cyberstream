package torrent

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebianTracker(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	tfiles, err := ParseFromFile(wdPath + `/debian.torrent`)
	assert.Nil(t, err)

	peerID := [20]byte{}
	_, err = rand.Read(peerID[:])
	assert.Nil(t, err)

	for _, tfile := range tfiles {
		_, err := tfile.Trackers()
		assert.Nil(t, err)
	}
}

func TestAnimeTracker(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	tfiles, err := ParseFromFile(wdPath + `/test.mkv.torrent`)
	assert.Nil(t, err)

	peerID := [20]byte{}
	_, err = rand.Read(peerID[:])
	assert.Nil(t, err)

	for _, tfile := range tfiles {
		_, err := tfile.Trackers()
		assert.Nil(t, err)
	}

}
