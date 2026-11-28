#!/usr/bin/env bash
# GB-6.07 videosidagi har chiqishni qayta yozadi → _video/chiqishlar.txt.
# Kerak: Go (GOTOOLCHAIN=go1.27.4), curl, Google Chrome (headless, _video/brauzer boshqaradi).
# Portlar: 8765 (bozor), 9090 (begona sayt) — boʻsh boʻlsin. Oxirida hamma jarayonlar toʻxtatiladi.
#
#   bash _video/chiqishlar.sh
#
# Bosqichlar (_video/stages): himoyasiz — darsgacha (cors va protectCSRF zanjirda yoʻq, SameSite=Lax);
# none — himoyasiz + SameSite=None; aks — CORS xatosi (Origin tekshirilmay qaytariladi). Yakuniy — cmd/api.
set -u
cd "$(dirname "$0")/.."
LAB=$PWD
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
T=$(mktemp -d)
OUT=$LAB/_video/chiqishlar.txt
API_PID=""
BEG_PID=""
cleanup() {
	[ -n "$API_PID" ] && kill "$API_PID" 2>/dev/null
	[ -n "$BEG_PID" ] && kill "$BEG_PID" 2>/dev/null
	wait 2>/dev/null
	rm -rf "$T"
}
trap cleanup EXIT

build() { # build NOM [bosqich] — lab nusxasiga bosqich fayllarini qoʻyib, cmd/api'ni quradi
	local name=$1 stage=${2:-}
	rm -rf "$T/src" && cp -r "$LAB" "$T/src"
	[ -n "$stage" ] && cp "$LAB/_video/stages/$stage/"*.go "$T/src/cmd/api/"
	(cd "$T/src" && go build -o "$T/api-$name" ./cmd/api)
}
serve() { # serve NOM — bozor'ni ishga tushiradi, log: $T/log-NOM
	"$T/api-$1" >"$T/log-$1" 2>&1 &
	API_PID=$!
	until curl -s -o /dev/null localhost:8765/health; do sleep 0.1; done
}
stop() { kill "$API_PID"; wait "$API_PID" 2>/dev/null; API_PID=""; }
log() { grep -v -e '"/health"' -e '"/favicon.ico"' -e 'POST "/login"' "$T/log-$1"; } # login/health shovqini olib tashlanadi
sec() { printf '\n## %s\n' "$*" >>"$OUT"; }
run() { printf '$ %s\n' "$*" >>"$OUT"; eval "$@" >>"$OUT" 2>&1; }
LOGIN='-H "Content-Type: application/json" -d '"'"'{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'"'"''

build final
build himoyasiz himoyasiz
build none none
build aks aks
go build -o "$T/begona" ./cmd/begona
(cd _video/brauzer && go build -o "$T/brauzer" .)
"$T/begona" >"$T/log-begona" 2>&1 &
BEG_PID=$!
until curl -s -o /dev/null localhost:9090/; do sleep 0.1; done

{
	echo "# GB-6.07 chiqishlari — $(date -u +%F) · $(go version) · $(curl --version | head -1 | cut -d' ' -f1-2) · $(google-chrome --version)"
	echo "# bash _video/chiqishlar.sh bilan yozilgan; har bosqich yangi jarayon (xotira — Seed'dan)."
} >"$OUT"

sec "1. hook: darsgacha (SameSite=Lax, CSRF tekshiruvisiz) — begona sahifa http://localhost:9090"
serve himoyasiz
run "$T/brauzer" csrf http://localhost:9090
stop
echo "--- bozor log'i:" >>"$OUT"; log himoyasiz >>"$OUT"

sec "2. Set-Cookie (yakuniy): __Host-session, Max-Age"
serve final
run curl -s -i -c jar.txt localhost:8765/login "$LOGIN"
run "cat jar.txt | tail -1"
run curl -s -b jar.txt localhost:8765/me
sec "3. HttpOnly: bozor sahifasida document.cookie (haqiqiy Chrome)"
run "$T/brauzer" httponly
stop

sec "4. darsgacha: begona sahifa boshqa SITE'da (127.0.0.1:9090), cookie SameSite=Lax"
serve himoyasiz
run "$T/brauzer" csrf http://127.0.0.1:9090
stop
echo "--- bozor log'i:" >>"$OUT"; log himoyasiz >>"$OUT"

sec "5. SameSite=None (CSRF tekshiruvisiz): 127.0.0.1:9090"
serve none
run curl -s -i localhost:8765/login "$LOGIN" '| grep Set-Cookie'
run "$T/brauzer" csrf http://127.0.0.1:9090
stop
echo "--- bozor log'i:" >>"$OUT"; log none >>"$OUT"

sec "6. yakuniy (protectCSRF): ikkala begona sahifa"
serve final
run "$T/brauzer" csrf http://localhost:9090
run "$T/brauzer" csrf http://127.0.0.1:9090
echo "--- bozor log'i:" >>"$OUT"; log final >>"$OUT"
sec "7. yakuniy: terminaldan curl (brauzer emas) va brauzer sarlavhalari bilan curl"
run curl -s -c jar.txt -o /dev/null localhost:8765/login "$LOGIN"
run curl -s -X POST -b jar.txt -H "'Sec-Fetch-Site: same-site'" -H "'Origin: http://localhost:9090'" localhost:8765/orders/41/cancel
run curl -s -X POST -b jar.txt localhost:8765/orders/41/cancel
stop
echo "--- bozor log'i:" >>"$OUT"; log final | tail -3 >>"$OUT"

sec "8. CORS: yakuniy — begona sahifa (localhost:9090) GET /me'ni oʻqishga urinadi"
serve final
run "$T/brauzer" cors http://localhost:9090
stop
echo "--- bozor log'i:" >>"$OUT"; log final | tail -2 >>"$OUT"

sec "9. CORS xatosi (aks): Origin tekshirilmay qaytariladi"
serve aks
run "$T/brauzer" cors http://localhost:9090
run curl -s -c jar.txt -o /dev/null localhost:8765/login "$LOGIN"
run curl -s -i -b jar.txt -H "'Origin: https://istalgan-sayt.example'" localhost:8765/me
stop

sec "10. CORS allow-list (yakuniy): ruxsatli va begona Origin, preflight"
serve final
run curl -s -c jar.txt -o /dev/null localhost:8765/login "$LOGIN"
run curl -s -i -b jar.txt -H "'Origin: http://localhost:3000'" localhost:8765/me
run curl -s -i -b jar.txt -H "'Origin: http://localhost:9090'" localhost:8765/me
run curl -s -i -X OPTIONS -H "'Origin: http://localhost:3000'" -H "'Access-Control-Request-Method: POST'" -H "'Access-Control-Request-Headers: content-type'" localhost:8765/login
run curl -s -i -X OPTIONS -H "'Origin: http://localhost:9090'" -H "'Access-Control-Request-Method: POST'" -H "'Access-Control-Request-Headers: content-type'" localhost:8765/login
run "$T/brauzer" cors http://localhost:9090
stop

sec "11. rate limit: POST /login, notoʻgʻri parol, 8 marta ketma-ket (yakuniy)"
serve final
BAD='-H "Content-Type: application/json" -d '"'"'{"email":"dilnoza@bozor.test","password":"taxmin-1"}'"'"''
run 'for i in $(seq 8); do curl -s -o /dev/null -w "%{http_code} " localhost:8765/login '"$BAD"'; done; echo'
run curl -s -i localhost:8765/login "$BAD"
sleep 13
run curl -s -o /dev/null -w "'%{http_code}\n'" localhost:8765/login "$BAD"
stop

sec "12. darsgacha (himoyasiz): rate limit yoʻq — 8 marta"
serve himoyasiz
run 'for i in $(seq 8); do curl -s -o /dev/null -w "%{http_code} " localhost:8765/login '"$BAD"'; done; echo'
stop
rm -f jar.txt
echo "tayyor: $OUT"
