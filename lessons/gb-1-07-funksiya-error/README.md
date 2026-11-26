# GB-1.07 · Funksiya va birinchi `error`: `os.Args` va `strconv.Atoi`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #7"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

Oldingi dars: [`gb-1-06-for-sikl`](../gb-1-06-for-sikl) · keyingi dars: [`gb-1-08-birinchi-test`](../gb-1-08-birinchi-test).

**Bu darsdan keyin:** bozor nom va narxni kod ichidan emas, buyruq qatoridan oladi: `./bozor add olma 12000`
(chek GB-1.06'dagidek hisoblanadi).
Narxni `parsePrice(s string) (int64, error)` funksiyasi tekshiradi, notoʻgʻri kirish xato oqimiga aniq xabar va
chiqish kodi 1 bilan rad etiladi. Vositalar: `os.Args` va `len`, `strconv.Atoi`, `(qiymat, error)`, `nil`,
`if err != nil`, `func` / parametr / natijalar / `return`, `fmt.Errorf`, `fmt.Fprintln(os.Stderr, …)`, `os.Exit(1)`.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi bozor (vazifa shu yerdan boshlanadi) |
| `main_test.go` | vazifani tekshiradi (`go test`): dasturni turli argumentlar bilan ishga tushiradi; test yozishni GB-1.08'da ochamiz |
| `TASK.md` | vazifa: miqdor buyruq qatoridan, `parseSold` |
| `solution/` | javob (xuddi shu test bilan) |
| `_video/NN-…/main.go` | videodagi har bir qadamning kodi; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi (`02-plus-xato` ataylab kompilyatsiya boʻlmaydi) |
| `_video/run.sh` | qadamni videodagidek ishga tushiradi: `./_video/run.sh 03-atoi add olma 12k`; `final` = papkadagi `main.go` |

`main.go` GB-1.06'ning `solution/main.go`'sidan boshlangan. Oʻzgarganlar: `import ( … )` (`os`, `strconv`),
`parsePrice`, `main` boshida ishga tushish qatorlaridan keyin `os.Args` tekshiruvi, `name` va `price` buyruq
qatoridan (paket darajasidagi `var name`/`var price` olib tashlandi), `price <= 0` tekshiruvi `parsePrice`'ga koʻchdi,
miqdor hali `var sold = 7` (vazifada buyruq qatoridan). Tekshiruvlar (GB-1.05), aksiya sikli (GB-1.06) va chek oʻzgarmagan.
Oldingi darsdan farq:

```bash
diff -u ../gb-1-06-for-sikl/solution/main.go main.go
```

## Videodagi buyruqlar va natijalari (go1.26.8 linux/amd64, `CGO_ENABLED=0`)

`os.Args` (`./_video/run.sh 01-args add olma 12000`):

```
4 [./bozor add olma 12000]
```

`os.Args[3] + 1` (`02-plus-xato`, `go build`):

```
# github.com/abdurahmon27/bozor
./main.go:10:14: invalid operation: price + 1 (mismatched types string and untyped int)
```

`strconv.Atoi(os.Args[3])` (`03-atoi`) — `12000`, `12k`, yigirmata `9`:

```
12000 <nil>
0 strconv.Atoi: parsing "12k": invalid syntax
9223372036854775807 strconv.Atoi: parsing "99999999999999999999": value out of range
```

Chegaradan chiqqanda `Atoi` 0 emas, `int`'ning eng katta qiymatini qaytaradi
([strconv.ParseInt](https://pkg.go.dev/strconv#ParseInt)): `err` nil boʻlmasa, qiymatga ishonmang.

`if err != nil` + `os.Exit(1)` (`04-iferr`):

```
$ ./bozor add olma 12k
xato: strconv.Atoi: parsing "12k": invalid syntax
$ echo $?
1
```

`parsePrice` bilan, bozor (`final`; `xato:` va `foydalanish:` qatorlari xato oqimida):

```
$ ./bozor add olma 12000
bozor ishga tushdi
versiya: 0.1.0
olma, 7 dona       60000.00
QQS                 7200.00
yetkazish          15000.00
jami               82200.00
aksiya: 2 dona bepul
$ ./bozor add olma -5
bozor ishga tushdi
versiya: 0.1.0
xato: narx musbat boʻlishi kerak, keldi -5
$ ./bozor add olma 12_000
bozor ishga tushdi
versiya: 0.1.0
xato: strconv.Atoi: parsing "12_000": invalid syntax
```

Uzunlik tekshirilmasa (`05-panic`, `./bozor add olma`), chiqish kodi 2:

```
bozor ishga tushdi
versiya: 0.1.0
panic: runtime error: index out of range [3] with length 3

goroutine 1 [running]:
main.main()
	/tmp/…/bozor/main.go:36 +0x933
```

Tekshiruv bilan (`final`, `./bozor add olma`), chiqish kodi 1:

```
bozor ishga tushdi
versiya: 0.1.0
foydalanish: bozor add NOM NARX
```

Nomda probel boʻlsa, qoʻshtirnoq ichida yozing: `./bozor add "oʻrik qoqi" 30000` (`os.Args[2]` = `oʻrik qoqi`).

## Vazifa

[TASK.md](TASK.md): miqdor (`./bozor add olma 12000 5`) va `func parseSold(s string) (int, error)`.
Tekshirish: `go test` (`ok` chiqquncha).

Ochiq qolgan holat (videoda aytiladi): Atoi qabul qiladigan, lekin ×100 da `int64`'dan oshib ketadigan narx
(`100000000000000000`): `parsePrice` manfiy son va `nil` qaytaradi, `main`'dagi `price > 10_000_000*tiyinPerSom`
tekshiruvi manfiy sonni oʻtkazib yuboradi. GB-1.08 uni jadvalli test bilan topadi va tuzatadi.

## Manbalar

- [os.Args](https://pkg.go.dev/os#pkg-variables) · [os.Exit](https://pkg.go.dev/os#Exit) · [fmt.Fprintln](https://pkg.go.dev/fmt#Fprintln) · [fmt.Errorf](https://pkg.go.dev/fmt#Errorf)
- [strconv.Atoi](https://pkg.go.dev/strconv#Atoi) · [strconv.ParseInt](https://pkg.go.dev/strconv#ParseInt) (ErrRange, `_` faqat base 0 da)
- Go spec: [Function declarations](https://go.dev/ref/spec#Function_declarations), [Return statements](https://go.dev/ref/spec#Return_statements), [Errors](https://go.dev/ref/spec#Errors), [Index expressions](https://go.dev/ref/spec#Index_expressions), [Logical operators](https://go.dev/ref/spec#Logical_operators), [Run-time panics](https://go.dev/ref/spec#Run_time_panics)
- [Effective Go: Multiple return values](https://go.dev/doc/effective_go#multiple-returns) · [Go FAQ: exceptions](https://go.dev/doc/faq#exceptions) · [Go Code Review Comments: Handle errors](https://go.dev/wiki/CodeReviewComments#handle-errors)

## Tekshirish (oʻqituvchi uchun)

```bash
gofmt -l . && go vet ./...
go test . ; echo "starter exit $?"     # FAIL, 1
go test ./solution/...                 # ok
```
