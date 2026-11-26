#!/usr/bin/env bash
# Butun o'lchov: 5 marta navbat bilan (xotira benchmark'lari + disk, fsync, tarmoq, internet),
# bir "shovqinli" daqiqa bitta qatorning hamma natijasini buzmasin.
# Ishlatish: ./olchov.sh [yadro]   — benchmark'lar uchun yadro (standart: 2, P-yadro).
# Natija: natija.txt. Bu noutbukda ~4 daqiqa. Internet qismi go.dev ga ~210 ta soʻrov yuboradi (har turda 31 ulanish + 11 HTTPS).
set -euo pipefail
cd "$(dirname "$0")"
yadro=${1:-2}
pin=()
if command -v taskset > /dev/null; then pin=(taskset -c "$yadro"); fi
export CGO_ENABLED=0
go build -o olcha.bin ./olcha
: > natija.txt
for i in $(seq 5); do
	echo "# tur $i $(date +%H:%M:%S) load $(cut -d' ' -f1 /proc/loadavg)" >> natija.txt
	"${pin[@]}" go test -run '^$' -bench . -count 2 ./xotira | grep ^Benchmark >> natija.txt
	./olcha.bin hammasi >> natija.txt
done
rm -f olcha.bin
go run ./jadval -tarqoq < natija.txt
