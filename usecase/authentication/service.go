package authentication

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) Authenticate(username, password string) error {
	return s.repo.Validate(username, password)
}

func (s *Service) NewAccessToken() {

}
