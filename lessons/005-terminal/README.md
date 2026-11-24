# 05 · Terminal: pipe, `grep` va bitta qatorda yechiladigan ishlar

Oldingi dars: [04 · `len("oʻzbek")` = 7?](../004-strings-utf8/) · Keyingi dars: SYS-0.02 · Jarayon va muhit: `ps`, `kill`, `$PATH`, env

| Fayl | Nima |
|---|---|
| `loyiha/` | mashq uchun kichik Go loyiha (TODO'lar bilan; `docs/o'quv reja.md` nomida boʻsh joy ham, oddiy apostrof ham bor) |
| `matn.txt` | 1986-yilgi pipeline uchun oʻzbekcha matn |
| `homework/javoblar.sh` | uy vazifasi javoblari |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/005-terminal`) ishga tushiriladi.
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

## Dars matni

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

## Manbalar
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

## Uy vazifasi
- 🟢 **Oson:** `loyiha` dagi barcha `.go` fayllarni sanang.
- 🟡 **Oʻrta:** loyihaning qaysi papkasida (`cmd`, `docs`, `internal`) eng koʻp TODO bor? (Ishora: `cut`.)
- 🟠 **Qiyin:** `matn.txt` dagi eng koʻp uchraydigan 3 soʻzni `oʻ` harfini buzmasdan toping (apostrofning hamma turi bilan).
- 🔴 **Guru:** `cmd > f 2>&1` va `cmd 2>&1 > f` farqi nima va nega?

Avval oʻzingiz urinib koʻring. Javoblar: [homework/javoblar.sh](homework/javoblar.sh) (`bash homework/javoblar.sh`):
Oson `4`; Oʻrta `5 internal, 1 docs, 1 cmd`; Qiyin `3 va, 3 oʻzbek, 2 uchun`; Guru: bash yoʻnaltirishlarni chapdan
oʻngga bajaradi. `> f 2>&1`: avval stdout → f, keyin stderr → stdout hozir qayerda boʻlsa (f): ikkalasi faylda.
`2>&1 > f`: avval stderr → stdout hozir qayerda boʻlsa (ekran), keyin faqat stdout → f: xato ekranda qoladi.
