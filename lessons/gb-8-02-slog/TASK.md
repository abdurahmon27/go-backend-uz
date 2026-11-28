# Vazifa · GB-8.02: login qatori — email'siz

Muvaffaqiyatli login — xavfsizlik uchun muhim hodisa: kim, qachon kirdi (OWASP Logging Cheat Sheet: autentifikatsiya
hodisalari log'ga yoziladi). Lekin log'da **email boʻlmasin** (shaxsiy maʼlumot) va parol hash'i ham boʻlmasin (secret).
Foydalanuvchini log'da uning `id`si va `role`i tanitadi; soʻrovning qolgan qatorlarini `request_id` bogʻlaydi.

**Kutilgan natija** (`POST /login`, admin bilan; bitta JSON qator, `time` qisqartirilgan):

```
{"time":"…","level":"INFO","msg":"login","request_id":"…","user":{"id":3,"role":"admin"}}
```

**Nima qilish kerak**

1. `cmd/api/auth.go`, `login` handler'i: cookie yozilgandan keyin (`setSessionCookie(w, token)` qatoridan keyin)
   soʻrovning logger'i bilan Info darajada bitta qator yozing: xabar `"login"`, atribut `"user"` — qiymati `u`
   (`store.User`, login qilgan foydalanuvchi). Soʻrovning logger'ini videodagi `loggerFrom` beradi: unda `request_id`
   allaqachon bor (`writeError`da qanday olinganiga qarang).
2. Testni ishga tushiring: qatorda **email chiqadi** — `"user":{"id":3,"email":"admin@bozor.test","role":"admin"}`.
   Sababi: JSON handler `User`ni `json.Marshal` qiladi, `Email` maydonida esa `json:"email"` tag'i bor.
   `internal/store/store.go`da `User`ga `LogValue() slog.Value` metodini yozing (videodagi `config.LogValue` kabi):
   `slog.GroupValue` ichida faqat ikkita atribut — `"id"` (`slog.Int64`) va `"role"` (`slog.String`).

`internal/store/tekshiruv_test.go` va `cmd/api/tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (papkaning ildizida)

2-qadam, Postgres kerak emas:

```bash
go test ./internal/store -run LogValue
```

Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestUserLogValue (0.00s)
    tekshiruv_test.go:22: User'da LogValue metodi yoʻq: JSON handler uni json.Marshal qiladi va email chiqadi (TASK.md, 2-qadam).
        hozir: {"time":"…","level":"INFO","msg":"login","user":{"id":3,"email":"admin@bozor.test","role":"admin"}}
    tekshiruv_test.go:26: log qatorida email yoki parol hash'i bor:
        …
FAIL	github.com/abdurahmon27/bozor/internal/store
```

1-qadam, haqiqiy Postgres bilan:

```bash
docker compose up -d --wait
go test ./cmd/api -run LoginLog
```

Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestLoginLogLine (0.12s)
    tekshiruv_test.go:31: log'da "login" qatori yoʻq (TASK.md, 1-qadam: login handler'ida, cookie yozilgandan keyin soʻrovning logger'i bilan Info darajada: xabar "login", atribut "user"). Hozirgi log:
        {"time":"…","level":"INFO","msg":"soʻrov","request_id":"…","method":"POST","path":"/login","status":200,"dur":…}
FAIL	github.com/abdurahmon27/bozor/cmd/api
```

Ikkalasi `ok` boʻlguncha ishlang. Keyin butun papka va tozalash:

```bash
go test ./...
docker compose down -v
```

Qoʻlda ham tekshiring: serverni `./api -port 8782 > api.log` bilan ishga tushirib (README, «Videodagi buyruqlar»),
seed'dagi admin bilan login qiling (`admin@bozor.test`, parol `qovun-anjir-tut-uzum-31` — faqat laboratoriya uchun,
`seed.sql`; baza: `goose up` va seed, README) va `grep '"msg":"login"' api.log` — `"user":{"id":3,"role":"admin"}`, email yoʻq.

**Maslahatlar**

- `func (u User) LogValue() slog.Value { … }` — qiymat qabul qiluvchi (receiver), `config.LogValue`dagidek. `store.go`ning
  `import`iga `"log/slog"` qoʻshiladi.
- `slog.Int64(key, v)` — `int64` atribut, `slog.String(key, v)` — matn atribut: `go doc log/slog.Int64`.
- Testlar `slog.LogValuer` interfeysini tekshiradi: metod nomi aynan `LogValue`, qaytaradigan turi `slog.Value`.

Javob: [`solution/cmd/api/auth.go`](solution/cmd/api/auth.go) va
[`solution/internal/store/store.go`](solution/internal/store/store.go). Tekshirish: `cd solution && go test ./...`.
