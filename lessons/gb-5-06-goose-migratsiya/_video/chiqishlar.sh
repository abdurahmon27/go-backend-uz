#!/usr/bin/env bash
# Videodagi barcha buyruqlarni ketma-ket, xuddi videodagidek ishga tushiradi va natijalarini chop etadi.
# Vaqtinchalik `bozor/` papkasida ishlaydi (compose loyihasi nomi = papka nomi, GB-5.02..5.04 videolaridagidek;
# shu uchun compose.yaml'dagi `name: gb506` qatori olib tashlanadi). Port 5486; oxirida `down -v` bilan tozalaydi.
# Kerak: PATH'da goose v3.30.2 (go install github.com/pressly/goose/v3/cmd/goose@v3.30.2) va GOTOOLCHAIN=go1.27.4.
#   ./_video/chiqishlar.sh | tee _video/chiqishlar.txt
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
t=$(mktemp -d) && mkdir -p "$t/bozor/migrations"
trap 'cd "$t/bozor" && docker compose down -v >/dev/null 2>&1; rm -rf "$t"' EXIT
cd "$t/bozor" || exit 1
cp "$here/products.sql" "$here/orders.sql" "$here/seed.sql" .
cp "$here/migrations/00001_products.sql" "$here/migrations/00002_orders.sql" migrations/
grep -v '^name: gb506' "$here/compose.yaml" > compose.yaml
if docker compose ls --format json | grep -q '"Name":"bozor"'; then echo "bozor compose loyihasi allaqachon bor; keyinroq urining" >&2; exit 1; fi
run() { echo "\$ $*"; "$@" 2>&1; echo "(chiqish kodi: $?)"; echo; }
sql() { run docker compose exec db psql -c "$1"; }

echo "### 0. oldingi darslardagi holat: qoʻlda yaratilgan jadvallar (videoda koʻrsatilmaydi)"
docker compose up -d --wait >/dev/null 2>&1
docker compose exec db psql -f /lab/products.sql -f /lab/orders.sql 2>&1
echo
echo "### 1. muammo: fayl ikkinchi marta"
run docker compose exec db psql -f /lab/orders.sql
echo "### 2. baza tozalanadi"
run docker compose down -v
run docker compose up -d --wait
echo "### 3. goose"
run goose -version
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5486/bozor
run ls migrations
run goose status
run goose up
run goose status
sql '\dt'
sql 'SELECT id, version_id, is_applied FROM goose_db_version ORDER BY id'
run goose up
echo "### 4. 3-migratsiya: xato"
run docker compose exec db psql -f /lab/seed.sql
run goose create -s orders_status sql
run ls migrations
cat migrations/00003_orders_status.sql
cat > migrations/00003_orders_status.sql <<'SQL'
-- +goose Up
ALTER TABLE orders ADD COLUMN status text NOT NULL;

-- +goose Down
ALTER TABLE orders DROP COLUMN status;
SQL
run goose up
run goose status
sql 'SELECT * FROM orders LIMIT 1'
echo "### 5. tuzatish"
cp "$here/migrations/00003_orders_status.sql" migrations/
run goose up
sql 'SELECT id, customer, status FROM orders'
echo "### 6. down"
run goose down
run goose status
sql 'SELECT * FROM orders LIMIT 1'
run goose up
