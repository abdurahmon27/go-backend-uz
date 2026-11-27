#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadamning main.go'si, go.mod va
# katalog fayllari (katalog.txt, katta.txt) vaqtinchalik bozor papkasiga koʻchiriladi va u yerda
# berilgan buyruq bajariladi. Oxirida "(chiqish kodi: N)" — buyruqning chiqish kodi.
#   ./_video/run.sh 02-read go run .
#   ./_video/run.sh final sh -c 'cat katalog.txt | go run .'
#   ./_video/run.sh final go run . katalog.txt        (final = papkadagi main.go)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 02-read}
shift
src="$here/_video/$step/main.go"
[ "$step" = final ] && src="$here/main.go"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
t=$tmp/bozor && mkdir -p "$t"
cp "$here/go.mod" "$here/katalog.txt" "$here/katta.txt" "$t/"
[ "$step" = 00-cat ] || cp "$src" "$t/"
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
