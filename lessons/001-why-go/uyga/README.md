# 01 · Uy vazifasi (2026-11-26 dagi yangi video)

Avval oʻzingiz urinib koʻring, keyin javobni oching. 🙂 Hamma oʻlchovlar `go1.26.8 linux/amd64`da.

## 🟢 Oson: oʻz ismingiz (build qilmasdan oldin taxmin qiling)
Starter: [`oson/main.go`](oson/main.go).

1. Avval `go build main.go`, keyin `./main` (Windows cmd: `main`): `Salom, dunyo!`.
2. `dunyo`ni oʻz ismingizga almashtiring va faylni saqlang. **Build qilmang.**
3. Taxmin qiling: `./main` endi nima chiqaradi? Keyin tekshiring.
4. `go build main.go` → `./main`.

<details><summary>Javob</summary>

```
$ ./main            # 3-qadam: hali eski fayl
Salom, dunyo!
$ go build main.go  # 4-qadam: qayta yasaymiz
$ ./main
Salom, Aziza!
```
Tayyor fayl (`main`) tarjima qilingan paytdagi koddan yasalgan: `main.go` oʻzgarsa ham, u oʻzgarmaydi.
</details>

## 🟡 Oʻrta: `for` bilan uch marta
Starter: [`orta/main.go`](orta/main.go) (Go 0 · 5-darsdagi `for`). Build qilishdan **oldin** taxmin qiling: fayl
`salom/main`dan (2 409 990 bayt) qanchaga katta boʻladi? Bir necha baytgami, minglab baytgami, million baytgami?

<details><summary>Javob</summary>

```
$ go build main.go && ls -l main
-rwxr-xr-x 1 haywan users 2410174 Nov 26 12:25 main
$ ./main
Salom!
Salom!
Salom!
```
Bor-yoʻgʻi **+184 bayt**. Sikl uchun bir necha kompyuter buyrugʻi kerak, xolos; faylning asosiy qismi baribir
runtime va `fmt`. (Aniq son papka yoʻliga qarab bir necha bayt farq qilishi mumkin.)
</details>

## 🔍 Qiziqish: faylni boshqa papkaga koʻchiring
`main` faylini (Windows'da `main.exe`) boshqa papkaga, masalan, Ish stoliga koʻchiring va oʻsha papkada terminal ochib
ishga tushiring: Linux/Mac `./main`, Windows cmd `main`. Ishlaydimi? Nega yonida `main.go` yoʻqligi muhim emas?

<details><summary>Javob</summary>

Ishlaydi (oʻlchangan: `cp main /tmp/x/ && cd /tmp/x && ./main` → `Salom, dunyo!`; hatto Go topilmaydigan
muhitda ham: `env -i ./main` → `Salom, dunyo!`). Tayyor fayl ichida bizning kod, `fmt` va runtime bor: unga na
`main.go`, na Go kerak. Faqat kompyuter turi bir xil boʻlsin (tizim va protsessor): Linux amd64 uchun yasalgan fayl
Windows'da yoki arm64'da ishlamaydi.
</details>

## ➕ Qoʻshimcha: boshqa tizim uchun fayl
Go bir kompyuterning oʻzida boshqa tizim uchun ham fayl yasay oladi. `GOOS` = maqsad tizim.

| Siz qayerdasiz | Buyruq | Natija |
|---|---|---|
| Linux / Mac | `GOOS=windows go build main.go` | `main.exe` |
| Windows, cmd | `set GOOS=linux` keyin `go build main.go` (oxirida `set GOOS=` bilan qaytaring) | `main` (Linux uchun) |
| Windows, PowerShell | `$env:GOOS="linux"; go build main.go` | `main` |

Linux'da oʻlchangan: `GOOS=windows go build main.go` → `main.exe`, **2 472 448 bayt** (ichida Windows uchun runtime
bor). Windows buyruqlari bu kompyuterda tekshirilmagan (Windows yoʻq); ular
[go.dev/doc/install/source#environment](https://go.dev/doc/install/source#environment) va cmd/PowerShell
qoidalariga asoslangan.


---

Eski videodagi (2026-11-25) uy vazifasi javoblari: [`../homework/README.md`](../homework/README.md).
