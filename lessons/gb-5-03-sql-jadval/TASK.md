# Vazifa · GB-5.03: bo'sh nomni rad eting va arzon mahsulotlarni tanlang

`products.sql` videodagi jadval: id, name, price, stock, active; name `NOT NULL UNIQUE`, price va stock uchun `CHECK`.
Lekin bo'sh nom (`''`) hozir o'tib ketadi: `NOT NULL` faqat `NULL`ni rad etadi, bo'sh matn `NULL` emas.

**Nima qilish kerak**

1. **`products.sql`.** `name` ustuniga `CHECK (name <> '')` qo'shing (`<>` — "teng emas"). Boshqa ustunlarga tegmang
   (`id`, `name`, `price`, `stock`, `active` — tekshiruv shu nomlarga tayanadi).
2. **`arzon.sql`.** Bitta `SELECT` yozing (odatdagidek `;` bilan tugating): **faol** (`active` ustuni rost) mahsulotlardan
   narxi **3 000 000 tiyindan kam** (`<`, 3 000 000 oʻzi kirmaydi) bo'lgan **eng arzon ikkitasining** `name` va `price`
   ustunlari, narx **o'sish** tartibida. Fayldagi izohda `SELECT` skeleti bor.

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test .
```

Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestConstraintlar (3.27s)
    --- FAIL: TestConstraintlar/boʻsh_nom (0.47s)
        tekshiruv_test.go:95: boʻsh nom: qator qabul qilindi, baza rad etishi kerak edi. products.sql'da name ustuniga CHECK (name <> '') qoʻshing (TASK.md, 1-qadam)
--- FAIL: TestArzon (0.00s)
    tekshiruv_test.go:108: arzon.sql'da SELECT yoʻq: faol mahsulotlardan narxi 3 000 000 tiyindan kam boʻlgan eng arzon ikkitasini tanlang (TASK.md, 2-qadam)
FAIL
```

Tekshiruv `tekshiruv_test.go`da: u `products.sql` va `arzon.sql`ni alohida vaqtinchalik sxemada bajaradi (sizning
`public` sxemangizga tegmaydi) va oxirida o'chiradi. Go dasturidan Postgres'ga ulanishni GB-5.05'da ko'ramiz;
bu yerda test shunchaki `docker compose exec db psql` buyrug'ini ishga tushiradi. Docker ishlamasa test aniq xabar bilan
yiqiladi, o'tkazib yuborilmaydi.

Tugagach: `docker compose down -v` (konteyner va uning ma'lumotlarini o'chiradi).

**Maslahatlar**

- Ustun qoidalari ustun tavsifining oxirida bo'sh joy bilan ketma-ket yoziladi: `name text NOT NULL UNIQUE CHECK (…),`.
- Shartlar `AND` bilan birlashadi; `boolean` ustun shartda o'zi yoziladi (`WHERE active AND …`) yoki `active = true`.
- `ORDER BY price` — o'sish tartibi (`DESC` kamayish). Narxlari bir xil qatorlar tartibi kafolatlanmaydi: ishonch uchun
  `ORDER BY price, id`. Natijani cheklash: `LIMIT 2`.
- Tekshiruv `arzon.sql`ni ikki xil maʼlumotda ishga tushiradi (ikkinchisida narx 2 999 999, 3 000 000 va 3 000 001 boʻlgan
  mahsulotlar bor), shuning uchun natijani qoʻlda yozib qoʻyish oʻtmaydi.

Javob: [`solution/`](solution/). Tekshirish: `cd solution && docker compose up -d --wait && go test .`
