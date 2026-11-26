#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadamning main.go'si
# vaqtinchalik bozor papkasiga (shu go.mod bilan) koʻchiriladi, keyin buyruq oʻsha yerda bajariladi.
#   ./_video/run.sh 03-xato              -> go run .
#   ./_video/run.sh 11-string-xato vet   -> go vet
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 03-xato}
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/bozor"
cp "$here/go.mod" "$here/_video/$step/main.go" "$tmp/bozor/"
cd "$tmp/bozor" || exit 1
case "${2:-run}" in
  vet) go vet ;;
  *) go run . ;;
esac
