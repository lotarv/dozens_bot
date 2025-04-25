package types

import "database/sql"

type User struct {
	ID           int            `json:"id" db:"id"`
	FullName     string         `json:"full_name" db:"full_name"`
	AvatarURL    sql.NullString `json:"avatar_url" db:"avatar_url"`
	Niche        string         `json:"niche" db:"niche"`
	AnnualIncome float64        `json:"annual_income" db:"annual_income"`
	TelegramID   string         `json:"telegram_id" db:"telegram_id"`
}
