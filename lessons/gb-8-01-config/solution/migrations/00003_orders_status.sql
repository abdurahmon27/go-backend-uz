-- +goose Up
ALTER TABLE orders
    ADD COLUMN status text NOT NULL DEFAULT 'yangi' CHECK (status <> '');

-- +goose Down
ALTER TABLE orders DROP COLUMN status;
