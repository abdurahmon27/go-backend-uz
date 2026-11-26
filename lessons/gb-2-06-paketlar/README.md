# GB-2.06 · Paketlar: katta harf, `internal/` va loyiha tuzilishi

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #14"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-2-05-metodlar`](../gb-2-05-metodlar) (metodlar: value yoki pointer receiver) · keyingi dars:
`gb-2-07` (interface: `io.Writer` va ikki haqiqiy implementatsiya).

**Bu darsdan keyin:** bitta `main.go`'ni paketlarga boʻlasiz (`cmd/bozor` — dastur, `internal/store` — katalog);
import yoʻli = `go.mod`'dagi modul nomi + papka ekanini, katta harf bilan boshlangan nom boshqa paketga koʻrinishini
(exported), kichik harfli nom faqat oʻz paketida qolishini bilasiz; tashqi modulni `go mod tidy` bilan qoʻshasiz,
`go.sum` nimani tekshirishini va `internal/` kimga import qilishga ruxsat berishini koʻrasiz.

Bu papka 2-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i bor. `Product` GB-2.03'dagidek; `Store`
(xotiradagi katalog: `Add`, `List`, `Sell`) shu darsda paydo boʻladi. 3-modul oxiridagi bozor CLI v0 (GB-3.06)
labs'ni birlashtiradi.

## Papka

| Fayl | Nima |
|---|---|
| `go.mod`, `go.sum` | modul `github.com/abdurahmon27/bozor`, `go 1.27`, `require golang.org/x/text v0.44.0` va uning xeshlari |
| `cmd/bozor/main.go` | dastur (`package main`): katalogni yaratadi, narxlarni soʻmda guruhlab chiqaradi |
| `internal/store/store.go` | `package store`: `Product`, `Store`, `Add`, `List`; vazifaning `Sell`'i (hozir hech narsa qilmaydi) |
| `internal/store/tekshiruv_test.go` | vazifani tekshiradi (`go test ./...`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: xuddi shu daraxt, oʻz `go.mod`'i bilan (alohida modul, shuning uchun ildizdagi `go test ./...` unga tegmaydi) |
| `_video/NN-…/` | videodagi har bir qadam (toʻliq kichik loyiha); `_` bilan boshlangan papkani `go` buyruqlari chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 01-paket go run ./cmd/bozor` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Boshlanishi: hammasi bitta `main.go`'da, `main` katalog ichidagi slice'ga toʻgʻridan-toʻgʻri yozadi
(`./_video/run.sh 00-birfayl go run .`):

```
-5
```

Paketlarga boʻlish: `mkdir -p cmd/bozor internal/store`, `mv main.go cmd/bozor/`, turlar va metodlar
`internal/store/store.go`'ga (`01-paket`). Ildizda endi Go fayl yoʻq:

```
$ go run .
no Go files in …/bozor
$ go run ./cmd/bozor
1 olma 1200000
2 non 400000
```

`main` boshqa paketdagi kichik harfli maydonga yozadi (`02-eksport`: `st.products[0].Stock = -5`):

```
# github.com/abdurahmon27/bozor/cmd/bozor
cmd/bozor/main.go:13:5: st.products undefined (cannot refer to unexported field products)
```

Chegaradagi teshik: `List` ichki slice'ni qaytaradi, natijaga yozish katalogni oʻzgartiradi (`06-teshik`:
`st.List()[0].Stock = -5`, keyin `fmt.Println(st.List()[0].Stock)`). Vazifaning 2-qadami shuni yopadi:

```
-5
```

Tez-tez uchraydigan xato: `Add` oʻrniga `add` (`03-kichik-harf`):

```
# github.com/abdurahmon27/bozor/cmd/bozor
cmd/bozor/main.go:12:5: st.add undefined (type store.Store has no field or method add, but does have method Add)
```

(Paketda oʻzingiz kichik harf bilan yozgan funksiyani tashqaridan chaqirsangiz, xabar faqat
`undefined: store.newID` boʻladi.)

Tashqi modul: `golang.org/x/text/language` va `…/message` import qilindi, `go.mod`'da hali yoʻq (`04-modul`):

```
$ go run ./cmd/bozor
cmd/bozor/main.go:5:2: no required module provides package golang.org/x/text/language; to add it:
	go get golang.org/x/text/language
cmd/bozor/main.go:6:2: no required module provides package golang.org/x/text/message; to add it:
	go get golang.org/x/text/message
$ go mod tidy
go: finding module for package golang.org/x/text/message
go: finding module for package golang.org/x/text/language
go: found golang.org/x/text/language in golang.org/x/text v0.44.0
go: found golang.org/x/text/message in golang.org/x/text v0.44.0
$ cat go.mod
module github.com/abdurahmon27/bozor

go 1.27

require golang.org/x/text v0.44.0
$ cat go.sum
golang.org/x/text v0.44.0 h1:YItdp/mQNGnsmrgpcdA5yQ3bnRxCfo/Tj/iByGW9GGo=
golang.org/x/text v0.44.0/go.mod h1:S4FR1waQ9LbjQTls97SyRjDaUUVkaM+cXwwuG/G+MwY=
$ go run ./cmd/bozor
1. olma: 12 000 soʻm
2. non: 4 000 soʻm
```

(`12 000` ichidagi boʻsh joy — oddiy probel emas, U+00A0 NO-BREAK SPACE: oʻzbek tili qoidasi shunday.)

`internal/`: x/text'ning ichki paketini import qilish (`05-internal`):

```
package github.com/abdurahmon27/bozor/cmd/bozor
	cmd/bozor/main.go:5:2: use of internal package golang.org/x/text/internal/number not allowed
```

Yakuniy kod (`./_video/run.sh final go run ./cmd/bozor`, papkaning oʻzida `go run ./cmd/bozor` ham shuni chiqaradi):

```
1. olma: 12 000 soʻm
2. non: 4 000 soʻm
```

## Vazifa

[TASK.md](TASK.md): `Store.Sell` qoldiq qoidasini bajarsin, `List` nusxa qaytarsin.
Tekshirish: `go test ./...` (hozir FAIL, vazifa bajarilgach `ok`). Javob bilan farq:
`diff -ru internal solution/internal`.

## Manbalar

- Go spec: Packages, Package clause — https://go.dev/ref/spec#Packages
- Go spec: Exported identifiers — https://go.dev/ref/spec#Exported_identifiers
- Go Modules Reference: modullar, `require`, `go.sum` — https://go.dev/ref/mod
- `go` buyrugʻi: Internal Directories — https://pkg.go.dev/cmd/go#hdr-Internal_Directories
- Effective Go: Package names — https://go.dev/doc/effective_go#package-names
- Go: Organizing a Go module (`cmd/`, `internal/`) — https://go.dev/doc/modules/layout
- `golang.org/x/text/message` — https://pkg.go.dev/golang.org/x/text/message
