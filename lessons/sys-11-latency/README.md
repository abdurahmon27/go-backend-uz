# SYS-1.01 · Nanosekunddan millisekundgacha: latency jadvalini oʻzimiz oʻlchaymiz

«Kompyuter qanday ishlaydi» modulining birinchi darsi. Oldingi darslar: [GO-1.07 · test va benchmark](../007-test-bench/),
[DSA-A.01 · RAM modeli va kesh](../dsa-a01-cache/) · Keyingi dars: SYS-1.02 · Protsessor nima qiladi: Go funksiyasi
qanday buyruqlarga aylanadi.

**Bitta gʻoya:** kompyuter ichida har bir ish oʻz kutish vaqtiga (kechikish, *latency*) ega, va bu vaqtlar million
barobar farq qiladi: protsessorning eng tez xotirasidan (keshdan) bitta son ~1,5 ns, internetdagi saytdan javob ~0,55 s.

| Fayl | Nima |
|---|---|
| `olcha/` | sekin ishlarni **bittalab** oʻlchaydi (har birini `time.Now` / `time.Since` bilan): `disk` (SSD'dan 4 KiB, `O_DIRECT`; xotiradagi nusxa — page cache), `fsync` (4 KiB yozish: oddiy va `f.Sync()` bilan), `tarmoq` (Wi-Fi orqali routergacha — bizda telefon hotspot'i — TCP ulanish = 1 borib-kelish), `internet` (go.dev gacha TCP ulanish; toʻliq HTTPS soʻrov `net/http/httptrace` bilan 4 bosqichga boʻlingan; faqat IPv4, birinchi urinish isitish uchun tashlanadi) |
| `xotira/` | nanosekundli ishlar — benchmark sikli bilan (`b.Loop`, GO-1.07): `time.Now` narxi, `sync.Mutex` Lock+Unlock, 512 MiB dan 1 MiB ketma-ket oʻqish |
| `jadval/` | `natija.txt` ni mashhur jadval bilan solishtiradi; `-tarqoq` (turlar orasidagi farq), `-xy biz\|mashhur` (video grafigi), `-tur N` (faqat N-tur) |
| `olchov.sh` | hammasi 5 marta navbat bilan → `natija.txt` (bu noutbukda ~3 daqiqa, go.dev ga ~210 ta soʻrov) |
| `natija.txt` | videodagi xom natija (2026-11-26 10:55–10:58, load 1,2–2,6) |
| `uyvazifa/`, `uyvazifa.txt` | uy vazifasi javoblari (Linux) va ularning uchta xom oʻlchovi |

Kesh va operativ xotiradan bitta oʻqish bu yerda **qayta oʻlchanmaydi**: ular DSA-A.01 da zanjir usuli bilan
oʻlchangan (`../dsa-a01-cache/natija.txt`, `Zanjir`, 10 oʻlchov mediani): 8 KiB (L1) → 1,52 ns, 256 KiB (L2) → 4,23 ns,
128 MiB (operativ xotira) → 105,8 ns. Bu — har bir oʻqish oldingisini kutadigan holatdagi *kechikish*; mustaqil,
ketma-ket oʻqishlar bir-birini kutmaydi va tezroq (DSA-A.01: L1 da 0,31 ns/element). `time.Now` oʻzi ~25 ns turadi
(yadro band boʻlsa 40–76 ns), shuning uchun bunday ishlarni bittalab oʻlchab boʻlmaydi.

## Videodagi buyruqlar

Hammasi shu papkadan. `olcha disk` birinchi marta shu papkada 1 GiB lik `sinov.dat` yaratadi (keyin oʻchirib qoʻying).
```bash
./olchov.sh                          # butun oʻlchov → natija.txt va tarqoqlik jadvali
go run ./olcha disk                  # SSD'dan 4 KiB (O_DIRECT) va xotiradagi nusxadan
go run ./olcha fsync                 # 4 KiB yozish: oddiy va fsync bilan
go run ./olcha tarmoq                # Wi-Fi: routergacha borib-kelish (standart: default gateway, 53-port)
grep Wi-Fi natija.txt                # 5 turdagi Wi-Fi natijasi
go run ./olcha internet | head -3    # go.dev: borib-kelish va sahifa soʻrovi
go run ./olcha internet              # + soʻrov bosqichlari
go run ./jadval < natija.txt         # mashhur jadval bilan solishtirish
go run ./jadval -tur 5 -xy biz < natija.txt   # video grafigi (terminal sahnalari ham 5-turdan)
taskset -c 2 go test -run '^$' -bench . -count 10 ./xotira
```
Natijalar `go1.26.8 linux/amd64`, Intel Core i5-12450H, HP SSD EX900 1TB (NVMe, ext4), Wi-Fi (MediaTek mt7921e) orqali
telefon hotspot'i, internet mobil tarmoq orqali. **Video ekranidagi terminal natijalari va grafik = `natija.txt` ning
5-turi.** NixOS'da `gcc` yoʻq, shuning uchun `CGO_ENABLED=0` bilan ishga tushirilgan (Go'ning oʻz DNS resolver'i;
cgo bilan DNS vaqti biroz boshqacha boʻlishi mumkin). `O_DIRECT` va `uyvazifa` faqat Linux'da; boshqa tizimlarda `olcha
disk` oddiy `open` bilan ishlaydi va "SSD" qatori xotiradagi nusxani ham oʻlchab qoʻyishi mumkin.

**Ustunlar.** `marta` — necha marta oʻlchandi; `odatda` — mediana (tartiblangan roʻyxatning oʻrtasi; juft sonda yuqori
oʻrtadagisi); `oraliq` — eng tez va eng sekin 10 % tashlab yuborilgandan keyin qolgani (p10 … p90). Tarmoqda 3 soniyada
javob kelmagan urinish vaqtga qoʻshilmaydi, alohida sanaladi.

**`O_DIRECT` haqiqatan diskka boradimi?** `/proc/diskstats` dagi farq: `olcha disk` paytida qurilma +5066 ta oʻqish,
+73 288 sektor (5000 × 4 KiB + 16 MiB isitish) bajardi — ha.

### Natija (`go run ./jadval < natija.txt`, 5 turning medianasi)
```
amal                  mashhur jadval bu noutbuk      nisbat
L1 keshdan o'qish             0.5 ns    1.52 ns  3.0× sekin
L2 keshdan o'qish               7 ns    4.23 ns    1.7× tez
mutex lock/unlock              25 ns    20.7 ns    1.2× tez
RAM'dan o'qish                100 ns     106 ns           ≈
RAM'dan 1 MiB                 250 µs     149 µs    1.7× tez
nusxadan 4 KiB                     —     605 ns
SSD'dan 4 KiB                 150 µs     141 µs           ≈
4 KiB yozish                       —    1.35 µs
4 KiB yozish + fsync               —     1.4 ms
tarmoq borib-kelish        500 µs DC    4.89 ms  9.8× sekin
internet borib-kelish   150 ms CA-NL     124 ms    1.2× tez
sahifa so'rovi                     —     555 ms
```
- `DC` = mashhur jadvaldagi «Round trip within same datacenter» (serverlar binosi ichida); bizniki — Wi-Fi orqali
  telefongacha. Bir xil narsa emas.
- `CA-NL` = «Send packet CA→Netherlands→CA» (masofa aniq); bizniki — telefon → mobil operator → go.dev javobini beradigan
  Google serverigacha (javob sarlavhasi `server: Google Frontend`; server qayerdaligini oʻlchamadik).
- SSD qatori Dean'ning 2009 slaydida yoʻq, u Bonér gist'idan (2012, «~1GB/sec SSD»).
- mutex va 1 MiB qatorlari 1–3-turlarda yadro band boʻlgani uchun shishgan (mutex 18–35 ns, 4–5-turlarda 11,5 ns;
  1 MiB 98–283 µs, 4–5-turlarda ~70 µs) — videoda ishlatilmaydi.

**Tarqoqlik** (`go run ./jadval -tarqoq`): SSD 135–146 µs, fsync 1,29–1,58 ms, Wi-Fi 4,18–5,00 ms, internet 121–128 ms,
sahifa 547–598 ms. Shu kunning oldingi oʻlchovlarida Wi-Fi medianasi 3–31 ms boʻlgan, eng yomon paytda (10:06–10:13)
141–246 ms va 5–7 % paket yoʻqolgan. Tarmoq — eng beqaror qavat. Oddiy yozish ham tarqoq: 0,67–6,2 µs, shuning uchun
fsync bilan nisbat turlarda 219× dan 2104× gacha.

## Dars matni

**Bitta gʻoya.** Agar protsessorning eng tez xotirasidan (keshdan) bitta sonni olish (1,52 ns) bir soniya davom etsa,
internetdagi saytdan javob (555 ms) 11,6 yil davom etadi.

**Birliklar.** 1 ms = soniyaning 1/1000 (koʻz yumib ochish 100–400 ms, BioNumbers BNID 100706); 1 µs = ms ning 1/1000;
1 ns = µs ning 1/1000 = soniyaning milliarddan biri.

**Kechikish (latency)** — dastur biror narsani soʻragan paytdan javob kelguncha oʻtgan vaqt. DSA-A.01 dagi oʻxshatish:
kesh = choʻntak (darhol), operativ xotira = ombor (borib kelish kerak); bugun yana: disk (SSD) = uzoqdagi arxiv,
tarmoq = qoʻshniga savol yuborib, javob kutish.

**Mashhur jadval qayerdan.** Peter Norvig «Teach Yourself Programming in Ten Years» maqolasiga 2001-yilda «Timing for
various operations on a typical 1GHz PC in summer 2001» jadvalini qoʻshgan (5 qator: L1 2 ns, operativ xotira 10 ns,
disk seek 8 ms; Wayback'da 2001-06-03 da yoʻq, 2001-10-17 da bor). Keng tarqalgan raqamlar Jeff Dean'ning 2009-yilgi
maʼruzalaridan: Stanford CS295 (PDF 2009-03, 13-sahifa) va LADIS 2009 keynote (24-sahifa, «Numbers Everyone Should
Know»): L1 0,5 ns, L2 7 ns, mutex 25 ns (Stanford versiyasida 100 ns), operativ xotira 100 ns, datacenter ichida
borib-kelish 500 µs, CA→NL→CA 150 ms. Norvig 2010-yil iyunida oʻz jadvalini shu raqamlarga yangilagan (Wayback
2010-06-11 → 2010-06-19). Jonas Bonér 2012-05-31 da ularni GitHub gist'ga yigʻgan («Latency Comparison Numbers
(~2012)»), u yerda SSD qatorlari ham bor (4K tasodifiy 150 µs, 1 MB ketma-ket 1 ms, «~1GB/sec SSD»).

**Qanday oʻlchaymiz.** `olcha` funksiyasi: `boshi := time.Now()` → ish → `time.Since(boshi)`, koʻp marta (30 dan 5000
gacha), keyin tartiblab mediana va p10…p90. Nanosekundli ishlar uchun benchmark sikli kerak (GO-1.07: `b.Loop`),
xotira raqamlari DSA-A.01 dan.

**Disk.** `O_DIRECT` bayrogʻi (open(2)) Linux'ning operativ xotiradagi nusxasini (page cache) chetlab oʻtadi. SSD'dan
tasodifiy 4 KiB — 146 µs (5-tur); xuddi shu oʻqish nusxadan (`pread` tizim chaqiruvi + nusxa koʻchirish) — 492 ns,
~300 barobar tez.

**Keng tarqalgan xato: «yozdim — demak diskda».** `write(2)`: muvaffaqiyatli `write` maʼlumot diskka yozilganini
kafolatlamaydi. `fsync(2)` qurilma «yozdim» deguncha kutadi (`O_SYNC`, `fdatasync` kabi boshqa yoʻllar ham bor; diskning
oʻz keshi va fayl tizimi sozlamalari ham rol oʻynaydi). 4 KiB yozish 732 ns, yozish + `f.Sync()` 1,29 ms — 1762 barobar
(5-tur; turlarda 219–2104×). Bu — shu noutbukning diski; elektr uzilishidan himoyalangan server disklari fsync'ni ancha
tez bajaradi. Shuning uchun baʼzi maʼlumotlar bazalari bir nechta tranzaksiyani bitta fsync bilan yozadi (PostgreSQL:
group commit, `commit_delay`). Batafsil: SYS-2.04.

**Tarmoq.** TCP ulanish = SYN ketadi, SYN-ACK qaytadi: bitta borib-kelish (round trip, RFC 9293 §3.5). Port yopiq
boʻlsa, RST qaytadi — bu ham bitta borib-kelish. Wi-Fi orqali telefonga: 4,89 ms (5-tur) — SSD'dan 33 barobar uzoq.

**Internet.** go.dev javobini beradigan Google serverigacha bitta borib-kelish 128 ms (5-tur); `https://go.dev/robots.txt`
ni har safar yangi ulanish bilan soʻrash 555 ms. Bosqichlar (medianalar, ular yigʻindiga teng emas): manzil topish (DNS,
javob telefondan) 13 ms, TCP ulanish 121 ms, TLS kelishuvi 154 ms (TLS 1.3 — bitta borib-kelish, RFC 8446 §2, yana
shifrlash hisob-kitobi), soʻrov → birinchi bayt 261 ms (borib-kelish + server ishi). Ochiq ulanishni qayta ishlatish
(keep-alive) DNS, TCP va TLS bosqichlarini olib tashlaydi: dasturning oldingi versiyasida oʻlchangan — ochiq ulanish bilan
~0,28 s, yangi ulanish bilan ~0,7 s (episodes/sys-1-01/research.md §3).

**Mashhur jadval bilan.** Kesh/xotira bir xil tartibda (1,5 va 106 ns), SSD 146 vs 150 µs, internet 128 vs 150 ms;
tarmoq ~10× sekin, lekin jadvaldagi raqam data markaz ichidagi tarmoq uchun, bizniki — Wi-Fi va telefon.

**Inson vaqtida** (keshdan bitta son = 1 s, 5-tur): operativ xotira 70 s, SSD 27 soat, Wi-Fi borib-kelish 37 kun,
internet borib-kelish 2,7 yil, saytdan javob 11,6 yil.

## Manbalar
- Peter Norvig, «Teach Yourself Programming in Ten Years», jadval: https://norvig.com/21-days.html#answers ·
  2001 versiyasi: https://web.archive.org/web/20011017203742/http://norvig.com/21-days.html ·
  2010 yangilanishi: https://web.archive.org/web/20100619021128/http://norvig.com/21-days.html
- Jeff Dean, «Designs, Lessons and Advice from Building Large Distributed Systems», LADIS 2009, 24-sahifa:
  https://www.cs.cornell.edu/projects/ladis2009/talks/dean-keynote-ladis2009.pdf
- Jeff Dean, «Software Engineering Advice from Building Large-Scale Distributed Systems» (Stanford CS295), 13-sahifa:
  https://static.googleusercontent.com/media/research.google.com/en//people/jeff/stanford-295-talk.pdf
- Jonas Bonér, «Latency Numbers Every Programmer Should Know» (gist, 2012-05-31): https://gist.github.com/jboner/2841832
- Colin Scott, «Numbers Every Programmer Should Know By Year» (yillar boʻyicha taxmin, oʻlchov emas):
  https://colin-scott.github.io/personal_website/research/interactive_latency.html
- open(2) `O_DIRECT`: https://man7.org/linux/man-pages/man2/open.2.html · fsync(2): https://man7.org/linux/man-pages/man2/fsync.2.html ·
  write(2): https://man7.org/linux/man-pages/man2/write.2.html · `os.File.Sync`: https://pkg.go.dev/os#File.Sync
- PostgreSQL, WAL sozlamalari (`commit_delay`, group commit): https://www.postgresql.org/docs/current/runtime-config-wal.html
- RFC 9293 (TCP) §3.5: https://www.rfc-editor.org/rfc/rfc9293#section-3.5 · RFC 8446 (TLS 1.3) §2: https://www.rfc-editor.org/rfc/rfc8446#section-2
- `net/http/httptrace`: https://pkg.go.dev/net/http/httptrace · `time.Since`: https://pkg.go.dev/time#Since · `testing.B.Loop`: https://pkg.go.dev/testing#B.Loop
- BioNumbers BNID 100706 (koʻz yumish 0,1–0,4 s): https://bionumbers.hms.harvard.edu/bionumber.aspx?id=100706
- DSA-A.01 kesh oʻlchovi: [../dsa-a01-cache/](../dsa-a01-cache/)

## Uy vazifasi

1. 🟢 **Oson.** `go run ./olcha disk` ni oʻz kompyuteringizda ishga tushiring: SSD qatori bizniki (146 µs) bilan qanday?
   Qoʻshimcha: `go run ./olcha -host <boshqa sayt> internet` — borib-kelish oʻzgaradimi?
2. 🟡 **Oʻrta.** 10 ta yozuvga bitta `fsync` qiling: bitta yozuvga qancha vaqt toʻgʻri keladi?
3. 🟠 **Qiyin.** Kompyuter oʻzi bilan oʻzi tarmoq orqali gaplashsa (localhost, TCP), borib-kelish qancha? Unix socket bilan-chi?
4. 🔴 **Guru.** SSD'dan 32 ta oʻqishni bir vaqtda soʻrang: bitta oʻqish kechikishi va soniyadagi oʻqishlar soni qanday oʻzgaradi?

### Javoblar (`go run ./uyvazifa guruh|unix|parallel`, xom natija: `uyvazifa.txt`, uchta oʻlchov)

1. Javob diskingizga bogʻliq. Bizda 135–146 µs (5 tur). Aylanadigan diskda (HDD) mashhur jadval bitta tasodifiy
   oʻqishga ~10 ms beradi (Dean: «Disk seek 10,000,000 ns»).
   Qoʻshimcha: `-host kun.uz` — borib-kelish 113–116 ms, toʻliq soʻrov 480–586 ms, go.dev (128 / 555 ms) ga yaqin. Bizning
   tarmoqda (telefon → mobil operator) kechikishning katta qismi yoʻlning boshida boʻlishi mumkin; `traceroute` bilan
   tekshirib koʻring (SYS-3 darslarida).
2. Har yozuvga fsync: 1,8–3,5 ms/yozuv; 10 ta yozuvga bitta fsync: 0,13–0,66 ms/yozuv — 5–19 barobar tez (2-oʻlchovda
   disk boshqa jarayonlar bilan band edi). Bu maʼlumotlar bazalaridagi *group commit* gʻoyasi: narx asosan fsync'da,
   yozuvlar sonida emas. Bahosi: bitta fsync oldidan 10 ta yozuv xavf ostida — elektr oʻchsa, ular birga yoʻqoladi.
3. TCP localhost: mediana 6,9–45 µs; Unix domain socket: 2,9–17 µs — har bir oʻlchovda Unix socket 1,4–2,7 barobar tez.
   Tarmoq kartasi ishlatilmaydi, lekin TCP qatlami (paket, tasdiq) baribir ishlaydi; Unix socket uni chetlab oʻtadi.
   Oʻlchovlar orasidagi katta farq (6,9 va 45 µs) protsessorning yuklamasi va quvvat holatiga bogʻliq boʻlsa kerak —
   tekshirilmagan. Har holda Wi-Fi'dagi ~5 ms dan yuz barobar tez.
4. 1 goroutine: 121–156 µs/oʻqish (yuklangan 2-oʻlchovda 412 µs), ~4 000–7 000 oʻqish/s; 32 ta: 318–370 µs,
   ~50 000–75 000/s. Bitta oʻqish 2–3 barobar sekinlashdi, lekin soniyadagi oʻqishlar 10 barobardan koʻproq oshdi: NVMe
   SSD bir vaqtda koʻp soʻrovni parallel bajaradi. Kechikish (latency) va oʻtkazish qobiliyati (throughput) — har xil narsa.
