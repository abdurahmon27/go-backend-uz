# Vazifa · GB-2.07: hisobot faqat `w`'ga yozilsin, `lineCounter` — `io.Writer`

Videoda `writeReport(w io.Writer, products []Product) error` bitta hisobotni `os.Stdout`'ga, `os.Create` ochgan
faylga va oʻzimiz yozgan `*byteCounter`'ga yozdi: uchalasida ham `Write(p []byte) (int, error)` metodi bor, shuning
uchun uchalasi ham `io.Writer`. `main.go`'da ikkita ish chala qolgan.

**Nima qilish kerak** (`main.go`'da):

1. **`writeLowStock`**: qoldigʻi (`Stock`) `limit`'dan kam mahsulotlar uchun bittadan qator yozadi, masalan
   `oʻrik (35000 soʻm) — 5 dona qoldi`. Hozir u qatorlarni `w`'ga emas, toʻgʻridan-toʻgʻri ekranga yozadi, shuning
   uchun faylga yoki testdagi `strings.Builder`'ga hech narsa tushmaydi. Hamma qator `w`'ga yozilsin. Yozish xato
   qaytarsa (masalan, disk toʻlgan), `writeLowStock` darhol shu xatoni qaytarsin; xato boʻlmasa — `nil`.
2. **`lineCounter`**: `byteCounter` kabi tur, lekin baytlarni emas, **qatorlarni** sanaydi. Unga `Write` metodini
   yozing, shunda `*lineCounter` `io.Writer` boʻladi: har chaqiruvda `p`'dagi qator oxirlari (`'\n'`) soni
   `lines`'ga qoʻshilsin, natija esa `len(p)` va `nil`. Uch mahsulotli hisobot yozilgach, `lines` = 3.

`writeReport`'ni, `Product`'ni va `String`'ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli `tekshiruv_test.go`: uni oʻzgartirmang. Hozir u shunday yiqiladi
(qisqartirilgan; tepadagi ikki qator — `fmt.Printf` testning ichidan ekranga yozgani):

```
oʻrik (35000 soʻm) — 5 dona qoldi
choy (25000 soʻm) — 8 dona qoldi
--- FAIL: TestVazifaKamQoldi (0.00s)
    tekshiruv_test.go:61: writeLowStock w'ga hech narsa yozmadi: fmt.Printf ekranga yozadi, w'ga fmt.Fprintf(w, …) yozadi (TASK.md, 1-qadam)
--- FAIL: TestVazifaQatorSanagich (0.00s)
    tekshiruv_test.go:78: *lineCounter io.Writer emas: unga Write(p []byte) (int, error) metodini yozing (TASK.md, 2-qadam)
FAIL
```

**Maslahatlar**

- `fmt.Fprintf` ikkita natija qaytaradi: yozilgan baytlar soni va `error`. Videodagi `writeReport`'ga qarang.
  Shablondagi `%v` (GB-1.04) `Product`'ni uning `String` metodi orqali yozadi, shablonni oʻzgartirmang.
- `p`'dagi qator oxirlarini sanash: `bytes.Count(p, []byte("\n"))` — `[]byte("\n")` matnni bayt slice'iga
  oʻgiradi, `bytes.Count` `p`'da shu bayt necha marta borligini qaytaradi (`import`'ga `"bytes"` qoʻshing).
- `Write` `lines`'ni oʻzgartiradi. Receiver qanday boʻlishi kerak (GB-2.05)? Test xabari ham shuni soʻraydi.

Javob: [`solution/main.go`](solution/main.go).
