# 06 · Funksiyalar: ikki qiymat qaytarish, closure va variadic

Oldingi dars: [05 · Go'da `while` yoʻq: bitta `for`](../005-for-switch/) · Keyingi dars: 07 · Taxmin qilmang, oʻlchang: birinchi test va benchmark

| Papka | Nima |
|---|---|
| `port/` | `PORT=9000`, server esa `8080` da: `:=` shadowing xatosi (dastur ishlaydi, lekin natija notoʻgʻri) |
| `portfix/` | tuzatilgan versiya |
| `ikki/` | ikki qiymat qaytarish: `(int, error)`, `_`, `strconv.Atoi` |
| `qiymat/` | funksiya qiymat sifatida: parametr, nomsiz funksiya, `%T` |
| `taqqos/` | `f == g`: **ataylab xato**, kompilyatsiya boʻlmaydi (funksiyani faqat `nil` bilan solishtirish mumkin; videoda yoʻq) |
| `hisob/` | closure: hisoblagich |
| `pauza/` | pauza-savol: closure qiymatni emas, oʻzgaruvchini ushlaydi |
| `sikl/`, `eski/` | sikldagi closure: Go 1.22+ da `0 1 2`, `go 1.21` modulda `3 3 3` (`eski/` alohida modul) |
| `mediana/` | variadic: `kechikish...` bilan yangi slice yaratilmaydi, `slices.Sort` chaqiruvchining slice'ini saralab yuboradi |
| `tengxato/` | `:=` oʻrniga `=`: **ataylab xato**, `undefined: err` |
| `ishlatilmagan/` | ichki `port` ishlatilmasa: **ataylab xato**, `declared and not used: port` |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
```bash
PORT=9000 go run ./port
go vet ./port                     # hech narsa demaydi
go run ./ikki
go run ./qiymat
go run ./hisob
go run ./pauza                    # pauza-savol javobi
go run ./sikl
cd eski && go run . && cd ..
go run ./mediana
go run ./tengxato                 # ataylab xato
go run ./ishlatilmagan            # ataylab xato
GOPROXY=file://$(go env GOMODCACHE)/cache/download GOSUMDB=off go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@v0.53.0
go vet -vettool=$(go env GOPATH)/bin/shadow ./port
PORT=9000 go run ./portfix
```
Natijalar `go1.26.8 linux/amd64` da olingan.

## Dars matni

**Bir nechta qiymat.** Go funksiyasi bir nechta qiymat qaytara oladi: `func bolish(a, b int) (int, error)`.
Go'da `try`/`catch` yoʻq: oddiy xatolar qiymat sifatida qaytadi, odatda oxirgi oʻrinda, hammasi joyida boʻlsa `nil`
([Go FAQ](https://go.dev/doc/faq#exceptions)). `panic`/`recover` bor, lekin ular favqulodda holatlar uchun (15-dars).
`10 / 3` butun sonlarda `3` (butun boʻlish). Kerak boʻlmagan qiymat oʻrniga `_` yoziladi.
Standart kutubxona ham shunday: `strconv.Atoi("12a")` → `0 strconv.Atoi: parsing "12a": invalid syntax`.

```
3 <nil>
nolga boʻlib boʻlmaydi
0 strconv.Atoi: parsing "12a": invalid syntax
```

Natijalarga nom ham berish mumkin: `func bolish(a, b int) (q int, err error)`, keyin boʻsh `return` ularni qaytaradi.
Qisqa funksiyalarda qulay, uzun funksiyada esa nima qaytayotganini koʻrish qiyinlashadi
([Go Code Review Comments: Named Result Parameters](https://go.dev/wiki/CodeReviewComments#named-result-parameters)).

**Funksiya qiymat.** Funksiyani oʻzgaruvchiga yozish, parametr qilib berish va qaytarish mumkin. Uning turi bor:
`func(string) string`. `hammasiga(sozlar, strings.ToUpper)` → `[OLMA NOK UZUM]`. Backend'da: `http.HandleFunc(pattern, handler)`
ikkinchi argument sifatida `func(ResponseWriter, *Request)` oladi. Lekin ikki funksiyani solishtirib
boʻlmaydi: `f == g` → `invalid operation: f == g (func can only be compared to nil)`
([spec: Comparison operators](https://go.dev/ref/spec#Comparison_operators)).

**Closure.** `hisoblagich()` ichidagi `n` funksiya tugagandan keyin ham yashaydi, chunki qaytgan nomsiz funksiya uni
ushlab qolgan. `hisoblagich()` ning har chaqiruvida yangi `n`: `a()` → `1 2 3`, `b()` → `1`. Spec: *"Function literals are closures: they may
refer to variables defined in a surrounding function. Those variables are then shared between the surrounding function
and the function literal, and they survive as long as they are accessible."*
`n` qayerda saqlanadi? Bu spetsifikatsiyada emas, kompilyatorda hal qilinadi (escape analysis, 11-dars).
`hisoblagich` ning alohida (inline qilinmagan) nusxasida `n` heap'ga koʻchadi, buni `-m` ham, `-m -l` ham koʻrsatadi:

```
$ go build -o /dev/null -gcflags='-m -l' ./hisob 2>&1 | grep heap | head -2
hisob/main.go:6:2: moved to heap: n
hisob/main.go:7:9: func literal escapes to heap
```

Lekin oddiy `-gcflags=-m` da `hisoblagich` `main` ichiga inline qilinadi va u yerda `func literal does not escape`:
`a` va `b` ning `n` lari `main` ning stekida qoladi (`-gcflags=-m=2`: `main capturing by ref: n`, heap'ga koʻchirish yoʻq).
Shuning uchun "closure = doim heap" degan qoida yoʻq.

**Pauza-savol.** `x := 1; f := func() int { return x * 10 }; x = 5; fmt.Println(f())` → **50**. Closure qiymatning
nusxasini emas, oʻzgaruvchining oʻzini ushlaydi.

**Sikldagi closure.** Shu sababli 05-darsdagi tuzoq closure bilan ham chiqar edi:

```go
for i := 0; i < 3; i++ {
	ishlar = append(ishlar, func() { fmt.Print(i, " ") })
}
```

Go 1.22+ (`go.mod`: `go 1.26`): `0 1 2`. `go 1.21` modulda: `3 3 3`. `go vet` ikkalasida ham jim (uning `loopclosure`
tekshiruvi faqat `go`/`defer`, errgroup `Go` va `t.Parallel` li `t.Run` ichidagi closure'larni, va faqat Go 1.22 dan
oldingi semantikali fayllarda tekshiradi; `append` ga berilgan closure'ni koʻrmaydi).

**Variadic.** `func mediana(vaqtlar ...int) int`: ichkarida `vaqtlar` oddiy `[]int`. `fmt.Println(a ...any)` ham
variadic. `mediana(90, 10, 50)` → 50 (bu chaqiruvda yangi slice yaratiladi). Tayyor slice'ni `kechikish...` bilan
bersangiz, **yangi slice yaratilmaydi** (spec: *"it is passed unchanged as the value for a ...T parameter. In this case no
new slice is created."*), va `slices.Sort` chaqiruvchining slice'ini saralab yuboradi (`mediana` soddalashtirilgan: boʻsh chaqiruvda `index out of
range` panic, juft sonli roʻyxatda esa haqiqiy mediana emas):

```
50
50
[10 50 90]
```

Go'da hamma argument nusxalanadi, lekin slice'ning nusxasi ham oʻsha massivga qaraydi (slice'lar: 08-dars). Funksiya
qanday chaqirilganini bilmaydi, shuning uchun variadic parametrni oʻzgartirmang; kerak boʻlsa `slices.Clone` bilan nusxa
oling. `...` eng koʻp `append(a, b...)` da uchraydi.

**`:=` shadowing.** `port/main.go`:

```go
port := 8080
if s := os.Getenv("PORT"); s != "" {
	port, err := strconv.Atoi(s) // if ichida YANGI port
	...
	fmt.Println("sozlandi:", port) // 9000
}
fmt.Println("server ishga tushdi:", port) // 8080
```

**Nega `:=`?** Chunki `err` yangi: `port, err = strconv.Atoi(s)` → `undefined: err` (`tengxato/`). `:=` esa `err`
bilan birga `port` ni ham yangi qiladi. Kompilyator nega jim? Ichki `port` `Println` da ishlatilgan; oʻsha qatorni
oʻchirsangiz: `declared and not used: port` (`ishlatilmagan/`).

Qoida ([spec: Short variable declarations](https://go.dev/ref/spec#Short_variable_declarations)): `:=` chap tomonda
kamida bitta yangi nom boʻlsa ishlaydi; **shu blokda** eʼlon qilingan nom qayta ishlatiladi (yangi oʻzgaruvchi emas),
**tashqi blokdagi** nom esa yangisi bilan toʻsiladi. `if` bloki yangi blok, shuning uchun `port` ham, `err` ham yangi.

Nomlangan natijalar bilan esa kompilyator yordam beradi: `func portOl() (port int, err error)` ichida, `port` toʻsilgan
blokda boʻsh `return` yozsangiz → `result parameter port not in scope at return` (spec: Return statements,
*implementation restriction*).

**Temporal, 2021.** [How Go shadowing and bad choices caused our first data loss bug in years](https://temporal.io/blog/go-shadowing-bad-choices):
tashqi `var err error`; `if request.IsNewBranch {` ichida `treeInfoDataBlob, err := …` yangi `err` yaratdi; bazaga yozish
xatosi keyin `err = h.session.ExecuteBatch(batch)` bilan **ichki** `err` ga tushdi; tashqi `err` `nil` qoldi va chaqiruvchiga
muvaffaqiyat qaytdi. Ular kodni qayta tuzib (ichki blokda `err` ni oʻrnatmaydigan qilib) tuzatdi; keyingi choralar
roʻyxatida testlar va shadowing'ni topadigan linter bor (guru vazifasi shu xatoning soddalashtirilgan shakli).

`go vet` buni koʻrmaydi (`cmd/vet/main.go`: `// shadow.Analyzer omitted: too noisy`). Go jamoasining `shadow` analizatori koʻradi, lekin u *"generates too many false positives and is
not yet enabled by default"* ([shadow](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow)):

```
port/main.go:12:3: declaration of "port" shadows declaration at line 10
```

`shadow@v0.53.0` Go 1.26 talab qiladi. VS Code'da gopls'ning `shadow` analizatorini sozlamalardan yoqish mumkin.

Tuzatish (`portfix/`): natijani yangi nomga oling (`p, err := …`) va tashqi oʻzgaruvchiga `=` bilan yozing (`port = p`).
Ikkinchi yoʻl: `var err error` deb eʼlon qilib, `port, err = strconv.Atoi(s)`. Faqat `:=` ni `=` ga almashtirish yetmaydi.

## Manbalar
- Go spec: [Function types](https://go.dev/ref/spec#Function_types), [Function literals](https://go.dev/ref/spec#Function_literals),
  [Passing arguments to `...` parameters](https://go.dev/ref/spec#Passing_arguments_to_..._parameters),
  [Short variable declarations](https://go.dev/ref/spec#Short_variable_declarations),
  [Declarations and scope](https://go.dev/ref/spec#Declarations_and_scope),
  [Comparison operators](https://go.dev/ref/spec#Comparison_operators), [For statements](https://go.dev/ref/spec#For_statements)
- [Go FAQ: Why does Go not have exceptions?](https://go.dev/doc/faq#exceptions)
- [Go blog: Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview)
- [`shadow` analizatori](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow) (golang.org/x/tools v0.53.0)
- Ryland Goldstein, [Temporal: Go shadowing data loss bug](https://temporal.io/blog/go-shadowing-bad-choices) (2021)
- [golang/go#377: proposal: spec: various changes to :=](https://github.com/golang/go/issues/377) (2009-dan beri ochiq)
- Teiva Harsanyi, [100 Go Mistakes #1: Unintended variable shadowing](https://100go.co/#unintended-variable-shadowing-1)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments#named-result-parameters)

## Uy vazifasi
- 🟢 **Oson:** slice'dagi eng kichik va eng katta sonni xato bilan birga bitta funksiyadan qaytaring (boʻsh slice'da xato).
- 🟡 **Oʻrta:** istalgan sonli soʻzni ajratgich bilan bitta qatorga qoʻshadigan variadic funksiya yozing.
- 🟠 **Qiyin:** `logla(nom, next)`: boshqa funksiyani oʻrab, har chaqiruvni sanab yozadigan funksiya qaytaring (middleware).
- 🔴 **Guru:** [`homework/guru/xato/`](homework/guru/xato/main.go): Temporal xatosining soddalashtirilgan shakli. `saqla()` xato qaytarishi kerak, lekin `nil` qaytaradi. Tuzating, `shadow` ham jim boʻlsin.

Avval oʻzingiz urinib koʻring. Javoblar: [homework/](homework/).
