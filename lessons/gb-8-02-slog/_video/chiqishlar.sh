#!/usr/bin/env bash
# GB-8.02 videosidagi har buyruqni shu kompyuterda qayta ishga tushiradi; chiqishlar chiqishlar.txt'da.
# Ishga tushirish (papka ildizidan): bash _video/chiqishlar.sh > _video/chiqishlar.txt 2>&1
# Kerak: Docker (compose loyihasi gb802, 127.0.0.1:5582), goose (go install github.com/pressly/goose/v3/cmd/goose@v3.30.2),
# curl, boʻsh 8782-port. Oxirida server va Postgres toʻxtatiladi (docker compose down -v).
# Bosqichlar (videodagi tartib):
#   0. eski  — GB-8.01'ning kodi (../gb-8-01-config/solution): log.Printf, xato matni toʻliq (hook'dagi ikki qator)
#   1. a     — main.go'da slog, config'da LogValue hali yoʻq ("config":{})
#   2. b     — LogValue bor; requestID logger'ni context'ga qoʻyadi; writeError hali err'ni toʻliq yozadi (sizish)
#   3. final — papkaning oʻzi: writeError 4xx'da faqat detail, store xatolarida email yoʻq
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0 PATH=$HOME/go/bin:$PATH
src=$(cd "$(dirname "$0")/.." && pwd)
t=$(mktemp -d)/bozor
mkdir -p "$t" && cp -r "$src"/. "$t"/ && rm -rf "$t"/solution "$t"/_video "$t"/.env
cd "$t"
unset DATABASE_URL PORT HOST
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5582/bozor
echo "== go version"; go version

echo; echo "== Postgres, migratsiyalar, seed"
docker compose -p gb802 down -v >/dev/null 2>&1
docker compose -p gb802 up -d --wait 2>&1 | tail -1
GOOSE_DRIVER=postgres GOOSE_DBSTRING=$DATABASE_URL GOOSE_MIGRATION_DIR=migrations goose up 2>&1 | tail -1
docker compose -p gb802 exec -T db psql -q -f /lab/seed.sql 2>&1

if ss -ltn | grep -q ':8782 '; then echo "8782 band"; exit 1; fi
start() { # $1 — sarlavha; server fonda, stdout → api.log (videoda: ./api -port 8782 > api.log)
	echo; echo "== $1"
	go build ./cmd/api || exit 1
	./api -port 8782 >api.log &
	pid=$!
	for _ in $(seq 100); do ss -ltn | grep -q ":8782 " && break; sleep 0.1; done
}
stop() { kill -TERM "$pid"; wait "$pid" 2>/dev/null; }
show() { echo; echo "\$ $1"; bash -c "$1" 2>&1; }
J="-H 'Content-Type: application/json'"
EMAIL="curl -s -o /dev/null -w '%{http_code} %header{x-request-id}\n' localhost:8782/login $J -d '{\"email\":\"nodira@bozor.test\",\"password\":\"x\"}'"
PAROL="curl -s -o /dev/null -w '%{http_code} %header{x-request-id}\n' localhost:8782/login $J -d '{\"email\":\"jasur@bozor.test\",\"password\":\"Samarqand\\ulugbek1394\"}'"

# ---- 0. eski: GB-8.01 kodi
mkdir -p ../eski && cp -r "$src"/../gb-8-01-config/solution/. ../eski/ && (
	cd ../eski && sed -i 's/5581/5582/' internal/pgtest/pgtest.go
	go build ./cmd/api || exit 1
	./api -port 8782 >api.log 2>&1 &
	pid=$!
	for _ in $(seq 100); do ss -ltn | grep -q ":8782 " && break; sleep 0.1; done
	echo; echo "== 0. eski (GB-8.01): log.Printf"
	show "$EMAIL"
	show "$PAROL"
	kill -TERM "$pid"; wait "$pid" 2>/dev/null
	show "cat api.log"
)

# ---- 1. a: slog, LogValue yoʻq
cp cmd/api/config.go ../config.go.final
python3 - <<'PY'
p='cmd/api/config.go'; s=open(p).read()
i=s.index('// LogValue —'); s=s[:i].rstrip()+'\n'
s=s.replace('\t"log/slog"\n','')
open(p,'w').write(s)
PY
start "1. a: slog, LogValue yoʻq"
echo '$ go build ./cmd/api && ./api -port 8782 > api.log'
stop
show "cat api.log"

# ---- 2. b: LogValue, request_id, writeError err'ni toʻliq yozadi
cp ../config.go.final cmd/api/config.go
cp cmd/api/problem.go ../problem.go.final; cp internal/store/store.go ../store.go.final
python3 - <<'PY'
import re
p='cmd/api/problem.go'; s=open(p).read()
i=s.index('func writeError('); j=s.index('\n}\n', i)+3
s=s[:i]+'''func writeError(w http.ResponseWriter, r *http.Request, err error) {
	p := problemFor(err)
	loggerFrom(r.Context()).Warn("xato", "status", p.Status, "err", err)
	writeProblem(w, p)
}
'''+s[j:]
open(p,'w').write(s)
p='internal/store/store.go'; s=open(p).read()
s=s.replace('return User{}, ErrUserNotFound','return User{}, fmt.Errorf("email %q: %w", email, ErrUserNotFound)')
open(p,'w').write(s)
PY
start "2. b: request_id bor, writeError err'ni toʻliq yozadi"
echo '$ go build ./cmd/api && ./api -port 8782 > api.log'
show "cat api.log"
show "$EMAIL"
id=$(grep -o '"request_id":"[A-Z0-9]*"' api.log | tail -1 | cut -d'"' -f4)
show "grep $id api.log"
show "$PAROL"
id=$(grep -o '"request_id":"[A-Z0-9]*"' api.log | tail -1 | cut -d'"' -f4)
show "grep $id api.log"
stop

# ---- 3. final
cp ../problem.go.final cmd/api/problem.go; cp ../store.go.final internal/store/store.go
start "3. final: 4xx — faqat detail"
echo '$ go build ./cmd/api && ./api -port 8782 > api.log'
show "$EMAIL"
show "$PAROL"
stop
show "cat api.log"
show "grep WARN api.log"
show "grep -c -e nodira -e ulugb api.log"

# ---- 4. jq: maydon boʻyicha filtr (Chuqurroq)
show "jq -c 'select(.level == \"WARN\") | {msg, status, detail}' api.log"

# ---- 5. go vet: juft boʻlmagan kalit (keng tarqalgan xato)
cat > cmd/api/vet_demo.go <<'GO'
package main

import "log/slog"

func vetDemo(logger *slog.Logger) {
	logger.Info("soʻrov", "status")
}
GO
show "grep -n Info cmd/api/vet_demo.go"
show "go vet ./cmd/api"
rm cmd/api/vet_demo.go

echo; echo "== tozalash"
docker compose -p gb802 down -v 2>&1 | tail -1
