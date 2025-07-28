-- +goose Up
-- +goose StatementBegin
ALTER TABLE declarations ADD COLUMN text TEXT not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE declarations DROP COLUMN text;
-- +goose StatementEnd
