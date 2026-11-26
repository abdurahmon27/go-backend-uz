#!/usr/bin/env bash
# Darsdagi har bir terminal natijasini qayta oʻlchaydi (go1.26.8 linux/amd64, docker 29, node 24).
# Ishlatish: bash olchov.sh  (3000-port boʻsh boʻlsin; internet kerak: npm, docker, https)
# Eslatma: bu kompyuterda gcc yoʻq, shuning uchun host'dagi build'lar CGO_ENABLED=0 bilan (videodagi prepare.ts ham shunday).
set -u
cd "$(dirname "$0")"
export CGO_ENABLED=0
h() { echo; echo "\$ $*"; }
serve() { "$@" > /tmp/olchov-srv.log 2>&1 & P=$!; sleep 0.7; }
stop() { pkill -P $P 2>/dev/null; kill $P 2>/dev/null; wait $P 2>/dev/null; sleep 0.3; }

h go mod init api;           (d=$(mktemp -d); cd "$d" && go mod init api 2>&1 && cat go.mod)
h npm install express;       (d=$(mktemp -d); cd "$d" && npm init -y >/dev/null && npm install express 2>&1 | grep added)

h go run .;                  go build -o /tmp/api-run . && serve /tmp/api-run
h curl localhost:3000/users/42;            curl -s localhost:3000/users/42; echo
h curl -i -X POST localhost:3000/health;   curl -s -i -X POST localhost:3000/health | tr -d '\r' | grep -v Date
stop; cat /tmp/olchov-srv.log
h ls -A after go run;        (d=$(mktemp -d); cp go.mod main.go "$d"; cd "$d"; sed -i 's/:3000/:3999/' main.go; timeout 3 go run . >/dev/null 2>&1; ls -A)

h Express POST/HEAD;         (cd node && serve node index.js && curl -s -i -X POST localhost:3000/health | tr -d '\r' | sed -n '1p;/<pre>/p' && curl -s -I localhost:3000/health | head -1; stop)

h cd xato '&&' go build;     (cd xato && go build; echo "exit $?"; ls -A)

h go build -o api .;          go build -o api . && ls -l api
h node binary;                ls -lL "$(command -v node)"; node -p process.arch
h ./api-arm;                  GOOS=linux GOARCH=arm64 go build -o api-arm . && bash -i -c ./api-arm 2>&1 | tail -1; rm -f api-arm

h CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api .;  GOOS=linux GOARCH=amd64 go build -o api . && ls -l api
h docker images;             docker build -q -t api-go . >/dev/null && (cd node && docker build -q -t api-node . >/dev/null) && docker images --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}' --filter reference=api-go --filter reference=api-node && docker image inspect api-go --format '{{.Size}}'
h docker run api-go;         docker run -d --rm --name api-go -p 3000:3000 api-go >/dev/null && sleep 0.7 && curl -s localhost:3000/health; echo; docker exec api-go sh 2>&1; docker stop api-go >/dev/null

h tuzoq Dockerfile;          grep -E '^(FROM|RUN)' tuzoq/Dockerfile
h docker run tuzoq;          docker build -q -t tuzoq -f tuzoq/Dockerfile . >/dev/null && docker run --rm tuzoq; echo "exit $?"
h tuzoq ldd;                 docker build -q --target build -t tuzoq-build -f tuzoq/Dockerfile . >/dev/null && docker run --rm tuzoq-build ldd /api
h tuzoq-ok ldd + run;        docker build -q --target build -t tuzoq-ok-build -f tuzoq/Dockerfile.ok . >/dev/null && docker run --rm tuzoq-ok-build ldd /api
                             docker build -q -t tuzoq-ok -f tuzoq/Dockerfile.ok . >/dev/null && docker run -d --rm --name tuzoq-ok -p 3001:3000 tuzoq-ok >/dev/null && sleep 0.7 && curl -s localhost:3001/health; echo; docker stop tuzoq-ok >/dev/null

h go version -m api;          go version -m api | grep -E 'go1|GOOS|GOARCH|CGO|vcs.revision|modified'
h go run vs build: /version;  serve go run ./uyga/javob/orta && curl -s localhost:3000/version; echo; stop
                             go build -o /tmp/api-v ./uyga/javob/orta && serve /tmp/api-v && curl -s localhost:3000/version; echo; stop
h chuqur https;              (cd chuqur/https && GOOS=linux GOARCH=amd64 go build -o https . && docker build -q -t chuqur-https . >/dev/null && docker run --rm chuqur-https; docker build -q -t chuqur-https-ok -f Dockerfile.ok . >/dev/null && docker run --rm chuqur-https-ok)
h -ldflags='-s -w' '(README)'; go build -ldflags="-s -w" -o /tmp/api-sw . && ls -l /tmp/api-sw
