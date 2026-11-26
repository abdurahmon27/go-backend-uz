#!/usr/bin/env bash
# baza.sh [usul] [necha marta] — poyga.sh bilan bir xil sinov, lekin Postgres'da.
#   usul json   — JSON bitta ustunda: Go oʻqiydi, qoʻshadi, qayta yozadi
#   usul insert — har buyurtma bitta INSERT
# Avval: cp .env.example .env (parol yozing) va docker compose up -d --wait
set -u
cd "$(dirname "$0")"
set -a; . ./.env; set +a
usul=${1:-insert}; marta=${2:-1}
CGO_ENABLED=0 go build -o bin/baza ./baza || exit 1
yoqoldi=0; eng_kam=2000; eng_kop=0
for ((k = 1; k <= marta; k++)); do
  ./bin/baza -tayyorla || exit 1
  if ((marta == 1)); then
    ./bin/baza -usul "$usul" -mijoz A & ./bin/baza -usul "$usul" -mijoz B & wait
    ./bin/baza -sana; exit
  fi
  ./bin/baza -usul "$usul" -mijoz A >/dev/null & ./bin/baza -usul "$usul" -mijoz B >/dev/null & wait
  jami=$(./bin/baza -sana | sed -E 's/json ustunida: ([0-9]+) ta · buyurtmalar jadvalida: ([0-9]+) .*/\1 \2/' | awk '{print $1 + $2}')
  ((yoqoldi += 2000 - jami))
  ((2000 - jami < eng_kam)) && eng_kam=$((2000 - jami))
  ((2000 - jami > eng_kop)) && eng_kop=$((2000 - jami))
done
echo "$marta urinish ($usul): yoʻqolgan buyurtma oʻrtacha $((yoqoldi / marta)), eng kami $eng_kam, eng koʻpi $eng_kop (2000 dan)"
