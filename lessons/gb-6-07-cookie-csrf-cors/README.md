# GB-6.07 · Cookie, CSRF, CORS va rate limit

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #49"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. `http.CrossOriginProtection` Go 1.25'dan beri bor.

Oldingi darslar: [`gb-6-06-buyurtma-kimniki`](../gb-6-06-buyurtma-kimniki) (egalik tekshiruvi, buyurtmalar 41/42) ·
[`gb-6-04-login-session`](../gb-6-04-login-session) (session cookie: `HttpOnly`, `Secure`, `SameSite` faqat nomlangan edi)
· keyingi dars: `gb-6-08` (nazorat nuqtasi: himoyalangan bozor).

**Bu darsdan keyin:** bozor brauzer tomondan ham himoyalangan — session cookie'ning har atributi nima qilishini bilasiz,
boshqa saytdan kelgan soxta POST'ni (CSRF) `http.CrossOriginProtection` bilan yopasiz, CORS'ni allow-list bilan
sozlaysiz va `POST /login`ga 429 + `Retry-After` qaytaradigan rate limit qoʻyasiz.

## Nega store xotirada

Bu laboratoriya — oʻzi mustaqil papka. Mavzu brauzer va server orasidagi qoidalar, Postgres emas, shuning uchun store
xotirada (`internal/store`: map + mutex), metodlari GB-6.04/6.06'dagi Postgres Store bilan bir xil shaklda
(`UserBySession`, `Order(ctx, id, userID)`, …). Docker kerak emas. Server har ishga tushganda xotira `Seed`dan
boshlanadi: Dilnoza (`dilnoza@bozor.test` / `tarvuz-osmon-qalam-ariq-27`, buyurtma 41) va Jasur
(`jasur@bozor.test` / `choy-non-qand-27`, buyurtma 42) — parollar faqat shu laboratoriya uchun.

Yangi route: `POST /orders/{id}/cancel` — buyurtmani bekor qiladi, faqat egasi (GB-6.06 sharti). **Body yoʻq**, shuning
uchun `requireJSON` yoʻq va oddiy HTML forma ham unga yetib keladi: CSRF misoli shu route'da.

```
cmd/api/        bozor API (:8765)
  auth.go         login, session cookie (__Host-session, Max-Age), requireUser
  browser.go      protectCSRF (http.CrossOriginProtection) va cors (allow-list)  ← yangi
  ratelimit.go    limiters (golang.org/x/time/rate), limitByIP, tooManyRequests   ← yangi
  middleware.go   zanjir: requestID → logRequests → recoverPanic → cors → protectCSRF → muxErrors → mux
cmd/begona/     «begona sayt» (:9090): CSRF formasi (/) va CORS sahifasi (/oqish)
_video/         videodagi chiqishlar: chiqishlar.sh → chiqishlar.txt, bosqichlar (stages/), brauzer/
```

## Ishga tushirish

```bash
go run ./cmd/api        # 1-terminal: bozor API: http://localhost:8765
go run ./cmd/begona     # 2-terminal: begona sayt: http://localhost:9090 va http://127.0.0.1:9090
```

Brauzer tajribalari videoda **haqiqiy Chrome** (headless) bilan qilingan. Uni `_video/brauzer` (alohida modul,
[chromedp](https://github.com/chromedp/chromedp) kutubxonasi) boshqaradi; Google Chrome oʻrnatilgan boʻlishi kerak:

```bash
cd _video/brauzer
go run . httponly                     # login → bozor sahifasida document.cookie va fetch("/me")
go run . csrf http://127.0.0.1:9090   # login → begona sahifa → «Sovgʻani olish» → buyurtma 41
go run . cors http://localhost:9090   # login → begona sahifadagi skript /me'ni oʻqishga urinadi
```

Hammasini qayta yozish (har bosqich alohida quriladi, oxirida jarayonlar toʻxtatiladi): `bash _video/chiqishlar.sh`.
Quyidagi chiqishlar — `_video/chiqishlar.txt`dan (go1.27.4 linux/amd64, curl 8.22.0, Google Chrome 155, 2026-11-28).

## 1. Cookie atributlari

```
$ curl -i -c jar.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'
HTTP/1.1 200 OK
Cache-Control: no-store
Content-Type: application/json
Set-Cookie: __Host-session=3Z8H3_sk7LTAV-SjJe_sZ-zxDUVndwMK7eq22LJxCF4; Path=/; Max-Age=86400; HttpOnly; Secure; SameSite=Lax
Vary: Origin
…
```

| atribut | brauzer nima qiladi | manba |
|---|---|---|
| `__Host-` prefiksi | cookie'ni faqat `Secure`, `Path=/` va `Domain`siz boʻlsa qabul qiladi: uni faqat shu host oʻrnatadi, qoʻshni subdomen ustidan yoza olmaydi | RFC 6265bis-22 §4.1.3.2 |
| `Max-Age=86400` | 86400 soniya (24 soat) saqlaydi; nisbiy, brauzer soatiga bogʻliq emas; `Expires`dan ustun | RFC 6265bis-22 §4.1.2.2 |
| `HttpOnly` | sahifa skripti (`document.cookie`) cookie'ni koʻrmaydi; soʻrovlarga baribir qoʻshadi | MDN Document.cookie |
| `Secure` | faqat HTTPS orqali yuboradi; `http://localhost` — istisno | MDN Set-Cookie |
| `SameSite=Lax` | boshqa **site**'dan boshlangan soʻrovga faqat yuqori darajadagi xavfsiz navigatsiyada (havola bosilganda GET) qoʻshadi; POST, iframe, fetch — yoʻq | RFC 6265bis-22 §5.6.7, OWASP CSRF |

Haqiqiy Chrome'da `HttpOnly`:

```
$ go run . httponly
login: 200 {"id":1,"email":"dilnoza@bozor.test"}
document.cookie: ""
fetch /me: 200 {"id":1,"email":"dilnoza@bozor.test"}
```

`curl` ham `http://localhost` uchun `Secure` va `__Host-` cookie'ni saqlaydi va qaytaradi (jar.txt'dagi qator:
`#HttpOnly_localhost	FALSE	/	TRUE	…	__Host-session	…`).

### Origin va site

- **origin** = sxema + host + port: `http://localhost:8765` va `http://localhost:9090` — ikki origin.
- **site** = sxema + roʻyxatdan oʻtgan domen (eTLD+1), **port hisobga olinmaydi** (MDN Glossary: Site):
  `localhost:9090` bozor bilan **bir site**; `app.bozor.uz` va `blog.bozor.uz` ham bir site; `127.0.0.1` — boshqa host,
  demak boshqa site.
- `SameSite` site'ga qaraydi, `Sec-Fetch-Site`/`Origin` sarlavhalari esa origin'ni ham ajratadi.

## 2. CSRF — haqiqiy brauzerda

Begona sahifa (`cmd/begona`, `formaPage`) — skriptsiz oddiy forma:

```html
<form method="post" action="http://localhost:8765/orders/41/cancel">
  <button>Sovgʻani olish</button>
</form>
```

Server log'iga bu darsda ikki maydon qoʻshildi: `site=` — `Sec-Fetch-Site` sarlavhasi (brauzer qoʻshadi: `same-origin`,
`same-site`, `cross-site` yoki `none`), `cookie=` — session cookie keldimi.

| begona sahifa | cookie | natija (darsgacha: `protectCSRF`siz) | log |
|---|---|---|---|
| `http://localhost:9090` (bir site) | Lax | **200, buyurtma bekor qilindi** | `site="same-site" cookie=true` |
| `http://127.0.0.1:9090` (boshqa site) | Lax | 401 | `site="cross-site" cookie=false` |
| `http://127.0.0.1:9090` | **None** | **200, bekor qilindi** | `site="cross-site" cookie=true` |

```
$ go run . csrf http://localhost:9090        # darsgacha
login: 200 {"id":1,"email":"dilnoza@bozor.test"}
sahifa: http://localhost:9090 → «Sovgʻani olish» bosildi
javob: {"id":41,"user_id":1,"status":"bekor qilindi","note":"Chilonzor 9-kvartal, 12-uy"}
GET /orders/41: {"id":41,"user_id":1,"status":"bekor qilindi","note":"Chilonzor 9-kvartal, 12-uy"}
--- bozor log'i:
2026/11/28 10:42:43 UYHM36Y4EXZZ4WAVOL6A6TIBK7 POST "/orders/41/cancel" 200 53.587µs site="same-site" cookie=true
```

Egalik tekshiruvi (GB-6.06) bu yerda yordam bermaydi: soʻrov Dilnozaning oʻz brauzeridan, oʻz cookie'si bilan keladi.
OWASP: `SameSite` — qoʻshimcha qatlam (defense in depth), asosiy CSRF himoyasi emas; u roʻyxatdan oʻtgan domenga
qaraydi, origin'ga emas — qoʻshni port yoki subdomen uni chetlab oʻtadi.

### Himoya: `http.CrossOriginProtection` (`browser.go`, `protectCSRF`)

Go 1.25'dan beri standart kutubxonada; OWASP CSRF Prevention Cheat Sheet uni "Built-In Or Existing CSRF
Implementations" boʻlimida nomlaydi. `Check(r)`:

1. `GET`, `HEAD`, `OPTIONS` — har doim oʻtadi (ular holatni oʻzgartirmasligi kerak: OWASP, "Lax only blocks unsafe methods");
2. `Sec-Fetch-Site: same-origin` yoki `none` — oʻtadi; boshqa qiymat (`same-site`, `cross-site`) — rad;
3. `Sec-Fetch-Site` yoʻq, `Origin` bor — `Origin`ning host'i `Host` sarlavhasiga teng boʻlsa oʻtadi, aks holda rad;
4. ikkalasi ham yoʻq — brauzer emas (curl, mobil ilova, boshqa server): oʻtadi.

`AddTrustedOrigin` — ataylab ruxsat berilgan origin'lar (bozor'ning oʻz veb ilovasi, `http://localhost:3000`).
Rad etilsa — `403` problem (GB-6.02), sabab log'ga:

```
$ go run . csrf http://localhost:9090        # yakuniy
javob: {"type":"about:blank","title":"Forbidden","status":403,"detail":"boshqa saytdan kelgan soʻrov","request_id":"CMFJKVOFHSJMJ2PCUFRXVPZ7XQ"}
GET /orders/41: {"id":41,"user_id":1,"status":"yangi","note":"Chilonzor 9-kvartal, 12-uy"}
$ go run . csrf http://127.0.0.1:9090        # yakuniy
javob: {"type":"about:blank","title":"Forbidden","status":403,"detail":"boshqa saytdan kelgan soʻrov","request_id":"WRCLJEC4JIA6L6EOVIAGORF5HQ"}
--- bozor log'i:
… xato 403: "boshqa saytdan kelgan soʻrov: cross-origin request detected from Sec-Fetch-Site header"

$ curl -X POST -b jar.txt localhost:8765/orders/41/cancel
{"id":41,"user_id":1,"status":"bekor qilindi","note":"Chilonzor 9-kvartal, 12-uy"}
$ curl -X POST -b jar.txt -H 'Sec-Fetch-Site: same-site' -H 'Origin: http://localhost:9090' localhost:8765/orders/41/cancel
{"type":"about:blank","title":"Forbidden","status":403,"detail":"boshqa saytdan kelgan soʻrov",…}
```

Sarlavhasiz curl'ning oʻtishi xavfsiz: CSRF — brauzer cookie'ni **oʻzi** qoʻshgani uchun xavf; begona sayt birovning
curl'iga cookie qoʻsha olmaydi.

## 3. CORS — nima emas

CORS — brauzerning **oʻqish** qoidasi: boshqa origin'dagi sahifa soʻrov yubora oladi, lekin javobni oʻqishi uchun server
`Access-Control-Allow-Origin`da aynan shu origin'ni aytishi kerak (Fetch Standard, CORS protocol). Bu server himoyasi
emas: soʻrov baribir bajariladi.

```
$ go run . cors http://localhost:9090
login: 200 {"id":1,"email":"dilnoza@bozor.test"}
sahifa: http://localhost:9090/oqish
oʻqib boʻlmadi: TypeError: Failed to fetch
--- bozor log'i:
2026/11/28 10:43:00 CJ3ZWK7JNYUASOQBSA2MFXUPA4 GET "/me" 200 32.95µs site="same-site" cookie=true
```

**Tez-tez uchraydigan xato** (`_video/stages/aks/browser.go`): roʻyxat oʻrniga soʻrovdagi `Origin`ni qaytarish,
`Access-Control-Allow-Credentials: true` bilan. Natija — har qanday sayt foydalanuvchi nomidan oʻqiydi:

```
$ go run . cors http://localhost:9090        # xato versiya
oʻqildi: {"id":1,"email":"dilnoza@bozor.test"}
$ curl -i -b jar.txt -H 'Origin: https://istalgan-sayt.example' localhost:8765/me
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: https://istalgan-sayt.example
…
```

`protectCSRF` bu yerda yordam bermaydi: GET'ni u doim oʻtkazadi. `Access-Control-Allow-Origin: *` bilan
`Allow-Credentials` brauzer tomonidan rad etiladi (OWASP CSRF, "Custom Headers and CORS"), shuning uchun bu xato
koʻpincha «*'ni almashtirish» uchun qilinadi.

**Toʻgʻrisi** (`browser.go`, `cors`): allow-list (`allowedOrigins` — `protectCSRF` bilan bitta roʻyxat), har javobda
`Vary: Origin`:

```
$ curl -i -b jar.txt -H 'Origin: http://localhost:3000' localhost:8765/me
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: http://localhost:3000
Cache-Control: no-store
Content-Type: application/json
Vary: Origin
…
$ curl -i -b jar.txt -H 'Origin: http://localhost:9090' localhost:8765/me
HTTP/1.1 200 OK
Cache-Control: no-store
Content-Type: application/json
Vary: Origin
…
```

**Preflight:** JSON body, `PATCH`, `DELETE` yoki maxsus sarlavhali soʻrovdan oldin brauzer `OPTIONS` bilan ruxsat
soʻraydi; javob ruxsat bermasa, haqiqiy soʻrov yuborilmaydi (Fetch Standard, CORS-preflight request):

```
$ curl -i -X OPTIONS -H 'Origin: http://localhost:3000' -H 'Access-Control-Request-Method: POST' -H 'Access-Control-Request-Headers: content-type' localhost:8765/login
HTTP/1.1 204 No Content
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type
Access-Control-Allow-Methods: POST, PATCH, DELETE
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Max-Age: 600
Vary: Origin
…
$ curl -i -X OPTIONS -H 'Origin: http://localhost:9090' -H 'Access-Control-Request-Method: POST' -H 'Access-Control-Request-Headers: content-type' localhost:8765/login
HTTP/1.1 405 Method Not Allowed
Allow: POST
…
```

## 4. Login'ga rate limit (`ratelimit.go`)

Har login argon2id hisoblaydi (GB-6.03: 46 MiB). Darsgacha cheklov yoʻq:

```
$ for i in $(seq 8); do curl -s -o /dev/null -w "%{http_code} " localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"taxmin-1"}'; done; echo
401 401 401 401 401 401 401 401
```

[`golang.org/x/time/rate`](https://pkg.go.dev/golang.org/x/time/rate) (Go jamoasining token bucket'i, `v0.16.0`):
har IP'ga chelak — 5 token, har 12 soniyada +1; token yoʻq → `429 Too Many Requests` (RFC 6585 §4) va `Retry-After`
(RFC 9110 §10.2.3). `wait`: `Reserve` navbat oladi, `Delay` — kutish; kutmaymiz, `Cancel` tokenni qaytaradi.

```
$ for i in $(seq 8); do curl … ; done; echo
401 401 401 401 401 429 429 429
$ curl -i localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"taxmin-1"}'
HTTP/1.1 429 Too Many Requests
Content-Type: application/problem+json
Retry-After: 12
…
{"type":"about:blank","title":"Too Many Requests","status":429,"detail":"urinishlar koʻp: keyinroq qayta urinib koʻring",…}
# 13 soniyadan keyin:
401
```

Kalit — `r.RemoteAddr`dagi IP (TCP ulanishning manzili). `X-Forwarded-For`ga ishonmaymiz: uni mijoz yozadi. Server proksi
(Caddy) ortida boʻlsa, `RemoteAddr` — proksining manzili: haqiqiy IP'ni proksidan olish GB-8.06'da.

## Vazifa

[`TASK.md`](TASK.md): email boʻyicha ikkinchi limit (hujumchi IP almashtirsa ham), `go test ./...` oʻtguncha.
Javob: [`solution/`](solution/). Farqni koʻrish: `diff -ru cmd solution/cmd`.

## Chuqurroq

- **Sarlavhasiz soʻrovlar.** `CrossOriginProtection` `Sec-Fetch-Site` ham, `Origin` ham boʻlmasa soʻrovni oʻtkazadi
  (fail-open): curl va mobil ilovalar ishlashi uchun. OWASP sezgir endpoint'lar uchun fail-safe'ni ham taklif qiladi. `Origin` ↔ `Host` solishtiruvida sxema koʻrilmaydi — HTTPS + HSTS kerak (GB-8.06);
  zamonaviy brauzerlar `Sec-Fetch-Site`ni 2023-yildan beri yuboradi, `Origin`ni esa POST'da har doim.
- **Fetch Metadata faqat ishonchli URL'larga** yuboriladi: `https`, `localhost`, `127.0.0.0/8` (OWASP, W3C Secure
  Contexts). Production'da HTTPS shart (GB-8.06); oddiy HTTP'da faqat `Origin` tekshiruvi qoladi.
- **`SameSite`siz cookie.** Chromium uni `Lax` deb hisoblaydi va yangi cookie'ni ~2 daqiqa davomida cross-site POST'ga
  ham qoʻshadi ("Lax-allowing-unsafe", RFC 6265bis-22 §5.6.7.2); Firefox va Safari default'da `Lax` qoʻymaydi. Shuning
  uchun `SameSite`ni har doim aniq yozing.
- **CSRF token** (synchronizer / signed double-submit) — eski brauzerlar yoki boshqa site'dan ataylab qabul qilinadigan
  formalar uchun klassik yoʻl; Fetch Metadata + Origin tekshiruvi uni zamonaviy brauzerlar uchun almashtiradi (OWASP).
- **XSS hamma CSRF himoyalarini yengadi** (OWASP): sahifaga tushgan skript bir xil origin'dan soʻrov yuboradi.
  `HttpOnly` faqat tokenni oʻgʻirlashni toʻxtatadi.
- **Limitlar xotirada** — bitta server jarayoniga tegishli; bir nechta nusxa boʻlsa, umumiy hisob (masalan Redis) kerak.
- **Map oʻsishi.** Kalitlar (IP, email) mijozdan keladi: `sweep`ni (TASK.md 🟡) production'da vaqti-vaqti bilan chaqirish
  shart (goroutine va `time.Ticker` — GB-7.01). Email uzunligi `decodeJSON`ning 1024 baytlik limiti bilan cheklangan.
- **IP kalitining chegaralari.** Bitta NAT ortidagi koʻp foydalanuvchi bitta chelakni boʻlishadi; IPv6'da bitta mijozda
  odatda butun /64 tarmoq boʻladi (production'da kalitni /64 boʻyicha olish mumkin); proksi ortida `RemoteAddr` —
  proksining manzili (GB-8.06).
- **Lockout DoS.** Email limitini hujumchi Dilnozani bloklash uchun ham ishlatishi mumkin (TASK.md maslahatlari, OWASP
  "Account Lockout").

## Manbalar

- OWASP Cross-Site Request Forgery Prevention Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html
- OWASP Authentication Cheat Sheet (Login Throttling, Account Lockout) — https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#login-throttling
- Go: `net/http.CrossOriginProtection` — https://pkg.go.dev/net/http#CrossOriginProtection · Go 1.25 release notes — https://go.dev/doc/go1.25#nethttppkgnethttp
- RFC 6265bis (draft-ietf-httpbis-rfc6265bis-22) — https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-rfc6265bis
- MDN: Set-Cookie — https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie · Site — https://developer.mozilla.org/en-US/docs/Glossary/Site · Sec-Fetch-Site — https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Sec-Fetch-Site · CORS — https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS
- Fetch Standard: CORS protocol — https://fetch.spec.whatwg.org/#http-cors-protocol
- `golang.org/x/time/rate` — https://pkg.go.dev/golang.org/x/time/rate
- RFC 6585 §4 (429) — https://www.rfc-editor.org/rfc/rfc6585#section-4 · RFC 9110 §10.2.3 (Retry-After) — https://www.rfc-editor.org/rfc/rfc9110#field.retry-after
