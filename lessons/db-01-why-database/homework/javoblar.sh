#!/usr/bin/env bash
# Uy vazifasi javoblarini ishga tushiradi. Avval: .env va docker compose up -d --wait (README).
set -u
cd "$(dirname "$0")/.."
set -a; . ./.env; set +a
export CGO_ENABLED=0
go build -o bin/fayl ./fayl && go build -o bin/baza ./baza && go build -o bin/tuzat ./homework/baza-tuzat && go build -o bin/fsync ./homework/fsync || exit 1

echo "== Oson: qulf usulini 100 marta kill -9"
./qulash.sh qulf 100

for u in sql forupdate; do
  echo "== $( [[ $u == sql ]] && echo Oʻrta || echo Qiyin ): Postgres, usul $u"
  ./bin/baza -tayyorla
  ./bin/tuzat -usul $u -mijoz A & ./bin/tuzat -usul $u -mijoz B & wait
  ./bin/baza -sana
done

echo "== Guru: qulf (fsync'siz) va qulf + fsync, ikki jarayon × 1000"
for p in fayl fsync; do
  rm -f buyurtmalar.json buyurtmalar.json.*
  bosh=$(date +%s%N)
  ./bin/$p -usul qulf -mijoz A >/dev/null & ./bin/$p -usul qulf -mijoz B >/dev/null & wait
  ms=$((($(date +%s%N) - bosh) / 1000000))
  echo "$p: $(./bin/fayl -sana), $ms ms"
done
rm -f buyurtmalar.json buyurtmalar.json.*
