# Vazifa · GB-4.04: bozor serverini oʻzingiz yigʻing — `newServer`

Videoda server `main` ichida qatorma-qator yozildi: `http.NewServeMux`, `mux.HandleFunc`, keyin `&http.Server{…}` va
`ListenAndServe`. Endi xuddi shuni oʻzingiz yozasiz, lekin `main`da emas, `newServer(addr string) *http.Server` funksiyasida:
test serverni `:8765`da ishga tushirmasdan, shu funksiya qaytargan `*http.Server`ni tekshira olishi uchun. `main` tayyor:
`newServer(":8765")`ni chaqiradi va `srv.ListenAndServe()` qiladi.

Hozir `newServer` boʻsh va `nil` qaytaradi, shuning uchun `go run .` darhol `xato: newServer hali yozilmagan (TASK.md)` deydi.

**Nima qilish kerak** (`main.go`, `newServer` funksiyasi; boshqa joyni oʻzgartirmang)

1. **Yoʻllar jadvali.** Yangi `ServeMux` yasang va unga uchta yozuv qoʻshing:
   - `/health` → `healthHandler`
   - `/products` → `catalogHandler`
   - `/products/1` → `productHandler` (olma). `/products/2` esa `404` boʻlib qolsin: faqat aynan `/products/1`.
     (GB-4.05'da `{id}` bilan hamma ID uchun bitta pattern yozamiz.)
2. **Server.** `*http.Server` qaytaring: `Addr` — funksiyaga kelgan `addr`, `Handler` — sizning mux'ingiz,
   `ReadHeaderTimeout` — 5 soniya. Fayl boshidagi `import (…)` ichiga `"time"`ni qoʻshing (GB-4.02).

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. `tekshiruv_test.go`ni oʻzgartirmang. U `newServer("127.0.0.1:0")`ni chaqiradi (`127.0.0.1` — shu
kompyuter, `:0` — boʻsh portni operatsion tizim tanlaydi; bu manzil faqat `Addr` maydonida oʻqiladi). Keyin qaytgan serverning
`Handler`ini `httptest.NewServer` bilan boʻsh portda ochadi (GB-4.03) va unga haqiqiy soʻrovlar yuboradi. `ReadHeaderTimeout`
maydonini esa toʻgʻridan-toʻgʻri oʻqiydi. Hozir u shunday yiqiladi:

```
--- FAIL: TestVazifaYollar (0.00s)
    tekshiruv_test.go:51: newServer nil qaytardi: mux'ni yasab, &http.Server{…} qaytaring (TASK.md, 2-qadam)
--- FAIL: TestVazifaOlma (0.00s)
    tekshiruv_test.go:63: newServer nil qaytardi: mux'ni yasab, &http.Server{…} qaytaring (TASK.md, 2-qadam)
--- FAIL: TestVazifaTimeout (0.00s)
    tekshiruv_test.go:80: newServer nil qaytardi: mux'ni yasab, &http.Server{…} qaytaring (TASK.md, 2-qadam)
FAIL
```

Test oʻtgach qoʻlda ham koʻring: `go run .`. Ikkinchi terminalda `curl http://localhost:8765/products/1` olmani,
`curl -i http://localhost:8765/products/2` esa `404 page not found`ni qaytarishi kerak. Serverni `Ctrl+C` bilan toʻxtating.

**Maslahatlar**

- 1-qadam: videodagi `mux1`–`mux4` sahnalari. Handler'ni qavssiz bering: funksiyani chaqirmaymiz, uning oʻzini beramiz.
- `/products/1`ni oxiri slesh bilan (`/products/`) yozsangiz, u `/products/` bilan boshlanadigan hamma yoʻlni oladi va
  `/products/2` ham olma qaytaradi: test buni ushlaydi.
- 2-qadam: videodagi `http.Server` struct literal'i (Addr, Handler, ReadHeaderTimeout). Faqat `5` yozilsa, bu 5 soniya emas,
  5 nanosoniya: `time.Duration` nanosoniyalarda sanaydi. Test buni ham ushlaydi.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
