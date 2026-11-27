# Vazifa · GB-4.06: `POST /orders` — `validateOrder` va `createOrder`

Videoda `POST /products` body'ni `decodeJSON` bilan qat'iy o'qidi (limit `1024` bayt, nomaʼlum kalit yoʻq), qiymatlarni
`validate` bilan tekshirdi va hamma xatoni `writeError` bilan bir xil JSON shaklida qaytardi. Endi buyurtma uchun xuddi
shunday endpoint kerak. `main.go`da `POST /orders` allaqachon `newMux()`ga ulangan, lekin `validateOrder` hech narsani
tekshirmaydi, `createOrder` esa hech narsa yozmaydi (javob `200` va bo'sh body).

Mijoz yuboradigan body (`orderInput`): `{"product_id": 1, "quantity": 3}`.

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan ikki joy)

1. **`validateOrder`.** `product_id` 0 dan katta, `quantity` 1 dan 100 gacha bo'lsin. Noto'g'ri maydon uchun jadvalga
   kalit (maydonning JSON'dagi nomi: `"product_id"`, `"quantity"`) va sabab matnini qo'ying. Birinchi xatoda to'xtamang:
   hamma noto'g'ri maydon jadvalda bo'lsin. Hammasi joyida bo'lsa, jadval bo'sh qaytsin.
2. **`createOrder`.** Tartib:
   1. body'ni `decodeJSON(w, r, &in)` bilan o'qing; `false` qaytsa darhol `return` (javob yozib bo'lingan);
   2. `validateOrder(in)`: jadval bo'sh bo'lmasa `422`, `code: "validation_failed"`, `Fields` — shu jadval; keyin `return`;
   3. `catalog`dan `ID == in.ProductID` mahsulotni toping; topilmasa `404`, `code: "not_found"`;
   4. topilsa `201` va `Order{ProductID, Quantity, Total}`; `Total` — narx × dona (tiyinda, `int64`).
   Hamma xato `writeError` orqali, muvaffaqiyat `writeJSON` orqali.

`Product`, `decodeJSON`, `writeJSON`, `writeError`, `createProduct`, `newMux`, `send`, `main` va test fayllarini
o'zgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` bo'lguncha ishlang. `tekshiruv_test.go` o'zgartirilmaydi: u `POST /orders`ni serversiz, `newMux()` orqali chaqiradi
(videodagidek). `handler_test.go` videodagi `POST /products` testlari: starterda ham o'tadi. Hozir vazifa testi shunday
yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaBuyurtma (0.00s)
    tekshiruv_test.go:23: buyurtma: status 200, kutilgan 201 (http.StatusCreated): decodeJSON, tekshiruv va writeJSON(w, http.StatusCreated, Order{…}) yozing; body: ""
--- FAIL: TestVazifaTekshiruv (0.00s)
    tekshiruv_test.go:55: ikkalasi notoʻgʻri: status 200, kutilgan 422 (http.StatusUnprocessableEntity): qiymatlar notoʻgʻri, JSON esa toʻgʻri
…
FAIL
```

**Maslahatlar**

- `createProduct`ni nusxalab oling: `decodeJSON` va `validate` qismi deyarli o'zgarmaydi, faqat `productInput`
  o'rniga `orderInput`, `validate` o'rniga `validateOrder`.
- `return` unutilmasin: test body'da ikkinchi JSON bo'lsa (xatodan keyin ham yozilgan bo'lsa) yiqiladi.
- Jami narx: `p.Price * int64(in.Quantity)` (`Quantity` — `int`, `Price` — `int64`).
- Mahsulotni qidirish: `for _, p := range catalog { … }`; topilganda javobni yozib `return` qiling, sikldan keyin 404.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.

**Ixtiyoriy (🟡):** `quantity` ustiga `stock` tekshiruvini qo'shing: qoldiq yetmasa `409` (`http.StatusConflict`),
`code: "insufficient_stock"`. Buning uchun o'zingiz `tekshiruv_test.go`ga o'xshash test yozing; qoldiqni kamaytirish
(bir vaqtda ikki buyurtma!) GB-4.07 va GB-5.07 mavzusi.
