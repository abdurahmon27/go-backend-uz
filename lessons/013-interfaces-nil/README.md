# 13 · Interface'lar: `nil` nega `nil`ga teng emas?

Oldingi dars: GO-1.12 · Metodlar: value receiver yoki pointer receiver? · Keyingi dars: GO-1.14 · Xato ham qiymat:
`%w`, `errors.Is`, `AsType` va `Join`

| Papka | Nima |
|---|---|
| `hook/` | `tekshir` nil `*SummaXato` ni `error` sifatida qaytaradi: **ataylab xato** (typed nil), `err != nil` doim true |
| `fix/` | tuzatilgani: xato boʻlmasa `return nil` |
| `stringer/` | yashirin bajarilish: `Buyurtma` va `Mijoz` `fmt.Stringer` ni hech qayerda eʼlon qilmasdan bajaradi |
| `ptrrecv/` | `String` pointer receiver bilan: `fmt.Println(b)` uni koʻrmaydi, `fmt.Println(&b)` koʻradi |
| `xato/` | **ataylab kompilyatsiya boʻlmaydi**: `var s fmt.Stringer = b` (method has pointer receiver) |
| `juftlik/` | interface ichidagi (tur, qiymat) juftligi `%T` / `%v` bilan; `unsafe.Sizeof(any)` |
| `assert/` | type assertion: comma-ok va panic |
| `tswitch/` | type switch |
| `savol/` | pauzadagi savol: `any(0) == any(int64(0))` |
| `homework/` | uy vazifasi javoblari |

`xato/` ataylab kompilyatsiya boʻlmagani uchun `go vet ./...` shu papkada xato beradi; qolganlarini tekshirish:
`go vet ./hook ./fix ./stringer ./ptrrecv ./juftlik ./assert ./tswitch ./savol ./homework/...`

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/013-interfaces-nil`) ishga tushiriladi.
```bash
sed -n 17,30p hook/main.go
go run ./hook | head -2                    # err: <nil>  /  xato bor!
go run ./stringer
sed -n 13,20p ptrrecv/main.go
go run ./ptrrecv                           # {150000 false}  /  buyurtma: 150000 soʻm, ...
go vet ./ptrrecv && echo vet: jim
go build -o /dev/null ./xato 2>&1 | grep -o 'Buyurtma does not.*'
go run ./juftlik
go run ./assert 2>&1 | head -3
go run ./tswitch
sed -n 658,676p $(go env GOROOT)/src/fmt/print.go   # fmt ichidagi type switch
go run ./savol                             # a == b: false
sed -n 17,30p hook/main.go
go run ./hook                              # ichida: *main.SummaXato, <nil>
sed -n 15p hook/main.go
sed -n 588,596p $(go env GOROOT)/src/fmt/print.go   # catchPanic: nil pointer → "<nil>"
sed -n 16,21p fix/main.go
go run ./fix
go vet ./hook && echo vet: jim
~/go/bin/staticcheck ./hook | head -2      # SA4023: this comparison is always true
```
Natijalar `go1.26.8 linux/amd64` va staticcheck `2026.2.1` da. Staticcheck oʻrnatish:
`go install honnef.co/go/tools/cmd/staticcheck@latest` (`~/go/bin` ga tushadi).

## Dars matni

**Hook.** `tekshir` toʻgʻri buyurtmada nil qaytaradi, ekranga `err: <nil>` chiqadi, lekin `if err != nil` baribir
ishlaydi. Bu Go FAQ'dagi [*Why is my nil error value not equal to nil?*](https://go.dev/doc/faq#nil_error) misolining
aynan oʻzi.

**Yashirin bajarilish.** Interface — metodlar roʻyxati. `fmt.Stringer`: `String() string`. Effective Go: *"Interfaces with
only one or two methods are common in Go code"*. Go'da `implements` yoʻq:
metodlari mos kelgan har qanday tur interface'ni oʻzi bajaradi (spec: *"A type T implements an interface I if T is not
an interface and is an element of the type set of I"*). Go FAQ: *"This property allows interfaces to be defined and
used without needing to modify existing code."* Kerak boʻlsa kompilyatsiya vaqtidagi tekshiruv:
`var _ fmt.Stringer = Buyurtma{}` ([Effective Go: Interface checks](https://go.dev/doc/effective_go#blank_implements)).

**Method set (GO-1.12 bilan bogʻliq).** `T` ning method set'ida faqat value receiver'li metodlar, `*T` nikida
ikkalasi. `String` pointer receiver bilan yozilsa, `Buyurtma` qiymati `Stringer` emas: `fmt.Println(b)` uni oddiy
struct sifatida chiqaradi (`{150000 false}`), **na kompilyator, na `go vet` ogohlantiradi**. `var s fmt.Stringer = b`
esa compile error: `Buyurtma does not implement fmt.Stringer (method String has pointer receiver)`.

**Ichida nima bor.** Interface qiymati ikki mashina soʻzi (64-bitli tizimda 2 × 8 = 16 bayt): tur va qiymat.
[`runtime/runtime2.go`](https://github.com/golang/go/blob/go1.26.8/src/runtime/runtime2.go#L184-L192):
`iface{tab *itab; data unsafe.Pointer}` (metodli interface; `itab` da tur va metodlar jadvali) va
`eface{_type *_type; data unsafe.Pointer}` (`any`). Interface faqat ikkalasi ham boʻsh boʻlganda `nil`. Go FAQ
buni (T, V) juftligi deydi: nil `*int` saqlangan interface `(T=*int, V=nil)` — `nil` emas.

**Type assertion.** `v, ok := x.(T)`: mos kelmasa `ok == false`, `v` — `T` ning zero value'si, panic yoʻq. `x.(T)`
ok'siz va mos kelmasa — panic: `interface conversion: interface {} is main.Buyurtma, not int`
([spec](https://go.dev/ref/spec#Type_assertions)).

**Type switch.** `switch v := x.(type)` turlarni solishtiradi. Bitta turli case'da `v` oʻsha turda, bir nechta turli
case'da `v` — `x` ning interface turida. `case nil` — interface boʻsh boʻlsa. `fmt` ham shunday qiladi:
[`print.go` `handleMethods`](https://github.com/golang/go/blob/go1.26.8/src/fmt/print.go#L622-L680) avval `Formatter` ni,
`%#v` uchun `GoStringer` ni, `%v`/`%s`/`%x`/`%X`/`%q` uchun esa type switch bilan avval `error`, keyin `Stringer` ni tekshiradi.

**Tenglik.** Ikki interface teng, agar dinamik turlari bir xil va qiymatlari teng boʻlsa yoki ikkalasi ham nil
boʻlsa ([spec: Comparison operators](https://go.dev/ref/spec#Comparison_operators)). Shuning uchun
`any(0) == any(int64(0))` — `false`. Turlar bir xil, lekin solishtirib boʻlmaydigan (slice, map, funksiya) boʻlsa —
runtime panic (uy vazifasi, qiyin).

**Typed nil.** `var xato *SummaXato` → `return xato`: `error` ichida `(T=*SummaXato, V=nil)`, shuning uchun
`err != nil` doim true.

**Nega ekranga `<nil>` chiqdi?** `fmt.Println` `Error()` ni chaqiradi, u nil pointer orqali `e.Summa` ni oʻqiydi va
panic qiladi. `fmt` panic'ni ushlaydi va argument nil pointer boʻlgani uchun `<nil>` yozadi
([`catchPanic`](https://github.com/golang/go/blob/go1.26.8/src/fmt/print.go#L588-L596): *"If it's a nil pointer, just
say "<nil>""*). Agar `Error()` receiver'ga tegmasa, typed nil uchun ham oʻz matnini chiqaradi — ya'ni `<nil>`
koʻrinishi tasodif. Haqiqatni `%T` aytadi: `*main.SummaXato`.

**Tuzatish** (Go FAQ): xato boʻlmasa aniq `return nil`; funksiya imzosida aniq tur (`*SummaXato`) emas, `error`.

**Kim topadi?** Kompilyator va `go vet` jim. [staticcheck SA4023](https://staticcheck.dev/docs/checks/#SA4023):
`this comparison is always true` … `tekshir never returns a nil interface value`.

**Qisqa.** Kichik interface (Rob Pike, [Go Proverbs](https://go-proverbs.github.io/): *"The bigger the interface,
the weaker the abstraction"*). Bu tuzoqni Go 2 da oʻzgartirish takliflari boʻlgan (Dave Cheney,
[Typed nils in Go 2](https://dave.cheney.net/2017/08/09/typed-nils-in-go-2); [golang/go#22729](https://github.com/golang/go/issues/22729)),
lekin hech biri qabul qilinmagan: #22729 2026-11-25 holatida hali ochiq (`LanguageChangeReview`), asosiy eʼtiroz — orqaga moslik.

## Manbalar
- [Go spec](https://go.dev/ref/spec) (go1.26): Interface types, Method sets, Type assertions, Type switches, Comparison operators
- [Go FAQ: nil error](https://go.dev/doc/faq#nil_error), [implements](https://go.dev/doc/faq#implements_interface),
  [guarantee satisfies interface](https://go.dev/doc/faq#guarantee_satisfies_interface)
- [Effective Go: Interface checks](https://go.dev/doc/effective_go#blank_implements)
- [A Tour of Go: Interface values with nil underlying values](https://go.dev/tour/methods/12)
- Russ Cox, [Go Data Structures: Interfaces](https://research.swtch.com/interfaces) (2009)
- Francesc Campoy, [Understanding nil](https://speakerdeck.com/campoy/understanding-nil) (GopherCon 2016)
- Go manba kodi (go1.26.8): [runtime/runtime2.go](https://github.com/golang/go/blob/go1.26.8/src/runtime/runtime2.go#L184-L192),
  [fmt/print.go](https://github.com/golang/go/blob/go1.26.8/src/fmt/print.go#L588-L676)
- [staticcheck SA4023](https://staticcheck.dev/docs/checks/#SA4023)
- [Go Proverbs](https://go-proverbs.github.io/)

## Uy vazifasi
- 🟢 **Oson:** `Mahsulot` turiga `String` metodini yozing va kompilyatsiya vaqtidagi tekshiruv qoʻshing.
- 🟡 **Oʻrta:** turda ham `Error`, ham `String` metodi bor. `fmt.Println` qaysi birini chaqiradi? Avval taxmin qiling.
- 🟠 **Qiyin:** ikki `any` ni panic'siz solishtiradigan `teng(a, b any) bool` yozing (`[]int{1}` ichida boʻlsa
  `a == b` panic beradi; `struct{ X any }` ichida boʻlsa ham).
- 🔴 **Guru:** `bosh(a any) bool`: `a == nil` yoki ichida nil pointer/map/slice/chan/func (typed nil) boʻlsa true.

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): `var _ fmt.Stringer = Mahsulot{}`; `String` ni oʻchirsangiz, kompilyatsiya boʻlmaydi.
- Oʻrta ([homework/orta](homework/orta/)): **`Error`**. `fmt` ning type switch'ida `case error` `case Stringer` dan oldin turadi.
- Qiyin ([homework/qiyin](homework/qiyin/)): `defer` + `recover` bilan panic'ni ushlab, `false` qaytarish.
  `reflect.TypeOf(a).Comparable()` yetarli emas: `struct{ X any }` "comparable", lekin `X` ichida slice boʻlsa
  baribir panic (jadvaldagi `struct ichida slice` holati). `TestPanic` panic'ni koʻrsatadi: `runtime error: comparing
  uncomparable type []int` (`go test -v ./homework/qiyin`). `recover` haqida batafsil: GO-1.15.
- Guru ([homework/guru](homework/guru/)): `reflect.ValueOf(a)` ning `Kind()` i pointer/map/slice/chan/func
  boʻlsa `IsNil()` (`reflect.ValueOf` hech qachon `Interface` kind qaytarmaydi: u ichidagi qiymatni oladi). Eslatma: bu diagnostika uchun; oddiy kodda typed nil'ni yaratmaslik (`return nil`) toʻgʻri yoʻl.
