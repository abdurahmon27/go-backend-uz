# GB-2.02 · Map: qidirish, comma ok va kafolatsiz tartib

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #10"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-2-01-slice`](../gb-2-01-slice) (slice: `len`, `cap`, `append` va umumiy massiv) · keyingi dars:
[`gb-2-03-struct-turlar`](../gb-2-03-struct-turlar) (struct: `Product` va `Order` turlari).

**Bu darsdan keyin:** mahsulotni nomi boʻyicha `map`'dan topasiz, katalogda yoʻq mahsulotni comma ok
(`price, ok := prices[name]`) bilan rad etasiz, buyurtmadagi donalarni `count[name]++` bilan sanaysiz va natijani
kalitlarni saralab har safar bir xil tartibda chiqarasiz. Vositalar: `map[string]int64`, map literal, `m[k]`,
`len`, `delete`, `make`, map ustida `range`, `slices.Sort`.

Bu papka Module 2'ning mustaqil laboratoriyasi: oʻz `main.go`'si bor, GB-1.08'ning `parsePrice` kodidan davom
etmaydi (bozor'ning hammasi GB-3.06 "bozor CLI v0" darsida bitta dasturga yigʻiladi). Nomlar kursdagidek: narx
`int64`, tiyinda, `tiyinPerSom = 100`; mahsulotlar `olma`, `non`, `choy`; ombor `stock`.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi bozor: `go run . olma non olma choy` buyurtmani `prices` map'idan comma ok bilan tekshiradi, donalarni sanaydi, chekni alifbo tartibida chiqaradi; `reserve` hozircha boʻsh (vazifa) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa: `reserve(stock, order) error` — buyurtmani ombordan ayirish, xato boʻlsa omborga tegmaslik |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/` | videodagi har bir qadamning `main.go` fayli; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi (`06-nil` ataylab panic beradi) |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 04-tartib go run . olma non olma choy` |
| `_video/tartib1000.sh` | `04-tartib`ni 1000 marta ishga tushirib, qaysi tartib necha marta chiqqanini sanaydi |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Hook: comma ok'siz bozor (`./_video/run.sh 00-hook go run . olma anor`). Katalogda `anor` yoʻq, `prices["anor"]`
esa `0` qaytaradi, anor tekinga ketadi:

```
anor  1 dona
olma  1 dona
jami: 12000 soʻm
```

Map asoslari (`_video/01-map`, `go run .`): `prices["olma"]`, `choy` qoʻshilgandan keyin `len(prices)`,
yoʻq kalit `prices["anor"]`:

```
1200000
3
0
```

Comma ok (`_video/02-ok`, soʻzlar soni GB-1.07'dagidek tekshiriladi): `go run . olma` → `olma 12000 soʻm`;
`go run .` → `foydalanish: bozor NOM` / `exit status 1`; `go run . anor`:

```
xato: anor katalogda yoʻq
exit status 1
```

Nol qiymat va `delete` (`_video/03-delete`, `go run .`): `prices["olma"] = 0` dan keyin kalit joyida (`ok` true),
`delete(prices, "non")` dan keyin yoʻq (`ok` false):

```
2 true
1 false
```

Yoʻq kalitni `delete` qilish xato emas, hech narsa qilmaydi (Go spec: *"If the map m is nil or the element m[k]
does not exist, delete is a no-op"*).

Map ustida `range` (`_video/04-tartib`, `go run . olma non olma choy`), bir xil buyruq, ikki xil natija:

```
olma 2        choy 1
non 1         olma 2
choy 1        non 1
```

1000 marta (`./_video/tartib1000.sh`), ikki oʻlchov:

| tartib | 1-oʻlchov (videoda) | 2-oʻlchov |
|---|---|---|
| `olma non choy` | 743 | 747 |
| `choy olma non` | 137 | 118 |
| `non choy olma` | 120 | 135 |

Sababi (Go spec): *"The iteration order over maps is not specified and is not guaranteed to be the same from one
iteration to the next."* Runtime har `range` boshida tasodifiy boshlanish joyini tanlaydi
(`internal/runtime/maps/table.go`, `Iter.Init`: `it.entryOffset = rand()`). Sakkiztagacha kalitli kichik map'da
kalitlar 8 oʻrinli bitta guruhda qoʻshilgan tartibda turadi; boshlanish 8 oʻrindan biri, 3 ta kalit bor, boʻsh
5 oʻrin va `olma`ning oʻrni (jami 6/8 = 75 %) `olma`dan boshlanadigan tartibni beradi. Bu runtime tafsiloti: unga
tayanmang.

Saralangan chek (`_video/05-sort`), har safar bir xil. `slices.Sort` matnlarni bayt boʻyicha solishtiradi: lotin
harflari alifbo tartibida, lekin `oʻ`/`gʻ` (U+02BB) alifbodagi oʻrnida emas (`olma` < `oʻrik` < `p…`); oʻzbekcha
matn va baytlar GB-2.10'da. `%-5s` ustuni baytlarni emas, belgilarni sanaydi, lekin 5 belgidan uzun nom ustunni
buzadi: haqiqiy jadval chiqishi uchun keyinroq `text/tabwriter`. Go 1.23'dan beri xuddi shu ish bir qatorda
`slices.Sorted(maps.Keys(count))` bilan ham yoziladi; iterator'lar GB-3.05'da.

```
choy  1 dona
non   1 dona
olma  2 dona
```

Yigʻilgan bozor (`main.go`): `go run . olma anor` → `xato: anor katalogda yoʻq` / `exit status 1`;
`go run . olma non olma choy`:

```
choy  1 dona
non   1 dona
olma  2 dona
jami: 53000 soʻm
```

Tez-tez uchraydigan xato: `var stock map[string]int` — nil map (`_video/06-nil`, `go run .`). Oʻqish ishlaydi,
yozish panic beradi:

```
0 0
panic: assignment to entry in nil map

goroutine 1 [running]:
main.main()
	/tmp/…/bozor/main.go:8 +0xa7
exit status 2
```

Tuzatish (`_video/07-nil-tuzat`): `stock := make(map[string]int)` → `0 0` / `40`.

## Vazifa

[`TASK.md`](TASK.md): `reserve(stock map[string]int, order []string) error` — buyurtmani ombordan ayirish:
`anor omborda yoʻq`, `choy: omborda 1 dona, buyurtmada 2`, xato boʻlsa ombor oʻzgarmaydi. Tekshirish:

```bash
go test
```

Boshlangʻich kodda (`return nil`) test yiqiladi, masalan:

```
--- FAIL: TestVazifaAyiradi (0.00s)
    tekshiruv_test.go:47: reserve'dan keyin non qoldigʻi 1 boʻlishi kerak, hozir 2 (TASK.md, 1-qadam)
    tekshiruv_test.go:47: reserve'dan keyin olma qoldigʻi 3 boʻlishi kerak, hozir 5 (TASK.md, 1-qadam)
--- FAIL: TestVazifaRadEtadi (0.00s)
    --- FAIL: TestVazifaRadEtadi/yoq (0.00s)
        tekshiruv_test.go:69: reserve([olma anor]) xato qaytarishi kerak: "anor omborda yoʻq"; hozir xato yoʻq (TASK.md, 2-qadam)
…
```

Javob bilan tekshirish: `go test ./solution/...` → `ok`.

Papkaning oʻz tekshiruvi (GB-WRITER §4):

```bash
gofmt -l . && go vet ./...                 # hech narsa chiqmaydi
go test . ; echo "starter exit $?"         # FAIL, TASK.md xabarlari, exit 1
go test ./solution/...                     # ok
t=$(mktemp -d) && cp -r . "$t" && cp solution/*.go "$t"/ && (cd "$t" && go test ./...)   # hammasi ok
```

## Keyingi darslarda

Map haqida bu darsda yoʻq, lekin ishda kerak boʻladiganlar: xatoni chaqiruvchi `errors.Is` bilan tekshirishi
(GB-2.08), `maps`/`slices` iterator'lari (GB-3.05), bir vaqtda ishlaydigan soʻrovlar bitta map'ga yozganda
`sync.Mutex` (GB-4.07). Kalit turi `==` bilan solishtiriladigan boʻlishi shart (slice kalit boʻla olmaydi, Go spec
Map types). GB-2.03'da `prices` va `stock` alohida map'lar oʻrniga `Product` turi (narx va qoldiq uning maydonlari) keladi.

## Manbalar

- Go spec: [Map types](https://go.dev/ref/spec#Map_types) (nil map), [Index expressions](https://go.dev/ref/spec#Index_expressions)
  (yoʻq kalit → zero value; `v, ok = a[x]`), [Deletion of map elements](https://go.dev/ref/spec#Deletion_of_map_elements),
  [Making slices, maps and channels](https://go.dev/ref/spec#Making_slices_maps_and_channels),
  [For statements with range clause](https://go.dev/ref/spec#For_range) (tartib belgilanmagan),
  [Run-time panics](https://go.dev/ref/spec#Run_time_panics)
- [Go maps in action](https://go.dev/blog/maps) (go.dev/blog)
- [`slices.Sort`](https://pkg.go.dev/slices#Sort), [`cmp.Ordered`](https://pkg.go.dev/cmp#Ordered) (matnlar bayt boʻyicha)
- [Go 1.24 release notes: Swiss table map](https://go.dev/doc/go1.24#runtime)
- Runtime manbasi: [`internal/runtime/maps/table.go`](https://github.com/golang/go/blob/go1.27.4/src/internal/runtime/maps/table.go)
