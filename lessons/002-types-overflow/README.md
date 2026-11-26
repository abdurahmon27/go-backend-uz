# 02 · Turlar va overflow: Go rad etadi, `x++` esa jim aylanadi

Video: *Golang oʻzbek tilida #2* · Lab: https://haywan.uz/lab/youtube/go#go-1-02 ·
Oldingi dars: [01 · Go nega bor?](../001-why-go/) · Keyingi dars: [03 · Konstantalar va `iota`](../003-const-iota/)

Boshqa tilda (JS/TS, Python, PHP) ishlaydigan, Go'ga yangi dasturchilar uchun. Avval "Go: JS/Python/PHP dasturchilari
uchun" kursi (GO-X): `go run`, statik turlar, `var` va `:=`, zero value, `error` qiymatlar.
Hamma natijalar shu kompyuterda olingan (2026-11-26): `go1.26.8 linux/amd64`, Node v24.10.0, Python 3.13.15, tsc 5.9.3.
Go dasturlarini https://go.dev/play da ham ishga tushirish mumkin.

**Bitta gʻoya:** Go'da butun sonning chegarasini siz tanlagan tur belgilaydi (`int8` = −128 … 127). Oʻzgaruvchidagi
qiymat chegaradan oʻtsa, two's complement boʻyicha jim aylanadi (Python'da chegara yoʻq, JS'da aniqlik yoʻqoladi).
Konstanta ifoda esa aniq hisoblanadi; yakuniy qiymat turga sigʻmasa, Go uni kompilyatsiya paytida rad etadi.

| Papka | Nima | Natija |
|---|---|---|
| `hook/` | `var id int64 = math.MaxInt64; id++` (+ Python va Node bilan solishtirish) | `-9223372036854775808` |
| `hook/id.js` | frontend 64 bitlik ID'ni JSON'dan son qilib oʻqiydi | `1234567890123456800` |
| `jsonid/` | Go: `json:"id,string"` — ID JSON'da string boʻlib ketadi | `{"id":"1234567890123456789"}` |
| `eslash/` | takrorlash: `var layk int; layk + 1` (JS'da `NaN`) | `1` |
| `sizes/` | butun son turlari: hajm va chegaralar | `int8 1 bayt: -128 ... 127` … |
| `aylanish/` | `var x int8 = 127; x++` | `-128` |
| `binary/` | bitlar: two's complement | `127 = 01111111`, `-128 = 10000000` |
| `sigmaydi/` | `var x int8 = 128` — ataylab kompilyatsiya boʻlmaydi | `cannot use 128 (untyped int constant) as int8 value … (overflows)` |
| `const64/` | `var id int64 = math.MaxInt64 + 1` — ataylab kompilyatsiya boʻlmaydi | `cannot use math.MaxInt64 + 1 (untyped int constant 9223372036854775808) as int64 value … (overflows)` |
| `ifoda/` | oʻrtadagi sinov: `var a int8 = 100 + 50 - 60` (oraliqda 150) | `90` |
| `sinov/` | oʻzgaruvchi bilan: `var a int8 = 100; a = a + 50` | `-106` |
| `layk/` | xato: `int8` bilan 200 ta layk | `-56` |
| `layk-int/` | tuzatilgan: `int` | `200` |
| `convert/` | Chuqurroq: konversiya `int8(big)` | `44`, `-56` |
| `uyga/` | uy vazifasi: boshlangʻich fayllar va `javob/` | |

## Videodagi buyruqlar
```bash
python3 -c 'print(9223372036854775807 + 1)'   # 9223372036854775808  (Python int: chegara yoʻq)
node -p '9223372036854775807 + 1'             # 9223372036854776000  (JS number = float64, yaxlitlandi)
go run hook/main.go                           # -9223372036854775808 (Go int64: aylandi)
node -p 'let layk; layk + 1'                  # NaN
go run eslash/main.go                         # 1  (zero value)
go run sizes/main.go
go run aylanish/main.go                       # -128
node -p 'const a = new Int8Array([127]); a[0]++; a[0]'   # -128 (JS ham 8 bitda aylanadi)
go run binary/main.go                         #  127 = 01111111 / -128 = 10000000 / -1 = 11111111
go run sigmaydi/main.go                       # sigmaydi/main.go:6:15: cannot use 128 (untyped int constant) as int8 value in variable declaration (overflows)
go run const64/main.go                        # const64/main.go:10:17: cannot use math.MaxInt64 + 1 (untyped int constant 9223372036854775808) as int64 value in variable declaration (overflows)
go run ifoda/main.go                          # 90
go run layk/main.go                           # -56
go run layk-int/main.go                       # 200
node hook/id.js                               # {"id":1234567890123456789} / 1234567890123456800
go run jsonid/main.go                         # {"id":"1234567890123456789"}
go run convert/main.go                        # 44 / -56
node -p 'new Int8Array([300])[0]'             # 44
```
`sigmaydi/` va `const64/` ataylab kompilyatsiya boʻlmaydi, shuning uchun `go vet ./...` bu ikki papkada xato beradi.

TypeScript tekshiruvi (`tsc 5.9.3`, `strict: true`): `const id: number = 9223372036854775807 + 1` xatosiz kompilyatsiya
boʻladi — TS son chegarasini tekshirmaydi, hamma son `number`, qiymat esa ishlash paytida yaxlitlanadi.
JS'da aniq katta butun son kerak boʻlsa — `BigInt`: `9223372036854775807n + 1n` → `9223372036854775808n`.

## Dars matni (qisqa)

**Hook.** Bir xil son, bir xil `+ 1`: Python `9223372036854775808` (aniq), Node `9223372036854776000` (yaxlitlandi),
Go `int64` oʻzgaruvchida `-9223372036854775808` (manfiy, xatosiz).

**Takrorlash (GO-X).** JS'da `let layk; layk + 1` → `NaN`. Go'da `undefined` yoʻq: `var layk int` → zero value `0`,
`layk + 1` → `1`.

**Chegaralar.** Python `int` — [cheklanmagan aniqlik](https://docs.python.org/3/library/stdtypes.html#numeric-types-int-float-complex).
JS `number` — float64, butun son faqat 2⁵³ gacha aniq ([`Number.MAX_SAFE_INTEGER`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER) = 9007199254740991;
`2**53 + 1` → `9007199254740992`). Go — tur nomidagi raqam bitlar soni: `int8`, `int16`, `int32`, `int64`; `int` bu
kompyuterda 64 bit ([spec](https://go.dev/ref/spec#Numeric_types): 32 yoki 64). `:=` butun son literal uchun doim `int`
tanlaydi, shuning uchun `int8` kerak boʻlsa `var x int8 = 127`.

**Aylanish va bitlar.** `x++` dan keyin `-128`: eng kattasidan keyin eng kichigi. `0111 1111 + 1 = 1000 0000`; eng chap
bitning vazni −128 (*two's complement*), shuning uchun `1000 0000` = −128, `1111 1111` = −1. Bu *overflow*.
[Go spec](https://go.dev/ref/spec#Integer_overflow): signed overflow qonuniy, natija aniq belgilangan, runtime panic
yoʻq. JS'dagi `Int8Array` ham xuddi shunday aylanadi (`127` → `++` → `-128`).

**Rad etish.** `var x int8 = 128`: `128` konstanta, qiymati kompilyatsiya paytida maʼlum, Go uni `int8` ga sigʻishini
tekshiradi ([Constants](https://go.dev/ref/spec#Constants): "It is an error if the constant value cannot be represented
as a value of the respective type"). Konstanta ifodalar [aniq hisoblanadi](https://go.dev/ref/spec#Constant_expressions)
("intermediate values … may require precision significantly larger than supported by any predeclared type"): xato
faqat yakuniy qiymat turga qoʻyilganda va sigʻmaganda chiqadi. Shuning uchun `var a int8 = 100 + 50 - 60` → `90`
(`ifoda/`), `var id int64 = math.MaxInt64 + 1` esa rad etiladi (`const64/`), `const big = 1 << 100` ham qonuniy
(turi yoʻq konstanta, keyingi darsda). Ifodada oʻzgaruvchi qatnashsa (`x++`, `a + 50`), natija ishlash paytida chiqadi
va jim aylanadi: `100 + 50 = 150 → 150 − 256 = −106` (`sinov/`).

**Koʻp uchraydigan xato.** "Joyni tejash" uchun layk hisoblagichiga `int8`: 200 ta layk → `-56`. Xato ham, log ham yoʻq.
Tuzatish: `int`. [A Tour of Go](https://go.dev/tour/basics/11): *"use int unless you have a specific reason to use a
sized or unsigned integer type"*.

**Frontend.** Backend `int64` ID yuborsa, JS uni JSON'dan son qilib oʻqiganda oxirgi raqamlar oʻzgaradi
(`node hook/id.js` → `1234567890123456800`). [X API hujjati](https://docs.x.com/fundamentals/x-ids): *"In JavaScript,
integers are limited to 53 bits. This causes precision loss with large IDs"*, *"Always use string IDs in your code"*
(2010-yildan beri `id_str`, [eʼlon](https://groups.google.com/g/twitter-api-announce/c/ahbvo3VTIYI)). Go'da:
`json:"id,string"` tegi `int64` ni JSON'da string qilib yozadi (`jsonid/`; structlar keyingi darslarda). Frontend
uni string yoki `BigInt` qilib saqlaydi.

**Hayotda: Boeing 787.** FAA majburiy koʻrsatmasi AD 2015-09-07 (faqat 787 modeli): generator boshqaruv bloklaridagi (GCU) dasturiy
hisoblagich 248 kun uzluksiz tok ostida turganda toʻlib qoladi; toʻrtta asosiy GCU bir vaqtda yoqilgan boʻlsa, samolyot
barcha AC elektr quvvatini yoʻqotishi va "loss of control of the airplane"ga olib kelishi mumkin edi. Vaqtinchalik
chora: quvvatni "at intervals not to exceed 120 days" oʻchirib-yoqish. AD 2018-20-15 yangi GCU dasturini talab qildi.
FAA hisoblagich turini aytmagan.

**Chuqurroq: konversiya.** `big := 300; int8(big)` → `44`: butun sonlar orasidagi konversiya pastki 8 bitni oladi va
[hech qanday belgi bermaydi](https://go.dev/ref/spec#Conversions) (`300 = 1 0010 1100` → `0010 1100` = 44). JS'da
`new Int8Array([300])[0]` ham `44`. Bazadan `int64` oʻqib `int32` ga aylantiradigan kodda ham shunday. `float` → `int` sigʻmasa natija implementatsiyaga bogʻliq.

## Uy vazifasi
- 🟢 **Oson** (`uyga/oson/main.go`, videodagi `jsonid`): postda `AuthorID int64` ham bor. Uni frontend buzmaydigan qilib
  yuboring. Tekshirish: `go run ./uyga/oson | node uyga/check.js` — `author_id` aniq `9007199254740993` boʻlib kelishi kerak
  (hozir `number 9007199254740992`).
- ➕ **Qoʻshimcha** (`uyga/int16/main.go`): `int16` (−32768 … 32767) qaysi qadamda manfiy boʻladi? Sikl ichida
  `if son < 0` bilan `i` va `son` ni chiqaring va `break` bilan chiqing.
- 🟡 **Oʻrta** (`uyga/orta/main.go`): `safeAdd(a, b int64) (int64, error)` — yigʻindi `int64` ga sigʻmasa, natija oʻrniga
  `ErrOverflow` qaytarsin. `main` dagi toʻrtta chaqiruv `kutilgan:` izohidagidek chiqishi kerak.
- Qiyinroq vazifalar (2038-yil muammosi, `abs(int8(-128))`): [homework/](homework/).

<details><summary>Javoblar</summary>

🟢 [`uyga/javob/oson`](uyga/javob/oson/main.go): `` AuthorID int64 `json:"author_id,string"` `` →
`{"id":"1234567890123456789","author_id":"9007199254740993"}`, `check.js`: `author_id string 9007199254740993`.

➕ [`uyga/javob/int16`](uyga/javob/int16/main.go) → `32768 -32768`: 32768-qadamda `int16` 32767 dan −32768 ga aylanadi.

🟡 [`uyga/javob/orta`](uyga/javob/orta/main.go) → `5 <nil>`, `0 int64 overflow`, `0 int64 overflow`, `-1 <nil>`.
Overflow faqat ishorasi bir xil ikki son qoʻshilganda boʻladi: musbat + musbat manfiy chiqsa yoki manfiy + manfiy
musbat (yoki 0) chiqsa. Go'da `a + b` panic bermaydi, shuning uchun avval qoʻshib, keyin ishorani tekshirish xavfsiz.
</details>

## Manbalar
- Go spec: [Integer overflow](https://go.dev/ref/spec#Integer_overflow), [Constants](https://go.dev/ref/spec#Constants),
  [Constant expressions](https://go.dev/ref/spec#Constant_expressions), [Numeric types](https://go.dev/ref/spec#Numeric_types),
  [The zero value](https://go.dev/ref/spec#The_zero_value), [Conversions](https://go.dev/ref/spec#Conversions)
- A Tour of Go, Basic types: https://go.dev/tour/basics/11
- ECMAScript ToInt8 (`Int8Array`): https://tc39.es/ecma262/#sec-toint8
- encoding/json (`,string`): https://pkg.go.dev/encoding/json#Marshal
- Python, Numeric types: https://docs.python.org/3/library/stdtypes.html#numeric-types-int-float-complex
- MDN, `Number.MAX_SAFE_INTEGER`: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
- X API, X IDs: https://docs.x.com/fundamentals/x-ids · Twitter API eʼloni (2010-10-19): https://groups.google.com/g/twitter-api-announce/c/ahbvo3VTIYI
- FAA AD 2015-09-07, Boeing 787: https://www.govinfo.gov/content/pkg/FR-2015-05-01/html/2015-10066.htm ·
  AD 2018-20-15: https://www.govinfo.gov/content/pkg/FR-2018-10-16/html/2018-22152.htm
- CMU 15-213, *Bits, Bytes & Integers*: https://www.cs.cmu.edu/~213/lectures/02-bits-bytes-ints.pdf

---

## Ilova: darsning birinchi versiyasi (2026-11-25)

Video 2026-11-26 da dasturchilar uchun qayta yozildi. Birinchi versiyaning matni va
dasturlari shu yerda qoldi: oʻlchamlar jadvali, `int` va `int64` aralashtirish, `x < x+1`, konversiya savoli,
Ariane 5, Boeing hisobi va `uint` ombor bug'i.

| Papka | Nima |
|---|---|
| `sizes/` | butun son turlari: hajm va chegaralar |
| `mix/` | `int + int64`: ataylab kompilyatsiya bo'lmaydi |
| `zero/` | zero value |
| `constover/` | `int8(127) + 1`: ataylab kompilyatsiya bo'lmaydi |
| `wrap/` | `x++` → `-128`, `0 - 1` → `255` |
| `lt/` | `x < x+1` har doim ham true emas |
| `binary/` | two's complement: bitlarni ko'ramiz |
| `convert/` | pauza-savol: `int8(300)`, keyin `int8(200)` |
| `boeing/` | 248 kun: int32 hisoblagich qachon to'ladi |
| `bug/`, `fix/` | debugging mashqi: ombordagi qoldiq `uint` bilan |
| `homework/` | eski uy vazifasi javoblari |

```bash
go run ./sizes
go run ./mix         # invalid operation: a + b (mismatched types int and int64)
go run ./zero        # 0 0 false / ""
go run ./wrap        # -128 / 255
go run ./lt          # true / false
go run ./convert     # 44 / -56
go run ./boeing      # 248.55 kun / -2147483648
go run ./bug         # qoldiq: 18446744073709551614
go run ./fix         # yetarli emas: bor 3 ta, so'ralgan 5 ta
```

### Birinchi versiyaning matni

**Hook.** `int8(127) + 1` ifodasini Go kompilyatsiya qilmaydi. Lekin xuddi shu hisobni o'zgaruvchida qilsangiz, dastur
hech qanday xatosiz `-128` chiqaradi. Bu ikki xatti-harakat qayerdan keladi?

**Bu faqat nazariya emas.** 2015-yilda AQSh aviatsiya boshqarmasi (FAA) Boeing 787 uchun majburiy ko'rsatma chiqardi:
generator boshqaruv bloklaridagi (GCU) dasturiy hisoblagich 248 kun to'xtovsiz tok ostida tursa, to'lib qoladi. To'rtta
asosiy GCU bir vaqtda yoqilgan bo'lsa, samolyot barcha AC elektr quvvatini yo'qotishi mumkin. Vaqtinchalik yechim:
kamida har 120 kunda elektr quvvatini o'chirib, qayta yoqish.

**Turlar.** `int8`, `int16`, `int32`, `int64`: raqam bitlar sonini bildiradi. `int8` 1 bayt, −128 … 127. Har birining
`u` bilan boshlanadigan unsigned jufti bor: `uint8` 0 … 255. `byte` = `uint8`, `rune` = `int32`. Oddiy `int` bu
kompyuterda 8 bayt, boshqa tizimda 4 bayt bo'lishi mumkin (spec: 32 yoki 64 bit). `int` va `int64` bir xil hajmda bo'lsa ham, Go uchun ular
boshqa-boshqa tur: aralashtirish uchun aniq konversiya kerak, shunda kod bitta kompyuterda ishlab, boshqasida sinmaydi.

| Tur | Bayt | Oraliq |
|---|---|---|
| `int8` | 1 | −128 … 127 |
| `int16` | 2 | −32 768 … 32 767 |
| `int32` | 4 | −2 147 483 648 … 2 147 483 647 |
| `int64` | 8 | −9 223 372 036 854 775 808 … 9 223 372 036 854 775 807 |
| `uint8` (`byte`) | 1 | 0 … 255 |

**Zero value.** Qiymat berilmagan o'zgaruvchi hech qachon tasodifiy axlat saqlamaydi: sonlar uchun `0`, `bool` uchun
`false`, `string` uchun bo'sh string `""` (pointer, slice, map uchun `nil`: keyingi darslarda). C'da funksiya ichidagi oddiy (static bo'lmagan) o'zgaruvchi uchun bunday kafolat yo'q; Go'da bu spec'da yozilgan.

**Konstantadagi overflow.** `int8(127) + 1` to'liq konstantalardan iborat, Go uni kompilyatsiya paytida hisoblaydi.
Natija 128 `int8`ga sig'maydi, va Go dasturni build qilmaydi. (Nega konstanta bunday qattiq tekshiriladi, keyingi
darsda.)

**O'zgaruvchidagi overflow.** `var x int8 = 127; x++` → `-128`. `var u uint8 = 0; u--` → `255`. Panic ham, warning ham
yo'q. `x` konstanta emas, o'zgaruvchi: kompilyator natijani bilsa ham (bu dasturda u hatto oldindan hisoblaydi),
[spec](https://go.dev/ref/spec#Integer_overflow) o'zgaruvchidagi overflow'ni xato emas, aniq qoida deb belgilagan:
signed overflow runtime panic bermaydi va natija aniq; unsigned turlarda hisob 2ⁿ bo'yicha qoldiq bilan (n = bitlar
soni). Bu kompilyator xatosi emas, lekin *sizning* kodingizda bug bo'lishi mumkin.

Spec yana talab qiladi: kompilyator overflow bo'lmaydi deb faraz qila olmaydi, ya'ni `x < x+1` har doim ham true
emas (`lt/`: `5` da true, `math.MaxInt32` da false). C'da signed overflow undefined behavior (C11 §6.5 ¶5), va
kompilyator bunday tekshiruvni olib tashlashi mumkin; Go'da bunday bo'lishi mumkin emas.

**Nega aynan −128?** Go butun sonlarni two's complement ko'rinishida saqlaydi. `int8`da eng chap bitning vazni −128:
```
 127 = 0111 1111
+  1
-128 = 1000 0000      (−128 + 0)
  -1 = 1111 1111      (−128 + 127)
```

**Pauza-savol.** `big := 300; int8(big)` nima chiqaradi? `44`. Butun sonlar orasidagi konversiya pastki 8 bitni oladi
va ularni two's complement bo'yicha o'qiydi. 300 = `1 0010 1100`, pastki 8 bit `0010 1100` = 44. 200 = `1100 1000`:
eng chap bit 1, demak `−128 + 72 = −56`. [Spec](https://go.dev/ref/spec#Conversions): butun sonlar orasidagi
konversiya har doim qiymat beradi, overflow haqida hech qanday belgi bermaydi. (float → int konversiyada sig'masa,
natija implementatsiyaga bog'liq.)

**Ariane 5 (1996).** Konversiya xatosining eng mashhur misoli: 501-parvozda dastur 64 bitli float qiymatni 16 bitli
butun songa aylantirdi, son sig'madi. Ada'da bu xato jim o'tmadi: Operand Error chiqdi, uni ushlaydigan kod yo'q edi,
navigatsiya kompyuteri o'chdi, va raketa parvozning birinchi daqiqasida yo'qotildi. Bu kod
Ariane 4 dan olingan va Ariane 5 ning parvoz ma'lumotlari bilan sinovdan o'tkazilmagan edi
([Inquiry Board report](https://csapp.cs.cmu.edu/2e/docs/ariane5rep.html)).

**248 kun.** FAA ko'rsatmasida hisoblagich turi aytilmagan. Faraz qilaylik, u `int32` va har 10 millisekundda bittaga
oshadi: 2³¹ / 100 sekund = **248.55 kun**. Aynan shu (2³¹-) qadamda hisoblagich `-2147483648` ga sakraydi. Chegara har
qanday tilda shu; undan keyin nima bo'lishi tilga bog'liq: Go va Java'da manfiyga aylanadi, C'da undefined behavior.

**Debugging.** Omborda 3 ta telefon, mijoz 5 ta buyurtma qildi; dasturchi "qoldiq manfiy bo'lmaydi" deb `uint`
tanlagan. `3 - 5` → `18446744073709551614`: `uint` noldan pastga tusha olmaydi; bu kompyuterda `uint` 64 bit, shuning uchun
ayirish 2⁶⁴ bo'yicha aylanadi (32 bitli tizimda natija `4294967294` bo'lardi).
Bu son JSON javobda frontend'ga ketsa, sayt shuncha telefon bor deb ko'rsatadi (`database/sql` esa bunday `uint64`ni
bazaga yozishni rad etadi: *uint64 values with high bit set are not supported*). Yechim: ayirishdan oldin tekshirish
(`fix/`). [Go Tour](https://go.dev/tour/basics/11) qoidasi: maxsus sabab bo'lmasa, `int` ishlating. Unda natija `-2`: bu ham
xato, lekin uni darhol ko'rasiz. `uint` "manfiy bo'lmaydi" degani emas, "manfiy bo'lsa aylanadi" degani. Real backend'da bu tekshiruv va ayirish bazadagi bitta so'rovda
bo'lishi kerak (`UPDATE … SET stock = stock - $1 WHERE stock >= $1`), aks holda ikki so'rov bir vaqtda tekshiruvdan
o'tib ketishi mumkin: bu haqda ma'lumotlar bazasi darslarida.

### Manbalar (birinchi versiya)
- Go spec: [Numeric types](https://go.dev/ref/spec#Numeric_types), [Integer overflow](https://go.dev/ref/spec#Integer_overflow),
  [The zero value](https://go.dev/ref/spec#The_zero_value), [Conversions](https://go.dev/ref/spec#Conversions),
  [Constants](https://go.dev/ref/spec#Constants)
- FAA AD 2015-09-07, Boeing 787 (Federal Register, 2015-05-01): https://www.govinfo.gov/content/pkg/FR-2015-05-01/html/2015-10066.htm
- CMU 15-213, *Bits, Bytes & Integers*: https://www.cs.cmu.edu/~213/lectures/02-bits-bytes-ints.pdf
- C11 standarti (N1570) §6.7.9 ¶10 (initsializatsiya qilinmagan lokal o'zgaruvchi) va §6.5 ¶5 (signed overflow = UB): https://www.open-std.org/jtc1/sc22/wg14/www/docs/n1570.pdf
- Ariane 501 Inquiry Board report (1996): https://csapp.cs.cmu.edu/2e/docs/ariane5rep.html
- A Tour of Go, Basic types: https://go.dev/tour/basics/11


### Uy vazifasi (birinchi versiya)
- 🟢 `int16` bilan overflow qilib, natijani tushuntiring. · 🟡 Unix vaqti `int32`da saqlansa, qaysi sanada toʻladi? ·
  🟠 ikkita `int64` qoʻshilganda overflow boʻlganini aniqlaydigan funksiya. · 🔴 `abs(int8(-128))` nima chiqadi, nega?
  Javoblar: [homework/](homework/).
