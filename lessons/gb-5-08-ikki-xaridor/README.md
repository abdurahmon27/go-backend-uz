# GB-5.08 · Ikki xaridor, bitta mahsulot: lost update va qatorni qulflash

**Go backend: noldan guru'gacha** kursining 40-darsi. Oxirgi bitta o'rik uchun ikki xaridor bir vaqtda «sotib olish»ni
bosadi. O'qi–tekshir–yoz ketma-ketligi `READ COMMITTED` transaction'da ikkalasiga ham «olindi» deydi (lost update, oversell).
Dars buni ikki psql seansida qadam-baqadam ko'rsatadi, Go testida 20 parallel xaridor bilan takrorlaydi va ikki usul bilan
tuzatadi: shartli `UPDATE ... WHERE stock >= $1` (+ `RowsAffected`) va `SELECT ... FOR UPDATE`.

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go bo'lsa,
> `go run`/`go test` birinchi marta kerakli versiyani o'zi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qo'ygan bo'lsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-5-07` (transaction: buyurtma qoldiqni kamaytiradi) · keyingi dars: `gb-5-09` (sqlc).
Bu papka mustaqil laboratoriya: oldingi darslarning kodi kerak emas (`bozor.sql` GB-5.07'dagi jadvallar va to'rt mahsulot,
qo'shimcha: `coupons` jadvali va `BAHOR` promo-kodi, vazifa uchun).
Dars tugatilgan `PlaceOrder`ga tegmaydi: u buyurtma qatorlarini yozish qoidasini GB-5.07'dan oladi. Bu laboratoriya faqat
**qoldiq qoidasini** (`Buy`) ajratib oladi, shuning uchun u `orders`/`order_items`ga yozmaydi.

## Papka

| fayl | nima |
|---|---|
| `compose.yaml` | Postgres 18 (`postgres:18-alpine`, kursning Module 5 bo'yicha bir xil tanlovi), loyiha nomi `gb508`, `127.0.0.1:5458`, nomli volume; Docker Compose 2.20+ kerak (`start_interval`) |
| `bozor.sql` | GB-5.07'dagi jadvallar va to'rt mahsulot + `coupons` (BAHOR: 3 ta ishlatish) |
| `store.go` | videoning oxiridagi holat: `Buy` (shartli UPDATE) va `BuyLocked` (FOR UPDATE). Ikkalasi `qty <= 0` bo'lsa `ErrInvalidQuantity` qaytaradi (GB-5.07'dagidek; videoda «Miqdor tekshiruvi» deb qisqartirilgan) |
| `naive.go`, `naive_test.go` | videodagi birinchi, XATO versiya `buyNaive` (o'qi, tekshir, yoz) va uni isbotlovchi test: `go test -run Naive -v` ortiqcha sotuvni ko'rsatadi. Faqat ko'rsatish uchun |
| `coupon.go` | VAZIFA: `UseCoupon` (xato versiya, parallel chaqiruvda buziladi) + ixtiyoriy `UseCouponLocked` |
| `store_test.go`, `coupon_test.go` | tekshiruv: har test alohida vaqtinchalik sxemada, 20 va 10 parallel chaqiruv, 5 takror |
| `solution/` | to'liq javob (o'z `go.mod`i bilan, alohida ishga tushadi) |
| `_video/` | videodagi buyruqlar: `chiqishlar.sh` → `chiqishlar.txt`, `sessions.py` (ikki psql seans), `stages/` (video bosqichlari) |

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, pgx v5.13.0, PostgreSQL 18.7).
To'liq jurnal: [`_video/chiqishlar.txt`](_video/chiqishlar.txt). Video `bozor` nomli papkada yozilgan (compose loyiha nomi = papka nomi,
konteyner `bozor-db-1`, terminalda `cwd: bozor`); bu papkadagi `compose.yaml` `name: gb508` beradi (`gb508-db-1`), buyruqlar bir xil.

**Ikki psql seansini o'zingiz ochish:** ikkita terminal oching va ikkalasida shu papkada `docker compose exec db psql` ishga tushiring
(Postgres'ga ikkita alohida ulanish). Buyruqlarni videodagi tartibda navbat bilan yuboring; `_video/sessions.py` shu ketma-ketlikni
avtomatik bajaradi (`python3 _video/sessions.py naive|shartli|forupdate`, compose papkasida).

```bash
docker compose up -d --wait
docker compose exec db psql -f /lab/bozor.sql
docker compose exec db psql -c "UPDATE products SET stock = 1 WHERE id = 4"     # UPDATE 1 (shartga mos qator yangilandi)
```

**1. Lost update, ikki seans (A va B).** Ikkalasi `BEGIN`, ikkalasi `SELECT stock ... WHERE id = 4` → `1`, A `UPDATE ... SET stock = 0`
va `COMMIT`, keyin B ham `UPDATE ... SET stock = 0` → `UPDATE 1` (xato yo'q) va `COMMIT`. Ikkala «xarid» o'tdi, qoldiq `0`.
`SHOW transaction_isolation` → `read committed` (Postgres'ning standarti).

**2. Shartli UPDATE.** A: `BEGIN; UPDATE products SET stock = stock - 1 WHERE id = 4 AND stock >= 1;` → `UPDATE 1` (COMMIT yo'q).
B: shu buyruq, javob **kelmaydi** (A qatorni band qilgan). A `COMMIT` qiladi, B javobi: `UPDATE 0` (Postgres B'ning `WHERE`ini
qatorning yangi qiymatiga qarab qayta tekshirdi). Manba: [Read Committed](https://www.postgresql.org/docs/18/transaction-iso.html#XACT-READ-COMMITTED).

**3. FOR UPDATE.** A: `BEGIN; SELECT stock ... FOR UPDATE` → `1`. B: `BEGIN; SELECT stock ... FOR UPDATE` javobsiz kutadi.
A `UPDATE ... SET stock = 0; COMMIT`; B javobi: `0`. Dars bu seansni diagramma bilan ko'rsatadi, to'liq jurnal
`_video/chiqishlar.txt`da (`seans-forupdate`).

**4. Go testi** (`go test`, 20 xaridor, omborda 5 ta). Birinchi versiyaning yiqilishini o'zingiz ko'rish uchun `store.go`dagi `Buy`
tanasini `naive.go`dagi `buyNaive` tanasi bilan almashtiring va `go test` ni ishga tushiring (yoki `go test -run Naive -v`: ortiqcha
sotuvni sanab beradi):

```
birinchi versiya (o'qi, tekshir, yoz)        → FAIL: 20 tasiga "olindi" dedik; kutilgan: 5
shartli UPDATE (Buy)                          → ok
FOR UPDATE (BuyLocked)                        → ok
FOR UPDATE, lekin pool.QueryRow orqali (xato) → FAIL: 19–20 tasiga "olindi" (qulf darhol ochiladi)
```

Xato versiyada «20» yoki «19» chiqadi: parallel testda aniq son har safar boshqacha, lekin 5 emas. Videodagi chiqishlar bosqichma-bosqich
`_video/stages/` fayllari bilan olingan (`01-store.go` birinchi `Buy`, `02-store.go` shartli `Buy`, `04-store-pool-mistake.go` xato).

Bosqich fayllarini `store.go` o'rniga qo'yib sinasangiz, `store_test.go`dagi `{"BuyLocked", ...}` qatorini bo'sh qoldiring
(`_video/chiqishlar.sh` shuni qiladi; qatorni o'chirmang, bo'shatib qo'ying: xabardagi satr raqami o'zgarmaydi).

## Chuqurroq (videoga kirmagan)

- **`CHECK (stock >= 0)` oxirgi to'siq.** `UPDATE ... SET stock = stock - 1` (shartsiz) qoldiq `0` bo'lganda xato beradi:
  `ERROR: new row for relation "products" violates check constraint "products_stock_check"` (SQLSTATE 23514). GB-5.07'dagi
  `PlaceOrder` (`stock = stock - $1`) shu xatoni `ErrNotEnoughStock`ga aylantiradi, ya'ni ortiqcha sotmaydi; bu dars esa xatoga
  yetkazmasdan qoidani o'zi tekshiradi (shartli `UPDATE` yoki `FOR UPDATE`). `CHECK` baribir kerak: u ilova xatosiga qarshi oxirgi himoya.
- **Optimistik lock** (`version` ustuni; bozor.sql'da bu ustun yo'q, bu misol): `UPDATE ... SET stock = $1, version = version + 1 WHERE id = $2 AND version = $3`;
  `RowsAffected() == 0` bo'lsa, qayta o'qib qayta urinish. Uzoq qulf yo'q (UPDATE qatorni faqat bir lahza qulflaydi), lekin takrorlash kerak; to'qnashuv kam bo'lsa mos.
  Umumiy nom: optimistic offline lock (https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html). Kursda keyinroq (GB-9) ko'riladi.
- **Bir nechta qatorni qulflash** (savatda ikki mahsulot): har doim bir xil tartibda (masalan `ORDER BY id`) qulflang, aks holda
  ikki transaction bir-birini kutib qotadi (deadlock; Postgres bittasini xato bilan to'xtatadi). Manba: [Deadlocks](https://www.postgresql.org/docs/18/explicit-locking.html#LOCKING-DEADLOCKS).
- **Yuqoriroq izolyatsiya.** `REPEATABLE READ`/`SERIALIZABLE`da Postgres bunday to'qnashuvni `could not serialize access`
  (SQLSTATE 40001) xatosi bilan to'xtatadi va transaction'ni qayta urinish kerak. Biz standart `READ COMMITTED` + aniq qoida bilan ishlaymiz.
- **Pool tuzog'i.** `BuyLocked`ning xato versiyasida (`_video/stages/04-…`) `Begin` bitta ulanishni oladi, `s.pool.QueryRow` ikkinchisini so'raydi.
  Pool kichik bo'lsa (20 xaridor, 20 ulanish), hamma birinchi ulanishni ushlab, ikkinchisini kutishi mumkin: biz buni 20 ulanishli pool'da
  `context deadline exceeded` bilan ko'rdik (GB-5.05'dagi `rows.Close` tuzog'iga o'xshash). Shuning uchun test pool'i 40 ulanishli (`store_test.go`);
  ishlaydigan `Buy`/`BuyLocked` bitta ulanish oladi.

## Vazifa

[`TASK.md`](TASK.md): `UseCoupon`ni shartli `UPDATE` bilan tuzating (🟡 `UseCouponLocked` ixtiyoriy). Tekshirish:

```bash
docker compose up -d --wait
go test .
```

Tugagach `docker compose down -v`.

## Manbalar

- PostgreSQL, tranzaksiya izolyatsiyasi, Read Committed (UPDATE kutadi va WHERE qayta tekshiriladi):
  https://www.postgresql.org/docs/18/transaction-iso.html#XACT-READ-COMMITTED
- PostgreSQL, qator qulflari (`FOR UPDATE`): https://www.postgresql.org/docs/18/explicit-locking.html#LOCKING-ROWS ·
  `SELECT` buyrug'i, locking clause: https://www.postgresql.org/docs/18/sql-select.html#SQL-FOR-UPDATE-SHARE
- PostgreSQL, `UPDATE`: https://www.postgresql.org/docs/18/sql-update.html
- pgx, `CommandTag.RowsAffected`: https://pkg.go.dev/github.com/jackc/pgx/v5/pgconn#CommandTag.RowsAffected ·
  `Tx`: https://pkg.go.dev/github.com/jackc/pgx/v5#Tx
- Go, `sync.WaitGroup.Go`: https://pkg.go.dev/sync#WaitGroup.Go · `sync.Mutex`: https://pkg.go.dev/sync#Mutex
