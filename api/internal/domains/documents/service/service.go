package service

import "github.com/lotarv/dozens_bot/internal/domains/documents/types"

type repository interface {
	GetRules() (types.Document, error)
}

type DocumentsService struct {
	repo repository
}

func New(repo repository) *DocumentsService {
	return &DocumentsService{
		repo: repo,
	}
}

func (s *DocumentsService) GetRules() (types.Document, error) {
	return s.repo.GetRules()
}
