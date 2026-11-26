#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadamning main.go'si
# vaqtinchalik bozor papkasiga (shu go.mod bilan) koʻchiriladi, keyin go run . oʻsha yerda bajariladi.
#   ./_video/run.sh 01-manfiy      -> go run .
# 06-teng-xato'da kod main.go.txt nomi bilan turadi: unda ataylab sintaksis xatosi bor, gofmt -l . va
# go vet ./... uni oʻqimasin. run.sh uni main.go qilib koʻchiradi.
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-manfiy}
t=$(mktemp -d)/bozor && mkdir -p "$t"
src="$here/_video/$step/main.go"
[ -f "$src" ] || src="$src.txt"
cp "$here/go.mod" "$t/" && cp "$src" "$t/main.go"
cd "$t" || exit 1
go run .
