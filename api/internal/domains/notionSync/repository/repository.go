package repository

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
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
	INSERT INTO members (fio, avatar_url, niche, annual_income, username)
	VALUES ($1,$2,$3,$4,$5)
	ON CONFLICT (username) DO UPDATE
	SET fio = EXCLUDED.fio,
		avatar_url = EXCLUDED.avatar_url,
		niche = EXCLUDED.niche,
		annual_income = EXCLUDED.annual_income`
	for _, member := range members {
		_, err := r.db.Exec(query, member.FIO, member.AvatarUrl, member.Niche, member.AnnualIncome, member.Username)
		if err != nil {
			slog.Error("Failed to sync member", "username", member.Username, "error", err)
			return fmt.Errorf("failed to sync member %s :%w", member.Username, err)
		}
	}
	slog.Info("Members synchronized", "count", len(members))
	return nil
}
