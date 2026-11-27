# GB-4.05 · Route'lar va `{id}`: `GET /products/{id}`, 404 va 405

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #29"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. Pattern'dagi metod va `{id}` Go 1.22'dan beri bor.

Oldingi dars: [`gb-4-04-birinchi-server`](../gb-4-04-birinchi-server) (`newServer`: `http.NewServeMux`, `mux.HandleFunc`,
`http.Server`, `:8765`) · keyingi dars: [`gb-4-06-json-body`](../gb-4-06-json-body) (`POST /products`: JSON kirish, limit,
validatsiya, bir xil xato javobi). Bu darsga kerak boʻlganlar: [`gb-4-03-handler-httptest`](../gb-4-03-handler-httptest)
(handler imzosi, `httptest`, javob tartibi, `http.Handler`/`ServeHTTP`), [`gb-4-01-http-matn`](../gb-4-01-http-matn)
(soʻrov qatori, status kodlari, `curl`), [`gb-3-06-bozor-cli`](../gb-3-06-bozor-cli) (`strconv.ParseInt(s, 10, 64)`),
[`gb-2-04-pointer`](../gb-2-04-pointer) (`findProduct` → `*Product` yoki `nil`).

**Bu darsdan keyin:** pattern'ga metod yozasiz (`"GET /products/{id}"`), yoʻldagi oʻzgaruvchi qismni `r.PathValue("id")`
bilan oʻqib, son boʻlmasa `400`, topilmasa `404` qaytarasiz; mux oʻzi yozadigan `404 page not found` va
`405 Method Not Allowed` + `Allow: GET, HEAD`ni oʻzingiznikidan ajratasiz; route'larni `newMux()`ga yigʻib, serversiz
`newMux().ServeHTTP(rec, req)` bilan test qilasiz; `"/"` pattern'i nega hamma yoʻlni yutishini va `"GET /{$}"` bilan
tuzatishni bilasiz.

Bu papka mustaqil: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i, tashqi modul yoʻq.
`Product`, `olma`, `catalog`, `healthHandler`, `catalogHandler`, `newServer` GB-4.03/4.04'dan, `findProduct` GB-2.04'dan;
yangi: `productHandler` (endi `{id}` bilan, GB-4.04'dagi `/products/1` oʻrniga), `homeHandler`, `newMux` (route'lar
`newServer`dan ajratildi: testlar ham shuni ishlatadi), `stockHandler` (vazifa). `main`dagi `srv == nil` tekshiruvi
(GB-4.04 vazifasi uchun edi) olib tashlandi.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `GET /{$}`, `GET /health`, `GET /products`, `GET /products/{id}`; `stockHandler` vazifa uchun boʻsh |
| `handler_test.go` | videoda yozilgan test (`TestProduct`) va yana uchtasi (400/404, 405 + `Allow`, `"/"` yutmasligi): starterda ham oʻtadi |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayllari) |
| `_video/NN-…/` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 01-route go run .` |
| `_video/chiqishlar.txt` | videodagi har bir buyruqning toʻliq (qisqartirilmagan) chiqishi |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`, curl 8.22.0)

Server birinchi terminalda (`./_video/run.sh 01-route go run .`, toʻxtatish — Ctrl+C), `curl` ikkinchisida. Videoda
`Date`, `Content-Length`, `X-Content-Type-Options` va boʻsh qator `…` bilan yashirilgan; toʻliq matn `_video/chiqishlar.txt`da.

Pattern'da metod (`"GET /products"`): boshqa metodga javobni **mux** yozadi, handler chaqirilmaydi. `GET` pattern'i
`HEAD`ni ham oladi:

```
$ curl -i -X POST http://localhost:8765/products
HTTP/1.1 405 Method Not Allowed
Allow: GET, HEAD
…
Method Not Allowed
```

Hech bir pattern'ga mos kelmaydigan yoʻl — mux'ning `404`i:

```
$ curl -i http://localhost:8765/narsalar
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
…
404 page not found
```

`GET /products/{id}`: `{id}` bitta segment; `r.PathValue("id")` matn qaytaradi, `strconv.ParseInt` uni songa aylantiradi:

```
$ curl -i http://localhost:8765/products/3
HTTP/1.1 200 OK
Content-Type: application/json
…
{"id":3,"name":"oʻrik","price":3500000,"stock":5}
$ curl -i http://localhost:8765/products/99
HTTP/1.1 404 Not Found
…
mahsulot topilmadi
$ curl -i http://localhost:8765/products/abc
HTTP/1.1 400 Bad Request
…
id butun son boʻlishi kerak
```

Ikki xil `404`: `404 page not found` — mux'niki (yoʻl yoʻq), `mahsulot topilmadi` — bizniki (yoʻl bor, mahsulot yoʻq).
`/products/3/` (oxirida slesh) ham mux'ning `404`i: `{id}` faqat bitta segment.

Tuzoq: handler'ni testda toʻgʻridan-toʻgʻri chaqirsangiz, `PathValue`ni toʻldiradigan mux yoʻq — `""` qaytadi
(`02-test-direct`):

```
$ go test
--- FAIL: TestProduct (0.00s)
    handler_test.go:14: status 400, kutilgan 200: id butun son boʻlishi kerak
FAIL
exit status 1
FAIL	github.com/abdurahmon27/bozor	0.002s
```

Tuzatish: route'lar `newMux()` funksiyasida (`newServer` uni `Handler`ga beradi), test `newMux().ServeHTTP(rec, req)` chaqiradi (`03-test-mux`):

```
$ go test -v
=== RUN   TestProduct
--- PASS: TestProduct (0.00s)
PASS
ok  	github.com/abdurahmon27/bozor	0.002s
```

Klassik xato — `mux.HandleFunc("/", homeHandler)` (`04-catchall`): slesh bilan tugaydigan pattern shu bilan
boshlanadigan hamma yoʻlni oladi, metodsiz pattern esa hamma metodni. Natijada xato yoʻl ham, notoʻgʻri metod ham `200`:

```
$ curl -i http://localhost:8765/produkts/3
HTTP/1.1 200 OK
…
bozor API
$ curl -i -X POST http://localhost:8765/products/3
HTTP/1.1 200 OK
…
bozor API
```

Tuzatish — `mux.HandleFunc("GET /{$}", homeHandler)` (`05-exact`): `{$}` — yoʻl shu yerda tugaydi, faqat aynan `/`:

```
$ curl -i http://localhost:8765/produkts/3
HTTP/1.1 404 Not Found
…
$ curl -i -X POST http://localhost:8765/products/3
HTTP/1.1 405 Method Not Allowed
…
```

Chuqurroq — bir xil soʻrovlarni oladigan ikki pattern (`06-conflict`, `"GET /products/{name}"` qoʻshilgan) dastur
ishga tushishidayoq panic beradi:

```
$ go run .
panic: pattern "GET /products/{name}" (registered at …) conflicts with pattern "GET /products/{id}" (…):
	GET /products/{name} matches the same requests as GET /products/{id}
…
exit status 2
```

Oldingi darsdan nima oʻzgargani: `diff -ru ../gb-4-04-birinchi-server/solution .` (oldingi darsdan nima oʻzgargani).

## Nima uchun shunday

- **Ustunlik:** ikki pattern bitta soʻrovga mos kelsa, *aniqrogʻi* (uning har soʻrovini ikkinchisi ham oladi, aksi emas) yutadi; shuning uchun `"/"`
  turganda ham `GET /products/3` `productHandler`ga boradi. Roʻyxatdan oʻtkazish tartibi ahamiyatsiz.
- **Metodsiz pattern** (`"/products"`) hamma metodni oladi: unda `405`ni handler oʻzi yozishi kerak (GB-4.03'dagidek). Metodli
  pattern'da buni mux qiladi va `Allow`ni ham oʻzi qoʻyadi (`Allow` `405`da majburiy: RFC 9110 §15.5.6).
- **`http.Error`** `Content-Type: text/plain; charset=utf-8` va `X-Content-Type-Options: nosniff` qoʻyadi, `Content-Length`ni
  oʻchiradi, statusni yozadi, keyin xabar + yangi qator — GB-4.03'dagi toʻgʻri tartib. Hamma endpoint'da bir xil JSON xato
  shakli — GB-4.06.
- **`r.PathValue`** faqat soʻrov mux orqali (`ServeHTTP`) oʻtganda toʻladi; mux'siz testda `req.SetPathValue("id", "3")` bilan
  qoʻlda qoʻyish ham mumkin, lekin `newMux().ServeHTTP` pattern'ning oʻzini ham tekshiradi.
- `{id}` son ekanini tekshirmaydi: `/products/abc` mos keladi, tekshiruv handler'da (`ParseInt` → `400`). `ParseInt`
  `+3` va `03`ni ham `3` deb oʻqiydi (bitta mahsulotga bir nechta URL), juda katta son esa `400` oladi (son, lekin
  `int64`ga sigʻmaydi). Ishlab turgan API'da `id < 1`ni ham rad etish va kanonik shaklni talab qilish mumkin; bu darsda
  soddalik uchun yoʻq.
- **`http.Error` handler'ni toʻxtatmaydi**: undan keyin `return` boʻlmasa, keyingi `Write` xabar ortidan body'ga qoʻshiladi
  (`tekshiruv_test.go` buni ushlaydi).
- **`go.mod`dagi `go` qatori muhim:** metodli pattern va `{id}` `go 1.22` yoki undan yangi modulda ishlaydi; eski `go`
  qatorida mux eski qoidalar bilan ishlaydi (GODEBUG `httpmuxgo121`, ServeMux hujjatining Compatibility boʻlimi).

## Manbalar

- [`net/http`: ServeMux](https://pkg.go.dev/net/http#ServeMux) (Patterns, wildcard'lar, `{$}`, Precedence, conflict → panic,
  GET HEAD'ni ham oladi), [Request.PathValue](https://pkg.go.dev/net/http#Request.PathValue),
  [Error](https://pkg.go.dev/net/http#Error), [NewServeMux](https://pkg.go.dev/net/http#NewServeMux),
  [ServeMux.ServeHTTP](https://pkg.go.dev/net/http#ServeMux.ServeHTTP).
- [Go 1.22 release notes: enhanced routing patterns](https://go.dev/doc/go1.22#enhanced_routing_patterns),
  [Go blog: Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements) (405 + Allow).
- [RFC 9110: 405 Method Not Allowed](https://www.rfc-editor.org/rfc/rfc9110#name-405-method-not-allowed),
  [HEAD](https://www.rfc-editor.org/rfc/rfc9110#name-head).
- [curl man page](https://curl.se/docs/manpage.html) (`-i`, `-X`).
