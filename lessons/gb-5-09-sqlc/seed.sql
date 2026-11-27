-- Sinov maʼlumoti: toʻrt mahsulot (narxlar tiyinda), uch buyurtma va ularning toʻrt qatori.
INSERT INTO products (name, price, stock, active)
VALUES ('olma', 1200000, 40, true),
       ('non', 400000, 3, true),
       ('choy', 2500000, 8, true),
       ('oʻrik', 3500000, 2, true);
INSERT INTO orders (customer, note)
VALUES ('Dilnoza', NULL),
       ('Jasur', 'ertalab yetkazing'),
       ('Malika', NULL);
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES (1, 1, 2, 1200000),
       (1, 2, 3, 400000),
       (2, 4, 1, 3500000),
       (2, 1, 1, 1200000);
