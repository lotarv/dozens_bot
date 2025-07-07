package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type BotRepository struct {
	db *sqlx.DB
}

func New(storage *storage.Storage) *BotRepository {
	return &BotRepository{
		db: storage.DB(),
	}
}

func (r *BotRepository) GetDocumentsAmount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM documents").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (r *BotRepository) GetDocumentNotionId(db_uuid string) (string, error) {
	var document_notion_id string
	err := r.db.QueryRow("SELECT document_notion_id FROM documents WHERE id = $1", db_uuid).Scan(&document_notion_id)
	if err != nil {
		return "", err
	}
	return document_notion_id, nil
}

func (r *BotRepository) GetMemberNotionId(username string) (string, error) {
	var author_notion_id string
	err := r.db.QueryRow("SELECT notion_database_id FROM members WHERE username = $1", username).Scan(&author_notion_id)
	if err != nil {
		return "", err
	}
	return author_notion_id, nil
}

func (r *BotRepository) GetUserState(userID int64) (string, error) {
	var state string
	err := r.db.QueryRow(`SELECT current_state FROM user_state WHERE telegram_id = $1`, userID).Scan(&state)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return state, err
}

func (r *BotRepository) SetUserState(userID int64, state string) error {
	_, err := r.db.Exec(`
		INSERT INTO user_state (telegram_id, current_state, updated_at)
		VALUES ($1, $2, now())
		ON CONFLICT (telegram_id) DO UPDATE
		SET current_state = EXCLUDED.current_state,
			updated_at = EXCLUDED.updated_at
	`, userID, state)
	return err
}

func (r *BotRepository) ResetUserState(userID int64) error {
	_, err := r.db.Exec(`DELETE FROM user_state WHERE telegram_id = $1`, userID)
	return err
}

func (r *BotRepository) GetDozenByCode(code string) (int, error) {
	var dozenID int
	err := r.db.QueryRow(`SELECT id FROM dozens WHERE code = $1`, code).Scan(&dozenID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return dozenID, err
}
