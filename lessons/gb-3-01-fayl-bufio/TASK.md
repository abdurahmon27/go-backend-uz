# Vazifa · GB-3.01: `loadCatalog` izohlarni oʻtkazib yuborsin va xato qatorini aytsin

Videoda `loadCatalog(r io.Reader)` katalogni `bufio.Scanner` bilan qatorma-qator oʻqidi: har qator
`parseLine`'dan oʻtib `Product` boʻldi, oxirida `sc.Err()` qaytarildi. Haqiqiy katalog fayllarida izohlar va boʻsh
qatorlar boʻladi, xato boʻlsa esa odam uni qaysi qatordan qidirishni bilishi kerak.

**Nima qilish kerak** (`main.go`'da, faqat `loadCatalog`'ning tanasini oʻzgartirasiz):

1. **Izoh va boʻsh qatorlar.** Boʻsh qatorni, faqat probellardan iborat qatorni va (boshidagi probellardan keyin)
   `#` bilan boshlanadigan qatorni tashlab oʻting: ular mahsulot emas, xato ham emas.

   ```
   # bozor katalogi, 2026-11-27
   1;olma;12000;40

     # ichimliklar
   4;choy;25000;8
   ```

   Bu matndan ikkita mahsulot chiqishi kerak: `olma` va `choy`.
2. **Qator raqami.** Biror qator xato boʻlsa, xato matni `«N-qator: »` bilan boshlansin, `N` — fayldagi qator raqami
   (1 dan boshlab, izoh va boʻsh qatorlar ham sanaladi). Masalan yuqoridagi matnning 4-qatori
   `4;non;abc;100` boʻlsa: `4-qator: narx: strconv.Atoi: parsing "abc": invalid syntax`
   (`narx: …` qismini `parseLine` allaqachon yozadi).
   - Asl xato **`%w` bilan** ichida qolsin (GB-2.08): chaqiruvchi `errors.As` / `errors.Is` bilan uni topa olsin.
   - Scanner'ning oʻz xatosiga ham (`sc.Err()`, masalan videodagi `token too long`) qator raqami qoʻshilsin: bu
     Scanner oʻqiy olmagan qator, yaʼni oxirgi oʻqilgan qatordan keyingisi. `katta.txt` uchun:
     `3-qator: bufio.Scanner: token too long`, va `errors.Is(err, bufio.ErrTooLong)` `true`.

`parseLine`'ni, `parsePrice`'ni, `writeReport`'ni va `run`'ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli `tekshiruv_test.go`: uni oʻzgartirmang. `TestVideo` videodagi `katalog.txt`'ni
ochadi, vazifa testlari esa `loadCatalog`'ni `strings.NewReader(…)` bilan chaqiradi (videodagidek). Hozir u shunday
yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaIzoh (0.00s)
    tekshiruv_test.go:74: loadCatalog izohli katalogda xato qaytardi («4 ta maydon kerak, bor 1»): boʻsh qatorlarni va # bilan boshlangan qatorlarni tashlab oʻting (TASK.md, 1-qadam)
--- FAIL: TestVazifaQatorRaqami (0.00s)
    --- FAIL: TestVazifaQatorRaqami/narx_son_emas (0.00s)
        tekshiruv_test.go:146: xato «2-qator: » bilan boshlanishi kerak, hozir «narx: strconv.Atoi: parsing "abc": invalid syntax» (TASK.md, 2-qadam)
    …
    --- FAIL: TestVazifaQatorRaqami/uzun_qator (0.00s)
        tekshiruv_test.go:146: xato «3-qator: » bilan boshlanishi kerak, hozir «bufio.Scanner: token too long» (TASK.md, 2-qadam)
FAIL
```

Qoʻlda ham tekshirib koʻring: `go run . katalog.txt` (4 mahsulot), `go run . katta.txt` (`xato: 3-qator: …`),
`cat katalog.txt | go run .`.

**Maslahatlar**

- Qator sanagich: sikldan oldin `n := 0`, sikl boshida `n++` (GB-1.06). Izoh qatorini `continue` bilan oʻtkazing,
  lekin `n++`'dan **keyin**: aks holda keyingi qatorlarning raqami siljiydi.
- `strings.TrimSpace(s)` — boshidagi va oxiridagi probellarni olib tashlaydi; `strings.HasPrefix(s, "#")` — `s`
  `#` bilan boshlanadimi. Ikkalasi ham `strings` paketida (u `main.go`'da allaqachon import qilingan).
- Xatoni oʻrash: `fmt.Errorf("%d-qator: %w", n, err)`. `%v` bilan yozsangiz, matn bir xil chiqadi, lekin
  `errors.As` ichidagi `*strconv.NumError`'ni topa olmaydi: test buni tekshiradi.
- `sc.Err()`'ni oxirida endi alohida tekshiring: `if err := sc.Err(); err != nil { … }`, xatoni `n+1` bilan oʻrang,
  aks holda `return products, nil`.

Javob: [`solution/main.go`](solution/main.go).
