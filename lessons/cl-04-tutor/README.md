# CL-1.04 · AI — ustoz, javob mashinasi emas: bilim illyuziyasisiz oʻrganish

Oldingi dars: [CL-1.03 · Prompt — bu texnik topshiriq: kontekst, misol, cheklov](../cl-03-prompt/) ·
Keyingi dars: CL-1.05 · AI javobini tekshirish: manba, ishga tushirish va ikkinchi fikr

**Darsning bitta gʻoyasi:** AI javobni yozib bersa-yu, siz faqat koʻchirsangiz, keyin oʻzingiz yecha olmasligingiz mumkin. Avval oʻzingiz yozasiz (xato boʻlsa ham),
Claude esa ustoz boʻladi: tayyor kodni bermaydi, xato qayerda ekanini koʻrsatadi va savol beradi.

Bitta kichik Go xatosini Claude Sonnet 5.5 ga ikki xil berdik — «tuzatib ber» (javob mashinasi) va ustoz qoidalari
bilan — har birini **uch marta**. Javoblar **oʻzgartirilmasdan** saqlangan, Go dasturi esa ularni oʻlchaydi va Claude
bergan kodni **ishga tushirib** tekshiradi.

| Fayl / papka | Nima |
|---|---|
| `bozor/main.go` | oʻquvchining urinishi: `errors.Is(err, &MaydonXato{})` → narx manfiy boʻlsa ham `500` (`400` boʻlishi kerak). `go vet` jim |
| `urinish/main.go` | oʻquvchining 2-urinishi, **ataylab xato**: `errors.As(err, mx)` (`&` yoʻq) → `panic`, `go vet` topadi |
| `tuzatildi/main.go` | oʻquvchining oʻzi yozgan tuzatish: `errors.AsType[*MaydonXato](err)` → `400 notoʻgʻri maydon: narx` |
| `promptlar/1-javob.txt` | «Bu Go kodni tuzatib ber…» + kod + natija (boʻsh tizim matni) |
| `promptlar/ustoz-tizim.txt` | ustoz qoidalari — tizim matni sifatida beriladi (Claude.ai'da: Project koʻrsatmasi, CL-1.02) |
| `promptlar/2-ustoz.txt` | 1-xabar: kod, natija va `<urinishim>` |
| `promptlar/3-urinish.txt` | 2-xabar: oʻquvchi `errors.As` ni topdi, lekin `panic` |
| `promptlar/4-tuzatish.txt` | 3-xabar: `errors.AsType` bilan tuzatish + «menga bitta savol ber» |
| `promptlar/5-claude-code.txt` | Claude Code uchun vazifa (Learning va default uslub) |
| `takrorla.sh` | suhbatni N marta (standart 3) boshidan oxirigacha oʻtkazadi, 2-xabardan boshlab `--resume` |
| `uslub.sh` | Claude Code'ni `outputStyle` bilan `bozor/` nusxasida ishga tushiradi, javob va diff'ni saqlaydi |
| `javoblar/` | **12 ta chat javobi** (1-javob 3 + 2-ustoz 3×3), **6 ta Claude Code javobi** (+ diff), 1 ta pilot |
| `olchov/` | javoblarni oʻlchaydi: tayyor kod, savol, kodni qoʻyib ishga tushirish (400) |
| `homework/sinov/` | Claude bergan uchta tekshiruv savoliga javobni ishga tushirib koʻrsatadi |

## Tajriba qanday oʻtkazildi

- **Sana va model:** 2026-11-26, `claude-sonnet-5-5`, effort `high`, Claude Code CLI 2.1.332 (`claude -p`, *print* rejimi).
- **Muhit:** har suhbat yangi, boʻsh vaqtinchalik papkada; `--tools ""` (Claude fayl oʻqiy olmaydi, buyruq ishlata olmaydi),
  `--setting-sources ""`, MCP va slash buyruqlarsiz. Ustoz suhbatida `--system-prompt "$(cat promptlar/ustoz-tizim.txt)"`,
  javob mashinasida `--system-prompt ""`. CL-1.03 dagidek CLI baribir qisqa tizim matni qoʻshadi.
- **Suhbat:** 2- va 3-xabar oʻsha suhbatning davomi (`--resume <session_id>`). Bu xabarlarni biz (odam) yozdik, bir marta,
  va uchala suhbatga bir xil yubordik.
- **Pilot:** `javoblar/0-pilot/1-1.txt` — 2–3-xabarni yozishdan oldin Claude qanday savol berishini koʻrish uchun bitta
  sinov. Videoda koʻrsatilmaydi va oʻlchovga kirmaydi.
- **Halol eslatma — ikki farq:** ustoz promptida ham qoidalar, ham `<urinishim>` izohi bor; javob mashinasida faqat
  kod (urinish va izoh kod ichida ikkalasida ham bor). Natijani faqat qoidalarga bogʻlab boʻlmaydi.
- **Claude Code ichidan ishga tushirsangiz:** skriptlar `CLAUDECODE` va `CLAUDE_CODE_*` oʻzgaruvchilarini oʻzi oʻchiradi.

Oʻzingiz qaytarish (Claude Code CLI va `jq` kerak; javoblaringiz biznikidan farq qiladi):
```bash
./takrorla.sh 1-javob - promptlar/1-javob.txt
./takrorla.sh 2-ustoz promptlar/ustoz-tizim.txt promptlar/2-ustoz.txt promptlar/3-urinish.txt promptlar/4-tuzatish.txt
./uslub.sh Learning 3; ./uslub.sh default 3
```

## Videodagi buyruqlar

Hammasi shu papkadan (`lessons/cl-04-tutor`), `go1.26.8 linux/amd64`. Videoda Claude javoblari kod kartochkasida
koʻchirma sifatida koʻrsatiladi (`javoblar/1-javob/1-1.txt`, `javoblar/2-ustoz/1-1.txt`, `1-2.txt`, `1-3.txt`):
markdown belgilari (`, **) olib tashlangan va uzun qatorlar boʻlingan, soʻzlar oʻzgarmagan.
```bash
go run ./bozor        # narx=-5 → 500
go run ./tuzatildi    # narx=-5 → 400 notoʻgʻri maydon: narx
go run ./olchov
go run ./olchov -code # videoda yoʻq, pastda
```

Natija:
```
usul             tayyor kod  savol  ishladi
javob mashinasi  3/3         0/3    3/3
ustoz, 1-xabar   0/3         3/3    —
ustoz, 2-xabar   0/3         3/3    —
ustoz, 3-xabar   0/3         3/3    —

usul             tayyor kod  savol  ishladi  TODO(human) Insight
Code: default    3/3         0/3    3/3      0/3         0/3
Code: Learning   3/3         0/3    3/3      0/3         2/3
```
Ustunlar:
- `tayyor kod`: javobdagi ```go blokida (Claude Code'da — Claude tahrirlagan faylda) tayyor tuzatish bor: `errors.As(err, &…)` yoki `errors.AsType[…]`.
- `savol`: javob oʻquvchiga savol beradi (kod bloklaridan tashqarida «?» bor).
- `ishladi`: oʻsha kod `bozor/main.go` ga qoʻyiladi (chat: `func holat` bloki; Claude Code: diff `patch` bilan) va dastur `narx=-5 → 400` chiqaradi.
- `TODO(human)`: Claude faylda sizga joy qoldirdimi. `Insight`: javobda `Insight` bloki bormi.
- `-v`: har javobdagi soʻzlar soni. Ustoz javoblari 66–110 soʻz, hammasi qoidadagi 150 dan kam.

Nima koʻrdik:
1. **Javob mashinasi** 3/3 tayyor `holat` funksiyasini yozdi (`errors.As(err, &me)`), tushuntirish bilan. Kod ishladi, lekin biz hech narsa yozmadik. Qiziq: uchalasi ham Go 1.26 dagi `errors.AsType` ni emas, eski `errors.As` ni berdi ([GO-1.14](../014-errors-values/): `errors.As` hujjati endi *«For most uses, prefer AsType»* deydi).
2. **Ustoz, 1-xabar:** uchala javob xatoli qatorni koʻrsatdi, `==` pointer'larni solishtirishini tushuntirdi, `errors.As` nomini **aytmadi** (0/3) va «turni qidiradigan funksiya bormi?» deb soʻradi.
3. **Ustoz, 2-xabar:** uchalasi `errors.As(err, mx)` qatorini va `mx` nil ekanini koʻrsatdi, `&mx` ni yozmadi, qanday uzatishni soʻradi.
4. **Ustoz, 3-xabar:** tuzatishni tasdiqladi va har biri bitta tekshiruv savoli berdi (uy vazifasida).
5. **Claude Code Learning uslubi** (videoda yoʻq — sinov tor edi): hujjatga koʻra u «real design decision» boʻlgan joyda `TODO(human)` qoldirib, sizni **kutadi**; print rejimi kuta olmaydi, vazifamiz esa oddiy «tuzat» edi. Bizning bitta vazifamizda, print rejimida, 3/3 kodni oʻzi yozdi, `TODO(human)` 0/3, `Insight` 2/3. Uslub haqiqatan yoqilgan edi (`/output-style` shu sozlamada «Output style: Learning» dedi). Bitta qoʻshimcha interaktiv sinov (saqlanmagan) ham kodni oʻzi yozdi. Hujjatning oʻzi: *«An output style gives Claude instructions to follow. It doesn't guarantee that something always happens or never happens.»*

**Cheklovlar.** Bitta model (Sonnet 5.5), bitta vazifa, har usul 3 marta; odamlarda oʻrganishni biz oʻlchamadik — buning uchun quyidagi tadqiqotlar bor. Bu qonun emas, takrorlab koʻrish mumkin boʻlgan tajriba.

## Manbalar

Tadqiqotlar (toʻliq matni oʻqildi, 2026-11-26):
- Bastani va boshq., [Generative AI without guardrails can harm learning: Evidence from high school mathematics](https://doi.org/10.1073/pnas.2422633122), PNAS 2025 ([ochiq matn, PMC](https://pmc.ncbi.nlm.nih.gov/articles/PMC12232635/)): ~1000 oʻquvchi, Turkiya, uch guruh (GPT Base, GPT Tutor, AI'siz); mashqda +48 % (GPT Base) va +127 % (GPT Tutor); AI olib qoʻyilgach GPT Base −17 %, GPT Tutor'da (ishora + oʻqituvchi yozgan yechimlar) zarar deyarli yoʻq, lekin foyda ham yoʻq, va ular imtihonni yaxshiroq yozdik deb oʻylashgan; GPT Base'da eng koʻp xabar «What is the answer?»; oʻquvchilar yomonlashganini sezmadi.
- Rozenblit & Keil, [The misunderstood limits of folk science: an illusion of explanatory depth](https://pmc.ncbi.nlm.nih.gov/articles/PMC3062901/), Cognitive Science 2002: molniya (zipper), unitaz, qulf…; 1-tadqiqot: 16 aspirant, 1–7 baho → qadam-baqadam tushuntirish → baho tushdi (T1 va T2 farqi p < .002).
- Roediger & Karpicke, [Test-enhanced learning](https://psychology.ecu.edu/wp-content/pv-uploads/sites/216/2019/03/Roediger-Karpicke-2006.pdf), Psychological Science 2006: 2-tajriba: 5 daqiqadan keyin faqat oʻqish (SSSS) 83 %, bir oʻqish + uch eslab yozish (STTT) 71 %; bir haftadan keyin 40 % va 61 %; oʻqiganlar oʻziga koʻproq ishongan.
- Bertsch, Pesta, Wiscott & McDaniel, [The generation effect: A meta-analytic review](https://mcdaniel97.github.io/Publications/Bertsch%20et%20al.%202007.pdf), Memory & Cognition 2007: 86 tadqiqot, 445 effekt, oʻrtacha d = 0,40.
- Shen & Tamkin, [How AI Impacts Skill Formation](https://arxiv.org/abs/2601.20245), arXiv 2026 · [Anthropic maqolasi](https://www.anthropic.com/research/AI-assistance-coding-skills): 52 dasturchi, Trio (Python), test: AI bilan 50 %, AI'siz 67 %; eng katta farq — debugging savollari; faqat tushuntirish soʻraganlar (n = 7) va kod olib keyin «nega?» deb soʻraganlar (n = 2) yuqori ball; mualliflar bu guruhlar uchun sabab-natija bogʻlanishini daʼvo qilmaydi; yordamchi GPT-4o edi.
- Qoʻshimcha: Fisher, Goddu & Keil, [Searching for explanations: How the Internet inflates estimates of internal knowledge](https://www.apa.org/pubs/journals/releases/xge-0000070.pdf), JEP: General 2015 — internetda qidirish ham «bilaman» hissini oshiradi.
- Benchmark: Liu va boshq. (Harvard CS50), [Teaching CS50 with AI](https://cs.harvard.edu/malan/publications/V1fp0567-liu.pdf), SIGCSE 2024 — «guide students toward solutions rather than offer them outright».

Anthropic rasmiy sahifalari:
- [Introducing Claude for Education](https://www.anthropic.com/news/introducing-claude-for-education) (2025-04-02): Learning mode — «Asking "How would you approach this problem?" instead of providing immediate solutions».
- [Claude Code: Output styles](https://code.claude.com/docs/en/output-styles): Learning uslubi `Insight` bloklari va `TODO(human)`; `outputStyle` sozlamasi.

Go: [errors.Is](https://pkg.go.dev/errors#Is), [errors.As](https://pkg.go.dev/errors#As), [errors.AsType](https://pkg.go.dev/errors#AsType), [GO-1.14](../014-errors-values/).

## Dars matni

**hook.** Mashqda qirq sakkiz foiz yaxshi, imtihonda esa oʻn yetti foiz yomon, va bu bitta guruhning natijasi. Turkiyada mingga yaqin oʻquvchi uch guruhga boʻlindi. Bir guruh matematika mashqini ChatGPT'ga oʻxshash yordamchi bilan ishladi. Imtihonda AI'ni olib qoʻyishdi, va bu guruhning bahosi AI'siz oʻqiganlarnikidan oʻn yetti foizga past chiqdi. Oʻquvchilarning oʻzi esa buni sezmadi.

**goya.** Bugungi darsning bitta gʻoyasi: AI javobni yozib bersa-yu, siz faqat koʻchirsangiz, keyin oʻzingiz yecha olmay qolishingiz mumkin. Shuning uchun avval oʻzingiz yozasiz, xato boʻlsa ham. Sakson oltita tajriba ham shuni koʻrsatadi: oʻzingiz yaratgan narsa shunchaki oʻqiganingizdan yaxshiroq esda qoladi. Claude esa ustoz boʻladi: xato qayerda ekanini koʻrsatadi va savol beradi.

**illyuziya.** Avval bitta savol: kurtkadagi molniya qanday ishlaydi? Yel universitetida oʻn olti aspirant shunday narsalar boʻyicha bilimini bittadan yettigacha baholadi. Keyin ularni qadam-baqadam tushuntirib yozdi, va baholar tushib ketdi. Buni bilim illyuziyasi deyishadi: tushunaman degan his bor, lekin tushuntirib berolmaysiz.

**qoida.** Turk tadqiqotida uchinchi guruh ham bor edi. Ularning yordamchisi javobni aytmay, ishora berardi, oʻqituvchilar unga yechim va koʻp uchraydigan xatolarni ham yozib qoʻygan. U guruhda imtihondagi zarar deyarli yoʻqoldi, lekin foyda ham boʻlmadi. Ular esa imtihonni yaxshiroq yozdik, deb oʻylashdi. Demak, ustozdan keyin ham oʻzingizni tekshirasiz.

**vazifa.** Endi tajribaga oʻtamiz, kursdagi kichik onlayn doʻkon dasturida. U har bir buyurtmaga raqam bilan javob beradi. Doʻkon peshtaxtasini tasavvur qiling. Ikki yuz, buyurtma qabul qilindi. Toʻrt yuz, blankani notoʻgʻri toʻldirdingiz. Besh yuz, doʻkonning oʻzida nosozlik. Manfiy narx yozilsa, javob toʻrt yuz boʻlishi kerak.

**vazifa-run.** Dasturni ishga tushiramiz, har qatorda buyurtma, uning narxi va doʻkonning javobi. Birinchi qatorda narx toʻgʻri, javob ikki yuz, nil esa xato yoʻq degani. Ikkinchi qatorda narx minus besh va xato matni: narx musbat boʻlishi kerak. Lekin javob besh yuz, nosozlik.

**urinish.** Oʻquvchimiz xayoliy, uning xabarlarini biz yozdik, Claude javoblari esa oʻzgarishsiz. Birinchi qatordagi izohda u yozgan: oʻtgan safar errors.Is yordam bergan edi. Keyingi qatorlar bilan u shunday demoqchi: agar xato MaydonXato boʻlsa, ya'ni narx kabi katak notoʻgʻri boʻlsa, toʻrt yuz qaytar. Lekin qator boshqa narsani soʻraydi.

**pasport.** Qator nimani soʻrashini bir misolda koʻramiz. Pasportingiz yoʻqoldi, uni sumkaga, sumkani chamadonga solib topilmalar boʻlimiga olib kelishdi. Go'da xato ham shunday qavat-qavat oʻraladi, bu zanjir. errors.Is soʻraydi: ichida aynan mening pasportim bormi? Oʻquvchi esa yangi, boʻsh pasport yasab, shu bormi, deb soʻradi. Kerakli savol boshqa: ichida umuman pasport bormi?

**javob.** Birinchi yoʻl, javob mashinasi. Claude Sonnet besh nuqta beshga kodni berib, tuzatib ber, deymiz, uch marta. Uchala javobda ham tayyor kod bor, mana birinchisi. Birinchi qator topilgan pasport uchun boʻsh choʻntak tayyorlaydi. Ikkinchi qator errors.As bilan kerakli savolni beradi: ichida pasport bormi? Boʻlsa, toʻrt yuz qaytaradi.

**javob-xulosa.** Keling, xulosa qilamiz: nima boʻldi? Kod ishladi, Claude tushuntirish ham yozdi. Lekin biz hech narsa yozmadik va hech narsani eslamadik, tushuntirishni oʻqib oʻtib ketish esa oson. Turk tadqiqotida yordamchiga eng koʻp yozilgan xabar ham shu edi: javobi nima?

**ustoz-qoida.** Ikkinchi yoʻl, bu ustoz. Claude'ga oldindan olti qatorli qoida beramiz, Claude.ai'da u ikkinchi darsdagi Project koʻrsatmasiga yoziladi. Birinchi qator: sen ustozsan, oʻquvchi kodni oʻzi yozsin. Ikkinchisi: tayyor kodni ham, funksiya nomini ham yozma. Uchinchisi: nima toʻgʻri ekanini ayt, xatoli qatorni koʻrsat. Toʻrtinchisi: bitta savol ber, beshinchisi esa tuzatishini tekshir. Oxirgisi: qisqa va oʻzbekcha yoz.

**rasmiy.** Anthropic dastlab universitetlar uchun e'lon qilgan Claude for Education'da ham Learning rejimi bor. Rasmiy e'longa koʻra, u tayyor yechim oʻrniga, bu masalaga qanday yondashardingiz, deb soʻraydi.

**ustoz-1.** Endi oʻsha kodni ustoz qoidalari bilan beramiz. Claude xatoli qatorni aniq koʻrsatdi, lekin tuzatmadi. Keyin savol berdi: errors paketida, ya'ni Go'ning tayyor vositalari orasida, zanjirdan aniq qiymatni emas, turni qidiradigan funksiya bormi? Qoidadagidek, nomini aytmadi, uni oʻquvchi hujjatdan topadi.

**ustoz-2.** Oʻquvchi hujjatdan errors.As'ni topdi va oʻzi yozdi. Birinchi qator topilgan pasport uchun choʻntak tayyorlaydi. Lekin ikkinchi qatorda u choʻntakni emas, boʻsh qoʻlni uzatdi, va dastur toʻxtab qoldi. Claude yana tuzatmadi: xato shu qatorda, hatto uning ikkinchi qismida ekanini aytdi va bitta savol berdi.

**ustoz-3.** Oʻquvchi bu savolga javob bermadi, Go kursining oʻn toʻrtinchi darsidagi boshqa yoʻlni esladi: errors.AsType. Bu funksiya choʻntak soʻramaydi. Birinchi qator: ichida shu turdagi xato bormi, boʻlsa uni ber. Ikkinchi qator toʻrt yuz va notoʻgʻri katak nomini qaytaradi.

**ustoz-3-run.** Dasturni yana ishga tushiramiz. Birinchi qatorda narx toʻgʻri, javob ikki yuz. Ikkinchisida minus besh uchun endi toʻrt yuz va qaysi katak notoʻgʻri ekani: narx.

**sinov.** Oxirida oʻquvchi, meni sinab koʻr, dedi. Claude tuzatishni tasdiqladi va savol berdi: errors.Is nega false, ya'ni yoʻq, AsType esa nega true, ya'ni ha qaytardi? Oʻz soʻzlaring bilan tushuntir. Bu molniya tajribasining oʻzi: tushuntira olsangiz, tushungansiz.

**olchov.** Keling, xulosa qilamiz: nima boʻldi? Oʻn ikki javobni Go dasturi bilan sanadik. Birinchi qator, javob mashinasi: uchala javobda tayyor kod bor, savol yoʻq, kodni qoʻyib koʻrganda uchalasi ishladi. Keyingi uch qator, ustozning uch xabari: toʻqqiz javobning birortasida tayyor kod yoʻq, hammasida savol bor. Bu oʻrganishni emas, qoidaga amal qilishni oʻlchaydi.

**eslash.** Bir tajribada bir guruh matnni faqat qayta-qayta oʻqidi. Ikkinchi guruh bir marta oʻqib, keyin uch marta kitobni yopib, esida qolganini yozdi, buni eslab yozish deymiz. Besh daqiqadan keyin oʻqiganlar yaxshiroq esladi, sakson uch foizga yetmish bir, va oʻzlariga koʻproq ishondi. Bir haftadan keyin esa ular qirq foiz, eslab yozganlar oltmish bir foiz esladi.

**dasturchilar.** Anthropic tadqiqotchilari ellik ikki dasturchiga notanish vosita bilan vazifa berdi, yarmiga AI yordamchi berildi, u GPT toʻrt o edi. Keyingi testda AI guruhining bali oltmish yettidan ellikka tushdi. Faqat tushuntirish soʻragan yetti kishi esa yuqori ball oldi.

**savol.** Endi sizga bitta savol beraman. Siz javob mashinasidan kod oldingiz va u ishladi. Uni haqiqatan tushunganingizni qanday bilasiz? Videoni pauza qilib, oʻylab koʻring.

**javob2.** Chatni yoping va oʻsha qatorni esdan yozib koʻring, yoki uni doʻstingizga, hatto stol ustidagi oʻyinchoqqa tushuntirib bering. Yozolmasangiz yoki tushuntirolmasangiz, bu bilim illyuziyasi edi. Unda ustozga qaytib, yana soʻrang.

**homework.** Uy vazifasi, toʻrt darajada. Oson: qoidalarni Project'ga qoʻyib, oʻz savolingizni bering. Oʻrta: Claude bergan uchta savolga oʻzingiz javob yozing, keyin sinov dasturini ishga tushiring. Qiyin: qoidalarga meni uchta savol bilan sina, degan qator qoʻshing. Guru: oʻlchov dasturiga yuz ellik soʻz qoidasini qoʻshing.

**xulosa.** Xoʻsh, bugun nima boʻldi? Tushunaman degan his aldashi mumkin, uni tushuntirib yoki eslab yozib tekshiring. Javob mashinasi tez, lekin kodni yozgan siz emassiz. Ustoz bilan kodni oʻzingiz yozasiz, Claude esa xatoni koʻrsatib, savol beradi.

**outro.** Keyingi darsda Claude bergan javobni qanday tekshirishni koʻramiz. Darsni koʻrganingiz uchun rahmat, koʻrishguncha.

## Uy vazifasi

1. **Oson.** `promptlar/ustoz-tizim.txt` ni Claude.ai'dagi Project koʻrsatmasiga qoʻying (CL-1.02) va oʻz Go savolingizni bering: avval oʻz urinishingizni yozing.
2. **Oʻrta.** Claude 3-xabarda bergan uchta tekshiruv savoliga **avval oʻzingiz** javob yozing (chatni va bu README'ni yopib), keyin `go run ./homework/sinov` bilan solishtiring:
   - (1-suhbat) `errors.Is(err, &MaydonXato{})` nega `false`, `AsType[*MaydonXato]` esa nega `true`?
   - (2-suhbat) `tekshir` natijasini yana `fmt.Errorf("handler: %w", …)` bilan oʻrasak, `holat` baribir 400 qaytaradimi? `%v` bilan-chi?
   - (3-suhbat) xato pointersiz, `MaydonXato{…}` qiymat sifatida qaytsa, `errors.AsType[*MaydonXato](err)` nima qaytaradi?
3. **Qiyin.** Qoidalarga «oxirida meni uchta qisqa savol bilan sina, javobini aytma» qatorini qoʻshing va `takrorla.sh` bilan sinang. Claude savollarni oxirida beradimi, javobini aytib yubormaydimi?
4. **Guru.** `olchov` ga «≤ 150 soʻz» ustunini qoʻshing: ustoz javobi qoidadagi chegaradan oshdimi?

### Javoblar

1. **Natija sizniki.** Tekshiring: Claude tayyor kod yozdimi? Yozsa, qoidani «hatto soʻrasam ham» deb kuchaytiring.
2. `go run ./homework/sinov`:
   ```
   1. errors.Is(err, &MaydonXato{}): false · AsType[*MaydonXato]: true
   2. %w bilan yana oʻralgan: true · %v bilan: false
   3. pointersiz xato: AsType[*QiymatXato]: false · AsType[QiymatXato]: true
   ```
   - (1) `errors.Is` zanjirdagi har bir xatoni `==` bilan solishtiradi; pointer'lar uchun bu «aynan shu manzilmi?» degani, yangi `&MaydonXato{}` esa boshqa manzil. `AsType` esa **tur**ni qidiradi.
   - (2) `%w` zanjirni saqlaydi — qavatlar soni muhim emas. `%v` xatoni oddiy matnga aylantiradi, zanjir uziladi → 500.
   - (3) Zanjirda `MaydonXato` (qiymat) bor, `*MaydonXato` (pointer) yoʻq, shuning uchun `false`; `AsType[MaydonXato]` kerak (Error() qiymat qabul qiluvchida boʻlishi sharti bilan). Tur aynan mos kelishi kerak.
3. Prompt va javoblarni oʻzingiz saqlang (`./takrorla.sh qiyin promptlar/ustoz-tizim.txt …`). Tekshirish: javobda «?» bormi va tayyor javob yoʻqmi — `olchov` ga yangi qator qoʻshing.
4. `olchov/main.go` da `sozlar` roʻyxati bor: `len(strings.Fields(matn)) <= 150` sanab, yangi ustun chiqaring. Bizda 9/9 (66–110 soʻz).
