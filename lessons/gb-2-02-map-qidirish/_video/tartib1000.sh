#!/usr/bin/env bash
# 04-tartib dasturini bir marta quradi, 1000 marta "olma non olma choy" bilan ishga tushiradi va
# chiqqan tartiblarni sanaydi (videodagi ustunli diagramma shu natijadan).
set -eu
here=$(cd "$(dirname "$0")/.." && pwd)
t=$(mktemp -d)
cp "$here/go.mod" "$here/_video/04-tartib/main.go" "$t/"
cd "$t"
go build -o tartib .
for i in $(seq 1000); do ./tartib olma non olma choy | awk '{printf "%s ", $1}'; echo; done | sort | uniq -c | sort -rn
