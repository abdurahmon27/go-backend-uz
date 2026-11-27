# GB-4.02 · Go'da HTTP client: Timeout, status va `resp.Body.Close`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #26"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-4-01-http-matn`](../gb-4-01-http-matn) (HTTP almashinuvi matn sifatida: `curl -v`) · keyingi dars:
`gb-4-03` (handler — oddiy funksiya, `httptest`). Bu darsga kerak boʻlganlar:
[`gb-3-02-json-saqlash`](../gb-3-02-json-saqlash) (`json.NewDecoder`, struct tag'lar, `&x`),
[`gb-2-09-defer-panic`](../gb-2-09-defer-panic) (`defer f.Close()` xato tekshiruvidan keyin),
[`gb-2-08-xatolar`](../gb-2-08-xatolar) (`%w`, `errors.As`), [`gb-2-02-map-qidirish`](../gb-2-02-map-qidirish)
(map'da yoʻq kalit — nol qiymat).

**Bu darsdan keyin:** `net/http` client'ini toʻgʻri ishlatasiz: `http.Client{Timeout: …}` bilan kutishni cheklaysiz;
javob kelganda ham `resp.StatusCode`ni oʻzingiz tekshirasiz (`404` xato emas!); `resp.Body`ni `defer resp.Body.Close()`
bilan yopasiz (xato tekshiruvidan **keyin**) va JSON'ni struct'ga oʻqiysiz. bozorda bu narxni dollarda ham
koʻrsatadigan kurs funksiyasi: `fetchRate`.

Bu papka mustaqil: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i, tashqi modul yoʻq.
`Product` (`ID`, `Name`, `Price` tiyinda, `Stock`) GB-2.03/GB-3.02'dan; yangi: `rateResponse`, `fetchRate`, `StatusError`.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `fetchRate` (Timeout, status, `Body.Close`, Decode), `run`, `main`; `StatusError` vazifa uchun tayyor tur |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang. Haqiqiy internetga chiqmaydi: har test oʻz soxta serverini ochadi |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 01-get go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Kurs API'si: `https://open.er-api.com/v6/latest/USD` (kalit shart emas; [shartlari](https://www.exchangerate-api.com/docs/free):
kurslarni koʻrsatsangiz, [Rates By Exchange Rate API](https://www.exchangerate-api.com) havolasi kerak; soʻrovlar soni cheklangan,
artiq boʻlsa `429` keladi va ~20 daqiqadan keyin tiklanadi: shuning uchun testlar bu API'ga chiqmaydi).
Kurs har kuni oʻzgaradi, sizda boshqa son chiqadi. Javob (qisqartirilgan): `{"result":"success","base_code":"USD",
"rates":{"EUR":0.911313,…,"UZS":11966.840025,…}}`.

```
$ go run .                      (01-get: http.Get, resp.Status va resp.StatusCode)
200 OK 200
```

`http.Get` standart client'ni ishlatadi, uning `Timeout`i 0 (cheklov yoʻq). Javob bermaydigan server: boshqa terminalda
`nc -l 8089` (ulanishni qabul qiladi, hech narsa yozmaydi). `nc` dasturlari har xil: LibreSSL/OpenBSD netcat'da shunday, GNU
netcat'da `nc -l -p 8089`. `02-sukut`, `http.Get("http://localhost:8089")`:

```
$ go run .                      (Control C bilan toʻxtatildi)
soʻrov yuborildi
^Csignal: interrupt
```

`http.Client{Timeout: 2 * time.Second}` (`03-timeout`):

```
$ go run .
xato: Get "http://localhost:8089": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
exit status 1
```

Mavjud boʻlmagan manzil (`04-404`, `…/v6/nothing`): server javob berdi, shuning uchun `err == nil`:

```
$ go run .
err: <nil>
404 Not Found 404
```

`resp.Body`ni `defer`dan keyin `json.NewDecoder(resp.Body).Decode(&data)` (`05-decode`):

```
$ go run .
USD 11966.840025
```

Eng koʻp uchraydigan xato (`06-vet`): `defer resp.Body.Close()` err tekshiruvidan **oldin**. `Get` xato bersa, `resp` nil
va `resp.Body` panic beradi (`panic: runtime error: invalid memory address or nil pointer dereference`); `go vet`
buni topadi:

```
$ go vet .
main.go:13:8: using resp before checking for errors
```

Tayyor `main.go` (`./_video/run.sh final go run .`):

```
$ go run .
1 USD = 11966.84 soʻm
olma: 12000 soʻm = 1.00 USD
```

Uchinchi tuzoq (`07-xato200`): API valyuta kodini bilmasa (`…/latest/XXX`), status `200`, body
`{"result":"error","error-type":"unsupported-code"}` (hujjatda `result` va `error-type` bor; `unsupported-code` qiymatini
biz 2026-11-27 kuni shu API'da kuzatdik); videodagi `fetchRate` jim nol qaytaradi (`Rates` boʻsh, map'da yoʻq
kalit nol beradi). Vazifa shuni tuzatadi:

```
$ go run .
0 <nil>
```

Chuqurroq (videoda kod bilan): `io.LimitReader(resp.Body, 1<<20)` — tashqi serverning cheksiz katta javobidan himoya (vazifaning
3-qadami); bitta `http.Client`ni qayta ishlating (u bir vaqtda ishlatishga xavfsiz).

Narxni dollarga aylantirish faqat **koʻrsatish** uchun `float64`da; saqlash va hisob-kitob tiyinda `int64` (GB-1.03).

## Manbalar

- `net/http`: [Client](https://pkg.go.dev/net/http#Client) (Timeout: ulanish, yoʻnaltirish va body oʻqishni qamrab oladi;
  0 = cheklov yoʻq; standart transport ulanish/TLS uchun oʻz cheklovlariga ega (30 s / 10 s), javob sarlavhalarini kutishga yoʻq; client'larni qayta ishlatish kerak), [Get](https://pkg.go.dev/net/http#Get) (non-2xx javob xato emas;
  xato: ulanish, yoʻnaltirishlar koʻpligi, protokol; caller `resp.Body`ni yopishi kerak), [Response](https://pkg.go.dev/net/http#Response) (Body).
- [RFC 9110, status kodlari](https://www.rfc-editor.org/rfc/rfc9110#name-status-codes).
- `go vet`: [httpresponse](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/httpresponse).
- [ExchangeRate-API (open access)](https://www.exchangerate-api.com/docs/free).
