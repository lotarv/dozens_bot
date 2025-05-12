package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lotarv/dozens_bot/internal/domains/members/types"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type MembersRepository struct {
	storage *sqlx.DB
}

func New(storage *storage.Storage) *MembersRepository {
	return &MembersRepository{
		storage: storage.DB(),
	}
}

func (s *MembersRepository) GetMembers() ([]types.Member, error) {
	var members []types.Member
	err := s.storage.Select(&members, `SELECT * FROM members`)
	if err != nil {
		return nil, err
	}
	return members, nil
}
