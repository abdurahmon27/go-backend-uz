# GB-1.06 · Takrorlash: Go'dagi yagona `for` va uning shakllari

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #6"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

Oldingi dars: [`gb-1-05-qaror-if-switch`](../gb-1-05-qaror-if-switch) · keyingi dars:
[`gb-1-07-funksiya-error`](../gb-1-07-funksiya-error) (funksiya, `os.Args`, `strconv.Atoi`).

**Bu darsdan keyin:** bozor buyurtma summasini sikl bilan hisoblaydi: "har uchinchi dona bepul" aksiyasi har bir
donaga alohida qoʻllanadi. Vositalar: uch qismli `for` (boshlanish; shart; qadam), `i++` va `+=`, `for i := range
sold` (Go 1.22+), faqat shartli va cheksiz `for`, `continue` va `break`.

**Bitta gʻoya:** `for` sarlavhasidagi boshlanish bir marta ishlaydi, shart har aylanishdan *oldin* tekshiriladi, qadam
har aylanishdan *keyin*; `continue` shu aylanishning qolgan qismini tashlab qadamga oʻtadi, `break` esa butun
sikldan chiqadi.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi bozor (vazifa shu yerdan boshlanadi) |
| `main_test.go` | vazifani tekshiradi (`go test`); test yozishni GB-1.08'da ochamiz |
| `TASK.md` | vazifa: aksiyaga chegara (bitta buyurtmada koʻpi bilan 2 dona bepul) |
| `solution/` | javob (xuddi shu test bilan) |
| `_video/NN-…/main.go` | videodagi har bir qadamning kodi; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi, chunki ikkitasida ataylab kompilyatsiya xatosi bor |
| `_video/run.sh` | qadamni videodagidek ishga tushiradi: `./_video/run.sh 03-range-xato` (`go run .`) |

Oldingi darsdan farq (`main.go` GB-1.05'ning `solution/main.go`'sidan boshlangan; faqat `subtotal` sikl bilan
hisoblanadi va oxirida aksiya qatori qoʻshildi):

```bash
diff -u ../gb-1-05-qaror-if-switch/solution/main.go main.go
```

## Videodagi buyruqlar va natijalari (go1.26.8 linux/amd64, `CGO_ENABLED=0`)

Uch qismli `for`, uch dona (`./_video/run.sh 01-sikl`): shart 4 marta tekshiriladi (`i` = 1, 2, 3 → true;
4 → false), tana 3 marta bajariladi:

```
599700
```

`i` faqat `for` ichida yashaydi (`./_video/run.sh 09-scope`):

```
# github.com/abdurahmon27/bozor
./main.go:11:24: undefined: i
```

`for i := range sold`, `sold := 3` (`./_video/run.sh 02-range`): noldan `sold-1` gacha.

```
0
1
2
```

Ishlatilmagan sanagich (`./_video/run.sh 03-range-xato`) va tuzatish `for range sold` (`04-range-tuzat` → `599700`):

```
# github.com/abdurahmon27/bozor
./main.go:10:6: declared and not used: i
```

Faqat shartli `for`: 50 000 soʻmga nechta 1 999 soʻmlik choy sigʻadi (`./_video/run.sh 05-shart`):

```
25 dona
```

`continue`: har uchinchi dona bepul, 7 dona → 5 donaning puli (`06-continue`) va xuddi shu joyda `break`
(`07-break`) → 2 donaning puli:

```
999500
399800
```

Sikl ichidagi `switch`'da `break` faqat `switch`'dan chiqadi (`08-switch-break`, `sold := 4`): `3` ham chiqadi.

```
1
2
3
4
```

Bozor (`go run .`):

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
choy 199900 true 40 1.5
choy 33
choy, 7 dona        9995.00
QQS                 1199.40
yetkazish          15000.00
jami               26194.40
aksiya: 2 dona bepul
```

## Vazifa

[`TASK.md`](TASK.md): har uchinchi dona bepul, lekin bitta buyurtmada koʻpi bilan 2 dona (`const maxFree = 2`).
Tekshirish: shu papkada `go test` (hozir `FAIL`, vazifa bajarilsa `ok`). Javob: [`solution/`](solution/).

Butun papkani bir buyruqda tekshirish (javobni nusxalab):

```bash
t=$(mktemp -d) && cp -r . "$t" && cp solution/*.go "$t"/ && (cd "$t" && go test ./...)
```

## Manbalar

- Go spec: [For statements](https://go.dev/ref/spec#For_statements),
  [For clause](https://go.dev/ref/spec#For_clause), [For range](https://go.dev/ref/spec#For_range),
  [Break statements](https://go.dev/ref/spec#Break_statements),
  [Continue statements](https://go.dev/ref/spec#Continue_statements),
  [Declarations and scope](https://go.dev/ref/spec#Declarations_and_scope),
  [Variable declarations](https://go.dev/ref/spec#Variable_declarations) (ishlatilmagan oʻzgaruvchi),
  [IncDec statements](https://go.dev/ref/spec#IncDec_statements),
  [Assignment statements](https://go.dev/ref/spec#Assignment_statements) (`+=`)
- [Go 1.22 Release Notes](https://go.dev/doc/go1.22#language): butun son ustida `range`
- [A Tour of Go: For](https://go.dev/tour/flowcontrol/1)
