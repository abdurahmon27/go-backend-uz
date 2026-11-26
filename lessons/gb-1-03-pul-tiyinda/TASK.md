# Vazifa: choy sotuvi, yetkazib berish bilan

`main.go` hozir videodagi holatda: olma, 12 000 soʻm, `sold := 3`, QQS 12 %. Oʻtgan darsdagi qatorlar
(`bozor ishga tushdi`, `versiya: 0.1.0`, `stock`, `active`, `weight`) joyida qoladi. Sotuvni quyidagiga oʻzgartiring:

- mahsulot: `choy`, narxi **1 999 soʻm** (dasturda tiyinda, `var price int64`)
- sotilgan: **7 dona** (`sold := 7`, turi `int` boʻlib qolsin)
- QQS: **12 %**, faqat mahsulot summasiga
- yetkazib berish: **15 000 soʻm**, QQS qoʻshilmaydi, jamiga qoʻshiladi

Dasturning birinchi qatori oʻzgarmaydi (`bozor ishga tushdi`), **oxirgi uch qatori** aynan shunday boʻlsin:

```
choy 7 dona
QQS: 1679 soʻm 16 tiyin
jami: 30672 soʻm 16 tiyin
```

Pul hisobi tiyinda, butun sonlarda: hisobda `int64(sold)` ishlatilsin (test buni ham tekshiradi).

## Tekshirish

```bash
go test
```

`ok` chiqsa, vazifa bajarildi. Xato boʻlsa, test qaysi qator qanday boʻlishi kerakligini oʻzbekcha yozadi.
`main_test.go` faylini oʻzgartirmang: uni GB-1.08 darsida toʻliq oʻqiymiz.

## Maslahatlar

1. Yetkazib berish narxi uchun ham `const` yozing va uni ham tiyinga oʻtkazing (`… * tiyinPerSom`). Uni
   `delivery := …` bilan alohida oʻzgaruvchiga olsangiz, untyped constant default turni, `int`'ni oladi:
   `mismatched types` xatosini darsdagidek oʻqing va tuzating.
2. QQS qatori `jami` qatori bilan bir xil usulda chiqadi: `vat / tiyinPerSom` — soʻm, `vat % tiyinPerSom` — tiyin.
3. Avval koʻpaytiring, boʻlish eng oxirida: `subtotal * vatPercent / 100`.

Javob `solution/` papkasida. Avval oʻzingiz urinib koʻring.
