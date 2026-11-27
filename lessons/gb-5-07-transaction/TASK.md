# Vazifa · GB-5.07: ErrNotEnoughStock, Restock va miqdorni tekshirish

`store.go`da `PlaceOrder` videodagidek bitta transaction'da yozilgan va ishlaydi: biror qadam xato bersa hech narsa yozilmaydi.
`ErrNotEnoughStock` va `ErrInvalidQuantity` e'lon qilingan, lekin hali ishlatilmaydi; `Restock` yozilmagan (hozir «hali yozilmagan» xatosini qaytaradi).

**Nima qilish kerak**

1. **`PlaceOrder`: `ErrNotEnoughStock`.** `UPDATE products SET stock = stock - …` `CHECK (stock >= 0)`ni buzsa, Postgres xato kodi `23514`
   (SQLSTATE) bilan javob beradi. Videoning «Chuqurroq» qismidagidek (`errors.As` + `*pgconn.PgError` + `Code`) shu xatoni taniting va
   `ErrNotEnoughStock` qaytaring (chaqiruvchi `errors.Is` bilan tekshiradi). Transaction oldingidek bekor bo'lishi kerak.
2. **`Restock(ctx, items []OrderItem) error`: tovar qabul qilish.** Har `OrderItem` uchun mahsulot qoldig'i `Quantity` ga **oshadi**.
   Hammasi bitta transaction'da: `items` ichida mavjud bo'lmagan mahsulot bo'lsa `ErrNotFound` qaytaring va **hech bir** mahsulotning
   qoldig'i o'zgarmasin.
3. **Miqdorni tekshirish.** `Quantity <= 0` bo'lsa `PlaceOrder` ham, `Restock` ham `ErrInvalidQuantity` qaytarsin: **`Begin`dan oldin**,
   shuning uchun hech narsa yozilmaydi. (Tekshirmasangiz, `Restock`ga `-2` berilsa qoldiq jimgina kamayib ketadi.)

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test .
```

Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestPlaceOrderNotEnoughStock (0.04s)
    store_test.go:154: PlaceOrder(non x5) = ERROR: new row for relation "products" violates check constraint "products_stock_check" (SQLSTATE 23514); kutilgan: ErrNotEnoughStock (TASK.md, 1-qadam: errors.As bilan *pgconn.PgError, Code 23514)
--- FAIL: TestRestock (0.04s)
    store_test.go:165: Restock = Restock hali yozilmagan; kutilgan: nil (TASK.md, 2-qadam)
--- FAIL: TestRestockAllOrNothing (0.04s)
    store_test.go:179: Restock(99-mahsulot) = Restock hali yozilmagan; kutilgan: ErrNotFound (TASK.md, 2-qadam)
--- FAIL: TestInvalidQuantity (0.04s)
    store_test.go:190: PlaceOrder(miqdor 0) = ERROR: new row for relation "order_items" violates check constraint "order_items_quantity_check" (SQLSTATE 23514); kutilgan: ErrInvalidQuantity (TASK.md, 3-qadam)
FAIL
```

Test har safar alohida vaqtinchalik sxema yaratadi (`bozor.sql` shu sxemaga qo'llanadi) va oxirida uni o'chiradi: sizning `public`
sxemangizga tegmaydi. Postgres boshqa manzilda bo'lsa `DATABASE_URL`ni qo'ying
(`postgres://bozor:bozor-lab@localhost:5497/bozor` — bu laboratoriyaning standart manzili). Postgres ishlamasa test aniq xabar
bilan yiqiladi, o'tkazib yuborilmaydi. Test faylini o'zgartirmaysiz; uning ichini GB-5.10 ochadi.

Tugagach: `docker compose down -v` (konteyner va uning ma'lumotlarini o'chiradi).

**Maslahatlar**

- 1-qadam: `var pgErr *pgconn.PgError` (import: `github.com/jackc/pgx/v5/pgconn`), keyin `if errors.As(err, &pgErr) && pgErr.Code == "23514" { … }`
  `UPDATE`ning `Scan` xatosidan keyin, `if err != nil`dan oldin.
- 2-qadam: `PlaceOrder`ni nusxalab boshlang: `Begin`, `defer tx.Rollback(ctx)`, tsiklda
  `tx.QueryRow(… "UPDATE products SET stock = stock + $1 WHERE id = $2 RETURNING stock" …).Scan(&stock)`, `pgx.ErrNoRows` → `ErrNotFound`,
  oxirida `return tx.Commit(ctx)`.
- 3-qadam: funksiyaning boshida `items` bo'ylab tsikl: `if it.Quantity <= 0 { return … ErrInvalidQuantity }`.
- Hamma buyruq `tx` orqali ketadi, `s.pool` orqali emas (videodagi xato).

Javob: [`solution/`](solution/). Tekshirish: `cd solution && docker compose up -d --wait && go test .`
