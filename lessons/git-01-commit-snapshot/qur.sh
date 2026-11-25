#!/usr/bin/env bash
# ombor/ ni noldan quradi: ikkita commit. Muallif va vaqt qat'iy (vaqt.sh). git >= 2.29.
set -e; cd "$(dirname "$0")"; rm -rf ombor; mkdir ombor; cd ombor
rf=--ref-format=files; git init -h 2>&1 | grep -q ref-format || rf=   # eski git'da bu bayroq yo'q
git init -q --template= -b main --object-format=sha1 $rf
git config user.name Lola
git config user.email lola@example.com
git config commit.gpgsign false; git config core.hooksPath /dev/null   # imzo va hook’lar hash’ni o’zgartiradi
git config core.pager cat; git config color.ui false; git config core.abbrev 7
git config core.autocrlf false; git config core.excludesFile /dev/null; git config core.attributesFile /dev/null
git config log.decorate false; git config branch.sort refname
git config merge.ff true; git config merge.stat true

cp ../fayllar/README.md .
git add README.md
. ../vaqt.sh; git commit -q -m "README"

cp ../fayllar/main.go .
git add main.go
. ../vaqt.sh; git commit -q -m "narx dasturi"
git log --oneline
