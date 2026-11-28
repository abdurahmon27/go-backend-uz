# GB-7.05 · Mutex yoki channel: qarorni oʻlchab qabul qilamiz

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #55"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-7-04-worker-errgroup`](../gb-7-04-worker-errgroup) (worker pool, `errgroup`, `SetLimit`) · keyingi dars:
`gb-7-06` (concurrent kodni test qilish: leak, vaqt va `synctest`). Bu darsga kerak boʻlganlar:
[`gb-4-07-mutex-race`](../gb-4-07-mutex-race) (data race, `sync.Mutex`, `-race`, qulf nusxalanmaydi),
[`gb-7-03-select-context`](../gb-7-03-select-context) (`select`, vazifa uchun),
[`gb-7-02-channel`](../gb-7-02-channel) (channel, `close`, `range`, yopiq channel'dan oʻqish),
[`gb-7-01-goroutine-hayoti`](../gb-7-01-goroutine-hayoti) (`go`, `wg.Go`/`wg.Wait`, har goroutine'ning egasi va tugash yoʻli).

**Bu darsdan keyin:** bitta umumiy holatni (mahsulot qoldigʻi) ikki usulda himoyalay olasiz — `sync.Mutex` bilan va holatning
yagona egasi boʻlgan goroutine bilan (soʻrov va javob channel'lari); egasi goroutine'da qulfsiz ham data race yoʻqligini
`go run -race` bilan va nega yoʻqligini Go xotira modelidan bilasiz; ikkala usulning bitta amal narxini sakkiz goroutine
ostida oʻlchaysiz; `sync/atomic`dagi `atomic.Int64` bilan qulfsiz sanagich yozasiz va uning chegarasini koʻrasiz: ikki atomik
amal (`Load`, keyin `Add`) bitta qadam emas, qoldiq manfiy boʻladi, `-race` esa jim (race condition, data race emas).

Bu papka mustaqil (Modul 7 laboratoriyasi): oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i,
tashqi modul yoʻq. Qoldiq bu yerda soddalashtirilgan: katalog oʻrniga faqat `map[int64]int` (mahsulot ID → dona). Yangi nomlar
(keyingi darslar uchun): `LockedStore`, `OwnedStore`, `NewOwnedStore`, `loop`, `sellReq`, `stockReq`, `sells`, `stocks`,
`reply`, `Close`, `Requests`, `measure`, `sold`.

## Papka

| Fayl | Nima |
|---|---|
| `store.go` | `LockedStore` (map + `sync.Mutex`) va `OwnedStore` (egasi goroutine) — videoning oxiridagi kod + vazifa uchun boʻsh `Stock`, `Requests` |
| `main.go` | `measure`: bitta amalni 8 goroutine'dan 800 000 marta chaqirib, oʻrtacha vaqtini chiqaradi; mutex, egasi va `atomic.Add` |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `store.go` (+ xuddi shu `main.go` va test) |
| `_video/NN-…/` | videodagi har qadamning fayllari; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 01-ikki go run .` |
| `_video/race.sh` | `-race` uchun (Linux'da gcc kerak; bu mashinada yoʻq, Docker'da): `./_video/race.sh 01-ikki go run -race .` |
| `_video/06-bitta/` | videoda yoʻq: xuddi shu oʻlchov 1 va 8 goroutine bilan (quyidagi jadval) |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, 12 mantiqiy CPU)

Vaqtlar har ishga tushirishda boshqacha; quyida videodagi ishga tushirish, oraliqlar keyingi boʻlimda.

**01-ikki** — `LockedStore` va `OwnedStore`, har biri 8 goroutine × 100 000 `Sell(1, 1)`. Avval `-race` bilan (Docker,
`race.sh`): ogohlantirish yoʻq; vaqtlar `-race` ostida bir necha barobar sekin, ularga qaralmaydi:

```
$ go run -race .
mutex 784ns
egasi 1.634µs
```

Keyin oddiy `go run .` (`ns` — nanosekund, `µs` — mikrosekund = 1000 ns). Chiqadigan son — jami vaqt / 800 000, ya'ni bitta
chaqiruvga toʻgʻri keladigan vaqt (8 goroutine bir vaqtda ishlaydi, shuning uchun bu bitta chaqiruvning oʻz kechikishi emas,
oʻtkazuvchanlikning teskarisi):

```
$ go run .
mutex 59ns
egasi 519ns
```

**02-atomic** — `var sold atomic.Int64` va `measure("atomic", func() { sold.Add(1) })` (= papkadagi `main.go`):

```
$ go run .
mutex 52ns
egasi 465ns
atomic 20ns
sotildi: 800000
```

`-race` bilan ham ogohlantirish yoʻq (`./_video/race.sh 02-atomic go run -race .` → `mutex 519ns`, `egasi 1.527µs`,
`atomic 305ns`, `sotildi: 800000`).

**03-atomic-xato** — qoldiqni `atomic.Int64` bilan "sotish": `if stock.Load() > 0 { stock.Add(-1) }`, 100 000 dona, 800 000
urinish. Videoda faqat qoldiq qatori koʻrsatilgan (vaqt qatori qisqartirilgan: 100 000 dan keyin har urinish bitta `Load`):

```
$ go run .
atomic 3ns
qoldiq: -4
$ go run -race .              # Docker, race.sh
atomic 233ns
qoldiq: -6
```

Qoldiq har safar boshqa: 30 ta ishga tushirishda 29 marta manfiy (−1…−7), bir marta `0`. Xato tasodifiy — shuning uchun
xavfli: oddiy test uni koʻpincha koʻradi, lekin kafolat yoʻq. `-race` jim: har
murojaat atomik, **data race yoʻq**. Xato boshqa: `Load` va `Add` har biri boʻlinmas, ikkalasi birga emas — ikki goroutine
bir paytda `1`ni oʻqiydi va ikkalasi ham sotadi (**race condition**: natija goroutine'lar tartibiga bogʻliq). Uni `-race`
emas, test ushlaydi. Tekshirish va oʻzgartirish bitta qadam boʻlishi kerak boʻlsa: mutex, egasi goroutine yoki CAS (quyida).

## Chuqurroq (videoda yoʻq)

**04-cas** — `atomic` bilan toʻgʻri variant: `CompareAndSwap(old, new)` qiymat hali `old` boʻlsa `new` yozadi va `true`
qaytaradi; oraliqda boshqa goroutine oʻzgartirgan boʻlsa `false`, biz qaytadan oʻqiymiz:

```go
func sellOne() bool {
	for {
		old := stock.Load()
		if old <= 0 {
			return false
		}
		if stock.CompareAndSwap(old, old-1) {
			return true
		}
	}
}
```

```
$ go run .
cas 6ns
qoldiq: 0
$ go run -race .              # Docker
cas 245ns
qoldiq: 0
```

Bu faqat **bitta son** uchun ishlaydi. Ikki qiymatni birga oʻzgartirish kerak boʻlsa (qoldiq va buyurtmalar roʻyxati) —
yana mutex yoki egasi goroutine. `sync/atomic` hujjati ham shuni aytadi: "Except for special, low-level applications,
synchronization is better done with channels or the facilities of the sync package."

**05-topshirilgan** — egasi goroutine'ning klassik xatosi: map `NewOwnedStore`ga topshirilgandan keyin `main` uni oʻzi
oʻqiydi. Endi map'ga ikki goroutine tegadi va `-race` darhol koʻradi (bu `01-ikki`dagi jimlik `-race` ishlamagani emasligini
ham koʻrsatadi):

```
$ go run -race .              # Docker, qisqartirilgan
qoldiq: 1000
==================
WARNING: DATA RACE
Write at 0x00c0000a60f0 by goroutine 8:
  runtime.mapassign_fast64()
  …
  main.(*OwnedStore).loop()
      /home/dev/bozor/store.go:50 +0x133
…
Previous read at 0x00c0000a60f0 by main goroutine:
  runtime.mapaccess2_fast64()
  …
```

Qoida: egasiga topshirilgan ma'lumotga topshirgan kod boshqa tegmaydi; tashqariga qiymat kerak boʻlsa, egasi **nusxa**
qaytaradi (map'ning oʻzini emas).

## Oʻlchovlar: oraliqlar va nimani anglatadi

`measure` — 8 goroutine, har biri 100 000 chaqiruv, jami 800 000; vaqt `time.Since(start)/800_000` (bitta chaqiruvga toʻgʻri
keladigan vaqt). Bu benchmark emas (benchmark'lar `go test -bench` / `b.Loop` bilan GB-11.01'da): qizdirish yoʻq, bitta
mahsulot (bitta «issiq» kalit), har chaqiruv muvaffaqiyatli, egasi narxiga har chaqiruvdagi `make(chan bool)` ham kiradi. Lekin
bitta mashinada nisbatni koʻrsatishga yetadi.

| amal | 8 goroutine (`01-ikki`/`02-atomic`/`06-bitta`, 5–8 ishga tushirish) | 1 goroutine (`06-bitta`, 3 ishga tushirish) |
|---|---|---|
| `atomic.Int64.Add` | 14–22 ns | — |
| `LockedStore.Sell` (mutex) | 34–80 ns | 21–31 ns |
| `OwnedStore.Sell` (egasi) | 448–739 ns | 636–714 ns |

Nega egasi qimmat: har `Sell` ikki channel uzatishi (soʻrov va javob) va odatda ikki marta kutayotgan goroutine'ni uygʻotish
(`runtime/chan.go`: kutayotgan qabul qiluvchi topilsa, qiymat unga toʻgʻridan-toʻgʻri oʻtadi va u `goready` bilan
uygʻotiladi), ustiga har chaqiruvda yangi `make(chan bool)`. Boʻsh mutex'ni olish esa bitta `CompareAndSwap`
(`internal/sync/mutex.go`, `Lock`ning "fast path"i) — bu 1 goroutine qatorini tushuntiradi; 8 goroutine'da qulf ba'zan
band boʻladi va `Lock` sekinroq yoʻlga (`lockSlow`: qisqa aylanish, keyin kutish) oʻtadi, shuning uchun 21–31 ns emas, 34–80 ns.

Amalda: 500 ns ham soniyasiga ~2 million `Sell`. Handler'da JSON yoki baza bundan ancha uzoq boʻlsa, farq umumiy vaqtda
koʻrinmaydi. Tanlovni avvalo kodning soddaligi hal qiladi ([Go wiki: MutexOrChannel](https://go.dev/wiki/MutexOrChannel):
"Use whichever is most expressive and/or most simple"); oʻlchov tezlik qayerda muhimligini koʻrsatadi.

| bozor'da | nima | nega |
|---|---|---|
| qoldiq map'i (qisqa oʻqish/yozish) | `sync.Mutex` | holat; GB-4.08 katalog API ham shunday |
| narx roʻyxatlari bosqichdan bosqichga | channel | ma'lumot egaligi uzatiladi (GB-7.02 pipeline) |
| soʻrovlar / sotuvlar soni | `atomic.Int64` | bitta son, bitta amal |
| holat + tartibli ish (masalan har oʻzgarishni ketma-ket faylga yozish) | egasi goroutine | ish va holat bitta goroutine'da, tartib oʻz-oʻzidan |

## `-race`: Linux va macOS talablari

Linux'da race detector uchun `cgo` yoqilgan va C kompilyator (odatda `gcc`) oʻrnatilgan boʻlishi kerak; macOS'da Go 1.20'dan beri
Xcode shart emas ([Go 1.20 eslatmalari](https://go.dev/doc/go1.20#cgo): "can be built and run without Xcode"). Bu mashinada `gcc` yoʻq, shuning uchun `-race` natijalari Docker'da olingan (`_video/race.sh`: gcc bor
`golang:1.26.8` image'i, Go oʻzi shu mashinadagi `go1.27.4` toolchain, faqat oʻqish uchun ulanadi). Yechim `-race` bilan toza:
`./_video/race.sh solution go test -race -count=1 .` → `ok`. `-race` faqat ishga tushgan kodda sodir boʻlgan data race'ni
topadi; mantiq xatosini (03-atomic-xato) emas.

## Vazifa

[TASK.md](TASK.md): `OwnedStore`ga qulfsiz `Stock(id)` (ikkinchi soʻrov channel'i + `loop`da `select`, yopiq `sells`da
`return`) va ixtiyoriy `Requests()` (`atomic.Int64`); `go test .` oʻtguncha. Bu papka mustaqil laboratoriya,
shuning uchun oldingi darsga nisbatan `diff` yoʻq.

## Manbalar

- [Go wiki: Use a sync.Mutex or a channel?](https://go.dev/wiki/MutexOrChannel) ("Use whichever is most expressive and/or most simple"; Channel: passing ownership of data, distributing units of work, communicating async results; Mutex: caches, state)
- [Effective Go: Share by communicating](https://go.dev/doc/effective_go#sharing) ("This approach can be taken too far. Reference counts may be best done by putting a mutex around an integer variable")
- [Go Proverbs](https://go-proverbs.github.io/) (Rob Pike, Gopherfest SV 2015) · [Go blog: Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [The Go Memory Model](https://go.dev/ref/mem): [#chan](https://go.dev/ref/mem#chan), [#locks](https://go.dev/ref/mem#locks), [#atomic](https://go.dev/ref/mem#atomic)
- [sync/atomic](https://pkg.go.dev/sync/atomic) ([Int64](https://pkg.go.dev/sync/atomic#Int64): "The zero value is zero", "must not be copied after first use") · [sync.Mutex](https://pkg.go.dev/sync#Mutex) · [sync.WaitGroup.Go](https://pkg.go.dev/sync#WaitGroup.Go)
- Go manbasi (go1.27.4): [internal/sync/mutex.go `Lock`](https://github.com/golang/go/blob/go1.27.4/src/internal/sync/mutex.go#L61-L68), [runtime/chan.go `send`](https://github.com/golang/go/blob/go1.27.4/src/runtime/chan.go#L318)
- [Data Race Detector](https://go.dev/doc/articles/race_detector) (Linux: cgo + C kompilyator; vaqt 2–20× koʻproq)
- [Go spec: Select statements](https://go.dev/ref/spec#Select_statements) · [Receive operator](https://go.dev/ref/spec#Receive_operator) · [Composite literals](https://go.dev/ref/spec#Composite_literals)
