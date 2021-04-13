package torrent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFromFile(t *testing.T) {
	tfile, err := ParseFromFile("/home/wta/Projects/hypertube/debian.torrent")
	assert.Nil(t, err)

	trackers, err := tfile.Trackers()
	assert.Nil(t, err)

	if len(trackers) != 0 {
		t.Errorf("expected empty trackers, got %d len", len(trackers))
		return
	}
}
