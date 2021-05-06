package fortytwo

type Reader interface{}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	GetAuthorizeURI(state string) (string, error)
	GetAccessToken(code, state string) (*Token, error)
}
