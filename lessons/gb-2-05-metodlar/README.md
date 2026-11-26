# GB-2.05 · Metodlar: value receiver yoki pointer receiver

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #13"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-2-04-pointer`](../gb-2-04-pointer) (pointer: nusxa yoki umumiy qiymat) · keyingi dars: `gb-2-06` (paketlar: katta harf,
`internal/` va loyiha tuzilishi). Domen turlari: [`gb-2-03-struct-turlar`](../gb-2-03-struct-turlar).

**Bu darsdan keyin:** turga metod yozasiz (`func (o *Order) Total() int64`) va uni `order.Total()` deb
chaqirasiz; value receiver nusxa, pointer receiver manzil olishini, Go `order.Pay()`'ni oʻzi `(&order).Pay()`
qilishini bilasiz; receiver'ni qoida boʻyicha tanlaysiz; int asosidagi `OrderStatus`'ga `String()` yozib, `fmt`
holatni nomi bilan chiqarishini koʻrasiz.

Bu papka 2-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i, oʻz `package main`'i; turlar
GB-2.03'dagidek. 3-modul oxiridagi bozor CLI v0 (GB-3.06) labs'ni birlashtiradi.
Oldingi dars bilan farqni koʻrish: `diff -ru ../gb-2-04-pointer/solution solution` (turlar GB-2.03'dan).

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: turlar, `main`, `String`, `Total`, `Pay`, `AddItem`; oxirida vazifaning `Cancel`'i (ataylab value receiver bilan) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 02-pay-value go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

`Total` metod sifatida (`./_video/run.sh 01-total go run .`): olma 1 200 000 tiyin × 3:

```
3600000
```

`Pay` value receiver bilan (`02-pay-value`: `func (o Order) Pay()`): holat oʻzgarmaydi, chunki `o` — nusxa:

```
0
```

`Pay` pointer receiver bilan (`03-pay-pointer`: `func (o *Order) Pay()`; `main` oʻzgarmagan, Go `(&order).Pay()`
qiladi):

```
1
```

`AddItem` pointer receiver bilan, olma 3 + non 2 (`04-additem`):

```
2 4400000
```

Tez-tez uchraydigan xato: `AddItem`'da yulduzcha yoʻq (`05-xato`: `func (o Order) AddItem(…)`). Kompilyator ham,
`go vet` ham jim:

```
$ go run .
0 0
$ go vet .
$
```

Yakuniy kod (`./_video/run.sh final go run .`, papkadagi `main.go` ham shuni chiqaradi):

```
2 4400000 toʻlangan
```

## Vazifa

[TASK.md](TASK.md): `AddItem` bir mahsulotni bitta qatorga birlashtirsin, `Cancel` ishlasin (receiver va qoida).
Tekshirish: `go test` (hozir FAIL, vazifa bajarilgach `ok`).

## Manbalar

- Go spec: Method declarations — https://go.dev/ref/spec#Method_declarations
- Go spec: Calls (`x.m()` — `(&x).m()` qisqartmasi) — https://go.dev/ref/spec#Calls
- Go spec: Method sets — https://go.dev/ref/spec#Method_sets
- Go FAQ: methods on values or pointers — https://go.dev/doc/faq#methods_on_values_or_pointers
- Go Code Review Comments: Receiver Type — https://go.dev/wiki/CodeReviewComments#receiver-type
- pkg.go.dev `fmt` (`String() string` metodi) — https://pkg.go.dev/fmt
