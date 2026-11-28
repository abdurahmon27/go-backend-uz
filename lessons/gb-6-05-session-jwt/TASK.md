# Vazifa · GB-6.05: logout JWT'ni ham bekor qilsin (jti denylist)

Videoda koʻrdingiz: `POST /logout`dan keyin eski cookie bilan `/me` darhol `401`, **oʻsha JWT bilan esa `200`** —
serverda JWT haqida yozuv yoʻq, logout oʻchiradigan narsa yoʻq, token `exp`gacha (bu laboratoriyada 15 daqiqa)
yaroqli. Endi buni tuzatasiz va tuzatishning narxini oʻz qoʻlingiz bilan koʻrasiz: **denylist** — bekor qilingan
tokenlarning `jti`lari roʻyxati, uni har JWT soʻrovida tekshirish kerak (RFC 7519 §4.1.7 `jti`; OWASP JWT Cheat Sheet,
"JWT denylist").

**Nima qilish kerak**

1. **`internal/store/store.go` → `RevokeToken` va `TokenRevoked`.** `Store`da `revoked map[string]time.Time` maydoni
   bor (kalit — `jti`, qiymat — tokenning `exp` vaqti).
   - `RevokeToken(ctx, jti, expires)` — `s.revoked[jti] = expires` va `nil`. Boshqa metodlardagidek avval
     `s.mu.Lock()` va `defer s.mu.Unlock()`: soʻrovlar bir vaqtda keladi (GB-4.07).
   - `TokenRevoked(ctx, jti)` — `jti` map'da bormi: comma-ok bilan oʻqing (GB-2.02) va `ok, nil` qaytaring.
2. **`cmd/api/auth.go` → `logout`.** `VAZIFA` izohi turgan joyda: soʻrov `Authorization: Bearer …` bilan kelgan boʻlsa
   (`bearer(r)`), tokenni `token.Parse(s, a.jwtKey)` bilan oching. `Parse` xato bermasa,
   `a.st.RevokeToken(r.Context(), c.JTI, c.Expires)` chaqiring; `RevokeToken` xato qaytarsa — `writeError(w, err)` va
   `return`. `Parse` xato bersa (token muddati oʻtgan, imzosi notoʻgʻri), hech narsa qilmang: bunday token baribir
   kirita olmaydi, logout esa GB-6.04'dagidek `204` qaytaradi. Cookie qismiga tegmang.
3. **`cmd/api/auth.go` → `authenticate`.** JWT yoʻlida, `token.Parse`dan keyin: `a.st.TokenRevoked(r.Context(), c.JTI)`
   `true` qaytarsa — `errBadToken` bilan oʻralgan xato (`fmt.Errorf("%w: …", errBadToken, …)`), shunda `problemFor` uni
   `401` qiladi va `WWW-Authenticate` sarlavhasi qoʻyiladi.

🟡 **Ixtiyoriy: `PurgeRevoked(ctx, now)`.** `exp`i `now`dan oldin boʻlgan yozuvlarni `s.revoked`dan oʻchiring va nechtasi
oʻchganini qaytaring. Muddati oʻtgan token `token.Parse`dan baribir oʻtmaydi, uning yozuvi keraksiz — tozalanmasa, map
abadiy oʻsadi. Yozmasangiz, bu tekshiruv `SKIP` boʻladi.

`tekshiruv_test.go` fayllarini oʻzgartirmang.

**Tekshirish** (papkaning ildizida):

```bash
go test ./...
```

Hozir (vazifadan oldin) ikki test yiqiladi:

```
--- FAIL: TestLogoutJWTniBekorQiladi (0.02s)
    tekshiruv_test.go:19: logout'dan keyin oʻsha JWT bilan /me: status 200, kutilgan 401 — jti denylist'ga tushmadi yoki authenticate uni tekshirmayapti (TASK.md, 2–3-qadam)
--- FAIL: TestRevokeToken (0.00s)
    tekshiruv_test.go:15: RevokeToken: vazifa: RevokeToken hali yozilmagan (TASK.md) — TASK.md, 1-qadam: s.revoked'ga yozing va nil qaytaring
```

Hammasi `ok` boʻlsa — vazifa bajarildi. Qoʻlda ham koʻring: 1-terminalda `go run ./cmd/api`, 2-terminalda videodagi
login va `T=$(jq -r .token login.json)`, keyin `curl -s -X POST -H "Authorization: Bearer $T" localhost:8765/logout` va
`curl -s -H "Authorization: Bearer $T" localhost:8765/me` — endi `401`, server logida esa `JWT yaroqsiz: jti … bekor qilingan`.

**Maslahatlar**

- `c.JTI` (`string`) va `c.Expires` (`time.Time`) — `token.Parse` qaytargan `token.Claims`ning maydonlari
  (`internal/token/token.go`, `type Claims`). `Store` struct'i va `New()` — `internal/store/store.go`ning boshida.
- `logout` — `cmd/api/auth.go` oxirida: hozir u cookie'dagi session'ni oʻchiradi (`endCurrentSession`, GB-6.04) va
  cookie'ni tozalaydi. Siz JWT qismini `VAZIFA` izohi turgan joyga qoʻshasiz.
- `logout` GB-6.04'dagidek `requireUser`siz (eskirgan cookie ham tozalansin): token shu yerda birinchi marta
  tekshiriladi. Yaroqsiz Bearer bilan logout ham `204` boʻlishini `api_test.go` tekshiradi.
- Oʻylab koʻring: endi har JWT soʻrovi store'dan bitta qidiruv qiladi — session'dagidek. Session'ga nisbatan JWT'dan nima
  qoldi? (README, "Xulosa").
- Javob: `solution/` (oʻz `go.mod`i bilan; `cd solution && go test ./...`).
