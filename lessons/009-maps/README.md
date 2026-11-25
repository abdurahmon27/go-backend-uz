# 09 · Map: tartib nega tasodifiy va `recover` nega qutqarmaydi

Oldingi dars: [08 · Slice'lar: `append` tuzogʻi va eskirgan 1024 qoidasi](../008-slices/) · Keyingi dars: [10 · Struct'lar va xotira: bir xil maydonlar, 24 va 16 bayt](../010-structs-memory/)

| Papka | Nima |
|---|---|
| `tartib/` | bitta map'ni 100 000 marta `range` qilamiz: 5 xil tartib, bittasi ≈50 % |
| `asos/` | literal, yozish, yoʻq kalit → nol qiymat, comma ok, `delete`, `len` |
| `nilmap/` | **ataylab xato**: `nil` map'ga yozish — oddiy panic, `recover` uni ushlaydi |
| `chop/` | pauza-savol: `fmt.Println(map)` kalitlarni saralab chiqaradi |
| `saralash/` | tartib kerak boʻlsa: `slices.Sorted(maps.Keys(m))` |
| `poyga/` | **ataylab xato**: ikki goroutine bitta map'ga yozadi → `fatal error: concurrent map writes`, `recover` ojiz |
| `mutex/` | tuzatish: `sync.Mutex` |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/009-maps`) ishga tushiriladi.
```bash
go run ./tartib                      # 100000 ta sikl, 5 xil tartib, biri ≈50 %
go run ./asos                        # 0 / 0 false / 2
go run ./nilmap                      # qutqarildi: assignment to entry in nil map
sed -n 13,18p nilmap/main.go
go run ./chop                        # ikki qator bir xil, kalitlar saralangan
sed -n 11p saralash/main.go
go run ./saralash
sed -n 19,23p poyga/main.go
bash ./poyga/demo.sh                 # haqiqiy fatal error + exit status 2
sed -n 15,17p mutex/main.go
go run ./mutex                       # tugadi: 1000000
```
Natijalar `go1.26.8 linux/amd64`. `poyga` ning xatoni aniqlashi best effort: baʼzan boshqa xabar bilan yiqiladi yoki davom etadi.
Videoda `bash ./poyga/demo.sh` yordamchi skriptini ochiq ishlatamiz. U dasturni bir marta kompilyatsiya qiladi,
`fatal error: concurrent map writes` va chiqish kodi 2 olinguncha qayta ishga tushiradi (koʻpi bilan 20 urinish).
Faqat haqiqiy ishga tushirishdan olingan xabarni va chiqish kodini chiqaradi, oʻzi ham 2 kodi bilan tugaydi.
Uzun goroutine stack trace videodan olib tashlanadi. Yigʻish 10 soniya, har bir urinish yarim soniya bilan cheklangan. `-race` bu kompyuterda ishlamadi: Linux'da race detector cgo va C
kompilyator (gcc) talab qiladi, bizda gcc yoʻq. Uni 22-darsda koʻramiz.

## Dars matni

**Asoslar.** `m[k] = v` qoʻshadi yoki oʻzgartiradi. Yoʻq kalit xato emas: qiymat turining nol qiymati qaytadi
(`narx["olma"]` → `0`). "Narxi 0" va "mahsulot yoʻq"ni ajratish uchun comma ok: `n, bor := narx["olma"]` → `0 false`.
`delete(m, k)` oʻchiradi, `len(m)` kalitlar soni. `make` qilinmagan map `nil`: undan oʻqish mumkin (boʻsh map kabi),
yozish esa run-time panic (`assignment to entry in nil map`) — oddiy panic, `defer` + `recover` uni ushlaydi.
`defer`, `panic`, `recover` 15-darsda batafsil.

**Tartib nega tasodifiy.** Spec: "The iteration order over maps is not specified and is not guaranteed to be the same
from one iteration to the next." [Go 1](https://go.dev/doc/go1#iteration) da tartib ataylab "unpredictable" deb
eʼlon qilindi: "Code should not assume that the elements are visited in any particular order." Go 1.1 dagi yangi map
kichik map'larda (≤ 8) buni unutgan edi; [Go 1.3](https://go.dev/doc/go1.3): "Go 1.3 reintroduces random iteration
for small maps". Go 1.24 da map [Swiss table](https://go.dev/blog/swisstable)'ga oʻtdi; qoidalar oʻzgarmadi.

**Qanday?** (`internal/runtime/maps/table.go`, go1.26.8) Har bir `range` boshida `Iter.Init`:
`it.entryOffset = rand()`, `it.dirOffset = rand()`. Kichik map (≤ 8 kalit) bitta 8 oʻrinli guruhda turadi; kalitlar
birinchi boʻsh oʻringa yoziladi (`putSlotSmall`), yaʼni oʻchirish boʻlmasa — qoʻshilgan tartibda. `Iter.Next`:
`k := (entryIdx + entryOffset) % 8` — sikl tasodifiy oʻrindan boshlanib, guruhni aylana boʻylab yuradi, boʻsh
oʻrinlarni oʻtkazib. 5 kalit = 0–4 oʻrinlar; boshlanish 5, 6, 7 boʻsh oʻrinlarga tushsa, aylana 0-oʻringa qaytadi.
Shuning uchun 8 ta boshlanishdan 4 tasi `non …` tartibini beradi (≈50 %), qolgan 4 tartib ≈12,5 % dan:
```
100000 ta sikl, 5 xil tartib:
choy olma non sut tuxum   12.5%
non sut tuxum choy olma   50.0%
olma non sut tuxum choy   12.5%
sut tuxum choy olma non   12.5%
tuxum choy olma non sut   12.5%
```
(foizlarning oʻnli qismi har safar biroz oʻzgaradi). Tartib **tasodifiy**, lekin **aralashtirilmagan**: 5! = 120 emas,
5 ta, teng ehtimolli ham emas. Ikki ishga tushirish bir xil tartib berish ehtimoli ≈31 % (0,5² + 4·0,125²). Bunga ham
tayanmang: bu implementatsiya tafsiloti.

**`fmt` tuzogʻi.** [Go 1.12](https://go.dev/doc/go1.12#fmt) dan beri `fmt` map'ni kalitlari boʻyicha saralab
chiqaradi ("to ease testing"). Debug'da map tartibli koʻrinadi, `range` esa baribir tasodifiy. Tartib kerak boʻlsa:
`slices.Sorted(maps.Keys(m))` (ikkalasi Go 1.23; `maps.Keys` slice emas, iterator qaytaradi — eski
`golang.org/x/exp/maps.Keys` dan farqli). Saralash **bayt tartibida** (`<`), oʻzbek alifbosidagidek emas: `oʻ`, `gʻ`,
`sh`, `ch` alifboda `z` dan keyin, `sort` uchun esa `ch` — `c` harfidagi soʻz.

**Concurrent write — fatal error.** `m[i] = g` (`map[int]int`) aslida `runtime_mapassign_fast64` ni chaqiradi
(`go tool objdump -s main.yoz`: `CALL runtime.mapassign_fast64`; `internal/runtime/maps/runtime_fast64.go:196`).
Unda oʻn qator ichida ikkala holat bor:
```go
if m == nil {
	panic(errNilAssign)            // nil map: oddiy panic, recover ushlaydi
}
...
if m.writing != 0 {
	fatal("concurrent map writes") // recover ushlamaydi
}
```
(`string` kalit uchun `runtime_faststr.go:265`, xuddi shunday.) `fatal` (`runtime/panic.go`) — `throw` bilan bir xil,
"used when user code is expected to be at fault for the failure, such as racing map writes": xabar + stack trace,
keyin `fatalthrow` → `exit(2)`. `defer`lar ishlamaydi, `recover` ushlay olmaydi. Videodagi dasturda `recover` aynan
yiqilayotgan goroutine ichida (`yoz`); boshqa goroutine'dagi `recover` oddiy panic'ni ham ushlay olmasdi, shuning
uchun `nilmap` va `poyga` bir xil tuzilgan. "Buzilgan map bilan davom etish xavfli" — bu bizning izohimiz; manba
faqat "dasturchi aybi" deydi.

**Best effort.** [Go 1.6](https://go.dev/doc/go1.6#runtime): "lightweight, best-effort detection of concurrent misuse
of maps" — kafolat emas. `m.writing ^= 1` oddiy (atomik boʻlmagan) almashtirish: ikki yozuvchi uni bir vaqtda
almashtirsa, tekshiruv buni sezmasligi mumkin. Bizda `GOMAXPROCS=1` bilan 100 dan 2 ta ishga tushirish yiqilmay
tugadi (critic oʻlchovi); 300 ta oddiy ishga tushirishdan 2 tasi `small map with no empty slot (concurrent map writes?)`,
1 tasi `unexpected fault address` (SIGSEGV) bilan yiqildi — buzilgan map haqiqatan xavfli. Ancha yaxshiroq vosita: [race detector](https://go.dev/doc/articles/race_detector) (`go test -race`),
lekin u ham faqat ishga tushgan kod yoʻlidagi race'larni topadi ("only finds races that happen at runtime"). [FAQ](https://go.dev/doc/faq#atomic_maps): "Map access
is unsafe only when updates are occurring" — faqat oʻqish xavfsiz.

**Haqiqiy dunyo.** [Uber: Data Race Patterns in Go](https://www.uber.com/blog/data-race-patterns-in-go/)
(M. Chabbi, M. K. Ramanathan, 2022): qariyb 2000 data race topilgan, ~1100 tasi tuzatilgan (210 dasturchi, 6 oy) va
tahlil qilingan (patternlar chastota boʻyicha tartiblanmagan; sarlavha "frequent" deydi). "Concurrent accesses to Go's built-in, thread-unsafe maps cause frequent data races"; dasturchilar
`table[key]` sintaksisini alohida elementlarga murojaat deb tushungan.

## Manbalar
- Go spec: [Map types](https://go.dev/ref/spec#Map_types), [For statements](https://go.dev/ref/spec#For_statements)
- [Go maps in action](https://go.dev/blog/maps) (A. Gerrand, 2013); [Faster Go maps with Swiss Tables](https://go.dev/blog/swisstable) (M. Pratt, 2025)
- Release notes: [Go 1](https://go.dev/doc/go1#iteration), [1.3](https://go.dev/doc/go1.3), [1.6](https://go.dev/doc/go1.6#runtime), [1.12](https://go.dev/doc/go1.12#fmt), [1.24](https://go.dev/doc/go1.24)
- Go manba (go1.26.8): `internal/runtime/maps/table.go` (`Iter`), `internal/runtime/maps/runtime_fast64.go`, `runtime_faststr.go` (`mapassign`), `map.go` (`putSlotSmall`), `runtime/panic.go` (`fatal`)
- [Go FAQ: atomic maps](https://go.dev/doc/faq#atomic_maps); [Data Race Detector](https://go.dev/doc/articles/race_detector)
- [Uber Engineering: Data Race Patterns in Go](https://www.uber.com/blog/data-race-patterns-in-go/) (2022)

## Uy vazifasi
- 🟢 **Oson:** comma ok bilan narxi 0 boʻlgan mahsulotni ("bepul") yoʻq mahsulotdan ajrating.
- 🟡 **Oʻrta:** matndagi soʻzlarni map bilan sanang va natijani saralab chiqaring.
- 🟠 **Qiyin:** 9 ta kalitli map'da nechta xil tartib chiqishini oʻlchang va natijani tushuntiring.
- 🔴 **Guru:** `poyga` dasturini mutex'siz tuzating. (Maslahat: map'ni bir vaqtda faqat bitta goroutine ushlasin.)

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): `n, bor := narx[nom]`; `!bor` → "yoʻq", `n == 0` → "bepul".
- Oʻrta ([homework/orta](homework/orta/)): `sanoq[soz]++` — yoʻq kalit 0 qaytargani uchun alohida tekshiruv kerak emas;
  chiqarish `slices.Sorted(maps.Keys(sanoq))` bilan: `choy 1 / non 3 / sut 2`.
- Qiyin ([homework/qiyin](homework/qiyin/)): bizda (100 000 sikl) 5 → 5, 8 → 8, 9 → 9, 16 → 16 xil tartib.
  9 kalit endi bitta guruhga sigʻmaydi, lekin iterator baribir butun jadval boʻylab tasodifiy oʻrindan aylanadi:
  tartiblar soni kalitlar soniga teng, `n!` emas. Katta map'larda (bir nechta jadval) `dirOffset` ham qoʻshiladi.
- Guru ([homework/guru](homework/guru/)): har bir goroutine oʻz map'ini toʻldiradi va slice'ning **oʻz elementiga**
  yozadi (`qismlar[g]`, turli elementlar — race yoʻq), `wg.Wait()` dan keyin `main` ularni birlashtiradi. Map'lar hech
  qachon bir vaqtda ikki goroutine qoʻlida boʻlmaydi.
