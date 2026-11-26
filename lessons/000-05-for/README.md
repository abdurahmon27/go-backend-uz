# Go 0 · 05 · `for`: kompyuter takrorlashdan charchamaydi

Video: *go 0: noldan boshlaymiz #5* · Lab: https://haywan.uz/lab/youtube/go#go-0-05 ·
Oldingi dars: [Go 0 · 04 · `if`: dastur qaror qabul qiladi](../000-04-if/) ·
Keyingi qadam: [Go 1 · 01 · Go nega bor? `go build` ichida nima bor](../001-why-go/)

Hech qachon kod yozmaganlar uchun. Go 0 mavsumining oxirgi darsi. Hamma natijalar `go1.26.8 linux/amd64`da olingan
(2026-11-26). Har bir dasturni https://go.dev/play da ham ishga tushirish mumkin: faylni nusxalab qoʻying va **Run** bosing.

**Bitta gʻoya:** sikl bir ishni qayta-qayta bajaradi: sanagich bir marta boshlanadi, har aylana oldidan savol
beriladi, har aylanadan keyin qadam sanagichni oshiradi. Qadam unutilsa sikl tugamaydi va dasturni **Ctrl+C** bilan
toʻxtatamiz; sikldan ichkaridan esa **`break`** bilan chiqamiz.

| Papka | Nima | Natija |
|---|---|---|
| `sanash/` | `for i := 1; i <= 3; i++` | `1` `2` `3` `tamom` |
| `salom/` | oʻrtadagi sinov: `i <= 4` → nechta `Salom`? | `Salom` × 4 |
| `yigindi/` | `jami = jami + i`, i: 1 … 5 | `15` |
| `cheksiz/` | ataylab xato: qadam (`i++`) unutilgan — hech qachon toʻxtamaydi | `1 1 1 …` → **Ctrl+C** |
| `loyiha/` | Go 0 loyihasi: telefonga necha oyda pul yigʻiladi (ming soʻmda) | `1 oy: 200` … `5 oy: 1000` |
| `uyga/` | uy vazifasi uchun boshlangʻich fayllar | |
| `javob/` | uy vazifasi javoblari | |

## Videodagi buyruqlar
```bash
go run sanash/main.go     # 1 2 3 tamom (har biri alohida qatorda)
go run salom/main.go      # Salom × 4
go run yigindi/main.go    # 15
go run cheksiz/main.go    # 1 1 1 … toʻxtamaydi → Ctrl+C bosing
go run loyiha/main.go     # 1 oy: 200 … 5 oy: 1000
```

## `for` qatorining uch qismi
```go
for i := 1; i <= 3; i++ {
	fmt.Println(i)
}
```
| Qism | Kodda | Qachon ishlaydi |
|---|---|---|
| boshlanish | `i := 1` | **bir marta**, eng boshida |
| savol | `i <= 3` | **har aylanadan oldin**; `true` boʻlsa jingalak qavslar ichi bajariladi, `false` boʻlsa sikl tugaydi |
| qadam | `i++` (= `i = i + 1`) | **har aylanadan keyin** |

Go'ning rasmiy qoidasi: "the init statement is executed once before evaluating the condition for the first
iteration; the post statement is executed after each execution of the block"
([For statements](https://go.dev/ref/spec#For_statements)). Go'da boshqa sikl soʻzi yoʻq: *"Go has only one looping
construct, the `for` loop"* ([A Tour of Go](https://go.dev/tour/flowcontrol/1)).

Tajribali dasturchilar odatda noldan sanaydi: `for i := 0; i < 3; i++` ham uch marta aylanadi (0, 1, 2). Noldan
sanash roʻyxatlar (slice) bilan ishlaganda kerak boʻladi, uni Go 1 mavsumida koʻramiz.

## Cheksiz sikl, Ctrl+C va `break`
- `cheksiz/main.go` — `for i := 1; i <= 3; {`: qadam yoʻq. Go buni qabul qiladi (for qatorining istalgan qismi boʻsh
  boʻlishi mumkin), lekin `i` hech qachon oshmaydi, `1 <= 3` doim `true` va sikl toʻxtamaydi. Tuzatish: `i++` ni
  qaytaring.
- `for {` — umuman qism yozilmasa ham sikl cheksiz: "If the condition is absent, it is equivalent to the boolean value
  true" ([For statements](https://go.dev/ref/spec#For_statements)). Bunday sikldan faqat `break` bilan chiqiladi.
- **Ctrl+C** (Mac'da ham Ctrl, Cmd emas) — terminal ishlayotgan dasturga *toʻxta* ishorasini (SIGINT) yuboradi. Bu
  kompyuterda bir soniyada 890 000 dan ortiq `1` chiqdi (3 marta oʻlchandi), keyin `go run` quyidagini yozdi:
  ```
  1
  1
  ^Csignal: interrupt
  ```
  `^C` — terminal bosilgan tugmani koʻrsatishi. Windows'da oxirgi qator boshqacha boʻlishi mumkin (bu kompyuterda
  tekshirilmagan). Ekranga chizadigan oddiy terminal bundan sekinroq, lekin baribir soniyasiga minglab qator.
- **go.dev/play**da Ctrl+C kerak emas: server dasturni bir necha soniyadan keyin oʻzi toʻxtatadi va
  `timeout running program` deb yozadi (2026-11-26 da tekshirildi).
- **`break`** — siklni ichkaridan, darhol tugatadi, hatto `if` guruhi ichida boʻlsa ham butun `for` tugaydi va dastur
  sikldan keyingi qatorga oʻtadi ([Break statements](https://go.dev/ref/spec#Break_statements)).

## Uy vazifasi
- 🟢 **Oson** (`uyga/oson/main.go`): bu videodagi loyiha. Har oy `200` emas, `300` yigʻing. Ishga tushirishdan oldin
  taxmin qiling: necha oyda `900` ga yetadi?
- 🟡 **Oʻrta** (`uyga/orta/main.go`): 7 ning koʻpaytirish jadvalini sikl bilan chiqaring: `7 x 1 = 7` … `7 x 10 = 70`.
  Kodda yulduzcha `*` koʻpaytiradi: `7 * 3` → `21`. Fayldagi bitta qatorni sikl ichiga oling.
- 🔴 **Qiyin** (`uyga/qiyin/main.go`): bu dastur `1 + 2 + 3 + 4 + 5` ni hisoblashi kerak edi, lekin `5` chiqaradi.
  Yashirin xatoni toping. Maslahat: har aylanada `jami` yorligʻi ostida nima boʻlayotganini qogʻozga chizing.

<details><summary>Javoblar</summary>

🟢 `jami = jami + 300` → `1 oy: 300`, `2 oy: 600`, `3 oy: 900`. **3 oyda**: `900 >= 900` — *katta yoki teng*, teng
boʻlgani uchun ham `true`. Agar `>` yozilganida, 4-oyda (`1200`) toʻxtardi. Fayl: `javob/oson/main.go`.

🟡
```go
for i := 1; i <= 10; i++ {
	fmt.Println(7, "x", i, "=", 7*i)
}
```
Fayl: `javob/orta/main.go`. Println vergul bilan berilgan qiymatlar orasiga oʻzi boʻsh joy qoʻyadi.

🔴 `jami = 0` sikl **ichida**: har aylanada savat boʻshatiladi va faqat oxirgi `i` (5) qoladi. Bu qatorni oʻchiring
(`jami := 0` sikldan oldin bir marta yetarli) → `15`. Fayl: `javob/qiyin/main.go`. Bu yangi boshlovchilarning eng
koʻp uchraydigan sikl xatolaridan biri: boshlangʻich qiymat sikldan **oldin** beriladi.
</details>

## Manbalar
- Go spec (go1.26): [For statements](https://go.dev/ref/spec#For_statements) ·
  [IncDec statements](https://go.dev/ref/spec#IncDec_statements) ("increment … by the untyped constant 1") ·
  [Break statements](https://go.dev/ref/spec#Break_statements)
- `fmt.Println`: "Spaces are always added between operands and a newline is appended." https://pkg.go.dev/fmt#Println
- A Tour of Go, For / For is Go's "while" / Forever: https://go.dev/tour/flowcontrol/1 (manba matni:
  https://github.com/golang/website/blob/master/_content/tour/flowcontrol.article)
- Harvard CS50x lecture 1 notes (meow × 3, `while (true)`, control-C): https://cs50.harvard.edu/x/notes/1/
- os/signal: "A SIGHUP, SIGINT, or SIGTERM signal causes the program to exit." https://pkg.go.dev/os/signal
- POSIX terminal: INTR (Ctrl+C) → SIGINT: https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap11.html#tag_11_01_09
- MIT 6.100L finger exercise ("hello world" N marta): https://ocw.mit.edu/courses/6-100l-introduction-to-cs-and-programming-using-python-fall-2022/pages/finger-exercises/

## Dars matni

**hook.** Doskaga bir gapni yuz marta yozish jazosini tasavvur qiling: qoʻlingiz tez charchaydi. Kompyuter esa million marta yozsa ham charchamaydi. Bugun unga buni qanday buyurishni oʻrganamiz.

**eslash.** Avval oldingi darsni eslaymiz: harorat ikki nuqta teng oʻttiz. Agar harorat yigirma beshdan katta boʻlsa issiq, aks holda salqin yoziladi. Ekranga qaysi soʻz chiqadi? Pauza qilib, oʻylab koʻring.

**eslash-javob.** Javob: issiq. harorat ostida oʻttiz turibdi, savol: oʻttiz yigirma beshdan kattami? Javob true, shuning uchun if guruhi bajarildi. Lekin if savolni bir marta beradi, bugun esa savolni qayta-qayta beradigan buyruqni oʻrganamiz.

**muammo.** Deylik, birdan beshgacha sonlarni chiqarish kerak: hozircha har biriga alohida Println yozamiz. Million son kerak boʻlsa-chi? Uni hech kim qoʻlda yozmaydi.

**aylana.** Hayotdan misol: murabbiy uch aylana yuguring deydi. Har aylana oldidan soʻraysiz: aylana raqami uchdan kichik yoki tengmi? Ha boʻlsa, yugurasiz va sanoqqa bir qoʻshasiz. Bir ishni shunday qayta-qayta bajarish oʻzbekcha sikl, inglizcha loop deyiladi.

**sanagich.** Aylanalarni barmoq bukib sanaysiz. Dasturda buni alohida oʻzgaruvchi qiladi: u oʻzbekcha sanagich, inglizcha counter deyiladi, nomi odatda i. Sanagich ham oddiy yorliq: har aylanada uning ostidagi qiymat bittaga oshadi.

**for-qator.** Endi Goʻda yozamiz, ekranda tanish qolipning faqat ichi. Sikl for soʻzi bilan boshlanadi, keyin nuqtali vergul bilan ajratilgan uch qism keladi. Birinchisi boshlanish: i ikki nuqta teng bir. Ikkinchisi savol: i kichik yoki teng uchmi? true boʻlsa, sikl davom etadi.

**qadam.** Uchinchi qism qadam: i plyus plyus, yaʼni i ni bittaga oshir, xuddi ikkinchi darsdagi ochko boʻlsin ochko plyus bir kabi. Jingalak qavslar ichidagi Println har aylanada bajariladi, sikl tugagach esa tamom yoziladi.

**bashorat1.** Taxmin qiling, ekranda nima chiqadi: bir, ikki, uch va tamommi, bir, ikki va tamommi, yoki noldan uchgacha va tamommi? Pauza qilib, oʻylab koʻring.

**run1.** Javob: bir, ikki, uch va tamom. go run bilan sanash papkasidagi dasturni ishga tushiramiz: har bir son alohida qatorda chiqdi.

**sikl-xotira.** Endi xotirada nima boʻldi? Avval i birga teng, javob true, ekranda bir. Keyin i ikkiga oshadi, ekranda ikki. Soʻng uchga, ekranda uch. Endi i toʻrtga oshadi, javob false, sikl tugaydi va tamom yoziladi. Qisqasi: boshlanish bir marta, savol va qadam esa har aylanada ishlaydi.

**savol.** Endi oʻzingizni tekshiring: i birdan boshlanadi, savol esa i kichik yoki teng toʻrtmi? Ichida Salom soʻzi chiqariladi. Ekranda nechta Salom boʻladi: uchtami, toʻrttami yoki beshtami?

**savol-javob.** Javob: toʻrtta, go run ham shuni koʻrsatadi. i bir, ikki, uch va toʻrt boʻlganda javob true, beshda esa false.

**yigindi-kod.** Endi sikl bilan hisob qilamiz: birdan beshgacha hamma sonlarning yigʻindisini topamiz. Birinchi qatorda, sikldan oldin, jami nolga teng. Sikl ichida jami boʻlsin jami plyus i: har aylanada navbatdagi son qoʻshiladi. Sikl tugagach, jamini chiqaramiz.

**bashorat2.** Ekranga nima chiqadi: beshmi, oʻn beshmi, yoki har aylanadagi jamining hammasimi? Ishora: Println jingalak qavslarning ichidami yoki tashqarisidami? Pauza qilib, oʻylab koʻring.

**yigindi-run.** Javob: oʻn besh, faqat bitta qator, go run ham shuni koʻrsatadi. Println jingalak qavslardan tashqarida, shuning uchun sikl tugagach bir marta ishlaydi.

**yigindi-xotira.** Kompyuter birdan beshgacha qanday qoʻshadi? Yorliqlar: jami va sanagich i. Avval jami nolga teng. i bir boʻlganda jami ham bir. i ikki boʻlganda jami uch. Keyin jami olti, oʻn va oʻn besh boʻladi. i olti boʻlganda javob false, sikl tugaydi va ekranga oʻn besh chiqadi. Qisqasi: jami yorligʻi joyida qoldi, faqat uning qiymati almashdi.

**cheksiz-kod.** Endi sikldagi eng koʻp uchraydigan xato. Bu sanash dasturining oʻzi, faqat for qatorida qadam unutilgan: i plyus plyus yoʻq. Go bunday kodni qabul qiladi, tarjima xatosi boʻlmaydi.

**bashorat3.** Bu dastur nima chiqaradi: bir, ikki, uch va tamommi, birni toʻxtovsiz chiqaradimi, yoki hech narsa chiqarmaydimi? Ishora: i qachon oshadi? Pauza qilib, oʻylab koʻring.

**cheksiz-run.** Javob: birni toʻxtovsiz chiqaradi. i oshmaydi, shuning uchun javob doim true. Bunday sikl cheksiz sikl deyiladi. go run bilan ishga tushirsak, bir soniyada yuz minglab bir chiqadi. Toʻxtatish uchun Ctrl tugmasini bosib turib, C ni bosamiz: oxirgi yozuvlar shuni bildiradi.

**tuzat.** Xotirada i ostida doim bir turibdi, chunki qadam yoʻq. Javob har safar true, ekran esa birlar bilan toʻladi. Tuzatish oson: qadamni qaytaramiz va dastur bir, ikki, uch, tamom chiqarib, oʻzi toʻxtaydi.

**loyiha.** Endi Go nol darslarini bitta loyihaga yigʻamiz. Telefon toʻqqiz yuz ming soʻm, siz esa har oy ikki yuz ming soʻm yigʻasiz. Necha oyda pul yetadi? Sikl bir yilga yoziladi, lekin pul yetgan oyda toʻxtashi kerak.

**break.** Buning uchun sikldan chiqish qoidasini yozamiz. Yugurishdan oldin oʻzingizga qoida qoʻyasiz: suv tugasa, aylanalar tugamasa ham yoʻlakdan chiqasiz. Sikldan shunday ichkaridan chiqish break deyiladi.

**loyiha-kod.** Kod yigʻindi dasturiga oʻxshaydi, pul esa ming soʻmda. Sanagich oy oʻn ikkigacha, yaʼni bir yilgacha sanaydi. Har aylanada jamiga ikki yuz qoʻshiladi. Println vergul bilan oy raqamini, matnni va jamini bitta qatorga chiqaradi: bu plyus emas, son va matn yonma-yon tura oladi.

**loyiha-if.** Endi qoida: agar jami toʻqqiz yuzdan katta yoki teng boʻlsa, break ishlaydi. break if guruhi ichida boʻlsa ham, butun siklni tugatadi.

**bashorat4.** Ekrandagi oxirgi qator qaysi: toʻrt oy sakkiz yuzmi, besh oy mingmi, yoki oʻn ikki oymi? Pauza qilib, oʻylab koʻring.

**loyiha-run.** Javob: besh oy. go run koʻrsatadi: toʻrtinchi oyda jami hali yetmaydi, beshinchisida esa ming, yaʼni million soʻm boʻldi va break siklni toʻxtatdi. Println ifʼdan oldin turgani uchun ming ham chiqdi. Qisqasi, bu yerda oʻzgaruvchi, son va matn, if va for birga ishladi.

**takror.** Uchta savol bilan qaytaramiz. Birinchi: for qatoridagi uch qism nima? Ikkinchi: qadam unutilsa nima boʻladi? Uchinchisi ikkinchi darsdan: jami boʻlsin jami plyus i qanday bajariladi? Pauza qilib, javob bering.

**takror-javob.** Endi javoblarni tekshiramiz. Birinchisi: boshlanish, savol va qadam. Ikkinchisi: sikl cheksiz boʻladi, uni Ctrl C toʻxtatadi, qadamni qaytarib tuzatamiz. Uchinchisi: avval oʻng tomon hisoblanadi, keyin natija jami yorligʻi ostiga qoʻyiladi.

**uyga.** Endi uyga vazifa beraman. Oson: loyihada har oy uch yuz yigʻing, necha oyda yetadi? Oʻrta: yettining koʻpaytirish jadvalini sikl bilan chiqaring, yulduzcha belgisi koʻpaytiradi. Qiyin: yigʻindi dasturining buzilgan nusxasini tuzating. Fayllar va javoblar tavsifdagi havolada.

**outro.** Tabriklayman, Go nol mavsumi tugadi! Keyingi qadam: Go asoslari mavsumining birinchi darsi, Go nega bor va go build ichida nima bor. Hujjatlar tavsifdagi lab havolasida. Koʻrishguncha.
