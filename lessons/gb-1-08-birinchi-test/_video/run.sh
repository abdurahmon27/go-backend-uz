#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasidagi main.go va
# price_test.go vaqtinchalik bozor papkasiga (shu go.mod bilan) koʻchiriladi, buyruq oʻsha yerda bajariladi.
#   ./_video/run.sh 00-gb107 go run . add olma 100000000000000000 1
#   ./_video/run.sh 03-jadval go test
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 03-jadval}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$here/_video/$step/"*.go "$t/"
cd "$t" || exit 1
"$@"
