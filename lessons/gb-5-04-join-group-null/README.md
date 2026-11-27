# GB-5.04 · JOIN, GROUP BY va NULL: buyurtmalar va ularning jami

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #36"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 dan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: GB-5.03 (CREATE TABLE, INSERT, SELECT va constraint) · keyingi dars: GB-5.05 (Go → Postgres: pgx).

**Bu darsdan keyin:** `orders` va `order_items` jadvallarini `REFERENCES` (foreign key) bilan yaratasiz, ikki jadvalni
`JOIN` va `LEFT JOIN` bilan bogʻlaysiz, `GROUP BY` + `SUM` bilan har buyurtmaning jamini hisoblaysiz, `NULL` nega `=` bilan
taqqoslanmasligini koʻrasiz va `IS NULL`, `COALESCE` yozasiz.

Bu papka 5-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`),
oʻz `compose.yaml`'i (bitta Postgres: `postgres:18-alpine`, `127.0.0.1:5464`). `products` jadvali GB-5.03 nihoyasidagi bilan
bir xil (`id`, `name`, `price` tiyinda, `stock`, `active`); `orders` va `order_items` shu darsda yaratiladi.

## Papka

| Fayl | Nima |
|---|---|
| `compose.yaml` | bitta Postgres; `psql` uchun `PGUSER`/`PGDATABASE`; papka konteynerda `/lab` ostida koʻrinadi |
| `products.sql` | GB-5.03'dagi `products` jadvali va toʻrt mahsulot (olma, non, choy, oʻrik) |
| `orders.sql` | `orders` va `order_items` jadvallari (`REFERENCES`) |
| `seed.sql` | uchta buyurtma (Dilnoza, Jasur, Malika) va toʻrt qator; Malikaning buyurtmasida qator yoʻq |
| `sotilgan.sql`, `izohsiz.sql` | vazifa: har biriga bitta `SELECT` yoziladi |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `sotilgan.sql`, `izohsiz.sql` (+ xuddi shu fayllar, test va `compose.yaml`: loyiha `gb504-solution`, port 5466, shu bois ikkala papkani birga ishlatsa boʻladi), oʻz `go.mod`'i bor |
| `_video/` | `chiqishlar.sh` (videodagi hamma buyruqni ketma-ket ishga tushiradi) va `chiqishlar.txt` (uning haqiqiy natijasi) |

Bu darsda `products.sql`, `orders.sql`, `seed.sql` oʻzgarmaydi. Boshqa 5-modul darslarida (GB-5.05 dan) `orders` va
`order_items` jadvallariga `migrations/` orqali oʻzgartirish kiritiladi.

## Videodagi buyruqlar va natijalari (postgres:18-alpine = PostgreSQL 18.7, Docker Compose 5.4.0)

Video `bozor/` papkasida yozilgan (GB-5.02, 5.03 kabi), shuning uchun compose loyihasi `bozor`, konteyner `bozor-db-1`; bu
papkada `compose.yaml`dagi `name: gb504` tufayli `gb504-db-1`. Har bir soʻrov `sorov.sql` fayliga yoziladi va
`docker compose exec db psql -f /lab/sorov.sql` bilan ishga tushadi (hammasi `_video/chiqishlar.txt`da; ishga tushirish:
`./_video/chiqishlar.sh`).

Baza, `products`, ikki yangi jadval va maʼlumot:

```bash
docker compose up -d --wait
docker compose exec db psql -f /lab/products.sql   # CREATE TABLE, INSERT 0 4
docker compose exec db psql -f /lab/orders.sql     # CREATE TABLE, CREATE TABLE
docker compose exec db psql -f /lab/seed.sql       # INSERT 0 3, INSERT 0 4
```

`orders`: `id | customer | note` (Dilnoza: izoh NULL, Jasur: «ertalab yetkazing», Malika: izoh NULL);
`order_items`: 1-buyurtmada 2 ta qator (olma ×2, non ×3), 2-buyurtmada 2 ta (oʻrik ×1, olma ×1), 3-buyurtmada yoʻq.

Tashqi kalit — mavjud boʻlmagan buyurtma (99) bilan qator rad etiladi:

```
ERROR:  insert or update on table "order_items" violates foreign key constraint "order_items_order_id_fkey"
DETAIL:  Key (order_id)=(99) is not present in table "orders".
```

`JOIN` (INNER) toʻrt qator beradi, Malika yoʻq (`sorov.sql`: `SELECT orders.customer, order_items.quantity FROM orders JOIN order_items ON order_items.order_id = orders.id;`):

```
 customer | quantity 
----------+----------
 Dilnoza  |        2
 Dilnoza  |        3
 Jasur    |        1
 Jasur    |        1
(4 rows)
```

`GROUP BY`siz `SUM`:

```
psql:/lab/sorov.sql:4: ERROR:  column "orders.id" must appear in the GROUP BY clause or be used in an aggregate function
LINE 1: SELECT orders.id, orders.customer,
               ^
```

`GROUP BY orders.id` + `ORDER BY orders.id` bilan (hali `JOIN`): Malika yoʻq.

```
 id | customer |  jami   
----+----------+---------
  1 | Dilnoza  | 3600000
  2 | Jasur    | 4700000
(2 rows)
```

`WHERE note = NULL` hech narsa bermaydi, `IS NULL` esa ikki buyurtmani beradi (Dilnoza, Malika):

```
 id | customer 
----+----------
(0 rows)
```

`LEFT JOIN` + `GROUP BY`: Malikaning yagona qatorida `quantity` NULL, `SUM` NULL qiymatlarni oʻtkazib yuboradi, shuning uchun jami NULL (boʻsh); `COALESCE(SUM(…), 0)` bilan 0:

```
 id | customer |  jami   
----+----------+---------
  1 | Dilnoza  | 3600000
  2 | Jasur    | 4700000
  3 | Malika   |       0
(3 rows)
```

Eslatma: `SUM(bigint)` Postgres'da `numeric` turini qaytaradi (butun son boʻlsa ham); Go'ga oʻqishda `::bigint` bilan oʻtkazish GB-5.05'da.
Eslatma: Postgres tashqi kalit ustuniga (`order_id`) indeksni oʻzi qoʻymaydi; kerak boʻlsa `CREATE INDEX` bilan oʻzingiz yaratasiz (bu darsda yoʻq).

Chuqurroq: `WHERE note <> 'ertalab yetkazing'` ham 0 qator beradi (NULL `<>` matn → NULL); `IS DISTINCT FROM` Dilnoza va
Malikani beradi.

Oldingi darsdan farq: `products.sql` GB-5.03 `solution/products.sql` jadvali bilan bir xil; bu yerda `DROP TABLE IF EXISTS`
qatori yoʻq, toʻrt mahsulotning `INSERT`i (GB-5.03'dagi `seed.sql`) shu faylning oxirida:
`diff -u ../gb-5-03-sql-jadval/solution/products.sql products.sql`.

## Vazifa

[`TASK.md`](TASK.md): `sotilgan.sql` (har mahsulotning sotilgan donasi, sotilmaganda 0) va `izohsiz.sql` (izohi `NULL`
buyurtmalar). Tekshirish:

```bash
docker compose up -d --wait
go test .
docker compose down -v      # tugagach: konteyner va ma'lumotni oʻchirish
```

Javob: [`solution/`](solution/) (`cd solution && docker compose up -d --wait && go test .`).

## Manbalar

- PostgreSQL 18: [Joins Between Tables](https://www.postgresql.org/docs/18/tutorial-join.html),
  [Joined tables](https://www.postgresql.org/docs/18/queries-table-expressions.html#QUERIES-JOIN),
  [GROUP BY](https://www.postgresql.org/docs/18/queries-table-expressions.html#QUERIES-GROUP),
  [Aggregate functions](https://www.postgresql.org/docs/18/functions-aggregate.html),
  [Foreign keys](https://www.postgresql.org/docs/18/ddl-constraints.html#DDL-CONSTRAINTS-FK),
  [Comparison operators](https://www.postgresql.org/docs/18/functions-comparison.html),
  [COALESCE](https://www.postgresql.org/docs/18/functions-conditional.html#FUNCTIONS-COALESCE-NVL-IFNULL).
