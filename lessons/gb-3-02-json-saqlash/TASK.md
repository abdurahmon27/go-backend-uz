# Vazifa · GB-3.02: buyurtma JSON'i va ishonchli `loadProducts`

Videoda katalog `bozor.json`'ga yozildi (`saveProducts`: `json.MarshalIndent` + `os.WriteFile`) va qayta oʻqildi
(`loadProducts`: `os.ReadFile` + `json.Unmarshal(data, &products)`). `Product`'ning kalitlari tag'lar bilan berilgan.
Ikki ish chala qolgan.

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan joylar)

1. **`Order` va `OrderItem` (GB-2.03 turlari)**: hozir ularda tag yoʻq, shuning uchun JSON kalitlari Go nomlari
   bilan chiqadi (`"ID"`, `"ProductID"` …). Jamoa kelishgan format:

   | maydon | JSON kaliti |
   |---|---|
   | `Order.ID` | `id` |
   | `Order.Items` | `items` |
   | `Order.Status` | `status` (nol boʻlsa ham yozilsin: `StatusNew` = 0) |
   | `OrderItem.ProductID` | `product_id` |
   | `OrderItem.Price` | `price` |
   | `OrderItem.Quantity` | `quantity` |

   Natija shunday boʻlishi kerak:

   ```
   {"id":7,"items":[{"product_id":1,"price":1200000,"quantity":3}],"status":0}
   ```

2. **`loadProducts`**:
   - **(a)** bozor birinchi marta ishga tushganda `bozor.json` hali yoʻq. Hozir `loadProducts` bunda xato qaytaradi
     (`open bozor.json: no such file or directory`) va `main` toʻxtaydi. Kerak: fayl yoʻq boʻlsa **boʻsh katalog va
     `nil` xato**. Boshqa `os.ReadFile` xatolari (masalan, huquq yoʻq) avvalgidek qaytsin.
   - **(b)** JSON buzuq boʻlsa, xabarda qaysi fayl ekani yoʻq:
     `json: cannot unmarshal string into Go struct field .0.price of type int64`. Kerak: xato **fayl yoʻli bilan
     boshlansin** (`bozor.json: json: cannot unmarshal …`) va asl JSON xatosi zanjirda qolsin, yaʼni chaqiruvchi uni
     `errors.As` bilan topa olsin. Xato boʻlganda katalog `nil` qaytsin (hozirgidek).

`Product`'ni, `saveProducts`'ni, `run`'ni va `main`'ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. Test fayli `tekshiruv_test.go` oʻzgartirilmaydi. Hozir u shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaOrderJSON (0.00s)
    tekshiruv_test.go:67: json.Marshal(Order):
          hozir    {"ID":7,"Items":[{"ProductID":1,"Price":1200000,"Quantity":3}],"Status":0}
          kutilgan {"id":7,"items":[{"product_id":1,"price":1200000,"quantity":3}],"status":0}
        Order va OrderItem maydonlariga json tag'larini qoʻying; status'da omitempty boʻlmasin (TASK.md, 1-qadam)
…
--- FAIL: TestVazifaFaylYoq (0.00s)
    tekshiruv_test.go:93: fayl yoʻq: loadProducts xato qaytardi: open /tmp/…/bozor.json: no such file or directory; kutilgan: boʻsh katalog va nil — errors.Is(err, os.ErrNotExist) (TASK.md, 2-qadam)
--- FAIL: TestVazifaBuzuqFayl (0.00s)
    --- FAIL: TestVazifaBuzuqFayl/tur_mos_emas (0.00s)
        tekshiruv_test.go:123: xato «json: cannot unmarshal string into Go struct field .0.price of type int64»
              fayl nomi bilan boshlanishi kerak: «/tmp/…/bozor.json: …» — fmt.Errorf bilan path'ni qoʻshing (TASK.md, 2-qadam)
…
FAIL	github.com/abdurahmon27/bozor	0.002s
```

**Maslahatlar**

- 1-qadam: videodagi `Product` tag'lariga qarang: `` `json:"kalit"` ``, maydon turidan keyin, backtick ichida,
  ikki nuqtadan keyin probelsiz. `status`'ga `omitempty` qoʻyilsa, `StatusNew` (0) yozilmay qoladi — videodagi savol.
- 2(a): `os.ReadFile`'ning "fayl yoʻq" xatosini `errors.Is(err, os.ErrNotExist)` taniydi (GB-2.08'dagi `errors.Is`).
  Bu tekshiruv `if err != nil`'dan **oldin** turishi kerak.
- 2(b): GB-2.08'dagi `fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)` shaklini eslang: bu safar `%d` oʻrnida yoʻl
  (`%s`), `%w` oʻrnida `Unmarshal`'ning xatosi. `%v` bilan yozsangiz, test "asl JSON xatosi yoʻq" deydi.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
