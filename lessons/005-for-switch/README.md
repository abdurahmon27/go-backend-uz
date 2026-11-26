# 05 · Go'da `while` yoʻq: bitta `for` va uning hamma shakllari

Oldingi dars: [04 · `len("oʻzbek")` = 7?](../004-strings-utf8/) · Keyingi dars: [06 · Funksiyalar: ikki qiymat qaytarish, closure va variadic](../006-functions/)

| Papka | Nima |
|---|---|
| `shakllar/` | `for` ning toʻrt shakli: uch qismli, faqat shart, cheksiz, `range 3` |
| `nusxa/` | pauza-savol: `range` qiymati nusxa, slice oʻzgarmaydi |
| `tanlov/` | `switch`: bir nechta qiymat, `fallthrough`, `default` |
| `break/` | `switch` ichidagi `break` sikldan chiqmaydi (dastur ishlaydi, lekin notoʻgʻri) |
| `label/` | yechim: `break qidiruv` |
| `manzil/` | sikl oʻzgaruvchisining manzili: Go 1.22+ da `a.uz b.uz c.uz` |
| `eski/` | xuddi shu kod, lekin alohida modul, `go.mod` da `go 1.21`: `c.uz c.uz c.uz` |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
```bash
go run ./shakllar
go run ./nusxa        # pauza-savol javobi
go vet ./nusxa        # hech narsa demaydi
$(go env GOPATH)/bin/staticcheck ./nusxa ./break   # go install honnef.co/go/tools/cmd/staticcheck@2026.2.1
go run ./tanlov
go run ./break
go run ./label
go run ./manzil
grep '^go ' ../../go.mod eski/go.mod
cd eski && go run . && go vet . && cd ..
```
Natijalar `go1.26.8 linux/amd64` da olingan.

## Dars matni

**Bitta `for`.** Go'da 25 ta kalit soʻz bor, `while` ular orasida yoʻq ([spec: Keywords](https://go.dev/ref/spec#Keywords)).
Hamma sikl `for` bilan yoziladi ([spec: For statements](https://go.dev/ref/spec#For_statements)):

```go
for i := 0; i < 3; i++ { … } // uch qismli: boshlanish; shart; qadam
for n < 100 { … }            // faqat shart: boshqa tillardagi while
for { … }                    // cheksiz: break yoki return bilan chiqiladi
for i := range 3 { … }       // son ustida: 0, 1, 2 (Go 1.22+)
```

```
0 1 2 | 128
-22
0 1 2
```

Shart atrofida qavs `( )` kerak emas (`for (n < 100)` ham kompilyatsiya boʻladi, lekin `gofmt` qavsni olib tashlaydi; uch qismli siklni qavsga olsangiz esa sintaksis xatosi), jingalak qavs `{ }` esa majburiy. `range` Go'da boshidan bor, Go 1.22 dan beri u butun son ustida ham yuradi.

**`range` nima ustida yuradi** ([spec: For range](https://go.dev/ref/spec#For_range)):

| Ifoda | 1-qiymat | 2-qiymat |
|---|---|---|
| slice, array | indeks | element (nusxa!) |
| string | rune boshlangan bayt oʻrni | rune (04-dars) |
| map | kalit | qiymat; tartib kafolatlanmagan (09-dars) |
| son `n` (Go 1.22+) | `0 … n-1` | — |
| funksiya (Go 1.23+) | iterator'lar (18-dars) | |
| channel | qiymat (goroutine darslarida) | — |

**Pauza-savol.** `for _, n := range narxlar { n *= 2 }` dan keyin `narxlar` = `[100 200 300]`: oʻzgarmadi.
`range` har qadamda elementning **nusxasini** `n` ga yozadi. Oʻzgartirish uchun indeks orqali yozing:
`for i := range narxlar { narxlar[i] *= 2 }` → `[200 400 600]`.
`go vet` bu yerda jim. `staticcheck` (2026.2.1) esa topadi: `nusxa/main.go:8:3: this value of n is never used (SA4006)`.

**`switch`.** C'dan va Java'ning klassik `switch` idan farqli, har bir `case` oʻzi tugaydi, keyingisiga oʻtmaydi, `break` shart emas.
Bitta `case` da bir nechta qiymat: `case 6, 7:`. Keyingi `case` ham bajarilsin desangiz, `fallthrough` deb ochiq
yozasiz (u keyingi case'ning shartini **tekshirmaydi**). Ifodasiz `switch { … }` birinchi `true` boʻlgan `case` ni
bajaradi: uzun `if/else if` zanjirining oʻrniga ([spec: Switch statements](https://go.dev/ref/spec#Switch_statements)).

```
1 hafta boshi, ish kuni
3 ish kuni
7 dam olish
9 bunday kun yoʻq
```

**`switch` ichidagi `break`.** *"A "break" statement terminates execution of the innermost "for", "switch", or "select"
statement within the same function"* ([spec: Break statements](https://go.dev/ref/spec#Break_statements)).
`switch` sikl ichida boʻlsa, `break` faqat `switch` dan chiqadi, sikl davom etadi. `go vet` jim, `staticcheck` topadi
(`break/main.go:11:4: ineffective break statement. Did you mean to break out of the outer loop? (SA4011)`). Xuddi shu
tuzoq `for { select { … break } }` da ham bor (goroutine darslari):

```
juft: 4
manfiy: -1 → toʻxtaymiz
juft: 8
manfiy: -5 → toʻxtaymiz
tugadi
```

`staticcheck` ning ikkinchi qatori (`S1023: redundant break statement`) esa bu `break` umuman ortiqcha ekanini aytadi.

Tarixdan: 1990-yil 15-yanvarda AT&T'ning uzoq masofali telefon tarmogʻi qariyb toʻqqiz soat nosoz ishladi.
RISKS Digest'da tarqalgan bir xatga koʻra (AT&T'ning rasmiy bayonoti emas), C dagi `do … while` ichida `switch`, uning ichida `if`, unda esa
`if` uchun moʻljallangan `break` boʻlgan, u esa `switch` dan chiqib ketgan ([RISKS 9.69](https://catless.ncl.ac.uk/risks/9/69)).

**Label.** Ishlatilmagan label kompilyatsiya xatosi (`label qidiruv defined and not used`). Koʻpincha yana bir yoʻl:
siklni alohida funksiyaga chiqarib, `return` qilish. `qidiruv:` deb siklga nom bering va `break qidiruv` yozing: endi aynan shu `for` tugaydi.
`continue qidiruv` esa tashqi siklning keyingi qadamiga oʻtadi ([spec: Labeled statements](https://go.dev/ref/spec#Labeled_statements)).

```
juft: 4
manfiy: -1 → toʻxtaymiz
tugadi
```

**Go 1.22: har iteratsiyada yangi oʻzgaruvchi.** `manzil/` sikl oʻzgaruvchisi `d` ning manzilini (`&d`) saqlaydi
(pointer'lar: 11-dars). Bizning modul (`go 1.26`): `a.uz b.uz c.uz`. `eski/` (`go 1.21`): `c.uz c.uz c.uz`.
Go 1.22 gacha butun sikl uchun **bitta** `d` boʻlgan, har qadam uning ustidan yozgan, uchala manzil bitta oʻzgaruvchiga
qaragan. Go 1.22 dan beri: *"Each iteration has its own separate declared variable (or variables)"* (spec). Faqat **sikl oʻzi
eʼlon qilgan** oʻzgaruvchi: `var j int; for j = 0; j < 3; j++ { … &j }` Go 1.26 da ham `3 3 3` (biz tekshirdik).
Bu oʻzgarish faqat `go.mod` dagi `go` qatori `1.22` yoki undan yuqori boʻlgan modullarda ishlaydi
([Go blog: Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview)). Eski modulga koʻchirilgan kod xatoni qaytaradi.
`go vet` (`loopclosure`) ikkala holatda ham jim: u faqat `go`/`defer` (va errgroup, `t.Parallel`) ichidagi closure'larni,
faqat Go 1.22 dan oldingi semantikali fayllarda tekshiradi; `&d` ni slice'ga qoʻshishni umuman koʻrmaydi.
`go 1.22+` modulda `go fix` (Go 1.26) eski `d := d` qatorlarini oʻzi oʻchiradi (biz tekshirdik).

**Let's Encrypt, 2020.** Sertifikat berishdan oldin server (Boulder) har bir domenning CAA yozuvini qayta tekshirishi
kerak edi. Kod `for k, v := range m` ichida `&v` ni saqlagan, shuning uchun bir nechta domenli soʻrovda bitta domen
N marta tekshirilgan, qolganlari esa qayta tekshirilmagan. Eng qizigʻi: kodda `k` uchun izoh bilan nusxa olingan
(`// Make a copy of k because it will be reassigned with each loop.` → `kCopy := k`), `v` esa unutilgan.
(Let's Encrypt saytlarga HTTPS uchun bepul sertifikat beradi; CAA, [RFC 8659](https://www.rfc-editor.org/rfc/rfc8659),
domen egasi qaysi tashkilot sertifikat bera olishini yozadigan DNS yozuvi.) `go vet` ham, `staticcheck` ham bu kodda jim;
Let's Encrypt ham staticcheck bu variantni topmasligini yozgan (gosec topar edi). Incident report: *"The proximate cause of the bug was a common
mistake in Go: taking a reference to a loop iterator variable."* Xato 2020-yil 29-fevralda topildi (03:08 UTC),
tuzatish 05:22 UTC da ishga tushdi. 3 048 289 ta sertifikatga taʼsir qildi, 1 711 396 tasi muddat ichida bekor
qilindi, qolganlari 2020-yil 29-maygacha muddati tugab yoki bekor qilinib yoʻqoldi
([Mozilla 1619047](https://bugzilla.mozilla.org/show_bug.cgi?id=1619047), [1619179](https://bugzilla.mozilla.org/show_bug.cgi?id=1619179),
[Let's Encrypt eʼloni](https://community.letsencrypt.org/t/2020-02-29-caa-rechecking-bug/114591)).

**Go jamoasi qanday qaror qildi.** Russ Cox: Google ichida ~100 000 testda yangi semantika sinaldi, 58 tasi oʻtmadi,
ulardan faqat 2 tasida kod eski xulqqa haqiqatan tayangan, qolganlarining deyarli hammasi yashirin xatolar edi
(36 tasi `t.Parallel` bilan notoʻgʻri testlar, ikkinchi eng koʻp sabab esa aynan `&v` ni slice'ga qoʻshish)
([discussion #56010](https://github.com/golang/go/discussions/56010)).

## Manbalar
- Go spec: [Keywords](https://go.dev/ref/spec#Keywords), [For statements](https://go.dev/ref/spec#For_statements),
  [For range](https://go.dev/ref/spec#For_range), [Switch statements](https://go.dev/ref/spec#Switch_statements),
  [Break statements](https://go.dev/ref/spec#Break_statements), [Labeled statements](https://go.dev/ref/spec#Labeled_statements)
- David Chase, Russ Cox: [Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview) (Go blog, 2023)
- [Go 1.22 Release Notes](https://go.dev/doc/go1.22), [LoopvarExperiment](https://go.dev/wiki/LoopvarExperiment)
- [RISKS Digest 9.69](https://catless.ncl.ac.uk/risks/9/69) (AT&T, 1990)
- Mozilla Bugzilla [1619047](https://bugzilla.mozilla.org/show_bug.cgi?id=1619047), [1619179](https://bugzilla.mozilla.org/show_bug.cgi?id=1619179);
  [Let's Encrypt: 2020.02.29 CAA Rechecking Bug](https://community.letsencrypt.org/t/2020-02-29-caa-rechecking-bug/114591)
- Russ Cox, [discussion #56010: redefining for loop variable semantics](https://github.com/golang/go/discussions/56010)
- Bill Kennedy, [For Range Semantics](https://ardanlabs.com/blog/2017/06/for-range-semantics.html) (Ardan Labs, 2017; Go 1.22 dan oldin yozilgan)

## Uy vazifasi
- 🟢 **Oson:** 1 dan 15 gacha FizzBuzz'ni ifodasiz `switch` bilan yozing (3 ga boʻlinsa `Fizz`, 5 ga `Buzz`, ikkalasiga `FizzBuzz`).
- 🟡 **Oʻrta:** `switch` bilan sonning har bir raqamini oʻzbekcha soʻzga aylantiring (`2026` → `ikki nol ikki olti`).
- 🟠 **Qiyin:** ikki oʻlchamli jadvalda birinchi `0` ni toping va label bilan ikkala sikldan chiqing.
- 🔴 **Guru:** `eski/` dagi xatoni `go.mod` ga tegmasdan tuzating (uch usul bor).

Avval oʻzingiz urinib koʻring. Javoblar: [homework/](homework/).
