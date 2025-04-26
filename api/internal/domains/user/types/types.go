package types

import "database/sql"

type User struct {
	ID           int64          `json:"id" db:"id"`
	FullName     string         `json:"full_name" db:"full_name"`
	AvatarURL    sql.NullString `json:"avatar_url" db:"avatar_url"`
	Niche        string         `json:"niche" db:"niche"`
	AnnualIncome float64        `json:"annual_income" db:"annual_income"`
}
