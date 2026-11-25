# 15 · `defer`, `panic`, `recover`: nimani qutqarib boʻlmaydi

Oldingi dars: GO-1.14 · Xato ham qiymat: `%w`, `errors.Is`, `AsType` va `Join` · Keyingi dars: GO-1.16 · `go.mod`dagi bitta qator dasturingiz xulqini oʻzgartiradi

| Papka | Nima |
|---|---|
| `sikl/` | 2000 ta faylni sikl ichida `defer f.Close()` bilan ochadi: **ataylab xato**, `ulimit -n 1024` da `too many open files` |
| `sikltuzat/` | sikl tanasi alohida `bittasi()` funksiyasida: har bir fayl darhol yopiladi |
| `tartib/` | uchta `defer`: teskari tartib (LIFO) |
| `vaqt/` | `defer fmt.Println(time.Since(start))` — **ataylab xato** (`go vet` ushlaydi), closure bilan toʻgʻrisi |
| `panik/` | indeks chegaradan tashqarida: `defer` baribir ishlaydi, `exit status 2` |
| `qutqar/` | `recover` + nomlangan natija: panic oddiy xatoga aylanadi |
| `goroutine/` | `main` dagi `recover` boshqa goroutine'dagi panic'ni ushlay olmaydi |
| `server/` | `net/http` handler'dagi panic'ni ushlaydi, handler ochgan goroutine'dagisini esa yoʻq (`demo.sh`) |
| `fatal/deadlock`, `fatal/stack`, `fatal/xotira` | `recover` ushlay olmaydigan `fatal error` lar |
| `exit/` | `log.Fatal` → `os.Exit`: `defer` ishlamaydi; `exit/tuzat/` — `run() error` bilan tuzatilgani |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/015-defer-panic-recover`) ishga tushiriladi.
```bash
sed -n 9,20p sikl/main.go
bash -c 'ulimit -n 1024; go run ./sikl'          # ochildi: 1019 ta fayl / xato: open log-1020.txt: too many open files
sed -n 6,17p tartib/main.go
go run ./tartib                        # ... lock qaytarildi / fayl yopildi / baza yopildi
sed -n 9,17p vaqt/main.go
go run ./vaqt                          # closure:  100ms / argument: 0s
go vet ./vaqt                          # call to time.Since is not deferred
sed -n 9,16p sikltuzat/main.go
bash -c 'ulimit -n 1024; go run ./sikltuzat'     # ochildi: 2000 ta fayl
sed -n 6,13p panik/main.go
go run ./panik                         # defer: baza yopildi / panic: ... / exit status 2
go run ./qutqar
sed -n 9,23p goroutine/main.go
go run ./goroutine                     # panic: assignment to entry in nil map / exit status 2
sed -n 20,27p server/main.go
bash ./server/demo.sh
go run ./fatal/deadlock 2>&1 | grep -E '^(fatal|qutqarildi|exit)'
go run ./fatal/stack 2>&1 | grep -E '^(runtime: goroutine|fatal|qutqarildi|exit)'
(ulimit -v 1000000; go run ./fatal/xotira) 2>&1 | grep -E '^(runtime: out|fatal|qutqarildi|exit)'
sed -n 9,15p exit/main.go
go run ./exit                          # log xabari bor, "defer:" qatori yoʻq / exit status 1
go run ./exit/tuzat                    # avval "defer:" qatori, keyin log / exit status 1
```
Natijalar `go1.26.8 linux/amd64` da. Video buyruqlarni `bash -c` bilan ishga tushiradi. **zsh'da** `ulimit -n` faqat
yumshoq limitni oʻzgartiradi va Go uni qayta koʻtaradi, shuning uchun `bash -c '…'` shaklidan foydalaning.
`server/demo.sh` da `CGO_ENABLED=0`: mashinada C kompilyatori boʻlmasa, `net` paketi cgo bilan yigʻilmaydi.
`go vet ./...` `vaqt/` da ataylab xato beradi; qolganlarini tekshirish:
`CGO_ENABLED=0 go vet ./sikl ./sikltuzat ./tartib ./panik ./qutqar ./goroutine ./server ./fatal/... ./exit ./homework/...`

## Dars matni

**Hook.** `sana` 2000 ta log faylni ochadi va har birini `defer f.Close()` bilan "yopadi". Ochiq fayllar limiti
1024 boʻlsa, 1019 tadan keyin `too many open files` (Linux'da 3 tasi stdin/stdout/stderr, 2 tasi runtime'ning poller'i: `eventpoll` va `eventfd`).
Eslatma: Go dastur ishga tushganda yumshoq (soft) `RLIMIT_NOFILE` ni qattiq (hard) limitdan bittaga kamgacha koʻtaradi
(`syscall/rlimit.go`: `nlim.Cur = nlim.Max - 1`), shuning uchun demo bash'dagi `ulimit -n` bilan ikkalasini ham 1024 ga tushiradi. Real hayotda
limit kattaroq, lekin uzoq yoki cheksiz siklda (masalan, navbatdan oʻqiyotgan worker) fayllar baribir toʻplanadi.

**Tartib.** `defer` chaqiruvni oʻrab turgan funksiya tugaguncha kechiktiradi; bir nechta `defer` teskari tartibda
(LIFO) ishlaydi. Spec: *"deferred functions are invoked immediately before the surrounding function returns, in the
reverse order they were deferred"*. Funksiya `return` bilan, oxiriga yetib yoki panic bilan tugasa ham ishlaydi.

**Argument vaqti.** *"Each time a "defer" statement executes, the function value and parameters to the call are
evaluated as usual and saved anew but the actual function is not invoked."* Shuning uchun
`defer fmt.Println(time.Since(start))` funksiya boshidagi ≈0 ni saqlaydi. Closure ichidagi kod esa oxirida ishlaydi va
100ms ni koʻradi. `go vet` ning `defers` tekshiruvi aynan shu holatni ushlaydi (`go tool vet help defers`). Xuddi shu
qoida method value uchun ham: value receiver'li `defer b.Chiqar()` da `b` defer qatorida nusxalanadi (GO-1.12 uy
vazifasi); pointer receiver'da esa `&b` saqlanadi, ya'ni keyingi oʻzgarishlar koʻrinadi.

**Sikldagi defer.** `defer` sikl qadami tugaganda emas, funksiya tugaganda ishlaydi. Kompilyator ham, `go vet` ham jim.
[staticcheck SA9001](https://staticcheck.dev/docs/checks/#SA9001) shu mavzuda: *"Deferred calls run when the
surrounding function returns, not at the end of each loop iteration."* (u asosan kanal ustidagi `range` siklini
tekshiradi). Narxi ham bor: Go 1.14 dan beri oddiy `defer` deyarli bepul (open-coded defer), lekin sikl ichidagi
`defer` bu yoʻldan foydalana olmaydi
([walk/stmt.go](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/walk/stmt.go#L116-L120)).
Tuzatish: sikl tanasini alohida funksiyaga chiqarish (`sikltuzat/`).

**panic.** Indeks chegaradan tashqarida → runtime panic. Oddiy ish toʻxtaydi, `defer` lar baribir ishlaydi, keyin
`panic: ...` xabari, stack trace va `exit status 2`.

**recover.** Faqat `defer` qilingan funksiyaning oʻzida, toʻgʻridan-toʻgʻri chaqirilganda ishlaydi (yordamchi funksiya ichidan
chaqirilsa `nil`, uy vazifasi "Guru") (spec: *"The return
value of recover is nil when the goroutine is not panicking or recover was not called directly by a deferred
function."*). Nomlangan natija (`err`) ni defer'dagi closure oʻzgartira oladi, shuning uchun panic oddiy xatoga
aylanadi. Lekin kutilgan xatolar uchun panic emas, `error` qaytaring (GO-1.14); `recover` chegarada kerak: handler,
goroutine boshi. Go 1.21 dan beri `panic(nil)` ham `*runtime.PanicNilError` boʻlib keladi, ya'ni `recover() != nil`.

**Goroutine chegarasi.** `go f()` funksiyani yangi goroutine'da, `main` bilan parallel ishga tushiradi (batafsil:
GO-1.21). Panic faqat oʻz goroutine'sining chaqiruvlari boʻylab yuqoriga chiqadi. Goroutine boshiga
yetganda butun dastur toʻxtaydi (spec: *"the program is terminated and the error condition is reported"*). `main` dagi
`recover` yordam bermaydi.

**net/http.** Handler doc: *"If ServeHTTP panics, the server … recovers the panic, logs a stack trace to the server
error log, and either closes the network connection or sends an HTTP/2 RST_STREAM"*. Lekin handler ochgan
goroutine'dagi panic'ni hech kim ushlamaydi: butun server yiqiladi. Har bir fon goroutine'ga oʻz `recover` i kerak
(Effective Go'dagi `safelyDo` misoli; uy vazifasi "Qiyin").

**Fatal error'lar.** Deadlock (`all goroutines are asleep`), stack overflow (64 bitli tizimda 1 GB'dan oshsa), out of
memory — bular panic emas. Runtime `throw` (runtime'ning ichki xatosi) yoki `fatal` (xato foydalanuvchi kodida,
masalan GO-1.09 dagi `concurrent map writes`) chaqiradi; ikkalasi `fatalthrow` orqali xabar va stack trace chiqarib,
`exit(2)` bilan jarayonni toʻxtatadi. `defer` lar ishlamaydi, `recover` ga navbat kelmaydi. Eslatma: deadlock faqat
*hamma* goroutine qotganda aniqlanadi; serverda bitta goroutine qotsa, runtime hech narsa demaydi.

**os.Exit / log.Fatal.** `os.Exit`: *"The program terminates immediately; deferred functions are not run."*
`log.Fatal` = `Print` + `os.Exit(1)`. Yechim (`exit/tuzat/`): ishni `run() error` funksiyasiga chiqaring, uning `defer`
lari ishlab boʻlgach `main` `os.Exit(1)` ni chaqiradi.

## Manbalar
- [Go spec: Defer statements](https://go.dev/ref/spec#Defer_statements), [Handling panics](https://go.dev/ref/spec#Handling_panics), [Run-time panics](https://go.dev/ref/spec#Run_time_panics)
- Andrew Gerrand, [Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover) (Go blog, 2010)
- [Effective Go: Recover](https://go.dev/doc/effective_go#recover)
- [Go 1.14 release notes](https://go.dev/doc/go1.14) (defer "almost zero overhead"), [Go 1.21](https://go.dev/doc/go1.21) (`PanicNilError`)
- [net/http.Handler](https://pkg.go.dev/net/http#Handler), [os.Exit](https://pkg.go.dev/os#Exit), [log.Fatal](https://pkg.go.dev/log#Fatal)
- Go manba kodi (go1.26.8): [runtime/panic.go](https://github.com/golang/go/blob/go1.26.8/src/runtime/panic.go#L1199-L1241) (`throw`, `fatal`),
  [runtime/stack.go](https://github.com/golang/go/blob/go1.26.8/src/runtime/stack.go#L1200-L1208) (stack overflow),
  [runtime/proc.go](https://github.com/golang/go/blob/go1.26.8/src/runtime/proc.go#L6468) (deadlock),
  [cmd/compile/internal/walk/stmt.go](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/walk/stmt.go#L116-L120) (sikldagi defer)
- [staticcheck SA9001](https://staticcheck.dev/docs/checks/#SA9001)

## Uy vazifasi
- 🟢 **Oson:** [homework/oson](homework/oson/main.go) nima chiqaradi? Avval qogʻozda yozing, keyin ishga tushiring.
- 🟡 **Oʻrta:** faylga yozadigan `Saqla(nom, data)` yozing: `f.Close()` qaytargan xato yoʻqolmasin.
- 🟠 **Qiyin:** `xavfsizGo(nom, f)` yozing: `f` ni yangi goroutine'da ishga tushirsin, uning panic'i serverni yiqitmasin.
- 🔴 **Guru:** [homework/guru](homework/guru/main.go) dagi toʻrtta variantdan (`a`, `b`, `c`, `d`) qaysilari panic'ni ushlaydi? Nega?

Avval oʻzingiz urinib koʻring. Javoblar:
- Oson: `B: 2`, `i = 2`, `i = 1`, `i = 0`, `A: 1`. `B` closure, shuning uchun `x` ning oxirgi qiymatini koʻradi; `A` ning
  argumenti defer qatorida (`x = 1`) hisoblangan; sikldagi defer'lar teskari tartibda.
- Oʻrta ([homework/orta](homework/orta/saqla.go)): nomlangan natija + `defer func() { err = errors.Join(err, f.Close()) }()`.
  `saqlaXato` dagi oddiy `defer f.Close()` Close xatosini tashlab yuboradi. Test `/dev/full` bilan yozish xatosini tekshiradi
  (`go test -v ./homework/orta`).
- Qiyin ([homework/qiyin](homework/qiyin/main.go)): `recover` goroutine'ning **ichida**, uning birinchi `defer` ida
  boʻlishi kerak. Natija: `fon: panic ushlandi: assignment to entry in nil map`, `server tirik`.
- Guru ([homework/guru](homework/guru/main.go)): faqat `a` va `b`. `b` da `qutqar` ning oʻzi defer qilingan, `recover` uni
  toʻgʻridan-toʻgʻri chaqiradi. `c` da `recover` defer qilingan funksiyaning ichidagi boshqa funksiyadan chaqirilgan →
  `nil`. `d` da (`defer recover()`) `recover` ning oʻzi defer qilingan, uni defer qilingan funksiya chaqirmagan → `nil`
  (spec: *"recover was not called directly by a deferred function"*).
