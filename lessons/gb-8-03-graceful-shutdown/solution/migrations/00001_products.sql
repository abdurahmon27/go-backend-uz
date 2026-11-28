-- +goose Up
CREATE TABLE products (
    id     bigint
           GENERATED ALWAYS AS IDENTITY
           PRIMARY KEY,
    name   text    NOT NULL UNIQUE CHECK (name <> ''),
    price  bigint  NOT NULL CHECK (price > 0),
    stock  integer NOT NULL DEFAULT 0 CHECK (stock >= 0),
    active boolean NOT NULL DEFAULT true
);

-- +goose Down
DROP TABLE products;
