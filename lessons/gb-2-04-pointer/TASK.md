# Vazifa · GB-2.04: `sellByID` va `cheapest`

Videoda `sell(p *Product, quantity int)` qoldiqni pointer orqali kamaytirdi, `findProduct` esa katalogdagi
elementning manzilini (`&products[i]`) yoki topilmasa `nil` qaytardi. `main.go` oxirida shular ustida ishlaydigan
ikkita funksiya bor, lekin ular hali boʻsh: ikkalasi ham doim `nil` qaytaradi.

**Nima qilish kerak** (`main.go`'da, faqat shu ikki funksiyaning tanasini yozasiz):

1. `func sellByID(products []Product, id int64, quantity int) error` katalogdan `id` raqamli mahsulotni topib,
   undan `quantity` dona sotsin: **katalogdagi elementning oʻzi** kamaysin. Uch holatda sotmasin, qoldiqqa
   tegmasin va xato qaytarsin:
   - mahsulot topilmasa (masalan `id` = 9): panic emas, xato;
   - `quantity` noldan katta boʻlmasa (`0`, `-5`);
   - `quantity` qoldiqdan katta boʻlsa (choy 7 ta, 8 ta soʻralsa).
   Qoldiqni toʻliq sotish mumkin: choy 7 ta, 7 ta soʻralsa qoldiq `0`, xato yoʻq.
2. `func cheapest(products []Product) *Product` eng arzon (`Price` eng kichik) mahsulotning **katalogdagi
   manzilini** qaytarsin; katalog boʻsh boʻlsa `nil`. Test qaytgan pointer orqali yozib, katalogdagi element
   oʻzgarganini tekshiradi: nusxaning manzili oʻtmaydi.

Xato matni sizning ixtiyoringizda (`fmt.Errorf`, GB-1.07), test faqat xato bor-yoʻqligini tekshiradi.

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli `tekshiruv_test.go`: uni oʻzgartirmang. Hozir u shunday yiqiladi
(qisqartirilgan):

```
--- FAIL: TestVazifaSellByID (0.00s)
    --- FAIL: TestVazifaSellByID/olma_3_dona (0.00s)
        tekshiruv_test.go:44: sellByID(id 1, 3 dona) dan keyin olma qoldigʻi 40; kutilgan 37 (TASK.md, 1-qadam)
…
    --- FAIL: TestVazifaSellByID/yoq_mahsulot (0.00s)
        tekshiruv_test.go:37: sellByID(id 9, 1 dona): xato kutilgan, nil qaytdi (TASK.md, 1-qadam)
…
--- FAIL: TestVazifaCheapest (0.00s)
    --- FAIL: TestVazifaCheapest/uch_mahsulot (0.00s)
        tekshiruv_test.go:63: cheapest(olma, non, choy) = nil; kutilgan non'ning manzili (TASK.md, 2-qadam)
FAIL
```

**Maslahatlar**

- `sellByID`: avval `p := findProduct(products, id)`, keyin videodagidek `if p == nil { return … }`. Qolgan ikki
  shart ham `if` bilan: `quantity <= 0` va `quantity > p.Stock`. Hammasi joyida boʻlsa `sell(p, quantity)`
  va `return nil`.
- `cheapest`: `var best *Product` (zero value'si `nil`) va `for i := range products`. `best == nil` yoki
  `products[i].Price < best.Price` boʻlsa, `best = &products[i]`.
- `for _, p := range products` + `&p` ishlamaydi: videoning «Chuqurroq» qismida koʻrdingiz, `p` elementning
  nusxasi.

Javob: [`solution/main.go`](solution/main.go).
