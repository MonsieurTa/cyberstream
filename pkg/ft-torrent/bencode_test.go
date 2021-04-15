package torrent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBencodeTorrent(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	r, _ := os.Open(wdPath + `/test.mkv.torrent`)

	defer r.Close()
	_, err = ReadTorrentFile(r)
	assert.Nil(t, err)
}
