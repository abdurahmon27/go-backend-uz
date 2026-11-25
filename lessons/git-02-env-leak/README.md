# GIT-1.02 · `.env`ni push qildingizmi? Kalit endi faqat sizniki emas

Oldingi dars: [GIT-1.01 · commit — snapshot, branch — bitta hash](../git-01-commit-snapshot/) · Keyingi dars: GIT-1.03 · GitHub flow: branch, pull request va odobli code review

| Fayl | Nima |
|---|---|
| `qur.sh` | ikki mashq repo'sini noldan quradi: `loyiha/` (`.env` commit qilingan, keyin "oʻchirilgan") va `tuzoq/` (`.env` birinchi commit'dan kuzatiladi) |
| `vaqt.sh` | commit vaqtini qatʼiy qiladi (2026-02-02 09:00 +05:00): hash'lar hamma joyda bir xil |
| `fayllar/` | `main.go` (kalitni `TOLOV_KALIT` muhit oʻzgaruvchisidan oʻqiydi), `env.txt`, `env.example` |
| `homework/javoblar.sh` | uy vazifasi javoblari: `bash homework/javoblar.sh` |

Kalit **soxta** (`soxta_9f2c7a41e8b3`): hech qayerda ishlamaydi va hech qaysi provayder formatiga oʻxshamaydi.
`loyiha/`, `tuzoq/`, `toza/` ichida oʻz `.git` papkasi bor, shuning uchun ular `.gitignore`da. git ≥ 2.29 kerak;
`qur.sh` SHA-1, `files` ref formati va `commit.gpgsign=false`ni majburlaydi, shuning uchun global sozlamangiz
(SHA-256, reftable, imzo) natijani oʻzgartirmaydi.

## Videodagi buyruqlar

```bash
./qur.sh
cd loyiha && git log --oneline && ls -A        # 36288b8 kalitni o'chirdim … papkada .env yoʻq
cd loyiha && git show HEAD~1:.env              # TOLOV_KALIT=soxta_9f2c7a41e8b3 — tarixda turibdi
cd loyiha && git log --oneline -- .env         # .env'ga tekkan 2 commit
cd loyiha && git grep TOLOV_KALIT= $(git rev-list --all)   # hamma commit'larda qidirish
cd tuzoq && echo .env > .gitignore && echo DEBUG=1 >> .env && git status --short
#  M .env        <- .gitignore yordam bermadi: fayl allaqachon kuzatilmoqda
cd tuzoq && git rm -q --cached .env && git add .gitignore && . ../vaqt.sh && git commit -qm ".env kuzatuvdan chiqdi" && git status --short --ignored && ls -A
# !! .env       <- endi eʼtiborsiz, fayl diskda qoldi
cd tuzoq && git check-ignore -v .env           # .gitignore:1:.env	.env
cd tuzoq && git show HEAD~1:.env               # eski commit'da kalit baribir bor
```

Katta repo'da `$(git rev-list --all)` argumentlar chegarasidan (ARG_MAX) oshib ketishi mumkin. Oʻrniga
`git log --all -p -S 'TOLOV_KALIT='` yoki `git rev-list --all | xargs git grep 'TOLOV_KALIT='` ishlating.

Tarixni tozalash (videoda kod sifatida koʻrsatilgan):

```bash
# git-filter-repo >= 2.47 kerak (Ubuntu 24.04 apt'da 2.38: pip install git-filter-repo); yangi clone'da ishlang
git clone git@github.com:SIZ/loyiha.git && cd loyiha
git log --all --name-only --format= | grep -E '(^|/)\.env' | sort -u      # tarixdagi HAMMA env fayllar
git filter-repo --sensitive-data-removal --invert-paths --path .env --path api/.env   # har biri toʻliq yoʻl bilan
git log --all -p -S 'SIR_QIYMATI'     # HAR BIR sir qiymati uchun boʻsh chiqishi kerak: yoʻlni emas, mazmunni tekshiring
# jamoa push'ni toʻxtatsin: force push sizning clone'ingizdan keyin push qilingan ishni oʻchiradi
git push --force --mirror origin      # refs/pull xatolari kutilgan; branch protection vaqtincha oʻchiriladi
```

`--path .env` faqat ildizdagi `.env`ni oʻchiradi: `api/.env`, `.env.local` yoki nomi oʻzgargan fayl qoladi (GitHub
hujjati: "the full git path … not just the filename"). `--path-glob` ham ishlaydi. Kalit boshqa fayllarga ham koʻchgan
boʻlsa (`config.yaml`, testlar), `--path` yetmaydi:
`git filter-repo --sensitive-data-removal --replace-text ../kalitlar.txt` (GitHub hujjati).

Bizda (git-filter-repo 2.47.0) `loyiha`ning yangi clone'ida natija: `You rewrote 2 (of 3) commits.` Faqat `.env`ni
qoʻshgan va oʻchirgan ikki commit boʻsh qolib, tarixdan tushib ketdi; birinchi commit (`4590f4e`) unga tegmagani uchun
hash'i oʻzgarmadi. Eski repo'ning oʻzida: `Aborting: Refusing to destructively overwrite repo history since this does
not look like a fresh clone.`

## Dars matni

**Oʻchirish ≠ unutish.** Har commit — butun loyiha surati (GIT-1.01). `git rm .env` + commit faqat *keyingi*
suratlardan oʻchiradi; eski commit'da kalit qoladi va `git show <commit>:.env` uni qaytaradi. Hujumchilar eski
commit'lardagi kalitlarni ham topadi (Unit 42 tajribasida kalit "hidden behind a past commit" edi).

**Qanchalik tez?**
- NDSS 2019 (Meli va boshq., maʼlumot 2017–2018): tadqiqotchilar sinov satrini push qilgan, ularning skaneri uni
  mediana **20 soniyada** topgan (0,5 s dan 4 daqiqadan koʻproqqacha). 100 000 dan ortiq repo'da kalit topilgan.
  Kalitlarning **81 %** kuzatuv davomida (~16 kun) repo'dan oʻchirilmagan. Oʻchirganlarning **hech biri** tarixni
  qayta yozmagan.
- Unit 42 (Palo Alto Networks, 2023) ataylab qoldirgan AWS kalitini hujumchilar **5 daqiqa** ichida topib ishlatgan.
  AWS kalitni 2 daqiqada karantin siyosati bilan cheklagan; tadqiqotchilar bu siyosatni ataylab almashtirgach,
  hujumchilar bir nechta regionda EC2 serverlar ochib, kriptovalyuta qazishgan. Provayderga xabar borishi kalit bekor
  qilindi degani emas.
- GitGuardian 2026: 2022-yilda **ishlashi tasdiqlangan** sizib chiqqan kalitlarning **64 %** dan koʻprogʻi 2026-yil
  yanvarida hali ham ishlagan; 2025-yilda ochiq commit'larda 28,65 mln yangi sir.

**1-qadam: rotate.** GitHub hujjati: "as a first step you need to revoke and/or rotate that secret". Kalit allaqachon
skanerlarda, fork'larda va hamkasblarning clone'larida boʻlishi mumkin. Tarixni tozalash ularga yetmaydi, repo'ni
yopish yoki oʻchirish ham. Toyota kodni 15-sentabrda yopgan, kalitni esa 17-sentabrda almashtirgan: oradagi ikki kun
kalit hali ishlagan.

Tartib:
1. **Fayldagi hamma sirlarni** roʻyxatga oling (`.env`da odatda bir nechta: DB paroli, JWT siri, API kalitlar).
2. Provayderda yangi kalit yarating va ilovaga bering, keyin eskisini bekor qiling (revoke). Tez deploy qila
   olmasangiz, **eskisini darhol bekor qiling**: ishlamay turish oʻgʻirlikdan arzon.
3. Provayder jurnallarini (AWS CloudTrail, DB loglari) **birinchi push vaqtidan** boshlab tekshiring.
4. Begona yaratilgan resurslarni (serverlar, foydalanuvchilar, kalitlar) topib oʻchiring, hisob-fakturani koʻring.
5. Shundan keyin, kerak boʻlsa, tarixni tozalang.

**2-qadam (kerak boʻlsa): tarixni tozalash.** Qachon kerak: almashtirib boʻlmaydigan maʼlumot sizgan boʻlsa
(mijozlar roʻyxati, shaxsiy maʼlumot), oson almashtirilmaydigan kalitlar, yoki skaner ogohlantirishlarini toʻxtatish.
Cheklovlar (GitHub hujjati va bizning sinov):
- force push sizning clone'ingizdan keyin push qilingan ishni (yangi commit'lar, yangi branch'lar) **oʻchirib
  yuboradi**: oldin jamoa push'ni toʻxtatsin;
- force push faqat **serverdagi** tarixni almashtiradi. Hamkasbning eski clone'i oʻzgarmaydi: u `git pull` qilsa
  "ahead 2" boʻladi, keyin oddiy `git push` kalitni serverga **qaytaradi**. Hamma qayta clone qilishi kerak;
- fork'lardagi eski commit'lar qoladi;
- pull request havolalari va keshlangan koʻrinishlarni faqat GitHub Support oʻchiradi, va u faqat kalitni
  almashtirish bilan xavfni yopib boʻlmaydigan hollarda yordam beradi;
- NDSS 2019: tarixi qayta yozilgan commit'ni ham hash'i orqali GitHub'dan ochish mumkin boʻlgan;
- kalit almashtirilgan boʻlsa, "rewrite the history … may not be warranted".

**`.gitignore` tuzogʻi.** `.gitignore` faqat hali **kuzatilmagan** fayllarga taʼsir qiladi. Kuzatilayotgan faylni
chiqarish: `git rm --cached .env` + commit (fayl diskda qoladi). Shundan keyingina yangi kalitni `.env`ga yozing,
aks holda keyingi `git commit -a` uni ham sizdiradi. Hamkasblar bu commit'ni `pull` qilganda ularning `.env` fayli
**oʻchadi**: oldindan nusxa olib qoʻyishsin. Eski commit'lardagi kalit baribir tarixda.

**Toʻgʻri odat.**
- `.gitignore` birinchi commit'dan **oldin**: `.env*` va `!.env.example` (`.env.local`, `.env.production` ham
  yashirinadi), `*.pem`, `id_*`.
- Repo'ga faqat `.env.example` (nomlar bor, qiymat yoʻq).
- Sirlar: kompyuteringizda yoki server/CI'ning maxfiy sozlamalarida (GitHub Actions secrets).
- `git add .` eʼtiborsiz qoldirilmagan hamma narsani qoʻshadi: avval `git status` / `git diff --cached`.
- Commit'dan oldin skaner (gitleaks, trufflehog, pre-commit hook) — GIT-5.04, push protection chuqurroq — GIT-6.03.

**GitHub nima qiladi.**
- **Secret scanning** ochiq repo'larda bepul va avtomatik ishlaydi. Hamkor provayderning kaliti topilsa, unga xabar
  yuboriladi.
- **Push protection for users** standart boʻyicha yoqilgan, lekin faqat **ochiq** repo'larga push'da ishlaydi. Xatoda
  kalit turi, commit va `fayl:qator` koʻrsatiladi. Uni sabab koʻrsatib chetlab oʻtish mumkin; foydalanuvchi darajasida
  chetlab oʻtish ogohlantirish yaratmaydi. Generic satrlarni (bizning `soxta_…` kabi) tanimasligi mumkin.
- Yopiq repo'larda bu himoya faqat pullik tashkilot tariflarida (GitHub Team / Enterprise Cloud + Secret Protection).

**Hodisalar.**
- **Uber, 2016 (FTC shikoyati):** AWS kaliti *yopiq* GitHub repo'sidagi kodda ochiq yozilgan. Hujumchilar (oʻz
  aytishicha) boshqa sizishlardagi parollar bilan GitHub'ga kirgan. Uber muhandislardan ikki bosqichli himoyani talab
  qilmagan va parolni qayta ishlatishni taqiqlamagan. ~25,6 mln ism va email, 22,1 mln telefon raqami. 2014-yilda ham
  kalit GitHub'da ochiq joylashtirilgan edi. **Saboq:** GitHub hisobingizda 2FA yoqing.
- **Toyota T-Connect, 2022 (Toyota bayonoti):** pudratchi 2017-yil dekabrda manba kodni ochiq GitHub repo'siga
  yuklagan. Kod 2022-yil 15-sentabrgacha ochiq turgan; 296 019 ta email va mijoz raqami sizgan boʻlishi mumkin.
  Kalit 17-sentabrda almashtirilgan.

## Manbalar

- GitHub Docs: [Removing sensitive data from a repository](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/removing-sensitive-data-from-a-repository) · [About push protection](https://docs.github.com/en/code-security/secret-scanning/introduction/about-push-protection) · [Push protection from the command line](https://docs.github.com/en/code-security/secret-scanning/working-with-secret-scanning-and-push-protection/working-with-push-protection-from-the-command-line) · [About secret scanning](https://docs.github.com/en/code-security/secret-scanning/introduction/about-secret-scanning) · [Secret scanning partner program](https://docs.github.com/en/code-security/secret-scanning/secret-scanning-partnership-program/secret-scanning-partner-program)
- Meli, McNiece, Reaves. [How Bad Can It Git? Characterizing Secret Leakage in Public GitHub Repositories](https://www.ndss-symposium.org/ndss-paper/how-bad-can-it-git-characterizing-secret-leakage-in-public-github-repositories/), NDSS 2019 ([PDF](https://www.ndss-symposium.org/wp-content/uploads/2019/02/ndss2019_04B-3_Meli_paper.pdf))
- Unit 42: [EleKtra-Leak — exposed IAM keys on GitHub](https://unit42.paloaltonetworks.com/malicious-operations-of-exposed-iam-keys-cryptojacking/) (2023)
- GitGuardian: [The State of Secrets Sprawl 2026](https://blog.gitguardian.com/the-state-of-secrets-sprawl-2026/)
- FTC: [Uber revised complaint (2018)](https://www.ftc.gov/system/files/documents/cases/1523054_uber_technologies_revised_complaint_0.pdf) · Toyota: [bayonot, 2022-10-07 (yapon tilida)](https://global.toyota/jp/newsroom/corporate/38095972.html)
- git: [gitignore](https://git-scm.com/docs/gitignore), [git-rm](https://git-scm.com/docs/git-rm), [git-check-ignore](https://git-scm.com/docs/git-check-ignore), [git-log `-S`](https://git-scm.com/docs/git-log), [git-filter-repo](https://github.com/newren/git-filter-repo)

## Uy vazifasi

1. 🟢 **Oson.** `./qur.sh`, keyin `loyiha` tarixidan kalitni oʻzingiz toping.
2. 🟡 **Oʻrta.** `tuzoq`da `.gitignore`ni shunday yozingki, `.env` bilan boshlangan hamma fayllar yashirinsin,
   `.env.example` esa koʻrinsin. `.env.local` va `.env.production` bilan tekshiring.
3. 🟠 **Qiyin.** Kalit qaysi commit'da paydo boʻlganini `git log -S` bilan toping.
4. 🔴 **Guru.** `loyiha`ni clone qiling, `git filter-repo` bilan `.env`ni tozalang va kalit qolmaganini tekshiring.

<details><summary>Javoblar (<code>bash homework/javoblar.sh</code>)</summary>

1. `git show HEAD~1:.env` yoki `git grep TOLOV_KALIT= $(git rev-list --all)` →
   `cdd4176…:.env:TOLOV_KALIT=soxta_9f2c7a41e8b3`.
2. `.gitignore`: `.env*` va `!.env.example`. `git status --short --ignored` → `!! .env`, `!! .env.local`,
   `!! .env.production`; `git ls-files` → `.env.example .gitignore main.go`. `git show HEAD~1:.env` esa hali ham
   kalitni chiqaradi: tarix oʻzgarmadi.
3. `git log --oneline -S TOLOV_KALIT=` satr soni oʻzgargan commit'larni beradi (`=` muhim: `-S TOLOV_KALIT` `main.go`dagi
   `os.Getenv("TOLOV_KALIT")`ni ham topib, `4590f4e`ni ham qoʻshadi; eng ishonchlisi kalit qiymati bilan qidirish): `36288b8` (oʻchirdi) va `cdd4176`
   (qoʻshdi). Birinchisi: `git log --reverse -S TOLOV_KALIT= | head` → `cdd4176 lokal sozlamalar`.
4. `git clone -q loyiha toza && cd toza && git filter-repo --sensitive-data-removal --invert-paths --path .env`.
   Tarixda faqat `4590f4e` qoladi; `git log --all -p -S soxta_9f2c7a41e8b3` boʻsh. Haqiqiy repo'da bundan keyin
   jamoa bilan kelishib `git push --force --mirror origin`, hamma qayta clone qiladi. Kalit esa baribir **rotate**
   qilinadi, bu birinchi qadam edi.

</details>
