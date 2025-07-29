-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_dozen (
    dozen_id INT NOT NULL,
    username TEXT NOT NULL,
    PRIMARY KEY (username),
    CONSTRAINT fk_dozen FOREIGN KEY (dozen_id) REFERENCES dozens(id) ON DELETE CASCADE,
    CONSTRAINT fk_username FOREIGN KEY (username) REFERENCES members(username) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_dozen;
-- +goose StatementEnd
