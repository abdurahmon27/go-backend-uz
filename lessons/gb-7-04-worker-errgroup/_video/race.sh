#!/usr/bin/env bash
# `go test -race` / `go run -race` cgo (C kompilyator, gcc) talab qiladi. Bu mashinada gcc yoʻq, shuning uchun -race
# natijalari Docker'da olingan: gcc bor golang:1.26.8 image'i, lekin Go oʻzi shu mashinadagi go1.27.4 toolchain
# (~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.4.linux-amd64, faqat oʻqish uchun ulanadi) va modul keshi
# (golang.org/x/sync uchun, faqat oʻqish uchun). Sizda gcc boʻlsa, oddiy `go test -race` yetadi (README).
#   ./_video/race.sh 04-worker go run -race .
#   ./_video/race.sh solution go test -race -count=1 .      (final = starter, solution = javob)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 04-worker}
shift
tc=${GOTOOLCHAIN_DIR:-$HOME/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.4.linux-amd64}
mc=${GOMODCACHE_DIR:-$HOME/go/pkg/mod}
t=$(mktemp -d)
cp "$here/go.mod" "$here/go.sum" "$t/"
if [ "$step" = final ]; then cp "$here"/*.go "$t/"; elif [ "$step" = solution ]; then cp "$here"/solution/*.go "$t/"; else cp "$here/supplier.go" "$here/_video/$step/"* "$t/"; fi
chmod -R a+rwX "$t"
docker run --rm -u "$(id -u):$(id -g)" -v "$tc":/usr/local/go:ro -v "$mc":/gomod:ro -v "$t":/home/dev/bozor -w /home/dev/bozor \
	-e HOME=/tmp -e GOCACHE=/tmp/gocache -e GOMODCACHE=/gomod -e GOFLAGS="-buildvcs=false -mod=mod" -e GOPROXY=off \
	-e GOROOT=/usr/local/go -e GOTOOLCHAIN=local -e CGO_ENABLED=1 -e PATH=/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin \
	golang:1.26.8 "$@"
echo "(chiqish kodi: $?)"
