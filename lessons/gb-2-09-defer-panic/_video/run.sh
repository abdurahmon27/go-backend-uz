#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadamning main.go'si vaqtinchalik
# bozor papkasiga (shu go.mod bilan) koʻchiriladi va u yerda berilgan buyruq bajariladi.
# Oxirida "(chiqish kodi: N)" — buyruqning chiqish kodi; vaqtinchalik papka oʻchiriladi.
#   ./_video/run.sh 02-defer bash -c 'ulimit -n 1024; go run .'
#   ./_video/run.sh 06-panic go run .
#   ./_video/run.sh final go run .        (final = papkadagi main.go)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 02-defer}
shift
src="$here/_video/$step/main.go"
[ "$step" = final ] && src="$here/main.go"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
t=$tmp/bozor && mkdir -p "$t"
cp "$here/go.mod" "$src" "$t/"
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
