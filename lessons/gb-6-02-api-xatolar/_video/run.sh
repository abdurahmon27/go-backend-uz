#!/usr/bin/env bash
# _video/run.sh STEP  — videodagi qadamni vaqtinchalik «bozor/» papkasida tayyorlaydi va yoʻlini chiqaradi.
#   oldin   = GB-6.01'ning solution'i (darsdan oldingi holat: {"error": …} shakli)
#   final   = shu papka (videoning oxiri = starter)
#   unutgan = final + getProduct'da writeError'dan keyin return yoʻq (keng tarqalgan xato)
#   solution
# Hammasi: GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0.
set -euo pipefail
here=$(cd "$(dirname "$0")/.." && pwd)
step=$1
t=$(mktemp -d)/bozor
case $step in
  oldin)    cp -r "$here/../gb-6-01-middleware/solution" "$t" ;;
  final)    mkdir -p "$t" && cp -r "$here"/{go.mod,bozor.json,cmd,internal} "$t"/ ;;
  unutgan)  mkdir -p "$t" && cp -r "$here"/{go.mod,bozor.json,cmd,internal} "$t"/ && cp "$here/_video/unutgan/handlers.go" "$t/cmd/api/handlers.go" ;;
  solution) cp -r "$here/solution" "$t" ;;
esac
echo "$t"
