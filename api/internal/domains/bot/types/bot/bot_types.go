package bot_types

type NotionConfig struct {
	DocumentsDBID string
	ReportsDBID   string
}
type UserState string

const (
	StateIdle          UserState = "idle"
	StateEnteringDozen UserState = "entering_dozen_name"
	StateJoiningDozen  UserState = "entering_join_code"
)
