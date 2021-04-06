package fortytwo

import (
	"golang.org/x/oauth2"
)

type Reader interface{}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	GetAuthorizeURI() (string, error)
	GetAccessToken(code, state string) (*oauth2.Token, error)
}
