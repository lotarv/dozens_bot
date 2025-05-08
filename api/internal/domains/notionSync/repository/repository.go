package repository

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	document_types "github.com/lotarv/dozens_bot/internal/domains/documents/types"
	member_types "github.com/lotarv/dozens_bot/internal/domains/members/types"
)

type NotionSyncRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *NotionSyncRepository {
	return &NotionSyncRepository{
		db: db,
	}
}

func (r *NotionSyncRepository) SyncMembersWithNotion(members []member_types.Member) error {
	query := `
	INSERT INTO members (fio, avatar_url, niche, annual_income, username, notion_database_id)
	VALUES ($1,$2,$3,$4,$5,$6)
	ON CONFLICT (username) DO UPDATE
	SET fio = EXCLUDED.fio,
		avatar_url = EXCLUDED.avatar_url,
		niche = EXCLUDED.niche,
		annual_income = EXCLUDED.annual_income,
		notion_database_id = EXCLUDED.notion_database_id`
	for _, member := range members {
		_, err := r.db.Exec(query, member.FIO, member.AvatarUrl, member.Niche, member.AnnualIncome, member.Username, member.Notion_database_id)
		if err != nil {
			slog.Error("Failed to sync member", "username", member.Username, "error", err)
			return fmt.Errorf("failed to sync member %s :%w", member.Username, err)
		}
	}
	slog.Info("Members synchronized", "count", len(members))
	return nil
}

func (r *NotionSyncRepository) SyncDeclarationsWithNotion(declarations []document_types.Declaration) error {
	query := `
		INSERT INTO declarations (id, author_notion_id, creation_date, end_date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	for _, declaration := range declarations {
		_, err := r.db.Exec(query, declaration.ID, declaration.AuthorNotionID, declaration.CreationDate, declaration.EndDate)
		if err != nil {
			slog.Error("failed to synchronize declaration: ", "declaration_id", declaration.ID, "author", declaration.AuthorNotionID, "error", err)
			continue
		}
	}
	return nil
}
