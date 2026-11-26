#!/usr/bin/env bash
# eski/ (go 1.21) ning nusxasiga go 1.22 talab qiladigan kutubxonani qo'shib, uni ishlatamiz.
# Router kodiga tegmaymiz. go get go.mod'dagi go qatorini o'zi ko'taradi va ServeMux xulqi o'zgaradi.
# Asl eski/ papkasi o'zgarmaydi: hammasi vaqtinchalik nusxada.
set -eu
proxy="file://$(cd "$(dirname "$0")" && pwd)/proxy"
nusxa=$(mktemp -d)
trap 'rm -rf "$nusxa"' EXIT
cp "$(dirname "$0")"/../eski/{go.mod,main.go} "$nusxa"
cd "$nusxa"
go run .
# yangi fayl: kutubxonani ishlatadi, router'ga tegmaydi
printf 'package main\n\nimport "example.com/ramka"\n\nvar _ = ramka.Nomi()\n' > ramka.go
GOPROXY="$proxy" GONOSUMDB=example.com go get example.com/ramka@v1.0.0 2>&1 | grep -v downloading
go run .
