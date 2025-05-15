package repository

import (
	"fmt"
	"log/slog"
	"strings"

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

func (r *DocumentsRepository) GetRawReports(username string) ([]types.Report, error) {
	member_notion_id, err := r.getMemberNotionId(username)
	if err != nil {
		return nil, err
	}

	var reports []types.Report
	err = r.db.Select(&reports, "SELECT * FROM reports WHERE author_notion_id = $1 ORDER BY creation_date DESC", member_notion_id)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *DocumentsRepository) getMemberNotionId(username string) (string, error) {
	var member_notion_id string
	err := r.db.QueryRow("SELECT notion_database_id FROM members WHERE username=$1", username).Scan(&member_notion_id)
	if err != nil {
		return "", err
	}
	return member_notion_id, nil
}

func (r *DocumentsRepository) GetUserAvatarUrl(username string) (string, error) {
	var avatar_url string
	err := r.db.QueryRow("SELECT avatar_url FROM members WHERE username=$1", username).Scan(&avatar_url)
	if err != nil {
		return "", err
	}
	return avatar_url, nil
}

func (r *DocumentsRepository) GetReportDocuments(reports []types.Report) (map[string]types.Document, error) {
	if len(reports) == 0 {
		return map[string]types.Document{}, nil
	}

	// Собираем список нужных document ID
	ids := make([]interface{}, 0, len(reports))
	for _, report := range reports {
		ids = append(ids, report.DocumentID)
	}

	// Строим SQL IN (...), например: WHERE id IN ($1, $2, $3...)
	query := "SELECT id, document_notion_id, text FROM documents WHERE document_notion_id IN ("
	params := make([]string, len(ids))
	for i := range ids {
		params[i] = fmt.Sprintf("$%d", i+1)
	}
	query += strings.Join(params, ",") + ")"
	slog.Info("query", "query", query, "ids", ids)

	// Выполняем запрос
	rows, err := r.db.Query(query, ids...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Сканируем результат
	documents := make(map[string]types.Document)
	for rows.Next() {
		var doc types.Document
		if err := rows.Scan(&doc.ID, &doc.DocumentNotionID, &doc.Text); err != nil {
			return nil, err
		}
		documents[doc.DocumentNotionID] = doc
	}

	return documents, nil
}
