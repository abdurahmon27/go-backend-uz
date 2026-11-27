#!/usr/bin/env bash
# Videodagi barcha buyruqlarni ketma-ket, xuddi videodagidek ishga tushiradi va natijalarini chop etadi.
# Vaqtinchalik `bozor/` papkasida ishlaydi (compose loyihasi nomi = papka nomi, GB-5.02..5.08 videolaridagidek;
# shu uchun compose.yaml'dagi `name: gb510` qatori olib tashlanadi). Port 5510; oxirida `down -v` bilan tozalaydi.
# «final» = papkaning oʻzi, tekshiruv_test.go'larsiz (videoning oxiri, vazifadan oldin).
# Kerak: PATH'da goose v3.30.2 (go install github.com/pressly/goose/v3/cmd/goose@v3.30.2), GOTOOLCHAIN=go1.27.4,
# boʻsh 8765-port.
#   ./_video/chiqishlar.sh | tee _video/chiqishlar.txt
set -u
export GOTOOLCHAIN=${GOTOOLCHAIN:-go1.27.4} CGO_ENABLED=0
here=$(cd "$(dirname "$0")/.." && pwd)
t=$(mktemp -d) && mkdir -p "$t/bozor"
api=
cleanup() {
	[ -n "$api" ] && kill "$api" 2>/dev/null && wait "$api" 2>/dev/null
	cd "$t/bozor" && docker compose down -v >/dev/null 2>&1
	rm -rf "$t"
}
trap cleanup EXIT
if docker compose ls --format json | grep -q '"Name":"bozor"'; then echo "bozor compose loyihasi allaqachon bor; keyinroq urining" >&2; exit 1; fi
cd "$t/bozor" || exit 1
(cd "$here" && cp -r go.mod go.sum migrate.go seed.sql migrations cmd internal "$t/bozor/")
grep -v '^name: gb510' "$here/compose.yaml" > compose.yaml
find . -name tekshiruv_test.go -delete
run() { echo "\$ $*"; "$@" 2>&1; echo "(chiqish kodi: $?)"; echo; }

echo "### 1. papka"
run ls migrations
run ls cmd/api internal/pgtest internal/store
echo "### 2. baza, migratsiya, seed"
run docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5510/bozor
run goose up
run docker compose exec db psql -f /lab/seed.sql
echo "### 3. server"
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5510/bozor
go build -o /tmp/gb510-api ./cmd/api # go run ./cmd/api bilan bir xil, faqat jarayonni toʻxtatish oson
echo '$ go run ./cmd/api'
/tmp/gb510-api > "$t/api.log" 2>&1 &
api=$!
for _ in $(seq 50); do curl -s localhost:8765/health >/dev/null && break; sleep 0.1; done
cat "$t/api.log"; echo
run curl -i localhost:8765/orders -d '{"customer":"Aziza","items":[{"product_id":1,"quantity":2},{"product_id":3,"quantity":1}]}'
run curl -i localhost:8765/orders -d '{"customer":"Jasur","items":[{"product_id":1,"quantity":2},{"product_id":4,"quantity":5}]}'
run curl -s localhost:8765/products/1
run curl -i localhost:8765/orders -d '{"customer":"Jasur","items":[{"product_id":1,"quantity":0}]}'
run docker compose exec db psql -c "SELECT id, customer, status FROM orders WHERE id > 3"
run curl -s localhost:8765/orders/4
kill "$api"; wait "$api" 2>/dev/null; api=
echo "server log:"; cat "$t/api.log"; echo
echo "### 4. testlar (final)"
run go test ./...
echo "### 5. xato: hamma test bitta public sxemada (_video/stages/pgtest-shared.go)"
cp internal/pgtest/pgtest.go "$t/pgtest.go.orig"
cp "$here/_video/stages/pgtest-shared.go" internal/pgtest/pgtest.go
run go test -count=5 ./internal/store
run go test -count=5 ./...
cp "$t/pgtest.go.orig" internal/pgtest/pgtest.go
run go test -count=5 ./...
echo "### 6. vazifa: starter (tekshiruv_test.go bilan)"
(cd "$here" && cp cmd/api/tekshiruv_test.go "$t/bozor/cmd/api/" && cp internal/store/tekshiruv_test.go "$t/bozor/internal/store/")
run go test ./...
