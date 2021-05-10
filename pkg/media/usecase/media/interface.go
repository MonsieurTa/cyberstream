package media

type Reader interface{}
type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
}
