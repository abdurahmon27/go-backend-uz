#!/usr/bin/env bash
# GB-5.07: videodagi har buyruqni ketma-ket ishga tushiradi va chiqishlarini chiqishlar.txt'ga yozadi.
# Ishlatish: GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0 ./chiqishlar.sh > chiqishlar.txt 2>&1
# Vaqtinchalik "bozor" papkasida ishlaydi (compose loyiha nomi = papka nomi, videodagidek) va oxirida docker compose down -v bilan tozalaydi.
set -u
HERE="$(cd "$(dirname "$0")" && pwd)"
LAB="$HERE/.."
TMP="$(mktemp -d)"
W="$TMP/bozor"
mkdir -p "$W"
cd "$W"
grep -v '^name: gb507' "$LAB/compose.yaml" > compose.yaml
cp "$LAB/bozor.sql" "$LAB/bekor.sql" "$LAB/xato.sql" "$LAB/go.mod" "$LAB/go.sum" "$LAB/store.go" "$LAB/notx.go" .
step() { echo; echo "### $1"; shift; echo "\$ $*"; "$@" 2>&1; echo "[exit $?]"; }
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5497/bozor
PSQL=(docker compose exec db psql)

step up docker compose up -d --wait
step bozor "${PSQL[@]}" -f /lab/bozor.sql

# 1. transaction'siz: yarim yozilgan buyurtma
cp "$HERE/stages/01-main-notx.go" main.go
step notx go run .
step holat1 "${PSQL[@]}" -c 'SELECT id, customer FROM orders WHERE id > 3' -c 'SELECT order_id, product_id, quantity FROM order_items WHERE order_id > 3'
step holat2 "${PSQL[@]}" -c 'SELECT id, name, stock FROM products WHERE id <= 2 ORDER BY id'

# 2. psql'da transaction
step bekor "${PSQL[@]}" -t -f /lab/bekor.sql
step bekor-jadval "${PSQL[@]}" -f /lab/bekor.sql
step xato "${PSQL[@]}" -f /lab/xato.sql

# 3. bazani qaytarish va transaction bilan
step reset "${PSQL[@]}" -f /lab/bozor.sql
cp "$HERE/stages/02-main-tx-fail.go" main.go
step txfail go run .
step holat3 "${PSQL[@]}" -c 'SELECT id, customer FROM orders WHERE id > 3' -c 'SELECT order_id, product_id, quantity FROM order_items WHERE order_id > 3'
step holat4 "${PSQL[@]}" -c 'SELECT id, name, stock FROM products WHERE id <= 2 ORDER BY id'
cp "$HERE/stages/03-main-final.go" main.go
step txok go run .
step holat5 "${PSQL[@]}" -c 'SELECT id, customer FROM orders WHERE id > 3' -c 'SELECT order_id, product_id, quantity, unit_price FROM order_items WHERE order_id > 3'
step holat6 "${PSQL[@]}" -c 'SELECT id, name, stock FROM products WHERE id <= 2 ORDER BY id'

# 4. xato: item INSERT'i tx oʻrniga pool orqali
step reset2 "${PSQL[@]}" -f /lab/bozor.sql
cp store.go store.go.ok
cp "$HERE/stages/04-store-pool-mistake.go" store.go
step mistake go run .
step holat7 "${PSQL[@]}" -c 'SELECT id, name, stock FROM products WHERE id <= 2 ORDER BY id'
cp store.go.ok store.go; rm store.go.ok

# 5. Chuqurroq: xato kodini *pgconn.PgError dan oʻqish (store.go = videodagi PlaceOrder, ErrNotEnoughStock'siz)
step reset3 "${PSQL[@]}" -f /lab/bozor.sql
cp "$HERE/stages/05-main-pgerr.go" main.go
step pgerr go run .

step down docker compose down -v
cd /; rm -rf "$TMP"
