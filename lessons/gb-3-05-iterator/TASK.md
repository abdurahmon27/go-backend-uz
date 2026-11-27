# Vazifa · GB-3.05: `LowStock` iteratori va `SortedNames`

Videoda `Store`'ga `All() iter.Seq[Product]` metodini yozdik: u ichki `products` slice'ini ochmasdan va
nusxalamasdan mahsulotlarni birma-bir `yield`'ga beradi, `yield` `false` qaytarsa (sikl `break` qildi) darhol
`return` qiladi. Map kalitlarini esa `slices.Sorted(maps.Keys(count))` bilan bir qatorda tartibladik.

**Nima qilish kerak** (`internal/store/store.go`'da, faqat `LowStock` va `SortedNames`'ning tanasini
oʻzgartirasiz):

1. **`LowStock(limit int) iter.Seq[Product]`** — qoldigʻi (`Stock`) `limit`'dan **kam** (`<`) mahsulotlarni
   katalog tartibida beradigan iterator. Videodagi katalogda `LowStock(5)` `non` va `oʻrik`'ni beradi.
   - Sikl `break` qilsa, iterator darhol toʻxtashi shart: `yield` `false` qaytargandan keyin uni yana chaqirish
     — videodagi `range function continued iteration …` panic'i.
   - Nusxa yigʻmang: iterator `s.products`'ni `range` paytida oʻqisin (videodagi `All` kabi jonli koʻrinish).
     Test buni tekshiradi: `LowStock(5)`'ni olib, keyin katalogga `anor` (qoldiq 1) qoʻshadi va u ham chiqishini kutadi.
   - `main.go`'da shunday ishlatiladi:

     ```go
     for p := range st.LowStock(5) {
     	fmt.Println("kam qolgan:", p.Name, p.Stock)
     }
     ```
2. **`SortedNames() []string`** — mahsulot nomlari, tartiblangan (`slices.Sort` tartibi: baytlar boʻyicha),
   **yangi** slice'da. `Names()` iteratori papkada tayyor: undan `slices.Sorted` bilan yigʻing. Videodagi
   katalog uchun: `[choy non olma oʻrik]`. Natijaga yozish katalogni oʻzgartirmasligi kerak.

`All`, `Names`, `Add`, `List`'ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test ./...
```

`ok` chiqquncha ishlang. Test fayli `internal/store/tekshiruv_test.go`: uni oʻzgartirmang. U iteratorni ham
`for … range` bilan, ham `yield`'ni oʻzi berib oddiy funksiya kabi chaqiradi: shu yoʻl bilan `yield` `false`
qaytargandan keyin iterator yana chaqirmayotganini sanaydi. Hozir u shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaLowStock (0.00s)
    --- FAIL: TestVazifaLowStock/beshdan_kam (0.00s)
        tekshiruv_test.go:77: LowStock(5) [] berdi, kutilgan [non oʻrik]: qoldigʻi limit'dan kam (<) mahsulotlarni katalog tartibida bering (TASK.md, 1-qadam)
    …
--- FAIL: TestVazifaSortedNames (0.00s)
    tekshiruv_test.go:112: SortedNames() [] qaytardi, kutilgan [choy non olma oʻrik]: Names() iteratoridan slices.Sorted bilan yigʻing (TASK.md, 2-qadam)
FAIL
```

**Maslahatlar**

- `LowStock`'ning tuzilishi `All`'niki bilan bir xil: `return func(yield func(Product) bool) { … }`, ichida
  `for _, p := range s.products`. Qoldigʻi yetarli mahsulotni `continue` bilan oʻtkazing (GB-1.06), qolganini
  `yield`'ga bering.
- `if !yield(p) { return }` — har `yield`'dan keyin. Shusiz testning `yield false` tekshiruvi yiqiladi.
- `slices.Sorted(seq)` tartiblanadigan qiymatlar (sonlar, `string`'lar) ustidagi `iter.Seq`'dan yangi tartiblangan
  slice yasaydi: `maps.Keys(m)` ham, sizning `s.Names()` ham `iter.Seq[string]`. (`slices.Sorted(st.All())`
  kompilyatsiya boʻlmaydi: `Product`'ni `<` bilan solishtirib boʻlmaydi.) `slices` papkada allaqachon import qilingan.
- Qoʻlda ham tekshiring: `main.go`'ga yuqoridagi siklni qoʻshib, `go run ./cmd/bozor`.

Javob: [`solution/internal/store/store.go`](solution/internal/store/store.go).
