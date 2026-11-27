# GB-3.03 · Jamoa asboblari: `go vet`, `staticcheck` va `govulncheck`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #21"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** 2–3-modul darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. **Bu darsda bunga bogʻliq bitta nozik joy
> bor**: asboblarni shu papkaning ichida oʻrnating (pastda, "Asboblarni oʻrnatish").

Oldingi dars: [`gb-3-02-json-saqlash`](../gb-3-02-json-saqlash) (JSON va struct tag'lar: mahsulotlarni
`bozor.json`ga saqlash) · keyingi dars: [`gb-3-04-generics`](../gb-3-04-generics) (generics: tur parametri qachon
foyda beradi).

**Bu darsdan keyin:** har commit'dan oldin kodni uchta statik tahlil asbobi bilan tekshirasiz va har birining
xabarini oʻqib tuzatasiz:

| Asbob | Nimani topadi | Videoda topgani |
|---|---|---|
| `go vet ./...` | Go bilan keladi: shubhali konstruksiyalar (35 ta tekshiruv; `go test` ulardan faqat 11 tasini ishlatadi, test fayli yoʻq paketda ham) | `OldPrice` `price` JSON tag'ini takrorlaydi → `encoding/json` ikkala maydonni jim tashlab yuboradi (`structtag`) |
| `staticcheck ./...` | 100+ tekshiruv, koʻpi vet'da yoʻq | `SA4006`: birinchi `Sell`'ning `err`'i tekshirilmay ustiga yozilgan → qisqa `if err := …; err != nil` |
| `govulncheck ./...` | bogʻliqliklardagi (va standart kutubxonadagi) maʼlum zaifliklar, vuln.go.dev; kodingiz zaif funksiyani chaqiradimi | `golang.org/x/text v0.3.7`: GO-2022-1059, `main` `language.MatchStrings`'ni chaqiradi |

Bu papka 2–3-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i bor. `Store` (`Add`, `List`, `Sell`),
`ErrNotFound`, `ErrBadQuantity` va `*StockError` GB-2.08'ning javobidagidek; `Product`'ga JSON tag'lari (GB-3.02) va
`OldPrice` maydoni (`json:"old_price"`, chegirmagacha narx, tiyin) qoʻshilgan; katalog GB-2.06'dagidek
`golang.org/x/text` bilan soʻmda chiqadi. 3-modul oxiridagi bozor CLI v0 (GB-3.06) labs'ni birlashtiradi.

## Papka

| Fayl | Nima |
|---|---|
| `go.mod`, `go.sum` | modul `github.com/abdurahmon27/bozor`, `go 1.27`; vazifa uchun `golang.org/x/text v0.3.7` (PR eski versiyani qaytargan) |
| `cmd/bozor/main.go` | dastur: katalogni tanlangan tilda chiqaradi (`x/text`), buyurtmani sotadi, katalogni JSON'da yozadi |
| `cmd/bozor/tekshiruv_test.go` | vazifaning 3-qadami: dastur qaysi `x/text` bilan yigʻilganini `debug.ReadBuildInfo` bilan oʻqiydi; oʻzgartirmang |
| `internal/store/store.go` | `Product` (JSON tag'lari bilan), `Store`, `Add`, `List`, `Sell`, `ErrNotFound`, `ErrBadQuantity`, `StockError` |
| `internal/store/rename.go` | jamoadoshning PR'i: `Rename` (vazifa) |
| `internal/store/store_test.go` | jamoaning eski testi (`Sell`), boshidanoq oʻtadi |
| `internal/store/tekshiruv_test.go` | vazifaning 1–2-qadami (`Rename`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: xuddi shu daraxt, oʻz `go.mod`'i bilan (alohida modul, shuning uchun ildizdagi `go test ./...` unga tegmaydi) |
| `_video/NN-…/` | videodagi har bir qadam (toʻliq kichik loyiha); `_` bilan boshlangan papkani `go` buyruqlari chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 00-boshlash go vet ./...` |

Video qadamlari: `00-boshlash` (ikki xato, x/text v0.3.7) → `01-tag` (tag tuzatildi) → `02-err` (qisqa `if`) →
`03-til` (til argumentdan, `MatchStrings`) → `04-yangilash` (x/text v0.44.0, videoning oxiri). Ildizdagi kod =
`04-yangilash` + PR (`rename.go`, `go.mod` yana v0.3.7) + vazifa testlari.

## Asboblarni oʻrnatish

Shu papkaning ichida (yoki istalgan `go 1.27` modulining ichida):

```bash
go install honnef.co/go/tools/cmd/staticcheck@2026.2.2
go install golang.org/x/vuln/cmd/govulncheck@v1.9.0
export PATH="$PATH:$(go env GOPATH)/bin"     # "command not found" boʻlsa
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc   # doimiy (zsh'da ~/.zshrc)
```

- **Nega papka ichida.** Statik tahlil asbobi Go manbasini oʻzi build qilingan Go'ning tur tekshiruvchisi bilan
  oʻqiydi. Kompyuteringizdagi Go eskiroq (masalan 1.26) boʻlsa, papkadan tashqarida `go install` asbobni Go 1.26
  bilan quradi va u Go 1.27 standart kutubxonasini oʻqiy olmaydi (videodagi "Chuqurroq" xatosi). Papka ichida esa
  `go` buyrugʻi `go.mod`'dagi Go 1.27'ga oʻtadi va asbob shu bilan quriladi. Go'ni yangilaganingizdan keyin ham
  asboblarni qayta `go install` qiling.
- **Nega versiyalar qotirilgan.** Jamoada hamma bir xil tekshiruvni koʻrsin. Shuningdek, **2026-11-27 holati:**
  `govulncheck@latest` (v1.10.0, 2026-11-17) kompyuterdagi `go` eskiroq boʻlib, Go 1.27'ga avtomatik oʻtganda
  standart kutubxonani yuklay olmadi (`internal/poll/splice_linux.go:237:21: use of promoted field … requires go1.27 or
  later`); Go 1.27 toʻgʻridan-toʻgʻri oʻrnatilgan boʻlsa, v1.10.0 ham ishlaydi. v1.9.0 ikkala holda ham ishladi.
  Yangi versiya chiqqanda qayta tekshirib, versiyani oshiring.
- staticcheck'da reliz nomi `2026.2.2`, Go moduli versiyasi `v0.8.2` — bu bitta reliz. `go install` bozor'ning
  `go.mod`'iga tegmaydi.
- Jamoada asbobni `go.mod`'ga yozish ham mumkin (Go 1.24'dan beri): `go get -tool honnef.co/go/tools/cmd/staticcheck@2026.2.2`,
  keyin `go tool staticcheck ./...`. Versiya `go.mod`'da qotadi va asbob har doim joriy Go bilan quriladi. Bu darsda
  `go install` ishlatildi: `go.mod` faqat bozor'ning oʻz bogʻliqliklarini saqlaydi.
- `govulncheck` bazani vuln.go.dev'dan oladi: internet kerak.

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`, staticcheck 2026.2.2, govulncheck v1.9.0, vuln DB 2026-11-24)

Boshlanishi (`00-boshlash`): kompilyatsiya boʻladi, test oʻtadi, lekin buyurtma "qabul qilindi" (olma 50 dona, qoldiq
40) va JSON'da narx yoʻq:

```
$ go run ./cmd/bozor
1. olma: 12 000 soʻm
2. non: 4 000 soʻm
buyurtma qabul qilindi
[{"id":1,"name":"olma","stock":40},{"id":2,"name":"non","stock":97}]
$ go test ./...
?   	github.com/abdurahmon27/bozor/cmd/bozor	[no test files]
ok  	github.com/abdurahmon27/bozor/internal/store	0.003s
$ go vet ./...
internal/store/store.go:16:2: struct field OldPrice repeats json tag "price" also at store.go:15
```

`go test` vet'ning faqat 11 ta tekshiruvini ishlatadi (`go help test`: atomic, bools, buildtag, directive, errorsas,
ifaceassert, nilfunc, printf, stdversion, stringintconv, tests). `structtag` ular orasida yoʻq. `go tool vet help`
hammasini (35 ta) koʻrsatadi; `go vet` sukut boʻyicha hammasini ishlatadi. Ikki maydonda bir xil tag boʻlsa,
`encoding/json` ikkalasini ham tashlaydi va xato bermaydi; `Unmarshal` ham: `{"id":1,"price":5}` → `Price` 0, `err` nil.

Tag tuzatildi (`01-tag`: `OldPrice int64 \`json:"old_price"\``; GB-3.02 qoidasi: songa `omitempty` qoʻyilmaydi):

```
$ go vet ./... && go run ./cmd/bozor
1. olma: 12 000 soʻm
2. non: 4 000 soʻm
buyurtma qabul qilindi
[{"id":1,"name":"olma","price":1200000,"old_price":0,"stock":40},{"id":2,"name":"non","price":400000,"old_price":0,"stock":97}]
```

staticcheck (`01-tag`; vet bu xatoni koʻrmaydi):

```
$ staticcheck ./...
bash: staticcheck: command not found          # go/bin PATH'da boʻlmasa
$ export PATH="$PATH:$(go env GOPATH)/bin"
$ staticcheck ./...
cmd/bozor/main.go:26:2: this value of err is never used (SA4006)
$ staticcheck -explain SA4006
A value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?

Available since
    2017.1

Online documentation
    https://staticcheck.dev/docs/checks#SA4006
```

Qisqa `if` bilan (`02-err`): `err` faqat shu `if` ichida, birini ikkinchisining ustiga yozib boʻlmaydi:

```go
if err := st.Sell(1, 50); err != nil {
	fmt.Fprintln(os.Stderr, "xato:", err)
	os.Exit(1)
}
```

```
$ staticcheck ./... && go run ./cmd/bozor
1. olma: 12 000 soʻm
2. non: 4 000 soʻm
xato: olma: qoldiq 40, soʻralgan 50
exit status 1
```

govulncheck: `x/text v0.3.7`, `main` faqat `message.NewPrinter(language.Uzbek)`'ni ishlatadi (`02-err`):

```
$ govulncheck ./...
=== Symbol Results ===

No vulnerabilities found.

Your code is affected by 0 vulnerabilities.
This scan also found 1 vulnerability in packages you import and 2
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.
$ echo $?
0
```

`-show verbose` uchala darajani koʻrsatadi: funksiya (Symbol) darajasida hech narsa, paket darajasida GO-2022-1059,
modul darajasida GO-2026-6629 (fixed v0.41.0) va GO-2026-5970 (fixed v0.39.0); boshida "Govulncheck scanned the
following 2 modules and the go1.27.4 standard library". Til argumentdan tanlanadi, `language.MatchStrings`
chaqiriladi (`03-til`, `main.go`'ning 21–27-qatorlari):

```
$ go doc language.MatchStrings
package language // import "golang.org/x/text/language"

func MatchStrings(m Matcher, lang ...string) (tag Tag, index int)
    MatchStrings parses and matches the given strings until one of them matches
    the language in the Matcher. A string may be an Accept-Language header as
    handled by ParseAcceptLanguage. The default language is returned if no other
    language matched.
$ govulncheck ./...
=== Symbol Results ===

Vulnerability #1: GO-2022-1059
    Denial of service via crafted Accept-Language header in
    golang.org/x/text/language
  More info: https://pkg.go.dev/vuln/GO-2022-1059
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.3.7
    Fixed in: golang.org/x/text@v0.3.8
    Example traces found:
      #1: cmd/bozor/main.go:26:33: bozor.main calls language.MatchStrings

Your code is affected by 1 vulnerability from 1 module.
This scan also found 0 vulnerabilities in packages you import and 2
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.
$ echo $?
3
$ go get golang.org/x/text@latest
go: upgraded golang.org/x/text v0.3.7 => v0.44.0
$ govulncheck ./...
No vulnerabilities found.
```

Commit'dan oldingi qator (`04-yangilash`, videoning oxiri):

```
$ go vet ./... && staticcheck ./... && govulncheck ./... && go test ./...
No vulnerabilities found.
?   	github.com/abdurahmon27/bozor/cmd/bozor	[no test files]
ok  	github.com/abdurahmon27/bozor/internal/store	0.006s
```

Oraliq savol (videoda faqat javobi): test fayli yoʻq `cmd/bozor`'da `fmt.Printf("%s\n", p.Name, p.Stock)` boʻlsa,
`go build` oʻtadi, `go test` esa yiqiladi, chunki `printf` 11 talik roʻyxatda va `go test` uni har paketda ishlatadi:

```
# github.com/abdurahmon27/bozor/cmd/bozor
cmd/bozor/main.go:48:31: fmt.Printf call needs 1 arg but has 2 args
FAIL	github.com/abdurahmon27/bozor/cmd/bozor [build failed]
ok  	github.com/abdurahmon27/bozor/internal/store	0.009s
FAIL
```

### Chuqurroq: Go yangilangandan keyin eski staticcheck

Go 1.26.8 bilan qurilgan staticcheck 2026.2.2'ni Go 1.27.4 loyihasida ishga tushirsak
(`GOTOOLCHAIN=go1.26.8 go install …@2026.2.2`):

```
$ go version
go version go1.27.4 linux/amd64
$ staticcheck ./...
-: …/src/internal/poll/splice_linux.go:237:21: unknown field rfd in struct literal of type splicePipe (compile)
```

Xato sizning kodingizda emas, Go'ning oʻz manbasida: Go 1.27 standart kutubxonasi yangi til imkoniyatini ishlatadi
(struct literal'da ichki struct maydoniga toʻgʻridan-toʻgʻri murojaat; staticcheck 2026.2 release notes: "struct
initializers with direct references to embedded fields"), eski Go bilan qurilgan asbobning tur tekshiruvchisi uni
tanimaydi. Yechim: `go install …@2026.2.2`'ni yangi Go bilan (shu papka ichida) qayta ishga tushirish.

## Vazifa

[TASK.md](TASK.md): jamoadoshning `Rename` PR'ini `go vet`, `staticcheck` va `govulncheck` bilan tekshirib, har
birining topganini tuzating. Tekshirish: `go test ./...` (hozir FAIL, vazifa bajarilgach `ok`). Javob bilan farq:
`diff -ru internal solution/internal`, `diff go.mod solution/go.mod`.

## Manbalar

- `go vet` — https://pkg.go.dev/cmd/vet · `go help test` (vet subset) — https://pkg.go.dev/cmd/go#hdr-Test_packages
- `encoding/json` Marshal (bir xil nomli maydonlar, `omitempty`) — https://pkg.go.dev/encoding/json#Marshal
- staticcheck — https://staticcheck.dev/docs/ · SA4006 — https://staticcheck.dev/docs/checks/#SA4006 ·
  2026.2 release notes — https://staticcheck.dev/changes/2026.2/
- Go zaifliklarini boshqarish — https://go.dev/doc/security/vuln/ · govulncheck — https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck ·
  rasmiy qoʻllanma — https://go.dev/doc/tutorial/govulncheck
- GO-2022-1059 (CVE-2022-32149) — https://pkg.go.dev/vuln/GO-2022-1059
- `go install`, `go get`, GOPATH, `go get -tool` — https://pkg.go.dev/cmd/go · https://go.dev/ref/mod#go-get
- Go spec: If statements (qisqa `if`) — https://go.dev/ref/spec#If_statements
- Go Code Review Comments: Handle errors — https://go.dev/wiki/CodeReviewComments#handle-errors
- Bash: Command Substitution, Single Quotes, Redirections, Startup Files — https://www.gnu.org/software/bash/manual/bash.html
