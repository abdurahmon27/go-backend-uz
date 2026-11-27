# Vazifa · GB-5.05: ByName va Cheapest

`store.go`da `Store` bor: `Get` (bitta mahsulot id bo'yicha), `ListActive` (faol mahsulotlar) va `OrderTotal` (buyurtma jami)
videodagidek yozilgan.
Ikkita metod hali yozilmagan: ular hozir «hali yozilmagan» xatosini qaytaradi.

**Nima qilish kerak**

1. **`ByName(ctx, name)`** — nom bo'yicha bitta mahsulot (`Product`). Mahsulot topilmasa `ErrNotFound` qaytaring
   (`pgx.ErrNoRows`ni o'zingiz aylantiring, `Get`dagidek). Nomni SQL matniga **qo'shmang**: `$1` parametridan foydalaning.
   Test `x' OR 'a'='a` va `olma'; DROP TABLE products; --` kabi nomlarni ham sinaydi.
2. **`Cheapest(ctx, n)`** — eng arzon `n` ta **faol** mahsulot, narx **o'sish** tartibida (narxlari teng bo'lsa `id` bo'yicha).
   `n`ni ham `$1` parametri qiling (`LIMIT $1`). Natija bo'sh bo'lsa xatosiz, `nil` bo'lmagan bo'sh slice (`[]Product{}`) qaytsin (JSON'da `null` emas, `[]` bo'lishi uchun). `rows.Close` va
   `rows.Err`ni unutmang.

Ikkala metodda ham `ctx`ni `QueryRow`/`Query`ga bering: test bekor qilingan context bilan chaqirib, xato kutadi.

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test .
```

Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestByName (0.04s)
    store_test.go:130: ByName("choy") = {ID:0 Name: Price:0 Stock:0}, ByName hali yozilmagan; kutilgan: {3 choy 2500000 8}, nil (TASK.md, 1-qadam)
--- FAIL: TestCheapest (0.04s)
    store_test.go:154: Cheapest(2) = [], Cheapest hali yozilmagan; kutilgan: [non olma], nil (TASK.md, 2-qadam)
FAIL
```

Videoda `pool_max_conns=1` bilan `export` qilgan bo'lsangiz, `unset DATABASE_URL` qiling yoki yangi terminal oching
(test o'zi ham pool'ni to'rtta ulanish bilan ochadi). Test har safar alohida vaqtinchalik sxema yaratadi (`products.sql`, `seed.sql`, `orders.sql` va `orders_seed.sql` shu sxemaga qo'llanadi) va oxirida
uni o'chiradi: sizning `public` sxemangizdagi jadvalga tegmaydi. Postgres boshqa manzilda bo'lsa `DATABASE_URL`ni qo'ying
(`postgres://bozor:bozor-lab@localhost:5475/bozor` — bu laboratoriyaning standart manzili). Postgres ishlamasa test aniq
xabar bilan yiqiladi, o'tkazib yuborilmaydi.

Tugagach: `docker compose down -v` (konteyner va uning ma'lumotlarini o'chiradi).

**Maslahatlar**

- `Get`ning kodiga qarang: `QueryRow(ctx, "SELECT … WHERE name = $1", name).Scan(…)`, keyin `errors.Is(err, pgx.ErrNoRows)`.
- `Cheapest` uchun `ListActive`ni nusxalang va so'rovni o'zgartiring: `WHERE active ORDER BY price, id LIMIT $1`.
- `Query` xatosidan keyin darhol `defer rows.Close()`, tsikldan keyin `return out, rows.Err()`.

Javob: [`solution/`](solution/). Tekshirish: `cd solution && docker compose up -d --wait && go test .`
