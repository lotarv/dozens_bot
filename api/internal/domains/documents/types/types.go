package types

type Declaration struct {
	ID             string `json:"id"`
	AuthorNotionID string `json:"author_notion_id"`
	CreationDate   string `json:"creation_date"`
	EndDate        string `json:"end_date"`
}

type Report struct {
	ID             string `json:"id" db:"id"`
	AuthorNotionID string `json:"author_notion_id" db:"author_notion_id"`
	CreationDate   string `json:"creation_date" db:"creation_date"`
}

type Document struct {
	ID               int    `json:"id"`
	DocumentNotionID string `json:"document_notion_id" db:"document_notion_id"`
	Text             string `json:"text" db:"text"`
}
