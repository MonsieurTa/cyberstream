package torrent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBencodeTorrent(t *testing.T) {
	r, _ := os.Open("/home/wta/Projects/hypertube/test.mkv.torrent")

	defer r.Close()
	_, err := Open(r)
	assert.Nil(t, err)
}
