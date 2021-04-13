package torrent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFromFile(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	tfile, err := ParseFromFile(wdPath + `/debian.torrent`)
	assert.Nil(t, err)

	trackers, err := tfile.Trackers()
	assert.Nil(t, err)

	if len(trackers) != 0 {
		t.Errorf("expected empty trackers, got %d len", len(trackers))
		return
	}
}
