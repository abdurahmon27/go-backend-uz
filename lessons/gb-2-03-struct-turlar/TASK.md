# Vazifa · GB-2.03: `orderTotal` va `canCancel`

Videoda bozor'ning domen turlari paydo boʻldi: `Product`, `OrderItem`, `Order` va `OrderStatus`
(`StatusNew`, `StatusPaid`, `StatusShipped`, `StatusCancelled`). `main.go` oxirida shu turlar ustida
ishlaydigan ikkita funksiya bor, lekin ular hali boʻsh: biri doim `0`, ikkinchisi doim `false` qaytaradi.

**Nima qilish kerak** (`main.go`'da, faqat shu ikki funksiyaning tanasini yozasiz):

1. `func orderTotal(o Order) int64` buyurtmaning jami summasini **tiyinda** qaytarsin: har bir qatorda
   `Price × Quantity`, hammasi qoʻshiladi. Masalan, olma 1 200 000 tiyindan 3 dona va non 400 000 tiyindan
   2 dona: `3 600 000 + 800 000 = 4 400 000`. Qatorsiz buyurtmaning summasi `0`.
2. `func canCancel(o Order) bool` buyurtmani bekor qilish mumkinmi, aytsin: holati `StatusNew` yoki
   `StatusPaid` boʻlsa `true`, `StatusShipped` yoki `StatusCancelled` boʻlsa `false`.
   `var empty Order` ham yangi buyurtma (zero value'si `StatusNew`), uni ham bekor qilish mumkin.

Kirish qiymatlari toʻgʻri deb olinadi: manfiy son yoki juda katta summa (overflow) tekshiruvi bu vazifada
yoʻq, buyurtmani tekshirish keyingi darslarda.

Turlarni va holatlar roʻyxatini oʻzgartirmang: test holatlarning raqamlarini ham tekshiradi (0, 1, 2, 3).

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli `tekshiruv_test.go`: uni oʻzgartirmang. Hozir u shunday yiqiladi
(qisqartirilgan):

```
--- FAIL: TestVazifaOrderTotal (0.00s)
    --- FAIL: TestVazifaOrderTotal/bitta_qator (0.00s)
        tekshiruv_test.go:45: orderTotal({ID:0 Items:[{ProductID:1 Price:1200000 Quantity:3}] Status:0}) = 0; kutilgan 3600000 tiyin (TASK.md, 1-qadam)
…
--- FAIL: TestVazifaCanCancel (0.00s)
    --- FAIL: TestVazifaCanCancel/yangi (0.00s)
        tekshiruv_test.go:67: canCancel(Status: 0) = false; kutilgan true (TASK.md, 2-qadam)
…
```

Qoʻlda ham sinab koʻrish mumkin: `main` ichida `fmt.Println(orderTotal(o), canCancel(o))` yozib, `go run .`.

**Maslahatlar**

- `o.Items` — slice: `for _, it := range o.Items { … }` har bir qatorni `it`'ga beradi (GB-1.08'dagi
  `for _, tc := range tests` kabi).
- `it.Price` — `int64`, `it.Quantity` — `int`: koʻpaytirishdan oldin `int64(it.Quantity)` deb oʻgiring
  (GB-1.07'dagi `int64(som)`). Yigʻindi uchun `var total int64` va `+=` (GB-1.06).
- `canCancel` bitta `return` bilan yoziladi: ikkita `==` taqqoslash va `||` (GB-1.05).

Javob: [`solution/main.go`](solution/main.go).
