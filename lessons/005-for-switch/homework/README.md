# 05 · Uy vazifasi javoblari

Avval oʻzingiz urinib koʻring. 🙂 Hammasi `go1.26.8` da ishga tushirilgan.

## 🟢 Oson
[`oson/main.go`](oson/main.go). Ifodasiz `switch` birinchi `true` boʻlgan `case` ni bajaradi, shuning uchun
`i%15 == 0` **birinchi** turishi kerak: aks holda 15 uchun `i%3 == 0` oldin ishlab, `Fizz` chiqadi.

```
1 2 Fizz 4 Buzz Fizz 7 8 Fizz Buzz 11 Fizz 13 14 FizzBuzz
```

## 🟡 Oʻrta
[`orta/main.go`](orta/main.go). `strconv.Itoa` sonni matnga aylantiradi, `range` uning rune'lari ustida yuradi,
`switch` har bir raqamni soʻzga aylantiradi (`default` kutilmagan belgilar uchun):

```
2026: ikki nol ikki olti
1990: bir toʻqqiz toʻqqiz nol
-7: minus yetti
```

## 🟠 Qiyin
[`qiyin/main.go`](qiyin/main.go). Oddiy `break` faqat ichki sikldan chiqadi; `break tashqi` ikkalasidan ham:

```
birinchi 0: qator 1, ustun 1
```

## 🔴 Guru
[`guru/`](guru/) alohida modul, `go.mod` da `go 1.21` (eski semantika). Ishga tushirish: `cd homework/guru && go run .`

```
a.uz b.uz c.uz ← d := d
a.uz b.uz c.uz ← &domenlar[i]
```

3-usul: [`guru/buildtag/main.go`](guru/buildtag/main.go), `cd homework/guru && go run ./buildtag`:

```
a.uz b.uz c.uz ← //go:build go1.22
```

1. `d := d`: sikl ichida yangi oʻzgaruvchi yaratiladi, har manzil oʻz nusxasiga qaraydi. Go 1.22 gacha shunday
   yozilardi; `go 1.22+` modulda bu qator ortiqcha va `go fix` uni oʻchiradi.
2. `&domenlar[i]`: sikl oʻzgaruvchisi emas, slice elementining manzili. Bu yerda nusxa umuman yoʻq
   (lekin slice keyin `append` bilan oʻssa, manzil eski massivga qarab qolishi mumkin: 08 va 11-darslar).
3. Fayl boshidagi `//go:build go1.22` qatori **shu fayl** uchun til versiyasini koʻtaradi: modul `go 1.21` da qolsa ham,
   bu faylda har iteratsiyada yangi `d` boʻladi. (Teskarisi ham ishlaydi: `go 1.26` modulda `//go:build go1.21`
   faylni eski semantikaga qaytaradi; biz ikkalasini tekshirdik.) Bu haqda batafsil: 16-dars, `go.mod` dagi `go` qatori.
