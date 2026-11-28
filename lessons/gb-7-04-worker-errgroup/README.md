# GB-7.04 · Chegaralangan worker'lar va `errgroup`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #54"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-7-03-select-context`](../gb-7-03-select-context) (`select`, `context`, `checkStock(ctx, id)`, ketma-ket
`sumStock`) · keyingi dars: `gb-7-05` (mutex yoki channel). Bu darsga kerak boʻlganlar:
[`gb-7-02-channel`](../gb-7-02-channel) (channel, `close`, `range`, buffered channel),
[`gb-7-01-goroutine-hayoti`](../gb-7-01-goroutine-hayoti) (`go`, `wg.Go`, har goroutine oʻz indeksiga yozadi),
[`gb-6-07-cookie-csrf-cors`](../gb-6-07-cookie-csrf-cors) (`429 Too Many Requests`, rate limit),
[`gb-6-03-parollar`](../gb-6-03-parollar) (`go get golang.org/x/…`), [`gb-4-07-mutex-race`](../gb-4-07-mutex-race) (`-race`),
[`gb-2-08-xatolar`](../gb-2-08-xatolar) (`%w`, `errors.Join`).

**Bu darsdan keyin:** har mahsulotga alohida goroutine ochish nega buzilishini haqiqiy raqamlarda koʻrasiz (yetkazib beruvchi
1000 soʻrovdan 990 tasini rad etadi; 100 000 ID'da 25 mingdan ortiq goroutine va 123 MB); bitta `jobs` channel ustida oʻnta
worker (worker pool) yozasiz; `golang.org/x/sync/errgroup` bilan (`Group`, `SetLimit`, `Go`, `Wait`) xuddi shu chegarani
xato qaytarish bilan birga olasiz; `errgroup.WithContext` bilan birinchi xatoda qolgan soʻrovlarni bekor qilasiz; guruh
context'ini `Wait`dan keyin ishlatish xatosini (`context canceled`) taniysiz.

Bu papka mustaqil (Modul 7 laboratoriyasi): oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i,
bitta tashqi modul: `golang.org/x/sync v0.25.0` (`go.sum` bilan). Yetkazib beruvchi **soxta** (`supplier.go`, tarmoqsiz):
GB-7.03'dagi yetkazib beruvchining tezroq va katta nusxasi (olma 1, non 2, oʻrik 3 qoldigʻi oʻsha: 40, 120, 5; qolgan ID'lar uchun
`ID % 50`). Bekor qilishdagi soʻrovlar soni va cheklovsiz parallel ishlar rejalashtirishga bogʻliq. Har soʻrov 20 ms, bir vaqtda 10 tadan koʻpi rad etiladi
(`ErrTooMany`, "429 too many requests"), faqat 1…1000 ID'lar mavjud (`ErrNotFound`). `saveStocks` ham soxta: faqat
context'ni tekshiradi; haqiqiy Postgres'ga yozish GB-7.07 nazorat nuqtasida (bekor qilinadigan importer).

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi `importStocks` (`errgroup.WithContext` + `SetLimit`, `gctx`) + vazifa uchun `fetchPrices` (cheklovsiz starter) |
| `supplier.go` | soxta yetkazib beruvchi: `checkStock`, `checkPrice`, `saveStocks`, `catalogIDs`, hisoblagichlar (`stats`) va `report` |
| `import_test.go` | videodagi `importStocks`ning testi: starterda ham oʻtadi |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu `supplier.go` va test fayllari) |
| `_video/NN-…/` | videodagi har bir qadamning `main.go`'si; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 04-worker go run .` |
| `_video/race.sh` | `-race` buyruqlari uchun (quyidagi eslatma): `./_video/race.sh 04-worker go run -race .` |

`report` uch qator chiqaradi: (1) yetkazib beruvchi sanagan soʻrovlar, eng koʻpi bilan bir vaqtda nechtasi ishlagani,
nechtasi rad etilgani va importning vaqti; (2) natija: `qoldiqlar: N` yoki xatoning birinchi qatori (`errors.Join` xatosi koʻp
qatorli boʻladi: "va yana N ta"); (3) dastur ishlagan paytdagi eng koʻp goroutine soni va Go runtime operatsion tizimdan olgan
xotira (`runtime.MemStats.Sys`).

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, 12 yadro)

Har qadam `./_video/run.sh <qadam> go run .` bilan. Vaqt, xotira, cheklovsiz parallel ishlar va bekor qilishdagi soʻrovlar soni sizda boshqacha boʻlishi mumkin.

```
$ ./_video/run.sh 01-ketma go run .          # ketma-ket, sumStock kabi
soʻrov: 1000 · bir vaqtda: 1 · rad: 0 · vaqt: 20.218s
qoldiqlar: 1000
goroutine: 1 · xotira: 7 MB

$ ./_video/run.sh 02-cheksiz go run .        # har ID uchun wg.Go
soʻrov: 1000 · bir vaqtda: 1000 · rad: 990 · vaqt: 30ms
xato: 429 too many requests (va yana 989 ta)
goroutine: 1001 · xotira: 12 MB

$ ./_video/run.sh 03-xotira go run .         # xuddi shu, 100 000 ID
soʻrov: 100000 · bir vaqtda: 23580 · rad: 99990 · vaqt: 249ms
xato: 429 too many requests (va yana 99989 ta)
goroutine: 25151 · xotira: 123 MB
```

`03-xotira`da "bir vaqtda", goroutine soni va xotira har safar boshqacha (shu mashinada 23–36 ming, 123–146 MB): sikl
yangilarini boshlayotganda eski goroutine'lar tugay boshlaydi. Muhimi: cheklovsiz kodda bu sonlar **ID'lar soni bilan birga oʻsadi** (1 000 000 ID'da shu
mashinada 67 995 goroutine va 370 MB).

```
$ ./_video/race.sh 04-worker go run -race .  # jobs channel + 10 worker
soʻrov: 1000 · bir vaqtda: 10 · rad: 0 · vaqt: 2.024s
qoldiqlar: 1000
goroutine: 11 · xotira: 8 MB

$ GOTOOLCHAIN=go1.27.4 go get golang.org/x/sync/errgroup
go: downloading golang.org/x/sync v0.25.0
go: added golang.org/x/sync v0.25.0

$ ./_video/run.sh 05-errgroup go run .       # var g errgroup.Group + SetLimit
soʻrov: 1000 · bir vaqtda: 10 · rad: 0 · vaqt: 2.034s
qoldiqlar: 1000
goroutine: 12 · xotira: 12 MB

$ ./_video/run.sh 06-xato go run .           # ids[499] = 9999, oddiy Group
soʻrov: 1000 · bir vaqtda: 10 · rad: 0 · vaqt: 2.043s
xato: mahsulot 9999: mahsulot topilmadi

$ ./_video/run.sh 07-withctx go run .        # errgroup.WithContext
soʻrov: 509 · bir vaqtda: 10 · rad: 0 · vaqt: 1.026s
xato: mahsulot 9999: mahsulot topilmadi

$ ./_video/run.sh 08-wait-ctx go run .       # Wait'dan keyin guruh ctx'i bilan saveStocks
soʻrov: 1000 · bir vaqtda: 10 · rad: 0 · vaqt: 2.027s
xato: saqlash: context canceled

$ ./_video/run.sh 09-gctx go run .           # gctx: tuzatilgan (= main.go'dagi importStocks)
soʻrov: 1000 · bir vaqtda: 10 · rad: 0 · vaqt: 2.024s
qoldiqlar: 1000
```

(`06`…`09`ning uchinchi qatori videoda qisqartirilgan.) Nimalarga e'tibor bering:

- **Worker pool'da goroutine 11** (10 worker + `main`), errgroup'da **11 yoki 12**: errgroup goroutine'i oʻrnini (limit
  "slot"ini) `defer` bilan, funksiyasi tugashi bilanoq boʻshatadi; yangisi boshlanganda eskisi hali chiqib ulgurmagan boʻlishi mumkin.
- **Bu safar 509 soʻrov** (kuzatilgan: 504–509; kafolatlangan oraliq emas): aniq son rejalashtirishga bogʻliq. Sikl qolgan ~490 ta `g.Go`ni baribir chaqiradi
  (errgroup siklni toʻxtatmaydi), lekin `checkStock` context bekor ekanini koʻrib, soʻrov yubormasdan qaytadi. Bekor qilish faqat
  context'ni hurmat qiladigan funksiyada ishlaydi (GB-7.03). `Wait` faqat **birinchi** xatoni qaytaradi; bekor qilinganlarning
  `context canceled` xatolari tashlanadi.
- **`context canceled`:** `errgroup.WithContext` bergan context birinchi xatoda **yoki `Wait` qaytganda** bekor qilinadi
  ([hujjat](https://pkg.go.dev/golang.org/x/sync/errgroup#WithContext)). `g, ctx := errgroup.WithContext(ctx)` yozilsa, `:=`
  faqat `g`ni yangi yaratadi: `ctx` funksiya parametri bilan bir xil blokda, shuning uchun parametrga guruh context'i yoziladi
  ([Go spec: Short variable declarations](https://go.dev/ref/spec#Short_variable_declarations)) va `Wait`dan keyingi
  `saveStocks(ctx, …)` bekor qilingan context oladi. Shuning uchun guruh context'iga alohida nom (`gctx`) beramiz.
- **`SetLimit(0)`**: ichida sigʻimi 0 boʻlgan channel, `g.Go` hech qachon qaytmaydi. Oddiy dasturda boshqa goroutine boʻlmasa,
  runtime `fatal error: all goroutines are asleep - deadlock!` bilan toʻxtatadi (`goroutine 1 [chan send]: …errgroup.(*Group).Go`);
  serverda yoki testda (boshqa goroutine'lar tirik) esa osilib qoladi. **panic** errgroup'da ushlanmaydi, butun dasturni toʻxtatadi ([v0.25.0 manbasi](https://cs.opensource.google/go/x/sync/+/refs/tags/v0.25.0:errgroup/errgroup.go)).

## `-race`: Linux va macOS talablari

Linux'da race detector uchun `cgo` yoqilgan va C kompilyator (odatda `gcc`) oʻrnatilgan boʻlishi kerak; macOS'da Go 1.20'dan beri
Xcode shart emas. Bu mashinada `gcc` yoʻq, shuning uchun `-race` natijalari Docker'da olingan (`_video/race.sh`: gcc bor
`golang:1.26.8` image'i, Go oʻzi shu mashinadagi `go1.27.4` toolchain, modul keshi faqat oʻqish uchun ulanadi). Hamma qadam va
yechim `-race` bilan toza: `./_video/race.sh solution go test -race -count=1 .` → `ok`.

## Vazifa

[TASK.md](TASK.md): `fetchPrices(ctx, ids, limit)`ni `errgroup.WithContext` + `SetLimit` bilan qayta yozing; `go test .` oʻtguncha.

## Manbalar

- [golang.org/x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup): `Group`, `WithContext`, `Go`, `SetLimit`, `TryGo`, `Wait`
- [errgroup.go v0.25.0](https://cs.opensource.google/go/x/sync/+/refs/tags/v0.25.0:errgroup/errgroup.go)
- [Go spec: Go statements](https://go.dev/ref/spec#Go_statements) · [Program execution](https://go.dev/ref/spec#Program_execution)
- [Go memory model](https://go.dev/ref/mem) · [sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup)
- [Go blog: Pipelines, bounded parallelism](https://go.dev/blog/pipelines#bounded-parallelism)
- [context](https://pkg.go.dev/context) · [errors.Join](https://pkg.go.dev/errors#Join) · [race detector](https://go.dev/doc/articles/race_detector)
- [RFC 6585 §4: 429 Too Many Requests](https://www.rfc-editor.org/rfc/rfc6585#section-4)
