# 16 · Uy vazifasi javoblari

Avval oʻzingiz urinib koʻring. 🙂 Hammasi `go1.26.8` (va `go1.27.4`) da ishga tushirilgan.

## 🟢 Oson
`go list -f '{{.DefaultGODEBUG}}' .` faqat Go toolchain'ining standart qiymatlaridan **farq qiladigan** sozlamalarni
chiqaradi. `go` qatori toolchain bilan bir xil boʻlsa (`go 1.26` + Go 1.26), roʻyxat boʻsh. Bu darsdagi `eski/`
(`go 1.21`) uchun 30 ta:

```
asynctimerchan=1,containermaxprocs=0,cryptocustomrand=1,decoratemappings=0,gotestjsonbuildtext=1,gotypesalias=0,
httpcookiemaxnum=0,httplaxcontentlength=1,httpmuxgo121=1,httpservecontentkeepheaders=1,multipathtcp=0,randseednop=0,
rsa1024min=0,tls10server=1,tls3des=1,tlsmlkem=0,tlsrsakex=1,tlssecpmlkem=0,tlssha1=1,tlsunsafeekm=1,updatemaxprocs=0,
urlmaxqueryparams=0,urlstrictcolons=0,winreadlinkvolume=0,winsymlink=0,x509keypairleaf=0,x509negativeserial=1,
x509rsacrt=0,x509sha256skid=0,x509usepolicies=0
```
Har birining tavsifi: https://go.dev/doc/godebug#history (masalan `randseednop=0`: `math/rand.Seed` yana ishlaydi, Go 1.24;
`rsa1024min=0`: 1024 bitdan qisqa RSA kalitlarga ruxsat, Go 1.24).

## 🟡 Oʻrta
Toʻrt usul, `go 1.26` qatoriga tegmasdan:
1. `go.mod` ga `godebug httpmuxgo121=1` qatori (Go 1.23+ toolchain kerak; faqat asosiy modulniki oʻqiladi).
2. `go.mod` ga `godebug default=go1.21`: til versiyasi 1.26 qoladi, GODEBUG standart qiymatlari esa 1.21 dagidek
   (30 ta sozlama, jumladan `httpmuxgo121=1`). Bu darsda aytilgan, lekin koʻrsatilmagan.
3. **Darsda aytilmagan usul:** `main` paketidagi faylning eng tepasiga, `package` dan oldin: `//go:debug httpmuxgo121=1`
   → [`orta/main.go`](orta/main.go). Go 1.21 dan beri ishlaydi (eski Go'lar uni shunchaki eʼtiborsiz qoldiradi).
4. Ishga tushirishda: `GODEBUG=httpmuxgo121=1 go run .` (bu kompilyatsiya qilingan standart qiymatlardan ham ustun).

```
$ cd homework/orta && go run .
GET /salom → 404
POST /salom → 404
$ go list -f '{{.DefaultGODEBUG}}' .
httpmuxgo121=1
```

## 🟠 Qiyin
[`qiyin/qiyin.sh`](qiyin/qiyin.sh) (`mvs/dokon` ning vaqtinchalik nusxasida):

```
go: downgraded example.com/log v1.3.0 => v1.2.0
go: removed example.com/yetkazish v1.0.0
--- go.mod:
module dokon

go 1.26

require example.com/tolov v1.0.0

require example.com/log v1.2.0 // indirect
--- go build:
main.go:7:2: no required module provides package example.com/yetkazish; to add it:
	go get example.com/yetkazish
```
`yetkazish v1.0.0` log ≥ v1.3.0 ni talab qiladi. Log'ni v1.2.0 ga tushirish uchun `go get` `yetkazish` ning ham log'ni
≤ v1.2.0 talab qiladigan eski versiyasini qidiradi. Bunday versiya yoʻq, shuning uchun `yetkazish` ni olib tashlaydi.
Kod esa uni hali import qiladi va build toʻxtaydi. MVS hech qachon talabdan **past** versiyani tanlamaydi.
(Aksincha, `go get example.com/log@latest` → `upgraded example.com/log v1.3.0 => v1.5.0`, ikkala kutubxona ham v1.5.0 ni oladi.)

## 🔴 Guru
Savol: [`guru/`](guru/) (`go 1.22` + `godebug asynctimerchan=1`) ni Go 1.27 da build boʻladigan qiling va kodni taymerning yangi xulqiga moslang.

```
$ cd homework/guru && go run .
go1.26.8 → eski signal keldi
$ GOTOOLCHAIN=go1.27.4 go run .
go: error loading go.mod:
go.mod:5: removed GODEBUG "asynctimerchan" set to old value "1" (https://go.dev/doc/godebug#go-127)
```
Go 1.27 `asynctimerchan` ni butunlay olib tashladi. Go 1.27 release notes: olib tashlangan sozlama `go.mod` yoki
`//go:debug` da **oxirgi standart qiymat** bilan yozilsa (`asynctimerchan=0`) qabul qilinadi, eski qiymat bilan
(`=1`) esa go buyrugʻi xato beradi. Bu yaxshi: dastur jimgina boshqa xulq bilan ishlashdan koʻra build toʻxtagani afzal.

Tuzatish: `godebug` qatorini oʻchiring yoki `asynctimerchan=0` qiling (Go 1.27.4 buni qabul qiladi, Go 1.26.8 da esa
endi `kanal boʻsh` chiqadi), keyin kodni yangi xulqqa moslang: `Reset`/`Stop` dan keyin kanalni qoʻlda boʻshatish
(`select { case <-t.C: default: }`) endi kerak emas, `len(t.C)` esa doim 0.
`go 1.22` qatorining oʻzi esa (`godebug` qatorisiz) Go 1.27 da xato bermaydi, faqat yangi xulq bilan ishlaydi
(darsdagi `taymer/`).
