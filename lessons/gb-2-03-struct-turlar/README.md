# GB-2.03 · Struct: `Product` va `Order` turlarini yasaymiz

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #11"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-2-02` (map: qidirish, comma-ok) · keyingi dars: `gb-2-04` (pointer: nusxa yoki umumiy qiymat).
Birinchi modulning oxirgi holati: [`gb-1-08-birinchi-test`](../gb-1-08-birinchi-test).

**Bu darsdan keyin:** Go'da struct turini eʼlon qilasiz (`type Product struct { … }`), uni maydon nomlari bilan
toʻldirasiz (struct literal), yozilmagan maydonning zero value'sini va `q := p` struct'ni toʻliq nusxalashini
bilasiz, buyurtma holatini `iota` bilan yozasiz.

Bu papka 2-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i, oʻz `package main`'i. bozor'ning
domen turlari bu yerda birinchi marta paydo boʻladi; 3-modul oxiridagi bozor CLI v0 (GB-3.06) ularni birlashtiradi.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `Product`, `OrderItem`, `Order`, `OrderStatus` + `iota`, `main`; oxirida vazifaning ikki boʻsh funksiyasi (`orderTotal`, `canCancel`) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi (`05-xato` ataylab kompilyatsiya boʻlmaydi) |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 03-nusxa go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Birinchi struct, literal va `%+v` (`./_video/run.sh 01-product go run .`):

```
olma 1200000
{ID:1 Name:olma Price:1200000 Stock:40}
```

Yozilmagan maydonlar zero value oladi (`02-nol`: `Product{Name: "non", Price: 400_000}`):

```
{ID:0 Name:non Price:400000 Stock:0}
```

`q := p` struct'ni toʻliq nusxalaydi (`03-nusxa`: `p.Stock = 40; q := p; q.Stock = 0; fmt.Println(p.Stock, q.Stock)`):

```
40 0
```

Buyurtma va zero value'li buyurtma (`04-order`, papkadagi `main.go` ham shuni chiqaradi):

```
{ID:1 Items:[{ProductID:1 Price:1200000 Quantity:3}] Status:0}
true 0
```

`Status:0`: `%v` holatning sonini yozadi; nomini chiqarish (`String()` metodi) GB-2.05'da.

Tez-tez uchraydigan xato: maydon nomlarisiz literal (`05-xato`: `Product{1, "olma", 1_200_000}`):

```
# github.com/abdurahmon27/bozor
./main.go:14:35: too few values in struct literal of type Product
```

Nomsiz literal'da hamma maydon eʼlondagi tartibda boʻlishi shart; toʻrttasi yozilsa kompilyatsiya boʻladi, lekin
`Product`'ga yangi maydon qoʻshilgan kuni har bir shunday qator buziladi. Shuning uchun maydon nomlari bilan
yozing. (`go vet` boshqa paketdagi struct'ning nomsiz literal'ini ham ogohlantiradi: `composites` tekshiruvi.)

Chuqurroq: saqlanadigan `iota` roʻyxatining oʻrtasiga qiymat qoʻshish (`06-iota`: `StatusNew` dan keyin
`StatusReserved`; `fmt.Println(StatusPaid, StatusShipped)`):

```
2 3
```

Oldin `1` deb saqlangan "toʻlangan" buyurtma endi "band qilingan" boʻlib oʻqiladi. Qoida: yangi qiymat faqat
roʻyxat oxiriga. `tekshiruv_test.go`'dagi `TestHolatRaqamlari` raqamlarni shu sababli qotiradi.

## Vazifa

[TASK.md](TASK.md): `orderTotal` (buyurtma summasi tiyinda) va `canCancel` (yangi va toʻlangan → `true`).
Tekshirish: shu papkada `go test`, `ok` chiqquncha. Hozir (boshlangʻich kod) `FAIL` bilan yiqiladi.

Papkani tekshirish (muallif uchun):

```bash
gofmt -l . && go vet ./...                 # hech narsa chiqmaydi
go test . ; echo "starter exit $?"         # FAIL, exit 1
go test ./solution/...                     # ok
```

## Manbalar

- Go spec: [Struct types](https://go.dev/ref/spec#Struct_types), [Composite literals](https://go.dev/ref/spec#Composite_literals),
  [The zero value](https://go.dev/ref/spec#The_zero_value), [Iota](https://go.dev/ref/spec#Iota),
  [Constant declarations](https://go.dev/ref/spec#Constant_declarations), [Type definitions](https://go.dev/ref/spec#Type_definitions),
  [Assignability](https://go.dev/ref/spec#Assignability), [Calls](https://go.dev/ref/spec#Calls)
- [pkg.go.dev/fmt](https://pkg.go.dev/fmt) (`%v`, `%+v`)
- [Effective Go: Composite literals, Constants](https://go.dev/doc/effective_go#composite_literals)
- [Go FAQ: pass by value](https://go.dev/doc/faq#pass_by_value)
