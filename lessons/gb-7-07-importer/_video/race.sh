#!/usr/bin/env bash
# `go test -race` / `go build -race` cgo (C kompilyator, gcc) talab qiladi. Bu mashinada gcc yoʻq, shuning uchun -race
# natijalari Docker'da olingan: gcc bor golang:1.26.8 image'i, lekin Go oʻzi shu mashinadagi go1.27.4 toolchain
# (~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.4.linux-amd64, faqat oʻqish uchun ulanadi) va modul keshi (faqat
# oʻqish uchun). --network host: konteyner ichidagi testlar localhost:5517'dagi Postgres'ga ulanadi.
# Sizda gcc boʻlsa, oddiy `go test -race ./...` yetadi (README).
#   ./_video/race.sh final    go test -race -count=1 ./...     (final = papkadagi kod, solution = javob)
#   ./_video/race.sh solution go test -race -count=1 ./...
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?final yoki solution}
shift
tc=${GOTOOLCHAIN_DIR:-$HOME/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.4.linux-amd64}
mc=${GOMODCACHE_DIR:-$HOME/go/pkg/mod}
t=$(mktemp -d)
if [ "$step" = solution ]; then src=$here/solution; else src=$here; fi
cp -r "$src"/cmd "$src"/internal "$src"/migrations "$src"/go.mod "$src"/go.sum "$src"/migrate.go "$t/"
[ -f "$here/katta.txt" ] && cp "$here/katta.txt" "$t/"
chmod -R a+rwX "$t"
docker run --rm --network host -u "$(id -u):$(id -g)" -v "$tc":/usr/local/go:ro -v "$mc":/gomod:ro -v "$t":/home/dev/bozor \
	-w /home/dev/bozor -e HOME=/tmp -e GOCACHE=/tmp/gocache -e GOMODCACHE=/gomod -e GOFLAGS="-buildvcs=false -mod=mod" \
	-e GOPROXY=off -e GOROOT=/usr/local/go -e GOTOOLCHAIN=local -e CGO_ENABLED=1 \
	-e DATABASE_URL="${DATABASE_URL:-postgres://bozor:bozor-lab@localhost:5517/bozor}" \
	-e PATH=/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin golang:1.26.8 "$@"
echo "(chiqish kodi: $?)"
rm -rf "$t"
