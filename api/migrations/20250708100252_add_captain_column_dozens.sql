-- +goose Up
-- +goose StatementBegin
ALTER TABLE dozens
ADD COLUMN captain BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE dozens
DROP COLUMN captain;
-- +goose StatementEnd
