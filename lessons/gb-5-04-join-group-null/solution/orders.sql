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
