package state

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	return &Service{
		repo,
	}
}

func (s *Service) Validate(state string) error {
	return s.repo.Exist(state)
}
func (s *Service) Save(state string) {
	s.repo.Save(state)
}
func (s *Service) Delete(state string) {
	s.repo.Delete(state)
}
