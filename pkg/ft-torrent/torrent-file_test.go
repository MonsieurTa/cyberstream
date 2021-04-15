package torrent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFromFile(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	tfiles, err := ParseFromFile(wdPath + `/debian.torrent`)
	assert.Nil(t, err)

	for _, tfile := range tfiles {
		_, err := tfile.Trackers()
		assert.Nil(t, err)
	}
}
