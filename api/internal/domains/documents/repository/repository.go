package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lotarv/dozens_bot/internal/domains/documents/types"
)

type DocumentsRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *DocumentsRepository {
	return &DocumentsRepository{
		db: db,
	}
}

func (r *DocumentsRepository) GetRules() (types.Document, error) {
	var rulesDoc types.Document
	err := r.db.Get(&rulesDoc, `SELECT * FROM documents WHERE id=$1`, 10)
	return rulesDoc, err
}
