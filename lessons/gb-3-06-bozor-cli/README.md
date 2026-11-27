# GB-3.06 · Nazorat nuqtasi: bozor CLI v0 (add, list, order, JSON, testlar)

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #24"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** 2–3-modul darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. `range` funksiya ustida (`st.All()`) Go
> 1.23'dan beri bor.

Oldingi dars: [`gb-3-05-iterator`](../gb-3-05-iterator) (`iter.Seq`, `All`) · keyingi dars:
[`gb-4-01-http-matn`](../gb-4-01-http-matn) (HTTP soʻrov — bu matn).

## Bu papka — 4-modulning boshlangʻich nuqtasi

2–3-modul darslari mustaqil lab'lar edi. Bu papka ularni **bitta bozor dasturiga** yigʻadi va kurs shu yerdan
oʻsadi: 4-modul (`net/http`, `cmd/api`) xuddi shu `internal/store`'ni HTTP orqali beradi. Keyingi darslar
boshlanadigan kod — shu papkaning `solution/`'i (vazifa bajarilgan holat).

**Bu darsdan keyin:** `cmd/bozor` + `internal/store` tuzilishidagi CLI'ni oʻqiysiz va kengaytirasiz; `run(path,
args, w) error` nega butun CLI'ni test qilinadigan qilishini bilasiz; chiqish kodi 1 (xato) va 2 (notoʻgʻri
chaqiruv) farqini qoʻyasiz; `os.WriteFile` nega yarim fayl qoldirishi mumkinligini va `os.CreateTemp` + `Sync` +
`os.Rename` bilan atomar saqlashni bilasiz; boʻsh katalog `null` emas `[]` boʻlib yozilishi kerakligini va fayldan
oʻqilgan mahsulot tekshirilishini koʻrasiz; `go vet` va `staticcheck` jim qolgan haqiqiy bug'ni (range nusxasi)
jadvalli test ushlashini koʻrasiz.

## Nima qayerdan keldi

| Qism | Dars |
|---|---|
| `cmd/bozor` + `internal/store`, `Store`, kichik harfli `products` | GB-2.06 |
| `Product{ID, Name, Price, Stock}`, narx tiyinda (`tiyinPerSom`) | GB-2.03, GB-1.03 |
| `find` → `&s.products[i]` (pointer) | GB-2.04 |
| `ErrNotFound`, `ErrBadQuantity`, `*StockError`, `errors.Is/As`, `message(err)` | GB-2.08 |
| `run() error` + `main`'da `os.Exit`, nomli natija `err` + `defer` | GB-2.09 |
| `io.Writer`, `fmt.Fprintf`, testda `strings.Builder` | GB-2.07 |
| `validate`: nom boʻsh emas, UTF-8, ≤ 20 belgi (`maxNameRunes`), narx > 0, qoldiq ≥ 0 | GB-3.01 `parseLine`, GB-2.10 `checkName` |
| JSON tag'lar, `MarshalIndent`, `os.ReadFile`, fayl yoʻq → boʻsh katalog | GB-3.02 `saveProducts`/`loadProducts` |
| `All() iter.Seq[Product]` | GB-3.05 |
| `parsePrice` (soʻm → tiyin, chegaralar), `os.Args`, stderr, chiqish kodi | GB-1.07 |
| jadvalli testlar, `t.Run` | GB-1.08 |
| `go vet ./...`, `staticcheck ./...` toza | GB-3.03 |

Birlashtirishdagi qarorlar (lab'lar orasidagi farqlar):

- **Bitta ombor — `bozor.json`.** GB-3.01'dagi `katalog.txt` + `loadCatalog` matn formatidan voz kechildi; uning
  tekshiruv qoidalari `validate`'ga koʻchdi va `Load` har mahsulotni shundan oʻtkazadi. GB-3.02'dagi
  `loadProducts`/`saveProducts` endi `store.Load(path)` va `(*Store).Save(path)`.
- **Bitta katalog**: olma (12 000 soʻm, 40), non (4 000, 3), choy (25 000, 8), oʻrik (35 000, 2) — GB-3.05'dagi,
  narxlar GB-3.01'dagi. Lab'larda qoldiqlar har xil edi (masalan, GB-3.04: non 3, choy 0).
- `Product`'da faqat kursning toʻrt maydoni bor; lab'lardagi `Note` (GB-3.02) va `OldPrice` (GB-3.03) oʻsha
  darslarning misoli edi, CLI v0'ga kerak emas. `golang.org/x/text` ham kerak emas (narx `%d soʻm`), shuning uchun
  `go.mod`'da bogʻliqlik yoʻq.
- `Add` endi `(Product, error)` qaytaradi (qoidaga zid mahsulot qoʻshilmaydi), ID = eng katta ID + 1;
  `Sell` endi `(Product, error)` qaytaradi (sotuvdan keyingi mahsulot: `order` narx va qoldiqni chiqaradi).
- `List()` + `slices.Clone` olib tashlandi: CLI faqat `All()` bilan aylanadi. Lab'lardagi boshqa qismlar ham CLI
  v0'ga kerak emas, shuning uchun koʻchirilmadi: `Filter`/`Sum` (GB-3.04), `Names`/`LowStock`/`SortedNames`
  (GB-3.05), `Rename` (GB-3.03), `Order`/`OrderItem`/`OrderStatus` va metodlari (GB-2.03/2.05). `order` buyrugʻi
  hozircha bitta mahsulot sotadi va buyurtmani saqlamaydi; buyurtmalar ombori 5-modulda (Postgres) keladi.
- **Chiqish kodlari:** GB-1.07'da har qanday xato 1 edi. CLI v0 notoʻgʻri chaqiruvni (`errUsage`: nomaʼlum buyruq,
  argumentlar soni) **2** bilan ajratadi, qolgan xatolar 1. Go'ning `flag` paketi ham (`ExitOnError` → `os.Exit(2)`),
  `go` buyrugʻi ham (`go foo` → 2) shunday qiladi.
- **`message`** (GB-2.08): `ErrNotFound` → `xato: bunday mahsulot yoʻq`, `ErrBadQuantity` → `xato: dona soni
  notoʻgʻri` (GB-2.08'dagi matn), `*StockError` → `xato: oʻrik: faqat 2 dona qoldi`, qolganlari `xato: ` + xato matni;
  foydalanish xatosi oʻz matni bilan. GB-2.08'dan farqi — hamma xabar `xato:` bilan boshlanadi.
- **Avval saqlash, keyin chiqarish:** `add`/`order`(/`restock`) natija qatorini qaytaradi (`(string, error)`),
  `run` avval `Save` qiladi, faqat muvaffaqiyatda qatorni `w`'ga yozadi: saqlanmagan buyurtma «bajarildi» deb
  chiqmaydi (`TestNoOutputWithoutSave`).

## Papka

| Fayl | Nima |
|---|---|
| `cmd/bozor/main.go` | `main`, `run(path, args, w)`, `add`, `list`, `order` (natija qatorini qaytaradi), `message`, `parsePrice` |
| `cmd/bozor/main_test.go` | CLI testlari: `run`'ni `t.TempDir()`dagi `bozor.json` bilan chaqiradi |
| `internal/store/store.go` | `Product`, xatolar, `validate`, `Store`, `Add`, `All`, `find`, `Sell`; `Restock` — vazifaning boshlanishi |
| `internal/store/file.go` | `Load(path)`, `(*Store).Save(path)` (atomar) |
| `internal/store/store_test.go` | `TestSell`, `TestAdd`, `TestLoad` (jadvalli), `TestSellStockError`, `TestSaveLoad`, `TestSaveKeepsOld` (Save xato bersa eski fayl butun), `TestSaveEmpty` |
| `*/tekshiruv_test.go` | vazifani tekshiradi; oʻzgartirmang (`solution/internal/store/restock_test.go` — vazifaning 2-qadami, oʻz testingiz) |
| `TASK.md` | vazifa: `bozor restock` |
| `solution/` | javob, oʻz `go.mod`'i bilan: `cd solution && go test ./...` |
| `_video/01-nusxa/` | videodagi bug: `Sell` range nusxasiga yozadi |
| `_video/02-null/` | `Save` nil tekshiruvisiz: boʻsh katalog `null` boʻlib yoziladi |
| `_video/run.sh` | qadamni yangi vaqtinchalik `bozor/` papkasida ishga tushiradi; `final` = shu papka, `tekshiruv_test.go`'siz (videoda vazifa hali yoʻq); oxirida `(chiqish kodi: N)` (quyida tushirilgan) |

`_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi.

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

```
$ ./_video/run.sh final bash      # keyingi buyruqlar shu vaqtinchalik bozor/ papkasida
$ ls cmd/bozor internal/store
cmd/bozor:
main.go  main_test.go

internal/store:
file.go  store.go  store_test.go
$ go build ./cmd/bozor
$ ./bozor list
katalog boʻsh: bozor add NOM NARX QOLDIQ
$ ./bozor add olma 12000 40
qoʻshildi: 1. olma, 12000 soʻm, 40 dona
$ ./bozor add non 4000 3
qoʻshildi: 2. non, 4000 soʻm, 3 dona
$ ./bozor add choy 25000 8
qoʻshildi: 3. choy, 25000 soʻm, 8 dona
$ ./bozor add oʻrik 35000 2
qoʻshildi: 4. oʻrik, 35000 soʻm, 2 dona
$ ./bozor list
 1. olma                    12000 soʻm   40 dona
 2. non                      4000 soʻm    3 dona
 3. choy                    25000 soʻm    8 dona
 4. oʻrik                   35000 soʻm    2 dona
$ head -6 bozor.json
[
  {
    "id": 1,
    "name": "olma",
    "price": 1200000,
    "stock": 40
$ ./bozor order 1 3
buyurtma: olma × 3 = 36000 soʻm, qoldiq 37
$ ./bozor list
 1. olma                    12000 soʻm   37 dona
 …
$ ./bozor order 4 5
xato: oʻrik: faqat 2 dona qoldi
$ echo $?
1
$ ./bozor sell 1 1
foydalanish: bozor add NOM NARX QOLDIQ | list | order ID DONA
$ echo $?
2
```

`oʻrik` qatori `olma` bilan bir chiziqda: `fmt`'da kenglik (`%-20s`) baytlar emas, rune'lar bilan sanaladi
(«Width and precision are measured in units of Unicode code points, that is, runes», pkg.go.dev/fmt). `oʻrik` —
5 rune, 6 bayt (GB-2.10). Go'ning oʻzi ham notoʻgʻri chaqiruvda 2 qaytaradi: `go foo` → `go foo: unknown command`,
`echo $?` → `2`.

Yozish oʻrtasida toʻxtagan `os.WriteFile`'ni qoʻlda yasaymiz (`truncate -s 60` faylni 60 baytgacha kesadi):

```
$ truncate -s 60 bozor.json
$ ./bozor list
xato: bozor.json: unexpected end of JSON input
```

`Load` fayldan oʻqilgan mahsulotni ham tekshiradi (oʻrik'ning qoldigʻini qoʻlda `-5` qilsak):

```
$ sed -i 's/"stock": 2$/"stock": -5/' bozor.json
$ ./bozor list
xato: bozor.json: 4-mahsulot: mahsulot notoʻgʻri: qoldiq manfiy boʻlmasin, keldi -5
```

Testlar (`final`, vazifa tekshiruvlarisiz):

```
$ go test ./...
ok  	github.com/abdurahmon27/bozor/cmd/bozor	0.025s
ok  	github.com/abdurahmon27/bozor/internal/store	0.009s
```

Boʻsh katalog: `Save`'dan nil tekshiruvi olib tashlansa (`_video/02-null`):

```
$ ./_video/run.sh 02-null go test ./internal/store
--- FAIL: TestSaveEmpty (0.00s)
    store_test.go:157: boʻsh katalog "null" boʻlib yozildi, kutilgan []
FAIL
FAIL	github.com/abdurahmon27/bozor/internal/store	0.013s
FAIL
```

Yangi `encoding/json/v2` paketi nil slice'ni ham `[]` qilib yozadi (`json.Marshal([]int(nil))`: v1 → `null`,
v2 → `[]`, go1.27.4'da tekshirildi); kurs `encoding/json` (v1) bilan ishlaydi, farqlar GB-3.02'ning Chuqurroq
qismida.

Bug: jamoadosh `Sell`'ni `for _, p := range s.products` bilan «soddalashtirdi» (`_video/01-nusxa`). `go vet` ham,
`staticcheck` (2026.2.2) ham jim; buyurtma xabari toʻgʻri, katalog esa oʻzgarmaydi; testlar ushlaydi:

```
$ go vet ./...
$ staticcheck ./...
$ ./bozor order 1 3
buyurtma: olma × 3 = 36000 soʻm, qoldiq 37
$ ./bozor list
 1. olma                    12000 soʻm   40 dona
 …
$ go test ./...
--- FAIL: TestOrderSaved (0.01s)
    main_test.go:63: order'dan keyin list'da olma qoldigʻi 37 emas:
         1. olma                    12000 soʻm   40 dona
         …
FAIL
FAIL	github.com/abdurahmon27/bozor/cmd/bozor	0.025s
--- FAIL: TestSell (0.00s)
    --- FAIL: TestSell/oddiy (0.00s)
        store_test.go:47: Sell(1, 3) dan keyin katalogda qoldiq 40, kutilgan 37
    --- FAIL: TestSell/hammasi (0.00s)
        store_test.go:47: Sell(4, 2) dan keyin katalogda qoldiq 2, kutilgan 0
FAIL
FAIL	github.com/abdurahmon27/bozor/internal/store	0.010s
FAIL
```

`range` har qadamda elementning **nusxasini** `p`'ga yozadi (Go spec, For statements with range clause), shuning
uchun `p.Stock -= quantity` katalogga yetmaydi; `find` esa `&s.products[i]` — elementning manzilini qaytaradi
(GB-2.04).

## Atomar saqlash: `Save` nima qiladi

1. `json.MarshalIndent` (nil katalog → `[]Product{}`, aks holda `null` yoziladi: «a nil slice encodes as the null
   JSON value», encoding/json).
2. `os.CreateTemp(filepath.Dir(path), ".bozor-*.json")` — **shu papkada** yangi fayl (`*` oʻrniga tasodifiy son),
   `f.Chmod(0o644)` (CreateTemp faylni `0600` bilan ochadi; `bozor.json` GB-3.02'dagidek `0644` boʻlib qoladi).
3. `Write`, `Sync` (baytlar diskka tushguncha kutadi: Sync'siz rename'dan keyin elektr oʻchsa, yangi nom hali
   diskka yozilmagan faylni koʻrsatishi mumkin), `Close`.
4. `os.Rename(tmp, path)` — Linux'da `rename(2)` bitta fayl tizimi ichida `bozor.json`'ni atomar almashtiradi:
   oʻquvchi yo eski, yo yangi faylni koʻradi. Boshqa fayl tizimiga (`/tmp` alohida boʻlsa) rename `EXDEV`
   («invalid cross-device link») bilan rad etiladi — shuning uchun vaqtinchalik fayl `bozor.json` yonida.
5. Xato boʻlsa `defer` vaqtinchalik faylni yopib oʻchiradi (nomli natija `err`, GB-2.09).

Eslatmalar: toʻliq «elektr oʻchsa ham» kafolat uchun rename'dan keyin
papkaning oʻzini ham `Sync` qilish kerak (Linux'da yangi nom papka yozuvida turadi); CLI v0 buni qilmaydi — eng
yomon holda eski katalog qoladi, yarim fayl emas. `go doc os.Rename`: «Even within the same directory, on non-Unix
platforms Rename is not an atomic operation» — atomarlik kafolati Linux/macOS uchun.

## Vazifa

[`TASK.md`](TASK.md): `Store.Restock`, uning jadvalli testi (`internal/store/restock_test.go`) va `bozor restock
ID DONA`. Tekshirish: `go test ./...`.

v0'ning ataylab qoldirilgan cheklovlari (keyingi modullar uchun): bitta jarayon uchun yozilgan — ikki `bozor`
bir vaqtda ishlasa, oxirgi `rename` yutadi (fayl qulfi yoʻq; 4-modulda mutex, 5-modulda Postgres tranzaksiyasi);
`Load` notanish JSON kalitlarini tashlaydi; `Stock`/`Price × dona` uchun yuqori chegara faqat `parsePrice` va
`Restock`'da; `strconv` xatolari foydalanuvchiga xom matn bilan chiqadi.

## Folder self-test (shu kompyuterda)

```bash
gofmt -l . && go vet ./...                                   # chiqish yoʻq
go test ./... ; echo "starter exit $?"                       # FAIL (TASK.md xabarlari), exit 1
(cd solution && go test ./...)                               # ok
t=$(mktemp -d) && cp -r . "$t" && cp -r solution/cmd solution/internal "$t"/ && (cd "$t" && go test ./...)   # ok
```

## Manbalar

- Go spec, For statements with range clause: https://go.dev/ref/spec#For_range
- pkg.go.dev `fmt` (Width): https://pkg.go.dev/fmt
- pkg.go.dev `os.WriteFile` («truncates it before writing … a failure mid-operation can leave the file in a
  partially written state»): https://pkg.go.dev/os#WriteFile
- pkg.go.dev `os.CreateTemp`, `(*os.File).Sync`, `os.Rename`: https://pkg.go.dev/os
- rename(2): https://man7.org/linux/man-pages/man2/rename.2.html
- pkg.go.dev `encoding/json` Marshal (nil slice → null): https://pkg.go.dev/encoding/json#Marshal
- pkg.go.dev `testing` `T.TempDir`: https://pkg.go.dev/testing#T.TempDir
- github.com/google/renameio (ishda shu usulning kutubxonasi): https://pkg.go.dev/github.com/google/renameio/v2

GB-3.05 bilan farq: `diff -ru ../gb-3-05-iterator/solution solution`.
