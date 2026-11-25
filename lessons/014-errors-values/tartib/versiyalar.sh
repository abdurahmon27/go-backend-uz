#!/usr/bin/env bash
# Bitta dasturni olti xil Go versiyasida ishga tushiradi.
# GOTOOLCHAIN kerakli versiyani proxy.golang.org dan yuklab oladi (birinchi marta).
cd "$(dirname "$0")"
for v in go1.20.14 go1.21.4 go1.21.5 go1.24.13 go1.25.0 local; do
  out=$(GOTOOLCHAIN=$v go run -trimpath . 2>&1 | grep -v '^go: downloading' | head -1)
  [ "$v" = local ] && v=$(go env GOVERSION)
  printf '%-10s %s\n' "$v" "$out"
done
