-- bozor: products jadvali va uchta mahsulot (SQL'ni GB-5.03'da oʻzimiz yozamiz)
CREATE TABLE products (
    id    bigint PRIMARY KEY,
    name  text    NOT NULL,
    price bigint  NOT NULL,
    stock integer NOT NULL
);

INSERT INTO products (id, name, price, stock) VALUES
    (1, 'olma', 1200000, 40),
    (2, 'non', 400000, 25),
    (3, 'oʻrik', 3500000, 5);
