# X.01 · 20 qatorlik HTTP API: Go'da birinchi xizmat va bitta binary

Kurs: **Go: JS/Python/PHP dasturchilari uchun** (GO-X, 1-dars) · Lab: https://haywan.uz/lab/youtube/go#go-x-01 ·
Keyingi dars: [X.02 · Turlar runtime'da ham bor](../0x-02-types/) · Bogʻliq: [SYS-0.05 Docker](../sys-05-docker/)
(multi-stage, scratch, cgo tuzogʻi batafsil) · [01 · go build ichida nima bor](../001-why-go/) (binary nega shuncha katta).

Kim uchun: JS/TS (React, Node), Python yoki PHP'da ishlaydigan dasturchi, Go'ga yangi. Dasturlashni noldan
boshlayotgan boʻlsangiz, avval bepul oʻzbekcha kurslardan birini koʻring (Mohirdev "Go dasturlash tili",
dFarkhod/golang-darslari, UzbekDevs).

Hamma natijalar shu kompyuterda olingan: **go1.26.8 linux/amd64**, Node 24.10.0, npm 11.6.1, Express 5.3.0,
Docker 29.8.1 (2026-11-26). Hammasini qayta oʻlchash: `bash olchov.sh` (natija: `olchov.txt`).
Hajmlar oʻnlik MB'da (docker kabi): 8 659 683 bayt = 8.66 MB.

**Bitta gʻoya:** Go xizmati serverga **bitta fayl** boʻlib chiqadi. `go build` bizning kod, standart kutubxona va
Go runtime'ni bitta binary'ga yigʻadi, shuning uchun serverda Node ham, `node_modules` ham, Go ham kerak emas.
Lekin binary **bitta OS va bitta CPU** uchun yasaladi (`GOOS`/`GOARCH`), va `FROM scratch`'da ishlashi uchun
**cgo oʻchiq** boʻlishi kerak (`CGO_ENABLED=0`).

| Go | JS/TS | Python | PHP |
|---|---|---|---|
| `go mod init api` → `go.mod` (2 qator) | `npm init` → `package.json` | `requirements.txt` / `pyproject.toml` | `composer.json` |
| modul = `go.mod`'li loyiha; paket = bitta papka, bitta import yoʻli | paket = `package.json`'li loyiha | paket = papka | — |
| `go` qatori (eng kichik versiya, standart xulqni ham tanlaydi) | `"engines"` (faqat maslahat) | `requires-python` | `"require": {"php": …}` |
| `go.sum` (xeshlar) | `package-lock.json` | lock fayllar | `composer.lock` |
| standart kutubxona (`net/http`, `encoding/json`); tashqi paket → umumiy modul keshi (`go env GOMODCACHE`) | `node_modules/` (Express: 65 paket) | `venv/` | `vendor/` |
| `go run .` (kompilyatsiya + ishga tushirish) | `tsx index.ts` | `python app.py` | `php -S` / php-fpm |
| `go build -o api .` → mashina kodi + runtime, 8.66 MB | `tsc` → JS (baribir `node` kerak: 129 MB) | — | — |
| serverga: `api` (1 fayl) | `node` + `node_modules` + kod (Node SEA: bitta fayl, lekin ichida butun `node`) | `python3` + `venv` + kod | `php-fpm` + `vendor/` + kod |

## Papkalar

| Fayl / papka | Nima | Natija |
|---|---|---|
| `go.mod`, `main.go` | darsdagi API: `GET /health`, `GET /users/{id}` (Go 1.22 router) | `{"status":"ok"}` · `{"id":"42"}` · `POST /health` → `405` + `Allow: GET, HEAD` |
| `node/` | xuddi shu API Express 5 bilan (`index.js`), `xato.js` (aniqlanmagan nom: server ishlaydi, `/health` → 500), `Dockerfile` (`node:24-slim`) | `POST /health` → `404 Cannot POST /health` · image 235 MB |
| `xato/main.go` | `encoding/json` import qilingan, lekin ishlatilmaydi (ataylab buzilgan) | `./main.go:4:2: "encoding/json" imported and not used` — binary yasalmaydi |
| `Dockerfile` | `FROM scratch` + tayyor binary | image 8.66 MB = binary hajmi |
| `tuzoq/Dockerfile` | tuzoq: binary `golang` image'da cgo yoqiq holda yigʻiladi | `exec /api: no such file or directory` · `ldd`: `libc.so.6`, `/lib64/ld-linux-x86-64.so.2` |
| `tuzoq/Dockerfile.ok` | tuzatish: `CGO_ENABLED=0` | `not a dynamic executable` · `{"status":"ok"}` |
| `chuqur/https/` | Chuqurroq: scratch'dan HTTPS soʻrov; `Dockerfile.ok` sertifikatlar bilan | `x509: certificate signed by unknown authority` → `status: 200 OK` |
| `uyga/` | uy vazifasi: boshlangʻich fayllar (`oson/`, `orta/`) va `javob/` | |
| `olchov.sh`, `olchov.txt` | hamma oʻlchovlar skripti va uning natijasi | |

`xato/` ataylab kompilyatsiya boʻlmaydi, shuning uchun tekshiruv: `go vet . ./uyga/... ./chuqur/...`.
Eslatma: `go build ./xato` lesson papkasidan ishlamaydi (`build output "xato" already exists and is a directory`) —
shuning uchun `cd xato && go build`.

## Videodagi buyruqlar

```bash
go mod init api && cat go.mod             # module api / go 1.26.8
go run .                                  # listening on :3000
curl localhost:3000/users/42              # {"id":"42"}
curl -i -X POST localhost:3000/health     # 405 Method Not Allowed, Allow: GET, HEAD
(cd node && npm install && node index.js) # Express: POST /health → 404, HEAD /health → 200
(cd xato && go build)                     # ./main.go:4:2: "encoding/json" imported and not used
go build -o api . && ls -l api            # 8659683
GOOS=linux GOARCH=arm64 go build -o api-arm . && ./api-arm   # amd64'da: Exec format error

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api .
docker build -t api-go . && (cd node && docker build -t api-node .)
docker images --format 'table {{.Repository}}:{{.Tag}}\t{{.Size}}' --filter reference=api-go --filter reference=api-node
docker run -d --rm --name api-go -p 3000:3000 api-go >/dev/null && sleep 1 && curl -s localhost:3000/health

docker build -t tuzoq -f tuzoq/Dockerfile . && docker run --rm tuzoq   # exec /api: no such file or directory
docker build --target build -t tuzoq-build -f tuzoq/Dockerfile . && docker run --rm tuzoq-build ldd /api
docker build -t tuzoq-ok -f tuzoq/Dockerfile.ok . && docker run -d --rm --name tuzoq-ok -p 3001:3000 tuzoq-ok >/dev/null && sleep 1 && curl -s localhost:3001/health   # CGO_ENABLED=0 → ishlaydi

go version -m api                         # go1.26.8, GOOS, GOARCH, CGO_ENABLED, vcs.revision, vcs.modified
```

Deploy retsepti (bitta qator, CI'da ham):

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api .
scp api server:/srv/api/    # yoki: COPY api /api  →  FROM scratch
```

## Dars matni (qisqa)

**Hook.** Ikki endpoint'li Express API'ni Docker image qildim (node:24-slim): 235 MB. Go'da xuddi shu API 8.66 MB,
image ichida Node ham, Go ham, hatto shell ham yoʻq.

**3 tushuncha.** `go.mod` ↔ `package.json`; `go run` / `go build` ↔ `tsx` / `tsc`; bitta static binary (hech narsaga
bogʻlanmagan fayl) ↔ serverdagi `node` + `node_modules`.

**go.mod.** `go mod init api` = `npm init`. Modul = `go.mod`'li loyiha, paket = bitta papka va bitta import yoʻli.
`module api`: modul nomi, ichki paketlar shu nom bilan import qilinadi (kutubxona boʻlsa odatda `github.com/siz/api`).
`go 1.26.8`: bu modul talab qiladigan eng kichik Go versiyasi, `"engines"`ga oʻxshaydi, lekin Go uni majburiy qiladi
(eskiroq Go build qilmaydi yoki `GOTOOLCHAIN=auto` bilan yangisini yuklaydi) va til hamda standart kutubxona xulqini ham
tanlaydi (masalan, `go 1.21` deb yozilsa, `"GET /health"` pattern'i eski router qoidasida ishlaydi va `/health` 404
beradi — oʻlchandi; GO-1.16 darsida batafsil). `go mod init` mahalliy Go versiyasini yozadi; eski 1.26.x bilan ham
ishlatish uchun uni `go 1.26.0` qilib qoʻyish mumkin. Express'da birinchi ish `npm install` (65 paket); Go'da HTTP
server va JSON standart kutubxonada, `require` qatori yoʻq. Tashqi modulni `go get` aniq versiyasi bilan `go.mod`'ga
yozadi (`go.mod` oʻzi versiyani qatʼiy belgilaydi), `go.sum`'ga esa faqat xeshlarni — `package-lock`'dagi `integrity` kabi;
kodi loyiha ichida emas, umumiy modul keshida (pnpm store kabi).

**Kod.** `package main` ichidagi `func main` — kirish nuqtasi. `http.HandleFunc("GET /health", …)` = `app.get("/health", …)`,
faqat metod pattern ichida (Go 1.22+). Handler joyida yozilgan funksiya (arrow function kabi):
`(w http.ResponseWriter, r *http.Request)` — avval javob, keyin soʻrov, Express'dagi `(req, res)`'dan teskari; tur nomdan
keyin (TS'dagi `r: Request` kabi). `*` — pointer (X.04). `w.Header().Set(...)` + `json.NewEncoder(w).Encode(...)` =
`res.json(...)`; `map[string]string` ≈ `Record<string, string>`. `{id}` = `:id`, `r.PathValue("id")` = `req.params.id`.
`http.ListenAndServe(":3000", nil)` faqat xato boʻlsa qaytadi (masalan port band), `log.Fatal` uni yozib, `os.Exit(1)`
qiladi; `nil` = handler berilmadi (JS'dagi `null` kabi), shunda server standart router `http.DefaultServeMux`'ni oladi (yuqoridagi `HandleFunc`'lar shunga yoziladi).

> Real xizmatda: `mux := http.NewServeMux()` (global router oʻrniga) va timeout'li server:
> `srv := &http.Server{Addr: ":3000", Handler: mux, ReadHeaderTimeout: 5 * time.Second}` — oddiy `ListenAndServe`'da
> timeout yoʻq. Yana: `log.Println("listening")` port band boʻlsa ham chiqadi, chunki u `ListenAndServe`'dan oldin.

**Bashorat 1.** `POST /health` → **405 Method Not Allowed**, `Allow: GET, HEAD` (router yoʻl borligini, metod mos
emasligini biladi; `GET` pattern `HEAD`'ni ham oladi — Express ham). Express 5 notoʻgʻri metodga `404 Cannot POST /health`.

**Oʻrtadagi savol (xato).** JSON qoʻlda yozildi (`w.Write`), `"encoding/json"` import'i qolib ketdi → `go build`:
`./main.go:4:2: "encoding/json" imported and not used`, binary yoʻq. Go kompilyatorida ogohlantirish yoʻq, faqat xato
(Go FAQ: "the Go compiler does not report warnings, only errors that prevent compilation"; "unused imports just slow
down compilation"). Muharrirda gopls (VS Code Go) saqlashda import'larni odatda oʻzi tartiblaydi.
(`node/xato.js` — taqqoslash uchun: Node'da aniqlanmagan nom server ishga tushganidan keyin, soʻrov oʻsha qatorga
yetganda `500` beradi.)

**Binary.** `go build -o api .` → 8 659 683 bayt: bizning kod + `net/http` + `encoding/json` + Go runtime (GC va
scheduler — Node'dagi V8 va event loop oʻrnida). Mashina kodi bitta OS (`GOOS`) va bitta CPU (`GOARCH`) uchun;
berilmasa, host'niki olinadi (`linux`/`amd64`; Node'da `process.arch` → `x64`). Serverga: Node'da `node` (129 MB) +
`node_modules` + kod; Python'da interpretator + venv; PHP'da php-fpm + `vendor/`; Go'da bitta fayl. Node SEA ham bitta
fayl beradi, lekin u `node` binary'sining nusxasi + sizning skriptingiz.
Eslatma: binary "boshqa hech narsasiz" ishlaydi, lekin oddiy fayllar (HTML, SQL migratsiyalar) ichiga oʻzi kirmaydi —
ular uchun `//go:embed`.

**Bashorat 2.** Apple chipli MacBook'da (arm64 Go) `GOOS=linux go build`, keyin amd64 VPS → **exec format error**:
`GOARCH` berilmagan, host'niki (arm64) boʻlib qoladi. Docker'dagi `--platform linux/amd64` ogohlantirishi bilan bir xil
gap. Tuzatish: `GOOS=linux GOARCH=amd64`. (zsh'da xabar: `zsh: exec format error: ./api`.)

**Docker.** `FROM scratch` — boʻm-boʻsh image (shell, libc, sertifikatlar yoʻq), faqat binary: 8.66 MB. `docker exec … sh`
ishlamaydi (`exec: "sh": executable file not found`).

**Tuzoq.** `golang` image'da (yoki gcc oʻrnatilgan har qanday Linux'da, CI'da) `go build` → cgo yoqiq → `net` paketi
cgo resolver bilan libc'ga bogʻlanadi; binary ELF interpreter (dinamik yuklovchi, skriptdagi shebang kabi) sifatida `/lib64/ld-linux-x86-64.so.2`'ni soʻraydi
(`readelf -l`). Scratch'da u yoʻq → `execve` ENOENT → `exec /api: no such file or directory`. Tuzatish:
`CGO_ENABLED=0` (net paketi sof Go resolver'ni oladi). Mac'dan boshqa OS uchun build qilganda cgo oʻzi oʻchadi
(cross-compile), lekin CI'da yoqiq qolishi mumkin — yozib qoʻying. Istisno: C kutubxonasiga bogʻliq paketlar
(masalan `mattn/go-sqlite3`; oʻrniga sof Go `modernc.org/sqlite`). JS'dagi oʻxshashi: glibc uchun yigʻilgan native
addon (masalan `bcrypt`, glibc host'dan koʻchirilgan `node_modules`) Alpine (musl) image'da ishlamaydi.

**Chuqurroq.** `go version -m api` binary ichidagi build maʼlumotini koʻrsatadi: Go versiyasi, `GOOS`/`GOARCH`,
`CGO_ENABLED`, `vcs.revision` (git commit) va `vcs.modified` (commit'dan keyin oʻzgargan fayllar). Bu faqat `go build`'da
yoziladi: `go run` va `.git`'siz Docker build'da revision boʻlmaydi. Scratch'ning ikkinchi tuzogʻi: tashqi HTTPS
(toʻlov tizimi, Telegram API) → `x509: certificate signed by unknown authority`; tuzatish — `ca-certificates.crt`'ni
koʻchirish (`chuqur/https/Dockerfile.ok`), `distroless/static` yoki `import _ "golang.org/x/crypto/x509roots/fallback"`.
Qoʻshimcha: `-ldflags="-s -w"` binary'ni 8 659 683 → 6 004 898 baytgacha kichraytiradi (GO-1.01).

**Xulosa savollari.** 1) `node_modules` yoʻq: kutubxonalar qayerda? → standart kutubxonada yoki modul keshida, `go build`
hammasini binary'ga yigʻadi. 2) MacBook → scratch image: build qatori? → `CGO_ENABLED=0 GOOS=linux GOARCH=amd64`.
3) Ishlatilmagan import? → xato, binary yasalmaydi.

## Uy vazifasi (javoblar: `uyga/javob/`)

- 🟢 **Oson** (`uyga/oson/main.go`): portni `PORT` muhit oʻzgaruvchisidan oʻqing (`os.Getenv`, Node'dagi
  `process.env.PORT`), boʻsh boʻlsa `3000` (kerakli ikki qator boshlangʻich faylda izohda). Keyin `GOOS=linux GOARCH=arm64 go build -o api-arm ./uyga/oson` va
  `go version -m api-arm | grep GOARCH` → `GOARCH=arm64`.
- 🟡 **Oʻrta** (`uyga/orta/main.go`): `GET /version` → `{"go":"go1.26.8","revision":"<commit>"}`.
  `version()` funksiyasi boshlangʻich faylda tayyor (`runtime/debug.ReadBuildInfo()`; `:=`, `range` X.02 darsida);
  siz uni route'ga ulaysiz va `go run .` bilan `go build` natijasini solishtirasiz: `go run`'da `revision` → `unknown`,
  `go build`'da commit (oʻlchandi).
- 🔍 **Qiziqish:** `tuzoq/Dockerfile.ok` asosida multi-stage Dockerfile yozing, oxirgi bosqichga sertifikatlarni ham
  qoʻshing (`chuqur/https/Dockerfile.ok`) va image hajmini `docker images` bilan solishtiring.

## Manbalar

- go.mod: https://go.dev/ref/mod (modules/packages, `module`, `go` directive, `go.sum`, module cache) · GODEBUG: https://go.dev/doc/godebug
- `go run` / `go build`: https://pkg.go.dev/cmd/go
- `net/http.ServeMux` patterns (Go 1.22): https://pkg.go.dev/net/http#ServeMux · 405 + `Allow`:
  https://github.com/golang/go/blob/go1.26.8/src/net/http/server.go#L2699-L2708 · `ListenAndServe` (nil → DefaultServeMux)
- Go FAQ — unused imports, no warnings: https://go.dev/doc/faq#unused_variables_and_imports
- `GOOS`/`GOARCH` (default = host): https://go.dev/doc/install/source#environment
- cgo default qoidalari: https://pkg.go.dev/cmd/cgo · `net` resolver (cgo vs pure Go): https://pkg.go.dev/net#hdr-Name_Resolution
- `execve(2)` ENOENT (ELF interpreter): https://man7.org/linux/man-pages/man2/execve.2.html
- Docker `scratch`: https://docs.docker.com/build/building/base-images/ · distroless: https://github.com/GoogleContainerTools/distroless
- `go version -m`, `runtime/debug.BuildInfo`: https://pkg.go.dev/runtime/debug#BuildInfo
- Node SEA: https://nodejs.org/api/single-executable-applications.html · Node addons: https://nodejs.org/api/addons.html
- Go FAQ (runtime har bir binary'da): https://go.dev/doc/faq#Why_is_my_trivial_program_such_a_large_binary
