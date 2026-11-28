# GB-7.01 · Goroutine va uning hayoti: kim boshlaydi, kim kutadi

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #51"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. `wg.Go` Go 1.25'da qoʻshilgan.

Oldingi dars: `gb-6-08` (nazorat nuqtasi: himoyalangan bozor) · keyingi dars: `gb-7-02` (channel: natijani goroutine'lar
orasida uzatish). Bu darsga kerak boʻlganlar: [`gb-4-07-mutex-race`](../gb-4-07-mutex-race) (server har ulanishni
goroutine'da ishlatadi, data race, `go test -race`, `sync`), [`gb-2-09-defer-panic`](../gb-2-09-defer-panic) (`defer`,
`panic`, `recover`, nomsiz funksiya), [`gb-2-01-slice`](../gb-2-01-slice) (`make`, `append`), [`gb-6-01-middleware`](../gb-6-01-middleware) (`recoverPanic`).

**Bu darsdan keyin:** funksiyani `go` kalit soʻzi bilan yangi goroutine'da boshlaysiz; `main` qaytganda tugamagan
goroutine'lar xatosiz yoʻqolishini bilasiz (200 marta oʻlchab koʻrdik); ularni `time.Sleep` bilan emas,
`sync.WaitGroup` (`wg.Go` + `wg.Wait`) bilan kutasiz; ishdagi eski shaklni (`wg.Add`, `defer wg.Done()`) oʻqiysiz;
natijalarni `append` bilan emas, har goroutine oʻz indeksiga yozib yigʻasiz; har goroutine'ning egasi borligini va
handler ichidagi egasiz `go` nimaga olib kelishini tushunasiz.

Bu papka mustaqil (Modul 7 laboratoriyasi): oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz
`package main`'i, tashqi modul yoʻq. bozor'dagi ehtiyoj: kechasi bir nechta yetkazib beruvchidan narx roʻyxatini olish
(Modul 7 oxirida `cmd/import`, GB-7.07). Tarmoq oʻrniga `time.Sleep`: har yetkazuvchi bir soniyada javob beradi.
Bu laboratoriya nomlari (bozor'ning `Product`/`Order` turlari bu yerda kerak emas): `fetch(name) int`, `offers` (yetkazuvchi → mahsulotlar soni), `fetchAll(suppliers, fetch) []int`,
yetkazuvchilar `chorsu`, `oloy`, `siyob` (testda yana `urgut`, `qoʻqon`), `sendReceipt(o)` (videoda, faqat misol).

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | starter: `fetch`, `offers`, `main` va vazifa: ketma-ket `fetchAll`, boʻsh `fetchAllAddDone` (🟡) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/` | videodagi har bir qadam (`main.go`); `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 04-wg go run .` |
| `_video/race.sh` | `-race` buyruqlari uchun (quyidagi eslatma): `./_video/race.sh 05-append go run -race .` |
| `_video/olchov.sh` | videodagi ikki oʻlchov: `./_video/olchov.sh go`, `./_video/olchov.sh append` |

Videoda nima oʻzgarganini koʻrish: `diff -u _video/01-ketma/main.go _video/04-wg/main.go`.

## Videodagi qadamlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Natijalar har ishga tushirishda biroz boshqa: vaqt va goroutine'lar tartibi kafolatlanmaydi.

**`01-ketma`** — `fetch` uch marta ketma-ket:

```
$ go run .
boshlandi: chorsu
tayyor: chorsu
boshlandi: oloy
tayyor: oloy
boshlandi: siyob
tayyor: siyob
jami: 3.002459679s
```

**`02-go`** — chaqiruv oldida `go`, kutish yoʻq. `main` darhol qaytadi, dastur tugaydi:

```
$ go run .
jami: 8.335µs
$ go run .
boshlandi: chorsu
jami: 6.407µs
```

`µs` — mikrosekund, soniyaning milliondan biri. `./_video/olchov.sh go` (200 marta): videodagi oʻlchovda 110 marta faqat
`jami`, 32 marta bitta, 37 marta ikkita, 21 marta uchta `boshlandi`; boshqa oʻlchovlarda 85/41/51/23 va 108/67/11/14.
`tayyor` 200 martadan **birortasida ham** chiqmadi. Oʻlchovlar orasidagi farq katta (yana biri: 66/37/66/31). Sababi Go spec'da ([Program execution](https://go.dev/ref/spec#Program_execution)):
"When that function invocation [main] returns, the program exits. It does not wait for other (non-main) goroutines to complete."

**`03-sleep`** — `main` oxirida `time.Sleep(500 * time.Millisecond)`:

```
$ go run .
boshlandi: siyob
boshlandi: oloy
boshlandi: chorsu
jami: 500.824023ms
```

Goroutine'lar boshlandi, lekin `main` yarim soniyada qaytdi. Ikki soniya qoʻysak, hammasi chiqadi, lekin bir soniya behuda;
yetkazuvchi uch soniyada javob bersa, yana yoʻqoladi. `Sleep` — taxmin.

**`04-wg`** — `var wg sync.WaitGroup`, `wg.Go(func() { fetch(name) })`, `wg.Wait()`:

```
$ go run .
boshlandi: siyob
boshlandi: chorsu
boshlandi: oloy
tayyor: chorsu
tayyor: siyob
tayyor: oloy
jami: 1.003058399s
```

`go doc sync.WaitGroup.Go`: "Go calls f in a new goroutine and adds that task to the WaitGroup. When f returns, the task is
removed from the WaitGroup. … The function f must not panic. … the return from f «synchronizes before» the return of any
Wait call that it unblocks" ([Go memory model](https://go.dev/ref/mem)). `wg.Go` [Go 1.25'da](https://go.dev/doc/go1.25) qoʻshilgan.

**`08-adddone`** — ishdagi eski shakl, natija xuddi shunday (`jami: 1.000833322s`):

```go
wg.Add(1)
go func() {
	defer wg.Done()
	fetch(name)
}()
```

`wg.Add(1)` `go`dan **oldin** turadi (goroutine ichida emas): hisob nol paytda `Wait` ishga tushsa, u darhol qaytishi mumkin
(`go doc sync.WaitGroup.Add`; u yerda: "Callers should prefer WaitGroup.Go"). `wg.Go` shu ikki qatorni oʻzi bajaradi, bitta farq
bilan: `f` panic qilsa, `wg.Go` `Done`ni chaqirmaydi va panic'ni qayta tashlaydi (`src/sync/waitgroup.go`; pastda `07-panic`).
`sync.WaitGroup` GB-4.07'ning testida (`orderParallel`) birinchi marta koʻrindi; bu darsda uning ichi ochiladi.
Xato qaytaradigan ishlar uchun `WaitGroup` yetmaydi (xatoni yigʻish yoʻli yoʻq): `errgroup`, GB-7.04.

**`05-append`** — `fetch` endi mahsulotlar sonini qaytaradi, natija `counts = append(counts, fetch(name))` bilan
yigʻiladi. `go run -race .` (Docker, `race.sh`), qisqartirilgan:

```
boshlandi: siyob
…
==================
WARNING: DATA RACE
Read at 0x00c000124030 by goroutine 9:
  main.main.func1()
      /home/dev/bozor/main.go:23 +0x54

Previous write at 0x00c000124030 by goroutine 8:
  main.main.func1()
      /home/dev/bozor/main.go:23 +0xc9
…
mahsulotlar: [40 5]
jami: 1.01250767s
Found 2 data race(s)
exit status 66
```

`main.main.func1` — `main` ichidagi birinchi nomsiz funksiya; 23-qator — `append` qatori. `66` — race topilganda dastur
qaytaradigan kod (`GORACE` `exitcode`, sukut boʻyicha 66, [hujjat](https://go.dev/doc/articles/race_detector#Options)).
`-race`'siz (`./_video/olchov.sh append`, 200 marta): bir oʻlchovda 193 marta uchala son, 6 marta ikkita, 1 marta bitta son
qoldi (7 tasida yoʻqoldi); boshqa oʻlchovda 180 / 13 / 7 (20 tasida). Tartib tasodifiy (`[5 40 120]`, `[120 40 5]`, …).

**`06-index`** — `counts := make([]int, len(suppliers))`, `for i, name := range suppliers`, `counts[i] = fetch(name)`:

```
$ go run -race .
…
tayyor: chorsu
mahsulotlar: [40 120 5]
jami: 1.000380247s
```

Har goroutine oʻz elementiga yozadi (har element alohida xotira joyi), `Wait` qaytgandan keyin `counts`ni oʻqish xavfsiz.
Go 1.22'dan beri `for` siklining har aylanishida `i` va `name` yangi oʻzgaruvchi ([Go 1.22](https://go.dev/doc/go1.22#language)),
shuning uchun har nomsiz funksiya oʻz `i`/`name`ini koʻradi (`go.mod`dagi `go 1.27` buni yoqadi).

**`07-panic`** (Chuqurroq) — `main`da `defer` + `recover`, panic esa `wg.Go` goroutine'ida:

```
$ go run .
panic: chorsu: javob buzuq [recovered, repanicked]

goroutine 7 [running]:
…
exit status 2
```

`ushlandi` chiqmadi: `recover` faqat oʻz goroutine'idagi panic'ni ushlaydi ([Go spec: Handling panics](https://go.dev/ref/spec#Handling_panics)).
`[recovered, repanicked]`: `wg.Go` panic'ni ushlab, `Done` chaqirmasdan qayta tashlaydi (aks holda `Wait` ochilib, dastur
panic tugashidan oldin chiqib ketishi mumkin edi). Shuning uchun GB-6.01'dagi `recoverPanic` handler goroutine'ini himoya
qiladi, lekin handler ichida `go` bilan boshlangan goroutine'dagi panic butun serverni yiqitadi.

## `-race`: Linux va macOS talablari

Linux'da race detector uchun `cgo` yoqilgan va C kompilyator (odatda `gcc`) oʻrnatilgan boʻlishi kerak; macOS'da Go 1.20'dan
beri kerak emas ([Go 1.20](https://go.dev/doc/go1.20#cgo)). Bu mashinada gcc yoʻq, shuning uchun videodagi `-race` natijalari
Docker'da olingan (`_video/race.sh`: gcc bor `golang:1.26.8` image'i, Go oʻzi shu mashinadagi `go1.27.4` toolchain), hisobotdagi
`/home/dev/bozor/` — konteynerdagi papka; sizda oʻz papkangiz chiqadi. Batafsil: [`gb-4-07-mutex-race`](../gb-4-07-mutex-race) README.

## Vazifa

[TASK.md](TASK.md): `fetchAll`ni `wg.Go` + `wg.Wait` bilan yozing, har goroutine oʻz indeksiga yozsin; 🟡 `fetchAllAddDone`
eski shaklda. Tekshirish: `go test .` (starterda FAIL, xabar oʻzbekcha). Javob: `go test ./solution/`.

## Har goroutine'ning egasi

Videodagi `createOrder` misoli (`go sendReceipt(o)` — handler javob yozib qaytadi, goroutine'ni hech kim kutmaydi): server
toʻxtasa, yarim yoʻldagi chek yoʻqoladi; panic boʻlsa, server yiqiladi; nechtasi ishlayotganini hech kim bilmaydi. Qoida:
goroutine'ni boshlagan kod uning qachon tugashini biladi va uni kutadi. Toʻgʻri yechimlar keyinroq: serverni toʻxtatishda
ishni tugatish — GB-8.03 (graceful shutdown), yoʻqolmaydigan fon ishi (bazadagi navbat) — GB-9.01, bekor qilish (`context`) —
GB-7.03.

## Manbalar

- Go spec: [Go statements](https://go.dev/ref/spec#Go_statements), [Program execution](https://go.dev/ref/spec#Program_execution), [Handling panics](https://go.dev/ref/spec#Handling_panics)
- [`sync.WaitGroup`](https://pkg.go.dev/sync#WaitGroup), [`WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go), [Go 1.25 release notes](https://go.dev/doc/go1.25)
- [The Go Memory Model](https://go.dev/ref/mem)
- [Go 1.22: per-iteration loop variables](https://go.dev/doc/go1.22#language), [Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview)
- [Data Race Detector](https://go.dev/doc/articles/race_detector)
- [`time.Sleep`](https://pkg.go.dev/time#Sleep), [`time.Since`](https://pkg.go.dev/time#Since)
