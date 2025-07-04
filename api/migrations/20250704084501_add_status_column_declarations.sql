-- +goose Up
-- +goose StatementBegin
ALTER TABLE declarations
ADD COLUMN status VARCHAR(36) DEFAULT 'In progress';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE declarations
DROP COLUMN status;
-- +goose StatementEnd
