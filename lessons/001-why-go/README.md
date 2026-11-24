# 01 · Go nega bor? `go build` ichida nima bor

Video: *Golang o'zbek tilida #1* · Lab: https://haywan.uz/lab/youtube/go#go-1-01 · Keyingi dars:
[02 · Turlar va overflow](../002-types-overflow/)

Hamma raqamlar `go1.26.8 linux/amd64`da o'lchangan. Boshqa versiya yoki OS'da baytlar biroz farq qiladi.

| Papka | Nima |
|---|---|
| `hello/` | videodagi dastur (7 qator) |
| `empty/` | bo'sh `func main() {}`: runtime'ning o'zi qancha joy oladi |
| `nomain/` | pauza-savol: `func main` → `func Main` (ataylab build bo'lmaydi) |
| `bug/` | debugging mashqi: ikkita xato (ataylab kompilyatsiya bo'lmaydi) |
| `asm/main.main.s` | `go tool objdump -s main.main` to'liq chiqishi (20 ta instruction) |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
```bash
go build -o empty/empty ./empty && go build -o hello/hello ./hello
ls -l empty/empty hello/hello                    # 1829980 va 2409998 bayt
go run ./hello                                   # Salom, dunyo!
go run ./nomain                                  # function main is undeclared in the main package
go tool nm hello/hello | grep -c ' [Tt] '                     # 1941 ta funksiya
go tool nm hello/hello | grep -c ' [Tt] runtime\.'            # 1376 tasi runtime'niki
go tool nm -size hello/hello | grep ' main\.main$'            # main.main: 83 bayt
go tool objdump -s 'main.main$' hello/hello                   # 20 ta instruction
ldd hello/hello                                  # not a dynamic executable
GOOS=windows go build -o hello/hello.exe ./hello && ls -l hello/hello.exe
go run ./bug                                     # 2 ta compile error
```

## Dars matni

**Hook.** `main` funksiyasi bo'm-bo'sh dastur (`func main() {}`) 1 829 980 bayt. Bitta `Println` qo'shsak, 2 409 998
bayt. Shu baytlarning ichida nima bor? (Linux/amd64, Go 1.26.8. Boshqa Linux'da build qilingan papkaga qarab bir necha bayt farq
qiladi: fayl yo'llari ham binary ichida saqlanadi. Windows yoki macOS'da farq bir necha o'n KB. Windows PowerShell'da
hajmni `dir` bilan ko'ring.)

**Bu kurs.** Har darsda kod yozamiz, uni ishga tushiramiz va ichida nima bo'layotganini o'lchaymiz. Kodni bugun AI ham
yozadi. Lekin bu kod nega shunday ishlashini, qayerda sinishini va qaysi yechim yaxshiroq ekanini tushunish
engineer'ning ishi bo'lib qoladi.

**Go qayerdan chiqdi.** 2007-yil sentabr. Google'da Rob Pike katta C++ dasturida ishlayapti, va bitta build distributed
build tizimida ham 45 daqiqa davom etadi. 21-sentabrda u Ken Thompson va Robert Griesemer bilan doskaga yangi tilning
maqsadlarini yoza boshlaydi. 2009-yil 10-noyabrda Go ochiq e'lon qilinadi, 2012-yil martda Go 1 chiqadi.
[Go FAQ](https://go.dev/doc/faq#creating_a_new_language)ga ko'ra maqsad: interpretatsiya qilinadigan tillar kabi
yozish oson, lekin compile qilinadigan tillar kabi tez va xavfsiz til; katta dastur ham bitta kompyuterda bir necha
soniyada build bo'lsin; concurrency va garbage collection tilning o'zida. Pike
[2012-yilgi maqolasi](https://go.dev/talks/2012/splash.article)da yana: sodda, o'qish oson sintaksis. Docker Engine, Kubernetes, Terraform, Prometheus va etcd Go'da yozilgan: ularning kodini ochib o'qish mumkin.

**Kod.** `package main`: Go'da har bir fayl qaysidir package'ga tegishli, `main` esa maxsus nom: bu library emas,
ishga tushadigan dastur. `import "fmt"`: matnni formatlab chiqaradigan standart package. `func main()`: bizning kodimiz shu
yerdan boshlanadi. `fmt.Println` ekranga bitta qator yozadi.

**`go run`.** Kodni compile qiladi, binary yasaydi va uni ishga tushiradi. Go 1.24'dan beri bu binary build
cache'da saqlanadi, shuning uchun keyingi `go run` tezroq boshlanadi.

**Pauza-savol.** `package main` qoladi, lekin `func main`ni `func Main` deb o'zgartirsak? Build bo'lmaydi, xatoni
compiler emas, *linker* aytadi: `runtime.main_main·f: function main is undeclared in the main package`. Demak dastur aslida
runtime'dan boshlanadi: runtime avval o'zini tayyorlaydi, keyin bizning `main`ni chaqiradi
([runtime/proc.go](https://github.com/golang/go/blob/go1.26.8/src/runtime/proc.go#L136)).

**`go build` ichida** ([cmd/compile/README.md](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/README.md)):
1. Parsing: matn token'larga bo'linadi, har bir fayldan syntax tree quriladi.
2. Type checking: turlar tekshiriladi.
3. IR (compiler'ning ichki ko'rinishi) va middle end: inlining, devirtualization, escape analysis.
4. SSA: yana optimizatsiyalar, keyin machine code (amd64, arm64, …).
5. Linker: bizning kodga `runtime` va `fmt`ni qo'shadi, hech kim chaqirmaydigan funksiyalarni tashlab yuboradi
   ([deadcode.go](https://github.com/golang/go/blob/go1.26.8/src/cmd/link/internal/ld/deadcode.go)).

Compiler qarorlarini ko'rish: `go build -gcflags=-m -o /dev/null ./hello` → `inlining call to fmt.Println`,
`"Salom, dunyo!" escapes to heap`.

**Hajm.** 2 409 998 bayt. Binary'da 1941 ta funksiya bor, ulardan 1376 tasi `runtime` package'iniki (71 %).
`main.main` esa 83 bayt. Faylning faqat chorak qismi machine code (uning uchdan ikki qismidan ko'pi runtime). Yana
chorak qismi `.gopclntab`: har bir funksiyaning nomi, qator raqamlari va stack xaritasi. GC va stack o'sishi shu
jadvalsiz ishlamaydi, panic paytidagi stack trace ham shundan. Yana chorak qismi DWARF: debugger uchun, dastur
ishlashi uchun kerak emas (`-ldflags=-w` bilan o'chadi). Qolgan oltidan bir qismi: turlar haqidagi ma'lumot
(`.rodata`, 222 338 bayt) va symbol jadvali.

Runtime: har bir Go dasturi ichiga qo'shiladigan yordamchi tizim. Garbage collector keraksiz xotirani tozalaydi,
scheduler goroutine'larni operatsion tizim thread'lariga taqsimlaydi. [Go FAQ](https://go.dev/doc/faq#Why_is_my_trivial_program_such_a_large_binary)ga
ko'ra, C'dagi static Hello World ≈ 750 KB; Go'niki kattaroq, chunki ichida runtime,
turlar haqidagi ma'lumot, pclntab va debug ma'lumoti bor.

| Bo'lim (`size -A`) | Bayt |
|---|---|
| `.gopclntab` (nomlar, qatorlar) | 677 386 |
| `.text` (machine code) | 643 537 |
| `.debug_*` (DWARF, jami) | ≈ 676 000 |
| `.rodata` | 222 338 |

**Machine code.** `main.main`ning birinchi instruction'i `CMPQ SP, 0x10(R14)`: stack yetarlimi? Goroutine kichik stack bilan
boshlanadi, kamida 2 KB ([`stackMin = 2048`](https://github.com/golang/go/blob/go1.26.8/src/runtime/stack.go#L78);
Go 1.19'dan beri boshlang'ich hajm o'rtacha foydalanishga qarab moslashadi). Bu Go'ning Plan 9 assembly sintaksisi:
`MOVQ SP, BP` = "SP'ni BP'ga ko'chir" (Intel sintaksisiga teskari tartib). `Println` ko'rinmaydi: compiler uni inline
qilgan.
Joy yetmasa, oxiridagi `runtime.morestack` stack'ni kattalashtiradi va `main` boshidan qayta boshlanadi. Bu
tekshiruvni compiler o'zi qo'shadi.

**Static binary.** `ldd`: *not a dynamic executable*. Binary tizimdagi umumiy library'larga bog'lanmagan, uni boshqa
amd64 Linux serverga ko'chirsangiz ham ishlaydi. `GOOS=windows` bilan shu Linux'ning o'zida Windows uchun `.exe` (2 473 472
bayt) chiqadi. Eslatma: cgo ishlatadigan dasturlar (masalan, ba'zi tizimlarda `net` yoki `os/user`) dynamic bo'lishi
mumkin; bizning `hello` faqat `fmt`ni ishlatadi.

**Debugging.** `bug/main.go`da `os` import qilingan va funksiya ichidagi `name` e'lon qilingan, lekin ikkalasi ishlatilmagan
(package darajasidagi o'zgaruvchi va funksiya parametri ishlatilmasa, xato emas). Go buni
warning emas, compile error deydi. Sababi [Go FAQ](https://go.dev/doc/faq#unused_variables_and_imports)da: ishlatilmagan
o'zgaruvchi bug belgisi bo'lishi mumkin, ishlatilmagan import esa build'ni sekinlashtiradi.

## Manbalar
- Go FAQ, loyiha tarixi: https://go.dev/doc/faq#history
- Go FAQ, nega yangi til: https://go.dev/doc/faq#creating_a_new_language
- Go FAQ, nega binary katta: https://go.dev/doc/faq#Why_is_my_trivial_program_such_a_large_binary
- Rob Pike, *Go at Google: Language Design in the Service of Software Engineering* (2012): https://go.dev/talks/2012/splash.article
- Go 1 chiqdi (2012-03-28): https://go.dev/blog/go1
- Go 1.24 release notes, `go run` cache: https://go.dev/doc/go1.24#go-command
- Go compiler bosqichlari: https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/README.md
- Goroutine stack'ining boshlang'ich hajmi: https://github.com/golang/go/blob/go1.26.8/src/runtime/stack.go#L78
- Go FAQ, ishlatilmagan o'zgaruvchi va import: https://go.dev/doc/faq#unused_variables_and_imports
- CMU 15-213, *Linking* ma'ruzasi (static va dynamic linking haqida chuqurroq): https://www.cs.cmu.edu/~213/schedule.html

## Uy vazifasi
- 🟢 **Oson:** ismingiz va yoshingizni chiqaradigan dastur yozing.
- 🟡 **O'rta:** dasturni Windows, Linux va macOS uchun build qiling va hajmlarini solishtiring.
- 🟠 **Qiyin:** `-ldflags=-s` bilan build qiling. Nima olib tashlandi, va nega `.gopclntab` qoldi?
- 🔴 **Guru:** `go tool objdump` bilan `main.main`ni oching va har bir instruction nima qilishini tushuntiring.

Avval o'zingiz urinib ko'ring. Javoblar: [homework/](homework/).
