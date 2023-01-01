-- +goose Up
-- +goose StatementBegin

CREATE TABLE users
(
    created_at  TIMESTAMP DEFAULT now(),
    user_id     BIGINT PRIMARY KEY,
    diary_id    BIGINT,
    username    VARCHAR(30),
    password    VARCHAR(30)
);

CREATE INDEX users_id_idx ON users(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS users_id_idx;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
