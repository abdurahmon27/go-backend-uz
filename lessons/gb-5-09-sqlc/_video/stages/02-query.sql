-- name: GetProduct :one
SELECT id, name, price, stock FROM products
WHERE id = $1;

-- name: ListActiveProducts :many
SELECT id, name, price, stock FROM products
WHERE active
ORDER BY id;
