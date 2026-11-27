#!/usr/bin/env bash
# Videodagi barcha buyruqlarni ketma-ket, xuddi videodagidek ishga tushiradi va natijalarini chop etadi.
# Vaqtinchalik `bozor/` papkasida ishlaydi (compose loyihasi nomi = papka nomi, GB-5.02/5.03 videolaridagidek;
# shu uchun compose.yaml'dagi `name: gb504` qatori olib tashlanadi). Port 5464; oxirida `down -v` bilan tozalaydi.
# Videoda har soʻrov `sorov.sql` fayliga yoziladi va `psql -f /lab/sorov.sql` bilan ishga tushadi.
#   ./_video/chiqishlar.sh | tee _video/chiqishlar.txt
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
t=$(mktemp -d) && mkdir -p "$t/bozor"
trap 'cd "$t/bozor" && docker compose down -v >/dev/null 2>&1; rm -rf "$t"' EXIT
cd "$t/bozor" || exit 1
cp "$here/products.sql" "$here/orders.sql" "$here/seed.sql" .
grep -v '^name: gb504' "$here/compose.yaml" > compose.yaml
if docker compose ls --format json | grep -q '"Name":"bozor"'; then echo "bozor compose loyihasi allaqachon bor; keyinroq urining" >&2; exit 1; fi
run() { echo "\$ $*"; "$@" 2>&1; echo "(chiqish kodi: $?)"; echo; }
sql() { run docker compose exec db psql -c "$1"; }
# q: soʻrovni sorov.sql'ga yozadi (stdin) va uni psql -f bilan ishga tushiradi
q() { echo "-- sorov.sql:"; tee sorov.sql; run docker compose exec db psql -f /lab/sorov.sql; }

run docker compose up -d --wait
run docker compose exec db psql -f /lab/products.sql
run docker compose exec db psql -f /lab/orders.sql
run docker compose exec db psql -f /lab/seed.sql
sql "INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES (99, 1, 1, 1200000)"
sql 'SELECT * FROM orders'
sql 'SELECT * FROM order_items'

q <<'SQL'
SELECT orders.customer, order_items.quantity
FROM orders
JOIN order_items ON order_items.order_id = orders.id;
SQL

q <<'SQL'
SELECT orders.id, orders.customer,
       SUM(order_items.quantity * order_items.unit_price) AS jami
FROM orders
JOIN order_items ON order_items.order_id = orders.id;
SQL

q <<'SQL'
SELECT orders.id, orders.customer,
       SUM(order_items.quantity * order_items.unit_price) AS jami
FROM orders
JOIN order_items ON order_items.order_id = orders.id
GROUP BY orders.id
ORDER BY orders.id;
SQL

q <<'SQL'
SELECT orders.customer, order_items.quantity
FROM orders
LEFT JOIN order_items ON order_items.order_id = orders.id;
SQL

sql "SELECT id, customer FROM orders WHERE note = NULL"
sql "SELECT id, customer FROM orders WHERE note IS NULL"

q <<'SQL'
SELECT orders.id, orders.customer,
       SUM(order_items.quantity * order_items.unit_price) AS jami
FROM orders
LEFT JOIN order_items ON order_items.order_id = orders.id
GROUP BY orders.id
ORDER BY orders.id;
SQL

q <<'SQL'
SELECT orders.id, orders.customer,
       COALESCE(SUM(order_items.quantity * order_items.unit_price), 0) AS jami
FROM orders
LEFT JOIN order_items ON order_items.order_id = orders.id
GROUP BY orders.id
ORDER BY orders.id;
SQL

sql "SELECT id, customer FROM orders WHERE note <> 'ertalab yetkazing'"
sql "SELECT id, customer FROM orders WHERE note IS DISTINCT FROM 'ertalab yetkazing'"
