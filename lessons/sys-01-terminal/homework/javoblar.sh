#!/usr/bin/env bash
# SYS-01 · Uy vazifasi javoblari. Ishga tushirish (sys-01-terminal papkasidan): bash homework/javoblar.sh
set -u

echo "== Oson: nechta .go fayl?"
find loyiha -name '*.go' | wc -l

echo "== Oʻrta: qaysi boʻlimda (cmd, docs, internal) eng koʻp TODO?"
# loyiha/internal/savat/savat.go:TODO  ->  2-maydon (/ boʻyicha) = internal
grep -ro TODO loyiha | cut -d/ -f2 | sort | uniq -c | sort -rn

echo "== Qiyin: top-3 soʻz, oʻ buzilmasin"
# 1) o/g dan keyingi ' ‘ ’ ni ʻ ga keltiramiz (GO-1.04 dagi tozala() ning sed varianti)
# 2) UTF-8 locale'da [[:alpha:]] ʻ (U+02BB, Lm) ni ham harf deb biladi
sed "s/\([oOgG]\)['‘’]/\1ʻ/g" matn.txt | grep -oE "[[:alpha:]]+" | tr A-Z a-z | sort | uniq -c | sort -rn | head -3

echo "== Guru: > f 2>&1  va  2>&1 > f"
grep -c TODO loyiha/cmd/api/main.go yoq > f1.txt 2>&1
echo "-- A (> f 2>&1) fayl ichida:"; cat f1.txt
echo "-- B (2>&1 > f) ekranga:"
grep -c TODO loyiha/cmd/api/main.go yoq 2>&1 > f2.txt
echo "-- B fayl ichida:"; cat f2.txt
rm -f f1.txt f2.txt
