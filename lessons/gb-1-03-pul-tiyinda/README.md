# GB-1.03 · Butun sonlar va pul: narx tiyinda, `float64` nega emas

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi «Golang backend #3» (kanal playlist'ida)
· yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

Oldingi papka: [`gb-1-02-ozgaruvchi-tur`](../gb-1-02-ozgaruvchi-tur) · keyingi: GB-1.04 (matn va `fmt.Printf`).

## Bu darsda

- `int32` 4 bayt, `int64` 8 bayt; `int` bu kompyuterda (linux/amd64) 8 bayt, `GOARCH=386` bilan 4 bayt.
- `int32`'da 1 500 000 000 × 2 → `-1294967296`: overflow, xatosiz va panic'siz (Go spec: *Integer overflow*).
- `12000.0 * 1.12` → `13440.000000000002`: `float64` ikkilik kasr, 1.12 xotirada `1.12000000000000010658…`.
- Pul eng kichik birlikda, butun son bilan: narx tiyinda, `int64` (Payme Merchant API ham `amount`'ni tiyinda oladi).
- `int64` va `int` aralashsa: `invalid operation … (mismatched types int64 and int)`; tuzatish `int64(sold)`.
- `const` (funksiyadan tashqarida ham) va untyped constant: `subtotal * vatPercent` conversion'siz ishlaydi, `x := 100` esa `int`.
- Butun boʻlish `/` kasrni tashlaydi, `%` qoldiq: soʻm va tiyin.
- Klassik xato: `subtotal * (vatPercent / 100)` → QQS 0 (12 / 100 = 0).
- Chuqurroq: `int64(19.99 * 100)` → `1998`: float → int conversion kasrni tashlaydi.

## bozor qayerda

`main.go` GB-1.02 yechimidan boshlanadi (`diff -ru ../gb-1-02-ozgaruvchi-tur/solution .`). Oʻzgarishlar:
`var price int64` (avval `int`), `price = 12_000 * tiyinPerSom` (tiyinda), fayl boshida ikki `const`, oxirida sotuv
summasi: `subtotal`, `vat`, `total` va ikki yangi `Println`. `weight` (ogʻirlik, pul emas) `float64` boʻlib qoladi.

## Buyruqlar va haqiqiy natijalar (go1.26.8 linux/amd64, CGO_ENABLED=0)

```console
$ go run ./misollar/olcham
int32: 4 bayt, -2147483648 … 2147483647
int64: 8 bayt, -9223372036854775808 … 9223372036854775807
int:   8 bayt (bu kompyuterda)

$ GOARCH=386 go run ./misollar/olcham
…
int:   4 bayt (bu kompyuterda)

$ go run ./misollar/int32
-1294967296

$ go run ./misollar/float
13440.000000000002

$ go run ./misollar/kesish
1998

$ go run .
bozor ishga tushdi
versiya: 0.1.0
 0 false
olma 1200000 true 40 1.5
olma 37
olma 3 dona
jami: 40320 soʻm 0 tiyin
```

Videoda `main.go`'ning oraliq holatlari ham bor. Ularni qayta koʻrish uchun `main.go`'ni vaqtincha oʻzgartiring:

1. **Turlar aralashsa** (`const`'lar hali yoʻq). GB-1.02 yechimida `var price int` → `var price int64`,
   `price = 12000` → `price = 1_200_000 // tiyin`, `fmt.Println(name, stock)`'dan keyin:
   ```go
   	subtotal := price * sold
   	fmt.Println(subtotal)
   ```
   ```console
   $ go run .
   # github.com/abdurahmon27/bozor
   ./main.go:25:14: invalid operation: price * sold (mismatched types int64 and int)
   ```
   `price * int64(sold)` bilan oxirgi qator: `3600000`.
2. **QQS nolga aylanadi.** `vat := subtotal * vatPercent / 100` oʻrniga `vat := subtotal * (vatPercent / 100)`:
   oxirgi qator `jami: 36000 soʻm 0 tiyin`.
3. **Tekshiruv savoli.** `12_000` → `1_999`, `sold := 3` → `sold := 1`: oxirgi qatorlar `olma 1 dona`,
   `jami: 2238 soʻm 88 tiyin`.

`misollar/*` — har biri oʻz `package main`'iga ega alohida dastur (bitta modul ichida). `misollar/olcham` kodida
darsda hali oʻtilmagan `math` va `unsafe` paketlari bor; videoda faqat uning natijasi koʻrsatiladi.

## Vazifa

[TASK.md](TASK.md): choy, 1 999 soʻm, `sold := 7`, QQS 12 % va 15 000 soʻm yetkazib berish. Tekshirish:

```console
$ go test            # boshlangʻich holatda FAIL, oʻzbekcha xabar bilan
$ go test ./...      # vazifa bajarilgach hammasi ok
```

Javob: [`solution/`](solution) (`go run ./solution`).

## Manbalar

- Go spec: [Numeric types](https://go.dev/ref/spec#Numeric_types), [Integer literals](https://go.dev/ref/spec#Integer_literals),
  [Integer overflow](https://go.dev/ref/spec#Integer_overflow), [Integer operators](https://go.dev/ref/spec#Integer_operators),
  [Constants](https://go.dev/ref/spec#Constants), [Constant expressions](https://go.dev/ref/spec#Constant_expressions),
  [Conversions](https://go.dev/ref/spec#Conversions), [Operator precedence](https://go.dev/ref/spec#Operator_precedence),
  [Constant declarations](https://go.dev/ref/spec#Constant_declarations), [Arithmetic operators](https://go.dev/ref/spec#Arithmetic_operators)
- Payme Merchant API, CheckPerformTransaction: `amount` — «Сумма платежа (в тийинах)»:
  <https://developer.help.paycom.uz/metody-merchant-api/checktransaction>
- Stripe: amounts in the smallest currency unit: <https://docs.stripe.com/currencies>
- Oʻzbekistonda QQS 12 % (2023-yil 1-yanvardan): <https://kun.uz/en/84167791>
