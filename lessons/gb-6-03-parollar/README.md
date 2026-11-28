# GB-6.03 · Parollar: argon2id, salt va hech qachon ochiq matn emas

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #45"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-6-02-api-xatolar`](../gb-6-02-api-xatolar) (RFC 9457 xato javoblari, `problemFor`, `writeError`) ·
keyingi dars: `gb-6-04` (login va bekor qilinadigan session'lar). Bu papka GB-6.02'ning `solution/`idan boshlanadi.

**Bu darsdan keyin:** bozor'da roʻyxatdan oʻtish bor (`POST /users`), parol esa bazaga faqat **argon2id satri** boʻlib
yoziladi: `$argon2id$v=19$m=19456,t=2,p=1$<salt>$<kalit>`. Har foydalanuvchiga `crypto/rand`dan yangi 16 baytli salt;
parametrlar — [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#argon2id)ning
minimal sozlamasi. Tekshirish (`Verify`) salt va parametrlarni **satrning oʻzidan** oʻqiydi va kalitni
`crypto/subtle.ConstantTimeCompare` bilan solishtiradi. Sizib chiqqan jadvalga oflayn hujumni shu kompyuterda oʻlchaymiz.

**Foydalanuvchilar xotirada** (`internal/store/users.go`, katalog kabi): bu dars parol saqlash haqida, Postgres emas.
Postgres'da xuddi shu satr `users` jadvalining `password_hash text` ustunida turadi; kod oʻzgarmaydi.

## Papka

```
gb-6-03-parollar/
  go.mod, go.sum               module github.com/abdurahmon27/bozor, go 1.27; YANGI: golang.org/x/crypto (+ x/sys)
  internal/password/           YANGI: Hash, Verify, NeedsRehash (vazifa), decode, ErrBadHash
    password_test.go           format, salt, Verify, buzuq satrlar (panic yoʻq, «har parol mos» yoʻq)
    tekshiruv_test.go          vazifa tekshiruvi; oʻzgartirmang
  internal/store/users.go      YANGI: User (PasswordHash `json:"-"`), Users.Add, Users.ByEmail, ErrEmailTaken
  cmd/api/users.go             YANGI: POST /users — userInput, validateUser (15–64 belgi), register
  cmd/api/handlers.go          api{st, users}; route POST /users
  cmd/api/problem.go           problemFor: + ErrEmailTaken → 409
  cmd/api/tekshiruv_test.go    GB-6.02 vazifasining tekshiruvi (bajarilgan holda, oʻtadi)
  cmd/taxmin/                  oflayn hujum oʻlchovi: sizgan jadval × nomzodlar.txt
  sizdi/                       «sizgan» jadval uch xil: ochiq.txt, sha256.txt, argon2id.txt (oʻylab topilgan)
  nomzodlar.txt                1000 ta nomzod parol (bozor soʻzlari × raqamlar; oʻylab topilgan, haqiqiy roʻyxat emas)
  cmd/bozor/, qolgan cmd/api   GB-6.02 bilan bir xil
  TASK.md                      vazifa
  solution/                    javob (oʻz go.mod'i bilan, alohida ishlaydi)
  _video/                      record.sh → chiqishlar.txt; sizdir/ (sizdi/ va nomzodlar.txt'ni yasaydi); xato/ (base64 xatosi)
```

GB-6.02 bilan farq: `diff -ru ../gb-6-02-api-xatolar/solution . -x _video -x solution -x sizdi -x nomzodlar.txt`.

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`, curl 8.22.0, Intel Core i5-12450H)

Hammasi `_video/record.sh` bilan shu kompyuterda olingan (toʻliq: `_video/chiqishlar.txt`).

**Sizib chiqqan jadval** — `sizdi/` fayllarini `go run ./_video/sizdir` yasaydi (emaillar `.test` domenida, RFC 2606):

```
$ cat sizdi/sha256.txt
dilnoza@bozor.test 385426486641a520e865c4e6b24d379c8f8afa1eb5de2f932ba90f0622f06289
jasur@bozor.test 385426486641a520e865c4e6b24d379c8f8afa1eb5de2f932ba90f0622f06289
malika@bozor.test f8fe182b2bef6749451694c3042b04e5cd5f5aa6aeecbfb0c2744dfbac6549a5

$ go run ./cmd/taxmin sizdi/sha256.txt
sizdi/sha256.txt: 3 qator · nomzodlar.txt: 1000 parol
dilnoza@bozor.test  topildi: choy2024 (71-nomzod) · 9µs
jasur@bozor.test    topildi: choy2024 (71-nomzod) · 8µs
malika@bozor.test   topilmadi · 94µs
1142 taxmin · 139µs · sekundiga 8 229 029 · har taxmin 0 B xotira
```

Salt yoʻq: bir xil parol — bir xil hash. sha256 bitta yadroda sekundiga ~8 mln taxmin (shu kompyuterda, bitta oʻlchov) va
bu oʻlchovda har taxmin uchun heap'dan xotira ajratilmadi (`runtime.MemStats.TotalAlloc`).

```
$ cat sizdi/argon2id.txt
dilnoza@bozor.test $argon2id$v=19$m=19456,t=2,p=1$fd+QrkyR3ap9uAJkzw4Ocw$1tRXqM9ibZAPh1bCRLKEoQtHxfvyHk/HuicsoJA/CsQ
jasur@bozor.test $argon2id$v=19$m=19456,t=2,p=1$rUPhzHtvN8st9ZvrLtD39A$QDmgSweLTZ3B4crIK4E1dj13e/GF9VEcFyNE49sErHY
malika@bozor.test $argon2id$v=19$m=19456,t=2,p=1$GX6RH2ZGTsRb4K5+L5A5Ew$nsERR+E5HtqLtBfmTmTOsYGRpK9blc8HIdMX/k2/SRY

$ go run ./cmd/taxmin sizdi/argon2id.txt
sizdi/argon2id.txt: 3 qator · nomzodlar.txt: 1000 parol
dilnoza@bozor.test  topildi: choy2024 (71-nomzod) · 1.3s
jasur@bozor.test    topildi: choy2024 (71-nomzod) · 1.3s
malika@bozor.test   topilmadi · 19.7s
1142 taxmin · 22.3s · sekundiga 51 · har taxmin 19.0 MiB xotira
```

`taxmin` argon2id satri uchun bozor'ning oʻz `password.Verify`ini chaqiradi. Raqamlar — bitta oʻlchov, shu kompyuterda;
kompyuter band paytida (2026-11-28, load average ≈ 15) qayta oʻlchaganda sha256 sekundiga 7 101 283, argon2id 27 chiqdi
(`_video/qayta-2026-11-28.txt`): mutlaq son oʻzgaradi, farq (yuz ming martadan ortiq) qoladi. Zaif parol (`choy2024`) baribir topiladi —
argon2id uni faqat **sekinlashtiradi** (~160 000 marta). Bir milliard nomzod: sha256 bilan ~2 daqiqa, argon2id bilan
~227 kun — bitta yadroda, **har foydalanuvchi uchun alohida** (salt har xil). Shuning uchun uzunlik qoidasi ham kerak.

**Nega xotira narxi muhim.** Hujumchi videokartada minglab yadro bilan bir vaqtda taxmin qiladi; sha256'ga xotira
kerak emas. argon2id'ning har taxmini 19 MiB oladi: 24 GiB kartaga (masalan, 24 GB xotirali hozirgi oʻyin kartalari)
bir vaqtda ≈ 1293 taxmindan koʻpi sigʻmaydi (24 · 1024 / 19; karta xotirasi boshqa narsaga ham ketadi, haqiqiy son kamroq). RFC 9106 §1: argon2 — *memory-hard* funksiya, maqsadi maxsus qurilmalarda (GPU, ASIC) taxmin narxini oshirish.

**Roʻyxatdan oʻtish** (server birinchi terminalda: `go run ./cmd/api` → `bozor API: http://localhost:8765`):

```
$ curl -i localhost:8765/users -H "Content-Type: application/json" -d '{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}'
HTTP/1.1 201 Created
Content-Type: application/json
…

{"id":1,"email":"dilnoza@bozor.test"}

$ curl -s localhost:8765/users -H "Content-Type: application/json" -d '{"email":"Dilnoza@bozor.test","password":"boshqa-uzun-ibora-2025"}'
{"type":"about:blank","title":"Conflict","status":409,"detail":"bu email bilan foydalanuvchi bor","request_id":"MPOP6LIHWO4DKRWFVGFX742BV4"}

$ curl -s localhost:8765/users -H "Content-Type: application/json" -d '{"email":"jasur@bozor.test","password":"choy2024"}'
{"type":"/problems/validation","title":"Qiymatlar notoʻgʻri","status":422,"errors":{"password":"15–64 belgi boʻlsin"},"request_id":"L45WTPQR4YBCK3R5ZLF5DK3NTY"}
```

Server log'ida roʻyxatdan oʻtishning narxi koʻrinadi: `POST "/users" 201 29.967944ms` — deyarli hammasi argon2id.
`409` ham ~38 ms: `Hash` `Add`dan oldin hisoblanadi.

**Tez-tez uchraydigan xato:** `decode` salt'ni `base64.StdEncoding` bilan oʻqisa (`_video/xato/password.go.txt`),
`Hash` esa `RawStdEncoding` (oxirida `=` yoʻq) bilan yozgan boʻlsa:

```
$ go test ./internal/password
…
--- FAIL: TestVerify (0.02s)
    password_test.go:39: toʻgʻri parol: Verify = false, hash satri notoʻgʻri: salt: illegal base64 data at input byte 20; kutilgan true, nil
…
FAIL
```

16 bayt salt = 22 belgi base64. `StdEncoding` 4 belgilik boʻlaklarni oʻqiydi: 20 belgidan keyin qolgan 2 belgi toʻliq
boʻlak emas, `==` kutiladi. Tuzatish: yozishda ham, oʻqishda ham bitta `b64` (`RawStdEncoding`).

## Xavfsizlik eslatmalari (videodan tashqari)

- **Parametrlarni oʻz serveringizda oʻlchang** — OWASP jadvali *minimal* sozlama; OWASP parametrlarni ishlaydigan
  serverda sinab tanlashni aytadi (bitta `Hash` shu yerda ≈ 20–45 ms, yuklamaga qarab). Xotira koʻp boʻlsa — yuqoriroq qator.
- **Parol uzunligi** — [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#implement-proper-password-strength-controls)
  va NIST SP 800-63B: ikkinchi omil (MFA) yoʻq boʻlsa kamida 15 belgi; eng katta uzunlik **kamida** 64 belgi boʻlishi kerak
  (bozor aynan 64 oladi; 128 ham toʻgʻri). Belgilar `rune`da sanaladi. Bitta parolning turli Unicode yozuvlari (masalan
  «oʻ» bitta belgi yoki harf + belgi) turli hash beradi; normallashtirish (NFC/NFKC, `golang.org/x/text`) bozor'da hali yoʻq.
  Tarkib qoidalari ("bitta katta harf…") yoʻq — NIST ularni tavsiya qilmaydi. Mashhur/sizgan parollar roʻyxati bilan
  tekshirish (blocklist) — keyingi qadam.
- **`409` email band** — bu javob kim roʻyxatdan oʻtganini aytadi (account enumeration). OWASP buni roʻyxatdan oʻtishda
  toʻliq yopish qiyinligini aytadi; toʻliq yechim — emailga tasdiqlash havolasi. bozor hozircha oddiy `409` bilan.
- **Xotira va DoS** — har `POST /users` ~20 ms CPU va 19 MiB xotira oladi. Login va roʻyxatdan oʻtishga soʻrov cheklovi
  (rate limit) — GB-6.07.
- **Log** — `writeError` log'ga `err.Error()` yozadi; validatsiya xatosi matnida parol yoʻq (`validatsiya: 1 ta maydon notoʻgʻri`).
  Lekin `decodeJSON`ning JSON sintaksis xatosi body'dan **bitta belgini** keltirishi mumkin (`invalid character '…' in string`):
  parol maydonidagi buzuq belgi log'ga tushadi. Toʻliq yechim — log'ga xato turini yozish, matnini emas (tuzilgan log, GB-8.02).
  Parolning oʻzi hech qachon log'ga yozilmasin.
- **Buzuq satr** — `decode` `t=0`/`p=0`ni (IDKey `panic` qilardi), 16 baytdan qisqa kalitni (boʻsh kalit bilan
  `ConstantTimeCompare` "teng" derdi), juda katta parametrlarni (`m` > 256 MiB, `t` > 16: `m=4294967295` jarayonni
  `fatal error: out of memory` bilan oʻldirardi, buni `recover` ushlamaydi) va `p=1abc` kabi ortiqcha belgini rad etadi:
  `password_test.go`'dagi `TestVerifyBadHash`. Videoda `decode`ning shu oxirgi tekshiruvi koʻrsatilmagan.
- **Doimiy vaqt** — `Verify`da solishtiriladigan kalit hujumchining oʻz parolidan hisoblanadi, shuning uchun bu yerda
  `ConstantTimeCompare` — qoʻshimcha ehtiyot. Sir toʻgʻridan-toʻgʻri solishtiriladigan joyda (session token, GB-6.04) u
  majburiy. Uzunlik farq qilsa, funksiya darhol 0 qaytaradi: vaqt faqat uzunlikni aytadi.
- **Email** — `Add` emailni yozilganicha saqlaydi va `strings.EqualFold` bilan solishtiradi; `@` tekshiruvi minimal.
  Postgres'da (GB-6.04) bitta shakl — `lower(email)` boʻyicha unikal indeks.

## Vazifa

[TASK.md](TASK.md): parametrlarni OWASP'ning birinchi qatoriga (`m=47104`, `t=1`, `p=1`) oʻtkazish va `NeedsRehash` —
eski parametrli satrni topadigan funksiya. Tekshirish: `go test ./...`.

```
$ go test ./...          # starter
ok  	github.com/abdurahmon27/bozor/cmd/api	0.086s
ok  	github.com/abdurahmon27/bozor/cmd/bozor	0.082s
?   	github.com/abdurahmon27/bozor/cmd/taxmin	[no test files]
--- FAIL: TestVazifaParams (0.02s)
--- FAIL: TestVazifaNeedsRehash (0.02s)
FAIL
FAIL	github.com/abdurahmon27/bozor/internal/password	0.213s
ok  	github.com/abdurahmon27/bozor/internal/store	0.251s
FAIL

$ cd solution && go test ./...
ok  	github.com/abdurahmon27/bozor/cmd/api	0.176s
ok  	github.com/abdurahmon27/bozor/cmd/bozor	0.117s
?   	github.com/abdurahmon27/bozor/cmd/taxmin	[no test files]
ok  	github.com/abdurahmon27/bozor/internal/password	0.414s
ok  	github.com/abdurahmon27/bozor/internal/store	0.238s
```

## Manbalar

- OWASP Password Storage Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
  (argon2id: "m=19456 (19 MiB), t=2, p=1", "m=47104 (46 MiB), t=1, p=1", …; salting; "Fast hashing algorithms such as
  SHA‑256 are not suitable"; upgrading the work factor)
- OWASP Authentication Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html
  (15 / 8 belgi, 64 belgigacha, NIST SP 800-63B)
- RFC 9106, Argon2 Memory-Hard Function — https://www.rfc-editor.org/rfc/rfc9106 (§1 memory-hard, §3.1 parametrlar:
  salt 16 bayt, kalit 32 bayt tavsiya; §4 parametrlarni tanlash)
- PHC string format — https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md (`$id$v=…$params$salt$hash`, base64 toʻldiruvchisiz)
- Go: `golang.org/x/crypto/argon2` (`IDKey`, `Version`) — https://pkg.go.dev/golang.org/x/crypto/argon2 ·
  `crypto/rand.Read` (Go 1.24'dan beri xato qaytarmaydi) — https://pkg.go.dev/crypto/rand#Read ·
  `crypto/subtle.ConstantTimeCompare` — https://pkg.go.dev/crypto/subtle#ConstantTimeCompare ·
  `encoding/base64.RawStdEncoding` — https://pkg.go.dev/encoding/base64 · `fmt.Sscanf` — https://pkg.go.dev/fmt ·
  `encoding/json` tag `"-"` — https://pkg.go.dev/encoding/json#Marshal
- RFC 2606 (`.test` domeni) — https://www.rfc-editor.org/rfc/rfc2606
