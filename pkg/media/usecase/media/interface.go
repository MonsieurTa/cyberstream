package media

type Reader interface{}
type Writer interface {
	AddMagnet()
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	StreamMagnet(magnet string) (string, <-chan bool, error)
}
