#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadamning main.go'si
# vaqtinchalik bozor papkasiga (shu go.mod bilan) koʻchiriladi, keyin buyruq oʻsha yerda bajariladi.
#   ./_video/run.sh 02-range   -> go run .
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 02-range}
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$here/_video/$step/main.go" "$t/"
cd "$t" || exit 1
go run .
