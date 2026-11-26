# 04 · `len("oʻzbek")` = 7? Bayt, rune va UTF-8

Video: *Golang oʻzbek tilida #4* · Lab: https://haywan.uz/lab/youtube/go#go-1-04 ·
Oldingi dars: [03 · Konstantalar va `iota`](../003-const-iota/) · Keyingi dars: [05 · Go'da `while` yoʻq](../005-for-switch/) ·
Avval: "Go: JS/Python/PHP dasturchilari uchun" kursi ([0x-01](../0x-01-http-api/) …), [01](../001-why-go/), [03](../003-const-iota/).

Boshqa tilda (JS/TS, Python, PHP) ishlaydigan, Go'ga yangi dasturchilar uchun (2026-11-26 da qayta yozildi).
Natijalar: `go1.26.8 linux/amd64`, Node v24.10.0, Python 3.13.15, PHP 8.4.26. Har bir Go dasturini
https://go.dev/play da ham ishga tushirish mumkin.

**Bitta gʻoya:** Go'da string — baytlar ketma-ketligi (kod ichidagi matn UTF-8 da): `len`, indeks va `==` baytlar
bilan ishlaydi (PHP `strlen` kabi). Python sanaydigan birlik (code point) Go'da `rune`, JS esa UTF-16 birliklarini
sanaydi. Shuning uchun bitta oʻzbekcha soʻzga tillar har xil uzunlik beradi, koʻzga bir xil soʻzlar esa teng chiqmaydi.

| Til | Ifoda | `"oʻzbek"` | `"😀"` | Nimani sanaydi |
|---|---|---|---|---|
| JavaScript / TS | `s.length` | 6 | 2 | UTF-16 birliklari |
| Python | `len(s)` | 6 | 1 | code point'lar |
| PHP | `strlen($s)` · `mb_strlen($s)` | 7 · 6 | 4 · 1 | baytlar · belgilar |
| Go | `len(s)` | 7 | 4 | baytlar |
| Go | `utf8.RuneCountInString(s)` | 6 | 1 | rune'lar (code point'lar) |
| JavaScript | `[...s].length` | 6 | 1 | code point'lar |

`oʻzbek` = **5 harf · 6 rune · 7 bayt** (`6f ca bb 7a 62 65 6b`). `rune` — `int32` ning boshqa nomi, bitta code point.

| Papka | Nima | Natija |
|---|---|---|
| `hook/` | `len("oʻzbek")` | `7` |
| `tillar/tekshir.sh` | bitta soʻz toʻrt tilda (node, python3, php, go — borlari) | `6 · 6 · 7 · 7` |
| `uzunlik/` | `len` va `utf8.RuneCountInString` | `7` / `6` |
| `indeks/` | `s[1]` — bayt (son), `string([]rune(s)[1])` — belgi | `202` / `ca` / `ʻ` |
| `emoji/` | oraliq savol: `len("😀")` | `4` / `1` |
| `tort/` | bitta soʻz, toʻrt xil apostrof | `len` 6 · 7 · 8 · 8, rune hammasida 6 |
| `teng/` | xato: `"Oʻzbekiston" == "O'zbekiston"` | `false` |
| `normal/` | tuzatish: `strings.NewReplacer` faqat o/g dan keyin (`'` va egri `’ ‘`) | `true` / `ma'no` / `true true` |
| `tuzat/` | xato tuzatish: `strings.ReplaceAll` hamma `'` ni almashtiradi | `true` / `maʻno` (buzildi) |
| `tozala/` | toʻliqroq qoida: o/g dan keyin `ʻ`, boshqa harflar orasida `ʼ`, egri `’ ‘` ham | pastda |
| `lexuz/` | lex.uz sahifasida oʻ/gʻ dan keyingi belgilarni sanash (tarmoq kerak) | `U+02BB 50 · U+2018 13 · U+2019 1` |
| `bitlar/` | Chuqurroq: `ʻ` qanday qilib `ca bb` boʻladi | `ca bb` |
| `kesish/` | Chuqurroq: `s[:2]` belgini yarmidan kesadi | `"o\xca"` / `false` / `"oʻ"` / `"o"` |
| `sinov/` | qoʻshimcha mashq: `len("gʻoz")` | `5` |
| `uyga/` | uy vazifasi: boshlangʻich fayllar va `javob/` | |

Qolgan papkalar (`uch/`, `aylanish/`, `bayt/`, `qidiruv/`, `yangi/`, `kirill/`, `homework/`) birinchi versiyadan
qoldi, ular pastdagi **Ilova**da.

## Videodagi buyruqlar
```bash
node -p '"oʻzbek".length'                  # 6
go run hook/main.go                        # 7
python3 -c 'print(len("oʻzbek"))'          # 6
php -r 'echo strlen("oʻzbek");'            # 7   (yoki: ./tillar/tekshir.sh)
go run uzunlik/main.go                     # 7 / 6
go run indeks/main.go                      # 202 / ca / ʻ
go run emoji/main.go                       # 4 / 1     (node -p '"😀".length' → 2)
go run tort/main.go                        # o'zbek len=6 · oʻzbek len=7 · o’zbek len=8 · o‘zbek len=8, rune=6
go run teng/main.go                        # false
go run normal/main.go                      # true / ma'no / true true
go run bitlar/main.go                      # U+02BB = 699 = 1010111011 (10 bit) / 110|01010  10|111011  = ca bb
go run kesish/main.go                      # "o\xca" / false / "oʻ" / "o"
curl -sL https://lex.uz/uz/docs/8548754 | go run ./lexuz   # 2026-11-26: U+02BB 50, U+2018 13, U+2019 1
```

## Dars matni (qisqa)

**Hook.** Bitta oʻzbekcha soʻz: `node` uzunligini 6 deydi, Go'ning `len`i 7. Frontend'dagi 5 belgilik limit va Go
backend'dagi tekshiruv shu farq tufayli bitta ismga ikki xil javob beradi.

**Eslash (03-dars).** `time.Sleep(2)` qancha kutadi? — 2 nanosekund: `time.Duration` nanosekundlarni sanaydi. Bugungi
savol ham birlik haqida: `len` qaytargan 7 — yettita nima?

**Toʻrt til.** Python `len` → 6, PHP `strlen` → 7. Go'da string — baytlar ketma-ketligi (kod ichidagi matn UTF-8 da),
`len` baytlarni sanaydi, PHP `strlen` kabi. Python sanagan birlik Go'da **rune**. JS'niki uchinchi xil: UTF-16 birliklari.

**UTF-8.** Har bir belgining Unicode raqami bor (JS `codePointAt`) — **code point**. UTF-8 uni 1–4 baytga yozadi:
0…127 → 1, …2047 → 2, …65535 → 3, …1114111 → 4 bayt. `ʻ` (U+02BB = 699) ikki bayt: `ca bb`.

**Rune.** Python sanagan oltita narsa — code point'lar: o, ʻ, z, b, e, k. Go'da bitta code point turi `rune` (`int32`).
Rune ham harf emas: `oʻ` bitta harf, lekin ikki rune.

```go
s := "oʻzbek"
fmt.Println(len(s))                    // 7 — baytlar
fmt.Println(utf8.RuneCountInString(s)) // 6 — rune'lar (JS: [...s].length → 6)
```

**Indeks.** JS'da `"oʻzbek"[1]` → `"ʻ"`. Go'da `s[1]` → `202`: bitta **bayt** (`byte` turi, son), `ʻ` ning birinchi
bayti `ca`. Belgi kerak boʻlsa: `string([]rune(s)[1])` → `ʻ`. JS'da baytlarni `new TextEncoder().encode(s)` beradi
(`[1]` → 202) (Go spec: https://go.dev/ref/spec#Index_expressions).

**Oraliq savol.** `len("😀")`? U+1F600 = 128 512 > 65 535 → **4** bayt, rune 1. JS `.length` → 2: UTF-16 bu belgini
ikki birlikka (surrogate pair) boʻladi.

**Toʻrt xil apostrof.** Teskari vergulni oddiy klaviaturada topish qiyin, shuning uchun bitta soʻz toʻrt xil yoziladi.
Unicode'ning til maʼlumotlari (CLDR `uz.xml`) oʻ va gʻ uchun `ʻ` U+02BB ni koʻrsatadi.

```
o'zbek  len=6  rune=6     ' U+0027 = 39     klaviatura apostrofi
oʻzbek  len=7  rune=6     ʻ U+02BB = 699    teskari vergul (MODIFIER LETTER TURNED COMMA)
o’zbek  len=8  rune=6     ’ U+2019 = 8217   oʻng egri apostrof (RIGHT SINGLE QUOTATION MARK)
o‘zbek  len=8  rune=6     ‘ U+2018 = 8216   chap egri apostrof (LEFT SINGLE QUOTATION MARK)
```

**Hayotda.** lex.uz'dagi bitta qonun sahifasining (https://lex.uz/uz/docs/8548754, 2026-11-26) koʻrinadigan matnida
(menyular, sarlavha va footer bilan) oʻ/gʻ dan keyin: `ʻ` 50 marta, `‘` 13, `’` 1; bitta qatorda (sahifa sarlavhasi)
ikki xil: "Oʻzbekiston Respublikasining Qonuni, 21.10.2026 yildagi O‘RQ-1184-son". `‘` larning koʻpi sayt
interfeysida (menyular), `’` — til tanlash tugmasi ("O’zb"). `lexuz/` `<head>`, `<script>`, `<style>` va teglarni
olib tashlab sanaydi.

**Koʻp uchraydigan xato.** Bazada `Oʻzbekiston` (ʻ, 12 bayt), qidiruvda `O'zbekiston` (', 11 bayt). `==` → **false**:
Go string'larni bayt boʻyicha solishtiradi (https://go.dev/ref/spec#Comparison_operators). JS `===` ham false;
Unicode normalizatsiyasi (`normalize("NFC")`, `"NFKC"`) ham yordam bermaydi — bular boshqa-boshqa belgilar.

**Tuzatish.** Har bir `'` ni `ʻ` ga almashtirish xato (`tuzat/`): *maʼno* dagi belgi — **tutuq belgisi** `ʼ` (U+02BC),
`strings.ReplaceAll("ma'no", "'", "ʻ")` → `maʻno`. Faqat o/g dan keyingi apostroflarni (egrilarini ham) almashtiramiz (`normal/`):

```go
r := strings.NewReplacer(
	"o'", "oʻ", "g'", "gʻ", "O'", "Oʻ", "G'", "Gʻ",
	"o’", "oʻ", "g’", "gʻ", "O’", "Oʻ", "G’", "Gʻ",
	"o‘", "oʻ", "g‘", "gʻ", "O‘", "Oʻ", "G‘", "Gʻ",
)
fmt.Println(r.Replace(qidiruv) == saytda) // true
fmt.Println(r.Replace("ma'no"))           // ma'no — tegilmadi
```

Normalizatsiyani bitta funksiyaga yigʻing va uni saqlashdan oldin ham, qidirishdan oldin ham chaqiring. Toʻliqroq
qoida `tozala/` da: o/g dan keyin → `ʻ`, boshqa harflar orasida → `ʼ` (soʻz boshidagi `‘` qoʻshtirnoqqa tegmaydi); uning ham chegaralari bor
(`Go'da` → `Goʻda`), pastdagi Ilovaga qarang.

**Chuqurroq.** `ʻ` = 699 = `1010111011` (10 bit). UTF-8 ning ikki baytli qolipi `110xxxxx 10xxxxxx`:
`110|01010 10|111011` = `ca bb`. `110…`/`1110…`/`11110…` — belgi boshi, `10…` — davomi, shuning uchun dastur istalgan
UTF-8 baytdan belgining boshini topadi (RFC 3629 §3). Baytlab kesish: JS `"oʻzbek".slice(0, 2)` → `"oʻ"`, Go `s[:2]` →
`"o\xca"` (`%q` buzilgan baytni `\xca` deb koʻrsatadi), `utf8.ValidString` → `false`; `string([]rune(s)[:2])` = `"oʻ"`.

Videoga sigʻmaganlar: `for i, r := range s` string ustida rune beradi, `i` esa bayt oʻrni (`aylanish/`, keyingi
darsda `range` shakllari); `%c` bilan har bir baytni chiqarish → `oÊ»zbek` (mojibake, `bayt/`); saralash bayt boʻyicha
(`qidiruv/`); 2027-yildan kuchga kiradigan yangi alifbo (OʻRQ-1184, Ö, Ğ) va `ö` ning ikki xil yozilishi (`yangi/`);
lotin `A` va kirill `А` (`kirill/`).

**Takror (savollar).** 1) `len("oʻzbek")`: Go 7, Python 6 — nega? — Go'ning `len`i baytlarni sanaydi (`ʻ` ikki bayt),
Python code point'larni. 2) `s[1]` nima qaytaradi? — bitta bayt (`202`); belgi uchun `[]rune(s)[1]`. 3) (01-darsdan)
`scratch` image'da `unknown time zone` chiqsa? — `import _ "time/tzdata"` yoki distroless image.

## Uy vazifasi

- 🟢 **Oson** — [`uyga/oson/main.go`](uyga/oson/main.go): roʻyxatdan oʻtish formasi. Frontend'da
  `<input maxLength={5}>` "Gʻani" ni oʻtkazadi (`"Gʻani".length === 5`), Go'dagi `ismTekshir` esa `len(ism) <= 5` bilan
  "juda uzun" deydi. Ishga tushirishdan **oldin** har bir ism uchun natijani yozing, keyin backend frontend bilan bir
  xil sanaydigan qilib tuzating.
- 🟡 **Oʻrta** — [`uyga/orta/main.go`](uyga/orta/main.go): bildirishnoma uchun `qisqa(s, n)` matnning birinchi `n` ta
  belgisini qaytarishi va kesilgan boʻlsa `…` qoʻshishi kerak. Hozir u baytlab kesadi: `"To\xca…"`, `"Gʻalab…"`.
  Rune boʻyicha kesadigan qiling.

### Javoblar
- 🟢 [`uyga/javob/oson`](uyga/javob/oson/main.go): `len("Gʻani")` = 6 (`ʻ` ikki bayt), rune'lar 5. Tekshiruv
  `utf8.RuneCountInString(ism) <= 5` boʻlsin → `Anvar qabul · Gʻani qabul · Oʻktam juda uzun` (Oʻktam frontend'da ham
  6). Eslatma: emoji kabi U+FFFF dan katta belgilarda JS `.length` 2 beradi, rune soni 1 — bunday ismlar uchun
  frontend'da `[...s].length` ishlating.
- 🟡 [`uyga/javob/orta`](uyga/javob/orta/main.go): `r := []rune(s)`; `len(r) <= n` boʻlsa `s`, aks holda
  `string(r[:n]) + "…"` → `"Toʻ…"` · `"Gʻalaba"` · `"Payme…"`.

Qoʻshimcha mashqlar (birinchi versiyadan, qiyinroq): [`homework/`](homework/).

## Manbalar
- Go qoidalari: String types (string = baytlar, `len` = baytlar) https://go.dev/ref/spec#String_types · Index
  expressions (`s[i]` = bayt) https://go.dev/ref/spec#Index_expressions · Numeric types (`rune` = `int32`)
  https://go.dev/ref/spec#Numeric_types · Comparison operators https://go.dev/ref/spec#Comparison_operators
- `unicode/utf8` https://pkg.go.dev/unicode/utf8 · `strings.NewReplacer` https://pkg.go.dev/strings#NewReplacer ·
  `strings.ReplaceAll` https://pkg.go.dev/strings#ReplaceAll
- JS: ECMAScript String type (16-bit elements) https://tc39.es/ecma262/#sec-ecmascript-language-types-string-type ·
  Python: str = code points https://docs.python.org/3/library/stdtypes.html#text-sequence-type-str · PHP: `strlen`
  baytlarni sanaydi https://www.php.net/manual/en/function.strlen.php
- RFC 3629 (UTF-8) https://www.rfc-editor.org/rfc/rfc3629
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
