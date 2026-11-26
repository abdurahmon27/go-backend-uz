# GB-2.10 · UTF-8: bayt, rune va oʻzbekcha matn

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #18"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-2-09` (`defer`, `panic` va `recover`) · keyingi dars: `gb-3-01` (fayllar, `bufio` va
`io.Reader`). Bu darsga kerak boʻlganlar: [`gb-1-04-matn-printf`](../gb-1-04-matn-printf) (`string`, `Printf`
feʼllari, kenglik `%-10s`; u yerda `rune` soʻzi shu darsga qoldirilgan),
[`gb-2-01-slice`](../gb-2-01-slice) (`len`, `for range`), [`gb-2-02-map-qidirish`](../gb-2-02-map-qidirish)
(`os.Args[1:]`), [`gb-2-07-interface-writer`](../gb-2-07-interface-writer) (`byte`, `Write(p []byte)`).

**Bu darsdan keyin:** `len` string'da baytlarni sanashini, UTF-8 bitta belgini 1–4 baytga yozishini
(`oʻrik` = 5 rune, 6 bayt) koʻrasiz; `for i, r := range name` bilan matnni belgilar (rune'lar) boʻyicha
aylanasiz; `unicode/utf8` paketidagi `RuneCountInString` va `ValidString`'ni ishlatasiz va bozor chekidagi
mahsulot nomini rune chegarasida kesib, harf yoʻqolishi va singan belgi xatosini tuzatasiz.

Bu papka 2-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`,
`go 1.27`), oʻz `package main`'i. 3-modul oxiridagi bozor CLI v0 (GB-3.06) modul darslarini birlashtiradi.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: rune'lar bilan kesadigan `short` va chek (`main`); oxirida vazifaning ikki boʻsh funksiyasi (`shortLabel`, `checkName`) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 09-tuzat go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Muammo: chek nomni `name[:10]` bilan kesadi, `len` bilan tekshiradi (`./_video/run.sh 00-chek go run .`).
`oʻrik qoqi` 10 rune, lekin 11 bayt: oxirgi harf yoʻqoldi; `paxta yogʻi`da kesish `ʻ`ning oʻrtasiga tushdi
(terminal yolgʻiz `ca` baytini `�` qilib chizadi):

```
[olma      ]
[oʻrik qoq ]
[paxta yog�]
```

`len("oʻrik")` (`01-len`) va baytlarning oʻzi `fmt.Printf("% x\n", name)` (`02-bayt`):

```
6
6f ca bb 72 69 6b
```

`o` = `6f`, `ʻ` (U+02BB, MODIFIER LETTER TURNED COMMA) = `ca bb`, `r i k` = `72 69 6b`. Indeks bayt beradi
(`03-indeks`, `fmt.Println(name[0], name[1])`):

```
111 202
```

`for i, r := range name { fmt.Printf("%d %c\n", i, r) }` (`04-range`): `i` — rune boshlangan bayt, `r` — rune:

```
0 o
1 ʻ
3 r
4 i
5 k
```

`rune` (= `int32`), rune literal `'ʻ'` va `%d %c %U` (`05-rune`):

```
699 ʻ U+02BB
```

`import ( "fmt" "unicode/utf8" )` va `utf8.RuneCountInString(name)` (`06-count`):

```
5
```

Oraliq savol: `fmt.Println(len("…"), utf8.RuneCountInString("…"))` (`07-quiz`):

```
3 1
```

Xato: baytlar bilan kesish (`08-kesish`, `cut := short("oʻrik", 2)`, `%q` va `utf8.ValidString`):

```
"o\xca"
false
```

Tuzatish: `short` rune'larni sanaydi va `name[:i]` bilan rune chegarasida kesadi (`09-tuzat`, papkadagi
`main.go` ham shu). `Printf` kengligi rune'larda oʻlchanadi, shuning uchun ustun tekis:

```
[olma      ]
[oʻrik qoqi]
[paxta yogʻ]
```

Chuqurroq: `[]rune` va `string(...)` (`10-runes`, `runes := []rune("paxta yogʻi")`,
`fmt.Println(len(runes), string(runes[:10]))`):

```
11 paxta yogʻ
```

Chuqurroq: oddiy apostrof va `ʻ` (`11-apostrof`, `"oʻrik" == "o'rik"`, `len` ikkalasiniki):

```
false
6 5
```

## Vazifa

[TASK.md](TASK.md): `shortLabel` (sigʻmasa `width-1` rune + `…`) va `checkName` (boʻsh, UTF-8 emas yoki 20
rune'dan uzun → xato). Hali qamrab olinmagan (keyingi darslar): faqat probellardan iborat nom, `oʻ`/`gʻ`'ni `o`
va `ʻ` orasidan kesmaslik, foydalanuvchi yozgan apostrof turlarini (`'` `‘` `’` `ʼ`) bitta `ʻ`ga keltirish. Tekshirish: shu papkada `go test`, `ok` chiqquncha. Hozir (boshlangʻich kod) `FAIL` bilan
yiqiladi.

Oldingi darsdan nima oʻzgargani: 2-modul darslari mustaqil lab (koordinator qoidasi), shuning uchun bu papka
gb-2-09'ning snapshot'i emas; `Product` turi bu darsga kerak emas, nomlar oddiy `[]string`. Toʻliq bozor GB-3.06'da
birlashadi.

Papkani tekshirish (muallif uchun):

```bash
gofmt -l . && go vet ./...                 # hech narsa chiqmaydi
go test . ; echo "starter exit $?"         # FAIL, exit 1
go test ./solution/...                     # ok
```

## Manbalar

- Go spec: [Source code representation](https://go.dev/ref/spec#Source_code_representation) (manba UTF-8),
  [String types](https://go.dev/ref/spec#String_types) (`len` = baytlar soni, string oʻzgarmas),
  [Rune literals](https://go.dev/ref/spec#Rune_literals), [Numeric types](https://go.dev/ref/spec#Numeric_types)
  (`rune` = `int32`), [Index expressions](https://go.dev/ref/spec#Index_expressions) (`s[i]` → `byte`),
  [Slice expressions](https://go.dev/ref/spec#Slice_expressions),
  [For statements with range clause](https://go.dev/ref/spec#For_range) (string ustida: bayt indeksi + rune, notoʻgʻri
  UTF-8 → `0xFFFD`), [Conversions to and from a string type](https://go.dev/ref/spec#Conversions_to_and_from_a_string_type)
- Go blog: [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings) (Rob Pike)
- [pkg.go.dev/unicode/utf8](https://pkg.go.dev/unicode/utf8) (`RuneCountInString`, `ValidString`, `RuneError`)
- [pkg.go.dev/fmt](https://pkg.go.dev/fmt#hdr-Printing) (`%x` va probel bayrogʻi, `%c`, `%U`, `%q`; kenglik
  rune'larda oʻlchanadi)
- [RFC 3629: UTF-8](https://www.rfc-editor.org/rfc/rfc3629) · [Unicode U+02B0–U+02FF](https://www.unicode.org/charts/PDF/U02B0.pdf)
  (U+02BB MODIFIER LETTER TURNED COMMA)
