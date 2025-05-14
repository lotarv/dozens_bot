package repository

import (
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

func (r *BotRepository) GetDocumentNotionId(db_id int) (string, error) {
	var document_notion_id string
	err := r.db.QueryRow("SELECT document_notion_id FROM documents WHERE id = $1", db_id).Scan(&document_notion_id)
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
