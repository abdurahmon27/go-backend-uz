# 17 · Generics va generic metodlar: qachon kerak emas

Oldingi dars: GO-1.16 · `go.mod`dagi bitta qator dasturingiz xulqini oʻzgartiradi · Keyingi dars: GO-1.18 · Iterator'lar:
`range` endi funksiya ustida ham ishlaydi. Asos: GO-1.12 (metodlar, receiver) va GO-1.13 (interface'lar).

> **Go versiyasi.** Generic metodlar **Go 1.27** dan bor. Bu papka alohida modul: `go.mod` da `go 1.27` va
> `toolchain go1.27.4`. Kompyuteringizda eskiroq Go boʻlsa ham, `GOTOOLCHAIN=auto` (standart) go1.27.4 ni oʻzi yuklab
> oladi. Asosiy `go-backend-uz` moduli `go 1.26` da qoladi. `eski/` — yana bitta alohida modul (`go 1.26`), xuddi shu
> kod u yerda kompilyatsiya boʻlmasligini koʻrsatish uchun.

| Papka | Nima |
|---|---|
| `hook/` | `Royxat[T]` ning `Aylantir[U any]` generic metodi — Go 1.27 da ishlaydi |
| `eski/` | `hook/main.go` ning aynan nusxasi, `go 1.26` modulida: **kompilyatsiya boʻlmaydi** |
| `katta/` | ikki marta yozilgan bir xil kod → `Katta[T cmp.Ordered]` |
| `cheklov/` | oʻz constraint'imiz `Son interface{ ~int64 \| ~float64 }` va `type Ball int64` |
| `tildesiz/` | **ataylab kompilyatsiya boʻlmaydi**: `~` siz `Ball` endi `Son` emas |
| `royxat/` | 1.18–1.26 usuli: `Filtr` metod, `Aylantir` esa paket darajasidagi funksiya |
| `zanjir/` | 1.27 usuli: `r.Filtr(katta).Aylantir(somda)` |
| `randn/` | standart kutubxonadagi generic metod: `(*rand.Rand).N` |
| `imetod/` | **ataylab kompilyatsiya boʻlmaydi**: interface metodi type parameter eʼlon qila olmaydi |
| `ibajar/` | **ataylab kompilyatsiya boʻlmaydi**: generic metod interface metodini bajarmaydi |
| `faq/` | Go FAQ'dagi `Empty.Nop` misoli: type assertion `false`, `reflect` generic metodni koʻrmaydi |
| `savol/` | pauzadagi savol: `Katta(3, 2.5)` |
| `hid/` → `toza/` | ortiqcha generic'ning uchta hidi va tozalangani |
| `tezlik/` | interface, generic va oddiy tur: benchmark + `olch.sh` (7 marta, mediana) |
| `homework/` | uy vazifasi javoblari |

Tekshirish (shu papkadan):
```bash
go vet ./hook ./katta ./cheklov ./royxat ./zanjir ./randn ./faq ./savol ./hid ./toza ./tezlik ./homework/...
go test ./tezlik ./homework/...
```
`tildesiz/`, `imetod/`, `ibajar/` ataylab kompilyatsiya boʻlmagani uchun `go vet ./...` xato beradi.
`gofmt` ham 1.27 niki boʻlishi kerak: `$(go env GOROOT)/bin/gofmt -l $(ls -d */ | grep -v imetod)` (1.26 dagi gofmt
generic metodni tanimaydi; `imetod/` ataylab parse boʻlmaydi).

## Videodagi buyruqlar
Hammasi `lessons/017-generics` dan.
```bash
sed -n '7p;11,17p' hook/main.go
cd eski && GOTOOLCHAIN=local go version && GOTOOLCHAIN=local go build -o /dev/null . 2>&1 | grep -v '^#'
#   go version go1.26.8 linux/amd64
#   ./main.go:11:28: syntax error: method must have no type parameters
go version && go run ./hook            # go1.27.4 · [150000 soʻm 99000 soʻm 42000 soʻm]
go run ./katta                         # 7 1.5 olma / 150000 / 7 olma
go doc builtin.max | sed -n 3p         # func max[T cmp.Ordered](x T, y ...T) T
go run ./cheklov                       # 245 / 3.75
go build -o /dev/null ./tildesiz 2>&1 | grep -v '^#'
#   tildesiz/main.go:22:18: Ball does not satisfy Son (possibly missing ~ for int64 in Son)
sed -n 18,23p $(go env GOROOT)/src/cmp/cmp.go
grep '^go ' eski/go.mod                # go 1.26
cd eski && GOTOOLCHAIN=go1.27.4 go build -o /dev/null . 2>&1 | grep -v '^#'
#   ./main.go:11:29: generic method requires go1.27 or later (-lang was set to go1.26; check go.mod)
sed -n 210,213p $(go env GOROOT)/src/math/rand/v2/rand.go
cd $(go env GOROOT)/src && grep -rEl --include='*.go' '^func \([^)]*\) [A-Za-z_]+\[' . | grep -v -e _test.go -e /testdata/
#   ./math/rand/v2/rand.go   (go1.27.4: testlardan tashqari yagona generic metod)
go run ./randn                         # 7 / 92465 / 2.353284004s / func(uint8) uint8
sed -n 5,7p imetod/main.go
go build -o /dev/null ./imetod 2>&1 | grep 'type parameters'
#   imetod/main.go:6:10: interface method must have no type parameters
sed -n 7,13p ibajar/main.go
go build -o /dev/null ./ibajar 2>&1 | grep -v '^#'
#   Bosh does not implement MatnNop (wrong type for method Nop)
#       have Nop[T any](T) T
#       want Nop(string) string
go run ./faq                           # false / false / true / reflect koʻrgan metodlar: 1
go run ./savol                         # 3 float64
diff <(go run ./hid) <(go run ./toza) && echo natija bir xil
./tezlik/olch.sh                       # videoda saqlangan natija: 1593 / 1558 / 287 ns/op
go build -gcflags=-m ./tezlik 2>&1 | grep 'inlining call'
#   tezlik/tezlik.go:33:14: inlining call to (*Buyurtma).Jami     (faqat JamiOddiy ichida)
go test -c -o /tmp/tezlik.test ./tezlik && go tool nm /tmp/tezlik.test | grep JamiGeneric
#   R .dict.JamiGeneric[*Buyurtma]
#   T JamiGeneric[go.shape.*uint8]
```
Natijalar: `go1.27.4 linux/amd64` (eski uchun `go1.26.8`), i5-12450H, 2026-11-26. Benchmark (mediana, 7 marta,
`-cpu 1`, 1000 ta element, hammasi `0 allocs/op`): yengil yuklamada Interface 1593–1644, Generic 1558–1560, Oddiy
287–302 ns/op (≈ 5–6×); mashina band boʻlganda (load ≈ 20) Interface 2218–2509, Generic 2050–2348, Oddiy 486–604 ns/op
(≈ 4×). Interface va generic bir-biriga yaqin (goh biri, goh ikkinchisi oldinda), oddiy tur doim bir necha baravar tez.

## Dars matni

**Hook.** [Go FAQ](https://go.dev/doc/faq#generic_methods) da 2026-11-26 holatida hali ham shunday yozilgan:
*"We do not anticipate that Go will ever add generic methods."* Lekin [Go 1.27](https://go.dev/doc/go1.27#language):
*"Go 1.27 now supports generic methods: a method declaration may declare its own type parameters."* Bir xil fayl
go1.26.8 da `syntax error: method must have no type parameters`, go1.27.4 da ishlaydi.

**Type parameter va constraint.** Ian Lance Taylor, [*When To Use Generics*](https://go.dev/blog/when-generics) (2022):
*"avoid type parameters until you notice that you are about to write the exact same code multiple times."*
`func Katta[T cmp.Ordered](a, b T) T`: `T` — type parameter (kvadrat qavsda), `cmp.Ordered` — constraint. Turni
kompilyator argumentlardan topadi (type inference) yoki `Katta[int64](…)` deb oʻzingiz berasiz; `T` oʻrniga aniq tur
qoʻyiladi (instantiation, [spec](https://go.dev/ref/spec#Instantiations)). Eslatma: Go 1.21 dan beri tilda tayyor
`max`/`min` bor ([Go 1.21](https://go.dev/doc/go1.21#language)) — `Katta` faqat misol.

**Constraint — interface.** Unda metodlardan tashqari turlar roʻyxati (union) boʻlishi mumkin.
[Spec](https://go.dev/ref/spec#General_interfaces): *"The type set of a term of the form ~T is the set of all types
whose underlying type is T."* `type Ball int64` → `~int64` ga kiradi, `int64` ga emas. Kompilyator `~` yoʻqligini oʻzi
aytadi: `possibly missing ~ for int64 in Son`. `cmp.Ordered` ning har bir termida `~` bor.

**Generic metodlar (Go 1.27).** 1.18 dan beri generic turning metodlari bor, lekin ular faqat receiver'dagi
parametrlarni ishlata olardi. Yangi tur kerak boʻlsa (`Aylantir` → `Royxat[U]`), FAQ maslahati: *"use top-level
functions with type parameters"*. 1.27 da: `func (r Royxat[T]) Aylantir[U any](f func(T) U) Royxat[U]` — spec'dagi
`List.Apply` misoli bilan bir xil shakl ([Method declarations](https://go.dev/ref/spec#Method_declarations): *"If the
method declaration specifies type parameters (possibly in addition to type parameters declared by the receiver
specification), the method name denotes a generic method [Go 1.27]."*). Sabablar — [taklif #77273](https://github.com/golang/go/issues/77273)
(R. Griesemer): metod — turning nomlar maydonidagi funksiya; `x.a().b().c()` chapdan oʻngga oʻqiladi.
Standart kutubxonada: `func (r *Rand) N[Int intType](n Int) Int` (`math/rand/v2`). go1.27.4 manbasida (testlar va
`testdata` dan tashqari) bu yagona generic metod:
`cd $(go env GOROOT)/src && grep -rEl --include='*.go' '^func \([^)]*\) [A-Za-z_]+\[' . | grep -v -e _test.go -e /testdata/`
→ `./math/rand/v2/rand.go`.

**go.mod sharti.** Yangi toolchain yetmaydi, **til versiyasi** ham 1.27 boʻlishi kerak: `eski/` ni go1.27.4 bilan yigʻsak ham,
`go.mod` da `go 1.26` turgani uchun: `generic method requires go1.27 or later (-lang was set to go1.26; check go.mod)`.
Kutubxonada generic metod yozsangiz, uning `go.mod` i `go 1.27` boʻladi va foydalanuvchilar ham Go 1.27 ga oʻtishi kerak
(GO-1.16: `go` qatori; https://go.dev/doc/toolchain).

**Cheklov.** [Go 1.27 release notes](https://go.dev/doc/go1.27#language): *"methods of interfaces may not declare type
parameters nor can interface methods be implemented by generic methods."* Spec (Basic interfaces): *"Interface methods
cannot declare type parameters, but they may use type parameters from the interface declaration."* Generic metodni
type argument'siz qiymat sifatida olib boʻlmaydi: `f := q.Juft` → `cannot use generic function q.Juft without
instantiation`; `q.Juft[string]` yoki `var f func(string) (int, string) = q.Juft` ishlaydi.

**Nega shunday?** Go FAQ generic metodlar uchun toʻrt yoʻlni sanaydi: link'dan keyin qayta kompilyatsiya, JIT, sekin
zaxira kod, yoki generic metodlar interface'ni umuman bajarmasligi — oxirgisi haqida: *"Disallowing generic methods from
satisfying interfaces is unacceptable from a design point of view."* Taklif #77273 (*"A change of view"*): *"concrete
methods are a language feature that is useful in itself, irrespective of interfaces."* Go 1.27 aynan toʻrtinchi yoʻlni
tanladi. FAQ'dagi `Empty.Nop` misoli 1.27 da: `x.(interface{ Nop(string) string })` → `false`;
`reflect.TypeOf(x).NumMethod()` generic metodni sanamaydi (taklif: *"Generic methods also won't be accessible via
reflection"*).

**Qachon kerak emas** (Taylor, *When To Use Generics*):
Quyidagilar — *maslahat*, qatʼiy qoida emas (Taylor: *"consider"*). "Type parametrni olib tashlaganda hech narsa
yoʻqolmasa, u ortiqcha edi" — bu darsning oʻz amaliy tekshiruvi.
1. *"If all you need to do with a value of some type is call a method on that value, use an interface type, not a type
   parameter."* → `Chiqar[T fmt.Stringer]` oʻrniga `Chiqar(x fmt.Stringer)`.
2. `Yoz[T any](x T)` — `fmt.Println` baribir `any` oladi; `T` hech narsa bermaydi.
3. *"if the implementation is different for each type, then use an interface type and write different method
   implementations, don't use a type parameter."* → ichida `switch any(x).(type)` boʻlgan generic.
4. *"don't change from interface types to type parameters just for speed, because it probably won't run any faster."*

**Tezlik.** Go FAQ ([generics_implementation](https://go.dev/doc/faq#generics_implementation)): *"The standard Go
compiler ordinarily emits a single instantiation for every type argument with the same shape"*.
[Dictionaries design (go1.18)](https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md):
barcha pointer turlar bitta `go.shape.*uint8` shaklida; type parameter metodini chaqirish *"implemented as a conversion of
the receiver to the type bound interface"*. `go tool nm` buni koʻrsatadi: `JamiGeneric[go.shape.*uint8]` (tana) va
`.dict.JamiGeneric[*Buyurtma]` (lugʻat). Asosiy farq esa inline'da: `-gcflags=-m` boʻyicha `Jami` faqat `JamiOddiy`
ichida inline qilinadi (chaqiruv yoʻqoladi); generic va interface versiyalarida u bilvosita chaqiruv boʻlib qoladi.
Bular **standart (gc) kompilyator** tafsilotlari — FAQ: *"Different compilers will make different choices"*.

## Manbalar
- Go 1.27 release notes — Changes to the language: https://go.dev/doc/go1.27#language
- Go spec (Language version go1.27, May 26, 2026): Method declarations, Basic/General interfaces, Instantiations,
  Type inference: https://go.dev/ref/spec
- golang/go#77273 *spec: generic methods for Go* (R. Griesemer, 2026-01-22, Proposal-Accepted, Go1.27):
  https://github.com/golang/go/issues/77273
- Go FAQ: generic_methods, generics_implementation: https://go.dev/doc/faq#generic_methods
- I. L. Taylor, *When To Use Generics* (2022): https://go.dev/blog/when-generics
- R. Griesemer, I. L. Taylor, *An Introduction To Generics* (2022): https://go.dev/blog/intro-generics
- Type Parameters Proposal, *No parameterized methods*: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
- Generics implementation — dictionaries (go1.18): https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md
- Manba kodi (go1.27.4): `src/cmp/cmp.go` L18–23, `src/math/rand/v2/rand.go` L210–213

## Uy vazifasi
- 🟢 **Oson:** `Indeks[T comparable](s []T, x T) int` yozing (topilmasa `-1`) va `slices.Index` bilan solishtiring.
- 🟡 **Oʻrta:** `type Quti[T any] struct{ v T }`. Qaysi qatorlar Go 1.27 da kompilyatsiya boʻladi?
  1. `func (q Quti[T]) Ol() T { return q.v }`
  2. `func (q Quti[T]) Juft[U any](u U) (T, U) { return q.v, u }`
  3. `type I interface{ M[T any](T) }`
  4. `f := q.Juft`
  5. `var o Oluvchi[int] = q` (bunda `type Oluvchi[T any] interface{ Ol() T }`)
  6. `f := q.Juft[string]`
- 🟠 **Qiyin:** `Royxat[T]` ga `Guruhla[K comparable](kalit func(T) K) map[K]Royxat[T]` generic metodini yozing va
  test qiling (bir marta `bool`, bir marta `string` kalit bilan).
- 🔴 **Guru:** `JamiGeneric[T Summali]` ni `*Buyurtma`, `*Mijoz`, `Narx` (`int64` asosli) va `Chegirma` (`int64` asosli)
  bilan chaqiring. Kompilyator nechta tana va nechta lugʻat yasaydi? `go tool nm` bilan tekshiring.

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): oddiy `for` + `==` (`comparable` kerak, chunki `==` ishlatiladi).
  `slices.Index` (Go 1.21+) aynan shuni qiladi — real kodda tayyorini ishlating.
- Oʻrta ([homework/orta](homework/orta/)): **1, 2, 5, 6** kompilyatsiya boʻladi. 3 — `interface method must have no type
  parameters`; 4 — `cannot use generic function q.Juft without instantiation`. 5: generic *interface* (`Oluvchi[T]`)
  boshqa narsa — uning metodi oddiy, faqat interface'ning `T` sini ishlatadi, shuning uchun `Quti[int]` uni bajaradi.
- Qiyin ([homework/qiyin](homework/qiyin/)): `K` — metodning oʻz type parametri; `go test ./homework/qiyin`.
- Guru ([homework/guru](homework/guru/)): **2 ta tana, 4 ta lugʻat** (standart gc kompilyatori, go1.27.4).
  ```
  go test -c -o /tmp/guru.test ./homework/guru && go tool nm /tmp/guru.test | grep JamiGeneric
  R ..dict.JamiGeneric[*…Buyurtma]   R ..dict.JamiGeneric[*…Mijoz]
  R ..dict.JamiGeneric[…Chegirma]    R ..dict.JamiGeneric[…Narx]
  T JamiGeneric[go.shape.*uint8]     T JamiGeneric[go.shape.int64]
  ```
  Ikkala pointer tur bitta `go.shape.*uint8` tanasini, `Narx` va `Chegirma` esa bitta `go.shape.int64` tanasini
  boʻlishadi; turlar orasidagi farq (qaysi `Jami`) lugʻatda.
