-- +goose Up
-- +goose StatementBegin
CREATE TABLE piggy_banks (
    id SERIAL PRIMARY KEY,
    dozen_code INTEGER UNIQUE NOT NULL,
    balance INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE piggy_banks;
-- +goose StatementEnd
