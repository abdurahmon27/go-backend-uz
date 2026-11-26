# 18 · Iterator'lar: `range` endi funksiya ustida ham ishlaydi

Oldingi dars: GO-1.17 · Generics va generic metodlar: qachon kerak emas · Keyingi dars: GO-1.19 · `io.Reader` va
`io.Writer`: Go'ning eng chiroyli ikki qatori. Asos: GO-1.06 (funksiya qiymat sifatida, closure) va GO-1.17 (generics,
`cmp.Ordered`).

> **Go versiyasi.** Range over func **Go 1.23** dan bor (Go 1.22 da faqat `GOEXPERIMENT=rangefunc` bilan). Bu papka
> asosiy modulda (`go 1.26`). `eski/` — alohida modul (`go 1.22`): xuddi shu kod u yerda kompilyatsiya boʻlmaydi.

| Papka | Nima |
|---|---|
| `hook/` | `Tolanganlar(bb) iter.Seq[Buyurtma]` — slice qaytarmaydi, `for range` uning ustida aylanadi |
| `buzuq/` | xuddi shu, lekin `yield` natijasi tekshirilmaydi → `break` dan keyin **panic** (`go vet` jim) |
| `eski/` | `hook/` nusxasi, `go.mod` da `go 1.22`: `requires go1.23 or later` |
| `iz/` | iterator va sikl tanasi navbatma-navbat ishlashini print bilan koʻrsatadi |
| `seq2/` | `iter.Seq2`: `slices.All`; eski `sync.Map.Range` ham Seq2 shaklida — toʻgʻridan-toʻgʻri `range` qilsa boʻladi |
| `birmarta/` | bir martalik iterator: `strings.SplitSeq` ikkinchi marta boshidan boshlamaydi |
| `daraxt/` | generic ikkilik qidiruv daraxti, rekursiv `All()` (spec'dagi `Tree.walk` shakli) |
| `tezlik/` | bir xil yigʻindi 4 usulda + `olch.sh` (7 marta, mediana) |
| `tartib/` | map tartibsiz; `maps.Keys` — `iter.Seq`, slice emas; `slices.Sorted(maps.Keys(m))` |
| `saralash/` | `int(b.Reyting - a.Reyting)` xatosi va `cmp.Compare` |
| `nan/` | `0/0 = NaN` qoʻlda yozilgan taqqoslashni buzadi; `cmp.Compare` + `cmp.Or` |
| `savol/` | pauzadagi savol: iterator ichidagi `defer` va `break` |
| `homework/` | uy vazifasi javoblari (testlar bilan) |

Tekshirish (shu papkadan):
```bash
go vet ./... && go test ./tezlik ./homework/...
```

## Videodagi buyruqlar
Hammasi `lessons/018-iterators` dan, go1.26.8, i5-12450H, 2026-11-26.
```bash
grep -n -e 'func Tolanganlar' -e 'range Tolanganlar' hook/main.go
go run ./hook                       # buyurtma 1 150000 / buyurtma 3 42000 / kichik summa, toʻxtaymiz
diff hook/main.go buzuq/main.go | tail -5
go vet ./buzuq && echo 'vet: jim'   # vet: jim
go run ./buzuq 2>&1 | head -4
#   panic: runtime error: range function continued iteration after function for loop body returned false
go run ./buzuq 2>&1 | sed -n '4p;6,7p'
#   goroutine 1 [running]:
#   main.main-range1(...)           ← sikl tanasidan yasalgan funksiya
gopls check buzuq/main.go           # buzuq/main.go:22:4-12: yield may be called again after returning false
grep '^go ' eski/go.mod             # go 1.22
cd eski && go build -o /dev/null .
#   cannot range over Tolanganlar(bb) (value of func type iter.Seq[Buyurtma]): requires go1.23 or later
#   (-lang was set to go1.22; check go.mod)
cd eski && GOTOOLCHAIN=go1.22.12 GOEXPERIMENT=rangefunc go run .   # ishlaydi (tajribaviy)
go run ./iz
go run ./daraxt                     # 7000 42000 99000 150000
./tezlik/olch.sh                    # videodagi diagramma: 6 ta oʻlchov medianasi 246 · 270 · 265 · 2145 ns/op (Seq: 2 allocs/op)
go build -gcflags=-m ./tezlik 2>&1 | grep -E 'inlining call to (Barchasi|JamiOz.Barchasi.func1|JamiOz-range1)|moved to heap: j|tezlik.go:53:2: func literal'
go run ./seq2                       # 0 150000 / 1 42000 / Aziz 150000
go run ./birmarta                   # 1-marta: Aziz, Bobur, Dilnoz / 2-marta: Dilnoz
go run ./tartib                     # map tartibi tasodifiy — har ishga tushirishda boshqacha boʻlishi mumkin
go run ./saralash
go run ./nan
go run ./savol                      # fayl ochildi / birinchi / fayl yopildi / sikldan keyin
```
Benchmark (mediana, 7 marta, `-cpu 1`, 1000 ta `int64`): yengil yuklamada Oddiy 232–278, Values 230–297, Oz 247–273,
Seq 1897–2348 ns/op; mashina band boʻlganda (load ≈ 12–18) hammasi koʻtariladi, lekin nisbat saqlanadi: dastlabki uchtasi
bir-biriga yaqin (farq ≈ 10 %), Seq 7–9 baravar sekin. Oddiy/Values/Oz — `0 allocs/op`, Seq — `2 allocs/op` (40 B): bu bitta yigʻindiga (1000 element) ikkita obyekt — `j` va sikl tanasi closure'i (`-memprofilerate 1` bilan tekshirilgan). Asosiy narx — 1000 ta haqiqiy chaqiruv.

## Dars matni

**Hook.** `Tolanganlar` slice emas, funksiya qaytaradi, `for range` esa uning ustida aylanadi. Iterator ichidagi
`if !yield(b) { return }` ni oddiy `yield(b)` ga almashtirsak, `go vet` jim, lekin `break` dan keyin:
`panic: runtime error: range function continued iteration after function for loop body returned false`.

**Nega kerak?** Go 1.23 gacha har paket oʻz usulida aylanardi: `sync.Map.Range(func)`, `filepath.Walk(func)`,
`bufio.Scanner` uchun `for sc.Scan() { sc.Text() }`. [Go blog](https://go.dev/blog/range-functions) (I. L. Taylor, 2024):
*"when you start using a new Go container package you may have to learn a new looping mechanism."*
[Go 1.22](https://go.dev/doc/go1.22) — `GOEXPERIMENT=rangefunc` bilan tajriba; [Go 1.23](https://go.dev/doc/go1.23#language) —
tilning bir qismi: `range` `func(func() bool)`, `func(func(K) bool)`, `func(func(K, V) bool)` turlarini qabul qiladi.
Til versiyasi `go.mod` dagi `go` qatoridan olinadi (GO-1.16), shuning uchun `go 1.22` modulni go1.26 ham rad etadi.

**`iter.Seq`.** [`iter`](https://pkg.go.dev/iter): `type Seq[V any] func(yield func(V) bool)` — generic tur (GO-1.17).
Hujjatga koʻra "Seq" *seek* deb oʻqiladi. Iterator — `yield` ni har element bilan chaqiradigan funksiya; `yield` `false`
qaytarsa, iterator darhol qaytishi kerak. `Tolanganlar` slice'ni closure sifatida ushlab qoladi (GO-1.06).
`slices.Values` ham aynan shunday yozilgan ([slices/iter.go](https://github.com/golang/go/blob/go1.26.8/src/slices/iter.go#L37-L45)).

**Kim qachon ishlaydi (`iz/`).** Iterator qiymat beradi → sikl tanasi ishlaydi → `yield` `true` qaytaradi → …
`break` boʻlsa `yield` `false` qaytaradi, iterator toʻxtaydi, qolgan qiymatlar umuman hisoblanmaydi. Hammasi bitta
goroutine'da.

**Kompilyator nima qiladi.** [`rangefunc/rewrite.go`](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/rangefunc/rewrite.go)
(*Theory of Operation*): `for x := range f { … }` → `f(func(x T) bool { …; return true })`. Tanadagi `continue` →
`return true`, `break` → `return false`, `return` → `#next = -1; return false`, keyin `f` qaytgach tashqi funksiya ham
qaytadi. Stack trace'dagi `main.main-range1` — aynan shu yasalgan funksiya.

**Panic qoidasi.** [Spec](https://go.dev/ref/spec#For_range): *"If the loop body terminates (such as by a break statement),
yield returns false and must not be called again."* Har sikl uchun kompilyator `#stateK` holat oʻzgaruvchisini yaratadi;
yasalgan `yield` uni tekshiradi va qoidabuzarlikda panic qiladi. Runtime'dagi toʻrtta xabar
([runtime/panic.go L371–374](https://github.com/golang/go/blob/go1.26.8/src/runtime/panic.go#L371-L374)):
`…after function for loop body returned false` (bizning holat), `…after loop body panic`, `…after whole loop exit`,
`range function recovered a loop body panic and did not resume panicking`. `go vet` bu xatoni tekshirmaydi (go1.26.8), lekin muharrirdagi **gopls** (v0.22.0, `yield` analizatori) koʻrsatadi: `yield may be called again after returning false`.

**Daraxt.** `aylan` bitta qator: `d.chap.aylan(yield) && yield(d.qiymat) && d.ong.aylan(yield)` — `&&` tufayli `break`
dan keyin rekursiya darhol toʻxtaydi. Slice, channel yoki goroutine kerak emas. Spec'dagi `Tree.walk` misoli ham shu gʻoyada (u yerda `t == nil ||` va ikki qiymat — Seq2 shakli).

**Tezlik.** Kompilyator oddiy iteratorni ham, sikl tanasini ham inline qiladi (`inlining call to Barchasi`,
`inlining call to JamiOz-range1`) — oddiy sikl bilan bir xil tezlik, `0 allocs/op`. Iterator qiymat sifatida inline
qilinmaydigan funksiyaga berilsa (`JamiSeq`, `//go:noinline`): `moved to heap: j`, `func literal escapes to heap`, har
element — haqiqiy chaqiruv; bizning juda kichik tanada 7–9 baravar sekin. [Go wiki](https://go.dev/wiki/RangefuncExperiment):
*"For more complex bodies or iterators, the overhead of the function calls is insignificant."* Bu — standart (gc)
kompilyator va go1.26.8 xulqi; wiki: *"In any given release the compiler may or may not implement this series of optimizations."*

**`iter.Seq2`.** `range` uch xil funksiyani qabul qiladi: `func(func() bool)`, `func(func(V) bool)`, `func(func(K, V) bool)`.
Ikki qiymatlisi — `iter.Seq2[K, V]`: `slices.All` (indeks, qiymat), `maps.All` (kalit, qiymat). `sync.Map.Range` (Go 1.9 dan
beri) aynan shu shaklda: `var _ iter.Seq2[any, any] = m.Range` kompilyatsiya boʻladi, `for k, v := range m.Range` ishlaydi.

**Bir martalik iterator'lar.** [iter hujjati](https://pkg.go.dev/iter#hdr-Single_Use_Iterators): baʼzi iterator'lar ketma-ketlikni faqat
bir marta beradi. Go 1.24 dagi `strings.SplitSeq`, `strings.Lines` (va `bytes` dagilari) hujjatida: *"It returns a
single-use iterator."* Sababi — closure oʻzi ushlagan `s` satrini qisqartirib boradi (GO-1.06). `birmarta/`: ikkinchi `range`
faqat `Dilnoz` ni beradi. Qayta kerak boʻlsa — `slices.Collect` bilan slice'ga yigʻing yoki iterator'ni qayta yarating.

**Xatolar.** `iter.Seq[T]` da xato uchun joy yoʻq. Ikki odatiy yoʻl: `iter.Seq2[T, error]` yoki `bufio.Scanner` kabi
alohida `Err()` metodi.

**`maps` + `slices`.** Spec: *"The iteration order over maps is not specified and is not guaranteed to be the same from
one iteration to the next."* `maps.Keys(m)` — `iter.Seq[K]`, slice emas (`%T` → `iter.Seq[string]`).
Barqaror tartib: `slices.Sorted(maps.Keys(m))` (iter hujjatidagi misol). Slice kerak boʻlsa — `slices.Collect`.

**`cmp`.** [`slices.SortFunc`](https://pkg.go.dev/slices#SortFunc): `cmp(a, b)` *"should return a negative number when
a < b, a positive number when a > b and zero when a == b"*, va u *strict weak ordering* boʻlishi shart; saralash stabil
emas (`SortStableFunc` bor).
1. `int(b.Reyting - a.Reyting)`: `4.8 − 4.2 = 0.6` → `int` da `0` (spec: kasr qismi tashlanadi) → "teng". Barcha
   ayirmalar < 1 boʻlsa, hech narsa surilmaydi. Butun sonlarda esa ayirma toʻlib ketishi (overflow) mumkin.
2. `0/0 = NaN`: `<` va `>` ikkalasi `false` → qoʻlda yozilgan funksiya NaN'ni hammaga teng deb biladi, bu strict weak
   ordering emas → natija buziladi.
3. Toʻgʻrisi: [`cmp.Compare`](https://pkg.go.dev/cmp#Compare) — *"a NaN is considered less than any non-NaN, a NaN is
   considered equal to a NaN"*. Bir nechta kalit: [`cmp.Or`](https://pkg.go.dev/cmp#Or) (Go 1.22) — birinchi nol boʻlmagan
   natija: `cmp.Or(cmp.Compare(b.Reyting, a.Reyting), strings.Compare(a.Nom, b.Nom))`.

**`defer` va `break` (savol).** `break` → `yield` `false` → iterator `return` qiladi → uning `defer` i ishlaydi → shundan
keyin `range` tugaydi: `fayl ochildi / birinchi / fayl yopildi / sikldan keyin`. Bu kafolat `range` bilan ishlaydi: `iter.Pull` da
`stop` chaqirilmasa, iterator oʻrtada qoladi va `defer` i ishlamaydi (`homework/guru` `TestStop`). Wiki: *"If an iterator function defers a
call, the call runs when the iterator function returns."* Sikl tanasidagi `defer` esa tashqi funksiya qaytganda ishlaydi.

## Uy vazifasi
- 🟢 **Oson:** `Oraliq(a, b int) iter.Seq[int]` — `a` dan `b` gacha (ikkalasi ham kiradi). `break` bilan test qiling.
- 🟡 **Oʻrta:** generic adapterlar: `Filtr[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V]` va
  `Aylantir[V, W any](seq iter.Seq[V], f func(V) W) iter.Seq[W]`. `slices.Values` bilan zanjir qiling.
- 🟠 **Qiyin:** buyurtmalardan mijozlar boʻyicha jami summani hisoblang va mijozlarni jami summa kamayishi boʻyicha,
  teng boʻlsa ism boʻyicha qaytaring (`maps.Keys`, `slices.SortedFunc`, `cmp.Or`). Natija har safar bir xil boʻlsinmi?
- 🔴 **Guru:** `Zip[A, B any](a iter.Seq[A], b iter.Seq[B]) iter.Seq2[A, B]` — ikki iterator'ni yonma-yon yurgizing.
  Nega `range` bunga yetmaydi va nega `stop` ni `defer` qilish kerak?

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da (`go test ./homework/...`):
- Oson ([homework/oson](homework/oson/)): oddiy `for` + `if !yield(i) { return }`. `TestBreak` — `break` dan keyin panic yoʻq.
- Oʻrta ([homework/orta](homework/orta/)): ichki `for v := range seq` — tashqi `yield` `false` qaytarganda `return`
  qilinadi, ichki iterator esa oʻz-oʻzidan toʻxtaydi (uning `yield` i ham `false` qaytaradi).
- Qiyin ([homework/qiyin](homework/qiyin/)): `slices.SortedFunc(maps.Keys(jami), func(a, b string) int { return
  cmp.Or(cmp.Compare(jami[b], jami[a]), strings.Compare(a, b)) })`. Map tartibi tasodifiy, lekin taqqoslash toʻliq
  tartib beradi — test 50 marta tekshiradi.
- Guru ([homework/guru](homework/guru/)): `range` bitta iterator'ni boshqaradi (push). Ikkitasini navbatma-navbat
  tortish uchun `iter.Pull` ularni pull iterator'ga (`next`, `stop`) aylantiradi. Hujjat: *"If clients do not consume
  the sequence to completion, they must call stop, which allows the iterator function to finish and return."*
  `TestStop`: `stop` chaqirilmasa, iterator'ning `defer` i ishlamaydi.

## Manbalar
- Go spec, For statements with range clause: https://go.dev/ref/spec#For_range
- `iter` paketi: https://pkg.go.dev/iter · [iter.go](https://github.com/golang/go/blob/go1.26.8/src/iter/iter.go)
- Go 1.23 / 1.22 release notes: https://go.dev/doc/go1.23#language · https://go.dev/doc/go1.22
- I. L. Taylor, *Range Over Function Types* (2024-08-20): https://go.dev/blog/range-functions
- Go wiki, *Rangefunc Experiment*: https://go.dev/wiki/RangefuncExperiment
- Proposal #61405 (R. Cox, 2023): https://github.com/golang/go/issues/61405 · muhokama #56413 (2022): https://github.com/golang/go/discussions/56413
- Kompilyator: https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/rangefunc/rewrite.go
- `cmp`: https://pkg.go.dev/cmp · `slices`: https://pkg.go.dev/slices · `maps`: https://pkg.go.dev/maps
