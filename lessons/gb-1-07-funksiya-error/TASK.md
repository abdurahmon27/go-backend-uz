# Vazifa · GB-1.07: miqdor buyruq qatoridan, `parseSold`

Hozir bozor nom va narxni buyruq qatoridan oladi (`./bozor add olma 12000`), lekin sotilgan dona soni
`main` ichida `sold := 7` deb qoʻlda yozilgan. Uni ham buyruq qatoridan oling.

**Nima qilish kerak** (`main.go`'da):

1. Buyruq toʻrtinchi soʻz sifatida miqdorni oladi, chek shu miqdor bilan hisoblanadi:

   ```
   $ ./bozor add olma 12000 5
   bozor ishga tushdi
   versiya: 0.1.0
    0 false
   olma 1200000 true 40 1.5
   olma 35
   olma, 5 dona       48000.00
   QQS                 5760.00
   yetkazish          15000.00
   jami               68760.00
   aksiya: 1 dona bepul
   ```

2. Miqdorni `func parseSold(s string) (int, error)` funksiyasi oʻqiydi:
   - son boʻlmasa (`besh`), `strconv.Atoi` qaytargan xatoni oʻzgartirmasdan qaytaradi;
   - `1` dan kichik boʻlsa (`0`, `-2`), oʻz xatosini qaytaradi, matni aynan shunday:
     `miqdor kamida 1 boʻlishi kerak, keldi 0` (oxirida kelgan son).
   `main`'dagi eski `sold <= 0` tekshiruvi endi kerak emas: bu qoida `parseSold`'ga koʻchadi
   (`parsePrice` va `price <= 0` bilan videoda shunday qildik).
3. Xato boʻlsa, `main` uni narx xatosi kabi xato oqimiga (`os.Stderr`) `xato: …` deb yozadi va `os.Exit(1)`
   bilan chiqadi.
4. Soʻzlar soni toʻgʻri boʻlmasa, foydalanish qatori yangilanadi:
   `foydalanish: bozor add NOM NARX MIQDOR` (xato oqimiga, chiqish kodi 1).

Ishga tushish qatorlari, narx qoidasi (`parsePrice`) va chek oʻzgarmaydi.

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli (`main_test.go`) dasturni turli argumentlar bilan ishga tushiradi va natijani
solishtiradi; uni oʻzgartirmang, ichini GB-1.08 darsida ochamiz. Qoʻlda ham sinab koʻring:

```bash
go build && ./bozor add olma 12000 5
./bozor add olma 12000 besh; echo $?
```

**Maslahatlar**

- `parseSold` `parsePrice`'ning tuzilishini takrorlaydi, faqat natija `int` va tiyinga koʻpaytirish yoʻq.
- `sold := 7` qatori oʻrniga `sold, err := parseSold(os.Args[4])` yozing. `err` narx qatorida eʼlon qilingan,
  `sold` esa yangi: chapda kamida bitta yangi nom boʻlgani uchun `:=` ishlaydi, `err`'ga shunchaki yangi qiymat
  beriladi (videoda `price, err := parsePrice(…)` bilan xuddi shunday). `parseSold` ichidagi oʻzgaruvchiga boshqa nom
  bering (masalan `n`).
- `len(os.Args)` endi nechta boʻlishi kerak? Dastur nomi ham hisobga kiradi.

Javob: [`solution/main.go`](solution/main.go).
