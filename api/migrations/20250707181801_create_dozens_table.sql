-- +goose Up
-- +goose StatementBegin
CREATE TABLE dozens (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dozens;
-- +goose StatementEnd
