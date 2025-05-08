package types

type Member struct {
	FIO          string `json:"fio"`
	AvatarUrl    string `json:"avatar_url"`
	Niche        string `json:"niche"`
	AnnualIncome int64  `json:"annual_income"`
	Username     string `json:"username"`
}
