package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lotarv/dozens_bot/internal/domains/user/types"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type UsersRepository struct {
	db *sqlx.DB
}

func New(storage *storage.Storage) *UsersRepository {
	return &UsersRepository{
		db: storage.DB(),
	}
}

func (r *UsersRepository) GetAll(ctx context.Context) ([]types.User, error) {
	var users []types.User

	err := r.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}
