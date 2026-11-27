# Vazifa · GB-6.02: qolgan xatolarni `problemFor`ga qoʻshing

Videoda `problemFor` (`cmd/api/problem.go`) `404`, `409` va body/yoʻl xatolarini (`400`) tanidi, qolgan hammasi
`default`ga tushib **xavfsiz `500`** boʻldi. Hozir uch holat notoʻgʻri statusda qaytadi:

| Soʻrov | Hozir | Kerak |
|---|---|---|
| `POST /products` `{"name":"","price":0}` (validate → `fieldErrors`) | `500` | `422` + `errors` |
| `POST /products` 21 belgili nom (store → `ErrBadProduct`) | `500` | `422` + `detail: productRules` |
| `POST /products` 2000 baytli body (`*http.MaxBytesError`) | `400` | `413` |

**Nima qilish kerak** (`problemFor`dagi `VAZIFA` izohi oʻrniga):

1. **`fieldErrors` → `422`.** `errors.As(err, &fe)` (`var fe fieldErrors`) case'i. Javob: `Type: "/problems/validation"`,
   `Title: "Qiymatlar notoʻgʻri"`, `Status: http.StatusUnprocessableEntity`, `Errors: fe` (RFC 9457 kengaytmasi:
   maydon → nima notoʻgʻri; `detail` kerak emas).
2. **Store qoidalari → `422`.** `store.ErrBadProduct` **yoki** `store.ErrBadQuantity` (`errors.Is`, bitta `case`da vergul
   bilan, GB-1.05). Oʻsha `type` va `title`, `Detail: productRules`. Store xatosining **matnini** (`err.Error()`) yozmang:
   u ichki qatlam matni, ertaga unga boshqa tafsilot qoʻshilishi mumkin.
3. **`*http.MaxBytesError` → `413`.** `errors.As(err, &tooBig)` (`var tooBig *http.MaxBytesError`), javob
   `newProblem(http.StatusRequestEntityTooLarge, fmt.Sprintf("body %d baytdan oshmasin", tooBig.Limit))`.

`tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (papkaning ildizida):

```bash
go test ./...
```

Hammasi `ok` boʻlguncha ishlang. Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaFieldErrors (0.00s)
    tekshiruv_test.go:38: POST /products {"name":"","price":0}: status 500, kutilgan 422. fieldErrors hozir problemFor'da default'ga tushyapti: errors.As(err, &fe) case'ini qoʻshing (TASK.md, 1-qadam). Body: {"type":"about:blank","title":"Internal Server Error","status":500,"request_id":"…"}
--- FAIL: TestVazifaStoreRules (0.00s)
    tekshiruv_test.go:52: POST /products (21 belgili nom): status 500, kutilgan 422. store.ErrBadProduct va store.ErrBadQuantity uchun case qoʻshing (TASK.md, 2-qadam). Body: …
--- FAIL: TestVazifaTooLarge (0.00s)
    tekshiruv_test.go:72: POST /products (2000 baytli body): status 400, kutilgan 413. *http.MaxBytesError uchun case qoʻshing va uni errBadJSON case'idan OLDIN qoʻying: switch birinchi mos case'ni oladi (TASK.md, 3-qadam). Body: …
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api
```

**Maslahatlar**

- `switch {` birinchi **rost** `case`ni oladi (GB-1.05). `decodeJSON` `MaxBytesError`ni ham `errBadJSON` bilan oʻraydi
  (`fmt.Errorf("%w: %w", errBadJSON, err)`), shuning uchun `errors.Is(err, errBadJSON)` uning uchun ham rost. `413` case'i
  `errBadJSON` case'idan **keyin** tursa, katta body `400` boʻlib qoladi.
- `errors.As`ning ikkinchi argumenti — **koʻrsatkich** (`&fe`, `&tooBig`), GB-2.08'dagidek.
- `fieldErrors`ning `Error()` metodi qiymat receiver'ida, shuning uchun nishon `var fe fieldErrors` (koʻrsatkich emas):
  `*fieldErrors` nishoni hech qachon mos kelmaydi va xato `500`da qoladi.
- `problemFor`dagi `errBadID` case'i — `pathID` uchun (yoʻldagi `{id}` son emas), `errBadJSON` kabi API sentinel'i; unga tegmang.
- 413'ning `title`i Go'ning `http.StatusText(413)` — `"Request Entity Too Large"` (eski nom; RFC 9110'da "Content Too Large").
- `go test ./cmd/api -run TestVazifa -v` faqat vazifa testlarini koʻrsatadi.

Javob: `solution/cmd/api/problem.go` (avval oʻzingiz urinib koʻring).
