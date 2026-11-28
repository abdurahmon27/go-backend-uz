#!/usr/bin/env bash
# GB-8.03 videosidagi har buyruqni shu kompyuterda qayta ishga tushiradi; chiqishlar chiqishlar.txt'da.
# Ishga tushirish (papka ildizidan): bash _video/chiqishlar.sh > _video/chiqishlar.txt 2>&1
# Kerak: Docker (compose loyihasi gb803, 127.0.0.1:5583), goose (go install github.com/pressly/goose/v3/cmd/goose@v3.30.2),
# curl, nc (OpenBSD netcat), boʻsh 8783-port. Oxirida server va Postgres toʻxtatiladi (docker compose down -v).
#
# Videoda uchta terminal: 1 — server (./api -port 8783 ; echo "kod: $?"), 2 — curl, 3 — pkill. Bu skript ularni
# fon jarayonlari bilan takrorlaydi (set -m: fon jarayoni SIGINT'ni eʼtiborsiz qoldirmasin, terminaldagidek).
# curl -sS: terminalda curl progress chiqarmaydi, faqat javob yoki xato; skriptda chiqish terminal emas, shuning uchun -sS.
# Server log'i (JSON) videoda qisqartirilgan: "time" va request_id oʻrnida «…».
#
# Bosqichlar (videodagi tartib):
#   0. eski   — GB-8.02'ning kodi (../gb-8-02-slog/solution) + namoyish uchun sekin /health: SIGTERM soʻrovni uzadi
#   1. yangi  — papkaning oʻzi + sekin /health: soʻrov tugaydi, yangisi rad etiladi; muddat tugashi; ikkinchi Ctrl+C
#   2. chek   — «go sendReceipt» (xato) va «a.bg.Go» (papkaning oʻzi)
#   3. timeout — ReadHeaderTimeout: sarlavhasini tugatmagan mijoz 5 s'da uziladi
#   4. vazifa — go test: starter yiqiladi, solution oʻtadi
set -u
set -m
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0 PATH=$HOME/go/bin:$PATH
src=$(cd "$(dirname "$0")/.." && pwd)
w=$(mktemp -d)
t=$w/bozor
mkdir -p "$t" && cp -r "$src"/. "$t"/ && rm -rf "$t"/_video "$t"/.env
cd "$t"
unset PORT HOST
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5583/bozor
echo "== go version"; go version

echo; echo "== Postgres, migratsiyalar, seed"
docker compose -p gb803 down -v >/dev/null 2>&1
docker compose -p gb803 up -d --wait 2>&1 | tail -1
GOOSE_DRIVER=postgres GOOSE_DBSTRING=$DATABASE_URL GOOSE_MIGRATION_DIR=migrations goose up 2>&1 | tail -1
docker compose -p gb803 exec -T db psql -q -f /lab/seed.sql 2>&1

if ss -ltn | grep -q ':8783 '; then echo "8783 band"; exit 1; fi

# sekin — /health'ga namoyish uchun vaqtincha kutish (videodagi «sekin-kod» sahnasi).
sekin() {
	python3 - "$1" <<'PY'
import sys
p=sys.argv[1]+'/cmd/api/handlers.go'; s=open(p).read()
old='''func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")'''
new='''func healthHandler(w http.ResponseWriter, r *http.Request) {
	d, _ := time.ParseDuration(r.FormValue("sekin")) // namoyish uchun, vaqtincha
	time.Sleep(d)
	fmt.Fprintln(w, "ok")'''
assert old in s; open(p,'w').write(s.replace(old,new))
PY
}
# start — 1-terminal: server fonda; log → api.log, chiqish kodi → kod.txt
start() {
	echo; echo "== $1"
	go build ./cmd/api || exit 1
	# dfl: SIGINT'ni standart holatga qaytarib ishga tushiradi (skript fonda ishlaganda SIGINT eʼtiborsiz meros qoladi,
	# Go esa stop()'dan keyin aynan meros holatni tiklaydi; terminalda SIGINT eʼtiborsiz emas).
	(python3 -c 'import os,signal,sys; signal.signal(signal.SIGINT, signal.SIG_DFL); os.execv(sys.argv[1], sys.argv[1:])' \
		./api -port 8783 >api.log 2>&1; echo "kod: $?" >>api.log) &
	for _ in $(seq 100); do ss -ltn | grep -q ":8783 " && break; sleep 0.1; done
}
waitapi() { while pgrep -x api >/dev/null; do sleep 0.1; done; sleep 0.2; }
log() { echo "--- 1-terminal: ./api -port 8783 ; echo \"kod: \$?\""; cat api.log; }
show() { echo "\$ $1"; bash -c "$1" 2>&1; }

# ---- 0. eski
mkdir -p "$w/eski" && cp -r "$src"/../gb-8-02-slog/solution/. "$w/eski"/ && sekin "$w/eski" && (
	cd "$w/eski" && start "0. eski (GB-8.02 kodi): SIGTERM"
	show "curl -sS 'localhost:8783/health?sekin=3s'" & sleep 1
	show "pkill -e -x api"
	wait; waitapi; log
)

# ---- 1. yangi
cp -r "$t" "$w/yangi" && sekin "$w/yangi" && (
	cd "$w/yangi"
	start "1a. yangi: sekin soʻrov tugaydi, yangisi rad etiladi"
	show "curl -sS 'localhost:8783/health?sekin=3s'" & sleep 1
	show "pkill -e -x api"; sleep 0.3
	show "curl -sS localhost:8783/health"
	wait; waitapi; log

	start "1b. yangi: muddat tugaydi (sekin=15s, shutdownTimeout 8 s)"
	SECONDS=0
	show "curl -sS 'localhost:8783/health?sekin=15s'" & sleep 1
	show "pkill -e -x api"
	wait; waitapi; log; echo "(${SECONDS} s)"

	start "1c. yangi: ikkinchi Ctrl+C (SIGINT ikki marta)"
	show "curl -sS 'localhost:8783/health?sekin=15s'" & sleep 1
	pkill -INT -x api; sleep 0.5; pkill -INT -x api
	wait; waitapi; log
)

# ---- 2. chek
cp -r "$t" "$w/chek" && (
	cd "$w/chek"
	python3 - <<'PY'
p='cmd/api/handlers.go'; s=open(p).read()
old='a.bg.Go(func() { sendReceipt(logger, o) })'
assert old in s; open(p,'w').write(s.replace(old,'go sendReceipt(logger, o)'))
PY
	start "2a. chek: go sendReceipt (egasiz)"
	show "./buyurtma.sh 8783 && pkill -x api"
	waitapi; log
)
(
	cd "$t"
	start "2b. chek: a.bg.Go (egasi server)"
	show "./buyurtma.sh 8783 && pkill -x api"
	waitapi; log
)

# ---- 3. ReadHeaderTimeout
(
	cd "$t"
	start "3. ReadHeaderTimeout"
	show "time printf 'GET /health HTTP/1.1\r\n' | nc localhost 8783"
	pkill -x api; waitapi
)

# ---- 4. vazifa
echo; echo "== 4. vazifa"
show "go test ./cmd/api -run WaitBG"
cp "$src"/solution/cmd/api/receipt.go cmd/api/receipt.go
show "go test ./cmd/api -run WaitBG"

docker compose -p gb803 down -v >/dev/null 2>&1
rm -rf "$w"
