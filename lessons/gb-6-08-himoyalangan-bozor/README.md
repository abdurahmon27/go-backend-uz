# GB-6.08 · Nazorat nuqtasi: himoyalangan bozor

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #50"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

6-modulning nazorat nuqtasi. GB-5.10'dagi Postgres'dagi katalog va buyurtmalar API'si endi 6-modulning hamma qatlami bilan
bitta API: roʻyxatdan oʻtish va argon2id (GB-6.03), login, logout va bekor qilinadigan session'lar (GB-6.04), `requireUser`,
egasi SQL shartida tekshiriladigan buyurtmalar va bir xil `404` (GB-6.06), faqat admin oʻzgartiradigan katalog (GB-6.06),
RFC 9457 problem javoblari (GB-6.02), middleware zanjiri (GB-6.01), CORS, CSRF tekshiruvi va login'ga rate limit
(GB-6.07). Testlar mock'siz, haqiqiy Postgres'da, va qoidalarni **ikki xaridor, admin va anonim** bilan isbotlaydi.

Oldingi darslar: [`gb-6-07-cookie-csrf-cors`](../gb-6-07-cookie-csrf-cors) (cookie, CSRF, CORS, rate limit) ·
[`gb-6-06-buyurtma-kimniki`](../gb-6-06-buyurtma-kimniki) (egalik, `requireAdmin`) ·
[`gb-5-10-buyurtma-api`](../gb-5-10-buyurtma-api) (5-modul nazorat nuqtasi) · keyingi dars: `gb-7-01` (goroutine).
Bu papka **GB-6.06'ning `solution/`idan** boshlanadi va unga GB-6.07'ning `ratelimit.go` (IP va email limiti),
`browser.go` (`cors`, `protectCSRF`) va cookie oʻzgarishlarini (`__Host-session`, `Max-Age`) qoʻshadi. Farqlarni koʻrish:
`diff -ru ../gb-6-06-buyurtma-kimniki/solution solution`.

## Spetsifikatsiya: kim nima qila oladi

| guruh | route'lar | himoya | yoʻq boʻlsa |
|---|---|---|---|
| ochiq | `GET /health`, `GET /products`, `GET /products/{id}`, `POST /users`, `POST /login`, `POST /logout` | `POST /login`: IP boʻyicha 5 urinish, keyin har 12 s bittadan; email boʻyicha 10 urinish, keyin har daqiqada bittadan | `429` + `Retry-After` |
| login kerak | `GET /me`, `POST /logout-all`, `POST /orders`, `GET /orders/{id}`, `PATCH /orders/{id}`, **`GET /orders`** (vazifa) | `requireUser`: cookie → session qatori → foydalanuvchi | `401` |
| buyurtma | `GET`/`PATCH /orders/{id}` | egasi SQL shartida: `WHERE id = $1 AND user_id = $2` — admin uchun ham | `404` (begona = yoʻq) |
| admin | `GET /admin/users`, `POST /admin/products`, `DELETE /admin/products/{id}`, **`GET /admin/orders`** (vazifa) | alohida `admin` mux, **bir marta** oʻralgan: `requireUser` ichida `requireAdmin` | `401` / `403` |
| parol | `users.password_hash` | faqat argon2id satri (`$argon2id$v=19$m=47104,t=1,p=1$…`) | — |
| brauzer | hamma route | `protectCSRF` (`http.CrossOriginProtection`): begona origin'dan POST/PATCH/DELETE → `403`; `cors`: faqat `http://localhost:3000` | `403` |

Har xato javobi `application/problem+json` (GB-6.02). Pul tiyinda (`int64`): `3600000` = 36 000 soʻm.

**Nima oʻzgardi (GB-6.06'dan):** katalogni oʻzgartirish `POST /products` va `DELETE /products/{id}`dan
`POST /admin/products` va `DELETE /admin/products/{id}`ga koʻchdi. GB-6.06'da har admin route alohida oʻralardi; birini
unutsangiz, oddiy xaridor unga kirardi (OWASP API5:2023). Endi `/admin/` ostidagi hamma route bitta `ServeMux`da:

```go
admin := http.NewServeMux()
admin.HandleFunc("GET /admin/users", a.listUsers)
admin.Handle("POST /admin/products", requireJSON(http.HandlerFunc(a.createProduct)))
admin.HandleFunc("DELETE /admin/products/{id}", a.deleteProduct)
mux.Handle("/admin/",
	a.requireUser(requireAdmin(muxErrors(admin))))
```

Pattern `"/admin/"` slash bilan tugaydi va metodsiz: shu yoʻl ostidagi hamma yoʻl, har metod
([net/http ServeMux](https://pkg.go.dev/net/http#ServeMux): "a pattern that ends in a slash … matches all paths that
begin with it"). Ichki mux'da yoʻl toʻliq yoziladi (`StripPrefix` yoʻq). Natija — yoʻq admin yoʻli ham avval
`401`/`403`dan oʻtadi: `GET /admin/yoq` anonimga `401`, xaridorga `403`, faqat admin'ga `404`. Admin boʻlmagan
foydalanuvchi admin boʻlimida nima borligini bilmaydi. Tashqi `muxErrors` `/admin/` pattern'ini topadi, shuning uchun
ichki mux'ning oʻz `404`/`405`lari uchun `muxErrors(admin)` ham kerak.

Zanjir (`cmd/api/middleware.go`), tashqaridan ichkariga:

```go
return requestID(
	logRequests(
		recoverPanic(
			cors(
				protectCSRF(
					muxErrors(newMux(st)))))))
```

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, curl 8.22.0, goose v3.30.2,
PostgreSQL 18 `postgres:18-alpine`). Toʻliq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt) (`_video/chiqishlar.sh`).
Skript `curl`ni `-s` bilan ishlatadi (chiqish faylga ketadi); terminalda `curl -i` aynan shu qatorlarni chiqaradi.
goose har qator boshiga sana va vaqtni yozadi; quyida (va videoda) u olib tashlangan. `…` — qisqartirilgan sarlavhalar.

```bash
docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5518/bozor
goose up                                         # 00001 … 00007, "successfully migrated database to version: 7"
docker compose exec db psql -f /lab/seed.sql     # 3 foydalanuvchi, 4 mahsulot, buyurtmalar 41 (Dilnoza) va 42 (Jasur)
export DATABASE_URL=$GOOSE_DBSTRING
go run ./cmd/api                                 # bozor API: http://localhost:8765
```

Ikkinchi terminalda bitta toʻliq sessiya — yangi foydalanuvchi Aziza:

```
$ curl -i localhost:8765/users -H "Content-Type: application/json" -d '{"email":"aziza@bozor.test","password":"anor-behi-yongoq-uzum-5"}'
HTTP/1.1 201 Created
…
{"id":4,"email":"aziza@bozor.test","role":"xaridor"}

$ docker compose exec db psql -c "SELECT id, email, left(password_hash, 31) FROM users WHERE id = 4"
 id |      email       |              left
----+------------------+---------------------------------
  4 | aziza@bozor.test | $argon2id$v=19$m=47104,t=1,p=1$
(1 row)

$ curl -s -c aziza.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"aziza@bozor.test","password":"anor-behi-yongoq-uzum-5"}'
{"id":4,"email":"aziza@bozor.test","role":"xaridor"}

$ curl -i -b aziza.txt localhost:8765/orders/41          # Dilnozaniki
HTTP/1.1 404 Not Found
…
{"type":"about:blank","title":"Not Found","status":404,"detail":"bunday buyurtma yoʻq","request_id":"LBOT75RM3YIOWF44JDFNIA55UX"}

$ curl -i -b aziza.txt localhost:8765/admin/products -H "Content-Type: application/json" -d '{"name":"anor","price":1800000,"stock":5}'
HTTP/1.1 403 Forbidden
…
{"type":"about:blank","title":"Forbidden","status":403,"detail":"bu amal uchun ruxsat yoʻq","request_id":"NLQZXEYBIUKJRJCUMASJC6RS45"}

$ curl -i localhost:8765/admin/yoq                       # cookie'siz, yoʻq yoʻl
HTTP/1.1 401 Unauthorized
$ curl -i -b aziza.txt localhost:8765/admin/yoq          # xaridor
HTTP/1.1 403 Forbidden

$ curl -i -b aziza.txt localhost:8765/orders -H "Content-Type: application/json" -d '{"items":[{"product_id":3,"quantity":1}]}'
HTTP/1.1 201 Created
Location: /orders/43
…
{"id":43,"user_id":4,"status":"yangi","items":[{"product_id":3,"quantity":1,"unit_price":2500000}],"total":2500000}

$ curl -i -b aziza.txt -X POST localhost:8765/logout
HTTP/1.1 204 No Content
Set-Cookie: __Host-session=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Lax
$ curl -i -b aziza.txt localhost:8765/orders/43          # aziza.txt'da token turibdi (-b faylni faqat oʻqiydi)
HTTP/1.1 401 Unauthorized
…
{"type":"about:blank","title":"Unauthorized","status":401,"detail":"avval login qiling","request_id":"D6CVCHG64WGXJ44AUXBKXSSEKF"}

$ for p in x1 x2 x3 x4 x5 anor-behi-yongoq-uzum-5; do curl -s -o /dev/null -w "%{http_code} " localhost:8765/login -H "Content-Type: application/json" -d '{"email":"aziza@bozor.test","password":"'$p'"}'; done; echo
401 401 401 401 401 429
$ curl -i localhost:8765/login -H "Content-Type: application/json" -d '{"email":"aziza@bozor.test","password":"anor-behi-yongoq-uzum-5"}'
HTTP/1.1 429 Too Many Requests
Content-Type: application/problem+json
Retry-After: 12
```

Oltinchi urinish **toʻgʻri** parol bilan, lekin `429`: limit urinishlarni sanaydi, parol natijasini emas. Aks holda
hujumchi `401` va `200` farqini koʻrib, taxminni davom ettirardi (OWASP Authentication Cheat Sheet, "Protect Against
Automated Attacks"). Skript bu sikldan oldin 13 s kutadi: Aziza'ning login'i olgan token chelakka qaytsin.

## Testlar: qoidalarning isboti

`cmd/api/himoya_test.go` (haqiqiy Postgres, har test oʻz sxemasida, `internal/pgtest`):

```go
type rule struct {
	who                string // "anonim" (cookie'siz), "dilnoza", "jasur", "admin"
	method, path, body string
	want               int
}

var accessRules = []rule{
	{"anonim", "GET", "/orders/41", "", 401},
	{"dilnoza", "GET", "/orders/42", "", 404},
	{"jasur", "GET", "/orders/41", "", 404},
	{"admin", "GET", "/orders/41", "", 404}, // /orders/{id} — «mening buyurtmam»: admin uchun ham
	{"dilnoza", "POST", "/admin/products", productBody, 403},
	{"dilnoza", "GET", "/admin/yoq", "", 403},
	// … jami 16 qator + vazifaning 3 qatori
}
```

| test | nimani isbotlaydi |
|---|---|
| `TestAccess` | jadvalning har qatori butun zanjirdan (`newHandler`) oʻtadi: 401 · 403 · 404 · 409 · 200/201 |
| `TestPasswordNotStored` | `POST /users`dan keyin `users.password_hash`da parol yoʻq, satr `$argon2id$` bilan boshlanadi; javobda hash yoʻq; login ishlaydi |
| `TestLogoutRevokes` | logout'dan keyin oʻsha cookie bilan `GET /orders/41` → `401`; `sessions` jadvali boʻsh |
| `TestLoginRateLimit` | bitta IP'dan 5 × `401`, 6-urinish toʻgʻri parol bilan → `429` va `Retry-After` |

```
$ go test -v -run 'TestAccess$|TestPasswordNotStored|TestLogoutRevokes|TestLoginRateLimit' ./cmd/api
=== RUN   TestAccess
--- PASS: TestAccess (0.20s)
=== RUN   TestPasswordNotStored
--- PASS: TestPasswordNotStored (0.14s)
=== RUN   TestLogoutRevokes
--- PASS: TestLogoutRevokes (0.11s)
=== RUN   TestLoginRateLimit
--- PASS: TestLoginRateLimit (0.21s)
PASS
ok  	github.com/abdurahmon27/bozor/cmd/api	0.708s
```

`api_test.go` — GB-6.04/6.06'dan qolgan testlar (cookie atributlari, login xatolari, `/me`, logout, qayta login, muddati
oʻtgan session, `requireJSON`, `logout-all`, begona buyurtmaning javobi yoʻq buyurtmanikiga teng, mass assignment `400`,
`requireAdmin` tashqarida → `500`); `note_test.go` — GB-6.06 vazifasi (`PATCH` faqat egasiga); store'da `store_test.go`,
`sessions_test.go`.

**Tez-tez uchraydigan xato: himoya testini bitta foydalanuvchi bilan yozish** (`_video/stages/bola`). `bitta_test.go`dagi
`TestOrdersOneUser` faqat Dilnoza bilan: oʻz buyurtmasi `200`, yoʻq buyurtma `404`, cookie'siz `401`. Store'dan
`AND user_id = $2` olib tashlansa (BOLA, OWASP API1:2023):

```
$ go test -run TestOrdersOneUser ./cmd/api
ok  	github.com/abdurahmon27/bozor/cmd/api	0.155s
$ go test -run 'TestAccess$' ./cmd/api
--- FAIL: TestAccess (0.17s)
    himoya_test.go:61: dilnoza GET /orders/42: 200, kutilgan 404
    himoya_test.go:61: jasur GET /orders/41: 200, kutilgan 404
    himoya_test.go:61: admin GET /orders/41: 200, kutilgan 404
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api	0.218s
```

Bitta foydalanuvchi begona buyurtmani hech qachon soʻramaydi: BOLA faqat ikkinchi foydalanuvchi bilan koʻrinadi. OWASP
Authorization Cheat Sheet ham ruxsat mantigʻi uchun unit va integration testlarni talab qiladi ("Create Unit and Integration
Test Cases for Authorization Logic"): har rol va har "begona obyekt" holati alohida qator.

## Ishga tushirish

Kerak: Go 1.27, Docker (Compose 2.20+), curl, goose (`go install github.com/pressly/goose/v3/cmd/goose@v3.30.2`, GB-5.06).
Postgres `127.0.0.1:5518`da (shu papkaning `compose.yaml`i, `gb608` loyihasi, nomli volume), API `:8765`da. Buyruqlar
yuqorida. Laboratoriya parollari (faqat shu yerda; bazada argon2id hash'i):

| email | parol | rol | buyurtma |
|---|---|---|---|
| `dilnoza@bozor.test` | `tarvuz-osmon-qalam-ariq-27` | xaridor | 41 |
| `jasur@bozor.test` | `choy-non-qand-27` | xaridor | 42 |
| `admin@bozor.test` | `qovun-anjir-tut-uzum-31` | admin | — |

Login: `curl -s -c dilnoza.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'`,
keyin `curl -s -b dilnoza.txt localhost:8765/orders/41`. Toʻxtatish va maʼlumotni oʻchirish: `docker compose down -v`.

**GB-6.07'dan farqi:** GB-6.07 laboratoriyasi xotiradagi store va `POST /orders/{id}/cancel` (CSRF misoli uchun) bilan
edi; bu yerda store Postgres'da (GB-6.06) va bekor qilish route'i yoʻq — CSRF himoyasi (`protectCSRF`) zanjirda hamma
route uchun turibdi.

**Bu laboratoriyadagi cheklovlar (keyingi modullar):** email limiti parolni tekshirishdan oldin ishlaydi: kimdir bir
daqiqada 10 marta Dilnozaning email'i bilan urinsa, Dilnozaning oʻzi ham bir muddat `429` oladi (hisobni bloklash xavfi,
OWASP Authentication Cheat Sheet "Account Lockout"; muqobil — faqat muvaffaqiyatsiz urinishlarni sanash va kechikishni
oshirib borish). `GET /admin/orders` roʻyxatni toʻliq qaytaradi — ishlab turgan tizimda sahifalash kerak. limit chelaklari bitta server jarayonining xotirasida — bir nechta
nusxada har biri oʻz hisobini yuritadi (umumiy hisob: Redis yoki Postgres, GB-8 dan keyin); chelaklar map'dan
oʻchirilmaydi (GB-6.07'ning 🟡 `sweep` vazifasi, uni chaqiradigan goroutine — GB-7.01); proksi ortida haqiqiy IP —
GB-8.06. CORS roʻyxatida faqat laboratoriyaning `http://localhost:3000` veb ilovasi (GB-6.07).

## Papka

```
gb-6-08-himoyalangan-bozor/
  go.mod, go.sum              module github.com/abdurahmon27/bozor, go 1.27; pgx v5.13.0, goose v3.30.2,
                              golang.org/x/crypto v0.59.0, golang.org/x/time v0.16.0 (GB-6.07)
  compose.yaml                Postgres 18, 127.0.0.1:5518, nomli volume (pgdata), loyiha gb608
  migrations/                 00001…00007 (GB-5.10, GB-6.04, GB-6.06) — oʻzgarishsiz
  seed.sql                    3 foydalanuvchi, 4 mahsulot, buyurtmalar 41 (Dilnoza) va 42 (Jasur)
  migrate.go                  migratsiyalar embed bilan (testlar uchun, GB-5.06)
  cmd/api/handlers.go         route'lar: ochiq, requireUser, admin mux (YANGI); VAZIFA: GET /orders, GET /admin/orders
  cmd/api/middleware.go       zanjir: requestID → logRequests → recoverPanic → cors → protectCSRF → muxErrors → mux
  cmd/api/auth.go, users.go   login (+ email limiti), requireUser, logout, register (GB-6.03/6.04/6.07)
  cmd/api/authz.go            requireAdmin, errForbidden (GB-6.06)
  cmd/api/ratelimit.go        limiters, limitByIP, tooManyRequests (GB-6.07)
  cmd/api/browser.go          cors, protectCSRF (GB-6.07, oʻzgarishsiz)
  cmd/api/problem.go, json.go problem+json (GB-6.02), decodeJSON (GB-4.06)
  cmd/api/himoya_test.go      YANGI: accessRules, TestAccess, TestPasswordNotStored, TestLogoutRevokes, TestLoginRateLimit
  cmd/api/tekshiruv_test.go   vazifa tekshiruvi — oʻzgartirmang
  cmd/api/api_test.go, note_test.go   GB-6.04/6.06 testlari
  internal/store              users, sessions (store.go), mahsulotlar va buyurtmalar (orders.go; VAZIFA: OrderSummary, Orders, AllOrders)
  internal/password, pgtest   GB-6.03, GB-5.10
  TASK.md                     vazifa
  solution/                   javob (oʻz go.mod'i bilan, alohida ishlaydi)
  _video/                     chiqishlar.sh → chiqishlar.txt (videodagi hamma buyruq va chiqish),
                              stages/bola (BOLA'li store va bitta foydalanuvchili test — videodagi xato)
```

## Vazifa

[`TASK.md`](TASK.md): `GET /orders` (faqat oʻz buyurtmalari, `[]` boʻlsa ham) va `GET /admin/orders` (admin mux'da,
oʻramsiz), va ularning `accessRules`dagi qatorlari. `go test ./...` oʻtguncha.

## Manbalar

- OWASP API Security Top 10 2023: [API1 BOLA](https://owasp.org/API-Security/editions/2023/en/0xa1-broken-object-level-authorization/),
  [API2 Broken Authentication](https://owasp.org/API-Security/editions/2023/en/0xa2-broken-authentication/) (login'ga
  rate limit), [API5 BFLA](https://owasp.org/API-Security/editions/2023/en/0xa5-broken-function-level-authorization/)
- [OWASP Authorization Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authorization_Cheat_Sheet.html)
  (deny by default, validate on every request, authorization tests)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html),
  [Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html),
  [Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110) §15.5.2 401, §15.5.4 403 (404 bilan mavjudlikni yashirish),
  §10.2.3 Retry-After; [RFC 6585 §4](https://www.rfc-editor.org/rfc/rfc6585#section-4) 429;
  [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) problem details
- Go: [net/http ServeMux](https://pkg.go.dev/net/http#ServeMux),
  [CrossOriginProtection](https://pkg.go.dev/net/http#CrossOriginProtection),
  [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate), [testing](https://pkg.go.dev/testing)
- [PostgreSQL 18: string functions](https://www.postgresql.org/docs/18/functions-string.html),
  [aggregate functions](https://www.postgresql.org/docs/18/functions-aggregate.html)
- Hamma chiqishlar shu kompyuterda: go1.27.4 linux/amd64, `CGO_ENABLED=0`, curl 8.22.0, goose v3.30.2, PostgreSQL 18
  (`postgres:18-alpine`), `_video/chiqishlar.sh`.
