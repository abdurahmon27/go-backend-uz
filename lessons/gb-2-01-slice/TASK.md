# Vazifa · GB-2.01 · `discounted`: aksiya narxlari katalogni buzmasin

Videoda `afterOrder` xatosini tuzatdik: funksiya qoldiqni faqat koʻrsatishi kerak edi, lekin parametr orqali
ombordagi `stock` massiviga yozardi. Parametr slice sarlavhasining nusxasi, massiv esa bitta.

`main.go`dagi `discounted` xuddi shu xatoga ega. U aksiya narxlarini hisoblaydi: har bir narx `percent` foizga
arzonroq (narxlar tiyinda, `int64`, GB-1.03). Hozir u natijani `prices`ning oʻziga yozadi, shuning uchun
katalogdagi asl narxlar ham tushib ketadi. `go run .` buni koʻrsatadi:

```
narx oldin:  [1200000 500000 2500000]
aksiya:      [1080000 450000 2250000]
narx keyin:  [1080000 450000 2250000]
```

**Vazifa:** `discounted`ni shunday tuzatingki, aksiya narxlari **alohida massivda** boʻlsin, `prices` esa
oʻzgarmasin. Funksiya nomi, parametrlari va qaytaradigan turi oʻsha qoladi (`main` va testlar uni shunday chaqiradi).
Toʻgʻri natija:

```
narx oldin:  [1200000 500000 2500000]
aksiya:      [1080000 450000 2250000]
narx keyin:  [1200000 500000 2500000]
```

## Tekshirish

Shu papkada:

```bash
go test
```

`ok` chiqsa, vazifa bajarildi. `go test` ikki test faylini birga ishga tushiradi: `stock_test.go` (videodagi
`afterOrder` tuzatishini tekshiradi, hozir ham oʻtadi) va `tekshiruv_test.go` (vazifa; uni oʻzgartirmang). Ikkala
faylda `slices.Equal` bor: u ikki slice'ning uzunligi va elementlari bir xilligini tekshiradi, chunki slice'larni
`==` bilan solishtirib boʻlmaydi (kompilyator buni rad etadi). (`go test ./...` esa `solution/` papkani ham
tekshiradi; u yerda javob turgani uchun `ok` chiqadi.) Hozir (vazifadan oldin) natija shunday:

```
--- FAIL: TestVazifaKatalog (0.00s)
    tekshiruv_test.go:27: discounted katalogni oʻzgartirdi: prices endi [1080000 450000 2250000], [1200000 500000 2500000] boʻlib qolishi kerak (TASK.md)
--- FAIL: TestVazifaAlohidaMassiv (0.00s)
    tekshiruv_test.go:40: sale[0] = 1 yozilganda prices[0] ham 1 boʻldi: natija katalog bilan bitta massivda (TASK.md, make)
FAIL
exit status 1
FAIL	github.com/abdurahmon27/bozor	0.002s
```

## Maslahatlar

1. Videodagi `afterOrder` tuzatishini qayta koʻring: natija uchun `make` bilan yangi slice, uzunligi `prices`niki
   bilan teng. Turi `[]int` emas, `[]int64` boʻlsin. (Ishda koʻrasiz: standart kutubxonadagi `slices.Clone(prices)`
   ham alohida massivdagi nusxa qaytaradi; videodagi `make` yoʻli bilan ishlash yetarli.)
2. Sikl faqat yangi slice'ga yozsin, `prices`ni faqat oʻqisin. Oxirida yangi slice'ni qaytaring.
3. `TestVazifaAlohidaMassiv` natijaning birinchi elementini oʻzgartirib koʻradi: `prices[0]` oʻzgarsa, natija hali
   ham katalog bilan bitta massivda.

Javob: [`solution/main.go`](solution/main.go). Avval oʻzingiz urinib koʻring.
