#!/usr/bin/env bash
# _video/record.sh [BOLIM…] — videodagi buyruqlarni shu kompyuterda ishga tushiradi (GB-6.03).
# Har boʻlim _video/out/<boʻlim>.txt'ga yoziladi, chiqishlar.txt — hammasining yigʻindisi.
# Boʻlimlar: meta goget sizdi taxmin server xato testlar (argumentsiz — hammasi).
# go1.27.4, CGO_ENABLED=0, curl. Server :8765 da (boʻsh ekani tekshiriladi), oxirida toʻxtatiladi.
# sizdi/ fayllari qayta yasalmaydi (videoda koʻringanlari shu): yangidan yasash — go run ./_video/sizdir.
set -uo pipefail
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
here=$(cd "$(dirname "$0")/.." && pwd)
dir=$here/_video/out; mkdir -p "$dir"
out=/dev/null
log() { printf '%s\n' "$*" >> "$out"; }
run() { log "\$ $*"; eval "$@" >> "$out" 2>&1; log "(exit $?)"; log ""; }
free() { ! ss -ltn | grep -q ":$1 "; }
copy() { local t; t=$(mktemp -d)/bozor; mkdir -p "$t"; cp -r "$here"/{go.mod,go.sum,bozor.json,cmd,internal,nomzodlar.txt,sizdi} "$t"/; echo "$t"; }
section() { out=$dir/$1.txt; : > "$out"; }

meta() { section 1-meta; log "### go version"; run go version; run "curl --version | head -1"; run "lscpu | grep 'Model name'"; }

goget() {
  section 2-goget
  log "### go get (GB-6.02 solution nusxasi, boʻsh modul keshi: yuklash qatorlari koʻrinsin)"
  local g; g=$(mktemp -d)/bozor; cp -r "$here/../gb-6-02-api-xatolar/solution" "$g"
  local tc; tc=$(go env GOROOT)/bin/go
  run "cd $g && GOTOOLCHAIN=local GOMODCACHE=$(mktemp -d) GOFLAGS=-modcacherw $tc get golang.org/x/crypto/argon2"
}

sizdi() {
  section 3-sizdi; local d; d=$(copy)
  log "### sizib chiqqan jadvallar"
  run "cd $d && cat sizdi/ochiq.txt"; run "cd $d && cat sizdi/sha256.txt"; run "cd $d && cat sizdi/argon2id.txt"
}

taxmin() {
  section 4-taxmin; local d; d=$(copy)
  log "### taxmin"
  run "cd $d && go run ./cmd/taxmin sizdi/sha256.txt"; run "cd $d && go run ./cmd/taxmin sizdi/argon2id.txt"
}

server() {
  section 5-server; local d; d=$(copy)
  log "### server: POST /users"
  (cd "$d" && go build -o "$d/api" ./cmd/api)
  free 8765 || { echo ":8765 band" >&2; exit 1; }
  (cd "$d" && exec ./api > "$d/server.log" 2>&1) & local srv=$!
  for _ in $(seq 50); do curl -s localhost:8765/health >/dev/null && break; sleep 0.1; done
  run "curl -s -i localhost:8765/users -H \"Content-Type: application/json\" -d '{\"email\":\"dilnoza@bozor.test\",\"password\":\"tarvuz-osmon-qalam-ariq-27\"}'"
  run "curl -s localhost:8765/users -H \"Content-Type: application/json\" -d '{\"email\":\"Dilnoza@bozor.test\",\"password\":\"boshqa-uzun-ibora-2025\"}'"
  run "curl -s localhost:8765/users -H \"Content-Type: application/json\" -d '{\"email\":\"jasur@bozor.test\",\"password\":\"choy2024\"}'"
  kill "$srv"; wait "$srv" 2>/dev/null
  log "--- 1-terminal (server log):"; cat "$d/server.log" >> "$out"; log ""
  free 8765 && log "(:8765 boʻsh)"; log ""
}

xato() {
  section 6-xato; local x; x=$(copy); cp "$here/_video/xato/password.go.txt" "$x/internal/password/password.go"
  log "### xato: decode'da base64.StdEncoding"
  run "cd $x && go test ./internal/password"
}

testlar() {
  section 7-testlar; local s; s=$(copy)
  log "### testlar"
  run "cd $s && go test ./... 2>&1 | grep -v '^    '"
  run "cd $s && go test ./internal/password"
  run "cd $here/solution && go test -count=1 ./..."
}

for b in "${@:-meta goget sizdi taxmin server xato testlar}"; do for x in $b; do "$x"; done; done
cat "$dir"/*.txt > "$here/_video/chiqishlar.txt"
