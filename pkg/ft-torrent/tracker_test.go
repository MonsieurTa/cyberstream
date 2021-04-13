package torrent

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebianTracker(t *testing.T) {
	tfile, err := ParseFromFile("/home/wta/Projects/hypertube/debian.torrent")
	assert.Nil(t, err)

	peerID := [20]byte{}
	_, err = rand.Read(peerID[:])
	assert.Nil(t, err)

	trackers, err := tfile.Trackers()
	assert.Nil(t, err)

	if len(trackers) != 0 {
		t.Errorf("expected empty trackers, got %d len", len(trackers))
	}
}

func TestAnimeTracker(t *testing.T) {
	tfile, err := ParseFromFile("/home/wta/Projects/hypertube/test.mkv.torrent")
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
