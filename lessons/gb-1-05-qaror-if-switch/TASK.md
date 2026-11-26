# Vazifa: nom, sotuv va narx qoidasi

`main.go` hozir notoʻgʻri miqdorni (`sold <= 0`) va ombordagidan koʻp buyurtmani (`sold > stock`) rad etadi. Lekin
mahsulotning oʻzi tekshirilmaydi: `name` boʻsh boʻlsa, chek `, 7 dona` bilan boshlanadi; `active` false boʻlsa
(mahsulot sotuvda emas), u baribir sotiladi; `price` berilmasa (zero value, `0`) yoki manfiy boʻlsa, mahsulot
tekinga yoki "minus" pulga sotiladi.

**Nima qilish kerak:** `main()` da `sold := 7` qatoridan keyin, **miqdor tekshiruvidan (`if sold <= 0`) oldin**,
uchta qoida qoʻshing, aynan shu tartibda. Har biri videodagidek: shart, xabar, `return`.

1. `name` boʻsh matn boʻlsa (`""`):

   ```
   xato: mahsulot nomi boʻsh
   ```

2. `active` false boʻlsa (`!` bilan yozing), xabarda mahsulot nomi boʻlsin, masalan `name` "choy" boʻlsa:

   ```
   xato: choy sotuvda emas
   ```

3. `price` nol yoki manfiy boʻlsa, **yoki** 10 000 000 soʻmdan katta boʻlsa (narxga ortiqcha nollar yozilib
   qolgan), narx tiyinda chiqsin. Masalan `price` 0 boʻlsa:

   ```
   xato: narx notoʻgʻri: 0 tiyin
   ```

   Aynan 10 000 000 soʻm (1 000 000 000 tiyin) va 1 tiyin hali toʻgʻri narx.

Bir nechta xato boʻlsa, birinchi qoidaning xabari chiqadi: masalan nom ham boʻsh, narx ham 0 boʻlsa, faqat
`xato: mahsulot nomi boʻsh`. Xabardan oldingi toʻrt qator (`bozor ishga tushdi`, `versiya: 0.1.0` va oʻtgan
darslardagi ikki sinov qatori) joyida qoladi.

**Tekshirish:** shu papkada

```bash
go test
```

`ok` chiqsa, vazifa bajarildi. Hozir test `FAIL` beradi, masalan:

```
--- FAIL: TestQoidalar (0.13s)
    main_test.go:210: nom boʻsh (name="" price=199900 active=true stock=40 sold=7):
          kerak:
            "xato: mahsulot nomi boʻsh"
          hozir:
            " 33"
            ", 7 dona           13993.00"
            "QQS                 1679.16"
            "yetkazish          15000.00"
            …
```

`main_test.go` faylini oʻzgartirmang: u `main.go`'dagi `name = …`, `price = …`, `active = …`, `stock := 40` va
`sold := 7` qatorlaridagi qiymatlarni almashtirib, dasturni har xil buyurtma bilan ishga tushiradi (test yozishni
GB-1.08'da ochamiz). Shuning uchun bu besh qatorni oʻchirmang. Videodagi qoidalar ham
(miqdor, ombor, yetkazish chegaralari) tekshiriladi, ularni buzmang. Xabardagi `ʻ` oʻrniga oddiy `'` yozsangiz ham
test qabul qiladi.

**Maslahatlar**
1. Matn ham `==` bilan solishtiriladi: `name == ""`.
2. `!active` "active emas" degani: `active` false boʻlganda true.
3. Narx qoidasi bitta `if` bilan yoziladi: ikki shartni `||` bilan ulang. Chegara tiyinda: `10_000_000*tiyinPerSom`.
   Aynan 10 000 000 soʻm rad etilmasligi uchun `>` va `>=` farqiga qarang.
4. Xohlasangiz, uchala qoidani bitta ifodasiz `switch`ʼga yigʻing: har bir `case`
   bitta qoida, `switch` birinchi true `case`'ni bajaradi, tartib oʻz-oʻzidan saqlanadi.

Javob: [`solution/main.go`](solution/main.go). Avval oʻzingiz urinib koʻring.
