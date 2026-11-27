# Vazifa · GB-5.09: qolgan uchta soʻrovni sqlc'ga koʻchirish

`store.go`da `Get`, `ListActive` va `PlaceOrder` videodagidek sqlc yozgan kod orqali ishlaydi (`s.q.GetProduct`, `qtx.CreateOrder` …).
Uchta metod esa hali qoʻlda yozilgan, GB-5.05'dagidek: SQL matni Go string'ida va `Scan`:
`ByName`, `Cheapest` va `OrderTotal`. (`Restock` ham GB-5.07'dagidek qoʻlda; uni koʻchirish — ixtiyoriy 🟡 qadam.)

**Nima qilish kerak**

1. **`query.sql`ga uchta soʻrov qoʻshing**, aynan shu nomlar va natija turlari bilan (test shularni qidiradi):
   - `-- name: ProductByName :one` — `ByName`dagi `SELECT … WHERE name = $1`;
   - `-- name: CheapestProducts :many` — `Cheapest`dagi `SELECT … WHERE active ORDER BY price, id LIMIT $1`;
   - `-- name: OrderTotal :one` — `OrderTotal`dagi `SELECT COALESCE(SUM(quantity * unit_price), 0) … WHERE order_id = $1`.
2. **`sqlc generate`** (shu papkada). `db/query.sql.go`da uchta yangi metod paydo boʻladi. `OrderTotal`ning imzosi
   `OrderTotal(ctx context.Context, orderID int64) (int64, error)` boʻlishi kerak — 3-maslahatni oʻqing.
3. **`store.go`dagi uchta metodni `s.q` orqali qayta yozing.** Imzolar va xatti-harakat oʻzgarmaydi: `ByName` yoʻq nom uchun
   `ErrNotFound` qaytaradi, `Cheapest` faol mahsulot boʻlmasa **boʻsh slice** (`nil` emas) qaytaradi, `OrderTotal` qatorsiz buyurtma
   uchun `0`. Bu uch metodda `Scan` ham, `s.pool.Query`/`QueryRow` ham qolmasin: SQL faqat `query.sql`da. (Test metodlarni
   qaysi `.go` faylda boʻlmasin topadi va kodini `go/parser` bilan oʻqiydi.)

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test .
```

Hozir shunday yiqiladi (qolgan testlar oʻtadi: qoʻlda yozilgan kod ham ishlaydi, lekin vazifa — uni sqlc'ga koʻchirish):

```
--- FAIL: TestQueriesMovedToSqlc (0.00s)
    store_test.go:52: query.sql'da «-- name: ProductByName :one» soʻrovi yoʻq (TASK.md, 1-qadam)
    store_test.go:64: ByName hali qoʻlda yozilgan (Scan, pool.QueryRow): s.q.ProductByName orqali yozing (TASK.md, 3-qadam)
    store_test.go:52: query.sql'da «-- name: CheapestProducts :many» soʻrovi yoʻq (TASK.md, 1-qadam)
    store_test.go:64: Cheapest hali qoʻlda yozilgan (pool.Query, Scan): s.q.CheapestProducts orqali yozing (TASK.md, 3-qadam)
    store_test.go:52: query.sql'da «-- name: OrderTotal :one» soʻrovi yoʻq (TASK.md, 1-qadam)
    store_test.go:64: OrderTotal hali qoʻlda yozilgan (Scan, pool.QueryRow): s.q.OrderTotal orqali yozing (TASK.md, 3-qadam)
FAIL
FAIL	github.com/abdurahmon27/bozor	0.284s
FAIL
```

Test `query.sql`, `db/query.sql.go` va `store.go`ni oʻqiydi, keyin metodlarni haqiqiy Postgres'da ishga tushiradi. Har test alohida
vaqtinchalik sxema yaratadi, unga `migrations/*.sql`ning `-- +goose Up` qismini va `seed.sql`ni qoʻllaydi, oxirida sxemani oʻchiradi:
`goose up` shart emas va sizning `public` sxemangizga tegmaydi. Postgres boshqa manzilda boʻlsa `DATABASE_URL`ni qoʻying
(`postgres://bozor:bozor-lab@localhost:5409/bozor` — bu laboratoriyaning standart manzili). Postgres ishlamasa test aniq xabar
bilan yiqiladi, oʻtkazib yuborilmaydi. Test faylini oʻzgartirmaysiz; uning ichini GB-5.10 ochadi.

Tugagach: `docker compose down -v` (konteyner va uning maʼlumotlarini oʻchiradi).

**Maslahatlar**

1. sqlc oʻrnatilmagan boʻlsa: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1` (videodagidek), keyin `sqlc version` → `v1.31.1`.
   Birinchi oʻrnatish bir necha daqiqa oladi. C kompilyatori (gcc) boʻlmasa Go cgo'siz quradi (Go 1.20'dan beri oʻzi shunday qiladi) va sqlc
   Postgres parser'ining WebAssembly versiyasini ishlatadi; gcc boʻlsa C versiyasi quriladi. Ikkalasi ham ishlaydi.
   `query.sql` yoki migratsiyalar oʻzgargach har safar `sqlc generate`; `db/` papkasini qoʻlda tahrirlamang.
2. `Cheapest(ctx, n int)`: sqlc `LIMIT $1` uchun `int32` parametr yasaydi (`limit int32`), demak `int32(n)` kerak. Ikki holatda
   `int32(n)`dan oldin xato qaytaring (masalan `fmt.Errorf("Cheapest: n = %d …", n)`; test bu holatni tekshirmaydi): `n` 2 147 483 647 dan
   (`math.MaxInt32`, `PlaceOrder`dagi kabi) katta boʻlsa — int32'ga sigʻmaydi va kesilib ketadi (GB-1.03); `n` manfiy boʻlsa — Postgres
   `LIMIT`ga manfiy son bermaydi (`ERROR: LIMIT must not be negative`). `:many` qator boʻlmasa `nil` qaytaradi (videodagi `var items`),
   shuning uchun natijani `make([]Product, 0, len(rows))` slice'iga yigʻing: `Get`/`ListActive`dagi `Product{…, Stock: int(r.Stock)}` kabi.
3. `OrderTotal` — sqlc'ning chegarasi. Oddiy ustunlarning turini sqlc sxemadan biladi, **hisoblangan ifodalarning** turini esa oʻzi taxmin
   qiladi va bu taxmin Postgres'nikidan farq qilishi mumkin (shu yerda oʻlchangan, `_video/chiqishlar.txt`, `ifoda-turlari`, `sum-null`):
   - `COALESCE(SUM(quantity * unit_price), 0)` → sqlc `interface{}` qaytaradi, `store.go` kompilyatsiya boʻlmaydi:
     `cannot use 1st function result (value of type interface{}) as int64 value in return statement: need type assertion`;
   - `SUM(quantity * unit_price) AS total` (COALESCE'siz) → sqlc `int64` deydi, lekin qatorsiz buyurtmada Postgres `NULL` qaytaradi va test
     ishlaganda yiqiladi: `can't scan into dest[0] (col: total): cannot scan NULL into *int64`. Postgres'da `SUM` `bigint` ustida aslida
     `numeric` beradi (GB-5.05), sqlc buni ham bilmaydi.
   Yechim — natija turini SQL'ning oʻzida aniq yozish: `qiymat::bigint` Postgres'da qiymatni `bigint` turiga oʻtkazadi.
   `SELECT COALESCE(SUM(quantity * unit_price), 0)::bigint AS total …`; `AS total` natija ustuniga nom beradi. Shundan keyin metod
   `(int64, error)` qaytaradi va `NULL` boʻlmaydi. Qoida: hisoblangan ustundan keyin generatsiya qilingan imzoni oʻqing.
4. `ByName`: `Get` bilan bir xil shakl: `pgx.ErrNoRows` → `ErrNotFound`, keyin `Product{…}`.

**🟡 4-qadam, ixtiyoriy: `Restock` sqlc bilan.** `Restock` GB-5.07'dagidek qoʻlda yozilgan (`tx.QueryRow(…).Scan(&stock)`) va ishlaydi
(`TestRestock` oʻtadi); `TestRestockUsesSqlc` u qoʻlda ekan, xabar bilan oʻtkazib yuboriladi (SKIP). Koʻchirish: `query.sql`ga
`-- name: AddStock :execrows` (`UPDATE products SET stock = stock + @quantity WHERE id = @id`) yozing: `:execrows` yangilangan qatorlar sonini
`int64` qilib qaytaradi (GB-5.08'dagi `RowsAffected`). Tsiklda `qtx := s.q.WithTx(tx)` orqali `qtx.AddStock(…)`; son `0` boʻlsa mahsulot yoʻq:
`ErrNotFound` (transaction bekor boʻladi). Miqdor tekshiruvi oʻsha. Tekshirish: `go test -run Restock -v` (ikkala test PASS).

Javob: [`solution/`](solution/). Tekshirish: `cd solution && docker compose up -d --wait && go test ./...`
