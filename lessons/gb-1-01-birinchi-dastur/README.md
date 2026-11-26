# GB-1.01 · Go: birinchi dastur. Modul, `package main`, `go run` va `go build`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · Golang backend #1

Oldingi dars: yoʻq (kursning birinchi darsi) ·
Keyingi dars: [GB-1.02 · qiymat, oʻzgaruvchi va zero value](../gb-1-02-ozgaruvchi-tur/)

Video: YouTube'dagi "Golang backend" playlist'ida (#1) · Yozma dars: https://www.haywan.uz/lab/youtube/gb/gb-1-01

## Darsda nima boʻladi

Boʻsh papkadan bozor moduli yaratiladi, `main.go` qatorma-qator yoziladi, `go run .` va `go build` bilan ishga
tushiriladi. `go build` yasagan bitta fayl Go topilmaydigan muhitda ham ishlaydi. Keyin eskirgan binary tuzogʻi,
ikki haqiqiy kompilyator xatosi (`fayl:qator:ustun`) va `go fmt ./...`.

Bu papka videoning oxiridagi holat: `main.go` bitta qator chiqaradi. Vazifa: [TASK.md](TASK.md).

| Fayl | Nima |
|---|---|
| `go.mod` | modul: `module github.com/abdurahmon27/bozor` va `go 1.26.8` |
| `main.go` | videodagi dastur (7 qator), vazifa shu yerdan boshlanadi |
| `main_test.go` | vazifani tekshiradi (`go test`); test yozishni GB-1.08'da oʻrganamiz |
| `solution/` | javob: ikki qator chiqaradigan `main.go` + oʻsha test |

## Go'ni oʻrnatish

Rasmiy yoʻriqnoma: https://go.dev/doc/install · yuklab olish: https://go.dev/dl

- **Linux:** `.tar.gz` arxivni yuklab oling, keyin
  `sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz` va
  `export PATH=$PATH:/usr/local/go/bin` qatorini `~/.profile` fayliga qoʻshing. Yangi terminal oynasi `~/.profile`'ni
  har doim qayta oʻqimaydi: `source ~/.profile` qiling yoki tizimdan chiqib qayta kiring.
- **macOS:** `.pkg` faylni oching (Apple chipli kompyuterda `darwin-arm64`, Intel'da `darwin-amd64`), keyin terminalni
  qayta oching.
- **Windows:** `.msi` faylni oching va koʻrsatmalar boʻyicha boring, keyin yangi terminal oching. Windows'da
  `go build` `bozor.exe` faylini yasaydi. Videodagi buyruqlar WSL'da (u yerda Linux) va Git Bash'da bir xil ishlaydi
  (Git Bash'da `./bozor` `bozor.exe`'ni ishga tushiradi); `cmd`/PowerShell'da `.\bozor.exe` deb yozasiz, `ls -lh`
  va `env -i` u yerda yoʻq.

Tekshirish: `go version`. Sizda yangiroq versiya boʻlishi mumkin (2026-yil noyabrida eng yangisi `go1.27.4`):
darsdagi hamma narsa bir xil ishlaydi, faqat `go.mod`'dagi `go` yozuvi sizdagi oʻrnatilgan versiyani koʻrsatadi
(`go mod init` uni oʻzi yozadi). `command not found` chiqsa, PATH yangilanmagan: yuqoridagi `source` qadamini qiling.

**`go` yozuvi eskiroq Go'da.** Bu papkadagi `go.mod` `go 1.26.8` talab qiladi. Sizdagi Go eskiroq boʻlsa, odatiy
`GOTOOLCHAIN=auto` sozlamasida Go kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.26.8`, internet kerak);
Debian/Ubuntu paketidagi Go'da (`GOTOOLCHAIN=local`) esa `go.mod requires go >= 1.26.8` xatosi chiqadi. Ikkala holatda
ham yechim: Go'ni go.dev/dl'dan yangilang yoki vazifani oʻz papkangizda `go mod init` bilan boshlang (TASK.md, 1-qism).

Muharrir: istalgan matn muharriri ishlaydi; VS Code + rasmiy Go kengaytmasi (`golang.go`) faylni saqlashda
`gofmt`'ni oʻzi ishga tushiradi. **Diqqat:** u saqlashda ishlatilmagan importlarni ham oʻchiradi, shuning uchun
`import "os"` xatosini koʻrish uchun faylni formatlamasdan saqlang (buyruqlar palitrasida "Save without Formatting")
yoki boshqa muharrirdan foydalaning.

## Buyruqlar va natijalar

Videodagi hamma buyruq, tartibi bilan, boʻsh papkada qayta ishga tushirildi: `go1.26.8 linux/amd64`,
`CGO_ENABLED=0`, 2026-11-26. Videoda chiqishdagi uzun yoʻl `…/bozor` deb qisqartirilgan.

```console
$ go version
go version go1.26.8 linux/amd64
$ mkdir bozor && cd bozor
$ go mod init github.com/abdurahmon27/bozor
go: creating new go.mod: module github.com/abdurahmon27/bozor
$ cat go.mod
module github.com/abdurahmon27/bozor

go 1.26.8
```

`github.com/abdurahmon27/bozor` oʻrniga oʻz GitHub nomingizni yozing (`github.com/<siz>/bozor`); manzil hozircha
mavjud boʻlishi shart emas. Bu papka videodagi yoʻlni saqlaydi.

`main.go` (7 qator):

```go
package main

import "fmt"

func main() {
	fmt.Println("bozor ishga tushdi")
}
```

```console
$ go run .
bozor ishga tushdi
$ ls
go.mod  main.go
$ go build
$ ls
bozor  go.mod  main.go
$ ls -lh bozor
-rwxr-xr-x 1 haywan users 2.3M Nov 26 23:27 bozor
$ ./bozor
bozor ishga tushdi
$ env -i go version
env: ‘go’: No such file or directory
$ env -i ./bozor
bozor ishga tushdi
```

- `go run .` paketni kompilyatsiya qiladi, natijani Go build cache'iga qoʻyadi (Go 1.24'dan beri; papkasi
  `go env GOCACHE`, bizda `~/.cache/go-build`, doimiy kesh) va darhol ishga tushiradi; loyiha papkasida fayl qolmaydi.
  `go run`'ga paket yoki fayl aytish shart: argumentsiz `go run` `go: no go files listed` deydi. `go run main.go`
  ham ishlaydi, lekin faqat shu bitta faylni oladi: paketda fayllar koʻpaysa, `go run .` kerak.
- `go build` argumentsiz joriy papkadagi paketni quradi va `bozor` faylini yasaydi: nomi paket yoʻlining oxirgi
  qismi. Hajmi taxminan 2,4 million bayt (bizda 2 410 070; yoʻl va build sozlamasiga qarab bir necha bayt farq
  qiladi); ichi boʻsh `func main() {}` ham taxminan 1,8 million bayt, katta qismi Go runtime.
- Fayl hech qanday tizim kutubxonasiga bogʻlanmagan: `ldd ./bozor` → `not a dynamic executable`.
  Qaysi koddan qurilganini `go version -m ./bozor` koʻrsatadi (git repo ichida qurilsa, commit ham yoziladi).
- `env -i` buyruqni boʻsh muhitda (PATH'siz) ishga tushiradi: `go` topilmaydi, `bozor` esa baribir ishlaydi.
  Istisno: PATH boʻlmasa, tizim `/bin` va `/usr/bin`'ni qidiradi; Go `apt install golang-go` bilan `/usr/bin`'ga
  oʻrnatilgan boʻlsa, `env -i go version` ishlaydi. Bu `bozor`'ning Go'ga muhtoj emasligini oʻzgartirmaydi.
- Fayl faqat shu operatsion tizim va protsessorda ishlaydi (`linux/amd64`). Boshqa tizim uchun qurish
  (`GOOS`/`GOARCH`) serverga chiqarish boʻlimida.

Eskirgan binary (`main.go`'dagi matn `bozor yangilandi` qilingandan keyin):

```console
$ go run .
bozor yangilandi
$ ./bozor
bozor ishga tushdi
$ go build && ./bozor
bozor yangilandi
```

Ikki xato (har birini `main.go`'da vaqtincha qiling, keyin qaytaring):

| Oʻzgarish | Xabar | Nima deyapti |
|---|---|---|
| 6-qatorda yopuvchi `"` yoʻq | `./main.go:6:34: newline in string` | 6-qator, 34-ustun: matn yopilmasdan qator tugadi |
| | `./main.go:6:34: syntax error: unexpected newline in argument list; possibly missing comma or )` | birinchi xatoning oqibati: qavs ichidagi argumentlar roʻyxati kutilmaganda tugadi |
| 4-qatorga `import "os"` | `./main.go:4:8: "os" imported and not used` | ishlatilmagan import xato, ogohlantirish emas: dastur qurilmaydi |

Har ikki holatda xabar oldidan `# github.com/abdurahmon27/bozor` qatori ham chiqadi: qaysi paket kompilyatsiya
qilinayotgani (videoda yashirilgan). Ustun **baytlarda** sanaladi: bu qatorlarda har bir belgi (Tab ham) bitta bayt,
lekin `oʻ` kabi harflar bir necha bayt, shuning uchun ularli qatorlarda raqam belgilar sonidan katta chiqadi.
Odatda importlar bitta guruhga yoziladi: `import (` … `"fmt"` … `"os"` … `)`; bu shaklni keyingi darslarda koʻramiz.

`go fmt` (tartibsiz yozilgan `main.go` bilan):

```console
$ go fmt ./...
main.go
$ go fmt ./...
$
```

Birinchi marta oʻzgartirgan faylini yozdi, ikkinchi marta oʻzgartiradigan narsa qolmadi. `./...` = joriy papka va
uning ichidagi hamma papkalardagi paketlar.

## Vazifa va tekshirish

[TASK.md](TASK.md): videodagi yoʻlni oʻz papkangizda takrorlang va dastur ikki qator chiqarsin
(`bozor ishga tushdi`, `versiya: 0.1.0`). Papkangizda:

```console
$ go test
--- FAIL: TestIkkiQator (0.00s)
    main_test.go:53: 2-qator «versiya: 0.1.0» boʻlishi kerak, hozir dastur faqat 1 qator chiqaryapti (TASK.md)
FAIL
exit status 1
FAIL	github.com/abdurahmon27/bozor	0.002s
```

Vazifa bajarilgach `ok` chiqadi. Namunaviy javob `solution/` ichida: `go test ./solution/` uni tekshiradi.

## Manbalar

- Go Modules Reference — modul, modul yoʻli, `go` direktivasi: https://go.dev/ref/mod#modules-overview ,
  https://go.dev/ref/mod#go-mod-file-go
- Go 1.26 release notes (`go mod init` 1.26.0'da `go 1.25.0` yozardi; 1.26.1'dan boshlab yana oʻrnatilgan versiyani
  yozadi, manba: `cmd/go/internal/modload/init.go`, `gover.Local()`): https://go.dev/doc/go1.26
- Go toolchains (`GOTOOLCHAIN`): https://go.dev/doc/toolchain
- Go spec — Package clause, Import declarations, Program execution, Exported identifiers: https://go.dev/ref/spec
- `fmt.Println`: https://pkg.go.dev/fmt#Println
- Go 1.24 release notes — `go run` natijasi build cache'da: https://go.dev/doc/go1.24
- `go help build`, `go help run`, `go help fmt`, `go help packages`: https://pkg.go.dev/cmd/go
- Go FAQ — ishlatilmagan import nega xato: https://go.dev/doc/faq#unused_variables_and_imports ; binary nega katta:
  https://go.dev/doc/faq#Why_is_my_trivial_program_such_a_large_binary
- Ustun = bayt: `cmd/compile/internal/syntax` pozitsiyalari, `go/token.Position`: https://pkg.go.dev/go/token#Position
- Effective Go — Formatting: https://go.dev/doc/effective_go#formatting · gofmt: https://go.dev/blog/gofmt
- Go release history (asosiy reliz va tuzatish versiyalari): https://go.dev/doc/devel/release
