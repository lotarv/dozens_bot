-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN username varchar(32) UNIQUE NOT NULL
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN username;
-- +goose StatementEnd
