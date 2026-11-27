#!/usr/bin/env bash
# Videodagi barcha SQL buyruqlarini ketma-ket, xuddi videodagidek ishga tushiradi va natijalarini chop etadi.
# Vaqtinchalik `bozor/` papkasida ishlaydi: compose loyihasi nomi papka nomidan, ya'ni `bozor` (GB-5.02 videosidagidek;
# compose.yaml'dagi `name: gb503` qatori shu uchun olib tashlanadi). Port 5463; oxirida `down -v` bilan tozalaydi.
#   ./_video/chiqishlar.sh | tee _video/chiqishlar.txt
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
t=$(mktemp -d) && mkdir -p "$t/bozor"
trap 'cd "$t/bozor" && docker compose down -v >/dev/null 2>&1; rm -rf "$t"' EXIT
cd "$t/bozor" || exit 1
cp "$here/seed.sql" .
grep -v '^name: gb503' "$here/compose.yaml" > compose.yaml
if docker compose ls --format json | grep -q '"Name":"bozor"'; then echo "bozor compose loyihasi allaqachon bor; keyinroq urining" >&2; exit 1; fi
run() { echo "\$ $*"; "$@" 2>&1; echo "(chiqish kodi: $?)"; echo; }
sql() { run docker compose exec db psql -c "$1"; }
run docker compose up -d --wait
sql '\dt'
cp "$here/_video/01-oddiy/products.sql" .
run docker compose exec db psql -f /lab/products.sql
sql '\dt'
run docker compose exec db psql -f /lab/seed.sql
sql 'SELECT * FROM products'
sql 'SELECT name, price FROM products'
sql 'SELECT name, price FROM products WHERE stock < 5'
sql 'SELECT name, price FROM products WHERE stock < 5 AND price > 1000000'
sql 'SELECT name, price FROM products WHERE stock < 5 ORDER BY price DESC'
sql 'SELECT name, price FROM products ORDER BY price DESC LIMIT 2'
sql 'SELECT name FROM products WHERE name = "olma"'
sql "SELECT name FROM products WHERE name = 'olma'"
sql "INSERT INTO products (name, price) VALUES ('sut', 900000) RETURNING id"
sql "INSERT INTO products (name, price) VALUES ('', -5) RETURNING *"
cp "$here/_video/02-qat/products.sql" .
run docker compose exec db psql -f /lab/products.sql
run docker compose exec db psql -f /lab/seed.sql
sql "INSERT INTO products (name, price) VALUES ('', -5)"
sql "INSERT INTO products (name, price) VALUES ('choy', 2600000)"
sql "INSERT INTO products (name, price) VALUES (NULL, 900000)"
sql "INSERT INTO products (name, price) VALUES ('sut', 900000) RETURNING id"
sql "INSERT INTO products (name, price, stock) VALUES ('qand', 1100000, -1)"
sql 'SELECT * FROM products ORDER BY id'
