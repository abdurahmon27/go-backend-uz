# GB-6.05 · Session yoki JWT: halol solishtiruv

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #47"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-6-04-login-session`](../gb-6-04-login-session) (login va bekor qilinadigan session'lar) · keyingi
dars: `gb-6-06` (bu buyurtma kimniki? authorization va BOLA). Bu **6-modulning (foydalanuvchilar va xavfsizlik) beshinchi darsi**.

**Bu darsdan keyin:** JWT'ning har baytini tushunasiz — uni kutubxonasiz yasaysiz (`header.payload.imzo`, base64url,
HMAC-SHA256), payload'ni kalitsiz oʻqiysiz, [`golang-jwt/jwt/v5`](https://pkg.go.dev/github.com/golang-jwt/jwt/v5) bilan
toʻrt qoida asosida tekshirasiz va session bilan bir xil soʻrovlarda solishtirasiz: logout'dan keyin eski cookie `401`,
oʻsha JWT esa `200`. Vazifa: `jti` denylist — logout JWT'ni ham bekor qilsin, `go test ./...` oʻtguncha.

**Nega store xotirada.** Session'lar `map`da (mutex bilan), lekin qoida [`gb-6-04`](../gb-6-04-login-session)'dagi
Postgres `sessions` jadvali bilan bir xil: tokenning oʻzi emas, SHA-256 hash'i saqlanadi; yozuvni oʻchirish = session'ni
bekor qilish. Bu darsning mavzusi ikki yondashuvning farqi, bazaning oʻzi emas, va papka Docker'siz `go test` bilan
tekshirilsin. `internal/store`dagi metodlar GB-6.04'dagi Postgres `Store` bilan bir xil shaklda (`ctx` birinchi).
Server qayta ishga tushsa, xotiradagi session'lar ham, JWT kaliti ham (`token.NewKey()`, har ishga tushishda yangi)
yoʻqoladi — eski cookie'lar ham, eski JWT'lar ham oʻladi.

**Solishtirish uchun** `POST /login` ikkalasini beradi: session tokeni cookie'da (GB-6.04), JWT javob body'sida
(`"token"`). `GET /me` ikkalasini qabul qiladi: `Authorization: Bearer <JWT>` boʻlsa — JWT yoʻli, aks holda cookie.
Ishda bittasi tanlanadi (pastdagi "Xulosa").

## Papka

```
gb-6-05-session-jwt/
  go.mod, go.sum               module github.com/abdurahmon27/bozor, go 1.27;
                               github.com/golang-jwt/jwt/v5 v5.3.1, golang.org/x/crypto v0.59.0 (argon2id, GB-6.03)
  internal/token/qolda.go      YANGI: SignHS256 — JWT'ni kutubxonasiz yasash (b64 = base64url, HMAC-SHA256)
  internal/token/token.go      YANGI: NewKey, Issue (kutubxona bilan imzolash), Parse (toʻrt qoida bilan tekshirish)
  internal/token/token_test.go YANGI: qoʻlda yasalgan token oʻtadi; boshqa kalit, sub almashgan, alg none, HS384,
                               exp yoʻq, muddati oʻtgan, iss/aud boshqa, sub son emas — rad
  internal/store/              xotiradagi foydalanuvchilar, session'lar va (vazifa) denylist; Seed — GB-6.04'ning ikki
                               foydalanuvchisi (dilnoza@bozor.test / tarvuz-osmon-qalam-ariq-27, jasur@bozor.test /
                               choy-non-qand-27), oʻsha argon2id satrlari bilan
  internal/password/           GB-6.03/6.04: argon2id (GB-6.04 bilan bir xil)
  cmd/qolda/                   YANGI: qoʻlda yasalgan token, uch qismi alohida qatorda
  cmd/tekshir/                 YANGI: qoʻlda yasalgan token va uning buzilgan nusxalari token.Parse'dan oʻtadi
  cmd/api/auth.go              login (cookie + JWT), authenticate (Bearer → JWT, aks holda cookie), logout, logout-all
  cmd/api/users.go             GB-6.03/6.04: POST /users (roʻyxatdan oʻtish), oʻzgarmagan
  cmd/api/handlers.go          /health, POST /users, POST /login, GET /me, POST /logout, POST /logout-all
  cmd/api/problem.go           GB-6.02: yaroqsiz JWT ham 401 "avval login qiling" (sababi faqat log'da);
                               JWT 401'da WWW-Authenticate: Bearer error="invalid_token" (RFC 6750 §3)
  cmd/api/middleware.go        GB-6.01/6.04: requestID → logRequests → recoverPanic → muxErrors → mux
  cmd/api/api_test.go          ikki yoʻl bilan /me, yaroqsiz JWT → 401 + WWW-Authenticate, session logout, yangi kalit,
                               /users va /logout-all (logout-all JWT'ga taʼsir qilmaydi)
  */tekshiruv_test.go          vazifa tekshiruvi; oʻzgartirmang
  TASK.md                      vazifa
  solution/                    javob (oʻz go.mod'i bilan)
  _video/                      chiqishlar.sh → chiqishlar.txt (videodagi hamma buyruq va chiqish),
                               stages/token-qoidasiz.go (klassik xato: qoidasiz Parse)
```

GB-6.04 bilan farq (Postgres oʻrniga xotira, JWT qoʻshildi): `diff -ru ../gb-6-04-login-session/solution/cmd/api cmd/api`.
`cmd/api`dagi oʻzgarishlar: `api`da `jwtKey`; `login` JWT ham beradi (`loginOutput`); `requireUser` endi `authenticate`ni
chaqiradi; `problemFor`da `errBadToken` → 401. `json.go`, `users.go`, `middleware.go` (faqat `newHandler`ning argumenti) —
GB-6.04'dagidek.

## JWT ichida: qoʻlda

```go
// internal/token/qolda.go
func b64(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func SignHS256(header, payload string, key []byte) string {
	input := b64([]byte(header)) + "." + b64([]byte(payload))
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(input))
	return input + "." + b64(mac.Sum(nil))
}
```

- **JWT** (RFC 7519) — nuqta bilan ulangan uch qism: `header` (`{"alg":"HS256","typ":"JWT"}`), `payload` (claim'lar:
  `sub` — kim haqida, `iss` — kim bergan, `aud` — kim uchun, `exp` — qachon oʻladi; kutubxona yana `iat` — berilgan vaqt
  va `jti` — tokenning noyob ID'sini qoʻshadi, RFC 7519 §4.1) va **imzo**.
- **base64url** (RFC 7515 §2, RFC 4648 §5): alifbo `A–Z a–z 0–9 - _`, oxirida `=` yoʻq — token URL'da ham, sarlavhada ham
  buzilmaydi. Odatda JWT `eyJ` bilan boshlanadi: bu `{"` va harfning base64url'i (header boʻsh joysiz yozilganda).
- **Imzo** = `HMAC-SHA256(kalit, base64url(header) + "." + base64url(payload))` (RFC 7518 §3.2). Kalitsiz uni hisoblab
  boʻlmaydi; kalit kamida 32 bayt boʻlsin (RFC 7518 §3.2: hash chiqishidan qisqa emas) — `token.NewKey()`.
- `exp` — Unix vaqt (1970-01-01 UTC'dan beri soniyalar, RFC 7519 §2 "NumericDate"): `time.Now().Add(15 * time.Minute).Unix()`.
- **Payload shifrlanmagan**: uni kalitsiz istalgan kishi oʻqiydi (`echo <payload> | basenc --base64url -d`). JWT'ga parol
  yoki boshqa sir qoʻyilmaydi. HS256 imzosi faqat "buni kalitni biladigan kimdir yozgan va u oʻzgarmagan"ni kafolatlaydi
(RFC 7515; OWASP: "only proves that the token was produced by some holder of the shared secret").
- Tekshiruvda server shu HMAC'ni qayta hisoblaydi va uchinchi qism bilan solishtiradi — bazaga murojaat qilmasdan.
  (Kutubxona solishtirishni doimiy vaqtda qiladi: `hmac.Equal`, qaysi baytda farq borligi vaqtdan bilinmasin.)

## Tekshirish: kutubxona va toʻrt qoida

```go
// internal/token/token.go
	_, err := jwt.ParseWithClaims(s, rc, func(*jwt.Token) (any, error) {
		return key, nil
	},
		jwt.WithValidMethods([]string{"HS256"}), // alg'ni token emas, server tanlaydi (RFC 8725 §3.1)
		jwt.WithIssuer(Issuer),                  // iss = "bozor"            (RFC 8725 §3.8)
		jwt.WithAudience(Audience),              // aud = "bozor-api"        (RFC 8725 §3.9)
		jwt.WithExpirationRequired(),            // exp boʻlishi shart, oʻtmagan boʻlsin
	)
```

`go run ./cmd/tekshir` (shu kompyuterda, `_video/chiqishlar.txt`):

```
qoʻlda yasalgan:         ok, sub=1
sub 1 → 2, imzo eski:    token signature is invalid: signature is invalid
alg none, imzo yoʻq:     token signature is invalid: signing method none is invalid
exp yoʻq:                token has invalid claims: token is missing required claim: exp claim is required
aud boshqa xizmat:       token has invalid claims: token has invalid audience
muddati oʻtgan:          token has invalid claims: token is expired
alg HS512, imzo toʻgʻri: token signature is invalid: signing method HS512 is invalid
```

**Klassik xato — qoidasiz `Parse`** (`_video/stages/token-qoidasiz.go`, faqat sinov uchun):

```
alg none, imzo yoʻq:     token signature is invalid: token is unverifiable: 'none' signature type is not allowed
exp yoʻq:                ok, sub=1
aud boshqa xizmat:       ok, sub=1
muddati oʻtgan:          token has invalid claims: token is expired
alg HS512, imzo toʻgʻri: ok, sub=1
```

`golang-jwt/jwt/v5` `none`ni oʻzi rad etadi (`jwt.UnsafeAllowNoneSignatureType` berilmasa) va `exp` bor boʻlsa uni
tekshiradi, lekin **`exp`siz tokenni** (abadiy token), **boshqa `aud`ni** va **header tanlagan boshqa algoritmni** qabul
qiladi. Toʻrt option shularni yopadi.

Qoidasiz `Parse`ning xavfi aniq nimada: `exp`siz token (abadiy) va boshqa xizmat uchun berilgan token (`aud`)
qabul qilinadi. `HS512`ning oʻtishi bitta HMAC kaliti bor joyda hali hujum emas, lekin algoritmni token tanlayapti
(RFC 8725 §3.1); ochiq kalitli algoritm bilan bu **alg confusion**ga aylanadi (pastda, "Xulosa").

Nega kutubxona: RFC 8725 (JWT Best Current Practices) tekshiruvchilar yanglishgan joylarni sanaydi — header'dagi `alg`ga
ishonish (§2.1, `none`, alg confusion; §3.1 — ruxsat etilgan algoritmlarni chaqiruvchi beradi), `iss`/`aud`ni
tekshirmaslik (§3.8, §3.9). `exp` bor boʻlsa, undan keyin token qabul qilinmasligi shart (RFC 7519 §4.1.4); `exp`ning
oʻzini **talab qilish** — bizning qoidamiz: `exp`siz token hech qachon oʻlmaydi. `github.com/golang-jwt/jwt` — faol
yuritiladigan JWT kutubxonasi (eski `dgrijalva/jwt-go`ning davomi); versiya aniq
qoʻyilgan: `go get github.com/golang-jwt/jwt/v5@v5.3.1`.

## Solishtiruv: bir xil soʻrovlar (shu kompyuterda, `_video/chiqishlar.txt`)

```bash
go run ./cmd/api                                  # 1-terminal: bozor API: http://localhost:8765
# 2-terminal:
curl -s -c jar.txt localhost:8765/login -H "Content-Type: application/json" \
  -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}' | tee login.json
T=$(jq -r .token login.json)                      # JWT shell oʻzgaruvchisida
echo "$T" | cut -d. -f2 | basenc --base64url -d   # payload: kalitsiz oʻqiladi
curl -s -b jar.txt localhost:8765/me              # {"id":1,"email":"dilnoza@bozor.test"}  (store'dan)
curl -s -H "Authorization: Bearer $T" localhost:8765/me   # {"id":1}  (faqat tokendagi sub)
cp jar.txt eski.txt                               # "oʻgʻirlangan" cookie nusxasi
curl -s -i -X POST -b jar.txt -c jar.txt -H "Authorization: Bearer $T" localhost:8765/logout   # 204
curl -s -b eski.txt localhost:8765/me             # 401: session yozuvi oʻchgan
curl -s -H "Authorization: Bearer $T" localhost:8765/me   # {"id":1} — 200: JWT exp'gacha yaroqli
printf %s "$T" | wc -c                            # 233 (session tokeni: 43)
```

| | session (GB-6.04) | JWT (bu dars) |
|---|---|---|
| serverda nima saqlanadi | token hash'i → foydalanuvchi, muddat | hech narsa (faqat kalit) |
| har soʻrovda | qidiruv: token hash'i → session → foydalanuvchi (GB-6.04: JOIN'li bitta SELECT; bu yerda ikki map) | imzo va claim'larni tekshirish, qidiruvsiz |
| logout / bloklash | yozuvni oʻchirish: **keyingi soʻrovdayoq 401** (`/logout-all` — hamma qurilmadagi session'lar) | **exp'gacha ishlaydi**; tuzatish — denylist yoki qisqa exp + refresh |
| oʻgʻirlansa | logout yoki admin uni darhol oʻldiradi | exp'gacha (bu yerda 15 daqiqa) ishlaydi |
| oʻlcham (shu yerda oʻlchangan) | cookie qiymati 43 belgi, `Cookie: session=…` 59 bayt | token 233 belgi, `Authorization: Bearer …` 255 bayt |
| kim nimani biladi | `/me` foydalanuvchini store'dan oladi (id, email) | faqat tokendagini (`sub` → id); email kerak boʻlsa — bazadan (payload'ni hamma oʻqiydi, shaxsiy maʼlumot tokenga qoʻyilmaydi — OWASP) |
| hammani chiqarish | jadvalni tozalash | kalitni almashtirish (hamma tokenlar birdan oʻladi) |

**Logout har credential'ni alohida bekor qiladi.** `POST /login` ikkalasini beradi, shuning uchun videoda logout'ga cookie
ham, `$T` ham yuboriladi. Faqat cookie bilan logout JWT'ga tegmaydi, faqat Bearer bilan logout (vazifadan keyin) —
session'ga; `/logout-all` ham faqat session'larni oʻchiradi (`api_test.go`: `TestLogoutHarBiriAlohida`,
`TestUsersVaLogoutAll`). JWT'larni "hamma qurilmadan" oʻldirish uchun foydalanuvchi boʻyicha yozuv kerak ("shu vaqtdan
oldin berilgan tokenlar yaroqsiz") — yana server holati. Ishda bitta usul tanlanadi, bu ikkilik yoʻqoladi.

**JWT'ni bekor qilishning narxi** (OWASP JWT Cheat Sheet, "JWT revocation" va "JWT denylist"):
- **denylist** (`jti`): logout `jti`ni `exp`gacha saqlaydi, har JWT soʻrovi uni tekshiradi — yana store'ga qidiruv, session'dagidek (vazifa).
  OWASP kalitni `(jti, iss)` juftligi qilishni tavsiya qiladi; bozor'da token beruvchi bitta (`iss = bozor`), shuning uchun `jti`
  yetarli. OWASP'ning oʻzi ham aytadi: denylist bilan tizim "won't be completely stateless anymore" va oddiy session
  tizimi ("a plain session system") koʻpincha soddaroq;
- **qisqa `exp` + refresh token**: access token bir necha daqiqa yashaydi, yangisi refresh token bilan olinadi; refresh
  token esa serverda saqlanadi va bekor qilinadi — bu yana session.

## Xulosa: bozor uchun tanlov

- **bozor: session.** Bitta API va bitta Postgres: soʻrovlar baribir bazaga boradi, primary key boʻyicha bitta qidiruv
  qimmat emas. Logout, "hamma qurilmadan chiqish" va admin bloklashi darhol ishlashi kerak. Brauzer uchun cookie
  (`HttpOnly`, `Secure`, `SameSite`, GB-6.07) tayyor mexanizm. Bu laboratoriyada faqat 24 soatlik mutlaq muddat bor;
  harakatsizlik (idle) muddati yoʻq — OWASP Session Management: kam xavfli ilovaga 15–30, yuqori xavfliga 2–5 daqiqa.
- **JWT qachon oʻrinli:** xizmatdan xizmatga qisqa umrli token — tekshiruvchi xizmat token beruvchining
  bazasiga ulanmaydi; tashqi login xizmati (OpenID Connect) bergan qisqa tokenlar. Koʻp xizmatli tizimda HS256 oʻrniga
  ochiq kalitli algoritm (masalan EdDSA, RS256): maxfiy kalit faqat beruvchida, tekshiruvchilarda ochiq kalit. Shunda
  ham `WithValidMethods` shart: aks holda **alg confusion** — header'da `HS256` yozib, hammaga maʼlum ochiq kalit HMAC
  kaliti qilinadi (RFC 8725 §2.1, §3.1).

**401 va `WWW-Authenticate`.** Yaroqsiz JWT'ga javob `WWW-Authenticate: Bearer error="invalid_token"` (RFC 6750 §3).
Cookie'siz yoki eskirgan session'ga GB-6.04'dagi `Cookie realm="bozor"` qoladi: RFC 9110 §11.6.1 401'da shu sarlavhani
talab qiladi, lekin cookie session uchun roʻyxatdan oʻtgan sxema yoʻq — bu bozor'ning izohli tanlovi, standart emas.
`bearer()` sxema nomini katta-kichik harfsiz solishtiradi (RFC 9110 §11.1) va token oldida bir necha boʻsh joyga ruxsat
beradi (RFC 6750 §2.1); Bearer bor-u yaroqsiz boʻlsa, cookie'ga oʻtilmaydi — 401.

## Vazifa

[TASK.md](TASK.md): `jti` denylist — `RevokeToken`/`TokenRevoked` (store), `logout` va `authenticate` (cmd/api);
🟡 `PurgeRevoked`. Tekshirish: `go test ./...`. Javob: [`solution/`](solution).

Vazifadan keyin (`_video/chiqishlar.txt`, 12-qism):

```
$ curl -s -o /dev/null -w '%{http_code}\n' -X POST -H "Authorization: Bearer $T" localhost:8765/logout
204
$ curl -s -H "Authorization: Bearer $T" localhost:8765/me
{"type":"about:blank","title":"Unauthorized","status":401,"detail":"avval login qiling","request_id":"…"}
# server logi: xato 401: "JWT yaroqsiz: jti RXSGSXAMWVCJ57YM4TLYXSPB34 bekor qilingan"
```

### Qoʻlda tekshirish

```bash
go run ./cmd/api            # 1-terminal
# 2-terminal: yuqoridagi login, T=…, keyin:
curl -s -X POST -H "Authorization: Bearer $T" localhost:8765/logout
curl -s -H "Authorization: Bearer $T" localhost:8765/me    # vazifadan keyin: 401
```

## Manbalar

- RFC 7519 — JSON Web Token (JWT): https://www.rfc-editor.org/rfc/rfc7519 (§3 tuzilish, §4.1 claim'lar, §6 `alg: none`)
- RFC 7515 — JSON Web Signature (base64url §2, imzo tekshiruvi §5.2): https://www.rfc-editor.org/rfc/rfc7515
- RFC 7518 §3.2 — HMAC with SHA-2 (HS256, kalit uzunligi): https://www.rfc-editor.org/rfc/rfc7518#section-3.2
- RFC 8725 — JSON Web Token Best Current Practices: https://www.rfc-editor.org/rfc/rfc8725
- RFC 6750 — Bearer Token Usage (§2.1 `Authorization: Bearer`, §3 `WWW-Authenticate`): https://www.rfc-editor.org/rfc/rfc6750
- OWASP JSON Web Token Cheat Sheet: https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_Cheat_Sheet.html
- OWASP Session Management Cheat Sheet: https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
- golang-jwt/jwt/v5: https://pkg.go.dev/github.com/golang-jwt/jwt/v5 · Go `crypto/hmac`: https://pkg.go.dev/crypto/hmac ·
  `encoding/base64`: https://pkg.go.dev/encoding/base64 · `strings.CutPrefix`: https://pkg.go.dev/strings#CutPrefix

Hammasi shu kompyuterda tekshirilgan: go1.27.4 linux/amd64, `CGO_ENABLED=0`, curl 8.22.0, coreutils 9.11 (`basenc`),
jq 1.8.2; `bash _video/chiqishlar.sh > _video/chiqishlar.txt` (server :8765'da, oxirida toʻxtatiladi).
