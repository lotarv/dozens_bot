-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_state (
    telegram_id BIGINT PRIMARY KEY,
    current_state TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_state;
-- +goose StatementEnd
