#!/usr/bin/env bash
# GB-6.05 videosidagi buyruqlar va ularning haqiqiy chiqishlari.
# Ishga tushirish (shu papkada): bash _video/chiqishlar.sh > _video/chiqishlar.txt
# Kerak: curl, jq, basenc (coreutils), 8765 porti boʻsh. Server videoda `go run ./cmd/api` bilan ishga tushadi;
# skript xuddi shu kodni `go build` qilib, binary'ni fonda ishga tushiradi va oxirida toʻxtatadi (chiqishi bir xil).
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
lab=$(pwd)
work=$(mktemp -d)/bozor
mkdir -p "$work"
cp -r cmd internal go.mod go.sum "$work"/
cd "$work"

show() { echo "\$ $*"; eval "$@" 2>&1; }
section() { echo; echo "### $*"; }
start_api() {
	if ss -ltn | grep -q ':8765 '; then echo "8765 band"; exit 1; fi
	go build -o api ./cmd/api || exit 1
	./api >>server.log 2>&1 &
	api_pid=$!
	for _ in $(seq 50); do curl -s localhost:8765/health >/dev/null && break; sleep 0.1; done
}
stop_api() { kill "$api_pid"; wait "$api_pid" 2>/dev/null; }

section "0. versiyalar"
show go version
show "curl --version | head -1"

section "1. qoʻlda yasalgan JWT"
show go run ./cmd/qolda | tee qolda.txt
p=$(grep payload: qolda.txt | awk '{print $2}')
h=$(grep header: qolda.txt | awk '{print $2}')
show "echo $h | basenc --base64url -d; echo"
show "echo $p | basenc --base64url -d; echo"
show "echo $p | wc -c"

section "2. go get (yangi modulda, toza kesh bilan emas: faqat chiqish matni)"
tmpmod=$(mktemp -d) && (cd "$tmpmod" && go mod init github.com/abdurahmon27/bozor >/dev/null 2>&1 && echo '$ go get github.com/golang-jwt/jwt/v5@v5.3.1' && go get github.com/golang-jwt/jwt/v5@v5.3.1 2>&1)
show grep jwt go.mod

section "3. kutubxona bilan tekshiruv"
show go run ./cmd/tekshir

section "4. klassik xato: qoidasiz Parse (_video/stages/token-qoidasiz.go)"
cp internal/token/token.go token.go.asl
cp "$lab/_video/stages/token-qoidasiz.go" internal/token/token.go
show go run ./cmd/tekshir
cp token.go.asl internal/token/token.go

section "5. server, login: session ham, JWT ham"
start_api
echo '$ go run ./cmd/api'
head -1 server.log
show "curl -s -c jar.txt localhost:8765/login -H \"Content-Type: application/json\" -d '{\"email\":\"dilnoza@bozor.test\",\"password\":\"tarvuz-osmon-qalam-ariq-27\"}' | tee login.json"
show 'T=$(jq -r .token login.json)'
T=$(jq -r .token login.json)
show 'grep session jar.txt | cut -f7'
show 'echo "$T" | cut -d. -f2 | basenc --base64url -d; echo'

section "6. /me ikki yoʻl bilan"
show curl -s -b jar.txt localhost:8765/me
show "curl -s -H \"Authorization: Bearer \$T\" localhost:8765/me"

section "7. logout (cookie ham, JWT ham bilan), keyin eski nusxalar"
show cp jar.txt eski.txt
show "curl -s -i -X POST -b jar.txt -c jar.txt -H \"Authorization: Bearer \$T\" localhost:8765/logout"
show curl -s -b eski.txt localhost:8765/me
show "curl -s -H \"Authorization: Bearer \$T\" localhost:8765/me"

section "8. har soʻrovdagi oʻlcham: cookie va Authorization sarlavhasi"
show "grep session eski.txt | cut -f7 | tr -d '\n' | wc -c"
show "printf %s \"\$T\" | wc -c"
show "printf %s \"Cookie: session=\$(grep session eski.txt | cut -f7)\" | wc -c"
show "printf %s \"Authorization: Bearer \$T\" | wc -c"

section "9. yaroqsiz JWT bilan /me (-i)"
show "curl -s -i -H \"Authorization: Bearer salom\" localhost:8765/me"
stop_api

section "10. server log'i"
cat server.log

section "11. testlar (vazifadan oldin)"
show "go test ./... 2>&1 | grep -v 'no test files'"

section "12. yechim bilan: logout'dan keyin JWT"
cp -r "$lab/solution/cmd" "$lab/solution/internal" .
start_api
show "curl -s -c jar.txt localhost:8765/login -H \"Content-Type: application/json\" -d '{\"email\":\"dilnoza@bozor.test\",\"password\":\"tarvuz-osmon-qalam-ariq-27\"}' > login.json"
T=$(jq -r .token login.json)
show "curl -s -o /dev/null -w '%{http_code}\n' -X POST -H \"Authorization: Bearer \$T\" localhost:8765/logout"
show "curl -s -H \"Authorization: Bearer \$T\" localhost:8765/me"
stop_api
tail -3 server.log

ss -ltn | grep -q ':8765 ' && echo "8765 hali band!" || echo "(8765 boʻsh)"
