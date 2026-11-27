# Vazifa · GB-6.01: `requireJSON` middleware'i

Videoda uchta middleware yozildi (`requestID`, `logRequests`, `recoverPanic`) va ular `newHandler`da bitta zanjirga
ulandi. Endi toʻrtinchisini oʻzingiz yozasiz va zanjirdagi toʻgʻri joyiga qoʻyasiz.

Muammo: `POST /products` va `POST /orders` body'ni JSON deb oʻqiydi, lekin `Content-Type`ni tekshirmaydi (GB-4.08:
«bu server Content-Type'ni tekshirmaydi»). `curl -d '…'` sarlavhasiz yuborsa, curl oʻzi
`Content-Type: application/x-www-form-urlencoded` qoʻyadi, server esa baribir JSON deb oʻqiydi. Toʻgʻri javob —
**`415 Unsupported Media Type`** (GB-4.01, RFC 9110 §15.5.16).

**Nima qilish kerak**

1. **`cmd/api/middleware.go` → `requireJSON`.** Hozir u `next`ni oʻzgarishsiz qaytaradi. Uni middleware shaklida yozing
   (`return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { … })`):
   - soʻrov metodi `POST` boʻlsa va `Content-Type` **`application/json` boʻlmasa** (yoki umuman yoʻq boʻlsa) —
     `writeError(w, http.StatusUnsupportedMediaType, apiError{Code: "unsupported_media_type", Message: …})` va **`return`**;
   - `application/json; charset=utf-8` ham JSON: parametrlari bilan kelgan turni ham qabul qiling;
   - qolgan hamma soʻrov (`GET`, `DELETE`, toʻgʻri `POST`) — `next.ServeHTTP(w, r)`.

   Keyin uni ishga ulang. Joyini oʻzingiz tanlang, sharti: **415 javobi ham `X-Request-Id` olsin va server logiga
   tushsin** (videodagi tartib qoidasi). Ikki yoʻl bor:
   - `newHandler`dagi `VAZIFA` izohi turgan zanjirga (oddiyroq; lekin shunda `POST /health` yoki yoʻq yoʻlga ham
     405/404 oʻrniga 415 boradi, chunki tekshiruv route tanlanishidan oldin ishlaydi);
   - yoki `newMux`da faqat body oladigan ikki route'ni oʻrab (javobdagi yoʻl):
     `mux.Handle("POST /products", requireJSON(http.HandlerFunc(a.createProduct)))` — `mux.Handle` `http.Handler`
     oladi, `HandleFunc` esa oddiy funksiya (`go doc net/http.ServeMux.Handle`).

2. 🟡 **Ixtiyoriy: `muxErrors`.** Mux'ning oʻz javoblari hali oddiy matn: `GET /narsalar` → `404 page not found`,
   `PUT /products/1` → `Method Not Allowed` (GB-4.08). `muxErrors`ni shunday yozingki, bu ikkisi ham bizning JSON
   shaklida chiqsin (`code`: `not_found` va `method_not_allowed`), 405'dagi `Allow` sarlavhasi saqlansin, oddiy
   route'lar esa oʻzgarmasin. Yozmasangiz, bu tekshiruv `SKIP` boʻladi.

`tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (papkaning ildizida):

```bash
go test ./...
```

Hammasi `ok` boʻlguncha ishlang. Hozir shunday yiqiladi:

```
--- FAIL: TestVazifaRequireJSON (0.00s)
    tekshiruv_test.go:33: requireJSON: POST, Content-Type "" → status 204 (soʻrov next'ga oʻtib ketdi), kutilgan 415 (http.StatusUnsupportedMediaType). requireJSON hali yozilmagan boʻlishi mumkin: TASK.md, 1-qadam
--- FAIL: TestVazifaZanjir (0.00s)
    tekshiruv_test.go:57: newHandler: POST /products, Content-Type text/plain → 201, kutilgan 415: requireJSON'ni ulang — newHandler zanjiriga yoki newMux'dagi POST route'larga (TASK.md, 1-qadam)
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api	0.004s
```

Qoʻlda ham koʻring (`go run ./cmd/api`, ikkinchi terminalda):

```bash
curl -i localhost:8765/orders -d '{"product_id":1,"quantity":1}'                                       # 415
curl -i localhost:8765/orders -H 'Content-Type: application/json' -d '{"product_id":1,"quantity":1}'  # 201
```

**Maslahatlar**

- **Soʻrov** sarlavhalari: `r.Header.Get("Content-Type")` — `r.Header` maydon (qavssiz). Videodagi `w.Header()` esa
  **javob** sarlavhalari va metod (qavsli); `r.Header().Get` kompilyatsiya boʻlmaydi.
- `Content-Type` qiymatini `r.Header.Get("Content-Type")` beradi. Undan turni ajratish: `mime.ParseMediaType`
  (`go doc mime.ParseMediaType`) — `"application/json; charset=utf-8"` uchun `"application/json"` va parametrlarni
  qaytaradi; boʻsh yoki buzuq qiymatga xato beradi (xato ham — 415).
- Xatoni yozgandan keyin `return` unutilmasin: aks holda soʻrov baribir `next`ga oʻtadi, handler ham javob yozadi
  va server logida `http: superfluous response.WriteHeader call` chiqadi (GB-4.03).
- Tartib: zanjir tashqaridan ichkariga ishlaydi. `requireJSON` `requestID`dan tashqarida boʻlsa, 415 javobida ID
  boʻlmaydi; `logRequests`dan tashqarida boʻlsa, 415 logga tushmaydi.
- 🟡 uchun: `mux.Handler(r)` soʻrovni bajarmasdan uning pattern'ini qaytaradi (`go doc net/http.ServeMux.Handler`):
  pattern boʻsh boʻlsa, javobni mux'ning oʻzi (404 yoki 405) yozadi. Shunda writer'ni `statusRecorder`ga oʻxshash
  oʻz turingiz bilan oʻrang: `WriteHeader(404/405)`da bizning `writeError`ni chaqiring, mux keyin yozadigan matnni
  `Write`da tashlab yuboring.

Javob: [`solution/`](solution/). Javobni tekshirish: `cd solution && go test ./...`.
