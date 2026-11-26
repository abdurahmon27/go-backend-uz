# DSA-A.01 · RAM modeli va haqiqiy protsessor: kesh qanchalik muhim

DSA trekining birinchi darsi. Oldingi Go darslari: [07 · test va benchmark](../007-test-bench/),
[08 · slice'lar](../008-slices/), [11 · pointer'lar](../011-pointers-escape/) · Keyingi dars: DSA-A.02 · Big-O'ni
oʻlchab topamiz: doubling usuli.

Bitta O(n) ish — `n` ta `int64` ni qoʻshish — toʻrt xil xotira tartibida, `n` = 2¹⁰ … 2²⁶ (8 KiB … 512 MiB):

| | mustaqil oʻqishlar | zanjir (keyingi manzil oldingi oʻqishdan) |
|---|---|---|
| **ketma-ket tartib** | `Ketma` | `TZanjir` (`s[i] = i+1`) |
| **aralash tartib** | `Aralash` (aralashtirilgan indekslar boʻyicha) | `Zanjir` (indeksli bogʻlangan roʻyxat, bitta tasodifiy aylana) |

| Papka / fayl | Nima |
|---|---|
| `kesh/kesh.go` | uchta sikl: `Ketma`, `Aralash`, `Zanjir` (`TZanjir` = `Zanjir` + `TartibliZanjir` slice'i) |
| `kesh/royxat.go` | haqiqiy pointer'li roʻyxat: `RoyxatYangi` (tugunlar birin-ketin) va `RoyxatSochilgan` (tugunlar xotirada aralash) |
| `kesh/tayyor.go` | slice, aralash tartib, ikki xil zanjirni tayyorlash; `Qadam` (har k-element) |
| `kesh/*_test.go` | toʻgʻrilik testi (beshala usulda yigʻindi n(n−1)/2), `TestYonmaYon`, benchmark'lar (`olcha_test.go`: 17 oʻlcham, `ns/el` metrikasi) |
| `olchov.sh [yadro]` | butun oʻlchov: 10 marta `-count 1`, navbat bilan → `natija.txt` (bu noutbukda 17 daqiqa) |
| `jadval/` | `go test` natijasini ekranga sigʻadigan grafik/jadvalga aylantiradi; kesh oʻlchamlari `/sys/devices/system/cpu/cpuN/cache` dan (yoki `-l1 48K -l2 1280K -l3 12M`) |
| `natija.txt` | videodagi xom natija (3-oʻlchov, 10 namuna) · `natija-1.txt`, `natija-2.txt` — oldingi ikki oʻlchov (5 namuna, `TZanjir` siz) |
| `qadam-royxat.txt` | pauza savoli va bogʻlangan roʻyxat: `-count 10` |
| `uyvazifa/` | uy vazifasi javoblari · `uyvazifa.txt` xom natija |

## Videodagi buyruqlar
Hammasi shu papkadan ishga tushiriladi.
```bash
go run ./jadval -n 16777216 < natija.txt     # hook: 2^24 element, uch usul, vaqt va nisbat
go run ./jadval -kesh                        # shu kompyuterning keshlari (P va E yadro)
./olchov.sh                                  # 17 daqiqa: natija.txt + tarqoqlik
go run ./jadval < natija.txt                 # terminal grafik (sparkline): ns/element, 2^10 … 2^26
go run ./jadval -xy Zanjir < natija.txt      # videodagi chart: "n median" qatorlari (Ketma, Aralash, Zanjir)
go run ./jadval -daraja < natija.txt         # har kesh darajasi ichida: eng kami … eng koʻpi
go run ./jadval -daraja < natija.txt | cut -c1-16,39-   # faqat zanjir ustuni
go run ./jadval -ikki < natija.txt           # 2×2: tartib × bogʻliqlik
taskset -c 2 go test -run '^$' -bench 'Qadam|Royxat' -count 10 ./kesh > qadam-royxat.txt
go run ./jadval -median < qadam-royxat.txt   # pauza javobi va bogʻlangan roʻyxat
go test -run YonmaYon -v ./kesh | grep bayt  # birin-ketin yaratilgan tugunlar yonma-yonmi?
```
Natijalar: `go1.26.8 linux/amd64`, Intel Core i5-12450H (4 unumdor + 4 tejamkor yadro), Linux 6.18, 2026-11-26.
Hamma oʻlchov `taskset -c 2` bilan bitta unumdor (P) yadroda: P va E yadrolarning keshi har xil
(`jadval -kesh`: P — L1d 48K, L2 1280K; E — L1d 32K, L2 2M; L3 12M umumiy; qator 64 bayt).

### Natijalar (`natija.txt`, 10 namuna mediani, ns/element)

```
qayer slice           ketma    aralash     zanjir
L1    8K…32K           0.31  0.41…0.42    1.5…1.7
L2    64K…1M      0.32…0.33  0.44…0.83    2.5…7.1
L3    2M…8M       0.32…0.49    1.2…3.8      17…81
DRAM  16M…512M    0.58…0.59   5.9…10.0     94…137

n = 2^26 (512M), ns/element    mustaqil   zanjir
ketma-ket tartib                0.58      1.7
aralash tartib                  10.0      137
```
`-daraja` har kesh uchun slice (8·n bayt) shu keshga sigʻadigan, lekin oldingisiga sigʻmaydigan **hamma** oʻlchamlarni
oladi va medianlarning eng kami … eng koʻpini koʻrsatadi. Zinapoyalar `|` dan biroz oldin boshlanadi: kesh hech
qachon bitta slice'ga toʻliq tegmaydi (boshqa maʼlumot, ikkinchi hyperthread, assotsiativlik). L3 ichida zanjir 17 dan
81 ns gacha oʻsadi: L3 toʻla boshlaydi, uni boshqa dasturlar ham ishlatadi, TLB xatolari (manzil tarjimasi; Drepper
§3.3.2, §4) qoʻshiladi.

**2×2 nimani koʻrsatadi.** Ikki mexanizm bor. *Prefetcher* tartibni payqab, qatorlarni oldindan olib keladi: ketma-ket
zanjir bogʻliq boʻlsa ham 1,7 ns. *Parallel kutish* (out-of-order + bir vaqtda bir nechta kesh xatosi, "nonblocking
cache") mustaqil oʻqishlarni ustma-ust kutadi: aralash mustaqil 10 ns, aralash zanjir 137 ns. 137 / 10 ≈ 13–14:
taxminan shuncha oʻqish bir vaqtda "yoʻlda" (taxminiy hisob).

**Ketma nega kesh ichida tekis?** Sikl skalyar: `go tool objdump` da `ADDQ 0(CX)(AX*8), BX` — har qoʻshish oldingi
yigʻindini kutadi (taktiga bitta), kesh esa bundan tezroq beradi. Faqat DRAM da xotira oʻzi chegaraga aylanadi.

### Oʻlchovlar va shovqin
Oʻlchov paytida kompyuterda boshqa ishlar ham bor edi (ovoz sintezi, video render navbati, boshqa benchmark'lar;
yuklama 3–13). Navbat toʻxtatilmadi. `taskset` faqat bizning jarayonni cpu2 ga bogʻlaydi, cpu2 va uning
hyperthread qoʻshnisi cpu3 ni boshqalardan band qilmaydi, `b.Loop` esa devor soati bilan oʻlchaydi. Shuning uchun:
10 namuna, navbat bilan, median; tarqoqlik ikki xil oʻlchandi (`jadval -tarqoq`):

| | IQR/median (median · eng koʻpi) | (max−min)/median (median · eng koʻpi) |
|---|---|---|
| ketma | 6,1 % · 11,7 % | 31,0 % · 66,4 % |
| aralash | 8,5 % · 40,2 % | 43,7 % · 77,3 % |
| zanjir | 6,4 % · 39,7 % | 35,5 % · 116,7 % |
| t.zanjir | 6,7 % · 10,2 % | 32,8 % · 39,6 % |

Uch oʻlchovda zinapoyalar bir joyda (`natija-1.txt`, `natija-2.txt` bilan `jadval -daraja`): zanjir L1 1,3–1,7, DRAM
(512 MiB) 124 / 164 / 137 ns; ketma DRAM 0,51–0,73. 2²⁴ da (hook) nisbat: 12× / 16× / 15× va 150× / 195× / 179×.

### Oʻlchov tuzoqlari
- **Hech yozilmagan xotira.** `make([]int64, n)` dan keyin elementlarga yozmasangiz, Linux bu sahifalarni bitta umumiy
  nol sahifaga ulaydi va oʻqish keshdan keladi. Shuning uchun `Sonlar`/`toldir` har elementga yozadi.
- **`b.Loop` tanasiga sikl yozmang.** `go doc testing.B.Loop`: tanadagi oʻzgaruvchilar "tirik" ushlanadi. Siklni
  toʻgʻridan-toʻgʻri `for b.Loop() { … }` ichiga yozganimizda `sum` stack'da yashadi (`ADDQ BX, 0x18(SP)`) va sikl
  bir necha barobar sekinlashdi. Oʻlchanadigan kod funksiyada turadi.
- **Tayyorlash axlati.** 512 MiB slice'lar tayyorlangandan keyin `runtime.GC()` — aks holda GC sikl paytida ishlashi
  mumkin (`taskset` ostida `GOMAXPROCS=1`).
- **Hybrid protsessor.** `taskset` siz Linux benchmark'ni E-yadroga oʻtkazishi mumkin — u yerda kesh boshqacha.
- `Aralash` qoʻshimcha `[]int32` tartib slice'ini (4·n bayt) ketma-ket oʻqiydi; u prefetcher bilan keladi, lekin
  keshdan joy oladi. `|` chiziqlar faqat asosiy slice (8·n) boʻyicha.
- **Go GC obyektlarni koʻchirmaydi** (`runtime/mgc.go`: "non-generational and non-compacting"): sochilgan tugunlar
  sochilganicha qoladi. Yangi (toza) heap'da esa birin-ketin yaratilgan 16 baytli tugunlarning 99,8 % i roppa-rosa
  16 bayt narida (`TestYonmaYon`).

## Dars matni

**hook.** Bir xil sonlar, bir xil qoʻshish, bir xil Big-O, lekin bir sikl ikkinchisidan yuz barobardan koʻproq sekin ishlaydi. Oʻn olti million sonni uch xil tartibda qoʻshdim, har bir son roppa-rosa bir marta oʻqiladi. Sizningcha, qaysi biri eng sekin? Birinchisi slice boʻylab ketma-ket yuradi va qariyb oʻn millisekundda tugaydi. Ikkinchisi xuddi shu sonlarni aralash tartibda oʻqiydi va oʻn besh barobar sekin. Uchinchisi ham aralash sakraydi, lekin bogʻlangan roʻyxatdek keyingi manzilni oldingi sondan oladi. U qariyb ikki soniya ishlaydi. Farqni xotira yaratadi: kesh va kutish.

**trek.** Bu DSA trekining birinchi darsi: algoritmlar va ma’lumotlar tuzilmalari, hammasi Goʻda. Bu trekda qoidamiz oddiy: Big-O bahosini taxmin qilmaymiz, benchmark bilan oʻlchab tekshiramiz. Avval esa Big-O qaysi modelga tayanishini va haqiqiy protsessor undan qayerda farq qilishini koʻramiz.

**ram.** Algoritmlar boʻyicha klassik darslik, Kormen va hammualliflarining kitobi, algoritmlarni RAM modelida tahlil qiladi. Bu yerda RAM operativ xotira emas, random-access machine, ya’ni ixtiyoriy murojaatli mashina degani. Bu modelda har bir oddiy amal oʻzgarmas vaqt oladi: qoʻshish ham, xotiradan oʻqish ham. Kitob ochiq yozadi: model kesh va virtual xotirani hisobga olmaydi, lekin odatda haqiqiy tezlikni yaxshi bashorat qiladi. Bugun u koʻrmaydigan narsani oʻlchaymiz.

**qavatlar.** Haqiqiy protsessorda xotira bir necha qavatdan iborat. Har bir yadroning ichida kichik va juda tez L1 kesh bor. Undan keyin kattaroq va sekinroq L2, keyin hamma yadrolar boʻlishadigan L3. Eng pastda esa operativ xotira. Buni choʻntak, sumka, shkaf va omborga oʻxshatish mumkin: yaqini kichik, lekin darhol qoʻl yetadi. Ma’lumot keshga bayt-bayt emas, oltmish toʻrt baytlik boʻlaklarda keladi. Bu boʻlak kesh qatori, inglizcha cache line deyiladi.

**keshlar.** Bu noutbukdagi qavatlarni Linux’ning oʻzidan oʻqiymiz. Protsessorda unumdor va tejamkor yadrolar bor. Unumdor yadroda L1 qirq sakkiz kilobayt, L2 ming ikki yuz sakson kilobayt. L3 oʻn ikki megabayt va hamma yadrolarga umumiy. Tejamkor yadroning keshi boshqacha, shuning uchun oʻlchovni taskset bilan bitta unumdor yadroga bogʻladik.

**kod.** Endi tajribaning kodiga qaraymiz. Sakkizinchi darsda koʻrgan edik: slice elementlari bitta massivda yonma-yon turadi. Uchta funksiya ham shu slice’dagi n ta sonni qoʻshadi. Ketma elementlarni xotiradagi tartibda oʻqiydi. Aralash ham har bir elementni bir marta oʻqiydi, lekin indeksni oldindan aralashtirilgan tartib slice’idan oladi. Zanjirda har bir element keyingi elementning indeksini saqlaydi: bu pointer oʻrniga indeksli bogʻlangan roʻyxat. Zanjir ham aralash sakraydi, lekin keyingi manzil oldingi oʻqish tugamaguncha noma’lum.

**olchov.** Oʻlchov yettinchi Go darsidagi b.Loop bilan yozilgan. Har bir oʻlcham alohida ichki benchmark: sakkiz kilobaytdan besh yuz oʻn ikki megabaytgacha. Slice’ni tayyorlash vaqti hisobga kirmaydi, chunki b.Loop birinchi chaqiruvda taymerni nolga qaytaradi. Tayyorlashdan qolgan axlatni esa garbage collector oldindan yigʻib qoʻyadi. Oxirida ReportMetric bitta elementga ketgan nanosekundni chiqaradi.

**grafik.** Oʻlchov skripti hammasini oʻn marta navbat bilan takrorlaydi, bunga oʻn yetti daqiqa ketdi. Kompyuter bu paytda boshqa ishlar bilan ham band edi, shuning uchun oʻnta natijaning medianasini olamiz. Gorizontal oʻqda elementlar soni, vertikal oʻqda bitta elementga ketgan vaqt, ikkalasi ham logarifmik shkalada. Ketma-ket oʻqish deyarli tekis qoladi. Aralash tartib ikkinchi kesh toʻlib borgani sari koʻtariladi. Zanjir esa har bir chegaradan keyin sakraydi: L1, L2, L3. Zinapoya belgidan biroz oldin boshlanadi, chunki kesh bitta slice’ga toʻliq tegmaydi.

**ketma.** Avval ketma-ket oʻqishni koʻraylik. L1 ichida bitta element uchdan bir nanosekundga yaqin vaqt oladi, operativ xotirada esa yarim nanosekunddan sal koʻp. Slice oltmish besh ming marta kattalashdi, vaqt esa ikki barobarga ham yetmay oshdi. Kesh ichida chegarani siklning oʻzi qoʻyadi: har bir qoʻshish oldingi yigʻindini kutadi. Operativ xotirada esa ikki narsa yordam beradi. Birinchisi kesh qatori: bitta qator sakkizta int64 olib keladi. Ikkinchisi prefetcher degan mexanizm: protsessor ketma-ket yurishni payqaydi va keyingi qatorlarni oldindan yuklab qoʻyadi.

**zanjir.** Zanjirda manzara butunlay boshqa. L1 ichida bitta qadam bir yarim nanosekund. L2 ichida ikki yarimdan yettigacha, L3 ichida oʻn yettidan saksongacha oʻsib boradi. Operativ xotirada esa toʻqsondan yuz oʻttiz yettigacha oshadi. Har bir oʻqish oldingisini toʻliq kutadi, shuning uchun biz taxminan har bir qavatning kechikishini koʻramiz. Katta oʻlchamda bunga manzil tarjimasi, ya’ni TLB xatolari ham qoʻshiladi. Umumiy L3’ni esa bu paytda boshqa dasturlar ham ishlatayotgan edi.

**ikki.** Endi asosiy savol: zanjirni aynan nima sekinlashtiradi? Buni ajratish uchun toʻrtinchi variantni ham oʻlchadik: zanjir, lekin keyingi element doim yonidagisi. Ketma-ket tartibda zanjir ikki nanosekunddan ham kam, aralash tartibda esa yuz oʻttiz yetti. Demak bu yerda ikki xil mexanizm bor. Prefetcher’ga tartib kerak, parallel kutish uchun esa oʻqishlar mustaqil boʻlishi kerak. Shuning uchun aralash tartibda zanjir mustaqil oʻqishlardan oʻn uch barobar sekinroq.

**savol.** Endi sizga bitta savol beraman. Yuz yigirma sakkiz megabaytlik slice’da har bir elementni emas, faqat har sakkizinchi elementni qoʻshamiz. Ish sakkiz barobar kamaydi. Vaqt ham sakkiz barobar kamayadimi? Videoni pauza qilib, kesh qatorini eslab oʻylab koʻring.

**javob.** Javob: yoʻq, vaqt sakkiz barobar kamaymaydi. Kichik, L1’ga sigʻadigan slice’da qadam sakkiz haqiqatan qariyb sakkiz barobar tez. Katta slice’da esa bor-yoʻgʻi bir yarim barobar tez. Sabab yana kesh qatori. Har sakkizinchi element yangi qatorda turadi, demak protsessor baribir har bir qatorni operativ xotiradan olib keladi. Bir yarim barobar qoʻshishlar kamligidan keladi, asosiy narx esa oʻsha qatorlar.

**royxat.** Endi haqiqiy, pointer’li bogʻlangan roʻyxat. Oʻn birinchi Go darsida koʻrgan edik: funksiyadan tashqarida yashaydigan qiymat heap’ga tushadi. Bir million tugunni birin-ketin yaratsak, toza heap’da Goʻning allocator’i ularni odatda yonma-yon qoʻyadi. Test buni tekshiradi: deyarli har bir tugundan keyingisi roppa-rosa oʻn olti bayt narida. Shunda bitta tugun ikki nanosekundga yaqin, sochilgan tugunlar toʻqson nanosekund, slice esa yarim nanosekunddan kam. Goʻning garbage collector’i obyektlarni koʻchirmaydi: bir marta sochilgan tugunlar sochilganicha qoladi.

**xulosa.** Endi qisqa xulosa qilamiz. Big-O amallar soni qanday oʻsishini aytadi, lekin bitta amalning narxini aytmaydi. Bir xil chiziqli sikl xotira tartibiga qarab bitta elementga uchdan bir nanosekunddan yuz nanosekunddan koʻproqqacha vaqt oldi. Keshni hisobga oladigan modellar ham bor, ularda narx kesh xatolari soni bilan oʻlchanadi. Goʻda amaliy ma’nosi: qiymatlarni slice’da yonma-yon saqlang va benchmark’ni haqiqiy oʻlchamda oʻtkazing, chunki kichik test L1’da yashaydi.

**homework.** Uy vazifasi, toʻrt darajada. Oson: oʻlchov skriptini oʻz kompyuteringizda ishga tushiring va kesh chegaralarini toping. Oʻrta: qadamni birdan oʻttiz ikkigacha oʻzgartiring: vaqt qanday oʻzgaradi va nega? Qiyin: int64 oʻrniga int32 slice’ni ketma-ket qoʻshing, operativ xotirada necha barobar farq chiqadi? Guru: toʻrt maydonli struct’lar slice’idan faqat bitta maydonni yigʻing va alohida slice bilan solishtiring. Javoblar GitHub’da, link pastda.

**outro.** Bugun koʻrdik: bir xil Big-O ostida yuz barobardan katta farq yashirinishi mumkin, uni faqat oʻlchab topamiz. Keyingi darsda oʻlchamni ikki barobar oshirib, vaqt nisbatidan Big-O darajasini topadigan asbob yozamiz. Bugungi kesh chegaralari u yerda ham chiqadi. Darsni koʻrganingiz uchun rahmat, koʻrishguncha.

## Manbalar
- T. Cormen, C. Leiserson, R. Rivest, C. Stein. *Introduction to Algorithms*, 3-nashr, §2.2 "Analyzing algorithms":
  RAM modeli — har bir koʻrsatma (arifmetika, maʼlumot koʻchirish, boshqaruv) oʻzgarmas vaqt oladi; xotira
  ierarxiyasi (kesh, virtual xotira) modellashtirilmaydi; shunga qaramay RAM tahlili odatda haqiqiy tezlikni yaxshi
  bashorat qiladi.
- Keshni hisobga oladigan modellar: A. Aggarwal, J. Vitter, "The input/output complexity of sorting and related
  problems", CACM 1988 (I/O modeli); M. Frigo, C. Leiserson, H. Prokop, S. Ramachandran, "Cache-oblivious
  algorithms", FOCS 1999 (ideal kesh modeli). Bu modelda `Ketma` ≈ n/8 kesh xatosi (64 / 8), `Zanjir` ≈ n.
- R. Bryant, D. O'Hallaron. *Computer Systems: A Programmer's Perspective*, 3-nashr: §5.7–5.8 (out-of-order, latency
  bound), 6-bob (xotira ierarxiyasi, lokallik, kesh qatorlari, §6.6.1 "memory mountain"), §9.6 (TLB); CMU 15-213:
  [The Memory Hierarchy](https://www.cs.cmu.edu/~213/lectures/09-memory-hierarchy.pdf),
  [Cache Memories](https://www.cs.cmu.edu/~213/lectures/10-cache-memories.pdf).
- U. Drepper, *What Every Programmer Should Know About Memory* (2007): [LWN 2-qism, CPU keshlari, §3.3.2
  "Measurements of Cache Effects"](https://lwn.net/Articles/252125/) — `struct l { struct l *n; long int pad[NPAD]; }`
  roʻyxati boʻylab ketma-ket va tasodifiy yurish, prefetch; [5-qism](https://lwn.net/Articles/255364/);
  [PDF](https://www.akkadia.org/drepper/cpumemory.pdf).
- J. Hennessy, D. Patterson, *Computer Architecture: A Quantitative Approach*, 6-nashr, §2.3 (nonblocking caches:
  hit under miss, miss under miss).
- I. Ostrovsky, [Gallery of Processor Cache Effects](https://igoro.com/archive/gallery-of-processor-cache-effects/)
  (2010), 1-misol: har elementni va har 16-elementni koʻpaytirish deyarli bir xil vaqt (80 va 78 ms).
- Golden Cove kesh kechikishlari (L1D 5 takt, L2 ≈ 15 takt): [Chips and Cheese](https://old.chipsandcheese.com/2022/02/11/going-armchair-quarterback-on-golden-coves-caches/).
- Go: [`testing.B.Loop`, `B.ReportMetric`](https://pkg.go.dev/testing#B.Loop);
  [`runtime/malloc.go`](https://github.com/golang/go/blob/go1.26.8/src/runtime/malloc.go#L27-L35) va
  [`nextFreeFast`](https://github.com/golang/go/blob/go1.26.8/src/runtime/malloc.go#L969-L981): manzil =
  `base + index·elemsize`, yangi span'dan ketma-ket olingan tugunlar yonma-yon;
  [`runtime/mgc.go`](https://github.com/golang/go/blob/go1.26.8/src/runtime/mgc.go#L7-L9): GC "non-compacting".
- Linux: [sysfs `cpu/cpuN/cache/indexM`](https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-devices-system-cpu).
- Qoʻshimcha oʻqish (videoga kirmadi): Drepper sahifa ichida aralashtirish bilan TLB narxini DRAM kechikishidan
  ajratadi; Mike Acton, "Data-Oriented Design and C++" (CppCon 2014) — guru vazifasining nomi.

## Uy vazifasi
1. **Oson.** `./olchov.sh N` ni oʻz kompyuteringizda ishga tushiring (`N` — `lscpu -e` dagi yadro raqami) va
   `go run ./jadval -cpu N < natija.txt` bilan grafik chizing. Zinapoyalar qaysi oʻlchamda?
2. **Oʻrta.** `uyvazifa`: qadam k = 1, 2, 4, 8, 16, 32 (2²⁴ element). Vaqt qanday oʻzgaradi va nega?
3. **Qiyin.** `[]int32` va `[]int64` ni ketma-ket yigʻing (2²⁴ element, DRAM). Necha barobar farq?
4. **Guru.** `[]Nuqta{X, Y, Z, W}` (32 bayt) dan faqat `X` ni yigʻish (AoS) va alohida `[]int64` (SoA) — necha
   barobar farq? (Bu "data-oriented design" gʻoyasi.)

### Javoblar
Oʻlchov: `taskset -c 2 go test -run '^$' -bench . -count 5 ./uyvazifa` ikki marta (10 namuna, yuklama 3–4,5),
`go run ./jadval -median < uyvazifa.txt`:
```
Qadam/k=1         8.9 ms  [8.6 ms … 9.2 ms]
Qadam/k=2         7.1 ms  [6.9 ms … 7.5 ms]
Qadam/k=4         6.5 ms  [6.2 ms … 7.7 ms]
Qadam/k=8         6.5 ms  [5.9 ms … 7.6 ms]
Qadam/k=16        5.5 ms  [5.2 ms … 6.6 ms]
Qadam/k=32        3.2 ms  [3.2 ms … 3.3 ms]
Int64             8.0 ms  [7.8 ms … 8.2 ms]
Int32             5.9 ms  [5.6 ms … 6.2 ms]
AoS                24 ms  [24 ms … 26 ms]
SoA               8.2 ms  [7.9 ms … 9.9 ms]
```
1. **Oson.** Bu noutbukda: `natija.txt` va yuqoridagi grafik. Sizda zinapoyalar L1d/L2/L3 hajmlaringizga mos joyda
   boʻladi (8·n ≈ kesh hajmi), lekin biroz oldinroq boshlanadi.
2. **Oʻrta.** k = 1 → 8 da qoʻshishlar 8 barobar kamaydi, vaqt esa atigi 1,4 barobar (8,9 → 6,5 ms), k = 4 va 8 deyarli
   bir xil: har qatorga baribir tegamiz. k = 16 da har ikkinchi qator tashlanadi, lekin vaqt faqat ~15 % kamaydi:
   L2 ning "spatial prefetcher"i har olingan qatorni juftiga — 128 baytlik bloklarga — toʻldiradi, tashlangan qator
   ham baribir keladi (Intel Optimization Reference Manual, L2 spatial prefetcher). k = 32 da toʻrt qatordan uchtasi
   tashlanadi (128 baytlik juftlarning yarmi) → vaqt k = 8 dan 2 barobar kam (3,2 ms). Qoida: vaqt ≈ xotiradan
   olingan baytlar, qoʻshishlar soni emas.
3. **Qiyin.** `[]int32` (64 MiB) 5,9 ms, `[]int64` (128 MiB) 8,0 ms: baytlar 2 barobar kam, vaqt esa 1,36 barobar.
   Bitta yadro xotira kanalini toʻliq band qila olmaydi: sikl ham (har qoʻshish oldingisini kutadi) vaqt oladi.
4. **Guru.** AoS (`[]Nuqta`, 32 bayt, kerakli `X` faqat 8 bayt) 24 ms, SoA (`[]int64`) 8,2 ms: 2,9 barobar. Bitta
   64 baytlik qatorda 2 ta `Nuqta` yoki 8 ta `X`: AoS 4 barobar koʻp bayt olib keladi. "Data-oriented design":
   bir vaqtda ishlatiladigan maydonlarni yonma-yon saqlang.
