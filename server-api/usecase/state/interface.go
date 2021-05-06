package state

type Reader interface {
	Exist(state string) error
}
type Writer interface {
	Save(state string)
	Delete(state string)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Validate(state string) error
	Save(state string)
	Delete(state string)
}
