package service

import (
	"context"

	"github.com/lotarv/dozens_bot/internal/domains/user/types"
)

type repository interface {
	GetAll(ctx context.Context) ([]types.User, error)
	GetUserByID(ctx context.Context, userID int64) (*types.User, error)
	UpdateUser(ctx context.Context, user *types.User) error
	CreateUser(ctx context.Context, user *types.User) error
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

func (s *UserService) CreateUser(ctx context.Context, user *types.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, user *types.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*types.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}
