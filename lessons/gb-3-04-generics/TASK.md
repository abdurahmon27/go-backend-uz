# Vazifa · GB-3.04: `SumBy` va toʻlangan buyurtmalar summasi

Videoda bir xil kodning ikki nusxasi bitta generic funksiyaga yigʻildi: `Filter[T any]` mahsulot va buyurtmalar
uchun, `Sum[N Number]` `int` va `int64` uchun. Lekin `Sum` faqat sonlar slice'ini oladi, bozor'da esa bizda
`[]Product` va `[]Order` bor: «jami qoldiq», «katalog qiymati», «toʻlangan summa» har biri yana alohida sikl
boʻlib qoladi. `main.go`'da shu ish uchun ikki funksiya chala.

**Nima qilish kerak**

1. **`SumBy`** (`main.go`). Signatura tayyor:

   ```go
   func SumBy[T any, N any](items []T, value func(T) N) N
   ```

   Tur parametri ikkita, vergul bilan ajratilgan, har birining oʻz constraint'i bor: `T` — elementlar turi,
   `N` — yigʻiladigan son turi. `SumBy` har bir element uchun `value(x)`'ni chaqirib, natijalarni qoʻshib
   qaytarsin; boʻsh slice uchun `0`. Masalan:

   ```go
   SumBy(products, func(p Product) int { return p.Stock })                   // 43
   SumBy(products, func(p Product) int64 { return p.Price * int64(p.Stock) }) // 49200000
   ```

   Tanani yozsangiz, kompilyator videodagi xatoni beradi: `operator + not defined on total (variable of type N
   constrained by any)`. `N`'ning constraint'ini tanlang (videoda koʻrgan constraint'lardan qaysi birida `+` bor?). `T`'niki `any` qoladi:
   elementning oʻzi bilan `SumBy` hech narsa qilmaydi, faqat `value`'ga beradi.
2. **`paidTotal(orders []Order) int64`**: faqat toʻlangan (`StatusPaid`) buyurtmalarning summasi, tiyinda.
   Buyurtma summasi — `o.Total()` (GB-2.05). Oddiy `for` ham toʻgʻri boʻlardi, lekin bu yerda ikkala funksiyani ishlatib koʻring: `Filter` va `SumBy` bilan ikki qator.
   Tayyor boʻlsa, `go run .`'ning oxirgi qatori `jami: 4400000 tiyin` boʻladi.

`Filter`, `Sum`, `Number`, `Contains` va `main`'ni oʻzgartirmang. Test oʻz maʼlumoti bilan ishlaydi (10-buyurtma testda bor, `main`'da yoʻq), shuning uchun undagi summa `go run .`'dagidan boshqa.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. Test fayli (`tekshiruv_test.go`) oʻzgartirilmaydi. U `SumBy`'ni `Product`, `Order` va
`string` bilan chaqiradi: generic boʻlmagan funksiya bilan test kompilyatsiya ham boʻlmaydi. Hozir shunday yiqiladi:

```
--- FAIL: TestVazifaSumBy (0.00s)
    tekshiruv_test.go:55: SumBy(mahsulotlar, Stock) = 0; kutilgan 43 (40 + 3 + 0). Har bir x uchun value(x)'ni total'ga qoʻshing (TASK.md, 1-qadam)
    tekshiruv_test.go:59: SumBy(mahsulotlar, Price × Stock) = 0; kutilgan 49200000 tiyin (TASK.md, 1-qadam)
    tekshiruv_test.go:63: SumBy([olma non], len) = 0; kutilgan 7 (TASK.md, 1-qadam)
--- FAIL: TestVazifaPaidTotal (0.00s)
    tekshiruv_test.go:77: paidTotal = 0; kutilgan 6900000 tiyin (7 va 10-buyurtmalar, 8 toʻlanmagan) (TASK.md, 2-qadam)
FAIL
```

**Maslahatlar**

- 1-qadam: tana `Sum`'nikiga juda oʻxshaydi: `for _, x := range items { … }`, faqat `n` oʻrnida `value(x)`.
  Constraint kvadrat qavs ichida, `N`'dan keyin turadi.
- 2-qadam: `Filter(orders, func(o Order) bool { … })` — toʻlanganlar; keyin `SumBy(…, func(o Order) int64 { … })`.
  `N` bu yerda `int64` boʻladi: kompilyator uni `value` funksiyasining natijasidan chiqaradi.
- Testdagi uchinchi tekshiruv `len(s)` bilan: `len("olma")` — 4 bayt (GB-2.10), shuning uchun `olma` + `non` = 7.

Javob: [`solution/main.go`](solution/main.go).
