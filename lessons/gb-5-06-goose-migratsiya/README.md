# GB-5.06 · goose migratsiyalari: sxemani versiya bilan o'zgartirish

**Go backend: noldan guru'gacha** kursining 38-darsi. bozor jadvallari (`products`, `orders`, `order_items`, GB-5.03 va GB-5.04)
endi qo'lda ishga tushiriladigan SQL fayllar emas, raqamlangan **migratsiya** fayllari: `migrations/00001_…sql`, har birida
`-- +goose Up` va `-- +goose Down` bo'limi. Goose (`go install` bilan o'rnatiladi) ularni bazaga bir martadan qo'llaydi,
qaysilari qo'llanganini `goose_db_version` jadvalida saqlaydi, xato bergan migratsiyani tranzaksiyada qaytaradi.

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go bo'lsa,
> `go run`/`go test` birinchi marta kerakli versiyani o'zi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qo'ygan bo'lsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-5-05` (Go → Postgres: pgx) · keyingi dars: `gb-5-07` (Transaction: buyurtma va qoldiq).
Bu papka mustaqil laboratoriya: oldingi darsning Go kodi kerak emas. Migratsiyalar GB-5.03/5.04'dagi jadvallarni yaratadi, shuning uchun GB-5.05'dagi `Store` shu sxemaga to'g'ridan-to'g'ri mos keladi.
Videoda compose loyihasi `bozor` deb ataladi (`bozor-db-1`), papkadagi `compose.yaml`da esa `name: gb506` (`gb506-db-1`): buyruqlar bir xil, faqat konteyner nomi boshqacha.

## Papka

| fayl | nima |
|---|---|
| `compose.yaml` | Postgres 18 (`postgres:18-alpine`), loyiha nomi `gb506`, `127.0.0.1:5486`, nomli volume |
| `products.sql`, `orders.sql` | oldingi darslardagi qo'lda ishga tushiriladigan fayllar (videoda `orders.sql` ikkinchi marta xato beradi) |
| `seed.sql` | to'rt mahsulot, uch buyurtma va ularning to'rt qatori (migratsiyalardan keyin qo'llanadi) |
| `migrations/` | `00001_products.sql`, `00002_orders.sql`, `00003_orders_status.sql` (`DEFAULT 'yangi'` va `CHECK` bilan): videoning oxiridagi holat |
| `migrate.go` | migratsiyalarni `go:embed` bilan binary'ga o'raydi va goose kutubxonasi uchun `NewProvider` beradi (Postgres sessiya qulfi bilan: bir necha nusxa bir vaqtda ishga tushsa ham migratsiyani bittasi qo'llaydi). Chuqurroq; GB-8.06'da ishlatamiz |
| `migrate_test.go` | tekshiruv: har test alohida vaqtinchalik sxemada haqiqiy Postgres'da goose bilan Up, Down, Up qiladi |
| `solution/` | to'liq javob (`migrations/00004_products_category.sql` bilan; o'z `go.mod`i, alohida ishga tushadi) |
| `_video/` | videodagi buyruqlar: `chiqishlar.sh` → `chiqishlar.txt` |

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, goose v3.30.2, pgx v5.13.0, Postgres 18).
To'liq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt). Goose har chiqish qatori oldiga sana va vaqt yozadi
(`2026/11/27 18:39:52 OK   00001_products.sql (10.27ms)`); videoda bu qism qisqartirilgan.

```bash
docker compose exec db psql -f /lab/orders.sql         # ERROR: relation "orders" already exists (qo'lda fayl ikkinchi marta)
docker compose down -v && docker compose up -d --wait  # bo'sh baza
go install github.com/pressly/goose/v3/cmd/goose@v3.30.2
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5486/bozor
goose status        # Pending -- 00001_products.sql, Pending -- 00002_orders.sql
goose up            # OK 00001…, OK 00002…, goose: successfully migrated database to version: 2
goose up            # goose: no migrations to run. current version: 2
docker compose exec db psql -f /lab/seed.sql
goose create -s orders_status sql   # Created new file: migrations/00003_orders_status.sql
goose up            # NOT NULL ustun, qatorlar bor: ERROR: column "status" of relation "orders" contains null values (SQLSTATE 23502)
                    # DEFAULT 'yangi' (va CHECK) qo'shilgach: OK 00003…, version 3
goose down          # OK 00003_orders_status.sql (faqat oxirgisi)
```

Xato bergan migratsiyadan keyin `orders` jadvalida `status` ustuni yo'q: goose odatda har migratsiya faylini bitta tranzaksiyada ishga
tushiradi (`-- +goose NO TRANSACTION` bilan o'chiriladi; shunda xatodan oldingi buyruqlar qaytarilmaydi), Postgres esa `ALTER TABLE`ni ham qaytara oladi.
Goose fayl mazmunini tekshirmaydi: qo'llangan faylni tahrir qilsangiz u jim qoladi (`no migrations to run`), shuning uchun test `00001`–`00003` fayllarining hash'ini tekshiradi.

## Vazifa

[`TASK.md`](TASK.md): `00004` migratsiya: `products.category` (`text NOT NULL DEFAULT 'boshqa'`, `CHECK (category <> '')`), Down ustunni o'chiradi.

```bash
docker compose up -d --wait
go test .
```

Tugagach `docker compose down -v`.

## Manbalar

- goose: https://pressly.github.io/goose/ (SQL migratsiyalar, `-- +goose Up/Down`, `NO TRANSACTION`: `documentation/annotations/`; CLI: `documentation/cli-commands/`;
  muhit o'zgaruvchilari: `documentation/environment-variables/`) · repo: https://github.com/pressly/goose (v3.30.2, `internal/sqlparser/parser.go`: `useTx = true` odatiy)
- `go install`: https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies · `embed`: https://pkg.go.dev/embed
- Postgres: `ALTER TABLE` https://www.postgresql.org/docs/18/sql-altertable.html · xato kodlari (23502 `not_null_violation`) https://www.postgresql.org/docs/18/errcodes-appendix.html

Video: kanaldagi «Go backend» playlist'ida (#38). Yozma dars: haywan.uz/lab (gb: Go backend).
