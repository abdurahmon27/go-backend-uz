# Vazifa · GB-2.10: `shortLabel` va `checkName`

Videoda `short(name, limit)` nomni `limit` ta **rune**'gacha qisqartirdi: `for i := range name` har bir rune
boshlangan baytning indeksini beradi, kesish shu indeksda boʻladi, shuning uchun belgi oʻrtasidan kesilmaydi.
`main.go` oxirida shu mavzudagi ikkita funksiya bor, lekin ular hali boʻsh: `shortLabel` nomni oʻzgartirmay
qaytaradi, `checkName` doim `nil` qaytaradi.

**Nima qilish kerak** (`main.go`'da, faqat shu ikki funksiyaning tanasini yozasiz):

1. `func shortLabel(name string, width int) string` — chek ustuni uchun yorliq:
   - nom `width` ta rune'ga sigʻsa, oʻzini qaytarsin (`"oʻrik qoqi"`, 10 → `"oʻrik qoqi"`: 10 rune, 11 bayt);
   - sigʻmasa, birinchi `width-1` ta rune va oxirida uch nuqta belgisi `…` (natija jami `width` ta rune):
     `"sariq sabzi"`, 10 → `"sariq sab…"`, `"tovuq goʻshti"`, 10 → `"tovuq goʻ…"`, `"oʻrik"`, 3 → `"oʻ…"`.
   - `width` kamida 1 deb hisoblang (0 va manfiy qiymatni test bermaydi).
2. `func checkName(name string) error` — yangi mahsulot nomini tekshirsin va uch holatda xato qaytarsin:
   - nom boʻsh (`""`);
   - nom toʻgʻri UTF-8 emas (masalan `"o\xca"`: rune oʻrtasidan kesilgan bayt);
   - nom **20 rune**'dan uzun. Baytlar emas: `"oʻgʻil bola koʻylagi"` 20 rune (23 bayt), u oʻtishi kerak;
     `"olma qizil, 500 gramm"` 21 rune (21 bayt), u rad etilishi kerak.
   Boshqa hollarda `nil`. Xato matni sizning ixtiyoringizda, test faqat xato bor-yoʻqligini tekshiradi.

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli `tekshiruv_test.go`: uni oʻzgartirmang. Hozir u shunday yiqiladi:

```
--- FAIL: TestVazifaShortLabel (0.00s)
    tekshiruv_test.go:29: shortLabel("sariq sabzi", 10) = "sariq sabzi"; kerak "sariq sab…" (TASK.md, 1-qadam)
    tekshiruv_test.go:29: shortLabel("tovuq goʻshti", 10) = "tovuq goʻshti"; kerak "tovuq goʻ…" (TASK.md, 1-qadam)
    tekshiruv_test.go:29: shortLabel("paxta yogʻi", 9) = "paxta yogʻi"; kerak "paxta yo…" (TASK.md, 1-qadam)
    tekshiruv_test.go:29: shortLabel("oʻrik", 3) = "oʻrik"; kerak "oʻ…" (TASK.md, 1-qadam)
--- FAIL: TestVazifaCheckName (0.00s)
    tekshiruv_test.go:55: checkName(""): xato kutilgan (boʻsh nom), nil qaytdi (TASK.md, 2-qadam)
    tekshiruv_test.go:55: checkName("o\xca"): xato kutilgan (UTF-8 emas: rune oʻrtasidan kesilgan), nil qaytdi (TASK.md, 2-qadam)
    tekshiruv_test.go:55: checkName("oʻrik qoqi, 500 gramm"): xato kutilgan (21 rune: uzun), nil qaytdi (TASK.md, 2-qadam)
    tekshiruv_test.go:55: checkName("olma qizil, 500 gramm"): xato kutilgan (21 rune, 21 bayt: uzun), nil qaytdi (TASK.md, 2-qadam)
FAIL
```

**Maslahatlar**

- Rune'lar soni: `utf8.RuneCountInString(name)`, toʻgʻri UTF-8'mi: `utf8.ValidString(name)`. Ikkalasi ham
  `unicode/utf8` paketida: import'ni videodagidek ikki paketli qiling (`import ( "fmt" "unicode/utf8" )`, har biri
  alohida qatorda).
- `shortLabel` uchun yangi sikl shart emas: tayyor `short(name, width-1)` birinchi `width-1` ta rune'ni beradi,
  unga `"…"`ni `+` bilan qoʻshing.
- `len(name)` bu yerda notoʻgʻri oʻlchov: u baytlarni sanaydi. `"…"`ning oʻzi uch bayt, lekin bitta rune.
  Uch nuqta belgisini (`…`, U+2026) shu fayldan nusxalang: uchta oddiy nuqta `...` emas.
- Eʼtibor: rune chegarasida kesish UTF-8'ni buzmaydi, lekin `g` va `ʻ` (yoki `o` va `ʻ`) orasidan kesilsa, `gʻ` harfi `g` boʻlib
  qoladi (`"paxta yogʻi"`, 10 → `"paxta yog…"`). Test bunday holatni bermaydi; ishda buni alohida qoida bilan
  hal qilasiz.
- `checkName`'da shartlar tartibi: avval boʻsh, keyin UTF-8, keyin uzunlik. Xato uchun `errors.New` yoki
  `fmt.Errorf` (GB-1.07, GB-2.08).

Javob: [`solution/main.go`](solution/main.go).
