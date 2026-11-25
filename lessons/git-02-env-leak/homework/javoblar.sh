#!/usr/bin/env bash
# GIT-1.02 uy vazifasi javoblari. Ishga tushirish: bash homework/javoblar.sh (dars papkasidan).
# Repo'larni noldan quradi (qur.sh). Guru qismi uchun git-filter-repo kerak (>= 2.47).
set -e; cd "$(dirname "$0")/.."
./qur.sh >/dev/null

echo "== Oson: kalit tarixda"
(cd loyiha && git show HEAD~1:.env && git grep TOLOV_KALIT= $(git rev-list --all))

echo; echo "== O'rta: .env* yashirin, .env.example ko'rinadi (+ videodagi tuzoq qadamlari)"
(cd tuzoq
 printf '.env*\n!.env.example\n' > .gitignore
 git rm -q --cached .env; cp ../fayllar/env.example .env.example; touch .env.local .env.production
 git add .gitignore .env.example; . ../vaqt.sh; git commit -qm ".env kuzatuvdan chiqdi, .env.example qo'shildi"
 git status --short --ignored
 git ls-files
 echo -n "HEAD~1:.env -> "; git show HEAD~1:.env)

echo; echo "== Qiyin: git log -S (kalit qatori qo'shilgan va o'chirilgan commit'lar)"
(cd loyiha && git log --oneline -S TOLOV_KALIT= && echo -n "birinchi paydo bo'lgan: " && git log --reverse --format='%h %s' -S TOLOV_KALIT= | head -1)

echo; echo "== Guru: yangi clone + filter-repo"
if ! git filter-repo --version >/dev/null 2>&1; then
  echo "git-filter-repo topilmadi (2.47+ kerak): https://github.com/newren/git-filter-repo/blob/main/INSTALL.md" >&2; exit 3
fi
rm -rf toza; git clone -q loyiha toza; cd toza
git filter-repo --sensitive-data-removal --invert-paths --path .env >/dev/null 2>&1
git log --oneline
echo ".env tegib o'tgan commit'lar: $(git log --all --oneline -- .env | wc -l)"
echo "tarixda kalit (mazmun bo'yicha): $(git log --all -p -S soxta_9f2c7a41e8b3 | wc -l) qator"
