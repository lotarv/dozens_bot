package types

type Document struct {
	ID           int    `json:"id"`
	Type         string `json:"type"`
	FileUrl      string `json:"file_url"`
	AuthorUrl    string `json:"author_url"`
	CreationDate string `json:"creation_date"`
}

type Declaration struct {
	ID             string `json:"id"`
	AuthorNotionID string `json:"author_notion_id"`
	CreationDate   string `json:"creation_date"`
	EndDate        string `json:"end_date"`
}

type Report struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}
