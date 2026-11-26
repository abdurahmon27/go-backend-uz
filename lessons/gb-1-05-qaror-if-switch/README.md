# GB-1.05 · Qaror: `if`, `else` va `switch`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #5"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

Oldingi dars: [`gb-1-04-matn-printf`](../gb-1-04-matn-printf) · keyingi dars: [`gb-1-06-for-sikl`](../gb-1-06-for-sikl).

**Bu darsdan keyin:** bozor notoʻgʻri miqdorni va ombordagidan koʻp buyurtmani chek chiqarishdan (va ombordan
ayirishdan) oldin rad etadi, yetkazish narxini summaga qarab tanlaydi. Vositalar: taqqoslash (`== != < <= > >=`,
natijasi `bool`), mantiqiy `&& || !`, `if` / `else if` / `else`, `return` bilan erta chiqish, qisqa eʼlonli `if` va
uning koʻrinish doirasi (scope), `switch` (birinchi true `case` bajariladi, keyingilariga oʻtilmaydi; `default`).

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi bozor (vazifa shu yerdan boshlanadi) |
| `main_test.go` | vazifani tekshiradi (`go test`): `main.go`'dagi besh qiymatni (`name`, `price`, `active`, `stock`, `sold`) almashtirib, dasturni har xil buyurtma bilan ishga tushiradi; test yozishni GB-1.08'da ochamiz |
| `TASK.md` | vazifa: nom, sotuv (`active`) va narx qoidasi |
| `solution/` | javob (xuddi shu test bilan) |
| `_video/NN-…/main.go` | videodagi har bir qadamning kodi; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi, chunki ikkitasida ataylab xato bor; `06-teng-xato`'da kod `main.go.txt` nomida, chunki sintaksis xatosini `gofmt` ham oʻqiy olmaydi |
| `_video/run.sh` | qadamni videodagidek ishga tushiradi: `./_video/run.sh 04-if` |

`main.go` GB-1.04'ning `solution/main.go`'si bilan bayt-baytigacha bir xil, faqat shu darsning oʻzgarishi qoʻshilgan:
miqdor va ombor qoidasi (`sold := 7`'dan keyin, `stock = stock - sold`'dan oldin) va yetkazish narxini tanlaydigan
`switch` (`delivery`, jami va chekdagi yetkazish qatori shundan olinadi):

```bash
diff -u ../gb-1-04-matn-printf/solution/main.go main.go
```

**Faqat shu bosqich uchun soddalashtirish:** xato xabari oddiy chiqishga (stdout) yoziladi va dastur 0 kodi bilan,
yaʼni "muvaffaqiyatli" tugaydi. Haqiqiy buyruq xatoni xato oqimiga (stderr) yozib, noldan farqli kod bilan chiqadi:
GB-1.07'da `fmt.Fprintln(os.Stderr, …)` va `os.Exit(1)`.

## Videodagi buyruqlar va natijalari (`CGO_ENABLED=0`, linux/amd64)

GB-1.04'dagi bozor, faqat `sold := -3` (`./_video/run.sh 01-manfiy`): kompilyator jim, ombor 43 ga oʻsadi, chek
manfiy:

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
choy 199900 true 40 1.5
choy 43
choy, -3 dona      -5997.00
QQS                 -719.-64
yetkazish          15000.00
jami                8283.36
```

`ok := sold > 0` (`02-bool`, sold = -3) → `false`; `03-mantiq` (sold 7, stock 40:
`ok := sold > 0 && sold <= stock`, `bad := sold <= 0 || sold > stock`, `fmt.Println(ok, bad, !ok)`) →
`true false false`.

`if sold <= 0 { fmt.Println("xato: …") }` `return`'siz (`04-if`): xato qatori chiqadi, keyin baribir ombor va chek:

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
choy 199900 true 40 1.5
xato: miqdor kamida 1 dona boʻlishi kerak
choy 43
choy, -3 dona      -5997.00
QQS                 -719.-64
yetkazish          15000.00
jami                8283.36
```

`return` bilan (`05-return`):

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
choy 199900 true 40 1.5
xato: miqdor kamida 1 dona boʻlishi kerak
```

Keng tarqalgan xato, shartda bitta `=` (`06-teng-xato`, `if sold = 0 {`):

```
# github.com/abdurahmon27/bozor
./main.go:26:10: syntax error: cannot use assignment sold = 0 as value
```

Ombor qoidasi, qisqa eʼlonli `if left := stock - sold; left < 0`, `sold := 50`, `stock := 40` (`07-ombor`):

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
choy 199900 true 40 1.5
xato: omborda 10 dona yetmaydi
```

`left` blokdan keyin (`08-scope`):

```
# github.com/abdurahmon27/bozor
./main.go:34:14: undefined: left
```

Yetkazish narxi `else if` (`09-elseif`) va `switch` (`10-switch`) bilan, oʻrik 25 000 soʻm × 24 = 600 000 soʻm
(ikkala shart ham true, faqat birinchisi bajariladi):

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
oʻrik 2500000 true 40 1.5
oʻrik 16
oʻrik, 24 dona    600000.00
QQS                72000.00
yetkazish              0.00
jami              672000.00
```

Qiymatli `switch name { case "olma", "oʻrik": … }` (`11-qiymat`, name = "oʻrik") → `meva`.

Yakuniy bozor (`go run .`, choy 7 dona):

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
choy 199900 true 40 1.5
choy 33
choy, 7 dona       13993.00
QQS                 1679.16
yetkazish          15000.00
jami               30672.16
```

## Vazifa

[TASK.md](TASK.md): `sold := 7`'dan keyin, miqdor tekshiruvidan oldin nom (`name == ""`), sotuv (`!active`) va narx
(`price <= 0 || price > 10 000 000 soʻm`) qoidalarini qoʻshing, `go test` oʻtguncha.

```
$ go test
--- FAIL: TestQoidalar (0.16s)
    main_test.go:210: nom boʻsh (name="" price=199900 active=true stock=40 sold=7):
          kerak:
            "xato: mahsulot nomi boʻsh"
          hozir:
            " 33"
            ", 7 dona           13993.00"
            "QQS                 1679.16"
            "yetkazish          15000.00"
            …
$ go test ./solution
ok  	github.com/abdurahmon27/bozor/solution
```

## Manbalar

- Go spec, taqqoslash operatorlari: https://go.dev/ref/spec#Comparison_operators
- Go spec, mantiqiy operatorlar: https://go.dev/ref/spec#Logical_operators
- Go spec, operatorlar ustuvorligi: https://go.dev/ref/spec#Operator_precedence
- Go spec, `if` (qisqa eʼlon, `else`): https://go.dev/ref/spec#If_statements
- Go spec, `switch` (ifodasiz switch = `true`, birinchi mos case, `fallthrough`): https://go.dev/ref/spec#Expression_switches
- Go spec, bloklar va koʻrinish doirasi: https://go.dev/ref/spec#Blocks
- Go spec, `return`: https://go.dev/ref/spec#Return_statements
- Go spec, butun sonli boʻlish va qoldiq (manfiy son): https://go.dev/ref/spec#Integer_operators
- Effective Go, `if` va `switch`: https://go.dev/doc/effective_go#if · https://go.dev/doc/effective_go#switch
- Go Code Review Comments, Indent Error Flow: https://go.dev/wiki/CodeReviewComments#indent-error-flow
