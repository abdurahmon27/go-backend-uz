-- name: GetProduct :one
SELECT id, name, price, stock FROM products
WHERE id = $1;

-- name: ListActiveProducts :many
SELECT id, name, price, stock FROM products
WHERE active
ORDER BY id;

-- name: CreateOrder :one
INSERT INTO orders (customer) VALUES ($1)
RETURNING id;

-- name: DecrementStock :one
UPDATE products SET stock = stock - @quantity
WHERE id = @id
RETURNING price;

-- name: AddOrderItem :exec
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES ($1, $2, $3, $4);

-- name: ProductByName :one
SELECT id, name, price, stock FROM products
WHERE name = $1;

-- name: CheapestProducts :many
SELECT id, name, price, stock FROM products
WHERE active
ORDER BY price, id
LIMIT $1;

-- name: OrderTotal :one
SELECT COALESCE(SUM(quantity * unit_price), 0)::bigint AS total
FROM order_items
WHERE order_id = $1;

-- name: AddStock :execrows
UPDATE products SET stock = stock + @quantity
WHERE id = @id;
