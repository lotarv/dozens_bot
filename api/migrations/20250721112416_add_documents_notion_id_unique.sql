-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents
ADD CONSTRAINT unique_document_notion_id UNIQUE (document_notion_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents
DROP CONSTRAINT unique_document_notion_id;
-- +goose StatementEnd
