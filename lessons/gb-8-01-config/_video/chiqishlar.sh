#!/usr/bin/env bash
# GB-8.01 videosidagi har buyruqni shu kompyuterda qayta ishga tushiradi va chiqishini chiqishlar.txt'ga yozadi.
# Ishga tushirish (papka ildizidan): bash _video/chiqishlar.sh > _video/chiqishlar.txt 2>&1
# Kerak: Docker (compose loyihasi gb801, 127.0.0.1:5581), bo'sh 8781-port. Oxirida server va Postgres to'xtatiladi.
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
src=$(cd "$(dirname "$0")/.." && pwd)
t=$(mktemp -d)/bozor
mkdir -p "$t" && cp -r "$src"/. "$t"/ && rm -rf "$t"/solution "$t"/_video "$t"/.env
cd "$t" && git init -q && git add -A && git -c user.name=lab -c user.email=lab@bozor.test commit -qm bozor
unset DATABASE_URL PORT HOST
show() { echo; echo "\$ $1"; bash -c "$1" 2>&1; }

echo "== go version"; go version
echo "== 1. build, sozlamasiz"
show 'go build ./cmd/api && ./api ; echo "kod: $?"'
show 'PORT=80a ./api ; echo "kod: $?"'
show 'DATABASE_URL=x ./api -port 70000 ; echo "kod: $?"'
show 'PORT= DATABASE_URL=x ./api ; echo "kod: $?"'
show './api -h ; echo "kod: $?"'
show 'PORT=9090 ./api -h'
show './api -h 2> /dev/null ; echo "kod: $?"'
show 'PORT= ./api ; echo "kod: $?"'

echo; echo "== 2. .env va git"
show 'cp .env.example .env'
sed -i 's/PAROL/bozor-lab/' .env   # videoda: muharrirda PAROL -> bozor-lab
show 'cat .env'
show 'cat .gitignore'
show 'set -a; . ./.env; set +a; echo $PORT'
show 'git check-ignore -v .env'
show 'git check-ignore -v .env.example ; echo "kod: $?"'
show 'git status --short'

echo; echo "== 3. Postgres"
echo '$ docker compose up -d --wait'; docker compose -p gb801 up -d --wait 2>&1 | tail -2

run_api() { # $1 = sarlavha, $2 = ko'rsatiladigan buyruq; serverni fonda ishga tushiradi, 2 s kutadi, ps/stat, keyin SIGTERM
  echo; echo "== $1"
  echo '$ set -a; . ./.env; set +a'
  echo "\$ $2"
  ( set -a; . ./.env; set +a; exec ./api -port 8781 ) > api.log 2>&1 &
  pid=$!
  sleep 2
  cat api.log
  echo; echo '$ ps -o args= -C api'; ps -o args= -C api
  echo; echo "\$ pgrep -x api  (shu kompyuterdagi api jarayonlari)"; pgrep -x api
  echo; echo '$ stat -c "%A %n" /proc/$pid/cmdline /proc/$pid/environ'
  stat -c "%A %n" /proc/$pid/cmdline /proc/$pid/environ
  echo; echo '$ curl -s localhost:8781/health'; curl -s localhost:8781/health; echo
  kill -TERM "$pid"; wait "$pid"; echo "(to'xtatildi: SIGTERM, videoda Ctrl+C) kod: $?"
}

# 3a. String() hali yoʻq (videodagi sizish)
cp cmd/api/config.go /tmp/config.go.$$
python3 - <<'PY'
p='cmd/api/config.go'; s=open(p).read()
i=s.index('// String —'); s=s[:i].rstrip()+'\n'
s=s.replace('\t"net/url"\n','')
open(p,'w').write(s)
PY
go build ./cmd/api && run_api "3a. String() yoʻq: sizish" "./api -port 8781"
# 3b. String() bilan
cp /tmp/config.go.$$ cmd/api/config.go && rm /tmp/config.go.$$
go build ./cmd/api && run_api "3b. String() bilan" "go build ./cmd/api && ./api -port 8781"

echo; echo "== 4. %#v (String chaqirilmaydi)"
cat > /tmp/gv_test.go.$$ <<'GO'
package main

import (
	"fmt"
	"testing"
)

func TestGoSyntax(t *testing.T) {
	c := config{port: 8781, databaseURL: "postgres://bozor:bozor-lab@localhost:5581/bozor"}
	fmt.Printf("%%v  %v\n%%+v %+v\n%%s  %s\n%%#v %#v\n", c, c, c, c)
}
GO
cp /tmp/gv_test.go.$$ cmd/api/gv_test.go && rm /tmp/gv_test.go.$$
go test ./cmd/api -run TestGoSyntax -v 2>&1 | grep '^%'
rm cmd/api/gv_test.go

echo; echo "== 5. Postgres to'xtatiladi"
docker compose -p gb801 down -v 2>&1 | tail -1
