# GB-5.03 · Kerakli SQL: `CREATE TABLE`, `INSERT`, `SELECT` va constraint

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #35"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 dan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: GB-5.02 (Postgres Docker Compose'da va psql bilan) · keyingi dars: GB-5.04 (JOIN, GROUP BY va NULL).

**Bu darsdan keyin:** Postgres'da `products` jadvalini qoʻlda yaratasiz (`bigint`, `text`, `integer`, `boolean`,
`GENERATED ALWAYS AS IDENTITY`, `PRIMARY KEY`), `INSERT` bilan qator qoʻshasiz, `SELECT … WHERE … ORDER BY … LIMIT` bilan
tanlaysiz, `RETURNING` bilan yangi `id`ni olasiz va `NOT NULL`, `UNIQUE`, `CHECK`, `DEFAULT` bilan notoʻgʻri qatorni
baza oʻzi rad etishini koʻrasiz.

Bu papka 5-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`),
oʻz `compose.yaml`'i (bitta Postgres: `postgres:18-alpine`, `127.0.0.1:5463`). Jadval 3-modul va 4-modul bozoridagi
`Product` bilan bir xil nomlarga ega: `id`, `name`, `price` (tiyinda, `bigint`), `stock`, `active`.

## Papka

| Fayl | Nima |
|---|---|
| `compose.yaml` | bitta Postgres; `psql` uchun `PGUSER`/`PGDATABASE`; papka konteynerda `/lab` ostida koʻrinadi |
| `products.sql` | videoning oxiridagi jadval (constraint'lar bilan); vazifaning 1-qadami shu faylda |
| `seed.sql` | toʻrt mahsulot: olma, non, choy, oʻrik (narxlar tiyinda) |
| `arzon.sql` | vazifaning 2-qadami: bitta `SELECT` shu yerga yoziladi |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `products.sql`, `arzon.sql` (+ xuddi shu test va `compose.yaml`) |
| `_video/` | videodagi qadamlar: `01-oddiy/products.sql` (faqat turlar, constraint yoʻq), `02-qat/products.sql` (constraint'lar bilan), `chiqishlar.sh` (hamma buyruqni ketma-ket ishga tushiradi) va `chiqishlar.txt` (uning haqiqiy natijasi) |

## Videodagi buyruqlar va natijalari (postgres:18-alpine = PostgreSQL 18.7, Docker Compose 5.4.0)

Bazani koʻtarish va boʻsh ekanini koʻrish (`psql` konteyner ichida, `PGUSER`/`PGDATABASE` `compose.yaml`'da).
Video `bozor/` papkasida yozilgan (GB-5.02 kabi), shuning uchun compose loyihasi `bozor`, konteyner `bozor-db-1`;
bu papkada `compose.yaml`dagi `name: gb503` tufayli `gb503-db-1` boʻladi (boshqa lablar bilan toʻqnashmaslik uchun):

```bash
docker compose up -d --wait
docker compose exec db psql -c '\dt'
```
```
 Container bozor-db-1 Healthy
Did not find any tables.
```

Jadval (`_video/01-oddiy/products.sql`, faqat turlar) va toʻrt qator (`seed.sql`):

```bash
# videoda shu paytda products.sql — _video/01-oddiy/products.sql (papkadagi fayl esa oxirgi, constraint'li versiya:
# u DROP TABLE IF EXISTS bilan boshlanadi, birinchi ishga tushirishda NOTICE chiqadi)
docker compose exec db psql -f /lab/products.sql      # CREATE TABLE
docker compose exec db psql -f /lab/seed.sql          # INSERT 0 4
docker compose exec db psql -c 'SELECT * FROM products'
```
```
 id | name  |  price  | stock | active 
----+-------+---------+-------+--------
  1 | olma  | 1200000 |    40 | t
  2 | non   |  400000 |     3 | t
  3 | choy  | 2500000 |     8 | t
  4 | oʻrik | 3500000 |     2 | t
(4 rows)
```

Tanlash (`WHERE`, `AND`, `ORDER BY … DESC`, `LIMIT`): `stock < 5` → non, oʻrik; `… AND price > 1000000` → oʻrik;
`stock < 5 ORDER BY price DESC` → oʻrik, non; `ORDER BY price DESC LIMIT 2` → oʻrik, choy. Hammasi
`_video/chiqishlar.txt`da.

Klassik xato — matn qoʻsh tirnoqda (SQL'da qoʻsh tirnoq ustun nomi):

```
ERROR:  column "olma" does not exist
LINE 1: SELECT name FROM products WHERE name = "olma"
                                               ^
```

`RETURNING` va muammo: constraint'siz jadval boʻsh nomni va manfiy narxni qabul qiladi:

```
 id | name | price | stock | active 
----+------+-------+-------+--------
  6 |      |    -5 |       | 
(1 row)

INSERT 0 1
```

Jadvalni constraint'lar bilan qayta yaratgach (`_video/02-qat/products.sql`) xuddi shu qator rad etiladi — faqat `CHECK`,
chunki boʻsh matn `NULL` emas, `NOT NULL`dan oʻtadi:

```
ERROR:  new row for relation "products" violates check constraint "products_price_check"
DETAIL:  Failing row contains (5, , -5, 0, t).
```

Ikkinchi `choy` (`UNIQUE`) va `NULL` nom (`NOT NULL`):

```
ERROR:  duplicate key value violates unique constraint "products_name_key"
DETAIL:  Key (name)=(choy) already exists.

ERROR:  null value in column "name" of relation "products" violates not-null constraint
DETAIL:  Failing row contains (7, null, 900000, 0, t).
```

Chuqurroq: rad etilgan `INSERT`lar ham identity hisoblagichidan raqam oladi, shuning uchun keyingi qator `id` = 8 oladi
(hisoblagich orqaga qaytmaydi; id'lar orasida bo'shliq bo'lishi mumkin).

## Vazifa

[`TASK.md`](TASK.md): `products.sql`'da `name` uchun `CHECK (name <> '')` va `arzon.sql`'da bitta `SELECT`.
Tekshirish:

```bash
docker compose up -d --wait
go test .
docker compose down -v      # tugagach: konteyner va ma'lumotni oʻchirish
```

Javob: [`solution/`](solution/) (`cd solution && docker compose up -d --wait && go test .`).

## Manbalar

- PostgreSQL 18: [CREATE TABLE](https://www.postgresql.org/docs/18/sql-createtable.html),
  [Constraints](https://www.postgresql.org/docs/18/ddl-constraints.html),
  [Identity columns](https://www.postgresql.org/docs/18/ddl-identity-columns.html),
  [INSERT](https://www.postgresql.org/docs/18/sql-insert.html),
  [SELECT](https://www.postgresql.org/docs/18/sql-select.html),
  [ORDER BY](https://www.postgresql.org/docs/18/queries-order.html),
  [RETURNING](https://www.postgresql.org/docs/18/dml-returning.html),
  [psql](https://www.postgresql.org/docs/18/app-psql.html).
