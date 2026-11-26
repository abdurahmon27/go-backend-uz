# MATH-0.01 · Muhandisga qancha matematika kerak? Halol xarita

Muhandislar uchun matematika yoʻnalishining birinchi darsi · Keyingi dars: MATH-0.02 · 0.1 + 0.2 yana: `float64` ichida nima bor

| Papka | Nima |
|---|---|
| `fanout/` | bitta sahifa N ta soʻrovni parallel yuboradi va eng sekinini kutadi; har soʻrov 100 tadan 1 holatda ≥ 1 s. Sahifalarning qanchasi sekin? (`go run ./fanout [N]`, standart N = 100; `go run ./fanout jadval` — 1 − 0.99ᴺ jadvali) |
| `latency/` | **sintetik** model, ikki reliz: `eski` (har soʻrov bazaga, ~35 ms) va `kesh` (96 % keshdan ~8 ms, 4 % sekin yoʻl ~400 ms). 100 000 soʻrovning oʻrtacha/p50/p99 si va 10 000 ta 100 soʻrovli sahifaning p50/p99 si |
| `abtest/` | tugma testini 10 000 marta simulyatsiya qiladi (`go run ./abtest [n]`, standart n = 2000 har guruhda). A/A: tugmalar bir xil (5 %) — necha marta "±10 %" koʻrinadi va z-test necha marta adashadi. A/B: yangi tugma rostdan +10 % (5.5 %) — z-test necha marta topadi. 10 goroutine, har biri oʻz PCG generatori bilan: natija har safar bir xil |
| `homework/` | uy vazifasi javoblari (`go run ./homework`) |

Hamma tasodifiy sonlar `math/rand/v2` PCG generatoridan, doimiy seed bilan (`fanout`, `latency`: `(1, 2)`; `abtest`: `(1, ishchi)`): natija har safar bir xil. Notoʻgʻri argument (`fanout 0`, `abtest 5`) — foydalanish xabari va chiqish kodi 2.

## Videodagi buyruqlar va natijalar
Hammasi shu papkadan (`lessons/math-01-map`). Natijalar `go1.26.8 linux/amd64` da olingan.
```
$ go run ./fanout
so'rovlar: 1 s dan sekini  1.00 %
sahifalar (100 so'rovli): 1 s dan sekini 63.5 %
qo'lda hisob: 1 - 0.99^100 = 63.4 %
$ go run ./latency
eski so'rov:  o'rtacha   37 ms | p50   35 ms | p99   70 ms
eski sahifa: 100 so'rov   | p50   73 ms | p99  107 ms
kesh so'rov:  o'rtacha   25 ms | p50    8 ms | p99  492 ms
kesh sahifa: 100 so'rov   | p50  531 ms | p99  918 ms
$ go run ./abtest
har guruhda 2000 kishi, 10000 tajriba
A/A, tugmalar bir xil:
  yangi tugma >= +10 %: 24.4 %
  yangi tugma <= -10 %: 22.4 %
  z-test "farq bor":     4.8 %
A/B, yangi tugma rostdan +10 %:
  z-test "farq bor":    10.9 %
qo'lda hisob: shovqin ~ sqrt(2*0.05*0.95/2000) = 0.69 foiz punkti
$ go run ./abtest 30400
har guruhda 30400 kishi, 10000 tajriba
A/A, tugmalar bir xil:
  yangi tugma >= +10 %:  0.4 %
  yangi tugma <= -10 %:  0.2 %
  z-test "farq bor":     5.0 %
A/B, yangi tugma rostdan +10 %:
  z-test "farq bor":    79.2 %
qo'lda hisob: shovqin ~ sqrt(2*0.05*0.95/30400) = 0.18 foiz punkti
$ go run ./fanout 10
so'rovlar: 1 s dan sekini  0.99 %
sahifalar (10 so'rovli): 1 s dan sekini  9.4 %
qo'lda hisob: 1 - 0.99^10 = 9.6 %
$ go run ./fanout jadval
N =   1:   1.0 % sahifa sekin
N =  10:   9.6 % sahifa sekin
N =  20:  18.2 % sahifa sekin
N =  50:  39.5 % sahifa sekin
N = 100:  63.4 % sahifa sekin
N = 200:  86.6 % sahifa sekin
N = 500:  99.3 % sahifa sekin
```
A/A da z-test ~5 % adashadi, A/B da (n = 30 400) ~79 % topadi: bu Kohavi formulasi vaʼda qilgan 95 % ishonch va 80 %
"power" (topish ehtimoli). n = 2000 da esa rostdan +10 % yaxshi tugma faqat ~11 % holatda topiladi.

## Dars matni

**Hook.** Har bir soʻrov 100 tadan 1 holatda 1 soniyadan sekin. Sahifa 100 ta soʻrovni parallel yuborib, eng
sekinini kutsa, sahifalarning 63 % sekin ochiladi. Misol Dean & Barroso, *The Tail at Scale* (CACM 2013, 76-bet)
maqolasidan: "If a user request must collect responses from 100 such servers in parallel, then 63% of user requests
will take more than one second." (Maqolada bu "99th-percentile latency of one second" deb atalgan; bizning
taʼrifimizda bu "1 % soʻrov 1 s dan sekin" degani — p99 ning oʻzi esa bu modelda ~50 ms.)

**Qoʻlda hisob.** Bitta soʻrov tez: 0.99. Soʻrovlar bir-biriga bogʻliq boʻlmasa (mustaqil), ehtimollar
koʻpaytiriladi: hammasi tez = 0.99¹⁰⁰ ≈ 0.366. Kamida bittasi sekin = 1 − 0.366 ≈ 63.4 %. Simulyatsiya 63.5 % berdi.
Ehtiyot boʻling: kod ham mustaqillik bilan yozilgan, shuning uchun ular mos kelishi faqat arifmetikani tekshiradi.
Haqiqiy tizimda sekinliklar koʻpincha bogʻliq (bitta GC pauzasi, bitta tiqilgan switch) — taxminni faqat haqiqiy
oʻlchov tekshiradi.

**Google servisida (Dean & Barroso, Table 1).** Bitta tasodifiy soʻrovning p99 i 10 ms, hamma soʻrovlar tugashini
kutishning p99 i 140 ms. Yechim — *hedged requests*: javob kechiksa, xuddi shu soʻrovni boshqa nusxaga ham yuborish.
Ularning testida (1000 ta kalit, 100 server, 10 ms dan keyin hedge) p99.9 1800 ms dan 74 ms ga tushgan, soʻrovlar
atigi 2 % koʻpaygan.

**Uchta oʻlchov.** Oʻrtacha = yigʻindi / soni. Mediana (p50) = saralangan roʻyxatning oʻrtasi. p99 = soʻrovlarning
99 % undan tez yoki unga teng; undan sekin 1 % — "dum". Biz "eng yaqin rang" (*nearest rank*) usulini ishlatamiz:
saralangan `xs` da `xs[ceil(p·n) − 1]`. [06-dars · Funksiyalar](../006-functions/) dagi `mediana` esa `xs[len/2]`
ni oladi: juft n da u yuqori oʻrtani, bizniki pastki oʻrtani tanlaydi (100 000 soʻrovda farq sezilmaydi).

**Qaror: qaysi reliz?** (Sintetik model: sonlarni biz tanladik, shakli real — kesh topmagan soʻrov sekin yoʻlga
ketadi.) `kesh` oʻrtacha (25 vs 37 ms) va mediana (8 vs 35 ms) boʻyicha yutadi, p99 boʻyicha esa 7 barobar
yutqazadi (492 vs 70 ms). Sahifa 100 ta soʻrov yuborsa: 1 − 0.96¹⁰⁰ ≈ 98.3 % sahifada kamida bitta sekin javob bor;
oʻlchov: sahifa p50 i keshli relizda 531 ms, eskisida 73 ms. Agar miss arzon boʻlsa (masalan ~40 ms), kesh hamma
oʻlchovda yutadi — qarorni oʻrtacha emas, dum hal qiladi. SRE kitobi: latency ni oʻrtacha bilan emas, oraliqlar
boʻyicha sanalgan gistogramma bilan kuzating.

**A/B test.** Har tugmani 2000 kishi koʻrdi: eski 100, yangi 110 xarid (+10 %). A/A simulyatsiyada (tugmalar bir xil)
24.4 % tajribada yangi tugma "≥ +10 %", 22.4 % da "≤ −10 %" koʻrinadi: qariyb yarmida 10 % lik farq tasodifdan.
Ikki ulushning z-testi (|z| > 1.96) esa A/A da ~5 % adashadi, lekin rostdan +10 % yaxshi tugmani n = 2000 da faqat
~11 % holatda topadi. Kohavi va boshq. (2009) A/A testni shunday tavsiflaydi: ikkala guruh bir xil tajribani koʻradi.

**Qoʻlda hisob.** Ikki ulush farqining standart xatosi (NIST e-Handbook 7.3.3, teng n): √(2·p(1−p)/n) =
√(2·0.05·0.95/2000) ≈ 0.69 foiz punkti. "+10 %" esa 5 % ning 10 % i = 0.5 foiz punkti: shovqindan kichik.
95 % ishonch va 80 % topish uchun shovqin farqdan ~2.8 (1.96 + 0.84) barobar kichik boʻlishi kerak: 0.5 / 2.8 ≈
0.18 pp. Shovqin 1/√n kabi kamayadi: (0.69 / 0.18)² ≈ 15 barobar koʻp odam. Kohavi va boshq. 2009, formula (2)
(ular van Belle 2002 ga tayanadi): n = 16σ²/Δ², Δ — sezmoqchi boʻlgan eng kichik farq: 16 · 0.0475 / 0.005² =
**30 400 kishi har guruhga**. Simulyatsiya: z-test A/A da 5.0 % adashadi, rost +10 % ni 79.2 % topadi.
[07-dars](../007-test-bench/) dagi `benchstat` ham shu savolga javob beradi: farq haqiqiymi yoki shovqinmi.
Har kuni natijaga qarash (peeking), koʻp metrika va effect size — MATH-4.07, 4.08, 4.09 da.

**Oltita modul.**
| Modul | Nima | Qaysi ishga |
|---|---|---|
| MATH-0 Kirish | float64, yaxlitlash | sonli xatolar |
| MATH-1 Chiziqli algebra | vektor, matritsa, proyeksiya, xos vektor, SVD | qidiruv (embedding + kosinus), PageRank, PCA |
| MATH-2 Hosila va gradient | sezgirlik, zanjir qoidasi, gradient descent, autodiff | model oʻqitish |
| MATH-3 Ehtimollik | Bayes, kutilma, taqsimotlar, CLT, entropiya | hash toʻqnashuvi, Bloom filter, simulyatsiya |
| MATH-4 Statistika | percentile, ishonch oraligʻi, bootstrap, A/B test | latency, A/B test, SLO |
| MATH-5 ML'ga koʻprik | chiziqli va logistik regressiya | ML trekiga kirish |

**Kimga qancha.** Backend muhandisiga MATH-0, MATH-3, MATH-4 (sonlar, ehtimollik, statistika). ML ga oʻtmoqchi
boʻlsangiz: yana MATH-1, MATH-2, MATH-5. 56 darsdan 20 tasi ixtiyoriy (`plan/tracks/math.yaml` da `core: false`).
Bugungi hisoblar uchun faqat foiz va daraja kerak edi.

**Nimani oʻtmaymiz.** Big-O va graflar — DSA yoʻnalishida; sonlar nazariyasi (kriptografiya) — bu yoʻnalishda yoʻq.
Rasmiy isbot yozish (MIT 6.042 da "definitions, proofs, sets, functions, relations" kursning taxminan uchdan biri).
Integral va determinantni qoʻlda hisoblash (maʼnosi qoladi, hisobni kompyuter qiladi). Statistik testlar jadvali
(Downey: aslida bitta test — nol gipotezani simulyatsiya qilib, natijani solishtirish; bugungi A/A aynan shunday).

**Usul.** Har mavzuda bitta rasm, bitta qoʻlda hisob va bitta kod. Mos kelmasa — formulada yoki kodda xato bor;
mos kelsa ham, taxminni (masalan, mustaqillikni) faqat haqiqiy oʻlchov tekshiradi. Grafiklar Python/NumPy'da, simulyatsiya va algoritmlar Go'da.

## Uy vazifasi
- **Oson:** 20 soʻrovli sahifa, har soʻrov 100 tadan 1 holatda 1 s dan sekin. Necha foiz sahifa sekin? Qoʻlda toping, keyin `go run ./fanout 20` bilan tekshiring.
- **Oʻrta:** har soʻrov 1000 tadan 1 holatda 1 s dan sekin. 100 soʻrovli sahifa?
- **Qiyin:** 5 % konversiyada +5 % nisbiy oʻsishni (5 % → 5.25 %) koʻrish uchun maqoladagi formula bilan har guruhga necha kishi kerak?
- **Guru:** `kesh` relizida p99 < 100 ms boʻlishi uchun kesh soʻrovlarning necha foizini ushlashi kerak?

### Javoblar (`go run ./homework`)
```
oson:  1 - 0.99^20 = 18.2 % (tekshiring: go run ./fanout 20)
o'rta: 1 - 0.999^100 = 9.5 % (100 tadan 1 bo'lganda 63.4 % edi)
qiyin: n = 16*0.0475/0.0025^2 = 121600 kishi har guruhda
guru:  kesh 96.0 % -> p99  492 ms
guru:  kesh 98.0 % -> p99  399 ms
guru:  kesh 99.0 % -> p99  185 ms
guru:  kesh 99.5 % -> p99   17 ms
```
- Oson: `fanout 20` simulyatsiyasi 18.0 % beradi.
- Oʻrta: 63.4 % → 9.5 %. 100 soʻrovli sahifa uchun har soʻrovning p99 emas, p99.9 i muhim.
- Qiyin: 121 600; Kohavi maqolasining oʻzi xuddi shu misolda "less than 122,000 users" deydi.
- Guru: miss ulushi 1 % dan kam boʻlishi kerak (99 % da hali chegarada, 185 ms), yaʼni kesh > 99 %. Sababi: p99 —
  eng sekin 1 %; misslar 1 % dan koʻp boʻlsa, p99 albatta miss ichiga tushadi.

## Manbalar
- Google SRE book, ch. 6 "Monitoring Distributed Systems" — *Worrying About Your Tail*: https://sre.google/sre-book/monitoring-distributed-systems/
- Kohavi, Longbotham, Sommerfield, Henne. *Controlled experiments on the web: survey and practical guide*. Data Mining and Knowledge Discovery, 2009 — formula (2) `n = 16σ²/Δ²`, A/A test: https://ai.stanford.edu/~ronnyk/2009controlledExperimentsOnTheWebSurvey.pdf
- NIST/SEMATECH e-Handbook 7.3.3 — ikki ulushni solishtirish: https://www.itl.nist.gov/div898/handbook/prc/section3/prc33.htm
- Allen Downey, *There is still only one test* (2016): http://allendowney.blogspot.com/2016/06/there-is-still-only-one-test.html
- Stanford CS109 *Probability for Computer Scientists* (Independence, CLT, Bootstrapping, MLE, Logistic Regression): https://web.stanford.edu/class/cs109/
- MIT 18.06SC *Linear Algebra* (Strang): https://ocw.mit.edu/courses/18-06sc-linear-algebra-fall-2011/
- MIT 6.042J *Mathematics for Computer Science*: https://ocw.mit.edu/courses/6-042j-mathematics-for-computer-science-spring-2015/
- Dean, Barroso. *The Tail at Scale*. CACM 56(2), 2013, 74–80-betlar — 63 % misoli, Table 1, hedged requests: https://research.google/pubs/the-tail-at-scale/ (toʻliq matn: https://web.archive.org/web/2019id_/https://www.barroso.org/publications/TheTailAtScale.pdf)
- Brin, Page. *The Anatomy of a Large-Scale Hypertextual Web Search Engine*, 1998, §2.1.1 (PageRank — xos vektor): https://snap.stanford.edu/class/cs224w-readings/Brin98Anatomy.pdf
- pgvector README — kosinus masofa (`<=>`): https://github.com/pgvector/pgvector
- `math/rand/v2`: https://pkg.go.dev/math/rand/v2
