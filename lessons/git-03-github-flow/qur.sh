#!/usr/bin/env bash
# GitHub flow mashqi uchun uchta repo'ni noldan quradi. git >= 2.29.
#   server.git/  GitHub o'rnida turgan bare repo (ishchi fayllarsiz, faqat tarix)
#   aziz/        reviewer: main'ni boshlagan, keyin README'ni yangilagan
#   lola/        muallif: README yangilanishidan oldin clone qilgan
# Muallif va vaqt qat'iy (vaqt.sh): hash'lar hamma joyda bir xil.
set -e; cd "$(dirname "$0")"; rm -rf server.git lola aziz
rf=--ref-format=files; git init -h 2>&1 | grep -q ref-format || rf=   # eski git'da bu bayroq yo'q
git init -q --bare --template= -b main --object-format=sha1 $rf server.git
sozla() {   # $1 ism, $2 email: imzo, hook, rang va global sozlamalar natijani o'zgartirmasin
  git config user.name "$1"; git config user.email "$2"; git config commit.gpgsign false
  git config core.hooksPath /dev/null; git config core.pager cat; git config color.ui false
  git config core.abbrev 7; git config core.autocrlf false; git config core.excludesFile /dev/null
  git config core.attributesFile /dev/null; git config log.decorate false; git config branch.sort refname
  git config pull.ff only; git config merge.stat true; git config push.default simple
  git config advice.detachedHead false; git remote set-url origin ../server.git 2>/dev/null || true; }

# aziz main'ni boshlaydi va serverga push qiladi
mkdir aziz; cd aziz; git init -q --template= -b main --object-format=sha1 $rf
git remote add origin ../server.git; sozla Aziz aziz@example.com
cp ../fayllar/README.md README.md; cp ../fayllar/v0/main.go main.go; printf 'module narx\n\ngo 1.22\n' > go.mod
git add .; . ../vaqt.sh; git commit -q -m "narx dasturi"
git push -q -u origin main 2>/dev/null; cd ..

# lola shu holatda clone qiladi
git clone -q --template= server.git lola; (cd lola; sozla Lola lola@example.com)

# aziz keyin main'da README'ni yangilaydi (lola buni hali fetch qilmagan)
cd aziz; cp ../fayllar/README.aziz.md README.md
git add README.md; . ../vaqt.sh; git commit -q -m "README: ishga tushirish"
git push -q origin main 2>/dev/null; cd ..
echo "tayyor: server.git/ (GitHub o'rnida), lola/ (muallif), aziz/ (reviewer)"
