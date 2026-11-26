# 03 · Konstanta va `iota`: oʻzgarmas nomlar, `time.Sleep(2)` va jim buziladigan raqamlar

Video: *Golang o'zbek tilida #3* · Lab: https://haywan.uz/lab/youtube/go#go-1-03 ·
Oldingi dars: [02 · Turlar va overflow](../002-types-overflow/) · Keyingi dars: [04 · Matn, byte, rune va UTF-8](../004-strings-utf8/) ·
Avval koʻring: Go 0 darslari ([000-01](../000-01-birinchi-dastur/) … [000-05](../000-05-for/)) va 02-dars.

Go 0 va 02-darsdan keyin koʻriladigan dars. Hamma natijalar `go1.26.8 linux/amd64`da olingan (2026-11-26).
Har bir dasturni https://go.dev/play da ham ishga tushirish mumkin: faylni nusxalab qoʻying va **Run** bosing.

**Bitta gʻoya:** konstanta oʻzgarmaydi, lekin uning maʼnosi qayerda turganiga bogʻliq. Turi yoʻq konstanta turini
tushgan katagidan oladi (`time.Sleep(2)` → ikki nanosekund), `iota` esa raqamini `const ( … )` ichidagi qator oʻrnidan
oladi (roʻyxat oʻrtasiga qator qoʻshilsa, faylda saqlangan raqamlarning maʼnosi jim oʻzgaradi).

| Papka | Nima | Natija |
|---|---|---|
| `ozgarmas/` | konstantaga yangi qiymat — ataylab tarjima boʻlmaydi | `cannot assign to soniyalar (neither addressable nor a map index expression)` |
| `suv/` | bitta turi yoʻq konstanta `int8` va `float64` katakka tushadi | `12 12.5` |
| `suv-xato/` | `int` oʻzgaruvchi esa `float64` katakka tushmaydi — ataylab xato | `cannot use n (variable of type int) as float64 value in variable declaration` |
| `kutish/` | `time.Duration` katagidagi `2` → `2ns`; `2 * time.Second` → `2s` | `2ns` / `2s` |
| `sleep2/` | `time.Sleep(2)` haqiqatan ham deyarli kutmaydi | `true` / `2ns` |
| `sinov/` | oʻrtadagi sinov: `const narx = 300` → `int` ✅, `int8` ❌ | `cannot use narx (untyped int constant 300) as int8 value in variable declaration (overflows)` |
| `holat/` | `iota`: `Nomalum`, `Yangi`, `Tolangan` | `1 2` |
| `xato/` | xato: `Bekor` roʻyxat oʻrtasiga qoʻshildi, saqlangan `2` maʼnosini oʻzgartirdi | `false true` |
| `tuzatish/` | raqamlar qoʻlda, yangi holat oxirida | `true false` |
| `big/`, `float/` | Chuqurroq: `1 << 100 >> 98`, `0.1 + 0.2` | `4` · `true` / `false` / `0.30000000000000004` |
| `uyga/` | uy vazifasi: boshlangʻich fayllar va `javob/` | |

Qolgan papkalar (`huge/`, `edge/`, `limit/`, `maxuint/`, `maxfix/`, `pi/`, `duration/`, `durfix/`, `status/`,
`sizes/`, `bug/`, `fix/`, `homework/`) darsning birinchi versiyasidan qoldi, ular pastdagi **Ilova**da.

## Videodagi buyruqlar
```bash
go run ozgarmas/main.go    # ozgarmas/main.go:7:2: cannot assign to soniyalar (neither addressable nor a map index expression)
go run suv/main.go         # 12 12.5
go run kutish/main.go      # 2ns / 2s (keyin 2 soniya kutadi)
go run sinov/main.go       # sinov/main.go:8:15: cannot use narx (untyped int constant 300) as int8 value in variable declaration (overflows)
go run holat/main.go       # 1 2
go run xato/main.go        # false true
go run tuzatish/main.go    # true false
go run big/main.go         # 4
go run float/main.go       # true / false / 0.30000000000000004
```

## Dars matni (qisqa)

**Hook.** Sinf jurnalida baholar ism bilan emas, tartib raqami bilan yozilgan: uchinchi oʻquvchiga besh. Roʻyxat
oʻrtasiga yangi oʻquvchi qoʻshilsa, oʻsha besh baho endi kimniki? Go dasturlarida ham shunday jim xato bor.

**Eslash.** 02-darsda `var x int8 = 128` tarjima boʻlmadi: 128 kodda tayyor yozilgan son, Go uni oldindan tekshiradi.
Kodda tayyor yozilgan, oʻzgarmaydigan son — *konstanta*
([Go qoidasi](https://go.dev/ref/spec#Constants): literallar ham konstanta).

**Konstanta (`const`).** Oʻzgaruvchi — stadion tablosi, son almashadi. Konstanta — toshga oʻyilgan yozuv:
`const soniyalar = 60`. Unga yangi qiymat berish (`soniyalar = 61`) tarjima xatosi: `cannot assign to soniyalar`
([Assignment statements](https://go.dev/ref/spec#Assignment_statements): chap tomon oʻzgaruvchi boʻlishi kerak).

**Turi yoʻq konstanta (untyped constant).** Tur yozilmagan konstantaning turi u katakka tushganda aniqlanadi — suv
idish shaklini olgandek: `const yuz = 100` ham `int8`, ham `float64` katakka tushadi (`suv/`), `n := 100` esa `int`
oʻzgaruvchi va `float64` katakka tushmaydi (`suv-xato/`). Idish kichik boʻlsa, Go suvni toʻkmaydi: `const narx = 300`
`int8` ga tarjima xatosi (`sinov/`). Katak aytilmasa (masalan `fmt.Println(yuz)`), Go *default* turni tanlaydi: butun
son uchun `int`. Batafsil: Rob Pike, [*Constants*](https://go.dev/blog/constants).

**`time.Sleep(2)`.** `time.Sleep` qavs ichidagi sonni `time.Duration` katagiga qoʻyadi, u vaqtni **nanosekund**da
sanaydi ([time.go](https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L915)). Turi yoʻq `2` jim moslashadi:
dastur ikki soniya emas, ikki nanosekund soʻraydi — `Sleep` kamida shuncha kutadi, bu kompyuterda 1000 marta oʻlchaganda
0,4–27 mikrosekund ([pkg.go.dev/time#Sleep](https://pkg.go.dev/time#Sleep)). Ikki soniya: `2 * time.Second`.
Kompilyator ham, `go vet` ham jim; [staticcheck SA1004](https://staticcheck.dev/docs/checks/#SA1004) buni ushlaydi.
Teskari holat: `n := 2; n * time.Second` tarjima boʻlmaydi (`int` va `time.Duration` mos emas, `duration/`),
yechim `time.Duration(n) * time.Second` (`durfix/`).

**`iota`.** `const ( … )` ichida `iota` = qator raqami, 0 dan boshlanadi. Faqat nom yozilgan qatorda Go oldingi
yozuvni (bu yerda `iota`ni) takrorlaydi, `iota` har qatorda bittaga oshadi
([Iota](https://go.dev/ref/spec#Iota), [Constant declarations](https://go.dev/ref/spec#Constant_declarations)).
`Nomalum = 0`: qiymat berilmagan katakda nol turadi (zero value), shuning uchun nol maʼnoli holat boʻlmasligi kerak
(Protobuf style guide ham birinchi qiymatni `_UNSPECIFIED` qiladi).

**Koʻp uchraydigan xato.** Kecha `Tolangan = 2` faylga `2` deb yozildi. Bugun `Bekor` `Yangi`dan keyin qoʻshildi:
endi `2 = Bekor`, `Tolangan = 3`. `saqlangan == Tolangan` → `false`, `== Bekor` → `true`, xato chiqmaydi (`xato/`).
Xulosa: faylga yoki bazaga yozilgan raqamning maʼnosi hech qachon oʻzgarmasin. Saqlanadigan holatlarga raqamni qoʻlda
yozing (`tuzatish/`); `iota` qoldirsangiz, yangi holatni faqat oxiriga qoʻshing va hech birini oʻchirmang (oʻchirilgan
oʻrniga `_` qoldiring). Protobuf ham oʻchirilgan raqamlarni
[qayta ishlatmaslikni](https://protobuf.dev/best-practices/dos-donts/#reserve-deleted-numbers) talab qiladi.

**Chuqurroq (ixtiyoriy).** Go qoidasi: butun konstanta kamida 256 bit bilan saqlanadi. `const big = 1 << 100` (2¹⁰⁰,
101 bit) hech bir *butun son* turiga sigʻmaydi, lekin konstanta sifatida yashaydi; `big >> 98` = 2¹⁰⁰ / 2⁹⁸ = `4`
([Constant expressions](https://go.dev/ref/spec#Constant_expressions) — aynan shu misol). Konstantalar aniq
hisoblanadi: `0.1 + 0.2 == 0.3` konstantalar bilan `true`, `float64` oʻzgaruvchilar bilan `false`, yigʻindi
`0.30000000000000004` (0.1 va 0.2 ikkilik sanoqda aniq yozilmaydi va yaxlitlanadi).

## Uy vazifasi
- 🟢 **Oson** (`uyga/oson/main.go`, videodagi `xato` dasturi): `Bekor` qatorini roʻyxat oxiriga koʻchiring.
  Ishga tushirishdan oldin natijani taxmin qiling.
- 🟡 **Oʻrta** (`uyga/orta/main.go`): hafta kunlarini `iota` bilan yozing, `0` ni `Nomalum` (nomaʼlum kun) ga
  qoldiring. `Juma` necha?
- Qiyinroq vazifalar (bit flaglar, `1.5 * time.Second`, `time.Month` nega 1 dan): [homework/](homework/).

<details><summary>Javoblar</summary>

🟢 `true false`: `Tolangan` yana `2`, `Bekor` = `3`. Eski raqamlar joyida qoldi. [`uyga/javob/oson`](uyga/javob/oson/main.go)

🟡 [`uyga/javob/orta`](uyga/javob/orta/main.go) → `1 5 7`: `Dushanba` 1, `Juma` 5, `Yakshanba` 7.
</details>

## Manbalar
- Go qoidalari (spec): [Constants](https://go.dev/ref/spec#Constants), [Constant expressions](https://go.dev/ref/spec#Constant_expressions),
  [Iota](https://go.dev/ref/spec#Iota), [Constant declarations](https://go.dev/ref/spec#Constant_declarations),
  [Assignment statements](https://go.dev/ref/spec#Assignment_statements), [The zero value](https://go.dev/ref/spec#The_zero_value)
- Rob Pike, *Constants* (Go blog, 2014): https://go.dev/blog/constants
- `time.Duration`, `time.Second`: https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L915 · https://pkg.go.dev/time#Sleep
- Protobuf: [enum zero value](https://protobuf.dev/programming-guides/style/), [oʻchirilgan raqamlar](https://protobuf.dev/best-practices/dos-donts/#reserve-deleted-numbers)
- staticcheck SA1004: https://staticcheck.dev/docs/checks/#SA1004

---

## Ilova: darsning birinchi versiyasi (2026-11-25)

Video 2026-11-26 da qayta yozildi (Go 0 mavsumidan keyin, noldan boshlovchilar uchun). Birinchi versiyaning matni va
dasturlari shu yerda qoldi: gc kompilyatorining 512 bitlik chegarasi, assembly, default tur tuzogʻi
(`math.MaxUint64`), `math.Pi`, `n * time.Second`, `type Status int`, KB/MB/GB va Protobuf.

| Papka | Nima |
|---|---|
| `big/` | `const big = 1 << 100`, `big >> 98` → 4 |
| `huge/`, `edge/` | gc kompilyatorining chegarasi: `1 << 511` ishlaydi, `1 << 512` ataylab kompilyatsiya bo'lmaydi |
| `limit/`, `maxuint/`, `maxfix/` | 1-tuzoq: untyped konstantaning default turi `int` (`limit`, `maxuint` ataylab xato) |
| `pi/` | `math.Pi` → `float32` va `float64` |
| `float/` | 2-tuzoq: `0.1 + 0.2 == 0.3` ham true, ham false |
| `duration/`, `durfix/`, `sleep2/` | 3-tuzoq: `n * time.Second` (`duration` ataylab xato) va `time.Sleep(2)` = 2 ns |
| `status/`, `sizes/` | `iota`: enum va KB/MB/GB/TB |
| `bug/`, `fix/` | production tuzog'i: `iota` qayta raqamlaydi |
| `homework/` | uy vazifasi javoblari |

### Videodagi buyruqlar
```bash
go run ./big        # 4
go build -gcflags=-S -o /dev/null ./big 2>&1 | grep -A1 'stmp_0 SRODATA'   # faqat 8 bayt: 04 00 ... 00
go run ./huge       # 2
go run ./edge       # constant shift overflow
go run ./limit      # cannot use big (untyped int constant 1267650600228229401496703205376) as int value ... (overflows)
go run ./maxuint    # cannot use math.MaxUint64 (untyped int constant 18446744073709551615) as int value ... (overflows)
go run ./maxfix     # 18446744073709551615
go run ./pi         # float32: 3.14159274... / float64: 3.14159265358979311600
go run ./float      # true / false / 0.30000000000000004
go run ./duration   # invalid operation: n * time.Second (mismatched types int and time.Duration)
go run ./durfix     # 2s
go run ./sleep2     # true / 2ns
go run ./status     # 0 1 2 3 / true
go run ./sizes      # 1024 1048576 1073741824 1099511627776
go run ./bug        # to'langanmi? false / bekor qilinganmi? true
go run ./fix        # to'langanmi? true
```

### Dars matni

**Hook.** `1 << 100` bu 101 bitlik son, Go'dagi eng katta butun tur esa 64 bit. Bunday sonni konstanta qilib yozsak,
dastur kompilyatsiya bo'ladimi? (Aynan shu misol [Go spec](https://go.dev/ref/spec#Constant_expressions)da bor,
[Go Tour](https://go.dev/tour/basics/16)da esa o'xshashi: `Big >> 99`.)

**Pauza-savol.** `const big = 1 << 100; fmt.Println(big >> 98)` kompilyatsiya bo'ladimi? Ha, va `4` chiqaradi:
2¹⁰⁰ / 2⁹⁸ = 4. Katta son dastur xotirasiga hech qachon tushmadi: `-gcflags=-S` chiqishida `main` uchun faqat bitta 8 baytlik qiymat bor,
`04 00 00 00 00 00 00 00`.

**Untyped konstanta.** `big`ning turi yo'q, u oddiy matematik son kabi yashaydi.
[Spec](https://go.dev/ref/spec#Constants): raqamli konstantalar aniq qiymatga ega va overflow bo'lmaydi; har bir
kompilyator butun konstantalarni **kamida 256 bit** bilan saqlashi shart. Untyped konstanta turga sig'ishi faqat tur olganda tekshiriladi
(kompilyatorning 512 bitlik chegarasi bundan mustasno).
*Typed* konstanta esa har bir amaldan keyin o'z turiga sig'ishi shart. Oldingi darsdagi `int8(127) + 1` shuning uchun
xato edi: `int8(127)` typed konstanta, `+ 1` dan keyin 128 `int8`ga sig'maydi.

**Amaldagi chegara.** Go'ning o'z kompilyatori butun konstantalarni 512 bitgacha saqlaydi
([`types2/const.go`: `const prec = 512`](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/types2/const.go#L41)):
`1 << 511` ishlaydi, `1 << 512` → `constant shift overflow`.

**Standart library.** [`math.Pi`](https://github.com/golang/go/blob/go1.26.8/src/math/const.go#L13) 63 ta raqam bilan
yozilgan, `float64` esa ~16 raqamni saqlaydi. `Pi` untyped bo'lgani uchun, aniqlik faqat oxirgi qadamda, turga
yozilganda yo'qoladi: `float32`da 7 ta raqam to'g'ri, `float64`da 16 ta (`pi/`).

**1-tuzoq: default tur.** `fmt.Println(big)`: `Println`ga tur kerak, untyped butun konstantaning default turi `int`,
101 bitlik son `int`ga sig'maydi. Xuddi shu `fmt.Println(math.MaxUint64)` bilan ham bo'ladi. Yechim:
`uint64(math.MaxUint64)`. `big` esa hech bir butun turga sig'maydi: `float64(big)` mumkin, runtime'da aniq katta son
kerak bo'lsa, `math/big` package bor.

**2-tuzoq: `0.1 + 0.2`.** Konstantalar bilan `0.1 + 0.2 == 0.3` → `true` (kompilyator aniq hisoblaydi). `float64`
o'zgaruvchilar bilan → `false`, yig'indi `0.30000000000000004`: 0.1 va 0.2 ikkilik sanoqda aniq yozilmaydi va
yaxlitlanadi. Batafsil: Rob Pike, [*Constants*](https://go.dev/blog/constants).

**3-tuzoq: `n * time.Second`.** `time.Second` `Duration` turidagi konstanta, `type Duration int64`
([time.go](https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L915)). `2 * time.Second`da `2` untyped,
u jimgina `Duration` bo'ladi. `n := 2` esa `int`, va Go turli turlarni aralashtirmaydi (02-dars). `n` `int64` bo'lsa
ham ishlamaydi: `Duration` `int64` emas, `int64` asosidagi alohida tur. Yechim:
`time.Duration(n) * time.Second` (time package hujjatidagi misol ham shu,
[time.go L932](https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L932)).

Teskari tomoni xavfliroq: `time.Sleep(2)` xatosiz kompilyatsiya bo'ladi. Untyped `2` `Duration`ga aylanadi, birligi
nanosekund: dastur 2 sekund emas, atigi 2 nanosekund so'raydi, ya'ni deyarli umuman kutmaydi (`Sleep` kamida shuncha kutadi, aniq vaqt OS'ga bog'liq) (`sleep2/`). Kompilyator
ham, `go vet` ham jim; [staticcheck SA1004](https://staticcheck.dev/docs/checks/#SA1004) esa buni ushlaydi.

**`iota`.** `const` blokidagi nechanchi e'lon (ConstSpec) ekanini bildiradi, har bir `const` blokida 0 dan
boshlanadi. Birinchi e'londa `iota` yozilsa, keyingilari oldingi ifoda va turni avtomatik takrorlaydi:
```go
type Status int

const (
	Unknown Status = iota // 0: qiymat berilmagan Status shu bo'ladi (zero value)
	New                   // 1
	Paid                  // 2
	Shipped               // 3
)
```
Nolinchi qiymat ma'noli holat bo'lmasligi kerak: qiymat berilmagan `Status` 0 ga teng. Protocol Buffers style guide ham
birinchi enum qiymatini `_UNSPECIFIED` qilishni tavsiya qiladi.
Go'da alohida enum turi yo'q: `Status` bor-yo'g'i `int` asosidagi tur, shuning uchun `Status(42)` ham qonuniy.
`_ = iota; KB = 1 << (10 * iota)` bilan KB, MB, GB, TB bir formuladan chiqadi
([Effective Go](https://go.dev/doc/effective_go#constants) misoli; aniqrog'i KiB, MiB, …).

**Production tuzog'i.** `Status` bazaga son sifatida yozilgan, `Paid = 2`. Kimdir `Cancelled`ni `New`dan keyin
qo'shadi: undan keyingi hamma qiymat bittaga suriladi, endi `2` = `Cancelled`, `3` = `Paid`. To'langan buyurtmalar
bekor qilingan, jo'natilganlari esa to'langan bo'lib ko'rinadi, xato chiqmaydi. Qoida: bazaga yoki tarmoqqa chiqqan
raqam hech qachon o'zgarmaydi va qayta ishlatilmaydi. Shuning uchun bunday qiymatlarga raqamni qo'lda yozing (`fix/`),
yoki bazaga son emas, holat nomini yozing. O'chirilgan qiymat o'rnida `_` qoldiring.
Protobuf'da ham har bir enum qiymatiga raqam qo'lda yoziladi, va
[qo'llanma](https://protobuf.dev/best-practices/dos-donts/#reserve-deleted-numbers) o'chirilgan raqamlarni qayta
ishlatmaslikni qat'iy tavsiya qiladi.

### Manbalar
- Go spec: [Constants](https://go.dev/ref/spec#Constants), [Iota](https://go.dev/ref/spec#Iota),
  [Constant declarations](https://go.dev/ref/spec#Constant_declarations)
- Rob Pike, *Constants* (Go blog, 2014): https://go.dev/blog/constants
- gc konstanta aniqligi (512 bit): https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/types2/const.go#L41
- `math/const.go`: https://github.com/golang/go/blob/go1.26.8/src/math/const.go
- `time.Duration` va konstantalar: https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L915
- Protobuf: o'chirilgan enum raqamlari (https://protobuf.dev/best-practices/dos-donts/#reserve-deleted-numbers), enum
  zero value (https://protobuf.dev/programming-guides/style/)
- A Tour of Go, Numeric Constants: https://go.dev/tour/basics/16 · Effective Go, Constants: https://go.dev/doc/effective_go#constants

### Uy vazifasi
- 🟢 **Oson:** hafta kunlari uchun `iota` bilan konstantalar yozing.
- 🟡 **O'rta:** `iota` bilan bir nechta bit flag yozib, ularni bitli YOKI (`|`) bilan birlashtiring.
- 🟠 **Qiyin:** nega `1.5 * time.Second` kompilyatsiya bo'lmaydi, `1500 * time.Millisecond` esa ishlaydi?
- 🔴 **Guru:** `time` package'ida `iota` qayerlarda ishlatilganini toping va nega `Month` 1 dan boshlanishini tushuntiring.

Avval o'zingiz urinib ko'ring. Javoblar: [homework/](homework/).
