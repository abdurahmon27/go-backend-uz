#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasi (go.mod, cmd/bozor,
# internal/store) vaqtinchalik bozor papkasiga koʻchiriladi va u yerda berilgan buyruq bajariladi.
# Oxirida "(chiqish kodi: N)" — buyruqning chiqish kodi.
#   ./_video/run.sh 02-iz go run ./cmd/bozor
#   ./_video/run.sh final go run ./cmd/bozor          (final = papkaning oʻzi)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 02-iz}
shift
base=$(mktemp -d) && t=$base/bozor && mkdir -p "$t"
trap 'rm -rf "$base"' EXIT
if [ "$step" = final ]; then
	cp -r "$here/go.mod" "$here/cmd" "$here/internal" "$t/"
else
	cp -r "$here/_video/$step/." "$t/"
fi
cd "$t" || exit 1
"$@"
code=$?
echo "(chiqish kodi: $code)"
exit $code
