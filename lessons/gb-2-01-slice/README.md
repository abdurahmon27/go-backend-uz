# GB-2.01 · Slice: `len`, `cap`, `append` va umumiy massiv

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #9"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

Oldingi dars: [`gb-1-08-birinchi-test`](../gb-1-08-birinchi-test) · keyingi dars: `gb-2-02` (map: qidirish,
comma-ok va kafolatsiz tartib).

**Bu darsdan keyin:** slice xotirada qanday turishini bilasiz: uch maydonli sarlavha (`ptr` massivning manzili,
`len`, `cap`) va uning orqasidagi massiv. `before := stock` va funksiyaga berish faqat sarlavhani nusxalaydi,
massiv bitta qoladi. `append` joy boʻlsa oʻsha massivga yozadi, joy boʻlmasa yangi massiv ochadi va yangi sarlavha
qaytaradi. Funksiya natijasini `make` bilan alohida massivga qurib, chaqiruvchining maʼlumotini himoya qilasiz.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi bozor laboratoriyasi: tuzatilgan `afterOrder` (natija `make` bilan alohida massivda) + vazifa uchun `discounted` (hali xato bilan) |
| `stock_test.go` | `TestAfterOrder`: qoldiq toʻgʻri va `stock` oʻzgarmagan |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa: `discounted` katalog narxlarini oʻzgartirmasin |
| `solution/` | javob: `main.go` (+ xuddi shu ikki test fayli) |
| `_video/NN-…/` | videodagi har bir qadamning `main.go` fayli; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi (`05-xato`da ataylab kompilyatsiya xatosi bor) |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 03-append go run .` |

Bu Module 2 darslari kabi mustaqil laboratoriya: bozor domenidagi kichik `main.go` (Module 1 nomlari: `stock`,
`sold`, narx `int64` tiyinda). `Product` turi GB-2.03'da keladi, shuning uchun bu yerda mahsulotlar tartib raqami
bir xil boʻlgan slice'lar: `stock[0]` olma, `stock[1]` non, `stock[2]` choy.

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`)

Xato (`./_video/run.sh 00-bug go run .`): `afterOrder` parametr orqali `stock` massiviga yozadi.

```
ombor oldin: [40 12 7]
koʻrish:     [38 11 7]
ombor keyin: [38 11 7]
```

Slice literal, `len`, `cap`, indeks (`01-slice`):

```
[40 12 7] 3 3
40 7
```

Sarlavha nusxasi, massiv bitta (`02-nusxa`): `before := stock; stock[0] = 38; fmt.Println(before[0])` → `38`.

`append` va manzil (`03-append`; `%p` slice qaraydigan massivning manzilini chiqaradi, har ishga tushirishda
boshqacha):

```
3 3 0x103325296000
4 6 0x10332529e000
5 6 0x10332529e000
```

Birinchi `append`da joy yoʻq (len = cap = 3): yangi massiv, cap 6. Ikkinchisida joy bor: oʻsha massiv, len 5.

`append` yangi massiv ochganda eski sarlavha eski massivda qoladi (`04-append-nusxa`):
`before := stock; stock = append(stock, 25); stock[0] = 38; fmt.Println(before[0], stock[0])` → `40 38`.

Tez-tez uchraydigan xato, natijani yozmaslik (`./_video/run.sh 05-xato go build`):

```
# github.com/abdurahmon27/bozor
./main.go:7:2: append(stock, 25) (value of type []int) is not used
```

Tuzatish (`06-tuzat`): `left := make([]int, len(stock))`, sikl faqat `left`ga yozadi.

```
ombor oldin: [40 12 7]
koʻrish:     [38 11 7]
ombor keyin: [40 12 7]
```

`go test` (`stock_test.go` bilan) → `PASS` / `ok  github.com/abdurahmon27/bozor`.

Oraliq savol (`08-make`): `s := make([]int, 0, 4); s = append(s, 1)` → `len 1, cap 4`.

Chuqurroq (`07-ikki-append`): joy bor sarlavhadan ikki `append` bitta elementga yozadi:

```go
stock := []int{40, 12, 7}
stock = append(stock, 25) // len 4, cap 6
a := append(stock, 3)
b := append(stock, 9)
fmt.Println(a[4], b[4]) // 9 9
```

Kichik slice'lar taxminan ikki baravar oʻsadi (runtime'ning tanlovi, spec faqat "yetarlicha katta" deydi). Bu
kompyuterda (go1.27.4; aniq sonlar Go versiyasiga bogʻliq) boʻsh slice'ga 2000 marta `append` qilinganda cap:
`4 8 16 32 64 128 256 512 848 1280 1792 2560` (256 dan keyin oʻsish sekinlashadi; batafsil GB-11.05).

Eslatmalar:
- `go.mod`dagi `go 1.27`: Module 2 Go 1.27'da yozildi (Module 1 papkalarida `go 1.26.8`). Kompyuteringizda
  eskiroq Go boʻlsa, `go` buyrugʻi kerakli versiyani oʻzi yuklab oladi (toolchain), internet kerak.
- `afterOrder` `stock` va `sold` uzunligi teng deb oladi va manfiy qoldiqni rad etmaydi: darsda faqat massiv
  boʻlishish mavzusi. Kirishni tekshirish qoidasi (GB-1.05, GB-1.07) bozor CLI v0'da (GB-3) qoʻshiladi.
- `discounted`dagi `p*percent/100` butun sonli boʻlish: chegirma tiyingacha pastga yaxlitlanadi. Darsdagi
  narxlarda qoldiq yoʻq.
- `solution/` alohida modul emas, shu modulning ichki paketi: `go run ./solution` va `go test ./solution/...`.

## Vazifa

[`TASK.md`](TASK.md): `discounted` aksiya narxlarini alohida massivda qaytarsin, `prices` oʻzgarmasin.
Tekshirish: shu papkada `go test`. Javob: [`solution/main.go`](solution/main.go).

```bash
diff -u main.go solution/main.go   # javob nimani oʻzgartiradi
```

## Manbalar

- Go spec: [Slice types](https://go.dev/ref/spec#Slice_types), [Length and capacity](https://go.dev/ref/spec#Length_and_capacity),
  [Appending and copying slices](https://go.dev/ref/spec#Appending_and_copying_slices), [For range](https://go.dev/ref/spec#For_range),
  [Array types](https://go.dev/ref/spec#Array_types), [Composite literals](https://go.dev/ref/spec#Composite_literals)
- [pkg.go.dev builtin: append, make](https://pkg.go.dev/builtin) · [fmt: `%p`](https://pkg.go.dev/fmt) ·
  [slices.Equal](https://pkg.go.dev/slices#Equal)
- go.dev/blog: [Go Slices: usage and internals](https://go.dev/blog/slices-intro) ·
  [Arrays, slices (and strings): The mechanics of 'append'](https://go.dev/blog/slices)
- Tour of Go: [Slices](https://go.dev/tour/moretypes/7)
