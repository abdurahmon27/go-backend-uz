# GB-5.07 · Transaction: buyurtma va qoldiq bitta butun ish

**Go backend: noldan guru'gacha** kursining 39-darsi. Buyurtma qo'yish uchta yozuv: `orders` qatori, `order_items` qatorlari va
`products.stock` kamayishi. Transaction'siz, ikkinchi yozuvdan keyingi xato yarim yozilgan buyurtma qoldiradi; `BEGIN` / `COMMIT` /
`ROLLBACK` (psql) va `pool.Begin` / `tx.Commit` / `defer tx.Rollback(ctx)` (pgx) bilan esa hammasi birga saqlanadi yoki hech narsa saqlanmaydi.

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go bo'lsa,
> `go run`/`go test` birinchi marta kerakli versiyani o'zi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qo'ygan bo'lsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-5-06` (goose migratsiyalari) · keyingi dars: `gb-5-08` (ikki xaridor, bitta mahsulot: lock).
Bu papka mustaqil laboratoriya: oldingi darsning kodi kerak emas. `store.go` GB-5.05 yechimidagi `Store` (`Product`, `Get`, `ListActive`, `OrderTotal`, `ByName`, `Cheapest`) ustiga shu darsning kodini qo'shadi; `bozor.sql` — GB-5.03/5.04 jadvallarining birlashgan nusxasi (`orders`, `order_items`, `products`). GB-5.06 migratsiyalari `orders.status` va `products.category` ustunlarini qo'shadi: bu darsning kodi ularga tegmaydi, shuning uchun ular bilan ham ishlaydi.

## Papka

| fayl | nima |
|---|---|
| `compose.yaml` | Postgres 18 (`postgres:18-alpine`), loyiha nomi `gb507`, `127.0.0.1:5497`, nomli volume |
| `bozor.sql` | `products`, `orders`, `order_items` jadvallari va boshlang'ich qatorlar; boshida `DROP TABLE IF EXISTS`, shuning uchun bazani boshlang'ich holatga qaytaradi |
| `bekor.sql`, `xato.sql` | videodagi ikkita psql tajribasi: `ROLLBACK` va transaction ichidagi xato |
| `store.go` | GB-5.05'dagi `Store` + `OrderItem`, `ErrNotEnoughStock`, `ErrInvalidQuantity`, `PlaceOrder` (transaction bilan, videodagidek) + tugallanmagan joylar (vazifa): xato kodini aylantirish, `Restock`, miqdorni tekshirish |
| `notx.go` | videodagi **xato** versiya `PlaceOrderNoTx`: bir xil buyruqlar, lekin transaction'siz. Faqat ko'rsatish uchun |
| `main.go` | videoning oxiridagi `main`: bitta buyurtma qo'yadi va natijani chop etadi |
| `store_test.go` | tekshiruv: har test alohida vaqtinchalik sxemada haqiqiy Postgres'da ishlaydi |
| `solution/` | to'liq javob (o'z `go.mod`i bilan, alohida ishga tushadi) |
| `_video/` | videodagi buyruqlar: `chiqishlar.sh` → `chiqishlar.txt`, har bosqichdagi `main.go` va xato versiya (`stages/`) |

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, pgx v5.13.0, Postgres 18). Videoda papka nomi `bozor`, shuning uchun konteyner `bozor-db-1`; sizda papka `gb-5-07-transaction`, konteyner `gb507-db-1`. To'liq jurnal:
[`_video/chiqishlar.txt`](_video/chiqishlar.txt).

```bash
docker compose up -d --wait
docker compose exec db psql -f /lab/bozor.sql
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5497/bozor
go run .      # PlaceOrderNoTx, olma x2 + non x5 (ombor: 3):
              # 0 ERROR: new row for relation "products" violates check constraint "products_stock_check" (SQLSTATE 23514)
```

Transaction'siz yarim yozilgan buyurtma (xatodan keyin):

```
 id | customer           order_id | product_id | quantity           id | name | stock
  4 | Aziza                     4 |          1 |        2            1 | olma |    38
```

`psql` ichida (`bekor.sql`: `BEGIN` → `UPDATE` → `ROLLBACK`; `-t` sarlavha va `(N rows)` qatorini yashiradi, `BEGIN`/`UPDATE 1`/`ROLLBACK` qoladi) va xato bilan (`xato.sql`):

```
BEGIN / UPDATE 1 / 36 / ROLLBACK / 38

BEGIN
psql:/lab/xato.sql:2: ERROR:  new row for relation "products" violates check constraint "products_stock_check"
DETAIL:  Failing row contains (2, non, 400000, -2, t).
psql:/lab/xato.sql:3: ERROR:  current transaction is aborted, commands ignored until end of transaction block
ROLLBACK
```

`PlaceOrder` (transaction bilan) xuddi shu buyurtmada: `0 ERROR: …` va bazada hech narsa o'zgarmaydi (orders/order_items: `(0 rows)`,
olma 40, non 3). To'g'ri buyurtma (olma x2, non x1): `5 <nil>` (4 emas: bekor bo'lgan `INSERT` ham identity raqamini oladi, GB-5.03);
`order_items`da 2 qator, olma 38, non 2. Xato: `Begin`dan keyin bitta buyruqni `tx` o'rniga `s.pool` orqali yuborish:
`0 ERROR: insert or update on table "order_items" violates foreign key constraint "order_items_order_id_fkey" (SQLSTATE 23503)`.

## Vazifa

[`TASK.md`](TASK.md): `PlaceOrder`da qoldiq yetmasligini `ErrNotEnoughStock`ga aylantiring, `Restock`ni bitta transaction'da yozing va `Quantity <= 0`ni `ErrInvalidQuantity` bilan rad eting.

```bash
docker compose up -d --wait
go test .            # hozir yiqiladi, vazifa bajarilgach o'tadi
docker compose down -v
```

Yechim: [`solution/`](solution/). `Store` oldingi Go darsidan (GB-5.05) qanday o'sganini `diff -ru ../gb-5-05-pgx-pool/solution .` ko'rsatadi (GB-5.06 faqat migratsiya darsi, Go kodi o'zgarmagan).

## Manbalar

- PostgreSQL 18: [Transactions](https://www.postgresql.org/docs/18/tutorial-transactions.html), [COMMIT](https://www.postgresql.org/docs/18/sql-commit.html), [Error codes (23514, 23503)](https://www.postgresql.org/docs/18/errcodes-appendix.html), [RETURNING](https://www.postgresql.org/docs/18/dml-returning.html), [UPDATE](https://www.postgresql.org/docs/18/sql-update.html), [Transaction isolation](https://www.postgresql.org/docs/18/transaction-iso.html), [Sequence functions](https://www.postgresql.org/docs/18/functions-sequence.html)
- pgx v5: [`Tx`](https://pkg.go.dev/github.com/jackc/pgx/v5#Tx), [`pgxpool.Pool.Begin`](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#Pool.Begin), [`pgconn.PgError`](https://pkg.go.dev/github.com/jackc/pgx/v5/pgconn#PgError)
