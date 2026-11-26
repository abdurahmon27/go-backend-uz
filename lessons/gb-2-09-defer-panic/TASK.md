# Vazifa · GB-2.09: `handle` panic'ni xatoga aylantirsin, `reserve` — hammasi yoki hech biri

Videoda `defer f.Close()` faylni har qanday `return`'da yopdi, `safely` esa `defer` ichidagi nomsiz funksiyada
`recover` bilan panic'ni toʻxtatib, uni nomli natija `err`'ga yozdi. `main.go`'da ikkita ish chala qolgan.

**Nima qilish kerak** (`main.go`'da):

1. **`handle(name string, fn func() error) (err error)`** — bozor buyrugʻini (`fn`) bajaradigan chegara.
   - `fn` xato qaytarsa, `handle` aynan shu xatoni **oʻzgartirmay** qaytarsin; xato boʻlmasa — `nil`.
   - `fn` panic qilsa, dastur yiqilmasin: `handle` xato qaytarsin, uning matni
     `«<name>: panic: <panic qiymati>»` boʻlsin, masalan `sotish: panic: kutilmagan holat`, va chaqiruvchi uni
     `errors.Is(err, ErrPanic)` bilan taniy olsin (`ErrPanic` `main.go`'da eʼlon qilingan).

   Hozir `handle` shunchaki `return fn()` qiladi: panic toʻgʻridan-toʻgʻri yuqoriga chiqadi.
2. **`reserve(products []Product, items []OrderItem) (err error)`** — buyurtma qatorlari uchun qoldiqni
   (`Stock`) kamaytiradi. Qoida: **hammasi yoki hech biri**. Biror qatorda xato boʻlsa (dona soni musbat emas,
   mahsulot yoʻq yoki qoldiq yetmaydi), undan oldingi qatorlarda kamaytirilgan qoldiqlar joyiga qaytsin, xato esa
   oʻzgarmay qaytarilsin (GB-2.08'dagidek `errors.Is(err, ErrNotFound)`, `errors.Is(err, ErrBadQuantity)`,
   `errors.As(err, &se)` bilan tekshiriladi). Hozir `reserve` xatoni qaytaradi, lekin oldingi qatorlar kamaygancha
   qoladi: olma 40 → 37, keyin non'da xato, olma esa 37 boʻlib qoladi.

`saveReport`'ni, `safely`'ni, `findProduct`'ni, xato turlarini va xato matnlarini oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli `tekshiruv_test.go`: uni oʻzgartirmang (unda ham `defer` va `recover` bor:
`handle` panic'ni ushlamasa, test yiqilmaydi, aniq xabar beradi). Hozir u shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaHandle (0.00s)
    --- FAIL: TestVazifaHandle/sotish (0.00s)
        tekshiruv_test.go:89: handle("sotish") panic'ni ushlamadi («kutilmagan holat»): handle'da defer va recover kerak (TASK.md, 1-qadam)
    …
--- FAIL: TestVazifaReserve (0.00s)
    --- FAIL: TestVazifaReserve/qoldiq_yetmaydi (0.00s)
        tekshiruv_test.go:144: reserve xato qaytardi («non: qoldiq 100, soʻralgan 500»), lekin qoldiq: olma 37, non 100, oʻrik 5; kutilgan 40, 100, 5: xato boʻlsa, oldin kamaytirilgan qoldiqlarni defer bilan qaytaring (TASK.md, 2-qadam)
    …
FAIL
```

**Maslahatlar**

- 1-qadam `safely`'ning oʻzi, ikki farq bilan: `fn` endi `error` qaytaradi (`return fn()`), xato matni esa
  `fmt.Errorf("%s: %w: %v", name, ErrPanic, r)` — `%w` `ErrPanic`'ni ichida saqlaydi (GB-2.08), `%v` panic qiymatini
  matnga qoʻshadi.
- 2-qadam: `defer` ichidagi nomsiz funksiya oʻzgaruvchini funksiya **tugaganda** oʻqiydi (videodagi `stock` 37).
  Boshida `ok := false` yozing, oxirgi `return nil`'dan oldin `ok = true`; `defer`'dagi funksiya `!ok` boʻlsa
  qoldiqlarni tiklasin. Shunda har bir `return fmt.Errorf(…)`'ga tegmaysiz, va oʻrtada panic boʻlsa ham qoldiq
  tiklanadi (panic'da ham defer'lar ishlaydi).
- Nimani tiklash: eng oddiysi — boshida har bir mahsulotning `Stock`'ini `[]int` slice'ga `append` bilan yigʻish
  (GB-2.01), keyin `defer`'da ularni `products[i].Stock`'ga qaytarish. Bitta mahsulot ikki
  qatorda kelishi mumkin (testdagi `bir_mahsulot_ikki_qator`): bu usul uni ham toʻgʻri tiklaydi.

Javob: [`solution/main.go`](solution/main.go).
