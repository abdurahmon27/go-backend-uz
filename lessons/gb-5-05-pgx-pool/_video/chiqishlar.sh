#!/usr/bin/env bash
# GB-5.05: videodagi har buyruqni ketma-ket ishga tushiradi va chiqishlarini chiqishlar.txt'ga yozadi.
# Ishlatish: GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0 ./chiqishlar.sh > chiqishlar.txt 2>&1
# Vaqtinchalik "bozor" papkasida ishlaydi (compose loyiha nomi = papka nomi, videodagidek) va oxirida docker compose down -v bilan tozalaydi.
set -u
HERE="$(cd "$(dirname "$0")" && pwd)"
LAB="$HERE/.."
TMP="$(mktemp -d)"
W="$TMP/bozor"
mkdir -p "$W"
cd "$W"
grep -v '^name: gb505' "$LAB/compose.yaml" > compose.yaml
cp "$LAB/products.sql" "$LAB/seed.sql" "$LAB/orders.sql" "$LAB/orders_seed.sql" .
printf 'module github.com/abdurahmon27/bozor\n\ngo 1.27\n' > go.mod
step() { echo; echo "### $1"; shift; echo "\$ $*"; "$@" 2>&1; echo "[exit $?]"; }
sh_step() { echo; echo "### $1"; echo "\$ $2"; bash -c "$2" 2>&1; echo "[exit $?]"; }
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5475/bozor

step up docker compose up -d --wait
step jadval docker compose exec db psql -f /lab/products.sql -f /lab/seed.sql -f /lab/orders.sql -f /lab/orders_seed.sql
step select docker compose exec db psql -c 'SELECT id, name, price, stock FROM products WHERE active ORDER BY id'

cp "$HERE/stages/01-main-pool.go" main.go
step tidy go mod tidy
echo; echo '### go.mod'; cat go.mod

cp "$HERE/stages/02-store.go" store.go; cp "$HERE/stages/03-main-get1.go" main.go
step get1 go run .
cp "$HERE/stages/04-main-get99.go" main.go
step get99 go run .
cp "$HERE/stages/05-store.go" store.go
step get99b go run .

cp "$HERE/stages/07-store.go" store.go; cp "$HERE/stages/08-main-final.go" main.go
step list go run .
cp "$HERE/stages/07a-store-sum.go" store.go
step sumbug go run .
cp "$HERE/stages/07-store.go" store.go
step sumfix go run .

cp "$HERE/stages/06-main-slow.go" main.go
step slow go run .

mkdir leak && cp "$HERE/leak/main.go" leak/main.go
sh_step leak "export DATABASE_URL='postgres://bozor:bozor-lab@localhost:5475/bozor?pool_max_conns=1'; go run ./leak"
mkdir leak-fixed && cp "$HERE/leak-fixed/main.go" leak/main.go
sh_step leakfixed "export DATABASE_URL='postgres://bozor:bozor-lab@localhost:5475/bozor?pool_max_conns=1'; go run ./leak"

cp "$HERE/stages/08-main-final.go" main.go
sh_step lazy "DATABASE_URL=postgres://bozor:bozor-lab@localhost:5499/bozor go run ."

step down docker compose down -v
cd /; rm -rf "$TMP"
