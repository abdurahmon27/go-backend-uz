#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadamning main.go'si vaqtinchalik
# bozor papkasiga (shu go.mod bilan) koʻchiriladi, go build qilinadi, keyin ./bozor berilgan
# argumentlar bilan ishga tushadi va chiqish kodi yoziladi.
#   ./_video/run.sh 03-atoi add olma 12k
#   ./_video/run.sh final add olma 12000      (final = papkadagi main.go)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 03-atoi}
shift
src="$here/_video/$step/main.go"
[ "$step" = final ] && src="$here/main.go"
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$src" "$t/"
cd "$t" || exit 1
go build || exit $?
./bozor "$@"
echo "(chiqish kodi: $?)"
