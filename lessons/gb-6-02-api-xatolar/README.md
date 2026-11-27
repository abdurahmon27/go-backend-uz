# GB-6.02 · Bir xil API xatolari: mijozga xavfsiz, log'ga toʻliq

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #44"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-6-01-middleware`](../gb-6-01-middleware) (middleware: requestID, logRequests, recoverPanic, zanjir
tartibi) · keyingi dars: `gb-6-03` (parollar: argon2id, salt). Bu papka GB-6.01'ning `solution/`idan boshlanadi.

**Bu darsdan keyin:** bozor API'sining har xatosi mijozga bitta standart shaklda —
[RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) "problem details", `Content-Type: application/problem+json` — va
**biz tanlab yozgan** matn bilan boradi; toʻliq ichki xato (`err.Error()`) esa faqat server log'iga, GB-6.01'dagi
`X-Request-Id` bilan yoziladi. Xatoni status va xavfsiz matnga **bitta funksiya** — `problemFor` — `errors.Is`/`errors.As`
bilan aylantiradi (GB-2.08); tanilmagan xato xavfsiz `500` boʻladi.

**Katalog xotirada** (GB-4.08/GB-6.01'dagi `internal/store`): bu dars HTTP xato qatlami haqida, Postgres emas. GB-5.10'dagi
Postgres store oldiga qoʻyilsa, `problemFor`ga uning qoʻshimcha xatolari uchun yana uchta case kerak: `ErrOrderNotFound` (404,
detail "bunday ID'li buyurtma yoʻq"), `ErrDuplicate` (409) va `ErrInUse` (409); qoʻshilmasa ular xavfsiz `500`ga tushadi —
hech narsa sizmaydi, lekin status notoʻgʻri. Usul oʻsha: har yangi xato — bitta sentinel yoki tur va bitta case.

## Muammo: `err.Error()` javobda

GB-4.06'dagi `writeDecodeError` Decode xatosining matnini toʻgʻridan-toʻgʻri mijozga berardi:

```
$ curl -s localhost:8765/products -H "Content-Type: application/json" -d '{"name":"choy","price":"25000"}'
{"error":{"code":"bad_json","message":"json: cannot unmarshal string into Go struct field productInput.price of type int64"}}

$ curl -s localhost:8765/products -H "Content-Type: application/json" -d '[1]'
{"error":{"code":"bad_json","message":"json: cannot unmarshal array into Go value of type main.productInput"}}
```

Server begona mijozga oʻz ichki Go tur nomlarini (`productInput`, `int64`, `main.`) aytdi: hujumchi server nima bilan
yozilganini biladi. Bu xatolar sinfi — **CWE-209** (xato xabarida ichki maʼlumot). `err.Error()` — har `%w` qavatining matni
bitta qatorda (GB-2.08): tur nomi, fayl yoʻli, ertaga Postgres'ning jadval va constraint nomlari (pgx xatosi:
`… violates unique constraint "products_name_key" (SQLSTATE 23505)`). U **dasturchi uchun**: uni log'ga yozamiz. Mijozga —
faqat biz yozgan matn.

## Papka

```
gb-6-02-api-xatolar/
  go.mod                       module github.com/abdurahmon27/bozor, go 1.27, tashqi modul yoʻq
  bozor.json                   katalog (API ishga tushganda oʻqiladi, xotirada turadi)
  cmd/api/problem.go           YANGI: problem (RFC 9457), newProblem, writeProblem, errBadJSON, errBadID,
                               fieldErrors, productRules, problemFor (vazifa shu yerda), writeError
  cmd/api/json.go              writeJSON, decodeJSON (Decode xatosi errBadJSON bilan ikkita %w orqali oʻraladi);
                               GB-4.06'dagi apiError, eski writeError va writeDecodeError olib tashlandi
  cmd/api/handlers.go          handler'lar: writeStoreError(w, err) → writeError(w, err); validate → fieldErrors
  cmd/api/middleware.go        GB-6.01: requestID, logRequests, recoverPanic (+ startWriter), requireJSON (newMux'dagi
                               POST route'larda), muxErrors — oʻz xato javoblari (500, 415, 404/405) endi writeProblem orqali
  cmd/api/*_test.go            api_test.go (TestRoutes, TestProblemShape, TestInternalError …), middleware_test.go
  cmd/api/tekshiruv_test.go    vazifa tekshiruvi; oʻzgartirmang
  cmd/bozor/, internal/store/  oʻzgarmagan (GB-6.01)
  TASK.md                      vazifa
  solution/                    javob (oʻz go.mod'i bilan, alohida ishlaydi)
  _video/                      run.sh, record.sh → chiqishlar.txt; unutgan/ (return unutilgan handler), panic/ (recall)
```

GB-6.01 bilan farq: `diff -ru ../gb-6-01-middleware/solution . -x _video -x solution` (yangi `problem.go`; `json.go`,
`handlers.go`, `middleware.go` va testlar oʻzgargan).

## Javob shakli (RFC 9457)

| Aʼzo | Maʼnosi | bozor'da |
|---|---|---|
| `type` | xato turi, URI (§3.1.1) | maxsus turi yoʻq boʻlsa `"about:blank"`; oʻz turlarimiz toʻliq yoʻl bilan: `/problems/insufficient-stock`, `/problems/validation` |
| `title` | tur uchun qisqa nom, holatdan holatga oʻzgarmaydi (§3.1.3) | `about:blank`da — HTTP status nomi: `http.StatusText(status)` (§4.2.1) |
| `status` | HTTP status nusxasi (§3.1.2) | haqiqiy status bilan bir xil (`TestProblemShape` tekshiradi) |
| `detail` | shu holat haqida mijozga matn; "debugging maʼlumoti emas" (§3.1.4) | faqat biz tanlagan maydonlardan; `err.Error()` hech qachon |
| `errors` | kengaytma (§3.2) | `422`: maydon → nima notoʻgʻri (`fieldErrors`) |
| `request_id` | kengaytma (§3.2) | GB-6.01'dagi `X-Request-Id`: log'dagi qator bilan bogʻlaydi |

`instance` (§3.1.5) ishlatilmaydi. RFC 9457 (2023) eski RFC 7807'ning oʻrnini egallagan.

- **`type` URI.** RFC absolute URI'ni afzal koʻradi; nisbiy boʻlsa, toʻliq yoʻl bilan (`/problems/…`). U identifikator:
  ochiladigan sahifa boʻlishi shart emas (bozor'da hozircha yoʻq), lekin RFC keyinchalik hujjat sahifasini qoʻyishni tavsiya qiladi.
- **`title` va Go.** `about:blank`da title — statusning tavsiya etilgan nomi (§4.2.1). Go'ning `http.StatusText` 404/500 uchun
  aynan shu nomlarni beradi, lekin 413 va 422 uchun eski nomlarni saqlagan (`Request Entity Too Large`,
  `Unprocessable Entity`; RFC 9110'da `Content Too Large`, `Unprocessable Content`). Biz Go'nikini ishlatamiz; testlar shuni kutadi.
- **409'dagi sonlar `detail`da.** RFC: mijoz `detail`ni parse qilmasin (§3.1.4); dastur qoldiqni oʻqishi kerak boʻlsa,
  `available`/`requested` kabi kengaytma qoʻshiladi. Bu darsda faqat odam oʻqiydigan matn.
- **Shakl oʻzgarishi buzuvchi.** `{"error":{"code":…}}` oʻqiydigan mijozlar yangi shaklda `code`ni topmaydi; ishda bunday
  oʻzgarish oldindan eʼlon qilinadi yoki versiya bilan chiqariladi.
- **Nomaʼlum kalit.** GB-4.06'da `json: unknown field "rang"` matni mijozga kalit nomini aytardi. Endi u umumiy `400`: bu xato
  turi yoʻq (`*errors.errorString`), kalitni matndan qirqib olish esa moʻrt. Narxi: mijoz qaysi kalit ortiqcha ekanini bilmaydi.

**Statuslar (`problemFor`):** `store.ErrNotFound` → 404 · `*store.StockError` → 409 (`/problems/insufficient-stock`, detail:
nom, qoldiq, soʻralgan) · `errBadID` → 400 · `*json.UnmarshalTypeError` (`Field` boʻsh emas) → 400, `"price: qiymat turi
notoʻgʻri"` (kalitni mijoz oʻzi yuborgan; Go turi aytilmaydi) · `errBadJSON` → 400 umumiy matn · **vazifa:** `fieldErrors` →
422 + `errors`, `ErrBadProduct`/`ErrBadQuantity` → 422 + `productRules`, `*http.MaxBytesError` → 413 · qolgan hammasi →
**500**, `detail`siz. Middleware'lar: panic → 500, `Content-Type` JSON emas → 415, mux'ning oʻzi → 404/405 (`Allow` saqlanadi).

**Log.** `writeError(w, err)` har xato uchun bitta qator yozadi: `<ID> xato <status>: "<err.Error()>"` (`%q`: matndagi yangi
qator `\n` boʻlib qoladi, GB-6.01), undan keyin GB-6.01'dagi `logRequests` qatori (metod, yoʻl, status, vaqt). Bu darsda 4xx ham
toʻliq yoziladi; koʻp trafikda 4xx'ni pastroq darajada yozish — `slog` darajalari, GB-8.02.

**`fieldErrors` haqida.** U map turi, `Error()` metodi bilan — demak `error`. `validate` uni aniq tur sifatida qaytaradi va handler
`len(f) > 0` bilan tekshiradi. Uni `error` oʻzgaruvchisiga solib `!= nil` bilan tekshirmang: boʻsh map ham nil emas, shart doim rost.

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, curl 8.22.0, `CGO_ENABLED=0`)

Hammasi `_video/record.sh` bilan shu kompyuterda olingan (toʻliq: `_video/chiqishlar.txt`). Server: birinchi terminalda
`go run ./cmd/api` (`bozor API: http://localhost:8765`), soʻrovlar ikkinchisida. `curl -i`: status va sarlavhalar ham.

```
$ curl -s localhost:8765/products -H "Content-Type: application/json" -d '{"name":"choy","price":"25000"}'
{"type":"about:blank","title":"Bad Request","status":400,"detail":"price: qiymat turi notoʻgʻri","request_id":"VQQZMDGNWMAB2MT32INXIC6VJB"}
```

Birinchi terminalda (server log'i) — xuddi shu ID bilan:

```
2026/11/27 21:54:51 VQQZMDGNWMAB2MT32INXIC6VJB xato 400: "body JSON sifatida oʻqilmadi: json: cannot unmarshal string into Go struct field productInput.price of type int64"
2026/11/27 21:54:51 VQQZMDGNWMAB2MT32INXIC6VJB POST "/products" 400 144.523µs
```

```
$ curl -i localhost:8765/orders -H "Content-Type: application/json" -d '{"product_id":4,"quantity":5}'
HTTP/1.1 409 Conflict
Content-Type: application/problem+json
X-Request-Id: 2VWWA73LU2WUTIAPGEXDJPSHBC
…
{"type":"/problems/insufficient-stock","title":"Omborda yetarli emas","status":409,"detail":"oʻrik: omborda 2 dona, soʻralgan 5","request_id":"2VWWA73LU2WUTIAPGEXDJPSHBC"}

$ curl -s localhost:8765/narsalar
{"type":"about:blank","title":"Not Found","status":404,"detail":"bunday yoʻl yoʻq","request_id":"3VMJYYNPACIVQ5MQE7EPDTIIWI"}

$ curl -i -X PUT localhost:8765/products/1
HTTP/1.1 405 Method Not Allowed
Allow: DELETE, GET, HEAD
Content-Type: application/problem+json
…
{"type":"about:blank","title":"Method Not Allowed","status":405,"detail":"bu yoʻl uchun bunday metod yoʻq","request_id":"D3VYDXSSFKLTPHEYQSYFAZ3U6E"}
```

Hali tanilmagan xato (vazifaning 1-qadami) — xavfsiz `500`, sabab log'da:

```
$ curl -s localhost:8765/products -H "Content-Type: application/json" -d '{"name":"","price":0}'
{"type":"about:blank","title":"Internal Server Error","status":500,"request_id":"6X4XV6QP6MBW5QLR5PKPBLQJ4M"}

2026/11/27 21:54:51 6X4XV6QP6MBW5QLR5PKPBLQJ4M xato 500: "validatsiya: 2 ta maydon notoʻgʻri"
2026/11/27 21:54:51 6X4XV6QP6MBW5QLR5PKPBLQJ4M POST "/products" 500 96.227µs
```

**Tez-tez uchraydigan xato:** `writeError`dan keyin `return` unutilsa (`_video/unutgan/handlers.go`, `getProduct`; server shu
nusxadan xuddi shunday `go run ./cmd/api` bilan ishga tushirilgan), handler
davom etadi: mijoz **ikki** JSON oladi, net/http log'da ogohlantiradi, `logRequests` esa notoʻgʻri status yozadi
(`statusRecorder` oxirgi `WriteHeader`ni eslab qoladi, mijozga esa birinchisi — 404 — ketgan):

```
$ curl -s localhost:8765/products/99
{"type":"about:blank","title":"Not Found","status":404,"detail":"bunday ID'li mahsulot yoʻq","request_id":"DAIXJUQK5OKJ3CWDLXS3L4OSKR"}
{"id":0,"name":"","price":0,"stock":0}

2026/11/27 21:54:51 DAIXJUQK5OKJ3CWDLXS3L4OSKR xato 404: "mahsulot 99: mahsulot topilmadi"
2026/11/27 21:54:51 http: superfluous response.WriteHeader call from main.(*statusRecorder).WriteHeader (middleware.go:48)
2026/11/27 21:54:51 DAIXJUQK5OKJ3CWDLXS3L4OSKR GET "/products/99" 200 170.472µs
```

Tuzatish: `writeError(w, err)`dan keyin darhol `return`. (GB-6.01'dagi `statusRecorder` har `WriteHeader`da `status`ni qayta
yozadi; faqat birinchisini eslab qolsa, log bunday holatda ham toʻgʻri statusni koʻrsatardi — bu GB-6.01 kodining kelgusi tuzatishi.)

## Vazifa

[TASK.md](TASK.md): `problemFor`ga uchta case — `fieldErrors` → 422 + `errors`, `ErrBadProduct`/`ErrBadQuantity` → 422 +
`productRules`, `*http.MaxBytesError` → 413 (`errBadJSON` case'idan **oldin**). Tekshirish: `go test ./...`.

```
$ go test ./...          # starter
--- FAIL: TestVazifaFieldErrors (0.00s)
--- FAIL: TestVazifaStoreRules (0.00s)
--- FAIL: TestVazifaTooLarge (0.00s)
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api	0.004s
ok  	github.com/abdurahmon27/bozor/cmd/bozor	0.027s
ok  	github.com/abdurahmon27/bozor/internal/store	0.210s
FAIL

$ cd solution && go test ./...
ok  	github.com/abdurahmon27/bozor/cmd/api	0.007s
ok  	github.com/abdurahmon27/bozor/cmd/bozor	0.033s
ok  	github.com/abdurahmon27/bozor/internal/store	0.214s
```

## Manbalar

- RFC 9457, Problem Details for HTTP APIs — https://www.rfc-editor.org/rfc/rfc9457 (§3.1 aʼzolar, §3.2 kengaytmalar,
  §4.2.1 `about:blank`, §5 xavfsizlik: "avoid making implementation details such as a stack dump available")
- OWASP Error Handling Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Error_Handling_Cheat_Sheet.html
  ("a generic response is returned … the error details are logged server side")
- OWASP Logging Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html
- CWE-209: Generation of Error Message Containing Sensitive Information — https://cwe.mitre.org/data/definitions/209.html
- RFC 9110 (status kodlari 400, 404, 405 + Allow, 409, 413, 415, 422, 500) — https://www.rfc-editor.org/rfc/rfc9110
- Go: `errors.Is`/`errors.As` — https://pkg.go.dev/errors · `fmt.Errorf` (bir nechta `%w`, Go 1.20) — https://pkg.go.dev/fmt#Errorf ·
  `json.UnmarshalTypeError` (`Field` — "the full path from root node to the value") — https://pkg.go.dev/encoding/json#UnmarshalTypeError ·
  `http.StatusText`, `http.MaxBytesError`, `ResponseWriter.WriteHeader` — https://pkg.go.dev/net/http
