# Vazifa · GB-8.03: `waitBG` muddatni hurmat qilsin

Videoda `run` toʻxtashda shunday tartibda ishlaydi: `srv.Shutdown(shCtx)` → `waitBG(shCtx, &bg)` → `defer pool.Close()`.
`shCtx` — `shutdownTimeout` (8 soniya) muddatli context. `srv.Shutdown` bu muddatni hurmat qiladi: muddat tugasa,
`context deadline exceeded` bilan qaytadi. `waitBG` esa hozir faqat `bg.Wait()`: bitta chek osilib qolsa (email servisi
javob bermasa), `run` ham cheksiz osilib qoladi va `docker stop` 10 soniyadan keyin jarayonni `SIGKILL` bilan oʻldiradi —
log'da sababi ham qolmaydi.

**Spetsifikatsiya** (`cmd/api/receipt.go`, `func waitBG(ctx context.Context, bg *sync.WaitGroup) error`)

```
hamma fon ishi tugadi, ctx hali tirik     → nil (ish tugaguncha kutadi)
ctx muddati oldin tugadi                  → ctx.Err(), yaʼni context.DeadlineExceeded — darhol, ishni kutmasdan
```

`run` allaqachon `waitBG`ning xatosini `fon ishlari: …` deb oʻrab qaytaradi, `main` uni yozib, `1` kodi bilan chiqadi:
boshqa joyni oʻzgartirish shart emas. `cmd/api/toxtash_test.go`ni oʻzgartirmang.

**Tekshirish** (papkaning ildizida; Postgres kerak emas):

```bash
go test ./cmd/api -run WaitBG
```

Hozir shunday yiqiladi:

```
--- FAIL: TestWaitBGDeadline (2.00s)
    toxtash_test.go:57: waitBG muddat (100 ms) tugagach 2 s oʻtdi, hali ham chekni kutyapti: ctx.Done() yopilganda kutishni toʻxtatib, ctx.Err() qaytarishi kerak (TASK.md)
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/api	2.106s
FAIL
```

`ok` boʻlguncha ishlang. Keyin butun papka (API testlari haqiqiy Postgres'da):

```bash
docker compose up -d --wait
go test ./...
docker compose down -v
```

**Ishoralar**

1. `bg.Wait()`ning muddati yoʻq va uni toʻxtatib boʻlmaydi. Shuning uchun uni **alohida goroutine'da** kutib, tugaganda
   bitta channel'ni yoping (`close`, GB-7.02): yopilgan channel'dan qabul darhol qaytadi. Bu channel orqali qiymat
   yuborilmaydi, faqat yopiladi: uning turi `chan struct{}` (`struct{}` — maydonsiz, hech narsa saqlamaydigan tur),
   yaratish `done := make(chan struct{})`.
2. Endi ikkita narsani birga kutasiz: shu channel va `ctx.Done()`. Bu — `select` (GB-7.03): qaysi biri oldin tayyor
   boʻlsa, shuni oladi. `ctx.Done()` tarmogʻida `ctx.Err()` qaytaring.
3. Muddat tugaganda `bg.Wait()`ni kutayotgan goroutine qoladi: bu yerda bu normal, jarayon baribir chiqyapti (README,
   «waitBG'dan qolgan goroutine»).

Javob: [`solution/cmd/api/receipt.go`](solution/cmd/api/receipt.go).
