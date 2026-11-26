# 03 · Konstanta va `iota`: JS `const` emas, `time.Sleep(2)` va enum oʻrtasidagi jim bug

Video: *Golang oʻzbek tilida #3* · Lab: https://haywan.uz/lab/youtube/go#go-1-03 ·
Oldingi dars: [02 · Turlar va overflow](../002-types-overflow/) · Keyingi dars: [04 · `len("oʻzbek") = 7?` Bayt, rune va UTF-8](../004-strings-utf8/)

Boshqa tilda (JS/TS, Python, PHP) ishlaydigan, Go'ga yangi dasturchilar uchun. Avval "Go: JS/Python/PHP dasturchilari
uchun" kursi (GO-X: `go run`, statik turlar, yashirin konversiya yoʻqligi, `:=` → `int`, zero value, `error`) va 02-dars.
Hamma natijalar shu kompyuterda olingan (2026-11-26): `go1.26.8 linux/amd64`, Node v24.10.0, Python 3.13.15,
PHP 8.4.26, tsc 5.9.3. Go dasturlarini https://go.dev/play da ham ishga tushirish mumkin.

**Bitta gʻoya:** Go'da `const` — kompilyatsiya paytida maʼlum qiymat, va tur yozilmagan konstanta (*untyped constant*)
turini ishlatilgan joydan oladi: `time.Sleep(2)` shuning uchun ikki nanosekund soʻraydi. `iota` esa qiymatini `const`
blokidagi eʼlon oʻrnidan oladi, shuning uchun roʻyxat oʻrtasiga qoʻshilgan holat bazadagi raqamlarning maʼnosini jim
oʻzgartiradi.

| Papka | Nima | Natija |
|---|---|---|
| `eslash/` | 02-darsdan: `int8` 127 + 3 | `-126` |
| `anchors/` | JS / Python / PHP / TS solishtirishlar (`const.js`, `const.py`, `const.php`, `sleep.py`, `sleep.js`, `literal.ts`, `alias.ts`, `enum1.ts`, `enum2.ts`) | pastda |
| `notconst/` | `const maxSize = 10 * 1024 * 1024` ✅, `math.Sqrt(2)` ❌, `time.Now()` ❌ | `… is not constant` |
| `untyped-xato/` | `int` oʻzgaruvchi `float64` ga oʻtmaydi — ataylab xato | `cannot use n (variable of type int) as float64 value …` |
| `untyped/` | bitta untyped constant → `int8` va `float64`; `:=` → default type `int` | `12 12.5` / `int` |
| `duration-type/` | `type Duration int64` alohida tur: `int` oʻzgaruvchi oʻtmaydi — ataylab xato | `cannot use n (variable of type int) as time.Duration value …` |
| `sleep2/` | `time.Sleep(2)` | `2ns true` |
| `kutish/` | `2` → `2ns`, `2 * time.Second` → `2s` (keyin 2 soniya kutadi) | `2ns` / `2s` |
| `config/` | `n := 2; time.Sleep(n * time.Second)` — ataylab xato | `mismatched types int and time.Duration` |
| `config-fix/`, `wrongfix/` | `time.Duration(n) * time.Second` ✅ · `time.Duration(n)` ❌ (2ns) | `2s` · `2s` / `2ns` |
| `sinov/` | oʻrtadagi sinov: untyped `a` ✅, typed `const b int` ❌ | `b * time.Second (mismatched types …)` |
| `holat/` | `type Holat int` + `iota`; zero value = `Nomalum` | `1 2 true` |
| `xato/` | xato: `Bekor` roʻyxat oʻrtasiga qoʻshildi, saqlangan `2` maʼnosini oʻzgartirdi | `false true` |
| `tuzatish/` | raqamlar qoʻlda, yangi holat oxirida | `true false` |
| `ochirish/` | takror savoli: `Yangi` oʻchirilsa `2` = `Bekor`; `_` bilan raqam band qoladi | `false true` / `true false` |
| `float/` | Chuqurroq: `0.1 + 0.2 == 0.3` — untyped / typed `float64` / oʻzgaruvchi | `true` / `false` / `false 0.30000000000000004` |
| `addr/` | qoʻshimcha: `&Port` — konstantaning manzili yoʻq (Go qoidasi: addressable emas) | `cannot take address of Port (untyped int constant 8080)` |
| `uyga/` | uy vazifasi: boshlangʻich fayllar va `javob/` | |

Qolgan papkalar (`ozgarmas/`, `big/`, `huge/`, `edge/`, `limit/`, `maxuint/`, `maxfix/`, `pi/`, `duration/`, `durfix/`,
`status/`, `sizes/`, `bug/`, `fix/`, `homework/`) darsning oldingi versiyalaridan qoldi (pastdagi **Ilova**).

## Videodagi buyruqlar
```bash
go run eslash/main.go        # -126
node -p 'const start = Date.now(); typeof start'   # number
go run notconst/main.go      # notconst/main.go:11:16: math.Sqrt(2) (value of type float64) is not constant
                             # notconst/main.go:12:16: time.Now() (value of struct type time.Time) is not constant
php anchors/const.php        # Fatal error: Constant expression contains invalid operations in …/anchors/const.php on line 5
                             # (PHP: nix shell nixpkgs#php -c php anchors/const.php)
tsc --strict anchors/literal.ts   # anchors/literal.ts(5,7): error TS2322: Type 'number' is not assignable to type '3'.
python3 anchors/sleep.py     # 2.0   (time.sleep(2) = 2 soniya)
node anchors/sleep.js        # 2     (setTimeout(…, 2) = 2 millisekund)
go run sleep2/main.go        # 2ns true
go run kutish/main.go        # 2ns / 2s
go run config/main.go        # config/main.go:10:13: invalid operation: n * time.Second (mismatched types int and time.Duration)
go run config-fix/main.go    # 2s
go run sinov/main.go         # sinov/main.go:12:7: invalid operation: b * time.Second (mismatched types int and time.Duration)
tsc anchors/enum1.ts && node anchors/enum1.js   # 1 2
go run xato/main.go          # false true
tsc anchors/enum2.ts && node anchors/enum2.js   # false true
go run tuzatish/main.go      # true false
go run float/main.go         # true / false / false 0.30000000000000004
```
(`tsc` 5.9.3: `npx -p typescript@5.9.3 tsc --strict --types --lib es2022,dom anchors/enum1.ts`; `.js` fayllar
`anchors/` ichida paydo boʻladi.)

## Dars matni (qisqa)

**Hook.** Bazada toʻlangan buyurtmalar `status = 2` bilan saqlangan. Kimdir Go'dagi holatlar roʻyxati oʻrtasiga yangi
holat qoʻshdi, kompilyator jim. Ertasi kuni oʻsha buyurtmalar qanday koʻrinadi? (Javob `xato/` da.)

**Eslash (02-dars).** `var x int8 = 127; x += 3` → `-126`: oʻzgaruvchi ishlash paytida jim aylanadi. Konstantani esa
Go kompilyatsiya paytida tekshiradi; 02-darsdagi xabarda *untyped int constant* degan soʻz bor edi.

**1. `const` — JS `const` emas.** JS'da `const` faqat nomni qayta bogʻlashni taqiqlaydi, qiymat ishlash paytida
hisoblanadi (`Date.now()`). Python'da konstanta — kelishuv (katta harfli nom, `MAX_RETRIES = 4` ham ishlaydi;
`typing.Final` ni faqat mypy kabi tekshiruvchi tekshiradi). Go'da `const` qiymati **kompilyatsiya paytida** maʼlum
boʻlishi shart va u faqat son (rune ham), string yoki bool boʻladi
([Constants](https://go.dev/ref/spec#Constants)). Konstanta ifodada konstantalar va operatorlar qatnashadi; oddiy
funksiya chaqiruvi emas (`len("abc")`, `min`/`max` kabi baʼzi built-in'lar konstantalar bilan istisno,
[Constant expressions](https://go.dev/ref/spec#Constant_expressions)). Shuning uchun `const root2 = math.Sqrt(2)` va
`const start = time.Now()` — `is not constant`. PHP'ning `const` i ham shunday: `const START = time();` →
`Fatal error: Constant expression contains invalid operations`, fayl umuman ishga tushmaydi. Ishlash paytidagi qiymat
uchun Go'da oddiy oʻzgaruvchi (`var` / `:=`).

**2. Untyped constant.** Go'da yashirin konversiya yoʻq (`untyped-xato/`). Tur yozilmagan konstantaning esa hali turi
yoʻq — *untyped constant*. U turini ishlatilgan joydan oladi (`int8`, `float64`); tur aytilmagan joyda (`:=`,
`fmt.Println`) *default type* oladi: butun son uchun `int`, kasr uchun `float64`, `'a'` uchun `rune`
([Constants](https://go.dev/ref/spec#Constants)). TypeScript'da oʻxshash gʻoya: `const a = 3` ning turi `3` (literal
type), `let b = a` da `number` ga kengayadi (`anchors/literal.ts`).

**`time.Duration`.** `func Sleep(d Duration)`, `type Duration int64` — nanosekundlar soni
([time.go](https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L915)). `type` Go'da **yangi tur** yaratadi:
TS'dagi `type Ms = number` faqat taxallus (`anchors/alias.ts`), Go'da esa `int` oʻzgaruvchini `Duration` ga qoʻyib
boʻlmaydi (`duration-type/`). Python `time.sleep(2)` = 2 soniya, JS `setTimeout(f, 2)` = 2 ms, Go `time.Sleep(2)` =
**2 nanosekund soʻraydi** (xatosiz: `2` untyped, `Duration` boʻladi). Oʻlchovda 1000 ta `time.Sleep(2)` 0,3–14 µs da
qaytdi (Sleep kamida *d* kutadi). Ikki soniya: `2 * time.Second`; `time.Second` — *typed constant*, turi `Duration`,
qiymati 1 000 000 000 ([time.go L934](https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L934)).
staticcheck [SA1004](https://staticcheck.dev/docs/checks/#SA1004) `time.Sleep(2)` kabi kichik **konstantani** ushlaydi,
lekin `time.Sleep(time.Duration(n))` ni ushlamaydi.

**Ishdagi holat: config.** `n := 2` (default type `int`) → `n * time.Second` — `mismatched types int and time.Duration`.
Tuzatish: `time.Duration(n) * time.Second` ([time.go L930](https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L930)
dagi rasmiy misol). Notoʻgʻri tuzatish: `time.Duration(n)` — yana 2 ns (`wrongfix/`). Config'da `"2s"` kabi matn boʻlsa,
`time.ParseDuration("2s")` (npm'dagi `ms('2s')` ga oʻxshash). Nega shunday: Go yashirin konversiyani ataylab olib
tashlagan ([FAQ](https://go.dev/doc/faq#conversions)), untyped constant esa literallarni qulay qoldiradi
([Rob Pike, Constants](https://go.dev/blog/constants)).

**Sinov.** `const a = 3` (untyped) → `a * time.Second` ✅; `const b int = 3` (typed) → `b * time.Second` ❌.

**3. `iota`.** TS `enum Status { Unknown, New, Paid }` → `New = 1`, `Paid = 2`. Go'da `enum` kalit soʻzi yoʻq:
`type Holat int` + `const` blokida `iota`. `iota` = blokdagi eʼlonning (ConstSpec) tartib raqami, 0 dan; izoh va boʻsh
qatorlar sanalmaydi, `_` esa raqam oladi. Faqat nom yozilgan eʼlon oldingi ifodani (`Holat = iota`) takrorlaydi
([Iota](https://go.dev/ref/spec#Iota), [Constant declarations](https://go.dev/ref/spec#Constant_declarations)).
Bu tur yopiq emas: `Holat(42)` ham kompilyatsiya boʻladi (barcha holatlar tekshirilishini `exhaustive` linter beradi).
`0` = `Nomalum`, chunki holat berilmagan oʻzgaruvchida zero value turadi (Protobuf'da ham birinchi qiymat
`…_UNSPECIFIED = 0`, [style guide](https://protobuf.dev/programming-guides/style/#enums)).

**Koʻp uchraydigan xato.** Kecha `Tolangan = 2` bazaga `2` deb yozildi. Bugun `Bekor` `Yangi`dan keyin qoʻshildi:
`Bekor = 2`, `Tolangan = 3`, kompilyator jim, saqlangan `2` endi «bekor» (`xato/` → `false true`; TS enum'da ham
`false true`). **Qoida:** bazaga yoki JSON'ga ketadigan raqamning maʼnosi oʻzgarmasin. (1) raqamni qoʻlda yozing
(`tuzatish/`), (2) `iota` qoldirsangiz, yangi holatni faqat oxiriga qoʻshing, oʻchirilgan holat oʻrniga `_` qoldiring
(`ochirish/`), (3) raqam emas, nom saqlang (`"paid"`; uy vazifasi 🟡). Protobuf ham raqamni qoʻlda yozdiradi va
oʻchirilgan raqamni `reserved` qiladi ([dos and don'ts](https://protobuf.dev/best-practices/dos-donts/#reserve-deleted-numbers)).
Keyinroq (metodlar darsi): `func (h Holat) String() string` yoki `stringer` vositasi `fmt.Println(Tolangan)` da nom chiqaradi.

**Chuqurroq.** JS `0.1 + 0.2 === 0.3` → `false`. Go'da **untyped** konstantalar bilan `true` (konstanta ifodalar aniq
hisoblanadi), `const c, d float64` (typed) bilan `false` (har biri `float64` ga yaxlitlanadi), `float64` oʻzgaruvchilar
bilan `false`, yigʻindi `0.30000000000000004` (`float/`). Qoʻshimcha: `&Port` — `cannot take address` (konstanta
addressable emas, `addr/`); Go 1.26 dan `new(Port)` ishlaydi (qiymat nusxasi bilan yangi oʻzgaruvchi).

## Uy vazifasi
- 🟢 **Oson** (`uyga/oson/main.go`: `Nomalum, Yangi, Tolangan, Bekor` — `iota` bilan, toʻgʻri tartibda): pul qaytarilgan buyurtmalar uchun `Qaytarilgan` holatini
  qoʻshing, shunda bazadagi `2` (`Tolangan`) va `3` (`Bekor`) maʼnosini saqlasin. Boshlangʻich fayl hozir kompilyatsiya
  boʻlmaydi (`undefined: Qaytarilgan`) — holatni qoʻshganingizda `true false false` chiqishi kerak.
- 🟡 **Oʻrta** (`uyga/orta/main.go`): bazaga raqam emas, nom saqlang. `HolatNomi(h Holat) string` va
  `HolatNomdan(s string) (Holat, error)` yozing; nomaʼlum nomga `error` qaytsin.
  Kutilgan natija: `paid 2 <nil>` va `0 nomaʼlum holat: "refund"`.
- Qiyinroq vazifalar (bit flaglar, `1.5 * time.Second`, `time.Month` nega 1 dan): [homework/](homework/).

<details><summary>Javoblar</summary>

🟢 `Qaytarilgan` roʻyxat **oxiriga** (`= 4`), eski raqamlar joyida → `true false false`. Qoʻlda yozilgan raqam bilan
ham boʻladi (`Qaytarilgan Holat = 4`). [`uyga/javob/oson`](uyga/javob/oson/main.go)

🟡 [`uyga/javob/orta`](uyga/javob/orta/main.go) → `paid 2 <nil>` / `0 nomaʼlum holat: "refund"`. Bazada `"paid"`
saqlansa, `const` blokidagi tartib oʻzgarsa ham maʼno buzilmaydi.
</details>

## Manbalar
- Go qoidalari (spec): [Constants](https://go.dev/ref/spec#Constants), [Constant expressions](https://go.dev/ref/spec#Constant_expressions),
  [Iota](https://go.dev/ref/spec#Iota), [Constant declarations](https://go.dev/ref/spec#Constant_declarations),
  [Address operators](https://go.dev/ref/spec#Address_operators), [The zero value](https://go.dev/ref/spec#The_zero_value)
- Rob Pike, *Constants* (Go blog, 2014): https://go.dev/blog/constants · Go FAQ: https://go.dev/doc/faq#conversions
- `time.Duration`, `time.Second`: https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L915 · https://pkg.go.dev/time#Sleep
- JS `const`: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/const ·
  `setTimeout`: https://developer.mozilla.org/en-US/docs/Web/API/Window/setTimeout
- Python `time.sleep`: https://docs.python.org/3/library/time.html#time.sleep · PEP 8 constants: https://peps.python.org/pep-0008/#constants
- PHP constants: https://www.php.net/manual/en/language.constants.syntax.php
- TypeScript enums: https://www.typescriptlang.org/docs/handbook/enums.html
- Protobuf: [enum zero value](https://protobuf.dev/programming-guides/style/#enums), [oʻchirilgan raqamlar](https://protobuf.dev/best-practices/dos-donts/#reserve-deleted-numbers)
- staticcheck SA1004: https://staticcheck.dev/docs/checks/#SA1004

---

## Ilova: ikkinchi versiya (2026-11-26, Go 0 tomoshabinlari uchun)

Bu versiya hayotiy oʻxshatishlar (sinf jurnali, tosh, suv) bilan yozilgan edi; video dasturchilar uchun qayta
yozilganda ular olib tashlandi. Dasturlari: `ozgarmas/` (`const soniyalar = 60; soniyalar = 61` →
`cannot assign to soniyalar`), `untyped/` va `untyped-xato/` (oldingi `suv/`, `suv-xato/`), `big/` (`1 << 100 >> 98` → `4`).

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
go run ./sleep2     # (oʻsha versiyada: true / 2ns; endi sleep2 `2ns true` chiqaradi)
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
