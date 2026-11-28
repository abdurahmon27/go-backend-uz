# Vazifa · GB-7.04: `fetchPrices`ni `errgroup` bilan yozing

Videoda `importStocks` mingta mahsulot qoldigʻini yetkazib beruvchidan bir vaqtda koʻpi bilan oʻntadan soʻradi va birinchi
xatoda qolganini bekor qildi (`errgroup.WithContext` + `SetLimit`). `main.go`da yana bir funksiya bor: **`fetchPrices`**
mahsulotlar narxini (tiyinda) `checkPrice` bilan soʻraydi, lekin videodagi ikkinchi variant kabi **har ID uchun alohida
goroutine** ochadi: limit yoʻq, xatoda qolgan soʻrovlar toʻxtamaydi, xato qaysi ID'dan kelgani yozilmaydi.

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan joy): `fetchPrices(ctx, ids, limit)`ni qayta yozing.

1. **Limit.** Yetkazib beruvchiga bir vaqtda koʻpi bilan `limit` ta soʻrov borsin (`g.SetLimit(limit)`).
2. **Tartib.** Natija `ids` tartibida: `prices[i]` — `ids[i]`ning narxi.
3. **Birinchi xato.** Xato boʻlsa, funksiya `nil` va `narx <ID>: <xato>` koʻrinishidagi xatoni qaytarsin
   (`fmt.Errorf("narx %d: %w", id, err)`, shunda `errors.Is(err, ErrNotFound)` ishlaydi) va **qolgan soʻrovlar bekor
   qilinsin**: `errgroup.WithContext(ctx)` bergan context'ni `g.Go` ichida `checkPrice`ga bering.
4. **Notoʻgʻri limit.** `limit` 1 dan kichik boʻlsa, hech narsa soʻramasdan darhol `errBadLimit` qaytaring. Sababi: `SetLimit(0)`
   bilan `g.Go` hech qachon qaytmaydi (videoning "Chuqurroq" qismi: oddiy dasturda runtime `deadlock` deb toʻxtatadi, serverda
   esa soʻrov osilib qoladi), manfiy limit esa limitni umuman olib tashlaydi.

`importStocks`, `supplier.go`, `import_test.go` va `tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. Hozir shunday yiqiladi:

```
--- FAIL: TestVazifaLimit (0.01s)
    tekshiruv_test.go:45: limit 5 edi, lekin yetkazib beruvchiga bir vaqtda 200 ta soʻrov ketdi.
        	Har ID uchun alohida goroutine ochilyapti: errgroup'ning SetLimit(limit)'i kerak.
--- FAIL: TestVazifaBirinchiXato (0.01s)
    tekshiruv_test.go:74: xatodan keyin ham soʻrovlar ketaverdi: jami 200 ta (200 tadan), qolganlari bekor qilinmadi.
        	errgroup.WithContext(ctx) bergan ctx'ni g.Go ichida checkPrice'ga bering.
--- FAIL: TestVazifaNolLimit (0.01s)
    tekshiruv_test.go:96: fetchPrices(…, 0): errBadLimit kutilgan, keldi <nil>
FAIL
FAIL	github.com/abdurahmon27/bozor	0.224s
FAIL
```

Mashinangizda C kompilyator (`gcc`) boʻlsa, `go test -race .` ham toza chiqishi kerak (har goroutine faqat oʻz `prices[i]`
joyiga yozadi). `gcc` yoʻq boʻlsa, `-race` ishlamaydi (`go: -race requires cgo`): bu normal, `go test .` yetadi.

**Maslahatlar**

- Videodagi `importStocks` (shu `main.go`da) deyarli tayyor andoza: `stocks` → `prices`, `checkStock` → `checkPrice`,
  `workers` → `limit`, xato matni `"narx %d: %w"`.
- `g, gctx := errgroup.WithContext(ctx)`: `g.Go` ichida `gctx`. `checkPrice(ctx, id)` deb yozsangiz, xatodan keyin ham
  soʻrovlar ketaveradi va `TestVazifaBirinchiXato` buni aytadi.
- `TestVazifaNolLimit` 3 soniyada "tugamadi" desa, `limit < 1` tekshiruvi `SetLimit`dan **oldin** turishi kerak.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
