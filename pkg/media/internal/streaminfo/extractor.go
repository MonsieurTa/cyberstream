package streaminfo

import (
	"errors"
	"fmt"

	"github.com/anacrolix/torrent"
)

var (
	STREAM_FILE_EXTS     = map[string]bool{".mp4": true, ".mkv": true}
	SUBTITLES_FILES_EXTS = map[string]bool{".srt": true, ".vtt": true}
)

// Need to call GetInfo() beforehand
func Extract(t *torrent.Torrent) (*StreamInfo, error) {
	file, err := biggestFile(t.Files())
	if err != nil {
		return nil, err
	}

	streamFile := &TorrentFile{file}
	err = validateStreamTorrentFile(streamFile)
	if err != nil {
		return nil, err
	}

	return &StreamInfo{
		InfoHash:       t.InfoHash().HexString(),
		StreamFile:     streamFile,
		SubtitlesFiles: getSubtitlesTorrentFiles(t.Files()),
	}, nil
}

func biggestFile(files []*torrent.File) (*torrent.File, error) {
	size := int64(-1)
	index := -1

	if len(files) == 0 || files == nil {
		return nil, errors.New("biggestFile(): bad parameter")
	}

	for i, f := range files {
		filesize := f.Length()
		if filesize > size {
			size = filesize
			index = i
		}
	}
	return files[index], nil
}

func validateStreamTorrentFile(file *TorrentFile) error {
	fileExt := file.Ext()
	_, ok := STREAM_FILE_EXTS[fileExt]
	if !ok {
		return fmt.Errorf("invalid file ext: %s", fileExt)
	}
	return nil
}

func getSubtitlesTorrentFiles(files []*torrent.File) []*TorrentFile {
	rv := make([]*TorrentFile, 0, len(files))
	for _, f := range files {
		tf := &TorrentFile{f}
		_, ok := SUBTITLES_FILES_EXTS[tf.Ext()]
		if ok {
			rv = append(rv, tf)
		}
	}
	return rv
}
