# Go 0 · 02 · Oʻzgaruvchi — qiymatga yopishtirilgan yorliq

Oldingi dars: [GO-0.01 · Dastur nima? Birinchi dasturimizni qatorma-qator oʻqiymiz](../000-01-birinchi-dastur/) ·
Keyingi dars: [GO-0.03 · Son va matn](../000-03-son-va-matn/)

Hech qachon kod yozmaganlar uchun. Videodagi buyruqlar terminal shu dars papkasida ochilganda ishlaydi
(oldingi darsdagi `go run main.go` bilan bir xil, faqat oldiga papka nomi qoʻshiladi). Kompyuteringizda Go boʻlmasa, har bir `main.go` ni
[go.dev/play](https://go.dev/play) ga koʻchirib, **Run** tugmasini bosing.

| Papka | Nima |
|---|---|
| `hisob/` | `ochko := 3` → chiqar → `ochko = ochko + 1` → yana chiqar: `3`, `4` |
| `savol/` | oʻrtadagi savol: `yosh := 15`, `yosh = 16` → faqat `16` |
| `nusxa/` | `rekord := ochko` nusxa oladi; `ochko = 0` dan keyin: `0 3` |
| `xato/` | **ataylab xato**: yorliq yaratilmagan, lekin `=` yozilgan → `undefined: ochko` |
| `xato2/` | **ataylab xato**: bor yorliqqa yana `:=` → `no new variables on left side of :=` |
| `uyga/yashil`, `uyga/sariq` | uy vazifasi uchun boshlangʻich fayllar |
| `javob/` | uy vazifasi javoblari |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/000-02-ozgaruvchi`) ishga tushiriladi.
```bash
go run hisob/main.go     # 3 va 4
go run savol/main.go      # 16
go run nusxa/main.go      # 0 3
go run xato/main.go       # xato/main.go:6:2: undefined: ochko  (va 7:14 da yana bir marta)
go run xato2/main.go      # xato2/main.go:7:8: no new variables on left side of :=
```
Natijalar `go1.26.8 linux/amd64` da, 2026-11-26.

## Dars matni (qisqacha)

**Oʻzgaruvchi.** Dastur ishlayotganda qiymatlar (son, matn) kompyuter xotirasida turadi. Qiymatga nom-yorliq
yopishtiramiz: bu **oʻzgaruvchi** (variable). Bitta yorliq ostida doim bitta qiymat turadi. Go'ning rasmiy qoidasi:
*"A variable is a storage location for holding a value"* ([spec: Variables](https://go.dev/ref/spec#Variables)) —
yorliq bir joyga yopishgan, `=` shu joydagi qiymatni almashtiradi.

Nega "quti" emas, "yorliq"? Hermans va boshq. (2018) 496 boshlovchida ikki oʻxshatishni solishtirgan: "quti"
oʻxshatishini olganlar oʻzgaruvchi bir nechta qiymatni saqlaydi degan xatoga koʻproq tushgan
([blog](https://www.felienne.com/archives/6063)).

**`:=` — yangi yorliq.** `ochko := 3` = yangi `ochko` oʻzgaruvchisi, qiymati 3. Toʻliq nomi: *short variable
declaration*; `var ochko = 3` ning qisqa shakli
([spec](https://go.dev/ref/spec#Short_variable_declarations)). Faqat funksiya ichida ishlaydi.

**`=` — qiymat berish (assignment).** Matematikadagi "teng" emas. `ochko = ochko + 1` ni "ochko endi ochko + 1
boʻlsin" deb oʻqing. Go uni ikki bosqichda bajaradi: avval oʻng tomon hisoblanadi (3 + 1 = 4), keyin natija chap
tomondagi oʻzgaruvchiga yoziladi; eski qiymat yoʻqoladi
([spec: Assignment statements](https://go.dev/ref/spec#Assignment_statements), *"The assignment proceeds in two
phases"*). Tenglikni tekshirish boshqa belgi: `==` (GO-0.04).

**Nusxa.** `rekord := ochko` — ochkoning *qiymatidan nusxa*. Keyin `ochko = 0` faqat ochkoni oʻzgartiradi:
`fmt.Println(ochko, rekord)` → `0 3`. `Println` vergul bilan berilgan qiymatlar orasiga boʻsh joy qoʻyadi
([fmt.Println](https://pkg.go.dev/fmt#Println)). (Bu sonlar uchun. Slice va map'larda nusxa ham bir xil
maʼlumotga qarab qoladi — GO-1.08, GO-1.09.)

**Ikki xato.**
- Yorliq hali yoʻq, lekin `=` → `undefined: ochko` ("ochko degan nom topilmadi"). Yechim: birinchi marta `:=`.
- Yorliq bor, lekin yana `:=` → `no new variables on left side of :=`. Yechim: `=`. (Chap tomonda kamida bitta
  *yangi* nom boʻlsa, `:=` ruxsat etiladi — spec'dagi "redeclaration" qoidasi, keyinroq kerak boʻladi.)
- Videoda yoʻq, lekin uy vazifasida duch kelishingiz mumkin: yaratilgan, lekin hech qayerda ishlatilmagan
  oʻzgaruvchi ham xato: `declared and not used: x`. Yechim: uni ishlating (masalan, `fmt.Println` ga bering) yoki oʻchiring.

**Qoʻshtirnoq.** `fmt.Println(ochko)` qiymatni (`3`) chiqaradi, `fmt.Println("ochko")` esa `ochko` soʻzini.
Son va matn farqi — keyingi dars.

## Uy vazifasi

🟢 **Oson** (`uyga/yashil/main.go`): `ochko` ni `yosh` deb qayta nomlang (hamma joyda), `3` oʻrniga oʻz yoshingizni
yozing. Dastur hozirgi va keyingi yilgi yoshingizni chiqarsin. Qoʻshimcha: `+ 1` ni `+ 5` ga almashtiring —
ekranda nima oʻzgaradi?

🟡 **Qiyinroq** (`uyga/sariq/main.go`): `a := 3`, `b := 8`. Izoh oʻrniga qatorlar yozib, `a` va `b` ning
qiymatlarini almashtiring, shunda dastur `8 3` chiqarsin. Maslahat: uchinchi yorliq kerak.

### Javoblar
- 🟢 `javob/yashil/main.go` → `15`, `16` (15 oʻrnida sizning yoshingiz).
- 🟡 `javob/sariq/main.go`: `vaqtincha := a`, `a = b`, `b = vaqtincha` → `8 3`. Nega `a = b; b = a` ishlamaydi?
  `a = b` dan keyin 3 yoʻqoladi, ikkala yorliq ostida 8 qoladi: `8 8`. (Go'da qisqa yoʻl ham bor:
  `a, b = b, a` — spec'ning ikki bosqichli qoidasi tufayli toʻgʻri ishlaydi.)

## Manbalar
- Go spec (go1.26): [Variables](https://go.dev/ref/spec#Variables),
  [Short variable declarations](https://go.dev/ref/spec#Short_variable_declarations),
  [Assignment statements](https://go.dev/ref/spec#Assignment_statements)
- [fmt.Println](https://pkg.go.dev/fmt#Println) · [A Tour of Go: Short variable declarations](https://go.dev/tour/basics/10)
- Hermans, Swidan va boshq. (2018), *Thinking out of the Box: Comparing Metaphors for Variables in Programming Education*
- Ma, Ferguson, Roper, Wood (2007), *Investigating the viability of mental models held by novice programmers*, SIGCSE
- [Harvard CS50x, lecture 1 notes](https://cs50.harvard.edu/x/notes/1/) (`counter = counter + 1`)
