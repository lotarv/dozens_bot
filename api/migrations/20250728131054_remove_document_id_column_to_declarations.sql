-- +goose Up
-- +goose StatementBegin
ALTER TABLE declarations
DROP COLUMN document_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE declarations
ADD COLUMN document_id VARCHAR(36);
-- +goose StatementEnd
