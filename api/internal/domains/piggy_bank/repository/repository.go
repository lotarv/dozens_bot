package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lotarv/dozens_bot/internal/domains/piggy_bank/types"
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

func (r *PiggyBankRepository) ChangeBankBalance(ctx context.Context, piggyBankID int, amount int, reason string, username string) error {
	tx, err := r.Storage.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // на случай ошибки

	//1 Обновляем баланс
	updateQuery := `
	UPDATE piggy_banks
	SET balance = balance + $1
	WHERE ID = $2
	`

	_, err = tx.ExecContext(ctx, updateQuery, amount, piggyBankID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	//2 Добавляем запись о транзакции

	insertQuery := `
	INSERT INTO piggy_bank_transactions (piggy_bank_id, amount, reason,author_username)
	VALUES ($1, $2, $3, $4)
	`
	_, err = tx.ExecContext(ctx, insertQuery, piggyBankID, amount, reason, username)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return tx.Commit()
}

func (r *PiggyBankRepository) GetPiggyBank(ctx context.Context, piggyBankID int) (types.PiggyBank, error) {
	var piggy types.PiggyBank

	//1. Получаем баланс
	err := r.DB().QueryRowContext(ctx, `
		SELECT balance
		FROM piggy_banks
		WHERE id = $1`, piggyBankID).Scan(&piggy.Balance)
	if err != nil {
		return piggy, fmt.Errorf("failed to get balance: %w", err)
	}

	//2. Получаем транзакции
	var rawTxns []struct {
		Amount    int       `db:"amount"`
		Reason    int       `db:"reason"`
		CreatedAt time.Time `db:"created_at"`
		FullName  string    `db:"fio"`
		AvatarUrl string    `db:"avatar_url"`
	}
	err = r.DB().SelectContext(ctx, &rawTxns, `
		SELECT 
			t.amount,
			t.reason,
			t.created_at,
			m.fio,
			m.avatar_url
		FROM piggy_bank_transactions t
		JOIN members m ON m.username = t.author_username
		WHERE t.piggy_bank_id = $1
		ORDER BY t.created_at DESC`, piggyBankID)
	if err != nil {
		return piggy, fmt.Errorf("failed to get transactions: %w", err)
	}

	for _, row := range rawTxns {
		piggy.Transactions = append(piggy.Transactions, types.Transaction{
			Amount: row.Amount,
			Reason: row.Reason,
			Date:   row.CreatedAt,
			Member: types.TransactionMember{
				FullName:  row.FullName,
				AvatarUrl: row.AvatarUrl,
			},
		})
	}

	return piggy, nil
}
