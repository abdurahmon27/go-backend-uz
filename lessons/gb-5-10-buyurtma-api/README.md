# GB-5.10 · Nazorat nuqtasi: buyurtmalar API, haqiqiy Postgres'da testlar

**Go backend: noldan guru'gacha** kursining 42-darsi, 5-modulning nazorat nuqtasi. GB-4.08'dagi katalog API
ma'lumotni xotirada saqlardi: server to'xtasa, hamma buyurtma yo'qolardi. Bu papkada o'sha API (`cmd/api`, `newMux`,
`writeStoreError`, jadvalli testlar) endi Postgres ustida ishlaydi: goose migratsiyalari (GB-5.06), `pgxpool` (GB-5.05),
`POST /orders` buyurtmani bitta transaction'da yozadi (GB-5.07) va qoldiq qoidasi shartli `UPDATE`da (GB-5.08).
Testlar mock'siz, haqiqiy Postgres'da, har test o'z sxemasida, va ikki qoidani isbotlaydi: buyurtma yo to'liq
yoziladi, yo umuman yozilmaydi; bir vaqtdagi buyurtmalarda ham qoldiq minusga tushmaydi.

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go bo'lsa,
> `go run`/`go test` birinchi marta kerakli versiyani o'zi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qo'ygan bo'lsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-5-09` (sqlc) · keyingi dars: `gb-6-01` (middleware) · API'ning xotiradagi versiyasi: [`../gb-4-08-katalog-api`](../gb-4-08-katalog-api).
Bu papka mustaqil laboratoriya: GB-5.05..5.08 dagi bir xil jadvallar va namunaviy ma'lumot (`seed.sql`). So'rovlar qo'lda
yozilgan pgx so'rovlari; sqlc (GB-5.09) ixtiyoriy.

## Spetsifikatsiya

| route | javob |
|---|---|
| `GET /health` | `200 ok` |
| `GET /products` | `200`, hamma mahsulot ID tartibida (`[]`, bo'sh bo'lsa) |
| `GET /products/{id}` | `200` · `400 bad_id` · `404 not_found` |
| `POST /products` `{"name","price","stock"}` | `201` + `Location` · `422 validation_failed` · `409 conflict` (nom band, `UNIQUE`) |
| `DELETE /products/{id}` | `204` · `404` · `409 conflict` (mahsulot buyurtmalarda bor: foreign key) |
| `POST /orders` `{"customer","items":[{"product_id","quantity"}]}` | `201` + `Location: /orders/{id}` va buyurtma · `404` (mahsulot yo'q) · `409 insufficient_stock` · `422` (bo'sh savat, dona 1..100 emas, mijoz nomi bo'sh) · `400 bad_json` |
| `GET /orders/{id}` | **vazifa** ([TASK.md](TASK.md)): `200` · `400` · `404` |

Pul tiyinda (`int64`, GB-1.03): `1200000` = 12 000 so'm. Buyurtma qatorida `unit_price` — sotuv paytidagi narx
(keyin mahsulot narxi o'zgarsa ham buyurtma o'zgarmaydi).

## Papka

| fayl | nima |
|---|---|
| `compose.yaml` | Postgres 18 (`postgres:18-alpine`), loyiha nomi `gb510`, `127.0.0.1:5510`, nomli volume; Docker Compose 2.20+ (`start_interval`) |
| `migrations/` | GB-5.06'dagi uchta goose migratsiyasi: `products`, `orders` + `order_items`, `orders.status` |
| `seed.sql` | jadvallarni `TRUNCATE … RESTART IDENTITY` bilan bo'shatib, to'rt mahsulot va uch buyurtma yozadi |
| `migrate.go` | `package bozor`: migratsiyalar va `seed.sql` `go:embed` bilan; `NewProvider` (goose kutubxona sifatida, testlar uchun) |
| `cmd/api/` | server: `main.go` (`run`, `pgxpool`, `Ping`), `handlers.go` (route'lar, `writeStoreError`), `json.go` (GB-4.06) |
| `internal/store/` | `Store`: `List`, `Get`, `Add`, `Delete`, `PlaceOrder` (+ `whyNoRow`), `Order` (VAZIFA) |
| `internal/pgtest/` | `pgtest.New(t)`: har test uchun yangi sxema + migratsiya + seed, test oxirida `DROP SCHEMA … CASCADE` |
| `*_test.go` | `store_test.go`, `api_test.go` — darsdagi testlar; `tekshiruv_test.go` (ikkala paketda) — vazifa tekshiruvi |
| `solution/` | to'liq javob (o'z `go.mod`i bilan, alohida ishga tushadi) |
| `_video/` | videodagi buyruqlar: `chiqishlar.sh` → `chiqishlar.txt`; `stages/pgtest-shared.go` — videodagi XATO versiya |

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, pgx v5.13.0, goose v3.30.2,
PostgreSQL 18). To'liq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt). Video `bozor` nomli papkada yozilgan
(compose loyiha nomi = papka nomi, konteyner `bozor-db-1`); bu papkadagi `compose.yaml` `name: gb510` beradi
(`gb510-db-1`), buyruqlar bir xil. goose har qator boshiga sana va vaqtni yozadi; quyida (va videoda) u olib tashlangan.
Skript `curl`ni chiqishi faylga ketgani uchun `-s` bilan ishlatadi; terminalda `curl -i` aynan shu qatorlarni chiqaradi.

```bash
docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5510/bozor
goose up
# OK   00001_products.sql (5.1ms)
# OK   00002_orders.sql (9.27ms)
# OK   00003_orders_status.sql (2.33ms)
# goose: successfully migrated database to version: 3
docker compose exec db psql -f /lab/seed.sql
# TRUNCATE TABLE · INSERT 0 4 · INSERT 0 3 · INSERT 0 4
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5510/bozor
go run ./cmd/api
# bozor API: http://localhost:8765
```

Ikkinchi terminalda:

```bash
curl -i localhost:8765/orders -d '{"customer":"Aziza","items":[{"product_id":1,"quantity":2},{"product_id":3,"quantity":1}]}'
# HTTP/1.1 201 Created
# Location: /orders/4
# {"id":4,"customer":"Aziza","status":"yangi","items":[{"product_id":1,"quantity":2,"unit_price":1200000},{"product_id":3,"quantity":1,"unit_price":2500000}],"total":4900000}

curl -i localhost:8765/orders -d '{"customer":"Jasur","items":[{"product_id":1,"quantity":2},{"product_id":4,"quantity":5}]}'
# HTTP/1.1 409 Conflict
# {"error":{"code":"insufficient_stock","message":"oʻrik: qoldiq 2, soʻralgan 5"}}

curl -s localhost:8765/products/1
# {"id":1,"name":"olma","price":1200000,"stock":38}     ← 36 emas: Jasurning transaction'i butunlay bekor bo'ldi

docker compose exec db psql -c "SELECT id, customer, status FROM orders WHERE id > 3"
#  id | customer | status
#   4 | Aziza    | yangi

curl -s localhost:8765/orders/4
# 404 page not found                                   ← vazifa: bu route hali yo'q
```

Testlar (`go test ./...`: ikkala paket `ok`). **Klassik xato** (`_video/stages/pgtest-shared.go`): hamma test bitta
`public` sxemada, har test uni o'chirib qayta quradi. `go test ./...` har paketning test dasturini alohida quradi va
ularni **parallel** ishga tushiradi (`-p`, standarti `GOMAXPROCS`, odatda mantiqiy CPU soni: bu mashinada 8 yadro, 12 oqim → 12), shuning uchun bir paketning testi
ikkinchisining jadvallarini o'chirib yuboradi. `go test -count=5 ./...` bilan deyarli har safar yiqiladi, har safar boshqa
testlar va boshqa xabarlar bilan. `go test -count=5 ./internal/store` yolg'iz esa `ok` (paket ichida testlar ketma-ket):

```
$ go test -count=5 ./internal/store
ok  	github.com/abdurahmon27/bozor/internal/store	2.134s
$ go test -count=5 ./...
    --- FAIL: TestRoutes/bitta (0.04s)
        api_test.go:55: migratsiya qoʻllanmadi: failed to list migrations: ERROR: relation "goose_db_version" does not exist (SQLSTATE 42P01)
…
    store_test.go:161: mahsulot 1 qoldigʻi oʻqilmadi: ERROR: relation "products" does not exist (SQLSTATE 42P01)
```

Yechim — bu papkadagi `pgtest.New`: har test o'z sxemasida (`CREATE SCHEMA test_<vaqt>`, `search_path`), oxirida o'chiriladi.
`go test -p 1 ./...` ham xatoni yashiradi, lekin testlarni sekinlashtiradi va izolyatsiya muammosini hal qilmaydi.

`go test -race ./...` (gcc kerak; bu mashinada gcc yo'q, Docker'dagi `golang:1.26.8` + go1.27.4 toolchain bilan
tekshirildi, GB-4.08 `_video/race.sh` kabi): ikkala paket `ok`, race yo'q.

## Chuqurroq (videoga kirmagan)

- **Nega `CHECK (stock >= 0)` ham turibdi?** Qoldiq qoidasini shartli `UPDATE` tekshiradi (`… AND stock >= $1`), `CHECK`
  esa ilova xatosiga qarshi oxirgi to'siq (GB-5.08 Chuqurroq). `PlaceOrder` `CHECK`ga yetmaydi: shart bajarilmasa qator
  o'zgarmaydi, transaction tirik qoladi va `whyNoRow` shu transaction ichida sababni o'qiy oladi. `CHECK` xatosidan
  keyin esa transaction buziladi (GB-5.07: `current transaction is aborted`), unda boshqa so'rov yuborib bo'lmaydi.
- **ID tartibida qulflash.** `PlaceOrder` savat qatorlarini `ProductID` bo'yicha saralaydi (`slices.SortFunc`, nusxada):
  `UPDATE` har qatorni transaction oxirigacha band qiladi; ikki savat [olma, choy] va [choy, olma] tartibida kelsa,
  bir-birini kutib qolardi (deadlock: Postgres bittasini `40P01` bilan to'xtatadi → 500). Bir xil tartib bu holatni yo'q qiladi.
- **Bekor bo'lgan buyurtmaning ID'si.** Jasurning `INSERT INTO orders` buyrug'i 5-ID'ni olgan edi; `ROLLBACK`dan keyin ham
  identity ketma-ketligi orqaga qaytmaydi, keyingi buyurtma 6 bo'ladi (GB-5.07 `ok-run`). ID'lar bo'shliqsiz bo'lishi kafolatlanmaydi.
- **`r.Context()`.** Har handler store'ga so'rovning context'ini beradi: mijoz ulanishni uzsa (yoki HTTP/2'da so'rovni
  bekor qilsa), context bekor bo'ladi ([net/http: Request.Context](https://pkg.go.dev/net/http#Request.Context)). pgx v5
  standart holatda kutishni to'xtatadi va ulanishni yopadi (`DeadlineContextWatcherHandler`); Postgres uzilgan ulanishning
  ochiq transaction'ini bekor qiladi. Postgres'ga alohida bekor qilish so'rovini yuborish uchun
  `pgconn.CancelRequestContextWatcherHandler` bor. Har so'rovga muddat qo'yish (`http.TimeoutHandler`
  yoki middleware) GB-6.01'dan keyin.
- **`main` → `run`.** `os.Exit` `defer`larni bajarmaydi ([os.Exit](https://pkg.go.dev/os#Exit)); shuning uchun ishni `run() error`
  qiladi (`defer pool.Close()` ishlaydi), `main` esa faqat xatoni chop etib chiqadi.
- **goose sessiya qulfi.** GB-5.06'dagi `NewProvider` Postgres advisory lock'ini ishlatardi. U butun bazaga bitta qulf, bu
  yerda esa testlar har biri o'z sxemasida parallel migratsiya qiladi va qulfni kutib sekinlashardi; shuning uchun bu
  laboratoriyada u yo'q. Bir nechta server nusxasi bir vaqtda migratsiya qiladigan joyda (GB-8.06) qulf qaytadi.
- **`(cached)`.** `go test` muvaffaqiyatli natijani keshlaydi: kod va `DATABASE_URL` o'zgarmasa, ikkinchi marta testlar
  ishlamaydi, `ok … (cached)` chiqadi, Postgres o'chiq bo'lsa ham. Bazaning holati kesh kalitiga kirmaydi; integratsion testni
  albatta ishga tushirish uchun `go test -count=1 ./...` (CI'da shunday). Videodagi `-count=5` ham keshni chetlab o'tadi.
- **Noma'lum mahsulot: 404 yoki 422?** `POST /orders` savatida yo'q `product_id` bo'lsa API 404 qaytaradi: GB-4.08'dagi
  qoida (store `ErrNotFound` → 404, bitta joyda) saqlangan. RFC 9110 bo'yicha 404 so'ralgan resurs (`/orders`) yo'qligini
  bildiradi, shuning uchun ko'p API'lar body'dagi noto'g'ri havolaga 422 beradi (`items[i].product_id`). Ikkala tanlov ham
  uchraydi; muhimi, spetsifikatsiyada yozilgan va test bilan qotirilgan bo'lsin.
- **Savatda bir mahsulot ikki marta** — `validateOrder` rad etadi (422 `items[i].product_id`): aks holda ikkinchi qatorning
  `StockError`i transaction ichida allaqachon kamaygan qoldiqni ko'rsatardi. Bo'sh savatni store'ning o'zi ham rad etadi.
- **`active` ustuni** (GB-5.03) bu laboratoriyada ishlatilmaydi: seed'dagi hamma mahsulot faol. Sotuvdan olingan mahsulotni
  sotmaslik (`… AND active` shartli `UPDATE`da) — keyingi qadam.
- **Nomlar.** GB-5.07'dagi `ErrNotEnoughStock`/`ErrInvalidQuantity` o'rniga GB-4.08'dagi `*StockError`/`ErrBadQuantity`:
  API xabari (`oʻrik: qoldiq 2, soʻralgan 5`) `StockError`dagi nom va qoldiqqa tayanadi. Konstruktor `store.New(pool)`.
  GB-4.08'dagi `cmd/bozor` CLI va JSON fayl ombori (`file.go`) bu laboratoriyaga olinmagan: endi ma'lumot Postgres'da.
- **goose CLI.** Videodagi `goose up` uchun: `go install github.com/pressly/goose/v3/cmd/goose@v3.30.2` (GB-5.06).
- **Testcontainers.** Ko'p jamoalar testlar uchun Postgres'ni testning o'zidan Docker orqali ko'taradi
  (testcontainers-go). Bu kursda Postgres'ni `docker compose` ko'taradi (GB-5.02), testlar esa unga ulanadi: qadamlar
  kamroq, CI'da ham shunday ishlaydi (GB-8.05).

## Vazifa

[`TASK.md`](TASK.md): `GET /orders/{id}` (`store.Order`, route, `404`). Tekshirish:

```bash
docker compose up -d --wait
go test ./...
```

Tugagach `docker compose down -v`. Vazifadan oldingi holatni ko'rish (vazifa tekshiruvisiz): `tekshiruv_test.go` fayllarini
vaqtincha olib qo'ying. Nima o'zgarganini ko'rish: `diff -ru ../gb-4-08-katalog-api/solution .` — API'ning o'zi (xotiradagi
store → Postgres). GB-5.05..5.09 papkalari bitta paketli laboratoriyalar, ularning `PlaceOrder`i shu yerdagi
`internal/store`ga ko'chgan: `diff ../gb-5-07-transaction/solution/store.go internal/store/store.go`.

## Manbalar

- PostgreSQL 18: [UPDATE … RETURNING](https://www.postgresql.org/docs/18/sql-update.html) ·
  [TRUNCATE](https://www.postgresql.org/docs/18/sql-truncate.html) · [Schemas, search_path](https://www.postgresql.org/docs/18/ddl-schemas.html) ·
  [DROP SCHEMA](https://www.postgresql.org/docs/18/sql-dropschema.html) · [Foreign keys](https://www.postgresql.org/docs/18/ddl-constraints.html#DDL-CONSTRAINTS-FK) ·
  [Error codes](https://www.postgresql.org/docs/18/errcodes-appendix.html) · [Deadlocks](https://www.postgresql.org/docs/18/explicit-locking.html#LOCKING-DEADLOCKS)
- pgx v5: [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool) · [Tx](https://pkg.go.dev/github.com/jackc/pgx/v5#Tx) ·
  [stdlib.OpenDBFromPool](https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib#OpenDBFromPool)
- goose: [Provider](https://pkg.go.dev/github.com/pressly/goose/v3#Provider) · [CLI](https://pressly.github.io/goose/documentation/cli-commands/)
- Go: [`go test` flags (`-count`)](https://pkg.go.dev/cmd/go#hdr-Testing_flags) · [`-p` (build flags)](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies) ·
  [testing: T.Cleanup](https://pkg.go.dev/testing#T.Cleanup) · [net/http: Request.Context](https://pkg.go.dev/net/http#Request.Context)
- RFC 9110: [201](https://www.rfc-editor.org/rfc/rfc9110#name-201-created) · [Location](https://www.rfc-editor.org/rfc/rfc9110#name-location) ·
  [409](https://www.rfc-editor.org/rfc/rfc9110#name-409-conflict) · [422](https://www.rfc-editor.org/rfc/rfc9110#name-422-unprocessable-content)
- Video: kanaldagi "Go backend" playlist'i · yozma dars: haywan.uz/lab
