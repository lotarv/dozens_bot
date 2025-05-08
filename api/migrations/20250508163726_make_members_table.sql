-- +goose Up
-- +goose StatementBegin
CREATE TABLE members (
    fio VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(2048),
    niche VARCHAR(255),
    annual_income BIGINT DEFAULT 0,
    username VARCHAR(100) UNIQUE NOT NULL,
    notion_database_id VARCHAR(50) unique not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE members;
-- +goose StatementEnd
