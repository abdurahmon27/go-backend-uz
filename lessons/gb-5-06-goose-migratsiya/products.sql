-- GB-5.03'dagi products jadvali va toʻrt mahsulot (narxlar tiyinda)
CREATE TABLE products (
    id     bigint
           GENERATED ALWAYS AS IDENTITY
           PRIMARY KEY,
    name   text    NOT NULL UNIQUE CHECK (name <> ''),
    price  bigint  NOT NULL CHECK (price > 0),
    stock  integer NOT NULL DEFAULT 0 CHECK (stock >= 0),
    active boolean NOT NULL DEFAULT true
);
INSERT INTO products (name, price, stock, active)
VALUES ('olma', 1200000, 40, true),
       ('non', 400000, 3, true),
       ('choy', 2500000, 8, true),
       ('oʻrik', 3500000, 2, true);
