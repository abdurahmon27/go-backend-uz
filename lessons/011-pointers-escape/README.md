# 11 · Pointer'lar: bu oʻzgaruvchi nega heap'ga qochdi?

Oldingi dars: [10 · Struct'lar va xotira](../010-structs-memory/) · Keyingi dars: GO-1.12 · Metodlar: value receiver yoki pointer receiver?

| Papka | Nima |
|---|---|
| `hook/` | `yarat` (qiymat qaytaradi) va `yaratP` (pointer qaytaradi), ikkalasi `//go:noinline` + benchmark |
| `inl/` | `hook/` bilan bir xil, faqat `yaratP` tepasidagi `//go:noinline` qatori oʻchirilgan: inlining'dan keyin pointer variant ham ajratmaydi |
| `pointer/` | nusxa vs pointer: `tola(b)` asl qiymatni oʻzgartirmaydi, `tolaP(&b)` oʻzgartiradi; `nil` |
| `escape/` | beshta holat, `go build -gcflags=-m ./escape` uchun |
| `savol/` | pauzadagi savol: `summa(&a)` va `eslab(&b)` (ikkalasi `//go:noinline`) — faqat `b` heap'ga chiqadi |
| `katta/` | 4 KiB'lik struct: qiymat bilan ≈50 ns, pointer bilan ≈1 ns (pointer haqiqatan tezroq) |
| `slaysp/` | slice elementiga pointer + `append`: **ataylab xato**, oʻzgarish yoʻqoladi |
| `slaysfix/` | `append` dan keyin indeks orqali |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/011-pointers-escape`) ishga tushiriladi.
```bash
sed -n 14,24p hook/buyurtma.go
go test -bench . -benchmem ./hook | grep ns/op        # Qiymat: 0 allocs, Pointer: 16 B/op, 1 allocs/op
go test -bench . -benchmem ./inl | grep ns/op         # izoh yoʻq: ikkalasi 0 allocs/op
go run ./pointer
go build -gcflags=-m ./escape 2>&1 | grep -v -e inlin -e argument
go build -gcflags=-m=2 ./escape 2>&1 | grep -A3 'b escapes'
sed -n 20,28p savol/savol.go
go build -gcflags=-m ./savol 2>&1 | grep moved        # faqat b (eslab'ga berilgan)
go test -gcflags=-m -run '^$' ./inl 2>&1 | grep -E '^inl/.*(yaratP|heap)'   # inlining call to yaratP + baribir "moved to heap: b"
cd ../006-functions && go build -o /dev/null -gcflags=-m ./hisob 2>&1 | grep -E 'moved|func literal|inlining call to hisob'
sed -n 12,16p slaysp/main.go
go run ./slaysp                                       # p.Tolandi: true, s[0].Tolandi: false
go vet ./slaysp && echo vet: jim
go run ./slaysfix
```
Natijalar `go1.26.8 linux/amd64`, Intel i5-12450H da. Nanosekundlar har safar biroz oʻzgaradi, `B/op` va
`allocs/op` oʻzgarmaydi. `-gcflags=-m` faqat koʻrsatilgan paketga qoʻllanadi (barcha paketlarga: `-gcflags=all=-m`).

## Dars matni

**Hook.** Bitta izoh qatorini (`//go:noinline`) oʻchirsak, heap'dan ajratish yoʻqoladi. `yarat` 16 baytli `Buyurtma` ni qiymat sifatida qaytaradi: 0 allocs/op, ≈0,7 ns. `yaratP` uning manzilini
qaytaradi: lokal `b` funksiyadan keyin ham kerak, shuning uchun kompilyator uni heap'ga qoʻyadi: 16 B/op,
1 allocs/op, ≈13 ns. `//go:noinline` kompilyatorga funksiyani chaqiruv joyiga
koʻchirishni (inline) taqiqlaydi: shunda har bir funksiyaning oʻzini oʻlchaymiz.

**Pointer.** Go'da hamma argument nusxalanadi (GO-1.06). `&b` — `b` ning manzili, `*p` — manzildagi qiymat.
`p.Tolandi` pointer orqali `(*p).Tolandi` bilan bir xil (spec: Selectors, misol `p.z // (*p).z`). Pointer'ning zero value'si `nil`; `nil` pointer
orqali `*p` yoki `p.Maydon` — runtime panic (`invalid memory address or nil pointer dereference`). Go'da pointer
arifmetikasi yoʻq (Go FAQ: *Why is there no pointer arithmetic?*): xavfsizlik uchun.

**Stack va heap.** Har bir chaqiruv stack'da frame oladi; funksiya qaytganda frame boʻshaydi. Heap'dagi qiymatni
garbage collector kuzatadi. Go FAQ: *"When possible, the Go compilers will allocate variables that are local to a
function in that function's stack frame. However, if the compiler cannot prove that the variable is not referenced
after the function returns, then the compiler must allocate the variable on the garbage-collected heap to avoid
dangling pointer errors."* Effective Go: *"unlike in C, it's perfectly OK to return the address of a local variable"*.
Bill Kennedy (2017) buni "qiymatni yuqoriga (chaqiruvchiga) ulashish" deb tushuntiradi; pastga (chaqirilgan
funksiyaga) ulashish odatda heap talab qilmaydi (`p does not escape`). Juda katta qiymat esa ulashilmasa ham heap'ga
tushishi mumkin (FAQ: *"if a local variable is very large, it might make more sense to store it on the heap"*).
C'da lokal oʻzgaruvchining manzilini qaytarish kompilyatsiya boʻladi, lekin funksiyadan keyin u pointer bilan
ishlash aniqlanmagan xatti-harakat (C11 §6.2.4p2: *"If an object is referred to outside of its lifetime, the
behavior is undefined"*). Kompilyatorning `escape.go` faylidagi ikki qoida: (1) heap'dagi obyekt stack'dagi
obyektga pointer saqlamaydi, (2) stack'dagi obyektga pointer undan uzoq yashamaydi.

**`-gcflags=-m`** (GO-1.01 da bir marta koʻrgan edik):
```
escape/escape.go:15:2: moved to heap: b                   # 1: return &b
escape/escape.go:20:12: p does not escape                 # 2: faqat oʻqildi
escape/escape.go:27:12: leaking param: p                  # 3: global oʻzgaruvchiga yozildi
escape/escape.go:33:15: b.Summa escapes to heap           # 4: fmt.Println(...any)
escape/escape.go:39:16: make([]byte, 64) does not escape  # 5a
escape/escape.go:40:15: make([]byte, 100000) escapes to heap  # 5b: > 64 KiB
escape/escape.go:41:17: make([]byte, n) does not escape   # 5c: Go 1.25+
```
4-holat haqida: kompilyator `b.Summa` ni "qochdi" deb belgilaydi (interface'ga joylanadi). Runtime 0–255 oraligʻidagi
sonlar uchun tayyor jadvaldan (`staticuint64s`) foydalanadi va xotira ajratmaydi (256 dan boshlab va manfiy sonlarda
1 alloc), shuning uchun aniq narxni
benchmark aytadi. `... argument does not escape` qatori `Println` oladigan `[]any` slice'i haqida (u stack'da qoladi),
videoda `grep -v argument` bilan yashirilgan.
5b: `cmd/compile/internal/ir/cfg.go`: `MaxImplicitStackVarSize = 64 * 1024` (`new`, `&T{}`, `make` bilan yaratiladigan
qiymatlar uchun; `var x T` uchun chegara 128 KiB). 5c: Keith Randall, *Allocating on the Stack* (Go blog,
2026-02-27): Go 1.25 dan beri oʻlchami nomaʼlum `make` uchun kompilyator stack'da kichik (hozircha 32 bayt) joy
ajratadi, sigʻmasa heap'dan oladi; Go 1.26 da `append` uchun ham. Eski qoida ("oʻlchami nomaʼlum boʻlsa — heap",
masalan ardanlabs/gotraining'dagi roʻyxat) Go 1.25 dan beri har doim ham toʻgʻri emas.

**Leaking param chaqiruvchiga taʼsir qiladi (`savol/`).** Kompilyator har bir funksiya uchun "parametr nima
boʻladi" degan qisqa xulosa (parameter tags) saqlaydi. `summa(&a)`: `a` stack'da. `eslab(&b)`: `b` heap'ga
(`moved to heap: b`), chunki `eslab` parametrni global oʻzgaruvchiga yozadi. `summa`/`eslab` da `//go:noinline` bor:
aks holda ular inline qilinadi va natija xulosadan emas, ichiga koʻchirilgan tanadan chiqadi (natija bir xil).
`slaysp` dagi `bir joymi? false`: `p` va `&s[0]` endi har xil manzil.

**Katta struct (`katta/`).** 4096 baytli `Hisobot` ni qiymat bilan uzatish ≈50 ns, pointer bilan ≈0,9 ns. Go
wiki CodeReviewComments, *Pass Values*: kichik qiymatlar uchun pointer'ni baytlarni tejash uchun ishlatmang, lekin
*"This advice does not apply to large structs"*.

**`-m=2`** sababni ham yozadi: `b escapes to heap in yangi: flow: ~r0 ← &b: from &b (address-of) … from return &b (return)`.

**Inlining qarorni oʻzgartiradi (`inl/`).** `yaratP` tepasidagi `//go:noinline` olib tashlansa, `yaratP` benchmark ichiga inline
qilinadi va buyurtma stack'da qoladi: 0 allocs/op. `-m` esa baribir `moved to heap: b` deydi: bu funksiyaning
alohida (inline qilinmagan) nusxasi haqida. Lekin inline qilishning chegarasi bor: kompilyator odatda kichik funksiyalarni inline qiladi (`inlineMaxBudget = 80`; PGO issiq funksiyalarga kattaroq byudjet beradi): bir nechta `if` va sikli bor
konstruktor direktivasiz ham heap'ga chiqadi. Xulosa: escape analysis har bir chaqiruv joyida (inlining'dan keyin)
qaror qiladi; `-m` chiqishini oʻqing, lekin narxni benchmark bilan oʻlchang. GO-1.06 dagi `hisob` ham xuddi shunday: alohida `hisoblagich` da `moved to heap: n` va `func literal escapes to
heap`, `main` ichida inline qilinganda esa `func literal does not escape` (GO-1.06 vaʼda qilgan javob).
**Diqqat:** `does not escape` har doim "xotira ajratilmaydi" degani emas: `make([]byte, n)` uchun `-m` "does not
escape" deydi, lekin `n > 32` boʻlsa runtime baribir heap'dan oladi (uy vazifasi guru). Oxirgi soʻz — benchmark.

**`&s[0]` + `append`.** `len == cap` boʻlsa, `append` yangi massiv ajratadi (spec: *"If the capacity of s is not
large enough … append allocates a new, sufficiently large underlying array"*). Eski massivga olingan pointer eski
massivda qoladi: `p.Tolandi = true` slice'ga yetib bormaydi. Kompilyator ham, `go vet` ham jim. Yechim: slice
elementiga pointer'ni uzoq saqlamang; `append` dan keyin indeks orqali murojaat qiling yoki pointer'ni qayta oling;
yoki `[]*Buyurtma` saqlang (har bir element alohida heap obyekti — oʻz narxi bor). Slice oʻsishi haqida batafsil:
GO-1.08.

**Qachon pointer.** Qiymatni oʻzgartirish kerak boʻlsa yoki nusxalash haqiqatan qimmat boʻlsa. Kichik struct'ni
nusxalash koʻpincha arzonroq, heap'ga chiqqan har bir qiymat esa GC uchun ish. GO-1.12 da metodlar uchun xuddi shu
savol: value receiver yoki pointer receiver.

## Manbalar
- [Go FAQ: How do I know whether a variable is allocated on the heap or the stack?](https://go.dev/doc/faq#stack_or_heap)
- [Effective Go: Composite literals](https://go.dev/doc/effective_go#composite_literals) (*"return the address of a local variable"*)
- [Go spec: Address operators](https://go.dev/ref/spec#Address_operators), [Selectors](https://go.dev/ref/spec#Selectors),
  [Appending to and copying slices](https://go.dev/ref/spec#Appending_and_copying_slices)
- [Go wiki: CodeReviewComments — Pass Values](https://go.dev/wiki/CodeReviewComments#pass-values)
- C11 (N1570) §6.2.4 Storage durations of objects: <https://www.open-std.org/jtc1/sc22/wg14/www/docs/n1570.pdf>
- Keith Randall, [Allocating on the Stack](https://go.dev/blog/allocation-optimizations) (Go blog, 2026)
- William Kennedy, [Language Mechanics On Stacks And Pointers](https://ardanlabs.com/blog/2017/05/language-mechanics-on-stacks-and-pointers.html),
  [Language Mechanics On Escape Analysis](https://ardanlabs.com/blog/2017/05/language-mechanics-on-escape-analysis.html) (2017)
- Go manba kodi: [cmd/compile/internal/escape/escape.go](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/escape/escape.go),
  [cmd/compile/internal/ir/cfg.go](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/ir/cfg.go)
- [cmd/compile](https://pkg.go.dev/cmd/compile) (`-m` bayrogʻi)

## Uy vazifasi
- 🟢 **Oson:** butun sonning manzilini qaytaradigan funksiya yozing. `-gcflags=-m` nima deydi? Avval taxmin qiling.
- 🟡 **Oʻrta:** `chop` funksiyasini `fmt.Println` siz yozing, shunda `-m` hech narsa qochdi demasin.
- 🟠 **Qiyin:** `&s[0]` + `append` xatosini test bilan ushlang va tuzating (ikkala holat: `cap` yetadi / yetmaydi).
- 🔴 **Guru:** oʻlchami nomaʼlum `make([]int64, n)` ning nechta elementi stack'ga sigʻadi? Benchmark bilan toping.

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): `moved to heap: x`.
- Oʻrta ([homework/orta](homework/orta/)): `strconv.AppendInt` stack'dagi `[24]byte` ga yozadi, keyin
  `os.Stdout.Write`: `-m` faqat `append does not escape` deydi.
- Qiyin ([homework/qiyin](homework/qiyin/)): xato faqat `cap` yetmaganda chiqadi, shuning uchun jadvalda ikkala
  holat bor. `qosh` (indeks orqali) ikkalasidan oʻtadi. `TestQosh` ni `qoshXato` bilan ishlatsangiz, "cap yetmaydi"
  holati yiqiladi; `TestXatoBor` xato haqiqatan borligini koʻrsatadi (`go test -v ./homework/qiyin`).
- Guru ([homework/guru](homework/guru/)): **4 ta**. Stack'dagi bufer 32 bayt, `int64` 8 bayt: `n = 4` → 0 B/op,
  0 allocs/op; `n = 5` (40 bayt) → 48 B/op, 1 allocs/op (48 baytlik size class). `make([]byte, n)` uchun chegara 32.
