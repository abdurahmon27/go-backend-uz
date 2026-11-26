# GB-2.07 · Interface: `io.Writer` va ikki haqiqiy implementatsiya

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #15"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-2-06-paketlar`](../gb-2-06-paketlar) (paketlar: katta harf, `internal/` va loyiha tuzilishi) · keyingi dars: `gb-2-08`
(error'lar chuqur: `%w`, `errors.Is` va `errors.As`). Metodlar va receiver: [`gb-2-05-metodlar`](../gb-2-05-metodlar).

**Bu darsdan keyin:** `go doc` bilan standart kutubxonadagi turni oʻqiysiz; interface — metodlar roʻyxati ekanini,
tur uni hech qanday eʼlonsiz, faqat kerakli metodga ega boʻlib bajarishini bilasiz; bitta
`writeReport(w io.Writer, …)` funksiyasi bilan hisobotni ekranga (`os.Stdout`), faylga (`os.Create`) va oʻzingiz
yozgan `*byteCounter`'ga yozasiz; «method Write has pointer receiver» xatosini oʻqib tuzatasiz; `String() string`
metodi bilan `fmt.Stringer`'ni bajarasiz va `fmt` uni qanday topishini (`x.(fmt.Stringer)`) koʻrasiz.

Bu papka 2-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i, oʻz `package main`'i; `Product` turi
GB-2.03'dagidek. 3-modul oxiridagi bozor CLI v0 (GB-3.06) labs'ni birlashtiradi.
Metodlar darsi bilan farqni koʻrish: `diff -ru ../gb-2-05-metodlar/solution solution`.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `Product` + `String`, `writeReport`, `byteCounter`, `main`; oxirida vazifaning `lineCounter` va `writeLowStock`'i (chala) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 02-fayl go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Standart kutubxona hujjati terminalda (qisqartirilgan: `package …` qatori va hujjat matni olib tashlangan; toʻliq matnni oʻzingiz koʻring):

```
$ go doc fmt.Fprintln
func Fprintln(w io.Writer, a ...any) (n int, err error)
$ go doc io.Writer
type Writer interface {
	Write(p []byte) (n int, err error)
}
$ go doc os.File.Write
func (f *File) Write(b []byte) (n int, err error)
$ go doc fmt.Stringer
type Stringer interface {
	String() string
}
```

`writeReport(os.Stdout, products)` (`./_video/run.sh 01-stdout go run .`):

```
olma    12000 soʻm   40 dona
non      4000 soʻm  100 dona
oʻrik   35000 soʻm    5 dona
```

Xuddi shu funksiya faylga (`02-fayl`: `os.Create("hisobot.txt")`, `writeReport(f, products)`, `f.Close()`):
terminalda hech narsa yoʻq, hisobot faylda:

```
$ go run .
$ cat hisobot.txt
olma    12000 soʻm   40 dona
non      4000 soʻm  100 dona
oʻrik   35000 soʻm    5 dona
```

Oʻzimizning `*byteCounter` (`03-counter`): `io.Writer` hech qayerda tilga olinmaydi, lekin `Write` metodi bor.
Hisoblagich fayldagi bayt soni bilan bir xil (`soʻm`'dagi `ʻ` ikki bayt; UTF-8 GB-2.10'da):

```
$ go run .
hisobot: 91 bayt
$ wc -c hisobot.txt
91 hisobot.txt
```

Tez-tez uchraydigan xato (`04-xato`: `writeReport(c, products)`, `&c` emas):

```
$ go run .
# github.com/abdurahmon27/bozor
./main.go:65:20: cannot use c (variable of struct type byteCounter) as io.Writer value in argument to writeReport: byteCounter does not implement io.Writer (method Write has pointer receiver)
```

`String()`'siz va `String()` bilan `fmt.Println(products)` (`05-oldin`, keyin `final` = papkadagi `main.go`):

```
hisobot: 91 bayt
[{1 olma 1200000 40} {2 non 400000 100} {3 oʻrik 3500000 5}]
hisobot: 91 bayt
[olma (12000 soʻm) non (4000 soʻm) oʻrik (35000 soʻm)]
```

## Vazifa

[TASK.md](TASK.md): `writeLowStock` hamma qatorni `w`'ga yozsin va yozish xatosini qaytarsin; `*lineCounter`
`io.Writer` boʻlsin va qatorlarni sanasin. Tekshirish: `go test` (hozir FAIL, vazifa bajarilgach `ok`).

## Manbalar

- Go spec: Interface types — https://go.dev/ref/spec#Interface_types
- Go spec: Method sets — https://go.dev/ref/spec#Method_sets
- Go spec: Type assertions — https://go.dev/ref/spec#Type_assertions
- pkg.go.dev `io.Writer` — https://pkg.go.dev/io#Writer
- pkg.go.dev `os.Create`, `os.File.Write`, `os.Stdout` — https://pkg.go.dev/os
- pkg.go.dev `fmt` (`Fprintf`, `Stringer`) — https://pkg.go.dev/fmt
- Effective Go: Interfaces — https://go.dev/doc/effective_go#interfaces
- Go Code Review Comments: Interfaces — https://go.dev/wiki/CodeReviewComments#interfaces
- `go doc` buyrugʻi — https://pkg.go.dev/cmd/go#hdr-Show_documentation_for_package_or_symbol
