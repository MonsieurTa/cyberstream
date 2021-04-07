package fortytwo

type Reader interface{}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	GetAuthorizeURI() (string, error)
	GetAccessToken(code, state string) (*Token, error)
}
