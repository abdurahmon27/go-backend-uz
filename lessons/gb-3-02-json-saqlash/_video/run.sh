#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasidagi fayllar
# (main.go, kerak boʻlsa bozor.json) vaqtinchalik bozor papkasiga shu go.mod bilan koʻchiriladi
# va u yerda berilgan buyruq bajariladi.
#   ./_video/run.sh 01-marshal go run .
#   ./_video/run.sh 09-buzuq go run .
#   ./_video/run.sh final go run .        (final = papkadagi main.go va bozor.json)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-marshal}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$t/"
if [ "$step" = final ]; then
	cp "$here/main.go" "$here/bozor.json" "$t/"
else
	cp "$here/_video/$step/"* "$t/"
fi
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
