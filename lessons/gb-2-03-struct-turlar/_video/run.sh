#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadamning main.go'si vaqtinchalik
# bozor papkasiga (shu go.mod bilan) koʻchiriladi va u yerda berilgan buyruq bajariladi.
#   ./_video/run.sh 01-product go run .
#   ./_video/run.sh final go run .        (final = papkadagi main.go)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-product}
shift
src="$here/_video/$step/main.go"
[ "$step" = final ] && src="$here/main.go"
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$src" "$t/"
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
