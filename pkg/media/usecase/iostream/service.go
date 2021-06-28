package iostream

import (
	"io"
	"io/ioutil"
	"os"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	return &Service{repo}
}

func (s *Service) Save(dirName, fileName string, r io.ReadCloser) error {
	key := dirName + "/" + fileName

	if s.repo.PlaylistRequestExists(key) {
		s.repo.SetPlaylistReceived(key)
	}

	b, err := ioutil.ReadAll(r)
	defer r.Close()
	if err != nil {
		return err
	}

	dirpath := os.Getenv("STATIC_FILES_PATH") + "/hls/" + dirName
	err = os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		return err
	}

	filepath := dirpath + "/" + fileName
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	f, err := os.OpenFile(filepath, flags, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) WaitMasterPlaylist(dirName, fileName string) {
	<-s.repo.GetPlaylist(dirName + "/" + fileName)
}
