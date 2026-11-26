# Vazifa · GB-2.08: chaqiruvchi xatoni tekshira olsin

Videoda `store.Sell` ikki xatoni chaqiruvchi tekshira oladigan qildi: mahsulot topilmasa
`fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)` (sentinel, `errors.Is` bilan topiladi), qoldiq yetmasa
`&StockError{…}` (typed error, `errors.As` bilan olinadi). Ikki ish chala qolgan.

**Nima qilish kerak**

1. **`internal/store/store.go`, `Sell`**: dona soni (`quantity`) musbat boʻlmasa, `Sell` hozir oddiy matnli xato
   qaytaradi: `dona soni musbat boʻlishi kerak, hozir 0`. Chaqiruvchi uni `errors.Is` bilan taniy olmaydi. Faylda
   `ErrBadQuantity` sentinel'i allaqachon eʼlon qilingan: `Sell` shu xatoni `%w` bilan oʻrab qaytarsin, matnda
   soʻralgan son ham qolsin (masalan `dona soni musbat boʻlishi kerak, hozir -2`). Qoldiq oʻzgarmasin.
2. **`cmd/bozor/main.go`, `message`**: `message(err)` xatoni foydalanuvchiga koʻrsatiladigan matnga aylantiradi.
   Hozir u har qanday xatoga `xato: …` qaytaradi. Kerak:

   | xato (zanjirning istalgan qavatida) | `message` qaytaradi |
   |---|---|
   | `store.ErrNotFound` | `bunday mahsulot yoʻq` |
   | `store.ErrBadQuantity` | `dona soni notoʻgʻri` |
   | `*store.StockError` | `olma: faqat 40 dona qoldi` (nom va qoldiq xatoning maydonlaridan) |
   | boshqa har qanday xato | `xato: ` + `err.Error()` (hozirgidek) |

   Testdagi baʼzi xatolar yana bir qavat oʻralgan (`buyurtma 7: …`): ularni ham tanish kerak.

`main`'ni, `StockError`'ni va `ErrNotFound`'ni oʻzgartirmang.

**Tekshirish** (shu papkada, ildizida):

```bash
go test ./...
```

Ikkala paket `ok` boʻlguncha ishlang. Test fayllari (`internal/store/tekshiruv_test.go`,
`cmd/bozor/tekshiruv_test.go`) oʻzgartirilmaydi. Hozir ular shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaMessage (0.00s)
    --- FAIL: TestVazifaMessage/topilmadi (0.00s)
        tekshiruv_test.go:44: message("mahsulot 9: mahsulot topilmadi") = "xato: mahsulot 9: mahsulot topilmadi"; kutilgan "bunday mahsulot yoʻq" — errors.Is(err, store.ErrNotFound) (TASK.md, 2-qadam)
…
FAIL	github.com/abdurahmon27/bozor/cmd/bozor	0.002s
--- FAIL: TestVazifaDonaSoni (0.00s)
    tekshiruv_test.go:60: Sell(1, 0) xatosi "dona soni musbat boʻlishi kerak, hozir 0", lekin errors.Is(err, ErrBadQuantity) = false: ErrBadQuantity'ni fmt.Errorf'da %w bilan oʻrab qaytaring (TASK.md, 1-qadam)
…
FAIL	github.com/abdurahmon27/bozor/internal/store	0.002s
```

**Maslahatlar**

- 1-qadam: `ErrBadQuantity`'ning matni `dona soni musbat boʻlishi kerak` — hozirgi xato matnining birinchi qismi.
  Shablonda `%w` xatoning matnini qoʻyadi, sonni esa `%d` (videodagi `"mahsulot %d: %w"`'ga qarang).
- 2-qadam: tartib — avval ikki sentinel (`errors.Is`), keyin tur (`errors.As`), oxirida qolgan hammasi. Har bir
  tekshiruv `if … { return … }`. `==` ishlamaydi: xato oʻralgan.
- `errors.As` uchun: `var se *store.StockError`, keyin `errors.As(err, &se)`; matnni `fmt.Sprintf` bilan
  `se.Name` va `se.Stock`'dan yigʻing. `errors.AsType[*store.StockError](err)` ham toʻgʻri.

Javob: [`solution/`](solution/) (`solution/internal/store/store.go`, `solution/cmd/bozor/main.go`).
