# 08 · Slice'lar: `append` tuzogʻi va eskirgan 1024 qoidasi

Oldingi dars: [07 · Taxmin qilmang, oʻlchang: birinchi test va benchmark](../007-test-bench/) · Keyingi dars: [09 · Map: tartib nega tasodifiy va `recover` nega qutqarmaydi](../009-maps/)

| Papka | Nima |
|---|---|
| `savat/` | **ataylab xato**: ikki savat bitta umumiy roʻyxatdan `append` qilinadi, Alining tuxumi Valining choyi bilan almashadi |
| `sabab/` | `len`, `cap`, ikkala savat bitta massivdami, slice oʻzi necha bayt (`unsafe.Sizeof`) |
| `pauza/` | pauza-savol: roʻyxat literal bilan yozilsa (`cap == len`), tuzoq yoʻq |
| `tuzat/` | tuzatish: `asos[:len(asos):len(asos)]` (yoki `slices.Clip`) |
| `osish/` | `[]int` ga bittalab `append`: `cap` 128 dan 2560 gacha qanday oʻzgaradi |
| `boshi/` | boʻsh slice'ga birinchi `append`: Go 1.26 da `cap` 1 emas, 4 (slice funksiyadan chiqmaydi) |
| `qochish/` | xuddi shu, lekin slice `Println`'ga beriladi (heap'ga qochadi): `1 2 4 4 8` |
| `oldindan/` | benchmark: boʻsh slice'dan oʻstirish vs `make([]int, 0, n)` |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/008-slices`) ishga tushiriladi.
```bash
sed -n 6,10p savat/main.go
go run ./savat                       # ali va vali: ikkalasida choy
go run ./sabab                       # len 3, cap 4, bitta massiv, 24 bayt
go run ./pauza                       # literal: hammasi toʻgʻri
sed -n 8p tuzat/main.go
go run ./tuzat                       # tuzatilgan: tuxum va choy
grep -n 'm+sepSave : m+sepSave' $(go env GOROOT)/src/bytes/bytes.go
go run ./osish                       # 256 → 512 → 848 → 1280 ...
grep -E '^//\s+4[89] ' $(go env GOROOT)/src/internal/runtime/gc/sizeclasses.go | awk '{print "class", $2 ":", $3, "bayt"}'
go run ./boshi                       # 4 4 4 4 8
go run ./qochish                     # 1 2 4 4 8 [0 1 2 3 4]
go test -bench . -benchmem -count 3 ./oldindan | grep ns/op | column -t
```
Natijalar `go1.26.8 linux/amd64` da olingan. Benchmark'da nanosekundlar har safar oʻzgaradi (bizda kompyuter band
paytda 3 000 dan 17 000 ns/op gacha sakradi), `B/op` va `allocs/op` esa har safar bir xil: 25152 B, 9 ta va 8192 B, 1 ta.

## Dars matni

**Tuzoq.** `asos := []string{"non", "sut"}` (len 2, cap 2), keyin `asos = append(asos, "suv")`: joy yetmadi,
`append` yangi massivni **ikki barobar** katta ajratdi, endi len 3, cap 4. `ali := append(asos, "tuxum")` boʻsh
4-joyga yozadi; `vali := append(asos, "choy")` **xuddi shu** joyga yozadi. Ikkala natija bitta massivga qaraydi:
`[non sut suv choy]` ikki marta. Kompilyator ham, `go vet` ham jim.

**Slice ichida nima bor.** Slice — kichik tuzilma: massivga pointer, `len` (nechta element) va `cap` (massivda
shu joydan boshlab nechta joy). amd64 da har biri 8 bayt, jami 24 bayt (`unsafe.Sizeof`). Elementlar slice ichida
emas, alohida massivda. Slice'ni funksiyaga bersangiz yoki boshqa oʻzgaruvchiga yozsangiz, shu 24 bayt nusxalanadi,
massiv esa umumiy qoladi. `cap` — slice boshidan massiv oxirigacha joylar soni (spec: "the sum of the length of the slice and the length of the array beyond the slice"), massivning butun hajmi emas (06-darsdagi `mediana` aynan shuning uchun chaqiruvchining slice'ini saralab yubordi).

**`append` qoidasi** (Go spec): sigʻim yetmasa — yangi, "sufficiently large" massiv va nusxa; yetsa — oʻsha massivning oʻzi. Qanchalik katta (×2 va h.k.) — spec'da yoʻq, bu runtime'ning tanlovi (`nextslicecap`).
Shuning uchun bir xil kod `cap > len` boʻlganda buziladi, `cap == len` boʻlganda ishlaydi (`pauza/`). `cap` esa
oʻsish qoidasiga va Go versiyasiga bogʻliq — testda oʻtib, boshqa maʼlumotda buziladigan xato shunday paydo boʻladi.

**Tuzatish.** Toʻliq slice ifodasi `a[low:high:max]` natijaning sigʻimini `max-low` qiladi.
`asos[:len(asos):len(asos)]` (= `slices.Clip(asos)`, Go 1.21) dan keyin boʻsh joy yoʻq, har bir `append` yangi
massiv ajratadi. Umumiy holda `a[low:high:max]` da `cap = max-low`; bu yerda `max = len`, shuning uchun `cap = len`.
Boshqa yoʻl: **har bir savat** oʻz nusxasini olsin: `ali := append(slices.Clone(asos), "tuxum")`. Diqqat: `asos` ni
**bir marta** klonlab, ikkala savatni shu klondan `append` qilsangiz, tuzoq qaytishi mumkin — `slices.Clone` hujjati:
"The result may have additional unused capacity".

**Standart kutubxona ham.** 2017-yilda [golang/go#21149](https://github.com/golang/go/issues/21149): `bytes.Split`
natijasidagi bitta boʻlakka `append` qilinsa, qoʻshni boʻlak buzilishi mumkin edi (avval ajratgich ustiga yoziladi,
kattaroq `append` qoʻshni boʻlakka yetadi). [Go 1.10](https://go.dev/doc/go1.10#bytes):
"`Fields`, `FieldsFunc`, `Split`, `SplitAfter` … each returned subslice to have capacity equal to its length, so that
appending to one cannot overwrite adjacent data". Bugungi manba: `a[i] = s[: m+sepSave : m+sepSave]`.

**1024 qoidasi.** Eski maqolalarda: "1024 gacha ×2, keyin ×1,25". Bu [Go 1.17 gacha](https://github.com/golang/go/blob/go1.17/src/runtime/slice.go)
toʻgʻri edi. 2021-yil sentabrda golang-nuts'da [“slices grow at 25% after 1024 but why 1024?”](https://groups.google.com/g/golang-nuts/c/UaVlMQ8Nz3o)
savoli berildi (Rob Pike: "It was just arbitrary, I'm sure."); Keith Randall'ning commit'i shu muhokamaga havola qiladi ("See discussion at …") va
[formulani almashtirdi](https://go.googlesource.com/go/+/2dda92ff6f9f07eeb110ecbf0fc2d7a0ddd27f9d) (Go 1.18):
```go
const threshold = 256
if oldCap < threshold {
	return doublecap
}
for {
	newcap += (newcap + 3*threshold) >> 2   // = newcap/4 + 192
	...
```
256 dan kichik — ×2; keyin koeffitsiyent asta-sekin 2 dan 1,25 ga tushadi (commit jadvali: 256 → 2,0; 512 → 1,63;
1024 → 1,44; 2048 → 1,35; 4096 → 1,30). Eslatma: Go 1.17 da ham haqiqiy `cap`lar ×1,25 dan biroz katta edi —
size class yaxlitlashi tufayli (commit: "both before and now"). Bizning oʻlchov (`osish/`):
```
  128 →   256  ×2.00
  256 →   512  ×2.00
  512 →   848  ×1.66
  848 →  1280  ×1.51
 1280 →  1792  ×1.40
 1792 →  2560  ×1.43
```
**Nega 832 emas, 848?** Formula 512 → 832 beradi, `832 × 8 = 6656` bayt. Xotira ajratuvchi faqat tayyor oʻlchamlarni
(size class) beradi: 6528 yoki 6784 (`internal/runtime/gc/sizeclasses.go`, class 48 va 49). `growslice`
`roundupsize` bilan yuqoriga yaxlitlaydi: 6784 bayt = **848** int. Keyingisi: 848 → 1252 → 10016 bayt → 10240 = 1280.
Shuning uchun koeffitsiyent tekis kamaymaydi: 1280 → 1792 ×1,40, keyin 1792 → 2560 ×1,43 (formula 2432 → 19456 bayt →
class 20480 → 2560). Commit ham shuni aytadi: "the real growth factor … is somewhat larger because we round up to the next size class".

**Boshi: nega 4?** `boshi/` da boʻsh slice'ga birinchi `append` dan keyin `cap` = 4 (`4 4 4 4 8`). Go 1.26 da, agar
slice funksiyadan chiqmasa (escape analysis), kompilyator boʻsh/nil slice'ning birinchi `append`i uchun stekda kichik
zaxira joy (32 bayt) ajratadi
([go.dev/blog/allocation-optimizations](https://go.dev/blog/allocation-optimizations), K. Randall, 2026;
`runtime/slice.go` `growsliceBuf`: "we try to just bump up to the next size class"). Slice "qochsa" boshqa ketma-ketlik
chiqadi: `qochish/` (slice `Println`'ga beriladi) → `1 2 4 4 8`; lokal slice oxirida global'ga yozilsa `1 2 3 4 8`, toʻgʻridan-toʻgʻri global'ga `append` qilinsa `1 2 4 4 8`.
Funksiyadan **qaytariladigan** slice esa ham stekdagi zaxiradan boshlaydi: oxirida kompilyator `runtime.move2heap`
qoʻyadi va slice heap'ga bir marta koʻchadi (blog: "for escaping slices … runtime.move2heap"). Shuning uchun
`oldindan/` dagi `yigNil` 11 emas, 9 marta ajratadi. Xulosa: `cap`ning aniq
qiymatiga hech qachon tayanmang; u Go 1.18 da ham, 1.26 da ham oʻzgardi.

**Oldindan joy soʻrash** (`oldindan/`). 1000 ta son: boʻsh slice'dan — 9 marta ajratish
(8, 16, 32, 64, 128, 256, 512, 848, 1280 elementlik massivlar, jami `3144 × 8 = 25152` bayt; dastlabki 4 element stekdagi zaxira joyda turadi),
`make([]int, 0, 1000)` — 1 marta, 8192 bayt (8000 → size class 8192). Funksiyalar `//go:noinline`: aks holda
kompilyator ularni benchmark ichiga koʻchiradi va oldindan soʻralgan slice'ni stekka joylaydi (0 allocs), bu esa
natija funksiyadan qaytadigan haqiqiy dasturdagi holat emas (pragma'siz bizda: 12 allocs va 0 allocs). Tinch
kompyuterda `Make` ≈ 2,3 barobar tez (≈1 400 va ≈3 100–3 700 ns/op).

Tekshirildi: Go 1.27.4 da ham `nextslicecap`, `growsliceBuf` va size class 48/49 oʻzgarmagan (bytes.go qatori 393).

## Manbalar
- Go spec: [Appending to and copying slices](https://go.dev/ref/spec#Appending_and_copying_slices), [Slice expressions](https://go.dev/ref/spec#Slice_expressions)
- [Go Slices: usage and internals](https://go.dev/blog/slices-intro) (A. Gerrand, 2011); [Arrays, slices (and strings): The mechanics of 'append'](https://go.dev/blog/slices) (R. Pike, 2013)
- Go manba (go1.26.8): `runtime/slice.go` (`growslice`, `nextslicecap`, `growsliceBuf`), `runtime/msize.go` (`roundupsize`), `internal/runtime/gc/sizeclasses.go`
- [commit 2dda92f: runtime: make slice growth formula a bit smoother](https://go.googlesource.com/go/+/2dda92ff6f9f07eeb110ecbf0fc2d7a0ddd27f9d) (2021); [Go 1.18 release notes](https://go.dev/doc/go1.18)
- [golang/go#21149](https://github.com/golang/go/issues/21149); [Go 1.10 release notes: bytes](https://go.dev/doc/go1.10#bytes)
- [Allocating on the Stack](https://go.dev/blog/allocation-optimizations) (K. Randall, 2026); [Go 1.26 release notes](https://go.dev/doc/go1.26)
- [pkg.go.dev/slices](https://pkg.go.dev/slices): `Clip`, `Clone`
- Teiva Harsanyi, [100 Go Mistakes](https://100go.co/) #25 "Unexpected side effects using slice append"

## Uy vazifasi
- 🟢 **Oson:** ishga tushirmasdan ayting, keyin tekshiring: `s := make([]int, 3, 5); t := append(s, 7); s[0] = 9; fmt.Println(t[0], len(t), cap(t))`.
- 🟡 **Oʻrta:** hech qachon chaqiruvchining massiviga yozmaydigan `qosh(s []string, v string) []string` yozing va test bilan isbotlang.
- 🟠 **Qiyin:** `[]byte` oʻsishini oʻlchang va nega `cap`lar `[]int` dan farq qilishini tushuntiring.
- 🔴 **Guru:** katta fayldan kichik boʻlak qaytaradigan funksiya butun faylni xotirada ushlab qoladi. Tuzating va isbotlang.

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): `9 4 5`. `s` ning cap'i 5, shuning uchun `append` oʻsha massivga yozdi;
  `s[0] = 9` `t[0]` ni ham oʻzgartiradi.
- Oʻrta ([homework/orta](homework/orta/)): `return append(s[:len(s):len(s)], v)`. Test `cap > len` boʻlgan roʻyxat
  bilan tekshiradi: ikkala natija toʻgʻri, asl massivning 4-joyi boʻsh qoldi (`go test ./homework/orta`).
- Qiyin ([homework/qiyin](homework/qiyin/)):
  ```
  []byte: 32 64 128 256 512 896 1408 2048
  []int:  4 8 16 32 64 128 256 512 848 1280 1792 2560
  ```
  Formula bir xil, lekin yaxlitlash **baytlarda**: `[]byte` da 512 → 832 bayt → size class 896; 896 → 1312 → 1408;
  1408 → 1952 → 2048. Boshidagi 32 — stekdagi 32 baytlik zaxira (`[]int` da bu 4 ta element).
- Guru ([homework/guru](homework/guru/)): `fayl[:i]` 64 MB massivga qaraydi, boʻlak tirik ekan massiv ham tirik
  (`heap ≈ 64 MB`). `slices.Clone` faqat kerakli baytlarni nusxalaydi, massiv GC'ga ketadi (`heap ≈ 0 MB`). Bu Go
  blogidagi `FindDigits` misolining oʻzi.
