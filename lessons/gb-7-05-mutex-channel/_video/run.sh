#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasidagi fayllar vaqtinchalik
# bozor papkasiga shu go.mod bilan koʻchiriladi va u yerda berilgan buyruq bajariladi.
#   ./_video/run.sh 01-ikki go run .
#   ./_video/run.sh final go test -count=1 .      (final = starter, solution = javob)
# `-race` bu mashinada (gcc yoʻq, CGO_ENABLED=0) ishlamaydi: uning uchun ./_video/race.sh
set -u
export GOTOOLCHAIN=${GOTOOLCHAIN:-go1.27.4} CGO_ENABLED=0
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-ikki}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$t/"
if [ "$step" = final ]; then
	cp "$here"/*.go "$t/"
elif [ "$step" = solution ]; then
	cp "$here"/solution/*.go "$t/"
else
	cp "$here/_video/$step/"* "$t/"
fi
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
