# Vazifa · GB-6.07: login'ga ikkinchi limit — email boʻyicha

Videoda `POST /login` IP boʻyicha cheklandi: bitta IP'dan 5 ta urinish, keyin har 12 soniyada bittadan. Lekin hujumchi
IP'ni almashtirsa (proksilar, botnet, IPv6 manzillar), har yangi IP'ning chelagi toʻla — Dilnozaning paroliga
taxminlar cheksiz davom etadi. OWASP Authentication Cheat Sheet ("Account Lockout"): muvaffaqiyatsiz urinishlar
hisoblagichi **hisobning oʻziga** bogʻlansin, manba IP'ga emas.

Ikkinchi limit qoʻshing: **email boʻyicha**, IP qanday boʻlishidan qatʼi nazar.

**Spetsifikatsiya**

```
bitta email'ga: 10 ta urinish, keyin har daqiqada bittadan (newLimiters(time.Minute, 10))
email harf kattaligisiz: dilnoza@bozor.test, Dilnoza@bozor.test, DILNOZA@BOZOR.TEST — bitta chelak
11-urinish (har safar boshqa IP'dan) → 429 Too Many Requests, Retry-After: <soniya>
boshqa email (jasur@bozor.test) → oddiy 401 (uning chelagi alohida)
tekshiruv argon2id'dan OLDIN: limitdan oshgan urinish serverni band qilmasin
```

**Nima qilish kerak** (`VAZIFA` izohlari turgan joylarda)

1. `cmd/api/handlers.go`: `api` struct'iga `loginEmail *limiters` maydoni; `newMux`da uni
   `newLimiters(time.Minute, 10)` bilan yarating (`loginIP` qanday yaratilgan boʻlsa, shunday).
2. `cmd/api/auth.go`, `login`: `decodeJSON`dan keyin, `checkPassword`dan **oldin** — email'ni kichik harfga
   oʻtkazing (`strings.ToLower`), `a.loginEmail.wait(...)` chaqiring; natija `0`dan katta boʻlsa
   `tooManyRequests(w, d)` va `return`. `limitByIP`ga qarang: shakli bir xil, faqat kalit boshqa.
3. 🟡 **Ixtiyoriy:** `cmd/api/ratelimit.go`'da `func (l *limiters) sweep(now time.Time) int` — chelagi **toʻla** (uzoq
   urinmagan) kalitlarni map'dan oʻchiradi va nechtasi oʻchganini qaytaradi. Aks holda map har yangi IP bilan cheksiz
   oʻsadi. Yordam: `lim.TokensAt(now) >= float64(l.burst)`. Tekshiruvi `TestSweep` (yozilmaguncha SKIP).

`*_test.go` fayllari — tekshiruv: ularni oʻzgartirmang.

**Tekshirish** (`lessons/gb-6-07-cookie-csrf-cors` papkasida; Docker kerak emas, store xotirada):

```bash
go test ./...
```

Hozir shunday yiqiladi:

```
--- FAIL: TestLoginEmailLimit (0.27s)
    tekshiruv_test.go:26: 11-urinish (dilnoza@bozor.test, yangi IP 10.0.0.11): status 401, kutilgan 429 (TASK.md: email boʻyicha limit — 10 ta urinish, harf kattaligi farq qilmaydi)
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api	0.674s
```

Vazifa tugagach `go test ./...` hamma paket uchun `ok` beradi.

**Maslahatlar**

- Test har urinishni boshqa IP'dan yuboradi (`10.0.0.1`, `10.0.0.2`, …): IP limiti ishlamaydi, faqat sizning limitingiz.
- Email body'da: uni `decodeJSON` `in.Email`ga oʻqigan. Shuning uchun bu tekshiruv middleware emas, `login`ning ichida.
- `strings.ToLower` uchun `auth.go`'ga `"strings"` importi kerak.
- Email limitining narxi bor: hujumchi Dilnozaning email'iga daqiqada bitta soʻrov yuborib, chelakni boʻsh ushlab
  turishi mumkin — Dilnoza toʻgʻri parol bilan ham 429 oladi (lockout DoS; OWASP Authentication Cheat Sheet, "Account
  Lockout": bunga qarshi chora oʻylansin, masalan parolni tiklash orqali kirish). Shuning uchun limit IP limitidan
  yumshoqroq (10 ta, keyin daqiqada bitta). Production'da qoʻshimcha choralar: faqat muvaffaqiyatsiz urinishlarni
  sanash, email+IP juftligi boʻyicha kalit, CAPTCHA.

Javob: [`solution/`](solution/) (`cmd/api/handlers.go`, `cmd/api/auth.go`, 🟡 `cmd/api/ratelimit.go`dagi `sweep`).
Tekshirish: `cd solution && go test ./...`
