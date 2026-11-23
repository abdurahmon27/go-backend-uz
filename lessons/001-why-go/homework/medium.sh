#!/usr/bin/env bash
# O'rta: bitta Linux mashinasida uch xil OS uchun build. GOOS/GOARCH = maqsad tizim, compiler o'zi shu yerda.
set -e
cd "$(dirname "$0")/.."
for t in linux/amd64 windows/amd64 darwin/arm64; do
  os=${t%/*}; arch=${t#*/}; ext=; [ "$os" = windows ] && ext=.exe
  GOOS=$os GOARCH=$arch go build -o "/tmp/hello-$os-$arch$ext" ./hello
done
ls -lh /tmp/hello-* | awk '{print $5, $9}'
