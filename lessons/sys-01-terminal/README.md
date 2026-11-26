# SYS 0 · 1 · Terminal: kichik buyruqlar va pipe — skriptsiz katta javob

Video: *Linux terminal oʻzbek tilida: pipe, grep va wc* (2026-11-26 yangi versiya; dasturlay oladigan, lekin Linux
terminaliga yangi dasturchilar uchun) · Lab:
https://haywan.uz/lab/youtube/sys#sys-0-01 · Oldingi: [Go 0 · 5 — for](../000-05-for/) · Keyingi:
[SYS-0.02 · Jarayon va muhit: `ps`, `kill`, `$PATH`, env](../sys-02-process-env/)

**Bitta gʻoya:** har biri bitta ish qiladigan kichik buyruqlarni tik chiziq (`|`, *pipe*) bilan bosqichma-bosqich
ulasak, ular skriptsiz katta savolga javob beradi: `grep` filtrlaydi, `wc -l` sanaydi — xuddi JavaScript'dagi
`lines.filter(isTodo).filter(isAziz).length` kabi.

| Fayl / papka | Nima |
|---|---|
| `eslatmalar/` | mashq papkasi: 3 fayl (dushanba, seshanba, chorshanba), 13 qator, 6 ta `TODO` |
| `yil/` | bir yillik eslatmalar: 12 fayl × 30 qator = 360 qator, 93 ta TODO (`yasash.py` yasaydi, natija har doim bir xil; har bir TODO alohida, hali qilinmagan ish) |
| `uyga/oson.sh` | uy vazifasi 🟢 uchun boshlangʻich fayl |
| `uyga/javoblar.sh` | uy vazifasi javoblari |
| `loyiha/`, `matn.txt`, `homework/` | eski (2026-11-25) videoning materiallari, pastdagi ilovada |

## Qayerda yozamiz?
- **Linux** va **Mac**: oddiy terminal (Mac: *Terminal* ilovasi). Mac'dagi `wc -l` sondan oldin boʻsh joy qoʻyadi
  (`       6`), bu normal.
- **Windows**: [Git for Windows](https://gitforwindows.org/) oʻrnating (Git Bash shu ichida): unda `bash`, `grep`,
  `wc` bor. Windows `cmd`da `grep` yoʻq (`'grep' is not recognized…`). `sys-01-terminal` papkasida sichqonchaning oʻng
  tugmasi → **Open Git Bash here** (Windows 11'da avval *Show more options*). WSL ham boʻladi.
- Terminal `sys-01-terminal` papkasida (ichida `eslatmalar/` turgan joyda) ochilishi kerak: `No such file or
  directory` desa, boshqa papkadasiz (`pwd` qayerdaligingizni aytadi).
- **Tik chiziq `|`**: klaviaturani inglizchaga oʻtkazing, **Shift + `\`** (teskari slesh; odatda Enter yonida).
  Ruscha layoutda bu tugma boshqa belgini beradi.
- Papkani olish: `git clone https://github.com/abdurahmon27/go-backend-uz` (yoki GitHub'da **Code → Download ZIP**),
  keyin `cd go-backend-uz/lessons/sys-01-terminal`.
- VS Code: terminal profili odatda PowerShell — `Terminal → New Terminal`, `+` yonidagi strelka → **Git Bash**.
  Git Bash'da nusxa qoʻyish: Shift+Insert yoki sichqonchaning oʻng tugmasi (Ctrl+V emas).

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/sys-01-terminal`) ishga tushiriladi.
```bash
grep -r TODO eslatmalar                          # 6 qator: fayl nomi : qatorning oʻzi
grep -r TODO eslatmalar | wc -l                  # 6
grep -r TODO eslatmalar | grep Aziz              # lentada nima bor? 3 qator
grep -r TODO eslatmalar | grep Aziz | wc -l      # 3
grep -r TODO eslatmalar | wc -l | grep Aziz      # hech narsa: wc faqat "6" ni beradi
grep -r TODO eslatmalar | grep -r Aziz eslatmalar | wc -l   # XATO: 5 (ikkinchi grep lentaga qaramaydi)
grep -r TODO yil | grep Aziz | wc -l             # 22
grep -r TODO eslatmalar arxiv | wc -l            # xato ekranda, wc baribir 6 (Chuqurroq)
```
Natijalar GNU bash 5.3.9, GNU grep 3.12, GNU coreutils 9.11 (`bash -c` orqali) bilan olingan. `grep -r` papkadagi
fayllarni fayl tizimi bergan tartibda oʻqiydi, shuning uchun qatorlar tartibi sizda boshqacha boʻlishi mumkin.

## Asosiy tushunchalar
- **Ishlab chiqarish liniyasi.** Har bosqich bitta ish qiladi va natijasini keyingi lentaga qoʻyadi; birinchi bosqich
  ombordan (papkadan) oʻqiydi, keyingilari faqat oldidagi lentadan. Terminalda bosqich = buyruq, lenta = tik chiziq `|`
  (*pipe*), butun zanjir = *pipeline*.
- **JavaScript bilan farqi.** `filter().filter().length` massivni bosqichma-bosqich toʻliq yigʻadi; pipeline'da esa
  bosqichlar alohida dasturlar, bir vaqtda ishlaydi va matnni oqim bilan uzatadi. Lentada faqat matn bor, tiplar yoʻq:
  `… | wc -l | grep Aziz` TypeError bermaydi, shunchaki hech narsa chiqarmaydi.
- **`|` (pipe)** chapdagi buyruq natijasini ekranga emas, oʻngdagi buyruqqa beradi.
- **`grep SOʻZ`** qatorlarni saralaydi: ichida soʻz bor qatorni butunligicha oʻtkazadi, qolganini tashlaydi.
  `-r` (*sozlama*, inglizcha *option*) — papkadagi hamma fayllarni oʻqi. Bir nechta fayl oʻqiganda har qator boshida
  fayl nomi va `:` turadi. grep katta-kichik harfni farqlaydi: `grep -r todo eslatmalar | wc -l` → `0` (`-i` farqlamaydi).
- **`wc -l`** (*word count*, `-l` = faqat qatorlar) kelgan qatorlarni sanaydi va bitta son yozadi. Fayl berilmasa,
  buyruq *stdin*dan (standart kirish), yaʼni tik chiziqdan oʻqiydi.
- **Fayl bermasangiz, buyruq lentadan oʻqiydi.** Fayl yoki papka bersangiz, lentaga qaramaydi (videodagi xato:
  `… | grep -r Aziz eslatmalar | wc -l` → 5, toʻgʻrisi 3). Diqqat: GNU grep'da `-r` faylsiz boʻlsa, joriy papkani
  oʻqiydi, lentani emas: `… | grep -r Aziz | wc -l` ham notoʻgʻri (shu papkadagi hamma fayllardagi Aziz, `yil/` ham).
  Shuning uchun tik chiziqdan keyingi grep'ga **na papka, na `-r`** beriladi. Hech qanday kirish boʻlmasa (`grep Aziz`
  yolgʻiz), grep klaviaturani kutadi va terminal "qotib qoladi": Ctrl+C bosing (GO-0.05).
- **Tartib muhim:** avval saralaymiz, keyin sanaymiz.

## Chuqurroq (ixtiyoriy)
- **Uchta oqim** (videoda qisqa Chuqurroq qismi). Har bir dasturning kirishi (*stdin*, 0), natijasi (*stdout*, 1) va xatolar uchun alohida chiqishi
  (*stderr*, 2) bor. `|` faqat stdout'ni keyingi buyruqning stdin'iga ulaydi (POSIX: "the standard output of command1
  shall be connected to the standard input of command2"); grep xatolarni stderr'ga yozadi, shuning uchun ular ekranda
  qoladi va sanalmaydi: `grep -r TODO eslatmalar arxiv | wc -l` → ekranda `grep: arxiv: No such file or directory`
  (yaʼni bunday fayl yoki papka yoʻq), wc esa baribir `6`. Xatoni ham lentaga qoʻshish: `2>&1 |` yoki bash'da `|&`.
  (Terminal xatoni odatda rang bilan belgilamaydi; videodagi qizil rang — bizning bezak.)
- **Yoʻnaltirish.** `> fayl` natijani faylga yozadi (fayl avval boʻshatiladi: `sort f > f` faylni yoʻq qiladi),
  `< fayl` faylni kirishga beradi, `2> fayl` faqat xatolarni yozadi.
- **Bir vaqtda.** Pipe'dagi buyruqlar navbat bilan emas, bir vaqtda ishlaydi; maʼlumot yadro ichidagi kichik bufer
  orqali oʻtadi (Linux'da 65 536 bayt). `yes | head -3` darhol tugaydi: `head` yopilgach, `yes` SIGPIPE signalini oladi
  ([pipe(7)](https://man7.org/linux/man-pages/man7/pipe.7.html)).
- **Qator ≠ topilma.** `wc -l` qatorlarni sanaydi: bitta qatorda ikkita TODO boʻlsa, bitta sanaladi. Har bir topilma
  uchun `grep -o` (har topilmani alohida qatorga chiqaradi). Misol (ilovadagi `loyiha/`): `grep -rn TODO loyiha | wc -l` → 6, lekin
  `grep -ro TODO loyiha | wc -l` → 7.
- **Fayl nomlari lentada kelsa.** grep lentadagi *matnni* oʻqiydi. Lentada fayl nomlari kelsa va oʻsha fayllarning
  ichini qidirmoqchi boʻlsangiz, `xargs` nomlarni buyruqqa qoʻshib beradi: `find eslatmalar -name '*.txt' | xargs grep -l Aziz`.
  Nomda boʻsh joy yoki apostrof boʻlsa, `find … -print0 | xargs -0 …` ishlating (ilovadagi `o'quv reja.md` misoli).
- **Tarix.** McIlroy 1964-yil 11-oktabrdagi xizmat xatida dasturlarni "like garden hose" (bogʻ shlangi kabi) ulashni
  taklif qilgan; Ritchie'ning yozishicha pipe Unix'ga (bugungi Linux va macOS'ning ajdodi) 1972-yilda qoʻshilgan.

## Uy vazifasi
- 🟢 **Oson:** Malika nechta ishni qilmagan? [`uyga/oson.sh`](uyga/oson.sh) dagi qatorda bitta soʻzni oʻzgartiring
  (`bash uyga/oson.sh`).
- 🟡 **Oʻrta:** `yil` papkasida Boburning **hamma** yozuvlari (qilinganlari bilan) nechta? Ishora: birinchi ishchini
  olib tashlang va papkani `grep Bobur`ga bering (endi u birinchi ishchi). Terminal qotib qolsa, Ctrl+C.
- 🔍 **Qiziqish:** `yil` papkasida kim eng koʻp ish qoldirgan? (Aziz, Malika, Bobur, Dilnoza.)

Avval oʻzingiz urinib koʻring. Javoblar: [`uyga/javoblar.sh`](uyga/javoblar.sh) (`bash uyga/javoblar.sh`):
Oson `2`; Oʻrta `100`; Qiziqish: Bobur `25` (Dilnoza 24, Aziz 22, Malika 22).

## Manbalar
- POSIX.1-2024, Shell Command Language §2.9.2 *Pipelines*: https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html#tag_19_09_02
- bash(1), *Pipelines* va *REDIRECTION*: https://man7.org/linux/man-pages/man1/bash.1.html
- GNU grep 3.12 manual: https://www.gnu.org/software/grep/manual/grep.html · wc(1): https://man7.org/linux/man-pages/man1/wc.1.html
- M. D. McIlroy, 1964-yilgi xizmat xati: https://www.nokia.com/bell-labs/about/dennis-m-ritchie/mdmpipe.html
- D. M. Ritchie, *The Evolution of the Unix Time-sharing System*: https://www.nokia.com/bell-labs/about/dennis-m-ritchie/hist.html
- Git Bash (Atlassian): https://www.atlassian.com/git/tutorials/git-bash · Git for Windows: https://gitforwindows.org/
- Software Carpentry, *The Unix Shell: Pipes and Filters*: https://swcarpentry.github.io/shell-novice/04-pipefilter.html
- MIT *The Missing Semester*, The Shell: https://missing.csail.mit.edu/2020/course-shell/

## Dars matni (video narration)

**Savol**

Bir yillik standup eslatmalari: oʻn ikki fayl, uch yuz oltmish qator. Aziz hali nechta ishni yopmagan? Skript yozmasdan, terminalda uchta kichik buyruqni ulab topamiz.
JavaScript’da bu savolni shunday yozardingiz: qatorlar massividan TODO borini filter qilish, keyin Aziz borini, oxirida length. Bugun xuddi shu zanjirni terminalda yigʻamiz. grep bu filter, wc minus l bu length, tik chiziq esa bosqichlarni ulaydi.

**Maʼlumot**

Maʼlumot oddiy matn fayllari. Kichik mashq papkasi eslatmalar: uchta fayl, har kunga bittasi, jami oʻn uch qator. Mana dushanba. Yopilgan ish oddiy yoziladi, ochiq ish esa TODO bilan boshlanadi.
Avval kurs papkasini yuklab oling, havola pastda. Linux va Mac’da oddiy terminal yetadi. Windows’da cmd va PowerShell’da grep yoʻq, shuning uchun Git Bash yoki WSL ishlating. Terminalni dars papkasida oching: ichida eslatmalar va yil papkalari bor. Fayl topilmadi desa, boshqa papkadasiz. Tik chiziq inglizcha layoutda Shift va teskari slesh.

**Pipe**

Buni ishlab chiqarish liniyasi deb tasavvur qiling. Har bosqich bitta ish qiladi va natijasini keyingi lentaga qoʻyadi. Birinchi bosqich ombordan, yaʼni papkadan oʻqiydi, keyingilari faqat oldidagi lentadan. Terminalda bosqich bu buyruq, lenta esa tik chiziq, inglizcha uni pipe deyishadi. Butun zanjir pipeline deyiladi. Massivdan farqi: bosqichlar alohida dasturlar, ular bir vaqtda ishlaydi va matnni oqim bilan uzatadi.

**1-bosqich: grep**

Birinchi bosqich grep. U qatorlarni filtrlaydi: naqsh bor qatorni butunligicha oʻtkazadi, qolganini tashlaydi. Minus r flagi papkani rekursiv oʻqiydi, keyin naqsh TODO, oxirida papka. grep katta-kichik harfni farqlaydi: naqshni kichik harflar bilan yozsangiz, hech narsa topilmaydi.
grep bu yerda uchta faylni oʻqiydi. Natijada har bir qator qanday koʻrinadi? A: faqat qatorning oʻzi. B: fayl yoʻli, ikki nuqta, keyin qator. C: faqat fayl nomlari. Pauza qiling va oʻylab koʻring.
Javob: B. Olti qator chiqdi, dushanbadan ikkitasi. Boshida fayl yoʻli, ikki nuqtadan keyin qatorning oʻzi: bir nechta fayl oʻqilsa, grep qaysi fayldan ekanini shunday koʻrsatadi. Keyingi bosqich yoʻq, shuning uchun natija ekranga chiqdi. Fayllar tartibi sizda boshqacha boʻlishi mumkin.

**2-bosqich: wc -l**

Ikkinchi bosqich wc minus l: word count, minus l flagi bilan faqat qatorlar soni. Uni tik chiziq bilan ulaymiz. wc’ga fayl bermaymiz: fayl berilmasa, u tik chiziqdan kelganini oʻqiydi. Buyruqning bu kirishi stdin, oʻzbekcha standart kirish.
Endi ekranda nima chiqadi? A: oʻsha olti qator. B: bitta son, olti. C: bitta son, oʻn uch. Pauza qiling.
Javob: B, olti. grep’ning qatorlari ekranga emas, lentaga tushdi, wc ularni sanab faqat sonni yozdi. Oʻn uch emas, chunki wc papkani oʻqimadi. Qisqasi, wc faqat lentada kelganini sanaydi.

**Oʻrtaga yangi bosqich: grep Aziz**

Endi asosiy savol: Azizning ishlari. JavaScript’dagi ikkinchi filter kabi, oʻrtaga yana bitta grep qoʻyamiz, naqsh Aziz. Unga minus r ham, papka ham bermaymiz: u qatorlarni lentadan oladi. Uchta bosqich, ikkita tik chiziq.
Ekranda qaysi son chiqadi? A: olti. B: uch. C: besh. Oldingi olti qatorni eslang. Pauza qiling.
Javob: B, uch. Avval wc’ni olib tashlab, lentada nima borligini koʻramiz: uchta qator, hammasi Azizning TODO qatorlari. Pipeline’ni shunday tekshirish qulay: oxirgi bosqichni olib, oraliq natijaga qaraysiz. Endi wc’ni qaytaramiz: uch.
Qisqa sinov. Bosqichlar tartibini almashtirsak-chi: wc oʻrtada, grep Aziz oxirida. Ekranda nima chiqadi? A: uch. B: olti. C: hech narsa. Pauza qiling.
Javob: C, hech narsa. Oxirgi grep’ni olib koʻramiz: wc lentaga faqat bitta qator qoʻydi, olti. Oxirgi grep shu qatorda Aziz naqshini topmaydi va hech narsa chiqarmaydi. JavaScript’da length’dan keyin filter TypeError beradi, terminal esa jim turadi, chunki lentada faqat matn bor, tiplar yoʻq: olti ham oddiy qator. Tartib muhim: avval filtrlash, keyin sanash.

**Koʻp uchraydigan xato**

Endi koʻp uchraydigan xato. Ikkinchi grep’ga ham odat boʻyicha minus r va papka nomini yozsak-chi? A: baribir uch. B: uchtadan koʻp. C: nol. Pauza qiling.
Javob: B, besh. Fayl yoki papka berilgan grep lentaga qaramaydi, omborga qaytib ketadi. U papkani boshidan oʻqib, Azizning hamma yozuvlarini, yopilganlarini ham chiqaradi, wc esa ularni sanaydi. Xato xabari yoʻq, son shunchaki notoʻgʻri. Qoida: tik chiziqdan keyingi bosqichga minus r ham, papka ham bermang. Mana toʻgʻrisi: uch.

**Bir yillik papka**

Endi boshidagi savol: bir yillik papka. Pipeline’da nimani oʻzgartirish kerak? Bir soniya oʻylang. Faqat papka nomini, eslatmalar oʻrniga yil. Natija: yigirma ikki. Azizda yigirma ikkita ochiq ish bor, oʻn ikki faylni bir qator buyruq bilan koʻrib chiqdik.

**Chuqurroq: xatolar lentaga tushmaydi**

Yana bir nozik joy. Xato xabarlari alohida oqimga, stderr’ga yoziladi, xuddi Node’dagi console error kabi. Oddiy tik chiziq esa faqat asosiy natijani uzatadi. Yoʻq papkani qoʻshsak, wc nima deydi? Bir soniya oʻylang. Xato ekranda qoladi, wc esa baribir olti deydi.

**Xulosa: 3 savol**

Uchta savol bilan xulosa qilamiz. Tik chiziq nima qiladi? Tik chiziqdan keyingi grep’ga nega papka bermaymiz? wc oʻrtada boʻlsa, nega terminal xato bermaydi? Pauza qiling va javob bering.
Javoblar. Tik chiziq chapdagi buyruq natijasini oʻngdagi buyruqning kirishiga beradi. Papka yoki minus r berilsa, grep lentaga qaramaydi. wc oʻrtada boʻlsa, lentaga oddiy matn, olti, tushadi: tiplar yoʻq, xato ham yoʻq.

**Uy vazifasi**

Uy vazifasi. Oson: eslatmalar papkasida Malikaning ochiq ishlari nechta? Bitta soʻzni oʻzgartiring. Oʻrta: yil papkasida Boburning hamma yozuvlari nechta? Birinchi bosqichni olib tashlang va minus r bilan yil papkasini grep Bobur’ga bering. Qiziqish uchun: yil papkasida toʻrt kishidan kimda eng koʻp ochiq ish bor? Tayyor fayl va javoblar README’da.
Bugun kichik buyruqlarni tik chiziq bilan ulab, skriptsiz javob oldik. Keyingi darsda ishlayotgan dasturni topish, toʻxtatish va uni muhit oʻzgaruvchilari bilan sozlashni koʻramiz. Koʻrishguncha.

## Ilova: eski video (2026-11-25)

Birinchi versiya (terminalni biladiganlar uchun, zich): 1986-yilgi oltita buyruq, stdin/stdout/stderr, `find | xargs` va apostrof tuzogʻi, yoʻnaltirish, `tr` va UTF-8 baytlari. Papkalar: `loyiha/`, `matn.txt`, `homework/`. Pastdagi matn oʻsha videoniki (oʻzgartirilmagan).

Oldingi dars: [04 · `len("oʻzbek")` = 7?](../004-strings-utf8/) · Keyingi dars: [SYS-0.02 · Jarayon va muhit: `ps`, `kill`, `$PATH`, env](../sys-02-process-env/)

| Fayl | Nima |
|---|---|
| `loyiha/` | mashq uchun kichik Go loyiha (TODO'lar bilan; `docs/o'quv reja.md` nomida boʻsh joy ham, oddiy apostrof ham bor) |
| `matn.txt` | 1986-yilgi pipeline uchun oʻzbekcha matn |
| `homework/javoblar.sh` | uy vazifasi javoblari |

### Videodagi buyruqlar
Hammasi shu papkadan (`lessons/sys-01-terminal`) ishga tushiriladi.
```bash
tr -cs A-Za-z '\n' < matn.txt | tr A-Z a-z | sort | uniq -c | sort -rn | head -5   # o, zbek?
yes | head -3                                   # pipe'dagi dasturlar bir vaqtda ishlaydi
find loyiha -type f | sort
grep -rn TODO loyiha
grep -rn TODO loyiha | wc -l                    # 6: qatorlar soni
grep -ro TODO loyiha | wc -l                    # 7: topilmalar soni
grep -rn 'TODO.*TODO' loyiha                    # bitta qatorda ikkita TODO
grep -ro TODO loyiha | sort | uniq -c
grep -ro TODO loyiha | sort | uniq -c | sort -rn
find loyiha -type f -not -name '*_test.go' | xargs grep -l TODO            # xargs: unmatched single quote
find loyiha -type f -not -name '*_test.go' -print0 | xargs -0 grep -l TODO
grep -rn TODO loyiha > todo.txt
wc -l < todo.txt
grep -rn TODO loyiha yoq-papka > todo.txt                # xato baribir ekranda: u stderr'da
grep -rn TODO loyiha yoq-papka > todo.txt 2> xato.txt
cat xato.txt
grep -rn TODO loyiha yoq-papka 2>&1 | wc -l    # 7: endi xato ham pipe'ga tushadi
tr -cs 'A-Za-zʻ' '\n' < matn.txt | tr A-Z a-z | sort | uniq -c | sort -rn | head -3   # bu matnda ishlaydi…
printf 'oʻzbek «salom»\n' | tr -cs 'A-Za-zʻ' '\n' | od -c                          # …lekin » ning bayti qoladi
```
Natijalar GNU bash 5.3, GNU grep 3.12, GNU findutils 4.10, GNU coreutils 9.11 (`bash -c` orqali) bilan olingan.
macOS'dagi BSD versiyalarda ham `-print0`/`-0` bor (POSIX.1-2024 da standart), `grep --include` esa eski BSD
grep'da boʻlmasligi mumkin.

### Dars matni

**1986.** Bitta qator, oltita buyruq: matndagi eng koʻp uchraydigan soʻzlarni sanaydi. Jon Bentley Donald Knuth'dan
shu vazifani bajaradigan dastur soʻragan edi. Knuth oʻz uslubini (*literate programming*) koʻrsatib, jurnalning bir
necha sahifasini egallagan puxta dastur yozdi. Doug McIlroy taqrizida bu taqdimotni ishonarli deb topdi ("convincing
as a demonstration of WEB"), lekin muhandislik nuqtai nazaridan rozi boʻlmadi ("an industrial-strength Fabergé egg")
va oʻrniga oltita buyruq yozdi (videoda oxirgi `sed ${1}q` oʻrniga `head -5`):
`tr -cs A-Za-z '\n' | tr A-Z a-z | sort | uniq -c | sort -rn | sed ${1}q`. U "joyida yozilgan va birinchi
urinishdayoq ishlagan".

**Pipe.** Shell har bir dasturga uchta stream beradi: stdin (0) kirish, stdout (1) natija, stderr (2) xatolar uchun.
`|` chapdagi dasturning stdout'ini oʻngdagining stdin'iga ulaydi; stderr pipe'dan oʻtmaydi. Dasturlar **bir vaqtda**
ishlaydi, maʼlumot yadro buferi orqali oʻtadi (vaqtinchalik fayl yoʻq): `yes | head -3` darhol tugaydi, chunki head
yopilgach `yes` SIGPIPE oladi. Gʻoya McIlroy'niki (1964-yil 11-oktabrdagi xizmat xati: dasturlarni "bogʻ shlangi"
kabi ulash), Dennis Ritchie'ning yozishicha pipe Unix'ga 1972-yilda qoʻshilgan. 1978-yilda McIlroy va hamkasblari
Unix quruvchilari orasida tarqalgan qoidalarni yozib qoʻyishgan: har bir dastur bitta ishni yaxshi qilsin; har bir
dasturning chiqishi boshqa, hali nomaʼlum dasturga kirish boʻlishini hisobga oling. (Bugungi `grep -r` bu qoidadan
biroz chekinadi: u papkalarni oʻzi aylanadi.)

**grep.** `grep -rn TODO loyiha`: `-r` papkalarni ichma-ich kezib chiqadi, `-n` qator raqamini qoʻshadi.
`| wc -l` qatorlarni sanaydi: 6. Lekin TODO'lar 7 ta: `savat.go` ning 15-qatorida ikkita TODO bor.
`wc -l` ham, `grep -c` ham **qatorlarni** sanaydi; har bir topilmani sanash uchun `grep -o` kerak.
(`-o TODO` "TODOs" ichidagi TODO ni ham topadi; faqat butun soʻz uchun `-w`.)

**Eng koʻp TODO qayerda?** `grep -ro TODO loyiha | sort | uniq -c | sort -rn`: `sort` bir xil qatorlarni yonma-yon
qoʻyadi, chunki `uniq` faqat **yonma-yon** turgan takrorlarni birlashtiradi (bu yerda grep natijasi allaqachon fayl
boʻyicha guruhlangan, lekin odatda `sort` shart); `uniq -c` sanaydi; `sort -rn` son
boʻyicha kamayish tartibida. Bu McIlroy pipeline'ining yuragi.

**find, xargs va apostrof.** Test fayllaridan tashqari fayllarda qidirish uchun `find -not -name '*_test.go'` qulay
(find nom, hajm, vaqt boʻyicha filtrlaydi). `xargs` esa odatda argumentlarni boʻsh joy va yangi qatordan boʻladi,
apostrof, qoʻshtirnoq va `\` ni maxsus belgi deb biladi. `o'quv reja.md` dagi apostrofning jufti yoʻq:
`xargs: unmatched single quote`, chiqish kodi 1, xargs shu joyda toʻxtaydi: bu fayl **va undan keyingi hamma fayllar**
natijadan tushib qoladi (qaysilari, `find` tartibiga bogʻliq). Yechim: `find -print0` nomlarni nol bayt
bilan ajratadi (nol bayt fayl yoʻlida boʻla olmaydi), `xargs -0` aynan shunday oʻqiydi. Muqobil:
`find … -exec grep -l TODO {} +`. Xatoni `2>/dev/null` bilan yashirsangiz, natija jimgina notoʻgʻri boʻladi.

**Yoʻnaltirish.** `>` stdout'ni faylga yozadi (fayl avval boʻshatiladi, shuning uchun `sort f > f` faylni yoʻq qiladi),
`<` faylni stdin'ga beradi. Yoʻq papka qoʻshilsa, xato baribir ekranda, chunki u stderr orqali chiqadi. `2> xato.txt`
(boʻsh joysiz `2>`) uni alohida faylga yuboradi; `2>&1` stderr'ni stdout ketayotgan joyga yuboradi.

**Nega `o` va `zbek`?** Birinchi `tr -cs A-Za-z '\n'` A–Z/a–z boʻlmagan hamma narsani yangi qatorga aylantiradi,
`-s` ketma-ket yangi qatorlarni bittaga siqadi. Apostrofning hech bir turi `A-Za-z` da yoʻq, shuning uchun `oʻzbek`
ikkiga boʻlinadi (inglizcha `don't` ham `don` + `t` boʻladi). McIlroy oʻsha taqrizda Knuth'ning dasturi 26 harfga
bogʻlanib qolganini tanqid qilgan: "Just how might one confidently change it to handle, say, Norwegian…?" Uning oʻz
`A-Za-z` qatorida ham xuddi shu cheklov bor. `ʻ` ni roʻyxatga qoʻshsak, bu matnda ishlaydi, lekin tasodifan: GNU `tr`
faqat bir baytli locale'larni toʻliq qoʻllaydi va `ʻ` ni ikki alohida bayt (`ca`, `bb`) deb oladi; `»` (`c2 bb`) ning
`bb` bayti ham (kirill `л` = `d0 bb` ham) oʻtib ketadi va matn buziladi. Ishonchli yechim: uy vazifasi (Qiyin).

### Manbalar
- Bentley, Knuth, McIlroy, *Programming Pearls: A Literate Program*, CACM 29(6), 1986, pp. 471–483,
  [doi:10.1145/5948.315654](https://doi.org/10.1145/5948.315654)
- D. M. Ritchie, [The Evolution of the Unix Time-sharing System](https://www.nokia.com/bell-labs/about/dennis-m-ritchie/hist.html);
  [McIlroy'ning 1964-yilgi xizmat xati](https://www.nokia.com/bell-labs/about/dennis-m-ritchie/mdmpipe.html)
- M. D. McIlroy, [A Research Unix Reader](https://mcilroy.cs.dartmouth.edu/reader.pdf)
- McIlroy, Pinson, Tague, *UNIX Time-Sharing System: Foreword*, BSTJ 57(6), 1978, [archive.org](https://archive.org/details/bstj57-6-1899)
- man sahifalari: [pipe(7)](https://man7.org/linux/man-pages/man7/pipe.7.html), [grep(1)](https://man7.org/linux/man-pages/man1/grep.1.html),
  [find(1)](https://man7.org/linux/man-pages/man1/find.1.html), [xargs(1)](https://man7.org/linux/man-pages/man1/xargs.1.html),
  [uniq(1)](https://man7.org/linux/man-pages/man1/uniq.1.html), [wc(1)](https://man7.org/linux/man-pages/man1/wc.1.html),
  [tr(1)](https://man7.org/linux/man-pages/man1/tr.1.html), [bash(1) REDIRECTION](https://man7.org/linux/man-pages/man1/bash.1.html)
- [POSIX.1-2024 xargs](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/xargs.html)
- MIT [The Missing Semester](https://missing.csail.mit.edu/): shell va data wrangling darslari

### Uy vazifasi
- 🟢 **Oson:** `loyiha` dagi barcha `.go` fayllarni sanang.
- 🟡 **Oʻrta:** loyihaning qaysi papkasida (`cmd`, `docs`, `internal`) eng koʻp TODO bor? (Ishora: `cut`.)
- 🟠 **Qiyin:** `matn.txt` dagi eng koʻp uchraydigan 3 soʻzni `oʻ` harfini buzmasdan toping (apostrofning hamma turi bilan).
- 🔴 **Guru:** `cmd > f 2>&1` va `cmd 2>&1 > f` farqi nima va nega?

Avval oʻzingiz urinib koʻring. Javoblar: [homework/javoblar.sh](homework/javoblar.sh) (`bash homework/javoblar.sh`):
Oson `4`; Oʻrta `5 internal, 1 docs, 1 cmd`; Qiyin `3 va, 3 oʻzbek, 2 uchun`; Guru: bash yoʻnaltirishlarni chapdan
oʻngga bajaradi. `> f 2>&1`: avval stdout → f, keyin stderr → stdout hozir qayerda boʻlsa (f): ikkalasi faylda.
`2>&1 > f`: avval stderr → stdout hozir qayerda boʻlsa (ekran), keyin faqat stdout → f: xato ekranda qoladi.
