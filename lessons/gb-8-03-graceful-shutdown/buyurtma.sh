#!/usr/bin/env bash
# buyurtma.sh — Jasur nomidan login qilib, bitta buyurtma yuboradi (GB-6.04: login va session cookie, GB-5.10: POST /orders).
# GB-8.03 videosida chekni koʻrsatish uchun: buyurtma yozilgach server chekni fonda yuboradi (~2 s).
# Ishlatish: ./buyurtma.sh [port]   (standart 8765)
set -eu
port=${1:-8765}
jar=$(mktemp)
trap 'rm -f "$jar"' EXIT
curl -sf -c "$jar" -o /dev/null "localhost:$port/login" \
	-H 'Content-Type: application/json' \
	-d '{"email":"jasur@bozor.test","password":"choy-non-qand-27"}'
curl -s -b "$jar" -o /dev/null -w 'POST /orders → %{http_code}\n' "localhost:$port/orders" \
	-H 'Content-Type: application/json' \
	-d '{"items":[{"product_id":1,"quantity":1}]}'
