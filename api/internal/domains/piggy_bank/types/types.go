package types

import "time"

type TransactionMember struct {
	FullName  string `json:"full_name" db:"fio"`
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
}

type Transaction struct {
	Member TransactionMember `json:"member"`
	Date   time.Time         `json:"created_at" db:"created_at"`
	Amount int               `json:"amount" db:"amount"`
	Reason string            `json:"reason" db:"reason"`
}

type PiggyBank struct {
	Balance      int           `json:"balance"`
	Transactions []Transaction `json:"transactions"`
}
