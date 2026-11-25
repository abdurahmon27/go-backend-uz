#!/usr/bin/env bash
# ./yigish.sh — videodagi hamma image'larni yigʻadi (birinchi marta golang:1.26.8 yuklanadi, ~300 MB).
set -e
cd "$(dirname "$0")"
docker build -q -t mehmon:bir    -f Dockerfile.bir .
docker build -q -t mehmon:kichik -f Dockerfile .
docker build -q -t mehmon-boshqa:cgo        -f boshqa/Dockerfile.cgo .
docker build -q -t mehmon-boshqa:cgo-build  -f boshqa/Dockerfile.cgo --target build .
docker build -q -t mehmon-boshqa:alpine     -f boshqa/Dockerfile.alpine .
docker build -q -t mehmon-boshqa:distroless -f boshqa/Dockerfile.distroless .
docker build -q -t mehmon-tuzoq             -f tuzoq/Dockerfile .
docker image ls --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}' --filter 'reference=mehmon*'
