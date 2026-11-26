#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasi vaqtinchalik
# bozor papkasiga koʻchiriladi va u yerda berilgan buyruq bajariladi.
#   ./_video/run.sh 01-paket go run ./cmd/bozor
#   ./_video/run.sh 04-modul sh -c 'go mod tidy && cat go.mod'
#   ./_video/run.sh final go run ./cmd/bozor       (final = papkaning oʻzi)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-paket}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
if [ "$step" = final ]; then
	cp -r "$here/go.mod" "$here/go.sum" "$here/cmd" "$here/internal" "$t/"
else
	cp -r "$here/_video/$step/." "$t/"
fi
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
