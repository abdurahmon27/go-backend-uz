# Vazifa: chekka yetkazib berish qatorini qoʻshing

`main.go` hozir chekni uch qator bilan chiqaradi: mahsulot, QQS, jami. Lekin `jami` ichida yetkazib berish narxi
ham bor (`deliverySom`, 15 000 soʻm), chekda esa u koʻrinmaydi: mijoz summalarni qoʻshib, jamiga yetmaydi.

**Nima qilish kerak:** QQS va jami orasiga `yetkazish` qatorini qoʻshing. Yuqoridagi qatorlarga (ishga tushish,
versiya, zero value, ombor) tegmang: test birinchi qator `bozor ishga tushdi` ekanini va oxirgi toʻrt qatorni
tekshiradi. U boshqa qatorlar bilan **bir xil
shablonda** boʻlsin: nom 14 belgilik ustunda chapda, summa soʻmda, nuqtadan keyin ikki raqamli tiyin, oʻng chetlari
bir chiziqda.

Kutilgan natija (`go run .`, oxirgi toʻrt qator):

```
choy, 7 dona       13993.00
QQS                 1679.16
yetkazish          15000.00
jami               30672.16
```

**Tekshirish:** shu papkada

```bash
go test
```

`ok` chiqsa, vazifa bajarildi. Hozir test `FAIL` beradi: `--- FAIL: TestChek` qatoridan keyin
`main_test.go:62:` (xabar test faylining qaysi qatoridan chiqqani) va xabarning oʻzi keladi: chek nechta qator
boʻlishi kerak, keyin `kerak:` va `hozir:` qatorlari. Qator soni toʻgʻri boʻlib, matni farq qilsa, xabar
`3-qator: kerak: "…" hozir: "…"` koʻrinishida boʻladi (qoʻshtirnoq `%q`'dan: probellarni koʻrsatadi).
`main_test.go` faylini oʻzgartirmang: u `main()` chiqargan matnni kutilgan chek bilan solishtiradi (test yozishni
GB-1.08'da ochamiz). `go test` vet'ning bir qismini ham ishga tushiradi, shuning uchun feʼl va argument mos kelmasa
(`%d` ga matn berilsa) yoki `string(son)` yozilsa, test ishga tushmasdan xato chiqadi.

**Maslahatlar**
1. `deliverySom` soʻmda, qolgan summalar tiyinda. Soʻm qismi `deliverySom`ning oʻzi, tiyin qismi esa 0. Yoki avval
   tiyinga aylantirib (`deliverySom * tiyinPerSom`), keyin boshqa qatorlardagidek `/` va `%` bilan ikkiga ajrating.
2. GB-1.03'dagi qoida: `int64` va `int` aralashsa, `mismatched types` chiqadi. Yangi oʻzgaruvchi kerak boʻlsa, turini
   `var … int64 = …` bilan aniq yozing.
3. Shablonni boshqa qatorlardan koʻchiring: `"%-14s %9d.%02d\n"`. Kenglikni oʻzgartirsangiz, ustun buziladi.
   Summani tayyor matn (`"15000.00"`) yoki qoʻlda son (`15000`) qilib yozmang: ikkinchi test buni tekshiradi.

Javob: [`solution/main.go`](solution/main.go). Avval oʻzingiz urinib koʻring.
