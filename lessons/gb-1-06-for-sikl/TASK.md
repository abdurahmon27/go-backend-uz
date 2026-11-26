# Vazifa: aksiyaga chegara qoʻying

`main.go` hozir "har uchinchi dona bepul" aksiyasini hisoblaydi: 7 donada 3- va 6-dona bepul. Lekin 30 dona non
olgan mijozga 10 dona bepul ketadi, doʻkon uchun bu juda qimmat.

**Yangi qoida:** har uchinchi dona hamon bepul, lekin **bitta buyurtmada koʻpi bilan 2 dona**. Ikkita bepul dona
berilgandan keyin qolgan hamma dona pullik. Chegarani nomli konstanta qiling:

```go
const maxFree = 2 // bitta buyurtmada koʻpi bilan shuncha dona bepul
```

Kutilgan natija (`qty = 12` boʻlsa: 3- va 6-dona bepul, 9- va 12-dona endi pullik):

```
choy, 12 dona      19990.00
QQS                 2398.80
yetkazish          15000.00
jami               37388.80
aksiya: 2 dona bepul
```

7 donalik buyurtma oʻzgarmaydi (u yerda baribir 2 dona bepul), 2 donada esa aksiya qatori chiqmaydi.

**Tekshirish:** shu papkada

```bash
go test
```

`ok` chiqsa, vazifa bajarildi. Hozir test `FAIL` beradi va qaysi buyurtmada nima kutilganini koʻrsatadi
(`kerak:` va `hozir:` qatorlari). `main_test.go` faylini oʻzgartirmang: u `name`, `price`, `qty`, `stock` ni
almashtirib, `main()` ni bir necha marta ishga tushiradi (test yozishni GB-1.08'da ochamiz).

**Maslahatlar**
1. Oʻzgarish faqat sikl ichida: qaysi holatda `continue` ishlashi kerak? Dona bepul boʻlishi uchun endi **ikki**
   shart bir vaqtda true boʻlishi kerak.
2. `free` bepul donalarni allaqachon sanaydi. GB-1.05'dagi `&&` ikkala tomon true boʻlgandagina true beradi.
3. `break` bu yerda notoʻgʻri: u qolgan donalarni ham hisobdan chiqarib yuboradi (videodagi savolni eslang).

Javob: [`solution/main.go`](solution/main.go). Avval oʻzingiz urinib koʻring.
