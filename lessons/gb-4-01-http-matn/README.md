# GB-4.01 · HTTP soʻrov — bu matn: `curl -v` bilan bitta almashinuv

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #25"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** 4-modul darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-3-06` (bozor CLI v0: add, list, order, JSON ombor) · keyingi dars: `gb-4-02` (Go'da HTTP client:
timeout, status va body'ni yopish).

**Bu darsdan keyin:** `curl -v` chiqishida soʻrov va javobni ajratasiz: method, yoʻl, versiya, header'lar, boʻsh
qator va body; status kodini sinfi bilan oʻqiysiz (`2xx` muvaffaqiyat, `4xx` xato soʻrovda, `5xx` xato serverda);
`Content-Type` va `Content-Length` header'larining nima uchun borligini bilasiz (`Content-Length` baytlarni sanaydi);
soʻrovni `printf` + `nc` bilan qoʻlda yozib koʻrasiz; Go'da soʻrov qatorini ajratadigan `parseRequestLine` va
`statusClass` funksiyalarini yozasiz.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | boshlangʻich kod: `Request` turi va ikkita boʻsh funksiya (`parseRequestLine`, `statusClass`) — vazifa shu yerda |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/server/main.go` | videodagi server: bozor katalogi HTTP orqali (ichki kodini darsda koʻrmaymiz). Uni GB-4.04'da qatorma-qator yozasiz; bu darsda faqat `curl` bilan soʻraymiz. `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/curl-chiqishlari.txt` | videodagi hamma buyruqning toʻliq (`*` qatorlarsiz) chiqishi |

Soʻrov va javob matnini oʻzingiz koʻrish uchun serverni bir terminalda ishga tushiring (port band boʻlsa, manzilni
argument qilib bering: `go run ./_video/server localhost:9000`, quyidagi buyruqlarda ham portni almashtiring):

```bash
go run ./_video/server
```

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`, curl 8.22.0)

Quyida har bir `curl -v` chiqishidan `*` bilan boshlangan qatorlar (ulanish haqida) va `[N bytes data]` qatorlari
olib tashlangan; qolgani oʻzgartirilmagan. `Date` qiymati har safar boshqacha boʻladi (hamma qadam bir soniyada ishlagani uchun bu yerda bir xil).

**1. GET** — `curl -v http://localhost:8765/products/1`: `>` — curl yuborgan soʻrov, `<` — serverning javobi.

```
> GET /products/1 HTTP/1.1
> Host: localhost:8765
> User-Agent: curl/8.22.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 27 Nov 2026 06:37:46 GMT
< Content-Length: 49
< 
{"id":1,"name":"olma","price":1200000,"stock":40}
```

**2. Qoʻlda yozilgan soʻrov** — `printf 'GET /products/1 HTTP/1.1\r\nHost: localhost:8765\r\nConnection: close\r\n\r\n' | nc localhost 8765`
(`\r\n` — HTTP'da qator oxiri, ikki bayt; `nc` matnni TCP ulanishiga yuboradi va javobni chiqaradi). `Connection: close` bo'lmasa server ulanishni ochiq qoldiradi va `nc` javobdan keyin ham kutib turadi (bu mashinada 3 soniyalik `timeout` bilan sinab koʻrildi: chiqish kodi 124); `nc` dasturlari har xil (bu yerda OpenBSD netcat; GNU variantida bayroqlar boshqacha boʻlishi mumkin):

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 27 Nov 2026 06:37:46 GMT
Content-Length: 49
Connection: close

{"id":1,"name":"olma","price":1200000,"stock":40}
```

**3. Yoʻq mahsulot** — `curl -v http://localhost:8765/products/99` (status `404`, `4xx` — xato soʻrovda):

```
> GET /products/99 HTTP/1.1
> Host: localhost:8765
> User-Agent: curl/8.22.0
> Accept: */*
> 
< HTTP/1.1 404 Not Found
< Content-Type: application/json
< Date: Fri, 27 Nov 2026 06:37:46 GMT
< Content-Length: 30
< 
{"error":"mahsulot topilmadi"}
```

**4. POST** — `curl -v -d '{"name":"oʻrik","price":3500000,"stock":5}' -H "Content-Type: application/json" http://localhost:8765/products`
(body 42 belgi, lekin `Content-Length: 43`: `ʻ` UTF-8'da ikki bayt, GB-2.10):

```
> POST /products HTTP/1.1
> Host: localhost:8765
> User-Agent: curl/8.22.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 43
> 
< HTTP/1.1 201 Created
< Content-Type: application/json
< Location: /products/3
< Date: Fri, 27 Nov 2026 06:37:46 GMT
< Content-Length: 50
< 
{"id":3,"name":"oʻrik","price":3500000,"stock":5}
```

**5. Klassik xato** — `-H` unutildi: curl `Content-Type: application/x-www-form-urlencoded` qoʻyadi, server `415` qaytaradi:

```
> POST /products HTTP/1.1
> Host: localhost:8765
> User-Agent: curl/8.22.0
> Accept: */*
> Content-Length: 43
> Content-Type: application/x-www-form-urlencoded
> 
< HTTP/1.1 415 Unsupported Media Type
< Content-Type: application/json
< Date: Fri, 27 Nov 2026 06:37:46 GMT
< Content-Length: 57
< 
{"error":"Content-Type application/json boʻlishi kerak"}
```

**6. Tuzatilgan** — `-H "Content-Type: application/json"` bilan qayta (server oldingi muvaffaqiyatli POST'ni eslab qolgan,
shuning uchun `id` 4):

```
> POST /products HTTP/1.1
> Host: localhost:8765
> User-Agent: curl/8.22.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 43
> 
< HTTP/1.1 201 Created
< Content-Type: application/json
< Location: /products/4
< Date: Fri, 27 Nov 2026 06:37:46 GMT
< Content-Length: 50
< 
{"id":4,"name":"oʻrik","price":3500000,"stock":5}
```

**7. Qoʻshimcha** (videoda yoʻq): server `GET /products/{id}` va `POST /products`ni biladi; `DELETE` esa `405`:

```
> DELETE /products/1 HTTP/1.1
> Host: localhost:8765
> User-Agent: curl/8.22.0
> Accept: */*
> 
< HTTP/1.1 405 Method Not Allowed
< Allow: GET, HEAD
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Fri, 27 Nov 2026 06:37:46 GMT
< Content-Length: 19
< 
Method Not Allowed
```

## Vazifa

[`TASK.md`](TASK.md): `parseRequestLine` va `statusClass`; tekshirish: `go test` (oʻtguncha). Javob: [`solution/`](solution).

```bash
gofmt -l . && go vet ./...
go test .             # boshlangʻichda yiqiladi (TASK.md xabari bilan)
go test ./solution/   # javob: ok
```

## Manbalar

- [RFC 9112 §2–§6: HTTP/1.1 xabar formati — soʻrov qatori (§3), status qatori (§4), header'lar (§5), boʻsh qator, body (§2, §6)](https://www.rfc-editor.org/rfc/rfc9112)
- [RFC 9110 §15: status kodlari sinflari va `200`, `201`, `400`, `404`, `415`, `500`](https://www.rfc-editor.org/rfc/rfc9110#name-status-codes)
- [RFC 9110 §8.3 Content-Type, §8.6 Content-Length (oktetlar = baytlar)](https://www.rfc-editor.org/rfc/rfc9110#name-content-type)
- [curl `-v/--verbose`, `-d/--data`, `-H/--header`](https://curl.se/docs/manpage.html)
- [Go `net/http` paketi](https://pkg.go.dev/net/http) (server GB-4.04'da)
