# GIT-1.01 · git: commit — bu snapshot, branch — faylga yozilgan bitta hash

Oldingi dars: [SYS-0.02 · Jarayon va muhit](../sys-02-process-env/) · Keyingi dars: [GIT-1.02 · `.env`ni push qildingizmi?](../git-02-env-leak/)

| Fayl | Nima |
|---|---|
| `qur.sh` | `ombor/` mashq repo'sini noldan quradi: 2 ta commit (`README`, `narx dasturi`) |
| `vaqt.sh` | commit vaqtini qatʼiy qiladi (2026-01-05 10:00 +05:00, har commit +10 daqiqa): hash'lar hamma joyda bir xil |
| `fayllar/` | repo'ga koʻchiriladigan fayllar (`README.md`, `main.go`) |
| `homework/javoblar.sh` | uy vazifasi javoblari: `bash homework/javoblar.sh` |

`ombor/` ichida oʻz `.git` papkasi bor, shuning uchun u `.gitignore`da: kurs repo'siga qoʻshilmaydi. Uni istalgan
payt `./qur.sh` bilan qayta yaratasiz.

**Talablar.** git ≥ 2.29 (`git init -b` 2.28 dan, `--object-format` 2.29 dan). `qur.sh` repo'ni SHA-1 va `files`
ref formatida yaratadi va unda `commit.gpgsign=false` qiladi: global sozlamangizda SHA-256, reftable yoki commit
imzosi yoqilgan boʻlsa ham hash'lar videodagidek chiqadi (biz shunday global sozlama bilan sinab koʻrdik).
Windows'da **Git Bash**'da ishlating: PowerShell 5 dagi `>` faylni UTF-16 da yozadi va qoʻlda yozilgan ref buziladi.

## Videodagi buyruqlar

Hammasi shu papkadan (`lessons/git-01-commit-snapshot`) ketma-ket ishga tushiriladi. Natijalar git 2.54 da oʻlchangan.

```bash
./qur.sh                                                   # daac3a4 narx dasturi / e657f41 README
cd ombor && git branch chegirma && cat .git/refs/heads/chegirma   # daac3a47cb72faf965c5385ab1002bc04bf89b45
cd ombor && wc -c .git/refs/heads/chegirma && tail -c 3 .git/refs/heads/chegirma | od -An -c   # 41 … 4 5 \n
cd ombor && git cat-file -p HEAD                           # tree, parent, author, committer, xabar
cd ombor && git log --format='%h  parent: %p  %s'          # zanjir: birinchi commit'da parent yoʻq
cd ombor && git cat-file -p 'HEAD^{tree}'                  # 100644 blob <hash> README.md / main.go
cd ombor && git cat-file -p HEAD:main.go                   # faylning oʻzi
cd ombor && git rev-parse HEAD:README.md HEAD~1:README.md  # ikkalasi 51c4fee9…: bitta obyekt
cd ombor && git cat-file --batch-all-objects --batch-check='%(objecttype) %(objectname)' | cut -c1-20   # 2 commit, 2 tree, 2 blob
printf 'non: 3000\n' | git hash-object --stdin             # 2ecbbc99c0932c7b1430e07e274c706f3afb62ef
printf 'blob 10\0non: 3000\n' | sha1sum                    # xuddi shu: hash = SHA-1("blob <uzunlik>\0" + mazmun); macOS: shasum
printf 'non: 3001\n' | git hash-object --stdin             # a3fad3e2…: butunlay boshqa
cd ombor && cat .git/HEAD                                  # ref: refs/heads/main
cd ombor && cat .git/refs/heads/main                       # daac3a47…
cd ombor && git rev-parse HEAD~1 > .git/refs/heads/eski && git branch   # faqat tushunish uchun!
cd ombor && ls .git/logs/refs/heads/                       # chegirma, main: eski'da reflog yoʻq
cd ombor && git switch -q chegirma && sed -i.bak 's/3000/0/' main.go && rm main.go.bak && git status --short   #  M main.go
cd ombor && git add main.go && . ../vaqt.sh && git commit -m "chegirma: 10 foiz"   # [chegirma 5a239b1]
cd ombor && git log --oneline --decorate --all             # 5a239b1 (HEAD -> chegirma) …
cd ombor && git switch -q main && git merge chegirma       # Fast-forward
cd ombor && cat .git/refs/heads/main .git/refs/heads/chegirma     # ikkalasi 5a239b10…
cd ombor && go run main.go                                 # non: 0 so'm
cd ombor && . ../vaqt.sh && git revert --no-edit HEAD      # [main 14b041d] Revert "chegirma: 10 foiz"
cd ombor && git log --oneline                              # 4 ta commit, hech biri yoʻqolmadi
cd ombor && go run main.go                                 # non: 3000 so'm
```

`. ../vaqt.sh` faqat vaqtni qatʼiy qilish uchun (shuning uchun `git revert` chiqishida ` Date:` qatori bor: muallif
vaqti alohida berilgan). Oʻz loyihangizda shunchaki `git commit` yozasiz, avval bir marta
`git config --global user.name "Ism"` va `git config --global user.email "siz@example.com"`.

## Dars matni

**Commit — snapshot.** Git har commit'da butun loyihaning suratini saqlaydi, oʻzgarishlar roʻyxatini emas: Pro Git
buni "a series of snapshots of a miniature filesystem" deydi. `git show` va `git diff` koʻrsatadigan diff ikki surat
orasidan har safar hisoblab chiqariladi. Commit'da yana muallif, vaqt, xabar va oldingi commit (`parent`) hash'i bor.

**Obyektlar.** `.git/objects` — "content-addressable" omborxona: har obyekt oʻz mazmunining hash'i bilan saqlanadi.
Uch asosiy tur:
- **blob** — fayl mazmuni (nomi ham, vaqti ham yoʻq);
- **tree** — papka: har qatorda rejim, tur, hash va nom. Rejim — ruxsatlar emas, fayl turi: `100644` oddiy fayl,
  `100755` bajariladigan fayl, `120000` symlink, `040000` ichki papka (tree);
- **commit** — `tree`, `parent`, `author` (kodni yozgan), `committer` (commit'ni yaratgan) va xabar. Vaqt Unix
  soniyalarida, yonida UTC'dan farq (`+0500`, bizda Toshkent).

Obyektlar zlib bilan siqilgan, shuning uchun `cat` emas, `git cat-file -p` bilan oʻqiladi. Mazmuni bir xil fayl
bitta obyekt boʻladi: ikkinchi commit'da README oʻzgarmagani uchun yangi blob yaratilmadi (6 ta obyekt).

**Packfile.** Yangi obyektlar avval alohida ("loose") fayl boʻlib yotadi. `git gc`, push va **clone** ularni
*packfile*'ga yigʻadi; clone qilingan repo'da obyektlar darhol paketda. Paket ichida oʻxshash obyektlar farq (delta)
koʻrinishida siqiladi. Bu faqat saqlash usuli: `git cat-file -p` baribir toʻliq tree va toʻliq faylni qaytaradi,
commit modeli snapshot boʻlib qoladi (Pro Git 10.4).

**Hash.** Git obyektning sarlavhasi (`<tur> <uzunlik>\0`) va mazmunidan SHA-1 oladi (40 ta hex belgi): videoda
`git hash-object` va `sha1sum` bir xil natija berdi. Commit mazmunida `parent` hash'i bor, shuning uchun eski commit
oʻzgarsa, undan keyingi hamma commit hash'lari oʻzgaradi. SHA-1 eskirgan (BreakingChanges: "only a matter of time
before SHA-1 will be considered broken completely"); Git 3.0 da yangi repo'lar uchun standart SHA-256 boʻlishi
rejalashtirilgan (sana yoʻq). SHA-256 repo'da hash 64 belgi, branch fayli esa 65 bayt. Hash'ni Go'da oʻzimiz
hisoblash — GIT-2.01.

**Branch — bitta hash yozilgan fayl.** `.git/refs/heads/main` faylida 40 belgilik hash va qator oxiri (`\n`): 41 bayt.
`git branch` bundan tashqari reflog faylini ham yaratadi (`.git/logs/refs/heads/chegirma`): branch qayerlarda
boʻlganining jurnali (GIT-3.04). `.git/HEAD` odatda hash emas, boshqa ref'ga havola (`ref: refs/heads/main`), yaʼni
hozir qaysi branch'dasiz. Commit'ni toʻgʻridan-toʻgʻri tanlasangiz (`git switch --detach`), HEAD'da hash turadi:
"detached HEAD".

Videodagi `> .git/refs/heads/eski` — faqat tushunish uchun tajriba. Hayotda `git branch` yoki `git update-ref`
ishlating:
- ular qiymatni tekshiradi va reflog yozadi (qoʻlda yozilgan `eski`da reflog yoʻq);
- ref'lar `git gc` yoki `git pack-refs --all` dan keyin `.git/packed-refs` faylida turishi mumkin (`refs/heads` boʻsh qoladi);
- Git 3.0 da yangi repo'larda ref'lar standart boʻyicha **reftable** formatida saqlanishi rejalashtirilgan.

**Merge.** `main` chegirma ochilgandan beri oldinga yurmagan edi, shuning uchun `git merge chegirma` **fast-forward**
qildi: yangi commit yoʻq, `main` fayliga chegirma hash'i yozildi va ishchi fayllar ham yangilandi. Ikkala branch'da
yangi commit boʻlsa, ikki parent'li merge commit paydo boʻladi — GIT-2.03.

**Revert.** `git revert` xatoli commit'ning teskarisini **yangi commit** sifatida yozadi; tarix oʻchmaydi, shuning
uchun jamoa ishlaydigan branch'da xavfsiz. `git reset` esa branch faylini orqaga suradi (tarixni qayta yozadi) —
GIT-3.03. **Diqqat:** revert eski commit'ni oʻchirmaydi. Commit'ga parol yoki kalit tushgan boʻlsa, u tarixda
qoladi — GIT-1.02.

## Manbalar

- [Pro Git 1.3: What is Git? (snapshots, not differences)](https://git-scm.com/book/en/v2/Getting-Started-What-is-Git%3F)
- [Pro Git 10.2: Git Objects](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects), [10.3: Git References](https://git-scm.com/book/en/v2/Git-Internals-Git-References), [10.4: Packfiles](https://git-scm.com/book/en/v2/Git-Internals-Packfiles)
- [git-cat-file](https://git-scm.com/docs/git-cat-file), [git-hash-object](https://git-scm.com/docs/git-hash-object), [git-switch](https://git-scm.com/docs/git-switch), [git-merge](https://git-scm.com/docs/git-merge), [git-revert](https://git-scm.com/docs/git-revert), [git-update-ref](https://git-scm.com/docs/git-update-ref), [git-pack-refs](https://git-scm.com/docs/git-pack-refs)
- [Git 3.0 BreakingChanges: SHA-256, reftable, `main`](https://git-scm.com/docs/BreakingChanges)
- [MIT Missing Semester: Version Control (Git)](https://missing.csail.mit.edu/2020/version-control/)
- Julia Evans: [Inside .git](https://jvns.ca/blog/2024/01/26/inside-git/), [Do we think of git commits as diffs, snapshots, and/or histories?](https://jvns.ca/blog/2024/01/05/do-we-think-of-git-commits-as-diffs--snapshots--or-histories/)

## Uy vazifasi

1. 🟢 **Oson.** `ombor`da yangi fayl qoʻshib commit qiling va uni `git cat-file -p HEAD` bilan oching.
2. 🟡 **Oʻrta.** Chegirmani toʻgʻri, 10 foiz qilib yozing: `non: 2700 so'm`.
3. 🟠 **Qiyin.** `git branch -d eski`: branch oʻchadi. Commit-chi?
4. 🔴 **Guru.** Hali push qilinmagan oxirgi commit xabarini `git commit --amend` bilan oʻzgartiring. Qaysi hash'lar
   oʻzgaradi, qaysilari qoladi?

<details><summary>Javoblar (<code>bash homework/javoblar.sh</code>)</summary>

1. `git add` + `git commit`, keyin `git cat-file -p HEAD`: `tree 530f9ba8…`, `parent 14b041d6…` (revert commit'i), `author Lola …`.
2. `sed -i.bak 's/3000/2700/' main.go`, `git commit -am "…"`, `go run main.go` → `non: 2700 so'm`.
3. `Deleted branch eski (was e657f41).` Commit esa joyida va yetib boriladi:
   `git merge-base --is-ancestor e657f41 main` muvaffaqiyatli, chunki u `main`ning ajdodi. Branch faqat havola edi.
   (`git cat-file -t` yetarli dalil emas: yetib boʻlmaydigan commit ham `gc` uni tozalaguncha `commit` deydi.)
   `-d` faqat merge qilingan branch'ni oʻchiradi; merge qilinmaganini `-D` oʻchiradi. Unda commit'ni odatda
   `git reflog` orqali topasiz, lekin qoʻlda yozilgan `eski`ning reflog'i yoʻq: faqat `git fsck --lost-found` yordam beradi.
4. Yangi commit hash'i boshqa (`914b512…` → `cdd8662…`), chunki xabar (va bizda committer vaqti) commit mazmunida.
   `tree` (`4cf5d8f7…`) va `parent` (`5f4efe00…`) oʻzgarmaydi: fayllar ham, oldingi tarix ham bir xil. Eski commit
   oʻchmaydi, `git reflog`da turadi. Push qilingan commit'ni amend qilmang: bu tarixni qayta yozish, hamkasblarga
   muammo (GIT-2.05, GIT-3.03).

</details>
