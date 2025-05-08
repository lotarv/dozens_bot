-- +goose Up
-- +goose StatementBegin
CREATE TABLE declarations (
    id VARCHAR(36) PRIMARY KEY,
    author_notion_id VARCHAR(36) NOT NULL,
    creation_date DATE NOT NULL,
    end_date DATE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
