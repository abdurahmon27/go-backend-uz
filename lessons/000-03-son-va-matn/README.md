# Go 0 · 03 · Son va matn: `2 + 2` bilan `"2" + "2"` nega har xil

Video: *go 0: noldan boshlaymiz #3* · Lab: https://haywan.uz/lab/youtube/go#go-0-03 ·
Oldingi dars: [Go 0 · 02 · Oʻzgaruvchi — qiymatga yopishtirilgan yorliq](../000-02-ozgaruvchi/) ·
Keyingi dars: Go 0 · 04 · `if`: dastur qaror qabul qiladi

Hech qachon kod yozmaganlar uchun. Hamma natijalar `go1.26.8 linux/amd64`da olingan (2026-11-26).
Har bir dasturni https://go.dev/play da ham ishga tushirish mumkin: faylni nusxalab qoʻying va **Run** bosing.

**Bitta gʻoya:** har bir qiymatning turi bor. Qoʻshtirnoq bilan yozilgan qiymat matn (`string`), qoʻshtirnoqsiz butun
son `int`. `+` faqat bir xil turdagi ikki qiymat bilan ishlaydi: sonlarni qoʻshadi, matnlarni yopishtiradi; turlarni
aralashtirsangiz, Go dasturni tarjima qilmaydi.

| Papka | Nima | Natija |
|---|---|---|
| `qoshish/` | `2 + 2` va `"2" + "2"` | `4`, `22` |
| `sinov/` | oʻrtadagi sinov: `soat := "10"`, `daqiqa := "5"`, `soat + daqiqa` | `105` |
| `aralash/` | `string + int` — ataylab tarjima boʻlmaydi | `invalid operation: gap + yosh (mismatched types string and int)` |
| `aylantir/` | `strconv.Itoa(yosh)` bilan tuzatilgan | `Yosh: 15` |
| `bolish/` | `7 / 2` (ikkalasi `int`) | `3` |
| `suv/` | `7.0 / 2.0` (ikkalasi `float64`) | `3.5` |
| `uyga/` | uy vazifasi uchun boshlangʻich fayllar | |

## Videodagi buyruqlar
```bash
go run qoshish/main.go    # 4 / 22
go run sinov/main.go      # 105
go run aralash/main.go    # aralash/main.go:8:14: invalid operation: gap + yosh (mismatched types string and int)
go run aylantir/main.go   # Yosh: 15
go run bolish/main.go     # 3
go run suv/main.go        # 3.5
```

## Xatoni qanday oʻqiymiz
```
aralash/main.go:8:14: invalid operation: gap + yosh (mismatched types string and int)
```
- `aralash/main.go:8:14` — fayl, **8-qator**, qatordagi **14-belgi**: `gap` soʻzi shu yerdan boshlanadi
  (qator boshidagi bitta tab ham bitta belgi).
- `invalid operation` — notoʻgʻri amal.
- `mismatched types string and int` — turlar mos emas: biri `string`, biri `int`.

## Uy vazifasi
- 🟢 **Oson** (`uyga/oson/main.go`): dastur `7 / 2` ni `3` deb chiqaradi. Faqat ikki qatorni oʻzgartirib, `3.5`
  chiqaradigan qiling.
- 🟡 **Oʻrta** (`uyga/orta/main.go`): dastur ikki qator chiqaradi: `Mening yoshim` va `15`. Uni bitta qatorga,
  `Mening yoshim 15` qilib chiqaring.
- 🔴 **Qiyin** (`uyga/qiyin/main.go`): `olma := 7` va `shisha := 2.0` ni boʻlib koʻring. Xatoni oʻqing va tuzating.

<details><summary>Javoblar</summary>

🟢 `olma := 7.0` va `bola := 2.0` (nuqta qoʻysak, ikkalasi `float64`) → `3.5`.
Diqqat: butun chiqadigan `float64` nuqtasiz koʻrinadi: `fmt.Println(6.0 / 2.0)` → `3`.

🟡
```go
package main

import "fmt"
import "strconv"

func main() {
	yosh := 15
	fmt.Println("Mening yoshim " + strconv.Itoa(yosh))
}
```
→ `Mening yoshim 15` (videodagi `"Yosh: "` kabi: matn oxiridagi boʻsh joy son bilan soʻz orasini ochadi).

🔴 Xato: `uyga/qiyin/main.go:8:14: invalid operation: olma / shisha (mismatched types int and float64)`.
`int` bilan `float64` ham har xil tur. Ikki tuzatish:
- videodagidek nuqta qoʻyish: `olma := 7.0` → `3.5`;
- yoki aniq aylantirish: `fmt.Println(float64(olma) / shisha)` → `3.5` (`float64(...)` butun sonni kasrli songa
  aylantiradi; Go buni oʻzi jim qilmaydi).
</details>

## Chuqurroq: qavsdagi son oʻzgaruvchi boʻlmasa
Videoda "ikki xil turdagi **oʻzgaruvchilar**" deymiz. Kodga toʻgʻridan-toʻgʻri yozilgan son (konstanta) moslashuvchan:
- `olma := 7; fmt.Println(olma / 2.0)` → `3` (xato yoʻq: `2.0` butun son sifatida ishlatildi, kasr baribir tashlandi);
- `olma / 2.5` → `2.5 (untyped float constant) truncated to int`;
- `"Yosh: " + 15` → `invalid operation: "Yosh: " + 15 (mismatched types untyped string and untyped int)`.
Qoida (Go spec, Operators): ikki tomondagi turlar bir xil boʻlishi kerak, "unless the operation involves shifts or
untyped constants".

## Manbalar
- Go spec (go1.26): Arithmetic operators — `+` "integers, floats, complex values, strings"; Integer operators —
  "x / y truncated towards zero"; String concatenation; Operators — "the operand types must be identical";
  Constants — standart tur `int`, `float64`, `string`. https://go.dev/ref/spec#Arithmetic_operators
- Go FAQ: Why does Go not provide implicit numeric conversions? https://go.dev/doc/faq#conversions
- MDN, Addition (+): `"2" + 2; // "22"` (jim aylantiradigan til misoli).
  https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Addition
- `strconv.Itoa`: https://pkg.go.dev/strconv#Itoa
- CS50x Week 1 notes (7 / 2 → 3): https://cs50.harvard.edu/x/notes/1/

## Dars matni
**hook.** Kalkulyatorda ikki plyus ikki har doim toʻrt boʻladi. Go dasturida esa xuddi shunga oʻxshash ikki plyus ikki yigirma ikki chiqarishi mumkin. Farq ekrandagi bitta kichik belgida yashiringan. Bugun oʻsha belgini topamiz.

**eslash.** Avval oldingi darsni eslaymiz. ochko, ikki nuqta teng, uch. Keyin ochko boʻlsin ochko plyus bir. Shundan keyin ekranga ochko nima chiqadi? Pauza qilib, oʻylab koʻring.

**eslash-javob.** Javob: toʻrt. Birinchi qatorda ochko yorligʻi uch qiymatiga yopishtirildi. Ikkinchi qatorda yorliq ostidagi qiymat toʻrtga almashdi, eskisi yoʻqoldi. Bugun yorliqlarga emas, qiymatlarning oʻziga qaraymiz.

**tur.** Tasavvur qiling, anketada ikki qator bor: yoshingiz va telefon raqamingiz. Ikkalasi ham raqamlardan iborat. Lekin yoshga bir qoʻshsangiz, keyingi yilgi yoshingiz chiqadi. Telefon raqami bilan esa hisob qilinmaydi, u faqat yoziladi.

**tur-soz.** Kompyuter ham har bir qiymatning xilini bilishi kerak: uni qoʻshadimi, yoki faqat yozadimi. Qiymatning xili dasturlashda tur, inglizcha type deyiladi. Go'da har bir qiymatning turi bor.

**int-string.** Go'da butun sonning asosiy turi int deyiladi: yosh ikki nuqta teng oʻn besh, bu int. Matn esa string deyiladi. Birinchi darsdagi Salom, dunyo kabi, kodda matnni qoʻshtirnoq ichida yozamiz. yozuv ikki nuqta teng qoʻshtirnoq ichida oʻn besh, bu string, xuddi telefon raqami kabi. Ikkalasida oʻn besh bor, lekin bu ikki xil qiymat.

**plyus-kod.** Endi plyusni ikkala turda sinaymiz. package main, import va func main birinchi darsdagi qolip, fmt.Println esa qavs ichidagini ekranga chiqaradi. Oltinchi qatorda qoʻshtirnoqsiz ikki plyus ikki, yettinchi qatorda esa qoʻshtirnoqdagi ikki plyus qoʻshtirnoqdagi ikki.

**bashorat1.** Bu dastur ekranga ikki qator chiqaradi. Toʻrt va toʻrtmi, toʻrt va yigirma ikkimi, yoki xatomi? Videoni pauza qilib, oʻylab koʻring.

**run1.** Javob: toʻrt va yigirma ikki. Boshidagi oʻsha kichik belgi mana shu qoʻshtirnoq edi. go run buyrugʻi qoshish papkasidagi main.go faylini tarjima qilib, darhol ishga tushiradi.

**plyus-nega.** Go'da ikki plyus ikki toʻrt beradi, qoʻshtirnoqdagi ikki plyus qoʻshtirnoqdagi ikki esa yigirma ikki. Sababi, qoʻshtirnoq bilan yozilgan qiymat matn boʻladi, plyus esa avval ikki tomonning turiga qaraydi. Ikkalasi int boʻlsa, ularni qoʻshadi. Ikkalasi string boʻlsa, yonma-yon yopishtiradi. Demak, bu yigirma ikki son emas, yonma-yon turgan ikkita ikki. Qisqasi, plyus faqat bir xil turdagi ikki qiymat bilan ishlaydi.

**sinov.** Endi kichik sinov: soat ikki nuqta teng qoʻshtirnoq ichida oʻn, daqiqa ikki nuqta teng qoʻshtirnoq ichida besh. Keyin ekranga soat plyus daqiqa chiqaramiz. Oʻn beshmi, yuz beshmi, yoki xatomi? Pauza qilib, oʻylab koʻring.

**sinov-javob.** Javob: bir, nol, besh. soat va daqiqa ostida qoʻshtirnoqli qiymat, yaʼni string bor, va plyus ularni yopishtirdi. Bu yuz besh soni emas, yonma-yon turgan uchta belgi.

**aralash.** Endi ikki turni aralashtiramiz. gap oʻzgaruvchisida qoʻshtirnoq ichida Yosh soʻzi, ikki nuqta va oxirida bitta boʻsh joy bor, bu string. yosh oʻzgaruvchisida qoʻshtirnoqsiz oʻn besh, yaʼni int. Sakkizinchi qatorda ularni plyus bilan birlashtirmoqchimiz.

**bashorat2.** Bu dastur ishga tushadimi? Esingizda boʻlsin, plyus faqat bir xil turdagi qiymatlar bilan ishlaydi, bu yerda esa string va int. Pauza qilib, oʻylab koʻring.

**aralash-run.** Javob: xato. go run bilan ishga tushirsak, kompilyator dasturni tarjima qilmadi, shuning uchun u ishga tushmadi. Sababi, plyusning ikki tomonidagi turlar har xil.

**xato-oqish.** Xabarni qismlarga boʻlib oʻqiymiz. Fayl nomidan keyingi birinchi son, sakkiz, bu qator raqami. Ikkinchisi, oʻn toʻrt, qatordagi belgi oʻrni: aynan gap soʻzi shu yerdan boshlanadi. Keyin inglizcha yozilgan: notoʻgʻri amal, qavs ichida sabab: turlar mos emas, biri string, biri int.

**nega-go.** Nega Go bunga ruxsat bermaydi? Baʼzi boshqa tillar sonni jimgina matnga aylantiradi. Unda doʻkon dasturida ikki ming plyus ikki ming toʻrt ming emas, ikki ming ikki ming degan matn boʻlib qolishi mumkin. Go esa turlarni tarjima paytida tekshiradi va xatoni dastur ishga tushmasdan oldin koʻrsatadi.

**aylantir.** Tuzatish uchun sonni matnga aylantiramiz. Go'da buning uchun strconv degan yana bir tayyor asboblar qutisi bor, uni ham fmt kabi import qilamiz. gap va yosh avvalgidek. strconv.Itoa qavs ichidagi butun sondan yangi matn yaratadi, yosh esa int boʻlib qoladi. Nima chiqadi, taxmin qiling?

**aylantir-run.** go run bilan ishga tushiramiz: ekranda Yosh, oʻn besh chiqdi. Endi plyusning ikki tomonida ham string bor. Qisqasi, turlar mos kelmasa Go xato beradi, strconv.Itoa esa sonni matnga aylantiradi.

**bolish.** Endi boʻlishga oʻtamiz: yetti olmani ikki bolaga teng boʻlamiz, lekin olmani kesmaymiz, har bola faqat butun olma oladi. olma yetti, bola ikki, ikkalasi ham int. Kodda qiya chiziq boʻlishni anglatadi.

**bashorat3.** Ekranga nima chiqadi: uch butun oʻndan beshmi, uchmi, yoki toʻrtmi? Ishora: int faqat butun son saqlaydi, olma esa kesilmaydi. Pauza qilib, oʻylab koʻring.

**bolish-run.** Javob: uch. Ikkala qiymat ham int, demak natija ham butun son. Go kasr qismini shunchaki tashlab yuboradi, yaxlitlamaydi. Hayotdagidek: har bir bolaga uchtadan olma tegadi, bittasi esa ortib qoladi.

**suv.** Agar bizga uch yarim kerak boʻlsa-chi? Yetti litr suvni ikki shishaga quyamiz. Kodda kasr son vergul bilan emas, nuqta bilan yoziladi. Nuqtali son Go'da float64 deyiladi. Nima chiqadi, taxmin qiling?

**suv-run.** go run bilan ishga tushiramiz: uch butun oʻndan besh. Turi float64 boʻlgani uchun kasr qismi saqlanib qoldi.

**ikki-son.** Rasmda ikkalasi yonma-yon: intʼda kasr tashlanadi, float64ʼda esa qoladi. int va float64 oʻzgaruvchilarni aralashtirsangiz ham, xuddi shu turlar mos emas xatosi chiqadi.

**takror.** Endi xulosa qilamiz, lekin savollar bilan. Birinchi savol: qoʻshtirnoqdagi oʻn besh qaysi turda? Ikkinchi savol: string bilan intʼni plyus bilan birlashtirsak nima boʻladi? Uchinchi savol: ikkalasi int boʻlsa, yetti boʻlish ikki nechchi? Pauza qilib, javob bering.

**takror-javob.** Javoblarni tekshiramiz: birinchisi string, chunki qoʻshtirnoq bor. Ikkinchisi turlar mos emas xatosi, tuzatish uchun strconv.Itoa ishlatamiz. Uchinchisi uch, chunki kasr qismi tashlanadi.

**uyga.** Uyga vazifaning oson qismi: uyga papkasidagi dastur uch chiqaradi, uni uch butun oʻndan besh chiqaradigan qiling. Oʻrta: Mening yoshim degan matn va yoshingizni bitta qatorda chiqaring. Qiyin: intʼni float64ʼga boʻlib, xatoni tuzating. Fayllar va javoblar tavsifdagi havolada.

**outro.** Keyingi darsda sonlarni taqqoslaymiz, va dastur if yordamida oʻzi qaror qabul qiladi. Koʻrishguncha.
