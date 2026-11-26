# 01 · Go nega bor? `go build` ichida nima bor

Video: *Golang oʻzbek tilida #1* (2026-11-26 yangi versiya, Go 0 mavsumidan keyin) · Lab:
https://haywan.uz/lab/youtube/go#go-1-01 · Oldingi: [Go 0 · 5 — for](../000-05-for/) · Keyingi:
[02 · Turlar va overflow](../002-types-overflow/)

**Bitta gʻoya:** `go build` kodimizni tarjima qilib, alohida tayyor faylga (*binary*) saqlaydi. Go bu fayl ichiga har
safar oʻz yordamchisini (*runtime*) ham joylaydi: shuning uchun eng kichik dastur ham qariyb ikki million bayt boʻladi
va Go oʻrnatilmagan kompyuterda ham ishlaydi. Kodni oʻzgartirsangiz, faylni qayta yasash kerak.

Hamma raqamlar `go1.26.8 linux/amd64`da oʻlchangan. Boshqa papkada bir necha bayt, boshqa tizimda bir necha oʻn KB
farq qiladi (fayl yoʻllari ham binary ichida saqlanadi).

| Papka | Nima |
|---|---|
| [`salom/`](salom/) | videodagi dastur (7 qator, Go 0 · 1 dagi Salom, dunyo) |
| [`bosh/`](bosh/) | ichi boʻsh `func main() {}`: runtime'ning oʻzi qancha joy oladi |
| [`uyga/`](uyga/) | uy vazifasi: starter fayllar va javoblar |
| `hello/`, `empty/`, `nomain/`, `bug/`, `asm/`, `homework/` | eski (2026-11-25) video materiallari, pastdagi ilovada |

## Videodagi buyruqlar
```bash
cd salom
go build main.go        # ekranda hech narsa chiqmaydi: tarjima main fayliga saqlandi
ls -l                   # main.go 74 bayt · main 2409990 bayt (Windows cmd: dir)
./main                  # Salom, dunyo!     (Windows cmd: main · PowerShell: .\main)

cd ../bosh
go build main.go && ls -l main     # 1829988 bayt: ichi boʻsh dastur ham ~1.8 million bayt

cd ../salom             # main.go'da dunyo → Aziza qilib saqlang, lekin build qilmang:
./main                  # Salom, dunyo!   ← eski fayl!
go build main.go && ./main         # Salom, Aziza!
```
Videodan keyin `salom/main.go`ni yana `Salom, dunyo!` holiga qaytaring.

## Asosiy tushunchalar
- **`go run`** = tarjima + darhol ishga tushirish (Go 0'da shuni ishlatdik). Tarjimani sizga bermaydi.
- **`go build`** = tarjimani faylga saqlaydi ("kitob qilib chop etadi"). Muvaffaqiyatli boʻlsa, jim turadi.
- **binary** (tayyor fayl): `main.go`dan yasalgan `main` (Windows'da `main.exe`). Ichida kompyuter tili.
- **runtime** (ichki yordamchi): har bir Go dasturi ichiga qoʻshiladigan yordamchi qism: xotirani tozalaydi
  (garbage collector), ishlarni navbatga qoʻyadi (scheduler), stack'ni boshqaradi
  ([Go FAQ](https://go.dev/doc/faq#runtime)).
- **static binary**: kerakli hamma narsa bitta faylda; Go oʻrnatilmagan, lekin **xuddi shu turdagi** (tizim va
  protsessor, masalan, amd64 Linux) kompyuterda ishlaydi (`ldd main` → *not a dynamic executable*). Eslatma: toʻliq
  static faqat Linux'da va cgo'siz; macOS'da Go fayllari tizim kutubxonasiga (libSystem), Windows'da tizim DLL'lariga
  bogʻlanadi — lekin Go'ning oʻzi hech qayerda kerak emas.
- **Eng koʻp xato:** `main.go`ni oʻzgartirib, `go build`ni qayta bajarmaslik: `./main` eski matnni chiqaradi.

## Nega Go?
2007-yil, Google: katta C++ dasturining bitta build'i ~45 daqiqa davom etgan
([Pike, 2012](https://go.dev/talks/2012/splash.article)); Pike'ning oʻzi aytishicha, Go gʻoyasi shunday kutishlardan
birida tugʻilgan (u buni Go'ning *origin myth*i deb ataydi). 2007-yil 21-sentabrda Robert Griesemer, Rob Pike va Ken
Thompson doskaga yangi tilning maqsadlarini yoza boshlashdi; 2009-yil 10-noyabrda Go ochiq eʼlon qilindi, 2012-yil
martda Go 1 chiqdi ([Go FAQ: history](https://go.dev/doc/faq#history), [go.dev/blog/go1](https://go.dev/blog/go1)).
Maqsad: interpretatsiya qilinadigan tillar kabi yozish oson, compile qilinadigan tillar kabi tez va xavfsiz, katta
dastur ham bitta kompyuterda bir necha soniyada build boʻlsin
([Go FAQ](https://go.dev/doc/faq#creating_a_new_language)). Bizda: `go` buyrugʻining oʻzi (319 paket) boʻsh cache bilan ~10 s'da
build boʻldi; bizning 7 qatorlik dastur birinchi marta 2.8 s, keyingisi ~0.1 s (build cache tayyor qismlarni eslab qoladi). Bugun Docker Engine, Kubernetes, Terraform, Prometheus Go'da yozilgan.

## Chuqurroq (ixtiyoriy)
`go build` = kompilyator (har bir paketni alohida tarjima qiladi) + **linker** (ulovchi: hammasini bitta faylga
tikadi, `runtime` va `fmt`ni qoʻshadi, hech kim chaqirmaydigan funksiyalarni tashlab yuboradi). `size -A main`:

| Boʻlim | Bayt | Nima |
|---|---|---|
| `.text` | 643 537 | kompyuter buyruqlari (machine code); uning 71 % runtime'niki, bizning `main.main` — 83 bayt |
| `.gopclntab` | 677 386 | funksiya nomlari va qator raqamlari jadvali: panic stack trace, GC va stack oʻsishi shundan foydalanadi |
| `.debug_*` (DWARF) | 676 118 | debugger uchun; `-ldflags=-w` bilan olib tashlanadi, dastur ishlashiga taʼsir qilmaydi |
| `.rodata` | 222 082 | turlar haqidagi maʼlumot, matnlar |

Batafsil (compiler bosqichlari, `go tool nm`, `objdump`, `ldd`, stack tekshiruvi) — pastdagi ilovada.

## Uy vazifasi
- 🟢 **Oson:** `dunyo` oʻrniga ismingiz → `go build main.go` → `./main`. Starter: [`uyga/oson/`](uyga/oson/)
- 🟡 **Oʻrta:** `for` bilan uch marta `Salom!`; build'dan oldin taxmin qiling: fayl qanchaga oʻsadi?
  Starter: [`uyga/orta/`](uyga/orta/)
- 🔍 **Qiziqish:** `main` faylini boshqa papkaga koʻchirib ishga tushiring — ishlaydimi? Qoʻshimcha: boshqa tizim uchun fayl yasang (`GOOS`).

Javoblar: [`uyga/README.md`](uyga/README.md).

## Manbalar
- Go FAQ, runtime: https://go.dev/doc/faq#runtime
- Go FAQ, nega binary katta: https://go.dev/doc/faq#Why_is_my_trivial_program_such_a_large_binary
- Go FAQ, loyiha tarixi: https://go.dev/doc/faq#history · nega yangi til: https://go.dev/doc/faq#creating_a_new_language
- Rob Pike, *Go at Google* (2012): https://go.dev/talks/2012/splash.article
- `go help build`, `go help run`: https://pkg.go.dev/cmd/go
- Go tutorial, *Compile and install the application*: https://go.dev/doc/tutorial/compile-install
- Go compiler bosqichlari: https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/README.md
- Harvard CS50x, lecture 1 (`make hello` → `./hello`): https://cs50.harvard.edu/x/notes/1/

## Dars matni (video narration)

**intro**

- `hook` Bir varaq retseptni tarjimaga bersangiz-u, tarjimon sizga qalin kitob qaytarsa, hayron qolasiz. Go yetti qatorlik dasturimizni tarjima qilganda aynan shunday qiladi: tayyor fayl asl matndan qariyb oʻttiz ikki ming marta katta. Bugun sababini topamiz.
- `eslash` Avval Go nol mavsumini eslaymiz. U yerda har safar go run deb yozdik. go run ichida qaysi ikki ish boʻlardi? Pauza qiling va tanlang.
- `eslash-javob` Javob: B. go run avval main.go’ni kompilyator bilan tarjima qiladi, keyin darhol ishga tushiradi. Bugun shu ikki ishni ajratamiz.
- `sozlar` Buning uchun bugun uchta soʻz kerak. go build, yaʼni tarjimani faylga saqlash. binary, yaʼni tayyor fayl. Va runtime, yaʼni har bir Go dasturi ichidagi yordamchi.

**go-build**

- `build` go run tarjimani oʻzi ishlatadi va sizga bermaydi. Baʼzan esa tarjimani saqlab qoʻyish kerak: keyin ishlatish yoki boshqa kompyuterga berish uchun. Tasavvur qiling: tarjimon matnni kitob qilib chop etadi va qoʻlingizga beradi. Go’da buni qiladigan buyruq go build deyiladi.
- `binary` Chop etilgan kitob bu yangi fayl. main.go’da odam oʻqiydigan harflar bor, yangi faylda esa kompyuter bajaradigan buyruqlar. Bu fayl inglizcha binary, oʻzbekcha tayyor fayl deyiladi. go build main.go desak, uning nomi main boʻladi.
- `savol1` Endi papkada go build main.go deb yozamiz. Ekranda nima chiqadi? A: Salom, dunyo. B: hech narsa. C: xato xabari. Pauza qiling va oʻylab koʻring.
- `javob1` Javob: B, hech narsa. go build faqat tarjima qiladi, ishga tushirmaydi, ish yaxshi boʻlsa esa jim turadi. Lekin papkada yangi fayl paydo boʻldi. Biz Linux’da ishlayapmiz, u ham Windows va Mac kabi kompyuterning asosiy tizimi. Bu yerda ls minus l buyrugʻi fayllarni hajmi bilan koʻrsatadi, Windows’da esa dir deb yozasiz.
- `hajm` Har qatorda oxirgi soʻz fayl nomi, sanadan oldingi son esa hajm, boshidagi belgilarga hozircha qaramang. Hajm baytda oʻlchanadi. Bayt bu sakkizta nol yoki birdan iborat kichik joy, unga bitta oddiy harf sigʻadi. main.go yetmish toʻrt bayt, boʻsh joylari bilan. Tayyor main esa qariyb ikki million toʻrt yuz ming bayt.
- `ishga` Endi tayyor faylni ishga tushiramiz: nuqta slesh main, yaʼni shu papkadagi main fayli. Windows’da esa shunchaki main deb yozasiz. Salom, dunyo chiqdi. Qisqasi, go build tarjimani faylga saqlaydi, kompyuter esa shu tayyor faylni bajaradi.

**runtime**

- `runtime` Bir qator ish uchun shuncha bayt nimaga kerak? Oshxonada retseptdan tashqari yordamchi ish ham bor: idish yuvish. Dastur ishlaganda ham shunday: keraksiz boʻlib qolgan qiymatlarni xotiradan tozalash kerak. Go bu yordamchi ishlarning buyruqlarini har bir tayyor faylga qoʻshadi, ular runtime deyiladi.
- `kitob` Demak, main faylidagi buyruqlar uch joydan keladi. Bizning yetti qatorimizning tarjimasi, bor-yoʻgʻi sakson uch bayt. Println’ni bergan fmt qutisi. Va eng kattasi, runtime: buyruqlarning toʻrtdan uch qismiga yaqini uniki.
- `bosh-kod` Bu ichi boʻsh dastur: tanish package main va func main bor, lekin na import, na Println bor. Jingalak qavslar orasida hech qanday buyruq yoʻq, u hech narsa qilmaydi.
- `savol2` Uni go build qilsak, tayyor fayl qancha boʻladi? A: deyarli nol bayt. B: yetmish bayt atrofida. C: million baytdan katta. Pauza qiling va oʻylab koʻring.
- `javob2` Javob: C. ls minus l koʻrsatadi: hech narsa qilmaydigan dastur ham qariyb bir million sakkiz yuz ming bayt. Println va fmt esa qariyb besh yuz sakson ming bayt qoʻshadi. Qisqasi, runtime har qanday Go dasturi ichida bor.
- `savol3` Endi tezkor savol. Tayyor main faylini doʻstingizning kompyuteriga koʻchiramiz. Unda ham Linux bor, lekin Go oʻrnatilmagan. Ishlaydimi? Pauza qiling va tanlang.
- `javob3` Javob: A, ishlaydi. Kerakli hamma narsa, runtime ham, fmt ham, bitta main faylining ichida, Go’ning oʻzi kerak emas. Faqat kompyuter xuddi shunday boʻlsin: masalan, Linux uchun yasalgan fayl Windows’da ishlamaydi.

**nega-go**

- `tarix` Go nega yaratilgan? Ikki ming yettinchi yilda Google'dagi katta dasturning bitta tarjimasi qirq besh daqiqa davom etgan. Hikoyaga koʻra, Go gʻoyasi shunday kutishda tugʻilgan: yozish oson boʻlsin, katta dastur ham bir necha soniyada tarjima qilinsin.

**xato**

- `xato-kod` Endi eng koʻp uchraydigan xato. main.go’dagi fmt nuqta Println qatorida dunyo soʻzini Aziza deb oʻzgartirdik va saqladik. Lekin go build’ni qayta yozmadik, darhol nuqta slesh main deb ishga tushiramiz.
- `savol4` Ekranda nima chiqadi? A: Salom, Aziza. B: Salom, dunyo. C: xato xabari. Kitobni eslang: u chop etilgandan keyin biz retseptni oʻzgartirdik. Pauza qiling va oʻylab koʻring.
- `xato-javob` Javob: B, eski Salom, dunyo. Tayyor fayl tarjima qilingan paytdagi koddan yasalgan. main.go’ni oʻzgartirsangiz, chop etilgan kitob oʻzgarmaydi, ular ikki alohida fayl.
- `tuzatish` Tuzatish oson: avval go build, keyin nuqta slesh main. Endi Salom, Aziza chiqdi. Qoida: kodni oʻzgartirdingizmi, faylni qayta yasang. go run buni har safar oʻzi qiladi, shuning uchun Go nol darslarida bu xato boʻlmagan.

**chuqurroq**

- `chuqur` Endi ixtiyoriy Chuqurroq qism: endi boshlayotgan boʻlsangiz, uni keyinroq koʻrsangiz ham boʻladi. Avval kompilyator main.go’ni, fmt qutisini va runtime’ni alohida tarjima qiladi. Tayyor qismlarni u eslab qoladi, keyingi safar ularni qayta tarjima qilmaslik uchun. Keyin ulovchi, inglizcha linker. U hammasini bitta main fayliga tikadi.
- `chuqur-ichi` Faylning qariyb chorak qismi kompyuter buyruqlari. Yana chorak qismi buyruqlar guruhlarining nomlari va qator raqamlari: dastur ishlayotganda xato bilan toʻxtasa, qaysi qatorda toʻxtaganini shundan aytadi, runtime ham undan foydalanadi. Yana chorak qismi dasturchi xato qidirganda kerak boʻladigan maʼlumot, usiz ham dastur ishlaydi.

**xulosa**

- `takror` Endi uchta savol bilan xulosa qilamiz. go run bilan go build farqi nima? Ichi boʻsh dastur nega million baytdan katta? Kodni oʻzgartirdingiz, endi nima qilasiz? Pauza qiling va javob bering.
- `takror-javob` Javoblarni tekshiramiz. go run tarjima qilib darhol ishga tushiradi, go build esa tarjimani tayyor faylga saqlaydi. Har bir tayyor faylning ichida runtime bor. Kod oʻzgarsa, go build’ni qayta bajaring.

**vazifa**

- `vazifa` Uy vazifasi uch qismli. Oson: dunyo oʻrniga ismingizni yozing va build qilmasdan natijani taxmin qiling. Oʻrta: for bilan uch marta Salom chiqaring, fayl qanchaga kattalashishini taxmin qiling. Qiziqish uchun: main faylini boshqa papkaga koʻchirib, shu papkada ishga tushiring. Javoblar dars papkasida.
- `outro` Keyingi darsda Go’da son qachon toʻlib qolishini koʻramiz. Darsni koʻrganingiz uchun rahmat, koʻrishguncha.


---

## Ilova: eski video (2026-11-25)

Birinchi versiya (allaqachon dasturlay oladiganlar uchun, zich): compiler bosqichlari, `go tool nm`, binary boʻlimlari,
`objdump`, `ldd`, ishlatilmagan import. Papkalar: `hello/`, `empty/`, `nomain/`, `bug/`, `asm/`, `homework/`.

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

### Videodagi buyruqlar
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

### Dars matni

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

### Manbalar
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

### Uy vazifasi
- 🟢 **Oson:** ismingiz va yoshingizni chiqaradigan dastur yozing.
- 🟡 **O'rta:** dasturni Windows, Linux va macOS uchun build qiling va hajmlarini solishtiring.
- 🟠 **Qiyin:** `-ldflags=-s` bilan build qiling. Nima olib tashlandi, va nega `.gopclntab` qoldi?
- 🔴 **Guru:** `go tool objdump` bilan `main.main`ni oching va har bir instruction nima qilishini tushuntiring.

Avval o'zingiz urinib ko'ring. Javoblar: [homework/](homework/).
