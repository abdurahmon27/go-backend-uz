# GB-1.02 · Qiymat, oʻzgaruvchi va zero value: `var`, `:=` va asosiy turlar

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · Golang backend #2

Oldingi dars: [GB-1.01 · birinchi dastur, go mod init, go run va go build](../gb-1-01-birinchi-dastur/) ·
Keyingi dars: [GB-1.03 · butun sonlar va pul: narx tiyinda, float64 nega emas](../gb-1-03-pul-tiyinda/)

Video: YouTube'dagi "Golang backend" playlist'ida (#2) · Yozma dars: https://www.haywan.uz/lab/youtube/gb/gb-1-02

## Darsda nima boʻladi

bozor'ning `main.go` fayliga birinchi mahsulot yozuvi qoʻshiladi. Avval `var` bilan qiymatsiz eʼlon qilinadi va
zero value koʻrinadi, keyin `=` bilan qiymat beriladi (olma: kilosi 12 000 soʻm, omborda 40 qadoq, bitta qadoq 1.5 kg), soʻng qisqa shakl `:=` va toʻrt asosiy tur: `int`, `float64`,
`string`, `bool`. Oxirida uchta haqiqiy kompilyator xatosi `fayl:qator:ustun` boʻyicha oʻqiladi.

Bu papka oldingi darsning `solution/` holatidan boshlanadi (versiya qatori bilan). Nima oʻzgarganini koʻrish:

```bash
diff -ru ../gb-1-01-birinchi-dastur/solution .
```

## Buyruqlar va natijalar

Hammasi shu papkada, `go1.26.8 linux/amd64`, 2026-11-26. Videoda terminal papkasi `bozor` deb koʻrinadi.
Yangiroq Go versiyasida natijalar bir xil; kompilyator xabarlarining matni versiyadan versiyaga biroz oʻzgarishi mumkin.

```console
$ go run .
bozor ishga tushdi
versiya: 0.1.0
 0 false
olma 12000 true 40 1.5
```

Uchinchi qator boʻsh joy bilan boshlanadi: bu `name`'dagi boʻsh matn (`string`'ning zero value'si), undan keyin
`Println` qoʻygan boʻsh joy. `price` va `active` hali qiymat olmagan: `0` va `false`.

Videodagi uchta xato (har birini `main.go`'ga vaqtincha qoʻshib, `go run .` bilan koʻring, keyin qaytaring):

| Oʻzgarish | Xabar | Nega |
|---|---|---|
| 20-qatorga `fmt.Println(price * weight)` (bitta qadoq narxi) | `./main.go:20:14: invalid operation: price * weight (mismatched types int and float64)` | Ikki oʻzgaruvchi ustidagi amalda turlar bir xil boʻlishi shart (qoʻshtirnoqsiz konstantalar istisno, GB-1.03); Go `int`'ni `float64`'ga oʻzi aylantirmaydi (conversion → GB-1.03) |
| 20-qatorga `stock := stock - 1` | `./main.go:20:8: no new variables on left side of :=` | `:=` eʼlon qiladi, `stock` esa bor. Bor oʻzgaruvchiga `stock = stock - 1` |
| 15-qatorda `price = "12000"` | `./main.go:15:10: cannot use "12000" (untyped string constant) as int value in assignment` | Qoʻshtirnoq matn yasaydi; `price`'ning turi eʼlonda `int` boʻlgan va oʻzgarmaydi |

Har bir xabardan oldin `go` paket nomini ham yozadi: `# github.com/abdurahmon27/bozor`.

## Qisqa qoidalar

- `var name string`: qiymatsiz eʼlon → zero value (`""`, `0`, `false`).
- `var stock int = 40` = `var stock = 40` = `stock := 40`: tur qiymatdan olinadi (`40` → `int`, `1.5` → `float64`).
- `:=` faqat funksiya ichida. Funksiyadan tashqarida oʻzgaruvchi faqat `var` bilan eʼlon qilinadi.
- Eʼlon qilingan oʻzgaruvchi ishlatilmasa: `declared and not used: …` va dastur qurilmaydi.
- Odatiy tanlov: qiymat bor boʻlsa `:=`, zero value bilan boshlanishi kerak boʻlsa `var`.
- `=` yangi oʻzgaruvchi yaratmaydi: bor oʻzgaruvchiga yozadi. Oʻng tomon avval hisoblanadi.

## Vazifa

[TASK.md](TASK.md): sotuvdan keyin qoldiqni hisoblash. Tekshiruv:

```bash
go test
```

Boshlangʻich holatda test yiqiladi va oʻzbekcha xabar beradi (`5-qator «olma 37» boʻlishi kerak, hozir dastur faqat 4 qator chiqaryapti (TASK.md)`).
Vazifa bajarilsa `ok`. `main_test.go` ichini GB-1.08'da ochamiz. Javob: [solution/](solution/)
(`go run ./solution`, `go test ./solution`).

## Manbalar

- Go spec: [Variables](https://go.dev/ref/spec#Variables), [Variable declarations](https://go.dev/ref/spec#Variable_declarations),
  [Short variable declarations](https://go.dev/ref/spec#Short_variable_declarations),
  [The zero value](https://go.dev/ref/spec#The_zero_value), [Constants](https://go.dev/ref/spec#Constants),
  [Operators](https://go.dev/ref/spec#Operators), [Assignment statements](https://go.dev/ref/spec#Assignment_statements)
- [Go FAQ: implicit numeric conversions](https://go.dev/doc/faq#conversions)
- [fmt.Println](https://pkg.go.dev/fmt#Println)
- [Google Go Style Guide: variable declarations](https://google.github.io/styleguide/go/best-practices#variable-declarations)
- [A Tour of Go: basics 8–14](https://go.dev/tour/basics/8) · [Go by Example: Variables](https://gobyexample.com/variables)
