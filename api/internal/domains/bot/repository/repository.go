package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
	bot_types "github.com/lotarv/dozens_bot/internal/domains/bot/types/bot"
	documents_repo "github.com/lotarv/dozens_bot/internal/domains/documents/repository"
	document_types "github.com/lotarv/dozens_bot/internal/domains/documents/types"
	bank_repo "github.com/lotarv/dozens_bot/internal/domains/piggy_bank/repository"
	user_types "github.com/lotarv/dozens_bot/internal/domains/user/types"
	"github.com/lotarv/dozens_bot/internal/storage"
)

type BotRepository struct {
	db *sqlx.DB
	UsersRepository
	PiggyBankRepository
	DocumentsRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user *user_types.User) error
	UpdateUser(ctx context.Context, user *user_types.User) error
	GetUserByID(ctx context.Context, userID int64) (*user_types.User, error)
}

type PiggyBankRepository interface {
	ChangeBankBalance(ctx context.Context, piggyBankID int, amount int, reason string, username string) error
}

type DocumentsRepository interface {
	CreateDeclaration(declaration document_types.DeclarationDB) error
}

func New(storage *storage.Storage, userRepo UsersRepository, bankRepo *bank_repo.PiggyBankRepository, docsRepo *documents_repo.DocumentsRepository) *BotRepository {
	return &BotRepository{
		db:                  storage.DB(),
		UsersRepository:     userRepo,
		PiggyBankRepository: bankRepo,
		DocumentsRepository: docsRepo,
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

func (r *BotRepository) DeleteUserState(userID int64) error {
	query := "DELETE FROM user_state WHERE telegram_id=$1"
	_, err := r.db.Exec(query, userID)
	return err
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

func (r *BotRepository) GetDozenByCode(code string) (bot_types.Dozen, error) {
	var dozen bot_types.Dozen
	err := r.db.Get(&dozen, `SELECT * FROM dozens WHERE code = $1`, code)
	if err != nil {
		return bot_types.Dozen{}, err
	}

	slog.Info("repository dozen", "dozen", dozen)
	return dozen, err
}

func (r *BotRepository) AddUserToDozen(userID int64, dozenID int) error {
	query := "INSERT INTO user_dozen (dozen_id, user_id) VALUES ($1, $2)"
	_, err := r.db.Exec(query, dozenID, userID)
	return err
}

func (r *BotRepository) CreateDozen(dozen bot_types.Dozen) error {
	_, err := r.db.Exec("INSERT INTO dozens (code, name, captain) VALUES ($1, $2, $3)",
		dozen.Code, dozen.Name, dozen.Captain)
	return err
}

func (r *BotRepository) GetUserDozen(userID int64) (bot_types.Dozen, error) {
	var dozen bot_types.Dozen

	err := r.db.Get(&dozen, `
		SELECT d.*
		FROM dozens d
		JOIN user_dozen ud ON d.id = ud.dozen_id
		WHERE ud.user_id = $1
	`, userID)

	if err != nil {
		return bot_types.Dozen{}, err
	}

	return dozen, nil
}

func (r *BotRepository) SaveDocument(id, encryptedText string) error {
	query := `
		INSERT INTO documents (id, text)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE SET text = EXCLUDED.text
	`
	_, err := r.db.Exec(query, id, encryptedText)
	return err
}

func (r *BotRepository) SetEncryptedText(uuidStr string, encryptedText string) error {
	query := `
		UPDATE documents 
		SET text=$1
		WHERE id=$2
	`
	_, err := r.db.Exec(query, encryptedText, uuidStr)
	return err
}

func (r *BotRepository) GetMembers() ([]bot_types.MemberDB, error) {
	var members []bot_types.MemberDB
	query := `SELECT fio, username FROM members`
	err := r.db.Select(&members, query)
	if err != nil {
		return []bot_types.MemberDB{}, err
	}
	return members, nil
}

func (r *BotRepository) CreateMeeting(dozenID int, startTime string, endTime *string, locationName, mapURL string) error {
	query := `
		INSERT INTO meetings (dozen_id, start_time, end_time, location_name, map_url)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, dozenID, startTime, endTime, locationName, mapURL)
	return err
}
