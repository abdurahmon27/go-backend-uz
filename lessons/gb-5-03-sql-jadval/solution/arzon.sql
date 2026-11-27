SELECT name, price
FROM products
WHERE active AND price < 3000000
ORDER BY price, id
LIMIT 2;
