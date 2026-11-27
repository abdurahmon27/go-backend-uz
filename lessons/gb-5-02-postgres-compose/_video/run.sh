#!/usr/bin/env bash
# GB-5.02 videosidagi hamma buyruqni yangi `bozor/` papkasida ishga tushiradi va natijani chiqaradi.
# Ishlatish: ./_video/run.sh > _video/chiqishlar.txt   (Docker, docker compose va psql kerak)
# 5482-port xostda boʻsh boʻlishi kerak. 5432 band boʻlsa, "xato" qadami xuddi videodagidek xato beradi.
set -u
SRC=$(cd "$(dirname "$0")/.." && pwd)
W=$(mktemp -d)/bozor
mkdir -p "$W" && cp "$SRC/compose.yaml" "$SRC/bozor.sql" "$W"/ && cd "$W" || exit 1
plain() { printf '\n$ %s\n' "$1"; bash -c "$1" 2>&1; }
# docker compose terminalda yangilanib turadigan jadval chiqaradi; ekranda qoladigan oxirgi holatini olamiz
tty() {
  printf '\n$ %s\n' "$1"
  script -qec "$1" /dev/null | sed 's/\x1b\[[0-9;?]*[a-zA-Z]//g' | tr -d '\r' \
    | awk '/^\[\+\] /{blk=""} {blk=blk $0 "\n"} END{printf "%s", blk}' | sed 's/ \+$//'
}
tty 'docker compose up -d'
sleep 6
plain 'docker compose logs db | tail -n 1'
plain "docker compose exec db psql -U bozor -d bozor -c '\\dt'"
tty 'docker compose exec db psql -U bozor -d bozor < bozor.sql'   # -T yoʻq: xato (savol)
plain 'docker compose exec -T db psql -U bozor -d bozor < bozor.sql'
plain "docker compose exec db psql -U bozor -d bozor -c '\\dt'"
plain "docker compose exec db psql -U bozor -d bozor -c '\\d products'"
plain "docker compose exec db psql -U bozor -d bozor -c 'SELECT name FROM products'"
tty 'docker compose down'
tty 'docker compose up -d'
sleep 4
plain "docker compose exec db psql -U bozor -d bozor -c 'SELECT name FROM products'"
tty 'docker compose down -v'
tty 'docker compose up -d'
sleep 6
plain "docker compose exec db psql -U bozor -d bozor -c '\\dt'"
tty 'docker compose down -v'
# xato: xostdagi 5432 band
sed -i 's/5482:5432/5432:5432/' compose.yaml
tty 'docker compose up -d'
tty 'docker compose down -v'
sed -i 's/5432:5432/5482:5432/' compose.yaml
# xostdan URL bilan
tty 'docker compose up -d'
sleep 6
plain "psql postgres://bozor:bozor@127.0.0.1:5482/bozor -c '\\conninfo'"
tty 'docker compose down -v'
