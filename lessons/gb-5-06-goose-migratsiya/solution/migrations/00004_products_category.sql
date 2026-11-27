-- +goose Up
ALTER TABLE products
    ADD COLUMN category text NOT NULL DEFAULT 'boshqa' CHECK (category <> '');

-- +goose Down
ALTER TABLE products DROP COLUMN category;
