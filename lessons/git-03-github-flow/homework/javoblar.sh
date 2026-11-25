#!/usr/bin/env bash
# GIT-1.03 uy vazifasi javoblari. Ishga tushirish: bash homework/javoblar.sh (dars papkasidan).
# Repo'larni noldan quradi (qur.sh), videodagi qadamlarni merge'gacha takrorlaydi, keyin javoblarni ko'rsatadi.
# GitHub'ga hech narsa yubormaydi: hammasi lokal server.git ustida.
set -e; cd "$(dirname "$0")/.."
./qur.sh >/dev/null
(cd lola && git switch -q -c chegirma && cp ../fayllar/v1/*.go . && git add . && . ../vaqt.sh \
  && git commit -q -m "chegirma: foiz bilan narx" && git push -q -u origin chegirma 2>/dev/null)
(cd aziz && git fetch -q)

echo "== O'rta: merge-base = uch nuqtali diff'ning boshlanish nuqtasi"
(cd aziz
 mb=$(git merge-base main origin/chegirma)
 echo "umumiy ajdod: $(git log -1 --format='%h %s' "$mb")"
 a=$(git diff "$mb" origin/chegirma | git hash-object --stdin)
 b=$(git diff main...origin/chegirma | git hash-object --stdin)
 echo "git diff \$(merge-base) origin/chegirma  -> $a"
 echo "git diff main...origin/chegirma          -> $b"
 [ "$a" = "$b" ] && echo "bir xil: A...B = merge-base(A,B)..B")

# videodagi tuzatish (Request changes'dan keyin)
(cd lola && cp ../fayllar/v2/*.go . && . ../vaqt.sh && git commit -qam "chegirma: foiz 0..100 oralig'ida" && git push -q 2>/dev/null)
(cd aziz && git fetch -q)

echo; echo "== Oson: GitHub'ning 'Create a merge commit' tugmasi (--no-ff)"
(cd aziz && git switch -q -c merge-tugmasi main && . ../vaqt.sh \
 && git merge -q --no-ff -m "Merge pull request #2 from lola/chegirma" origin/chegirma \
 && git log --graph --oneline && echo "parent'lar: $(git log -1 --format=%p)")

echo; echo "== Guru: 'Squash and merge' (git merge --squash + commit)"
(cd aziz && git switch -q -c squash-tugmasi main && git merge -q --squash origin/chegirma >/dev/null 2>&1 && . ../vaqt.sh \
 && git commit -q -m "chegirma: foiz bilan narx (#2)" \
 && git log --graph --oneline && echo "parent'lar: $(git log -1 --format=%p)  (bitta: merge commit emas)" \
 && echo "ikki tugma natijasi bir xil fayllarmi? $(git diff --quiet merge-tugmasi squash-tugmasi && echo ha)")

echo; echo "== Qiyin: o'z PR'ingizni approve qilish"
echo "Bo'lmaydi. GitHub Docs: \"Pull request authors cannot approve their own pull requests.\""
echo "O'z PR'ingizga faqat Comment qoldira olasiz. Repo'da majburiy approve yoqilgan bo'lsa (branch protection),"
echo "merge uchun boshqa odamning tasdig'i kerak: review aynan shuning uchun bor."
