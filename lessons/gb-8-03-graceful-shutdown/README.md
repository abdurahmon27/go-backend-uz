# GB-8.03 · Graceful shutdown: SIGTERM'dan keyin soʻrovni tashlamaslik

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #60"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

8-modulning ("serverga chiqaramiz") uchinchi darsi. Deploy paytida eski jarayon `SIGTERM` oladi (`docker stop`,
`systemctl stop`, `kill`). GB-8.02'gacha bozor bunda **darhol** oʻlardi: ishlayotgan soʻrov yarim yoʻlda uzilardi
(`curl: (52) Empty reply from server`, chiqish kodi `143`), handler fonda boshlagan chek esa yoʻqolardi. Endi:

1. `signal.NotifyContext` — `SIGINT` (Ctrl+C) yoki `SIGTERM` dasturni oʻldirmaydi, `sigCtx`ni bekor qiladi;
2. `ListenAndServe` alohida goroutine'da, `select` server xatosini (masalan, port band) yoki signalni kutadi;
3. `srv.Shutdown(shCtx)` 8 soniyalik muddat bilan: port yopiladi (yangi ulanish rad etiladi), boʻsh ulanishlar yopiladi,
   faol soʻrovlar tugashi kutiladi; muddat tugasa — `srv.Close()`, xato va `1` kodi;
4. chek (`sendReceipt`) endi egasi server boʻlgan `sync.WaitGroup` orqali (`a.bg.Go`), `run` uni `Shutdown`'dan keyin
   `waitBG` bilan kutadi;
5. `defer pool.Close()` — Postgres pool'i eng oxirida;
6. serverning toʻrt timeout'i: `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`.

Oldingi darslar: [`gb-8-02-slog`](../gb-8-02-slog) (bu papka uning `solution/`idan boshlanadi; JSON log, `request_id`) ·
[`gb-8-01-config`](../gb-8-01-config) (`run(args)`, `config`) · [`gb-7-01-goroutine-hayoti`](../gb-7-01-goroutine-hayoti)
(`go sendReceipt(o)` — «kim kutadi?», egasi qoidasi, `wg.Go`) · [`gb-7-03-select-context`](../gb-7-03-select-context)
(`select`, `ctx.Done()`) · [`gb-4-04-birinchi-server`](../gb-4-04-birinchi-server) (`http.Server`, `ReadHeaderTimeout`) ·
keyingi dars: `gb-8-04` (Docker).
Farqlarni koʻrish: `diff -ru ../gb-8-02-slog/solution .` — yangi `cmd/api/receipt.go`, `cmd/api/toxtash_test.go`,
`buyurtma.sh`; oʻzgargan `cmd/api/main.go` (signal, Shutdown, timeout'lar), `handlers.go` (`api.bg`, `createOrder`'da
chek), `middleware.go` (`newHandler(st, logger, bg)`), testlarda (`api_test.go`, `tekshiruv_test.go`)
`newHandler(…, new(sync.WaitGroup))`. Laboratoriya Postgres'i boshqa portda: `127.0.0.1:5583` (compose loyihasi `gb803`) —
shuning uchun `compose.yaml`, `.env.example`, `internal/pgtest/pgtest.go`, `cmd/api/config_test.go` va `cmd/api/log_test.go`da
faqat port (5582 → 5583) va loyiha nomi oʻzgargan.

`run`da `signal.NotifyContext` eng boshida, Postgres'ga ulanishdan oldin turadi: ulanish paytida kelgan signal ham
ulanishni bekor qiladi (ulanish context'i `sigCtx`dan). Videoda bu qatorlar alohida koʻrsatilgan; tartibi shu faylda.

## Toʻxtash tartibi

```
SIGTERM / Ctrl+C
  └─ sigCtx.Done() yopiladi → select'dan chiqamiz → stop() (ikkinchi Ctrl+C endi darhol oʻldiradi)
       └─ srv.Shutdown(shCtx)      1) port yopiladi  2) boʻsh ulanishlar yopiladi  3) faol soʻrovlar kutiladi
            └─ waitBG(shCtx, &bg)  handler'lar boshlagan cheklar
                 └─ run qaytadi → defer pool.Close() → main: kod 0 (xato boʻlsa: stderr'ga «xato: …», kod 1)
```

Nega aynan shu tartib: soʻrovlar va cheklar Postgres'ni ishlatadi, shuning uchun pool ular tugagandan keyin yopiladi.
`pgxpool.Pool.Close` hujjatda: «Blocks until all connections are returned to pool and closed» — pool'ni ishlatayotgan
soʻrov qolsa, `Close` uni kutadi. Shuning uchun muddat tugaganda `srv.Close()` chaqiriladi: u qolgan ulanishlarni uzadi,
ularning `r.Context()`i bekor boʻladi va `r.Context()` bilan ketgan SQL soʻrovlar ham toʻxtaydi.

Muddat tugagan yoʻlda (`srv.Shutdown` xato qaytarsa) fon ishlari kutilmaydi — muddat allaqachon tugagan — va run darhol
qaytadi. `srv.Close()` qolgan soʻrovlarning `r.Context()`ini bekor qiladi (`_video/closectx`da tekshirilgan), bozor'dagi har
SQL soʻrov `r.Context()` bilan ketadi, shuning uchun ular qaytadi va `pool.Close()` osilib qolmaydi. Context'siz
Postgres'ni band qilgan kod boʻlsa, `pool.Close()` uni kutib qoladi: har soʻrovga context bering (GB-5.05).

Chek goroutine'i oʻz panic'ini oʻzi ushlaydi (`receipt.go`): fon goroutine'idagi panic'ni `recoverPanic` middleware
ushlamaydi va u butun serverni yiqitadi (GB-7.01).

`shutdownTimeout` = 8 s: `docker stop` standart holda 10 s kutadi, keyin `SIGKILL` yuboradi (`--timeout`/`-t` bilan
oʻzgaradi; compose'da `stop_grace_period`). Bizning muddat undan qisqa boʻlishi kerak, aks holda `SIGKILL` muddatni kesadi
va log'da sabab qolmaydi. Kubernetes'da bu `terminationGracePeriodSeconds` (standart 30 s) — GB-12'da.

## Shutdown nimani qilmaydi

- **Hijack qilingan ulanishlar** (WebSocket va h.k.): `Shutdown` ularni yopmaydi ham, kutmaydi ham (hujjat).
  Ularni `srv.RegisterOnShutdown(f)` bilan ogohlantirasiz va oʻzingiz kutasiz.
- **Handler boshlagan goroutine'lar**: `Shutdown` faqat HTTP ulanishlarini biladi. Videodagi `go sendReceipt(…)` —
  jarayon chiqqanda chek jim yoʻqoladi. Tuzatish: goroutine'ning egasi (GB-7.01) — server: `a.bg.Go(…)` + `waitBG`.
- **Soʻrov context'ini bekor qilmaydi**: `Shutdown` paytida `r.Context()` tirik, handler oxirigacha ishlaydi. Agar
  uzoq handler'lar toʻxtash haqida bilishi kerak boʻlsa — `http.Server.BaseContext` bilan ularga signal context'idan
  hosil boʻlgan context bering (bozor'da kerak emas: soʻrovlar qisqa).
- **ListenAndServe'ni kutmaydi**: aksincha, `Shutdown` chaqirilishi bilan `ListenAndServe` darhol `http.ErrServerClosed`
  qaytaradi. Hujjat: «Make sure the program doesn't exit and waits instead for Shutdown to return». Shuning uchun
  `run` `errc`ni signal kelgandan keyin oʻqimaydi (buffer 1: goroutine bloklanmaydi) va `Shutdown` qaytishini kutadi.

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi `_video/chiqishlar.sh` bilan shu kompyuterda qayta ishga tushirilgan (skriptda `curl -sS`: chiqish terminal
emas, progress chiqmasin; terminalda oddiy `curl` xuddi shu qatorlarni beradi): go1.27.4 linux/amd64, `CGO_ENABLED=0`,
PostgreSQL 18 (`postgres:18-alpine`), API `:8783`da. Toʻliq chiqish: [`_video/chiqishlar.txt`](_video/chiqishlar.txt).
JSON log qatorlari bu yerda qisqartirilgan (`"time"` va `request_id` qiymatlari `…`).

**Namoyish uchun sekin soʻrov.** Videoda `/health`ga vaqtincha ikki qator qoʻshiladi (papkada ular **yoʻq**):

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
	d, _ := time.ParseDuration(r.FormValue("sekin")) // namoyish uchun, vaqtincha
	time.Sleep(d)
	fmt.Fprintln(w, "ok")
}
```

Tayyorgarlik (bir marta): `docker compose up -d --wait`, `goose up` va `seed.sql` (GB-5.06; `cmd/api/main.go`
boshidagi izohda), `cp .env.example .env` (PAROL → `bozor-lab`), `set -a; . ./.env; set +a`.

**1. Eski kod (GB-8.02) + sekin `/health`: SIGTERM soʻrovni uzadi.** Uch terminal:

```console
# 1-terminal
$ go build ./cmd/api && ./api -port 8783 ; echo "kod: $?"
…
{"time":"…","level":"INFO","msg":"bozor API ishga tushdi","addr":":8783"}
Terminated
kod: 143

# 2-terminal (uch soniyalik soʻrov)
$ curl 'localhost:8783/health?sekin=3s'
curl: (52) Empty reply from server

# 3-terminal (bir soniyadan keyin)
$ pkill -e -x api
api killed (pid 1397638)
```

`143` = 128 + 15 (`SIGTERM`ning raqami): jarayon signal bilan oʻldirilgan (Bash: «Exit Status»). `Terminated` — shell'ning
xabari. `pkill` signal koʻrsatilmasa `SIGTERM` yuboradi; `-x` — nomi aynan `api`; `-e` — kimga yuborganini yozadi.

**2. Yangi kod: eski soʻrov tugaydi, yangisi rad etiladi.**

```console
# 3-terminal
$ pkill -e -x api
api killed (pid 1397783)
$ curl localhost:8783/health
curl: (7) Failed to connect to localhost:8783 after 0 ms: Could not connect to server

# 2-terminal
$ curl 'localhost:8783/health?sekin=3s'
ok

# 1-terminal
…
{"time":"…","level":"INFO","msg":"toʻxtash boshlandi"}
{"time":"…",…,"msg":"soʻrov",…,"path":"/health","status":200,…}
{"time":"…","level":"INFO","msg":"toʻxtadi"}
kod: 0
```

**3. Muddat tugaydi** (`sekin=15s`, `shutdownTimeout` 8 s): 9 soniyadan keyin (1 s + 8 s)

```console
# 2-terminal
curl: (52) Empty reply from server
# 1-terminal
…
{"time":"…","level":"INFO","msg":"toʻxtash boshlandi"}
xato: shutdown: context deadline exceeded
kod: 1
```

Shu 8 soniya ichida ikkinchi marta Ctrl+C (`SIGINT`) bossangiz, `stop()` signalning standart ishini qaytargani uchun
jarayon darhol oʻladi: `kod: 130` (128 + 2). Skriptda bu `pkill -INT` bilan tekshirilgan (`1c` bosqichi).

**4. Chek: egasiz goroutine vs egasi server.** `./buyurtma.sh` Jasur nomidan login qilib, bitta buyurtma yuboradi
(GB-6.04 login + cookie, GB-5.10 `POST /orders`); `&&` — muvaffaqiyatli boʻlsa darhol `pkill`.

```console
$ ./buyurtma.sh 8783 && pkill -x api
POST /orders → 201
```

`go sendReceipt(logger, o)` bilan (xato) — 1-terminal:

```
{"time":"…",…,"msg":"soʻrov",…,"path":"/orders","status":201,…}
{"time":"…","level":"INFO","msg":"toʻxtash boshlandi"}
{"time":"…","level":"INFO","msg":"toʻxtadi"}
kod: 0
```

`a.bg.Go(func() { sendReceipt(logger, o) })` bilan (papkaning oʻzi):

```
{"time":"…","level":"INFO","msg":"toʻxtash boshlandi"}
{"time":"…",…,"msg":"chek yuborildi","request_id":"…","order_id":44}
{"time":"…","level":"INFO","msg":"toʻxtadi"}
kod: 0
```

Goroutine'ga `r.Context()` emas, `logger` beriladi: soʻrov context'i handler qaytganda bekor boʻladi, chek esa handler'dan
keyin ham ishlaydi. Chek qatorida buyurtma soʻrovining `request_id`si bor — logger soʻrovdan olingan (GB-8.02).

**5. `ReadHeaderTimeout`: sarlavhasini tugatmagan mijoz.**

```console
$ time printf 'GET /health HTTP/1.1\r\n' | nc localhost 8783

real	0m5.032s
user	0m0.002s
sys	0m0.001s
```

Faqat soʻrov qatori ketdi, sarlavhalar tugamadi (`\r\n\r\n` yoʻq). Server 5 soniya kutdi va ulanishni javobsiz yopdi
(OpenBSD `nc` stdin tugagach ham ulanishni ochiq ushlaydi). Timeout'siz bunday ulanishlar minglab ochilsa, server ular
bilan toʻlib qoladi: **slowloris** hujumi.

| maydon | bozor | nimani cheklaydi (go doc net/http.Server) |
|---|---|---|
| `ReadHeaderTimeout` | 5 s | soʻrov qatori va sarlavhalarni oʻqish |
| `ReadTimeout` | 10 s | butun soʻrov, body bilan |
| `WriteTimeout` | 30 s | sarlavhalar oʻqilgandan javob yozilib boʻlguncha; eng sekin handler'dan katta boʻlsin |
| `IdleTimeout` | 1 daqiqa | keep-alive ulanish keyingi soʻrovni kutishi (nol boʻlsa — `ReadTimeout`) |

`WriteTimeout` handler'ni toʻxtatmaydi: muddatdan keyin javob yozib boʻlmaydi, ulanish uziladi. Handler'ning oʻzini
cheklash uchun `http.TimeoutHandler` yoki `context.WithTimeout` (GB-7.03) kerak.

## waitBG'dan qolgan goroutine

`solution/`dagi `waitBG` `bg.Wait()`ni alohida goroutine'da kutadi. Muddat tugasa, `waitBG` qaytadi, u goroutine esa
`bg.Wait()`da qoladi — GB-7.06'dagi leak shakli. Bu yerda bu toʻgʻri: `waitBG` faqat jarayon chiqishidan oldin bir marta
chaqiriladi va jarayon bilan birga tugaydi. Uzoq yashaydigan kodda (masalan, har soʻrovda) bunday kutish leak boʻlardi.

## Papka

```
gb-8-03-graceful-shutdown/
  cmd/api/main.go             run: NotifyContext, ListenAndServe goroutine'da, select, Shutdown, waitBG; toʻrt timeout
  cmd/api/receipt.go          sendReceipt (simulyatsiya, 2 s), waitBG — VAZIFA
  cmd/api/handlers.go         api.bg; createOrder: a.bg.Go(func() { sendReceipt(logger, o) })
  cmd/api/toxtash_test.go     vazifa tekshiruvi (Postgres'siz): go test ./cmd/api -run WaitBG; + timeout'lar, chek
  cmd/api/…                   GB-8.02'dagi API (config, slog, handler'lar, middleware, testlar — haqiqiy Postgres'da)
  buyurtma.sh                 login + bitta buyurtma (videodagi chek namoyishi)
  compose.yaml                Postgres 18, 127.0.0.1:5583
  solution/                   javob (oʻz go.mod'i bilan, alohida ishlaydi)
  _video/                     chiqishlar.sh → chiqishlar.txt
```

## Vazifa

[`TASK.md`](TASK.md): `waitBG(ctx, bg)` muddat tugasa `ctx.Err()` bilan darhol qaytsin. Tekshirish:
`go test ./cmd/api -run WaitBG` (Postgres kerak emas); butun papka: `docker compose up -d --wait && go test ./...`.

## Manbalar

- Go: [net/http Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown) («first closing all open listeners, then
  closing all idle connections, and then waiting indefinitely for connections to return to idle»; «Make sure the program
  doesn't exit and waits instead for Shutdown to return»; «does not attempt to close nor wait for hijacked connections»),
  [Server.Close](https://pkg.go.dev/net/http#Server.Close), [Server.RegisterOnShutdown](https://pkg.go.dev/net/http#Server.RegisterOnShutdown),
  [Server (timeout maydonlari)](https://pkg.go.dev/net/http#Server), [ErrServerClosed](https://pkg.go.dev/net/http#ErrServerClosed)
- Go: [os/signal NotifyContext](https://pkg.go.dev/os/signal#NotifyContext) (stop «may restore the default behavior»),
  [os/signal: Default behavior](https://pkg.go.dev/os/signal#hdr-Default_behavior_of_signals_in_Go_programs) («A SIGHUP,
  SIGINT, or SIGTERM signal causes the program to exit»), [syscall.SIGTERM](https://pkg.go.dev/syscall#SIGTERM),
  [sync.WaitGroup.Go](https://pkg.go.dev/sync#WaitGroup.Go), [time.ParseDuration](https://pkg.go.dev/time#ParseDuration),
  [Go spec: Defer statements](https://go.dev/ref/spec#Defer_statements), [Go spec: Select statements](https://go.dev/ref/spec#Select_statements)
- [pgxpool.Pool.Close](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#Pool.Close) — «Blocks until all connections are
  returned to pool and closed»
- [docker container stop](https://docs.docker.com/reference/cli/docker/container/stop/) — `SIGTERM`, `--timeout` (standart
  10 s), keyin `SIGKILL`; [signal(7)](https://man7.org/linux/man-pages/man7/signal.7.html) (SIGINT 2, SIGKILL 9, SIGTERM 15);
  [pkill(1)](https://man7.org/linux/man-pages/man1/pkill.1.html) (standart signal `SIGTERM`, `-x`, `-e`)
- [Bash: Exit Status](https://www.gnu.org/software/bash/manual/bash.html#Exit-Status) (signal bilan oʻlgan buyruq: 128+N);
  [curl: exit codes](https://curl.se/libcurl/c/libcurl-errors.html) (7 `CURLE_COULDNT_CONNECT`, 52 `CURLE_GOT_NOTHING`)
- [RFC 9112 §9.3 Persistence](https://www.rfc-editor.org/rfc/rfc9112#section-9.3) (keep-alive);
  [Cloudflare: Slowloris](https://www.cloudflare.com/learning/ddos/ddos-attack-tools/slowloris/)
- [Dave Cheney: Never start a goroutine without knowing how it will stop](https://dave.cheney.net/2016/12/22/never-start-a-goroutine-without-knowing-how-it-will-stop)
- Alex Edwards, *Let's Go Further*, 11-bob (graceful shutdown) va 13.4 (fon goroutine'larini kutish) — shu mavzuning eng
  yaxshi yozma bayoni.
