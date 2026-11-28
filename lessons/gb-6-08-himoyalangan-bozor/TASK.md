# Vazifa · GB-6.08: `GET /orders` va `GET /admin/orders`

Himoyalangan bozor'da buyurtmani ID boʻyicha olish bor (`GET /orders/{id}`, faqat egasiga), lekin **roʻyxat** yoʻq:
xaridor oʻz buyurtmalarini, admin esa hammasini koʻra olmaydi. Hozir:

```bash
# login: cookie'lar dilnoza.txt va admin.txt'ga (parollar README'da)
curl -s -c dilnoza.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'
curl -s -c admin.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"admin@bozor.test","password":"qovun-anjir-tut-uzum-31"}'
curl -i -b dilnoza.txt localhost:8765/orders        # 405: bu yoʻlda faqat POST bor
curl -i -b admin.txt localhost:8765/admin/orders    # 404: admin mux'da bunday route yoʻq
```

Shu ikki route'ni yozing va ularning qoidalarini himoya jadvaliga qoʻshing. Qoʻllanma yoʻq: spetsifikatsiya va testlar
bor, testlar hal qiladi.

**Spetsifikatsiya**

```
GET /orders  (Dilnoza)      → 200 [{"id":41,"user_id":1,"status":"yangi","total":3600000}]
GET /orders  (Jasur)        → 200 [{"id":42,"user_id":2,"status":"yangi","total":6000000}]
GET /orders  (buyurtmasiz)  → 200 []
GET /orders  (cookie'siz)   → 401
GET /admin/orders (admin)   → 200 [{"id":41,…,"total":3600000},{"id":42,…,"total":6000000}]
GET /admin/orders (xaridor) → 403
GET /admin/orders (cookie'siz) → 401
```

Roʻyxatda buyurtma qatorlari (`items`) yoʻq, faqat jami: `total` = Σ `quantity × unit_price` (tiyin). Tartib — ID
boʻyicha. Buyurtma boʻlmasa `[]`, `null` emas.

**Nima qilish kerak**

1. **Store** (`internal/store/orders.go`, `VAZIFA` izohi turgan joy): `OrderSummary` turi (`ID`, `UserID`, `Status`,
   `Total`; JSON kalitlari `id`, `user_id`, `status`, `total`) va ikki metod:
   - `Orders(ctx context.Context, userID int64) ([]OrderSummary, error)` — faqat shu foydalanuvchining buyurtmalari:
     egasi SQL shartida (`WHERE o.user_id = $1`), `GET /orders/{id}`dagi `Order` metodidagidek;
   - `AllOrders(ctx context.Context) ([]OrderSummary, error)` — hamma buyurtmalar.
2. **Route'lar** (`cmd/api/handlers.go`, `newMux`dagi ikki `VAZIFA` izohi):
   - `GET /orders` — asosiy `mux`da, `a.requireUser(…)` bilan; egasi `currentUser(r).ID`dan (yoʻldan yoki body'dan emas);
   - `GET /admin/orders` — **`admin` mux'da**. Oʻram yozmang: admin mux allaqachon `requireUser` + `requireAdmin` bilan
     bir marta oʻralgan. Uni asosiy `mux`ga qoʻysangiz, test Dilnoza uchun `403` emas, `200` koʻradi.
3. **Himoya jadvali** (`cmd/api/himoya_test.go`, `accessRules`dagi `VAZIFA` izohi): yangi route'larning qatorlari —
   kamida `anonim GET /orders → 401`, `anonim GET /admin/orders → 401` va bitta xaridor (`dilnoza` yoki `jasur`)
   `GET /admin/orders → 403`. `TestAccess` ularni ham ishga tushiradi.

`cmd/api/tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (papkaning ildizida):

```bash
docker compose up -d --wait
go test ./...
```

Har test oʻz vaqtinchalik sxemasida ishlaydi (`internal/pgtest`), sizning bazangizga tegmaydi. Postgres boshqa manzilda
boʻlsa `DATABASE_URL`ni qoʻying (standart: `postgres://bozor:bozor-lab@localhost:5518/bozor`). Postgres ishlamasa test
aniq xabar bilan yiqiladi, oʻtkazib yuborilmaydi.

Hozir shunday yiqiladi (qisqartirilgan, `_video/chiqishlar.txt` §6):

```
--- FAIL: TestMyOrders (0.17s)
    tekshiruv_test.go:25: Dilnoza GET /orders: 405 {"type":"about:blank","title":"Method Not Allowed",…}
          kutilgan 200 [{"id":41,"user_id":1,"status":"yangi","total":3600000}] (TASK.md, 1–2-qadam: faqat oʻz buyurtmalari, boʻlmasa [])
…
--- FAIL: TestAdminOrders (0.14s)
    tekshiruv_test.go:39: admin GET /admin/orders: 404 {"type":"about:blank","title":"Not Found",…}
…
--- FAIL: TestAccessRulesCoverTask (0.00s)
    tekshiruv_test.go:65: accessRules'da qator yoʻq: anonim GET /orders → 401 (TASK.md, 3-qadam: himoya_test.go jadvali)
…
FAIL	github.com/abdurahmon27/bozor/cmd/api
```

Hammasi `ok` boʻlguncha ishlang. Tugagach: `docker compose down -v`.

**Maslahatlar**

- Jami uchun bitta soʻrov yetadi: `orders`ni `order_items` bilan `LEFT JOIN` (qatorsiz buyurtma ham chiqsin),
  `GROUP BY o.id` va `COALESCE(sum(i.quantity * i.unit_price), 0)` (GB-5.04). `sum(bigint)` Postgres'da `numeric`
  qaytaradi: `::bigint` bilan `int64`ga sigʻadigan turga qaytaring.
- `rows.Next` sikli GB-5.05'dagidek; `defer rows.Close()` va oxirida `rows.Err()`. Boʻsh natija `[]` boʻlishi uchun
  slice'ni `[]OrderSummary{}` bilan boshlang (`nil` slice JSON'da `null`, GB-3.06).
- `Orders` va `AllOrders` farqi faqat `WHERE`da: umumiy qismini yordamchi metodga chiqarsangiz boʻladi, lekin SQL'ni
  satrlarni ulab yasamang.
- Jadval qatori `rule{who, method, path, body, want}`: masalan `{"anonim", "GET", "/orders", "", 401}`.
- Laboratoriyada roʻyxat toʻliq qaytadi; ishlab turgan tizimda `AllOrders`ga sahifalash (`LIMIT`, keyingi sahifa kaliti) kerak boʻladi — bu vazifada talab qilinmaydi.

Javob: [`solution/`](solution/) (`internal/store/orders.go` oxiri, `cmd/api/handlers.go`dagi `listOrders` va
`allOrders`, `cmd/api/himoya_test.go`dagi uch qator). Tekshirish: `cd solution && docker compose up -d --wait && go test ./...`
