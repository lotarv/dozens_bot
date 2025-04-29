-- +goose Up
-- +goose StatementBegin
CREATE TABLE groups (
    chat_id BIGINT PRIMARY KEY,
    title TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE groups;
-- +goose StatementEnd
