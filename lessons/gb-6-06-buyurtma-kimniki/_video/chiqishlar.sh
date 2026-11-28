#!/usr/bin/env bash
# GB-6.06 videosidagi buyruqlar va ularning haqiqiy chiqishlari.
# Ishga tushirish (shu papkada): bash _video/chiqishlar.sh > _video/chiqishlar.txt
# Kerak: Docker, curl, goose (go install github.com/pressly/goose/v3/cmd/goose@v3.30.2), 8765 va 5516 portlari boʻsh.
# Server videoda `go run ./cmd/api` bilan ishga tushadi; skript esa xuddi shu kodni `go build` qilib, binary'ni
# fonda ishga tushiradi va oxirida toʻxtatadi (chiqishi bir xil). curl -s: chiqish faylga ketgani uchun progress
# satri chiqmasin; terminalda `curl` aynan shu qatorlarni chiqaradi.
# Bosqichlar (_video/stages): v0-bola — Order egasini tekshirmaydi; v1-admin-yoq — /admin/users'da requireAdmin yoʻq;
# v2-teskari — requireAdmin requireUser'dan tashqarida. Boshqa hamma joyda papkaning oʻz kodi (videoning oxiri).
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
lab=$(pwd)
work=$(mktemp -d)/bozor
mkdir -p "$work"
cp -r cmd internal migrations go.mod go.sum migrate.go seed.sql compose.yaml "$work"/
cd "$work"
cp cmd/api/handlers.go handlers.asl
cp internal/store/orders.go orders.asl

show() { echo "\$ $*"; eval "$@" 2>&1; }
section() { echo; echo "### $*"; }
start_api() {
	if ss -ltn | grep -q ':8765 '; then echo "8765 band"; exit 1; fi
	go build -o api ./cmd/api || exit 1
	: >server.log
	./api >>server.log 2>&1 &
	api_pid=$!
	for _ in $(seq 50); do curl -s localhost:8765/health >/dev/null && break; sleep 0.1; done
}
stop_api() { kill "$api_pid"; wait "$api_pid" 2>/dev/null; }
stage() { # stage <nom>: bosqich fayllarini qoʻyadi; "asl" — papkaning oʻz kodi
	cp handlers.asl cmd/api/handlers.go
	cp orders.asl internal/store/orders.go
	if [ "$1" != asl ]; then
		cp "$lab/_video/stages/$1/handlers.go" cmd/api/handlers.go
		[ -f "$lab/_video/stages/$1/orders.go" ] && cp "$lab/_video/stages/$1/orders.go" internal/store/orders.go
	fi
}
login_json() { echo "'{\"email\":\"$1@bozor.test\",\"password\":\"$2\"}'"; }

section "1. baza, migratsiyalar, seed"
docker compose down -v >/dev/null 2>&1
docker compose up -d --wait >/dev/null 2>&1
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5516/bozor
show goose up | sed -E 's#^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9:]{8} ##'
show docker compose exec -T db psql -f /lab/seed.sql
show "docker compose exec -T db psql -c \"SELECT id, user_id, note FROM orders\""
show "docker compose exec -T db psql -c \"SELECT id, email, role FROM users\""
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5516/bozor

section "2. XATO versiya (v0-bola): Dilnoza login, oʻz buyurtmasi 41, keyin 42"
stage v0-bola
start_api
echo '$ go run ./cmd/api'
head -1 server.log
show "curl -s -c jar.txt localhost:8765/login -H \"Content-Type: application/json\" -d $(login_json dilnoza tarvuz-osmon-qalam-ariq-27)"
show curl -s -b jar.txt localhost:8765/orders/41
show curl -s -b jar.txt localhost:8765/orders/42
show curl -s -b jar.txt localhost:8765/orders/43
stop_api
section "2b. XATO versiyada test"
show go test -run TestOrderOwner ./cmd/api

section "3. tuzatilgan versiya: 42 va 999"
stage asl
start_api
show curl -s -i -b jar.txt localhost:8765/orders/42
echo
show curl -s -i -b jar.txt localhost:8765/orders/999
echo
echo "--- server log (xato qatorlari):"
grep -E ' xato |GET' server.log
show curl -s -b jar.txt localhost:8765/orders/41
stop_api
section "3b. tuzatilgan versiyada test"
show go test -run TestOrderOwner ./cmd/api
show go test -run TestOrderOwner -v ./cmd/api

section "4. XATO versiya (v1-admin-yoq): Dilnoza /admin/users"
stage v1-admin-yoq
start_api
show curl -s -b jar.txt localhost:8765/admin/users
stop_api

section "5. tuzatilgan: Dilnoza 403, admin 200"
stage asl
start_api
show curl -s -i -b jar.txt localhost:8765/admin/users
echo
show "curl -s -c admin.txt localhost:8765/login -H \"Content-Type: application/json\" -d $(login_json admin qovun-anjir-tut-uzum-31)"
show curl -s -b admin.txt localhost:8765/admin/users
stop_api

section "6. XATO versiya (v2-teskari): admin /admin/users"
stage v2-teskari
start_api
show curl -s -i -b admin.txt localhost:8765/admin/users
echo
echo "--- server log (panic qatori):"
grep -E 'panic|GET' server.log | grep -v goroutine | head -3
stop_api

section "7. Chuqurroq: body'da user_id"
stage asl
start_api
show "curl -s -b jar.txt localhost:8765/orders -H \"Content-Type: application/json\" -d '{\"user_id\":2,\"items\":[{\"product_id\":3,\"quantity\":1}]}'"
show "curl -s -b jar.txt localhost:8765/orders -H \"Content-Type: application/json\" -d '{\"items\":[{\"product_id\":3,\"quantity\":1}]}'"
stop_api

section "8. vazifa: starter'da go test ./..."
show go test ./... | sed -E 's#\t[0-9.]+s$##'

section "tozalash"
docker compose down -v >/dev/null 2>&1 && echo "docker compose down -v: ok"
ss -ltn | grep -q ':8765 ' && echo "8765 hali band!" || echo "8765 boʻsh"
