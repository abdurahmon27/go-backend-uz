# DB-1.01 · Nega JSON fayl emas, baza? Ikki Go jarayoni bitta faylni buzadi

Maʼlumotlar bazasi trekining birinchi darsi. Oldingi: [SYS-0.06 · docker compose: Postgres va Go ilovasi](../sys-06-compose/)
(Postgres muhiti), [SYS-0.02 · jarayon va signal](../sys-02-process-env/), GO-1 asoslari. Keyingi: DB-1.02 · Relyatsion model:
jadval — bu toʻplam.

| Fayl | Nima |
|---|---|
| `fayl/main.go` | buyurtmalarni bitta JSON faylda saqlaydigan "baza": har buyurtmada oʻqi → qoʻsh → `os.WriteFile`. `-usul rename` (vaqtinchalik fayl + `rename`), `-usul qulf` (`flock` + `rename`), `-sana` (sanash va buzuq faylga tashxis) |
| `poyga.sh [usul] [marta]` | ikki jarayon (A, B) bir vaqtda 1000 tadan yozadi; koʻp marta boʻlsa: fayl necha marta buzildi, nechta buyurtma yoʻqoldi |
| `qulash.sh [usul] [marta]` | bitta jarayon toʻxtovsiz yozadi, tasodifiy paytda (0.1–0.9 s) `kill -9`; fayl butunmi? |
| `compose.yaml`, `.env.example` | `postgres:18-alpine`, healthcheck, volume (SYS-0.06 dagidek), port faqat `127.0.0.1:5440` |
| `baza/main.go` | xuddi shu buyurtmalar Postgres'da: `-usul json` (JSON bitta `jsonb` ustunda, Go oʻqiydi/qoʻshadi/yozadi) va `-usul insert` (har buyurtma bitta `INSERT`) |
| `baza.sh [usul] [marta]` | `poyga.sh` ning Postgres varianti |
| `qulash-baza.sh [marta]` | A va B toʻxtovsiz `INSERT … RETURNING id` qiladi va tasdiq olgan har bir id'ni yozib boradi, `docker compose kill db` (SIGKILL), bazani qayta koʻtarib har bir id joyidami, tekshiradi (`baza -tekshir`) |
| `homework/` | uy vazifasi javoblari: `bash homework/javoblar.sh` |

## Videodagi buyruqlar
Shu papkadan. Postgres qismi uchun avval `cp .env.example .env`, `.env` ga parol yozing va `docker compose up -d --wait`
(5440 port boʻsh boʻlsin). Skriptlar `.env` ni oʻzi oʻqiydi va `bin/` ga yigʻadi (`CGO_ENABLED=0`).
```bash
./poyga.sh                      # ikki jarayon: fayl buzildi, "boshida N ta buyurtma, ortida yana ~70 bayt eski yozuv"
tail -n 9 buyurtmalar.json      # }{ — ikki JSON yopishib qolgan
./poyga.sh oddiy 20             # 20 urinish: necha marta buzildi
./qulash.sh oddiy 100           # bitta jarayon + kill -9: nol bayt yoki yarim fayl
./qulash.sh rename 100          # rename bilan: har doim butun
./poyga.sh rename               # fayl butun, lekin ~1000 ta buyurtma yoʻq (lost update)
./poyga.sh rename 20
./poyga.sh qulf 20              # flock + rename: 2000 ta
./baza.sh json                  # JSON Postgres ustunida: baribir ~1000 ta
./baza.sh json 20
./baza.sh insert 20             # INSERT: har safar 2000
docker compose exec db psql -U ombor -c 'SELECT mijoz, count(*), sum(summa) FROM buyurtmalar GROUP BY mijoz ORDER BY mijoz'
./qulash-baza.sh                # kill -9 → "not properly shut down; automatic recovery", redo
./qulash-baza.sh 20 | tail -1
# tozalash
docker compose down -v; rm -rf bin buyurtmalar.json*
```

## Natijalar (bu kompyuterda)
Go 1.26.8, Linux 6.18 (NixOS), ext4, i5-12450H (12 ta mantiqiy yadro), Docker Engine 29.8.1, Compose 5.4.0,
`postgres:18-alpine` (PostgreSQL 18.7), pgx v5.13.0, 2026-11-26. Poyga natijalari har safar boshqacha — videoda bitta
yozib olingan ishga tushirish koʻrsatiladi, quyida hamma oʻlchovlar.

| Sinov | Natija |
|---|---|
| `poyga.sh oddiy 20` (ikki marta) | fayl **20/20** buzildi (ikkinchi marta ham 20/20); buzilganda boshida 1–289 (boshqa safar 5–175) ta butun buyurtma |
| `qulash.sh oddiy 100` (olti marta) | buzildi **10/100** (9 boʻm-boʻsh, 1 yarim) — videodagi; boshqa safarlar 2, 10, 19, 9, 11 |
| `qulash.sh rename 100`, `qulash.sh qulf 100` | 100/100 butun |
| `poyga.sh rename 20` (uch marta) | fayl 0/20 buzildi; 2000 dan oʻrtacha **983** (962, 984) ta yoʻqoldi, eng kami 964, eng koʻpi 999 |
| `poyga.sh qulf 20` (ikki marta) | 0 yoʻqoldi, 20/20 |
| `baza.sh json 20` (ikki marta) | oʻrtacha **999** yoʻqoldi (994–1000) |
| `baza.sh insert 20` (ikki marta) | 0 yoʻqoldi |
| `qulash-baza.sh 20` | har bir id alohida tekshirildi: tasdiqlangan id'lar **20/20** joyida (yoʻqolgan 0); 17 urinishda bazada bitta tasdiqsiz ortiqcha qator. Avvalgi (faqat sonni solishtirgan) ikki seriya: 40/40 |
| uy vazifasi (Guru) | `qulf`: 6.1 / 7.0 s, `qulf` + fsync: 15.4 / 13.6 s (ikki jarayon × 1000, ikki marta) |

## Dars matni (videodagi soʻzlar)

<!-- matn: lesson.yaml dan -->
**hook.** Ikki Go dasturi bitta JSON faylga buyurtma yozdi va faylni buzdi. A va B ikki alohida jarayon, har biri ming tadan buyurtma yozishi kerak edi. Ikkalasi ham faylni qayta-qayta oʻqishga urindi va oxiri toʻxtadi: fayl endi JSON emas. Boshida oltmish yettita butun buyurtma bor, ortida esa yetmish bayt eski yozuv qolgan. Go’dagi Mutex bu yerda yordam bermaydi: u bitta jarayon ichidagi goroutine'larni tartibga soladi, bular esa ikki alohida jarayon.

**trek.** Bu ma'lumotlar bazasi trekining birinchi darsi. Trekda SQL yozishni, Go’dan Postgres'ga ulanishni, index va transaction'larni oʻrganamiz, oxirida esa storage engine'ni Go’da oʻzimiz quramiz. Lekin avval eng oddiy savol: fayl nega yetmaydi? Karnegi Mellon universitetining ma'lumotlar bazasi kursi ham birinchi ma'ruzasida aynan shu savolni beradi. Bugun uch muammoni oʻz koʻzimiz bilan koʻramiz: bir vaqtda yozish, crash va soʻrov. Docker darsidagi mehmon serveri ham faylga yozardi, shuning uchun bu bizning muammomiz ham.

**kod.** Kod Go asoslari darslaridagi narsalardan iborat. Birinchi funksiya faylni toʻliq oʻqiydi va JSON'ni struct'ga aylantiradi. Har bir buyurtma uchun dastur faylni oʻqiydi, roʻyxat oxiriga bitta buyurtma qoʻshadi va hammasini qayta yozadi. os.WriteFile avval faylni boʻshatadi, keyin yangi matnni yozadi. Go hujjatida ochiq yozilgan: bu bir nechta tizim chaqiruvi, yarim yoʻlda uzilsa, fayl chala qoladi.

**buzilish.** Faylning oxiriga qarasak, yopuvchi qavsdan keyin darhol yangi ochuvchi qavs keladi. Ikkala jarayon faylni deyarli bir vaqtda ochdi va boʻshatdi. Biri uzunroq JSON'ni yozdi, ikkinchisi esa qisqaroq JSON'ni faylning boshidan ustiga yozdi. Uzun yozuvning dumi joyida qoldi. Yigirma marta takrorladik va fayl har safar buzildi: ba'zan birinchi buyurtmadayoq, eng kechi ikki yuz sakson toʻqqizinchisida.

**qulash.** Endi raqibni olib tashlaymiz, faylga bitta jarayon yozadi. U toʻxtovsiz yozadi, biz esa uni tasodifiy paytda kill minus toʻqqiz bilan oʻldiramiz. Jarayonlar darsidan eslang: bu signalni ushlab boʻlmaydi, dastur hech narsani tugatishga ulgurmaydi. Yuz marta sinadik va oʻn marta fayl buzildi. Toʻqqiz marta u nol baytli boʻlib qoldi, ya'ni hamma buyurtma yoʻqoldi. Kill fayl boʻshatilgan, lekin hali yozilmagan lahzaga toʻgʻri keldi. Boshqa yuztalik sinovlarda buzilgan fayllar ikkitadan oʻn toʻqqiztagacha boʻldi.

**rename.** Birinchi tuzatish: faylni joyida oʻzgartirmaymiz. Yangi JSON avval yonidagi vaqtinchalik faylga yoziladi, keyin rename uni eski faylning oʻrniga qoʻyadi. Linux'da rename atomik: boshqa jarayon eski yoki yangi faylni koʻradi, yarim faylni hech qachon koʻrmaydi. Ikkinchi tuzatish qulf: flock bilan hamma jarayon bitta umumiy qulf oladi. Oʻqish, qoʻshish va yozish qulf ichida bajariladi, keyin qulf boʻshaydi.

**rename-run.** Avval faqat rename bilan sinaymiz: yuz marta kill qildik, fayl yuz marta butun. Ikki jarayonni poygaga qoʻyamiz: ikkalasi ham ming tadan yozdim deydi, buzuq oʻqish nol. Lekin sanaymiz: faylda bor-yoʻgʻi ming oʻn uchta buyurtma. Yigirma urinishda ikki mingdan oʻrtacha toʻqqiz yuz sakson uchta buyurtma yoʻqoldi. Endi qulf bilan: yigirma urinish, yoʻqotish nol.

**lost.** Bu xatoning nomi lost update, ya'ni yoʻqolgan yangilanish. A faylni oʻqidi, unda yuzta buyurtma bor. B ham xuddi oʻsha yuztani oʻqidi. A oʻzinikini qoʻshib, yuz bir buyurtmani yozdi. B ham oʻzinikini qoʻshib, oʻzining yuz birtasini yozdi va A'ning buyurtmasini oʻchirib yubordi. Fayl butun, xato yoʻq, lekin buyurtma yoʻq. Bu buzilgan fayldan ham xavfliroq, chunki uni hech kim sezmaydi.

**royxat.** Qulf bilan ishladi, lekin bu hali boshlanishi. flock maslahat qulfi: uni ishlatmaydigan dastur faylga bemalol yozaveradi. Elektr oʻchishiga qarshi faylni ham, papkani ham fsync qilish, ya'ni diskka majburan tushirish kerak. Har bir buyurtma uchun butun fayl qayta yoziladi. Bitta buyurtmani topish uchun esa hammasini oʻqiysiz. Shularning hammasini hal qilsangiz, oʻzingiz ma'lumotlar bazasi yozgan boʻlasiz.

**savol.** Endi savol: JSON'ni fayldan olib, Postgres jadvalidagi bitta ustunga qoʻysak-chi? Go kodi xuddi oʻsha: oʻqi, qoʻsh, yoz. Postgres ikki ming buyurtmani saqlab qoladimi? Videoni pauza qiling va oʻylab koʻring.

**pg-json.** Postgres'ni oʻtgan darsdagidek compose bilan koʻtaramiz, bu safar Go jarayonlari xostda ishlaydi. Natija: json ustunida ikki mingdan bor-yoʻgʻi ming ikkita. Yigirma urinishda oʻrtacha toʻqqiz yuz toʻqson toʻqqizta buyurtma yoʻqoldi. Baza sehrli emas, u faqat oʻziga berilgan amalni himoya qiladi. Oʻqish bilan yozish orasidagi oraliq Go kodingizda qoldi, baza uni koʻrmaydi.

**pg-insert.** Toʻgʻri yoʻl: oʻzgarishni bazaning oʻziga beramiz. Har bir buyurtma bitta INSERT, ya'ni jadvalda bitta yangi qator, oldindan hech narsa oʻqilmaydi. Yigirma urinish, har safar ikki ming, yoʻqotish nol. Va uchinchi muammo, soʻrov: kim nechta buyurtma berganini va summasini bitta SQL hisoblaydi. Fayl dasturida buning uchun butun faylni oʻqib, Go’da sikl yozganmiz.

**pg-qulash.** Endi crash sinovi: A va B toʻxtovsiz yozyapti, biz esa Postgres container'ini kill minus toʻqqiz bilan oʻldiramiz. Ikkala jarayon ham xato oldi. Bazani qayta koʻtarsak, logda: baza toʻgʻri yopilmagan, avtomatik tiklash boshlandi, keyin redo, ya'ni qayta bajarish. Har bir tasdiqlangan buyurtmaning raqamini alohida tekshirdik: ming bir yuz oʻttiz toʻqqiztasi ham joyida. Yigirma marta takrorladik va tasdiqlangan buyurtma bir marta ham yoʻqolmadi. Bazadagi bitta ortiqcha qator esa commit boʻlgan, lekin javobi dasturga yetib bormagan buyurtma.

**wal.** Bunda bazani write-ahead log, qisqasi WAL qutqardi. Postgres oʻzgarishni avval jurnal oxiriga yozadi. Commit'da jurnal diskka tushishini kutadi va shundan keyingina muvaffaqiyat haqida javob qaytaradi. Jadval fayllari esa keyinroq, shoshilmasdan yangilanadi. Crash'dan keyin baza jurnalni oʻqiydi va jadvalga yetib bormagan oʻzgarishlarni qayta bajaradi, logdagi redo aynan shu. Biz faqat jarayonni oʻldirdik. Elektr oʻchishida nima boʻlishini va fsync'ni oltinchi modulda chuqur koʻramiz.

**qachon.** Demak, fayl hech qachon kerak emasmi? Kerak, lekin oʻz joyida. Bitta dastur yozadigan sozlama fayliga JSON, rename va fsync yetarli. Bitta kompyuterdagi ilovaga SQLite mos keladi. Uning sayti SQLite'ni server bazalari bilan emas, oddiy fayl ochish funksiyasi fopen bilan solishtiradi. Lekin SQLite bir paytda faqat bitta yozuvchiga ruxsat beradi. Koʻp yozuvchi va tarmoqdagi koʻp mijoz boʻlsa, Postgres kabi server bazasi kerak. Bu trekda Postgres bilan ishlaymiz.

**homework.** Uy vazifasi, toʻrt darajada. Oson: qulf usulini yuz marta kill qiling, fayl butun qoladimi va nega? Oʻrta: Postgres json usulini Go’da oʻqimasdan, bitta UPDATE bilan tuzating. Qiyin: oʻqish va yozishni bitta transaction'ga oling va qatorni SELECT FOR UPDATE bilan qulflang. Guru: rename usuliga fayl va papka fsync'ini qoʻshing va qancha sekinlashishini oʻlchang. Javoblar GitHub'dagi kurs repository'sida, link pastda.

**outro.** Bugun koʻrdik: ikki jarayon JSON faylni buzadi, kill uni boʻshatib qoʻyadi. rename buzilishni toʻxtatadi, lekin oʻqi, qoʻsh, yoz sikli buyurtmalarni baribir jimgina yoʻqotadi. Postgres'dagi INSERT ikki ming buyurtmani ham, crash'ni ham koʻtardi, chunki oʻzgarish bazaning ichida bajariladi va avval jurnalga yoziladi. Keyingi darsda relyatsion model: jadval nega aslida toʻplam. Unda koʻrishamiz, salomat boʻling.

<!-- /matn -->

## Chuqurroq

**Fayl nega buziladi.** `os.WriteFile` = `open(O_WRONLY|O_CREATE|O_TRUNC)` + `write` + `close`. Go hujjati: "Since
WriteFile requires multiple system calls to complete, a failure mid-operation can leave the file in a partially written
state." Ikki jarayon faylni deyarli bir vaqtda ochsa, ikkalasi ham uni boʻshatadi, keyin ikkalasi ham 0-offsetdan
yozadi. Qisqa JSON uzun JSON ustiga tushsa, uzunining dumi qoladi: `…]}{ "id": 68, …]}` — `encoding/json`:
`invalid character '{' after top-level value`. Dumi har safar ~70 bayt: bu aynan bitta buyurtmaning `MarshalIndent`
dagi uzunligi (bir jarayon boshqasidan bitta buyurtma koʻp yozgan). Bu poyga boʻlgani uchun videodagi son har
ishga tushirishda boshqacha.

**kill -9 va fayl.** SIGKILL'ni ushlab boʻlmaydi (SYS-0.02). Jarayon `O_TRUNC` bilan faylni boʻshatgan, lekin `write` ga
ulgurmagan boʻlsa — fayl 0 bayt (hamma buyurtma yoʻq). `write` yarmida toʻxtasa — yarim fayl. Bu *jarayon* crash'i:
yozilgan maʼlumot OS page cache'da qoladi. Elektr oʻchishi yoki OS crash'i boshqa masala: unda `fsync` qilinmagan
maʼlumot ham yoʻqolishi mumkin (DB-6.01).

**rename.** `rename(2)`: "If newpath already exists, it will be atomically replaced, so that there is no point at which
another process attempting to access newpath will find it missing." Oʻquvchi eski yoki yangi faylni koʻradi. Lekin
rename faqat *fayl almashinuvini* atomik qiladi, *oʻqi → qoʻsh → yoz* siklini emas: ikki jarayon bir xil eski faylni
oʻqib, har biri oʻz versiyasini qoʻyadi — keyingisi oldingisini oʻchiradi (lost update). Elektr oʻchishiga qarshi
vaqtinchalik faylni `fsync`, rename'dan keyin papkani ham `fsync` qilish kerak (Guru vazifasi).

**flock.** `flock(2)`: "flock() places advisory locks only" — qulfni hamma yozuvchi ishlatishi shart. Jarayon oʻlsa,
uning deskriptorlari yopiladi va qulf oʻzi boʻshaydi (`.lock` faylining borligi qulf emas). Qulf + rename bilan fayl
toʻgʻri ishlaydi, lekin har buyurtmada butun fayl qayta yoziladi (2000 buyurtmada ~6 s), qidirish uchun hammasini
oʻqish kerak, bir nechta fayl orasida atomiklik yoʻq — bu bazaning ishi.

**JSON Postgres'da ham yoʻqoladi.** `-usul json` da `SELECT doc` va `UPDATE … SET doc = $1` — ikki alohida buyruq
(har biri oʻz avtomatik transaction'ida). Orasidagi oraliq Go kodida, Postgres uni koʻrmaydi. DDIA 7-bob tilida bu
*read-modify-write* sikli va *lost update*. Tuzatish yoʻllari (uy vazifasi): bazadagi atomik oʻzgartirish (bitta
`UPDATE … jsonb_set`), aniq qulf (`SELECT … FOR UPDATE` transaction ichida). `INSERT` da oʻqish yoʻq — har buyurtma
oʻz qatori.

**kill -9 va Postgres.** `docker compose kill db` konteynerning PID 1 iga (postmaster) SIGKILL yuboradi va konteyner
toʻxtaydi. Qayta ishga tushganda log: `database system was not properly shut down; automatic recovery in progress`,
`redo starts at …`, `redo done at …`. Tekshiruv har bir id boʻyicha: mijoz `INSERT` nil xato qaytargandan keyingina
id'ni faylga yozadi, `baza -tekshir` bu id'larni `SELECT id FROM buyurtmalar` bilan solishtiradi. WAL: "changes to data files … must be written only after those changes have been
logged, that is, after WAL records describing the changes have been flushed to permanent storage" (PostgreSQL 18, 28.3).
`synchronous_commit = on` (standart): commit javobidan oldin WAL diskka flush qilinadi. Ortiqcha qator: commit
bazada boʻldi, lekin javob mijozga yetib bormadi — mijoz "xato" deb oʻylaydi. Har jarayonda bir paytda bitta `INSERT`
boʻlgani uchun ortiqcha 0, 1 yoki 2 boʻlishi mumkin. Qayta yuborsa — dublikat (idempotentlik,
keyingi modullarda). Diqqat: PostgreSQL hujjati serverni SIGKILL bilan toʻxtatmaslikni maslahat beradi (shared memory va
bola jarayonlar); biz crash'ni taqlid qilish uchun ataylab shunday qildik. Toʻgʻri tez toʻxtatish — SIGQUIT (immediate),
u ham keyingi ishga tushishda WAL'dan tiklaydi.

**SQLite qachon.** "SQLite does not compete with client/server databases. SQLite competes with fopen()." Bitta kompyuter,
bitta ilova, fayl formati — SQLite. "it will only allow one writer at any instant in time" — koʻp yozuvchi, tarmoq
orqali koʻp mijoz — Postgres kabi server bazasi.

## Manbalar
- CMU 15-445/645 Fall 2024, Lecture 01 (A. Pavlo): [Flat file strawman — data integrity, implementation, durability](https://15445.courses.cs.cmu.edu/fall2024/slides/01-relationalmodel.pdf)
- M. Kleppmann, *Designing Data-Intensive Applications*, 7-bob "Transactions" — "Preventing Lost Updates" (faqat tuzilma uchun): <https://dataintensive.net/>
- Go: [os.WriteFile (go1.26.8, file.go:927–931)](https://github.com/golang/go/blob/go1.26.8/src/os/file.go#L927-L931), [os.Rename](https://pkg.go.dev/os#Rename)
- Linux: [open(2) O_TRUNC](https://man7.org/linux/man-pages/man2/open.2.html), [rename(2)](https://man7.org/linux/man-pages/man2/rename.2.html), [flock(2)](https://man7.org/linux/man-pages/man2/flock.2.html), [fsync(2)](https://man7.org/linux/man-pages/man2/fsync.2.html)
- PostgreSQL 18: [WAL Introduction](https://www.postgresql.org/docs/18/wal-intro.html), [synchronous_commit](https://www.postgresql.org/docs/18/runtime-config-wal.html#GUC-SYNCHRONOUS-COMMIT), [Server Shutdown](https://www.postgresql.org/docs/18/server-shutdown.html), [Transaction Isolation (Read Committed)](https://www.postgresql.org/docs/18/transaction-iso.html), [SELECT … FOR UPDATE](https://www.postgresql.org/docs/18/sql-select.html#SQL-FOR-UPDATE-SHARE), [jsonb_set](https://www.postgresql.org/docs/18/functions-json.html)
- SQLite: [Appropriate Uses For SQLite](https://www.sqlite.org/whentouse.html), [Atomic Commit In SQLite](https://www.sqlite.org/atomiccommit.html)
- pgx: [pgx v5](https://pkg.go.dev/github.com/jackc/pgx/v5)

## Uy vazifasi
- 🟢 **Oson:** `./qulash.sh qulf 100` — fayl butun qoladimi? Nega? Jarayon qulfni ushlab turganida oʻlsa, ikkinchi jarayon abadiy kutib qoladimi?
- 🟡 **Oʻrta:** Postgres `-usul json` ni Go'da oʻqimasdan, bitta `UPDATE` bilan tuzating (`jsonb_set`, `||`, `jsonb_array_length`).
- 🟠 **Qiyin:** Go kodini oʻzgartirmasdan tuzating: oʻqish va yozishni bitta transaction'ga oling va qatorni `SELECT … FOR UPDATE` bilan qulflang.
- 🔴 **Guru:** `rename` usuliga vaqtinchalik fayl `fsync` va rename'dan keyin papka `fsync` qoʻshing. Qancha sekinlashdi? Nega papkani ham?

Avval oʻzingiz urinib koʻring. Javoblar: `bash homework/javoblar.sh` (`.env` va `docker compose up -d --wait` kerak):
- Oson: `100 marta kill -9 (qulf): fayl butun 100` — `qulf` ham rename bilan yozadi, kill qayerga tushmasin, joyida yo eski,
  yo yangi butun fayl turadi. Qulf `flock` — jarayon oʻlganda yadro uning fayl deskriptorlarini yopadi va qulf oʻzi
  boʻshaydi, shuning uchun hech kim abadiy kutmaydi (oddiy "lock fayli bor-yoʻqligi" sxemasida esa qolib ketardi).
- Oʻrta ([homework/baza-tuzat/main.go](homework/baza-tuzat/main.go), `-usul sql`): `UPDATE ombor SET doc = jsonb_set(doc,
  '{buyurtmalar}', doc->'buyurtmalar' || jsonb_build_object('id', jsonb_array_length(doc->'buyurtmalar') + 1, …))`.
  Postgres yangilanayotgan qatorni qulflaydi, ikkinchi `UPDATE` kutadi va Read Committed'da qatorning *yangi* versiyasini
  qayta oʻqib, uning ustiga qoʻshadi → `json ustunida: 2000 ta`. Bu DDIA'dagi "atomic write operation".
- Qiyin (`-usul forupdate`): `Begin` → `SELECT doc … FOR UPDATE` → Go'da qoʻshish → `UPDATE` → `Commit`. Ikkinchi
  jarayon `FOR UPDATE` da birinchisi commit qilguncha kutadi → 2000. Bu "explicit locking"; DB-5 da isolation darajalari bilan
  chuqurroq.
- Guru ([homework/fsync/main.go](homework/fsync/main.go)): `f.Sync()` (vaqtinchalik fayl maʼlumoti diskda), `Rename`,
  keyin `os.Open(dir).Sync()` (yangi nom — papka yozuvi — ham diskda; usiz elektr oʻchsa rename yoʻqolishi mumkin).
  Bizda ikki jarayon × 1000: 6.1 s → 15.4 s va 7.0 s → 13.6 s (~2–2.5×). Qiziqish uchun: Postgres 2000 ta `INSERT` ni, har commit'da WAL
  flush bilan, 1.9 s da yozdi. Bu teng sharoitdagi taqqoslash emas (bizning faylda har buyurtmada butun fayl qayta
  yoziladi va ikki fsync, Postgres'da esa jurnalga kichik yozuv qoʻshiladi) — sababini oʻlchamadik, bu taxmin.
