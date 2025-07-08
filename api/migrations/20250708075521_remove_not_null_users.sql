-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ALTER COLUMN full_name DROP NOT NULL,
    ALTER COLUMN niche DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    ALTER COLUMN full_name SET NOT NULL,
    ALTER COLUMN niche SET NOT NULL;
-- +goose StatementEnd
