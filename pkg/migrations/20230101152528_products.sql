-- +goose Up
-- +goose StatementBegin

CREATE TABLE products
(
    created_at       TIMESTAMP DEFAULT now(),
    product_id       BIGINT PRIMARY KEY,
	name             VARCHAR(30),
	unit_composition VARCHAR(30),
	unit             VARCHAR(10),
	amount_unit	     VARCHAR(10)
);

CREATE INDEX products_id_idx ON products(product_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS products_id_idx;
DROP TABLE IF EXISTS products;

-- +goose StatementEnd
