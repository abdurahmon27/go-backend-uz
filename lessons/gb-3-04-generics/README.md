# GB-3.04 · Generics: tur parametri qachon foyda beradi

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #22"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** 2–3-modul darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. Generics tilda Go 1.18'dan beri bor,
> `errors.AsType` — Go 1.26'dan beri.

Oldingi dars: [`gb-3-03-jamoa-asboblari`](../gb-3-03-jamoa-asboblari) (`go vet`, `staticcheck`, `govulncheck`) ·
keyingi dars: `gb-3-05` (iterator'lar: `range` funksiya ustida, `slices` va `maps`). `io.Writer` va interface:
[`gb-2-07-interface-writer`](../gb-2-07-interface-writer); `errors.AsType` va `StockError`:
[`gb-2-08-xatolar`](../gb-2-08-xatolar); `Order`, `OrderStatus`, `Total()`:
[`gb-2-03-struct-turlar`](../gb-2-03-struct-turlar), [`gb-2-05-metodlar`](../gb-2-05-metodlar).

**Bu darsdan keyin:** bir xil kod faqat turi boshqa holda takrorlanganini koʻrib, uni bitta generic funksiyaga
yigʻasiz (`Filter[T any]`); nega `[]Product`'ni `[]any` parametriga berib boʻlmasligini (xotiradagi koʻrinish
boshqa) bilasiz; kompilyator `T`'ni argumentlardan qanday chiqarishini (type inference) va `Filter[Product]`
yozuvi (instantiation) oddiy funksiya turini berishini `%T` bilan koʻrasiz; constraint tur parametri bilan
qaysi amallar mumkinligini belgilashini (`any` bilan `+` yoʻq, `int | int64` bilan bor) kompilyator xabarlaridan
oʻqiysiz; qiymatning faqat metodi chaqiriladigan joyda (`writeTotal(w io.Writer, …)`) tur parametri ortiqcha
ekanini bilasiz; `errors.AsType[E error]`'dagi kvadrat qavsni nega yozish shartligini (`cannot infer E`)
tushunasiz.

Bu papka 2–3-modul qoidasi boʻyicha **mustaqil lab**: oʻz `go.mod`'i bor, bitta `package main`.
`Product`, `OrderItem`, `Order`, `OrderStatus` (GB-2.03) va `(*Order).Total` (GB-2.05) oʻsha darslardan
(`OrderStatus` faqat `StatusNew`, `StatusPaid` bilan; `String()` bu labga kerak emas). Yangi: `Filter`, `Number`,
`Sum`, `Contains`, `writeTotal`, vazifaning `SumBy` va `paidTotal`'i; chuqurroq qismda `type Tiyin int64` (faqat `_video/11-…`,
`12-…`). Hamma lablarni 3-modul oxiridagi bozor CLI v0 (GB-3.06) birlashtiradi. GB-2.05 bilan farq:
`diff -u ../gb-2-05-metodlar/solution/main.go solution/main.go`.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `Filter[T any]`, `Number`, `Sum[N Number]`, `Contains[T comparable]`, `writeTotal(w io.Writer, …)`; vazifaning `SumBy` (hozir har doim 0) va `paidTotal` (hozir 0) funksiyalari |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` + xuddi shu test fayli; `go test ./solution/` · `go run ./solution` |
| `_video/NN-…/main.go` | videodagi har bir qadam (toʻliq dastur); `_` bilan boshlangan papkani `go` buyruqlari chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi vaqtinchalik `bozor/` papkasida ishga tushiradi: `./_video/run.sh 05-sum-any go build` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Har bir xato oldida `go` `# github.com/abdurahmon27/bozor` qatorini ham chiqaradi (paket nomi); videoda u
koʻrsatilmagan. `run.sh` oxirida yozadigan `(chiqish kodi: N)` qatori ham bu yerda tushirilgan.

Boshlanishi (`00-takror`): `filterProducts` va `filterOrders` — tanasi bir xil, faqat turi boshqa ikki funksiya.

```
$ go run .
kam qolgan: non 3
kam qolgan: choy 0
toʻlangan buyurtma: 7
toʻlangan buyurtma: 9
```

`any` bilan bitta funksiya (`01-any`, `filterAny(items []any, …)`): `[]Product` `[]any` emas.

```
$ go build
./main.go:59:19: cannot use products (variable of type []Product) as []any value in argument to filterAny
```

Go FAQ: *"the two types do not have the same representation in memory"* — har elementni yangi slice'ga
nusxalash kerak, natija ham `[]any` boʻladi va turni yana `x.(Product)` bilan tiklash kerak.

`Filter[T any]` (`02-filter`): natija `00-takror` bilan bir xil. Kompilyator `T`'ni argumentdan chiqaradi.
Turni oʻzimiz yozsak (`03-tur`, `fmt.Printf("%T\n", Filter[Product])`):

```
$ go run .
func([]main.Product, func(main.Product) bool) []main.Product
func([]main.Order, func(main.Order) bool) []main.Order
```

Kvadrat qavssiz `fmt.Printf("%T\n", Filter)` esa kompilyatsiya boʻlmaydi: `in call to fmt.Printf, cannot infer T`
(`any` parametridan `T`'ni chiqaradigan joy yoʻq).

`T` argumentlardan mos kelmasa (`04-mos-emas`: `Filter(orders, func(p Product) bool {…})`):

```
$ go build
./main.go:59:24: in call to Filter, type func(p Product) bool of func(p Product) bool {…} does not match inferred type func(Order) bool for func(T) bool
```

Tez-tez uchraydigan xato (`05-sum-any`): `Sum[N any]` ichida `total += n`.

```
$ go build
./main.go:40:3: invalid operation: operator + not defined on total (variable of type N constrained by any)
```

Constraint `type Number interface { int | int64 }` bilan (`06-sum`):

```
$ go run .
43 4400000
```

`==` ham hamma turda yoʻq (`07-contains`, `Contains[T any]` ichida `if v == x`): slice, map va funksiya
qiymatlarini `==` bilan solishtirib boʻlmaydi.

```
$ go build
./main.go:39:6: invalid operation: v == x (incomparable types in type set)
```

Oldindan eʼlon qilingan `comparable` constraint'i bilan (`07b-comparable`, `Contains[T comparable]`):

```
$ go run .
true false
```

Standart kutubxonada xuddi shu funksiya bor: `slices.Contains[S ~[]E, E comparable]` — ishda oʻshani ishlating.

(Roʻyxatda yoʻq tur esa `Sum([]float64{1.5, 2})` → `float64 does not satisfy Number (float64 missing in int | int64)`.)

Ortiqcha tur parametri (`08-writer`: `writeTotal[W io.Writer](w W, …)`) va usiz (`08b-writer`:
`writeTotal(w io.Writer, …)`) — natija harfma-harf bir xil:

```
$ go run .
jami: 4400000 tiyin
```

`errors.AsType` (GB-2.08) endi oʻqiladi:

```
$ go doc errors.AsType
package errors // import "errors"

func AsType[E error](err error) (E, bool)
    AsType finds the first error in err's tree that matches the type E, and if
    one is found, returns that error value and true. …
```

`E` faqat natijada, shuning uchun kvadrat qavssiz chaqiruv (`09-astype`) kompilyatsiya boʻlmaydi:

```
$ go build
./main.go:51:28: in call to errors.AsType, cannot infer E (declared at $GOROOT/src/errors/wrap.go:167:13)
```

`errors.AsType[*StockError](err)` bilan (`10-astype-ok`):

```
$ go run .
*main.StockError olma 40
```

Chuqurroq: `type Tiyin int64` (`11-tilde`) — `Tiyin` `int64` asosidagi yangi tur, `int | int64` roʻyxatida yoʻq:

```
$ go build
./main.go:53:17: Tiyin does not satisfy Number (possibly missing ~ for int64 in Number)
```

`~int | ~int64` bilan (`12-tilde-ok`) — asosiy turi `int64` boʻlgan hamma turlar:

```
$ go run .
1600000
```

Standart kutubxonada ham shunday: `go doc slices.Clone` → `func Clone[S ~[]E, E any](s S) S`. `slices` paketida
`Filter`, `Map` yoki `Sum` yoʻq (`go doc slices`), shuning uchun ularni oʻzingiz yozasiz yoki `for` qoldirasiz.

Papkaning oʻzi (`final`, vazifa bajarilmagan):

```
$ go run .
kam qolgan: non 3
kam qolgan: choy 0
toʻlangan buyurtma: 7
toʻlangan buyurtma: 9
43 4400000
true false
jami: 0 tiyin
```

## Qachon tur parametri, qachon interface

I. L. Taylor, *When To Use Generics* (Go blog, 2022):

- *"If you find yourself writing the exact same code multiple times, where the only difference between the copies
  is that the code uses different types, consider whether you can use a type parameter."* — `Filter`, `Sum`.
- *"If all you need to do with a value of some type is call a method on that value, use an interface type, not a
  type parameter."* — `writeTotal(w io.Writer, …)`.
- Metodning implementatsiyasi har tur uchun boshqa boʻlsa — interface (har turning oʻz metodi).
- *"Don't change from interface types to type parameters just for speed, because it probably won't run any faster."*

Taylor'ning oʻz misoli: `func ReadSome[T io.Reader](r T) ([]byte, error)` oʻrniga `func ReadSome(r io.Reader) ([]byte, error)` — *"Don't make that
kind of change. Omitting the type parameter makes the function easier to write, easier to read, and the execution
time will likely be the same."* Belgisi: `W` imzoda bir marta, faqat parametr turi sifatida uchraydi va boshqa
hech qanday turni bogʻlamaydi (`08-writer` → `08b-writer`, natija bir xil). `Filter`'da esa `T` kirish va chiqish
turini bogʻlaydi: olib tashlansa, natija turi (`[]Product`) yoʻqoladi.

## Chuqurroq (videoda yoʻq)

- `any` slice'ining elementi interface qiymati: 16 bayt (tur va qiymatga koʻrsatkich), `Product` esa 40 bayt
  (`unsafe.Sizeof(any(nil))`, `unsafe.Sizeof(Product{})`, go1.27.4 linux/amd64). Shuning uchun `[]Product`
  `[]any`'ga aylanmaydi — har element yangi joyga nusxalanadi.
- Go 1.27'dan beri tur parametri faqat funksiya va turda emas, metodda ham boʻla oladi (*generic methods*):
  `func (l List[E]) Apply[F any](f func(E) F) List[F]`. Interface metodlari esa tur parametri eʼlon qila olmaydi
  va generic metod interface metodini bajarmaydi (https://go.dev/doc/go1.27#language).
- Go 1.21'dan beri `var f func([]Product, func(Product) bool) []Product = Filter` kompilyatsiya boʻladi (tur
  oʻzgaruvchining turidan chiqariladi); Go 1.27'dan beri konversiya ham:
  `(func([]Product, func(Product) bool) []Product)(Filter)`. Videodagi `fmt.Printf("%T\n", Filter)` esa yoʻq — `any` parametridan
  `T`'ni chiqarib boʻlmaydi (`cannot infer T`).
- Tur parametri turda ham boʻladi (`type List[T any] struct { … }`); keyingi darsdagi `iter.Seq[V]` shunday tur.

## Vazifa

[TASK.md](TASK.md): `SumBy[T any, N any]`'ning tanasini yozib, `N`'ga constraint tanlang; `paidTotal`'ni `Filter`
va `SumBy` bilan yozing. Tekshirish: `go test .` (hozir FAIL, vazifa bajarilgach `ok`). Javobni tekshirish:
`go test ./solution/`.

## Manbalar

- Go spec (go1.27): Type parameter declarations — https://go.dev/ref/spec#Type_parameter_declarations ·
  Type constraints — https://go.dev/ref/spec#Type_constraints · General interfaces (`|`, `~T`, faqat constraint
  sifatida) — https://go.dev/ref/spec#General_interfaces · Instantiations — https://go.dev/ref/spec#Instantiations ·
  Type inference — https://go.dev/ref/spec#Type_inference · Arithmetic operators (tur parametri uchun operator
  roʻyxatdagi har bir turga mos boʻlishi kerak) — https://go.dev/ref/spec#Arithmetic_operators
- Go 1.18 release notes: Generics — https://go.dev/doc/go1.18#generics
- Go blog: An Introduction To Generics (R. Griesemer, I. L. Taylor, 2022) — https://go.dev/blog/intro-generics
- Go blog: When To Use Generics (I. L. Taylor, 2022) — https://go.dev/blog/when-generics
- Go FAQ: Can I convert a []T to an []interface{}? — https://go.dev/doc/faq#convert_slice_of_interface
- Go tutorial: Getting started with generics — https://go.dev/doc/tutorial/generics
- pkg.go.dev `errors.AsType` — https://pkg.go.dev/errors#AsType · Go 1.26 release notes — https://go.dev/doc/go1.26#errorspkgerrors
- pkg.go.dev `slices` — https://pkg.go.dev/slices
