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
