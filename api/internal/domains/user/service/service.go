package service

import (
	"context"
	"os"

	member_types "github.com/lotarv/dozens_bot/internal/domains/members/types"
	"github.com/lotarv/dozens_bot/internal/domains/user/types"
)

type repository interface {
	GetAll(ctx context.Context) ([]types.User, error)
	GetUserByID(ctx context.Context, userID int64) (*types.User, error)
	UpdateUser(ctx context.Context, user *types.User) error
	CreateUser(ctx context.Context, user *types.User) error
	GetMemberByUsername(username string) (member_types.Member, error)
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

func (s *UserService) GetEncryptionKey(ctx context.Context, userID int64) (string, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	//TODO:сделать ключ для каждой десятки (после того как избавимся от Notion)
	//Пока проверяем, что пользователь есть в members
	_, err = s.GetMemberByUsername(ctx, user.Username)
	if err != nil {
		return "", err
	}

	return os.Getenv("ENCRYPTION_KEY"), nil
}

func (s *UserService) GetMemberByUsername(ctx context.Context, username string) (member_types.Member, error) {
	return s.repo.GetMemberByUsername(username)
}
