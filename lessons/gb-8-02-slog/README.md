# GB-8.02 · `log/slog`: qidiriladigan structured log'lar

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #59"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

8-modulning ("serverga chiqaramiz") ikkinchi darsi. GB-8.01'dagi bozor API'sining log'i `log.Printf` matnidan
**`log/slog` JSON qatorlariga** oʻtadi: har yozuv — bitta JSON obyekt (`time`, `level`, `msg` va atributlar), shuning
uchun log'ni maydon boʻyicha qidirish mumkin. GB-6.01'dagi `requestID` middleware endi soʻrov context'iga **request_id
biriktirilgan logger** qoʻyadi: bitta soʻrovning hamma qatorida bir xil `request_id` bor. Darsda bozor kodidagi
**ikki haqiqiy sizish** koʻrsatiladi va tuzatiladi: JSON decode xatosi paroldan besh harfni (`\ulugb`), store xatosi esa
emailni log'ga yozardi.

Oldingi darslar: [`gb-8-01-config`](../gb-8-01-config) (bu papka uning `solution/`idan boshlanadi; `config.String`) ·
[`gb-6-01-middleware`](../gb-6-01-middleware) (`requestID`, `logRequests`, `X-Request-Id`) ·
[`gb-6-02-api-xatolar`](../gb-6-02-api-xatolar) (`writeError`: "log'ga toʻliq, mijozga xavfsiz") ·
[`gb-6-04-login-session`](../gb-6-04-login-session) (`context.WithValue`, `currentUser`, `UserByEmail`) ·
[`gb-3-02-json-saqlash`](../gb-3-02-json-saqlash) (kichik harfli maydon JSON'ga chiqmaydi) · keyingi dars: `gb-8-03` (graceful shutdown).

Farqlarni koʻrish: `diff -ru ../gb-8-01-config/solution .` — oʻzgargan: `cmd/api/main.go` (logger), `middleware.go`
(`requestID(logger, next)`, `loggerKey`, `loggerFrom`, `logRequests`, `recoverPanic`), `problem.go` (`writeError(w, r, err)`),
`config.go` (`LogValue`, `dbHost`), `internal/store/store.go` (xato matnlarida email yoʻq), 27 ta `writeError` chaqiruvi
(`w, r, err`), testlar (`quiet` logger); yangi: `cmd/api/log_test.go`, ikki `tekshiruv_test.go`; GB-8.01'ning
`tekshiruv_test.go`si endi `config_test.go`. Laboratoriya Postgres'i boshqa portda: `127.0.0.1:5582` (compose loyihasi
`gb802`).

## Asosiy qoidalar

```go
// cmd/api/main.go, run: bitta asosiy logger, aniq parametr bilan beriladi
handler := slog.NewJSONHandler(os.Stdout, nil) // har yozuv — bitta JSON qator; nil → daraja Info'dan
logger := slog.New(handler)
slog.SetDefault(logger)                         // slog.Default() ham shu JSON logger (loggerFrom'ning zaxirasi)
logger.Info("config", "config", cfg)            // cfg.LogValue(): faqat addr va db
srv := newServer(cfg.addr(), store.New(pool), logger)

// cmd/api/middleware.go: soʻrovning logger'i — context'da, request_id bilan
ctx := context.WithValue(r.Context(), loggerKey{}, logger.With("request_id", id))
next.ServeHTTP(w, r.WithContext(ctx))
// … istalgan joyda: loggerFrom(r.Context()).Info(…)

// cmd/api/problem.go: 4xx — faqat biz yozgan matn, 5xx — toʻliq xato
if p.Status >= 500 {
	logger.Error("ichki xato", "status", p.Status, "err", err)
} else {
	logger.Warn("mijoz xatosi", "status", p.Status, "detail", p.Detail)
}
```

| daraja | bozor'da | ichki qiymati |
|---|---|---|
| `Debug` | `"Postgres ulandi"` — standart sozlamada yozilmaydi | −4 |
| `Info` | har soʻrov qatori (`logRequests`), ishga tushish | 0 |
| `Warn` | mijozning xatosi (4xx): status + detail | 4 |
| `Error` | serverning xatosi (5xx) va panic: toʻliq `err`, stack | 8 |

**Nega soʻrovning logger'i context'da, asosiy logger esa parametr?** Asosiy logger bitta, `run` uni yasaydi va
`newServer` → `newHandler` → `requestID` ga aniq beradi (testda `slog.DiscardHandler` yoki `bytes.Buffer`ga yozadigan
logger beriladi). Soʻrovning logger'i esa har soʻrovda yangi (`With("request_id", id)`); `writeError` 27 joydan,
jumladan `decodeJSON`, `pathID` kabi oddiy funksiyalardan chaqiriladi va ularning hammasida `r` bor. `r.Context()`dan
olingan logger ID'ni allaqachon biladi — uni qoʻshishni unutib boʻlmaydi. `context` paketi hujjati: context qiymatlari
«request-scoped data» uchun. Kalit — oʻz turimiz `loggerKey{}` (`go doc context.WithValue`: kalit oʻrnatilgan tur
boʻlmasin). `loggerFrom` context'da logger topmasa `slog.Default()` qaytaradi; `run` uni `slog.SetDefault(logger)` bilan
oʻsha JSON logger qilgan, shuning uchun qator `request_id`'siz, lekin JSON'da va stdout'da (`SetDefault`siz `slog.Default()`
matnni stderr'ga yozardi va `> api.log`ga tushmasdi). Bu — tanlov: slog hujjati logger'ni context'da saqlashni buyurmaydi;
boshqa yoʻl — `ctx`dan qiymat oʻqiydigan oʻz `slog.Handler`ingiz va hamma joyda `InfoContext(ctx, …)`.

**Nega `config` JSON handler'da `{}` edi?** `JSONHandler` qiymatni `encoding/json` kabi yozadi (`go doc
log/slog.JSONHandler.Handle`): `config`ning maydonlari kichik harfda va JSON'ga chiqmaydi, `String()` esa chaqirilmaydi
(uni `TextHandler` ishlatadi). `LogValue() slog.Value` (`slog.LogValuer`) bor qiymatni handler yozishdan oldin shu metod
natijasiga almashtiradi — biz faqat kerakli maydonlarni tanlaymiz (`addr`, `db`).

**Videodan tashqari, kodda:** `recoverPanic` panic'ni `Error` darajada `panic` va `stack` atributlari bilan yozadi;
`writeJSON`dagi `Encode` xatosi log'ga yozilmaydi (status allaqachon ketgan, odatda mijoz uzilgan — yozsangiz, `Debug`
darajada); `TestMain` testlarda `slog.SetDefault(quiet)` qiladi. Ishga tushishdagi sozlama xatolari (`xato: …`, GB-8.01)
logger yasalishidan oldin chiqadi va oddiy matn boʻlib stderr'ga ketadi. Muvaffaqiyatsiz login ham log'da bor
(`Warn "mijoz xatosi"`, `status` 401, `path` `/login`, `request_id`) — kim ekanini (emailni) yozmasdan. `logRequests`dagi yoʻl endi `%q`siz: JSON handler satr
ichidagi yangi qatorni `\n` qilib qochiradi (GB-6.01'dagi log injection), har yozuv bitta qatorda qoladi.

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, PostgreSQL 18 `postgres:18-alpine`,
goose v3.30.2, curl 8.22.0, jq 1.8.2). Toʻliq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt)
(`bash _video/chiqishlar.sh`; bosqichlar skriptning boshida yozilgan). Baza: `docker compose up -d --wait`, `goose up`,
`docker compose exec db psql -f /lab/seed.sql` (GB-5.06), `DATABASE_URL=postgres://bozor:bozor-lab@localhost:5582/bozor`.

```console
# GB-8.01 kodi (log.Printf) — hook'dagi ikki qator:
2026/11/28 18:18:33 D3QGB6DSPPCGJKDRQY5EPR2TRS xato 401: "email yoki parol notoʻgʻri: email \"nodira@bozor.test\": foydalanuvchi topilmadi"
2026/11/28 18:18:33 MTUOSLB5TGTBGISF65BF2RYVLA xato 400: "body JSON sifatida oʻqilmadi: invalid escape sequence `\\ulugb` in string"

# 1-terminal: slog, config'da LogValue hali yoʻq
$ go build ./cmd/api && ./api -port 8782 > api.log
# 2-terminal:
$ cat api.log
{"time":"2026-11-28T18:18:34.317527911+05:00","level":"INFO","msg":"config","config":{}}
{"time":"2026-11-28T18:18:34.324389331+05:00","level":"INFO","msg":"bozor API ishga tushdi","addr":":8782"}

# LogValue bilan (vaqt bundan keyin "…" bilan qisqartirilgan):
{"time":"…","level":"INFO","msg":"config","config":{"addr":":8782","db":"localhost:5582/bozor"}}

$ grep -n Info cmd/api/vet_demo.go          # namuna: qiymatsiz kalit
6:	logger.Info("soʻrov", "status")
$ go vet ./cmd/api
cmd/api/vet_demo.go:6:2: call to slog.Logger.Info missing a final value

# request_id bor, writeError hali err'ni toʻliq yozadi — SIZISH:
$ curl -s -o /dev/null -w '%{http_code} %header{x-request-id}\n' localhost:8782/login -H 'Content-Type: application/json' -d '{"email":"nodira@bozor.test","password":"x"}'
401 F4JJID26UEPWPAKXR5OPRW7GUG
$ grep F4JJID26UEPWPAKXR5OPRW7GUG api.log
{"time":"…","level":"WARN","msg":"xato","request_id":"F4JJID26UEPWPAKXR5OPRW7GUG","status":401,"err":"email yoki parol notoʻgʻri: email \"nodira@bozor.test\": foydalanuvchi topilmadi"}
{"time":"…","level":"INFO","msg":"soʻrov","request_id":"F4JJID26UEPWPAKXR5OPRW7GUG","method":"POST","path":"/login","status":401,"dur":121904087}
$ curl -s -o /dev/null -w '%{http_code} %header{x-request-id}\n' localhost:8782/login -H 'Content-Type: application/json' -d '{"email":"jasur@bozor.test","password":"Samarqand\ulugbek1394"}'
400 HN7DXON3WSCNY24BLHN2XIWXPI
$ grep HN7DXON3WSCNY24BLHN2XIWXPI api.log
{"time":"…","level":"WARN","msg":"xato","request_id":"HN7DXON3WSCNY24BLHN2XIWXPI","status":400,"err":"body JSON sifatida oʻqilmadi: invalid escape sequence `\\ulugb` in string"}
{"time":"…","level":"INFO","msg":"soʻrov","request_id":"HN7DXON3WSCNY24BLHN2XIWXPI","method":"POST","path":"/login","status":400,"dur":153598}

# tuzatishdan keyin (papkaning oʻzi), oʻsha ikki soʻrov:
$ grep WARN api.log
{"time":"…","level":"WARN","msg":"mijoz xatosi","request_id":"…","status":401,"detail":"email yoki parol notoʻgʻri"}
{"time":"…","level":"WARN","msg":"mijoz xatosi","request_id":"…","status":400,"detail":"body bitta JSON obyekt boʻlsin, faqat maʼlum kalitlar bilan"}
$ grep -c -e nodira -e ulugb api.log
0
$ jq -c 'select(.level == "WARN") | {msg, status, detail}' api.log    # videoda yoʻq: maydon boʻyicha filtr
{"msg":"mijoz xatosi","status":401,"detail":"email yoki parol notoʻgʻri"}
{"msg":"mijoz xatosi","status":400,"detail":"body bitta JSON obyekt boʻlsin, faqat maʼlum kalitlar bilan"}
```

Nega decode xatosida parol harflari: JSON satrida `\` maxsus belgi boshlanishi, `\u` dan keyin esa toʻrtta oʻn oltilik
raqam kutiladi (RFC 8259 §7). Parolni JSON'ga qoʻlda (qochirmasdan) qoʻygan mijoz `\ulugbek…` yuboradi, `encoding/json`
xatosi esa notoʻgʻri ketma-ketlikni oʻzi koʻrsatadi: ``invalid escape sequence `\ulugb` in string``. Xato matnida qaysi
bayt boʻlishini biz tanlamaymiz — shuning uchun 4xx'da `err` emas, faqat biz yozgan `detail` log'ga boradi.

`dur` — `time.Duration`, JSON handler uni butun son, **nanosekundda** yozadi (`43463688` ≈ 43 ms: notanish email uchun
ham argon2id bir marta ishlaydi, GB-6.04).

Tugagach: serverni toʻxtating (`Ctrl+C`), `docker compose down -v`.

## Papka

```
gb-8-02-slog/
  cmd/api/main.go             asosiy logger: slog.NewJSONHandler(os.Stdout, nil) → newServer(…, logger)
  cmd/api/middleware.go       requestID(logger, next), loggerKey, loggerFrom, logRequests, recoverPanic
  cmd/api/problem.go          writeError(w, r, err): 4xx → Warn + detail, 5xx → Error + err
  cmd/api/config.go           LogValue (slog), String (fmt), dbHost
  cmd/api/log_test.go         videodagi qoidalar (Postgres'siz): go test ./cmd/api -run Slog
  cmd/api/tekshiruv_test.go   vazifa, 1-qadam (Postgres): go test ./cmd/api -run LoginLog
  internal/store/tekshiruv_test.go   vazifa, 2-qadam (Postgres'siz): go test ./internal/store -run LogValue
  cmd/api/…                   GB-8.01'dagi API (handler'lar, testlar — haqiqiy Postgres'da)
  compose.yaml                Postgres 18, 127.0.0.1:5582
  solution/                   javob (oʻz go.mod'i bilan, alohida ishlaydi)
  _video/                     chiqishlar.sh → chiqishlar.txt
```

## Vazifa

[`TASK.md`](TASK.md): muvaffaqiyatli login uchun `"login"` qatori va `store.User`ga `LogValue` (faqat `id` va `role`,
email ham, parol hash'i ham yoʻq). Tekshirish: `go test ./internal/store -run LogValue` (Postgres kerak emas), keyin
`docker compose up -d --wait && go test ./cmd/api -run LoginLog`; butun papka: `go test ./...`.

## Manbalar

- [Go `log/slog`](https://pkg.go.dev/log/slog) — `Logger`, `Logger.With` («includes the given attributes in each output
  operation»), [`JSONHandler.Handle`](https://pkg.go.dev/log/slog#JSONHandler.Handle) («formats its argument Record as a
  JSON object on a single line»; qiymatlar `encoding/json` kabi, `error` — `Error()` matni),
  [`Level`](https://pkg.go.dev/log/slog#Level) (Debug −4, Info 0, Warn 4, Error 8),
  [`HandlerOptions`](https://pkg.go.dev/log/slog#HandlerOptions) («If Level is nil, the handler assumes LevelInfo»),
  [`LogValuer`](https://pkg.go.dev/log/slog#LogValuer), [`Default`](https://pkg.go.dev/log/slog#Default),
  [`DiscardHandler`](https://pkg.go.dev/log/slog#DiscardHandler); [Go 1.21 release notes: log/slog](https://go.dev/doc/go1.21#slog);
  [Go blog: Structured Logging with slog](https://go.dev/blog/slog)
- [`go vet`: slog analyzer](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/slog) — juft boʻlmagan kalit-qiymat
- [Go `context.WithValue`](https://pkg.go.dev/context#WithValue) — «Use context Values only for request-scoped data»; kalit
  turi oʻzimizniki
- [OWASP Logging Cheat Sheet: Data to exclude](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html#data-to-exclude)
  — parollar, session ID'lari, shaxsiy maʼlumot (PII) log'ga toʻgʻridan-toʻgʻri yozilmaydi: olib tashlanadi, maskalanadi yoki
  hash qilinadi; autentifikatsiya hodisalari esa yoziladi (vazifa)
- [The Twelve-Factor App: XI. Logs](https://12factor.net/logs) — ilova log'ni stdout'ga yozadi, uni qayerga yuborishni muhit hal qiladi
- [RFC 8259 §7](https://www.rfc-editor.org/rfc/rfc8259#section-7) — JSON satridagi `\` va `\uXXXX`
- [RFC 9457 §3.1.4](https://www.rfc-editor.org/rfc/rfc9457#section-3.1.4) — `detail`: mijoz uchun yozilgan tushuntirish
- [curl `-w`](https://curl.se/docs/manpage.html#-w) (`%{http_code}`, `%header{name}`, curl 7.84+),
  [GNU grep](https://www.gnu.org/software/grep/manual/grep.html) (`-c`, `-e`, `-n`), [jq manual](https://jqlang.org/manual/)
  (`select`, obyekt yasash), [Bash: Redirecting Output](https://www.gnu.org/software/bash/manual/bash.html#Redirecting-Output)
- Hamma chiqishlar shu kompyuterda: go1.27.4 linux/amd64, `CGO_ENABLED=0`, PostgreSQL 18 (`postgres:18-alpine`),
  `_video/chiqishlar.sh`.
