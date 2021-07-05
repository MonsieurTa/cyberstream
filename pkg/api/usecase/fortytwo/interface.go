package fortytwo

type Reader interface{}

type Writer interface{}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	// GetAuthorizeURI generates the 42 authorize uri and generate + save the associated
	// state in memory
	GetAuthorizeURI() (string, error)

	// GetToken consume the previously generated state by GetAuthorizeURI
	// to retrieve the 42 token by calling their endpoint
	GetToken(code, state string) (*Token, error)
	GetUserInfo(accessToken string) (*UserInfo, error)
}
