#!/usr/bin/env bash
# Videodagi HAMMA chiqishni qaytadan oladi va _video/chiqishlar.txt'ga yozadi (go1.27.4, CGO_ENABLED=0;
# -race qismlari Docker'da: race.sh). Server :8765'da ishga tushadi (band boʻlsa — toʻxtaydi) va oxirida oʻchiriladi.
#   ./_video/record.sh
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
out="$here/_video/chiqishlar.txt"
export GOTOOLCHAIN=${GOTOOLCHAIN:-go1.27.4} CGO_ENABLED=0
if ss -ltn | grep -q ':8765 '; then echo ":8765 band — avval uni boʻshating" >&2; exit 1; fi
t=$(mktemp -d)/bozor && mkdir -p "$t"
(cd "$here" && cp -r go.mod bozor.json cmd internal "$t/") && find "$t" -name tekshiruv_test.go -delete
{
echo "# GB-4.08 videosidagi chiqishlar · $(go version) · $(curl --version | head -1 | cut -d' ' -f1-2) · $(date -u +%F)"
echo
echo '## tree-run (final)'; echo '$ ls cmd/api internal/store'; (cd "$t" && ls -C cmd/api internal/store)
echo
echo '## server-run … quiz-run: 1-terminal go run ./cmd/api, 2-terminal curl'
(cd "$t" && go build -o /dev/null ./cmd/api)          # kompilyatsiya keshga, server tez koʻtariladi
(cd "$t" && exec go run ./cmd/api > "$t/server.log" 2>&1) &
gopid=$!
for _ in $(seq 50); do ss -ltn | grep -q ':8765 ' && break; sleep 0.2; done
echo '$ go run ./cmd/api'; cat "$t/server.log"
# buyruq qatori videodagidek chiqadi; --no-progress-meter faqat bu yerda: chiqish faylga ketganda curl progress
# jadvalini chiqaradi, terminalda esa chiqarmaydi (videodagi holat)
c() {
	local a q=""
	for a in "$@"; do case $a in *[\{\ ]*) q+=" '$a'" ;; *) q+=" $a" ;; esac; done
	echo "\$ curl$q"; curl --no-progress-meter "$@"; echo
}
c -s localhost:8765/products/4
c -i localhost:8765/orders -d '{"product_id":4,"quantity":5}'
c -i localhost:8765/orders -d '{"product_id":4,"quantity":2}'
c -s localhost:8765/products/4
c -i localhost:8765/narsalar
c -i -X DELETE localhost:8765/products/1
echo '# videoda yoʻq (README):'
c -s localhost:8765/products -d '{"name":"","price":-5}'
c -s localhost:8765/products/abc
c -i localhost:8765/products -d '{"name":"anor","price":1800000,"stock":7}'
pid=$(ss -ltnp | grep ':8765 ' | grep -o 'pid=[0-9]*' | head -1 | cut -d= -f2)
[ -n "$pid" ] && kill -INT "$pid"; kill "$gopid" 2>/dev/null; wait "$gopid" 2>/dev/null
sleep 0.3; ss -ltn | grep -q ':8765 ' && echo "OGOHLANTIRISH: :8765 hali band" || echo "# server toʻxtadi, :8765 boʻsh"
echo
echo '## test-run (final)'; echo '$ go test ./...'; (cd "$t" && go test -count=1 ./...)
echo
echo '## bug-kod: 01-pointer bilan oddiy go test ham oʻtadi'; "$here/_video/run.sh" 01-pointer go test -count=1 ./...
echo
echo '## race-run, race-fail (01-pointer, Docker)'; echo '$ go test -race ./cmd/api'
"$here/_video/race.sh" 01-pointer go test -race -count=1 ./cmd/api 2>&1
echo
echo '## 01-pointer: race 5 marta ishga tushirildi (topilgan hisobotlar soni)'
for i in 1 2 3 4 5; do "$here/_video/race.sh" 01-pointer go test -race -count=1 ./cmd/api 2>&1 | grep -c 'WARNING: DATA RACE'; done
echo
echo '## tuzat-run (final, Docker)'; echo '$ go test -race ./...'; "$here/_video/race.sh" final go test -race -count=1 ./... 2>&1
echo '(~1 s: race runtime chiqishdan oldin atexit_sleep_ms=1000 kutadi; GORACE=atexit_sleep_ms=0 bilan ~0.01 s)'
echo
echo '## solution, -race (Docker)'; "$here/_video/race.sh" solution go test -race -count=1 ./... 2>&1
} > "$out" 2>&1
echo "yozildi: $out"
