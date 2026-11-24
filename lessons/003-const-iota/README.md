# 03 · Konstantalar va `iota`: `1 << 100` qanday sig'adi?

Video: *Golang o'zbek tilida #3* · Lab: https://haywan.uz/lab/youtube/go#go-1-03 ·
Oldingi dars: [02 · Turlar va overflow](../002-types-overflow/) · Keyingi dars: 04 · `len("oʻzbek")` = 7? Byte, rune va UTF-8

Hamma natijalar `go1.26.8 linux/amd64`da olingan.

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

## Videodagi buyruqlar
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

## Dars matni

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

## Manbalar
- Go spec: [Constants](https://go.dev/ref/spec#Constants), [Iota](https://go.dev/ref/spec#Iota),
  [Constant declarations](https://go.dev/ref/spec#Constant_declarations)
- Rob Pike, *Constants* (Go blog, 2014): https://go.dev/blog/constants
- gc konstanta aniqligi (512 bit): https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/types2/const.go#L41
- `math/const.go`: https://github.com/golang/go/blob/go1.26.8/src/math/const.go
- `time.Duration` va konstantalar: https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L915
- Protobuf: o'chirilgan enum raqamlari (https://protobuf.dev/best-practices/dos-donts/#reserve-deleted-numbers), enum
  zero value (https://protobuf.dev/programming-guides/style/)
- A Tour of Go, Numeric Constants: https://go.dev/tour/basics/16 · Effective Go, Constants: https://go.dev/doc/effective_go#constants

## Uy vazifasi
- 🟢 **Oson:** hafta kunlari uchun `iota` bilan konstantalar yozing.
- 🟡 **O'rta:** `iota` bilan bir nechta bit flag yozib, ularni bitli YOKI (`|`) bilan birlashtiring.
- 🟠 **Qiyin:** nega `1.5 * time.Second` kompilyatsiya bo'lmaydi, `1500 * time.Millisecond` esa ishlaydi?
- 🔴 **Guru:** `time` package'ida `iota` qayerlarda ishlatilganini toping va nega `Month` 1 dan boshlanishini tushuntiring.

Avval o'zingiz urinib ko'ring. Javoblar: [homework/](homework/).
