package streaminfo

import (
	"path/filepath"

	"github.com/anacrolix/torrent"
)

type TorrentFile struct {
	*torrent.File
}

type StreamInfo struct {
	InfoHash       string
	StreamFile     *TorrentFile
	SubtitlesFiles []*TorrentFile
}

func (tf *TorrentFile) Ext() string {
	return filepath.Ext(tf.Path())
}

func (si *StreamInfo) HasSubtitles() bool {
	return len(si.SubtitlesFiles) != 0
}
