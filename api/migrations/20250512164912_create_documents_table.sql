-- +goose Up
-- +goose StatementBegin
CREATE TABLE documents (
    id integer PRIMARY KEY,
    document_notion_id VARCHAR(36) NOT NULL,
    text TEXT
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE documents
-- +goose StatementEnd
