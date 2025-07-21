package service

import (
	"context"

	"github.com/lotarv/dozens_bot/internal/domains/piggy_bank/types"
)

type repository interface {
	GetPiggyBank(ctx context.Context, piggyBankID int) (types.PiggyBank, error)
}

type PiggyBankService struct {
	repo repository
}

func New(repo repository) *PiggyBankService {
	return &PiggyBankService{
		repo: repo,
	}
}

func (s *PiggyBankService) GetPiggyBank(ctx context.Context, bank_id int) (types.PiggyBank, error) {
	return s.repo.GetPiggyBank(ctx, bank_id)
}
