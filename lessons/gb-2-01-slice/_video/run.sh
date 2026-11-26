#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: qadam papkasidagi main.go vaqtinchalik
# bozor papkasiga (shu go.mod bilan) koʻchiriladi, buyruq oʻsha yerda bajariladi.
#   ./_video/run.sh 00-bug go run .
#   ./_video/run.sh 05-xato go build
# Videodagi Go versiyasi: GOTOOLCHAIN=go1.27.4, CGO_ENABLED=0 (skript ikkalasini oʻzi qoʻyadi).
set -u
export GOTOOLCHAIN=${GOTOOLCHAIN:-go1.27.4} CGO_ENABLED=0
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 03-append}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$here/_video/$step/"*.go "$t/"
cd "$t" || exit 1
"$@"
