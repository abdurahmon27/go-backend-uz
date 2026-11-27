# GB-3.05 · Iterator'lar: `range` funksiya ustida, `slices` va `maps`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #23"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** 2–3-modul darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. `range` funksiya ustida tilda, `iter`
> paketi standart kutubxonada Go 1.23'dan beri bor: `go.mod`'da `go 1.22` yoki pastroq boʻlsa, kompilyator
> funksiya ustidagi `range`'ni rad etadi.

Oldingi dars: [`gb-3-04-generics`](../gb-3-04-generics) (tur parametri `[T any]`) · keyingi dars: `gb-3-06`
(nazorat nuqtasi: bozor CLI v0). `Store`, `internal/store` va `List()` + `slices.Clone`:
[`gb-2-06-paketlar`](../gb-2-06-paketlar); map, `range` tartibi va `slices.Sort`:
[`gb-2-02-map-qidirish`](../gb-2-02-map-qidirish); nomsiz funksiya tashqi oʻzgaruvchini koʻrishi:
[`gb-2-09-defer-panic`](../gb-2-09-defer-panic).

**Bu darsdan keyin:** `Store`'ga `All() iter.Seq[Product]` metodini yozib, katalogni ichki slice'ni ochmasdan va
nusxalamasdan `for p := range st.All()` bilan aylanasiz; iterator va sikl tanasi navbatma-navbat ishlashini
(kompilyator tanani `yield` funksiyasiga aylantiradi) chiqishda koʻrasiz; `break` boʻlganda `yield` `false`
qaytarishini va iterator `if !yield(p) { return }` bilan toʻxtashi shartligini (aks holda `range function continued
iteration …` panic'i) bilasiz; `go doc iter.Seq`'ni oʻqiysiz; map kalitlarini `slices.Sorted(maps.Keys(m))` bilan
tartiblaysiz va `maps.Keys` nega slice emasligini (`len` xatosi) koʻrasiz.

Bu papka 2–3-modul qoidasi boʻyicha **mustaqil lab**: oʻz `go.mod`'i bor, tuzilishi GB-2.06'dagidek
`cmd/bozor` + `internal/store`. `Product`, `Store`, `Add`, `List` (GB-2.06, `slices.Clone` bilan) oʻsha darsdan;
yangi: `All`, `Names`, vazifaning `LowStock` va `SortedNames`'i. GB-2.06'dagi `Sell` va narxni tilga qarab chiqarish
(`golang.org/x/text`) bu labga kerak emas, shuning uchun olib tashlangan; ularni, shu iterator'lar bilan birga,
3-modul oxiridagi bozor CLI v0 (GB-3.06) birlashtiradi. GB-2.06 bilan farq: `diff -ru ../gb-2-06-paketlar/solution solution`.

## Papka

| Fayl | Nima |
|---|---|
| `cmd/bozor/main.go` | videoning oxiridagi `main`: katalogni `range st.All()` bilan chiqaradi, birinchi kam qolgan mahsulotda `break`, buyurtmadagi donalarni `slices.Sorted(maps.Keys(count))` tartibida chiqaradi |
| `internal/store/store.go` | `Product`, `Store`, `Add`, `List`, `All() iter.Seq[Product]`, `Names() iter.Seq[string]`; `LowStock` va `SortedNames` — vazifaning boshlanishi |
| `internal/store/tekshiruv_test.go` | vazifani tekshiradi (`go test ./...`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: xuddi shu papka, vazifa bajarilgan; `cd solution && go test ./...` |
| `_video/NN-…/` | videodagi har bir qadam (toʻliq `go.mod` + `cmd` + `internal`); `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi vaqtinchalik `bozor/` papkasida ishga tushiradi va oxirida `(chiqish kodi: N)` yozadi (quyida bu qator tushirilgan) |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

GB-2.06'dagi `List()` har chaqiruvda `slices.Clone(s.products)` qiladi (`00-list`):

```
$ ./_video/run.sh 00-list go run ./cmd/bozor
birinchi kam qolgan: non
```

`All()` slice emas, funksiya qaytaradi; `range` uning ustida aylanadi (`01-range`, imzo hali
`func(yield func(Product) bool)`):

```
$ ./_video/run.sh 01-range go run ./cmd/bozor
olma 40
non 3
choy 8
oʻrik 2
```

Iterator va sikl tanasi navbatma-navbat (`02-iz`: iterator `yield`'dan oldin, tana esa oʻzi chiqaradi):

```
$ ./_video/run.sh 02-iz go run ./cmd/bozor
iterator: olma beraman
sikl: olma
iterator: non beraman
sikl: non
iterator: choy beraman
sikl: choy
iterator: oʻrik beraman
sikl: oʻrik
```

`break` (`03-break`): `yield` `false` qaytaradi, iterator `return` qiladi, keyin `range` tugaydi:

```
$ ./_video/run.sh 03-break go run ./cmd/bozor
iterator: olma beraman
sikl: olma
iterator: non beraman
sikl: non
iterator: yield false → toʻxtayman
sikldan keyin
```

Tez-tez uchraydigan xato (`04-buzuq`): `yield` natijasi tekshirilmaydi. `go build` ham, `go vet` ham jim; ishlaganda
kompilyator yasagan `yield` oʻzi `false` qaytarganini eslab turadi va qayta chaqirilganda panic qiladi (stack
qatorlari videoda qisqartirilgan; vaqtinchalik papka yoʻli har safar boshqa):

```
$ ./_video/run.sh 04-buzuq go vet ./...
$ ./_video/run.sh 04-buzuq go run ./cmd/bozor
iterator: olma beraman
sikl: olma
iterator: non beraman
sikl: non
iterator: choy beraman
panic: runtime error: range function continued iteration after function for loop body returned false

goroutine 1 [running]:
main.main-range1(...)
	/tmp/…/bozor/cmd/bozor/main.go:16
github.com/abdurahmon27/bozor/internal/store.(*Store).All.func1(...)
	/tmp/…/bozor/internal/store/store.go:39
main.main()
	/tmp/…/bozor/cmd/bozor/main.go:16 +0x59a
exit status 2
```

(`main.main-range1` — kompilyator `main` ichidagi birinchi `range` siklining tanasidan yasagan funksiya. Runtime'da
yana ikki shunday xabar bor: `… after loop body panic` va `… after whole loop exit`, `src/runtime/panic.go`.
`gopls`'da `yield` tahlilchisi bor (`yield may be called again after returning false`), lekin u faqat `iter`'ni import
qilgan paketni tekshiradi; bu yerdagi `gopls` v0.22.0 esa `go 1.27` modulini tahlil qilmadi, `gopls check` jim qoldi.
Shuning uchun xatoga muharrirga emas, `if !yield(p) { return }` qoidasiga tayaning.)

`iter.Seq` (`05-seq`: `All() iter.Seq[Product]`, natija `01-range` bilan bir xil):

```
$ go doc iter.Seq
package iter // import "iter"

type Seq[V any] func(yield func(V) bool)
    Seq is an iterator over sequences of individual values. When called as
    seq(yield), seq calls yield(v) for each value v in the sequence, stopping
    early if yield returns false. See the iter package documentation for more
    details.
```

`slices.Values(s)` — standart kutubxonadagi aynan shu tana (`src/slices/iter.go`):

```go
func Values[Slice ~[]E, E any](s Slice) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}
```

Map kalitlarini tartiblash: GB-2.02 usuli (`06-eski`) va `slices.Sorted(maps.Keys(count))` (`08-sorted`), natija bir xil:

```
$ ./_video/run.sh 06-eski go run ./cmd/bozor
choy 1
non 2
olma 3
$ ./_video/run.sh 08-sorted go run ./cmd/bozor
choy 1
non 2
olma 3
```

`maps.Keys` slice emas, `iter.Seq` qaytaradi (`07-len`; birinchi qator videoda qisqartirilgan). Eski kodda va
qidiruv natijalarida `golang.org/x/exp/maps.Keys` uchraydi: u slice qaytarardi; standart `maps.Keys` — iterator:

```
$ ./_video/run.sh 07-len go run ./cmd/bozor
# github.com/abdurahmon27/bozor/cmd/bozor
cmd/bozor/main.go:18:18: invalid argument: maps.Keys(count) (value of func type iter.Seq[string]) for built-in len
$ go doc slices.Collect
func Collect[E any](seq iter.Seq[E]) []E
    Collect collects values from seq into a new slice and returns it. If seq is
    empty, the result is nil.
```

### Qoʻshimcha: iterator ichidagi `defer` (`10-defer`, videoda faqat aytiladi)

Iterator faylni ochadi va `defer` bilan yopadi; sikl birinchi qatordan keyin `break` qiladi. `yield` `false`
qaytaradi, iterator `return` qiladi, uning `defer`'i ishlaydi, shundan keyingina `range` tugaydi:

```
$ ./_video/run.sh 10-defer go run ./cmd/bozor
sikl: 1;olma;12000;40
iterator: fayl yopildi
sikldan keyin
```

(Demoda ochish xatosi tashlab yuborilgan. Xatoni ham beradigan iterator `iter.Seq2[string, error]` turida
yoziladi va `for line, err := range …` bilan aylanadi.)

### Oʻlchov: 10 000 mahsulot, birinchi kam qolgan mahsulot 2-oʻrinda (`09-olchov`)

```
$ ./_video/run.sh 09-olchov go test -run '^$' -bench . -benchmem -count 3 ./internal/store
cpu: 12th Gen Intel(R) Core(TM) i5-12450H
BenchmarkList-12    	   18474	     64991 ns/op	  401409 B/op	       1 allocs/op
BenchmarkList-12    	   20222	     60045 ns/op	  401408 B/op	       1 allocs/op
BenchmarkList-12    	   20527	     62590 ns/op	  401408 B/op	       1 allocs/op
BenchmarkAll-12     	174815608	         6.922 ns/op	       0 B/op	       0 allocs/op
BenchmarkAll-12     	209411606	         5.794 ns/op	       0 B/op	       0 allocs/op
BenchmarkAll-12     	204349036	         5.804 ns/op	       0 B/op	       0 allocs/op
```

Videodagi raqamlar — mediana: `List` ≈ 62.6 µs va 401 408 bayt nusxa (10 000 × 40 baytlik `Product`), `All` ≈ 5.8 ns,
0 bayt. Farq nusxadan va erta toʻxtashdan: `All` faqat ikki mahsulotga tegadi. Benchmark'lar GB-11'da ochiladi.

## Vazifa

[TASK.md](TASK.md): `LowStock(limit) iter.Seq[Product]` (qoldigʻi `limit`'dan kam mahsulotlar, `break`'da
toʻxtaydi) va `SortedNames() []string` (`slices.Sorted(s.Names())`). Tekshirish: `go test ./...` (hozir FAIL,
vazifa bajarilgach `ok`).

## Manbalar

- Go 1.23 release notes: range over function iterators, `iter` paketi — https://go.dev/doc/go1.23#language
- Go spec: For statements with range clause (iterator funksiyalar, `yield` `false`dan keyin chaqirilmaydi) —
  https://go.dev/ref/spec#For_range
- pkg.go.dev `iter` (`Seq`, `Seq2`, «conventionally named yield») — https://pkg.go.dev/iter
- pkg.go.dev `maps.Keys`, `slices.Sorted`, `slices.Collect`, `slices.Values` — https://pkg.go.dev/maps ·
  https://pkg.go.dev/slices
- Go blog: Range Over Function Types (I. L. Taylor, 2024) — https://go.dev/blog/range-functions
- Go wiki: Rangefunc Experiment (defer, panic, kompilyator tanani qanday oʻzgartiradi) —
  https://go.dev/wiki/RangefuncExperiment
- `src/runtime/panic.go` (`rangeDoneError`) — https://github.com/golang/go/blob/go1.27.4/src/runtime/panic.go
