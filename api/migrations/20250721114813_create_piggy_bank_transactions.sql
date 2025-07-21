-- +goose Up
-- +goose StatementBegin
CREATE TABLE piggy_bank_transactions (
    id SERIAL PRIMARY KEY,
    piggy_bank_id INTEGER NOT NULL REFERENCES piggy_banks(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    reason TEXT,
    author_username TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE piggy_bank_transactions;
-- +goose StatementEnd
