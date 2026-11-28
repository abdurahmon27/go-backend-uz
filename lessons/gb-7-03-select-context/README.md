# GB-7.03 · `select` va `context`: ishni vaqtida toʻxtatish

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #53"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-7-02-channel`](../gb-7-02-channel) (channel: yuborish, qabul, buffer, close) · keyingi dars: `gb-7-04`
(chegaralangan worker'lar va `errgroup`). Bu darsga kerak boʻlganlar: [`gb-7-01-goroutine-hayoti`](../gb-7-01-goroutine-hayoti)
(`go func(){…}()`), [`gb-5-05-pgx-pool`](../gb-5-05-pgx-pool) (`context.WithTimeout`, `defer cancel()`, `context deadline exceeded`),
[`gb-6-04-login-session`](../gb-6-04-login-session) (`r.Context()`), [`gb-4-07-mutex-race`](../gb-4-07-mutex-race) (`go run . &`,
`curl -s`, server logi), [`gb-4-05-route-id`](../gb-4-05-route-id) (`ServeMux`, `GET /stock`).

**Bu darsdan keyin:** `select` bilan bir nechta channel'ni birga kutasiz va `time.After` bilan timeout qoʻyasiz; `ctx.Done()` va
`ctx.Err()` bilan context'ni hurmat qiladigan funksiya yozasiz (`context.WithTimeout` → `context deadline exceeded`,
`context.WithCancel` → `context canceled`); handler'da `context.Background()` emas, `r.Context()` ishlatish kerakligini haqiqiy,
uzilgan `curl` bilan koʻrasiz; buffer'siz natija channel'i goroutine leak berishini oʻlchaysiz.

Bu papka mustaqil (Modul 7 laboratoriyasi): oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i,
tashqi modul yoʻq. Yangi nomlar: `supplier(id) chan int` (sekin yetkazib beruvchining taqlidi), `supplierDelay`, `stocks`
(1 olma 40, 2 non 120, 3 oʻrik 5), `checkStock(ctx, id)`, `stockHandler` (`GET /stock`), `sumStock` va `reportHandler`
(`GET /report`, vazifa). Soddalashtirishlar: route `GET /stock` qisqartirilgan (bozor API'sida bu
`GET /products/{id}/stock`, GB-4.05 vazifasi) va oʻrik doim ID 3 (GB-4.04/4.05'dagi katalog: oʻrik 3, 5 dona); yetkazib beruvchi — `time.Sleep` qiladigan goroutine; haqiqiy tizimda bu HTTP soʻrov
yoki baza soʻrovi boʻladi va ctx'ni oʻzi qabul qiladi (masalan, pgx, GB-5.05). Server logi `fmt.Println`; tuzilgan log (`slog`) GB-8.02'da.

## Papka

| Fayl | Nima |
|---|---|
| `supplier.go` | `stocks`, `supplierDelay`, `supplier`: videodagi yetkazib beruvchi (buffer 1) |
| `main.go` | videoning oxiridagi kod (`checkStock`, `r.Context()` bilan `stockHandler`) + vazifa: `sumStock` (boʻsh) va `reportHandler` (`context.Background()` bilan, muddatsiz), `newMux`, `main` |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang. Yetkazib beruvchini 100 ms'ga qisqartiradi |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go`, `supplier.go` (+ xuddi shu test fayli) |
| `_video/NN-…/` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 03-timeout go run .` |
| `_video/race.sh` | `-race` buyruqlari uchun (Linux'da gcc kerak; quyida): `./_video/race.sh 04-cancel go run -race .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64)

**01-kutish** — `main` qoldiqni oddiy qabul (`<-`) bilan kutadi: yetkazib beruvchi qancha kutdirsa, shuncha.

```
$ go run .
qoldiq: 5 3s
```

**02-select** — `select`: supplier'ning channel'i yoki `time.After(time.Second)`, qaysi biri oldin tayyor boʻlsa.

```
$ go run .
supplier javob bermadi
1s
```

**03-timeout** — `checkStock(ctx, id)`: `select` ikkinchi case'da `ctx.Done()`ni kutadi; `main`da `context.WithTimeout(…, time.Second)`.

```
$ go run .
0 context deadline exceeded 1s
```

**04-cancel** — `context.WithCancel`: boshqa goroutine bir soniyadan keyin `cancel()`ni chaqiradi.

```
$ go run .
0 context canceled 1s
```

**05-background** — server (`localhost:8093`), `stockHandler` ichida `checkStock(context.Background(), 3)` (klassik xato).
Server binary qilib quriladi va fonda ishga tushiriladi (`;` — buyruqlarni ketma-ket, `&` — faqat `./bozor`ni fonda;
`go run . &` emas: `kill %1` faqat `go run`ni toʻxtatib, uning ichidagi server portni band qilib qolardi).
`curl --max-time 1` bir soniyadan keyin ulanishni uzadi (`-s` jim, `-S` curl xatosini koʻrsatadi; 28 — curl'ning timeout kodi):

```
$ go build; ./bozor &
$ curl -sS --max-time 1 localhost:8093/stock
curl: (28) Operation timed out after 1000 milliseconds with 0 bytes received
javob yozildi 3s
```

Ikkinchi qator — server logi, ikki soniyadan keyin: mijoz ketgan, server esa yana ikki soniya ishladi.

**06-handler** — tuzatish: `checkStock(r.Context(), 3)`. Server mijoz ulanishi yopilganda `r.Context()`ni bekor qiladi
([`Request.Context`](https://pkg.go.dev/net/http#Request.Context)). Eski serverni toʻxtatib, yangisini quramiz:

```
$ kill %1
$ go build; ./bozor &
$ curl -sS --max-time 1 localhost:8093/stock
curl: (28) Operation timed out after 1000 milliseconds with 0 bytes received
toʻxtadi: context canceled 1s
```

Ikki jarayon bir terminalga bir paytda yozgani uchun bu ikki qator sizda aralashib chiqishi mumkin; millisekundlar (1000/1001)
har safar boshqacha. `curl`ni `--max-time`siz ishga tushirib, Ctrl+C bossangiz ham server `toʻxtadi: context canceled …` yozadi
(vaqt — Ctrl+C bosilgan payt). Handler qaytadi, lekin `supplier`ning goroutine'i 3 soniya uxlashda davom etadi: u ctx olmagan;
haqiqiy yetkazib beruvchi (HTTP yoki baza soʻrovi, masalan pgx) ctx'ni oʻzi oladi va bekor qilinadi.
Ishingiz tugagach serverni toʻxtating: `kill %1`.

**07-leak** ("Chuqurroq") — supplier'da buffer'siz channel (`make(chan int)`): timeout'dan keyin hech kim qabul qilmaydi, goroutine
`ch <- …`da abadiy qoladi. `runtime.NumGoroutine()` tirik goroutine'larni sanaydi:

```
$ go run .                      # make(chan int)
context deadline exceeded goroutine'lar: 2
$ go run .                      # make(chan int, 1)
context deadline exceeded goroutine'lar: 1
```

## `-race`

Linux'da race detector uchun `cgo` va C kompilyator (`gcc`) kerak (batafsil: [`gb-4-07-mutex-race`](../gb-4-07-mutex-race/README.md)).
Bu mashinada gcc yoʻq, shuning uchun tekshiruvlar Docker'da qilingan (`_video/race.sh`; gcc bor `golang:1.26.8` image'i, Go oʻzi
shu mashinadagi go1.27.4):

```
$ ./_video/race.sh 04-cancel go run -race .
0 context canceled 1s
$ ./_video/race.sh solution go test -race -count=3 .
ok  	github.com/abdurahmon27/bozor	3.876s
```

DATA RACE hisoboti yoʻq: `cancel()`ni boshqa goroutine'dan chaqirish xavfsiz ("A CancelFunc may be called by multiple goroutines
simultaneously"), natija esa goroutine'dan channel orqali keladi.

## Vazifa

[TASK.md](TASK.md): `sumStock(ctx, ids)` (har mahsulot `checkStock(ctx, id)` bilan, bekor boʻlsa darhol xato) va `GET /report`'ga
`context.WithTimeout(r.Context(), reportTimeout)` (muddat tugasa 504 Gateway Timeout, mijoz ketsa hech narsa). Tekshirish:

```bash
go test .
```

Javob: [`solution/`](solution). Oʻzgarishni koʻrish: `diff -u main.go solution/main.go`.

Papkaning oʻz-oʻzini tekshiruvi (shu mashinada):

```
$ gofmt -l . && go vet ./...          # chiqish yoʻq
$ go test . ; echo "starter exit $?"   # FAIL (TASK.md'dagi xabarlar), exit 1
$ go test ./solution/...
ok  	github.com/abdurahmon27/bozor/solution	0.954s
```

## Manbalar

- Go spec: [Select statements](https://go.dev/ref/spec#Select_statements) (tayyorlar orasidan tasodifiy tanlov; `default`siz select kutadi),
  [Receive operator](https://go.dev/ref/spec#Receive_operator) (yopilgan channel'dan qabul darhol qaytadi)
- [Go memory model](https://go.dev/ref/mem): channel orqali yuborish va yopish sinxronlashtiradi
- [`context`](https://pkg.go.dev/context): `Done`, `Err`, `WithCancel`, `WithTimeout`, "ctx birinchi parametr", "struct'da saqlamang",
  `CancelFunc` (birinchi chaqiruvdan keyingilari hech narsa qilmaydi), go vet `lostcancel`
- [`net/http` Request.Context](https://pkg.go.dev/net/http#Request.Context): server soʻrov context'ini mijoz ulanishi yopilganda bekor qiladi
- [`time.After`](https://pkg.go.dev/time#After), [`runtime.NumGoroutine`](https://pkg.go.dev/runtime#NumGoroutine)
- Go blog: [Context](https://go.dev/blog/context), [Pipelines and cancellation](https://go.dev/blog/pipelines)
- [curl man page](https://curl.se/docs/manpage.html): `--max-time`, `--show-error`, exit code 28
- [RFC 9110 §15.6.5: 504 Gateway Timeout](https://www.rfc-editor.org/rfc/rfc9110#name-504-gateway-timeout)
