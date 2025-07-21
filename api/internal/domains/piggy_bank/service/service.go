package service

type repository interface {
}

type PiggyBankService struct {
	repo repository
}

func New(repo repository) *PiggyBankService {
	return &PiggyBankService{
		repo: repo,
	}
}
