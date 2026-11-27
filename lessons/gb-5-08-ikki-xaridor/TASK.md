# Vazifa · GB-5.08: promo-kod parallel chaqiruvlarda

`coupon.go`da `UseCoupon(ctx, code)` bor: promo-kodni bir marta ishlatadi (`coupons.uses_left` bittaga kamayadi).
`bozor.sql`dagi `BAHOR` kodi **3 marta** ishlatiladi. Hozirgi versiya videodagi birinchi `Buy` kabi yozilgan:
transaction ichida `uses_left`ni **o'qiydi**, Go'da tekshiradi, keyin `uses_left = $1` deb **yozadi**. Bir vaqtda 10 kishi
urinsa, hammasi bir xil qiymatni o'qiydi va hammasiga «o'tdi» deydi.

**Nima qilish kerak**

1. **`UseCoupon`ni qayta yozing:** tekshiruv va kamaytirish **bitta** `UPDATE` buyrug'ida bo'lsin
   (`... SET uses_left = uses_left - 1 WHERE code = $1 AND uses_left > 0`) va `RowsAffected()` natijani aytsin.
   Bitta qator o'zgargan bo'lsa, `nil` qaytaring.
2. Nol qator o'zgargan bo'lsa, sabab ikki xil: promo-kod **tugagan** yoki **yo'q**. `UPDATE` buni aytmaydi
   (darsda: «nol qator: sabab noaniq»), shuning uchun nol qatordan keyin `uses_left`ni alohida o'qing:
   qator yo'q bo'lsa (`pgx.ErrNoRows`) `ErrCouponNotFound`, bor bo'lsa `ErrCouponUsedUp` qaytaring.
3. 🟡 **Ixtiyoriy:** `UseCouponLocked` yozing. Natija `UseCoupon` bilan bir xil bo'lsin, lekin boshqa yo'l bilan:
   transaction ichida `SELECT uses_left ... FOR UPDATE` (`tx.QueryRow`, `pool` emas!), Go'da tekshiruv, keyin `UPDATE`.
   Yozilmagan bo'lsa, bu test «o'tkazib yuborildi» (SKIP) deb chiqadi; yozsangiz, u ham parallel tekshiriladi.

`store_test.go`, `coupon_test.go` va `naive_test.go` — tekshiruv fayllari: ularni o'zgartirmang.

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test .
```

Hozir shunday yiqiladi:

```
--- FAIL: TestUseCouponConcurrent (0.06s)
    coupon_test.go:61: 1-urinish: promo-kod 3 marta, 10 kishi bir vaqtda: 10 tasiga oʻtdi; kutilgan: 3 (TASK.md)
FAIL
```

(«10 tasiga» o'rniga 6 dan 10 gacha chiqishi mumkin: parallel testda aniq son har safar boshqacha, lekin 3 emas.) Vazifa tugagach
`go test .` `ok` beradi. Test har safar alohida vaqtinchalik sxema yaratadi va oxirida o'chiradi: sizning `public`
sxemangizga tegmaydi. Postgres boshqa manzilda bo'lsa `DATABASE_URL`ni qo'ying (standart: `postgres://bozor:bozor-lab@localhost:5458/bozor`).
Postgres ishlamasa test aniq xabar bilan yiqiladi, o'tkazib yuborilmaydi. Tugagach: `docker compose down -v`.

**Maslahatlar**

- `Buy`ning kodiga qarang (`store.go`): `Exec`, keyin `tag.RowsAffected()`. `$1` parametrini ikki joyda ishlatish mumkin.
- `UPDATE` o'zi bo'linmas: shartli versiyada transaction kerak emas, `s.pool.Exec` yetadi.
- Nol qatordan keyingi o'qish `GB-5.05`dagi `Get`ga o'xshaydi: `QueryRow(...).Scan(&left)`, keyin `errors.Is(err, pgx.ErrNoRows)`.
- Ixtiyoriy qadam uchun `BuyLocked`ni nusxalab, jadval va ustun nomlarini almashtiring; `ErrNoRows` → `ErrCouponNotFound`.

Javob: [`solution/`](solution/) (`coupon.go`). Tekshirish: `cd solution && docker compose up -d --wait && go test .`
