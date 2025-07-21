-- +goose Up
-- +goose StatementBegin
ALTER TABLE reports
ADD CONSTRAINT fk_reports_document
FOREIGN KEY (document_id) REFERENCES documents(document_notion_id)
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE reports
DROP CONSTRAINT fk_reports_document;
-- +goose StatementEnd
