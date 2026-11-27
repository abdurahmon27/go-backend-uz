SELECT id, customer
FROM orders
WHERE note IS NULL
ORDER BY id;
