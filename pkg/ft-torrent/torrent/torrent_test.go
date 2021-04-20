package torrent

import (
	"os"
	"testing"

	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/common"
	"github.com/stretchr/testify/assert"
)

func TestParseFromFile(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	tfiles, err := ParseFromFile(wdPath + `/testfile/debian.torrent`)
	assert.Nil(t, err)

	for _, tfile := range tfiles {
		peerID, err := common.GeneratePeerID()
		assert.Nil(t, err)

		_, err = tfile.Trackers(peerID)
		assert.Nil(t, err)
	}
}

func TestDownload(t *testing.T) {
	wdPath, err := os.Getwd()
	assert.Nil(t, err)

	tfiles, err := ParseFromFile(wdPath + `/testfile/debian.torrent`)
	assert.Nil(t, err)

	tfile := tfiles[0]

	err = tfile.Download()
	assert.Nil(t, err)
}
