# GB-7.02 · Channel: natijani goroutine'lar orasida uzatish

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #52"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-7-01-goroutine-hayoti`](../gb-7-01-goroutine-hayoti) (`go`, main qaytsa goroutine'lar toʻxtashi, `wg.Go`/`wg.Wait`)
· keyingi dars: [`gb-7-03-select-context`](../gb-7-03-select-context) (`select`, `context` bilan toʻxtatish).
Bu darsga kerak boʻlganlar: [`gb-4-07-mutex-race`](../gb-4-07-mutex-race) (data race, deadlock xabari),
[`gb-3-01-fayl-bufio`](../gb-3-01-fayl-bufio) (`parseLine`, «ID;nom;narx;qoldiq»), [`gb-2-09-defer-panic`](../gb-2-09-defer-panic) (`defer`, `panic`),
[`gb-2-02-map-qidirish`](../gb-2-02-map-qidirish) (map, `v, ok :=`), [`gb-2-01-slice`](../gb-2-01-slice) (`make`, `len`, `cap`).

**Bu darsdan keyin:** goroutine natijasini channel orqali uzatasiz; buffer'siz channel'da yuborish va qabul qilish bir-birini
kutishini, hamma goroutine bloklansa runtime `fatal error: all goroutines are asleep - deadlock!` berishini va `[chan send]` /
`[chan receive]` qatorlarini oʻqishni bilasiz; `close` + `for range` bilan oqimni tugatasiz; buffered channel nimani
oʻzgartirishini (`len`/`cap`) koʻrasiz; `parse → inStock` pipeline'ini qurasiz va ikki yuboruvchini bitta `close` bilan toʻgʻri yopasiz.

Bu papka mustaqil (Modul 7 laboratoriyasi): oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i,
tashqi modul yoʻq. `Product` va `parseLine` GB-3.01'dan (`product.go`); yangi nomlar: `chorsu`, `oloy` (`data.go`: ikki yetkazib
beruvchining roʻyxati, GB-3.01 formatida), `parse`, `inStock`, `unique`, `importAll`, channel'lar `raw` / `stocked` / `ready`.
Modul oxiridagi importer (GB-7.07) shu bosqichlardan quriladi. Soddalashtirish: `parse` xato qatorni tashlab ketadi
(`continue`); xatoni pipeline'dan qaytarish GB-7.04'da (`errgroup`).

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `parse`, `inStock`, ikki yuboruvchili pipeline (`importAll`) + vazifa uchun boʻsh `unique` |
| `product.go` | `Product`, `parsePrice`, `parseLine` (GB-3.01) |
| `data.go` | `chorsu` (4 qator) va `oloy` (2 qator) roʻyxatlari |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu `product.go`, `data.go`, test) |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 02-deadlock go run .` |
| `_video/race.sh` | `-race` uchun (Linux'da gcc kerak; bu mashinada yoʻq, Docker'da): `./_video/race.sh solution go test -race .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64)

Fayl yoʻllari `/tmp/tmp.XXXX/bozor/main.go` qisqartirilib `…/bozor/main.go` deb yozilgan; `+0x36` kabi siljishlar va
goroutine raqamlari har buildda boshqacha.

**01-uzatish** — goroutine `counts`ga 40 yuboradi, `main` qabul qiladi:

```
$ go run .
40
```

**02-deadlock** — goroutine'siz: `main` 7-qatorda oʻzi yuboradi, 8-qatorda oʻzi qabul qilmoqchi. Buffer'siz channel'da yuborish
qabul qiluvchini kutadi, u esa hech qachon kelmaydi:

```
$ go run .
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	…/bozor/main.go:7 +0x36
exit status 2
```

**03-range** — `parse` yopmaydi, `for p := range products` oxirgi qiymatdan keyin abadiy kutadi, `import tugadi` chiqmaydi:

```
$ go run .
olma 40
oʻrik 0
non 100
choy 8
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.main()
	…/bozor/main.go:18 +0x199
exit status 2
```

**04-close** — `parse` boshida `defer close(out)`; sikl tugaydi va `import tugadi` chiqadi:

```
$ go run .
olma 40
oʻrik 0
non 100
choy 8
import tugadi
```

**05-buffer** — `make(chan int, 2)`: ikki yuborish kutmaydi, uchinchisi buffer toʻla boʻlgani uchun kutadi:

```
$ go run .
2 2
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
main.main()
	…/bozor/main.go:10 +0xe6
exit status 2
```

**06-quiz** — buffer'da 40, keyin `close`; ikki marta oʻqish:

```
$ go run .
40
0 false
```

**07-pipeline** — `parse(chorsu) → raw → inStock → ready → main` (oʻrikning qoldigʻi 0, oʻtmaydi):

```
$ go run .
olma 40
non 100
choy 8
```

**08-ikki** — klassik xato: ikkita `parse` bitta `raw`ga yozadi va **ikkalasi ham** uni yopadi. Natija har safar boshqacha.
Shu mashinada 6 ta ishga tushirishdan 1 tasi xatosiz, lekin `oloy` mahsulotlari jimgina yoʻqoldi (main qaytdi, qolgan goroutine
toʻxtadi), 5 tasi panic:

```
$ go run .
olma 40
non 100
choy 8
$ go run .
olma 40
non 100
panic: send on closed channel
	panic: close of closed channel

goroutine 8 [running]:
main.parse({0x57a7a8?, 0x0?, 0x0?}, 0x24e8cb6be230)
	…/bozor/main.go:12 +0xf6
created by main.main in goroutine 1
	…/bozor/main.go:29 +0x15c
exit status 2
```

Ikkinchi `panic` qatori: panic paytida `parse`ning `defer close(out)`i ishlab, yopiq channel'ni yana yopdi.

**09-wg** — tuzatish: `parse` yopmaydi; ikkala `parse` `wg.Go` bilan, alohida goroutine `wg.Wait()`dan keyin `raw`ni **bir marta**
yopadi (Go blog'idagi "merge" shakli). Tartib har safar boshqacha, non ikki marta keladi (vazifa shu):

```
$ go run .
olma 40
non 60
anor 25
non 100
choy 8
```

## Chuqurroq: channel jadvali (`_video/10-jadval`)

`./_video/run.sh 10-jadval go run . <holat>` (har holatning birinchi qatorlari):

| Holat | Natija |
|---|---|
| `nil-send` (nil channel'ga yuborish) | `fatal error: all goroutines are asleep - deadlock!` · `goroutine 1 [chan send (nil chan)]:` |
| `nil-recv` (nil channel'dan oʻqish) | `… deadlock!` · `goroutine 1 [chan receive (nil chan)]:` |
| `nil-close` | `panic: close of nil channel` |
| `closed-close` (qayta yopish) | `panic: close of closed channel` |
| `closed-send` (yopiq channel'ga yuborish) | `panic: send on closed channel` |
| `closed-recv` (yopiq channel'dan oʻqish) | `0 false` |

Yoʻnalishli tur (`_video/11-yonalish`: `inStock` ichida `close(in)`):

```
$ go vet .
# github.com/abdurahmon27/bozor
# [github.com/abdurahmon27/bozor]
vet: ./main.go:11:8: invalid operation: cannot close receive-only channel in (variable of type <-chan Product)
```

`chan<- Product`dan oʻqish ham kompilyatsiya xatosi (`_video/13-faqat-yuborish`):

```
$ go build .
# github.com/abdurahmon27/bozor
./main.go:5:9: invalid operation: cannot receive from send-only channel chan<- Product out (variable of type chan<- Product)
```

Boshqa goroutine tirik boʻlsa (`_video/12-osilish`: fonda `time.Sleep` siklidagi goroutine; serverda bu ulanish kutayotgan goroutine boʻladi),
runtime deadlock'ni aniqlamaydi: dastur xabarsiz abadiy osilib qoladi (`timeout 4 ./_video/run.sh 12-osilish go run .` →
4 soniyadan keyin `timeout` toʻxtatadi, chiqish kodi 124). Bunday osilib qolgan goroutine — goroutine leak (GB-7.06).

`-race` (Docker, `race.sh`): `./_video/race.sh solution go test -race -count=3 .` → `ok … 1.023s`;
`./_video/race.sh 09-wg go run -race .` → 5 qator, `WARNING: DATA RACE` yoʻq.

## Vazifa

[`TASK.md`](TASK.md): `unique` bosqichini yozing va `importAll`ga qoʻying; `go test .` `ok` boʻlguncha.

```bash
go test .                # starterda yiqiladi (oʻzbekcha xabar bilan)
go test ./solution/      # javob: ok
```

## Manbalar

- Go spec: [Channel types](https://go.dev/ref/spec#Channel_types), [Send statements](https://go.dev/ref/spec#Send_statements),
  [Receive operator](https://go.dev/ref/spec#Receive_operator), [For range](https://go.dev/ref/spec#For_range),
  [Close](https://go.dev/ref/spec#Close), [Length and capacity](https://go.dev/ref/spec#Length_and_capacity),
  [Program execution](https://go.dev/ref/spec#Program_execution)
- [The Go Memory Model: Channel communication](https://go.dev/ref/mem#chan)
- [`builtin.close`](https://pkg.go.dev/builtin#close) ("should be executed only by the sender, never the receiver")
- [Effective Go: Channels](https://go.dev/doc/effective_go#channels), [Go blog: Pipelines and cancellation](https://go.dev/blog/pipelines)
- [`sync.WaitGroup.Go`](https://pkg.go.dev/sync#WaitGroup.Go) (Go 1.25'dan)
