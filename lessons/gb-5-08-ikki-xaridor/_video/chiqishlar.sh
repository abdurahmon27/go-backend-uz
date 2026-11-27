#!/usr/bin/env bash
# GB-5.08: videodagi har buyruqni ketma-ket ishga tushiradi va chiqishlarini chiqishlar.txt'ga yozadi.
# Ishlatish: GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0 ./chiqishlar.sh > chiqishlar.txt 2>&1
# Vaqtinchalik "bozor" papkasida ishlaydi (compose loyiha nomi = papka nomi, videodagidek) va oxirida docker compose down -v bilan tozalaydi.
set -u
HERE="$(cd "$(dirname "$0")" && pwd)"
LAB="$HERE/.."
TMP="$(mktemp -d)"
W="$TMP/bozor"
mkdir -p "$W"
cd "$W"
grep -v '^name: gb508' "$LAB/compose.yaml" > compose.yaml
cp "$LAB/bozor.sql" "$LAB/go.mod" "$LAB/go.sum" .
step() { echo; echo "### $1"; shift; echo "\$ $*"; "$@" 2>&1; echo "[exit $?]"; }
PSQL=(docker compose exec db psql)

# stage: papka tarkibini bosqichga moslaydi (store.go + test; promo-kod fayllarisiz)
stage() { # $1 = store.go manbasi
  rm -f ./*.go
  cp "$1" store.go
  sed 's/.*{"BuyLocked", .*//' "$LAB/store_test.go" > store_test.go   # satr raqamlari papkadagi faylniki bilan bir xil qoladi
}

step up docker compose up -d --wait
step bozor "${PSQL[@]}" -f /lab/bozor.sql
step qoldiq1 "${PSQL[@]}" -c "UPDATE products SET stock = 1 WHERE id = 4"
step qoldiq2 "${PSQL[@]}" -c "SELECT name, stock FROM products WHERE id = 4"

echo; echo "### seans-naive"; python3 "$HERE/sessions.py" naive 2>&1

# 1. birinchi versiya: oʻqi, tekshir, yoz
stage "$HERE/stages/01-store.go"
step test-naive go test
# 2. shartli UPDATE
echo; echo "### seans-shartli"; python3 "$HERE/sessions.py" shartli 2>&1
stage "$HERE/stages/02-store.go"
step test-shartli go test
# 3. FOR UPDATE (butun fayl)
echo; echo "### seans-forupdate"; python3 "$HERE/sessions.py" forupdate 2>&1
stage "$LAB/store.go"; cp "$LAB/store_test.go" store_test.go
step test-locked go test
# 4. xato: FOR UPDATE pool orqali
cp "$HERE/stages/04-store-pool-mistake.go" store.go
step test-xato go test
# 5. papkaning oʻzi (boshlangʻich holat) va yechim
rm -f ./*.go; cp "$LAB"/*.go .
step starter go test
mkdir -p sol && cp "$LAB"/solution/*.go sol/ && cp "$LAB/go.mod" "$LAB/go.sum" sol/ && cp "$LAB/bozor.sql" sol/
(cd sol && step solution go test)

step down docker compose down -v
cd /; rm -rf "$TMP"
