# GB-2.09 · `defer`, `panic` va `recover`: har biri nima uchun

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #17"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-2-08-xatolar`](../gb-2-08-xatolar) (error'lar chuqur: `%w`, `errors.Is` va `errors.As`) ·
keyingi dars: `gb-2-10` (UTF-8: bayt, rune va oʻzbekcha matn). `writeReport` va `os.Create`:
[`gb-2-07-interface-writer`](../gb-2-07-interface-writer); `findProduct` va nil pointer panic'i:
[`gb-2-04-pointer`](../gb-2-04-pointer).

**Bu darsdan keyin:** erta `return`'da yopilmay qolgan fayl nima uchun «too many open files»ga olib kelishini
koʻrasiz va faylni `defer f.Close()` bilan har qanday `return`'da yopasiz; bir nechta `defer` qaysi tartibda
bajarilishini va `defer` argumentlari qachon hisoblanishini oldindan aytasiz; panic defer'larni bajarib yuqoriga
chiqishini va qachon panic emas, `error` qaytarish kerakligini bilasiz; `recover` va nomli natija bilan panic'ni
chegarada oddiy xatoga aylantirasiz va `recover` qayerda ishlamasligini (defer'siz, ichki funksiyada) koʻrasiz;
`os.Exit` defer'larni bajarmasligini va `run() error` bilan tuzatishni bilasiz.

Bu papka 2-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i bor. GB-2.06/GB-2.08'dagi `cmd/bozor` +
`internal/store` tuzilishi emas, GB-2.07'dagidek bitta `package main` (defer/panic/recover uchun paket kerak emas).
`Product` (maydonlari) va `writeReport` GB-2.07'dan (`String` metodisiz), `findProduct` va `sell` GB-2.04'dan,
`OrderItem` GB-2.03'dan, `ErrNotFound`, `ErrBadQuantity`, `StockError` GB-2.08'dan. 3-modul oxiridagi bozor CLI v0
(GB-3.06) labs'ni birlashtiradi. GB-2.07 bilan farqni koʻrish:
`diff -ru ../gb-2-07-interface-writer/solution solution`.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `writeReport`, `defer f.Close()` bilan `saveReport`, `findProduct`, `safely`, `run` + `main`; GB-2.08'ning xato turlari; vazifaning `ErrPanic`, `handle` va `reserve`'i (chala). Videoda `run` hisobotni `/dev/full`'ga yozgan; bu yerda `hisobot.txt`, shuning uchun `go run .` xatosiz ishlaydi |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi vaqtinchalik `bozor/` papkasida ishga tushiradi va oxirida `(chiqish kodi: N)` yozadi (quyidagi natijalarda bu qator tushirilgan): `./_video/run.sh 06-panic go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

`/dev/full` — Linux'dagi maxsus qurilma: unga yozish doim `no space left on device` (disk toʻla) xatosini beradi.
`ulimit -n 1024` shu terminalda bitta jarayon ochishi mumkin boʻlgan fayllar sonini 1024 ga tushiradi (xato tez
chiqishi uchun). bash'da `ulimit -n` qattiq (hard) limitni ham qoʻyadi. Bu muhim: Go 1.19'dan beri `os`'ni import
qilgan dastur ishga tushganda yumshoq (soft) limitni qattiqqacha oʻzi koʻtaradi
([Go 1.19 release notes](https://go.dev/doc/go1.19#os/exec)). zsh'da `ulimit -n` faqat yumshoq limitni qoʻyadi,
shuning uchun u yerda `ulimit -Hn 1024` yozing. Ikkalasi Linux uchun; macOS'da `/dev/full` yoʻq. Limit 1024, lekin
xato 1020-urinishda: 0, 1, 2 (stdin, stdout, stderr) va Go runtime'ning ikkita ichki fayli (`eventpoll`,
`eventfd`, birinchi `os.Create` paytida ochiladi) ham joy oladi.

Boshlanishi (`01-leak`): `saveReport` yozish xatosida `f.Close()`'gacha yetmaydi:

```
$ ./_video/run.sh 01-leak bash -c 'ulimit -n 1024; go run .'
1: write /dev/full: no space left on device
1020: open /dev/full: too many open files
```

`defer f.Close()` bilan (`02-defer`):

```
$ ./_video/run.sh 02-defer bash -c 'ulimit -n 1024; go run .'
1: write /dev/full: no space left on device
```

`defer` tartibi (`03-tartib`), argument shu zahoti (`04-arg`) va nomsiz funksiya (`05-closure`):

```
$ ./_video/run.sh 03-tartib go run .
hisobot boshlandi
hisobot tugadi
defer 2: buyurtma yopildi
defer 1: fayl yopildi
$ ./_video/run.sh 04-arg go run .
hozir: qoldiq 37
defer: qoldiq 40
$ ./_video/run.sh 05-closure go run .
hozir: qoldiq 37
defer: qoldiq 37
```

panic defer'larni bajarib yuqoriga chiqadi: avval `sell`'ning, keyin `main`'ning defer'i (`06-panic`; videoda
`[signal …]` va stack qatorlari qisqartirilgan):

```
$ ./_video/run.sh 06-panic go run .
defer: sell tugadi
defer: main tugadi
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x20 pc=…]

goroutine 1 [running]:
main.sell(…)
	/tmp/…/bozor/main.go:26 +0x51
main.main()
	/tmp/…/bozor/main.go:38 +0x…
exit status 2
```

`safely` — `defer` + `recover` + nomli natija (`07-recover`); `recover` defer'dagi funksiya chaqirgan boshqa
funksiyada, `toError`'da (`08-quiz`, savol); `recover` defer'siz, `fn`'dan oldin (`12-defersiz`, videoda faqat
aytiladi). Ikkalasida ham panic toʻxtamaydi:

```
$ ./_video/run.sh 07-recover go run .
defer: sell tugadi
xato: panic: runtime error: invalid memory address or nil pointer dereference
main oxiri
$ ./_video/run.sh 08-quiz go run .
defer: sell tugadi
panic: runtime error: invalid memory address or nil pointer dereference
…
exit status 2
$ ./_video/run.sh 12-defersiz go run .
defer: sell tugadi
panic: runtime error: invalid memory address or nil pointer dereference
…
exit status 2
```

Tez-tez uchraydigan xato: `os.Exit` defer'larni bajarmaydi (`09-exit`), tuzatish — `run() error` (`10-run`):

```
$ ./_video/run.sh 09-exit go run .
xato: write /dev/full: no space left on device
exit status 1
$ ./_video/run.sh 10-run go run .
defer: run tugadi
xato: write /dev/full: no space left on device
exit status 1
```

Chuqurroq: sikl ichidagi `defer` (`11-sikl`):

```
$ ./_video/run.sh 11-sikl bash -c 'ulimit -n 1024; go run .'
1020: open hisobot.txt: too many open files
```

## Vazifa

[TASK.md](TASK.md): `handle` panic'ni `ErrPanic`'ni `%w` bilan oʻragan xatoga aylantirsin va `fn`'ning xatosini
oʻzgartirmay qaytarsin; `reserve` xato boʻlganda oldin kamaytirilgan qoldiqlarni `defer` bilan joyiga qaytarsin.
Tekshirish: `go test` (hozir FAIL, vazifa bajarilgach `ok`).

## Manbalar

- Go spec: Defer statements — https://go.dev/ref/spec#Defer_statements
- Go spec: Handling panics — https://go.dev/ref/spec#Handling_panics
- Go spec: Run-time panics — https://go.dev/ref/spec#Run_time_panics
- Go spec: Function literals — https://go.dev/ref/spec#Function_literals
- Go blog: Defer, Panic, and Recover — https://go.dev/blog/defer-panic-and-recover
- Effective Go: Defer, Panic, Recover — https://go.dev/doc/effective_go#defer
- pkg.go.dev `os.Exit` («deferred functions are not run») — https://pkg.go.dev/os#Exit
- pkg.go.dev `os.File.Close` — https://pkg.go.dev/os#File.Close
- pkg.go.dev builtin `panic`, `recover` — https://pkg.go.dev/builtin#recover
- Go Code Review Comments: Don't Panic — https://go.dev/wiki/CodeReviewComments#dont-panic
