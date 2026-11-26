# CL-1.03 · Prompt — bu texnik topshiriq: kontekst, misol, cheklov

Oldingi dars: [CL-1.02 · Claude.ai'dan toʻgʻri foydalanish: model tanlash, Projects va fayllar](../cl-02-claude-ai/) ·
Keyingi dars: CL-1.04 · AI — ustoz, javob mashinasi emas: bilim illyuziyasisiz oʻrganish

Bitta Go xatosini (Go kursi, [GO-1.14](../014-errors-values/): `%w` bilan oʻralgan xato `==` bilan solishtiriladi)
Claude Sonnet 5.5 ga yettita prompt bilan berdik, har birini **uch marta**. Javoblar **oʻzgartirilmasdan** saqlangan,
Go dasturi esa ularni oʻlchaydi va Claude bergan tuzatish va testlarni **ishga tushirib** tekshiradi.

| Fayl / papka | Nima |
|---|---|
| `bozor/main.go` | xato: `topish` xatoni `%w` bilan oʻraydi, `holat` esa `err == ErrTopilmadi` deb tekshiradi → `GET /buyurtma/42 → 500` |
| `tuzatildi/main.go` | xuddi shu kod, `errors.Is(err, ErrTopilmadi)` bilan → `404` |
| `promptlar/1-yomon.txt` | bitta qator: «Go'da xato tekshirish ishlamayapti, nega?» |
| `promptlar/2-kontekst.txt` | + kim, qaysi loyiha, Go versiyasi, kutilgan natija, butun kod va `go run` natijasi |
| `promptlar/3-format.txt` | + natija formati (aynan 3 qator: `Sabab:` / `Tuzatish:` / `Tekshirish:`) va *nega* shunday format |
| `promptlar/4-misol.txt` | + boshqa xato (nil map) uchun bitta namunaviy javob |
| `promptlar/5-xml.txt` | 4 bilan bir xil mazmun, lekin `<kontekst>`, `<kod>`, `<natija>`, `<format>`, `<misol>`, `<savol>` teglarida |
| `promptlar/6-test.txt` | `holat` uchun table test soʻrash: `<kod>` + `<vazifa>` |
| `promptlar/7-test-misol.txt` | 6 + `<misol>`: GO-1.07 dagi `kabisa` testi ([007-test-bench/sodda/jadval_test.go](../007-test-bench/sodda/jadval_test.go)) |
| `takrorla.sh` | promptni Claude'ga N marta beradi, javobni `javoblar/<prompt>/<i>.txt`, metamaʼlumotni `<i>.json` ga yozadi. Javoblar shu skriptning oldingi nomi (`sora.sh`) va formatlanishi bilan olingan; bayroqlar bir xil |
| `javoblar/` | **21 ta haqiqiy javob** (7 prompt × 3), sana, model, token va vaqt bilan |
| `olchov/` | javoblarni oʻlchaydi: soʻz, format, tayyor qator, tuzatishni ishga tushirib 404 tekshiruvi; `-test`: uslub, PASS/FAIL |
| `homework/` | uy vazifasi promptlari, javoblari (`homework/javoblar/`) va `guru/` |

## Tajriba qanday oʻtkazildi

- **Sana va model:** 2026-11-26, `claude-sonnet-5-5`, effort `high`, Claude Code CLI 2.1.332 (`claude -p`, *print* rejimi).
- **Muhit:** `takrorla.sh` boʻsh vaqtinchalik papkada ishlaydi (CLAUDE.md va loyiha fayllari taʼsir qilmaydi), `--system-prompt ""`,
  `--tools ""` (Claude fayl oʻqiy olmaydi, buyruq ishlata olmaydi).
- **Halol eslatma 1 — CLI qoʻshimchasi:** CLI baribir qisqa tizim matni qoʻshadi. Bir qatorli promptda `input_tokens_jami` (input + cache yaratish + cache oʻqish) 581, demak qoʻshimcha taxminan 570 token. Yomon promptga kelgan javob oxiridagi «loyiha papkasini koʻrsating, men kodni oʻzim koʻrib chiqaman» gapi shundan boʻlsa kerak.
- **Halol eslatma 2 — Haiku yozuvi:** `.json` dagi `models` roʻyxatida `claude-haiku-5-5` ham bor. Bu CLI'ning oʻz kichik xizmat chaqiruvi, javobni Sonnet yozgan.
- **Nega uch marta:** javob har safar boshqacha boʻlishi mumkin ([CL-1.01](../cl-01-next-token/)).
  Hozirgi modellarda `temperature` yoʻq, shuning uchun takrorlash — yagona halol yoʻl.

Oʻzingiz qaytarish (Claude Code CLI va `jq` kerak; javoblaringiz biznikidan farq qilishi mumkin):
```bash
./takrorla.sh promptlar/3-format.txt 3      # javoblar/3-format/1..3.txt ustiga yozadi
MODEL=claude-haiku-5-5 ./takrorla.sh promptlar/1-yomon.txt 3
```
Claude Code ichidan ishga tushirsangiz: `env -u CLAUDECODE ./takrorla.sh …`.

## Videodagi buyruqlar

Hammasi shu papkadan (`lessons/cl-03-prompt`), `go1.26.8 linux/amd64`.
```bash
cat promptlar/1-yomon.txt
grep -E '^\*\*[0-9]' javoblar/1-yomon/1.txt        # 7 ta taxmin; haqiqiy sabab 4-oʻrinda
head -3 promptlar/2-kontekst.txt
go run ./bozor                                     # GET /buyurtma/42 → 500
grep -h "o'zingiz almashtiring" javoblar/2-kontekst/*.txt
tail -4 promptlar/3-format.txt
cat javoblar/3-format/1.txt
go run ./olchov
grep -hE '^func Test|:= \[\]struct' javoblar/6-test/1.txt javoblar/7-test-misol/1.txt
go run ./olchov -test
go run ./homework/guru
```

Natija:
```
prompt        soʻz (3 marta)  format  tayyor      404
1-yomon       340 374 345     0/3     0/3         —
2-kontekst    196 208 234     0/3     0/3         —
3-format      31 31 31        3/3     3/3         3/3
4-misol       28 32 29        3/3     3/3         3/3
5-xml         32 27 28        3/3     3/3         3/3

prompt        soʻz (3 marta)  nomlar  tuzatilgan  eski xato
6-test        262 250 262     1 1 1   PASS 3/3    FAIL 3/3
7-test-misol  97 98 107       3 3 3   PASS 3/3    FAIL 3/3
```
Ustunlar:
- `format`: javob aynan uch qator va `Sabab:`, `Tuzatish:`, `Tekshirish:` bilan boshlanadi.
- `tayyor`: javobda koʻchirib qoʻyiladigan `errors.Is(err, ErrTopilmadi)` bor.
- `404`: `Tuzatish:` qatori `bozor/main.go` dagi `if err == ErrTopilmadi {` oʻrniga qoʻyiladi, dastur ishga tushiriladi va `/buyurtma/42` haqiqatan 404 qaytaradi.
- `nomlar`: misoldagi (GO-1.07) nomlardan nechtasi javobga koʻchgan: `jadval`, `tc.nom`, `kutilgan` (0–3). Bu uslubning hammasi emas, faqat nomlar; qolgan farqni `homework/guru` va `wc -l` koʻrsatadi.
- `tuzatilgan`: Claude yozgan test `tuzatildi/` bilan oʻtadi.
- `eski xato`: oʻsha test `bozor/` (xatoli) bilan yiqiladi, yaʼni xatoni ushlaydi.

Nima koʻrdik:
1. **Yomon prompt**: 340–374 soʻz, har safar 7 ta taxmin, haqiqiy sabab 3- yoki 4-oʻrinda. Oxirida Claude kod, kutilgan va haqiqiy natija hamda Go versiyasini soʻradi.
2. **Kontekst**: sabab 3/3 toʻgʻri, javob ~200 soʻz, lekin tayyor qator 0/3. Promptda «Xatoni oʻzim tuzatib, sababini tushunmoqchiman» deyilgan edi, shuning uchun Claude yoʻl koʻrsatdi («…shartni oʻzingiz almashtiring»). Kontekst javobning ohangini ham oʻzgartiradi.
3. **Format**: 31 soʻz, 3 qator, tuzatish 3/3 ishlaydi.
4. **Misol va XML**: bu kichik, aniq vazifada sezilarli oʻzgarish **yoʻq**. Biz XML'ning alohida taʼsirini oʻlchay olmadik.
   Hujjatdagi «30 foizgacha» raqami XML haqida emas: u 20k+ tokenli, koʻp hujjatli promptlarda savolni oxiriga qoʻyish haqida.
5. **Misolning foydasi testda koʻrindi**: misolsiz — 2 ta funksiya, 11 holat, `tests`/`tt`, ~91 qator; misol bilan — 1 funksiya, 6–7 holat, `jadval`/`tc`, 31–32 qator. Ikkalasi ham toʻgʻri (tuzatilgan kodda PASS, xatoli kodda FAIL): misol shaklni oʻzgartirdi, toʻgʻrilikni emas.
   Diqqat: 7-promptda misoldan tashqari «testlar <misol> dagidek yoziladi» degan bitta qator ham bor. Bu qator misolsiz maʼnosiz, shuning uchun alohida sinamadik.

**Cheklovlar.** Bitta model (Sonnet 5.5, effort `high`), bitta vazifa, har prompt 3 marta. Bu qonun emas, takrorlab koʻrish mumkin boʻlgan tajriba. Effort, boshqa model yoki CLI oʻrniga API boshqa natija berishi mumkin. 1–2-promptlardagi «sabab toʻgʻri topildi» degan baho dastur bilan emas, qoʻlda oʻqib berilgan.

## Manbalar

Anthropic hujjatlari (2026-11-26 da oʻqildi):
- [Prompting best practices](https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices):
  - *Be clear and direct*: «brilliant but new employee who lacks context on your norms and workflows», «Show your prompt to a colleague with minimal context».
  - *Add context to improve performance*: qoidaning sababini tushuntiring, «Claude can generalize from the explanation».
  - *Use examples effectively*: «Include 3–5 examples», «Mirror your actual use case closely», «vary enough…», `<example>` teglari.
  - *Structure prompts with XML tags*: «parse complex prompts unambiguously», «consistent, descriptive tag names».
  - *Long context prompting*: 20k+ token: hujjat tepada, savol oxirida, «up to 30 percent».
  - *Migrating away from prefilled responses*.
- [Prompt engineering overview](https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/overview): avval muvaffaqiyat mezoni va uni tekshirish usuli. Biznikida bular `olchov` va `go test`.
- Prefill va `temperature`/`top_p`/`top_k` hozirgi modellarda 400 xato qaytaradi:
  - [Sonnet 5.5 migration guide](https://platform.claude.com/docs/en/models/sonnet-5-5/migration-guide)
  - [Opus 5.5](https://platform.claude.com/docs/en/models/opus-5-5/migration-guide)
  - [Haiku 5.5](https://platform.claude.com/docs/en/models/haiku-5-5/migration-guide)
  - [Fable 5.2](https://platform.claude.com/docs/en/models/fable-5-2/migration-guide)
- [Anthropic Prompt Engineering Interactive Tutorial](https://github.com/anthropics/prompt-eng-interactive-tutorial): yaxshi kurs, lekin Claude 3 Haiku uchun yozilgan. 5-bob prefill'ni oʻrgatadi, hozirgi modellarda u ishlamaydi.
- Atama tarixi: Brown va boshq., [Language Models are Few-Shot Learners](https://arxiv.org/abs/2005.14165) (2020); sharh: Schulhoff va boshq., [The Prompt Report](https://arxiv.org/abs/2406.06608) (2024).
- Go: [errors.Is](https://pkg.go.dev/errors#Is), [fmt.Errorf `%w`](https://pkg.go.dev/fmt#Errorf).

## Dars matni

**hook.** Bitta Go xatosi va bitta model, lekin ikki xil prompt: biri yettita taxminli roʻyxat oldi, ikkinchisi ishlaydigan tuzatish. Prompt, bu modelga yozgan topshirigʻimiz. Claude Sonnet besh nuqta beshga Go'da xato tekshirish ishlamayapti, nega, deb uch marta yozdik. Har safar taxminan uch yuz ellik soʻzli javob keldi, ichida yettita taxmin. Haqiqiy sabab uchinchi yoki toʻrtinchi oʻrinda yashiringan edi. Oxirida esa Claude kodni soʻradi. Model bir xil edi, oʻzgargani faqat prompt.

**reja.** Birinchi darsda koʻrdik: model faqat siz bergan matnni koʻradi va bir xil savolga har xil javob berishi mumkin. Shuning uchun bugun har bir promptni uch martadan beramiz. Promptga toʻrt narsani navbat bilan qoʻshamiz: kontekst, natija formati, misol va XML teglar. Har qadamda javob qanday oʻzgarishini Go dasturi bilan oʻlchaymiz.

**tajriba.** Avval tajriba qanday oʻtkazilganini koʻramiz. Takrorla skripti Claude Code'ning print rejimini ishlatadi: model Sonnet besh nuqta besh, effort yuqori, asboblar oʻchirilgan, papka boʻsh. Claude bizning promptdan tashqari CLI qoʻshadigan qisqa tizim matnini ham oladi, taxminan besh yuz yetmish token, uning mazmunini biz koʻrmaymiz. Har bir javob oʻzgartirilmasdan faylga yoziladi, yonida sana, model va tokenlar soni. Yigirma bitta javobning hammasi GitHub'da.

**xodim.** Anthropic'ning prompt yozish qoʻllanmasi Claude'ni aqlli, lekin ishga yangi kelgan xodimga oʻxshatadi. U sizning qoidalaringiz va ish jarayoningizni bilmaydi. Tekshirish usuli ham oddiy: promptni vazifadan bexabar hamkasbga koʻrsating. Bizning savolni olgan hamkasb darhol soʻraydi: qaysi kod, qaysi Go versiyasi, nima kutilgan va nima chiqdi? Claude ham javob oxirida aynan shularni soʻradi.

**kontekst.** Ikkinchi promptda hamkasb soʻragan hamma narsa bor. Kim yozyapti va qaysi loyiha, Go bir nuqta yigirma olti nuqta sakkiz, kutilgan javob toʻrt yuz toʻrt. Ostida butun kod va uning haqiqiy natijasi: qirq ikkinchi buyurtma uchun besh yuz qaytyapti. Xato Go kursining oʻn toʻrtinchi darsini koʻrganlarga tanish boʻlishi mumkin. Doimiy kontekstni har safar yozmaslik uchun esa oldingi darsdagi Project koʻrsatmasi bor.

**kontekst-javob.** Uchala javobni oʻqib chiqdik, uchalasi ham sababni toʻgʻri topdi: foiz dablyu bilan oʻralgan xato teng belgisi bilan solishtirilsa, topilmaydi, errors.Is kerak. Javob ikki yuz soʻz atrofiga qisqardi, lekin hech birida koʻchirib qoʻyadigan tayyor qator yoʻq. Bittasi hatto hujjat linkini berib, shartni oʻzingiz almashtiring dedi. Sababi promptning oʻzida: xatoni oʻzim tuzatib, tushunmoqchiman degan edik. Kontekst javobning aniqligini ham, ohangini ham oʻzgartiradi.

**format.** Uchinchi promptga natija formatini qoʻshamiz: aynan uch qator, sabab, tuzatish va tekshirish. Sababini ham yozdik: javobni terminalda oʻqiyman. Anthropic hujjati shunday deydi: qoida nega kerakligini tushuntirsangiz, Claude uni yaxshiroq umumlashtiradi. Natija oʻttiz bir soʻz, uch qator va tayyor tuzatish. Uch marta berdik, uchala javob deyarli soʻzma-soʻz bir xil chiqdi.

**olchov.** Anthropic prompt yozishdan oldin muvaffaqiyat mezonini va uni tekshirish usulini belgilashni maslahat beradi. Bizda bu Go'da yozilgan kichik dastur, u oʻn beshta javobni oʻlchaydi. Birinchi ustun soʻzlar soni: uch yuz qirqdan oʻttiz atrofiga tushdi. Format va tayyor qator faqat uchinchi promptdan boshlab paydo boʻldi. Oxirgi ustun eng muhimi: dastur Claude bergan qatorni kodga qoʻyib, ishga tushirdi va toʻqqizta javobning hammasida toʻrt yuz toʻrt oldi. Toʻrtinchi va beshinchi promptlar, ya'ni misol va XML, bu kichik vazifada hech narsani oʻzgartirmadi.

**misol.** Misol soʻz bilan tushuntirish qiyin boʻlgan narsani koʻrsatadi, masalan kursimizdagi test uslubini. Ikkinchi tajribada Claude'dan holat funksiyasi uchun table test soʻradik. Oltinchi promptda faqat kod va vazifa bor. Yettinchisida esa Go kursining yettinchi darsidagi kabisa testi va testni shunday yoz degan bitta qator qoʻshilgan. Anthropic hujjati uchtadan beshtagacha misolni tavsiya qiladi: ular haqiqiy vazifaga yaqin va har xil boʻlsin. Biz soddalik uchun bittasini oldik, ikkinchisi uy vazifasida.

**misol-run.** Endi natijaga birga qaraymiz. Misolsiz Claude toʻqson qatordan uzun test yozdi: ikkita funksiya, oʻn bitta holat, tests va tt degan nomlar. Misol bilan esa oʻttiz bir, oʻttiz ikki qator, bitta funksiya va misoldagi jadval, tc degan nomlar. Ikkalasi ham toʻgʻri ishlaydi: tuzatilgan kodda testlar oʻtdi, eski kodda esa xatoni ushladi. Misol javobning toʻgʻriligini emas, shaklini oʻzgartirdi.

**xml.** Oltinchi va yettinchi promptlarda XML teglarni ishlatdik: kod, misol va vazifa alohida teg ichida. Yettinchi promptda ikki xil Go kodi bor, teglar qaysi biri bizning kod, qaysi biri faqat namuna ekanini ajratadi. Hujjatga koʻra, teglar murakkab promptni aniq ajratishga yordam beradi, nomlari esa tushunarli va doim bir xil boʻlsin. Biz oʻzbekcha nomlarni ishlatdik va javoblar toʻgʻri chiqdi. Lekin XML'ning alohida ta'sirini bu tajribada koʻra olmadik.

**uzun.** Uzun promptlar uchun hujjatda yana bitta aniq maslahat bor. Matn yigirma ming tokendan oshsa, hujjatlarni promptning boshiga, savolni esa oxiriga qoʻying. Anthropic sinovlarida bu javob sifatini oʻttiz foizgacha oshirgan. Bu ayniqsa bir nechta hujjatli murakkab promptlarda seziladi. Bizning kichik promptimizda esa toʻrtinchi va beshinchi variant bir xil natija berdi.

**eski.** Internetdagi eski qoʻllanmalarda yana ikki usul uchraydi. Birinchisi prefill, ya'ni Claude javobining boshini oʻzingiz yozib berish. Anthropic'ning interaktiv qoʻllanmasi uni eski Claude uch Haiku modeli uchun oʻrgatgan. Hozirgi toʻrtta modelda esa bunday soʻrov toʻrt yuz kodli xato qaytaradi. Ikkinchisi temperature, top_p va top_k: ularni yuborsangiz ham endi xato qaytadi. Hujjat bitta yoʻlni aytadi: xohishingizni promptning oʻzida aniq yozing.

**savol.** Endi sizga bitta savol beraman. Uchinchi promptda ham xatoni oʻzim tuzatib, tushunmoqchiman degan gap turibdi, lekin Claude tayyor qatorni berdi. Nega shunday boʻldi, va oʻrganmoqchi boʻlsangiz, formatni qanday oʻzgartirardingiz? Videoni pauza qilib, oʻylab koʻring.

**javob.** Bizning taxminimiz shunday: prompt oʻziga zid edi. Bir joyda oʻzim tuzataman dedik, boshqa joyda tuzatish qatorini aniq soʻradik, va Claude formatdagi aniq talabni bajardi. Oʻrganmoqchi boʻlsangiz, Tuzatish oʻrniga Ishora qatorini soʻrang va tayyor kod yozma deng. Biz buni ham uch marta sinadik: uchala javobda sabab va yoʻnalish bor, tayyor kod yoʻq. Keyingi darsda Claude'ni aynan shunday ustoz sifatida ishlatamiz.

**homework.** Uy vazifasi, toʻrt darajada. Oson: oʻz Go savolingizni yomon va formatli promptda uch martadan bering, Claude.ai'da uchta yangi chatda ham boʻladi. Oʻrta: uchinchi promptda Tuzatish oʻrniga Ishora soʻrang va javoblarni solishtiring. Qiyin: yettinchi promptga ikkinchi, boshqacha misol qoʻshing va Claude qaysi uslubni olishini koʻring. Guru: oʻlchov dasturiga test holatlari sonini sanaydigan ustun qoʻshing. Javoblar GitHub'dagi kurs repository'sida.

**outro.** Bugun koʻrdik: kontekst Claude'ga vazifani tushuntiradi, format javobni ishlatsa boʻladigan qiladi, misol esa uslubni koʻrsatadi. Har bir javobni baribir ishga tushirib tekshirdik. Rol berish va loyiha qoidalari kabi usullarni Claude Code darslarida koʻramiz. Keyingi darsda Claude'ni javob mashinasi emas, ustoz sifatida ishlatamiz. Darsni koʻrganingiz uchun rahmat, koʻrishguncha.

## Uy vazifasi

1. **Oson.** Oʻz Go savolingizni ikki xil bering (Claude.ai'da uchta yangi chatda ham boʻladi): `1-yomon` kabi bitta qator va `3-format` kabi kontekst + format. Ikkalasini `takrorla.sh` bilan uch martadan ishga tushiring va soʻzlarni sanang (`wc -w javoblar/*/*.txt`).
2. **Oʻrta.** `3-format` da `Tuzatish:` oʻrniga `Ishora:` soʻrang (tayyor kod yozmasin) va javoblarni solishtiring.
3. **Qiyin.** `7-test-misol` ga ikkinchi, boshqacha uslubdagi misol qoʻshing (masalan, cl-02 dagi map'li `TestLuhn`). Claude qaysi uslubni oladi?
4. **Guru.** `olchov` ga test javobida nechta `Test` funksiya va nechta jadval holati borligini sanaydigan ustun qoʻshing.

### Javoblar

1. **Natija sizniki boʻladi.** Bizda bir qatorli savol 340–374 soʻz, kontekst + format esa 31 soʻz berdi. Claude Code yoʻq boʻlsa: Claude.ai'da uchta yangi chat oching (har biri toza kontekst), promptni yopishtiring, soʻzlarni sanang va tuzatishni Go kodingizga qoʻyib `go run` qiling. Har bir javobni ishga tushirib tekshiring (CL-1.05).
2. **Prompt va javoblar:** `homework/ishora.txt` va `homework/javoblar/ishora/1–3.txt` (2026-11-26, Sonnet 5.5).
   - Uchala javobda sabab va «qayerga qarash» bor: `holat` dagi `==` hamda `errors` paketidagi oʻralgan xatolar zanjirini tekshiradigan funksiyalar.
   - Tayyor `errors.Is(err, ErrTopilmadi)` qatori 0/3.
   - Bir javobda `Tekshirish:` = `go doc errors.Is`. Bu formatdagi «buyruq» talabini boshqacha tushunish.
   - Bu prompt «tayyor kod yozma» deb ham aytadi, shuning uchun natijani faqat «Tuzatish»ni olib tashlashga bogʻlab boʻlmaydi.
   - Tuzatishni endi oʻzingiz yozasiz. CL-1.04 shu usul haqida.
3. **Prompt va javoblar:** `homework/qiyin.txt` (ikkita misol `<misollar>` ichida) va `homework/javoblar/qiyin/1–3.txt`.
   - Uchala javob ham **struct jadval** uslubini oldi (`jadval`, `tc`, `kutilgan`), map'li uslubni olmadi.
   - Lekin uchalasida oʻzbekcha harflar `ʻ` (U+02BB) bilan yozildi (`oʻralgan`), xuddi cl-02 misolining izohlaridagidek. 7-prompt javoblarida oddiy `'` edi.
   - Bu hujjat ogohlantirgan «keraksiz andoza»ning (*unintended pattern*) kichik, lekin haqiqiy namunasi.
   - Uchala test `tuzatildi/` bilan oʻtdi va `bozor/` bilan yiqildi.
4. **Yechim:** `homework/guru/main.go` (`go run ./homework/guru`):
   ```
   6-test        Test funksiya: 2 2 2 · holat: 11 11 11
   7-test-misol  Test funksiya: 1 1 1 · holat: 6 6 7
   ```
