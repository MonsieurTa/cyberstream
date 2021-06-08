package file

import (
	"errors"
	"os"
	"path/filepath"
)

var mediaType map[string]string = map[string]string{
	".mp4": "video/mp4",
	".mkv": "application/x-mpegURL",
}

func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func GetMediaType(filename string) (string, error) {
	ext := filepath.Ext(filename)
	rv, ok := mediaType[ext]
	if !ok {
		return "", errors.New("could not recognize media type")
	}
	return rv, nil
}
