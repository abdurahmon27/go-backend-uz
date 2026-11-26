# Go 0 · 04 · `if`: dastur qaror qabul qiladi

Video: *go 0: noldan boshlaymiz #4* · Lab: https://haywan.uz/lab/youtube/go#go-0-04 ·
Oldingi dars: [Go 0 · 03 · Son va matn: `2 + 2` bilan `"2" + "2"` nega har xil](../000-03-son-va-matn/) ·
Keyingi dars: Go 0 · 05 · `for`: kompyuter takrorlashdan charchamaydi

Hech qachon kod yozmaganlar uchun. Hamma natijalar `go1.26.8 linux/amd64`da olingan (2026-11-26).
Har bir dasturni https://go.dev/play da ham ishga tushirish mumkin: faylni nusxalab qoʻying va **Run** bosing.

**Bitta gʻoya:** dastur ikki sonni solishtirib savol beradi, javob faqat ha (`true`) yoki yoʻq (`false`) boʻladi,
`if` esa shu javobga qarab ikki yoʻldan faqat bittasini tanlaydi.

| Papka | Nima | Natija |
|---|---|---|
| `taqqos/` | `katta := xarid >= 100000` (`xarid := 50000`) | `false` |
| `chegirma/` | `if xarid >= 100000 { … } else { … }`, `xarid := 100000` | `chegirma bor` |
| `sinov/` | oʻrtadagi sinov: xuddi shu, lekin `>` | `chegirma yoʻq` |
| `baho/` | `if` / `else if` / `else`, `ball := 75` | `yaxshi` |
| `xato/` | `if javob = 4` — ataylab tarjima boʻlmaydi | `syntax error: cannot use assignment javob = 4 as value` |
| `tuzat/` | `if javob == 4` | `toʻgʻri` |
| `tartib/` | savollar tartibi xato: avval `>= 70`, keyin `>= 90`, `ball := 95` | `yaxshi` (aslida `aʼlo` boʻlishi kerak edi) |
| `uyga/` | uy vazifasi uchun boshlangʻich fayllar | |

## Taqqoslash belgilari

| Belgi | Oʻqilishi | Misol (`a := 5`) |
|---|---|---|
| `>` | katta | `a > 3` → `true` |
| `<` | kichik | `a < 3` → `false` |
| `>=` | katta yoki teng (≥) | `a >= 5` → `true` |
| `<=` | kichik yoki teng (≤) | `a <= 4` → `false` |
| `==` | tengmi? | `a == 5` → `true` |
| `!=` | teng emasmi? | `a != 5` → `false` |

Bitta teng `=` savol emas: u yorliq ostiga qiymat qoʻyadi (`a = 4` — "a boʻlsin 4", [oldingi dars](../000-02-ozgaruvchi/)).

## Videodagi buyruqlar
```bash
go run taqqos/main.go     # false
go run chegirma/main.go   # chegirma bor
go run sinov/main.go      # chegirma yoʻq
go run baho/main.go       # yaxshi
go run xato/main.go       # xato/main.go:7:11: syntax error: cannot use assignment javob = 4 as value
go run tuzat/main.go      # toʻgʻri
go run tartib/main.go     # yaxshi
```

## Xatoni qanday oʻqiymiz
```
xato/main.go:7:11: syntax error: cannot use assignment javob = 4 as value
```
- `xato/main.go:7:11` — fayl, **7-qator**, qatordagi **11-belgi**: aynan bitta teng `=` (qator boshidagi tab ham bitta belgi).
- `syntax error` — yozilish xatosi: kompilyator kodni oʻqiyolmadi.
- `cannot use assignment javob = 4 as value` — qiymat berishni (assignment) javob oʻrnida ishlatib boʻlmaydi.
  Go kompilyatorining oʻzida bu xato uchun alohida xabar bor, chunki `==` oʻrniga `=` yozish juda koʻp uchraydi
  ([parser.go](https://github.com/golang/go/blob/go1.26.8/src/cmd/compile/internal/syntax/parser.go): "A common syntax
  error is to write '=' instead of '=='").

## Uy vazifasi
- 🟢 **Oson** (`uyga/oson/main.go`): chegirma 100 000 soʻmdan emas, 50 000 soʻmdan boshlansin. Oʻzgartirishdan oldin
  taxmin qiling: hozir nima chiqadi, keyin-chi?
- 🟡 **Oʻrta** (`uyga/orta/main.go`): baho dasturiga 55 va undan yuqori ball uchun `qoniqarli` degan pogʻona qoʻshing
  (`ball := 60` → `qoniqarli`). Yangi savolni qayerga qoʻyish kerak?
- 🔴 **Qiyin** (`uyga/qiyin/main.go`): `!=` (teng emasmi?) bilan pin kodni tekshiring: teng boʻlmasa `notoʻgʻri kod`,
  teng boʻlsa `xush kelibsiz`.

<details><summary>Javoblar</summary>

🟢 Hozir: `chegirma yoʻq` (70 000 < 100 000). `if xarid >= 50000 {` qilsak: `chegirma bor`.

🟡 Yangi `else if` ni `>= 70` dan keyin, `else` dan oldin qoʻyamiz (eng qattiq savol birinchi):
```go
	if ball >= 90 {
		fmt.Println("aʼlo")
	} else if ball >= 70 {
		fmt.Println("yaxshi")
	} else if ball >= 55 {
		fmt.Println("qoniqarli")
	} else {
		fmt.Println("qayta topshiring")
	}
```
→ `qoniqarli`. Agar `>= 55` ni eng tepaga qoʻysangiz, 95 ball ham `qoniqarli` boʻlib qoladi (videodagi tartib xatosi).

🔴
```go
	pin := 2580
	kiritilgan := 1111
	if kiritilgan != pin {
		fmt.Println("notoʻgʻri kod")
	} else {
		fmt.Println("xush kelibsiz")
	}
```
→ `notoʻgʻri kod`. `kiritilgan := 2580` qilsangiz → `xush kelibsiz`.
</details>

## Chuqurroq (videoda yoʻq)
- Go'da savol atrofida oddiy qavs shart emas, jingalak qavslar esa majburiy: `if xarid >= 100000 {`
  ([A Tour of Go](https://go.dev/tour/flowcontrol/5)).
- `else if ball >= 70` ga `ball < 90` deb qoʻshish shart emas: bu savolga faqat oldingi `ball >= 90` javobi `false`
  boʻlganda navbat keladi (CS50 ham shu savolni beradi: "not all these lines of code are required").
- `else` majburiy emas: `tuzat/` dasturida faqat `if` bor — javob `false` boʻlsa, dastur hech narsa yozmaydi.
- Matnni son bilan taqqoslab boʻlmaydi — oldingi darsdagi xato yana chiqadi:
  `yosh := "15"; fmt.Println(yosh == 15)` → `invalid operation: yosh == 15 (mismatched types string and untyped int)`.
- Taqqoslash natijasining turi `bool`: `fmt.Printf("%T\n", 5 > 3)` → `bool`.

## Manbalar
- Go'ning rasmiy qoidalari (spec): [Comparison operators](https://go.dev/ref/spec#Comparison_operators) ("Comparison
  operators compare two operands and yield an untyped boolean value"), [Boolean types](https://go.dev/ref/spec#Boolean_types),
  [If statements](https://go.dev/ref/spec#If_statements).
- Harvard CS50x, Week 1 notes (conditionals, `==` vs `=`): https://cs50.harvard.edu/x/notes/1/
- University of Helsinki, Python MOOC, part 1 "Conditional statements": https://programming-24.mooc.fi/part-1/5-conditional-statements
- Go by Example: If/Else: https://gobyexample.com/if-else

## Dars matni

**intro**

- `hook` — Doʻkonda xaridingiz yuz ming soʻmga yetsa, kassa chegirmani oʻzi beradi, kassir hech narsani hisoblamaydi. Bugun xuddi shunday oʻzi qaror qabul qiladigan birinchi dasturimizni yozamiz.
- `eslash` — Avval oldingi darsni eslaymiz: fmt.Println qoʻshtirnoqdagi uch plyus qoʻshtirnoqdagi uchni ekranga chiqarsa, nima chiqadi? Oltimi, xatomi, yoki oʻttiz uchmi? Pauza qilib, oʻylab koʻring.
- `eslash-javob` — Javob: oʻttiz uch. Qoʻshtirnoqdagi uch bu matn, yaʼni string, plyus esa matnlarni yopishtiradi. Bugun esa ikkita sonni bir-biri bilan solishtiramiz.

**savol**

- `kassir` — Tasavvur qiling, kassa har bir xaridga bitta savol beradi: xarid yuz ming soʻmdan katta yoki tengmi? Javob faqat ikki xil: ha boʻlsa chegirma beriladi, yoʻq boʻlsa toʻliq narx olinadi.
- `taqqoslash` — Bunday savolda ikki son solishtiriladi. Buni oʻzbekcha taqqoslash, inglizcha comparison deyiladi. Go'da maktabdagi katta va kichik belgilari ishlatiladi. Klaviaturada katta yoki teng belgisi yoʻq, uni ikki belgi bilan yozamiz: katta, keyin teng.
- `true-false` — Taqqoslashning javobi ham qiymat, lekin Go uni inglizcha soʻz bilan aytadi. Ha oʻrniga true, yaʼni rost. Yoʻq oʻrniga false, yaʼni yolgʻon. Har bir taqqoslash shu ikkisidan birini beradi.
- `javob-xotira` — Endi xotirada: xarid yorligʻi yuz yigirma mingga yopishadi. Ikkinchi qatorda Go avval oʻng tomondagi savolni hisoblaydi, javob true, va katta degan yangi yorliq shu javobga yopishadi. katta ostidagi qiymatning turi bool, yaʼni ha-yoʻq turi. Uchinchi qator ekranga true chiqaradi.
- `bashorat1` — Endi birinchi qatorni oʻzgartiramiz: xarid ellik ming. Ekranga nima chiqadi: falseʼmi, trueʼmi, yoki ellik mingmi? Pauza qilib, oʻylab koʻring.
- `run1` — Javob: false. go run bilan taqqos papkasidagi dasturni ishga tushiramiz. Ellik ming yuz mingdan kichik, shuning uchun javob false. Qisqasi, taqqoslash savol beradi, uning javobi esa true yoki false.

**if**

- `agar` — Javobni bilish yetmaydi, dastur unga qarab ish qilishi kerak: agar xarid yuz mingdan katta yoki teng boʻlsa, chegirma ber, aks holda toʻliq narx ol. Go'da if soʻzi agar degani, else esa aks holda degani.
- `if-kod` — Endi kodda: birinchi qatorda xarid aynan yuz ming. Ikkinchi qatorda if, savol va jingalak qavs, yaʼni birinchi darsdagidek guruhning boshi. Guruh ichidagi qator ekranga chegirma bor deb yozadi, lekin faqat javob true boʻlsa bajariladi.
- `else-kod` — Endi ikkinchi guruhni qoʻshamiz. Guruh oxiridagi jingalak qavsdan keyin else va yangi guruh keladi. Undagi qator ekranga chegirma yoʻq deb yozadi, lekin faqat javob false boʻlsa bajariladi.
- `bashorat2` — Xarid aynan yuz ming. Ekranga nima chiqadi: chegirma bormi, yoʻqmi, yoki ikkalasimi? Pauza qilib, oʻylab koʻring.
- `run2` — Javob: chegirma bor. go run bilan chegirma papkasidagi dasturni ishga tushiramiz, terminal ham shuni yozdi.
- `chegirma-rasm` — Nega ekranga chegirma bor chiqdi? Yuz ming yuz mingga teng, katta yoki teng belgisi esa teng boʻlganda ham true beradi. Shuning uchun if guruhi bajarildi, else guruhi esa oʻtkazib yuborildi. Qisqasi, ikki guruhdan doim faqat bittasi bajariladi.

**sinov**

- `sinov` — Endi kichik sinov: ikkinchi qatordagi katta yoki teng oʻrniga faqat katta belgisini qoʻyamiz, xarid yana yuz ming. Ekranga endi nima chiqadi? Pauza qilib, oʻylab koʻring.
- `sinov-javob` — Javob: chegirma yoʻq. go run bilan sinov papkasidagi dasturni ishga tushirsak, shu chiqadi: yuz ming yuz mingdan katta emas, javob false, shuning uchun else guruhi bajarildi.

**else if**

- `baho` — Baʼzan javob ikki xil emas, uch xil boʻladi. Masalan, test bali toʻqson va undan yuqori boʻlsa aʼlo, yetmish va undan yuqori boʻlsa yaxshi, qolgan hollarda qayta topshiring. Savollarni ketma-ket qoʻyamiz: else if, yaʼni aks holda, agar. Oxirgi else hech bir javob true boʻlmasa bajariladi.
- `bashorat3` — Ball yetmish besh boʻlsa, ekranga nima chiqadi: yaxshimi, aʼlomi, yoki qayta topshiringmi? Pauza qilib, oʻylab koʻring.
- `run3` — Javob: yaxshi. go run bilan baho papkasidagi dasturni ishga tushirsak, terminal ham shuni yozadi.
- `baho-rasm` — Go savollarni tepadan pastga bittalab tekshiradi. Yetmish besh toʻqsondan katta yoki tengmi? Javob false, shuning uchun keyingi savolga oʻtadi. Yetmish besh yetmishdan katta yoki tengmi? Javob true, ekranga yaxshi chiqadi. Qisqasi, birinchi true topilgach, qolgani umuman tekshirilmaydi.

**xatolar**

- `tengmi` — Oldingi darsda bitta teng boʻlsin degani edi: u yorliq ostiga qiymat qoʻyadi, savol bermaydi. Ikki son tengmi deb soʻrash uchun Go'da ikkita teng yoziladi. Teskari savol ham bor: undov belgisi va teng, yaʼni teng emasmi.
- `xato-kod` — Bu dasturda oʻquvchining javobi toʻrt, toʻrt boʻlsa ekranga toʻgʻri chiqishi kerak. Eslang: ifʼdan keyin savol turadi. Lekin yettinchi qatorda adashib bitta teng yozilgan, u esa savol emas.
- `xato-savol` — Bu dastur nima qiladi: toʻgʻri deb yozadimi, jim qoladimi, yoki xato beradimi? Pauza qilib, oʻylab koʻring.
- `xato-run` — Javob: xato. go run bilan ishga tushirsak, kompilyator dasturni tarjima qilmaydi. Xabarni oʻqiymiz: yettinchi qator, oʻn birinchi belgi, yaʼni aynan bitta teng. Inglizcha qismida yozilish xatosi deyilgan: assignment, yaʼni qiymat berishni, savol oʻrnida ishlatib boʻlmaydi. Tuzatish oson: ikkita teng yozsak, dastur toʻgʻri deb yozadi.
- `tartib-kod` — Ikkinchi koʻp uchraydigan xato savollarning tartibida boʻladi. Bu xuddi baho dasturi, faqat ikki savol oʻrin almashgan: avval yetmish, keyin toʻqson soʻraladi. Ball esa toʻqson besh.
- `tartib-savol` — Toʻqson besh ball uchun ekranga nima chiqadi: aʼlomi, yaxshimi, yoki ikkalasimi? Pauza qilib, oʻylab koʻring.
- `tartib-run` — Javob: yaxshi. go run bilan tartib papkasidagi dasturni ishga tushirsak, terminal ham shuni yozadi.
- `tartib-nega` — Toʻqson besh ball aʼlo boʻlishi kerak edi, nega dastur yaxshi dedi? Go else if savollarini tepadan pastga tekshiradi va birinchi true, yaʼni ha javobida toʻxtaydi. Toʻqson besh yetmishdan katta yoki tengmi? Javob true, ekranga yaxshi chiqadi va aʼlo savoliga navbat kelmaydi. Tuzatish: eng qattiq savolni, yaʼni toʻqsonni, birinchi qoʻying.

**takrorlash**

- `takror` — Xulosani savollar bilan qilamiz. Birinchi: besh kichik uch, javob nima? Ikkinchi: bitta teng va ikkita teng farqi nima? Uchinchi: else guruhi qachon bajariladi? Birinchi darsdan esa: xato xabaridagi birinchi son nimani bildiradi? Pauza qilib, javob bering.
- `takror-javob` — Endi javoblar: birinchisi false, chunki besh uchdan kichik emas. Ikkinchisi: bitta teng qiymat qoʻyadi, ikkita teng tengmi deb soʻraydi. Uchinchisi: else faqat javob false boʻlganda bajariladi. Toʻrtinchisi esa qator raqami.
- `uyga` — Uyga vazifaning osonida chegirma ellik mingdan boshlansin. Oʻrtachasida baho dasturiga qoniqarli degan yana bir daraja qoʻshing. Qiyinida bank kartasining pin kodini teng emasmi belgisi bilan tekshiring. Fayllar va javoblar tavsifdagi havolada.
- `outro` — Keyingi darsda for soʻzi bilan dastur bitta ishni koʻp marta takrorlaydi. Koʻrishguncha.

