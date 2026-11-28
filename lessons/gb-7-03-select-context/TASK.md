# Vazifa · GB-7.03: `sumStock` va `GET /report`'ning muddati

Videoda `checkStock(ctx, id)` yetkazib beruvchining javobini yoki `ctx.Done()`ni kutdi, `stockHandler` esa unga `r.Context()`ni berdi:
mijoz ketsa, ish toʻxtaydi. Endi bitta emas, bir nechta mahsulot. `main.go`da ikki joy `VAZIFA` deb belgilangan.

**Nima qilish kerak** (`main.go`)

1. **`sumStock(ctx, ids)`ni yozing.** `ids`dagi har mahsulot qoldigʻini `checkStock(ctx, id)` bilan **ketma-ket** soʻrang va
   yigʻindisini qaytaring (`total, nil`). `checkStock` xato qaytarsa, darhol `0` va **oʻsha xatoni** (oʻzgartirmay) qaytaring.
   `ctx`ni pastga uzating: ichkarida `context.Background()` ishlatmang, aks holda muddat ham, mijozning ketishi ham ishlamaydi.
2. **`reportHandler`ga muddat qoʻying.** Hozir u `sumStock(context.Background(), reportIDs)` chaqiradi. Uni shunday qiling:
   - `ctx, cancel := context.WithTimeout(r.Context(), reportTimeout)` va keyingi qatorda `defer cancel()`;
   - `sumStock(ctx, reportIDs)`;
   - xato `context.DeadlineExceeded` boʻlsa (`errors.Is`), `http.Error(w, "yetkazib beruvchi javob bermadi", http.StatusGatewayTimeout)`
     va `return` — 504: server oʻzi kutgan tashqi tizim (yetkazib beruvchi) vaqtida javob bermadi;
   - boshqa xato boʻlsa (mijoz ketdi, `context.Canceled`), hech narsa yozmay `return`;
   - aks holda `jami: N` (bu qator oʻzgarmaydi).

`supplier.go`, `checkStock`, `stockHandler`, `newMux`, `main` va `tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. Tekshiruvda yetkazib beruvchi 3 soniya emas, 100 millisekundda javob beradi (`TestMain`), shuning uchun
testlar taxminan bir soniya ishlaydi. Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifa1SumStock (0.00s)
    tekshiruv_test.go:28: sumStock(ctx, [1 2 3]) xato qaytardi: VAZIFA 1: sumStock hali yozilmagan
        	Har id uchun checkStock(ctx, id) chaqirib, qoldiqlarni qoʻshing (VAZIFA 1).
--- FAIL: TestVazifa1Muddat (0.00s)
    tekshiruv_test.go:46: muddati 150ms boʻlgan ctx bilan sumStock xatosi «VAZIFA 1: sumStock hali yozilmagan», kutilgan context.DeadlineExceeded (got = 0).
…
--- FAIL: TestVazifa2Muddat (0.00s)
    tekshiruv_test.go:95: reportTimeout 150ms, yetkazib beruvchi 3 × 100ms: status 500, kutilgan 504 Gateway Timeout.
…
FAIL
```

**Maslahatlar**

- `sumStock` — oddiy `for _, id := range ids` sikli. Alohida `select` shart emas: `checkStock` oʻzi `ctx.Done()`ni kutadi.
  Bekor qilingan ctx bilan u darhol `ctx.Err()` qaytaradi.
- Handler'dagi muddatning otasi `r.Context()`: shunda ctx `reportTimeout` tugaganda **yoki** mijoz ketganda bekor boʻladi
  (videoda: `WithTimeout`ning otasi `r.Context()`).
- Muddat bilan mijozning ketishini `errors.Is(err, context.DeadlineExceeded)` ajratadi (videoda: `Err` qiymatlari
  `context.DeadlineExceeded` va `context.Canceled`).

Qoʻshimcha (tekshirilmaydi): serverni ishga tushiring (`go run .`) va `curl localhost:8093/report` qiling. Uch mahsulot
ketma-ket 3 × 3 = 9 soniya, muddat esa 5 soniya: 504 qaytadi. Bir vaqtda soʻrash — keyingi dars (GB-7.04).
