# Vazifa · GB-7.07: partiya timeout'i (`Config.BatchTimeout`)

Videodagi importer Ctrl+C'da toza toʻxtaydi. Lekin baza javob bermay qolsa (tarmoq uzildi, jadval boshqa transaction
tomonidan qulflangan), worker `SaveBatch` ichida **abadiy** kutadi: Ctrl+C bosmaguncha import tugamaydi, hisobot ham
chiqmaydi. Kechasi cron ishga tushirgan importda Ctrl+C bosadigan odam yoʻq.

`internal/importer/importer.go`dagi `Config`da `BatchTimeout` maydoni bor, lekin u hozir **ishlatilmaydi**. Uni ishlating.

**Spetsifikatsiya**

```
Config{Workers: 1, BatchTimeout: 5 * time.Second}, 3-partiya osilib qoladi →
  Run aynan 5 s dan keyin qaytadi
  err.Error() "partiya 3: " bilan boshlanadi
  errors.Is(err, context.DeadlineExceeded) == true
  (testdagi soxta saqlovchi bilan: "partiya 3: context deadline exceeded";
   haqiqiy Postgres'da pgx oʻraydi: "partiya 3: timeout: context deadline exceeded")
  hisob: yozildi 2 partiya · rollback 1 partiya
BatchTimeout: 0 → cheklov yoʻq (10 s davom etadigan partiya ham yoziladi)
BatchTimeout'dan qisqa partiyalar toʻxtatilmaydi
```

**Nima qilish kerak** (`importer.go`, `work` funksiyasi, `VAZIFA` izohi):

1. `timeout > 0` boʻlsa, har partiya **oʻz** muddatli ctx'ini olsin: `context.WithTimeout(ctx, timeout)` (GB-7.03) va
   `SaveBatch`ga oʻsha ctx. `cancel`ni partiya tugashi bilan chaqiring: `defer` bilan alohida yordamchi funksiyada
   (`save(ctx, s, rows, timeout) error`) yoki toʻgʻridan-toʻgʻri `SaveBatch`dan keyin. Sikl ichida `defer cancel()` yozmang:
   u sikl emas, butun `work` tugaganda ishlaydi.
2. Xato hozirgidek `partiya N: …` koʻrinishida, `%w` bilan qaytsin va rollback hisobiga qoʻshilsin (bu qism tayyor:
   `fmt.Errorf("partiya %d: %w", …)` — faqat buzib qoʻymang).
3. `timeout == 0` — cheklov yoʻq: `context.WithTimeout(ctx, 0)` emas (u darhol tugaydi), oddiy `ctx`.

Ixtiyoriy (testlar tekshirmaydi): `cmd/import/main.go`ga `-timeout` bayrogʻini qoʻshing:
`flag.Duration("timeout", 0, "…")` → `importer.Config{Workers: …, BatchTimeout: *timeout}`; `-timeout 10s` kabi yoziladi.

**Tekshirish** (shu papkada; Postgres testlari uchun baza kerak):

```bash
docker compose up -d --wait
go test ./...
```

`ok` boʻlguncha ishlang. Hozir shunday yiqiladi:

```
--- FAIL: TestTaskBatchTimeout (0.00s)
    tekshiruv_test.go:48: 3-partiya osilib qoldi va Run 1 daqiqada ham qaytmadi: BatchTimeout (5s) ishlatilmayapti (TASK.md, 1-qadam)
FAIL
FAIL	github.com/abdurahmon27/bozor/internal/importer	0.360s
```

Test `synctest` ichida ishlaydi (GB-7.06): «5 s» va «1 daqiqa» soxta soatda, test bir lahzada tugaydi.

**Maslahatlar**

- `WithTimeout` qaytargan xato `context.DeadlineExceeded`; `%w` bilan oʻralsa, `errors.Is` uni topadi (GB-2.08).
- Postgres saqlovchisi (`pg.go`) hech narsani oʻzgartirishi shart emas: `Exec(ctx, …)` muddat tugaganda toʻxtaydi,
  transaction rollback boʻladi. Commit va Rollback'ning oʻz chegarasi bor (`endTimeout`, 5 s), shuning uchun osilgan
  baza `Commit`da ham abadiy kuttirmaydi.
- `go vet ./...` `cancel` chaqirilmagan joyni koʻrsatadi (`lostcancel` tekshiruvi, GB-7.03).

Javob: [`solution/`](solution) (`save` funksiyasi va `-timeout` bayrogʻi). Tugagach: `docker compose down -v`.
