#!/usr/bin/env bash
# Ikkita mashq repo'sini noldan quradi. Muallif va vaqt qat'iy (vaqt.sh). git >= 2.29.
#   loyiha/: .env commit qilingan, keyingi commit'da "o'chirilgan"
#   tuzoq/:  .env birinchi commit'dayoq kuzatila boshlagan
# Kalit soxta: soxta_... hech qayerda ishlamaydi.
set -e; cd "$(dirname "$0")"
rf=--ref-format=files; git init -h 2>&1 | grep -q ref-format || rf=   # eski git'da bu bayroq yo'q
yangi() { rm -rf "$1"; mkdir "$1"; cd "$1"; git init -q --template= -b main --object-format=sha1 $rf
  git config user.name Lola; git config user.email lola@example.com; git config commit.gpgsign false
  git config core.hooksPath /dev/null
git config core.pager cat; git config color.ui false; git config core.abbrev 7
git config core.autocrlf false; git config core.excludesFile /dev/null; git config core.attributesFile /dev/null
git config log.decorate false; git config branch.sort refname
git config merge.ff true; git config merge.stat true; }

yangi loyiha
cp ../fayllar/main.go .
git add main.go; . ../vaqt.sh; git commit -q -m "to'lov: kalit muhitdan"
cp ../fayllar/env.txt .env
git add .; . ../vaqt.sh; git commit -q -m "lokal sozlamalar"
git rm -q .env; . ../vaqt.sh; git commit -q -m "kalitni o'chirdim"
cd ..

yangi tuzoq
cp ../fayllar/main.go .; cp ../fayllar/env.txt .env
git add .; . ../vaqt.sh; git commit -q -m "birinchi commit"
cd ..
echo "tayyor: loyiha/ va tuzoq/"
