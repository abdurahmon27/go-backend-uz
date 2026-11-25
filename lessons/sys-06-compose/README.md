# SYS-0.06 · docker compose: Postgres va Go ilovasi bitta buyruq bilan

Oldingi dars: [SYS-0.05 · Docker: «mening kompyuterimda ishlaydi» muammosiga yechim](../sys-05-docker/) · Keyingi: SYS-1.01 ·
Nanosekunddan millisekundgacha. Shu muhit keyingi maʼlumotlar bazasi (DB) va backend lab'larida ishlatiladi.

| Fayl | Nima |
|---|---|
| `compose.yaml` | ikki servis: `db` (`postgres:18-alpine`, healthcheck, `pgdata` volume) va `app` (Go server, shu papkadagi Dockerfile'dan) |
| `.env.example` | nusxa oling: `cp .env.example .env` va `DB_PASSWORD` ni yozing (boʻsh qolsa compose `:?` xatosi bilan toʻxtaydi). `.env*` repoga tushmaydi, `.env.example` istisno (`.gitignore`, GIT-1.02) |
| `server/` | SYS-0.05 dagi mehmon serveri, endi hisob Postgres jadvalida (`pgx` v5): `DATABASE_URL` va `PORT` muhitdan, ulanishni `Ping` bilan tekshiradi, SIGTERM'da toza toʻxtaydi (`kut.go`) |
| `Dockerfile` | SYS-0.05 dagi multi-stage → `scratch`, endi `go mod download` va `USER 65532:65532` bilan (image 15.1 MB) |
| `tuzoq.yaml`, `tuzoq.sh` | xuddi shu loyiha, lekin healthcheck yoʻq va `depends_on: [db]`. `./tuzoq.sh` toza volume bilan koʻtaradi (server `connection refused` bilan yiqiladi), `./tuzoq.sh bor` — tayyor baza bilan (odatda ishlaydi) |
| `qolda.sh` | compose'siz: xuddi shu ikki container'ni oddiy `docker` buyruqlari bilan koʻtarish (solishtirish uchun) |
| `homework/` | uy vazifasi javoblari: `bash homework/javoblar.sh` |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/sys-06-compose`). Avval `cp .env.example .env` va `.env` ga parol yozing. 8077 va 8078
portlar boʻsh boʻlsin. Birinchi `docker compose up` app image'ini yigʻadi (`golang:1.26.8` va `postgres:18-alpine` yuklanadi).
**Go kodini oʻzgartirsangiz — `docker compose up -d --build`**: `up` mavjud image'ni qayta yigʻmaydi (biz tekshirdik:
salomlashuv matnini oʻzgartirdik, `up` dan keyin eski matn, `up --build` dan keyin yangisi).
```bash
docker compose up -d 2>&1 | grep -E 'Created|Healthy|Started' | uniq    # tarmoq, volume, db Healthy, keyin app
sleep 1; curl -s localhost:8077; curl -s localhost:8077                 # 1-, 2-mehmon (Postgres'da; server ulanib jadval yaratguncha 1 s)
cat .gitignore .env.example
docker compose config | grep DATABASE_URL                               # ${...} .env dan qoʻyilgan (parol ochiq: hech qayerga joylamang)
docker compose --env-file /dev/null up -d 2>&1 | sed -n '1p;$p'         # ogohlantirish + ${DB_PASSWORD:?} xatosi
docker compose exec db getent hosts db app                              # servis nomi = DNS nomi
docker compose ps --format 'table {{.Service}}\t{{.Status}}\t{{.Ports}}' # app: 127.0.0.1:8077->8077, db: faqat 5432/tcp
grep -A1 '  depends_on' tuzoq.yaml
./tuzoq.sh                                                              # toza volume: app Exited (1), connection refused
./tuzoq.sh bor                                                          # tayyor baza: app Up
grep -A4 healthcheck compose.yaml; grep -A2 depends_on compose.yaml
docker compose exec db psql -U mehmon -d mehmonxona -c 'SELECT id, vaqt::time(0) AS vaqt FROM mehmonlar'
docker compose down 2>&1 | grep Removed | uniq                          # volume qoladi
docker compose up -d >/dev/null 2>&1; sleep 1; curl -s localhost:8077   # 3-mehmon
DB_PASSWORD=yangi-parol docker compose up -d >/dev/null 2>&1; sleep 2; docker compose logs app --no-log-prefix | tail -1   # password authentication failed
docker compose down -v 2>&1 | grep Removed | uniq; docker compose up -d >/dev/null 2>&1; sleep 1; curl -s localhost:8077   # volume oʻchdi: 1
docker run --rm -e POSTGRES_PASSWORD=x -v eski:/var/lib/postgresql/data postgres:18-alpine 2>&1 | grep -E -A1 '^Error|Counter|suggested'; docker volume rm eski >/dev/null
# tozalash
docker compose down -v; docker compose -f tuzoq.yaml down -v; docker image rm mehmonxona-app tuzoq-app
```
Natijalar Docker Engine 29.8.1, Docker Compose 5.4.0, `postgres:18-alpine` (PostgreSQL 18.7), pgx v5.13.0, Go 1.26.8,
Linux 6.18 (NixOS), 2026-11-25. Hamma buyruqlar toza holatdan ikki marta ishga tushirildi: vaqt qiymatlaridan
(`time=`, `Up N seconds`, psql'dagi `vaqt`) tashqari natijalar bir xil. `tuzoq.yaml` toza volume bilan har safar
(8 ta urinishning 8 tasida) yiqildi; volume'da baza allaqachon boʻlganda 6 tasining 6 tasida ishladi — bu poyga, kafolat emas.

## Dars matni (videodagi soʻzlar)

<!-- matn: lesson.yaml dan -->
**hook.** Go server, Postgres, tarmoq va volume bitta buyruq bilan koʻtariladi, lekin bitta qator yetishmasa, server toza kompyuterda yiqiladi. docker compose up avval tarmoq va volume yaratdi. Keyin Postgres sogʻlom boʻlishini kutdi va shundan keyingina serverni ishga tushirdi. Server bir lahzada bazaga ulanib, jadval yaratadi. Endi mehmonlar faylga emas, Postgres jadvaliga yoziladi. Bugun compose faylini qatorma-qator oʻqiymiz va oʻsha qator nega kerakligini oʻz koʻzimiz bilan koʻramiz.

**qolda.** Compose boʻlmasa, xuddi shu ishni oddiy docker buyruqlari bilan qilardik. Repodagi qoʻlda skripti aynan shunday. Avval tarmoq va volume yaratiladi. Keyin Postgres container'i uchta muhit oʻzgaruvchisi va volume bilan ishga tushadi. Undan keyin baza tayyor boʻlguncha kutadigan sikl kerak. Oxirida server yigʻiladi va oʻsha tarmoqda, port bilan ishga tushadi. Hamma flagʻni va tartibni esda tutish kerak, hammasini tozalash uchun esa yana toʻrtta buyruq bor.

**compose.** Compose esa shu ishni bitta faylga yozadi, uning nomi compose.yaml. Unda buyruqlar emas, natija tasvirlanadi: qanday servislar, tarmoq va volume kerak. Eski qoʻllanmalardagi version qatori endi kerak emas. name loyiha nomi, container va volume nomlari shundan boshlanadi. services ostida ikkita servis bor: db va app. db rasmiy postgres image'idan ishga tushadi, versiyasi oʻn sakkiz, alpine ustida. environment container ichidagi muhit oʻzgaruvchilari: foydalanuvchi, parol va baza nomi. volumes qatori pgdata nomli volume'ni Postgres papkasiga ulaydi. app esa build nuqta, ya'ni shu papkadagi multi-stage Dockerfile'dan yigʻiladi. Go kodini oʻzgartirsangiz, up'ga minus minus build qoʻshing, aks holda eski image ishlayveradi. ports oʻtgan darsdagi minus p, oldidagi loopback manzili portni faqat shu kompyuterga ochadi. Faylning oxirida pgdata volume'i e'lon qilingan.

**env.** Dollar va jingalak qavs ichidagi nomlar oʻzgaruvchilar. Compose ularning qiymatini compose.yaml yonidagi .env faylidan oladi va shu joyga qoʻyadi. Shell'da xuddi shu nomli oʻzgaruvchi boʻlsa, u .env'dan ustun turadi. .env container'ga oʻz-oʻzidan tushmaydi: bizning faylda container'ga faqat environment'da yozilganlar yetadi. Nom ortidagi ikki nuqta va soʻroq belgisi esa qoʻriqchi: qiymat boʻlmasa, compose shu xabar bilan toʻxtaydi.

**env-run.** Parol repoga tushmasligi uchun git darsidagi odat: .gitignore .env bilan boshlanadigan hamma faylni yashiradi, .env.example esa istisno. Unda oʻzgaruvchi nomlari bor, parol boʻsh. docker compose config faylni qiymatlar qoʻyilgan holda koʻrsatadi, unda parol ochiq turadi, uni hech qayerga joylamang. Endi .env boʻlmasa nima boʻlishini sinaymiz: minus minus env-file /dev/null boʻsh fayl beradi. Qoʻriqchisiz oʻzgaruvchi uchun compose faqat ogohlantiradi va boʻsh satr qoʻyadi. Qoʻriqchili qatorda esa xato bilan toʻxtadi va bizning xabarimizni chiqardi.

**dns.** Server bazaga qaysi manzil bilan ulanadi? DATABASE_URL'da xost nomi db, ya'ni servis nomi. Compose loyiha uchun alohida tarmoq yaratadi, unda har bir servisni nomi bilan topish mumkin. db container'i ichida getent ikkala nomni IP manzilga aylantiradi. Bu yerda localhost yozish xato boʻlardi: oʻtgan darsda koʻrganimizdek, container'ning oʻz tarmogʻi bor, ichidagi localhost container'ning oʻzi. docker compose ps yana bir narsani koʻrsatadi. app porti strelka bilan xostga ulangan. Postgres'da strelka yoʻq: uning porti faqat ichki tarmoqda, server unga shu tarmoq orqali ulanadi.

**go.** Go tomonida oʻzgarish kam. Server ulanish manzilini DATABASE_URL muhit oʻzgaruvchisidan oladi: jarayonlar darsidagi qoida, sozlama koddan tashqarida turadi. Postgres bilan pgx kutubxonasi orqali gaplashamiz. pgxpool.New hali ulanmaydi, birinchi ulanish birinchi soʻrovda ochiladi. Shuning uchun Ping bilan ulanishni darhol oʻzimiz tekshiramiz. Baza javob bermasa, server xatoni yozadi va chiqish kodi bir bilan toʻxtaydi. Bu tafsilot birozdan keyin juda kerak boʻladi.

**savol.** Savol: healthcheck'ni olib tashlab, depends_on'ni oddiy roʻyxat qilib qoʻysak nima boʻladi? Compose baribir avval db'ni, keyin app'ni ishga tushiradi. Server bazaga ulana oladimi? Videoni pauza qiling va oʻylab koʻring.

**tuzoq-run.** depends_on bazani kutadi deb oʻylaysizmi? Mana healthcheck'siz compose fayli, depends_on unda oddiy roʻyxat. Toza volume bilan koʻtaramiz va server chiqish kodi bir bilan yiqildi: connection refused. depends_on faqat container ishga tushishini kutadi, baza tayyor boʻlishini emas. Postgres birinchi ishga tushishda bazani yaratadi, bizda bir yarim soniyaga yaqin TCP portda hech kim tinglamadi. Endi xuddi shu faylni tayyor bazali volume bilan koʻtaramiz: baza tezroq tayyor boʻldi va server ulandi. Bu klassik «menda ishlaydi» holati. Yechim healthcheck va condition service_healthy: compose app'ni baza sogʻlom boʻlmaguncha ishga tushirmaydi.

**health.** Endi healthcheck qatorlarini oʻqiymiz. test buyrugʻi pg_isready: u Postgres ulanishni qabul qilyaptimi, shuni tekshiradi. Ishga tushish davrida Docker uni har soniyada, keyin esa har oʻn soniyada ishga tushiradi. Minus h bilan pg_isready TCP orqali soʻraydi. Birinchi ishga tushishda Postgres vaqtincha faqat Unix socket'da tinglaydi, socket orqali soʻrasak, tayyor degan javob erta kelishi mumkin. condition service_healthy app'ni shu tekshiruv oʻtguncha ushlab turadi. Lekin healthcheck faqat ishga tushish tartibini beradi: baza keyinroq qayta ishga tushsa, server oʻzi qayta ulana bilishi kerak.

**volume.** Mehmonlarimiz aslida qayerda yashaydi? psql bilan jadvalga qaraymiz: ikki qator, ikkita mehmon. docker compose down container'lar va tarmoqni oʻchiradi, volume'ga esa tegmaydi. Qayta koʻtarsak, hisob uchdan davom etadi. Endi parolni shell orqali almashtiramiz va server password authentication failed bilan yiqildi. Postgres environment'dagi parolni faqat boʻsh volume'da oʻqiydi, bazada esa eski parol qolgan. Minus v qoʻshilsa, volume ham oʻchadi va baza noldan boshlanadi, mehmon yana birinchi. Ehtiyot boʻling: serverda down minus v bazadagi hamma ma'lumotni oʻchiradi.

**eski.** Yana bir yangi tuzoq, u Postgres oʻn sakkiz bilan keldi. Eski qoʻllanmalarda volume var lib postgresql data papkasiga ulanadi. Oʻn sakkizinchi versiyadan boshlab rasmiy image ma'lumotni versiya nomli ichki papkaga yozadi, bu keyingi versiyaga oʻtishni osonlashtiradi. PGDATA'ni oʻzgartirmagan boʻlsangiz va eski yoʻlga volume ulasangiz, container ishga tushmaydi. Xato aynan shu papkani koʻrsatadi va bitta volume'ni var lib postgresql papkasiga ulashni maslahat beradi. Bizning faylda ham xuddi shunday.

**homework.** Uy vazifasi, toʻrt darajada. Oson: app portini APP_PORT oʻzgaruvchisidan oling, u boʻlmasa sakkiz ming yetmish yetti qolsin. Oʻrta: tuzoq.yaml'ga restart on-failure qoʻshing va server oʻzi tiklanishini kuzating. Qiyin: baza qayta ishga tushsa ham ishlashi uchun Go server bazani oʻzi kutsin, qayta urinish va backoff bilan. Guru: parol muhit oʻzgaruvchisida turibdi, docker inspect uni koʻrsatadi va u logga tushib qolishi mumkin. Uni secret faylga oʻtkazing. Javoblar GitHub'dagi kurs repository'sida, link pastda.

**outro.** Bugun koʻrdik: compose.yaml ikki container, tarmoq va volume'ni bitta faylda tasvirlaydi, parol esa repoga emas, .env'ga yoziladi. depends_on faqat tartibni beradi, tayyorlikni esa healthcheck tekshiradi. Server fayldan Postgres'ga koʻchdi: keyingi ma'lumotlar bazasi va backend lab'lari aynan shu muhitda ishlaydi. Sys trekidagi keyingi dars esa latency jadvali, uni oʻzimiz oʻlchaymiz. Unda koʻrishamiz, salomat boʻling.
<!-- /matn -->

## Chuqurroq

**compose.yaml — buyruqlar emas, natija.** `qolda.sh` da 6 ta buyruq va kutish sikli bor, hammasini tozalash uchun yana 4 ta
(`down` ning oʻzi ikkitasiga teng: container'lar va tarmoq; `down -v` volume'ni ham).
`compose.yaml` esa qanday servislar, volume va tarmoq kerakligini tasvirlaydi; `up` yetishmayotganini yaratadi, `down`
oʻchiradi. Loyiha nomi (`name:`, boʻlmasa papka nomi) hamma nomlarning boshida: `mehmonxona-db-1`,
`mehmonxona_pgdata`, `mehmonxona_default`. `version:` kaliti eskirgan, fayl nomi `compose.yaml` afzal.

**.env va interpolatsiya.** `${DB_USER}` qiymati `compose.yaml` yonidagi `.env` dan (yoki shell muhitidan — u ustun)
olinadi. Oʻzgaruvchi topilmasa compose boʻsh satr qoʻyadi va faqat ogohlantiradi — `POSTGRES_USER` boʻsh boʻlsa,
image jimgina `postgres` foydalanuvchisini oladi. Shuning uchun `db` dagi uchala qiymatda `${…:?xabar}` qoʻriqchisi
bor: qiymat boʻlmasa compose xato bilan toʻxtaydi. Diqqat: `.env`
(interpolatsiya uchun) va servisdagi `env_file:` (container muhiti uchun) — ikki xil narsa. `.env` repoga tushmasin:
bu papkada oʻz `.gitignore` bor (GIT-1.02). `POSTGRES_*` oʻzgaruvchilari faqat **boʻsh** maʼlumot papkasida ishlaydi:
volume'da baza bor boʻlsa, `.env` dagi parolni oʻzgartirish bazadagi parolni oʻzgartirmaydi — server `password
authentication failed` bilan yiqiladi (videoda koʻrsatilgan). Yechim: `down -v` (maʼlumot oʻchadi!) yoki avval
`docker compose exec db psql -U mehmon -d mehmonxona -c "ALTER USER mehmon PASSWORD '…'"`, keyin `.env`.

**Tarmoq.** Compose `<loyiha>_default` tarmogʻini yaratadi, unda servis nomi DNS nomi (`getent hosts db app`). Container
ichidagi `localhost` — container'ning oʻzi (SYS-0.05), shuning uchun `DATABASE_URL` da xost `db`. `db` porti xostga
chiqarilmagan: bizning kompyuterda 127.0.0.1:5432 da xostning oʻz Postgres'i bor, `5432:5432` ziddiyat ham berardi.
Bazani xostdan koʻrish kerak boʻlsa: `docker compose exec db psql …` yoki `ports: ["127.0.0.1:5433:5432"]`.

**depends_on tuzogʻi.** Compose `depends_on` boʻyicha tartibni saqlaydi, lekin "On startup, Compose does not wait until a
container is "ready", only until it's running". Toza volume'da Postgres avval `initdb` qiladi, keyin vaqtinchalik serverni
`listen_addresses=''` bilan (faqat Unix socket) ishga tushirib, sozlab, toʻxtatadi va shundan keyingina TCP'da tinglaydi:
bizda container start → TCP 1.24–1.66 s. Go server esa millisekundlarda ulanishga urinadi va `connection refused` oladi.
Baza tayyor volume'da Postgres ~0.12 s da tinglaydi — server odatda ulguradi, shuning uchun xato faqat toza muhitda
(hamkasbda, CI'da, yangi serverda) chiqadi. Yechim: `healthcheck` + `condition: service_healthy`. Lekin u faqat
**ishga tushish** tartibini beradi: baza ishlash paytida qayta ishga tushsa, healthcheck ilovani qutqarmaydi — ilovaning
oʻzi ham qayta ulana bilishi kerak (Qiyin vazifa). Production orkestratorlarida (masalan Kubernetes) `depends_on` umuman
yoʻq.

**Nega `pg_isready -h 127.0.0.1`.** `-h` siz pg_isready Unix socket orqali soʻraydi. Init paytidagi vaqtinchalik server
socket'da "tayyor" deydi, TCP esa hali yopiq: 0.1 s qadamda kuzatuvda bir marta `sock=0 tcp=2` holati ushlandi (oyna
~0.1–0.2 s, ya'ni har safar emas, lekin boʻlishi mumkin). `-h 127.0.0.1` ilova ishlatadigan yoʻlni (TCP) tekshiradi.
`start_period: 30s` + `start_interval: 1s`: ishga tushish davrida har soniyada tekshiradi (app kutib qolmaydi), keyin
`interval: 10s` (standart 30 s); start davridagi muvaffaqiyatsiz tekshiruvlar `retries` (standart 3) ga hisoblanmaydi.
`start_interval` Docker Engine 25+ ni talab qiladi.

**Volume va Postgres 18.** `down` volume'ga tegmaydi, `down -v` uni oʻchiradi (serverda — bazaning hamma maʼlumoti!).
Postgres 18 dan boshlab rasmiy image `PGDATA=/var/lib/postgresql/18/docker` va `VOLUME /var/lib/postgresql`: volume'ni
`/var/lib/postgresql` ga ulang. Eski qoʻllanmalardagi `/var/lib/postgresql/data` ga ulasangiz, entrypoint
"Error: in 18+, …" bilan chiqadi. 17 va undan eski image'lar uchun eski yoʻl toʻgʻri.

**Go tomoni.** `pgxpool.New` hali ulanmaydi: `MinConns` standart 0, birinchi ulanish birinchi `Acquire`/soʻrovda ochiladi va
ulanish xatosi `New` dan qaytmaydi — shuning uchun `Ping`. Xatolar stderr'ga yoziladi, mijozga esa faqat "bazada xato"
(DB xatosining matni tashqariga chiqmasin). Image `USER 65532` (root emas) bilan ishlaydi. `daftar`
har soʻrovda `INSERT INTO mehmonlar DEFAULT VALUES RETURNING id` qiladi: SYS-0.05 dagi faylni `sync.Mutex` himoya
qilardi, endi bir vaqtdagi soʻrovlarni baza tartibga soladi (nega bu muhim — DB-1.01). `id` — identity ustuni, oʻchirilgan
yoki bekor qilingan qatorlardan keyin raqamlar teshikli boʻlishi mumkin; bu hisoblagich emas, qator raqami.

## Manbalar
- Docker Compose: [Overview](https://docs.docker.com/compose/), [Application model / compose.yaml](https://docs.docker.com/compose/intro/compose-application-model/),
  [Startup order](https://docs.docker.com/compose/how-tos/startup-order/), [Interpolation](https://docs.docker.com/compose/how-tos/environment-variables/variable-interpolation/),
  [Networking](https://docs.docker.com/compose/how-tos/networking/), [docker compose down](https://docs.docker.com/reference/cli/docker/compose/down/),
  [docker compose config](https://docs.docker.com/reference/cli/docker/compose/config/), [Secrets](https://docs.docker.com/compose/how-tos/use-secrets/)
- Docker: [Dockerfile HEALTHCHECK](https://docs.docker.com/reference/dockerfile/#healthcheck), [Bridge network (DNS by name)](https://docs.docker.com/engine/network/drivers/bridge/)
- Postgres: [Docker Hub postgres (POSTGRES_*, _FILE, 18+ PGDATA)](https://hub.docker.com/_/postgres), [docker-library/postgres PR #1259](https://github.com/docker-library/postgres/pull/1259),
  [docker-entrypoint.sh](https://github.com/docker-library/postgres/blob/master/docker-entrypoint.sh), [pg_isready](https://www.postgresql.org/docs/18/app-pg-isready.html)
- Go: [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool), [pool.go: NewWithConfig, Ping (v5.13.0)](https://github.com/jackc/pgx/blob/v5.13.0/pgxpool/pool.go#L228-L233)
- Kubernetes: [Init containers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/)

## Uy vazifasi
- 🟢 **Oson:** app portini `APP_PORT` oʻzgaruvchisidan oling, u boʻlmasa 8077 qolsin.
- 🟡 **Oʻrta:** `tuzoq.yaml` ga `restart: on-failure` qoʻshing: server oʻzi tiklanadimi? Necha marta qayta ishga tushdi?
- 🟠 **Qiyin:** baza qayta ishga tushsa ham (healthcheck'siz ham) ishlasin: Go server bazani oʻzi kutsin — `Ping` ni qayta urinish va backoff bilan yozing.
- 🔴 **Guru:** parol muhit oʻzgaruvchisida: `docker inspect mehmonxona-db-1` uni ochiq koʻrsatadi, bola jarayonlar meros oladi, xato paytida logga tushishi mumkin. Parolni secret faylga oʻtkazing, server ham undan oʻqisin.

Avval oʻzingiz urinib koʻring. Javoblar: `bash homework/javoblar.sh` (8078 va 9000 portlar boʻsh boʻlsin):
- Oson ([homework/oson.yaml](homework/oson.yaml)): `ports: ["127.0.0.1:${APP_PORT:-8077}:8077"]`. `.env` ga `APP_PORT=9000`
  yoki `APP_PORT=9000 docker compose up -d` (shell ustun) → `curl localhost:9000` → `Salom! Siz 1-mehmonsiz`.
  `:-` = "boʻsh yoki yoʻq boʻlsa, standart"; `:?` = "boʻsh yoki yoʻq boʻlsa, xato".
- Oʻrta ([homework/orta.yaml](homework/orta.yaml)): server bir necha marta `connection refused` bilan yiqiladi, Docker uni
  qayta ishga tushiradi (bizda `RestartCount` 4), soʻng `bazaga ulandim`. Ishlaydi, lekin loglar xatoga toʻla va har bir
  yiqilish haqiqiy xatoni ham yashiradi — healthcheck baribir kerak.
- Qiyin ([homework/qiyin/main.go](homework/qiyin/main.go)): `ulan()` `Ping` ni 100 ms dan boshlab har safar ikki baravar
  koʻproq kutib takrorlaydi (5 s dan oshmaydi, jami 30 s): bizda 4 urinish (100, 200, 400, 800 ms), keyin ulandi. Notoʻgʻri
  manzil (`ParseConfig` xatosi) qayta urinishsiz darhol qaytadi.
- Guru ([homework/guru.yaml](homework/guru.yaml), [homework/guru/main.go](homework/guru/main.go)): `secrets: db_parol: environment: DB_PASSWORD`
  qiymatni `.env` dan oladi, lekin container'ga `/run/secrets/db_parol` fayli boʻlib tushadi. `db` da
  `POSTGRES_PASSWORD_FILE=/run/secrets/db_parol` (image `_FILE` qoʻshimchasini qoʻllaydi), app'da `DATABASE_URL` parolsiz,
  Go esa parolni `DB_PASSWORD_FILE` dan oʻqib `cfg.ConnConfig.Password` ga qoʻyadi. `docker inspect` endi faqat
  `*_PASSWORD_FILE` ni koʻrsatadi. Bu tasodifiy sizishdan himoya (Docker hujjati: env orqali parol "unintentional
  information exposure" xavfi); `docker exec` qila oladigan odam faylni baribir oʻqiydi — compose secret'i Swarm'siz
  oddiy ulangan fayl. `docker compose config` ham endi parolni koʻrsatmaydi (`environment: DB_PASSWORD`): sir faqat `.env` faylida va container ichidagi `/run/secrets` da.
