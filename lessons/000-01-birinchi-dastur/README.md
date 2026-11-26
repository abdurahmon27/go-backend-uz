# Go 0 · 01 · Dastur nima? Birinchi dasturimizni qatorma-qator oʻqiymiz

Video: *Golang oʻzbek tilida noldan · Go 0 · 1* · Lab: https://haywan.uz/lab/youtube/go#go-0-01 ·
Keyingi dars: Go 0 · 02 · Oʻzgaruvchi — qiymatga yopishtirilgan yorliq ·
Go'ni allaqachon bilsangiz: [01 · Go nega bor? `go build` ichida nima bor](../001-why-go/)

Hech qachon kod yozmagan boʻlsangiz, shu yerdan boshlang. Natijalar `go1.26.8 linux/amd64`da olingan
(Go Playground'da `go1.27.4`, matnlar bir xil).

| Fayl | Nima |
|---|---|
| `main.go` | videodagi dastur (7 qator): ekranga `Salom, dunyo!` yozadi |
| `xato/main.go` | ataylab xato: `fmt.println` (kichik p) — ishga tushmaydi |
| `homework/oson/main.go` | uy vazifasi uchun boshlangʻich fayl (videodagi kodning nusxasi) |
| `homework/javoblar/` | uy vazifasi javoblari |

## Videodagi buyruqlar

```bash
go version                  # go version go1.26.8 linux/amd64 (sizda yangiroq boʻlishi mumkin)
go run main.go              # Salom, dunyo!
cd xato && go run main.go   # ./main.go:6:6: undefined: fmt.println (but have Println)
```

Brauzerda: https://go.dev/play — tayyor misolni oʻchiring, `main.go`dagi 7 qatorni yozing, **Run** tugmasini bosing.
Pastda `Salom, dunyo!` va `Program exited.` chiqadi.

## Go'ni oʻrnatish (batafsil)

Rasmiy yoʻriqnoma: https://go.dev/doc/install · yuklab olish: https://go.dev/dl

- **Windows:** go.dev/dl dan `.msi` faylni yuklab oling, oching va koʻrsatmalar boʻyicha **Next** bosib boring.
  Keyin ochiq terminallarni yoping. Start menyusida `cmd` deb qidiring, ochilgan oynada `go version` yozing.
- **Mac:** `.pkg` faylni yuklab oling (Apple chipli kompyuterlarda `darwin-arm64`, Intel'da `darwin-amd64`),
  oching va koʻrsatmalarga amal qiling. Terminal ilovasini qaytadan oching, `go version` yozing.
- **Linux:** `.tar.gz` arxivni yuklab oling va rasmiy sahifadagi buyruqlarni bajaring:
  `sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz`, keyin
  `export PATH=$PATH:/usr/local/go/bin` qatorini `~/.profile` fayliga qoʻshing, terminalni qayta oching.

Windows'da oʻrnatuvchi ruxsat soʻrasa («Bu ilovaga oʻzgarish kiritishga ruxsat berasizmi?») — **Ha** bosing.

`go version` javobidagi raqam — versiya (sizda yangiroq boʻlishi normal), oxirgi soʻz — kompyuteringiz turi
(`windows/amd64`, `darwin/arm64`, `linux/amd64`…). `go version` oʻrniga `command not found` yoki
`'go' is not recognized` chiqsa: terminalni (yoki kompyuterni) qayta oching — oʻrnatuvchi yangi sozlamani faqat yangi
oynalarga beradi.

**Terminalni ochish:** Windows — Start → `cmd`; Mac — Command + Probel, `Terminal` deb yozing; Linux (Ubuntu) —
Ctrl + Alt + T. Videodagi `❯` belgisi (Windows'da `C:\Users\…>`) — terminalning taklif belgisi, uni yozmaysiz.

## Kodni qayerga yozish va terminalni papkada ochish

1. Kompyuterda yangi papka oching, masalan `salom`.
2. Istalgan oddiy matn muharririda (eng qulayi bepul VS Code: https://code.visualstudio.com, yoki Windows'dagi
   Bloknot) yangi fayl oching, 7 qatorni yozing va papkaga `main.go` nomi bilan saqlang.
   - **Bloknot:** «Fayl turi»ni **Barcha fayllar** qiling, aks holda fayl `main.go.txt` boʻlib qoladi. Windows fayl
     oxirini yashiradi: Explorer → **Koʻrinish (View) → Fayl nomi kengaytmalari** belgisini qoʻying va tekshiring.
   - **Mac TextEdit:** Format → **Make Plain Text**, Edit → Substitutions → **Smart Quotes**ni oʻchiring. Aks holda
     qoʻshtirnoqlar `“ ”`ga aylanadi va Go `invalid character U+201C` deydi. VS Code bilan bu muammo yoʻq.
   - Qator boshidagi surilishni **Tab** bilan qoʻying. Toʻrtta boʻsh joy bilan yozsangiz ham dastur ishlaydi, faqat
     xato xabaridagi belgi raqami boshqacha chiqadi (videodagi xato `6:6` oʻrniga `6:9` boʻladi).
3. Terminalni shu papkada oching:
   - **Windows:** papkani Explorer'da oching, tepadagi manzil qatoriga `cmd` deb yozing va Enter bosing.
   - **Mac / Linux / Windows (umumiy usul):** terminalda `cd ` (oxirida boʻsh joy bilan) yozing, papkani sichqoncha
     bilan terminal oynasiga tortib tashlang va Enter bosing. `cd` — papkaga kirish buyrugʻi.
   - VS Code'da: papkani oching, menyudan **Terminal → New Terminal** — terminal oʻsha papkada ochiladi.
4. `go run main.go` yozing. `stat main.go: no such file or directory` (Mac/Linux) yoki «fayl topilmadi» maʼnosidagi
   xato (Windows) chiqsa — terminal boshqa papkada ochilgan yoki fayl nomi `main.go.txt` boʻlib qolgan.
   Windows'da xato xabarida yoʻl `./main.go` emas, `.\main.go` koʻrinishida chiqishi mumkin — maʼnosi bir xil.

## Uy vazifasi

- 🟢 **Oson.** `homework/oson/main.go`ni oching (yoki go.dev/play'ga yozing) va `Salom, dunyo!` oʻrniga oʻz
  ismingiz bilan salomlashing. Javob: `homework/javoblar/oson/main.go` — faqat qoʻshtirnoq ichidagi matn oʻzgaradi,
  qoʻshtirnoqlar qoladi.
- 🟡 **Oʻrta.** Ikkinchi `fmt.Println(...)` qatorini qoʻshib, ekranga ikki qator chiqaring. Javob:
  `homework/javoblar/orta/main.go`. Ikkala qator jingalak qavslar ichida boʻlishi kerak; kompyuter ularni tepadan
  pastga bajaradi, shuning uchun birinchi qator tepada chiqadi.
- 🔍 **Qiziqish uchun.** Oxirgi qoʻshtirnoqni oʻchiring: `fmt.Println("Salom, dunyo!)`. Javob (go1.26.8):
  ```
  ./main.go:6:29: newline in string
  ./main.go:6:29: syntax error: unexpected newline in argument list; possibly missing comma or )
  ```
  Fayl `main.go`, **6-qator**, 29-belgi — qator oxiri: matn tugamay turib qator tugadi («newline in string» = matn
  ichida qator tugadi). Ikkinchi qator shu bitta xatoning davomi: avval birinchi xabarni tuzating.
- ➕ **Yana bittasi.** `Println`ni `PRINTLN` deb yozing. Xato qaysi belgini koʻrsatadi? Javob:
  `./main.go:6:6: undefined: fmt.PRINTLN (but have Println)` — yana oltinchi belgi. Lekin bu safar u `P`, ya’ni toʻgʻri harf: Go notoʻgʻri harfni emas, notanish nomning *boshini* koʻrsatadi. Videodagi misolda nomning boshi va xato harf bir joyga tushgan.

## Manbalar

- Go'ni oʻrnatish: https://go.dev/doc/install · https://go.dev/dl
- Go Playground: https://go.dev/play
- Rasmiy boshlangʻich qoʻllanma: https://go.dev/doc/tutorial/getting-started
- Go qoidalari (spec): package, `func main`, katta-kichik harf — https://go.dev/ref/spec#Program_execution ·
  https://go.dev/ref/spec#Identifiers
- `fmt.Println`: https://pkg.go.dev/fmt#Println · `go run`: https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program
- Kompilyator = manba kod → mashina kodi: Harvard CS50x, lecture 1 — https://cs50.harvard.edu/x/notes/1/
- Xato xabaridagi ustun raqami baytlarda, 1 dan: https://pkg.go.dev/go/token#Position

## Dars matni

**Kirish.**
Telefoningizdagi kalkulyator ikki qoʻshuv ikki toʻrt boʻlishini qayerdan biladi? Kimdir unga nima qilishni qadamma-qadam yozib bergan. Bugun siz ham kompyuterga birinchi marta oʻzingiz yozib buyruq berasiz.
Boshlashdan oldin bitta savol. Siz hozirgacha kompyuter yoki telefonga qanday qilib buyruq bergansiz? Bir oz oʻylab koʻring.
Koʻpchilik sichqoncha bilan bosgan yoki barmogʻi bilan ekranga tekkan. Bugun esa buyruqni yozib beramiz, buning uchun uchta soʻz kerak: dastur, yaʼni retsept, kompilyator, yaʼni tarjimon, va terminal, yaʼni yozishma oynasi.

**Dastur va kompilyator.**
Kompyuterga har bir qadamni aytish kerak. Osh retseptini tasavvur qiling: piyozni toʻgʻrang, sabzini qovuring, guruchni soling. Kompyuter ham roʻyxatni tepadan pastga, qatorma-qator bajaradi. Kompyuter uchun buyruqlar roʻyxati dastur deyiladi, inglizcha program. Retseptni yozish bu kod, yaʼni dastur matnini yozish, oshni pishirish esa dasturni ishga tushirish. Bizning tilimiz Go deb ataladi.
Biz kodni odam oʻqiy oladigan soʻzlar bilan yozamiz, kompyuter esa faqat nol va birlardan iborat oʻz tilini tushunadi. Demak, oʻrtada tarjimon kerak. Kodni kompyuter tiliga tarjima qiladigan dastur compiler, oʻzbekcha kompilyator deyiladi. U juda talabchan: buyruqdagi bitta harf notoʻgʻri boʻlsa ham tarjima qilmaydi va xato qayerda ekanini aytadi.

**Kod qatorma-qator.**
Mana birinchi dasturimiz, u ekranga Salom, dunyo deb yozadi. Fayl nomi main.go, oxiridagi nuqta go Go kodi ekanini bildiradi. Birinchi soʻz package, yaʼni toʻplam: package main ishga tushadigan dastur degani, main esa inglizcha asosiy.
Ikkinchi qator boʻsh, uchinchisi import fmt. Import olib kelish degani, fmt esa format soʻzining qisqasi: Go bilan birga keladigan tayyor asboblar qutisi. Nomi qoʻshtirnoq ichida, chunki qoʻshtirnoq matnning boshi va oxirini belgilaydi.
Beshinchi qator: func main. Func soʻzi buyruqlar guruhini yaratadi, bu guruhning nomi main. Tepadagi qatorlar faylni tanishtiradi, buyruqlarimiz esa aynan shu guruhdan boshlanadi. Jingalak qavslar guruhning boshi va oxirini belgilaydi, oddiy qavslarni esa keyinroq koʻramiz.
Oltinchi qator eng muhimi. fmt nuqta Println, yaʼni fmt qutisidagi Println asbobi. Println inglizcha print line, bitta qator chiqar degani. Qavs ichida, qoʻshtirnoqlar orasida, ekranga chiqadigan matn turibdi. Qator boshidagi surilish Tab tugmasi bilan qoʻyiladi. Qisqasi, toʻrt qism bor: bu dastur, asboblar qutisi, boshlanish joyi va ekranga yozadigan buyruq.

**Brauzerda: go.dev/play.**
Endi ishga tushiramiz, eng oson yoʻli brauzerda. Goʻning rasmiy saytidagi go.dev/play sahifasini oching, hech narsa oʻrnatmaysiz. Bu Go Playground, yaʼni Go oʻyin maydonchasi. Tayyor misolni oʻchirib, yetti qatorimizni yozamiz va tepadagi Run, yaʼni ishga tushir tugmasini bosamiz. Natija pastda chiqadi.
Run bosilganda pastda nima chiqadi? A: Salom, dunyo, qoʻshtirnoqsiz. B: xuddi shu matn, qoʻshtirnoqlari bilan. C: butun oltinchi qator. Videoni pauza qiling va oʻylab koʻring.
Javob: A. Ekranda faqat Salom, dunyo chiqadi, qoʻshtirnoqlarsiz, chunki qoʻshtirnoq faqat matnning chegarasi edi. Program exited esa dastur ishini tugatdi degani.

**Kompyuterda: terminal va go run.**
Endi oʻz kompyuterimizda ishga tushiramiz. Kompyuterga yozib buyruq beradigan oyna bor, xuddi yozishmadek: siz bir qator yozib Enter bosasiz, kompyuter javobni pastga yozadi. Bu oyna terminal deyiladi. Masalan, go version Goʻning versiyasini aytadi, boshidagi belgini esa yozmaysiz.
Buning uchun avval Goʻni oʻrnatamiz. Birinchi: go.dev/dl sahifasidan faylni yuklab oling, Windows uchun msi, Mac uchun pkg. Ikkinchi: faylni oching va koʻrsatmalarga amal qiling. Uchinchi: Windows’da Start menyusida cmd deb qidirib, terminalni oching. Toʻrtinchi: go version deb yozing.
Endi yetti qatorni Bloknot yoki VS Code’da yozib, yangi papkaga main.go nomi bilan saqlang, Bloknot’da Barcha fayllar turini tanlang. Keyin papka manzili qatoriga cmd deb yozing, terminal shu papkada ochiladi. Mac uchun yoʻl README, yaʼni tushuntirish faylida, havola video tagida.
Endi go run main.go, yaʼni main.goʻni ishga tushir, deb yozamiz. Natija brauzerdagidek: Salom, dunyo. Fayl topilmadi degan xato chiqsa, terminal boshqa papkada ochilgan.
Tezkor savol: kompyuter bizning kodimizni tushunmaydi-ku, unda go run ichida nima boʻldi? Videoni pauza qiling va tanlang.
Javob: B. Go run avval kompilyator bilan kodni kompyuter tiliga tarjima qiladi, keyin uni darhol ishga tushiradi. Qisqasi, kompyuter bizning matnimizni emas, uning tarjimasini bajaradi.

**Ataylab xato.**
Endi ataylab xato qilamiz. Oltinchi qatordagi Println soʻzining katta P harfini kichik p qilamiz. Odam uchun bu deyarli bir xil soʻz.
Endi go run main.go nima qiladi? A: baribir Salom, dunyo chiqadi. B: Go xato xabarini beradi. C: Go xatoni oʻzi tuzatib, ishga tushiradi. Pauza qiling va oʻylab koʻring.
Javob: B. Kompilyator kichik p bilan yozilgan soʻzni tanimadi va tarjima qilmadi. Dastur umuman ishga tushmadi, ekranda faqat xato xabari. Sizda tepada sarlavha qatori ham boʻladi.
Xato xabari qoʻrqinchli koʻrinadi, lekin u xatoning aniq manzilini aytadi. Biz Println soʻzini kichik p bilan yozgandik. Birinchi boʻlak: nuqta slesh main.go, yaʼni shu papkadagi fayl. Birinchi olti, bu oltinchi qator. Ikkinchi olti esa oltinchi belgi: boshidagi Tab ham sanaladi, demak aynan kichik p. Keyin undefined, yaʼni Go bunday nomni bilmaydi. Qavs ichida maslahat: katta P bilan Println bor.
Harfni yana katta qilib, qayta ishga tushiramiz: Salom, dunyo qaytib keldi. Go katta va kichik harfni farqlaydi. Xato chiqsa, avval fayl va qatorni toping, keyin xabarni oʻqing, xato bir nechta boʻlsa, birinchisidan boshlang.

**Xulosa.**
Endi uchta savol bilan xulosa qilamiz: dastur nima, kompilyator nima qiladi va main.go, olti, olti nimani bildiradi? Pauza qiling va javob bering.
Javoblarni tekshiramiz. Dastur kompyuter tepadan pastga bajaradigan buyruqlar roʻyxati. Kompilyator kodni kompyuter tiliga tarjima qiladi, xato boʻlsa toʻxtaydi. Uchinchisi: main.go faylining oltinchi qatori, oltinchi belgisi.

**Uy vazifasi.**
Uy vazifasi uch qismli. Oson: Salom, dunyo oʻrniga oʻz ismingizni yozing, masalan Salom, Aziza. Oʻrta: ikkinchi Println bilan ikki qator chiqaring. Qiziqish uchun: oxirgi qoʻshtirnoqni oʻchiring va xato qatorini toping. Hammasi Playground’da ham ishlaydi, javoblar README faylida.
Bugun birinchi dasturingizni oʻqidingiz, ishga tushirdingiz va xatosini oʻqidingiz. Keyingi darsda dastur maʼlumotni qanday eslab qolishini koʻramiz, koʻrishguncha.
