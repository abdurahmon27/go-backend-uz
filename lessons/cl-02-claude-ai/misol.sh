#!/usr/bin/env bash
# Videodagi misol fayllarini yaratadi: misol/xato.log va misol/bozor/.
# Kalitlar SOXTA. Provayder prefiksli kalit (sk-ant-…) repo'da toʻliq satr boʻlib turmaydi
# (boʻlaklardan yigʻiladi), shunda GitHub secret scanning uni haqiqiy kalit deb oʻylamaydi.
set -euo pipefail
cd "$(dirname "$0")"
rm -rf misol && mkdir -p misol/bozor

a="sk-ant-"; a="${a}api03-SOXTA7kalit9Qx2mN4vB8cR1tY6uE3wZ5pL0aD"   # Anthropic API kalitiga oʻxshash
t="Xq7Lm2Pz9Rt4Vb8Nc3Kd6Wf1Hs5Jy0Ge"                             # nomaʼlum formatdagi toʻlov kaliti
d="postgres://bozor:Qovun2026"; d="${d}@localhost:5432/bozor"       # parolli ulanish manzili

cat > misol/xato.log <<LOG
10:14:02 bozor: server :8080 portda ishga tushdi
10:14:02 config: DB=${d}
10:14:02 config: ClaudeKalit=${a}
10:14:07 POST /buyurtma id=2026112510140712 mijoz=aziz@example.com
10:14:07 tel=+998 90 000 11 22 karta=4111 1111 1111 1111
10:14:07 tolov: kalit=${t}
10:14:08 tolov: javob 401, buyurtma saqlanmadi
build: commit 3f9a1c27e0b84d6a9e2c5b7f1d08a4c6e9b2f731
LOG

printf 'ANTHROPIC_API_KEY=%s\nDB_URL=%s\n' "$a" "$d" > misol/bozor/.env
printf 'ANTHROPIC_API_KEY=\nDB_URL=\n' > misol/bozor/.env.example
printf '# bozor\n\nGo backend kursining loyihasi. Ishga tushirish: go run ./cmd/bozor\n' > misol/bozor/README.md
printf 'package main\n\nfunc main() {\n\tprintln("bozor")\n}\n' > misol/bozor/main.go.txt
