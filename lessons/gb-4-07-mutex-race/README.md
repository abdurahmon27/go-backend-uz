# GB-4.07 · Soʻrovlar bir vaqtda ishlaydi: `sync.Mutex` va `-race`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #31"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: `gb-4-06` (POST /products: JSON kirish, limit, validatsiya) · keyingi dars: `gb-4-08` (katalog API nazorat
nuqtasi). Bu darsga kerak boʻlganlar: [`gb-4-03-handler-httptest`](../gb-4-03-handler-httptest) (handler, `httptest`),
[`gb-3-05-iterator`](../gb-3-05-iterator) (`iter.Seq`, `yield`: "Chuqurroq" qismida), [`gb-2-09-defer-panic`](../gb-2-09-defer-panic)
(`defer`), [`gb-2-02-map-qidirish`](../gb-2-02-map-qidirish) (map), [`gb-1-08-birinchi-test`](../gb-1-08-birinchi-test) (`go test`).

**Bu darsdan keyin:** server har ulanishni alohida goroutine'da ishlatishini, shuning uchun soʻrovlar bir vaqtda ishlashini `curl` bilan koʻrasiz;
xotiradagi katalogda (`map`) ming parallel soʻrov qanday qoldiqni buzishini va `go test` nega `fatal error: concurrent map …`
bilan qulashini bilasiz; `go test -race` hisobotini (`Read at` / `Previous write`) oʻqiysiz; `sync.Mutex`
(`Lock` + `defer Unlock`) bilan tuzatasiz; mutex'ni nusxalash xatosini `go vet` bilan topasiz.

Bu papka mustaqil (Modul 4 laboratoriyasi): oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i,
tashqi modul yoʻq. `Product` GB-3.02'dan; yangi: `Store` (xotiradagi katalog, `mu sync.Mutex`), `Sell`, `Stock`, `Restock`,
`Total`, `buy` (POST /buy), `slowHandler`. Keyingi dars (GB-4.08) bularni katalog API'ga birlashtiradi. Soddalashtirishlar (laboratoriya): `Sell` miqdorni tekshirmaydi
(musbat deb oladi; haqiqiy `store` paketida tekshiruv GB-2.08'da), `buy` har xatoga 409 Conflict beradi (haqiqiy API'da yoʻq mahsulot
404 boʻlardi), `map[int64]Product` — ID boʻyicha qidirish uchun (GB-3.06'dagi `internal/store` slice'i shu darsda qisqartirildi).
Mutex yoki channel? Bu darsda umumiy `map`ni mutex himoyalaydi; channel'lar bilan bogʻliq usullar GB-7'da.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod (`Store`, qulf bilan `Sell`, `buy`, `slowHandler`; videoda `Stock` ham qulflanadi) + `NewStore`, `main` va vazifa uchun `Restock` (qulfsiz) hamda `Total` (boʻsh) |
| `buy_test.go` | videoda yozilgan test: `orderParallel` va `TestBuyParallel` (1000 parallel soʻrov): starterda ham oʻtadi |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang. Qulfsiz kod ba'zan oʻzi qulaydi, shuning uchun parallel tekshiruv ikkinchi jarayonda ishlaydi |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayllari) |
| `_video/NN-…/` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 02-race go test` |
| `_video/race.sh` | `-race` buyruqlari uchun (quyidagi eslatma): `./_video/race.sh 02-race go test -race -count=1 .` |

## `-race`: Linux va macOS talablari

Linux'da race detector uchun `cgo` yoqilgan va C kompilyator (odatda `gcc`) oʻrnatilgan boʻlishi kerak. macOS'da Go 1.20'dan beri
race detector Xcode yoki cgo talab qilmaydi ([Go 1.20 eslatmalari](https://go.dev/doc/go1.20#cgo)); bu dars uchun Xcode oʻrnatish shart emas.
Linux'da C kompilyator yoʻq boʻlsa, Go `cgo`ni oʻzi oʻchirib qoʻyadi va `go test -race` shunday deydi (shu mashinada ham shunday):

```
$ go test -race
go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
```

`CGO_ENABLED=1` bilan majburlasangiz ham, kompilyatorsiz natija shu:

```
$ CGO_ENABLED=1 go test -race
# runtime/cgo
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH
```

Videodagi `-race` natijalari shu sababli Docker'da olingan: gcc bor `golang:1.26.8` image'i, lekin Go oʻzi shu mashinadagi
`go1.27.4` toolchain (`_video/race.sh`), shuning uchun hisobotdagi fayl yoʻli `/home/dev/bozor/main.go` (konteynerdagi papka):
sizda oʻz papkangiz yoʻli chiqadi. Linux'da C kompilyator oʻrnating (masalan, `apt install gcc`) va `cgo` yoqilganini tekshiring
(`go env CGO_ENABLED`: `1`); keyin oddiy `go test -race` yetadi.
Manzillar (`0x00c0…`) va goroutine raqamlari har ishga tushirishda boshqacha boʻladi; hisobotlar "…" bilan qisqartirilgan.
`-race` faqat ishga tushgan kodda sodir boʻlgan murojaatlarni topadi, shuning uchun test parallel chaqirishi kerak.

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64)

`localhost:8085` (erkin port tanladik; sizda band boʻlsa, boshqasini yozing): `01-slow` — `slowHandler`, soʻrov boshida `boshladi`,
bir soniyadan keyin `tugadi` yozadi. Server fonda (`&`), keyin bitta `curl` uchta bir xil manzilni bir vaqtda soʻraydi
(`--parallel-immediate` curl'ni uchala ulanishni darhol ochishga majbur qiladi; usiz curl ulanish qayta ishlatilishini
kutib, ikkinchi va uchinchi soʻrovni kechiktiradi). Chiqish — serverning logi (`curl` hech narsa chiqarmaydi, handler body yozmaydi):

```
$ go run . &
$ curl -s --parallel --parallel-immediate localhost:8085/slow localhost:8085/slow localhost:8085/slow
boshladi
boshladi
boshladi
tugadi
tugadi
tugadi
$ kill %1        # serverni toʻxtatish (interaktiv shell'da; `go run` ishlatgan dastur ham yopilmasa: pkill -f bozor)
```

Uchala soʻrov bir soniya ichida birga ishladi (ketma-ket boʻlsa — uch soniya; usiz `--parallel` ikki soniya olgan): server har yangi
ulanishni alohida goroutine'ga beradi (`net/http` `Server.Serve` hujjati: "creating a new service goroutine for each [connection]").
Bir ulanishdagi (keep-alive) soʻrovlar ketma-ket oʻqiladi; bir vaqtlilik turli ulanishlardan keladi.

Mutex'siz `Store` (`02-race`: `Sell` qulfsiz, `buy_test.go` — 1000 parallel `POST /buy`). `go test` `CGO_ENABLED=0` bilan:

```
$ go test
fatal error: concurrent map read and map write     (baʼzan: concurrent map writes; juda kamdan-kam oddiy FAIL: qoldiq nol emas)
…
FAIL	github.com/abdurahmon27/bozor	0.006s
FAIL
```

`go test -race` (Docker, `race.sh`), hisobot qisqartirilgan:

```
WARNING: DATA RACE
Read at 0x00c00016c690 by goroutine 13:
  github.com/abdurahmon27/bozor.(*Store).Sell()
      /home/dev/bozor/main.go:27 +0x97
…
Previous write at 0x00c00016c690 by goroutine 10:
  github.com/abdurahmon27/bozor.(*Store).Sell()
      /home/dev/bozor/main.go:32 +0x144
```

`main.go:27` — `p, ok := s.products[id]` (oʻqish), `main.go:32` — `s.products[id] = p` (yozish). Tuzatish (`03-mutex`):
`Store`da `mu sync.Mutex`, `Sell`da `s.mu.Lock()` va `defer s.mu.Unlock()`:

```
$ go test -race
ok  	github.com/abdurahmon27/bozor	1.024s
```

Bu yerdagi taxminan bir soniyaning asosiy sababi — race runtime chiqishdan oldin kutishi: `GORACE` sozlamasi
`atexit_sleep_ms` sukut boʻyicha `1000` millisoniya ([hujjat](https://go.dev/doc/articles/race_detector#Options)). Bu vaqt kodning
200 barobar sekinlashganini bildirmaydi. Kutishni oʻchirib, natijani keshdan olmasdan tekshirish mumkin:

```
$ GORACE=atexit_sleep_ms=0 go test -race -count=1 .
ok  	github.com/abdurahmon27/bozor	0.082s
```

`_video/race.sh 03-mutex` muhitida uch ishga tushirish: `ok` vaqti odatiy `-race` bilan 1,045–1,062 s, kutishsiz
0,067–0,082 s. Alohida, oldindan yigʻilgan test dasturlarida oʻrtacha wall vaqt kutishsiz race bilan 0,102 s, oddiy rejimda
0,036 s (taxminan 2,8 barobar; atigi uchta oʻlchov, umumiy benchmark emas). Tafsilotlar yozma darsda; race tekshiruvi
qoʻshimcha vaqt va xotira sarflaydi, lekin chiqishdagi kutish alohida xarajat.

Klassik xato (`04-nusxa`): `func (s Store) Sell(…)` — qabul qiluvchi nusxa, qulf ham nusxalanadi, himoya yoʻq (`go test` yana
qulaydi). `go vet` sababini koʻrsatadi:

```
$ go vet ./...
main.go:28:9: Sell passes lock by value: github.com/abdurahmon27/bozor.Store contains sync.Mutex
```

Chuqurroq (`05-deadlock`): GB-3.05'dagi `All` iteratorini qulf bilan oʻrab, sikl tanasidan `Sell`ni chaqirish. Tana `yield`
sifatida `All` qulfni ushlab turgan paytda ishlaydi, `Sell` esa oʻsha qulfni kutadi:

```
$ go run .
fatal error: all goroutines are asleep - deadlock!
…
goroutine 1 [sync.Mutex.Lock]:
…
main.(*Store).Sell(0x784cb554030, 0x3, 0x1)
…
main.(*Store).All.func1(0x784cb624e90)
```

Qoida: qulfni `yield`dan oldin qoʻyib yuboring, yoki nusxani qulf ichida olib, qulfsiz aylaning.

## Vazifa

[`TASK.md`](TASK.md): `Restock`ni qulf bilan himoyalang va `Total`ni yozing; `go test .` `ok` boʻlguncha.

```bash
go test .                # starterda yiqiladi (oʻzbekcha xabar bilan)
go test ./solution/      # javob: ok
go test -race .          # Linux: C kompilyator + cgo; macOS: Xcode shart emas
```

## Manbalar

- Go `net/http`: [`Server.Serve`](https://pkg.go.dev/net/http#Server.Serve) (har ulanish uchun alohida goroutine), [`ListenAndServe`](https://pkg.go.dev/net/http#Server.ListenAndServe)
- Go [`sync.Mutex`](https://pkg.go.dev/sync#Mutex) (nol qiymati — ochiq qulf; nusxalanmaydi), [`sync.WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go) (Go 1.25'dan; toʻliq GB-7.01'da)
- [Data Race Detector](https://go.dev/doc/articles/race_detector) (`atexit_sleep_ms=1000`; Linux: cgo + C kompilyator; xotira 5–10×, vaqt 2–20× koʻproq), [Go xotira modeli](https://go.dev/ref/mem)
- [Go FAQ: map'lar parallel foydalanish uchun xavfsizmi](https://go.dev/doc/faq#atomic_maps)
- [`go vet` (copylocks)](https://pkg.go.dev/cmd/vet), [RFC 9110 §15.5.10 — 409 Conflict](https://www.rfc-editor.org/rfc/rfc9110#name-409-conflict)
