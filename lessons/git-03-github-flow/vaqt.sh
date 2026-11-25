# Manba qiling: `. ../vaqt.sh` (lola/ yoki aziz/ ichidan). Commit vaqtini qat'iy qiladi:
# 2026-03-02 10:00 (+05:00) va repo'da ko'rinadigan har bir commit uchun +10 daqiqa.
# Shuning uchun commit hash'lari sizda ham, bizda ham bir xil chiqadi.
# muallif repo sozlamasidan olinsin (user.name/email), muhitdagi qiymatdan emas
unset GIT_AUTHOR_NAME GIT_AUTHOR_EMAIL GIT_COMMITTER_NAME GIT_COMMITTER_EMAIL
n=$(git rev-list --all --count 2>/dev/null || echo 0)
t=$((1772427600 + n * 600))
export GIT_AUTHOR_DATE="@$t +0500" GIT_COMMITTER_DATE="@$t +0500"
