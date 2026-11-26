# 01 · Uy vazifasi (2026-11-26, dasturchilar uchun versiya)

Avval taxmin qiling, keyin oʻlchang, keyin javobni oching. Hamma oʻlchovlar `go1.26.8 linux/amd64`, Docker 29.8.1.
Bu laptop'da gcc yoʻq, shuning uchun `net` import qiladigan dasturlar `CGO_ENABLED=0` bilan build qilinadi
(GO-X.01'dagi retsept). Papka yoʻli oʻzgarsa, baytlar bir necha bayt farq qiladi.

## 🟢 Oson: `import _ "net/http"` — hajm qancha?
Starter: [`../statustext/`](../statustext/). `fmt.Println(http.StatusText(404))` oʻrniga `fmt.Println("salom")` yozing va
importni `_ "net/http"` qiling (endi `net/http`dan hech narsa chaqirilmaydi). Build'dan **oldin** taxmin qiling:
binary 2.4 MB'ga (Salom, dunyo) tushadimi yoki ~5.1 MB'da qoladimi? Keyin `paketlar.sh` bilan `net/http` ulushini koʻring.

<details><summary>Javob</summary>

```
$ CGO_ENABLED=0 go build -o st main.go && ls -l st
-rwxr-xr-x 1 haywan users 5139103 Nov 26 16:49 st
$ ../paketlar.sh st 300 | grep ' net/http$'
    3729   0.3% net/http
```
Deyarli oʻzgarmadi (5 143 595 → 5 139 103). `StatusText`ning oʻzi arzimas edi; hajmni `net/http` va uning
bogʻliqliklarining `init` kodi va global oʻzgaruvchilari (`DefaultTransport` → `net`, TLS, `crypto/*`) beradi, ular esa
paket import qilinishi bilanoq ishlaydi. `import _` aynan shu: "faqat init uchun". Xulosa: Go'da import qilish
bepul emas, uni `paketlar.sh` bilan oʻlchang.
</details>

## 🟡 Oʻrta: `vaqt` → `scratch`, `main.go`ga tegmasdan
Starter: [`orta/`](orta/) (videodagi xato versiya + `Dockerfile`). `main.go`ni oʻzgartirmasdan shunday build qiling-ki,
`docker run` Toshkent soatini chiqarsin. Maslahat: `go doc time/tzdata`.

<details><summary>Javob</summary>

`go doc time/tzdata`: "This package will be automatically imported if you build with -tags timetzdata."
```
$ go build -tags timetzdata -o vaqt main.go && ls -l vaqt
-rwxr-xr-x 1 haywan users 2884802 Nov 26 16:34 vaqt
$ docker build -q -t vaqt . && docker run --rm vaqt
Toshkent soati: 16
```
Hajm `import _ "time/tzdata"` bilan deyarli bir xil (2 885 090 — farq papka yoʻlida). Qaysi biri yaxshi? `import _`
kodning oʻzida turadi va build buyrugʻini unutib boʻlmaydi; `-tags` esa kutubxona yoki bir nechta binary uchun qulay.
Hujjat tavsiyasi: tzdata'ni kutubxona emas, `main` paketi import qilsin. Uchinchi yoʻl: `FROM gcr.io/distroless/static-debian13`
(ichida tzdata bor; oʻlchangan: tuzatilmagan `vaqt` ham `Toshkent soati: 16`).
</details>

## 🔍 Qiziqish: `-s -w` dan keyin stack trace
Starter: [`qiziq/main.go`](qiziq/main.go) (indeks chegaradan chiqadi → panic). Oddiy va `-ldflags="-s -w"` bilan build
qiling, ikkala panic chiqishini solishtiring. Taxmin: funksiya nomlari va qator raqamlari yoʻqoladimi?

<details><summary>Javob</summary>

```
$ go build -o p main.go && ./p
panic: runtime error: index out of range [5] with length 3

goroutine 1 [running]:
main.boshla(...)
	/home/haywan/Developer/go-backend-uz/lessons/001-why-go/uyga/qiziq/main.go:5
main.main()
	/home/haywan/Developer/go-backend-uz/lessons/001-why-go/uyga/qiziq/main.go:9 +0xa
$ go build -ldflags="-s -w" -o p2 main.go && ./p2
(xuddi shu chiqish)
$ ls -l p p2
-rwxr-xr-x 1 haywan users 1830012 Nov 26 16:13 p
-rwxr-xr-x 1 haywan users 1208482 Nov 26 16:13 p2
```
Hech narsa yoʻqolmadi: stack trace `.gopclntab` jadvalidan chiqadi, uni runtime oʻzi ishlatadi va linker uni
tashlamaydi. `-w` faqat DWARF'ni (debugger, masalan Delve uchun), `-s` simvollar jadvalini (`go tool nm` endi "no
symbols" deydi) olib tashlaydi. Fayl yoʻllarini yashirish uchun esa `-trimpath` kerak.
</details>

## ➕ Qoʻshimcha: HTTPS soʻrov `scratch`da
Starter: [`qoshimcha/`](qoshimcha/) (`http.Get("https://go.dev")` + `FROM scratch`). `CGO_ENABLED=0 go build -o app main.go`,
`docker build`, `docker run` → nima chiqadi va qanday tuzatasiz?

<details><summary>Javob</summary>

```
$ docker run --rm https
xato: Get "https://go.dev": tls: failed to verify certificate: x509: certificate signed by unknown authority
```
Ildiz CA sertifikatlari ham tizimdan oʻqiladi (`crypto/x509/root_linux.go`: `/etc/ssl/certs/ca-certificates.crt`, …),
`scratch`da esa ular yoʻq. Ikki tuzatish (ikkalasi ham oʻlchangan → `200 OK`):
```dockerfile
FROM scratch
COPY --from=alpine:3.24 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY app /app
ENTRYPOINT ["/app"]
```
yoki `FROM gcr.io/distroless/static-debian13` (sertifikatlar va tzdata bilan, ~2 MB).
</details>
