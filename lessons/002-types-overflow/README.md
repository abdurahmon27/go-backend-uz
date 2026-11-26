# 02 · Turlar va overflow: Go rad etadi, `x++` esa jim aylanadi

Video: *Golang o'zbek tilida #2* · Lab: https://haywan.uz/lab/youtube/go#go-1-02 ·
Oldingi dars: [01 · Go nega bor?](../001-why-go/) · Keyingi dars: [03 · Konstantalar va `iota`](../003-const-iota/) ·
Avval koʻring: Go 0 darslari ([000-01](../000-01-birinchi-dastur/) … [000-05](../000-05-for/)).

Go 0 dan keyin koʻriladigan dars. Hamma natijalar `go1.26.8 linux/amd64`da olingan (2026-11-26).
Har bir dasturni https://go.dev/play da ham ishga tushirish mumkin: faylni nusxalab qoʻying va **Run** bosing.

**Bitta gʻoya:** son turi katakning kattaligini ham aytadi. `int8` katagiga faqat −128 … 127 sigʻadi. Sigʻmaydigan son
oʻzgaruvchisiz, faqat sonlardan chiqsa, Go dasturni tarjima qilmaydi; hisobda oʻzgaruvchi qatnashsa, son mashina
hisoblagichi kabi jim aylanadi (127 dan keyin −128). Qiymat berilmagan katakda doim shu turning nol qiymati turadi.

| Papka | Nima | Natija |
|---|---|---|
| `aylanish/` | `var x int8 = 127; x++` | `-128` |
| `sigmaydi/` | `var x int8 = 128` — ataylab tarjima boʻlmaydi | `cannot use 128 (untyped int constant) as int8 value in variable declaration (overflows)` |
| `sinov/` | oʻrtadagi sinov: `var a int8 = 100; a = a + 50` | `-106` |
| `nol/` | zero value: `int8`, `bool`, `string` | `0 false []` |
| `layk/` | xato: `int8` bilan 200 ta layk | `-56` |
| `layk-int/` | tuzatilgan: `int` | `200` |
| `constover/` | Chuqurroq: `int8(127) + 1` — ataylab tarjima boʻlmaydi | `int8(127) + 1 (constant 128 of type int8) overflows int8` |
| `binary/` | Chuqurroq: bitlar | `127 = 01111111`, `-128 = 10000000`, `-1 = 11111111` |
| `uyga/` | uy vazifasi: boshlangʻich fayllar va `javob/` | |

Qolgan papkalar (`sizes/`, `mix/`, `zero/`, `wrap/`, `lt/`, `convert/`, `boeing/`, `bug/`, `fix/`, `homework/`)
darsning birinchi versiyasidan qoldi, ular pastdagi **Ilova**da.

## Videodagi buyruqlar
```bash
go run aylanish/main.go    # -128
go run sigmaydi/main.go    # sigmaydi/main.go:6:15: cannot use 128 (untyped int constant) as int8 value in variable declaration (overflows)
go run sinov/main.go       # -106
go run nol/main.go         # 0 false []
go run layk/main.go        # -56
go run layk-int/main.go    # 200
go run constover/main.go   # constover/main.go:6:14: int8(127) + 1 (constant 128 of type int8) overflows int8
go run binary/main.go      #  127 = 01111111 / -128 = 10000000 / -1 = 11111111
```

## Dars matni (qisqa)

**Hook.** Mashinaning kilometr hisoblagichida hamma raqamlar 9. Yana bir kilometrdan keyin hamma gʻildirakcha aylanadi
va nollar chiqadi, hech qanday ogohlantirishsiz. Go'dagi baʼzi sonlar ham shunday aylanadi, faqat minusga.

**Katak.** Memory rasmidagi qiymat turadigan toʻrtburchak — *katak*. Oddiy `int` juda katta katak (bu kompyuterda
64 bit; Go qoidasi boʻyicha 32 yoki 64 bit boʻlishi mumkin — ikkalasida ham `int8` dan juda katta). `int8` esa kichik:
eng kichigi −128, eng kattasi 127 (256 ta son: yarmi manfiy, yarmi 0 va musbat). `:=` oddiy `int` tanlaydi, shuning
uchun `int8` katakni `var` bilan yaratamiz: `var x int8 = 127` — "x yorligʻi, turi `int8`, qiymati 127".

**Aylanish.** `x++` dan keyin `x` = −128. Sonlarni soat yuzi kabi doira boʻylab tasavvur qiling: 127 dan keyin eng
kichik son keladi. Bu *overflow* (toʻlib ketish). [Go qoidasi](https://go.dev/ref/spec#Integer_overflow):
oʻzgaruvchidagi overflow runtime panic bermaydi, natija aniq belgilangan; `go run` xato ham, ogohlantirish ham bermaydi.

**Rad etish.** `var x int8 = 128` tarjima boʻlmaydi: oʻngdagi 128 — oʻzgaruvchisiz yozilgan son (konstanta), Go uni
tarjima paytida tekshiradi ([Constants](https://go.dev/ref/spec#Constants): "It is an error if the constant value
cannot be represented as a value of the respective type"). `a = a + 50` da esa 50 ham kodda yozilgan, lekin hisobda
`a` oʻzgaruvchisi qatnashadi — natija dastur ishlayotganda chiqadi va jim aylanadi (`-106`: 100 dan 127 gacha 27 qadam,
keyin yana 23 qadam: −128, −127, …, −106).

Uchinchi jim yoʻl (videoda yoʻq): oʻzgaruvchidan **konversiya**. `big := 300; int8(big)` → `44`, `int8(200)` → `-56`:
butun sonlar orasidagi konversiya pastki 8 bitni oladi va [hech qanday belgi bermaydi](https://go.dev/ref/spec#Conversions)
(`convert/`).

**Nol qiymat (zero value).** `var y int8` → `0`, `var tayyor bool` → `false`, `var ism string` → `""` (boʻsh matn;
koʻrinishi uchun `"["+ism+"]"` deb chiqardik → `[]`). [Go qoidasi](https://go.dev/ref/spec#The_zero_value): qiymat
berilmagan oʻzgaruvchi hech qachon tasodifiy axlat saqlamaydi.

**Koʻp uchraydigan xato.** "Joyni tejash" uchun layk hisoblagichiga `int8`: 200 ta layk → `-56` (128-laykda −128 ga
aylanadi, qolgan 72 ta uni −56 ga olib keladi). Xato chiqmaydi — bug. Tuzatish: `int`.
[A Tour of Go](https://go.dev/tour/basics/11): *"use int unless you have a specific reason to use a sized or unsigned
integer type"*.

**Hayotda: Boeing 787.** FAA majburiy koʻrsatmasi AD 2015-09-07: generator boshqaruv bloklaridagi (GCU) dasturiy
hisoblagich 248 kun uzluksiz tok ostida turganda toʻlib qoladi; toʻrtta asosiy GCU bir vaqtda yoqilgan boʻlsa, samolyot
barcha AC elektr quvvatini yoʻqotishi va "loss of control of the airplane"ga olib kelishi mumkin edi. Vaqtinchalik
choralar: quvvatni "at intervals not to exceed 120 days" oʻchirib-yoqish. Keyinroq AD 2018-20-15 yangi GCU dasturini
talab qildi. FAA hisoblagich turini aytmagan: `boeing/` dagi `int32`, sekundiga 100 marta degan hisob (2³¹ / 100 s =
248.55 kun) — bizning farazimiz.

**Chuqurroq (ixtiyoriy).** 1 bit = 1 chiroqcha (0/1); `int8` = 8 bit. `127 = 0111 1111`, `+1` → `1000 0000`. Eng chap
bitning vazni −128 (two's complement), shuning uchun `1000 0000` = −128, `1111 1111` = −128 + 127 = −1.
`int8(127) + 1` ham tarjima boʻlmaydi: hamma son konstanta (`constover/`, keyingi darsda batafsil).

## Uy vazifasi
- 🟢 **Oson** (`uyga/oson/main.go`, videodagi `layk` dasturi): `200` ni avval `127` ga, keyin `128` ga almashtiring.
  Har safar ishga tushirishdan oldin natijani taxmin qiling.
- 🟡 **Oʻrta** (`uyga/orta/main.go`): `int16` (−32768 … 32767) qaysi qadamda manfiy boʻladi? Sikl ichida `if son < 0`
  bilan `i` va `son` ni chiqaring va `break` bilan chiqing.
- Qiyinroq vazifalar (2038-yil muammosi, `int64` qoʻshishda overflow'ni aniqlash, `abs(int8(-128))`):
  [homework/](homework/).

<details><summary>Javoblar</summary>

🟢 `i <= 127` → `127` (hali sigʻadi), `i <= 128` → `-128` (128-laykda doira boshiga). [`uyga/javob/oson`](uyga/javob/oson/main.go)

🟡 [`uyga/javob/orta`](uyga/javob/orta/main.go) → `32768 -32768`: 32768-qadamda `int16` 32767 dan −32768 ga aylanadi.
</details>

## Manbalar
- Go qoidalari (spec): [Integer overflow](https://go.dev/ref/spec#Integer_overflow), [Constants](https://go.dev/ref/spec#Constants),
  [The zero value](https://go.dev/ref/spec#The_zero_value), [Numeric types](https://go.dev/ref/spec#Numeric_types),
  [Conversions](https://go.dev/ref/spec#Conversions)
- A Tour of Go, Basic types: https://go.dev/tour/basics/11
- FAA AD 2015-09-07, Boeing 787 (Federal Register, 2015-05-01): https://www.govinfo.gov/content/pkg/FR-2015-05-01/html/2015-10066.htm ·
  AD 2018-20-15: https://www.govinfo.gov/content/pkg/FR-2018-10-16/html/2018-22152.htm
- CMU 15-213, *Bits, Bytes & Integers*: https://www.cs.cmu.edu/~213/lectures/02-bits-bytes-ints.pdf

---

## Ilova: darsning birinchi versiyasi (2026-11-25)

Video 2026-11-26 da qayta yozildi (Go 0 mavsumidan keyin, noldan boshlovchilar uchun). Birinchi versiyaning matni va
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
