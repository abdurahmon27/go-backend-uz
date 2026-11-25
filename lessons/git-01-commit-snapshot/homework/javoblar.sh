#!/usr/bin/env bash
# GIT-1.01 uy vazifasi javoblari. Ishga tushirish: bash homework/javoblar.sh (dars papkasidan).
# ombor/ ni noldan quradi, videodagi qadamlarni takrorlaydi, keyin har bir javobni ko'rsatadi.
set -e; cd "$(dirname "$0")/.."
./qur.sh >/dev/null
cd ombor
git branch chegirma; git rev-parse HEAD~1 > .git/refs/heads/eski
git switch -q chegirma; sed -i.bak 's/3000/0/' main.go; rm -f main.go.bak
. ../vaqt.sh; git commit -qam "chegirma: 10 foiz"
git switch -q main; git merge -q chegirma
. ../vaqt.sh; git revert --no-edit HEAD >/dev/null

echo "== Oson: yangi commit va uning ichi"
echo "sut: 9000 so'm" > sut.txt; git add sut.txt
. ../vaqt.sh; git commit -qm "sut narxi"
git cat-file -p HEAD

echo; echo "== O'rta: chegirma to'g'ri, 10 foiz (3000 -> 2700)"
sed -i.bak 's/3000/2700/' main.go; rm -f main.go.bak
. ../vaqt.sh; git commit -qam "chegirma: 10 foiz, to'g'ri"
go run main.go

echo; echo "== Qiyin: branch o'chadi, commit qoladi"
e=$(git rev-parse eski)
git branch -d eski
git merge-base --is-ancestor "$e" main && echo "$e main tarixida bor: commit yo'qolmadi"
echo "eski uchun reflog bormi? $(ls .git/logs/refs/heads/ | grep -cx eski || true) (qo'lda yozilgan edi)"

echo; echo "== Guru: --amend dan oldin va keyin (faqat push qilinmagan commit'da!)"
echo "commit: $(git rev-parse HEAD)  tree: $(git rev-parse 'HEAD^{tree}')  parent: $(git rev-parse HEAD~1)"
. ../vaqt.sh; git commit -q --amend -m "chegirma: o'n foiz"
echo "commit: $(git rev-parse HEAD)  tree: $(git rev-parse 'HEAD^{tree}')  parent: $(git rev-parse HEAD~1)"
echo "eski commit reflog'da: $(git reflog -2 --format='%h %gs' | tail -1)"
