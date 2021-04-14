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

	tfile, err := ParseFromFile(wdPath + `/debian.torrent`)
	assert.Nil(t, err)

	peerID := [20]byte{}
	_, err = rand.Read(peerID[:])
	assert.Nil(t, err)

	trackers, err := tfile.Trackers()
	assert.Nil(t, err)

	if len(trackers) != 1 {
		t.Errorf("expected default tracker, got %d len", len(trackers))
	}
}

func TestAnimeTracker(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	tfile, err := ParseFromFile(wdPath + `/test.mkv.torrent`)
	assert.Nil(t, err)

	peerID := [20]byte{}
	_, err = rand.Read(peerID[:])
	assert.Nil(t, err)

	trackers, err := tfile.Trackers()
	assert.Nil(t, err)

	if len(trackers) != 16 {
		t.Errorf("expected trackers of len 16, got %d", len(trackers))
	}
}
