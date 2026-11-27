# GB-5.05 · Go'dan Postgres'ga: pgx, pool, `Scan` va timeout

**Go backend: noldan guru'gacha** kursining 37-darsi. Bozor katalogi (`products` jadvali, GB-5.03) endi Go'dan o'qiladi:
`pgxpool` bilan ulanamiz, `$1` parametrli so'rov yuboramiz, `QueryRow`/`Query` natijasini `Scan` bilan `Product`ga
o'qiymiz, `pgx.ErrNoRows`ni o'z `ErrNotFound`imizga aylantiramiz va har so'rovga `context` bilan muddat qo'yamiz.

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go bo'lsa,
> `go run`/`go test` birinchi marta kerakli versiyani o'zi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qo'ygan bo'lsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-5-04` (JOIN, GROUP BY va NULL) · keyingi dars: `gb-5-06` (goose migratsiyalari).
Bu papka mustaqil laboratoriya: oldingi darslarning kodi kerak emas (`products.sql`, `seed.sql` GB-5.03 yechimining, `orders.sql`
va `orders_seed.sql` GB-5.04 yechimining nusxasi).
Bu dars Go kodini birinchi marta qo'shadi, shuning uchun oldingi `solution/` bilan `diff` qilinadigan Go kodi yo'q.

## Papka

| fayl | nima |
|---|---|
| `compose.yaml` | Postgres 18 (`postgres:18-alpine`), loyiha nomi `gb505`, `127.0.0.1:5475`, nomli volume |
| `products.sql`, `seed.sql` | GB-5.03'dagi jadval (`CHECK (name <> '')` bilan) va to'rt mahsulot |
| `orders.sql`, `orders_seed.sql` | GB-5.04'dagi `orders`, `order_items` va ularning qatorlari |
| `store.go` | `Product`, `ErrNotFound`, `Store`: `Get`, `ListActive`, `OrderTotal` (videodagidek) + ikkita yozilmagan metod (vazifa) |
| `main.go` | videoning oxiridagi `main`: pool yaratadi, so'rovlarga 2 soniyalik `qctx` beradi, `Get(1)`, `ListActive` va `OrderTotal`ni chop etadi |
| `store_test.go` | tekshiruv: har test alohida vaqtinchalik sxemada haqiqiy Postgres'da ishlaydi |
| `solution/` | to'liq javob (o'z `go.mod`i bilan, alohida ishga tushadi) |
| `_video/` | videodagi buyruqlar: `chiqishlar.sh` → `chiqishlar.txt`, har bosqichdagi fayllar (`stages/`, `leak/`) |

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, pgx v5.13.0, Postgres 18).
To'liq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt).

```bash
docker compose up -d --wait
docker compose exec db psql -f /lab/products.sql -f /lab/seed.sql -f /lab/orders.sql -f /lab/orders_seed.sql
export DATABASE_URL=postgres://bozor:bozor-lab@localhost:5475/bozor
go mod tidy          # go: found github.com/jackc/pgx/v5/pgxpool in github.com/jackc/pgx/v5 v5.13.0
go run .             # 1 olma 1200000 40   (id 99 bo'lsa: xato: no rows in result set → xato: mahsulot topilmadi)
```

Eslatma: video `bozor` nomli papkada yozilgan (compose loyiha nomi = papka nomi, konteyner `bozor-db-1`); bu papkadagi
`compose.yaml` `name: gb505` beradi (`gb505-db-1`), buyruqlar bir xil.

`SUM(quantity*unit_price)` (numeric) butun bo'lsa `int64`ga o'qiladi (`3600000 <nil>`); buyurtma yo'q bo'lsa `SUM` NULL beradi va
`Scan` xato qaytaradi: `can't scan into dest[0] (col: sum): cannot scan NULL into *int64`; `COALESCE(SUM(...), 0)` bilan `0 <nil>`.
Kasrli numeric'ni `int64`ga o'qib bo'lmaydi (`SELECT 1.5::numeric` → `cannot convert … to integer`).

Muddat (`pg_sleep(5)`, muddat 2 soniya) va `rows.Close` unutilgan holat (`pool_max_conns=1`):

```
xato: timeout: context deadline exceeded 2s
---
1 olma <nil> 0s
2  context deadline exceeded 2s        ← rows.Close yo'q: ulanish pool'ga qaytmadi
---
1 olma <nil> 0s
2 olma <nil> 0s                        ← defer rows.Close() bilan
```

Postgres'ga ulanib bo'lmasa (port noto'g'ri) `pgxpool.New` xato bermaydi, xato birinchi so'rovda keladi:
`failed to connect to … connection refused` (`pool.Ping(ctx)` ulanishni darhol tekshiradi; `_video/chiqishlar.txt`, bosqich `lazy`).

## Vazifa

[`TASK.md`](TASK.md): `ByName` va `Cheapest` metodlarini yozing. Tekshirish:

```bash
docker compose up -d --wait
go test .
```

Tugagach `docker compose down -v`.

## Manbalar

- pgx: https://pkg.go.dev/github.com/jackc/pgx/v5 · pool: https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool
- `Rows` (`Close`, `Err`, `Next`): https://pkg.go.dev/github.com/jackc/pgx/v5#Rows
- `context.WithTimeout`: https://pkg.go.dev/context#WithTimeout
- Postgres: ulanish URI'lari https://www.postgresql.org/docs/18/libpq-connect.html#LIBPQ-CONNSTRING-URIS ·
  parametrli so'rov (extended query) https://www.postgresql.org/docs/18/protocol-flow.html#PROTOCOL-FLOW-EXT-QUERY ·
  `pg_sleep` https://www.postgresql.org/docs/18/functions-datetime.html#FUNCTIONS-DATETIME-DELAY

Video: kanaldagi «Go backend» playlist'ida (#37). Yozma dars: haywan.uz/lab (gb: Go backend).
