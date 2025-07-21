package repository

import (
	"context"
	"time"

	"github.com/lotarv/dozens_bot/internal/storage"
)

type PiggyBankRepository struct {
	storage.Storage
}

func New(storage *storage.Storage) *PiggyBankRepository {
	return &PiggyBankRepository{
		Storage: *storage,
	}
}

func (r *PiggyBankRepository) AddTransactionTx(ctx context.Context, piggyBankID int, amount int, reason, author string) error {
	tx, err := r.Storage.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // на случай ошибки

	// Здесь можно добавить проверку баланса, лимитов и пр.

	query := `
		INSERT INTO piggy_bank_transactions (piggy_bank_id, amount, reason, author_username, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, query, piggyBankID, amount, reason, author, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}
