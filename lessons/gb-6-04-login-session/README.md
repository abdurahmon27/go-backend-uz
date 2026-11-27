# GB-6.04 · Login va bekor qilinadigan session'lar

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #46"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-6-03-parollar`](../gb-6-03-parollar) (argon2id, `internal/password`, `POST /users`) · keyingi dars:
`gb-6-05-session-jwt` (session yoki JWT). Postgres qismi va testlar usuli: [`gb-5-10-buyurtma-api`](../gb-5-10-buyurtma-api).

**Bu darsdan keyin:** bozor API'sida login ishlaydi. `POST /login` parolni argon2id bilan tekshiradi (GB-6.03), 32 bayt
tasodifiy **session tokeni** yasaydi (`crypto/rand`), Postgres'ga esa tokenning o'zini emas, faqat **SHA-256 hash'ini**
yozadi (`sessions` jadvali, muddat bilan) va tokenni cookie'da qaytaradi. `requireUser` middleware'i (GB-6.01 shakli)
har so'rovda cookie'ni hash qilib session egasini topadi va uni so'rov **context**'iga qo'yadi. `POST /logout` qatorni
o'chiradi: o'sha token bilan keyingi so'rov darhol `401`.

**Nima o'zgardi (GB-6.03 → GB-6.04).** GB-6.01–6.03'da foydalanuvchilar xotirada edi (`store.Users`), chunki u darslar
HTTP qatlami va parol haqida edi. Session — bazadagi qator, shuning uchun bu laboratoriya Postgres'da: GB-5.10'ning
`compose.yaml`, goose migratsiyalari 00001–00003 (o'zgarishsiz), `internal/pgtest`, `run()` + `pgxpool`; yangi
`00004_users.sql` (GB-6.03'dagi `users`: `id`, `email`, `password_hash`; email katta-kichik harfsiz unikal, GB-6.03'dagi
`EqualFold` qoidasi) va `00005_sessions.sql`. `POST /users` (GB-6.03'dagi `register`, `validateUser`, 422/409) endi
`store.AddUser` orqali shu jadvalga yozadi. `internal/password` va `cmd/api/middleware.go` — GB-6.03'ning `solution/`idagi
fayllar, **o'zgarishsiz** (`requireJSON` faqat body'li route'larda: `POST /users`, `POST /login`). `problem.go` — GB-6.02
shakli, shu laboratoriyadagi xatolar bilan (422, 409, 401, 413, 400, 500). Katalog va buyurtmalar kodi bu papkada yo'q (ular
GB-5.10'da Postgres ustida); hammasini bitta API'ga yig'ish — GB-6.08 nazorat nuqtasi.

## Spetsifikatsiya

| route | javob |
|---|---|
| `GET /health` | `200 ok` |
| `POST /users` `{"email","password"}` | `201 {"id","email"}` · `422` (email, parol 15–64 belgi) · `409` email band (GB-6.03) |
| `POST /login` `{"email","password"}` | `200` + `Set-Cookie: session=<token>; Path=/; Expires=…; HttpOnly; Secure; SameSite=Lax` + `Cache-Control: no-store` + `{"id","email"}` · `401` «email yoki parol notoʻgʻri» (email yo'q bo'lsa ham xuddi shu) · `400` · `415` |
| `GET /me` | `200 {"id","email"}` · `401` «avval login qiling» (cookie yo'q, token notanish yoki muddati o'tgan) |
| `POST /logout` | `204` + `Set-Cookie: session=; …; Max-Age=0` (session bo'lmasa ham: idempotent) |
| `POST /logout-all` | **vazifa** ([TASK.md](TASK.md)): `204`, foydalanuvchining hamma session'lari o'chadi · `401` |

Har `401`da `WWW-Authenticate: Cookie realm="bozor"` (pastda, "Chuqurroq").

Laboratoriya foydalanuvchilari (`seed.sql`, GB-6.03'dagi ismlar; parollar faqat shu yerda ishlatiladi):
`dilnoza@bozor.test` / `tarvuz-osmon-qalam-ariq-27`, `jasur@bozor.test` / `choy-non-qand-27`.

## Papka

| fayl | nima |
|---|---|
| `compose.yaml` | Postgres 18 (`postgres:18-alpine`), loyiha `gb604`, `127.0.0.1:5514`, nomli volume |
| `migrations/` | goose: 00001–00003 (GB-5.10), `00004_users.sql`, `00005_sessions.sql` |
| `seed.sql` | `TRUNCATE sessions, users RESTART IDENTITY` + ikki foydalanuvchi (argon2id satrlari) |
| `migrate.go` | `package bozor`: migratsiyalar va seed `go:embed` bilan, `NewProvider` (testlar uchun, GB-5.10) |
| `cmd/api/auth.go` | `newToken`, `hashToken`, `login`, `setSessionCookie`, `clearSessionCookie`, `checkPassword`, `requireUser`, `currentUser`, `endCurrentSession`, `logout` |
| `cmd/api/users.go` | `POST /users` (GB-6.03, endi Postgres'ga) |
| `cmd/api/handlers.go` | route'lar va `me` |
| `cmd/api/problem.go`, `json.go`, `middleware.go` | GB-6.02/6.01/6.03'dan |
| `internal/store/` | `AddUser`, `UserByEmail`, `CreateSession`, `UserBySession`, `DeleteSession`, `DeleteUserSessions` (VAZIFA), `DeleteExpiredSessions` (🟡) |
| `internal/password/` | argon2id `Hash`/`Verify` (GB-6.03) |
| `internal/pgtest/` | har test uchun yangi sxema (GB-5.10) |
| `solution/` | to'liq javob, o'z `go.mod`i bilan |
| `_video/` | `chiqishlar.sh` → `chiqishlar.txt` (videodagi hamma buyruq va chiqish, video tartibida); `stages/store-muddatsiz.go` — videodagi XATO versiya |

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan: `GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, curl 8.22.0, goose v3.30.2,
pgx v5.13.0, PostgreSQL 18. To'liq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt). Skript chiqishni faylga
yozgani uchun `curl -s` va `docker compose exec -T` bilan ishlaydi; terminalda `curl -i` va `exec db` aynan shu qatorlarni
chiqaradi. Server videoda `go run ./cmd/api` bilan, skriptda esa xuddi shu kod `go build` qilinib ishga tushadi.
Token, hash va request ID har ishga tushirishda boshqa bo'ladi. goose CLI: `go install github.com/pressly/goose/v3/cmd/goose@v3.30.2` (GB-5.06).

```bash
docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5514/bozor
goose up
# …
# OK   00004_users.sql (4.58ms)
# OK   00005_sessions.sql (5.14ms)
# goose: successfully migrated database to version: 5
docker compose exec db psql -f /lab/seed.sql
# TRUNCATE TABLE · INSERT 0 2
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5514/bozor
go run ./cmd/api
# bozor API: http://localhost:8765
```

Ikkinchi terminalda:

```
$ curl -i localhost:8765/me
HTTP/1.1 401 Unauthorized
Content-Type: application/problem+json
Www-Authenticate: Cookie realm="bozor"
…
{"type":"about:blank","title":"Unauthorized","status":401,"detail":"avval login qiling","request_id":"RKK6WX5KJUHTMHT3OTFRONNRQQ"}

# notoʻgʻri parol va yoʻq email: bir xil javob
$ curl -s localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-28"}'
{"type":"about:blank","title":"Unauthorized","status":401,"detail":"email yoki parol notoʻgʻri","request_id":"PRCW5K7GMTSUNQHHIEW52EJTOX"}
$ curl -s localhost:8765/login -H "Content-Type: application/json" -d '{"email":"yoq@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'
{"type":"about:blank","title":"Unauthorized","status":401,"detail":"email yoki parol notoʻgʻri","request_id":"Q3KFT5KVRWYCAYXSDAAL47REEK"}

# login: telefon (cookie jar.txt'ga) va noutbuk
$ curl -i -c jar.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'
HTTP/1.1 200 OK
Cache-Control: no-store
Content-Type: application/json
Set-Cookie: session=MrYRxQh-6_Oah9BsQ_0-d3QbB_ZWOsMXko-YxfruN88; Path=/; Expires=Sat, 28 Nov 2026 17:34:05 GMT; HttpOnly; Secure; SameSite=Lax
…
{"id":1,"email":"dilnoza@bozor.test"}
$ curl -s localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'
{"id":1,"email":"dilnoza@bozor.test"}

$ docker compose exec db psql -c "SELECT encode(token_hash, 'hex') AS token_hash, user_id, expires_at FROM sessions"
                            token_hash                            | user_id |          expires_at
------------------------------------------------------------------+---------+-------------------------------
 65f384e24f4858c7387680bf412f3a550d2db4873249a4ca2fc33d7a01d70249 |       1 | 2026-11-28 17:34:05.440928+00
 bbff80b15ffcc4a175e43ead82346468029c3f6fc780277b977e50aed24d3db0 |       1 | 2026-11-28 17:34:05.475103+00
(2 rows)

# server qayta ishga tushdi (1-terminal: Ctrl+C, go run ./cmd/api): session bazada, cookie ishlaydi
$ curl -s -b jar.txt localhost:8765/me
{"id":1,"email":"dilnoza@bozor.test"}

# baza sizib chiqdi: jadvaldagi qiymatni cookie qilib yuboramiz
$ curl -s -b session=65f384e24f4858c7387680bf412f3a550d2db4873249a4ca2fc33d7a01d70249 localhost:8765/me
{"type":"about:blank","title":"Unauthorized","status":401,"detail":"avval login qiling","request_id":"LMRNJHJLPFMI632U4PTZWOSMRT"}
```

**Klassik xato** (`_video/stages/store-muddatsiz.go`): `UserBySession` so'rovida faqat `WHERE s.token_hash = $1` qolsa,
muddati o'tgan session ham ishlaydi. Shu versiya bilan server qayta ishga tushirildi, keyin:

```
$ docker compose exec db psql -c "UPDATE sessions SET expires_at = now() - interval '1 minute'"
UPDATE 2
$ curl -s -b jar.txt localhost:8765/me
{"id":1,"email":"dilnoza@bozor.test"}
```

Shart (`AND s.expires_at > now()`) qaytarilib, server yana ishga tushirilgach, o'sha so'rov:

```
$ curl -s -b jar.txt localhost:8765/me
{"type":"about:blank","title":"Unauthorized","status":401,"detail":"avval login qiling","request_id":"USAJPGROBYCYYJRD5AI62AJQNT"}
```

Cookie'dagi `Expires` himoya emas: uni mijoz boshqaradi, o'g'irlangan token nusxasini hujumchi istagancha yuboradi.
Muddatni server o'z vaqti (`now()`) bilan tekshiradi.

Logout (avval qayta login, keyin nusxa):

```
$ curl -s -c jar.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'
{"id":1,"email":"dilnoza@bozor.test"}
$ cp jar.txt eski.txt
$ curl -i -X POST -b jar.txt -c jar.txt localhost:8765/logout
HTTP/1.1 204 No Content
Set-Cookie: session=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Lax
…
$ curl -s -b eski.txt localhost:8765/me
{"type":"about:blank","title":"Unauthorized","status":401,"detail":"avval login qiling","request_id":"6JMPDTCMLHVZPD7P4PL2VWGPQP"}
```

Server log'ida har 401'ning sababi toʻliq (mijozga emas, GB-6.02):

```
… RKK6WX5KJUHTMHT3OTFRONNRQQ xato 401: "http: named cookie not present"
… PRCW5K7GMTSUNQHHIEW52EJTOX xato 401: "email yoki parol notoʻgʻri: foydalanuvchi 1"
… Q3KFT5KVRWYCAYXSDAAL47REEK xato 401: "email yoki parol notoʻgʻri: email \"yoq@bozor.test\": foydalanuvchi topilmadi"
… LMRNJHJLPFMI632U4PTZWOSMRT xato 401: "session yoʻq yoki muddati oʻtgan"
```

## Chuqurroq (videoga kirmagan yoki qisqa aytilgan)

- **Nega hash, va nega SHA-256, argon2id emas.** OWASP Session Management Cheat Sheet ("Server-Side Session Token Storage"):
  saqlash joyining o'qilib qolishi (backup, replika, SQL injection) xavfi bo'lsa, tokenning o'zini emas, bir tomonlama
  tekshiruv qiymatini saqlang; token yoki verifier CSPRNG'dan kamida 128 bit (afzali 160+). Bizda 256 bit. Argon2id inson
  paroli kabi **kam entropiyali** sirni taxmin qilishni qimmatlashtiradi; 256 bitli tasodifiy tokenni taxmin qilib
  bo'lmaydi, shuning uchun bitta tez SHA-256 yetadi (u har so'rovda ishlaydi).
- **Doimiy vaqt (GB-6.03'dagi va'da).** Bazadagi `WHERE token_hash = $1` doimiy vaqtda solishtirmaydi. Lekin solishtiriladigan
  narsa — hash: javob vaqtidan saqlangan hash'ning baytlarini bilib olish eng yomoni "jadval o'qib olindi" degani, yuqorida
  ko'rdik, bu token bermaydi (hash'dan token tiklanmaydi). OWASP'ning boshqa varianti — "split token": ochiq lookup ID + maxfiy
  verifier; serverda verifier'ning SHA-256'i saqlanadi va `subtle.ConstantTimeCompare` bilan solishtiriladi.
- **Session fixation.** Hujumchi o'z bilgan tokenini qurbonning brauzeriga qo'ydiradi va qurbon login qilgach, o'sha token
  bilan kiradi. Bu yerda server faqat o'zi bergan tokenni qabul qiladi, har login **yangi** token oladi va so'rovda eski
  session cookie'si bo'lsa, uning qatorini o'chiradi (`endCurrentSession`; OWASP: "Renew the Session ID After Any Privilege
  Level Change", eski ID yo'q qilinsin). Test: `TestReloginDropsOldSession`.
- **Muddat.** `sessionTTL = 24 * time.Hour` — absolute timeout (eng uzun umr), bizning tanlov. OWASP misoli: ish kuni
  ilovasiga 4–8 soat; harakatsizlik (idle) muddati kam xavfli ilovaga 15–30 daqiqa, yuqori xavfliga 2–5 daqiqa. Idle timeout
  uchun `expires_at` so'rovlarda uzaytiriladi (har so'rovda yozmaslik uchun, masalan, muddat yarmidan o'tganda), absolute
  chegara esa alohida ustunda (`created_at`) qoladi — keyingi qadam. Muddat ilova vaqti bilan yoziladi (`time.Now()`),
  bazada esa `now()` bilan tekshiriladi: ikki soat orasidagi farq soniyalarda bo'lsa ahamiyatsiz, aks holda bitta soatni
  tanlang (`expires_at > $2` yoki `now() + interval`). `Expires` mijozning soatiga bog'liq; `Max-Age` bog'liq emas.
- **Javob vaqti.** Email yo'q bo'lsa ham `checkPassword` `dummyHash` bilan bir marta argon2id hisoblaydi (OWASP
  Authentication Cheat Sheet, "Authentication Responses"): aks holda yo'q email'ga javob tezroq keladi va vaqtning o'zi
  email bazada yo'qligini aytadi. O'lchov (`chiqishlar.txt` §3b, curl `time_total`, 5 tadan): noto'g'ri parol 24.8–26.2 ms,
  yo'q email 25.1–31.0 ms.
- **Login'ga hujum.** Har urinish 46 MiB xotira (argon2id) oladi: urinishlar soni cheklanmasa, bu ham parol taxmini, ham
  serverni bo'g'ish yo'li. Rate limit — GB-6.07.
- **`Cache-Control: no-store`.** Login javobi va `requireUser` o'tkazgan har javob: session ID yoki shaxsiy ma'lumot bor javob
  oraliq keshlarda qolmasin (OWASP SM CS, "Web Content Caching").
- **`401` va `WWW-Authenticate`.** RFC 9110 §11.6.1: `401` javobida `WWW-Authenticate` bo'lishi shart (MUST). Cookie session
  uchun ro'yxatdan o'tgan auth-scheme yo'q; `Cookie realm="bozor"` — sintaksisi to'g'ri challenge (auth-scheme — token),
  mijozga "cookie bilan kiring" deydi. Go sarlavha nomini `Www-Authenticate` ko'rinishiga keltiradi (`CanonicalHeaderKey`);
  HTTP sarlavha nomlari katta-kichik harfga sezgir emas.
- **Secure va localhost.** `Secure` cookie faqat HTTPS orqali yuboriladi; MDN: "The https: requirements are ignored when the
  Secure attribute is set by localhost" — brauzerlar va curl (bu yerda 8.22) `http://localhost`'ni istisno qiladi. Bu RFC 6265
  emas, mijozlar odati. Real serverda HTTPS — GB-8.06. Cookie nomi: OWASP `__Host-` prefiksini tavsiya qiladi
  (`__Host-session`: faqat Secure, `Path=/`, Domain'siz) — GB-6.07.
- **Logout idempotent.** `POST /logout` `requireUser`siz: session allaqachon eskirgan yoki o'chirilgan bo'lsa ham `204` va
  `Max-Age=0` — brauzerda eskirgan cookie qolmaydi. `Max-Age=0` (Go'da `MaxAge: -1`, `go doc net/http.Cookie`; RFC 6265 §5.2.2)
  faqat mijoz tomoni; session'ni bekor qiladigan narsa — bazadagi `DELETE` (OWASP: logout'da server tomonida bekor qilinsin).
- **`currentUser` panic qiladi.** Route `requireUser`siz ulansa, `currentUser` jim `User{}` (ID 0) qaytarmaydi: bu dasturchi
  xatosi, `recoverPanic` uni `500` qiladi (`TestCurrentUserWithoutMiddleware`).
- **Kalit turi.** `type userKey struct{}` — `go doc context.WithValue`: kalit `string` yoki boshqa ichki tur bo'lmasin.
- **Indekslar.** `token_hash` primary key: indeksi avtomatik. `sessions_user_id_idx` — `DELETE … WHERE user_id = $1` (vazifa)
  jadvalni oxirigacha o'qimasin. `users_email_lower_idx` — `lower(email)` bo'yicha unikal; `UserByEmail` ham
  `lower(email) = lower($1)` bilan qidiradi.
- **Log'da nima bor.** Token ham, hash ham log'ga yozilmaydi; request ID va xato sababi (yo'q email'ning manzili ham —
  OWASP Logging CS autentifikatsiya xatolarini yozishni tavsiya qiladi; ishda shaxsiy ma'lumot siyosatiga qarab).
- **Eski qatorlar.** Muddati o'tgan qatorlar jadvalda qoladi (`UserBySession` ularni qabul qilmaydi). Ularni tozalash —
  TASK.md'dagi 🟡 qism; doimiy fon ishlari — GB-9. Audit kerak bo'lsa, o'chirish o'rniga `revoked_at` ustuni ham ishlatiladi.

## Vazifa

[`TASK.md`](TASK.md): `POST /logout-all` ("hamma qurilmadan chiqish"). Tekshirish:

```bash
docker compose up -d --wait
go test ./...
```

Tugagach `docker compose down -v`. Nima o'zgarganini ko'rish: `diff -ru . solution` (vazifa) va
`diff -u ../gb-6-03-parollar/solution/cmd/api/users.go cmd/api/users.go` (`register` endi Postgres'ga yozadi).

## Manbalar

- OWASP: [Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
  (Session ID Entropy, Server-Side Session Token Storage, Cookies, Renew the Session ID, Session Expiration, Web Content Caching) ·
  [Authentication Cheat Sheet: Authentication Responses](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-responses) ·
  [Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- RFC 6265 [HTTP State Management (cookies)](https://www.rfc-editor.org/rfc/rfc6265) · RFC 9110
  [401](https://www.rfc-editor.org/rfc/rfc9110#name-401-unauthorized), [204](https://www.rfc-editor.org/rfc/rfc9110#name-204-no-content),
  [WWW-Authenticate](https://www.rfc-editor.org/rfc/rfc9110#name-www-authenticate) ·
  MDN [Set-Cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie)
- Go: [crypto/rand.Read](https://pkg.go.dev/crypto/rand#Read) · [encoding/base64.RawURLEncoding](https://pkg.go.dev/encoding/base64#pkg-variables) ·
  [crypto/sha256.Sum256](https://pkg.go.dev/crypto/sha256#Sum256) · [net/http: Cookie, SetCookie, Request.Cookie](https://pkg.go.dev/net/http#Cookie) ·
  [Request.WithContext](https://pkg.go.dev/net/http#Request.WithContext) · [context.WithValue](https://pkg.go.dev/context#WithValue) ·
  [ServeMux.Handle](https://pkg.go.dev/net/http#ServeMux.Handle)
- PostgreSQL 18: [bytea](https://www.postgresql.org/docs/18/datatype-binary.html) · [timestamptz, interval](https://www.postgresql.org/docs/18/datatype-datetime.html) ·
  [now()](https://www.postgresql.org/docs/18/functions-datetime.html#FUNCTIONS-DATETIME-CURRENT) · [encode](https://www.postgresql.org/docs/18/functions-binarystring.html) ·
  [Indexes on Expressions](https://www.postgresql.org/docs/18/indexes-expressional.html) ·
  [ON DELETE CASCADE](https://www.postgresql.org/docs/18/ddl-constraints.html#DDL-CONSTRAINTS-FK) · [DELETE](https://www.postgresql.org/docs/18/sql-delete.html)
- curl: [`-c` / `-b`](https://curl.se/docs/manpage.html)
- Video: kanaldagi "Go backend" playlist'i · yozma dars: haywan.uz/lab
