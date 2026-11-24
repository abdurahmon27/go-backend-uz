#!/usr/bin/env bash
# bash homework/javoblar.sh   (lessons/sys-02-process-env papkasidan)
exec 2>&1
set -m   # job control: fon jarayonlari SIGINT ni eʼtiborsiz qoldirmasin (interaktiv terminaldagi kabi)
echo "== Oson: sleep 100 ni Ctrl+C (SIGINT, 2) bilan toʻxtatish"
sleep 100 & kill -INT $!; wait $! 2>/dev/null; echo "chiqish kodi: $? (128 + 2)"

echo "== Oʻrta: PORT=9191"
CGO_ENABLED=0 go build -o srv ./server
PORT=9191 ./srv & sleep 0.3
curl -sS localhost:9191
kill $!; wait

echo "== Qiyin: os.LookupEnv"
go run ./homework/qiyin
ISM= go run ./homework/qiyin
ISM=Aziz go run ./homework/qiyin

echo "== Guru: ikkinchi SIGINT darhol chiqaradi"
CGO_ENABLED=0 go build -o srv-guru ./homework/guru
./srv-guru & g=$!; sleep 0.3
curl -sS localhost:8077 & sleep 0.3
kill -INT $g; sleep 0.3
kill -INT $g; wait $g 2>/dev/null; kod=$?
wait                                   # curl ham tugasin, tartib barqaror boʻlsin
echo "server chiqish kodi: $kod (128 + 2: Shutdown kutilmadi)"
