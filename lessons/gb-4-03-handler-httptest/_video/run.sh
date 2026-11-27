#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasidagi main.go
# vaqtinchalik bozor papkasiga shu go.mod bilan koʻchiriladi va u yerda berilgan buyruq bajariladi.
#   ./_video/run.sh 01-get go run .
#   ./_video/run.sh final go run .        (final = papkadagi main.go)
# 02-sukut va 03-timeout uchun avval boshqa terminalda: nc -l 8089
set -u
export GOTOOLCHAIN=${GOTOOLCHAIN:-go1.27.4} CGO_ENABLED=0
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-get}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$t/"
if [ "$step" = final ]; then
	cp "$here/main.go" "$t/"
else
	cp "$here/_video/$step/"* "$t/"
fi
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
