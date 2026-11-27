#!/usr/bin/env bash
# GB-6.04 videosidagi buyruqlar va ularning haqiqiy chiqishlari.
# Ishga tushirish (shu papkada): bash _video/chiqishlar.sh > _video/chiqishlar.txt
# Kerak: Docker, curl, goose (go install github.com/pressly/goose/v3/cmd/goose@v3.30.2), 8765 va 5514 portlari boʻsh.
# Server videoda `go run ./cmd/api` bilan ishga tushadi; skript esa xuddi shu kodni `go build` qilib, binary'ni
# fonda ishga tushiradi va oxirida toʻxtatadi (chiqishi bir xil). curl -s: chiqish faylga ketgani uchun progress
# satri chiqmasin; terminalda `curl -i` aynan shu qatorlarni chiqaradi.
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
lab=$(pwd)
work=$(mktemp -d)/bozor
mkdir -p "$work"
cp -r cmd internal migrations go.mod go.sum migrate.go seed.sql compose.yaml "$work"/
cd "$work"

show() { echo "\$ $*"; eval "$@" 2>&1; }
section() { echo; echo "### $*"; }
start_api() {
	if ss -ltn | grep -q ':8765 '; then echo "8765 band"; exit 1; fi
	go build -o api ./cmd/api || exit 1
	./api >>server.log 2>&1 &
	api_pid=$!
	for _ in $(seq 50); do curl -s localhost:8765/health >/dev/null && break; sleep 0.1; done
}
stop_api() { kill "$api_pid"; wait "$api_pid" 2>/dev/null; }

section "1. baza va migratsiyalar"
docker compose down -v >/dev/null 2>&1
docker compose up -d --wait >/dev/null 2>&1
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5514/bozor
show goose up | sed -E 's#^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9:]{8} ##'
show docker compose exec -T db psql -f /lab/seed.sql
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5514/bozor
start_api
echo '$ go run ./cmd/api'
head -1 server.log

section "2. cookie'siz /me"
show curl -s -i localhost:8765/me

section "3. notoʻgʻri parol va yoʻq email"
show "curl -s localhost:8765/login -H \"Content-Type: application/json\" -d '{\"email\":\"dilnoza@bozor.test\",\"password\":\"tarvuz-osmon-qalam-ariq-28\"}'"
show "curl -s localhost:8765/login -H \"Content-Type: application/json\" -d '{\"email\":\"yoq@bozor.test\",\"password\":\"tarvuz-osmon-qalam-ariq-27\"}'"

section "3b. javob vaqti: notoʻgʻri parol va yoʻq email (dummyHash bilan), 5 tadan"
for body in '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-28"}' '{"email":"yoq@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'; do
	echo "$body"
	for _ in 1 2 3 4 5; do curl -s -o /dev/null -w '%{time_total}s ' localhost:8765/login -H "Content-Type: application/json" -d "$body"; done
	echo
done

section "4. login: telefon (jar.txt) va noutbuk (cookie saqlanmaydi)"
show "curl -s -i -c jar.txt localhost:8765/login -H \"Content-Type: application/json\" -d '{\"email\":\"dilnoza@bozor.test\",\"password\":\"tarvuz-osmon-qalam-ariq-27\"}'"
echo
show "curl -s localhost:8765/login -H \"Content-Type: application/json\" -d '{\"email\":\"dilnoza@bozor.test\",\"password\":\"tarvuz-osmon-qalam-ariq-27\"}'"
show "docker compose exec -T db psql -c \"SELECT encode(token_hash, 'hex') AS token_hash, user_id, expires_at FROM sessions\""
show "grep session jar.txt | cut -f7"

section "5. cookie bilan /me, keyin server qayta ishga tushadi"
show curl -s -b jar.txt localhost:8765/me
stop_api
start_api
echo '$ go run ./cmd/api   (qayta)'
tail -1 server.log
show curl -s -b jar.txt localhost:8765/me

section "6. bazadagi hash'ni cookie qilib yuborish (baza sizib chiqdi)"
leaked=$(docker compose exec -T db psql -At -c "SELECT encode(token_hash, 'hex') FROM sessions LIMIT 1")
show curl -s -b "session=$leaked" localhost:8765/me
stop_api

section "7. klassik xato: muddat tekshirilmaydi (_video/stages/store-muddatsiz.go)"
cp internal/store/store.go store.go.asl
cp "$lab/_video/stages/store-muddatsiz.go" internal/store/store.go
start_api
show "docker compose exec -T db psql -c \"UPDATE sessions SET expires_at = now() - interval '1 minute'\""
show curl -s -b jar.txt localhost:8765/me
stop_api

section "8. tuzatish: shart qaytdi, server qayta ishga tushdi"
cp store.go.asl internal/store/store.go
start_api
show curl -s -b jar.txt localhost:8765/me

section "9. logout"
show "curl -s -c jar.txt localhost:8765/login -H \"Content-Type: application/json\" -d '{\"email\":\"dilnoza@bozor.test\",\"password\":\"tarvuz-osmon-qalam-ariq-27\"}'"
show cp jar.txt eski.txt
show curl -s -i -X POST -b jar.txt -c jar.txt localhost:8765/logout
show curl -s -b eski.txt localhost:8765/me
show "docker compose exec -T db psql -c \"SELECT count(*) FROM sessions WHERE expires_at > now()\""
stop_api

section "10. server log'i"
cat server.log

section "11. testlar (vazifadan oldin)"
show go test ./... | grep -v '^?'

docker compose down -v >/dev/null 2>&1
ss -ltn | grep -q ':8765 ' && echo "8765 hali band!" || echo "(8765 boʻsh, konteyner va volume oʻchirildi)"
