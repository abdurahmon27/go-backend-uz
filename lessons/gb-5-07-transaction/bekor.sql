BEGIN;
UPDATE products SET stock = stock - 2 WHERE id = 1;
SELECT stock FROM products WHERE id = 1;
ROLLBACK;
SELECT stock FROM products WHERE id = 1;
