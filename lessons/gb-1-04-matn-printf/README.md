# GB-1.04 · Matn va formatlangan chiqish: `string` va `fmt.Printf`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #4"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

Oldingi dars: [`gb-1-03-pul-tiyinda`](../gb-1-03-pul-tiyinda) · keyingi dars: `gb-1-05` (if, else, switch).

**Bu darsdan keyin:** bozor buyurtma chekini tekis ustunlar bilan chiqarasiz: summa soʻmda, nuqtadan keyin ikki
raqamli tiyin. Vositalar: `string` va `+`, `fmt.Printf` feʼllari (`%s %d %v %q`), kenglik (`%-14s %12s %9d %02d`)
va `fmt.Sprintf`.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi bozor (vazifa shu yerdan boshlanadi) |
| `main_test.go` | vazifani tekshiradi (`go test`); test yozishni GB-1.08'da ochamiz |
| `TASK.md` | vazifa: chekka `yetkazish` qatorini qoʻshish |
| `solution/` | javob (xuddi shu test bilan) |
| `_video/NN-…/main.go` | videodagi har bir qadamning kodi; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi, chunki uchtasida ataylab xato bor (`03-xato` kompilyatsiya boʻlmaydi, `11-string-xato` va `12-feyl-xato` vet'dan oʻtmaydi) |
| `_video/run.sh` | qadamni videodagidek ishga tushiradi: `./_video/run.sh 03-xato` (`go run .`), `./_video/run.sh 11-string-xato vet` (`go vet`) |

Oldingi darsdan farq: `main.go` GB-1.03'ning `solution/main.go`'si bilan bir xil, faqat oxiridagi uchta `Println`
chekka almashgan (yuqoridagi ishga tushish, zero value va ombor qatorlari oʻz joyida):

```bash
diff -ru ../gb-1-03-pul-tiyinda/solution .
```

## Videodagi buyruqlar va natijalari (go1.26.8 linux/amd64, `CGO_ENABLED=0`)

GB-1.03'dagi chek (`./_video/run.sh 01-println`, oxirgi uch qator; undan oldin besh qator: ishga tushish, versiya, zero value, mahsulot, ombor):

```
choy 7 dona
QQS: 1679 soʻm 16 tiyin
jami: 30672 soʻm 16 tiyin
```

Matn va son plyus bilan qoʻshilmaydi (`./_video/run.sh 03-xato`):

```
# github.com/abdurahmon27/bozor
./main.go:8:14: invalid operation: label + total (mismatched types string and int64)
```

`fmt.Printf("%s: %d tiyin\n", name, price)` (`04-printf`) va `\n`'siz ikki chaqiruv (`05-qator`):

```
choy: 199900 tiyin
```
```
choy: 199900 tiyinnon: 400000 tiyin
```

`%v`, `%s`, `%q` (`06-feyllar`; ikkinchi qator oxirida probel bor, lekin koʻrinmaydi):

```
7 dona, sotuvda: true
 choy 
" choy "
```

Tiyin: `%d.%d` (`07-pul`) va `%d.%02d` (`08-pul-tuzat`); choy, 8 dona: QQS = 191 904 tiyin:

```
QQS: 1919.4 soʻm
```
```
QQS: 1919.04 soʻm
```

Kenglik `%-14s %9d` (`09-jadval`, hozircha faqat soʻm qismi):

```
QQS                 1679
yetkazish          15000
jami               30672
```

`fmt.Sprintf("%s, %d dona", name, sold)` (`10-sprintf`): `choy, 7 dona`.

Yakuniy chek (`go run .`; birinchi besh qator GB-1.03'dagidek):

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
choy 199900 true 40 1.5
choy 33
choy, 7 dona       13993.00
QQS                 1679.16
jami               30672.16
```

Ikki keng tarqalgan xato:

```
$ ./_video/run.sh 11-string-xato        # "jami: " + string(total)
jami: �
$ ./_video/run.sh 11-string-xato vet
main.go:7:25: conversion from int64 to string yields a string of one rune, not a string of digits
$ ./_video/run.sh 12-feyl-xato          # fmt.Printf("%d\n", name)
%!d(string=choy)
$ ./_video/run.sh 12-feyl-xato vet
main.go:7:14: fmt.Printf format %d has arg name of wrong type string
```

Toʻgʻri yoʻl: `fmt.Sprintf("%d", total)` yoki toʻgʻridan-toʻgʻri `fmt.Printf`; `strconv.FormatInt(total, 10)` ham shu ishni qiladi (`strconv` paketi GB-1.07'da).

Cheklov: `%d.%02d` faqat manfiy boʻlmagan summa uchun toʻgʻri (`-5` tiyin → `0.-5`); manfiy summani rad etish GB-1.05'da.

## Vazifa

[TASK.md](TASK.md): QQS va jami orasiga `yetkazish` qatorini qoʻshing, `go test` oʻtguncha.

```bash
go test            # hozir FAIL: "chek 4 qator boʻlishi kerak, hozir 3 qator" + kerak/hozir
go test ./solution # ok
```

## Manbalar

- fmt paketi: feʼllar, kenglik, `0` va `-` bayroqlari, format xatolari: https://pkg.go.dev/fmt
- `fmt.Sprintf`: https://pkg.go.dev/fmt#Sprintf
- Go spec, string qoʻshish: https://go.dev/ref/spec#String_concatenation
- Go spec, butun sonli boʻlish va qoldiq: https://go.dev/ref/spec#Integer_operators
- Go spec, butun sondan string'ga oʻgirish: https://go.dev/ref/spec#Conversions_to_and_from_a_string_type
- `go vet` stringintconv: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/stringintconv
- `go test` ishga tushiradigan vet tekshiruvlari: `go help test` (https://pkg.go.dev/cmd/go#hdr-Test_packages)
