# GB-5.02 · Postgres Docker Compose'da va psql bilan ulanish

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #38"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: GB-5.01 (nega baza kerak: ikki jarayon bitta JSON faylni buzadi) · keyingi dars: GB-5.03 (SQL: `CREATE TABLE`,
`INSERT`, `SELECT`, constraint).

**Bu darsdan keyin:** `compose.yaml` bilan Postgres'ni bitta buyruqda koʻtarasiz, `psql` bilan ichiga qarab
(`\dt`, `\d products`, `SELECT`) jadvalni koʻrasiz, `down` va `down -v` farqini va ma'lumot volume'da qolishini
tekshirasiz, `ports`dagi ikki raqamni va band port xatosini oʻqiysiz, `compose.yaml`'ga ikkinchi servis qoʻshasiz.

## Talablar

- Docker Engine va Docker Compose v2 (`docker compose version`; [oʻrnatish](https://docs.docker.com/engine/install/), Linux'da
  foydalanuvchi `docker` guruhida boʻlishi kerak). Video: Docker 29.8.1, Compose 5.4.0.
- Go (`go version`) — faqat `go test` uchun. Video: go1.27.4 linux/amd64, `CGO_ENABLED=0`.
- `psql` kompyuteringizda boʻlishi shart emas: darsdagi hamma buyruq `docker compose exec` orqali konteyner ichidagi
  `psql`'ni ishlatadi. (Kompyuterda `psql` boʻlsa, ulanish URL'i bilan xostdan ham ulanasiz, pastga qarang.)
- Kompyuteringizda `5482` va `5483` portlari boʻsh boʻlishi kerak (band boʻlsa, `compose.yaml`'dagi chap raqamni
  almashtiring; testlar esa shu raqamlarni kutadi, shuning uchun band boʻlsa vazifada ham ularni boʻshating).

## Papka

| Fayl | Nima |
|---|---|
| `compose.yaml` | videoning oxiridagi fayl: `db` servisi (image, environment, ports, volumes). Vazifada shunga `db-test` qoʻshiladi |
| `bozor.sql` | `products` jadvali va uchta mahsulot (SQL'ni GB-5.03'da oʻzimiz yozamiz) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `go.mod` | oʻz moduli (`github.com/abdurahmon27/bozor`, `go 1.27`) |
| `TASK.md` | vazifa |
| `solution/` | javob: `compose.yaml` (`name: gb502-solution`: boshqa darslarning `solution/` papkalari bilan loyiha nomi toʻqnashmasin), xuddi shu `bozor.sql`, test fayli va oʻz `go.mod`'i: `cd solution && go test ./...` |
| `_video/run.sh` | videodagi hamma buyruqni yangi `bozor/` papkasida ishga tushiradi (`./_video/run.sh`); natijasi `_video/chiqishlar.txt` |

## Videodagi buyruqlar va natijalari (Docker 29.8.1, Compose 5.4.0, `postgres:18-alpine` = PostgreSQL 18.7)

Serverni koʻtarish (`docker compose up -d`):

```
[+] up 3/3
 ✔ Network bozor_default Created
 ✔ Volume bozor_pgdata   Created
 ✔ Container bozor-db-1  Started
```

Server tayyormi (`docker compose logs db | tail -n 1`; vaqt har safar boshqacha):

```
db-1  | 2026-11-27 12:21:41.271 UTC [1] LOG:  database system is ready to accept connections
```

Fayldan yuklashda `-T` kerak: busiz Docker konteynerga terminal ulaydi, kiritish esa fayldan keladi va xato beradi
(`docker compose exec db psql -U bozor -d bozor < bozor.sql`):

```
cannot attach stdin to a TTY-enabled container because stdin is not a terminal
```

`psql` bilan ichiga qarash (`exec` — ishlab turgan konteynerda buyruq; `-U` foydalanuvchi, `-d` baza, `-c` bitta buyruq):

```
$ docker compose exec db psql -U bozor -d bozor -c '\dt'
Did not find any tables.
$ docker compose exec -T db psql -U bozor -d bozor < bozor.sql
CREATE TABLE
INSERT 0 3
$ docker compose exec db psql -U bozor -d bozor -c '\dt'
          List of tables
 Schema |   Name   | Type  | Owner
--------+----------+-------+-------
 public | products | table | bozor
(1 row)
$ docker compose exec db psql -U bozor -d bozor -c '\d products'
              Table "public.products"
 Column |  Type   | Collation | Nullable | Default
--------+---------+-----------+----------+---------
 id     | bigint  |           | not null |
 name   | text    |           | not null |
 price  | bigint  |           | not null |
 stock  | integer |           | not null |
Indexes:
    "products_pkey" PRIMARY KEY, btree (id)
$ docker compose exec db psql -U bozor -d bozor -c 'SELECT name FROM products'
 name
-------
 olma
 non
 oʻrik
(3 rows)
```

`docker compose down` (volume qoladi), keyin `up -d` — `SELECT name FROM products` yana uchta qator beradi.
`docker compose down -v` (volume ham oʻchadi), keyin `up -d` — `\dt`:

```
Did not find any tables.
```

Klassik xato: `ports` chap raqami (kompyuter porti) band (`"127.0.0.1:5432:5432"`; bu kompyuterda 5432 boshqa Postgres'niki):

```
Error response from daemon: … failed to bind host port 127.0.0.1:5432/tcp: address already in use
```

Tuzatish: chap raqamni almashtiring (`"127.0.0.1:5482:5432"`); oʻng raqam (konteyner ichidagi port) oʻzgarmaydi.

Ulanish URL'i (kompyuterda `psql` boʻlsa): `postgres://bozor:bozor@127.0.0.1:5482/bozor`:

```
$ psql postgres://bozor:bozor@127.0.0.1:5482/bozor -c '\conninfo'
You are connected to database "bozor" as user "bozor" on host "127.0.0.1" at port "5482".
```

Go'dan shu URL bilan ulanish GB-5.05'da (pgx).

## Muhim eslatmalar

- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` faqat **boʻsh** volume'da ishlaydi; mavjud bazaga tegmaydi
  ([Docker Official Image: postgres](https://hub.docker.com/_/postgres)). Parolni almashtirish uchun yoki baza ichida
  `ALTER ROLE`, yoki `docker compose down -v` bilan yangidan boshlash kerak.
- Postgres 18 image'ida volume `/var/lib/postgresql`'ga ulanadi (`PGDATA` versiyaga bogʻliq: `/var/lib/postgresql/18/docker`).
  17 va undan eskilarida volume `/var/lib/postgresql/data`'ga ulanishi kerak: image oʻsha yoʻlda anonim volume eʼlon qiladi, shuning
  uchun `/var/lib/postgresql`'ga ulangan nomli volume'da maʼlumot saqlanmaydi (konteyner qayta yaratilganda yoʻqoladi). Manba:
  [postgres image README](https://hub.docker.com/_/postgres), "PGDATA" boʻlimi.
- `compose.yaml`'dagi parol (`bozor`) faqat shu kompyuterdagi oʻquv bazasi uchun: port faqat `127.0.0.1`'ga ochiq.
  Haqiqiy parollar GB-8.01'da (config va secret'lar).

## Vazifa

[TASK.md](TASK.md): `db-test` servisini qoʻshing, `docker compose up -d`, `bozor.sql`'ni yuklang, `go test` oʻtguncha.
Tekshirish:

```bash
docker compose up -d
docker compose exec -T db psql -U bozor -d bozor < bozor.sql
go test ./...
docker compose down -v   # ish tugagach tozalash
```

Boshlangʻich holatda `go test` Oʻzbekcha xabar bilan yiqiladi: `compose.yaml'da db-test servisi yoʻq …`. Test Postgres tayyor
boʻlishini 30 soniyagacha oʻzi kutadi.

`solution/`'ni tekshirish (oʻz `go.mod`'i bor; oldin starter'ning servislarini `docker compose down -v` bilan oʻchiring: portlar bir xil):

```bash
cd solution
docker compose up -d
docker compose exec -T db psql -U bozor -d bozor < bozor.sql
go test ./...
docker compose down -v
```

## Manbalar

- Docker Compose file reference: [services](https://docs.docker.com/reference/compose-file/services/), [volumes](https://docs.docker.com/reference/compose-file/volumes/)
- Docker: [Volumes](https://docs.docker.com/engine/storage/volumes/), [docker compose up](https://docs.docker.com/reference/cli/docker/compose/up/), [down](https://docs.docker.com/reference/cli/docker/compose/down/), [exec](https://docs.docker.com/reference/cli/docker/compose/exec/)
- [Docker Official Image: postgres](https://hub.docker.com/_/postgres)
- PostgreSQL 18: [psql](https://www.postgresql.org/docs/18/app-psql.html), [connection URIs](https://www.postgresql.org/docs/18/libpq-connect.html#LIBPQ-CONNSTRING-URIS)
