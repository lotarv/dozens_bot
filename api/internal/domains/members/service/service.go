package service

import "github.com/lotarv/dozens_bot/internal/domains/members/types"

type repository interface {
	GetMembers() ([]types.Member, error)
}

type MembersService struct {
	repo repository
}

func New(repo repository) *MembersService {
	return &MembersService{
		repo: repo,
	}
}

func (s *MembersService) GetMembers() ([]types.Member, error) {
	members, err := s.repo.GetMembers()
	if err != nil {
		return nil, err
	}
	return members, nil
}
