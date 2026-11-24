# 03 · Uy vazifasi javoblari

Avval o'zingiz urinib ko'ring. 🙂 Hamma natijalar `go1.26.8 linux/amd64`da olingan.

## 🟢 Oson: hafta kunlari
[`easy/main.go`](easy/main.go) → `1 5 7` va `0`.
`Monday Weekday = iota + 1` bilan dushanba 1 dan boshlanadi, `0` esa "kun tanlanmagan" bo'lib qoladi. Shunda qiymat
berilmagan `Weekday` (zero value) tasodifan dushanba bo'lib qolmaydi. Standart `time.Weekday` boshqacha: u yerda
`Sunday Weekday = iota`, ya'ni yakshanba 0.

## 🟡 O'rta: bit flaglar
[`medium/main.go`](medium/main.go):
```go
const (
	Read  Perm = 1 << iota // 001
	Write                  // 010
	Exec                   // 100
)
p := Read | Write          // 011
p&Write != 0               // true
p&^Write                   // 001  (&^ = "bit'ni o'chir")
```
Har bir flag alohida bit, shuning uchun ular `|` bilan birlashadi va `&` bilan tekshiriladi. `io/fs.FileMode` xuddi
shunday qurilgan: `ModeDir FileMode = 1 << (32 - 1 - iota)`
([io/fs/fs.go](https://github.com/golang/go/blob/go1.26.8/src/io/fs/fs.go#L192)).

## 🟠 Qiyin: `1.5 * time.Second`
```
fmt.Println(1.5 * time.Second)
# xato: 1.5 (untyped float constant) truncated to int64
```
`time.Second` `Duration` turidagi konstanta, `type Duration int64`. Untyped `1.5` avval `Duration`ga, ya'ni butun
songa aylanishi kerak, lekin `1.5` butun songa aniq aylanmaydi, va konstanta yaxlitlanmaydi: kompilyator rad etadi.
To'g'ri yo'llar ([`hard/main.go`](hard/main.go)):
```go
fmt.Println(1500 * time.Millisecond) // 1.5s
fmt.Println(time.Second * 3 / 2)     // 1.5s: avval 3 s, keyin 2 ga butun bo'lish
```

## 🔴 Guru: `time` package'idagi `iota`
[`time/time.go`](https://github.com/golang/go/blob/go1.26.8/src/time/time.go#L327):
- `January Month = 1 + iota` (L327): oylar 1 dan 12 gacha, odamlar yozganidek. `time.Date(2026, 11, 25, ...)`
  dagi `11` to'g'ridan-to'g'ri `November` bo'ladi.
- `Sunday Weekday = iota` (L355): hafta kunlari 0 dan.

[`guru/main.go`](guru/main.go):
```
January 1
Sunday 0
%!Month(0)
```
Qiymat berilmagan `time.Month` 0 ga teng, lekin 0-oy yo'q: `String()` buni `%!Month(0)` deb ko'rsatadi. Ya'ni Go
jamoasi bu yerda "zero value ma'noli bo'lsin" qoidasidan ko'ra "raqam odamlar ishlatadigan raqamga teng bo'lsin"
ni tanlagan. Ikkalasi ham ongli tanlov, muhimi: tanlovni bilib qilish.
