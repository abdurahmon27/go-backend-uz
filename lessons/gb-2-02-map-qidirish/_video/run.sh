#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasidagi main.go vaqtinchalik bozor
# papkasiga (shu go.mod bilan) koʻchiriladi, buyruq oʻsha yerda bajariladi.
#   ./_video/run.sh 00-hook go run . olma anor
#   ./_video/run.sh 04-tartib go run . olma non olma choy
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 04-tartib}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$here/_video/$step/"*.go "$t/"
cd "$t" || exit 1
"$@"
