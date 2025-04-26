package repository

import (
	"context"
	"fmt"

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

func (r *UsersRepository) CreateUser(ctx context.Context, user *types.User) error {
	query := `INSERT INTO users (id, full_name, avatar_url, niche, annual_income)
	VALUES (:id, :full_name, :avatar_url, :niche, :annual_income)`

	_, err := r.db.NamedExec(query, user)

	if err != nil {
		return fmt.Errorf("failed to create new user: %v", err)
	}
	return nil
}

func (r *UsersRepository) UpdateUser(ctx context.Context, user *types.User) error {
	query := `
        UPDATE users
        SET full_name = :full_name,
            avatar_url = :avatar_url,
            niche = :niche,
            annual_income = :annual_income,
        WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UsersRepository) GetUserByID(ctx context.Context, userID int64) (*types.User, error) {
	var user types.User
	err := r.db.Get(&user, `
	SELECT *
	FROM users
	WHERE id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get User by id: %v", err.Error())
	}

	return &user, nil
}

func (r *UsersRepository) GetAll(ctx context.Context) ([]types.User, error) {
	var users []types.User
	err := r.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err.Error())
	}
	return users, nil
}
