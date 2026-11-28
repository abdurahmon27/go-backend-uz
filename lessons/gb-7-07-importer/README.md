# GB-7.07 · Nazorat nuqtasi: bekor qilinadigan importer

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #57"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

7-modulning nazorat nuqtasi. `cmd/import` — katta katalog faylini (200 000 qator) bir nechta worker bilan Postgres'ga
yozadigan CLI. Modulning hamma qismi bitta dasturda: goroutine'lar va ularning egasi (GB-7.01), channel va pipeline
(GB-7.02), `select` va context bilan bekor qilish (GB-7.03), `errgroup` (GB-7.04), `atomic` hisoblagichlar (GB-7.05),
`synctest` bilan test (GB-7.06). Baza — GB-5.10'dagi sxema (`products`), Postgres 18 Docker'da.

Oldingi darslar: [`gb-7-06-leak-synctest`](../gb-7-06-leak-synctest) · [`gb-7-04-worker-errgroup`](../gb-7-04-worker-errgroup)
(errgroup; uning `saveStocks` vaʼdasi shu yerda `PG.SaveBatch` sifatida: import Postgres'ga yoziladi) · [`gb-7-03-select-context`](../gb-7-03-select-context) · [`gb-5-10-buyurtma-api`](../gb-5-10-buyurtma-api)
(sxema, `pgtest`) · keyingi dars: `gb-8-01` (config va secret'lar).
Bu papka mustaqil laboratoriya: 7-modul darslari parallel yozilgan, har biri oʻz `package main`ida. GB-7.04'dagi
`importStocks` (yetkazib beruvchidan qoldiq soʻrash) bu yerda yoʻq: importer faylni oʻqiydi. GB-5.10 bilan umumiy qismlar:
`migrations/` (aynan oʻsha uchta fayl), `migrate.go`, `internal/pgtest` (bu yerda seed'siz: jadvallar boʻsh).
Farqni koʻrish: `diff -ru ../gb-5-10-buyurtma-api/migrations migrations`.

## Spetsifikatsiya

```
go run ./cmd/import [-workers 4] katta.txt
qator        nom;narx;qoldiq   narx soʻmda → tiyin
partiya      500 qator = bitta transaction
nom bor      narx va qoldiq yangilanadi (INSERT … ON CONFLICT (name) DO UPDATE)
Ctrl+C       yangi partiya yoʻq, yarim partiya rollback (Exec'lari tugagan partiya commit boʻladi)
buzuq qator  hammasi toʻxtaydi, qator raqami bilan (undan oldin commit boʻlganlar bazada qoladi)
hisobot      yozildi = count(*)  (boʻsh jadvalga, nomlari takrorlanmaydigan fayl bilan)
```

«yozildi» — commit boʻlgan partiyalardagi qatorlar: qoʻshilganlar ham, upsert bilan yangilanganlar ham. Shuning uchun
`count(*)` bilan tenglik boʻsh jadvalda va nomlari takrorlanmaydigan faylda toʻgʻri (videodagi birinchi import); qayta
importda «yozildi» 200 000, jadvalda ham 200 000, lekin ularning bir qismi yangilangan qatorlar.

Qator formati GB-3.01'dagi `ID;nom;narx;qoldiq`ning ID'siz varianti: `products.id` `GENERATED ALWAYS AS IDENTITY`
(GB-5.03), ID'ni baza beradi, importda kalit — `name` (`UNIQUE`). Boʻsh qatorlar va `#` bilan boshlangan izohlar
tashlab oʻtiladi (GB-3.01 kabi).

## Qanday ishlaydi

```
katta.txt → read (1 goroutine) → batches (buffersiz channel) → work × N → Saver (Postgres: partiya = transaction)
            └──────────────── hammasi bitta errgroup'da, bitta ctx bilan ────────────────┘
```

- **Ctrl+C → ctx.** `signal.NotifyContext(context.Background(), os.Interrupt)`: SIGINT kelganda dastur oʻlmaydi, ctx
  bekor boʻladi. Uni `read` (`send` ichidagi `select`), har `work` va har `Exec` koʻradi; hamma qaytadi, `Wait` tugaydi,
  hisobot chiqadi. Exit code 1, xabar `import: toʻxtatildi (Ctrl+C)`.
- **Birinchi xato hammasini toʻxtatadi.** `errgroup.WithContext`: buzuq qator yoki baza xatosi `gctx`ni bekor qiladi.
- **Partiya yo toʻliq, yo yoʻq.** `PG.SaveBatch`: `Begin` → har qatorga `Exec(ctx, upsert, …)` → `Commit`. `Exec` ctx'ni
  oladi (Ctrl+C uni toʻxtatadi → rollback). `Commit` va `Rollback` esa alohida `end` ctx bilan:
  `context.WithTimeout(context.WithoutCancel(ctx), endTimeout)` — Ctrl+C'ni koʻrmaydi (Exec'lari tugagan partiyaning
  Commit'i yarmida uzilmaydi), lekin 5 s dan ortiq ham kutmaydi (baza yoki tarmoq javob bermasa). Hisob (`Stats`, `atomic.Int64`) faqat natija maʼlum boʻlgach oshadi,
  shuning uchun hisobotdagi «yozildi» bazadagi qatorlar soniga teng.
- **Leak yoʻq.** `Run` faqat hamma goroutine tugagach qaytadi (`g.Wait`); `read` hech qachon `select`siz yubormaydi.
  `Stats` `Run` qaytgandan keyin oʻqiladi (`Wait` hamma yozuvni koʻrsatadi: [Go memory model](https://go.dev/ref/mem)).

## Papka

| fayl | nima |
|---|---|
| `compose.yaml` | Postgres 18 (`postgres:18-alpine`), loyiha nomi `gb707`, `127.0.0.1:5517`, nomli volume `pgdata` |
| `migrations/`, `migrate.go` | GB-5.10'dagi uchta goose migratsiyasi (`products`, `orders` + `order_items`, `orders.status`) |
| `cmd/import/` | CLI: `-workers`, `signal.NotifyContext`, `pgxpool`, hisobot |
| `cmd/katalog/` | sinov fayli: `go run ./cmd/katalog 200000 > katta.txt` (birinchi qatorlar: olma, non, oʻrik, choy) |
| `internal/importer/` | `Product`, `parseLine` (`product.go`), `Run`/`read`/`send`/`work`, `Stats` (`importer.go`), `PG` saqlovchi (`pg.go`), `Generate` |
| `internal/importer/*_test.go` | `importer_test.go` (synctest: bekor qilish, buzuq qator), `pg_test.go` (haqiqiy Postgres), `tekshiruv_test.go` (VAZIFA) |
| `internal/pgtest/` | `pgtest.New(t)`: har test oʻz sxemasida (GB-5.10) |
| `solution/` | toʻliq javob (oʻz `go.mod`i bilan) |
| `_video/` | `chiqishlar.sh` → `chiqishlar.txt` (videodagi buyruqlar), `race.sh` → `race.txt` (`-race`, Docker), `ctrlc.py`, `stages/importer-xato.go` (videodagi XATO versiya) |

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, pgx v5.13.0, goose v3.30.2, PostgreSQL 18,
golang.org/x/sync v0.25.0). Toʻliq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt). Ctrl+C'ni
[`_video/ctrlc.py`](_video/ctrlc.py) bosadi: 2.5 s dan keyin butun jarayon guruhiga SIGINT (terminal ham shunday qiladi).
Vaqtlar mashina yukiga bogʻliq; sanoqlar (`yozildi` = `count(*)`) esa har safar bir-biriga teng.

```bash
docker compose up -d --wait
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations
export GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5517/bozor
export DATABASE_URL=$GOOSE_DBSTRING
goose up
# OK   00001_products.sql … goose: successfully migrated database to version: 3
go run ./cmd/katalog 200000 > katta.txt
head -4 katta.txt
# olma;12000;40 · non;4000;120 · oʻrik;35000;5 · choy;25000;8

go run ./cmd/import katta.txt        # 2.5 s dan keyin Ctrl+C
# ^Cyozildi: 48000 qator (96 partiya) · rollback: 2000 qator (4 partiya) · vaqt: 2.182s
# import: toʻxtatildi (Ctrl+C)
# exit status 1
docker compose exec db psql -c "SELECT count(*) FROM products"
# 48000                                ← hisobotdagi «yozildi» bilan teng: yarim partiya yoʻq

go run ./cmd/import katta.txt        # qayta, oxirigacha
# yozildi: 200000 qator (400 partiya) · rollback: 0 qator (0 partiya) · vaqt: 8.896s
docker compose exec db psql -c "SELECT count(*) FROM products"
# 200000                               ← 48000 tasi yangilandi (ON CONFLICT), qolgani qoʻshildi; fayl boshidan oʻqildi

sed '120000s/;[0-9]*;/;-5;/' katta.txt > buzuq.txt
sed -n 120000p buzuq.txt
# sut-120000;-5;199
go run ./cmd/import buzuq.txt
# yozildi: 117500 qator (235 partiya) · rollback: 2000 qator (4 partiya) · vaqt: 5.364s
# import: 120000-qator: narx: narx musbat boʻlishi kerak, keldi -5
```

Testlar (`go test -v -run 'TestCancel|TestBadLine|TestPostgres' ./internal/importer`): hammasi `PASS`;
`TestPostgresCancel` logi (raqamlar har safar boshqacha, tenglik doim): `yozildi: 2000 qator (4 partiya) · rollback: 1000 qator (2 partiya) · bazada: 2000 qator`.
`go test ./...` hozir vazifa tekshiruvi tufayli `FAIL` (pastda).

**`-race`** gcc talab qiladi; bu mashinada gcc yoʻq, shuning uchun [`_video/race.sh`](_video/race.sh) Docker'da
(`golang:1.26.8` image'i + shu mashinadagi go1.27.4 toolchain, `--network host`) ishga tushiradi. Natija:
[`_video/race.txt`](_video/race.txt) — testlar (papkadagi kod va `solution/`) `ok`, `go build -race` qilingan CLI
Ctrl+C bilan (60 000 yozildi, 4 partiya rollback) va 8 worker bilan oxirigacha: `WARNING: DATA RACE` yoʻq. Sizda gcc boʻlsa: `go test -race ./...`.

## Klassik xato: `select`siz yuborish

[`_video/stages/importer-xato.go`](_video/stages/importer-xato.go): `read` partiyani `send(ctx, out, b)` oʻrniga oddiy
`out <- b` bilan beradi. Oxirigacha boradigan import ishlaydi. Ctrl+C'da esa worker'lar qaytadi, `read` channel'ga
yuborishda abadiy turadi, `g.Wait` uni kutadi: **dastur hisobot bermaydi va chiqmaydi** (skriptda 8 s kutildi, keyin
`KILL`). Papkadagi kodda `context.AfterFunc(ctx, stop)` bor: birinchi Ctrl+C'dan keyin ushlash oʻchadi, ikkinchi Ctrl+C
dasturni odatdagidek darhol toʻxtatadi; usiz ikkinchi Ctrl+C ham yutilardi (tekshirilgan). `synctest` buni darhol koʻradi:

```
--- FAIL: TestCancel (0.00s)
panic: deadlock: all goroutines in bubble are blocked [recovered, repanicked]
goroutine 19 [sync.WaitGroup.Wait (durable), synctest bubble 1]:     ← Run, errgroup'ning Wait'i
goroutine 21 [chan send (durable), synctest bubble 1]:               ← read, out <- b (importer.go:107)
FAIL	github.com/abdurahmon27/bozor/internal/importer	0.007s
```

## Chuqurroq (videoga kirmagan)

- **Worker soni.** Boʻsh jadvalga 200 000 qator, shu mashinada (`chiqishlar.txt` 8-boʻlim):
  1 worker 22.3 s · 4 → 8.1 s · 8 → 5.5 s · 12 → 5.6 s. Boshqa oʻlchovlarda: 43.3 · 13.5 · 8.8 · 9.5 s (mashina band) va 14.3 · 5.9 · 4.2 · 4.3 s.
  Nisbat bir xil: 8 dan keyin tezlashmaydi. Chegara endi Go tomonida emas, Postgres tomonida (taxmin: bitta `UNIQUE`
  indeks va WAL yozuvi; oʻlchanmagan). Worker'lar pool
  ulanishlaridan (`pgxpool` standarti: `max(4, CPU soni)`, [pgxpool.Config](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#Config))
  koʻp boʻlsa, ortiqchasi ulanish kutadi; kerak boʻlsa `DATABASE_URL`ga `?pool_max_conns=…` (GB-5.05).
- **Tezroq yozish.** Qatorma-qator `Exec` oddiy va tushunarli, lekin eng tez yoʻl emas: `pgx.Batch` (bitta round trip),
  `COPY` (`pool.CopyFrom`) va vaqtinchalik jadval + bitta `INSERT … SELECT … ON CONFLICT` — round trip'lar soni
  kamayadi (bu laboratoriyada oʻlchanmagan; oʻlchash — GB-11 mavzusi).
  Bekor qilish va «partiya = transaction» qoidasi oʻzgarmaydi.
- **Fayl yo butunlay, yo umuman.** Buzuq qatorgacha commit boʻlgan partiyalar bazada qoladi (videoda 117 500 qator).
  Butun fayl bitta qadam boʻlishi kerak boʻlsa: avval butun faylni tekshiring (ikki oʻtish), yoki hammasini bitta
  transaction'da yozing (uzoq qulf, katta WAL), yoki vaqtinchalik jadvalga yuklab, oxirida bitta transaction'da koʻchiring.
- **Ikkinchi Ctrl+C.** `stop()` chaqirilmaguncha `NotifyContext` SIGINT'ni ushlaydi: ikkinchi Ctrl+C dasturni
  oʻldirmasdi. Shuning uchun `run`da `context.AfterFunc(ctx, stop)`: ctx bekor boʻlishi bilan `stop()` chaqiriladi va
  standart xatti-harakat qaytadi ([os/signal: NotifyContext](https://pkg.go.dev/os/signal#NotifyContext),
  [context.AfterFunc](https://pkg.go.dev/context#AfterFunc)). Graceful shutdown serverda — GB-8.03.
- **Hisobning chegaralari.** `Begin` xatosi ham «rollback» hisobiga tushadi (unda hech narsa yozilmagan). Commit
  bazada bajarilib, javobi tarmoqda yoʻqolsa, partiya «rollback» deb sanaladi, lekin bazada boʻladi; upsert tufayli
  qayta import buni tuzatadi.
- **ID'lardagi boʻshliqlar.** `ON CONFLICT … DO UPDATE` qatorni yangilasa ham identity ketma-ketligi bitta qiymatni
  sarflaydi: 200 000 qatorni qayta import qilish 200 000 ta ID'ni «yeydi». ID'lar boʻshliqsiz boʻlishi kafolatlanmaydi
  (GB-5.10 Chuqurroq); `bigint`da bu muammo emas.
- **`Commit` va bekor qilingan ctx.** `Commit(ctx)` bekor qilingan ctx bilan yuborilmaydi; pgx `Rollback` xatosida ulanishni
  yopadi ([pgx: Tx](https://pkg.go.dev/github.com/jackc/pgx/v5#Tx)). GB-5.07'dagi `defer tx.Rollback(ctx)` soʻrov
  ctx'i uchun yetarli edi; bu yerda ctx ataylab bekor qilinadi, shuning uchun ikkalasi `end` ctx bilan.
  `context.WithoutCancel` muddatni ham olib tashlaydi ([context.WithoutCancel](https://pkg.go.dev/context#WithoutCancel),
  Go 1.21), shuning uchun ustidan `WithTimeout(…, endTimeout)`. Exec paytida ctx bekor boʻlsa, pgx
  standart holatda kutishni toʻxtatib ulanishni yopadi (`DeadlineContextWatcherHandler`), Postgres esa uzilgan ulanishning
  ochiq transaction'ini bekor qiladi.
- **`SetLimit` nega yoʻq.** Goroutine'lar soni oldindan maʼlum (`Workers` + bitta `read`), ularning har biri channel
  yopilguncha ishlaydi (GB-7.04'dagi worker pool). `SetLimit` har ish uchun yangi goroutine ochiladigan joyda kerak.

## Vazifa

[`TASK.md`](TASK.md): `Config.BatchTimeout` — osilib qolgan partiya `partiya N: context deadline exceeded` bilan
toʻxtasin. Tekshirish:

```bash
docker compose up -d --wait
go test ./...
```

Hozir `TestTaskBatchTimeout` yiqiladi: `3-partiya osilib qoldi va Run 1 daqiqada ham qaytmadi: BatchTimeout (5s) ishlatilmayapti`.
Tugagach: `docker compose down -v` (konteyner va `pgdata` volume oʻchadi). Javob: [`solution/`](solution).

## Manbalar

- Go: [os/signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext) · [context](https://pkg.go.dev/context)
  ([WithoutCancel](https://pkg.go.dev/context#WithoutCancel), [WithTimeout](https://pkg.go.dev/context#WithTimeout)) ·
  [flag](https://pkg.go.dev/flag) · [sync/atomic.Int64](https://pkg.go.dev/sync/atomic#Int64) ·
  [testing/synctest](https://pkg.go.dev/testing/synctest) · [Select statements](https://go.dev/ref/spec#Select_statements) ·
  [The Go Memory Model](https://go.dev/ref/mem) · [Race detector](https://go.dev/doc/articles/race_detector) ·
  [Go blog: Pipelines and cancellation](https://go.dev/blog/pipelines)
- golang.org/x/sync: [errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)
- pgx v5: [Tx](https://pkg.go.dev/github.com/jackc/pgx/v5#Tx) · [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool)
- PostgreSQL 18: [INSERT … ON CONFLICT](https://www.postgresql.org/docs/18/sql-insert.html#SQL-ON-CONFLICT) ·
  [Transactions](https://www.postgresql.org/docs/18/tutorial-transactions.html) ·
  [Identity columns](https://www.postgresql.org/docs/18/ddl-identity-columns.html)
- POSIX: [SIGINT](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/signal.h.html) (terminal Ctrl+C'da butun jarayon guruhiga yuboradi)
