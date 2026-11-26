# 16 · go.mod'dagi bitta qator dasturingiz xulqini oʻzgartiradi

Oldingi dars: [15 · defer, panic, recover](../015-defer-panic-recover/) · Keyingi dars: GO-1.17 · Generics va generic metodlar: qachon kerak emas

Bogʻliq darslar: [01 · go build ichida nima bor](../001-why-go/) (`go run`, kompilyator) va
[05 · for va Go 1.22 sikl xatosi](../005-for-switch/) (`go 1.21` modulda `c.uz c.uz c.uz`; guru vazifasidagi `//go:build go1.22`).

Bu papkadagi deyarli har bir kichik papka **alohida modul** (oʻz `go.mod` fayli bor): `go` qatorining taʼsirini shunday
koʻrsatamiz. Ular kursning asosiy modulidan (`github.com/abdurahmon27/go-backend-uz`) chiqarib tashlangan.

| Papka | Nima |
|---|---|
| `yangi/`, `eski/` | bir xil `main.go` (`mux.HandleFunc("GET /salom", …)`), `go 1.26` va `go 1.21`: 200/405 va 404/404 |
| `qoida/` | `ilova` (`go 1.26`) + `kutubxona` (`go 1.21`, `replace` bilan): router 200 (GODEBUG ilovaniki), sikl `c.uz c.uz c.uz` (`-lang` kutubxonaniki) |
| `til/` | `go 1.21` modulda `for i := range 3` → kompilyatsiya xatosi (`-lang was set to go1.21; check go.mod`) |
| `sozlama/` | `go 1.21`, lekin `go.mod` da `godebug httpmuxgo121=0` → yangi router |
| `taymer/` | `go 1.22`: `Reset` dan keyin eski signal (Go 1.26) / kanal boʻsh (Go 1.27, `asynctimerchan` olib tashlangan) |
| `kelajak/` | `go 1.27.4`: Go 1.26 uni oʻzi yuklab oladi (`GOTOOLCHAIN=auto`) yoki rad etadi (`GOTOOLCHAIN=local`) |
| `mvs/manba/` | kichik modullar: `example.com/log` v1.2.0, v1.3.0, v1.5.0; `tolov` (log ≥ v1.2.0); `yetkazish` (log ≥ v1.3.0); `ramka` (`go 1.22`) |
| `mvs/qur/` | `manba/` dan `proxy/` ni quradi (GOPROXY fayl protokoli, `archive/zip`) |
| `mvs/proxy/` | `file://` proxy (qurilgan, repoda bor) |
| `mvs/dokon/` | `tolov` + `yetkazish` ni ishlatadigan modul: MVS `log v1.3.0` ni tanlaydi |
| `mvs/zaharli/` | `example.com/log v1.3.0` ning **almashtirilgan** kodi (versiya raqami oʻsha) |
| `mvs/hujum.sh` | `dokon` ni toza module cache bilan `zaharli/` dan build qiladi → `SECURITY ERROR` |
| `mvs/yuksal.sh` | `eski/` nusxasiga `ramka` (go 1.22) ni qoʻshib ishlatadi → `go get` `go` qatorini koʻtaradi → 404 dan 200 ga (hook) |
| `ichki/` | `ombor/internal/narx`: `ombor/hisob` import qila oladi, `kassa` — yoʻq |
| `asbob/` | `tool golang.org/x/tools/cmd/stringer` + `//go:generate go tool stringer -type=Status` (03-darsdagi `Status`) |
| `homework/` | uy vazifasi javoblari |

## Ishga tushirishdan oldin
`mvs/dokon` dagi buyruqlar `example.com/...` modullarini talab qiladi. Ular internetda yoʻq, shu papkadagi proxy'da bor.
Yangi klon qilingan repoda bir marta:

```bash
cd lessons/016-go-mod
export GOPROXY=file://$PWD/mvs/proxy,https://proxy.golang.org GONOSUMDB=example.com CGO_ENABLED=0
(cd mvs/dokon && go mod download)
```
`GONOSUMDB=example.com`: bu modullar ochiq emas, ularni `sum.golang.org` da tekshirib boʻlmaydi.
`kelajak/` (`go 1.27.4` qatori) va `taymer/` dagi `GOTOOLCHAIN=go1.27.4` birinchi marta Go 1.27.4 ni yuklaydi (arxiv qariyb 76 MB, ochilgandan keyin 228 MB, `$(go env GOMODCACHE)` ichida).

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/016-go-mod`) ishga tushiriladi. Natijalar `go1.26.8 linux/amd64` (va `go1.27.4`) da.
```bash
mvs/yuksal.sh                            # 404 → go: upgraded go 1.21 => 1.22 → 200 (hook)
cd eski && go build -n . 2>&1 | grep -o -- '-lang=go[0-9.]*'   # -lang=go1.21
cd til && go run .                       # cannot range over 3 … requires go1.22 or later (-lang was set to go1.21; check go.mod)
cd eski && go list -f '{{.DefaultGODEBUG}}' . | tr , '\n' | wc -l                  # 30
cd eski && go list -f '{{.DefaultGODEBUG}}' . | tr , '\n' | grep -E 'mux|tls10|maxprocs'   # containermaxprocs=0, httpmuxgo121=1, tls10server=1, updatemaxprocs=0
cd qoida/ilova && go build -a -n . 2>&1 | grep -oE -- '-p (main|example.com/kutubxona) -lang=go[0-9.]*'   # kutubxona go1.21, main go1.26
cd qoida/ilova && go run .               # router 200, kutubxona sikl c.uz c.uz c.uz, ilova sikl a.uz b.uz c.uz
cat sozlama/go.mod
cd sozlama && go run .                   # 200 / 405
cd sozlama && GODEBUG=httpmuxgo121=1 go run .   # 404 / 404: muhit oʻzgaruvchisi ustun
sed -n 10,12p taymer/main.go
cd taymer && go run .                    # go1.26.8 → eski signal keldi
cd taymer && GOTOOLCHAIN=go1.27.4 go run .   # go1.27.4 → kanal boʻsh
cd homework/guru && GOTOOLCHAIN=go1.27.4 go run . 2>&1 | tail -1   # removed GODEBUG "asynctimerchan" set to old value "1"
grep '^go ' kelajak/go.mod               # go 1.27.4
cd kelajak && go run .                   # (birinchi marta) go: downloading go1.27.4 (linux/amd64)
cd $(go env GOMODCACHE) && du -sh --si --apparent-size cache/download/golang.org/toolchain/@v/*go1.27.4*.zip golang.org/toolchain@*go1.27.4*   # 76M, 228M
cd kelajak && GOTOOLCHAIN=local go run .  # go: go.mod requires go >= 1.27.4 (running go 1.26.8; GOTOOLCHAIN=local)
cd mvs/dokon && GOPROXY=file://$PWD/../proxy go mod graph | grep -v ^dokon | grep ' example.com/log@'
cd mvs/dokon && GOPROXY=file://$PWD/../proxy go run .   # ikkalasida ham v1.3.0
cat mvs/proxy/example.com/log/@v/list    # v1.2.0 v1.3.0 v1.5.0
grep log mvs/dokon/go.sum
mvs/hujum.sh 2>&1 | grep -E 'mismatch|downloaded|go.sum:|SECURITY'
cd ichki && find . -name '*.go' | sort
cd ichki && go build ./ombor/... && echo ombor: ok
cd ichki && go build ./kassa              # use of internal package ichki/ombor/internal/narx not allowed
cd asbob && grep -n stringer main.go go.mod
cd asbob && go generate && go run .      # Unknown New Paid Shipped
```

## Dars matni

**Hook.** `mvs/yuksal.sh`: `eski/` (`go 1.21`) nusxasiga `go 1.22` deb yozilgan `example.com/ramka` kutubxonasini
qoʻshib ishlatamiz (router kodiga tegmaymiz, Go ham oʻsha 1.26.8). Avval `GET /salom` → 404, `POST /salom` → 404;
`go get` → `go: upgraded go 1.21 => 1.22`; keyin `GET` → 200, `POST` → 405. Go 1.22 da `ServeMux` naqshlari
kengaydi: `"GET /salom"` endi metod + yoʻl (release notes: *Enhanced routing patterns*). `go 1.21` qatori uchun Go
`httpmuxgo121=1` ni tanlaydi va eski router ishlaydi: `/` bilan boshlanmagan naqshni u host bilan boshlanadigan
naqsh deb oʻqiydi (`net/http/servemux121.go`: `if pattern[0] != '/' { mux.hosts = true }`), `"GET "` esa hech
qaysi host'ga mos kelmaydi → 404. Bu safar oʻzgarish yaxshi tomonga, lekin xuddi shu koʻtarilish `tls10server` ni
ham oʻchiradi: standart `tls.Config` bilan TLS 1.0/1.1 mijozlari serveringizga ulana olmay qoladi.

**go.mod.** Modul — birga versiyalanadigan paketlar toʻplami. `module` — modul nomi va import yoʻli; `go` — pastda;
`tool` — loyiha ishlatadigan dastur (Go 1.24+); `require` — kerakli modullar va ularning **eng kichik** versiyalari;
`go 1.17` dan beri `go.mod` bilvosita kerak boʻlgan modullarni ham yozadi (module graph pruning); `// indirect` — asosiy
modulning hech bir paketi bu modulni toʻgʻridan-toʻgʻri import qilmaydi (Go Modules Reference). Katta versiya
(`v2+`) — boshqa import yoʻli (`example.com/log/v2`), ya'ni MVS uchun boshqa modul.

**`go` qatorining uchta vazifasi.**
1. *Til versiyasi.* `go` buyrugʻi kompilyatorga `-lang=go1.21` beradi (`cmd/go/internal/work/gc.go`:
   `p.Module.GoVersion`). Bu **har bir modul uchun alohida**: kutubxona oʻz `go` qatori bilan kompilyatsiya qilinadi.
   Sikl oʻzgaruvchisi (05-dars) va `range 3` shunga bogʻliq. Bitta fayl uchun `//go:build go1.22` versiyani oʻzgartiradi (05-dars guru).
2. *GODEBUG standart qiymatlari.* Go eski kodni buzishi mumkin boʻlgan oʻzgarish kiritsa, eski xulqni qaytaradigan
   sozlama qoʻshadi. Yangi toolchain eski `go` qatorini koʻrsa, **shunday sozlamalarda** oʻsha versiyaning xulqini
   tanlaydi (boshqa tuzatishlar va runtime yaxshilanishlari baribir ishlaydi). Bu yerda **faqat siz build qilayotgan asosiy
   (main) modulning** `go` qatori hisobga olinadi (`cmd/go/internal/load/godebug.go`: `MainModules.GoVersion`);
   workspace (`go.work`) ishlatilsa — `go.work` dagi `go` qatori va `godebug` lari. `go list -f '{{.DefaultGODEBUG}}'`
   Go 1.26.8 da: `go 1.21` → 30 ta, `go 1.22` → 25, `go 1.23` → 17, `go 1.24` → 8, `go 1.25` → 3, `go 1.26` → 0.
   Roʻyxatda `httpmuxgo121=1`, `tls10server=1` (server tomonda TLS 1.0 ga ruxsat), `tlsrsakex=1`, `tls3des=1`,
   `containermaxprocs=0` va `updatemaxprocs=0` (Go 1.25 dagi konteyner CPU limitiga moslashadigan GOMAXPROCS oʻchiq) bor.
3. *Minimal toolchain.* Go 1.21 dan beri `go` qatori majburiy talab (go.dev/doc/toolchain).

**Ikki vazifa — ikki xil `go` qatori** (`qoida/`). `ilova` (`go 1.26`) `kutubxona` (`go 1.21`) ni ishlatadi:
`go build -a -n` → `-p example.com/kutubxona -lang=go1.21`, `-p main -lang=go1.26`. Kutubxonadagi `"GET /salom"` → 200
(GODEBUG ilovaniki), kutubxonadagi sikl → `c.uz c.uz c.uz` (til versiyasi kutubxonaniki), ilovadagi xuddi shu sikl →
`a.uz b.uz c.uz`. Ya'ni Go 1.22 dagi sikl tuzatishi eski `go` qatorli kutubxonalarga yetib bormaydi.

**GODEBUG ni oʻzgartirish** (go.dev/doc/godebug): `go.mod` dagi `godebug` qatori (Go **1.23+** toolchain kerak: Go 1.22
`go.mod:5: unknown directive: godebug` deydi; faqat asosiy modulniki oʻqiladi), `godebug default=go1.21` (til versiyasi
yangi, GODEBUG standart qiymatlari 1.21 dagidek), `main` paketidagi `//go:debug` izohi (Go 1.21+), dastur ishga tushganda
`GODEBUG` muhit oʻzgaruvchisi. Tartib: `go` qatori → `go.mod` `godebug` → `//go:debug` → `GODEBUG` muhit oʻzgaruvchisi (eng ustun).

**`go get` `go` qatorini koʻtaradi.** Go 1.21 dan beri modulning `go` qatori uning bogʻliqliklarinikidan past boʻla
olmaydi (hook). Xuddi shu narsa `asbob/` da ham boʻldi: `go get -tool golang.org/x/tools/cmd/stringer` uchta modul
qoʻshdi va `go 1.26` ni `go 1.26.0` ga oʻzgartirdi (x/tools v0.53.0 `go 1.26.0` talab qiladi). Shuning uchun kutubxona
muallifi `go` qatorini imkon qadar past tutishi kerak: u kutubxonadan foydalanuvchilarning hammasini koʻtaradi.
(`go mod init` Go 1.26.8 da `go 1.26.8` yozadi — kutubxona uchun buni pasaytirish kerak boʻlishi mumkin.)

**Sozlamalar abadiy emas.** GODEBUG sozlamalari kamida ikki yil (toʻrtta reliz) saqlanadi. Go 1.27 `asynctimerchan`
(1.23 da qoʻshilgan — roppa-rosa toʻrt reliz), `tls10server`, `tlsrsakex`, `tls3des`, `tlsunsafeekm`, `x509keypairleaf`,
`gotypesalias` ni olib tashladi. `taymer/` (`go 1.22`): Go 1.23 gacha taymer kanalida bitta joy (bufer) bor edi va
ishlab boʻlgan taymer signali `Reset` dan keyin ham unda qolardi. Go 1.26.8 → `eski signal keldi`, Go 1.27.4 →
`kanal boʻsh`. `go` qatorining oʻzi jim oʻzgaradi, `godebug asynctimerchan=1` deb aniq yozilgan boʻlsa esa Go 1.27
build'ni toʻxtatadi (`homework/guru`). Go 1.27.4 da `eski/` uchun `DefaultGODEBUG` 26 ta (7 tasi olib tashlangan, 3 ta yangi).

**Toolchain.** `go 1.27.4` modulda Go 1.26.8 `GOTOOLCHAIN=auto` (standart) bilan `golang.org/toolchain@v0.0.1-go1.27.4.linux-amd64`
modulini yuklab oladi (biz oʻlchadik: arxiv 75 846 323 bayt ≈ 76 MB, ochilgandan keyin 228 MB (`du --si --apparent-size`;
diskda bloklar bilan 266 MB), toza kesh bilan 76 s) va checksum'ni tekshiradi. `GOTOOLCHAIN=local` →
`go: go.mod requires go >= 1.27.4 (running go 1.26.8; GOTOOLCHAIN=local)`. `toolchain` qatori — tavsiya:
`GOTOOLCHAIN=auto` da standart Go undan eski boʻlsa, oʻsha versiyaga oʻtadi (`go 1.26` + `toolchain go1.27.4` →
`go1.27.4`), `GOTOOLCHAIN=local` da esa eʼtiborsiz qoladi (→ `go1.26.8`). Faqat asosiy modulning `toolchain` qatori oʻqiladi.

**MVS.** `tolov` log ≥ v1.2.0, `yetkazish` log ≥ v1.3.0, proxy'da v1.5.0 ham bor. MVS har modul uchun talab qilingan
versiyalarning eng kattasini oladi — bu hamma talabni qondiradigan eng kichik versiya: **v1.3.0**. Yangi versiya
chiqishi build'ni oʻzgartirmaydi, shuning uchun lock fayl kerak emas (Go Modules Reference, Russ Cox *Minimal Version
Selection*, 2018). Versiyalar faqat `go get` qilganda yoki yangi import qoʻshilganda (`go mod tidy` uni `@latest` ga
bogʻlaydi) oʻzgaradi.

**go.sum.** Har qatorda: modul, versiya (yoki `versiya/go.mod`), hash. Odatda har modul uchun ikki qator; faqat graf
uchun kerak boʻlgan modullar uchun faqat `/go.mod` qatori. `h1:` = SHA-256 (base64). go.sum lock fayl emas: versiyani
tanlamaydi, faqat kodni tekshiradi. `mvs/hujum.sh`: `zaharli/` proxy'da `log v1.3.0` kodi boshqacha →
`verifying example.com/log@v1.3.0: checksum mismatch … SECURITY ERROR`. go.sum'da hali yoʻq ochiq modullar
`sum.golang.org` (checksum database) bilan tekshiriladi. go.sum'ni commit qiling (Go wiki: Modules FAQ).

**internal/.** `internal` ichidagi paketni faqat `internal` ning ota papkasi ostidagi kod import qila oladi
(`go help gopath`, *Internal Directories*).

**tool.** Go 1.24 dan beri `go get -tool` `go.mod` ga `tool` qatorini yozadi; `go tool stringer` aynan shu versiyani
ishga tushiradi. Ilgari `//go:build tools` va blank import'li `tools.go` fayli ishlatilardi. `asbob/` 03-darsdagi
`Status` (`Unknown` = zero value) ni ishlatadi: `go generate` dan keyin `Unknown New Paid Shipped`.

## Manbalar
- Go Modules Reference: https://go.dev/ref/mod (go.mod, MVS, go.sum, checksum database, `// indirect`)
- Go Toolchains: https://go.dev/doc/toolchain
- Go, Backwards Compatibility, and GODEBUG: https://go.dev/doc/godebug (va `$(go env GOROOT)/doc/godebug.md`)
- Release notes: https://go.dev/doc/go1.22 (routing, TLS), https://go.dev/doc/go1.23 (timers), https://go.dev/doc/go1.24 (tool), https://go.dev/doc/go1.27 (olib tashlangan GODEBUG'lar)
- Russ Cox, Minimal Version Selection (2018): https://research.swtch.com/vgo-mvs
- Russ Cox, Go blog: Backward Compatibility, Go 1.21, and Go 2: https://go.dev/blog/compat · Forward Compatibility and Toolchain Management in Go 1.21: https://go.dev/blog/toolchain
- Manba kod (go1.26.8): `src/cmd/go/internal/work/gc.go`, `src/cmd/go/internal/load/godebug.go`, `src/internal/godebugs/table.go`, `src/net/http/servemux121.go`

## Uy vazifasi
- 🟢 **Oson:** oʻz loyihangizda `go list -f '{{.DefaultGODEBUG}}' .` ni ishga tushiring. Boʻsh chiqsa — nega?
- 🟡 **Oʻrta:** `yangi/` ni `go` qatoriga tegmasdan eski routerga oʻtkazing; darsda aytilmagan usulni ham toping.
- 🟠 **Qiyin:** `mvs/dokon` da `go get example.com/log@v1.2.0`: `yetkazish` bilan nima boʻladi va nega?
- 🔴 **Guru:** `homework/guru` (`go 1.22` + `godebug asynctimerchan=1`) ni Go 1.27 da build boʻladigan qiling va kodni taymerning yangi xulqiga moslang.

Avval oʻzingiz urinib koʻring. Javoblar: [homework/](homework/).
