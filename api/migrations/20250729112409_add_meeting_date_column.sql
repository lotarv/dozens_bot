-- +goose Up
-- +goose StatementBegin
ALTER TABLE meetings ADD COLUMN meeting_date date GENERATED ALWAYS AS (start_time::date) STORED;

CREATE UNIQUE INDEX unique_meeting_per_day_per_dozen ON meetings (dozen_id, meeting_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX unique_meeting_per_day_per_dozen;

ALTER TABLE meetings DROP COLUMN meeting_date;
-- +goose StatementEnd
