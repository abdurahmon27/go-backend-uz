CREATE TABLE products (
    id     bigint
           GENERATED ALWAYS AS IDENTITY
           PRIMARY KEY,
    name   text,
    price  bigint,
    stock  integer,
    active boolean
);
