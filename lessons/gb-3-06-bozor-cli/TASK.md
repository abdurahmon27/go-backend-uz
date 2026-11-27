# Vazifa · GB-3.06: `bozor restock` buyrugʻi va uning testi

Videoda bozor CLI v0'ni yigʻdik: `bozor add`, `list`, `order`, katalog `bozor.json`'da, `Save` atomar
(vaqtinchalik fayl + `os.Rename`), xatolar `message` orqali stderr'ga, chiqish kodi 1 yoki 2. Endi omborga mahsulot
keldi: qoldiqni oshiradigan buyruq kerak. Bu safar vazifa kattaroq va tutorial'siz: store'ga metod, uning testi,
CLI'ga buyruq va saqlash.

```
$ ./bozor restock 2 10
toʻldirildi: non, qoldiq 13
```

**Nima qilish kerak**

1. **`internal/store/store.go`: `func (s *Store) Restock(id int64, quantity int) (Product, error)`** — hozir u
   faqat xato qaytaradi. U `id`'li mahsulotning `Stock`'iga `quantity` dona qoʻshib, yangilangan mahsulotni
   qaytarsin.
   - `quantity` 1 dan 10 000 gacha boʻlsin. Aks holda katalogni oʻzgartirmasdan `ErrBadQuantity`'ni `%w` bilan
     oʻrab qaytaring (0 va manfiy son uchun matn `Sell`'dagidek: `…, hozir 0`; 10 000 dan koʻpi uchun matnni
     oʻzingiz tanlang, masalan `bir martada koʻpi bilan 10000`). 10 000 — xato yozilgan sondan himoya
     (`restock 2 1000000`).
   - Bunday ID boʻlmasa — `ErrNotFound`, `%w` bilan oʻralgan (`Sell`'dagidek).
   - `Sell`'dagi kabi mahsulotni `s.find(id)` bilan toping: u katalogdagi elementning manzilini beradi. Nusxaga
     yozsangiz (videodagi bug), test katalogdagi qoldiqni tekshiradi va yiqiladi.
2. **`internal/store/restock_test.go`: oʻzingizning `TestRestock`'ingiz** — GB-1.08'dagi jadval shaklida, kamida
   3 qator (masalan: oddiy, yoʻq ID, notoʻgʻri dona), har qator `t.Run` bilan. `store_test.go`'dagi `TestSell` va
   `katalog()` yordamchisi namuna (fayl `package store`'da boʻlsa, `katalog()` va `find` koʻrinadi).
3. **`cmd/bozor/main.go`: `bozor restock ID DONA`** buyrugʻi.
   - `run`'dagi `switch`'ga `case "restock":` qoʻshing va `msg, err = restock(st, args[1:])` deb chaqiring; `run`
     avval katalogni saqlaydi, keyin `msg`'ni chiqaradi (`order` kabi).
   - `func restock(st *store.Store, args []string) (string, error)` — `order` bilan bir xil tuzilish: aynan 2 ta
     argument (aks holda `errUsage`), ID — `strconv.ParseInt(args[0], 10, 64)` (oʻnlik, 64 bitli son,
     GB-3.01), dona — `strconv.Atoi`, xatolar `"ID: %w"` / `"dona: %w"`, keyin `st.Restock`.
   - Natija qatori: `toʻldirildi: <nom>, qoldiq <yangi qoldiq>` (`ʻ` — U+02BB, `oʻ` harfidagi belgi; oddiy `'`
     emas, shu yerdan nusxalang).
   - `errUsage` matniga `| restock ID DONA` qoʻshing, shunda `./bozor` foydalanish qatorida yangi buyruq koʻrinadi.

`Sell`, `Load`, `Save`, `message` va boshqa buyruqlarni oʻzgartirmang (`message` `ErrNotFound`'ni «bunday mahsulot
yoʻq», `ErrBadQuantity`'ni «dona soni notoʻgʻri» deb allaqachon chiqaradi).

**Tekshirish** (shu papkada):

```bash
go test ./...
```

`ok` chiqquncha ishlang. Tekshiruv fayllari `internal/store/tekshiruv_test.go` (1- va 2-qadam) va
`cmd/bozor/tekshiruv_test.go` (3-qadam): ularni oʻzgartirmang. Hozir ular shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaRestockBuyrugi (0.01s)
    tekshiruv_test.go:17: bozor restock 2 10: xato «foydalanish: bozor add NOM NARX QOLDIQ | list | order ID DONA», kutilgan chiqish «toʻldirildi: non, qoldiq 13» (TASK.md, 3-qadam)
…
--- FAIL: TestVazifaRestock (0.00s)
    --- FAIL: TestVazifaRestock/oddiy (0.00s)
        tekshiruv_test.go:46: Restock(2, 10) xatosi «Restock hali yozilmagan (TASK.md, 1-qadam)», kutilgan: xatosiz (TASK.md, 1-qadam)
…
--- FAIL: TestVazifaOzTestingiz (0.00s)
    tekshiruv_test.go:63: internal/store/restock_test.go oʻqilmadi (open restock_test.go: no such file or directory): Restock uchun jadvalli test yozing (TASK.md, 2-qadam)
FAIL
```

Qoʻlda ham tekshiring — videodagi toʻrt mahsulot qoʻshilgan `bozor.json` bilan (`add olma 12000 40`,
`add non 4000 3`, `add choy 25000 8`, `add oʻrik 35000 2`):

```bash
go build ./cmd/bozor
./bozor restock 2 10              # toʻldirildi: non, qoldiq 13
./bozor restock 2 0; echo $?      # xato: dona soni notoʻgʻri → 1
./bozor restock 2; echo $?        # foydalanish: … | restock ID DONA → 2
```

**Maslahatlar**

- `Restock` — `Sell`'ning teskarisi: avval `quantity`'ni tekshiring, keyin `s.find(id)`, keyin `p.Stock += quantity`
  va `return *p, nil`. Chegara uchun konstanta qulay: `const maxRestock = 10_000`.
- Avval testni yozing (2-qadam), keyin metodni: `go test ./internal/store -run TestRestock` (GB-1.08) sizning
  jadvalingizni ishga tushiradi.

Javob: [`solution/`](solution) (`solution/internal/store/store.go`, `solution/internal/store/restock_test.go`,
`solution/cmd/bozor/main.go`).
