package types

type Meeting struct {
	ID           int     `json:"id" db:"id"`
	DozenID      int     `json:"dozen_id" db:"dozen_id"`
	StartTime    string  `json:"start_time" db:"start_time"`
	EndTime      *string `json:"end_time" db:"end_time"`
	LocationName string  `json:"location_name" db:"location_name"`
	MapUrl       string  `json:"map_url" db:"map_url"`
	MeetingDate  string  `json:"meeting_date" db:"meeting_date"`
}
