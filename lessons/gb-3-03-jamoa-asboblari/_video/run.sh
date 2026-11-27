#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasi vaqtinchalik
# bozor papkasiga koʻchiriladi va u yerda berilgan buyruq bajariladi.
#   ./_video/run.sh 00-boshlash go run ./cmd/bozor
#   ./_video/run.sh 01-tag sh -c 'go vet ./... && go run ./cmd/bozor'
#   ./_video/run.sh final go test ./...          (final = papkaning oʻzi, vazifa bilan)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 00-boshlash}
shift
base=$(mktemp -d) && t=$base/bozor && mkdir -p "$t"
trap 'rm -rf "$base"' EXIT
if [ "$step" = final ]; then
	cp -r "$here/go.mod" "$here/go.sum" "$here/cmd" "$here/internal" "$t/"
else
	cp -r "$here/_video/$step/." "$t/"
fi
cd "$t" || exit 1
"$@"
code=$?
echo "(chiqish kodi: $code)"
exit $code
