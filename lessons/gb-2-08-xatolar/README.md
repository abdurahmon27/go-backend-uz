# GB-2.08 · Error'lar chuqur: `%w`, `errors.Is` va `errors.As`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #16"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. Darsdagi `errors.AsType` Go 1.26'dan beri bor.

Oldingi dars: [`gb-2-07-interface-writer`](../gb-2-07-interface-writer) (interface: `io.Writer` va ikki haqiqiy
implementatsiya) · keyingi dars: `gb-2-09` (`defer`, `panic` va `recover`: har biri nima uchun).
`store` paketi va `Sell`: [`gb-2-06-paketlar`](../gb-2-06-paketlar).

**Bu darsdan keyin:** `error` — `Error() string` metodli interface ekanini bilasiz; `errors.New` bilan sentinel
xato (`ErrNotFound`) eʼlon qilib, uni `fmt.Errorf`'dagi `%w` bilan kontekst qoʻshib oʻraysiz va chaqiruvchida
`errors.Is` bilan topasiz; `%w` oʻrniga `%v` yozilsa zanjir nega uzilishini koʻrasiz; maydonli oʻz xato turingizni
(`*StockError`, typed error) yozib, uni `errors.As` (yoki `errors.AsType`) bilan zanjirdan olasiz.

Bu papka 2-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i bor. Tuzilishi GB-2.06'dagidek
(`cmd/bozor` + `internal/store`); `Product`, `Store`, `Add`, `List` GB-2.06 javobidan (izohlari qisqartirilgan),
`Sell` oʻsha qoidalar bilan, lekin xatolari tekshiriladigan qilib qayta yozilgan (`p := &s.products[i]`). GB-2.06'dagi
`golang.org/x/text` (soʻmni guruhlab chiqarish) bu darsga kerak emas, shuning uchun tashqi modul yoʻq. 3-modul
oxiridagi bozor CLI v0 (GB-3.06) labs'ni birlashtiradi. GB-2.06 bilan farq:
`diff -ru ../gb-2-06-paketlar/solution/internal internal`.

## Papka

| Fayl | Nima |
|---|---|
| `go.mod` | modul `github.com/abdurahmon27/bozor`, `go 1.27` (tashqi modul yoʻq) |
| `internal/store/store.go` | `package store`: `ErrNotFound`, `StockError`, `Product`, `Store`, `Add`, `List`, `Sell`; vazifaning `ErrBadQuantity`'si (`Sell` uni hali qaytarmaydi) |
| `cmd/bozor/main.go` | dastur: `Sell` xatolarini `errors.Is` va `errors.As` bilan tekshiradi; vazifaning `message` funksiyasi (chala) |
| `internal/store/tekshiruv_test.go`, `cmd/bozor/tekshiruv_test.go` | vazifani tekshiradi (`go test ./...`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: xuddi shu daraxt, oʻz `go.mod`'i bilan (alohida modul, shuning uchun ildizdagi `go test ./...` unga tegmaydi) |
| `_video/NN-…/` | videodagi har bir qadam (toʻliq kichik loyiha); `_` bilan boshlangan papkani `go` buyruqlari chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 03-wrap go run ./cmd/bozor` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Boshlanishi (`00-matn`, GB-2.06 javobidagi `Sell`): ikki xil xato, lekin `main` uchun ikkalasi faqat matn:

```
$ go run ./cmd/bozor
mahsulot 9 topilmadi
olma: qoldiq 40, soʻralgan 41
```

`error`'ning oʻzi (hujjatdan faqat tur qoldirildi) va `errors.New` Go manbasida
(`$(go env GOROOT)/src/errors/errors.go`, 64–75-qatorlar, izoh qatori qisqartirilgan):

```
$ go doc builtin.error
type error interface {
	Error() string
}
```

```go
func New(text string) error {
	return &errorString{text}
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
```

Sentinel: `var ErrNotFound = errors.New("mahsulot topilmadi")`, `Sell` uni oʻzini qaytaradi, `main` `==` bilan
solishtiradi (`01-sentinel`). Xuddi shu matnli yangi `errors.New` esa teng emas (`02-yangi`):

```
$ go run ./cmd/bozor
bunday mahsulot yoʻq
false
```

`%w` bilan ID qoʻshiladi (`03-wrap`: `fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)`), `main` hali `==` bilan:

```
$ go run ./cmd/bozor
xato: mahsulot 9: mahsulot topilmadi
```

Nega `==` yolgʻon: interface ichidagi tur boshqa (`04-tur`, `fmt.Printf("%T\n", …)`):

```
$ go run ./cmd/bozor
xato: mahsulot 9: mahsulot topilmadi
*fmt.wrapError
*errors.errorString
```

`errors.Is(err, store.ErrNotFound)` (`05-is`) zanjirni `Unwrap` bilan pastga tushib tekshiradi
(`main`'da endi `if … else if err != nil`):

```
$ go run ./cmd/bozor
bunday mahsulot yoʻq
```

Tez-tez uchraydigan xato (`06-v`): `%w` oʻrniga `%v`. Matn harfma-harf oʻsha, lekin `errors.Is` — `false`:

```
$ go run ./cmd/bozor
xato: mahsulot 9: mahsulot topilmadi
```

Savol (`07-quiz`): `errors.Is` matnni emas, qiymatni solishtiradi — xuddi shu matnli yangi `errors.New` topilmaydi:

```
$ go run ./cmd/bozor
bunday mahsulot yoʻq
false
```

Typed error `*StockError{Name, Stock, Want}`; `main` xatoni `fmt.Errorf("buyurtma 7: %w", err)` bilan yana oʻraydi.
Oddiy type assertion `err.(*store.StockError)` faqat tashqi qavatga qaraydi (`08-assert`):

```
$ go run ./cmd/bozor
bunday mahsulot yoʻq
buyurtma 7: olma: qoldiq 40, soʻralgan 41
false
```

`errors.As(err, &se)` esa zanjir boʻylab topadi (`09-as`; uning `main`'i papkadagi `main` bilan bir xil), xuddi shu
ish `errors.AsType[*store.StockError](err)` bilan (`10-astype`); natija bir xil:

```
$ go run ./cmd/bozor
bunday mahsulot yoʻq
olma: faqat 40 dona qoldi
```

Chuqurroq: `errors.Join` (`11-join`) — nil xatolar tashlanadi, `errors.Is` har bir qismini tekshiradi:

```
$ go run ./cmd/bozor
mahsulot 9: mahsulot topilmadi
olma: qoldiq 40, soʻralgan 41
true
```

## Vazifa

[TASK.md](TASK.md): `Sell` musbat boʻlmagan dona sonida `ErrBadQuantity`'ni `%w` bilan oʻrab qaytarsin;
`message(err)` xatoni `errors.Is`/`errors.As` bilan foydalanuvchi matniga aylantirsin. Tekshirish:
`go test ./...` (hozir FAIL, vazifa bajarilgach ikkala paket `ok`). Javobni tekshirish: `(cd solution && go test ./...)`.

## Manbalar

- Go spec: Errors — https://go.dev/ref/spec#Errors · Comparison operators (interface qiymatlarini solishtirish) — https://go.dev/ref/spec#Comparison_operators
- pkg.go.dev `errors` (`New`, `Is`, `As`, `AsType`, `Join`) — https://pkg.go.dev/errors
- pkg.go.dev `fmt.Errorf` (`%w`) — https://pkg.go.dev/fmt#Errorf · `%T` — https://pkg.go.dev/fmt#hdr-Printing
- Go manbasi: `src/errors/errors.go`, `src/fmt/errors.go` (go1.27.4) — https://cs.opensource.google/go/go/+/refs/tags/go1.27.4:src/errors/
- Go blog: Working with Errors in Go 1.13 (D. Neil, J. Amsterdam) — https://go.dev/blog/go1.13-errors
- Go blog: Errors are values (R. Pike) — https://go.dev/blog/errors-are-values
- Release notes: Go 1.13 (wrapping) — https://go.dev/doc/go1.13#error_wrapping · Go 1.20 (`Join`) — https://go.dev/doc/go1.20#errors · Go 1.26 (`AsType`) — https://go.dev/doc/go1.26#errorspkgerrors
- Go Code Review Comments: Error Strings — https://go.dev/wiki/CodeReviewComments#error-strings
