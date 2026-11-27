#!/usr/bin/env bash
# _video/record.sh — videodagi hamma buyruqni shu kompyuterda ishga tushirib, chiqishlar.txt'ga yozadi.
# Server :8765 da (boʻsh ekani tekshiriladi), oxirida toʻxtatiladi. go1.27.4, CGO_ENABLED=0, curl.
set -uo pipefail
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
here=$(cd "$(dirname "$0")" && pwd)
out=$here/chiqishlar.txt
: > "$out"
log() { printf '%s\n' "$*" >> "$out"; }
run() { log "\$ $*"; eval "$@" >> "$out" 2>&1; log ""; }
free() { ! ss -ltn | grep -q ":$1 "; }

start() { # start STEP: binary'ni quradi va :8765'da ishga tushiradi
  local d; d=$("$here/run.sh" "$1")
  (cd "$d" && go build -o "$d/api" ./cmd/api)
  free 8765 || { echo ":8765 band" >&2; exit 1; }
  (cd "$d" && exec ./api > "$d/server.log" 2>&1) & SRV=$!
  for _ in $(seq 50); do curl -s localhost:8765/health >/dev/null && break; sleep 0.1; done
  SRVDIR=$d
}
stop() { kill "$SRV"; wait "$SRV" 2>/dev/null; log "--- 1-terminal (server log):"; cat "$SRVDIR/server.log" >> "$out"; log ""; free 8765 && log "(:8765 boʻsh)"; log ""; }

log "### go version"; run go version; run curl --version \| head -1

log "### oldin (GB-6.01 solution)"
start oldin
run "curl -s localhost:8765/products -H \"Content-Type: application/json\" -d '{\"name\":\"choy\",\"price\":\"25000\"}'"
run "curl -s localhost:8765/products -H \"Content-Type: application/json\" -d '[1]'"
run "curl -s -i localhost:8765/narsalar"
stop

log "### final (= starter)"
start final
run "curl -s localhost:8765/products -H \"Content-Type: application/json\" -d '{\"name\":\"choy\",\"price\":\"25000\"}'"
run "curl -s -i localhost:8765/products -H \"Content-Type: application/json\" -d '{\"name\":\"choy\",\"price\":\"25000\"}'"
run "curl -s localhost:8765/products -H \"Content-Type: application/json\" -d '[1]'"
run "curl -s localhost:8765/orders -H \"Content-Type: application/json\" -d '{\"product_id\":4,\"quantity\":5}'"
run "curl -s -i localhost:8765/orders -H \"Content-Type: application/json\" -d '{\"product_id\":4,\"quantity\":5}'"
run "curl -s localhost:8765/products/99"
run "curl -s localhost:8765/narsalar"
run "curl -s -i -X PUT localhost:8765/products/1"
run "curl -s localhost:8765/products -H \"Content-Type: application/json\" -d '{\"name\":\"\",\"price\":0}'"
stop

log "### unutgan (return yoʻq)"
start unutgan
run "curl -s localhost:8765/products/99"
stop

log "### testlar"
d=$("$here/run.sh" final)
run "cd $d && go test ./... 2>&1 | grep -v '^    \\|^\\s*$' "
run "cd $d && go test ./..."
d=$("$here/run.sh" solution)
run "cd $d && go test ./..."

log "### recall: panic, recovery yoʻq (:8766)"
free 8766 || { echo ":8766 band" >&2; exit 1; }
p=$(mktemp -d); cp "$here/panic/main.go" "$p/" && (cd "$p" && printf 'module panic\n\ngo 1.27\n' > go.mod && go build -o boom .)
("$p/boom" > "$p/log" 2>&1) & B=$!
for _ in $(seq 50); do ss -ltn | grep -q ":8766 " && break; sleep 0.1; done
run "curl -i localhost:8766/boom"
kill $B; wait $B 2>/dev/null
log "--- server log:"; head -3 "$p/log" >> "$out"
free 8766 && log "(:8766 boʻsh)"
