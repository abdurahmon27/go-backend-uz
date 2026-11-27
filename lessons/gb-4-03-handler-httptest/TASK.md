# Vazifa · GB-4.03: `catalogHandler`ni yozing

Videoda `productHandler` bitta mahsulotni (olma) JSON qilib berdi va `httptest.NewRecorder` bilan serversiz
tekshirildi. Endi butun katalogni beradigan handler kerak. `main.go`dagi `catalogHandler` hozir bo'sh: u hech narsa
yozmaydi, shuning uchun javob `200` va bo'sh body.

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan joy)

1. **GET.** `catalog` (`[]Product`) ni JSON massiv qilib yozing: avval `Content-Type: application/json` sarlavhasi,
   keyin `json.NewEncoder(w).Encode(...)` (videodagi `productHandler`dagidek).
2. **Boshqa metod.** `GET` bo'lmasa: `Allow: GET` sarlavhasi, keyin `w.WriteHeader(http.StatusMethodNotAllowed)`,
   keyin `return`. Body bo'lmasin. Tartibni unutmang: sarlavha `WriteHeader`dan **oldin**.
3. **Bo'sh katalog.** `catalog` bo'sh (`nil`) bo'lsa ham javob `[]` bo'lsin, `null` emas: nil slice JSON'da `null`
   bo'ladi (GB-3.06'da `[]` saqlagan edik). Bo'sh, lekin `nil` bo'lmagan slice kerak (`[]Product{}`).

`Product`, `healthHandler`, `productHandler`, `main` va `handler_test.go`ni o'zgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` bo'lguncha ishlang. `tekshiruv_test.go` o'zgartirilmaydi: u `catalogHandler`ni serversiz, `httptest.NewRecorder`
bilan chaqiradi (siz videoda qilgandek). Hozir u shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaGET (0.00s)
    tekshiruv_test.go:29: GET: Content-Type «», kutilgan «application/json»: sarlavhani Encode'dan OLDIN qoʻying
    tekshiruv_test.go:33: GET: body JSON massiv emas (EOF): catalog'ni json.NewEncoder(w).Encode bilan yozing, hozirgi body: ""
--- FAIL: TestVazifaMetod (0.00s)
    --- FAIL: TestVazifaMetod/POST (0.00s)
        tekshiruv_test.go:45: POST: status 200, kutilgan 405 (http.StatusMethodNotAllowed): WriteHeader'ni chaqiring va return qiling
…
--- FAIL: TestVazifaBosh (0.00s)
    tekshiruv_test.go:63: catalog boʻsh: body «», kutilgan «[]» (nil slice JSON'da null boʻladi: …)
FAIL
```

**Maslahatlar**

- Videodagi `productHandler`ni nusxalang va `olma` o'rniga `catalog`ni bering; `if r.Method != http.MethodGet { … }`
  bloki o'zgarmaydi.
- Sarlavhalar tartibi: `w.Header().Set(…)` → `w.WriteHeader(…)` → body. `Write`/`Encode`dan keyin qo'yilgan
  sarlavha yo'qoladi.
- 3-qadam: `products := catalog`, keyin `if products == nil { products = []Product{} }` va `Encode(products)`.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
