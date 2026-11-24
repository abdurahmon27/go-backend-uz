# 07 · Taxmin qilmang, oʻlchang: birinchi test va benchmark

Oldingi dars: [06 · Funksiyalar](../006-functions/) · Keyingi dars: GO-1.08 · Slice'lar: `append` tuzogʻi va eskirgan 1024 qoidasi

| Papka | Nima |
|---|---|
| `sodda/` | kabisa yilining birinchi, sodda versiyasi (`yil%4 == 0`) + oddiy test + jadval test. **Ataylab yiqiladi** (1900). |
| `kabisa/` | toʻgʻri versiya (Grigorian qoidasi), oʻsha testlar; `bn_test.go`: eski `b.N` uslubida boʻsh sikl va kabisa; `loop_test.go`: xuddi shu ikkalasi `b.Loop` bilan |
| `doimiy/` | tajriba: `b.Loop` ichida doimiy argument — Go 1.26 javobni kompilyatsiya paytida hisoblaydi (pastda) |
| `tozala/` | [04-darsdagi](../004-strings-utf8/) `tozala` oʻzgarishsiz + jadval test |
| `tez/` | "tez yoʻl" qoʻshilgan, lekin faqat `'` ni tekshiradi. **Ataylab yiqiladi** (`tog’`). |
| `tezfix/` | toʻgʻri tez yoʻl (`strings.ContainsAny`) + `BenchmarkEski` / `BenchmarkTez`, har biri ikki matnda (`b.Run`: `toza`, `apostrof`) |
| `homework/` | uy vazifasi javoblari |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/007-test-bench`) ishga tushiriladi.
```bash
go test -bench KabisaBN ./kabisa | grep ns/op      # 0.1–0.3 ns/op, roppa-rosa 1e9 marta: shubhali
go test -v -run TestKabisa ./sodda                   # bitta test: PASS
go test ./sodda                                      # jadval: 100_ga yiqiladi (1900)
grep return kabisa/kabisa.go
go test -v -run Jadval/100 ./kabisa                  # faqat bitta subtest
go test ./kabisa                                     # butun paket: ok
go test ./tozala ./tez                               # eski: ok, tez: jingalak holat yiqiladi
grep ContainsAny tezfix/tozala.go
go test ./tezfix
go test -bench . -benchmem ./tezfix | grep ns/op | column -t
go test -bench BN ./kabisa | grep ns/op              # BoshBN ≈ KabisaBN: kabisa oʻchirib tashlangan
go test -c -o k.test ./kabisa && go tool objdump -s 'KabisaBN$' k.test | awk 'NR>1 {print $4, $5, $6}'   # faqat sanagich
go test -bench Loop ./kabisa | grep ns/op            # KabisaLoop ≈ 2 × BoshLoop: haqiqiy ish
```
Natijalar `go1.26.8 linux/amd64`, Intel i5-12450H (max 4,4 GHz) da olingan. **Nanosekundlar har safar oʻzgaradi**
(ayniqsa kompyuterda boshqa ish ketayotganda: bizda yuklama 13 boʻlganda ham, 1,5 boʻlganda ham oʻlchadik),
`B/op` va `allocs/op` esa oʻzgarmaydi. Shuning uchun videoda nisbatlar va ajratishlar soni haqida gapiramiz.
Muvaffaqiyatli testlar natijasi keshlanadi: kod oʻzgarmasa, `go test` vaqt oʻrniga `(cached)` deydi;
keshni chetlab oʻtish uchun `-count=1`.

## Dars matni

**Benchmark yolgʻon gapirdi.** `BenchmarkKabisaBN` kabisa funksiyasi uchun nanosekundning bir necha oʻndan birini
koʻrsatadi, sikl esa roppa-rosa 1 000 000 000 marta aylangan (`b.N` ning yuqori chegarasi, `testing/benchmark.go`:
`n < 1e9`). 4,4 GHz da bitta takt ≈ 0,23 ns. Bu hali isbot emas: zamonaviy protsessor juda qisqa siklni bir taktdan
tezroq aylantira oladi (boʻsh sikl ham 0,13–0,30 ns). Shuning uchun "shubhali" deymiz, isbot esa boʻsh sikl bilan
solishtirish. Go blogi aynan shunday holat haqida: *"CPUs are fast, but not that fast!"*

**Birinchi test.** Test fayli nomi `_test.go` bilan tugaydi: `go build` uni dasturga qoʻshmaydi, `go test` esa
topib ishga tushiradi. Test funksiyasi `TestXxx(t *testing.T)`. `t.Errorf` xatoni yozadi va testni yiqilgan deb
belgilaydi, lekin **toʻxtatmaydi** (`t.Fatalf` esa shu testni toʻxtatadi). `-v` har bir testni koʻrsatadi, `-run` nomi
regexp'ga mos testlarni tanlaydi (regexp langarsiz: `TestKabisa` `TestJadval` ga mos kelmaydi).

**Jadval test va subtest.** Har bir qator: nom, kirish, kutilgan javob. Yangi holat = bitta yangi qator. `t.Run`
har bir qatorni alohida subtest qiladi: xabarda qaysi holat yiqilgani koʻrinadi, `-run Jadval/100` bilan bittasini
ishga tushirish mumkin (`/` test nomini subtest nomidan ajratadi). Tuzatgandan keyin faqat yiqilgan holatni emas, butun
paketni ham ishga tushiring. Subtest nomidagi boʻsh joy `_` ga aylanadi
(`4 ga` → `4_ga`). Struct va slice keyingi darslarda (GO-1.08, GO-1.10) batafsil.

**1900.** Grigorian qoidasi: 4 ga boʻlinadi → kabisa, lekin 100 ga boʻlinsa, 400 ga ham boʻlinishi shart: 2000
kabisa, 1900 va 2100 emas. Microsoft hujjatiga koʻra, Lotus 1-2-3 1900-yilni kabisa deb hisoblagan, bu dasturga hisobni
osonlashtirgan ("This made it easier for the program to handle leap years"), Excel esa u bilan
moslik uchun buni saqlab qolgan: 1900-yil 29-fevral Excel'da hali ham "bor"; tuzatish mavjud hisob-kitoblarni
(masalan, hafta kunlarini) oʻzgartirib yuboradi.

**tozala va tez yoʻl.** 04-darsda "tozalashni bitta joyda qiling va uni test bilan himoyalang" degan edik. Jadval:
klaviatura `'`, jingalak `’`, tutuq belgisi va qoʻshtirnoq ichidagi gap. Xabarda `%+q` ishlatilgan: u ASCII'dan
tashqari belgilarni `\u2019` koʻrinishida chiqaradi, chunki koʻz `’` va `ʻ` ni ajratmaydi. Tez yoʻl: apostrof yoʻq
boʻlsa, `s` ning oʻzini qaytaramiz. (Halol eslatma: oʻzbekcha matnda apostrof deyarli har gapda bor, bu darsning
narratsiyasida 108 gapdan 75 tasida, yaʼni har uchtadan ikkitasida. Tez yoʻl raqamlar, kodlar, kirill yozuvidagi matn uchun foydali (inglizchada ham apostrof koʻp:
*don't*); apostrofli matnda foyda yoʻq, buni videoda oʻlchaymiz.) Birinchi urinish faqat `'` ni qidirdi va `tog’` holatida yiqildi; tuzatish:
`strings.ContainsAny(s, "'`´‘’ʻʼ")`. Test boʻlmasa, bu xato ekranda deyarli koʻrinmasdi.

**Benchmark.** `BenchmarkXxx(b *testing.B)`, oʻlchanadigan kod `for b.Loop() { … }` ichida. Go siklni kerakli marta
aylantiradi (odatda jami ≈1 soniya, `-benchtime`, lekin `b.N` ko'pi bilan 1e9). `-bench .` hamma benchmark'larni, `-benchmem` xotirani ham
oʻlchaydi. Ustunlar: iteratsiyalar soni, `ns/op`, `B/op`, `allocs/op`. Bizning (48 belgili) matnda eski versiya har
chaqiruvda 2 marta xotira ajratadi (`[]rune(s)` va `string(r)`). Qisqa matnda (≤32 rune) `[]rune(s)` stekda qoladi
(`go build -gcflags=-m`: `([]rune)(s) does not escape`), shunda 1 marta. Apostrofsiz matnda tez versiya 0 marta
va bizda ≈2,5–3 barobar tez; apostrofli matnda (`Oʻzbekiston togʻlari, maʼno va qoʻshiq`) ikkalasi ham 208 B/op,
2 allocs/op: foyda yoʻq, vaqt farqi shovqin ichida (`ContainsAny` birinchi `ʻ` da toʻxtaydi).

**Nega eski benchmark yolgʻon gapirdi.** `for i := 0; i < b.N; i++ { kabisa(i) }` da natija ishlatilmaydi.
Kompilyator `kabisa` ni inline qiladi va hisobni butunlay oʻchiradi: assembly'da sikl ichida faqat sanagich qoladi,
`BenchmarkBoshBN` bilan (manzillardan tashqari) bir xil:
```bash
go test -c -o k.test ./kabisa
go tool objdump -s 'BenchmarkBoshBN$'   k.test | awk '{print $4, $5, $6}'   # XORL CX, CX / JMP / INCQ CX / CMPQ 0x210(AX), CX / JG / RET
go tool objdump -s 'BenchmarkKabisaBN$' k.test | awk '{print $4, $5, $6}'   # xuddi shunday
```
Shuning uchun `KabisaBN` ≈ `BoshBN`. `b.Loop` ning oʻzi ham bir oz vaqt oladi (hisoblagich xotirada), shuning uchun
`KabisaLoop` ni `BoshLoop` bilan solishtiramiz; ikkalasida ham `yil++` bor (u ham "tirik" saqlanadi). Bizda
≈1,7–2,1 barobar farq (≈0,5 va ≈0,9–1,0 ns/op): kabisaning haqiqiy ishi ≈0,4–0,5 ns, qolgani oʻlchov narxi. Go 1.24'da qoʻshilgan `b.Loop` sikl ichidagi chaqiruvlarning
argument va natijalarini "tirik" saqlaydi; Go 1.26 da bu `runtime.KeepAlive` bilan qilinadi va inlining endi
taqiqlanmaydi. `b.Loop` setup kodini ham bir marta bajaradi (b.N uslubida benchmark funksiyasi bir necha marta
chaqiriladi) va taymerni oʻzi boshqaradi.

**`b.Loop` ham sehr emas (`doimiy/`).** `for b.Loop() { kabisa(2024) }`: Go 1.26.8 `kabisa` ni inline qilib,
`kabisa(2024) = true` ni kompilyatsiya paytida hisoblaydi; siklda faqat `true` ni saqlash qoladi:
```bash
go test -c -o d.test ./doimiy
go tool objdump -s 'BenchmarkDoimiyLoop$' d.test | grep -c IMULQ    # 0: %4/%100/%400 hisobi yoʻq
go tool objdump -s 'BenchmarkKabisaLoop$' d.test | grep -c IMULQ    # 1: hisob bor
```
(Go `% 100` ni koʻpaytirish bilan hisoblaydi, shuning uchun `IMULQ` ni sanaymiz.) Shuning uchun videodagi
`BenchmarkKabisaLoop` da yil har aylanishda oʻzgaradi. (Go 1.24'da `b.Loop` sikl tanasiga
inlining'ni taqiqlardi, u yerda bu holat boshqacha boʻlgan; biz faqat 1.26.8 ni tekshirdik.)

**Shovqin.** Bitta oʻlchovga ishonmang: ikki versiyani solishtirganda har birini kamida 10 marta ishga tushiring va
[benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) bilan solishtiring:
```bash
go install golang.org/x/perf/cmd/benchstat@latest
go test -run '^$' -bench . -benchmem -count 10 ./tezfix > eski.txt   # oʻzgarishdan oldin
# ... kodni oʻzgartiring ...
go test -run '^$' -bench . -benchmem -count 10 ./tezfix > yangi.txt  # keyin
benchstat eski.txt yangi.txt
```
benchstat median va ishonch oraligʻini beradi; `~` "statistik farq yoʻq" degani. "Natija yoqguncha qayta-qayta
ishga tushirish" ham statistik xato (multiple testing).

## Manbalar
- [pkg.go.dev/testing](https://pkg.go.dev/testing): `TestXxx`, `BenchmarkXxx`, `B.Loop`, subtestlar;
  [`go help testflag`](https://pkg.go.dev/cmd/go/internal/test): `-run`, `-bench`, `-benchmem`, `-count`, `-benchtime`
- [go.dev/blog/testing-b-loop](https://go.dev/blog/testing-b-loop) (J. Shao, 2025);
  [Go 1.24 release notes](https://go.dev/doc/go1.24#testing); [Go 1.26 release notes](https://go.dev/doc/go1.26);
  taklif [golang/go#61515](https://github.com/golang/go/issues/61515) (A. Clements); [#73137](https://github.com/golang/go/issues/73137)
- [Go source: testing/benchmark.go `B.Loop`](https://go.dev/src/testing/benchmark.go),
  [testing/match.go](https://go.dev/src/testing/match.go) (subtest nomi)
- [go.dev/blog/subtests](https://go.dev/blog/subtests) (M. van Lohuizen, 2016); [Go wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests)
- Dave Cheney: [How to write benchmarks in Go](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go) (2013),
  [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html) (b.Loop'dan oldingi davr)
- [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- Microsoft: [Excel incorrectly assumes that the year 1900 is a leap year](https://learn.microsoft.com/en-us/troubleshoot/microsoft-365-apps/excel/wrongly-assumes-1900-is-leap-year)

## Uy vazifasi
- 🟢 **Oson:** `kabisa` jadvaliga 2100 va 2400 yillarini qoʻshing.
- 🟡 **Oʻrta:** tez yoʻlni ham oʻlchang: `ContainsAny` oʻrniga baytlarni oddiy sikl bilan tekshiring va solishtiring.
- 🟠 **Qiyin:** eski uslubdagi benchmark'ni `b.Loop`'siz tuzating: natijani paket darajasidagi oʻzgaruvchiga yozing.
- 🔴 **Guru:** eski `tozala` ni shunday oʻzgartiringki, matnda hech narsa oʻzgarmasa, yangi string yaratmasin;
  buni `-benchmem` bilan isbotlang.

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da:
- Oson ([homework/oson](homework/oson/)): 2100 → `false` (100 ga boʻlinadi, 400 ga yoʻq), 2400 → `true`. `go test ./homework/oson` → `ok`.
- Oʻrta ([homework/bayt](homework/bayt/)): `strings.ContainsAny` ASCII boʻlmagan belgilar toʻplamida matnni
  rune'ma-rune tekshiradi. `'` va `` ` `` dan boshqa hamma apostroflar UTF-8'da `0x80` dan katta baytdan boshlanadi,
  shuning uchun baytlar boʻyicha oddiy sikl (`c == '\'' || c == '`' || c >= 0x80` → sekin yoʻl) ham toʻgʻri, ham bizda
  ≈4,5 barobar tez: ≈32 va ≈145 ns/op, ikkalasi 0 allocs; testlar oʻtadi. Narxi: kirill matni ham (baytlari `≥ 0x80`)
  endi sekin yoʻlga tushadi. Har bir tezlashtirishning oʻz chegarasi bor, shuning uchun oʻlchaymiz.
- Qiyin ([homework/qiyin](homework/qiyin/)): `var natija bool` va `natija = kabisa(i)`. Natija paket
  oʻzgaruvchisiga yozilgani uchun kompilyator hisobni oʻchira olmaydi: `BenchmarkSinkBN` boʻsh sikldan sezilarli sekin
  (bizda ≈0,52 ns/op, `KabisaBN` ≈0,13). Bu Dave Cheney 2013-yilda tavsiya qilgan usul; bugun `b.Loop` afzal.
- Guru ([homework/guru](homework/guru/)): `ozgardi` bayrogʻi; hech bir belgi oʻzgarmasa `return s`. Apostroflari
  allaqachon toʻgʻri matnda 2 → **1** allocs/op (faqat `[]rune`), 208 → 160 B/op. Testlar oʻtadi
  (`go test ./homework/guru`).
