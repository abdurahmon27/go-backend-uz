# 12 · Metodlar: value receiver yoki pointer receiver?

Oldingi dars: [11 · Pointer'lar: bu oʻzgaruvchi nega heap'ga qochdi?](../011-pointers-escape/) · Keyingi dars: GO-1.13 · Interface'lar: `nil` nega `nil`ga teng emas?

| Papka | Nima |
|---|---|
| `hook/` | **ataylab kompilyatsiya boʻlmaydi**: `s[0].Tola()` ishlaydi, `m["a1"].Tola()` — `cannot call pointer method` |
| `metod/` | metod eʼloni, `Buyurtma.Jami(b)` (method expression), value receiver `Tola` (xato) va pointer receiver `TolaP` |
| `manzilsiz/` | **ataylab kompilyatsiya boʻlmaydi**: 5 ta addressable qiymat, `&T{}` istisnosi va 3 ta addressable boʻlmagan qiymatda pointer metod |
| `mapfix/` | map elementida pointer metod: nusxa → oʻzgartir → qaytarib yoz, yoki `map[string]*Buyurtma` |
| `stringer/` | `String()` pointer receiver'da: `fmt.Println(b)` uni chaqirmaydi, `fmt.Println(&b)` chaqiradi |
| `stringerxato/` | **ataylab kompilyatsiya boʻlmaydi**: `var s fmt.Stringer = b` — `method String has pointer receiver` |
| `embed/` | embedding'da metodlar koʻtariladi (GO-1.10 vaʼdasi) |
| `savol/` | pauzadagi savol: method value receiver'ni qachon oladi |
| `myuteks/` | `sync.Mutex` bor struct'da aralash receiver'lar — **ataylab deadlock**, `go vet` ushlaydi |
| `nilp/` | `nil` pointer orqali metod: pointer receiver `nil`ni tekshira oladi, value receiver panic beradi (**ataylab**) |
| `homework/` | uy vazifasi javoblari (`homework/orta` ham ataylab kompilyatsiya boʻlmaydi) |

Ataylab kompilyatsiya boʻlmaydigan papkalar sabab `go vet ./...` bu dars papkasida xato bilan tugaydi: har bir papkani
alohida ishga tushiring.

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/012-methods-receivers`) ishga tushiriladi.
```bash
sed -n '10p;12,18p' hook/main.go | grep -v '^$'
go build -o /dev/null ./hook 2>&1 | grep -v '^#'          # cannot call pointer method Tola on Buyurtma
go run ./metod                                             # Tola: false, TolaP: true
~/go/bin/staticcheck ./metod                               # SA4005 ineffective assignment
go build -o /dev/null ./manzilsiz 2>&1 | grep -v '^#'     # 31, 32, 33-qatorlar
go run ./mapfix
go run ./stringer                                          # {50000 false} va String() natijasi
go vet ./stringer && echo vet: jim
go build -o /dev/null ./stringerxato 2>&1 | grep -v '^#'  # method String has pointer receiver
go run ./embed
go run ./savol                                             # f(): 100, g(): 200
go run ./myuteks                                           # fatal error: all goroutines are asleep - deadlock!
go vet ./myuteks                                           # Balansi passes lock by value
go doc strings.Builder                                     # Do not copy a non-zero Builder.
go doc -all time | grep -E '^func \(t \*?Time\)'           # 43 ta value, 4 ta pointer receiver
go doc -all strings.Builder | grep -E '^func \(b'          # hammasi *Builder
```
Natijalar `go1.26.8 linux/amd64` va staticcheck 2026.2.1 bilan. Videoda `go vet ./myuteks` chiqishidagi uzun paket yoʻli
(`github.com/abdurahmon27/go-backend-uz/lessons/012-methods-receivers/myuteks.Hisob`) `sed` bilan `Hisob` ga qisqartirilgan.
staticcheck oʻrnatish: `go install honnef.co/go/tools/cmd/staticcheck@latest`.

## Dars matni

**Hook.** Bitta pointer metod (`func (b *Buyurtma) Tola()`) slice elementida ishlaydi, map elementida esa dastur
kompilyatsiya boʻlmaydi: `cannot call pointer method Tola on Buyurtma`. Sababi ikki qoida: method set va addressability.

**Metod va receiver.** Metod — receiver'i bor funksiya (spec: *Method declarations*). Receiver `func` va metod nomi
orasida yoziladi. Go FAQ: receiver *"behaves exactly as if it were an argument to the method"*. Buni method expression
koʻrsatadi: `Buyurtma.Jami(b)` — oddiy funksiya, birinchi argumenti receiver (spec: *Method expressions*).
GO-1.06 va GO-1.11 dagi qoida shu yerda ham: value receiver nusxa oladi. `Tola` nusxani oʻzgartiradi (`false`),
`TolaP` manzil oladi (`true`). Kompilyator ham, `go vet` ham jim; staticcheck: `SA4005 ineffective assignment`.

**Go oʻzi qoʻshadi.** Spec, *Calls*: *"If x is addressable and &x's method set contains m, x.m() is shorthand for
(&x).m()"*. Teskari yoʻnalish: pointer orqali value metod — `p.Jami()` = `(*p).Jami()` (spec: *Method values*). `p` nil
boʻlsa — panic (`nilp/`). Pointer receiver'li metod esa `nil` receiver bilan ham chaqiriladi va uni oʻzi tekshira oladi.

**Addressability.** Spec, *Address operators*: oʻzgaruvchi, pointer indirection (`*p`), slice elementi, addressable
struct maydoni, addressable massiv elementi. Composite literal addressable emas, lekin `&T{}` yozish mumkin — bu faqat `&`
uchun istisno (`(&T{}).M()` da receiver allaqachon pointer, addressability kerak ham emas). Addressable emas: map
elementi, funksiya natijasi, `T{}` ning oʻzi, konstanta. `manzilsiz/` da kompilyator faqat 31–33-qatorlarga xato beradi.

**Nega map elementi addressable emas.** Bir sababi — oʻsish: Go 1.24 dan beri map Swiss table (GO-1.09); jadval oʻsganda
runtime elementlarni yangi jadvalga koʻchiradi (`internal/runtime/maps/map.go`, *Growth*: ikki barobar katta jadvalga
almashtiradi yoki ikkiga boʻladi). Eski manzil eskirgan nusxaga qarab qolardi — GO-1.11 dagi `&s[0]` + `append`
xatosining aynan oʻzi. Slice'da bu xato runtime'da jim yuz beradi, map'da til uni taqiqlaydi. Yechim: nusxa → oʻzgartir
→ qaytarib yoz, yoki `map[string]*Buyurtma` (Go wiki *MethodSets*: struct elementli map'lar koʻpincha pointer saqlaydi).

**Method set va `fmt.Stringer`.** `Println` `...any` oladi: unga har qanday qiymat sigʻadi, `String` borligi dastur
ishlayotganda tekshiriladi.  Spec, *Method sets*: `T` ning method set'i — receiver'i `T` boʻlgan metodlar; `*T`
niki — `T` va `*T`. `String()` pointer receiver'da boʻlsa, `Buyurtma` `fmt.Stringer` emas: `fmt.Println(b)` maydonlarni
xom chiqaradi (`{50000 false}`), `fmt.Println(&b)` esa `String()` ni chaqiradi. `fmt` buni `handleMethods` ichidagi
`case Stringer:` type switch'i bilan tekshiradi. `Println` bilan `go vet` jim (`fmt.Printf("%s", b)` ni esa vet ushlaydi:
`format %s has arg b of wrong type`). `fmt.Stringer` turidagi oʻzgaruvchiga bersangiz, kompilyator
xato beradi: `Buyurtma does not implement fmt.Stringer (method String has pointer receiver)`. Go FAQ, *Why do T and *T
have different method sets?*: interface ichidagi qiymatning manzilini xavfsiz olib boʻlmaydi; olinsa ham, metod nusxani
oʻzgartirar va oʻzgarish yoʻqolardi (`io.Copy(buf, os.Stdin)` misoli). Interface'lar — GO-1.13 da.

**Embedding.** GO-1.10 da maydonlar koʻtarildi, metodlar ham koʻtariladi: `a.Kochir("Samarqand")` = `(&a.Manzil).Kochir(...)`,
`a.Toliq()` = `a.Manzil.Toliq()`. Receiver — ichki `Manzil`, `Mijoz` emas (Effective Go, *Embedding*). Bu meros emas.
Spec (*Struct types*): `S` ichida `T` boʻlsa, `S` va `*S` method set'iga `T` receiver'li metodlar kiradi, `*S` ga yana `*T`
receiver'lilari ham; `S` ichida `*T` boʻlsa, ikkalasiga ham hammasi (uy vazifasi guru).

**Method value (pauza).** `f := b.Jami` — receiver shu paytda hisoblanib saqlanadi (spec, *Method values*: *"evaluated and
saved"*). Value receiver: nusxa (`f() = 100`), pointer receiver: `&b` (`g() = 200`). `defer b.M()` ham xuddi shunday.

**Mutex.** `Hisob` da `sync.Mutex` bor. `Qosh` (pointer receiver) qulfni ushlab turib `Balansi` (value receiver) ni
chaqiradi: `Balansi` mutex'ning nusxasini oladi, nusxa ham qulflangan, `Lock` abadiy kutadi →
`fatal error: all goroutines are asleep - deadlock!`. `go vet` (copylocks): `Balansi passes lock by value` (videoda paket
yoʻli `sed` bilan qisqartirilgan). Agar `Hisob` ning hamma metodi value receiver boʻlsa, deadlock boʻlmaydi, lekin
yozuvlar `Tola` dagi kabi nusxaga ketadi va mutex hech narsani himoya qilmaydi. `go doc sync.Mutex`: *"A Mutex must
not be copied after first use."*

**Qaysi receiver.** Go FAQ (*Should I define methods on values or pointers?*) va Go wiki CodeReviewComments
(*Receiver Type*):
1. metod receiver'ni oʻzgartirsa → pointer;
2. struct'da `sync.Mutex` boʻlsa yoki struct katta boʻlsa → pointer;
3. kichik, oʻzgarmas qiymat (`time.Time` kabi) yoki oddiy tur → value: arzon va tushunarli;
4. bitta turda aralashtirmang (izchil method set);
5. ikkilansangiz → pointer.
Standart kutubxona: `time.Time` (24 bayt) ning 43 ta metodi value receiver, 4 tasi (`GobDecode`, `UnmarshalBinary`,
`UnmarshalJSON`, `UnmarshalText`) pointer receiver — ular vaqtning oʻzini oʻzgartiradi, asosli istisno. `strings.Builder` ning
hamma 9 ta metodi pointer receiver. Hujjat: *"Do not copy a non-zero Builder."* — yozuv metodlari `copyCheck` bilan
nusxani tekshiradi (`addr *Builder`). `Len` buni tekshirmaydi, u tur bo'yicha izchillik uchun pointer receiver.

## Manbalar
- Go spec (go1.26): [Method declarations](https://go.dev/ref/spec#Method_declarations), [Method sets](https://go.dev/ref/spec#Method_sets),
  [Calls](https://go.dev/ref/spec#Calls), [Address operators](https://go.dev/ref/spec#Address_operators),
  [Method expressions](https://go.dev/ref/spec#Method_expressions), [Method values](https://go.dev/ref/spec#Method_values),
  [Struct types](https://go.dev/ref/spec#Struct_types)
- [Go FAQ: Should I define methods on values or pointers?](https://go.dev/doc/faq#methods_on_values_or_pointers)
- [Go FAQ: Why do T and *T have different method sets?](https://go.dev/doc/faq#different_method_sets)
- [Go wiki: CodeReviewComments — Receiver Type](https://go.dev/wiki/CodeReviewComments#receiver-type)
- [Go wiki: MethodSets](https://go.dev/wiki/MethodSets)
- [Effective Go: Embedding](https://go.dev/doc/effective_go#embedding)
- Go manba kodi: [internal/runtime/maps/map.go — Growth](https://github.com/golang/go/blob/go1.26.8/src/internal/runtime/maps/map.go#L86-L102),
  [fmt/print.go — handleMethods](https://github.com/golang/go/blob/go1.26.8/src/fmt/print.go#L622)
- [staticcheck SA4005](https://staticcheck.dev/docs/checks/#SA4005), [go vet copylocks](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/copylock)
- A. Donovan, B. Kernighan, *The Go Programming Language* (2015), §4.3 (map elementining manzili), §6.2
- William Kennedy, [Methods, Interfaces and Embedded Types in Go](https://www.ardanlabs.com/blog/2014/05/methods-interfaces-and-embedded-types.html) (2014)

## Uy vazifasi
- 🟢 **Oson:** nega `3` emas, `0` chiqadi? Toping va tuzating.
  ```go
  type Hisoblagich struct{ n int }
  func (h Hisoblagich) Oshir()     { h.n++ }
  func (h Hisoblagich) Qiymat() int { return h.n }
  // var h Hisoblagich; h.Oshir(); h.Oshir(); h.Oshir(); fmt.Println(h.Qiymat())
  ```
- 🟡 **Oʻrta:** [homework/orta/main.go](homework/orta/main.go) dagi 1–8 raqamli qatorlardan qaysilari kompilyatsiya boʻladi?
  Avval oʻzingiz javob bering, keyin `go build -o /dev/null ./homework/orta` bilan tekshiring.
- 🟠 **Qiyin:** `b := Buyurtma{Summa: 100}; defer b.Chiqar(); defer func() { b.Chiqar() }(); b.Summa = 200`
  (`Chiqar` value receiver). Nima chiqadi va qaysi tartibda?
- 🔴 **Guru:** `type A struct{ Manzil }` va `type B struct{ *Manzil }`. `Kochir` pointer receiver'da. `A{}`, `&A{}`,
  `B{}`, `&B{}` dan qaysilari `interface{ Kochir(string) }` ni qanoatlantiradi? Kompilyator bilan isbotlang.

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): value receiver nusxani oshiradi. Yechim: `func (h *Hisoblagich) Oshir()` →
  `3`. Izchillik uchun `Qiymat()` ham pointer receiver. staticcheck xato variantni `SA4005` bilan ushlaydi.
- Oʻrta: 1, 2, 3, 4, 8 — ha; 5, 6, 7 — yoʻq (kompilyator 25, 26, 27-qatorlarni koʻrsatadi):
  5 `m["a"].X = 5` → `cannot assign to struct field m["a"].X in map` (map elementining maydoniga ham yozib boʻlmaydi);
  6 `f := m["a"].Sur` → method value ham `&` talab qiladi; 7 `ikki()[0].Sur(1)` → funksiya qaytargan massiv addressable
  emas, demak elementi ham. 3 va 4: value metod uchun addressable shart emas. 2 va 8: `mp["a"]` pointer, uning ortidagi
  qiymat addressable.
- Qiyin ([homework/qiyin](homework/qiyin/)): `200`, keyin `100`. `defer` LIFO. `defer b.Chiqar()` — method value,
  receiver (nusxa, 100) `defer` qatorida olinadi; closure esa `b` ni funksiya oxirida oʻqiydi (200).
- Guru ([homework/guru](homework/guru/)): `&A{}`, `B{}`, `&B{}` — ha; `A{}` — yoʻq: `A does not implement Kochiruvchi
  (method Kochir has pointer receiver)`. `B{}` ichidagi `*Manzil` nil boʻlsa, `Kochir` chaqiruvi panic beradi.
