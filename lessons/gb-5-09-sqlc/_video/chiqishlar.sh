#!/usr/bin/env bash
# GB-5.09: videodagi har buyruqni ketma-ket ishga tushiradi va chiqishlarini chop etadi.
# Ishlatish: GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0 ./_video/chiqishlar.sh > _video/chiqishlar.txt 2>&1
# Kerak: PATH'da goose v3.30.2 (GB-5.06). sqlc'ni skript o'zi `go install` bilan vaqtinchalik GOBIN'ga o'rnatadi.
# Vaqtinchalik "bozor" papkasida ishlaydi (compose loyiha nomi = papka nomi, videodagidek) va oxirida docker compose down -v bilan tozalaydi.
set -u
HERE="$(cd "$(dirname "$0")" && pwd)"
LAB="$HERE/.."
TMP="$(mktemp -d)"
W="$TMP/bozor"
mkdir -p "$W/migrations"
trap 'cd "$W" && docker compose down -v >/dev/null 2>&1; rm -rf "$TMP"' EXIT
cd "$W" || exit 1
if docker compose ls --format json | grep -q '"Name":"bozor"'; then echo "bozor compose loyihasi allaqachon bor; keyinroq urining" >&2; exit 1; fi
grep -v '^name: gb509' "$LAB/compose.yaml" > compose.yaml
cp "$LAB/seed.sql" "$LAB/go.mod" "$LAB/go.sum" .
cp "$LAB"/migrations/*.sql migrations/
export GOBIN="$TMP/bin" PATH="$TMP/bin:$PATH"
step() { echo; echo "### $1"; shift; echo "\$ $*"; "$@" 2>&1; echo "[exit $?]"; }
gocode() { rm -f ./*.go; for f in "$@"; do cp "$f" .; done; }

# 0. baza: compose, goose up (GB-5.06), sinov ma'lumoti
step up docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5409/bozor
step goose-up goose up
step seed docker compose exec db psql -f /lab/seed.sql
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5409/bozor

# 1. muammo: qo'lda yozilgan Get, ustun nomi xato (stock → qoldiq)
gocode "$HERE/stages/01-store-qolda.go" "$HERE/stages/01-main.go"
mv 01-store-qolda.go store.go; mv 01-main.go main.go
step build go build
step vet go vet .
step run-xato go run .
rm -f bozor

# 2. sqlc o'rnatish
step install go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
step version sqlc version

# 3. query.sql (ikki so'rov) + sqlc.yaml → generate
cp "$LAB/sqlc.yaml" .
cp "$HERE/stages/02-query.sql" query.sql
step generate1 sqlc generate
step ls-db ls db
echo; echo "### models.go"; cat db/models.go
echo; echo "### query.sql.go"; cat db/query.sql.go
echo; echo "### db.go"; cat db/db.go

# 4. Store sqlc orqali: Get
gocode "$HERE/stages/02-store-get.go" "$HERE/stages/02-main-get.go"
mv 02-store-get.go store.go; mv 02-main-get.go main.go
step run-get go run .

# 5. transaction so'rovlari: avval $1/$2, keyin @quantity/@id
cp "$HERE/stages/03-query-dollar.sql" query.sql
step generate2 sqlc generate
echo; echo "### DecrementStockParams (\$1, \$2)"; sed -n '/^type DecrementStockParams/,/^}/p' db/query.sql.go
cp "$LAB/query.sql" query.sql
step generate3 sqlc generate
echo; echo "### DecrementStockParams (@quantity, @id)"; sed -n '/^type DecrementStockParams/,/^}/p' db/query.sql.go
echo; echo "### AddOrderItemParams"; sed -n '/^type AddOrderItemParams/,/^}/p' db/query.sql.go
sed -n '/^func (q \*Queries) AddOrderItem/,/^}/p' db/query.sql.go
sed -n '/^const createOrder/,/^}/p' db/query.sql.go

# 6. yakuniy store.go + main.go: Get va PlaceOrder
gocode "$LAB/store.go" "$LAB/main.go"
step run-final go run .
step tekshir docker compose exec db psql -c 'SELECT order_id, product_id, quantity, unit_price FROM order_items WHERE order_id = 4'

# 7. xato: query.sql'da ustun nomi xato → sqlc generate yiqiladi, db/ o'zgarmaydi
before=$(cat db/*.go | sha256sum)
sed '0,/SELECT id, name, price, stock FROM products/s//SELECT id, name, price, qoldiq FROM products/' "$LAB/query.sql" > query.sql
head -3 query.sql
step generate-xato sqlc generate
after=$(cat db/*.go | sha256sum)
[ "$before" = "$after" ] && echo "db/*.go o'zgarmadi" || echo "db/*.go O'ZGARDI"
step build-eski go build
rm -f bozor
cp "$LAB/query.sql" query.sql

# 8. quiz: query.sql'ga yangi so'rov, lekin generate unutildi → Go kompilyatori
gocode "$HERE/stages/05-store-undefined.go" "$LAB/main.go"
mv 05-store-undefined.go store.go
printf '\n-- name: ProductByName :one\nSELECT id, name, price, stock FROM products\nWHERE name = $1;\n' >> query.sql
tail -3 query.sql
step build-undefined go build
cp "$LAB/query.sql" query.sql
rm -f bozor

# 9. papkaning o'zi (boshlang'ich holat) va yechim
gocode "$LAB"/*.go
rm -rf db && cp -r "$LAB/db" db && cp "$LAB/query.sql" query.sql
step starter go test .
mkdir -p sol && cp -r "$LAB"/solution/. sol/
(cd sol && step solution go test ./...)

# 10. README/TASK'dagi qo'shimcha o'lchovlar (videoda yo'q): ifoda turlari, emit_empty_slices, sqlc diff
M="$TMP/m" && mkdir -p "$M" && cp -r "$LAB"/solution/. "$M"/ && cd "$M"
sigs() { sqlc generate && grep -h "^func (q \*Queries) \(T\|Full\|L\)" db/query.sql.go; }
cp query.sql query.sql.orig
printf -- '-- name: T1 :one\nSELECT COALESCE(SUM(quantity * unit_price), 0) FROM order_items WHERE order_id = $1;\n\n-- name: T2 :one\nSELECT SUM(quantity * unit_price) AS total FROM order_items WHERE order_id = $1;\n\n-- name: T3 :one\nSELECT COALESCE(SUM(quantity * unit_price), 0)::bigint AS total FROM order_items WHERE order_id = $1;\n\n-- name: Full :one\nSELECT id, name, price, stock, active FROM products WHERE id = $1;\n' > query.sql
echo; echo "### ifoda-turlari"; sigs
sed -i 's|sql_package: pgx/v5|sql_package: pgx/v5\n        emit_empty_slices: true|' sqlc.yaml
printf -- '-- name: L :many\nSELECT id FROM products;\n' > query.sql
echo; echo "### emit_empty_slices"; sqlc generate && grep -n "items :=\|var items" db/query.sql.go
cp query.sql.orig query.sql; sed -i '/emit_empty_slices/d' sqlc.yaml; sqlc generate
echo; echo "### diff-toza"; sqlc diff; echo "[exit $?]"
printf -- '\n-- name: X :one\nSELECT 1;\n' >> query.sql
echo; echo "### diff-farq"; sqlc diff >/dev/null 2>&1; echo "sqlc diff → [exit $?]"
cp query.sql.orig query.sql
# OrderTotal: COALESCE'siz SUM — sqlc int64 deydi, Postgres NULL qaytaradi
sed -i 's|SELECT COALESCE(SUM(quantity \* unit_price), 0)::bigint AS total|SELECT SUM(quantity * unit_price) AS total|' query.sql
sqlc generate && grep -h "^func (q \*Queries) OrderTotal" db/query.sql.go
step sum-null go test -run TestOrderTotal .
cp query.sql.orig query.sql
sed -i 's|), 0)::bigint AS total|), 0) AS total|' query.sql
sqlc generate && grep -h "^func (q \*Queries) OrderTotal" db/query.sql.go
step coalesce-vet go vet .
cd "$W"

step down docker compose down -v
