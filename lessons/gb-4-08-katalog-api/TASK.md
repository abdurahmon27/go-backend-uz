# Vazifa · GB-4.08: `DELETE /products/{id}` — qoʻllanmasiz

Bu 4-modulning nazorat nuqtasi: yangi endpoint'ni oʻzingiz, videodagi boʻlaklardan yigʻasiz. Kerakli hamma narsa darsda
koʻrsatilgan: qulf (`Get`, `Sell`), `pathID`, `writeStoreError`, `newMux`, jadvalli test (`TestRoutes`).

**Nima qilish kerak**

1. **`internal/store/store.go` → `Delete`.** `func (s *Store) Delete(id int64) error` hozir faqat xato qaytaradi.
   Uni yozing: id'li mahsulotni katalogdan oʻchirsin, qolganlari **oʻz tartibida** qolsin. Mahsulot yoʻq boʻlsa —
   `ErrNotFound` (boshqa metodlardagidek `fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)`). `Delete` `products`ni oʻzgartiradi,
   shuning uchun boshida `s.mu.Lock()` va `defer s.mu.Unlock()` boʻlishi **shart**.
2. **`cmd/api/handlers.go` → route va handler.** `newMux`dagi `VAZIFA` izohi oʻrniga `DELETE /products/{id}` pattern'ini
   `a.deleteProduct` bilan roʻyxatdan oʻtkazing. `deleteProduct`ni yozing:
   - id son boʻlmasa — `400 bad_id` (`pathID` buni oʻzi yozadi);
   - mahsulot yoʻq boʻlsa — `404 not_found` (`writeStoreError`);
   - oʻchirildi — **`204 No Content`**: `w.WriteHeader(http.StatusNoContent)`, body **yozilmaydi**.
3. **`cmd/api/api_test.go` → oʻz testingiz.** `func TestDeleteProduct(t *testing.T)` yozing: kamida uchta holat
   (bor ID → 204, yoʻq ID → 404, `abc` → 400) va oʻchirilgan mahsulot `GET`da 404 boʻlishi. `send` va `testStore`
   shu faylda bor. (Tekshiruv faqat test borligini, `DELETE` yuborishini va `204`ni tekshirishini koʻradi; testning
   sifati — sizning ishingiz.)

`tekshiruv_test.go` fayllarini oʻzgartirmang.

**Tekshirish** (papkaning ildizida):

```bash
go test ./...
```

Hammasi `ok` boʻlguncha ishlang. Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaDeleteRoute (0.00s)
    tekshiruv_test.go:17: DELETE /products/2: 405 — route yoʻq: newMux'ga "DELETE /products/{id}" pattern'ini a.deleteProduct bilan qoʻshing (TASK.md, 2-qadam)
--- FAIL: TestVazifaOzTest (0.00s)
    tekshiruv_test.go:73: cmd/api'da func TestDeleteProduct(t *testing.T) yoʻq: DELETE /products/{id} uchun oʻz testingizni yozing (TASK.md, 3-qadam)
FAIL	github.com/abdurahmon27/bozor/cmd/api
…
--- FAIL: TestVazifaDelete (0.00s)
    tekshiruv_test.go:41: Delete(2) xato qaytardi: Delete hali yozilmagan (TASK.md, 1-qadam) — non (ID 2) katalogda bor, xato boʻlmasligi kerak
--- FAIL: TestVazifaDeleteQulf (0.00s)
    tekshiruv_test.go:73: Delete qulf band boʻlsa ham darhol qaytdi (xato: Delete hali yozilmagan (TASK.md, 1-qadam)): boshiga s.mu.Lock() va defer s.mu.Unlock() qoʻying (GB-4.07)
--- FAIL: TestVazifaDeleteParallel (0.00s)
    tekshiruv_test.go:124: 200 ta parallel Delete'dan keyin katalogda 400 ta mahsulot, kutilgan 200 ta (faqat juft ID'lar): …
FAIL	github.com/abdurahmon27/bozor/internal/store
```

Tekshiruv nimani koʻradi (`-race`siz ham ishlaydi):
- `TestVazifaDeleteQulf` qulfni oʻzi ushlab turib `Delete`ni chaqiradi: toʻgʻri `Delete` kutib turadi va katalogga tegmaydi;
  qulfsiz yoki `Lock`i oxirida boʻlgan `Delete` darhol tugaydi yoki katalogni oʻzgartirib qoʻyadi.
- `Delete`dan keyin qulf boʻsh boʻlishi kerak (`defer s.mu.Unlock()` unutilmasin).
- `Delete` 2 soniyada qaytmasa — test osilib qolmaydi, "deadlock" deb aytadi: qulf ichida `Get`, `Sell` yoki `All`ni
  chaqirmang (ular ham `Lock` oladi, GB-4.07).
- `TestVazifaDeleteParallel` 200 ta `Delete`ni `All` bilan bir vaqtda chaqiradi: qidirish va oʻchirish bitta qulf ostida
  boʻlmasa (masalan avval `s.Get(id)`, keyin qulfsiz oʻchirish), katalog buziladi.

Mashinangizda `gcc` (yoki Docker) boʻlsa, oxirida `-race` bilan ham tekshiring — topshirish sharti shu:

```bash
go test -race ./...
# gcc yoʻq, Docker bor:  docker run --rm -v "$PWD":/src -w /src golang:1.27 go test -race ./...
```

**Maslahatlar**

- Slice'dan elementni oʻchirish: indeksini toping (`for i := range s.products` yoki `slices.IndexFunc`), keyin
  `slices.Delete(s.products, i, i+1)` (`go doc slices.Delete`) — qolganlar tartibini saqlaydi.
- `deleteProduct` `getProduct`ga juda oʻxshaydi: `pathID` → store metodi → xato boʻlsa `writeStoreError` va `return`.
  Farqi oxirida: `writeJSON` emas, faqat `WriteHeader(http.StatusNoContent)`.
- Body'siz javob: `writeJSON(w, 204, nil)` emas. U `null` yozishga urinadi: haqiqiy server `Write`ni xato bilan rad etadi
  (`http: … status code does not allow body`), testdagi recorder esa `null`ni yozib qoʻyadi va tekshiruv yiqiladi.
  Faqat `w.WriteHeader(http.StatusNoContent)`.
- `TestDeleteProduct` uchun `TestRoutes` shaklini oling (204 holatini ham tekshiring):
  `[]struct{ name, path string; status int }`, `t.Run`, har qatorga `newMux(testStore(t))`.

Javob: [`solution/`](solution/). Javobni tekshirish: `cd solution && go test ./...`.
