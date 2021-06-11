package iostream

import "io"

type Reader interface {
	GetPlaylist(key string) chan bool
	PlaylistRequestExists(key string) bool
}

type Writer interface {
	SetPlaylistReceived(key string) error
	ClosePlaylistRequest(key string)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Save(dirName, fileName string, r io.ReadCloser) error
	WaitMasterPlaylist(dirName, fileName string)
}
