# Vazifa · GB-2.06: `Sell` qoidasi va `List` nusxasi

Videoda bozor ikki paketga boʻlindi: `cmd/bozor` (dastur) va `internal/store` (katalog). `Store`'ning
`products` maydoni kichik harf bilan, shuning uchun `main` unga yoza olmaydi: katalog faqat `store`'ning
metodlari orqali oʻzgaradi. Endi shu metodlar qoidani bajarishi kerak. `internal/store/store.go`'da ikkita ish
chala qolgan.

**Nima qilish kerak** (`internal/store/store.go`'da):

1. **`Sell(id int64, quantity int) error`**: hozir u hech narsa qilmaydi va har doim `nil` qaytaradi. Qoidani
   yozing:
   - `quantity` musbat boʻlmasa (0 yoki manfiy) — `error`;
   - `id`'li mahsulot katalogda boʻlmasa — `error`;
   - qoldiq (`Stock`) `quantity`'dan kam boʻlsa — `error`;
   - aks holda shu mahsulotning `Stock`'i `quantity`'ga kamayadi va `nil` qaytadi.
   Xato boʻlganda hech bir qoldiq oʻzgarmasin. Masalan, non (ID 2, qoldiq 100): `Sell(2, 3)` → `nil`, qoldiq 97;
   olma (ID 1, qoldiq 40): `Sell(1, 41)` → `error`, qoldiq 40.
2. **`List() []Product`**: hozir u `s.products`'ning oʻzini qaytaradi. Slice sarlavhasi nusxalanadi, lekin
   massiv umumiy (GB-2.01), shuning uchun `st.List()[0].Stock = -5` katalogdagi qoldiqni ham oʻzgartiradi:
   `main` kichik harfli maydon himoyasini shu yoʻl bilan chetlab oʻtadi. `List` **nusxa** qaytarsin: natijaga
   yozish katalogga taʼsir qilmasin.

`Product`, `Store`'ning maydonlari va `Add`'ni oʻzgartirmang.

**Tekshirish** (papkaning ildizida):

```bash
go test ./...
```

`ok` chiqquncha ishlang. Test fayli `internal/store/tekshiruv_test.go`: uni oʻzgartirmang (u `store` paketining
ichida, shuning uchun `products`'ni ham koʻradi). Hozir u shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaSell (0.00s)
    tekshiruv_test.go:31: Sell(2, 3)'dan keyin non qoldigʻi 100; kutilgan 97 (TASK.md, 1-qadam)
    --- FAIL: TestVazifaSell/qoldiqdan_kop (0.00s)
        tekshiruv_test.go:49: Sell(1, 41) = nil; kutilgan error (TASK.md, 1-qadam)
…
--- FAIL: TestVazifaList (0.00s)
    tekshiruv_test.go:68: List() natijasida olma qoldigʻi -5 qilindi va katalogdagi qoldiq ham -5 boʻldi; kutilgan 40: List nusxa qaytarsin (TASK.md, 2-qadam)
FAIL
```

Qoʻlda ham sinab koʻrish mumkin: `cmd/bozor/main.go`'da `fmt.Println(st.Sell(2, 3))` qoʻshib (va `"fmt"`'ni
import qilib), `go run ./cmd/bozor`.

**Maslahatlar**

- Mahsulotni indeks bilan qidiring: `for i := range s.products { if s.products[i].ID == id { … } }`.
  `for _, p := range s.products` ichida `p.Stock -= quantity` **ishlamaydi**: `p` — elementning nusxasi.
- Xato uchun `fmt.Errorf` (GB-1.07), masalan `fmt.Errorf("mahsulot %d topilmadi", id)`. `store.go`'ga
  `import "fmt"` qoʻshing.
- Nusxa uchun standart kutubxonadagi `slices` paketi: `slices.Clone(s.products)` yangi massivli nusxa qaytaradi
  (`import "slices"`). Yoki qoʻlda: `out := make([]Product, len(s.products))`, keyin `copy(out, s.products)`.

Javob: [`solution/internal/store/store.go`](solution/internal/store/store.go).
