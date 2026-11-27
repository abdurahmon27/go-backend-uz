# GB-5.09 · sqlc: SQL'ni siz yozasiz, Go kodi generatsiya qilinadi

**Go backend: noldan guru'gacha** kursining 41-darsi. GB-5.05 va GB-5.07'da har soʻrovni qoʻlda yozdik: SQL matni Go string'ida va
`Scan`, ustunlar bir xil tartibda. Ustun nomida xato boʻlsa, `go build` ham, `go vet` ham jim oʻtadi, xato faqat soʻrov ishlaganda chiqadi.
Bu darsda oʻsha soʻrovlar `query.sql` fayliga koʻchadi (`-- name: … :one/:many/:exec` izohlari bilan), **sqlc** ularni GB-5.06'dagi
goose migratsiyalaridan olingan sxemaga solishtiradi va `db/` papkasiga oddiy Go kodini yozadi. `Store` endi shu kodni chaqiradi.

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-5-08-ikki-xaridor`](../gb-5-08-ikki-xaridor) (lost update, shartli `UPDATE`, `FOR UPDATE`) · keyingi dars: `gb-5-10`
(nazorat nuqtasi: buyurtmalar API, haqiqiy Postgres'da testlar).
Video: kanaldagi «Go backend» playlist'ida (#41). Yozma dars: haywan.uz/lab (gb: Go backend).

Bu papka mustaqil laboratoriya (5-modul darslari kabi): oldingi darslarning kodi kerak emas. Sxema — GB-5.06'ning uchta migratsiyasi
(`migrations/`, baytma-bayt oʻsha; GB-5.06 vazifasidagi `00004_products_category.sql` bu yerda yoʻq, GB-5.07'dagidek), `seed.sql` — oʻsha toʻrt
mahsulot, uch buyurtma va toʻrt qator. `store.go` — **GB-5.07 yechimining** `Store`i (`Get`, `ListActive`, `ByName`, `Cheapest`, `OrderTotal`,
`PlaceOrder`, `Restock`): video `Get`, `ListActive` va `PlaceOrder`ni sqlc'ga oʻtkazadi, qolganlari vazifa. GB-5.08'ning `Buy`/`BuyLocked`/`UseCoupon`i
(qoldiq qoidasi, alohida laboratoriya) bu yerga olinmagan; ularni GB-5.10 (nazorat nuqtasi) birlashtiradi.
Nima oʻzgarganini koʻrish: `diff -u ../gb-5-07-transaction/solution/store.go store.go`.

## Papka

| fayl | nima |
|---|---|
| `compose.yaml` | Postgres 18 (`postgres:18-alpine`), loyiha nomi `gb509`, `127.0.0.1:5409`, nomli volume |
| `migrations/` | GB-5.06'ning migratsiyalari: `products`, `orders` + `order_items`, `orders.status`. sqlc sxemani shu yerdan oʻqiydi (faqat `-- +goose Up` qismini) |
| `seed.sql` | sinov maʼlumoti (GB-5.06'dagidek) |
| `sqlc.yaml` | sqlc sozlamasi: `engine: postgresql`, `schema: migrations`, `queries: query.sql`, Go kodi `db/` papkasiga, `sql_package: pgx/v5` |
| `query.sql` | soʻrovlar: `GetProduct`, `ListActiveProducts`, `CreateOrder`, `DecrementStock` (`@quantity`, `@id`), `AddOrderItem` |
| `db/` | **generatsiya qilingan** (`sqlc generate`): `db.go` (`DBTX`, `Queries`, `New`, `WithTx`), `models.go` (`Product`, `Order`, `OrderItem`), `query.sql.go`. Qoʻlda tahrirlanmaydi, git'ga commit qilinadi |
| `store.go` | videoning oxiridagi holat: `Store{pool, q}`, `NewStore`, `Get`, `ListActive`, `PlaceOrder` sqlc orqali; VAZIFA: qoʻlda yozilgan `ByName`, `Cheapest`, `OrderTotal`; 🟡 qoʻlda yozilgan `Restock` (GB-5.07) |
| `main.go` | videodagi dastur: `Get(1)` va Azizaning buyurtmasi |
| `store_test.go` | tekshiruv: `query.sql` va `db/query.sql.go` (nomlar, imzolar), `*Store` metodlari `go/parser` bilan (Scan va `pool.Query` qolmaganmi, `s.q.…` chaqiriladimi), keyin haqiqiy Postgres'da har test alohida vaqtinchalik sxemada |
| `solution/` | toʻliq javob (oʻz `go.mod`i va `db/`i bilan, alohida ishga tushadi) |
| `_video/` | videodagi buyruqlar: `chiqishlar.sh` → `chiqishlar.txt`, `stages/` (video bosqichlari) |

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, sqlc v1.31.1, goose v3.30.2, pgx v5.13.0, PostgreSQL 18).
Toʻliq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt). Video `bozor` nomli papkada yozilgan (compose loyiha nomi = papka nomi,
konteyner `bozor-db-1`); bu papkadagi `compose.yaml` `name: gb509` beradi (`gb509-db-1`), buyruqlar bir xil.

```bash
docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations                  # GB-5.06'dagidek
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5409/bozor
goose up                                    # OK 00001…00003, successfully migrated database to version: 3
docker compose exec db psql -f /lab/seed.sql                                 # INSERT 0 4 / INSERT 0 3 / INSERT 0 4
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5409/bozor
```

**1. Muammo: qoʻlda yozilgan SQL'dagi xato.** `_video/stages/01-store-qolda.go` — GB-5.05'dagi `Get`, `stock` oʻrniga `qoldiq`:

```
$ go build                 # chiqish yoʻq, exit 0
$ go vet .                 # chiqish yoʻq, exit 0
$ go run .
{0  0 0} ERROR: column "qoldiq" does not exist (SQLSTATE 42703)
```

**2. sqlc.**

```
$ go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
$ sqlc version
v1.31.1
$ sqlc generate            # chiqish yoʻq: xato yoʻq
$ ls db
db.go
models.go
query.sql.go
```

sqlc Postgres'ga ulanmaydi: SQL'ni Postgres'ning oʻz parser'i (libpg_query, `pg_query`) bilan tahlil qiladi. `CGO_ENABLED=0` bilan qurilgan
sqlc bu C kutubxonasini WebAssembly koʻrinishida ishlatadi (`internal/engine/postgresql/parse_wasi.go`, build tag `windows || !cgo`),
shuning uchun C kompilyatori kerak emas. `CGO_ENABLED=1` va gcc boʻlsa, cgo versiyasi quriladi; natija bir xil.

**3. Generatsiya qilingan kod** (toʻliq matn `db/` va `_video/chiqishlar.txt`da): `models.go`'da `Product{ID int64, Name string, Price int64, Stock int32, Active bool}`
(`integer` 4 bayt → `int32`), `Order.Note pgtype.Text` (`NULL` boʻlishi mumkin: `String` + `Valid`). `query.sql.go`'da `GetProductRow` (4 ustun:
hamma ustun tanlanmagani uchun alohida tur; `SELECT id, name, price, stock, active` boʻlsa sqlc `Product`ni qaytaradi) va
`GetProduct`: `QueryRow` + `Scan`, GB-5.05'dagi qoʻlda yozilgan kodning oʻzi.

**4. Store sqlc orqali.** `go run .` → `{1 olma 1200000 40} <nil>`.

**5. Parametr nomlari.** `UPDATE products SET stock = stock - $1 WHERE id = $2` uchun sqlc `DecrementStockParams{Stock int32; ID int64}` yasaydi:
`$1` nomini u turgan joydagi ustundan oladi, aslida esa bu miqdor. `@quantity` va `@id` bilan: `DecrementStockParams{Quantity int32; ID int64}`
(sqlc Postgres'ga baribir `$1`, `$2` yuboradi; `@name` — `sqlc.arg(name)`ning qisqa yozuvi).

**6. PlaceOrder.** `qtx := s.q.WithTx(tx)` — shu transaction ustidagi `Queries`. `go run .` → `{1 olma 1200000 40} <nil>` va `4 <nil>`;
`order_items`da toʻrtinchi buyurtmaning ikki qatori (`1 | 2 | 1200000`, `2 | 1 | 400000`).

**7. Xato sqlc'da.** `query.sql`da `stock` → `qoldiq`:

```
$ sqlc generate
# package db
query.sql:2:25: column "qoldiq" does not exist
```

exit 1, `db/*.go` oʻzgarmaydi (skript sha256 bilan tekshiradi).

Shundan keyin `go build` **oʻtadi** (exit 0): `db/` oʻzgarmagani uchun unda eski, toʻgʻri kod turibdi. Generate xatosini eʼtiborsiz qoldirsangiz,
dastur eski soʻrov bilan quriladi; shuning uchun CI'da `sqlc diff` (pastda).

**Oraliq savol: `sqlc generate` unutilsa.** `query.sql`ga `ProductByName` yozilib, `db/` yangilanmasa, Go kompilyatori toʻxtatadi:
`./store.go:132:16: s.q.ProductByName undefined (type *db.Queries has no field or method ProductByName)`.

Videodagi chiqishlar bosqichma-bosqich `_video/stages/` fayllari bilan olingan: `01-*` qoʻlda yozilgan `Get` (xato bilan), `02-*` faqat `Get`/`ListActive`
sqlc orqali, `03-query-dollar.sql` `$1`/`$2` versiyasi, `05-store-undefined.go` oraliq savol uchun.

## Chuqurroq (videoga kirmagan qismi)

- **`db/` commit qilinadi.** Generatsiya qilingan kod — oddiy Go: dasturni qurish uchun sqlc kerak emas. Har faylning boshida
  `// Code generated by sqlc. DO NOT EDIT.` — Go'ning umumiy konvensiyasi (https://go.dev/s/generatedcode; `go/ast.IsGenerated`, Go 1.21+): linter'lar va code review asboblari bunday faylni generatsiya qilingan deb taniydi.
- **`sqlc diff`** `db/` fayllari `query.sql` va sxemaga mosligini tekshiradi: farq boʻlsa uni koʻrsatadi va **exit 1** bilan chiqadi
  (`chiqishlar.txt`: `diff-toza` exit 0, `diff-farq` exit 1). GB-8.05'dagi CI'da shu buyruq «generate unutildi» xatosini PR bosqichida ushlaydi.
  `sqlc compile` fayl yozmasdan faqat tekshiradi, `sqlc vet` soʻrovlarga lint qoidalarini qoʻllaydi.
- **Takrorlanadigan oʻtkazish kodi.** Har soʻrov uchun alohida `…Row` turi (`GetProductRow`, `ListActiveProductsRow`, …) va har metodda
  `Product{…, Stock: int(r.Stock)}`. `SELECT` jadvalning hamma ustunini olsa (`id, name, price, stock, active`) sqlc `db.Product`ni qaytaradi
  (`chiqishlar.txt`, `Full`): shunda bitta `productFromDB(db.Product) Product` yetadi; `sqlc.embed(products)` ham shu maqsadda.
  Bu darsda soʻrovlar GB-5.05'dagidek qoldi, toʻgʻridan-toʻgʻri solishtirish uchun.
- **Transaction yordamchisi.** `PlaceOrder` va `Restock`da `Begin` / `defer Rollback` / `WithTx` / `Commit` takrorlanadi. Ishda koʻpincha
  `func (s *Store) execTx(ctx, fn func(q *db.Queries) error) error` yoziladi; bozor'da u GB-5.10'da paydo boʻladi.
- **Migratsiya fayllari tartibi.** sqlc `schema:` papkasidagi fayllarni nom boʻyicha (leksikografik) tartibda oʻqiydi; goose'ning
  `00001_…` nomlari shunga mos (https://docs.sqlc.dev/en/v1.31.1/howto/ddl.html).
- **`db` paket nomi** qisqa, lekin funksiya ichida `db` degan oʻzgaruvchi ochsangiz paketni yopib qoʻyadi; generatsiya qilingan `New(db DBTX)`
  ham shunday nom ishlatadi (u faqat funksiya ichida). Kattaroq loyihada `out: internal/db` odatiy joy.
- **Hisoblangan ifodalar — sqlc'ning taxmini.** Ustunlarning turi sxemadan aniq, lekin ifoda va agregatlarning turini sqlc oʻzi taxmin
  qiladi (Postgres'dan soʻramaydi). Oʻlchangan (`chiqishlar.txt`, `ifoda-turlari`): `COALESCE(SUM(quantity * unit_price), 0)` → `interface{}`,
  `SUM(quantity * unit_price)` → `int64` (aslida Postgres `numeric` va qatorsiz holatda `NULL` qaytaradi: test `cannot scan NULL into *int64`
  bilan yiqiladi, `sum-null`), `…::bigint` → `int64`. Qoida: hisoblangan ustunni `::tur` bilan yozing va generatsiya qilingan imzoni oʻqing.
  sqlc'ning `database:` + `analyzer` sozlamasi turlarni haqiqiy Postgres'dan soʻraydi (https://docs.sqlc.dev/en/v1.31.1/reference/config.html).
- **Xatoni Postgres emas, sqlc topadi.** `column "qoldiq" does not exist` — sqlc'ning oʻz xabari (sxemani `migrations/`dan oʻqib tuzgan katalogi
  boʻyicha); Postgres parser'i (libpg_query) faqat SQL matnini qismlarga ajratadi. Baza ishtirok etmaydi, shuning uchun SQLSTATE ham yoʻq.
- **`emit_empty_slices: true`** (`sqlc.yaml`, `gen.go` ostida): `:many` qator boʻlmasa `nil` emas, `[]T{}` qaytaradi (`chiqishlar.txt`,
  `emit_empty_slices`: `items := []int64{}`). JSON API'da `null` va `[]` farqi muhim boʻlsa (GB-3.02), shu sozlama yoki vazifadagi `make([]Product, 0, …)`.
- **`:execrows`** — `RowsAffected()` sonini qaytaradi (vazifaning 🟡 qismi), **`:execresult`** — butun `pgconn.CommandTag`.
- **Tur oʻzgartirish chegarasi.** Bizning `Product.Stock` va `OrderItem.Quantity` `int`, sxemada `integer` (`int32`). Oʻqishda `int(r.Stock)` xavfsiz
  (kengayadi), yozishda `int32(it.Quantity)` kesib yuborishi mumkin, shuning uchun `PlaceOrder` miqdorni `1 … math.MaxInt32` oraligʻida tekshiradi.
  Boshqa yoʻl — `sqlc.yaml`'da `overrides` bilan `int4`ni Go `int`ga bogʻlash yoki domen turlarini `int32` qilish.
- **sqlc versiyasini loyihaga bogʻlash.** Bu laboratoriya goose kabi `go install …@v1.31.1` ishlatadi. Go 1.24'dan beri `go get -tool
  github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1` versiyani `go.mod`ning `tool` qatoriga yozadi va `go tool sqlc generate` bilan ishlaydi
  (koʻpincha `//go:generate go tool sqlc generate` qatori bilan); narxi — sqlc'ning koʻp bogʻliqliklari `go.sum`ga tushadi.
- **Oʻrnatish vaqti va cgo.** Birinchi `go install` bir necha daqiqa oladi (bu yerda ~4 daqiqa, 74 MB). C kompilyatori boʻlmasa Go 1.20'dan beri
  cgo'siz quradi va sqlc parser'ning WebAssembly versiyasini oladi; gcc boʻlsa libpg_query C'da kompilyatsiya qilinadi. Natija bir xil.
- **sqlc Docker image'i** ham bor (`sqlc/sqlc`), lekin bu yerda kerak emas: `go install` C kompilyatorisiz ishladi.

## Vazifa

[`TASK.md`](TASK.md): `ByName`, `Cheapest`, `OrderTotal` soʻrovlarini `query.sql`ga koʻchirish va `store.go`dagi qoʻlda yozilgan `Scan`ni
oʻchirish (🟡 `Restock` `:execrows` bilan). Tekshirish:

```bash
docker compose up -d --wait
go test .
```

Tugagach `docker compose down -v`. Videodagi holat va yechim orasidagi farq: `diff -u query.sql solution/query.sql` va `diff -u store.go solution/store.go`.

## Manbalar

- sqlc hujjatlari (v1.31.1): https://docs.sqlc.dev/en/v1.31.1/index.html ·
  oʻrnatish: https://docs.sqlc.dev/en/v1.31.1/overview/install.html ·
  soʻrov izohlari (`:one`, `:many`, `:exec`, `:execrows`): https://docs.sqlc.dev/en/v1.31.1/reference/query-annotations.html ·
  sozlama (`version: "2"`, `sql_package`, `emit_empty_slices`): https://docs.sqlc.dev/en/v1.31.1/reference/config.html ·
  sxema migratsiya fayllaridan, «sqlc ignores down migrations»: https://docs.sqlc.dev/en/v1.31.1/howto/ddl.html ·
  parametr nomlari (`@name`, `sqlc.arg`): https://docs.sqlc.dev/en/v1.31.1/howto/named_parameters.html ·
  CLI (`generate`, `diff`): https://docs.sqlc.dev/en/v1.31.1/reference/cli.html ·
  turlar (`pgtype.Text`): https://docs.sqlc.dev/en/v1.31.1/reference/datatypes.html
- sqlc manba kodi: `integer` → `int32` https://github.com/sqlc-dev/sqlc/blob/v1.31.1/internal/codegen/golang/postgresql_type.go ·
  cgo'siz parser https://github.com/sqlc-dev/sqlc/blob/v1.31.1/internal/engine/postgresql/parse_wasi.go ·
  https://github.com/wasilibs/go-pgquery (libpg_query WebAssembly'da) · https://github.com/pganalyze/pg_query_go
- pgx: `pgtype.Text` https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype#Text · `Tx` https://pkg.go.dev/github.com/jackc/pgx/v5#Tx
- PostgreSQL: sonli turlar (`integer` 4 bayt, `bigint` 8 bayt) https://www.postgresql.org/docs/18/datatype-numeric.html ·
  `sum(bigint)` → `numeric` https://www.postgresql.org/docs/18/functions-aggregate.html ·
  xato kodlari (`42703` undefined_column) https://www.postgresql.org/docs/18/errcodes-appendix.html
- goose izohlari (`-- +goose Up/Down`): https://pressly.github.io/goose/documentation/annotations/
- Go: `go vet` https://pkg.go.dev/cmd/vet · tool bogʻliqliklari (Go 1.24) https://go.dev/doc/modules/managing-dependencies#tools
