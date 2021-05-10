package media

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{}
}
