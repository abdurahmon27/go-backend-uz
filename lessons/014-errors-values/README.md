# 14 · Xato ham qiymat: `%w`, `errors.Is`, `AsType` va `Join`

Oldingi dars: GO-1.13 · Interface'lar: `nil` nega `nil`ga teng emas? · Keyingi dars: GO-1.15 · `defer`, `panic`,
`recover`: nimani qutqarib boʻlmaydi

| Papka | Nima |
|---|---|
| `wrap/` | `%v` va `%w` farqi: matn bir xil, lekin `errors.Is` faqat `%w` da topadi |
| `zanjir/` | uch qavat oʻralgan xato, `errors.Unwrap` bilan qavatma-qavat ochish |
| `astype/` | `errors.As` (Go 1.13) va `errors.AsType` (Go 1.26) yonma-yon |
| `asxato/` | **ataylab xato**: `errors.As(err, sx)` (`&` yoʻq) — kompilyatsiya boʻladi, ishlaganda panic, `go vet` topadi |
| `astypexato/` | **ataylab kompilyatsiya boʻlmaydi**: `errors.AsType[SummaXato]` — `SummaXato` error emas (pointer receiver) |
| `tezlik/` | `errors.As` va `errors.AsType` benchmark'i |
| `join/` | `errors.Join` bilan bir nechta xatoni yigʻish |
| `savol/` | pauzadagi savol: `errors.Join(nil, nil) == nil` |
| `tartib/` | **ataylab xato**: natija `err` tekshirilmasdan oldin ishlatiladi; `versiyalar.sh` uni olti Go versiyasida ishga tushiradi. Alohida modul (`go 1.20`), eski toolchain'lar ham ishlata olishi uchun |
| `homework/` | uy vazifasi javoblari |

`astypexato/` ataylab kompilyatsiya boʻlmaydi, `asxato/` esa `go vet` dan ataylab oʻtmaydi. Qolganlarini tekshirish:
`go vet ./wrap ./zanjir ./astype ./tezlik ./join ./savol ./homework/...` va `(cd tartib && go vet .)`.

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/014-errors-values`) ishga tushiriladi.
```bash
sed -n 27,34p tartib/main.go
cd tartib && GOTOOLCHAIN=go1.24.13 go run -trimpath .    # xato: buyurtma topilmadi
cd tartib && GOTOOLCHAIN=go1.26.8 go run -trimpath . 2>&1 | head -1   # panic: ... nil pointer dereference
sed -n 315,319p $(go env GOROOT)/src/builtin/builtin.go  # type error interface { Error() string }
sed -n 10,17p wrap/main.go
go run ./wrap
go run ./zanjir
sed -n 54,75p $(go env GOROOT)/src/errors/wrap.go        # errors.Is ichida: Is metodi, Unwrap bilan pastga
go run ./astype
sed -n 18p asxato/main.go                                # & yoʻq
go run ./asxato 2>&1 | head -1                           # panic: errors: target must be a non-nil pointer
go vet ./asxato 2>&1 | grep -o 'second argument.*'
sed -n 17p astypexato/main.go                            # * yoʻq
go build -o /dev/null ./astypexato 2>&1 | grep -o 'SummaXato does not.*'
go test -bench . -benchmem ./tezlik | grep ns/op
go test -gcflags=-m -run XXX ./tezlik 2>&1 | grep moved  # moved to heap: sx
sed -n 21,30p join/main.go
go run ./join
go run ./savol                                           # err == nil: true
sed -n 20,29p $(go env GOROOT)/src/errors/join.go
sed -n 21,29p join/main.go
./tartib/versiyalar.sh
cd tartib && go vet . && echo vet: jim
cd tartib && ~/go/bin/staticcheck . && echo staticcheck: jim
```
Natijalar `go1.26.8 linux/amd64`, i5-12450H, staticcheck `2026.2.1`. `versiyalar.sh` `GOTOOLCHAIN` orqali
go1.20.14, go1.21.4, go1.21.5, go1.24.13, go1.25.0 ni `proxy.golang.org` dan yuklab oladi (birinchi marta, har biri
~70 MB). Mahalliy Go yangiroq boʻlsa (masalan go1.27), oxirgi qator oʻsha versiyani koʻrsatadi.

Benchmark (bitta ishga tushirish; 10 marta: `As` 88–115 ns, `AsType` 10.2–12.2 ns):
```
BenchmarkAs-12        	13093354	        96.46 ns/op	       8 B/op	       1 allocs/op
BenchmarkAsType-12    	107957600	        11.30 ns/op	       0 B/op	       0 allocs/op
```

## Dars matni

**Hook.** `tartib/` dagi kod Go 1.24.13 da `xato: buyurtma topilmadi` deb jim tugaydi, Go 1.26.8 da esa
`panic: runtime error: invalid memory address or nil pointer dereference`. Kod bir xil. Javob — darsning oxirida.

**Xato — qiymat.** Spec: `type error interface { Error() string }`, *"the nil value representing no error"*
([spec: Errors](https://go.dev/ref/spec#Errors)). Xato odatda oxirgi natija (staticcheck
[ST1008](https://staticcheck.dev/docs/checks/#ST1008)). `errors.New` matnli xato yaratadi; paket darajasidagi
`var ErrTopilmadi = errors.New(...)` — *sentinel* xato. Rob Pike: *"Errors are values. Values can be programmed"*
([Errors are values](https://go.dev/blog/errors-are-values), 2015). `error` interface ekanini va typed nil tuzogʻini
GO-1.13 da koʻrdik.

**`%v` va `%w`.** `fmt.Errorf("buyurtma %d: %v", id, err)` asl xatodan faqat matnni oladi (Go blog: *"discards
everything from the original error except the text"*). `%w` esa asl xatoni ichida saqlaydi: natijada `Unwrap() error`
metodi bor `*fmt.wrapError` boʻladi. Matn ikkalasida bir xil, lekin `errors.Is(err, ErrTopilmadi)` faqat `%w` da
`true`. Oʻralgan xatoda `==` har doim `false`: interface ichida endi boshqa tur (`*fmt.wrapError`) turibdi (GO-1.13). Go 1.13 dan beri
([release notes](https://go.dev/doc/go1.13#error_wrapping)).

**Zanjir.** `%w` bilan yozilgan har bir `fmt.Errorf` yangi qavat qoʻshadi: `toʻlov → buyurtma 42 → buyurtma topilmadi`. `errors.Unwrap` bitta
qavatni ochadi. `errors.Is` avval xatoning oʻzini solishtiradi (yoki uning `Is(error) bool` metodini chaqiradi),
keyin `Unwrap` bilan pastga tushadi ([`errors/wrap.go`](https://github.com/golang/go/blob/go1.26.8/src/errors/wrap.go#L54-L79)).

**Qachon oʻramaslik kerak.** Go blog: *"wrapping an error makes that error part of your API"*. Misol:
`errors.Is(err, sql.ErrNoRows)` ga mijozlar bogʻlanadi, bazani almashtirsangiz ularning kodi buziladi. Ichki
tafsilotni yashirish uchun `%v` — matn bir xil qoladi
([Working with Errors in Go 1.13: Whether to Wrap](https://go.dev/blog/go1.13-errors#whether-to-wrap)). Uber Go
style guide ham shunday: `%w` — chaqiruvchi asl xatoga muhtoj boʻlsa, `%v` — yashirish uchun
([Error Wrapping](https://github.com/uber-go/guide/blob/master/style.md#error-wrapping)).

**`errors.As` → `errors.AsType`.** Eski usul (Go 1.13): `var sx *SummaXato; if errors.As(err, &sx) {…}`. Yangi
(Go 1.26, [proposal #51945](https://github.com/golang/go/issues/51945)): `if sx, ok := errors.AsType[*SummaXato](err); ok {…}`
— type assertion `x.(T)` ga oʻxshaydi, `sx` faqat `if` ichida. Go 1.26 release notes: *"The new AsType function is a
generic version of As. It is type-safe, faster, and, in most cases, easier to use."*
([go1.26#errorspkgerrors](https://go.dev/doc/go1.26#errorspkgerrors)). `errors.As` doc'i ham endi *"For most uses, prefer AsType"* deydi.

**Type-safe nimani anglatadi.** `errors.As(err, sx)` (`&` unutilgan) kompilyatsiya boʻladi, ishlaganda
`panic: errors: target must be a non-nil pointer`. Uni kompilyator emas, `go vet` (`errorsas` analyzer) ushlaydi;
`go test` ham `errorsas` ni avtomatik ishga tushiradi (`go help test`: *"a high-confidence subset of the default go vet
checks … errorsas …"*).
`errors.AsType[SummaXato](err)` esa compile error: `SummaXato does not satisfy error (method Error has pointer
receiver)` — GO-1.12 dagi method set qoidasi. Eslatma: `AsType` faqat `error` ni bajaradigan turni rad etadi;
zanjirda yoʻq, lekin `error` boʻlgan tur (`*os.PathError`) baribir kompilyatsiya boʻladi va `ok == false` qaytaradi.
Narxi ham bor: `errors.As` istalgan interface'ni qabul qiladi (doc: *"doesn't require its target argument to implement
error"*), `AsType[interface{ Timeout() bool }]` esa compile error (`missing method Error`); buning uchun
`AsType[interface{ error; Timeout() bool }]` yozing.

**Tezroq.** Ikki qavat oʻralgan xatoda `As` ~90–100 ns va 1 allocation (8 B), `AsType` ~11 ns, 0 allocation
(i5-12450H, go1.26.8). Sabab: `As` `internal/reflectlite` bilan ishlaydi va `&sx` `any` ga uzatilgani uchun `sx` heap'ga
qochadi (`moved to heap: sx`, GO-1.11). `AsType` oddiy type assertion `err.(E)` qiladi.

**`errors.Join`.** Go 1.20: bir xatoda bir nechta xato (`Unwrap() []error`); `errors.Is`/`As` ularning har birini
tekshiradi; `fmt.Errorf` bir nechta `%w` qabul qiladi ([go1.20#errors](https://go.dev/doc/go1.20#errors)). `Join` nil
xatolarni tashlaydi; matnlar `\n` bilan qoʻshiladi.

**Savol.** `errors.Join(nil, nil) == nil` → `true`. `Join` nil boʻlmaganlarni sanaydi, `n == 0` boʻlsa aniq `return nil`
([`join.go`](https://github.com/golang/go/blob/go1.26.8/src/errors/join.go#L20-L29)) — GO-1.13 dagi qoida. Shuning uchun
`join/` da `summaXato` `error` turida eʼlon qilingan (guru vazifasiga qarang).

**Go 1.25: nil tekshiruvi.** `b, err := top(42); summa := b.Summa; if err != nil {…}` — `b` nil, `b.Summa` spec boʻyicha
panic berishi kerak. Oʻlchangan natija:

```
go1.20.14  panic: runtime error: invalid memory address or nil pointer dereference
go1.21.4   panic: runtime error: invalid memory address or nil pointer dereference
go1.21.5   xato: buyurtma topilmadi
go1.24.13  xato: buyurtma topilmadi
go1.25.0   panic: runtime error: invalid memory address or nil pointer dereference
go1.26.8   panic: runtime error: invalid memory address or nil pointer dereference
```

`b.Summa` ni oʻqishdan oldin kompilyator yashirin nil tekshiruvi qoʻyadi. Bug shu tekshiruvni `if err != nil` dan
keyinga surgan (dastur esa chiqib ketgani uchun tekshiruvgacha yetmagan); Go 1.25 tuzatdi
([release notes: nil pointer bug](https://go.dev/doc/go1.25#nil-pointer-bug), [golang/go#72860](https://github.com/golang/go/issues/72860),
CL [657715](https://go.dev/cl/657715) *"cmd/compile: don't move nilCheck operations during tighten"*).

**Aniqlik kiritish.** Release notes *"introduced in Go 1.21"* deydi. Oʻlchovimizda go1.21.0–go1.21.4 hali panic beradi
(release notes'dagi `os.Open` misoli ham), bug go1.21.5 dan boshlanadi. Sababi: bug'ni keltirib chiqargan CL
[537775](https://go.dev/cl/537775) (#63657 tuzatishi, Ian Lance Taylor #72860 da koʻrsatgan) 1.21 ga backport qilingan:
[#63743](https://github.com/golang/go/issues/63743), milestone Go1.21.5.

**Eski kodni sindirdimi?** Dastur oldin ham notoʻgʻri edi: bug uni yashirgan. Go jamoasi tuzatishni 1.23/1.24 tuzatish
versiyalariga qoʻshmadi: *"We've found instances of code that depend on the buggy behavior … We don't want a point
release to break anyone's (even incorrect) code"* (randall77, #72860, 2025-03-18). `go vet` ham, staticcheck ham bu
kodda jim (`staticcheck -checks all` ham). Uber'ning [NilAway](https://github.com/uber-go/nilaway) esa topadi:
`tartib/main.go:29:11: error: Potential nil panic detected … result 0 of `top()` lacking guarding; accessed field
`Summa`` (nilaway `v0.0.0-20261124161746-91ad9bf5d8e6`, videoda koʻrsatilmagan). Qoida: xatoni darhol, natijani
ishlatishdan oldin tekshiring.

## Manbalar
- [Go spec: Errors](https://go.dev/ref/spec#Errors) (go1.26)
- [pkg.go.dev/errors](https://pkg.go.dev/errors), [fmt.Errorf](https://pkg.go.dev/fmt#Errorf);
  manba kodi (go1.26.8): [errors/wrap.go](https://github.com/golang/go/blob/go1.26.8/src/errors/wrap.go),
  [errors/join.go](https://github.com/golang/go/blob/go1.26.8/src/errors/join.go)
- Release notes: [Go 1.13](https://go.dev/doc/go1.13#error_wrapping), [Go 1.20](https://go.dev/doc/go1.20#errors),
  [Go 1.25](https://go.dev/doc/go1.25#nil-pointer-bug), [Go 1.26](https://go.dev/doc/go1.26#errorspkgerrors)
- Go blog: [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors) (D. Neil, J. Amsterdam, 2019),
  [Errors are values](https://go.dev/blog/errors-are-values) (R. Pike, 2015)
- Issues: [#51945](https://github.com/golang/go/issues/51945) (AsType), [#72860](https://github.com/golang/go/issues/72860),
  [#63743](https://github.com/golang/go/issues/63743)
- [Uber Go Style Guide: Error Wrapping](https://github.com/uber-go/guide/blob/master/style.md#error-wrapping)
- [go vet errorsas](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/errorsas)

## Uy vazifasi
- 🟢 **Oson:** toʻlangan buyurtma uchun `ErrTolangan` sentinel xatosini yarating, uni `%w` bilan oʻrang va `errors.Is` bilan toping.
- 🟡 **Oʻrta:** `fmt.Errorf("saqlash: %w; %w", ErrDisk, ErrLog)`. `errors.Is` ikkalasini topadimi? `errors.Unwrap(err)` nima qaytaradi?
- 🟠 **Qiyin:** `MaydonXato{Maydon string; Err error}` turini yozing: `errors.Is` va `errors.AsType` uning ichiga ham qarasin. Test bilan tekshiring.
- 🔴 **Guru:** `join/` dagi `tekshir` da `summaXato` ni `*SummaXato` turida eʼlon qiling. Toʻgʻri buyurtma uchun nima chiqadi va nega?

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): `==` — `false`, `errors.Is` — `true`.
- Oʻrta ([homework/orta](homework/orta/)): `errors.Is` ikkalasini ham topadi. `errors.Unwrap` esa **`nil`**: bir nechta `%w`
  li xatoda `Unwrap() error` emas, `Unwrap() []error` bor, `errors.Unwrap` esa faqat birinchisini chaqiradi
  (doc: *"Unwrap only calls a method of the form "Unwrap() error""*).
- Qiyin ([homework/qiyin](homework/qiyin/)): `func (e *MaydonXato) Unwrap() error { return e.Err }`.
  `go test -v ./homework/qiyin` — PASS.
- Guru ([homework/guru](homework/guru/)): `err == nil: false`, `err: %!v(PANIC=Error method: runtime error: invalid memory
  address or nil pointer dereference)`. `Join` "nil" larni `err != nil` bilan tashlaydi, lekin nil `*SummaXato`
  `error` ichida typed nil — nil emas (GO-1.13). `Join` uni saqlaydi, `Error()` nil pointer'dan `e.Summa` ni oʻqib
  panic qiladi, `fmt` uni ushlab `%!v(PANIC=…)` deb yozadi (`joinError` pointer'i nil emas, shuning uchun `<nil>` emas).
