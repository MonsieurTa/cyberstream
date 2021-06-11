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

func (si *StreamInfo) UnwrapStreamFile() *torrent.File {
	return si.StreamFile.File
}

func (si *StreamInfo) UnwrapSubtitlesFiles() []*torrent.File {
	rv := make([]*torrent.File, 0, len(si.SubtitlesFiles))
	for _, f := range si.SubtitlesFiles {
		rv = append(rv, f.File)
	}
	return rv
}
