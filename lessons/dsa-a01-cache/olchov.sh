#!/usr/bin/env bash
# Butun o'lchov: 4 usul × 17 o'lcham (2^10 … 2^26), 10 marta, navbat bilan
# (bir "shovqinli" daqiqa bitta o'lchamning hamma natijasini buzmasin).
# Ishlatish: ./olchov.sh [yadro]   — yadro raqami lscpu -e dan (standart: 2, P-yadro).
# Natija: natija.txt. Bu noutbukda ~20 daqiqa.
set -euo pipefail
cd "$(dirname "$0")"
yadro=${1:-2}
pin=()
if command -v taskset > /dev/null; then pin=(taskset -c "$yadro"); fi
: > natija.txt
for i in $(seq 10); do
	"${pin[@]}" go test -run '^$' -bench 'Ketma|Aralash|Zanjir' -count 1 ./kesh >> natija.txt
done
go run ./jadval -cpu "$yadro" -tarqoq < natija.txt
