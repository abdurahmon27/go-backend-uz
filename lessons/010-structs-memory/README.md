# 10 · Struct'lar va xotira: bir xil maydonlar, 24 va 16 bayt

Oldingi dars: GO-1.09 · Map: tartib nega tasodifiy va `recover` nega qutqarmaydi · Keyingi dars: [11 · Pointer'lar: bu oʻzgaruvchi nega heap'ga qochdi?](../011-pointers-escape/)

| Papka | Nima |
|---|---|
| `olcham/` | maydonlari bir xil ikki struct: `Buyurtma` 24 bayt, `Ixcham` 16 bayt |
| `literal/` | struct eʼlon qilish, kalitli composite literal, zero value, embedding |
| `xarita/` | `unsafe.Sizeof/Alignof/Offsetof` bilan har bir baytning xaritasi (`·` = padding) |
| `savol/` | pauzadagi savolning javobi: `bool, int32, bool` = 12 bayt, qayta tartiblansa 8 |
| `million/` | 1 000 000 ta buyurtma uchun ajratilgan xotira (24,0 va 16,0 MB) + `go test` benchmark |
| `tezlik/` | xuddi shu hisob, ikkala variant navbatma-navbat 100 marta, eng yaxshi natija (videoda shu); `go run ./tezlik 100000` — keshga sigʻadigan hajm |
| `atom/` | `GOARCH=386` da `bool` dan keyingi `int64` ga `atomic.AddInt64`: **ataylab xato** (panic) |
| `atomfix/` | `atomic.Int64` (Go 1.19+): 32 bitli tizimda ham 8 baytga tekislanadi |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/010-structs-memory`) ishga tushiriladi.
```bash
sed -n 10,21p olcham/main.go
go run ./olcham                  # Buyurtma: 24 bayt, Ixcham: 16 bayt
go run ./literal
go run ./xarita                  # [T·······SSSSSSSSY·······] va [SSSSSSSSTY······]
go run ./savol                   # 12 va 8
go run ./million                 # 24.0 MB va 16.0 MB
go run ./tezlik                  # Ixcham / Buyurtma: ≈0.72–0.85 (bizda)
go run ./tezlik 100000           # Ixcham / Buyurtma: ≈0.9–1.2: farq kichik
GOARCH=386 go run ./xarita       # 32 bit: 16 va 12 bayt, align 4
GOARCH=386 go run ./atom 2>&1 | head -2   # panic: unaligned 64-bit atomic operation
GOARCH=386 go run ./atomfix
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest   # bizda v0.53.0
$(go env GOPATH)/bin/fieldalignment ./olcham 2>&1 | sed "s|$PWD/||"
```
`GOARCH=386` bilan yigʻilgan dastur 64 bitli Linux'da odatda toʻgʻridan-toʻgʻri ishlaydi (32 bitli dasturlarni
ishga tushirish yoqilgan boʻlsa). Natijalar `go1.26.8 linux/amd64`, Intel i5-12450H (L3 kesh 12 MiB) da olingan.

## Dars matni

**Struct va composite literal.** Struct bir nechta maydonni bitta turga birlashtiradi. Qiymat composite literal
bilan yaratiladi: `Mijoz{Ism: "Aziza", Balans: 50_000}`. Kalitli (maydon nomi bilan) literalda tartib muhim emas,
yozilmagan maydonlar zero value oladi (`a` da `Manzil` yozilmagan), `var b Mijoz` da esa hamma maydon. Kalitsiz `Mijoz{"Aziza", 50_000, Manzil{}}` ham bor, lekin u hamma maydonni
tartib bilan talab qiladi va struct'ga yangi maydon qoʻshilsa buziladi (`go vet` boshqa paketdagi struct uchun
kalitsiz literalni ogohlantiradi: `composites` tekshiruvi).

**Embedding.** `Manzil` alohida maydon nomisiz yozilgan: bu embedded field, tur nomi maydon nomi boʻladi (spec: *"The unqualified type name acts as the field name"*). Uning maydonlari `Mijoz` ga "koʻtariladi"
(spec: *promoted*): `a.Shahar` aslida `a.Manzil.Shahar`. Bu meros emas: `Mijoz` `Manzil` turiga aylanmaydi, faqat
koʻtarilgan maydonlar (Effective Go: *"There's an important way in which embedding differs from subclassing"*).
Metodlar ham xuddi shunday koʻtariladi va bu tur qaysi interface'larni bajarishini oʻzgartiradi: GO-1.12, GO-1.13.

**Padding.** `unsafe.Sizeof` turning oʻlchamini padding bilan birga beradi, `unsafe.Alignof` manzil nechaga karrali
boʻlishi kerakligini, `unsafe.Offsetof` maydon struct boshidan necha bayt narida turganini. `xarita/` natijasi:
```
Buyurtma  24 bayt, align 8  [T·······SSSSSSSSY·······]
Ixcham    16 bayt, align 8  [SSSSSSSSTY······]
```
`T` = `Tolandi` (bool), `S` = `Summa` (int64), `Y` = `Yetkazildi` (bool), `·` = padding. `int64` 8 ga karrali
manzildan boshlanadi, shuning uchun birinchi `bool` dan keyin 7 bayt boʻsh. Struct'ning alignment'i maydonlar
alignment'larining eng kattasi (8; eng *katta maydon* emas: `struct{A [16]byte; B int32}` ning align'i 4), oʻlchamini
esa kompilyator 8 ga karrali qilib yaxlitlaydi: `[]Buyurtma` da har bir elementning `Summa` si ham
toʻgʻri manzilda boʻlishi uchun. Go spec faqat minimal kafolatlarni beradi (sonli turlarning oʻlchami, struct alignment'i = maydonlar alignment'larining eng kattasi); maydonlarni manba tartibida joylash — `gc` kompilyatorining amaldagi
xatti-harakati. Spec maydonlar tartibini vaʼda qilmaydi (golang/go#10014), C bilan bir xil joylashuv kerak boʻlsa,
Go 1.23 dan beri `structs.HostLayout` belgisi bor.

**Nega tekislash kerak.** Protsessor xotirani bayt-bayt emas, boʻlaklab (soʻz, kesh qatori) oʻqiydi. CMU 15-213
qoidasi: K baytli oddiy tur K ga karrali manzilda turishi kerak (slaydda: *"Required on some machines; advised on
IA32"*). CS:APP3e §3.9.3: x86-64 tekislanmagan maʼlumot bilan ham toʻgʻri ishlaydi, lekin Intel unumdorlik uchun
tekislashni tavsiya qiladi. Tekislanmagan qiymat ikki boʻlak (kesh qatori, sahifa) chegarasini kesib oʻtishi mumkin.
32 bitli Linux'da C ham 8 baytli `double` uchun 4 baytlik tekislashni ishlatgan (CMU slaydi, IA32),
Go ham `GOARCH=386` da shunday.

**Savol javobi.** `struct { A bool; B int32; C bool }`: `A` 0-baytda, 3 bayt padding, `B` 4–7, `C` 8, jami 9 → 4 ga
karralisi **12**. `B` ni boshiga qoʻysak: 4 + 1 + 1 = 6 → **8**.

**Qachon muhim emas.** `go run ./tezlik 100000` (2,4 va 1,6 MB, keshga sigʻadi): nisbat kompyuter boʻsh paytda 1,1–1,2 (ixcham *sekinroq*), yuklama ostida 0,89–1,08: farq ancha kichik. Bir million elementda esa 0,72–0,85. Bu protsessorda P- va E-yadrolar bor: natija qaysi yadroda ishlaganiga ham bogʻliq. Kritik koʻrib chiqishda 1 000–10 000 elementda ham 1,2–1,4 chiqdi. Ehtimoliy sabab: kompilyator
24 va 16 baytli qadam uchun turlicha kod chiqaradi; keshga sigʻadigan maʼlumotda bu baytlar sonidan muhimroq.

**Bir million buyurtma.** `million/`: `make([]Buyurtma, 1e6)` 24,0 MB, `make([]Ixcham, 1e6)` 16,0 MB
(`runtime.MemStats.TotalAlloc`). Toʻlangan summani hisoblash: `tezlik/` (eng yaxshi 100 oʻlchov) bizda
≈1,1 va ≈0,9 ms, nisbat ≈0,8. Kompyuter boʻsh paytda `go test -bench . -count 10 ./million` medianasi ≈1,15 va
≈0,82 ms (nisbat ≈0,7) boʻldi; boshqa ogʻir ish (xotirani koʻp ishlatadigan jarayon) ketayotganda `go test`
natijalari juda sakraydi, hatto tartib almashadi. Shuning uchun videoda navbatma-navbat oʻlchaydigan `tezlik/`.
24 MB ham, 16 MB ham L3 keshga (12 MiB) sigʻmaydi: bu yerda asosiy ish xotiradan oʻqish, ixcham variant esa uchdan
bir qism kam bayt oʻqiydi. Vaqt baytlar bilan aniq proporsional emas (≈20–30 %, 33 % emas).

**32 bit.** `GOARCH=386` da `int64` ning alignment'i 4: `Buyurtma` 16, `Ixcham` 12 bayt. Shu sababli 32 bitli
tizimda `bool` dan keyingi `int64` 4-baytdan boshlanadi va `atomic.AddInt64` panic beradi:
`panic: unaligned 64-bit atomic operation` (bu protsessor xatosi emas: Go runtime'ining `atomic_386.s` dagi
`Xadd64` manzilni tekshiradi va `panicUnaligned` ni chaqiradi). `sync/atomic` hujjatining Bugs boʻlimi: *"On ARM, 386, and 32-bit MIPS,
it is the caller's responsibility to arrange for 64-bit alignment of 64-bit words accessed atomically via the
primitive atomic functions (types Int64 and Uint64 are automatically aligned)."* Go 1.19 release notes:
*"Int64 and Uint64 are automatically aligned to 64-bit boundaries in structs and allocated data, even on 32-bit
systems."* Eski yechim: 64 bitli maydonni struct boshiga qoʻyish. Lekin `sync/atomic` hujjati faqat *"The first word in an allocated struct, array, or slice; in a global variable; or in a local variable"* tekislanganini kafolatlaydi: ichma-ich struct yoki struct'lar massivida bu yetmasligi mumkin (golang/go#11891, GitLab hodisasi pastda).

**Haqiqiy hodisa.** GitLab CE 15.7.0 ga yangilagan Raspberry Pi OS 32-bit foydalanuvchilarida Gitaly
`panic: unaligned 64-bit atomic operation` bilan yiqildi; birinchi maslahat 15.6.3 ga qaytish edi
([gitaly#4701](https://gitlab.com/gitlab-org/gitaly/-/issues/4701), 2022-12-27). Sabab
([MR !5223](https://gitlab.com/gitlab-org/gitaly/-/merge_requests/5223)): `ObjectHash` ga `ProtoFormat` maydoni
qoʻshilgach, uni oʻz ichiga olgan `requestQueue` dagi `int64` hisoblagich siljidi. Hisoblagichni birinchi qoʻyish yetmadi (!5223: *"Moving `outstandingRequests` first in the struct isn't enough"*, chunki struct'lar massivida faqat birinchi element kafolatlanadi), shuning uchun hisoblagichlar alohida struct'ga ajratildi va tekislanganini tekshiradigan test qoʻshildi; ikkinchi joy [!5226](https://gitlab.com/gitlab-org/gitaly/-/merge_requests/5226) da tuzatildi (*"We need to reorder the fields to ensure the atomic values are 64-bit aligned"*); 15-7-stable uchun !5224 (2022-12-28…30).
`atomic.Int64` ga oʻtish taklifi ([gitaly#4702](https://gitlab.com/gitlab-org/gitaly/-/issues/4702)) 2026-11-25
holatiga hali ochiq.
Shunga oʻxshash: [allegro/bigcache#368](https://github.com/allegro/bigcache/issues/368) (32 bitli ARM, 2023).

**fieldalignment.** `golang.org/x/tools` dagi analizator: `Buyurtma has size 24 but the optimal size is 16 leading
to a waste of 8 bytes (33%)`. Uning hujjatida: *"the diagnostics produced by fieldanalyzer very rarely indicate a
significant problem, so the analyzer is not included in typical suites such as vet or gopls"*. `-fix` bayrogʻi
maydonlarni tartiblaydi, lekin bizning sinovda (`v0.53.0`) maydonlar yonidagi izohlarni (`// toʻlandimi?`)
oʻchirib yubordi: `-fix` dan keyin `git diff` ga qarang. Hujjatda yana: *"the most compact order is not always the
most efficient … 'false sharing'"*.

**Qachon tartiblash kerak.** Koʻpincha kerak emas: Dave Cheney ham ("Padding is hard", 2015) *"Do Go programmers
need to care about this minutiae ? Probably not."* deydi. Millionlab nusxasi xotirada turadigan struct'larda
(kesh yozuvlari, katta slice'lar) oʻlchang. Avval oʻqilishi oson tartib, keyin oʻlchov.

## Manbalar
- [Go spec: Size and alignment guarantees](https://go.dev/ref/spec#Size_and_alignment_guarantees),
  [Struct types](https://go.dev/ref/spec#Struct_types), [Composite literals](https://go.dev/ref/spec#Composite_literals)
- [pkg.go.dev/unsafe](https://pkg.go.dev/unsafe): `Sizeof` (*"For a struct, the size includes any padding introduced
  by field alignment"*), `Alignof`, `Offsetof`
- [Effective Go: Embedding](https://go.dev/doc/effective_go#embedding)
- [sync/atomic: Bugs](https://pkg.go.dev/sync/atomic#pkg-note-BUG), [Go 1.19 release notes: atomic types](https://go.dev/doc/go1.19#atomic_types)
- [fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment)
- [structs.HostLayout](https://pkg.go.dev/structs#HostLayout) (Go 1.23), [golang/go#10014](https://github.com/golang/go/issues/10014)
- [CMU 15-213: Machine-Level Programming IV: Data](https://www.cs.cmu.edu/afs/cs/academic/class/15213-s15/www/lectures/08-machine-data.pdf) (struct alignment)
- Dave Cheney, [Padding is hard](https://dave.cheney.net/2015/10/09/padding-is-hard) (2015)
- Go manba kodi: [cmd/compile/internal/types/size.go](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/types/size.go) (oxiridagi 0 baytli maydon, issue 9401)

## Uy vazifasi
- 🟢 **Oson:** `struct { A bool; B int64; C int32; D bool }` necha bayt? Avval qogʻozda, keyin `unsafe.Sizeof` bilan.
- 🟡 **Oʻrta:** `struct { Faol bool; Ism string; Admin bool; Yosh int32 }`: oʻlchami va eng ixcham tartib.
- 🟠 **Qiyin:** `struct{}` 0 bayt. Nega `struct { A int64; B struct{} }` 16 bayt, `struct { B struct{}; A int64 }` esa 8?
- 🔴 **Guru** (goroutine'lardan keyin)**:** ikki goroutine har biri oʻz `atomic.Int64` hisoblagichini oshiradi. Hisoblagichlar yonma-yon va
  64 bayt ajratilgan holatda benchmark yozing va solishtiring.

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): `A` 0, 7 bayt padding, `B` 8–15, `C` 16–19, `D` 20 → 21 → **24**.
  Ixcham tartib `B, C, A, D` = 14 → **16**.
- Oʻrta ([homework/orta](homework/orta/)): `string` = pointer + uzunlik, 16 bayt, align 8. `Faol` 0, `Ism` 8–23,
  `Admin` 24, `Yosh` 28–31 → **32**. Ixcham: `Ism, Yosh, Faol, Admin` = 22 → **24** (fieldalignment ham shuni aytadi).
- Qiyin ([homework/qiyin](homework/qiyin/)): oxirgi 0 baytli maydonning manzili (`&x.B`) struct'dan keyingi baytga,
  yaʼni boshqa obyektga ishora qilmasligi uchun kompilyator unga 1 bayt qoʻshadi, keyin 8 ga yaxlitlaydi: 9 → **16**
  (`cmd/compile/internal/types/size.go`, issue 9401). Boshida boʻlsa, qoʻshimcha bayt kerak emas: **8**.
- Guru ([homework/guru](homework/guru/), goroutine'lar va `sync/atomic` darslaridan keyin): `Yonma` (A, B
  yonma-yon) va `Ajratilgan` (orada `_ [56]byte`). `go test -bench . -count 5 ./homework/guru`: bizda `Ajratilgan`
  taxminan 5 barobar tez (mutlaq vaqtlar kompyuterga qarab koʻp oʻzgaradi).
  Ikki hisoblagich bitta 64 baytli kesh qatorida boʻlsa, har bir yozuv ikkinchi yadrodagi nusxani bekor qiladi
  (false sharing). Kesh qatori oʻlchami: `getconf LEVEL1_DCACHE_LINESIZE` → 64.

## Video uchun qayd etilgan vaqtlar

2026-11-25 kuni Intel i5-12450H, Go 1.26.8 linux/amd64 da, shu dars papkasida
ishlab chiqarish muhiti bilan o'lchandi:

`env -i HOME=$HOME PATH=$HOME/.local/bin:/run/current-system/sw/bin:/etc/profiles/per-user/$USER/bin bash -c 'go run ./tezlik'`

Har dastur ikkala variantni navbatma-navbat 100 marta o'lchab, eng yaxshi vaqtni chiqaradi.
Ikki ishlab chiqarish sinovidan ikkinchisi lesson.yaml dagi `out:` ga o'zgarishsiz yozildi:

| Buyruq | Buyurtma | Ixcham | Nisbat |
|---|---:|---:|---:|
| `go run ./tezlik` | 1426.5 µs | 1162.7 µs | 0.82 |
| `go run ./tezlik 100000` | 65.8 µs | 74.7 µs | 1.14 |

Birinchi ishlab chiqarish sinovida nisbatlar 0.73 va 1.14 edi. Barcha to'rtta
Bash sinovida ular 0.72–0.82 va 0.85–1.14 oralig'ida bo'ldi.
Bu qayd etilgan o'lchovlar, tezlik kafolati emas. Render paytida boshqa ishlar
ketishi mumkin; `tools/prepare.ts` da `out` bor buyruqlar qayta bajarilmaydi.

`fieldalignment@latest` o'rnatish buyrug'i offline ishlamaydi: GOPROXY=off
sinovida module lookup disabled xatosini berdi. Videodagi bo'sh `out` sabab
o'rnatish bajarilmaydi; analizatorning o'zi oldindan o'rnatilgan
`~/go/bin/fieldalignment` dan offline ishlaydi.
