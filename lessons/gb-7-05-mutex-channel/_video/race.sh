#!/usr/bin/env bash
# `go test -race` cgo (C kompilyator, gcc) talab qiladi. Bu mashinada gcc yoʻq, shuning uchun videodagi -race
# natijalari Docker'da olingan: gcc bor golang:1.26.8 image'i, lekin Go oʻzi shu mashinadagi go1.27.4 toolchain
# (~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.4.linux-amd64, faqat oʻqish uchun ulanadi).
# Sizning mashinangizda gcc oʻrnatilgan boʻlsa, oddiy `go test -race` yetadi (README).
#   ./_video/race.sh 01-ikki go run -race .
#   ./_video/race.sh final go test -race -count=1 .      (final = starter, solution = javob)
set -u
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 01-ikki}
shift
tc=${GOTOOLCHAIN_DIR:-$HOME/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.27.4.linux-amd64}
t=$(mktemp -d)
cp "$here/go.mod" "$t/"
if [ "$step" = final ]; then cp "$here"/*.go "$t/"; elif [ "$step" = solution ]; then cp "$here"/solution/*.go "$t/"; else cp "$here/_video/$step/"* "$t/"; fi
chmod -R a+rwX "$t"
docker run --rm -u "$(id -u):$(id -g)" -v "$tc":/usr/local/go:ro -v "$t":/home/dev/bozor -w /home/dev/bozor \
	-e HOME=/tmp -e GOCACHE=/tmp/gocache -e GOROOT=/usr/local/go -e GOTOOLCHAIN=local -e GOFLAGS=-buildvcs=false \
	-e CGO_ENABLED=1 -e PATH=/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin golang:1.26.8 "$@"
echo "(chiqish kodi: $?)"
