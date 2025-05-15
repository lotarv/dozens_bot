-- +goose Up
-- +goose StatementBegin
CREATE TABLE reports (
    id VARCHAR(36) PRIMARY KEY,
    document_id VARCHAR(36) NOT NULL,
    author_notion_id VARCHAR(36) NOT NULL,
    creation_date DATE NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reports;
-- +goose StatementEnd
