# 02 · Uy vazifasi javoblari

Avval o'zingiz urinib ko'ring. 🙂 Hamma natijalar `go1.26.8 linux/amd64`da olingan.

## 🟢 Oson: `int16` overflow
[`easy/main.go`](easy/main.go) → `-32768`.
`32767` ikkilikda `0111 1111 1111 1111`. Bir qo'shilganda ko'chish eng chap bitga o'tadi: `1000 0000 0000 0000`.
`int16`da eng chap bitning vazni `-32768`, demak natija `-32768`. Bu `int8`dagi `127 + 1 = -128` bilan aynan bir xil
qoida, faqat 16 bitda.

## 🟡 O'rta: `int32` Unix vaqti qachon to'ladi?
[`medium/main.go`](medium/main.go):
```
2038-01-19 03:14:07 +0000 UTC
1901-12-13 20:45:52 +0000 UTC
```
Unix vaqti: 1970-01-01 00:00:00 UTC dan beri o'tgan sekundlar. `int32`ning eng katta qiymati 2 147 483 647, ya'ni
**2038-yil 19-yanvar, 03:14:07 UTC**. Bir sekunddan keyin son `-2147483648` ga aylanadi, va bu 1901-yil 13-dekabr.
Bu "2038-yil muammosi": eski tizimlarda vaqt 32 bitli songa yozilgan. Masalan, MySQL'dagi `TIMESTAMP` turining oralig'i
aynan `'2038-01-19 03:14:07' UTC` da tugaydi ([MySQL hujjati](https://dev.mysql.com/doc/refman/8.4/en/datetime.html)).
Go'da `time.Time` ichida 64 bitli son ishlatiladi, shuning uchun bu muammo yo'q.

## 🟠 Qiyin: `int64` qo'shishda overflow'ni aniqlash
[`hard/main.go`](hard/main.go). G'oya: overflow faqat **bir xil ishorali** ikki son qo'shilganda bo'ladi, va bunda
natijaning ishorasi o'zgaradi:
```go
sum := a + b
overflow := (a >= 0) == (b >= 0) && (sum >= 0) != (a >= 0)
```
`a + b` ni oldin hisoblash xavfsiz: Go'da signed overflow aniq belgilangan, panic ham, "undefined behavior" ham yo'q
(C'dan farqi shu). Boshqa yo'l: `b > 0 && a > math.MaxInt64-b` kabi chegarani oldindan tekshirish.
Unsigned sonlar uchun standart library'da tayyor funksiya bor: `math/bits.Add64` ko'chish (carry) bitini qaytaradi.

## 🔴 Guru: `abs(int8(-128))`
[`guru/main.go`](guru/main.go):
```
5 -128
5 128
```
Oddiy `abs`: `if x < 0 { return -x }`. Lekin `int8` oralig'i nosimmetrik: −128 … 127. +128 sig'maydi, shuning uchun
`-(-128)` overflow bo'lib, yana `-128` ga aylanadi: "musbat" bo'lishi kerak bo'lgan natija manfiy chiqadi. Two's
complement'da `-x` = bitlarni teskari qilib, 1 qo'shish: `1000 0000` → `0111 1111` → `1000 0000`. Yechim: natijani
kattaroq turda qaytarish (`int16`), yoki `math.MinInt8`ni alohida tekshirish. Xuddi shu tuzoq `int64` va
`math.MinInt64` da ham bor.
