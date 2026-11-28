#!/usr/bin/env bash
# GB-7.07 videosidagi buyruqlar va ularning haqiqiy chiqishlari.
# Ishga tushirish (shu papkada): bash _video/chiqishlar.sh > _video/chiqishlar.txt 2>&1
# Kerak: Docker, python3 (Ctrl+C uchun _video/ctrlc.py), goose (go install github.com/pressly/goose/v3/cmd/goose@v3.30.2),
# 5517 porti boʻsh. Video `bozor` nomli papkada yozilgan; skript papkani vaqtinchalik joyga nusxalaydi.
# Ctrl+C: _video/ctrlc.py buyruqni yangi jarayon guruhida ishga tushirib, 2.5 s dan keyin butun guruhga SIGINT
# yuboradi (terminal Ctrl+C'da shunday qiladi) va «^C» yozadi. psql: `docker compose exec -T` (chiqish faylga ketadi);
# terminalda `exec` aynan shu qatorlarni chiqaradi. -race natijalari: _video/race.sh (Docker, gcc uchun).
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
lab=$(pwd)
work=$(mktemp -d)/bozor
mkdir -p "$work"
cp -r cmd internal migrations go.mod go.sum migrate.go compose.yaml "$work"/
cd "$work"

show() { echo "\$ $*"; eval "$@" 2>&1; }
section() { echo; echo "### $*"; }
count() { show "docker compose exec -T db psql -c \"SELECT count(*) FROM products\""; }

section "1. baza va migratsiyalar"
docker compose down -v >/dev/null 2>&1
show docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5517/bozor
export DATABASE_URL=$GOOSE_DBSTRING
show goose up | sed -E 's#^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9:]{8} ##'

section "2. katta katalog"
show "go run ./cmd/katalog 200000 > katta.txt"
show head -4 katta.txt
show wc -l katta.txt

section "3. import va Ctrl+C (2.5 s)"
go build ./cmd/import && rm -f import # kompilyatsiya keshda: go run darhol ishga tushadi
echo '$ go run ./cmd/import katta.txt'
python3 "$lab/_video/ctrlc.py" 2.5 go run ./cmd/import katta.txt 2>&1
echo "(chiqish kodi: $?)"
count

section "4. qayta, oxirigacha"
show go run ./cmd/import katta.txt
count

section "5. buzuq qator (120000-qator: narx -5)"
sed '120000s/;[0-9]*;/;-5;/' katta.txt > buzuq.txt
show sed -n 120000p buzuq.txt
show go run ./cmd/import buzuq.txt
echo "(chiqish kodi: $?)"

section "6. testlar"
show "go test -v -run 'TestCancel|TestBadLine|TestPostgres' ./internal/importer"
show go test ./...

section "7. xato: read select'siz yuboradi (stages/importer-xato.go)"
cp internal/importer/importer.go importer.asl
cp "$lab/_video/stages/importer-xato.go" internal/importer/importer.go
show "go test -run 'TestCancel\$' ./internal/importer"
echo "--- xuddi shu kod CLI'da: Ctrl+C dan keyin 8 s kutamiz"
go build -o import-xato ./cmd/import
docker compose exec -T db psql -qc "TRUNCATE products RESTART IDENTITY CASCADE" >/dev/null 2>&1
timeout -s KILL 8 python3 "$lab/_video/ctrlc.py" 2.5 ./import-xato katta.txt 2>&1
echo "(chiqish kodi: $? — 137 = 8 s dan keyin KILL bilan oʻldirildi, hisobot chiqmadi)"
cp importer.asl internal/importer/importer.go
show "go test -run 'TestCancel\$' ./internal/importer"

section "8. worker soni va vaqt (Chuqurroq; har safar boʻsh jadval)"
go build -o imp ./cmd/import
for w in 1 4 8 12; do
	docker compose exec -T db psql -qc "TRUNCATE products RESTART IDENTITY CASCADE" >/dev/null 2>&1
	echo "\$ ./imp -workers $w katta.txt"
	./imp -workers "$w" katta.txt
done

section "9. tozalash"
cp katta.txt "$lab/katta.txt" 2>/dev/null # race.sh uchun
show docker compose down -v
rm -rf "$work"
