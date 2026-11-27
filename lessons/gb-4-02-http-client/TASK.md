# Vazifa · GB-4.02: `fetchRate`ni ishonchli qiling

Videoda `fetchRate` ochiq kurs API'sidan (`open.er-api.com`) 1 USD necha so'm ekanini oldi: `http.Client{Timeout}`,
status tekshiruvi, `defer resp.Body.Close()` va `json.NewDecoder(resp.Body).Decode(&data)`. Uch holat chala qolgan.

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan joylar)

1. **Javob 200, lekin kurs yoʻq.** Videodagi uchinchi tuzoq: API valyuta kodini bilmasa, status `200` bilan
   `{"result":"error","error-type":"unsupported-code"}` qaytaradi; `fetchRate` esa jim `0, nil` qaytaradi.
   - **(a)** `result` `"success"` bo'lmasa, xato qaytaring; xabarda API'ning `error-type` qiymati bo'lsin
     (`rateResponse`ga `error-type` kalitli maydon qo'shing).
   - **(b)** `result` `"success"`, lekin so'ralgan `code` (masalan `"XXX"`) `rates`da yoʻq bo'lsa ham xato qaytaring;
     xabarda kod bo'lsin. Yo'q kalit map'da nol beradi: GB-2.02'dagi vergulli shaklni (`v, ok := m[k]`) ishlating.
2. **`*StatusError`.** Hozir `200` dan boshqa statusda `fetchRate` oddiy `fmt.Errorf` qaytaradi, chaqiruvchi
   `404` yoki `503` ekanini xabar matnidan ajratishga majbur. Kerak: `&StatusError{Code: resp.StatusCode}` qaytaring
   (`StatusError` turi `main.go`da tayyor), shunda chaqiruvchi `errors.As` bilan `Code`ni oladi (GB-2.08).
3. **Javob hajmi.** Tashqi serverning javobiga ishonib boʻlmaydi: u cheksiz katta boʻlishi mumkin (videoning
   «chuqurroq» qismi). `Decode`ni `io.LimitReader(resp.Body, 1<<20)` orqali qiling: 1 MB dan katta javob xato bersin.

`Product`, `StatusError`, `run` va `main`ni o'zgartirmang. Timeout va `defer resp.Body.Close()` joyida qolsin.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` bo'lguncha ishlang. Test fayli `tekshiruv_test.go` o'zgartirilmaydi; u haqiqiy internetga chiqmaydi, har test
o'z soxta serverini ochadi (`httptest`: buni GB-4.03'da tushuntiramiz). Hozir u shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaKodYoq (0.00s)
    tekshiruv_test.go:74: «XXX» kodi javobda yoʻq: fetchRate xato qaytarishi kerak, hozir 0 va nil — …
--- FAIL: TestVazifaResultXato (0.00s)
    tekshiruv_test.go:88: API "result":"error" dedi (status 200), fetchRate esa 0 va nil qaytardi; …
--- FAIL: TestVazifaStatusError (0.00s)
    --- FAIL: TestVazifaStatusError/404 (0.00s)
        tekshiruv_test.go:…: status 404: xato «kurs soʻrovi: javob 404 Not Found» *StatusError emas; …
--- FAIL: TestVazifaKattaJavob (0.0Xs)
    tekshiruv_test.go:…: javob 2 MB dan katta, fetchRate esa 12000.5 va nil qaytardi; …
FAIL
```

**Maslahatlar**

- 1(a): `rateResponse`ga `ErrorType string` maydoni va `json:"error-type"` tegi kerak (kalitdagi chiziqcha tegda
  bemalol yoziladi). Tekshiruv `Decode`dan keyin, `Rates`ga qarashdan oldin turadi.
- 1(b): `rate, ok := data.Rates[code]`; `!ok` bo'lsa `fmt.Errorf("kurs topilmadi: %s", code)` kabi xato.
- 2: `200` dan boshqa statusda body'ni o'qimay qaytarsangiz ham bo'ladi: `defer resp.Body.Close()` uni yopadi.
- 3: `body := io.LimitReader(resp.Body, 1<<20)` (`"io"`ni import qiling) va `Decode`ga `resp.Body` oʻrniga `body`ni bering.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
