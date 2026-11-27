BEGIN;
UPDATE products SET stock = stock - 5 WHERE id = 2;
SELECT 1;
COMMIT;
