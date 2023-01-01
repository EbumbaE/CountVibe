-- +goose Up
-- +goose StatementBegin

CREATE TABLE diary
(
    created_at   TIMESTAMP DEFAULT now(),
    diary_id     BIGINT,
    date         VARCHAR(10),
    meal_name    VARCHAR(30),
    product_id   BIGINT,
    amount       FLOAT
);

CREATE INDEX diary_diary_id_idx ON diary(diary_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS diary_diary_id_idx;
DROP TABLE IF EXISTS diary;

-- +goose StatementEnd
