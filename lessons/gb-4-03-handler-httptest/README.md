# GB-4.03 · Handler — oddiy funksiya: `httptest` bilan serversiz tekshirish

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #27"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-4-02-http-client`](../gb-4-02-http-client) (client tomoni: `err` va status, `resp.Body`) · keyingi dars:
`gb-4-04` (birinchi server: `ServeMux`, `ListenAndServe`). Bu darsga kerak boʻlganlar:
[`gb-4-01-http-matn`](../gb-4-01-http-matn) (soʻrov qatori, `Content-Type`, `405`),
[`gb-3-02-json-saqlash`](../gb-3-02-json-saqlash) (`Product`, tag'lar), [`gb-2-07-interface-writer`](../gb-2-07-interface-writer)
(`io.Writer`, `fmt.Fprintln`), [`gb-1-08-birinchi-test`](../gb-1-08-birinchi-test) (`TestXxx`, `t.Errorf`, `-v`, `-run`).

**Bu darsdan keyin:** handler imzosi `func(w http.ResponseWriter, r *http.Request)`ni oʻqiysiz, handler'ni server
ishga tushirmasdan `httptest.NewRequest` va `httptest.NewRecorder` bilan chaqirib, statusni, sarlavhani va JSON body'ni
test qilasiz; javob tartibini (sarlavhalar → status → body) va `WriteHeader` chaqirilmasa status nega `200` ekanini
bilasiz; GB-4.02'dagi test fayli (`http.HandlerFunc`, `httptest.NewServer`) endi sizga tanish.

Bu papka mustaqil: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i, tashqi modul yoʻq.
`Product` (`ID`, `Name`, `Price` tiyinda, `Stock`, json tag'lari) GB-3.02'dan; yangi: `healthHandler`, `productHandler`,
`catalog`, `catalogHandler` (vazifa).

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `healthHandler`, `productHandler`, `main` (handler'ni serversiz chaqiradi); `catalogHandler` vazifa uchun boʻsh |
| `handler_test.go` | videoda yozilgan testlar (`TestHealth`, `TestProduct`, `TestProductMetod`): starterda ham oʻtadi |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayllari) |
| `_video/NN-…/` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 02-recorder go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

`*http.Request` — soʻrov matnining struct koʻrinishi (`01-request`: `httptest.NewRequest(http.MethodGet, "/products/1", nil)`):

```
$ go run .
GET /products/1 HTTP/1.1
```

Handler'ni oddiy funksiya kabi chaqirish: `NewRecorder` — hech qayerga yubormaydigan `ResponseWriter`; `WriteHeader`
chaqirilmasa, birinchi `Write` statusni `200` qiladi, `Fprintln` esa oxiriga yangi qator qoʻshadi (`02-recorder`,
`%q` yashirin belgilarni koʻrsatadi):

```
$ go run .
200 "ok\n"
```

Test (`03-test`, `go test -v`):

```
$ go test -v
=== RUN   TestHealth
--- PASS: TestHealth (0.00s)
PASS
ok  	github.com/abdurahmon27/bozor	0.002s
```

`productHandler` (`04-product`): `main` shu handler'ni chaqiradi; narx tiyinda, kalitlar tag'lardan, `Encode` oxiriga
yangi qator qoʻshadi:

```
$ go run .
200 {"id":1,"name":"olma","price":1200000,"stock":40}
$ go test -v -run Product
=== RUN   TestProduct
--- PASS: TestProduct (0.00s)
=== RUN   TestProductMetod
--- PASS: TestProductMetod (0.00s)
PASS
ok  	github.com/abdurahmon27/bozor	0.003s
```

Eng koʻp uchraydigan xato — avval body, keyin status (`05-tartib`, `06-tartib-fix`). `Fprintln` birinchi `Write`'da `200`ni
joʻnatib boʻldi, keyingi `WriteHeader(404)` ta'sir qilmaydi:

```
$ go test
--- FAIL: TestNotFound (0.00s)
    handler_test.go:13: status 200, kutilgan 404
FAIL
exit status 1
FAIL	github.com/abdurahmon27/bozor	0.002s
```

Tuzatilgani (avval `WriteHeader`, keyin `Fprintln`) — `PASS`. Haqiqiy serverda esa bu xato testsiz jimgina mijozga
ketadi, faqat server log'ida shunday satr chiqadi (biz buni `httptest.NewServer` bilan koʻrdik; sana va `main.go:12`
sizda boshqacha boʻladi):

```
http: superfluous response.WriteHeader call from main.main.func1 (main.go:12)
```

`http.HandlerFunc` va `httptest.NewServer` (`07-server`) — GB-4.02'dagi testlar shundan foydalangan:

```
$ go run .
200 OK
```

## Nima uchun shunday

- Handler — `http.Handler` interface'ining (`ServeHTTP(w, r)`) oddiy funksiya shakli. `http.HandlerFunc` — funksiya turi:
  uning `ServeHTTP` metodi ichidagi funksiyani chaqiradi, shuning uchun oddiy funksiyani `Handler` sifatida berish mumkin.
- Handler'ning oʻzini tekshirish uchun `httptest.NewRecorder` yetarli (tez, port kerak emas). `httptest.NewServer` haqiqiy
  server ochadi (boʻsh portda): client kodini tekshirishda kerak (GB-4.02).
- `405` javobida `Allow` sarlavhasi majburiy ([RFC 9110 §15.5.6](https://www.rfc-editor.org/rfc/rfc9110#name-405-method-not-allowed)).
  GB-4.05'da `ServeMux` `GET /products` kabi pattern bilan metodni oʻzi tekshiradi; bu darsda `r.Method`ni qoʻlda oʻqidik.
- `Allow: GET` HEAD'ni ham rad etadi; ishlab turgan serverda `HEAD` ham `GET` kabi javob berishi kerak (RFC 9110 §9.3.2). GB-4.05'da `GET /products` patterni HEAD'ni ham qabul qiladi; bu darsda soddalik uchun GET'dan boshqa hamma metod `405`.
- `Encode(...)` xatosi videoda tekshirilmaydi: sarlavhalar ketib boʻlgan, statusni oʻzgartirib boʻlmaydi, xato faqat uzilgan
  mijozni bildiradi. Uni log'ga yozish — GB-8.02.

## Manbalar

- [`net/http`: ResponseWriter](https://pkg.go.dev/net/http#ResponseWriter) (Header, Write, WriteHeader; `WriteHeader` chaqirilmasa
  birinchi `Write` `200` yuboradi; `Write`/`WriteHeader`dan keyin sarlavhani oʻzgartirish ta'sir qilmaydi),
  [Handler](https://pkg.go.dev/net/http#Handler), [HandlerFunc](https://pkg.go.dev/net/http#HandlerFunc),
  [Request](https://pkg.go.dev/net/http#Request).
- [`net/http/httptest`](https://pkg.go.dev/net/http/httptest): `NewRequest`, `NewRecorder` (`Code`, `Body`, `Header`), `NewServer`.
- [`encoding/json`: NewEncoder](https://pkg.go.dev/encoding/json#NewEncoder).
- [RFC 9110, 405 Method Not Allowed](https://www.rfc-editor.org/rfc/rfc9110#name-405-method-not-allowed).
