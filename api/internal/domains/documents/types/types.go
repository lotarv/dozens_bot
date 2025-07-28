package types

type Declaration struct {
	ID             string `json:"id"`
	DocumentID     string `json:"document_id" db:"document_id"`
	AuthorNotionID string `json:"author_notion_id"`
	CreationDate   string `json:"creation_date"`
	EndDate        string `json:"end_date"`
	Status         string `json:"status"`
}

type Report struct {
	ID             string `json:"id" db:"id"`
	DocumentID     string `json:"document_id" db:"document_id"`
	AuthorNotionID string `json:"author_notion_id" db:"author_notion_id"`
	CreationDate   string `json:"creation_date" db:"creation_date"`
}

type Document struct {
	ID               string `json:"id"`
	DocumentNotionID string `json:"document_notion_id" db:"document_notion_id"`
	Text             string `json:"text" db:"text"`
}

type ReportItem struct {
	CreationDate string `json:"creation_date"`
	Text         string `json:"text"`
}

type ReportsResponse struct {
	Username  string       `json:"username"`
	AvatarUrl string       `json:"avatar_url"`
	Reports   []ReportItem `json:"reports"`
}

type DeclarationDocument struct {
	ID           string  `json:"id" db:"id"`
	Text         string  `json:"text" db:"text"`
	EndDate      *string `json:"end_date" db:"end_date"`
	CreationDate string  `json:"creation_date" db:"creation_date"`
	Status       string  `json:"status" db:"status"`
}

type DeclarationDB struct {
	ID             string  `json:"id" db:"id"`
	AuthorNotionID string  `json:"author_notion_id" db:"author_notion_id"`
	Text           string  `json:"text" db:"text"`
	EndDate        *string `json:"end_date" db:"end_date"`
	CreationDate   string  `json:"creation_date" db:"creation_date"`
	Status         string  `json:"status" db:"status"`
}
