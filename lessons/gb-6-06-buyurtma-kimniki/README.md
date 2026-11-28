# GB-6.06 · Bu buyurtma kimniki? Authorization va BOLA

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #48"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi darslar: [`gb-6-05-session-jwt`](../gb-6-05-session-jwt) (session yoki JWT; bozor session'da qoladi) ·
[`gb-6-04-login-session`](../gb-6-04-login-session) (login, `requireUser`, `currentUser`) · keyingi dars: `gb-6-07`
(cookie, CSRF, CORS, rate limit). Bu papka **GB-6.04'ning `solution/`idan** boshlanadi (Postgres'dagi foydalanuvchilar,
`POST /users` roʻyxatdan oʻtish, session'lar) va unga **GB-5.10'ning katalog va buyurtmalarini** (`GET/POST/DELETE
/products`, `POST /orders`, `GET /orders/{id}`, transaction'dagi `PlaceOrder`) qoʻshadi; xato javoblari GB-6.02'dagi
`problemFor` orqali. Farqlarni koʻrish:
`diff -ru ../gb-6-04-login-session/solution solution`.

**Bu darsdan keyin:** foydalanuvchi faqat **oʻz** buyurtmasini koʻradi va oʻzgartiradi. Egalik tekshiruvi SQL
soʻrovining oʻzida (`WHERE id = $1 AND user_id = $2`), egasi session'dan; begona buyurtma yoʻq buyurtma bilan **bir xil
`404`** oladi; `/admin/…` route'lari alohida rol tekshiruvi (`requireAdmin`) bilan yopiladi va oddiy xaridorga `403`.

## Ikki savol

| | savol | bozor'da | yoʻq boʻlsa |
|---|---|---|---|
| **authentication** | siz kimsiz? | `requireUser` (GB-6.04): cookie → session → foydalanuvchi context'da | `401` |
| **authorization** | shu ishga ruxsatingiz bormi? | obyekt: store metodidagi `AND user_id = $2` · funksiya: `requireAdmin` | `404` (obyekt) · `403` (funksiya) |

OWASP Authorization Cheat Sheet: authorization — soʻralgan amal shu subyektga ruxsat etilganini tekshirish; authentication
esa uning kimligini tekshirish. Login qilgan foydalanuvchiga hamma narsa ruxsat degani emas.

## Muammo: BOLA (OWASP API1:2023)

Har buyurtmaning egasi bazada yozilgan (`orders.user_id`, quyida), `GET /orders/{id}` esa login bilan himoyalangan —
lekin store faqat `WHERE id = $1` bilan oʻqirdi. Dilnoza (ID 1) login qilib, yoʻldagi `41`ni `42`ga almashtirdi va
**Jasurning** buyurtmasini manzili bilan oldi (`_video/stages/v0-bola`):

```
$ curl -s -b jar.txt localhost:8765/orders/41
{"id":41,"user_id":1,"status":"yangi","note":"Chilonzor 9-kvartal, 12-uy","items":[{"product_id":1,"quantity":2,"unit_price":1200000},{"product_id":2,"quantity":3,"unit_price":400000}],"total":3600000}
$ curl -s -b jar.txt localhost:8765/orders/42
{"id":42,"user_id":2,"status":"yangi","note":"Yunusobod 4-kvartal, 17-uy, ertalab","items":[{"product_id":4,"quantity":1,"unit_price":3500000},{"product_id":3,"quantity":1,"unit_price":2500000}],"total":6000000}
```

Bu — [OWASP API Security Top 10 2023](https://owasp.org/API-Security/editions/2023/en/0xa1-broken-object-level-authorization/)
roʻyxatining birinchi bandi, **Broken Object Level Authorization** (IDOR deb ham ataladi: OWASP Authorization Cheat
Sheet, CWE-639). ID'lar ketma-ket:
hujumchi 43, 44… ni ham sinaydi. OWASP'ning qoidasi: mijozdan kelgan ID bilan bazadagi yozuvga tegadigan **har**
funksiyada foydalanuvchining shu amalga ruxsatini tekshiring, va buni testlar bilan mahkamlang.

## Tuzatish

0. **Baza (egasi yozilgan, tekshirilmagan edi)** — `migrations/00006_orders_owner.sql`: `orders.user_id bigint NOT NULL
   REFERENCES users (id)`, eski `customer` matni oʻchadi (kim ekanini session aytadi), `orders_user_id_idx` indeksi.
   `NOT NULL` default'siz faqat boʻsh jadvalda ishlaydi (laboratoriyada `orders` boʻsh, `seed.sql` keyin toʻldiradi).
   **Ishlab turgan bazada** bu bitta fayl emas, bir necha relizga boʻlingan migratsiya (expand/contract):
   1) `ADD COLUMN user_id bigint` (NULL bilan) va `REFERENCES … NOT VALID`, keyin `VALIDATE CONSTRAINT` (jadval uzoq
   qulflanmaydi); 2) eski qatorlarga egasini yozish — erkin matnli `customer`ni `users.id`ga avtomatik moslab boʻlmaydi,
   bu biznes qarori; 3) `CHECK (user_id IS NOT NULL) NOT VALID` → `VALIDATE`, keyin `SET NOT NULL` (Postgres tekshirilgan
   CHECK'dan foydalanib butun jadvalni qayta skanerlamaydi); 4) indeks `CREATE INDEX CONCURRENTLY` bilan (goose'da
   `-- +goose NO TRANSACTION`); 5) eski ustunni (`customer`) keyingi relizda, uni oʻqiydigan kod qolmaganda oʻchirish.
   Bu fayldagi `Down` ham `customer` matnini qaytara olmaydi (faqat boʻsh ustun).
1. **Store** — `internal/store/orders.go`: `Order(ctx, id, userID)`, soʻrov `WHERE id = $1 AND user_id = $2`. Qator
   yoʻq → `ErrOrderNotFound` — begona ham, yoʻq ham. Buyurtmani oʻqiydigan yoki oʻzgartiradigan har metod egasini oladi:
   imzo egasini berishga majbur qiladi (ikkala parametr `int64`, shuning uchun ularning tartibini testlar tekshiradi).
2. **Handler** — `getOrder`: egasi `currentUser(r).ID`dan, ID esa yoʻldan. `POST /orders` ham buyurtmani session egasi
   nomidan yozadi; body'dagi `"user_id"` — nomaʼlum kalit, `400` (OWASP API3:2023, mass assignment).

```
$ curl -i -b jar.txt localhost:8765/orders/42
HTTP/1.1 404 Not Found
Content-Type: application/problem+json
Cache-Control: no-store
X-Request-Id: QIMZBTSJEJQDPN5YALXRKZR62B
Date: Fri, 27 Nov 2026 17:51:15 GMT
Content-Length: 131

{"type":"about:blank","title":"Not Found","status":404,"detail":"bunday buyurtma yoʻq","request_id":"QIMZBTSJEJQDPN5YALXRKZR62B"}

$ curl -i -b jar.txt localhost:8765/orders/999
HTTP/1.1 404 Not Found
Content-Type: application/problem+json
Cache-Control: no-store
X-Request-Id: TXIAA7JQBFNMNVSEOY5EO4YI5W
Date: Fri, 27 Nov 2026 17:51:15 GMT
Content-Length: 131

{"type":"about:blank","title":"Not Found","status":404,"detail":"bunday buyurtma yoʻq","request_id":"TXIAA7JQBFNMNVSEOY5EO4YI5W"}
```

Server log'ida esa kim nimani soʻragani (GB-6.02: mijozga xavfsiz matn, log'ga toʻliq xato):

```
2026/11/27 22:51:15 QIMZBTSJEJQDPN5YALXRKZR62B xato 404: "buyurtma 42, foydalanuvchi 1: buyurtma topilmadi"
2026/11/27 22:51:15 TXIAA7JQBFNMNVSEOY5EO4YI5W xato 404: "buyurtma 999, foydalanuvchi 1: buyurtma topilmadi"
```

**Nega `403` emas, `404`?** `403` "bu buyurtma bor, lekin sizniki emas" deydi: ID'larni sinab, qaysi buyurtmalar
borligini bilib olish mumkin. (Ketma-ket ID'lar buyurtmalar *sonini* baribir ochadi — yangi buyurtma 43 boʻladi; buni
404 yashirmaydi, quyidagi Chuqurroq'qa qarang.)
[RFC 9110 §15.5.4](https://www.rfc-editor.org/rfc/rfc9110#name-403-forbidden): "An origin server that wishes to "hide"
the current existence of a forbidden target resource MAY instead respond with a status code of 404 (Not Found)."
`/admin/users` uchun esa biz `403`ni tanlaymiz: admin boʻlimi borligi sir emas, foydalanuvchi esa kimligini biladi.

## Test

`cmd/api/api_test.go` (haqiqiy Postgres, har test oʻz sxemasida, `internal/pgtest`):

```go
func TestOrderOwner(t *testing.T) {
	h, _ := testAPI(t)
	dilnoza := login(t, h, dilnozaLogin)
	for id, want := range map[string]int{"41": 200, "42": 404, "999": 404} {
		rec := send(h, "GET", "/orders/"+id, "", dilnoza.Value)
		if rec.Code != want {
			t.Errorf("Dilnoza GET /orders/%s: %d, kutilgan %d", id, rec.Code, want)
		}
	}
}
```

```
$ go test -run TestOrderOwner ./cmd/api          # eski kod (v0-bola)
--- FAIL: TestOrderOwner (0.19s)
    api_test.go:294: Dilnoza GET /orders/42: 200, kutilgan 404
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api	0.278s
FAIL
$ go test -run TestOrderOwner ./cmd/api          # tuzatilgan kod
ok  	github.com/abdurahmon27/bozor/cmd/api	0.283s
```

Yana: `TestForeignOrderLooksMissing` (begona va yoʻq buyurtmaning javobi `request_id`dan boshqa hammasi bilan bir xil,
Jasurning manzili Dilnozaga chiqmaydi, Jasur 41'ni ham koʻrmaydi), `TestCreateOrderOwner` (egasi session'dan, body'dagi
`user_id` → `400`), `TestRegisterThenLogin` (yangi foydalanuvchi `xaridor`, body'dagi `"role"` → `400`), `TestAdminUsers`,
`TestProductsAdminOnly` (katalogni faqat admin oʻzgartiradi), `TestRequireAdminOutside`, `TestUpdateNoteInput`; store'da
`TestOrderOwner`, `TestPlaceOrderOwner`, `TestUsersRole`.

## Rol: `requireAdmin` (OWASP API5:2023)

`migrations/00007_users_role.sql`: `users.role text NOT NULL DEFAULT 'xaridor' CHECK (role IN ('xaridor', 'admin'))` —
yangi foydalanuvchi har doim xaridor (eng kam huquq). `store.User` endi `Role` maydoni bilan, `UserBySession` uni ham
oʻqiydi. `cmd/api/authz.go`:

```go
func requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if currentUser(r).Role != store.RoleAdmin {
			writeError(w, errForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

Route'lar: `GET /admin/users`, `POST /products`, `DELETE /products/{id}` — hammasi
`a.requireUser(requireAdmin(…))`: tashqi qavat kim ekanini aniqlab foydalanuvchini context'ga qoʻyadi, ichki qavat rolini
tekshiradi. `requireJSON` (body'li POST/PATCH) `requireUser`ning ichida: login qilmagan soʻrov avval `401` oladi.
Rol har soʻrovda bazadan oʻqiladi (`UserBySession`): admin huquqi olib qoʻyilsa, keyingi soʻrovdayoq ishlaydi.

**Keyingi qadam (GB-6.08):** har admin route'ni alohida oʻrash unutilishi mumkin. Ishonchliroq shakl — `/admin/` uchun
alohida `ServeMux`, u bir marta `a.requireUser(requireAdmin(admin))` bilan ulanadi: unga qoʻshilgan har route oʻzi
himoyalangan boʻladi (OWASP API5: "consistent and easy-to-analyze authorization module").

**Tez-tez uchraydigan xato** (`_video/stages/v1-admin-yoq`): yoʻlda `/admin` bor, lekin oʻram faqat `requireUser`.
Oddiy xaridor (Dilnoza) hamma foydalanuvchilarni oladi — **Broken Function Level Authorization**
([OWASP API5:2023](https://owasp.org/API-Security/editions/2023/en/0xa5-broken-function-level-authorization/): "Don't
assume that an API endpoint is regular or administrative only based on the URL path"):

```
$ curl -s -b jar.txt localhost:8765/admin/users
[{"id":1,"email":"dilnoza@bozor.test","role":"xaridor"},{"id":2,"email":"jasur@bozor.test","role":"xaridor"},{"id":3,"email":"admin@bozor.test","role":"admin"}]
```

Tuzatilganda Dilnoza `403 Forbidden` (`"detail":"bu amal uchun ruxsat yoʻq"`), admin — roʻyxat. **Tartib teskari**
boʻlsa (`requireAdmin(a.requireUser(…))`, `_video/stages/v2-teskari`), admin ham `500` oladi: `requireAdmin` birinchi
ishlaydi, context'da hali foydalanuvchi yoʻq, `currentUser` panic qiladi (GB-6.04), `recoverPanic` uni `500` qiladi, log'da
`panic: currentUser: route requireUser bilan oʻralmagan`. Xato jim oʻtib ketmaydi.

Chuqurroq: body'dagi `"user_id"` yoki `"role"` — **mass assignment** (OWASP API3:2023 Broken Object Property Level
Authorization'ning bir turi). Asosiy himoya — kirish struct'ida (`orderInput`, `userInput`) bunday maydon yoʻq;
`DisallowUnknownFields` uni jim tashlamay, ochiq `400` qiladi. Tasodifiy ID (UUID) taxminni qiyinlashtiradi va OWASP API1 uni
qoʻshimcha qatlam sifatida tavsiya qiladi, lekin yolgʻiz oʻzi yetarli emas
([OWASP Authorization Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authorization_Cheat_Sheet.html#ensure-lookup-ids-are-not-accessible-even-when-guessed-or-cannot-be-tampered-with):
"generally not sufficient by itself"): tekshiruv har soʻrovda kerak. Egalik — munosabatga asoslangan qoida, rol — RBAC;
OWASP murakkab tizimlar uchun ABAC/ReBAC'ni afzal koʻradi.

## Ishga tushirish

Kerak: Go 1.27, Docker, curl, goose (`go install github.com/pressly/goose/v3/cmd/goose@v3.30.2`, GB-5.06).
Postgres `127.0.0.1:5516`da (shu papkaning `compose.yaml`i, `gb606` loyihasi), API `:8765`da.

```bash
docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5516/bozor
goose up                                         # 00001 … 00007
docker compose exec db psql -f /lab/seed.sql     # 3 foydalanuvchi, 4 mahsulot, 41 va 42-buyurtmalar
export DATABASE_URL=$GOOSE_DBSTRING
go run ./cmd/api                                 # bozor API: http://localhost:8765
```

Laboratoriya parollari (faqat shu yerda; bazada argon2id hash'i, GB-6.03):

| email | parol | rol | buyurtma |
|---|---|---|---|
| `dilnoza@bozor.test` | `tarvuz-osmon-qalam-ariq-27` | xaridor | 41 |
| `jasur@bozor.test` | `choy-non-qand-27` | xaridor | 42 |
| `admin@bozor.test` | `qovun-anjir-tut-uzum-31` | admin | — |

Ikkinchi terminalda: `curl -s -c jar.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'`,
keyin `curl -s -b jar.txt localhost:8765/orders/41`. Toʻxtatish va maʼlumotni oʻchirish: `docker compose down -v`.

## Papka

```
gb-6-06-buyurtma-kimniki/
  go.mod, go.sum              module github.com/abdurahmon27/bozor, go 1.27; pgx v5.13.0, goose v3.30.2, x/crypto
  compose.yaml                Postgres 18, 127.0.0.1:5516, nomli volume (pgdata)
  migrations/                 00001…00005 (GB-5.10, GB-6.04), YANGI: 00006_orders_owner.sql, 00007_users_role.sql
  seed.sql                    3 foydalanuvchi (rollar bilan), 4 mahsulot, buyurtmalar 41 (Dilnoza) va 42 (Jasur)
  migrate.go                  migratsiyalar embed bilan (testlar uchun, GB-5.06)
  cmd/api/authz.go            YANGI: errForbidden, requireAdmin
  cmd/api/handlers.go         route'lar; mahsulotlar (yozish faqat admin), getOrder (egasi bilan), createOrder,
                              updateNote (VAZIFA), listUsers
  cmd/api/auth.go, users.go   GB-6.04/6.03: login, requireUser, currentUser, logout, logoutAll, register (POST /users)
  cmd/api/problem.go          GB-6.02: problemFor (+ 403, buyurtma 404, 409, 422)
  cmd/api/middleware.go       GB-6.01 zanjiri; requireJSON endi PATCH body'sini ham tekshiradi
  cmd/api/*_test.go           api_test.go, tekshiruv_test.go (vazifa tekshiruvi — oʻzgartirmang)
  internal/store/orders.go    YANGI: mahsulotlar va buyurtmalar (List/Get/Add/Delete, Order, PlaceOrder, UpdateNote — VAZIFA)
  internal/store/store.go     foydalanuvchilar (Role bilan, AddUser rolini bazadan oladi), session'lar, Users
  internal/password, pgtest   GB-6.03, GB-5.10 (har test oʻz sxemasida)
  TASK.md                     vazifa
  solution/                   javob (oʻz go.mod'i bilan, alohida ishlaydi)
  _video/                     chiqishlar.sh → chiqishlar.txt (videodagi hamma buyruq va chiqish), vazifa.sh → vazifa.txt,
                              stages/ (videodagi xato versiyalar: v0-bola, v1-admin-yoq, v2-teskari)
```

## Vazifa

[`TASK.md`](TASK.md): `PATCH /orders/{id}` (yetkazish izohini oʻzgartirish) UPDATE'ida xuddi shu xato bor — Jasur
Dilnozaning manzilini oʻzgartira oladi. `go test ./...` oʻtguncha tuzating. Starter'da (`_video/vazifa.txt`):

```
$ curl -s -i -X PATCH -b jasur.txt localhost:8765/orders/41 -H "Content-Type: application/json" -d '{"note":"Jasurning manzili"}' | head -1
HTTP/1.1 204 No Content
$ docker compose exec -T db psql -c "SELECT id, user_id, note FROM orders WHERE id = 41"
 id | user_id |       note        
----+---------+-------------------
 41 |       1 | Jasurning manzili
(1 row)
```

`solution/`da oʻsha soʻrov `HTTP/1.1 404 Not Found`, izoh `Chilonzor 9-kvartal, 12-uy` boʻlib qoladi.

## Manbalar

- OWASP API Security Top 10 2023: [API1 BOLA](https://owasp.org/API-Security/editions/2023/en/0xa1-broken-object-level-authorization/),
  [API3 BOPLA / mass assignment](https://owasp.org/API-Security/editions/2023/en/0xa3-broken-object-property-level-authorization/),
  [API5 BFLA](https://owasp.org/API-Security/editions/2023/en/0xa5-broken-function-level-authorization/)
- [OWASP Authorization Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authorization_Cheat_Sheet.html)
  (least privilege, deny by default, validate on every request, lookup IDs, logging, tests)
- [RFC 9110 §15.5.4 403 Forbidden](https://www.rfc-editor.org/rfc/rfc9110#name-403-forbidden),
  [RFC 5789 PATCH](https://www.rfc-editor.org/rfc/rfc5789)
- [PostgreSQL 18: ALTER TABLE](https://www.postgresql.org/docs/18/sql-altertable.html),
  [CREATE INDEX](https://www.postgresql.org/docs/18/sql-createindex.html)
- Hamma chiqishlar shu kompyuterda: go1.27.4 linux/amd64, `CGO_ENABLED=0`, curl 8.22.0, goose v3.30.2, PostgreSQL 18
  (`postgres:18-alpine`), `_video/chiqishlar.sh`.
