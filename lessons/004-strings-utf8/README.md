# 04 · `len("oʻzbek")` = 7? Harf, rune va bayt

Video: *Golang oʻzbek tilida #4* · Lab: https://haywan.uz/lab/youtube/go#go-1-04 ·
Oldingi dars: [03 · Konstantalar va `iota`](../003-const-iota/) · Keyingi dars: [05 · Go'da `while` yoʻq](../005-for-switch/) ·
Avval koʻring: Go 0 darslari ([000-01](../000-01-birinchi-dastur/) … [000-05](../000-05-for/)), [01](../001-why-go/) (bayt), [02](../002-types-overflow/) (int8).

Go 0 dan keyin koʻriladigan dars (2026-11-26 da qayta yozildi). Hamma natijalar `go1.26.8 linux/amd64`da olingan.
Har bir dasturni https://go.dev/play da ham ishga tushirish mumkin: faylni nusxalab qoʻying va **Run** bosing.

**Bitta gʻoya:** kompyuter har bir belgini raqam qilib baytlarda saqlaydi: oddiy klaviatura belgisi bitta bayt,
oʻ dagi teskari vergul (ʻ) esa ikki bayt oladi. `len` harflarni emas, baytlarni sanaydi, shuning uchun koʻzga bir xil
koʻringan soʻzlar ichida har xil baytlar boʻlishi va Go ularni teng emas deyishi mumkin.

| Soʻz | Maʼnosi |
|---|---|
| harf | alifbo harfi: `oʻ` bitta harf |
| belgi | ekrandagi bitta belgi: `oʻ` ekranda ikki belgi (`o` + `ʻ`) |
| Unicode | dunyodagi hamma belgilarning raqamlangan roʻyxati (Unicode 18.0: 172 808 ta belgi) |
| rune | belgining Unicode raqami (`ʻ` → 699) |
| UTF-8 | raqamni baytlarga yozish qoidasi: 1–4 bayt |
| `len(s)` | matndagi **baytlar** soni |
| `utf8.RuneCountInString(s)` | matndagi rune'lar soni |

`oʻzbek` = **5 harf · 6 rune · 7 bayt**.

| Papka | Nima | Natija |
|---|---|---|
| `uzunlik/` | `len` va `RuneCountInString` | `7` / `6` |
| `sinov/` | oʻrtadagi sinov: `len("gʻoz")` | `5` |
| `tort/` | bitta soʻz, toʻrt xil apostrof | `len` 6 · 7 · 8 · 8, rune hammasida 6 |
| `teng/` | xato: `"Oʻzbekiston" == "O'zbekiston"` | `false` |
| `tuzat/` | `strings.ReplaceAll` bilan tuzatish va uning cheklovi | `true` / `maʻno` |
| `lexuz/` | lex.uz sahifasida oʻ/gʻ dan keyingi belgilarni sanash (tarmoq kerak) | `U+02BB 50 · U+2018 13 · U+2019 1` |
| `bitlar/` | Chuqurroq: `ʻ` qanday qilib `ca bb` boʻladi | `ca bb` (699 = 10 bit) |
| `kesish/` | Chuqurroq: `s[:2]` belgini yarmidan kesadi | `"o\xca"` / `false` / `"oʻ"` / `"o"` |
| `uyga/` | uy vazifasi: boshlangʻich fayllar va `javob/` | |

Qolgan papkalar (`uch/`, `aylanish/`, `bayt/`, `qidiruv/`, `tozala/`, `yangi/`, `kirill/`, `homework/`) darsning
birinchi versiyasidan qoldi, ular pastdagi **Ilova**da (`tozala/` — toʻliqroq tozalash funksiyasi, videoda tilga olinadi).

## Videodagi buyruqlar
```bash
go run uzunlik/main.go   # 7 / 6
go run sinov/main.go     # 5
go run tort/main.go      # o'zbek len=6 rune=6 · oʻzbek len=7 · o’zbek len=8 · o‘zbek len=8
go run teng/main.go      # false
go run tuzat/main.go     # true / maʻno
go run bitlar/main.go    # U+02BB = 699 = 1010111011 (10 bit) / 110|01010  10|111011  = ca bb
go run kesish/main.go    # "o\xca" / false / "oʻ" / "o"
curl -sL https://lex.uz/uz/docs/8548754 | go run ./lexuz   # 2026-11-26: U+02BB 50, U+2018 13, U+2019 1
```

## Dars matni (qisqa)

**Hook.** Oʻzbek soʻzida nechta harf bor? Besh: oʻ, z, b, e va k. Go'dan shu soʻzning uzunligini soʻrasangiz, yetti deydi.

**Eslash.** Bitta baytga nima sigʻadi? *Bitta oddiy harf* (01-dars). Bayt sakkizta nol yoki bir; 02-darsdagi `int8`
ham bitta bayt, unga 256 xil son sigʻadi. Kompyuter harfni ham son qilib saqlaydi: har bir belgining raqami bor (`z` → 122).

**Unicode va UTF-8.** Dunyoda belgilar juda koʻp (lotin, kirill, xitoy yozuvi, emoji): hammasi bitta katta roʻyxatda
raqamlangan — **Unicode** (170 mingdan ortiq belgi), bitta baytga esa faqat 256 xil son sigʻadi. Shuning uchun katta
raqamli belgi bir nechta baytga yoziladi, xuddi avtoturargohda avtobus ikki-uch joy egallagani kabi. Oddiy ingliz
klaviaturasidagi belgilar bittadan bayt, `ʻ` (699) ikki bayt, `’` (8217) va `‘` (8216) uch baytdan. Belgini baytlarga
yozish qoidasi — **UTF-8**.

**len.** `oʻzbek`: o 1 + ʻ 2 + z b e k 4 = **7 bayt**. `len` (inglizcha *length*, uzunlik) harflarni emas, baytlarni
sanaydi. Ekranda esa oltita belgi: o, ʻ, z, b, e, k. Har bir belgining Unicode raqami Go'da **rune** deyiladi.

```go
s := "oʻzbek"
fmt.Println(len(s))                    // 7 — baytlar
fmt.Println(utf8.RuneCountInString(s)) // 6 — rune'lar (fayl boshida: import "unicode/utf8")
```

**Toʻrt xil apostrof.** Teskari vergulni oddiy klaviaturada topish qiyin, shuning uchun oʻzbek soʻzi toʻrt xil yoziladi.
Toʻgʻrisi `ʻ` U+02BB: Unicode'ning tillar maʼlumotlari (CLDR `uz.xml`) oʻ va gʻ uchun shuni koʻrsatadi.

```
o'zbek  len=6  rune=6     ' U+0027 klaviatura apostrofi
oʻzbek  len=7  rune=6     ʻ U+02BB teskari vergul (MODIFIER LETTER TURNED COMMA)
o’zbek  len=8  rune=6     ’ U+2019 oʻng egri apostrof (RIGHT SINGLE QUOTATION MARK)
o‘zbek  len=8  rune=6     ‘ U+2018 chap egri apostrof (LEFT SINGLE QUOTATION MARK)
```

**Hayotda.** lex.uz'dagi bitta qonun sahifasida (https://lex.uz/uz/docs/8548754, 2026-11-26) koʻrinadigan matnda
oʻ/gʻ dan keyin: `ʻ` 50 marta, `‘` 13, `’` 1; hatto bitta qatorda ikki xil: "Oʻzbekiston Respublikasining Qonuni,
21.10.2026 yildagi O‘RQ-1184-son". Hisoblash usuli: `lexuz/` `<head>`, `<script>`, `<style>` va teglarni olib tashlaydi
(menyular va sarlavhalar sanaladi). Qonun matnining oʻzi faqat `ʻ` ishlatadi.

**Koʻp uchraydigan xato.** Saytdagi roʻyxatda `Oʻzbekiston` (ʻ, 12 bayt), foydalanuvchi qidiruvga `O'zbekiston`
(', 11 bayt) yozdi. `saytda == qidiruv` → **false**: Go matnlarni baytlari boʻyicha solishtiradi
(Go qoidasi: https://go.dev/ref/spec#Comparison_operators). Qidiruv topa olmaydi.

**Tuzatish.**

```go
qidiruv = strings.ReplaceAll(qidiruv, "'", "ʻ") // har bir ' ni ʻ ga almashtir
fmt.Println(saytda == qidiruv)                  // true
fmt.Println(strings.ReplaceAll("ma'no", "'", "ʻ")) // maʻno — notoʻgʻri!
```

Cheklov: *maʼno* dagi belgi — **tutuq belgisi** `ʼ` (U+02BC), uni `ʻ` ga aylantirish xato. Shuning uchun tozalashni bir
joyda qiling (saqlashdan oldin ham, qidirishdan oldin ham). Toʻliqroq qoida `tozala/` da: o/g dan keyin → `ʻ`, boshqa
harflar orasida → `ʼ`; uning ham chegaralari bor (`Go'da` → `Goʻda`), pastdagi Ilovaga qarang.

**Chuqurroq (ixtiyoriy).** `ʻ` = 699 = `1010111011` (10 bit). UTF-8 ning ikki baytli qolipi `110xxxxx 10xxxxxx`:
`110|01010 10|111011` = `ca bb`. `110…` belgi boshini, `10…` davomini bildiradi, shuning uchun dastur istalgan baytga
qarab belgining boshimi yoki davomimi biladi (RFC 3629 §3). Baytlab kesish: `s[:2]` = `"o\xca"`, `utf8.ValidString` →
`false`; `string([]rune(s)[:2])` = `"oʻ"` (lekin `[:1]` = `"o"`: rune ham harf emas).

Videoga sigʻmaganlar (Ilovada batafsil): har bir baytni `%c` bilan chiqarish → `oÊ»zbek` (kodirovkasi adashgan
saytlardagi "gʻalati belgilar" ham shunday paydo boʻladi; `bayt/`); `for range` string ustida rune beradi, indeks esa
bayt oʻrni (`aylanish/`); saralash bayt boʻyicha, shuning uchun `oʻrik` olma bilan pomidor orasiga tushadi (`qidiruv/`);
2027-yil 22-apreldan kuchga kiradigan yangi alifbo qonuni (OʻRQ-1184, Ö, Ğ) va `ö` ning ikki xil yozilishi (`yangi/`);
lotin `A` va kirill `А` (`kirill/`).

**Takror (savollar).** 1) `len("oʻzbek")` nega 7? — `len` baytlarni sanaydi, `ʻ` ikki bayt. 2) Rune'larni nima
sanaydi? — `utf8.RuneCountInString`. 3) (01-darsdan) `go build` nima yasaydi? — kodni tarjima qilib, tayyor faylga
(binary) saqlaydi.

## Uy vazifasi

- 🟢 **Oson** — [`uyga/oson/main.go`](uyga/oson/main.go): `ism` ga oʻz ismingizni yozing. Ishga tushirishdan **oldin**
  ikkala sonni (bayt va rune) qogʻozga yozing. Keyin `ʻ` ni klaviaturadagi `'` ga almashtiring: qaysi son oʻzgaradi?
- 🟡 **Oʻrta** — [`uyga/orta/main.go`](uyga/orta/main.go): sayt ismni `len` bilan "koʻpi bilan 5" deb tekshiradi.
  `Gʻani` ekranda 5 belgi, lekin dastur `juda uzun` deydi. Nega? Tuzating.

### Javoblar
- 🟢 [`uyga/javob/oson`](uyga/javob/oson/main.go): `Gʻayrat` → `8` bayt, `7` rune; `G'ayrat` → `7`, `7`. Faqat bayt soni
  oʻzgaradi, chunki `ʻ` ikki bayt, `'` esa bitta.
- 🟡 [`uyga/javob/orta`](uyga/javob/orta/main.go): `len("Gʻani")` = 6 (`ʻ` ikki bayt). Belgilarni sanash uchun
  `utf8.RuneCountInString(ism) <= 5` (va `import "unicode/utf8"`) → `qabul`.

Qoʻshimcha mashqlar (birinchi versiyadan, qiyinroq): [`homework/`](homework/).

## Manbalar
- Go qoidalari: String types (`len` = baytlar) https://go.dev/ref/spec#String_types · Rune literals (rune = Unicode
  raqami) https://go.dev/ref/spec#Rune_literals · Source code (UTF-8) https://go.dev/ref/spec#Source_code_representation ·
  Comparison operators (string baytma-bayt) https://go.dev/ref/spec#Comparison_operators
- `unicode/utf8` https://pkg.go.dev/unicode/utf8 · `strings.ReplaceAll` https://pkg.go.dev/strings#ReplaceAll
- RFC 3629 (UTF-8) https://www.rfc-editor.org/rfc/rfc3629
- Unicode 18.0.0 (172 808 ta belgi) https://www.unicode.org/versions/latest/
- CLDR `uz.xml` (`{oʻ} {gʻ}` = U+02BB, `ʼ` = U+02BC) https://github.com/unicode-org/cldr/blob/main/common/main/uz.xml
- Rob Pike, "Strings, bytes, runes and characters in Go" https://go.dev/blog/strings
- lex.uz, OʻRQ-1184 https://lex.uz/uz/docs/8548754

---

# Ilova: darsning birinchi versiyasi (2026-11-25)

Quyidagi matn va papkalar birinchi (2026-11-25) videoga tegishli. Hammasi hali ham toʻgʻri, faqat zichroq —
Go'ni biladiganlar uchun. Eslatma: lex.uz soni oʻshanda butun HTML (koʻrinmaydigan `<head>` bilan) boʻyicha sanalgan
(60/16/1/2); hozirgi `lexuz/` faqat koʻrinadigan matnni sanaydi (50/13/1).

## 04 · `len("oʻzbek")` = 7? Uch xil apostrof, uch xil uzunlik

Oldingi dars: [03 · Konstantalar va `iota`](../003-const-iota/) · Keyingi dars: 05 · Go'da `while` yoʻq: bitta `for`

| Papka | Nima |
|---|---|
| `uch/` | bitta soʻz, toʻrt xil apostrof: `len` va rune soni |
| `bitlar/` | `ʻ` qanday qilib `ca bb` boʻladi (bit darajasida) |
| `aylanish/` | `for range` string ustida: indeks = bayt oʻrni |
| `bayt/` | pauza-savol: baytma-bayt `%c` bilan chiqarish |
| `kesish/` | `s[:2]` belgini yarmidan kesadi; `[]rune` bilan kesish va uning chegarasi |
| `qidiruv/` | qidiruv, soʻzlarga boʻlish va saralash qayerda sinadi |
| `tozala/` | apostroflarni bir xil qiladigan funksiya |
| `yangi/` | 2027-yilgi alifbo: `ö` ikki xil yozilishi mumkin |
| `kirill/` | lotin `A` va kirill `А`: koʻzga bir xil, baytlari boshqa |
| `lexuz/` | sahifadagi apostrof turlarini sanash: `curl -sL https://lex.uz/uz/docs/8548754 \| go run ./lexuz` (videoda ishlatilmaydi) |
| `homework/` | uy vazifasi javoblari |

### Videodagi buyruqlar
```bash
go run ./uch
go run ./bitlar
go run ./aylanish
go run ./bayt        # pauza-savol
go run ./kesish
go run ./qidiruv
go run ./tozala
go run ./yangi
go run ./kirill
```
Natijalar `go1.26.8 linux/amd64` da olingan.

### Dars matni

**Uzunlik.** "Oʻzbek" soʻzida besh harf bor (oʻ, z, b, e, k), rune'lar olti, Go'da `len` esa yetti deydi.
Apostrofni boshqasiga almashtirsangiz, `len` olti yoki sakkiz boʻladi:

```
o'zbek  len=6  rune=6  6f 27 7a 62 65 6b          U+0027 klaviatura apostrofi
oʻzbek  len=7  rune=6  6f ca bb 7a 62 65 6b       U+02BB MODIFIER LETTER TURNED COMMA
o’zbek  len=8  rune=6  6f e2 80 99 7a 62 65 6b    U+2019 RIGHT SINGLE QUOTATION MARK
o‘zbek  len=8  rune=6  6f e2 80 98 7a 62 65 6b    U+2018 LEFT SINGLE QUOTATION MARK
```

Rune soni hammasida olti, harflar esa besh. Demak `len` harflarni ham, rune'larni ham emas, **baytlarni** sanaydi. Go spetsifikatsiyasi aynan shunday
deydi: string, bu baytlar ketma-ketligi, uning uzunligi esa baytlar soni. String istalgan baytni saqlay oladi, lekin
Go kodining oʻzi UTF-8'da yoziladi, shuning uchun qoʻshtirnoq ichidagi matn UTF-8 baytlariga aylanadi. `byte` bu `uint8`, `rune` esa `int32`.

**UTF-8.** 1992-yilda Ken Thompson Rob Pike bilan birga, Nyu-Jersidagi kafeda dasturxon qogʻozi ustida oʻylab topgan (ikkalasi keyinchalik Go'ning uch muallifidan
ikkitasi boʻldi). Belgining Unicode raqami qanchalik katta boʻlsa, shuncha koʻp bayt kerak:

| Oraliq | Bayt | Misol |
|---|---|---|
| U+0000 – U+007F | 1 | ASCII, `'` U+0027 |
| U+0080 – U+07FF | 2 | `ʻ` U+02BB, `ʼ` U+02BC |
| U+0800 – U+FFFF | 3 | `’` U+2019, `‘` U+2018 |
| U+10000 – U+10FFFF | 4 | koʻp emojilar, masalan 😀 U+1F600 (lekin ✅ U+2705 — 3 bayt) |

**Bitlar.** `ʻ` = U+02BB = 699 = `1010111011` (10 bit). Ikki baytli qolip `110xxxxx 10xxxxxx` da 11 ta boʻsh joy bor,
shuning uchun oldiga bitta nol qoʻshamiz: `110|01010 10|111011` = `ca bb`. `110…` belgi boshini, `10…` davomini bildiradi, shuning uchun istalgan baytdan
belgi boshini topish mumkin.

**Qaysi biri toʻgʻri?** Unicode'ning tillar haqidagi maʼlumotlar bazasi (CLDR) `oʻ` va `gʻ` uchun **U+02BB** ni, tutuq belgisi
(*maʼno*) uchun **U+02BC** ni ishlatadi. Unicode ikkalasini ham harf (`Lm`) deb hisoblaydi. `‘ ’` esa Unicode'da
qoʻshtirnoq deb nomlangan (inglizchada apostrof uchun `’` tavsiya qilinadi, lekin oʻzbekchada `oʻ` ichidagi belgi
harfning qismi). Oddiy klaviaturada bularning hech biri yoʻq, Word esa odatda `'` ni avtomatik `’` ga almashtiradi.
Amalda: yangi alifbo qonunining lex.uz sahifasida (2026-11-25, menyu va sarlavhalar bilan) o/g dan keyin
U+02BB 60, U+2018 16, U+0027 1, U+2019 2 marta uchraydi (`<script>` bloklari hisobga olinmagan). Qonun matnining
oʻzida esa faqat U+02BB, lekin bitta qatorda ikkalasi ham bor: "Oʻzbekiston Respublikasining Qonuni, … O‘RQ-1184-son".
`‘` (U+2018) `ʻ` ga oʻxshaydi, shuning uchun koʻp joyda uning oʻrnida ishlatiladi; Word esa `’` qoʻyadi.
Inglizchada apostrof uchun Unicode faqat `’` (U+2019) ni tavsiya qiladi.

**Rune.** `for i, r := range s` matnni rune'lar boʻyicha oʻqiydi: `i` rune boshlangan baytning oʻrni, `r` rune'ning
oʻzi. Indekslar `0 1 3 4 5 6`: ikki yoʻq, chunki `ʻ` ikki bayt egalladi. `rune` bu `int32`: `ʻ` = 699, `z` = 122.

**Pauza-savol.** `for i := 0; i < len(s); i++ { fmt.Printf("%c", s[i]) }` nima chiqaradi? Javob: `oÊ»zbek`
(Windows-1251 kirill kodirovkasida xuddi shu baytlar `oК»zbek` boʻlib koʻrinadi).
`s[i]` bitta bayt, `%c` esa shu raqamdagi Unicode belgini chiqaradi: `0xCA` = 202 = `Ê`, `0xBB` = `»`. Kodlashi
buzilgan saytlarda xuddi shunday belgilarni koʻrgan boʻlishingiz mumkin: UTF-8 baytlari Latin-1 deb oʻqilganda shunday chiqadi.

**Kesish.** `s[:2]` = `"o\xca"`: belgining yarmi, `utf8.ValidString` false. Rune'lar boʻyicha kesish UTF-8'ni
buzmaydi: `string([]rune(s)[:2])` = `"oʻ"`. Lekin harfni buzishi mumkin: `string([]rune(s)[:1])` = `"o"`, chunki
`oʻ` ikki rune. Yana: `[]rune(s)` yangi slice ajratadi (katta matnda qimmat) va notoʻgʻri baytlarni
jimgina U+FFFD ga aylantiradi. (`strings.Index` qaytargan indekslar bayt indekslari, ular bilan kesish xavfsiz.)

**Qidiruv, boʻlish, saralash.**
- Bazada `Oʻzbekiston` (U+02BB), foydalanuvchi `O'zbekiston` (U+0027) deb qidiradi: `strings.Contains` → `false`.
- Harf boʻlmagan joydan boʻlsak (`unicode.IsLetter`): `[oʻzbek o zbek o zbek]`. U+02BB harf, `'` va `’` esa tinish
  belgisi, shuning uchun soʻz ikki boʻlakka boʻlinadi.
- `slices.Sort` baytlarni solishtiradi: `[gilos gʻoz olma oʻrik pomidor zarang]`. Hozirgi alifboda esa `oʻ gʻ sh ch`
  `z` dan keyin keladi (CLDR uz collation): toʻgʻri tartib `gilos olma pomidor zarang oʻrik gʻoz`.
  `golang.org/x/text/collate` (v0.42.0, CLDR 23 jadvallari) ham `language.Uzbek` bilan bu tartibni bermaydi: biz tekshirdik.

**Yechim: `tozala`.** Matnni bazaga yozishdan oldin ham, qidirishdan oldin ham bir xil koʻrinishga keltiring.
Biz bilgan apostrof turlari (`` ' ` ´ ‘ ’ ʻ ʼ ``) faqat **harfdan keyin** turgan boʻlsa oʻzgartiriladi: `o`/`g` dan
keyin (soʻz oxirida ham: `tog'`) → `ʻ` (U+02BB); boshqa ikki harf orasida → `ʼ` (U+02BC). Soʻz boshidagi
qoʻshtirnoqqa tegilmaydi. Natija: `Oʻzbekiston, togʻ, maʼno, Oʻzbek`, `U ‘Salom’ dedi` oʻzgarmaydi, qidiruv `true`.
Bu qoida mukammal emas: `Go'da` → `Goʻda`, `Rio’ga` → `Rioʻga` (qoʻshimchani ajratuvchi apostrof),
`‘Bingo’ oʻyini` → `‘Bingoʻ oʻyini` (o/g bilan tugagan soʻzdan keyingi yopuvchi qoʻshtirnoq). Tozalashni bitta joyda qiling va test bilan himoyalang (testlar: 07-dars).

**2027.** OʻRQ-1184 qonuni (2026-yil 21-oktabr, kuchga kirishi 2027-yil 22-aprel) bilan alifbo 28 harf + tutuq
belgisi `ʼ` boʻladi (`Ö`, `Ğ`, `Ş`, `Ç` bor, `oʻ`/`gʻ` yoʻq). Eski shakldagi hujjatlar amal qiladi (2-modda), oʻtish
dasturini Vazirlar Mahkamasi tasdiqlaydi (3-modda). Muammo yoʻqolmaydi, shakli oʻzgaradi: `ö` bitta belgi (U+00F6, 2 bayt) yoki `o` +
U+0308 (3 bayt) boʻlishi mumkin. Ekranda bir xil, `==` esa `false`. Yechim: Unicode normalization (NFC),
Go'da `golang.org/x/text/unicode/norm`.

**Kirill `А`.** Shu qonunning lex.uz'dagi alifbo jadvalida birinchi harf `А а` kirill U+0410/U+0430 bilan yozilgan
(`pе` da ham kirill `е`). Koʻzga bir xil, lekin 2 bayt, `==` false. Normalization buni tuzatmaydi (U+0410 ning
decomposition'i yoʻq); Unicode bunday juftlarni *confusable* deydi (UTS #39). Koʻzingizga emas, baytlarga ishoning.

### Manbalar
- Go spec: [String types](https://go.dev/ref/spec#String_types), [For range](https://go.dev/ref/spec#For_range),
  [Source code representation](https://go.dev/ref/spec#Source_code_representation),
  [Comparison operators](https://go.dev/ref/spec#Comparison_operators)
- Rob Pike, [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings) (Go blog, 2013)
- [Text normalization in Go](https://go.dev/blog/normalization) (Go blog)
- [RFC 3629: UTF-8](https://www.rfc-editor.org/rfc/rfc3629) (2003)
- Unicode CLDR: [`common/main/uz.xml`](https://github.com/unicode-org/cldr/blob/main/common/main/uz.xml) (exemplarCharacters),
  [`common/collation/uz.xml`](https://github.com/unicode-org/cldr/blob/main/common/collation/uz.xml) (oʻ gʻ sh ch → z dan keyin)
- Unicode chart [U02B0 Spacing Modifier Letters](https://www.unicode.org/charts/PDF/U02B0.pdf)
- Wikipedia: [Uzbek alphabet](https://en.wikipedia.org/wiki/Uzbek_alphabet), [Oʻ](https://en.wikipedia.org/wiki/O%CA%BB)
- [lex.uz: OʻRQ-1184, 21.10.2026](https://lex.uz/uz/docs/8548754)
- Rob Pike, [UTF-8 history](https://www.cl.cam.ac.uk/~mgk25/ucs/utf-8-history.txt) (placemat)
- Unicode [NamesList](https://www.unicode.org/Public/UCD/latest/ucd/NamesList.txt) (U+2019 izohi), [UTS #39 confusables](https://www.unicode.org/reports/tr39/)
- Microsoft: [AutoFormat as you Type](https://support.microsoft.com/en-us/office/195c72fc-5e3f-43d5-a0c3-bf3606e2e58e) (smart quotes)

### Uy vazifasi
- 🟢 **Oson:** ismingizda nechta bayt va nechta rune bor?
- 🟡 **Oʻrta:** oʻzbekcha soʻzdagi harflarni sanang: `oʻ`, `gʻ`, `sh`, `ch` bitta harf.
- 🟠 **Qiyin:** soʻzni rune boʻyicha teskari aylantiring va `oʻzbek` bilan sinang: apostrof qaysi harfga yopishdi?
- 🔴 **Guru:** soʻzlarni oʻzbek alifbosi tartibida saralang (`oʻ gʻ sh ch` z dan keyin).

Avval oʻzingiz urinib koʻring. Javoblar: [homework/](homework/).
