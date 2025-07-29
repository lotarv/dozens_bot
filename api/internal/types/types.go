package types

type ContextKey string

const (
	ContextKeyUserID      ContextKey = "user_id"
	ContextKeyCredentials ContextKey = "credentials"
)

type Dozen struct {
	ID      int    `db:"id"`
	Code    string `db:"code"`
	Name    string `db:"name"`
	Captain int64  `db:"captain"`
}
