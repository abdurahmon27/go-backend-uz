DROP TABLE IF EXISTS order_items, orders, products;
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
CREATE TABLE orders (
    id       bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer text   NOT NULL,
    note     text
);
CREATE TABLE order_items (
    id         bigint  GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id   bigint  NOT NULL REFERENCES orders (id),
    product_id bigint  NOT NULL REFERENCES products (id),
    quantity   integer NOT NULL CHECK (quantity > 0),
    unit_price bigint  NOT NULL CHECK (unit_price > 0)
);
INSERT INTO orders (customer, note)
VALUES ('Dilnoza', NULL),
       ('Jasur', 'ertalab yetkazing'),
       ('Malika', NULL);
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES (1, 1, 2, 1200000),
       (1, 2, 3, 400000),
       (2, 4, 1, 3500000),
       (2, 1, 1, 1200000);
