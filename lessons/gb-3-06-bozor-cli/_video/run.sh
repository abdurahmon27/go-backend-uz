#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasi (go.mod, cmd/bozor,
# internal/store) vaqtinchalik bozor papkasiga koʻchiriladi va u yerda berilgan buyruq bajariladi.
# Oxirida "(chiqish kodi: N)" — buyruqning chiqish kodi.
#   ./_video/run.sh 01-nusxa go test ./...
#   ./_video/run.sh final go test ./...                (final = papkaning oʻzi, vazifa tekshiruvlari
#                                                       tekshiruv_test.go'siz: videoda vazifa hali yoʻq)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-nusxa}
shift
base=$(mktemp -d) && t=$base/bozor && mkdir -p "$t"
trap 'rm -rf "$base"' EXIT
if [ "$step" = final ]; then
	(cd "$here" && tar -cf - --exclude=tekshiruv_test.go go.mod cmd internal) | tar -xf - -C "$t"
else
	cp -r "$here/_video/$step/." "$t/"
fi
cd "$t" || exit 1
"$@"
code=$?
echo "(chiqish kodi: $code)"
exit $code
