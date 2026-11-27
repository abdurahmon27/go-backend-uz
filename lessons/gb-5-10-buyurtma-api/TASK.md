# Vazifa · GB-5.10: `GET /orders/{id}`

API buyurtmani yoza oladi (`POST /orders`), lekin uni qaytarib o'qiy olmaydi: `curl -s localhost:8765/orders/4` hozir
`404 page not found` beradi (mux'ning o'z javobi, route yo'q). Shu endpoint'ni yozing. Qo'llanma yo'q: spetsifikatsiya va
testlar bor, testlar hal qiladi.

**Spetsifikatsiya**

```
GET /orders/1 → 200 OK
{"id":1,"customer":"Dilnoza","status":"yangi",
 "items":[{"product_id":1,"quantity":2,"unit_price":1200000},
          {"product_id":2,"quantity":3,"unit_price":400000}],
 "total":3600000}

GET /orders/99  → 404 {"error":{"code":"not_found",…}}
GET /orders/abc → 400 {"error":{"code":"bad_id",…}}
```

Javob `POST /orders` javobi bilan bir xil shaklda (`store.Order`): `POST` qaytargan `Location` manzilini `GET` qilsangiz,
aynan o'sha body qaytishi kerak. `items` qatorlari `order_items.id` tartibida, `total` = Σ `quantity × unit_price` (tiyin).

**Nima qilish kerak**

1. `internal/store/store.go`: `Order(ctx, id)` metodini yozing. Ikki so'rov: `orders`dan `customer, status`
   (`QueryRow` + `Scan`; qator yo'q bo'lsa `ErrOrderNotFound`, `%w` bilan o'rab), keyin `order_items`dan shu buyurtmaning
   qatorlari (`Query` + `rows.Next` + `Scan`, `ORDER BY id`) va `Total`. Buyurtmada qator bo'lmasa ham `Items` JSON'da `[]`
   bo'lsin, `null` emas (seed'dagi 3-buyurtma, Malika, qatorsiz: test buni tekshiradi).
2. `cmd/api/handlers.go`: `newMux`dagi `VAZIFA` izohi o'rniga `GET /orders/{id}` route'i va `getOrder` handler'i
   (`getProduct` kabi: `pathID` → `a.st.Order(r.Context(), id)` → `writeJSON(w, 200, o)`).
3. `writeStoreError`: `store.ErrOrderNotFound` ham `404 not_found` bo'lsin (`VAZIFA` izohi turgan joyda).

`*_test.go` fayllari — tekshiruv: ularni o'zgartirmang.

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test ./...
```

Har test o'z vaqtinchalik sxemasida ishlaydi (`internal/pgtest`): migratsiyalar va `seed.sql` o'sha sxemaga qo'llanadi,
test tugagach sxema o'chiriladi. Sizning `public` sxemangizga tegmaydi. Postgres boshqa manzilda bo'lsa `DATABASE_URL`ni
qo'ying (standart: `postgres://bozor:bozor-lab@localhost:5510/bozor`). Postgres ishlamasa test aniq xabar bilan yiqiladi,
o'tkazib yuborilmaydi.

Hozir shunday yiqiladi (qisqartirilgan, `_video/chiqishlar.txt` §6):

```
--- FAIL: TestGetOrder (0.20s)
    --- FAIL: TestGetOrder/bor (0.05s)
        tekshiruv_test.go:25: GET /orders/1: status 404, kutilgan 200 (TASK.md, 2–3-qadam): 404 page not found
…
FAIL	github.com/abdurahmon27/bozor/cmd/api	1.227s
--- FAIL: TestOrder (0.06s)
    tekshiruv_test.go:17: Order(1) = Order hali yozilmagan (TASK.md, 1-qadam); kutilgan: seed.sql'dagi Dilnoza buyurtmasi (TASK.md, 1-qadam)
…
FAIL	github.com/abdurahmon27/bozor/internal/store	0.847s
```

Vazifa tugagach `go test ./...` ikkala paket uchun `ok` beradi. Tugagach: `docker compose down -v`.

**Maslahatlar**

- `Get` (mahsulot) va `List` metodlariga qarang: birinchi so'rov `Get`ga, ikkinchisi `List`ga o'xshaydi.
  `defer rows.Close()` va oxirida `rows.Err()`ni unutmang (GB-5.05).
- `Items: []OrderItem{}` bilan boshlang, keyin `append`: bo'sh slice JSON'da `[]`, `nil` slice `null` (GB-3.06).
- `Total`ni Go'da hisoblash yetadi; xohlasangiz SQL'da `SUM(quantity * unit_price)` bilan ham bo'ladi (GB-5.04).
- 🟡 **Ixtiyoriy:** so'rovlarni `sqlc` bilan generatsiya qiling (GB-5.09). Testlar o'sha: ular `store.Order`ning natijasini
  tekshiradi, ichida qanday yozilganini emas.

Javob: [`solution/`](solution/) (`internal/store/store.go`dagi `Order`, `cmd/api/handlers.go`dagi `getOrder`).
Tekshirish: `cd solution && docker compose up -d --wait && go test ./...`
