# Vazifa · GB-6.04: "hamma qurilmadan chiqish" — `POST /logout-all`

Dilnoza telefonini yoʻqotdi. Telefonda ham, noutbukda ham bozor'ga kirgan: bazada uning ikkita session qatori bor.
`POST /logout` faqat **shu** cookie'ning session'ini o'chiradi — telefondagi token tirik qoladi. Noutbukdan bitta
so'rov bilan Dilnozaning **hamma** session'larini o'chiradigan endpoint yozing. Boshqa foydalanuvchilarning
session'lariga tegmasin.

**Spetsifikatsiya**

```
POST /logout-all   (Cookie: session=<Dilnozaning istalgan tokeni>)
→ 204 No Content, Set-Cookie: session=; …; Max-Age=0 …

keyin: Dilnozaning telefon tokeni ham, noutbuk tokeni ham → GET /me → 401
       Jasurning tokeni → GET /me → 200 (oʻzgarmaydi)
cookie'siz POST /logout-all → 401
```

**Nima qilish kerak**

1. `internal/store/store.go`: `DeleteUserSessions(ctx, userID)` — `VAZIFA` izohi turgan joyda. Bitta SQL:
   `DELETE FROM sessions WHERE user_id = $1`. `Exec` natijasidagi `RowsAffected()` (GB-5.08) — o'chgan qatorlar soni,
   shuni qaytaring. Xatoni `%w` bilan o'rang.
2. `cmd/api/handlers.go` va `cmd/api/auth.go`: `newMux`dagi `VAZIFA` izohi o'rniga `POST /logout-all` route'i,
   `requireUser` bilan o'ralgan (`GET /me` qanday ulangan bo'lsa, shunday), va `logoutAll` handler'i: foydalanuvchi
   `currentUser(r)`dan, `a.st.DeleteUserSessions(...)`, xato bo'lsa `writeError`, keyin `clearSessionCookie(w)` va
   `204`. `logout` handler'iga qarang: shakli bir xil, faqat qaysi qatorlar oʻchishi boshqa.
3. 🟡 **Ixtiyoriy:** `DeleteExpiredSessions(ctx)` — muddati o'tgan qatorlarni tozalash (`expires_at <= now()`).
   `UserBySession` ularni baribir qabul qilmaydi; bu faqat jadval cheksiz o'sib ketmasligi uchun. Tekshiruvi
   `TestDeleteExpiredSessions` (yozilmaguncha SKIP).

`*_test.go` fayllari — tekshiruv: ularni o'zgartirmang.

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test ./...
```

Har test o'z vaqtinchalik sxemasida ishlaydi (`internal/pgtest`, GB-5.10): migratsiyalar va `seed.sql` o'sha sxemaga
qo'llanadi, test tugagach sxema o'chiriladi. Postgres ishlamasa test aniq xabar bilan yiqiladi (o'tkazib yuborilmaydi):

```
store_test.go:14: Postgres'ga ulanib boʻlmadi (localhost:5514): avval shu papkada `docker compose up -d --wait` ishga tushiring.
```

Boshqa manzil uchun `DATABASE_URL`.

Hozir shunday yiqiladi (`_video/chiqishlar.txt` §11):

```
--- FAIL: TestLogoutAll (0.14s)
    tekshiruv_test.go:18: POST /logout-all: status 404, kutilgan 204 (TASK.md, 2-qadam): {"type":"about:blank","title":"Not Found","status":404,"detail":"bunday yoʻl yoʻq",…}
FAIL	github.com/abdurahmon27/bozor/cmd/api	1.175s
ok  	github.com/abdurahmon27/bozor/internal/password	0.266s
--- FAIL: TestDeleteUserSessions (0.07s)
    tekshiruv_test.go:27: DeleteUserSessions(1): DeleteUserSessions hali yozilmagan (TASK.md, 1-qadam)
FAIL	github.com/abdurahmon27/bozor/internal/store	0.338s
```

Vazifa tugagach `go test ./...` hamma paket uchun `ok` beradi. Tugagach: `docker compose down -v`.

**Maslahatlar**

- 1-qadam `DeleteSession`ga juda o'xshaydi; farqi — `WHERE` sharti va qaytadigan son (`tag.RowsAffected()`).
- `logoutAll`da tokenni cookie'dan qayta o'qish shart emas: kimligini `requireUser` context'ga qo'ygan, `currentUser(r).ID`.
- `DeleteUserSessions` ikki natija qaytaradi; handler'ga son kerak emas: `_, err := …` (GB-1.08).
- Route `requireUser`siz ulansa, `currentUser` panic qiladi va test `500` ko'radi: bu ataylab (README, "Chuqurroq").

Javob: [`solution/`](solution/) (`internal/store/store.go`dagi `DeleteUserSessions`, `cmd/api/auth.go`dagi `logoutAll`).
Tekshirish: `cd solution && docker compose up -d --wait && go test ./...`
