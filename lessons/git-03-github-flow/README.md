# GIT-1.03 · GitHub flow: branch, pull request va odobli code review

Oldingi dars: [GIT-1.02 · `.env`ni push qildingizmi?](../git-02-env-leak/) · Keyingi Git darsi: GIT-2.01 · `.git` ichiga kiramiz: commit hash'ini Go'da hisoblaymiz
Bu dars [GIT-1.01](../git-01-commit-snapshot/) ustiga quriladi: branch — bitta hash yozilgan ref, merge — shu ref'ni yangilash.

| Fayl | Nima |
|---|---|
| `qur.sh` | uchta repo'ni noldan quradi: `server.git/` (GitHub oʻrnida turgan bare repo), `aziz/` (reviewer), `lola/` (muallif) |
| `vaqt.sh` | commit vaqtini qatʼiy qiladi (2026-03-02 10:00 +05:00, har commit +10 daqiqa): hash'lar hamma joyda bir xil |
| `fayllar/v0/` | `main`dagi boshlangʻich dastur (`non: 3000 so'm`) |
| `fayllar/v1/` | Lolaning birinchi PR versiyasi: `chegirma(narx, foiz int) int` — foiz 100 dan oshsa narx manfiy |
| `fayllar/v2/` | review'dan keyingi versiya: foiz 0..100 tekshiriladi, xato qaytadi, doc comment va yangi testlar |
| `fayllar/pr.md` | PR tavsifi namunasi (Nima / Nega / Qanday tekshirdim / Qayerga qarash kerak) |
| `fayllar/izoh.md` | review izohlari namunasi: `issue`, `nit (non-blocking)`, `question`, `praise` |
| `fayllar/javob.md` | muallifning javoblari namunasi |
| `fayllar/gh.sh` | haqiqiy GitHub'dagi `gh` buyruqlari (**ishga tushirilmaydi**, faqat koʻrsatiladi) |
| `homework/javoblar.sh` | uy vazifasi javoblari: `bash homework/javoblar.sh` |

Hech narsa GitHub'ga yuborilmaydi: "server" — shu papkadagi `server.git`. `server.git/`, `lola/`, `aziz/` ichida
oʻz `.git` papkasi bor, shuning uchun ular `.gitignore`da. Talablar: git ≥ 2.29, Go ≥ 1.22. `qur.sh` SHA-1, `files`
ref formati va `commit.gpgsign=false`ni majburlaydi, global sozlamangiz natijani oʻzgartirmaydi.

## Videodagi buyruqlar

Hammasi shu papkadan ketma-ket. Natijalar git 2.54.0, go1.26.8 linux/amd64 da oʻlchangan (`go test` vaqti farq qiladi).

```bash
./qur.sh
ls -C server.git lola                              # server.git: config HEAD objects refs — ishchi fayllar yoʻq
for r in aziz lola; do echo "$r:"; git -C $r log --oneline; done
# aziz: 9c274de README: ishga tushirish / 3644e95 narx dasturi    lola: 3644e95 narx dasturi

cd lola && git switch -c chegirma                  # Switched to a new branch 'chegirma'
cd lola && cp ../fayllar/v1/*.go . && go test      # PASS
cd lola && git add . && . ../vaqt.sh && git commit -q -m "chegirma: foiz bilan narx" && git push -u origin chegirma
#  * [new branch]      chegirma -> chegirma
git -C server.git branch                           #   chegirma / * main

cd aziz && git fetch -q && git log --oneline main..origin/chegirma   # 360db8f chegirma: foiz bilan narx
cd aziz && git diff --stat main..origin/chegirma   # README.md | 2 -- ... 4 files changed   <- IKKI nuqta
cd aziz && git diff --stat main...origin/chegirma  # 3 files changed, 36 insertions(+), 3 deletions(-)  <- UCH nuqta = PR
cd aziz && git switch -q chegirma && go test && go run . 10   # PASS / non: 2700 so'm
cd aziz && go run . 150                            # non: -1500 so'm   <- review topgan xato

cd lola && cp ../fayllar/v2/*.go . && go test && go run . 150
# xato: foiz 0..100 oralig'ida bo'lishi kerak, berildi: 150
cd lola && . ../vaqt.sh && git commit -qam "chegirma: foiz 0..100 oralig'ida" && git push   # 360db8f..39feca2
cd aziz && git pull -q --ff-only && go run . 150

cd aziz && git switch -q main && . ../vaqt.sh && git merge -q --no-ff -m "Merge pull request #2 from lola/chegirma" origin/chegirma && git log --graph --oneline main~2..main
# *   26b4804 Merge pull request #2 from lola/chegirma
# |\
# | * 39feca2 chegirma: foiz 0..100 oralig'ida
# | * 360db8f chegirma: foiz bilan narx
# * 9c274de README: ishga tushirish
cd aziz && git push -q origin main && git push origin --delete chegirma    #  - [deleted]  chegirma
cd lola && git switch -q main && git pull -q --ff-only && git branch -d chegirma     # Deleted branch chegirma (was 39feca2).
```

Haqiqiy GitHub'da xuddi shu sikl (`fayllar/gh.sh`, GitHub CLI 2.102 `--help` boʻyicha tekshirilgan, ishga tushirilmagan):

```bash
gh pr create --base main --title "chegirma: foiz bilan narx" --body-file pr.md --reviewer aziz
gh pr checkout 2                                     # reviewer PR branch'ini oʻziga oladi
gh pr review 2 --request-changes --body-file izoh.md # umumiy review; qatorga izoh — veb sahifada
gh pr review 2 --approve --body "Rahmat, tayyor!"
gh pr merge 2 --merge --delete-branch                # --squash yoki --rebase ham bor (GIT-4.03)
git fetch origin pull/2/head:pr-2                    # PR ham ref: refs/pull/2/head (faqat oʻqish uchun)
```

## Dars matni

**Yomon izoh.** "Bu kod notoʻgʻri" — qayeri, nega va nima qilish kerakligi aytilmagan. Yaxshi izohda belgi
(`issue (blocking):`), aniq misol (`go run . 150 → -1500`) va buyruq emas, taklif bor.

**GitHub flow** (GitHub Docs): branch → commit va push → pull request → review izohlariga javob → merge → branch'ni
oʻchirish. Pull request — git buyrugʻi emas, GitHub xizmati; lekin u GIT-1.01dagi branch va merge ustiga qurilgan.
(Git'ning oʻzida `git request-pull` bor: u oʻzgarishlar xulosasini yozib, "shu yerdan olib keting" deb soʻraydi —
Linux yadrosi maintainer'lari pull request'ni shunday, email orqali yuboradi: [kernel docs](https://docs.kernel.org/maintainer/pull-requests.html).)

**Mashq.** `server.git` — bare repo: ishchi fayllar yoʻq, faqat tarix (Pro Git 4.2). Aziz `main`ga README qatorini
qoʻshgan, Lola esa undan oldin clone qilgan — bu farq diff'da kerak boʻladi.

**Branch.** Qisqa, ishni koʻrsatadigan nom (`chegirma`; GitHub Docs: `increase-test-timeout`). `git push -u` lokal
branch'ni serverdagisiga bogʻlaydi (upstream). Serverda yangi branch, `main` joyida.

**PR tavsifi.** Nima oʻzgardi, nega, qanday tekshirildi, qayerga qarash kerak (GitHub Docs: "Provide context and
guidance", "Review your own pull request first"). `Closes #1` — PR default branch'ga merge qilinganda issue oʻzi
yopiladi. Tugamagan ish — **draft** PR (merge qilib boʻlmaydi, code owner'lar avtomatik chaqirilmaydi).

**Ikki nuqta va uch nuqta.** `git diff main..chegirma` — `main`ning *hozirgi* holati bilan solishtiradi, shuning
uchun Azizning README qatorlari "oʻchirilgan"dek koʻrinadi. `git diff main...chegirma` — umumiy ajdoddan (merge
base) boshlab: faqat PR olib kelgan oʻzgarish. GitHub PR sahifasi uch nuqtali diff'ni koʻrsatadi (GitHub Docs).
Chalkashlik kamroq boʻlishi uchun GitHub `main`ni branch'ga tez-tez merge qilishni va PR'ni tez yakunlashni
maslahat beradi. Hajm: 36 qator qoʻshildi, 3 qator oʻchdi. Google: "100 lines is usually a reasonable size for a CL, and 1000 lines is usually too large".

**Diff'ni oʻqish yetmaydi.** Lola tavsifda 0 va 10 foizni tekshirganini yozgan; Aziz yozilmagan holatni sinaydi:
150 foiz → `non: -1500 so'm`.

**Review qarori** (GitHub Docs): *Comment* — umumiy fikr; *Approve* — merge'ga tayyor; *Request changes* —
merge'dan oldin tuzatish kerak. Bu yerda — Request changes. Lekin Google: "there is no such thing as 'perfect'
code — there is only *better* code": kod sogʻligʻini aniq yaxshilayotgan oʻzgarishni mayda did uchun ushlab turmang.

**Yaxshi izoh** (Google "How to write code review comments" + Conventional Comments):
1. Odam haqida emas, kod haqida ("Why did **you** use threads…" — yomon misol Google'dan).
2. Nega muhim va aniq misol.
3. Yechim — odatda taklif yoki savol (tuzatish muallifning ishi); Google: muammoni koʻrsatish va aniq koʻrsatma
   berish orasida muvozanat — baʼzan aniq koʻrsatma yoki kod foydaliroq.
4. Belgi: `issue:` — muammo; majburiyligini `(blocking)` / `(non-blocking)` aytadi (Conventional Comments'da bu
   "decoration"). `nit:` — Google'ning "Nit:" prefiksi (Conventional Comments'da `nitpick:`), mayda va PR'ni
   toʻxtatmaydi. `question:` — shubha boʻlsa. `praise:` — samimiy maqtov (yolgʻon maqtov emas).
5. GitHub'da ```` ```suggestion ```` bloki: muallif "Commit suggestion" bilan qabul qiladi, taklif qilgan odam
   commit'ga co-author boʻladi.

**Javob berish** (Google "How to handle reviewer comments"): shaxsiy qabul qilmang; jahl bilan hech qachon javob
yozmang ("will live forever in the code review tool"); tushunarsiz joyni izohda emas, kodda aniqlang; tuzatgach —
nima qildingiz va qaysi commit'da; rozi boʻlmasangiz — sabab va savol; PR doirasidan tashqari taklif — alohida issue
(GitHub Docs). Hal boʻlgach — **Resolve conversation**.

**Tuzatish** shu branch'ga push qilinadi: PR oʻzi yangilanadi (GitHub Docs). `go run` oxiridagi `exit status 2` —
dasturning chiqish kodi (SYS-0.02). Lola review'ni qayta soʻraydi (Re-request review), Aziz tuzatishni tekshiradi.
Reviewer tomoni: Google boʻyicha review soʻroviga javob "bir ish kunidan" oshmasin.

**Merge.** Himoyalangan branch'da (branch protection) merge uchun tasdiq soni va oʻtgan CI tekshiruvlari
talab qilinishi mumkin (GitHub Docs: About protected branches). GitHub'ning odatiy "Merge pull request" tugmasi `--no-ff` bilan ishlaydi: fast-forward mumkin boʻlsa ham
merge commit yaratiladi (GitHub Docs). Bizda `main`da Azizning commit'i bor edi, shuning uchun merge commit baribir
kerak edi. Keyin server'dagi branch oʻchiriladi (PR va commit'lar qoladi), Lola `main`ni pull qilib, oʻz
branch'ini `git branch -d` bilan oʻchiradi. GitHub tugmasi PR'ni yopadi va `Closes #1` yozilgan issue'ni ham yopadi;
mashqda biz faqat git qismini qoʻlda takrorlaymiz (`#1` — chegirma soʻralgan issue, misol uchun). Squash va rebase
tugmalari — GIT-4.03. `qur.sh` `pull.ff only` qoʻyadi; videoda ham `git pull --ff-only` aniq yozilgan.

## Manbalar

- GitHub Docs: [GitHub flow](https://docs.github.com/en/get-started/using-github/github-flow) ·
  [Pull request reviews](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/about-pull-request-reviews) ·
  [Reviewing proposed changes](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/reviewing-proposed-changes-in-a-pull-request) ·
  [Commenting on a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/commenting-on-a-pull-request) ·
  [Incorporating feedback](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/incorporating-feedback-in-your-pull-request) ·
  [Helping others review your changes](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/getting-started/helping-others-review-your-changes) ·
  [Comparing branches (three-dot)](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-comparing-branches-in-pull-requests) ·
  [About pull request merges](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges) ·
  [About protected branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches) ·
  [Linking a PR to an issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue) ·
  [About pull requests (draft, PR refs)](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests) ·
  [Checking out PRs locally](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/checking-out-pull-requests-locally)
- Google Engineering Practices: [How to write code review comments](https://google.github.io/eng-practices/review/reviewer/comments.html) ·
  [The Standard of Code Review](https://google.github.io/eng-practices/review/reviewer/standard.html) ·
  [How to handle reviewer comments](https://google.github.io/eng-practices/review/developer/handling-comments.html) ·
  [Small CLs](https://google.github.io/eng-practices/review/developer/small-cls.html) ·
  [Speed of code reviews](https://google.github.io/eng-practices/review/reviewer/speed.html)
- [Conventional Comments](https://conventionalcomments.org/)
- git: [git-diff (A...B)](https://git-scm.com/docs/git-diff) · [git-merge --no-ff](https://git-scm.com/docs/git-merge) ·
  [git-push -u](https://git-scm.com/docs/git-push) · [git-request-pull](https://git-scm.com/docs/git-request-pull) · [Linux kernel: pull requests](https://docs.kernel.org/maintainer/pull-requests.html) ·
  [Pro Git 4.2: bare repo](https://git-scm.com/book/en/v2/Git-on-the-Server-Getting-Git-on-a-Server)
- [GitHub CLI manual: gh pr](https://cli.github.com/manual/gh_pr)

## Uy vazifasi

1. **Oson.** `./qur.sh` dan boshlab videodagi buyruqlarni takrorlang va `git log --graph --oneline`da merge commit'ni oling.
2. **Oʻrta.** `git merge-base main origin/chegirma` bilan umumiy ajdodni toping. `git diff <ajdod> origin/chegirma`
   uch nuqtali `git diff main...origin/chegirma` bilan bir xil ekanini tekshiring.
3. **Qiyin.** Oʻz GitHub repo'ngizda haqiqiy PR oching (`gh pr create`) va uni oʻzingiz approve qilib koʻring. Nima boʻladi?
4. **Guru.** GitHub'ning "Squash and merge" tugmasini `git merge --squash` + `git commit` bilan takrorlang. Grafik va
   parent'lar soni `--no-ff` merge'dan qanday farq qiladi? Fayllar bir xilmi?

### Javoblar

`bash homework/javoblar.sh` (hammasi lokal, ikki marta ishga tushirildi — natija bayt-bayt bir xil):

```
== O'rta: merge-base = uch nuqtali diff'ning boshlanish nuqtasi
umumiy ajdod: 3644e95 narx dasturi
git diff $(merge-base) origin/chegirma  -> 5a24ca4c0c59320755e67421cdf02912246b0a4c
git diff main...origin/chegirma          -> 5a24ca4c0c59320755e67421cdf02912246b0a4c
bir xil: A...B = merge-base(A,B)..B

== Oson: GitHub'ning 'Create a merge commit' tugmasi (--no-ff)
*   26b4804 Merge pull request #2 from lola/chegirma
|\
| * 39feca2 chegirma: foiz 0..100 oralig'ida
| * 360db8f chegirma: foiz bilan narx
* | 9c274de README: ishga tushirish
|/
* 3644e95 narx dasturi
parent'lar: 9c274de 39feca2

== Guru: 'Squash and merge' (git merge --squash + commit)
* 92daa53 chegirma: foiz bilan narx (#2)
* 9c274de README: ishga tushirish
* 3644e95 narx dasturi
parent'lar: 9c274de  (bitta: merge commit emas)
ikki tugma natijasi bir xil fayllarmi? ha
```

- **Qiyin:** boʻlmaydi. GitHub Docs: "Pull request authors cannot approve their own pull requests." Oʻz PR'ingizga faqat
  Comment qoldira olasiz; majburiy approve yoqilgan repo'da merge uchun boshqa odamning tasdigʻi kerak.
- **Guru:** squash bitta oddiy commit (bitta parent) yaratadi, Lolaning ikki commit'i `main` tarixida koʻrinmaydi;
  fayllar esa ikkala usulda bir xil. GitHub Docs: squash qisqa umrli branch'lar uchun yaxshi; uzoq yashaydigan
  branch'da keyingi PR'lar allaqachon squash qilingan commit'larni qayta olib kelishi mumkin.
