#!/usr/bin/env bash
# GB-6.08 videosidagi buyruqlar va ularning haqiqiy chiqishlari.
# Ishga tushirish (shu papkada): bash _video/chiqishlar.sh > _video/chiqishlar.txt
# Kerak: Docker, curl, goose (go install github.com/pressly/goose/v3/cmd/goose@v3.30.2), 8765 va 5518 portlari boʻsh.
# Server videoda `go run ./cmd/api` bilan ishga tushadi; skript esa xuddi shu kodni `go build` qilib, binary'ni
# fonda ishga tushiradi va oxirida toʻxtatadi (chiqishi bir xil). curl -s: chiqish faylga ketgani uchun progress
# satri chiqmasin; terminalda `curl -i` aynan shu qatorlarni chiqaradi. Video `bozor` nomli papkada yozilgan.
# Bosqich (_video/stages/bola): store.Order egasini tekshirmaydi (BOLA) + bitta foydalanuvchili test.
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
lab=$(pwd)
work=$(mktemp -d)/bozor
mkdir -p "$work"
cp -r cmd internal migrations go.mod go.sum migrate.go seed.sql compose.yaml "$work"/
cd "$work"
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
json='-H "Content-Type: application/json"'

section "1. baza, migratsiyalar, seed, server"
docker compose down -v >/dev/null 2>&1
show docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5518/bozor
show goose up | sed -E 's#^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9:]{8} ##'
show docker compose exec -T db psql -f /lab/seed.sql
export DATABASE_URL=$GOOSE_DBSTRING
start_api
echo '$ go run ./cmd/api'
head -1 server.log

section "2. roʻyxatdan oʻtish va bazadagi parol"
show "curl -s -i localhost:8765/users $json -d '{\"email\":\"aziza@bozor.test\",\"password\":\"anor-behi-yongoq-uzum-5\"}'"
echo
show "docker compose exec -T db psql -c \"SELECT id, email, left(password_hash, 31) FROM users WHERE id = 4\""

section "3. login, begona buyurtma, admin yoʻli"
show "curl -s -c aziza.txt localhost:8765/login $json -d '{\"email\":\"aziza@bozor.test\",\"password\":\"anor-behi-yongoq-uzum-5\"}'"
show curl -s -i -b aziza.txt localhost:8765/orders/41
echo
show "curl -s -i -b aziza.txt localhost:8765/admin/products $json -d '{\"name\":\"anor\",\"price\":1800000,\"stock\":5}'"
echo
show curl -s -i localhost:8765/admin/yoq
echo
show curl -s -i -b aziza.txt localhost:8765/admin/yoq
echo

section "4. oʻz buyurtmasi, logout, eski cookie"
show "curl -s -i -b aziza.txt localhost:8765/orders $json -d '{\"items\":[{\"product_id\":3,\"quantity\":1}]}'"
echo
show curl -s -i -b aziza.txt localhost:8765/orders/43
echo
show curl -s -i -b aziza.txt -X POST localhost:8765/logout
echo
show cat aziza.txt
show curl -s -i -b aziza.txt localhost:8765/orders/43
echo

section "5. login'ga rate limit: 5 notoʻgʻri parol, 6-chisi toʻgʻri"
sleep 13 # Aziza'ning login'i olgan token qaytsin (chelak: 5 ta, har 12 soniyada bittadan)
show "for p in x1 x2 x3 x4 x5 anor-behi-yongoq-uzum-5; do curl -s -o /dev/null -w \"%{http_code} \" localhost:8765/login $json -d '{\"email\":\"aziza@bozor.test\",\"password\":\"'\$p'\"}'; done; echo"
show "curl -s -i localhost:8765/login $json -d '{\"email\":\"aziza@bozor.test\",\"password\":\"anor-behi-yongoq-uzum-5\"}'"
echo
echo "--- server log (oxirgi qatorlar):"
tail -4 server.log
stop_api

section "6. testlar: papkaning kodi (videoning oxiri = starter)"
show "go test -v -run 'TestAccess\$|TestPasswordNotStored|TestLogoutRevokes|TestLoginRateLimit' ./cmd/api"
show go test ./... | sed -E 's#\t[0-9.]+s$##'

section "7. XATO: bitta foydalanuvchili test va BOLA'li store"
cp "$lab/_video/stages/bola/orders.go" internal/store/orders.go
cp "$lab/_video/stages/bola/bitta_test.go" cmd/api/bitta_test.go
show go test -run TestOrdersOneUser ./cmd/api
show "go test -run 'TestAccess\$' ./cmd/api"
cp orders.asl internal/store/orders.go
rm cmd/api/bitta_test.go
show "go test -run 'TestAccess\$' ./cmd/api"

section "8. vazifa: solution"
cp -r "$lab/solution/cmd" "$lab/solution/internal" .
show go test ./... | sed -E 's#\t[0-9.]+s$##'

section "tozalash"
docker compose down -v >/dev/null 2>&1 && echo "docker compose down -v: ok"
ss -ltn | grep -q ':8765 ' && echo "8765 hali band!" || echo "8765 boʻsh"
