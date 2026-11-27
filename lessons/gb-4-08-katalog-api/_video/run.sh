#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: papka (starter) vaqtinchalik bozor/
# papkasiga koʻchiriladi, ustiga qadam papkasidagi fayllar yoziladi va u yerda buyruq bajariladi.
#   ./_video/run.sh final go test ./...                 (final = starter, tekshiruv_test.go'siz; starter = papkaning oʻzi)
#   ./_video/run.sh 01-pointer go test ./...            (qadam = starter + _video/01-pointer/)
#   ./_video/run.sh solution go test ./...              (javob)
# `-race` bu mashinada (gcc yoʻq, CGO_ENABLED=0) ishlamaydi: uning uchun ./_video/race.sh
set -u
export GOTOOLCHAIN=${GOTOOLCHAIN:-go1.27.4} CGO_ENABLED=0
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-pointer}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
if [ "$step" = solution ]; then
	cp -r "$here/solution/." "$t/"
else
	(cd "$here" && cp -r go.mod bozor.json cmd internal "$t/")
	case $step in
	starter) ;;                                                        # papka oʻzi, tekshiruv bilan
	final) find "$t" -name tekshiruv_test.go -delete ;;                # videoning oxiri: vazifa tekshiruvisiz
	*) find "$t" -name tekshiruv_test.go -delete; cp -r "$here/_video/$step/." "$t/" ;;
	esac
fi
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
