# 01 · Go nega bor? `go build` ichida nima bor

Video: *Golang oʻzbek tilida #1* (2026-11-26, dasturchilar uchun versiya) · Lab: https://haywan.uz/lab/youtube/go#go-1-01 ·
Oldingi: [GO-X.01 · 20 qatorlik HTTP API va bitta binary](../0x-01-http-api/) · Keyingi:
[02 · Turlar va overflow](../002-types-overflow/)

**Kim uchun:** boshqa stack'da (JS/TS, Python, PHP) ishlaydigan, Go'ga yangi dasturchi. GO-X.01'dagi `go build`,
static binary, `GOOS`/`GOARCH`, `FROM scratch` va `CGO_ENABLED=0` maʼlum deb olinadi.

**Bitta gʻoya:** Go binary'sining katta qismi bizning kod emas. Har bir dasturda runtime bor; import qilingan
paketlardan linker faqat chaqiriladigan funksiyalarni va har bir paketning `init` kodini qoldiradi; tizimdagi
maʼlumot (vaqt zonalari bazasi) esa umuman kirmaydi — uni `import _ "time/tzdata"` bilan ataylab qoʻshasiz.

Hamma raqamlar `go1.26.8 linux/amd64`, node v24.10.0, Docker 29.8.1'da oʻlchangan (2026-11-26). Boshqa papkada bir
necha bayt farq qiladi: fayl yoʻllari ham binary ichida (`-trimpath` buni olib tashlaydi).

| Papka / fayl | Nima |
|---|---|
| [`paketlar.sh`](paketlar.sh) | Go binary uchun "bundle analyzer": `go tool nm -size` natijasini paketlar boʻyicha qoʻshadi |
| [`salom/`](salom/) | Salom, dunyo (7 qator): 74 B → 2 409 990 B |
| [`bosh/`](bosh/) | ichi boʻsh `func main() {}`: runtime'ning doimiy narxi (1 829 988 B) |
| [`statustext/`](statustext/) | `net/http`dan faqat `http.StatusText(404)`: tree-shaking va `init` (5 143 595 B) |
| [`vaqt/`](vaqt/) | `time.LoadLocation("Asia/Tashkent")` + `FROM scratch` → `unknown time zone` |
| [`vaqt-tz/`](vaqt-tz/) | tuzatish: `import _ "time/tzdata"` (+410 400 B) |
| [`vaqt-local/`](vaqt-local/) | zonasiz `time.Now()`: `scratch`da xatosiz, jimgina `UTC` |
| [`uyga/`](uyga/) | uy vazifasi, starter'lar va javoblar |
| `hello/`, `empty/`, `nomain/`, `bug/`, `asm/`, `homework/` | 2026-11-25 dagi birinchi video materiallari (pastdagi ilova) |

## Videodagi buyruqlar
```bash
cd ../0x-01-http-api && go build -o api .   # GO-X.01 API: 8659683 bayt
../001-why-go/paketlar.sh api 6                  # runtime 17.7 % · crypto/tls 9.4 % · net/http 8.9 % …
../001-why-go/paketlar.sh api 400 | grep ' main$' # 1595: bizning 20 qator

cd ../001-why-go/salom
go build main.go && ls -l          # main.go 74 · main 2409990
../paketlar.sh main 4              # mashina kodi: runtime 71.4 % · fmt 4.8 % · reflect 4.2 %
go tool nm -size main | grep ' main.main$'   # 83 bayt

cd ../bosh
go build main.go && ls -l main     # 1829988: hech narsa qilmaydigan dastur (Salom, dunyo'ning 76 %)

cd ../statustext                   # bu laptop'da gcc yoʻq: net uchun CGO_ENABLED=0 shart
CGO_ENABLED=0 go build -o st main.go && ls -l st   # 5143595 (toʻliq api: 8381628)
../paketlar.sh st 300 | grep ' net/http$'          # 5262 bayt (GO-X.01 api'da 231802)
../paketlar.sh st 6                                # net 121834 · math/big · crypto/internal/fips140/... (47 ta crypto/* paket)

cd ../vaqt
go build -o vaqt main.go && ./vaqt                 # Toshkent soati: 16   (laptop: tizimda zoneinfo bor)
docker build -t vaqt . && docker run --rm vaqt     # xato: unknown time zone Asia/Tashkent
cd ../vaqt-tz
go build -o vaqt main.go && docker build -t vaqt . && docker run --rm vaqt   # Toshkent soati: 16
ls -l vaqt ../vaqt/vaqt                            # 2885090 vs 2474690 (+410400)
cd ../vaqt-local && go build -o vaqt main.go && docker build -t vaqt-local . && docker run --rm vaqt-local   # 11:48 UTC (laptop: 16:48 +05)

cd ../salom
go build -ldflags="-s -w" -o main-sw main.go && ls -l main-sw   # 1601698 (−33.5 %)
```

## Asosiy tushunchalar
- **runtime — doimiy narx.** Garbage collector, scheduler, map'lar, stack boshqaruvi har bir Go binary'si ichida
  ([Go FAQ](https://go.dev/doc/faq#runtime)). Node'da ular bitta `node` faylida (V8 + libuv) va serverga bir marta
  oʻrnatiladi; Go'da har bir xizmat oʻz runtime'ini olib yuradi. Ichi boʻsh dastur: 1.8 MB (uning ichida runtime
  kodi, uning jadvallari va DWARF debug maʼlumoti; `-s -w` bilan 1.2 MB).
- **linker = tree-shaking, lekin `init` qoladi.** Linker `main.main` va barcha paketlarning `init` vazifalaridan
  boshlab yetib boriladigan funksiyalarni qoldiradi, qolganini tashlaydi
  ([`cmd/link/internal/ld/deadcode.go`](https://github.com/golang/go/blob/go1.26.8/src/cmd/link/internal/ld/deadcode.go)).
  Paket `init` funksiyalari va global oʻzgaruvchilarni boshlash main'dan oldin albatta bajariladi
  ([spec: Package initialization](https://go.dev/ref/spec#Package_initialization)), shuning uchun ular va ular
  chaqiradigan hamma narsa qoladi — bundler'dagi side effect'li modul kabi. `statustext`: `net/http` kodi 232 KB (GO-X.01 api) → 5 KB,
  lekin `net`, `math/big` va 47 ta `crypto/*` paketi init va global oʻzgaruvchilar orqali kirdi (zanjir: `main.inittasks` →
  `net/http.init` → `DefaultTransport` → `net.(*Dialer).DialContext`). Ikkinchi qoida: interfeysga saqlangan turning
  metodlari, xuddi shu nomli interfeys metodi chaqirilsa, qoladi (`math/big` shu yoʻl bilan, `-ldflags=-dumpdep`). Katta
  global `map` initsializatorlarini esa linker ishlatilmasa tashlay oladi. JS'da `"sideEffects": false` bor; Go'da bunday bayroq yoʻq.
- **Tizim maʼlumoti binary'da yoʻq.** `time.LoadLocation` bazani tartib bilan qidiradi: `$ZONEINFO` → Unix tizim
  joylari (`/usr/share/zoneinfo`, … `/etc/zoneinfo`) → `$GOROOT/lib/time/zoneinfo.zip` → `time/tzdata` (agar import
  qilingan boʻlsa) (`go doc time.LoadLocation`). `scratch`da birortasi yoʻq. Zonasiz `time.Now()` esa xato bermaydi:
  `time.Local` jimgina UTC boʻladi (`vaqt-local/`) — bu yashirinroq bug. HTTPS ildiz sertifikatlari ham tizimdan
  oʻqiladi (`crypto/x509/root_linux.go`) — GO-X.01 Chuqurroq.
- **`import _ "time/tzdata"`** — paketni nomsiz, faqat `init`i uchun ulash (`import 'dotenv/config'` kabi). `_`siz:
  `"time/tzdata" imported and not used` (Go ishlatilmagan importni xato deydi). Uning
  `init`i ichki bazani `time` paketiga roʻyxatdan oʻtkazadi. Hujjat: "about 450 KB"; bizda +410 400 B. Muqobil:
  `go build -tags timetzdata` yoki `gcr.io/distroless/static` image (ichida tzdata va CA sertifikatlari bor; bizda
  `vaqt` tzdata'siz ham `Toshkent soati: 16` chiqardi). Node'da muammo kam: `node` ichidagi ICU oʻz bazasini olib
  yuradi (`node -p process.versions.tz` → `2025b`).

## Chuqurroq
`size -A salom/main`: `.text` 643 537 (26.7 %) · `.gopclntab` 677 386 (28.1 %: funksiya nomlari va qator raqamlari;
panic stack trace, GC stack map'lari) · DWARF `.debug_*` 676 118 (28.1 %) · `.rodata` 222 082. `-ldflags="-w"` DWARF'ni,
`-s` simvollar jadvalini olib tashlaydi (Go 1.22'dan `-s` `-w`ni ham anglatadi; [cmd/link](https://pkg.go.dev/cmd/link)):
2 409 990 → 1 601 698 B. `.gopclntab` qoladi, shuning uchun stack trace oʻzgarmaydi (`uyga/qiziq`).

Nega Go importni shunday qiladi? 2007-yil Google'da C++ binary'sida `#include`lar ochilgach "a blow-up of 2000 bytes
for every C++ source byte", build esa "took 45 minutes using a precursor distributed build system". Go'da "the
compiler reads the object file for B, not its source code". Pike buni "the origin myth for Go" deydi
([Go at Google, 2012](https://go.dev/talks/2012/splash.article)). Go FAQ maqsadi: "at most a few seconds to build a
large executable on a single computer" ([FAQ](https://go.dev/doc/faq#creating_a_new_language)). (Build cache esa keyinroq, Go 1.10'da qoʻshilgan.)
Bizda `go` buyrugʻi (319 paket) boʻsh cache bilan ~10 s'da build boʻldi.

## Uy vazifasi
- 🟢 **Oson:** [`statustext/`](statustext/)da `StatusText` chaqiruvini olib tashlang, `net/http`ni `_` bilan import qiling. Build'dan oldin hajmni taxmin qiling.
- 🟡 **Oʻrta:** [`uyga/orta/`](uyga/orta/) — `vaqt` servisini `scratch`da ishlating, lekin `main.go`ga tegmang.
  Maslahat: `go doc time/tzdata`.
- 🔍 **Qiziqish:** [`uyga/qiziq/`](uyga/qiziq/) — `-ldflags="-s -w"` bilan build qilib, panic stack trace'ni solishtiring.
- ➕ **GO-X.01 takrori:** [`uyga/qoshimcha/`](uyga/qoshimcha/) — HTTPS soʻrov `scratch`da: `x509: certificate signed by unknown authority`. Tuzating.

Javoblar: [`uyga/README.md`](uyga/README.md).

## Manbalar
- Go FAQ: runtime https://go.dev/doc/faq#runtime · binary hajmi https://go.dev/doc/faq#Why_is_my_trivial_program_such_a_large_binary · nega yangi til https://go.dev/doc/faq#creating_a_new_language
- Go spec, Package initialization https://go.dev/ref/spec#Package_initialization · Import declarations https://go.dev/ref/spec#Import_declarations
- Linker dead code: https://github.com/golang/go/blob/go1.26.8/src/cmd/link/internal/ld/deadcode.go · `cmd/link` flags: https://pkg.go.dev/cmd/link
- `time.LoadLocation` https://pkg.go.dev/time#LoadLocation · `time/tzdata` https://pkg.go.dev/time/tzdata
- `net/http.StatusText` https://pkg.go.dev/net/http#StatusText
- Docker, `scratch`: https://docs.docker.com/build/building/base-images/#create-a-minimal-base-image-using-scratch · distroless: https://github.com/GoogleContainerTools/distroless
- Node Intl / ICU: https://nodejs.org/api/intl.html
- Rob Pike, *Go at Google* (2012): https://go.dev/talks/2012/splash.article

## Dars matni (video narration)

**intro**

- `hook` Oldingi kursdagi yigirma qatorlik API binary’si sakkiz megabaytdan katta edi. Frontend’da sakkiz megabaytlik bundle koʻrsangiz, darhol bundle analyzer ochasiz. Go binary’sida bu megabaytlar nima, va ulardan qanchasi bizning kodimiz? Bugun shuni oʻlchaymiz.
- `eslash` Avval oldingi kursning birinchi darsidan bitta savol. API’ni golang image’da build qilib scratch’ga qoʻydik, u esa no such file or directory dedi, fayl esa joyida edi. Sabab nima edi? Pauza qiling.
- `eslash-javob` Javob: B. cgo yoqiq qolgan, binary tizimdan libc va uning yuklovchisini qidirdi. Demak binary’ga nima kirishi va nima kirmasligi deploy’ga toʻgʻridan-toʻgʻri ta’sir qiladi. Bugun binary ichini paketlar boʻyicha oʻlchaymiz.
- `sozlar` Bugun uchta tushuncha. Birinchisi, runtime’ning doimiy narxi. Ikkinchisi, linker, ya’ni paketlarni bitta faylga yigʻuvchi qism: u Vite yoki Rollup kabi tree-shaking qiladi. Uchinchisi, import pastki chiziq: paketni faqat yon ta’siri uchun ulash, import dotenv slesh config kabi.

**analyzer**

- `savol1` Avval taxmin qiling: bizning yigirma qatorimiz shu sakkiz megabaytning qanchasi? A: bir megabayt atrofida. B: yuz kilobayt atrofida. C: ikki kilobaytdan kam. Pauza qiling.
- `analyzer` Javob: C. Goʻda tayyor bundle analyzer yoʻq, lekin go tool nm har bir funksiyaning hajmini beradi. Dars papkasidagi paketlar.sh ularni paketlar boʻyicha qoʻshadi. API’da eng kattasi runtime, keyin crypto tls va net http. Bizning yigirma qatorimiz esa bor-yoʻgʻi bir ming besh yuz toʻqson besh bayt.
- `analyzer-fayl` E’tibor bering, bu faqat mashina kodi: ikki yarim megabaytdan sal koʻp, faylning uchdan biriga yaqin. Qolgani funksiya nomlari jadvali va debug ma’lumoti, ularni Chuqurroq qismida ochamiz. Endi eng kichik dasturda runtime’ning narxini oʻlchaymiz.

**runtime**

- `build` Yetti qatorlik Salom, dunyo. main.go yetmish toʻrt bayt, go build yasagan main esa ikki million toʻrt yuz ming baytdan sal koʻp.
- `runtime` Bu yerda mashina kodining yetmish bir foizi runtime’niki, bizning main.main esa sakson uch bayt. Runtime oldingi kursdan tanish: garbage collector va scheduler. Yangisi shu: API’da ham, Salom, dunyoʻda ham u bor, bu har bir Go dasturining doimiy narxi.
- `bosh-kod` Bu doimiy narxni oʻlchash uchun fmt’ni ham olib tashlaymiz: import yoʻq, func main boʻsh, dastur hech narsa qilmaydi.
- `savol2` Uning binary’si Salom, dunyoʻning qancha qismi boʻladi? A: yigirmadan biri, chunki hajmning koʻpi fmt edi. B: toʻrtdan uchi, chunki fmt va u tortgan paketlar chorakka yaqin. C: deyarli teng, chunki fmt arzimas. Pauza qiling.
- `javob2` Javob: B. Hech narsa qilmaydigan dastur ham bir million sakkiz yuz ming baytdan sal koʻp, ya’ni toʻrtdan uchi. Bu runtime, uning jadvallari va debug ma’lumoti. fmt’ning oʻzi kichik, lekin u reflect va strconv kabi paketlarni tortgan edi.

**tree-shaking**

- `st-kod` Endi frontend’dagi tree-shaking savoli. net http’ni import qilamiz, lekin undan faqat bitta kichik funksiya, StatusText’ni chaqiramiz: u toʻrt yuz toʻrt uchun Not Found matnini qaytaradi.
- `savol-st` Binary qancha boʻladi? A: Salom, dunyoʻga yaqin, chunki linker ishlatilmaganini tashlaydi. B: oʻrtada, chunki bir qismi qoladi. C: API’ga yaqin, chunki paket butunlay kiradi. Pauza qiling.
- `javob-st` Javob: B, qariyb besh million bir yuz qirq ming bayt. net http’ning oʻz kodi esa API’dagidan qirq barobardan koʻp kichik: linker chaqirilmagan funksiyalarni haqiqatan tashladi. Unda qoʻshimcha kod qayerdan keldi?
- `init` Javob init’da. JavaScript modulining yuqori darajadagi kodi import paytida ishlaydi. Goʻda bunga har bir paketning init’i va global oʻzgaruvchilarni boshlash toʻgʻri keladi: ular main’dan oldin albatta bajariladi. Linker ularni va ular chaqiradigan hamma narsani qoldiradi. Masalan, net http’dagi DefaultTransport global oʻzgaruvchisi net paketini tortadi, shunday qilib net va qirqdan ortiq crypto paketi kirdi. package.json’dagi sideEffects false kabi bayroq Goʻda yoʻq.

**from-scratch**

- `xato-kod` Endi production’dan bir holat. Servis Toshkent soatini koʻrsatadi: time.LoadLocation vaqt zonasini yuklaydi. Xatoni oldingi kursdagidek err orqali tekshiramiz. Laptop’da soat toʻgʻri chiqadi.
- `savol4` Uni CGO_ENABLED teng nol bilan build qilib, oldingi kursdagi FROM scratch image’ga qoʻyamiz. docker run nima chiqaradi? A: Toshkent soati. B: UTC soati. C: xato. Pauza qiling.
- `xato-javob` Javob: C, unknown time zone. Vaqt zonalari bazasi binary ichida emas: LoadLocation uni tizimning zoneinfo papkasidan oʻqiydi, scratch’da esa u yoʻq. B esa boshqa holatda chiqadi: zonasiz time.Now xatosiz, jimgina UTC koʻrsatadi. Node’da bu kam uchraydi, chunki Intl ishlatadigan ICU oʻz bazasini node faylining ichida olib yuradi.
- `tuzatish` Tuzatish import guruhiga bitta qator: pastki chiziq va time slesh tzdata. Pastki chiziqsiz Go kompilyatsiya qilmaydi, imported and not used deydi, chunki paketdan hech narsa chaqirmaymiz. Bizga faqat uning init’i kerak: u binary ichidagi bazani time paketiga roʻyxatdan oʻtkazadi.
- `tuzatish-run` Qayta build qilib, scratch’da ishga tushiramiz: Toshkent soati chiqdi. Narxi qariyb toʻrt yuz oʻn ming bayt. Yoki scratch oʻrniga distroless static image oling: unda vaqt zonalari bazasi bor.

**chuqurroq**

- `chuqur-ichi` Chuqurroq qism. Salom, dunyo faylida mashina kodi chorak qism. Yana chorak qismi funksiya nomlari va qator raqamlari jadvali: runtime undan stack trace va garbage collector uchun foydalanadi. Yana chorak qismi debug ma’lumoti.
- `chuqur-sw` Minus ldflags linker’ga flag beradi. Minus w debug ma’lumotini olib tashlaydi, minus s esa nm oʻqiydigan simvollar jadvalini ham, va Go bir nuqta yigirma ikkidan beri minus w’ni ham anglatadi. Fayl hajmining qariyb uchdan biri ketadi, xuddi source map’larsiz production build kabi. Runtime ishlatadigan nomlar jadvali qoladi, stack trace oʻzgarmaydi.
- `tarix` Goʻning import modeli ham shu tezlik uchun qurilgan. Google’dagi C plyus plyus dasturida include’lar ochilgach, kompilyatorga manbadan ikki ming marta koʻp bayt borgan, build qirq besh daqiqa olgan. Go kompilyatori esa import qilingan paketning manbasini emas, tayyor obyekt faylini oʻqiydi. Rob Pike Go gʻoyasi shunday kutishda tugʻilganini kelib chiqish afsonasi deydi.

**xulosa**

- `takror` Uchta savol bilan xulosa qilamiz. Boʻsh Go dasturi nega bir million sakkiz yuz ming bayt? Bitta StatusText uchun net http import qilsak, hajm nega ikki barobar oshadi? scratch’da unknown time zone chiqsa, nima qilasiz? Pauza qiling.
- `takror-javob` Javoblar. Har bir binary oʻz runtime’ini, jadvallari va debug ma’lumoti bilan olib yuradi. Linker chaqirilmagan funksiyalarni tashlaydi, lekin paketlarning init’i va global oʻzgaruvchilari qoladi. Vaqt zonalari tizimdan oʻqiladi: import pastki chiziq time slesh tzdata yoki distroless image.

**vazifa**

- `vazifa` Uy vazifasi dars papkasida. Oson: statustext’da StatusText chaqiruvini olib, net http’ni pastki chiziq bilan import qiling; build’dan oldin hajmni taxmin qiling. Oʻrta: vaqt servisini main.goʻga tegmasdan, faqat build buyrugʻi bilan tuzating. Qiziqish uchun: minus s minus w’dan keyin crash’dagi stack trace’ni solishtiring.
- `outro` Keyingi darsda Goʻdagi son turlari va overflow’ni koʻramiz. Darsni koʻrganingiz uchun rahmat, koʻrishguncha.

---

Oldingi versiya (2026-11-26, dasturlashni noldan boshlaganlar uchun, Go 0 asosida) git tarixida: `git log -- lessons/001-why-go/README.md` (commit b452013).

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
