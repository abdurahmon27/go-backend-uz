# Vazifa · GB-2.05: `AddItem` takrorni birlashtirsin, `Cancel` ishlasin

Videoda `Order` turi uchta metod oldi: `Total` (summani oʻqiydi), `Pay` va `AddItem` (buyurtmani oʻzgartiradi),
hammasi pointer receiver bilan (`func (o *Order) …`). `OrderStatus` esa `String` metodini oldi, shuning uchun
`fmt` holatni nomi bilan chiqaradi. `main.go`'da ikkita ish chala qolgan.

**Nima qilish kerak** (`main.go`'da):

1. **`AddItem`**: hozir u har chaqiruvda yangi qator qoʻshadi, shuning uchun `olma` ikki marta qoʻshilsa,
   buyurtmada ikkita olma qatori boʻladi. Shu mahsulot (`ProductID` bir xil) buyurtmada allaqachon boʻlsa,
   yangi qator ochilmasin: oʻsha qatorning `Quantity`'si `quantity`'ga oshsin. Masalan, olma 3, non 2,
   olma 1 → ikki qator: olma 4 dona, non 2 dona, `Total()` = 5 600 000 tiyin.
2. **`Cancel`**: metod bor, lekin u holatni umuman oʻzgartirmaydi. Avval sababini toping (videodagi `Pay`
   misolini eslang), keyin qoidani yozing: holati `StatusNew` yoki `StatusPaid` boʻlsa, holat `StatusCancelled`
   boʻlsin va `nil` qaytsin; `StatusShipped` yoki `StatusCancelled` boʻlsa, `error` qaytsin va holat
   oʻzgarmasin.

Turlarni, holatlar roʻyxatini va `String`'ni oʻzgartirmang. Kirish qiymatlari toʻgʻri deb olinadi
(`quantity` musbat): tekshiruv keyingi darslarda.

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli `tekshiruv_test.go`: uni oʻzgartirmang. Hozir u shunday yiqiladi
(qisqartirilgan):

```
--- FAIL: TestVazifaAddItem (0.00s)
    tekshiruv_test.go:38: olma 3, non 2, olma 1 qoʻshildi: len(Items) = 3; kutilgan 2 — olma ikkinchi qator boʻlib qoʻshilmasligi kerak (TASK.md, 1-qadam)
--- FAIL: TestVazifaCancel (0.00s)
    --- FAIL: TestVazifaCancel/yangi (0.00s)
        tekshiruv_test.go:75: Cancel()'dan keyin holat «yangi»; kutilgan «bekor qilingan» — receiver oʻzgarishni saqlaydimi? (TASK.md, 2-qadam)
…
```

Qoʻlda ham sinab koʻrish mumkin: `main`'da `order.AddItem(olma, 1)` qoʻshib, `go run .`.

**Maslahatlar**

- `for _, it := range o.Items { it.Quantity += quantity }` **ishlamaydi**: `it` — qatorning nusxasi, xuddi
  value receiver kabi. Asl qatorga indeks orqali yozing: `for i := range o.Items { … o.Items[i].Quantity … }`
  (indeks bilan yozish GB-2.01'da). `for i := range o.Items` faqat indeksni beradi: `i` = 0, 1, 2 …
  Topgach, qoʻshing va `return` yozing: natijasiz metodda yolgʻiz `return` metoddan shu yerda chiqadi,
  pastdagi `append`'gacha yetmaydi.
- `Cancel`'ning receiver'iga qarang: `Pay`'ning birinchi varianti bilan bir xil emasmi?
- Xato uchun `fmt.Errorf` (GB-1.07): `%v` holatni `String` orqali nomi bilan yozadi.

Javob: [`solution/main.go`](solution/main.go).
