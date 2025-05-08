package types

type Member struct {
	FIO                string `json:"fio" db:"fio"`
	AvatarUrl          string `json:"avatar_url" db:"avatar_url"`
	Niche              string `json:"niche" db:"niche"`
	AnnualIncome       int64  `json:"annual_income" db:"annual_income"`
	Username           string `json:"username" db:"username"`
	Notion_database_id string `json:"notion_database_id" db:"notion_database_id"`
}
