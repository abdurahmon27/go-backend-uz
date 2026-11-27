INSERT INTO orders (customer, note)
VALUES ('Dilnoza', NULL),
       ('Jasur', 'ertalab yetkazing'),
       ('Malika', NULL);
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES (1, 1, 2, 1200000),
       (1, 2, 3, 400000),
       (2, 4, 1, 3500000),
       (2, 1, 1, 1200000);
