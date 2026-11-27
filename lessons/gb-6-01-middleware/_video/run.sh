#!/usr/bin/env bash
# Videodagi hamma chiqish shu skriptdan: ./run.sh > chiqishlar.txt 2>&1
# Har variant: papkaning (starter) nusxasi, cmd/api/middleware.go variant bilan almashtiriladi,
# kerak boʻlsa getProduct'ga vaqtincha panic qatori qoʻshiladi; server :8765 da, curl, server toʻxtatiladi.
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
here=$(cd "$(dirname "$0")" && pwd)
lab=$(dirname "$here")
go version; curl --version | head -1

variant() { # $1 = variant papkasi ("" = starter), $2 = panic (1/0)
  dir=$(mktemp -d)/bozor
  mkdir -p "$dir"
  (cd "$lab" && tar --exclude=./_video --exclude=./solution -cf - .) | (cd "$dir" && tar xf -)
  [ -n "$1" ] && cp "$here/$1/middleware.go" "$dir/cmd/api/middleware.go"
  if [ "$2" = 1 ]; then
    python3 - "$dir/cmd/api/handlers.go" <<'PY'
import sys
p=sys.argv[1]; s=open(p).read()
a="""	id, ok := pathID(w, r)
	if !ok {
		return
	}
	p, err := a.st.Get(id)"""
assert a in s
s=s.replace(a, a.replace("	p, err := a.st.Get(id)", """	if id == 13 { // vaqtincha: bug oʻrnida
		panic("sinov: 13-mahsulotda bug")
	}
	p, err := a.st.Get(id)"""),1)
open(p,'w').write(s)
PY
  fi
  (cd "$dir" && go build -o api ./cmd/api) || exit 1
  (cd "$dir" && exec ./api > server.log 2>&1) &
  pid=$!
  for _ in $(seq 50); do (exec 3<>/dev/tcp/127.0.0.1/8765) 2>/dev/null && break; sleep 0.1; done
}
stop() { kill "$pid" 2>/dev/null; wait "$pid" 2>/dev/null; echo "--- server log:"; sed "s#$dir#/…/bozor#g" "$dir/server.log"; echo; }
c() { echo "\$ curl $*"; curl -sS "$@"; echo " [exit $?]"; }  # -sS: progress meter yoʻq (terminalda u baribir chiqmaydi)

# ONLY=v7 ./run.sh — faqat v7 (yoʻl %s bilan) qayta yoziladi.
if [ "${ONLY:-}" = v7 ]; then
  echo "=== v7: yakuniy middleware.go, yoʻl %s bilan (log injection)"
  variant v7-s 0
  c -s 'localhost:8765/products/1%0A2026/11/27%2021:00:00%20soxta-qator'
  stop
  exit 0
fi

echo "=== v0: middleware yoʻq, 13-mahsulotda panic"
variant v0-yoq 1
c -i localhost:8765/products/13
c -s localhost:8765/products/1
stop

echo "=== v1: logRequests (status yoʻq)"
variant v1-log 0
c -s localhost:8765/products/1
c -s localhost:8765/products/99
stop

echo "=== v2: statusRecorder, status: http.StatusOK"
variant v2-status 0
c -s localhost:8765/products/99
c -s localhost:8765/health
c -s localhost:8765/narsalar
stop

echo "=== v3: status boshlangʻich qiymatsiz"
variant v3-nol 0
c -s localhost:8765/health
c -s localhost:8765/products/99
stop

echo "=== v4: notoʻgʻri tartib recoverPanic(logRequests(mux)), panic"
variant v4-tartib 1
c -i localhost:8765/products/13
stop

echo "=== v5: yakuniy zanjir (starter), panic"
variant "" 1
c -i localhost:8765/products/13
stop

echo "=== v6: yakuniy zanjir (starter), panic qatori olib tashlangan"
variant "" 0
c -i localhost:8765/products/1
c -s localhost:8765/products/99
c -s 'localhost:8765/products/1%0A2026/11/27%2021:00:00%20soxta-qator'
stop

echo "=== v7: yakuniy middleware.go, yoʻl %s bilan (log injection)"
variant v7-s 0
c -s 'localhost:8765/products/1%0A2026/11/27%2021:00:00%20soxta-qator'
stop

echo "=== testlar"
cd "$lab" && gofmt -l . ; go vet ./... && go test ./... ; echo "starter exit $?"
(cd "$lab/solution" && go test ./...)
