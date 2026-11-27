-- Namunaviy maʼlumot (GB-5.05..5.08 dagi bilan bir xil): toʻrt mahsulot, uch buyurtma.
-- Jadvallarni goose yaratadi (migrations/); bu fayl ularni tozalab, qaytadan toʻldiradi.
TRUNCATE order_items, orders, products RESTART IDENTITY;
INSERT INTO products (name, price, stock)
VALUES ('olma', 1200000, 40),
       ('non', 400000, 3),
       ('choy', 2500000, 8),
       ('oʻrik', 3500000, 2);
INSERT INTO orders (customer, note)
VALUES ('Dilnoza', NULL),
       ('Jasur', 'ertalab yetkazing'),
       ('Malika', NULL);
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES (1, 1, 2, 1200000),
       (1, 2, 3, 400000),
       (2, 4, 1, 3500000),
       (2, 1, 1, 1200000);
