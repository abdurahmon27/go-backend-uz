#!/usr/bin/env bash
# Qiyin: dokon nusxasida example.com/log ni v1.2.0 ga tushiramiz.
# Asl mvs/dokon o'zgarmaydi: hammasi vaqtinchalik nusxada.
set -u
bu=$(cd "$(dirname "$0")" && pwd)
export GOPROXY="file://$bu/../../mvs/proxy" GONOSUMDB=example.com
nusxa=$(mktemp -d)
trap 'rm -rf "$nusxa"' EXIT
cp "$bu"/../../mvs/dokon/{go.mod,go.sum,main.go} "$nusxa"
cd "$nusxa"
go get example.com/log@v1.2.0 2>&1 | grep -v downloading
echo "--- go.mod:"
cat go.mod
echo "--- go build:"
go build -o /dev/null .
