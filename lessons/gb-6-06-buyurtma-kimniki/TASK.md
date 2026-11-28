# Vazifa · GB-6.06: `PATCH /orders/{id}` faqat oʻz buyurtmangizni oʻzgartirsin

Videoda `GET /orders/{id}` tuzatildi: `store.Order(ctx, id, userID)` soʻrovida `WHERE id = $1 AND user_id = $2`, egasi
session'dan (`currentUser(r).ID`), begona buyurtma — yoʻq buyurtma kabi `404`.

Papkada yana bir route bor: `PATCH /orders/{id}` — buyurtmaning yetkazish izohini (`note`, odatda manzil) oʻzgartiradi.
Uning `UPDATE`ida videodagi xato bor: kimning buyurtmasi ekani tekshirilmaydi.

```bash
# Jasur (ID 2) login qiladi va Dilnozaning 41-buyurtmasiga oʻz manzilini yozadi — hozir 204, izoh oʻzgaradi:
curl -s -c jasur.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"jasur@bozor.test","password":"choy-non-qand-27"}'
curl -i -X PATCH -b jasur.txt localhost:8765/orders/41 -H "Content-Type: application/json" -d '{"note":"Jasurning manzili"}'
```

Bu — OWASP API Security Top 10 2023'ning birinchi bandi (BOLA), faqat oʻqish emas, **oʻzgartirish** yoʻlida: hujumchi
boshqa birovning buyurtmasini oʻz manziliga yoʻnaltiradi.

**Nima qilish kerak:**

1. **Store** (`internal/store/orders.go`, `UpdateNote`): metod egasini ham olsin —
   `UpdateNote(ctx context.Context, id, userID int64, note string) error` — va uni SQL shartiga qoʻysin:
   `UPDATE orders SET note = $3 WHERE id = $1 AND user_id = $2`. Raqamlar argumentlar tartibiga mos:
   `$1` ← `id`, `$2` ← `userID`, `$3` ← `note`, ya'ni `Exec(ctx, sql, id, userID, note)`. Begona yoki yoʻq buyurtmada
   hech qator oʻzgarmaydi (`tag.RowsAffected() == 0`) → `ErrOrderNotFound` (`%w` bilan, hozirgidek).
2. **Handler** (`cmd/api/handlers.go`, `updateNote`): store'ga `currentUser(r).ID`ni bering (`VAZIFA` izohi turgan joy).
   Egani body'dan yoki yoʻldan olmang: uni faqat session biladi.

`cmd/api/tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (papkaning ildizida, Postgres ishlab turgan boʻlsin):

```bash
docker compose up -d --wait
go test ./...
```

Hammasi `ok` boʻlguncha ishlang. Hozir shunday yiqiladi:

```
--- FAIL: TestUpdateNoteOwner (0.20s)
    tekshiruv_test.go:28: Jasur PATCH /orders/41 (Dilnozaniki): status 204, kutilgan 404 — begona buyurtma yoʻq buyurtma kabi (TASK.md, 1–2-qadam)
    tekshiruv_test.go:31: Jasur Dilnozaning 41-buyurtmasini oʻzgartirdi: izoh «Jasurning manzili», kutilgan «Chilonzor 9-kvartal, 12-uy» — UPDATE'da egasi tekshirilmayapti (TASK.md, 1-qadam: WHERE id = $1 AND user_id = $2)
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api
```

Test uch narsani tekshiradi: Jasurning soʻrovi `404` oladi va izoh **oʻzgarmaydi**; Dilnozaning oʻz buyurtmasi `204` va
izoh oʻzgaradi; yoʻq buyurtma (`999`) ham `404`.

**Maslahatlar:**

- `store.Order` metodiga qarang: oʻsha shakl, faqat `QueryRow` oʻrnida `Exec` va `RowsAffected()` (GB-5.08).
- Imzoni oʻzgartirgandan keyin `go vet ./...` handler'dagi eski chaqiruvni koʻrsatadi (argument yetmaydi): imzo egasini
  berishga majbur qiladi. `id` va `userID` ikkalasi `int64` — ularning tartibini test tekshiradi.
- Testdagi Jasur va Dilnoza `seed.sql`dagi foydalanuvchilar (parollar README'da). Har test oʻz sxemasida ishlaydi
  (`internal/pgtest`), sizning bazangizga tegmaydi.

Javob: `solution/` papkasida (`internal/store/orders.go`, `cmd/api/handlers.go`).
