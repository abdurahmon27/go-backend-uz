# GB-2.04 · Pointer: nusxa yoki umumiy qiymat

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #12"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

Oldingi dars: [`gb-2-03-struct-turlar`](../gb-2-03-struct-turlar) (struct: `Product` va `Order`) · keyingi dars:
`gb-2-05` (metodlar: value receiver yoki pointer receiver).

**Bu darsdan keyin:** funksiyaga struct berilganda nusxa ketishini xotira rasmida koʻrasiz, `&olma` bilan
qiymatning manzilini olib `*Product` turidagi pointer'ga berasiz, `sell(p *Product, quantity int)` bilan
qoldiqni chaqiruvchining qiymatida kamaytirasiz, katalogdan elementning manzilini (`&products[i]`) qaytarasiz va
nil pointer panic'ini `if p == nil` bilan oldini olasiz.

Bu papka 2-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`,
`go 1.27`), oʻz `package main`'i. `Product` turi GB-2.03'dagidek; 3-modul oxiridagi bozor CLI v0 (GB-3.06)
modul darslarini birlashtiradi.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `Product`, `sell`, `findProduct`, `main` (nil tekshiruvi bilan); oxirida vazifaning ikki boʻsh funksiyasi (`sellByID`, `cheapest`) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi (`03-xato` ataylab kompilyatsiya boʻlmaydi, `06-nil` ataylab panic qiladi) |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 04-sell go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Muammo: `sell(p Product, quantity int)` nusxani kamaytiradi (`./_video/run.sh 01-nusxa go run .`):

```
sell ichida: 37
main'da: 40
```

Pointer: `p := &olma; p.Stock = 37; q := p; q.Stock = 35; fmt.Println(olma.Stock, p.Stock)` (`02-pointer`).
`q := p` manzilni nusxalaydi, ikkala pointer bitta `olma`'ni koʻrsatadi:

```
35 35
```

Parametr `*Product` boʻldi, chaqiruv esa eskicha `sell(olma, 3)` (`03-xato`):

```
# github.com/abdurahmon27/bozor
./main.go:19:7: cannot use olma (variable of struct type Product) as *Product value in argument to sell
```

Toʻgʻrisi `sell(&olma, 3)` (`04-sell`):

```
main'da: 37
```

Katalogdan topish: `findProduct` `&products[i]` qaytaradi, `sell(p, 1)` katalogdagi non'ni kamaytiradi (`05-find`,
`fmt.Println(catalog[1].Stock)`):

```
99
```

Katalogda yoʻq ID: `findProduct(catalog, 9)` → `nil`, keyin `p.Stock` (`06-nil`; videoda birinchi ikki qator
koʻrsatiladi, qolgani `…`):

```
true
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x20 pc=0x499f50]

goroutine 1 [running]:
main.main()
	/tmp/tmp.qUNOXSklzM/bozor/main.go:34 +0x110
exit status 2
```

Tuzatish: `if p == nil { fmt.Println("mahsulot topilmadi: 9"); return }` (`07-guard`, papkadagi `main.go` ham shu):

```
mahsulot topilmadi: 9
```

Chuqurroq: `for _, p := range products { if p.ID == id { return &p } }` — `p` elementning nusxasi, `&p` nusxaning
manzili (`08-range`, `fmt.Println(p.Stock, catalog[1].Stock)`):

```
99 100
```

## Vazifa

[TASK.md](TASK.md): `sellByID` (topilmasa, miqdor ≤ 0 yoki qoldiq yetmasa xato, qoldiqqa tegmaydi) va `cheapest`
(eng arzon mahsulotning katalogdagi manzili, boʻsh katalogda `nil`). Tekshirish: shu papkada `go test`, `ok`
chiqquncha. Hozir (boshlangʻich kod) `FAIL` bilan yiqiladi.

Oldingi darsdan nima oʻzgargani: 2-modul darslari mustaqil lab (koordinator qoidasi), shuning uchun bu papka
gb-2-03'ning snapshot'i emas: `Product` turi oʻsha, `Order` turlari esa bu darsga kerak emas. Toʻliq bozor GB-3.06'da
birlashadi.

Eʼtibor: `findProduct` qaytargan pointer katalog slice'ining massiviga ishora qiladi. Keyin `append` yangi massiv
ochsa (GB-2.01), eski pointer eski massivni oʻzgartiradi, katalog esa oʻzgarmaydi: pointer'ni saqlab qoʻymang,
katalogni oʻzgartirgandan keyin qaytadan toping.

Papkani tekshirish (muallif uchun):

```bash
gofmt -l . && go vet ./...                 # hech narsa chiqmaydi
go test . ; echo "starter exit $?"         # FAIL, exit 1
go test ./solution/...                     # ok
```

## Manbalar

- Go spec: [Address operators](https://go.dev/ref/spec#Address_operators), [Pointer types](https://go.dev/ref/spec#Pointer_types),
  [Selectors](https://go.dev/ref/spec#Selectors), [Calls](https://go.dev/ref/spec#Calls),
  [Assignability](https://go.dev/ref/spec#Assignability), [The zero value](https://go.dev/ref/spec#The_zero_value),
  [For statements with range clause](https://go.dev/ref/spec#For_range)
- [Go FAQ: When are function parameters passed by value?](https://go.dev/doc/faq#pass_by_value)
- [Go Code Review Comments: Pass Values](https://go.dev/wiki/CodeReviewComments#pass-values)
- [A Tour of Go: Pointers](https://go.dev/tour/moretypes/1), [Pointers to structs](https://go.dev/tour/moretypes/4)
- [pkg.go.dev/testing](https://pkg.go.dev/testing) (`func TestXxx(*testing.T)`)
