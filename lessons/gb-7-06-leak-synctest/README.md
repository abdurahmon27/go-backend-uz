# GB-7.06 · Concurrent kodni test qilish: leak, vaqt va `synctest`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #56"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. `testing/synctest` Go 1.25'da standart
> kutubxonaga kirgan, `goroutineleak` profili Go 1.27'da (1.26'da faqat `GOEXPERIMENT=goroutineleakprofile` bilan).

Oldingi dars: `gb-7-05` (mutex yoki channel: `OwnedStore`) · keyingi dars: `gb-7-07` (nazorat nuqtasi: bekor qilinadigan
importer). Bu darsga kerak boʻlganlar: [`gb-7-03-select-context`](../gb-7-03-select-context) (`select`, `ctx.Done()`,
`context.WithCancel`, `runtime.NumGoroutine`, buffer'siz channel'dagi leak), [`gb-7-02-channel`](../gb-7-02-channel)
(buffer'siz yuborish qabul qiluvchini kutadi, `close`, `range`), [`gb-7-01-goroutine-hayoti`](../gb-7-01-goroutine-hayoti)
(`go func(){…}()`, har goroutine'ning egasi), [`gb-1-08-birinchi-test`](../gb-1-08-birinchi-test) (`go test`, `-v`, `t.Errorf`),
[`gb-5-05-pgx-pool`](../gb-5-05-pgx-pool) (`t.Context()`, `context.WithTimeout`), [`gb-4-07-mutex-race`](../gb-4-07-mutex-race) (`-race`).

**Bu darsdan keyin:** goroutine leak'ni (abadiy kutib qolgan goroutine) `runtime.NumGoroutine` bilan oʻlchaysiz va
`runtime/pprof`ning goroutine profilidan qaysi qatorda turganini topasiz; oddiy test uni nega koʻrmasligini bilasiz;
`testing/synctest` bilan testni pufak (bubble) ichida ishga tushirib, leak'ni test yiqilishiga aylantirasiz; soxta soat
bilan vaqtni aniq tekshirasiz (`50ms` aniq `50ms`); buffer nega yetmasligini va `context.WithCancel` + `defer cancel()`
goroutine'larni chaqiruv bilan birga qanday tugatishini koʻrasiz; unutilgan `Close`ni ham shu usulda ushlaysiz.

Bu papka mustaqil (Modul 7 laboratoriyasi): oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz
`package main`'i, tashqi modul yoʻq. bozor'dagi ehtiyoj: mahsulot narxini uch yetkazib beruvchidan (chorsu, oloy, siyob)
birdan soʻrab, eng tez javobni olish. Tarmoq oʻrniga kutish: javob 100, 200 va 300 ms'da keladi. Yangi nomlar:
`prices` (ID → tiyin: olma 1 = 1 200 000), `supplierDelays`, `quote(ctx, delay, id)`, `fastestPrice(ctx, id)`; vazifa uchun
`priceOffer`, `priceOffers`, `askOffer`, `quoteTimeout`, `errNoOffer`, `cheapestPrice`. `store.go` — GB-7.05'dagi
`OwnedStore`ning qisqa nusxasi (faqat `NewOwnedStore`, `Sell`, `Close`), xato sahnasi uchun.

## Papka

| Fayl | Nima |
|---|---|
| `price.go` | videoning oxiridagi `quote` va `fastestPrice` (ctx + cancel, leak yoʻq) |
| `store.go` | `OwnedStore` (GB-7.05'dan, qisqartirilgan) |
| `main.go` | 3 × 100 soʻrov va tirik goroutine'lar soni |
| `price_test.go` | videodagi testlar: `TestFastestPrice`, `TestFastestPriceTimeout`, `TestOwnedStore` (hammasi `synctest.Test`) |
| `offer.go` | vazifa: `priceOffers`, `askOffer` tayyor, `cheapestPrice` boʻsh |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob (xuddi shu fayllar, `cheapestPrice` yozilgan) |
| `_video/NN-…/` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 04-synctest go test` |
| `_video/race.sh` | `-race` buyruqlari uchun (Linux'da gcc kerak; quyida): `./_video/race.sh 07-timeout go test -race .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64)

Stek yoʻllari `…/bozor/` deb qisqartirilgan (run.sh vaqtinchalik papkada ishlaydi); qolgani aynan.

**01-leak** — `fastestPrice(id)`: buffer'siz `ch`, har yetkazib beruvchiga goroutine, `return <-ch`. `main` 3 marta 100 soʻrov
qiladi, bir soniya kutib, `runtime.NumGoroutine()`ni chiqaradi:

```
$ go run .
goroutine'lar: 201
goroutine'lar: 401
goroutine'lar: 601
```

Har soʻrovda eng tez javob qabul qilinadi, qolgan ikki goroutine `ch <- …`da abadiy kutadi.

**02-profil** — `main` oxirida `pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)` (bir xil steklar guruhlangan matn):

```
$ go run .
…
goroutine'lar: 601
goroutine profile: total 601
600 @ 0x47d88a 0x41421c 0x413e17 0x4ded9c 0x4835c1
#	0x4ded9b	main.fastestPrice.func1+0x5b	…/bozor/price.go:24

1 @ 0x440e11 0x47cbfd 0x4ccab1 0x4cc785 0x4c96e9 0x4debfd 0x44aa27 0x4835c1
#	0x4ccab0	runtime/pprof.writeRuntimeProfile+0xb0	$GOROOT/src/runtime/pprof/pprof.go:848
…
#	0x4debfc	main.main+0x13c				…/bozor/main.go:19
```

`price.go:24` — `ch <- quote(d, id)`. Ishlayotgan serverda xuddi shu matn `net/http/pprof` bilan
`/debug/pprof/goroutine?debug=1` manzilida.

**03-test** — oddiy test (synctest'siz): leak bor, lekin test oʻtadi.

```
$ go test -v
=== RUN   TestFastestPrice
--- PASS: TestFastestPrice (0.10s)
PASS
ok  	github.com/abdurahmon27/bozor	0.103s
```

**04-synctest** — xuddi shu test `synctest.Test(t, func(t *testing.T) { … })` ichida:

```
$ go test
--- FAIL: TestFastestPrice (0.00s)
panic: deadlock: main bubble goroutine has exited but blocked goroutines remain [recovered, repanicked]

goroutine 7 [running]:
…
goroutine 11 [sleep (durable), synctest bubble 1]:
time.Sleep(0xbebc200)
	$GOROOT/src/runtime/time.go:366 +0x145
github.com/abdurahmon27/bozor.quote(...)
	…/bozor/price.go:16
github.com/abdurahmon27/bozor.fastestPrice.func1()
	…/bozor/price.go:24 +0x29
created by github.com/abdurahmon27/bozor.fastestPrice in goroutine 9
	…/bozor/price.go:23 +0x65

goroutine 12 [sleep (durable), synctest bubble 1]:
…
FAIL	github.com/abdurahmon27/bozor	0.004s
```

Pufakdagi soat test funksiyasi tugaganda toʻxtaydi, shuning uchun qolgan ikki goroutine `time.Sleep`da qotgan
(`0xbebc200` ns = 200 ms, `0x11e1a300` = 300 ms). Test 0.00s: soat soxta, 100 ms kutilmadi.

**05-buffer** — `ch := make(chan int64, len(supplierDelays))`: yuborish endi kutmaydi, lekin test baribir yiqiladi
(xuddi shu `deadlock`, `goroutine 9 [sleep (durable), synctest bubble 1]`): goroutine'lar `fastestPrice`dan uzoqroq yashaydi.

**06-ctx** — `quote(ctx, …)` `select` bilan `time.After(delay)` yoki `ctx.Done()`ni kutadi; `fastestPrice(ctx, id)`
`context.WithCancel` + `defer cancel()` qiladi, goroutine faqat xatosiz narxni buffer'li `ch`ga yuboradi:

```
$ go test -v
=== RUN   TestFastestPrice
--- PASS: TestFastestPrice (0.00s)
PASS
ok  	github.com/abdurahmon27/bozor	0.002s
$ go run .
goroutine'lar: 1
goroutine'lar: 1
goroutine'lar: 1
```

**07-timeout** — `TestFastestPriceTimeout`: `context.WithTimeout(t.Context(), 50*time.Millisecond)`, eng tez javob 100 ms;
xato `context.DeadlineExceeded` va oʻtgan vaqt **aniq** `50ms` boʻlishi kerak:

```
$ go test -v
=== RUN   TestFastestPrice
--- PASS: TestFastestPrice (0.00s)
=== RUN   TestFastestPriceTimeout
--- PASS: TestFastestPriceTimeout (0.00s)
PASS
ok  	github.com/abdurahmon27/bozor	0.002s
```

**08-close** — `OwnedStore` testda yopilmadi (`defer s.Close()` yoʻq):

```
$ go test
--- FAIL: TestOwnedStore (0.00s)
panic: deadlock: main bubble goroutine has exited but blocked goroutines remain [recovered, repanicked]
…
goroutine 38 [chan receive (durable), synctest bubble 1]:
github.com/abdurahmon27/bozor.(*OwnedStore).loop(...)
	…/bozor/store.go:24
created by github.com/abdurahmon27/bozor.NewOwnedStore in goroutine 37
	…/bozor/store.go:19 +0xc7
…
FAIL	github.com/abdurahmon27/bozor	0.004s
```

**09-close-fix** — `defer s.Close()` bilan: `ok  	github.com/abdurahmon27/bozor	0.002s`.

**10-leakprofil** ("Chuqurroq") — `pprof.Lookup("goroutineleak")` (Go 1.27): avval leak qidiradigan maxsus GC, keyin faqat
leak boʻlgan goroutine'lar (`main` yoʻq):

```
$ go run .
…
goroutineleak profile: total 600
600 @ 0x47d88a 0x41421c 0x413e17 0x4ded9c 0x4835c1
#	0x4ded9b	main.fastestPrice.func1+0x5b	…/bozor/price.go:24
```

`WriteTo(os.Stdout, 2)` bilan har goroutine alohida, holati `[chan send (leaked)]` deb chiqadi.
Cheklovi ([Go 1.26 release notes](https://go.dev/doc/go1.26), "goroutine leak profile"): leak deb faqat hech bir ishlayotgan goroutine'dan yetib boʻlmaydigan channel/qulfda
kutayotganlar belgilanadi; channel global oʻzgaruvchida yoki tirik goroutine qoʻlida boʻlsa, profil uni koʻrmaydi.

**12-xotira** (videoda yoʻq) — leak qancha xotira oladi (`runtime.GC()` dan keyin `MemStats`, kechikishlar 1/2/3 ms):

```
$ go run .
1 goroutine; StackInuse KB: 512 HeapInuse KB: 536
2001 goroutine; StackInuse KB: 4544 HeapInuse KB: 2336
22001 goroutine; StackInuse KB: 44768 HeapInuse KB: 18008
```

**13-mutex** (videodagi "synctest chegarasi") — goroutine `sync.Mutex`da abadiy kutadi: bu durable emas, synctest deadlock
demaydi, test `-timeout`gacha osilib turadi:

```
$ go test -timeout 5s
panic: test timed out after 5s
…
```

**11-soat** (videodagi savol) — pufakda `time.Sleep(time.Hour)`, boshqa goroutine yoʻq:

```
$ go test -v
=== RUN   TestHour
    price_test.go:14: 1h0m0s 2000-01-01 00:00:00 +0000 UTC
--- PASS: TestHour (0.00s)
```

## `-race`

Linux'da race detector uchun `cgo` va C kompilyator (`gcc`) kerak (batafsil: [`gb-4-07-mutex-race`](../gb-4-07-mutex-race/README.md)).
Bu mashinada gcc yoʻq (`go: -race requires cgo; enable cgo by setting CGO_ENABLED=1`), shuning uchun tekshiruvlar Docker'da
qilingan (`_video/race.sh`; gcc bor `golang:1.26.8` image'i, Go oʻzi shu mashinadagi go1.27.4):

```
$ ./_video/race.sh 07-timeout go test -race .
ok  	github.com/abdurahmon27/bozor	1.013s
$ ./_video/race.sh solution go test -race -count=3 .
ok  	github.com/abdurahmon27/bozor	1.020s
```

DATA RACE hisoboti yoʻq: natija goroutine'dan channel orqali keladi (yuborish qabuldan oldin sodir boʻladi, Go memory model).

## Vazifa

[TASK.md](TASK.md): `cheapestPrice(ctx)` — hamma yetkazib beruvchidan birdan soʻrab, `quoteTimeout` (1 s) ichida kelgan eng
arzon taklifni qaytarish, hech kim javob bermasa `errNoOffer`, chaqiruvchi ctx'i tugasa `ctx.Err()`, va ortda birorta
goroutine qoldirmaslik. Tekshirish:

```bash
go test .
```

Javob: [`solution/`](solution). Oʻzgarishni koʻrish: `diff -u offer.go solution/offer.go`.

Papkaning oʻz-oʻzini tekshiruvi (shu mashinada):

```
$ gofmt -l . && go vet ./...          # chiqish yoʻq
$ go test . ; echo "starter exit $?"   # FAIL (TASK.md'dagi xabarlar), exit 1
$ go test ./solution/...
ok  	github.com/abdurahmon27/bozor/solution	0.003s
$ t=$(mktemp -d) && cp -r . "$t" && cp solution/*.go "$t"/ && (cd "$t" && go test ./...)
ok  	github.com/abdurahmon27/bozor	0.004s
ok  	github.com/abdurahmon27/bozor/solution	0.003s
```

## Manbalar

- [`testing/synctest`](https://pkg.go.dev/testing/synctest): pufak (bubble), `Test` ("waits for all goroutines in the bubble
  to exit", deadlock'da test yiqiladi), `Wait`, "durably blocked" roʻyxati (bubble'da yasalgan channel, `time.Sleep`,
  `sync.Cond.Wait`, `WaitGroup.Wait`; mutex, I/O va syscall durable emas), soxta soat (2000-01-01 00:00 UTC, faqat hamma
  goroutine durable bloklanganda suriladi, root goroutine chiqqanda toʻxtaydi), `T.Context`
- Go blog: [Testing concurrent code with testing/synctest](https://go.dev/blog/synctest) (2025-02, Go 1.24'dagi eksperimental API —
  `synctest.Run`; Go 1.25'dan beri `synctest.Test`, gʻoya oʻsha)
- Release notes: [Go 1.25 — testing/synctest](https://go.dev/doc/go1.25#new-testingsynctest-package),
  [Go 1.26 — experimental goroutine leak profile](https://go.dev/doc/go1.26), [Go 1.27 — goroutine leak profile GA](https://go.dev/doc/go1.27)
- [`runtime/pprof`](https://pkg.go.dev/runtime/pprof): `Lookup`, `Profile.WriteTo` (`debug=1` matn, `debug=2` panic formatidagi
  steklar), `goroutine` va `goroutineleak` profillari; [`net/http/pprof`](https://pkg.go.dev/net/http/pprof)
- [`runtime.NumGoroutine`](https://pkg.go.dev/runtime#NumGoroutine), [`context`](https://pkg.go.dev/context) (`WithCancel`,
  `Done`, `Err`), [Go memory model: channel communication](https://go.dev/ref/mem#chan)
- Go spec: [Send statements](https://go.dev/ref/spec#Send_statements) (buffer'siz yuborish qabul qiluvchi tayyor boʻlguncha
  kutadi), [For range](https://go.dev/ref/spec#For_range) (channel yopilguncha)
- Go blog: [Pipelines and cancellation](https://go.dev/blog/pipelines) ("goroutines are not garbage collected; they must exit on their own")
- [go.uber.org/goleak](https://pkg.go.dev/go.uber.org/goleak): pufaksiz testlar uchun (haqiqiy tarmoq) test oxirida ortiqcha goroutine tekshiruvi
