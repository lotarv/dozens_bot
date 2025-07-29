-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS user_dozen;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Можно восстановить старую структуру, если нужно
CREATE TABLE user_dozen (
    dozen_id INT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (user_id),
    CONSTRAINT fk_dozen FOREIGN KEY (dozen_id) REFERENCES dozens(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd
