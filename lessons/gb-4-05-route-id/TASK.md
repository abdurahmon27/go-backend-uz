# Vazifa · GB-4.05: `GET /products/{id}/stock`

Videoda `GET /products/{id}` route'i bitta mahsulotni qaytardi: `r.PathValue("id")` → `strconv.ParseInt` → 400,
`findProduct` → 404, keyin JSON. Ombor xodimiga butun mahsulot emas, faqat qoldiq (`Stock`) kerak. `main.go`dagi
`stockHandler` hozir boʻsh va `newMux`da uning route'i yoʻq, shuning uchun bu yoʻlga mux oʻzi `404 page not found` beradi.

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan ikki joy)

1. **Route.** `newMux`ga `GET /products/{id}/stock` pattern'ini `stockHandler` bilan qoʻshing.
2. **Handler.** `stockHandler` id'ni yoʻldan oʻqisin (videodagi `productHandler`dagidek):
   - id son boʻlmasa (`/products/abc/stock`): `http.Error(w, "id butun son boʻlishi kerak", http.StatusBadRequest)` va `return`;
   - `findProduct` `false` qaytarsa (`/products/99/stock`): `http.Error(w, "mahsulot topilmadi", http.StatusNotFound)` va `return`;
   - topilsa: qoldiqni bitta qator matn qilib yozing: `fmt.Fprintln(w, p.Stock)` → olma uchun `40`.
3. **Metod.** `POST /products/1/stock` kelsa, javobni handler emas, mux bersin: `405` va `Allow`. Buning uchun
   pattern metod bilan boshlanishi kerak (handler ichida `r.Method`ni qoʻlda tekshirmang).

`Product`, `catalog`, `findProduct`, boshqa handler'lar va `handler_test.go`ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. `tekshiruv_test.go` oʻzgartirilmaydi: u soʻrovlarni serversiz, `newMux().ServeHTTP` orqali
yuboradi (videodagi tuzatilgan testdagidek). Hozir u shunday yiqiladi:

```
--- FAIL: TestVazifaStock (0.00s)
    tekshiruv_test.go:28: GET /products/1/stock: mux «404 page not found» qaytardi: newMux'da "GET /products/{id}/stock" route'ini stockHandler bilan roʻyxatdan oʻtkazing (TASK.md, 1-qadam)
--- FAIL: TestVazifaStockXato (0.00s)
    tekshiruv_test.go:43: GET /products/abc/stock: status 404, kutilgan 400: id son boʻlmasa http.Error(…, http.StatusBadRequest) va return (TASK.md, 2-qadam); body: "404 page not found\n"
    tekshiruv_test.go:49: GET /products/99/stock: 404 "404 page not found\n", kutilgan 404 va «mahsulot topilmadi»: findProduct nil qaytarsa http.Error(…, http.StatusNotFound) (TASK.md, 2-qadam)
--- FAIL: TestVazifaStockMetod (0.00s)
    tekshiruv_test.go:58: POST /products/1/stock: status 404, kutilgan 405: pattern metod bilan boshlansin, "GET /products/{id}/stock" (TASK.md, 3-qadam)
FAIL
```

**Maslahatlar**

- `{id}` bitta segment oladi, shuning uchun `/products/{id}/stock` va `/products/{id}` bir-biriga xalaqit bermaydi:
  birinchisi uch segmentli, ikkinchisi ikki segmentli yoʻllarni oladi.
- `productHandler`ning birinchi ikki blokini (`ParseInt` + 400, `findProduct` + 404) nusxalang; oxirgi ikki qator
  (`Content-Type` va `Encode`) oʻrniga bitta `fmt.Fprintln(w, p.Stock)`.
- Serverda ham tekshirib koʻring: `go run .`, ikkinchi terminalda `curl -i http://localhost:8765/products/1/stock`.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
