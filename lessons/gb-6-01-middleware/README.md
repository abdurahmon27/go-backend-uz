# GB-6.01 · Middleware: handler'ni oʻrab, log va recovery qoʻshamiz

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #43"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-5-10-buyurtma-api`](../gb-5-10-buyurtma-api) (buyurtmalar API, haqiqiy Postgres'da testlar) · keyingi
dars: `gb-6-02` (bir xil API xatolari). Bu **6-modulning (foydalanuvchilar va xavfsizlik) birinchi darsi**.

**Bu darsdan keyin:** `func(next http.Handler) http.Handler` shaklidagi middleware yozasiz va uni zanjirning toʻgʻri
joyiga qoʻyasiz: har soʻrovga log qatori (metod, yoʻl, status, vaqt), handler'dagi panic oʻrniga `500` va bizning JSON
xato shaklimiz, har soʻrovga `X-Request-Id`. Tartib nega muhimligini real natijada koʻrasiz: recovery log'dan tashqarida
boʻlsa, panic boʻlgan soʻrov log'dan yoʻqoladi. Vazifa: `requireJSON` (415) va 🟡 mux'ning 404/405'ini JSON qilish.

**Nega store xotirada.** Kod — [`gb-4-08-katalog-api`](../gb-4-08-katalog-api)'ning `solution/`i (xotiradagi katalog,
`bozor.json`), 5-modulning Postgres papkasi emas — ataylab: bu darsning mavzusi Postgres emas, middleware, va papka
Docker'siz `go test` bilan tekshirilsin. Middleware faqat `http.Handler`ni koʻradi. Shuning uchun
`cmd/api/middleware.go` [`gb-5-10`](../gb-5-10-buyurtma-api)'dagi Postgres versiyasiga ham oʻzgarishsiz qoʻshiladi
(`newServer`da `Handler: newHandler(st)`; tekshirildi: `go vet ./cmd/api` va `go build ./cmd/api` oʻtadi).

## Papka

```
gb-6-01-middleware/
  go.mod                       module github.com/abdurahmon27/bozor, go 1.27, tashqi modul yoʻq
  bozor.json                   katalog (GB-4.08): olma, non, choy, oʻrik
  cmd/api/middleware.go        YANGI: newHandler (zanjir), requestID, statusRecorder, logRequests, recoverPanic,
                               vazifa uchun requireJSON va 🟡 muxErrors (hozircha hech narsa qilmaydi)
  cmd/api/middleware_test.go   YANGI: log qatori, panic → 500, request ID, tartib, ErrAbortHandler, log injection
  cmd/api/main.go              newServer: Handler endi newHandler(st)
  cmd/api/api_test.go          GB-4.08 testlari; send body boʻlsa Content-Type: application/json qoʻyadi
  cmd/api/handlers.go, json.go GB-4.08 dagidek (DELETE /products/{id} bilan)
  cmd/api/tekshiruv_test.go    vazifa tekshiruvi; oʻzgartirmang
  cmd/bozor/, internal/store/  GB-4.08 dagidek (GB-4.08 vazifasining tekshiruvi endi store/delete_test.go)
  TASK.md                      vazifa
  solution/                    javob (oʻz go.mod'i bilan)
  _video/                      videodagi variantlar (v0…v7/middleware.go) va run.sh → chiqishlar.txt
```

GB-4.08 bilan farq: `diff -ru ../gb-4-08-katalog-api/solution . -x _video -x solution`.

## Zanjir

```go
func newHandler(st *store.Store) http.Handler {
	// tashqaridan ichkariga: requestID → logRequests → recoverPanic → mux
	return requestID(logRequests(recoverPanic(newMux(st))))
}
```

- **Qurish** ichkaridan: Go avval argumentlarni hisoblaydi — `newMux`, keyin `recoverPanic`, `logRequests`, eng oxiri
  `requestID` chaqiriladi. **Soʻrov** esa tashqaridan kiradi: avval `requestID` qaytargan funksiya, javob ichkaridan qaytadi.
- `logRequests` log qatorini `next.ServeHTTP`dan **keyin** yozadi. Panic shu chaqiruvdan uchib oʻtsa, bu qator bajarilmaydi —
  shuning uchun `recoverPanic` `logRequests`ning **ichida** turadi.
- Videodagi `statusRecorder` — `http.ResponseWriter`ni **embedding** bilan oʻraydi: `Header` va `Write` ichki writer'niki, `WriteHeader`ni
  oʻzimiz yozdik. `status` `http.StatusOK`dan boshlanadi: `WriteHeader` chaqirmay `Write` qilgan handler (masalan
  `healthHandler`) uchun javob 200. Videodagi variantda `Write` ichki writer'ga oʻtadi — boshlangʻich qiymatsiz logda `0` chiqadi.
- `requestID` — `crypto/rand.Text()` (Go 1.24+; kamida 128 bit tasodif, 26 belgili base32 matn) → javobning `X-Request-Id`
  sarlavhasi. `logRequests` va `recoverPanic` ID'ni `w.Header().Get("X-Request-Id")` bilan oʻqiydi (sarlavhalar bitta map,
  `statusRecorder`ning `Header()`i ham embedding orqali oʻsha map). Mijoz yuborgan `X-Request-Id`ga ishonmaymiz. Handler
  ichidagi log'larga ID'ni `context` orqali yetkazish va `slog` — GB-8.02.
- `recoverPanic` — GB-2.09'dagi `safely`, har soʻrov uchun: `recover()`, `debug.Stack()` logga, mijozga `500` va
  `{"error":{"code":"internal","message":"ichki xato"}}` (panic matni mijozga chiqmaydi).
- **`recoverPanic`: videodan tashqari ikki himoya** (papkadagi kodda bor, videoda vaqt uchun koʻrsatilmagan):
  1. `http.ErrAbortHandler` — net/http'ning "javobni ataylab uz" degan panic qiymati (`go doc net/http.ErrAbortHandler`):
     uni toʻxtatmaymiz, `panic(v)` bilan net/http'ga qaytaramiz.
  2. **Javob boshlangandan keyingi panic.** Handler status yoki body'ni yozib boʻlgach panic qilsa, 500'ni yuborib
     boʻlmaydi: oddiy versiyada mijoz `200` va `[{"id":1}{"error":…}` kabi buzuq body olardi, access log esa `500`
     derdi. Shuning uchun `recoverPanic` writer'ni `startWriter` bilan oʻraydi (`WriteHeader`/`Write` chaqirilganini eslaydi)
     va javob boshlangan boʻlsa, panic'ni ID bilan logga yozib, `panic(http.ErrAbortHandler)` bilan javobni uzadi
     (`TestPanicAfterWrite`). Bu holda `logRequests` qatori yozilmaydi; panic qatori logda bor.
- Handler oʻzi ishga tushirgan **goroutine**dagi panic'ni hech bir middleware ushlamaydi (recover faqat panic qilgan
  goroutine'da ishlaydi): u butun serverni yiqitadi. Bunday goroutine ichida oʻz `defer`/`recover`i boʻlsin (7-modul).
- Panic qiymati logga `%v` bilan yoziladi; qiymatda mijoz matni boʻlsa, u ham yangi qator olib kirishi mumkin —
  yoʻldagidek `%q` bilan yozish (`fmt.Sprint(v)`ni `%q`ga berish) xavfsizroq.
- `statusRecorder.Unwrap()` (va `startWriter`, `jsonErrorWriter`ning `Unwrap`i) ichki writer'ni qaytaradi:
  `http.NewResponseController(w)` (Flush, deadline'lar) oʻralgan writer orqali ham ishlashi uchun
  (`go doc net/http.NewResponseController`). Oʻralgan writer'ga `w.(http.Flusher)` kabi tekshiruv esa endi `false` beradi:
  shuning uchun `ResponseController` ishlating.
- **Videodagi `statusRecorder` soddaroq:** har `WriteHeader`da `status`ni almashtiradi, `Write` esa ichki writer'ga oʻtadi.
  Papkadagi variant birinchi yakuniy statusni saqlaydi: aniq `WriteHeader` yoki `Write`ning yashirin 200'i.
  Qayta tekshiruv (`go1.27.4`, `CGO_ENABLED=0`): `WriteHeader(201); WriteHeader(500)` — mijozga 201, eski recorder'da 500;
  `Write([]byte("ok\n")); WriteHeader(500)` — mijozga 200, eski recorder'da 500. Tuzatilganda log ham mos ravishda 201 va 200.
  Ikkinchi `WriteHeader` ichki writer'ga yetib boradi: net/http odatdagi `superfluous response.WriteHeader call`
  ogohlantirishini saqlaydi. `Write` ichki writer'ning Content-Type aniqlashini ham saqlaydi.
  `TestLogRequestsFirstStatus` ikkala yoʻlni tekshiradi. `_video/v2-status`, `v3-nol`, `v4-tartib` — videodagi tarixiy variantlar;
  yakuniy `v7-s` esa papkadagi kod bilan bir xil, faqat yoʻl `%s` orqali yoziladi.
- Javobda (`solution/`) `requireJSON` zanjirda emas, `newMux`da faqat `POST /products` va `POST /orders`ni oʻraydi: zanjir
  boshida tursa, `POST /health` yoki yoʻq yoʻlga ham 405/404 oʻrniga 415 berardi.
- Yoʻl logga `%q` bilan yoziladi: `%0A` (URL'dagi yangi qator) `r.URL.Path`da haqiqiy yangi qatorga ochiladi; `%s` bilan
  mijoz logga soxta qator qoʻsha oladi (log injection, OWASP).

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, curl 8.22.0)

Hammasi `_video/run.sh` dan (`./run.sh > chiqishlar.txt 2>&1`): har variant papkaning nusxasida `cmd/api/middleware.go`ni
`_video/vN-…/middleware.go` bilan almashtiradi, kerak boʻlsa `getProduct`ga vaqtincha panic qatorlarini qoʻshadi, serverni
`:8765`da ishga tushirib curl bilan soʻraydi va toʻxtatadi. Toʻliq chiqishlar: [`_video/chiqishlar.txt`](_video/chiqishlar.txt).
Stack'larda `…` — videoda qisqartirilgan qatorlar. Quyidagi sonlar videodan iqtibos; `chiqishlar.txt`dagi v5–v7
statusRecorder tuzatilgandan keyin qayta yozilgan (2026-11-27).

Vaqtincha panic (`getProduct`, `pathID`dan keyin; videoning oxirida olib tashlanadi):

```go
	if id == 13 { // vaqtincha: bug oʻrnida
		panic("sinov: 13-mahsulotda bug")
	}
```

**v0 — middleware yoʻq.** Mijoz javobsiz qoladi, server esa tirik:

```
$ curl -i localhost:8765/products/13
curl: (52) Empty reply from server
$ curl -s localhost:8765/products/1
{"id":1,"name":"olma","price":1200000,"stock":40}
```

Server logi (net/http panic'ni oʻzi recover qiladi, stack'ni yozadi va ulanishni yopadi — `go doc net/http.Handler`):

```
2026/11/27 21:45:08 http: panic serving [::1]:55184: sinov: 13-mahsulotda bug
goroutine 8 [running]:
net/http.(*conn).serve.func1()
…
main.(*api).getProduct(…)
…
net/http.(*ServeMux).ServeHTTP(…)
…
```

**v1 — `logRequests`, status yoʻq:**

```
2026/11/27 21:45:09 GET "/products/1" 17.938µs
2026/11/27 21:45:09 GET "/products/99" 157.936µs
```

**v3 — `statusRecorder`, `status` boshlangʻich qiymatsiz** (`rec := &statusRecorder{ResponseWriter: w}`):

```
2026/11/27 21:45:09 GET "/health" 0 13.35µs
2026/11/27 21:45:09 GET "/products/99" 404 140.397µs
```

**v2 — `status: http.StatusOK` bilan:**

```
2026/11/27 21:45:09 GET "/products/99" 404 162.28µs
2026/11/27 21:45:09 GET "/health" 200 11.725µs
2026/11/27 21:45:09 GET "/narsalar" 404 10.352µs
```

**v4 — notoʻgʻri tartib `recoverPanic(logRequests(newMux(st)))`.** Mijoz 500 oladi, lekin log qatori yoʻq:

```
$ curl -i localhost:8765/products/13
HTTP/1.1 500 Internal Server Error
Content-Type: application/json
…
{"error":{"code":"internal","message":"ichki xato"}}
```

```
2026/11/27 21:45:10 panic: sinov: 13-mahsulotda bug
goroutine 5 [running]:
…
main.logRequests.func1(…)
…
main.recoverPanic.func1(…)
…
```

**v5 — yakuniy zanjir (shu papka) + vaqtincha panic:**

```
$ curl -i localhost:8765/products/13
HTTP/1.1 500 Internal Server Error
Content-Type: application/json
X-Request-Id: ZNTHHUQNT2T72EFC7GLZFXVT5N
…
{"error":{"code":"internal","message":"ichki xato"}}
```

```
2026/11/27 21:45:10 ZNTHHUQNT2T72EFC7GLZFXVT5N panic: sinov: 13-mahsulotda bug
…
main.recoverPanic.func1(…)
…
main.logRequests.func1(…)
…
main.requestID.func1(…)
…
2026/11/27 21:45:10 ZNTHHUQNT2T72EFC7GLZFXVT5N GET "/products/13" 500 225.507µs
```

**v6 / v7 — log injection.** `%0A` — yangi qator, `%20` — boʻsh joy (v7 = shu papkaning `middleware.go`i, faqat
`logRequests`da `%q` → `%s`):

```
$ curl -s 'localhost:8765/products/1%0A2026/11/27%2021:00:00%20soxta-qator'
404 page not found
```

`%s` bilan (v7) logda ikkinchi, soxta qator:

```
2026/11/27 21:45:11 6ETS6UDWEJUSLDW2LUE74Q3I3H GET /products/1
2026/11/27 21:00:00 soxta-qator 404 13.918µs
```

`%q` bilan (shu papka) — bitta qator:

```
2026/11/27 21:45:11 322Q3G3WU6NE5N3RTLIGUN66TC GET "/products/1\n2026/11/27 21:00:00 soxta-qator" 404 33.138µs
```

**Testlar:**

```
$ go test ./...                         # starter: vazifa tekshiruvi yiqiladi (TASK.md)
--- FAIL: TestVazifaRequireJSON (0.00s)
…
FAIL	github.com/abdurahmon27/bozor/cmd/api	0.004s
ok  	github.com/abdurahmon27/bozor/cmd/bozor
ok  	github.com/abdurahmon27/bozor/internal/store
$ (cd solution && go test ./...)
ok  	github.com/abdurahmon27/bozor/cmd/api	0.003s
ok  	github.com/abdurahmon27/bozor/cmd/bozor
ok  	github.com/abdurahmon27/bozor/internal/store
```

## Vazifa

[`TASK.md`](TASK.md): `requireJSON` — `POST`da `Content-Type` `application/json` boʻlmasa `415 unsupported_media_type`,
zanjirning shunday joyida turadiki, 415 ham `X-Request-Id` oladi va logga tushadi; 🟡 `muxErrors` — mux'ning 404/405 javoblari JSON.
Tekshirish: `go test ./...`. Javob: [`solution/`](solution/).

## Manbalar

- `net/http.Handler` (panic: server recover qiladi, stack'ni logga yozadi, ulanishni yopadi; `ErrAbortHandler`):
  https://pkg.go.dev/net/http#Handler · https://pkg.go.dev/net/http#ErrAbortHandler
- `net/http.Server.Serve` (har ulanishga yangi goroutine): https://pkg.go.dev/net/http#Server.Serve
- `net/http.HandlerFunc`, `ResponseWriter` (birinchi `Write` → `WriteHeader(200)`): https://pkg.go.dev/net/http#ResponseWriter
- `net/http.NewResponseController` va `Unwrap`: https://pkg.go.dev/net/http#NewResponseController
- `net/http.ServeMux.Handler` (🟡): https://pkg.go.dev/net/http#ServeMux.Handler
- `crypto/rand.Text`: https://pkg.go.dev/crypto/rand#Text · Go 1.24 release notes: https://go.dev/doc/go1.24#cryptorandpkgcryptorand
- `runtime/debug.Stack`: https://pkg.go.dev/runtime/debug#Stack · `mime.ParseMediaType`: https://pkg.go.dev/mime#ParseMediaType
- Go spec: embedded fields https://go.dev/ref/spec#Struct_types · Handling panics https://go.dev/ref/spec#Handling_panics ·
  Calls (argumentlar avval) https://go.dev/ref/spec#Calls · Effective Go: Embedding https://go.dev/doc/effective_go#embedding
- RFC 9110 §15.5.16 415 Unsupported Media Type: https://www.rfc-editor.org/rfc/rfc9110#name-415-unsupported-media-type
- OWASP Log Injection: https://owasp.org/www-community/attacks/Log_Injection · OWASP Logging Cheat Sheet:
  https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html
- curl xato kodlari (52 = empty reply): https://curl.se/libcurl/c/libcurl-errors.html
