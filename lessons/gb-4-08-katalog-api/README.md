# GB-4.08 · Nazorat nuqtasi: katalog API (CRUD, testlar, `-race` toza)

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #32"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-4-07-mutex-race`](../gb-4-07-mutex-race) (soʻrovlar bir vaqtda: `sync.Mutex`, `-race`) · keyingi dars:
`gb-5-01` (nega baza kerak: ikki jarayon bitta JSON faylni buzadi). Bu papka **4-modulning nazorat nuqtasi** va
**5-modulning boshlanish nuqtasi**: Postgres darslari shu `cmd/api` + `internal/store`dan davom etadi.

**Bu darsdan keyin:** 4-modul boʻlaklarini — server va `ReadHeaderTimeout` ([`gb-4-04`](../gb-4-04-birinchi-server)),
pattern'lar va `{id}` ([`gb-4-05`](../gb-4-05-route-id)), JSON kirish/chiqish va bitta xato shakli ([`gb-4-06`](../gb-4-06-json-body)),
qulf va `-race` ([`gb-4-07`](../gb-4-07-mutex-race)), `httptest` ([`gb-4-03`](../gb-4-03-handler-httptest)) — GB-3.06'dagi
katalog ([`gb-3-06-bozor-cli`](../gb-3-06-bozor-cli), `internal/store`) ustida bitta API'ga yigʻasiz; store xatolarini HTTP
statusga bitta joyda aylantirasiz; har route'ni jadvalli test bilan tekshirasiz; qulf boʻlsa ham koʻrsatkich qulfdan
tashqariga chiqqani uchun yuzaga kelgan data race'ni `go test -race` bilan topib tuzatasiz. Vazifa: `DELETE /products/{id}`
endpoint'ini qoʻllanmasiz oʻzingiz qoʻshasiz.

## Papka

```
gb-4-08-katalog-api/
  go.mod                         module github.com/abdurahmon27/bozor, go 1.27, tashqi modul yoʻq
  bozor.json                     katalog: CLI (cmd/bozor) yozgan fayl, API ishga tushganda oʻqiydi
  cmd/bozor/                     GB-3.06'dagi CLI, oʻzgarmagan (add, list, order, restock)
  cmd/api/main.go                newServer (Handler: newMux(st), ReadHeaderTimeout) va main (store.Load → ListenAndServe)
  cmd/api/handlers.go            api struct, newMux (route'lar), handler'lar, pathID, writeStoreError
  cmd/api/json.go                GB-4.06: apiError, writeJSON, writeError, decodeJSON, writeDecodeError
  cmd/api/api_test.go            har route uchun jadvalli test (TestRoutes), TestParallel, TestServer …
  cmd/api/tekshiruv_test.go      vazifa tekshiruvi (2- va 3-qadam); oʻzgartirmang
  internal/store/store.go        katalog: Product, xatolar, Store{mu, products}, Add, All, Get, Sell, Restock, Delete (vazifa)
  internal/store/file.go         Load va atomar Save (GB-3.06)
  internal/store/*_test.go       store testlari (GB-3.06) + tekshiruv_test.go (vazifa, 1-qadam)
  TASK.md                        vazifa
  solution/                      javob (oʻz go.mod'i bilan, alohida ishlaydi)
  _video/                        videodagi qadamlar: run.sh, race.sh, 01-pointer/ (bug'li Get), record.sh → chiqishlar.txt
```

GB-3.06 bilan farqni koʻrish: `diff -ru ../gb-3-06-bozor-cli/solution .` (yangi: `cmd/api/`, `bozor.json`; `store.go`'da
`mu sync.Mutex`, har metodda `Lock`/`defer Unlock`, `All` nusxa ustida, yangi `Get`, vazifa uchun `Delete`; `file.go`'da
`Save` nusxani qulf ostida oladi; `validate`'da narx va qoldiqning yuqori chegarasi — `maxPrice`, `maxStock`, pastda).

## Route'lar

| Soʻrov | Javob |
|---|---|
| `GET /health` | `200`, `ok` (matn) |
| `GET /products` | `200`, katalog JSON massiv (boʻsh katalog: `[]`) |
| `GET /products/{id}` | `200` mahsulot · `400 bad_id` · `404 not_found` |
| `POST /products` | `201` yangi mahsulot (ID'ni server beradi) + `Location: /products/{id}` · `400 bad_json` · `413 too_large` · `422 validation_failed` |
| `POST /orders` | `201` `{"product_id","quantity","total"}` · `400` · `404 not_found` · `409 insufficient_stock` · `422` |
| `DELETE /products/{id}` | **vazifa:** `204 No Content` (body yoʻq) · `400 bad_id` · `404 not_found` |

Hamma xato bitta shaklda: `{"error": {"code": …, "message": …, "fields": …}}` (GB-4.06). Store xatosi statusga
**faqat** `writeStoreError`da aylanadi: `ErrNotFound` → 404, `*StockError` → 409, `ErrBadQuantity`/`ErrBadProduct` → 422,
qolgan hammasi → `500` va `"ichki xato"` (kutilmagan xatoning matni — masalan fayl yoʻli — mijozga chiqmaydi, faqat server
logiga `log.Printf` bilan yoziladi; tuzilgan log `slog` — GB-8.02).
Narx va jami summa **tiyinda** (GB-1.03): `3500000` = 35 000 soʻm.

**Yuqori chegaralar.** GB-3.06'da narxning yuqori chegarasi (10 mln soʻm) faqat CLI'da edi (`parsePrice`). API'ga esa istalgan son
kelishi mumkin: `"price": 9000000000000000000` va 2 dona buyurtma `int64`dan oshib, manfiy `total` berardi (GB-1.03: toʻlib
ketish). Shuning uchun chegaralar store'ga koʻchdi: `maxPrice` (1 000 000 000 tiyin = 10 mln soʻm) va `maxStock` (1 000 000 dona) —
`validate` va `Restock` tekshiradi, CLI ham, API ham shu qoidadan oʻtadi. Buyurtmada `quantity` ≤ 100 (`validateOrder`), demak
`price × quantity` ≤ 10¹¹ — `int64`ga bemalol sigʻadi. Buyurtma javobi `orderResponse` (GB-4.06'dagi `Order`): `Order` nomi
bozor'ning domen turi uchun saqlanadi (5-modul).

**Mux'ning oʻz javoblari oddiy matn boʻlib qoladi.** Hech bir pattern'ga mos kelmagan yoʻl — `404 page not found`
(`text/plain`), yoʻl bor-u metod yoʻq — `405 Method Not Allowed` + `Allow` (GB-4.05). Status va `Allow` toʻgʻri; ularni
ham JSON qilish uchun mux'ni oʻrab oladigan qatlam (middleware) kerak — GB-6.01. Metodsiz `"/"` pattern'i bilan "hamma
narsani tutish" yoʻli notoʻgʻri: u 405'ni ham yutib, 404 qilib yuboradi (GB-4.05 "`/` hammasini oladi").

**Katalog xotirada.** `main` ishga tushganda `bozor.json`ni `store.Load` bilan oʻqiydi; API'dagi oʻzgarishlar (yangi mahsulot,
buyurtma) faylga **yozilmaydi** va server toʻxtasa yoʻqoladi. Bu ataylab: bitta faylni ikki jarayon (CLI va API) yozsa nima
boʻlishi — GB-5.01'ning mavzusi, doimiy saqlash — Postgres (5-modul).

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, curl 8.22.0)

Papkaning tuzilishi:

```
$ ls cmd/api internal/store
cmd/api:
api_test.go  handlers.go  json.go  main.go

internal/store:
file.go  restock_test.go  store.go  store_test.go
```

(Videoda `tekshiruv_test.go` fayllari koʻrsatilmagan: ular vazifaniki.) Birinchi terminalda server:

```
$ go run ./cmd/api
bozor API: http://localhost:8765
```

Ikkinchi terminalda (`-s` — ortiqcha chiqishsiz, `-i` — status qatori va sarlavhalar, `-d` — body va POST, `-X` — metod;
`…` — videoda qisqartirilgan sarlavhalar: `Date`, `Content-Length`, `X-Content-Type-Options`):

```
$ curl -s localhost:8765/products/4
{"id":4,"name":"oʻrik","price":3500000,"stock":2}

$ curl -i localhost:8765/orders -d '{"product_id":4,"quantity":5}'
HTTP/1.1 409 Conflict
Content-Type: application/json
…
{"error":{"code":"insufficient_stock","message":"oʻrik: qoldiq 2, soʻralgan 5"}}

$ curl -i localhost:8765/orders -d '{"product_id":4,"quantity":2}'
HTTP/1.1 201 Created
…
{"product_id":4,"quantity":2,"total":7000000}

$ curl -s localhost:8765/products/4
{"id":4,"name":"oʻrik","price":3500000,"stock":0}

$ curl -i localhost:8765/narsalar
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
…
404 page not found

$ curl -i -X DELETE localhost:8765/products/1
HTTP/1.1 405 Method Not Allowed
Allow: GET, HEAD
…
Method Not Allowed
```

`-H "Content-Type: application/json"` yuborilmadi: bu server `Content-Type`ni tekshirmaydi, body'ni baribir JSON deb oʻqiydi
(GB-4.01'dagi namunaviy server tekshirgani uchun `415` bergan edi; tekshirish — kelishuv masalasi). Serverni birinchi terminalda
`Ctrl+C` bilan toʻxtating (GB-4.04). Port band boʻlsa: `listen tcp :8765: bind: address already in use` — boshqa terminalda
eski server ishlayapti. `go run ./cmd/api`ni papkaning ildizidan ishga tushiring: boshqa papkada `bozor.json` topilmaydi va
katalog jimgina boʻsh boshlanadi (`Load`: fayl yoʻq = boʻsh katalog, GB-3.06). Videoda koʻrsatilmagan yana uch soʻrov:

```
$ curl -s localhost:8765/products -d '{"name":"","price":-5}'
{"error":{"code":"validation_failed","message":"qiymatlar notoʻgʻri","fields":{"name":"boʻsh boʻlmasin","price":"0 dan katta boʻlsin"}}}
$ curl -s localhost:8765/products/abc
{"error":{"code":"bad_id","message":"id butun son boʻlishi kerak"}}
$ curl -i localhost:8765/products -d '{"name":"anor","price":1800000,"stock":7}'
HTTP/1.1 201 Created
Content-Type: application/json
Location: /products/5
…
{"id":5,"name":"anor","price":1800000,"stock":7}
```

Testlar (vazifa tekshiruvisiz — videoda shunday; `./_video/run.sh final go test ./...`):

```
$ go test ./...
ok  	github.com/abdurahmon27/bozor/cmd/api	0.003s
ok  	github.com/abdurahmon27/bozor/cmd/bozor	0.040s
ok  	github.com/abdurahmon27/bozor/internal/store	0.013s
```

## Bug: qulf bor, lekin koʻrsatkich qulfdan chiqib ketdi

Videoda `Get`ning birinchi versiyasi `*Product` qaytargan (`_video/01-pointer/internal/store/store.go`):

```go
func (s *Store) Get(id int64) (*Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.find(id)
	if p == nil {
		return nil, fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
	}
	return p, nil
}
```

Qulf bor, lekin `p` katalog slice'ining ichiga qaraydi (`&s.products[i]`, GB-2.04). `Get` qaytgach qulf boʻsh, `getProduct`
esa `writeJSON` → `Encode` orqali `p.Stock`ni **qulfsiz** oʻqiydi; shu paytda boshqa soʻrovning `Sell`i qulf ostida
`Stock`ga yozadi. Oddiy `go test ./...` oʻtadi (qoldiq ham toʻgʻri: yozuvlar qulf ostida), `go test -race` esa ushlaydi:

```
$ go test -race ./cmd/api          # ./_video/race.sh 01-pointer go test -race ./cmd/api
==================
WARNING: DATA RACE
Read at 0x00c00011f4c0 by goroutine 30:
  reflect.Value.Int()
      /usr/local/go/src/reflect/value.go:1461 +0x453
  encoding/json/v2.makeIntArshaler.func1()
…
  github.com/abdurahmon27/bozor/cmd/api.writeJSON()
      /home/dev/bozor/cmd/api/json.go:24 +0x1fb
  github.com/abdurahmon27/bozor/cmd/api.(*api).getProduct()
      /home/dev/bozor/cmd/api/handlers.go:61 +0xae
…
Previous write at 0x00c00011f4c0 by goroutine 31:
  github.com/abdurahmon27/bozor/internal/store.(*Store).Sell()
      /home/dev/bozor/internal/store/store.go:160 +0x6a4
  github.com/abdurahmon27/bozor/cmd/api.(*api).createOrder()
      /home/dev/bozor/cmd/api/handlers.go:143 +0x364
…
--- FAIL: TestParallel (0.00s)
    testing.go:1867: race detected during execution of test
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api	0.013s
FAIL
```

`reflect.Value.Int` — `encoding/json` struct maydonini shu orqali oʻqiydi (hozirgi Go'da `encoding/json` ichida
`encoding/json/v2` dvigateli ishlaydi, shuning uchun stack'da `v2` nomlari koʻrinadi; API — oʻsha `encoding/json`). Toʻliq
chiqish: `_video/chiqishlar.txt` (`./_video/record.sh` uni qaytadan yozadi). Bu testda race har safar topildi (5/5,
`chiqishlar.txt`): race detector shu ishga tushirishda **sodir boʻlgan** murojaatlarni koʻradi va ularning orasida
"kim kimdan oldin" (happens-before) munosabati yoʻqligini aniqlaydi; ikki murojaat aynan bir paytda boʻlishi shart emas.
Qaysi murojaat oldin sodir boʻlganiga qarab hisobot tartibi almashadi: bu safar `Read at …` / `Previous write at …`, boshqa
safar `Write at …` (Sell) / `Previous read at …` (getProduct) — mazmuni bir xil.

Tuzatish — `Get` **nusxa** qaytaradi (`Product`, `return *p, nil`), nusxa qulf ostida olinadi. Shu qoida `All`da ham:
qulf ostida `slices.Clone` bilan nusxa, `yield` esa qulfsiz (GB-4.07 Chuqurroq: `yield` qulf ostida boʻlsa va sikl tanasi
`Sell`ni chaqirsa — deadlock). `Save` ham nusxani qulf ostida olib, faylga qulfsiz yozadi. **Qoida: store'dan tashqariga
koʻrsatkich emas, nusxa chiqsin; `find` faqat qulf ostida chaqiriladi.**

```
$ go test -race ./...              # ./_video/race.sh final go test -race ./...
ok  	github.com/abdurahmon27/bozor/cmd/api	1.014s
ok  	github.com/abdurahmon27/bozor/cmd/bozor	1.046s
ok  	github.com/abdurahmon27/bozor/internal/store	1.017s
```

Har paketda ~1 s: bu sekinlik emas — race runtime chiqishdan oldin `atexit_sleep_ms=1000` kutadi (race detector maqolasi,
`GORACE` sozlamalari); `GORACE=atexit_sleep_ms=0 go test -race ./...` bilan ~0.01 s.

## `-race` uchun C kompilyator kerak

Linux'da (va Windows'da) race detector `cgo`ni, yaʼni C kompilyatorni talab qiladi (Linux'da odatda `gcc`); macOS'da
Go 1.20'dan beri cgo'siz ishlaydi. Kompilyator yoʻq boʻlsa: `go: -race requires cgo; enable cgo by setting CGO_ENABLED=1`. Videodagi `-race` natijalari shu sababli Docker'da olingan
(gcc bor `golang:1.26.8` image'i + shu mashinadagi `go1.27.4` toolchain, `_video/race.sh`), shuning uchun fayl yoʻllari
`/home/dev/bozor/…`. Linux'da `gcc` boʻlsa (`apt install gcc`), oddiy `go test -race ./...` yetadi.
Docker bilan (gcc'siz mashinada) — Go'ning rasmiy image'i:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.27 go test -race ./...
```

Manzillar (`0x00c0…`) va goroutine raqamlari har ishga tushirishda boshqacha.

## Vazifa

[`TASK.md`](TASK.md): `DELETE /products/{id}` — `store.Delete` (qulf bilan, `ErrNotFound`), route va `deleteProduct`
handler'i (`204 No Content`, `400`, `404`), oʻzingizning `TestDeleteProduct` testingiz. Tekshirish:

```bash
go test ./...            # hozir FAIL: tekshiruv_test.go nima yetishmayotganini oʻzbekcha aytadi
go test -race ./...      # gcc yoki Docker boʻlsa
```

Javob: [`solution/`](solution/) (`cd solution && go test ./...`).

## Chuqurroq

- **Nega `api` struct, global `catalog` emas.** GB-4.03…4.06 laboratoriyalarida katalog global oʻzgaruvchi edi. Global
  holat testlarni bir-biriga bogʻlaydi (biri buyurtma bersa, keyingisi oʻzgargan katalogni koʻradi) va `t.Parallel()` bilan
  testlarni parallel ishlatishni xavfli qiladi. Bu yerda har test `newMux(testStore(t))` bilan oʻz katalogini oladi. Handler'lar `*api` metodlari;
  `a.listProducts` — **metod qiymati** (Go spec: Method values): `a` bilan bogʻlangan oddiy funksiya, `http.HandlerFunc`
  turiga mos keladi.
- **Validatsiya ikki joyda.** `validate` (handler) maydonma-maydon xabar beradi (`fields`, GB-4.06), `store.Add` esa
  katalogning toʻliq qoidasini (UTF-8, nom ≤ 20 belgi, GB-3.06) yana tekshiradi va `ErrBadProduct` qaytaradi — u ham 422
  boʻladi (`TestRoutes/uzun_nom`). Qoidalarning manbai — store; handler'dagisi mijozga qulay xabar uchun.
- **Katalog ID'lari GB-3.06'dagidek:** olma 1, non 2, choy 3, oʻrik 4 (GB-4.05…4.07 laboratoriyalarida oʻrik 3 edi).
- **ID qayta berilishi.** `Add` yangi ID'ni "eng katta ID + 1" qiladi (GB-3.06). Vazifadagi `Delete` oxirgi mahsulotni
  oʻchirsa, keyingi `Add` xuddi shu ID'ni beradi: eski havola endi boshqa mahsulotni koʻrsatadi. Haqiqiy tizimda ID'lar
  qayta ishlatilmaydi — Postgres'da `GENERATED … AS IDENTITY` (5-modul) buni hal qiladi.
- **409 yoki 422?** Qoldiq yetmasligi — soʻrov toʻgʻri, lekin katalogning hozirgi holatiga zid: `409 Conflict` (RFC 9110
  §15.5.10). `422` — qiymatning oʻzi qoidaga zid (`quantity: 0`). Bu kursning kelishuvi; ba'zi API'lar qoldiq uchun ham 422 beradi.
- **`healthHandler`** matn qaytaradi (`ok`): yuklash balansiri va monitoring uchun shu yetadi (GB-8).
- **`ReadHeaderTimeout` yetarlimi?** U faqat sarlavhalarni oʻqishni cheklaydi (sekin mijoz hujumi, GB-4.04). Body uchun
  `ReadTimeout`, javob uchun `WriteTimeout`, boʻsh ulanishlar uchun `IdleTimeout` — server sozlamalari va toʻxtatish GB-8'da.

## Manbalar

- https://pkg.go.dev/net/http#ServeMux · #Server · #Request.PathValue · https://go.dev/blog/routing-enhancements
- https://pkg.go.dev/sync#Mutex · https://go.dev/doc/articles/race_detector · https://go.dev/ref/mem
- https://pkg.go.dev/slices#Clone · #Collect · #Delete · #IndexFunc · https://pkg.go.dev/iter
- https://go.dev/ref/spec#Method_values · https://pkg.go.dev/errors
- RFC 9110: https://www.rfc-editor.org/rfc/rfc9110#name-204-no-content · #name-409-conflict · #name-422-unprocessable-content · #name-500-internal-server-error
- https://go.dev/wiki/TableDrivenTests · https://pkg.go.dev/net/http/httptest
