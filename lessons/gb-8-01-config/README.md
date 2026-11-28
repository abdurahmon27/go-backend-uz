# GB-8.01 · Config va secret'lar: notoʻgʻri sozlamada aniq xato bilan toʻxtash

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #58"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

8-modulning ("serverga chiqaramiz") birinchi darsi. GB-6.08'dagi himoyalangan bozor API'sining sozlamasi koddan chiqadi:
port flag yoki muhit oʻzgaruvchisidan, Postgres manzili (secret) faqat muhitdan; hamma qiymat **ishga tushishda bir marta**
tekshiriladi va notoʻgʻri boʻlsa server hech narsa ochmay, hamma muammolarni sanab, `1` kodi bilan chiqadi. Startup
log'idagi `%+v` Postgres parolini chiqarardi: `String()` metodi endi faqat host va baza nomini beradi. `.env.example`
gitda, `.env` esa `.gitignore`da.

Oldingi darslar: [`gb-6-08-himoyalangan-bozor`](../gb-6-08-himoyalangan-bozor) (bu papka uning `solution/`idan boshlanadi) ·
[`gb-5-10-buyurtma-api`](../gb-5-10-buyurtma-api) (`main` → `run()` pattern'i) ·
[`gb-5-05-pgx-pool`](../gb-5-05-pgx-pool) (`DATABASE_URL`, `os.Getenv`) · [`gb-2-07-interface-writer`](../gb-2-07-interface-writer)
(`fmt.Stringer`) · [`gb-2-08-xatolar`](../gb-2-08-xatolar) (`errors.Join`) · keyingi dars: `gb-8-02` (`log/slog`).
Farqlarni koʻrish: `diff -ru ../gb-6-08-himoyalangan-bozor/solution .` — yangi `cmd/api/config.go`,
`cmd/api/tekshiruv_test.go`, `.env.example`, `.gitignore`; oʻzgargan `cmd/api/main.go`; GB-6.08'ning `tekshiruv_test.go`si
endi `royxat_test.go` (oddiy test). Laboratoriya Postgres'i boshqa portda: `127.0.0.1:5581` (compose loyihasi `gb801`).

## Sozlamalar

| muhit oʻzgaruvchisi | flag | qiymat | standart | secret? |
|---|---|---|---|---|
| `PORT` | `-port` | 1–65535 oraligʻidagi butun son | `8765` | yoʻq |
| `DATABASE_URL` | — | `postgres://foydalanuvchi:parol@host:port/baza` | — (majburiy) | **ha**: faqat muhitdan |
| `HOST` (vazifa) | `-host` | IP yoki `localhost`; boʻsh — hamma manzillar | `""` | yoʻq |

Ustunlik: **flag > muhit oʻzgaruvchisi > standart**. `PORT=` (qoʻyilgan, lekin boʻsh) — xato, standartga jim oʻtmaydi
(`os.LookupEnv` qoʻyilmagan va boʻshni ajratadi). Xato matnida secret'ning qiymati hech qachon yoʻq, faqat nomi.

Nega `DATABASE_URL` flag emas: jarayonning buyruq qatori (`ps`, `/proc/PID/cmdline`, huquqi `-r--r--r--`) odatda shu
kompyuterdagi hamma foydalanuvchiga koʻrinadi (`/proc` `hidepid` bilan mount qilinmagan boʻlsa); muhit (`/proc/PID/environ`, `-r--------`) faqat egasi va root'ga.
Muhit ham toʻliq himoya emas: OWASP muhit oʻzgaruvchilari log va dump'larga tushishi mumkinligini aytadi; bitta serverda
qiymatlarni faqat root oʻqiydigan fayldan systemd beradi (GB-8.06), katta tizimlarda secret manager ishlatiladi.

```go
// cmd/api/config.go (qisqartirilgan)
func loadConfig(args []string) (config, error) {
	fs := flag.NewFlagSet("api", flag.ContinueOnError)
	port := fs.String("port", envOr("PORT", "8765"), "HTTP port (muhitda: PORT)")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	var errs []error
	// … PORT: strconv.Atoi + 1–65535; DATABASE_URL: boʻsh boʻlmasin
	if len(errs) > 0 {
		return config{}, errors.Join(errs...)
	}
	return config{port: n, databaseURL: db}, nil
}

func (c config) String() string { // log'da: port=8781 db=localhost:5581/bozor
	db := "[yashirin]"
	if u, err := url.Parse(c.databaseURL); err == nil && u.Host != "" {
		db = u.Host + u.Path // foydalanuvchi, parol, ?parametrlar — yoʻq
	}
	return fmt.Sprintf("port=%d db=%s", c.port, db)
}
```

`fmt`ning `%v`, `%+v`, `%s`i `String()`ni chaqiradi; **`%#v` chaqirmaydi** va hamma maydonni (parol bilan) chiqaradi —
config'ni `%#v` bilan yozmang (`_video/chiqishlar.txt` §4). `log/slog` uchun xuddi shu ish `LogValue` bilan qilinadi (GB-8.02).

**`String()`ning chegarasi.** U config'ni *toʻgʻridan-toʻgʻri* chop etganda ishlaydi. `%d` kabi notoʻgʻri verb yoki config
boshqa struct'ning kichik harfli maydonida turgan holda `%+v` — fmt metodni chaqira olmaydi va maydonlarni xom chiqaradi
(fmt kichik harfli maydonlarning metodlarini chaqirmaydi; `type secret string` ham bu holatda yordam bermaydi). Qoida: config'ni
faqat oʻz log qatori orqali yozing, butun struct'larni log'ga tashlamang; GB-8.02'da slog uchun `LogValue`.

**Laboratoriya paroli.** `bozor-lab` ataylab ochiq: `compose.yaml`, `internal/pgtest` va `tekshiruv_test.go`da turadi —
bu faqat shu kompyuterdagi (`127.0.0.1`) laboratoriya bazasi. Haqiqiy serverning paroli hech qachon gitga tushmaydi.

**Hali tekshirilmaydiganlar** (ataylab, dars hajmi uchun): `DATABASE_URL`ning shakli `loadConfig`da emas, keyingi
`pgxpool.New`/`Ping`da tekshiriladi — baribir port ochilishidan oldin; ortiqcha pozitsion argument (`./api 9999`) e'tiborsiz
qoladi (`fs.NArg()` bilan rad etish mumkin). `HOST`ning standarti boʻsh (hamma manzillar) — GB-6.08'dagi xatti-harakat va
GB-8.04'dagi konteyner shuni talab qiladi; serverda `HOST=127.0.0.1` qoʻyiladi.

**`.env` va shell.** `set -a; . ./.env` faylni shell kodi sifatida bajaradi: `;`, `&`, `$`, boʻsh joyli qiymat bittalik
qoʻshtirnoqda boʻlishi shart (`.env.example`dagidek), aks holda qiymat kesiladi yoki buyruq boʻlib bajariladi.

## Videodagi buyruqlar (haqiqiy chiqishlar)

Hammasi shu kompyuterda ishga tushirilgan (`GOTOOLCHAIN=go1.27.4`, `CGO_ENABLED=0`, PostgreSQL 18 `postgres:18-alpine`).
Skript papkaning nusxasida `git init` qiladi (sizning bozor reponingiz kabi). Toʻliq jurnal:
[`_video/chiqishlar.txt`](_video/chiqishlar.txt) (`bash _video/chiqishlar.sh`).

```console
$ go build ./cmd/api && ./api ; echo "kod: $?"
xato: DATABASE_URL qoʻyilmagan: .env.example'ga qarang
kod: 1
$ PORT=80a ./api ; echo "kod: $?"
xato: port "80a": 1–65535 oraligʻidagi butun son kerak
DATABASE_URL qoʻyilmagan: .env.example'ga qarang
kod: 1
$ PORT=9090 ./api -h
Usage of api:
  -port string
    	HTTP port (muhitda: PORT) (default "9090")
$ ./api -h 2> /dev/null ; echo "kod: $?"
kod: 0
$ PORT= ./api ; echo "kod: $?"
xato: port "": 1–65535 oraligʻidagi butun son kerak
DATABASE_URL qoʻyilmagan: .env.example'ga qarang
kod: 1

$ cat .gitignore
.env
/api
$ git check-ignore -v .env
.gitignore:1:.env	.env
$ cp .env.example .env          # muharrirda: PAROL → bozor-lab (qoʻshtirnoqlar qoladi)
$ set -a; . ./.env; set +a      # .env'dagi qatorlarni muhitga eksport qiladi
$ echo $PORT
8765

$ docker compose up -d --wait
 Container gb801-db-1 Healthy
$ ./api -port 8781              # String() yoʻq paytda — SIZISH:
config: {port:8781 databaseURL:postgres://bozor:bozor-lab@localhost:5581/bozor}
bozor API: http://localhost:8781
$ go build ./cmd/api && ./api -port 8781   # String() bilan:
config: port=8781 db=localhost:5581/bozor
bozor API: http://localhost:8781

# ikkinchi terminalda (server ishlab turganda):
$ ps -o args= -C api
./api -port 8781
$ pid=$(pgrep -x api)
$ stat -c '%A %n' /proc/$pid/cmdline /proc/$pid/environ
-r--r--r-- /proc/990337/cmdline
-r-------- /proc/990337/environ
```

Tugagach: serverni `Ctrl+C` bilan toʻxtating, `docker compose down -v`.

## Papka

```
gb-8-01-config/
  cmd/api/config.go           config, envOr, loadConfig, addr, String (GB-8.01)
  cmd/api/main.go             main → run(os.Args[1:]); -h xato emas (kod 0)
  cmd/api/tekshiruv_test.go   vazifa tekshiruvi (Postgres'siz): go test ./cmd/api -run Config
  cmd/api/…                   GB-6.08'dagi API (handler'lar, middleware, testlar — haqiqiy Postgres'da)
  .env.example                gitda: nomlar va namuna (PAROL — haqiqiy emas)
  .gitignore                  .env va /api (binary; sleshsiz cmd/api papkasi ham yashirinardi)
  compose.yaml                Postgres 18, 127.0.0.1:5581
  solution/                   javob (oʻz go.mod'i bilan, alohida ishlaydi)
  _video/                     chiqishlar.sh → chiqishlar.txt
```

## Vazifa

[`TASK.md`](TASK.md): `HOST` sozlamasi (`-host`, IP yoki `localhost`, `net.JoinHostPort`). Tekshirish:
`go test ./cmd/api -run Config` (Postgres kerak emas); butun papka: `docker compose up -d --wait && go test ./...`.

## Manbalar

- [The Twelve-Factor App: III. Config](https://12factor.net/config) — «An app's config is everything that is likely to
  vary between deploys»; «strict separation of config from code»; «whether the codebase could be made open source at any
  moment, without compromising any credentials»; «The twelve-factor app stores config in environment variables».
- [OWASP Secrets Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
  — §8.3 secrets must «Never be logged»; §9.2 «Keys that were exposed should undergo immediate revocation»; §5.1
  «environment variables are generally accessible to all processes and may be included in logs or system dumps».
- [OWASP Logging Cheat Sheet: Data to exclude](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html#data-to-exclude)
  — «Authentication passwords», «Database connection strings», «Encryption keys and other primary secrets».
- Go: [os.LookupEnv](https://pkg.go.dev/os#LookupEnv), [flag](https://pkg.go.dev/flag) (`NewFlagSet`,
  `ContinueOnError`, `ErrHelp`), [errors.Join](https://pkg.go.dev/errors#Join), [net/url.URL](https://pkg.go.dev/net/url#URL),
  [net.JoinHostPort](https://pkg.go.dev/net#JoinHostPort), [net.ParseIP](https://pkg.go.dev/net#ParseIP),
  [fmt: Stringer, %#v](https://pkg.go.dev/fmt#hdr-Printing)
- Linux: [proc_pid_cmdline(5)](https://man7.org/linux/man-pages/man5/proc_pid_cmdline.5.html),
  [proc_pid_environ(5)](https://man7.org/linux/man-pages/man5/proc_pid_environ.5.html); Bash:
  [`set -a`](https://www.gnu.org/software/bash/manual/bash.html#The-Set-Builtin),
  [`VAR=qiymat buyruq`](https://www.gnu.org/software/bash/manual/bash.html#Simple-Command-Expansion)
- [git check-ignore](https://git-scm.com/docs/git-check-ignore), [gitignore](https://git-scm.com/docs/gitignore)
  (boshidagi `/` — faqat shu papka ildizida)
- [RFC 9293 §3.1](https://www.rfc-editor.org/rfc/rfc9293#section-3.1): TCP port — 16 bit
- Hamma chiqishlar shu kompyuterda: go1.27.4 linux/amd64, `CGO_ENABLED=0`, PostgreSQL 18 (`postgres:18-alpine`),
  `_video/chiqishlar.sh`.
