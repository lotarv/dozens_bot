-- +goose Up
-- +goose StatementBegin
CREATE TABLE reviews(
    id VARCHAR(36) PRIMARY KEY,
    meeting_notion_id VARCHAR(36) NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reviews
-- +goose StatementEnd
