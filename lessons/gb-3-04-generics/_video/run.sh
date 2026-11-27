#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadamning main.go'si va go.mod
# vaqtinchalik bozor papkasiga koʻchiriladi va u yerda berilgan buyruq bajariladi.
# Oxirida "(chiqish kodi: N)" — buyruqning chiqish kodi; skript ham shu kod bilan tugaydi.
#   ./_video/run.sh 02-filter go run .
#   ./_video/run.sh 05-sum-any go build
#   ./_video/run.sh final go run .        (final = papkadagi main.go)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 02-filter}
shift
src="$here/_video/$step/main.go"
[ "$step" = final ] && src="$here/main.go"
[ -f "$src" ] || { echo "qadam yoʻq: $step (ls _video)" >&2; exit 2; }
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
t=$tmp/bozor && mkdir -p "$t"
cp "$here/go.mod" "$src" "$t/"
cd "$t" || exit 1
"$@"
rc=$?
echo "(chiqish kodi: $rc)"
exit "$rc"
