# GB-4.06 · JSON kirish va chiqish: limit, validatsiya va xato javobi

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #30"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-4-05` (route'lar, `GET /products/{id}`, 404 va 405; videodagi eslatma savoli shunga tayanadi, bu papkada `GET /products/{id}` route'i yoʻq) ·
keyingi dars: `gb-4-07` (bir vaqtda soʻrovlar: `sync.Mutex` va `-race`). Video: YouTube'dagi "Golang backend #30" (kanaldagi "Go backend" playlist'i). Bu darsga kerak boʻlganlar: [`gb-4-03-handler-httptest`](../gb-4-03-handler-httptest)
(handler, `httptest`, javob tartibi: sarlavha → status → body), [`gb-3-02-json-saqlash`](../gb-3-02-json-saqlash)
(`Product`, tag'lar, `omitempty`), [`gb-2-08-xatolar`](../gb-2-08-xatolar) (`errors.As`).

**Bu darsdan keyin:** `POST` body'ini xavfsiz oʻqiysiz: hajm limiti (`http.MaxBytesReader`, `413`), ortiqcha kalitni rad
etish (`DisallowUnknownFields`, `400`), qiymatlarni tekshirish (`422` va `fields`), va hamma xatoni bitta JSON shaklida
(`{"error": {"code", "message", "fields"}}`) qaytarasiz; `writeError`'dan keyin `return` unutilsa nima boʻlishini
haqiqiy natijada koʻrdingiz.

Bu papka mustaqil: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i, tashqi modul yoʻq.
Bu Module 4 qoidasi: har dars mustaqil lab (oldingi darsning `solution/`i ustiga qurilmaydi); GB-4.08 nazorat nuqtasi hammasini birlashtiradi.
`Product` GB-3.02'dan; katalog xotirada (`catalog`), `addProduct` GB-3.06'dagi `Add`ning oddiy nusxasi (u bir vaqtda ikki
soʻrov uchun xavfsiz emas: GB-4.07).

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `apiError`, `writeJSON`, `writeError`, `decodeJSON` (limit + qat'iy oʻqish + bitta qiymat), `validate`, `createProduct`, `newMux`, `send`, `main` (u yerga oʻz soʻrovlaringizni yozib koʻring); `validateOrder` va `createOrder` vazifa uchun boʻsh |
| `handler_test.go` | `POST /products` testlari (`201`, `400`, `413`, `422`, `return` unutilmaganligi): starterda ham oʻtadi |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayllari) |
| `_video/NN-…/` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 02-qatiy go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Har bir qadamda `main` soʻrovlarni serversiz, `newMux()` orqali yuboradi va `status body` chiqaradi (`send`, GB-4.03'dagi
`httptest.NewRequest`/`NewRecorder`).

**1. Oddiy `Decode`** (`01-oddiy`, `01b-rang`): toʻgʻri body — `201`; buzuq JSON — `400`; `productInput`'da yoʻq kalit
(`rang`) — jimgina tashlanadi, javob `201` (`stock` nol):

```
$ go run .
201 {"id":4,"name":"choy","price":2500000,"stock":30}
400 {"error":{"code":"bad_json","message":"unexpected EOF"}}
$ go run .          # 01b-rang: {"name":"choy","price":2500000,"rang":"qora"}
201 {"id":4,"name":"choy","price":2500000,"stock":0}
```

**2. `DisallowUnknownFields`** (`02-qatiy`):

```
$ go run .
400 {"error":{"code":"bad_json","message":"json: unknown field \"rang\""}}
```

**3. `http.MaxBytesReader`** (`03-limit`: 2000 baytli body, limit 1024; hali hamma xato `400`) va **`*http.MaxBytesError`
→ `413`** (`04-413`; ikkinchi qator — buzuq JSON, hamon `400`):

```
$ go run .          # 03-limit
400 {"error":{"code":"bad_json","message":"http: request body too large"}}
$ go run .          # 04-413
413 {"error":{"code":"too_large","message":"body 1024 baytdan oshmasin"}}
400 {"error":{"code":"bad_json","message":"unexpected EOF"}}
```

**4. `decodeJSON` funksiyasi** (`05-ajrat`): uchala yoʻl:

```
$ go run .
201 {"id":4,"name":"choy","price":2500000,"stock":30}
400 {"error":{"code":"bad_json","message":"json: unknown field \"rang\""}}
413 {"error":{"code":"too_large","message":"body 1024 baytdan oshmasin"}}
```

**4b. Bitta qiymat** (`11-tail`; `decodeJSON`ga ikkinchi `Decode` qoʻshilgan, bu `main.go`dagi oxirgi shakl):

```
$ go run .          # {"name":"non","price":400000}{"x":1}  va  {…} + 2000 boʻsh joy
400 {"error":{"code":"bad_json","message":"json: unknown field \"x\""}}
413 {"error":{"code":"too_large","message":"body 1024 baytdan oshmasin"}}
```

(Birinchi javobdagi xabar ikkinchi qiymatni `struct{}`ga oʻqishga urinishdan chiqqan; status esa toʻgʻri.)

**5. Validatsiya, `422`** (`06-tekshiruv`: `{"name":"","price":-5}`; `07-stock`: `{"name":"choy","price":2500000,"stock":-1}`):

```
$ go run .          # 06-tekshiruv
422 {"error":{"code":"validation_failed","message":"qiymatlar notoʻgʻri","fields":{"name":"boʻsh boʻlmasin","price":"0 dan katta boʻlsin"}}}
$ go run .          # 07-stock
422 {"error":{"code":"validation_failed","message":"qiymatlar notoʻgʻri","fields":{"stock":"manfiy boʻlmasin"}}}
```

**6. Eng koʻp uchraydigan xato — `writeError`'dan keyin `return` yoʻq** (`08-return`). Status `422` (birinchi
`WriteHeader` hisobga olinadi, GB-4.03), lekin kod davom etdi: body'da ikkinchi JSON va katalogga notoʻgʻri mahsulot
qoʻshildi (3 → 4):

```
$ go run .
422 {"error":{"code":"validation_failed","message":"qiymatlar notoʻgʻri","fields":{"name":"boʻsh boʻlmasin","price":"0 dan katta boʻlsin"}}}
{"id":4,"name":"","price":-5,"stock":0}
4 mahsulot
```

**7. Chegara: mux'ning oʻz javoblari matn** (`09-mux`; xato shakli faqat bizning handler'larimiz uchun):

```
$ go run .
404 404 page not found
405 Method Not Allowed
```

## Chuqurroq (videoda yoʻq)

`10-qoshimcha` (ikkala soʻrov `decodeJSON`ning **limit + qat'iy oʻqish** qismi bilan, bitta-qiymat tekshiruvisiz; `11-tail`dan oldingi holat):

```
$ go run .
201 {"id":4,"name":"choy","price":2500000,"stock":0}
201 {"id":5,"name":"non","price":400000,"stock":0}
5 mahsulot
```

- Birinchi soʻrovda kalitlar `NAME` va `PRICE` (katta harf): `encoding/json` kalit nomini harf katta-kichikligiga qaramay
  moslaydi. Go 1.27'da `encoding/json` yangi `encoding/json/v2` ustiga qurilgan; v2'ning oʻz API'sida nom harfi farqlanadi
  (GB-3.02 Chuqurroq). Xato matnlari versiyalar orasida oʻzgarishi mumkin: bu papkadagi natijalar go1.27.4'dan.
- Ikkinchi soʻrov `{"name":"non","price":400000}{"x":1}`: bitta-qiymat tekshiruvisiz (videoda `11-tail`dan oldin)
  `Decode` faqat birinchi qiymatni oʻqiydi, qolganini eʼtiborsiz qoldiradi.
- `fields` boʻsh qiymatni (`stock` berilmagan) nol bilan farqlamaydi: bu uchun `*int` kabi koʻrsatkich maydon kerak.
  Bu darsda `stock: 0` — toʻgʻri qiymat.
- `400`dagi `message` — `Decode` xatosining matni. U baʼzan Go nomlarini koʻrsatadi, masalan
  `json: cannot unmarshal number into Go struct field productInput.name of type string`. Ishlab chiqarishda har xato sinfi
  (`SyntaxError`, `UnmarshalTypeError`, …) uchun qat'iy matn qaytaring; bu darsda soddalik uchun xatoning oʻz matni.
- Boshqa simplifikatsiyalar: oʻqish xatosining hammasi (mijoz uzilishi ham) `400`; `Content-Type` tekshirilmaydi (real
  API `415` qaytaradi); `201` javobida `Location` sarlavhasi yoʻq (RFC 9110 §15.3.2 yaratilgan resursni `Location` yoki
  soʻrov URI bilan koʻrsatadi); `addProduct`/`catalog` bir vaqtda xavfsiz emas (GB-4.07, `go test -race`).
- Bu lab qoidalari GB-3.06'ning `internal/store` qoidalaridan **qisqaroq**: nom uzunligi (20 belgi) va UTF-8 tekshiruvi
  yoʻq, buyurtma `quantity` 1..100 (CLI'dagi `Sell` faqat `quantity > 0` talab qiladi). GB-4.08'da API `internal/store`ga
  ulanganda qoidalar bitta joyga yigʻiladi.
- Server logida `ret-run` (`return` unutilgan) holati yana bir satr beradi:
  `http: superfluous response.WriteHeader call` (haqiqiy `net/http` serveri; `httptest.NewRecorder` log yozmaydi).

## Nima uchun shunday

- **Mijoz yuborgan narsaga ishonmang.** `http.MaxBytesReader` oʻrab olingan body'ni belgilangan baytdan keyin
  `*http.MaxBytesError` bilan toʻxtatadi va (iloji boʻlsa) ulanishni yopadi; limit `Decode`dan oldin qoʻyiladi.
- **`400` va `422`.** [RFC 9110 §15.5.1](https://www.rfc-editor.org/rfc/rfc9110#name-400-bad-request): `400` — server soʻrovni
  mijoz xatosi deb hisoblab qayta ishlamaydi; [§15.5.21](https://www.rfc-editor.org/rfc/rfc9110#name-422-unprocessable-content): `422` —
  soʻrov tushunarli, lekin mazmunini qayta ishlab boʻlmaydi. RFC "shakl yoki qiymat" boʻlinishini belgilamaydi: bu **kursning
  kelishuvi** (shakl → `400`, qiymat → `422`); baʼzi API'lar hamma holatda `400` beradi. Go'dagi konstanta nomi `http.StatusUnprocessableEntity` (RFC 9110 matnida sarlavha "Unprocessable Content").
- **`413`.** [RFC 9110 §15.5.14](https://www.rfc-editor.org/rfc/rfc9110#name-413-content-too-large): body server
  tayyor boʻlgandan katta. Go'da `http.StatusRequestEntityTooLarge`.
- **`201 Created`** ([§15.3.2](https://www.rfc-editor.org/rfc/rfc9110#name-201-created)): yangi resurs yaratildi;
  yaratilgan mahsulot body'da (ID bilan). `Location` sarlavhasi bu darsda yoʻq (yuqorida).
- **Xato matni.** `400`dagi `message` — `Decode` xatosining matni (Chuqurroq'da uning kamchiligi). Ichki xatolarni (baza, fayl)
  mijozga bermang: ularni log'ga yozing (GB-8.02) va umumiy xabar qaytaring.
- **Limit nimani chegaralaydi.** `Decode` bitta JSON qiymatni oxirigacha oʻqiydi va u xotirada saqlanadi, shuning uchun bitta
  juda katta qiymat xavfli; `MaxBytesReader` shuni chegaralaydi. Birinchi qiymatdan keyingi baytlarni esa `Decode` oʻqimaydi,
  shuning uchun `decodeJSON` ikkinchi `Decode` bilan body tugashini tekshiradi (`11-tail`).
- **`return`.** `writeError` faqat javob yozadi, funksiyani toʻxtatmaydi.

## Manbalar

- [`net/http`: MaxBytesReader](https://pkg.go.dev/net/http#MaxBytesReader), [MaxBytesError](https://pkg.go.dev/net/http#MaxBytesError),
  [Status konstantalari](https://pkg.go.dev/net/http#pkg-constants).
- [`encoding/json`: Decoder.DisallowUnknownFields](https://pkg.go.dev/encoding/json#Decoder.DisallowUnknownFields),
  [Unmarshal](https://pkg.go.dev/encoding/json#Unmarshal) (nomaʼlum kalit tashlanadi, nom harfga qaramay moslanadi).
- [Go 1.27 release notes: `encoding/json/v2`](https://go.dev/doc/go1.27).
- [RFC 9110: 400, 413, 422, 201](https://www.rfc-editor.org/rfc/rfc9110).
