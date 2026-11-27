# Vazifa · GB-5.06: 00004 migratsiya, `products.category`

`migrations/` papkasida videodagi uchta migratsiya bor (`00001_products.sql`, `00002_orders.sql`, `00003_orders_status.sql`).
To'rtinchisi hali yo'q: test hozir shuni aytib yiqiladi.

**Nima qilish kerak**

1. `migrations/` papkasida `00004_` bilan boshlanadigan yangi `.sql` fayl yarating (masalan `00004_products_category.sql`).
   Buni `goose create -s products_category sql` qilib beradi (videodagidek: `GOOSE_DRIVER`, `GOOSE_DBSTRING`,
   `GOOSE_MIGRATION_DIR` qo'yilgan bo'lsa), yoki faylni qo'lda yozing.
2. **`-- +goose Up`** bo'limida `products` jadvaliga `category` ustunini qo'shing: turi `text`, `NOT NULL`,
   `DEFAULT 'boshqa'`, va bo'sh matn bo'lmasin (`CHECK (category <> '')`). Jadvalda qatorlar bo'lishi mumkin:
   `DEFAULT`siz `NOT NULL` ustun qatorlari bor jadvalga qo'shilmaydi (videodagi xato).
3. **`-- +goose Down`** bo'limi ustunni o'chirsin (`DROP COLUMN`).

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test .
```

Test har safar alohida vaqtinchalik sxema yaratadi: 1–3-migratsiyalarni qo'llaydi, ikkita mahsulot yozadi, keyin sizning
4-migratsiyangizni qo'llaydi (mavjud mahsulotlarda `category = 'boshqa'` bo'lishi kerak, bo'sh matn rad etilishi kerak),
`Down` qiladi (ustun yo'qolishi, versiya 3 bo'lishi kerak) va yana `Up` qiladi. Oxirida sxemani o'chiradi: sizning
`public` sxemangizga va `goose_db_version` jadvalingizga tegmaydi. Postgres boshqa manzilda bo'lsa `DATABASE_URL`ni qo'ying
(`postgres://bozor:bozor-lab@localhost:5486/bozor` bu laboratoriyaning standart manzili). Postgres ishlamasa test aniq
xabar bilan yiqiladi, o'tkazib yuborilmaydi.

Hozir shunday yiqiladi:

```
--- FAIL: TestCategoryMigration (0.01s)
    migrate_test.go:171: 4-migratsiya yoʻq: migrations/ papkasida 00004_ bilan boshlanadigan .sql fayl yarating (TASK.md, 1-qadam)
FAIL
```

O'z bazangizda ham sinab ko'ring (ixtiyoriy):

```bash
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5486/bozor
goose up && goose down && goose up
```

Tugagach: `docker compose down -v` (konteyner va uning ma'lumotlarini o'chiradi).

**Maslahatlar**

- `00003_orders_status.sql`ga qarang: u aynan shu shakldagi migratsiya: `ALTER TABLE … ADD COLUMN … NOT NULL DEFAULT … CHECK (…)` va `ALTER TABLE … DROP COLUMN …`.
- Qo'llangan migratsiya faylini (`00001`–`00003`) o'zgartirmang: hamma o'zgarish yangi raqamli faylda bo'ladi; test buni ham tekshiradi.

Javob: [`solution/`](solution/). Tekshirish: `cd solution && docker compose up -d --wait && go test .`
