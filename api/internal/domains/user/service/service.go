package service

import (
	"context"

	"github.com/lotarv/dozens_bot/internal/domains/user/types"
)

type repository interface {
	GetAll(ctx context.Context) ([]types.User, error)
}

type UserService struct {
	repo repository
}

func New(repo repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetAll(ctx context.Context) ([]types.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
