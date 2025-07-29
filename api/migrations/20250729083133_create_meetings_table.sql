-- +goose Up
-- +goose StatementBegin
CREATE TABLE meetings (
    id serial primary key,
    dozen_id integer not null references dozens(id),
    start_time timestamp not null,
    end_time timestamp,
    location_name text not null,
    map_url text not null
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE meetings;
-- +goose StatementEnd
