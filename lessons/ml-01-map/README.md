# ML-0.01 · Nega ML oʻrganyapmiz: xarita va halollik qoidalari

«ML: birga oʻrganamiz» trekining birinchi darsi · Keyingi dars: ML-0.02 · Maʼlumotdan oʻrganish nima: loss, chiziq va test

| Fayl | Nima |
|---|---|
| `tajriba.py` | kichik tajriba (uv + numpy, faqat CPU): sintetik `ball = 40 + 5·soat + shovqin`, 10 ta train talaba, 10 ta yangi (test) talaba; toʻgʻri chiziq va 9-darajali koʻphad |
| `bashorat.md` | 5 ta bashorat, `tajriba.py` **birinchi marta ishga tushishidan oldin** commit qilingan (`a9200d1`) |
| `natija.txt` | birinchi ishga tushirishning xom natijasi, oʻzgartirilmagan (`d380da2`, 45 soniyadan keyin) |
| `bashorat2.md`, `natija2.txt` | pauza savoli uchun xuddi shu tartib: 100 ta train talaba (`5dedd5b` → `3c439f3`) |
| `homework/` | uy vazifasi javoblari: `nuqtalar.py` (seed 1 train nuqtalari → CSV + numpy chizigʻi), `guru/main.go` (chiziq Go'da noldan) |

## Halollik tartibi (biz nima qildik)

1. Kod va bashorat yozildi, **ishga tushirilmasdan** commit qilindi (faqat `python -m py_compile`, u kodni bajarmaydi).
2. Birinchi ishga tushirish → `natija.txt` (keyin hech qachon tahrirlanmaydi) → commit.
3. `ming` rejimi (1000 seed tahlili) natijadan **keyin** qoʻshildi, bu skriptda ham yozilgan. `seeds` chiqishi keyinroq
   video terminaliga sigʻishi uchun ixchamlashtirildi, raqamlar `natija.txt` dagi bilan bir xil.
4. Git sanasini qoʻlda oʻzgartirish mumkin (`git commit --date`, `GIT_COMMITTER_DATE`), shuning uchun bu kriptografik
   isbot emas, ochiq odat.

```bash
git log --reverse --format='%h %ad  %s' --date=iso -- .
```

**Ishonch belgilari:** *tekshirdim* — shu kompyuterda ishga tushirdik/oʻlchadik, natija ekranda · *oʻqidim* — nomi
aytilgan manbada yozilgan, oʻzimiz takrorlamadik · *taxminim* — na oʻlchov, na manba; xato chiqishi mumkin.

**Bashoratlar hisobi:** 1-toʻplam 3/5, 2-toʻplam 2/3, jami **5/8**.

| # | Bashorat (oldin) | Natija | |
|---|---|---|---|
| 1 | chiziq: train ≈ 4, test 4–5 | train 1.6, test 4.4 | ✗ |
| 2 | daraja 9: train 0.0 | 0.0 | ✓ |
| 3 | daraja 9 test ≈ 30 (10…100) | 52.1 | ✓ (lekin 1000 seed'da faqat 175 tasi shu oraliqda: omad) |
| 4 | eng katta xato train oraligʻidan tashqarida | 7.2 soat, oraliq 0.3…9.5: ichida | ✗ (1000 seed'da 418 tasi ichida) |
| 5 | 5 seed'ning 5 tasida daraja 9 yomonroq | 5 / 5 | ✓ (1000 seed'da 999) |
| 2.1 | 100 train: daraja 9 test har seed'da 4…6 | 2.4 3.3 2.2 4.4 3.4 | ✗ (1000 seed medianasi 4.2) |
| 2.2 | chiziq bilan farq < 1 | 0.4 0.5 0.0 0.5 0.1 | ✓ |
| 2.3 | kamida 3 seed'da daraja 9 biroz yomonroq | 4 / 5 | ✓ |

**Errata:** hozircha boʻsh. Videoda fakt sifatida aytilgan gap keyin xato chiqsa, shu yerga yoziladi va keyingi video
shu tuzatish bilan boshlanadi.

## Videodagi buyruqlar

Hammasi shu papkadan (`lessons/ml-01-map`). Kerak: [uv](https://docs.astral.sh/uv/) (u Python va `numpy==2.5.5` ni
oʻzi oladi). GPU kerak emas.

```bash
cat bashorat.md
git log --reverse --format='%ad  %s' --date=format:%H:%M:%S -- bashorat.md natija.txt
uv run --quiet tajriba.py            # chiziq 1.6 / 4.4 · daraja 9 0.0 / 52.1 · 7.2 soat: 79 → 233
uv run --quiet tajriba.py seeds      # daraja 9: 52.1 3904.1 28487.6 55.9 236.6 · 5 / 5 yomonroq
uv run --quiet tajriba.py ming       # 1000 seed: 3.5 · mediana 1820 · 175 · 418 · 999
cat bashorat2.md
uv run --quiet tajriba.py seeds 100  # daraja 9: 2.4 3.3 2.2 4.4 3.4 · 4 / 5
uv run --quiet tajriba.py ming 100   # mediana 4.2
```

Oʻlchangan: 2026-11-26, Intel i5-12450H, faqat CPU, uv 0.11.21, Python 3.12.13, numpy 2.5.5, go1.26.8 (faqat vazifa).


## Dars matni

**hook**

Oʻnta talabada nol xato qilgan model eng yomon holda yangi talabaga yetmish toʻqqiz oʻrniga ikki yuz oʻttiz uch ball berdi. Bu natijani koʻrishdan oldin besh bashorat yozib, git'ga commit qilgan edik. Ulardan ikkitasi xato chiqdi. Bu ML trekining birinchi darsi va biz modeldan emas, odatdan boshlaymiz: avval bashorat, keyin natija, keyin halol solishtirish.

**trek** · _reja · plan/tracks/ml.yaml: ML: birga oʻrganamiz_

Bu trek tayyor ekspert ma'ruzasi emas, birga oʻrganish jurnali: nimani sinadik, nima buzildi, nima hali chalkash. Modellarni Python va PyTorch'da oʻqitamiz, Go'da esa algoritmlarni noldan qayta yozamiz va modelni serve qilamiz. Ish oʻrinlari yoki AGI haqida bashorat yoʻq, faqat qurganimiz va oʻlchaganimiz.

**xarita** · _reja · plan/tracks/ml.yaml: ML-0 … ML-6_

Avval trek xaritasini koʻraylik. Boshlangʻich chiziqdan keyin klassik ML: chiziqli regressiya, Naive Bayes va har natijani baseline bilan solishtirish. Keyin neyron tarmoqlarni noldan yozamiz va oʻqitish nega buzilishini oʻrganamiz. Soʻng oʻzbek matnida tokenizer va kichik GPT, keyin modellardan muhandis sifatida foydalanish, oxirida maqolalarni birga oʻqiymiz.

**checkpoint** · _reja · plan/tracks/ml.yaml: C1 (ML-1.01) · C9 (ML-4.04) · C11 (ML-5.03)_

Beshta bosqich raqamli checkpoint bilan tugaydi, jami oʻn ikkita. Masalan, C1'da noldan yozgan chiziqli regressiyamiz sklearn bilan milliondan bir aniqlikda mos kelishi kerak. C9'da validatsiya loss'i bigram'dan MLP'ga, undan kichik GPT'ga qarab pasayishi kerak. C11'da Go'dan serve qilingan modelning javob vaqtini oʻlchaymiz. Oʻtmagan checkpoint ham ochiq e'lon qilinadi.

**manba1** · _oʻqidim · karpathy.ai/zero-to-hero · course.fast.ai_

Endi toʻrtta asosiy manbamizni aytaman. Birinchisi, Andrej Karpathy'ning Neural Networks: Zero to Hero kursi: backpropagation asoslaridan GPT'gacha, hammasi kodda. Talabi: yaxshi Python va boshlangʻich matematika, masalan, hosila va Gauss taqsimoti. Ikkinchisi, fast.ai'ning Practical Deep Learning for Coders kursi: bepul, bir yillik dasturlash tajribasi va maktab matematikasi yetarli.

**manba2** · _oʻqidim · cs229.stanford.edu · udlbook.github.io/udlbook_

Uchinchisi, Stanford'ning CS229 kursi: supervised va unsupervised learning, learning theory va reinforcement learning. Tengyu Ma va Andrew Ng yozgan konspektining PDF'ini hozir roʻyxatdan oʻtmasdan yuklab olsa boʻladi. Toʻrtinchisi, Simon Prince'ning Understanding Deep Learning kitobi: MIT Press ikki ming yigirma uchinchi yilda chiqargan, PDF'i bepul, Python notebook'lari bor.

**belgilar**

Har bir faktga uchta belgidan birini qoʻyamiz, ekrandagi manba satrida ham shu belgi turadi. Tekshirdim: shu kompyuterda ishga tushirdik yoki oʻlchadik. Oʻqidim: nomi aytilgan manbada yozilgan, lekin biz takrorlamadik. Taxminim: na oʻlchov, na manba bor, xato chiqishi mumkin. Masalan, manbalar haqidagi gaplar oʻqidim belgisida.

**odat** · _oʻqidim · Feynman, Cargo Cult Science (1974) · Karpathy, A Recipe for Training Neural Networks (2019) · cos.io/prereg_

Nega natijadan oldin bashorat yozamiz? Richard Feynman bir ming toʻqqiz yuz yetmish toʻrtinchi yilgi Caltech nutqida aytgan: birinchi qoida, oʻzingizni aldamang, eng oson aldanadigan odam esa oʻzingizsiz. Karpathy esa yozadi: neyron tarmoqni oʻqitish jim buziladi, koʻpincha shunchaki biroz yomonroq ishlaydi. Shuning uchun bashoratni yozamiz, commit qilamiz, keyin ishga tushiramiz va solishtiramiz, ilm-fandagi preregistration'ga oʻxshab.

**tajriba** · _tekshirdim · tajriba.py · numpy 2.5.5 Polynomial.fit · CS229 notes §8_

Endi tajribaning oʻziga oʻtamiz. Ma'lumotni oʻzimiz yasaymiz, shuning uchun qonuniyatni bilamiz: ball qirq plyus besh karra oʻqilgan soat, ustiga shovqin. Oʻnta talaba bilan oʻqitamiz, bu train. Model koʻrmagan oʻnta yangi talabada tekshiramiz, bu test. Modellar toʻgʻri chiziq va toʻqqizinchi darajali koʻphad: oʻnta nuqtaga oʻnta parametr, ya'ni koʻphad har nuqtadan oʻtadi. Xato esa oʻrtacha necha ballga adashgani.

**bashorat** · _tekshirdim · bashorat.md · git log (commit a9200d1 → natija)_

Ishga tushirishdan oldin besh bashorat yozdik. Chiziqning train xatosi toʻrt atrofida, test xatosi toʻrt bilan besh orasida. Koʻphadning train xatosi nol, test xatosi oʻttiz atrofida, oʻndan yuzgacha. Eng katta xato train oraligʻidan tashqarida boʻladi. Besh seed'ning beshtasida koʻphad yomonroq. git log koʻrsatadi: bashoratdan qirq besh soniya keyin natija fayli commit qilingan. Git vaqtini qoʻlda oʻzgartirsa boʻladi, bu isbot emas, ochiq odat.

**natija** · _tekshirdim · uv run tajriba.py · numpy 2.5.5 · seed 1 · CS229 notes §8_

Har bir nuqtadan oʻtadigan model yangi talabaga yetmish toʻqqiz oʻrniga ikki yuz oʻttiz uch ball berdi. Sintetik ma'lumotda oʻnta talaba bilan ikki model oʻqitdik, bu train, va ularni oʻnta yangi talabada tekshiramiz, bu test. Toʻgʻri chiziq train'da bir butun oʻndan olti, testda toʻrt butun oʻndan toʻrt ballga adashdi. Toʻqqizinchi darajali koʻphad train'da nol, testda qariyb ellik ikki ballga adashdi. Eng katta xato yetti butun oʻndan ikki soat oʻqigan talabada: haqiqiy bali yetmish toʻqqiz, koʻphad ikki yuz oʻttiz uch dedi. Train xatosi kichik, test xatosi katta boʻlsa, buni overfitting deyishadi.

**seeds** · _tekshirdim · uv run tajriba.py seeds_

Bitta natija bitta tasodif boʻlishi mumkin. Seed, ya'ni tasodifiy sonlarning boshlangʻich qiymatini oʻzgartirib, tajribani besh marta qaytaramiz. Chiziq har safar besh balldan kam adashdi. Koʻphad esa taxminan ellik ikki, uch ming toʻqqiz yuz, yigirma sakkiz ming besh yuz, ellik olti va ikki yuz oʻttiz yetti ballga adashdi. Beshtadan beshtasida koʻphad yomonroq.

**solishtir**

Endi bashoratlarni natija bilan solishtiramiz. Birinchisi xato: test xatosi toʻgʻri, lekin train xatosi toʻrt emas, bir butun oʻndan olti. Ikkinchisi toʻgʻri, koʻphad train'da nol. Uchinchisi oraliq boʻyicha toʻgʻri chiqdi. Toʻrtinchisi xato: yetti butun oʻndan ikki soat train oraligʻining ichida. Beshinchisi ham toʻgʻri chiqdi. Bashoratlar hisobi hozircha beshdan uchta.

**nega** · _tekshirdim · uv run tajriba.py ming (seed 1…1000)_

Sababini ming seed bilan qidiramiz, bu tahlil rejimini natijadan keyin qoʻshdik. Chiziqning train xatosi oʻrtacha uch butun oʻndan besh, demak birinchi seed'da shovqin shunchaki kichik chiqqan. Koʻphadning test medianasi ming sakkiz yuz yigirma, oʻndan yuzgacha oraliqqa mingtadan bir yuz yetmish beshtasi tushdi: uchinchi bashoratda omadimiz kelgan. Eng katta xato esa mingtadan toʻrt yuz oʻn sakkiz holatda train oraligʻining ichida chiqdi.

**errata** · _oʻqidim · UDL_Errata.pdf (udlbook.github.io/udlbook)_

Xato bashorat uyat emas, u ham ma'lumot: hisobga yoziladi va bosqich oxirida e'lon qilinadi. Errata boshqa narsa: videoda fakt deb aytgan gapimiz xato chiqsa, keyingi video shu tuzatish bilan boshlanadi. Understanding Deep Learning kitobining ham ochiq errata roʻyxati bor. Bizning roʻyxatimiz esa hozircha boʻsh.

**savol**

Endi bashorat qilish navbati sizda. Train talabalarini oʻntadan yuztaga koʻpaytiramiz, test esa yana oʻnta yangi talaba. Koʻphadning test xatosi qancha boʻladi va u chiziqdan yomonroqmi? Javobingizni yozib qoʻying, keyin videoni davom ettiring. Biz oʻz bashoratimizni ishga tushirishdan oldin commit qildik.

**javob** · _tekshirdim · bashorat2.md · uv run tajriba.py seeds 100 · ming 100_

Bizning bashoratimiz: koʻphad test xatosi har seed'da toʻrtdan oltigacha, chiziq bilan farq bir balldan kam, kamida uchta seed'da koʻphad yomonroq. Natijada farq haqiqatan bir balldan kam, koʻphad beshtadan toʻrttasida yomonroq. Lekin xatolar ikki butun oʻndan ikkidan toʻrt butun oʻndan toʻrtgacha, demak birinchi bashorat yana xato, ming seed'da esa mediana toʻrt butun oʻndan ikki. Umumiy hisob sakkizdan beshta.

**homework**

Uy vazifasi, toʻrt darajada, har birida avval bashorat yozing. Oson: shovqinni nolga tushiring. Oʻrta: uchinchi darajali koʻphad qoʻshing va seed'lar bilan solishtiring. Qiyin: ming seed'dan bittasida koʻphad yutdi, qaysi biri va nega? Guru: chiziqni Go'da noldan hisoblab, numpy bilan solishtiring. Javoblar GitHub'dagi kurs repository'sida, link pastda.

**bilmaymiz** · _oʻqidim · CS229 Lecture Notes (Ma, Ng), §8.2 The double descent phenomenon_

Har video hali bilmaganlarimiz bilan tugaydi. Birinchidan, uchinchi seed'da koʻphad nega aynan yigirma sakkiz ming ballga adashganini tushuntira olmaymiz. Ikkinchidan, koʻphad chiziqni yutgan yagona seed'da nima boshqacha boʻlganini bilmaymiz. Uchinchidan, CS229 konspektiga koʻra, parametrlar juda koʻp boʻlganda test xatosi yana tushishi mumkin, buni double descent deyishadi: oʻqidik, lekin hali koʻrmadik.

**outro**

Bugun trek xaritasini, checkpoint'larni, uchta ishonch belgisini va natijadan oldin bashorat yozish odatini koʻrdik. Keyingi darsda loss nima, fit nima va modelni nega u koʻrmagan ma'lumotda tekshirish shartligini koʻramiz. Darsni koʻrganingiz uchun rahmat, koʻrishguncha.

## Manbalar

Hammasi 2026-11-26 da rasmiy sahifalardan oʻqildi.

- Andrej Karpathy, *Neural Networks: Zero to Hero* — https://karpathy.ai/zero-to-hero.html («building neural networks,
  from scratch, in code»; backprop'dan GPT'gacha; talab: Python + boshlangʻich matematika). micrograd (MIT litsenziya):
  https://github.com/karpathy/micrograd
- fast.ai, *Practical Deep Learning for Coders* (2022) — https://course.fast.ai/ (bepul; «a year of experience is
  enough»; maktab matematikasi)
- Stanford *CS229: Machine Learning* — https://cs229.stanford.edu/ · Tengyu Ma va Andrew Ng konspekti:
  https://cs229.stanford.edu/main_notes.pdf (§8 overfitting taʼrifi, §8.1 bias-variance, §8.2 double descent)
- Simon J.D. Prince, *Understanding Deep Learning* (MIT Press, 2023) — https://udlbook.github.io/udlbook/ (bepul PDF,
  Python notebook'lar, 8-bob «Measuring performance»), errata: https://github.com/udlbook/udlbook/raw/main/UDL_Errata.pdf
- Richard Feynman, *Cargo Cult Science* (Caltech, 1974) — http://calteches.library.caltech.edu/51/2/CargoCult.htm
- Andrej Karpathy, *A Recipe for Training Neural Networks* (2019) — https://karpathy.github.io/2019/04/25/recipe/
- Center for Open Science, *Preregistration* — https://www.cos.io/initiatives/prereg
- scikit-learn, *Underfitting vs. Overfitting* (shu gʻoyaning klassik rasmi: daraja 1, 4, 15) —
  https://scikit-learn.org/stable/auto_examples/model_selection/plot_underfitting_overfitting.html
- numpy `Polynomial.fit` — https://numpy.org/doc/stable/reference/generated/numpy.polynomial.polynomial.Polynomial.fit.html
- git `--date` / `GIT_COMMITTER_DATE` — https://git-scm.com/docs/git-commit#_date_formats

## Uy vazifasi

Har birida **avval bashorat yozing** (iloji boʻlsa, commit qiling), keyin ishga tushiring.

1. **Oson.** `talabalar` ichida shovqinni nolga tushiring: `rng.normal(0, 0, n)`. Ikki model nima qiladi?
2. **Oʻrta.** `MODELLAR` ga `("daraja 3", 3)` qoʻshing va `uv run tajriba.py seeds` bilan solishtiring.
3. **Qiyin.** `ming` rejimida daraja 9 1000 tadan 999 tasida yomonroq. Qaysi seed'da u chiziqni yutdi va nega?
4. **Guru.** Toʻgʻri chiziqni Go'da formuladan noldan hisoblang va numpy bilan solishtiring (C1'ning kichik versiyasi).

### Javoblar

1. Ikkalasi ham train'da va test'da `0.0`. Shovqin yoʻq boʻlsa, 10 nuqtadan oʻtadigan yagona ≤9-darajali koʻphad shu
   chiziqning oʻzi, yodlab oladigan shovqin yoʻq. **Lekin** oxirgi qator `daraja 9 testda yomonroq: 5 / 5` deydi: aniq
   qiymatlar chiziqda 4e-15…3e-14, daraja 9 da 1e-13…5e-11 (float yaxlitlash xatosi). Ekrandagi `0.0` ikkalasini
   tenglashtiradi, `>` esa yoʻq. Saboq: float'larni tolerantlik bilan solishtiring (`np.isclose`).
2. `daraja 3`: `4.5 2.8 5.0 4.7 2.7`, chiziq: `4.4 3.3 4.8 3.4 2.8`. 4 parametr 10 nuqtadan kam: ba'zan biroz yaxshi,
   ba'zan biroz yomon, lekin daraja 9 kabi portlamaydi.
3. Seed **669**: chiziq 4.00, daraja 9 2.90. U yerda 10 ta test soatining hammasi train oraligʻi (0.69…9.32) ichida.
   Bu yetarli sabab emas: 1000 seed'ning 25 % ida ham hamma test nuqtasi ichkarida, ularda ham daraja 9 medianasi 204.
   Toʻliq sababini **bilmaymiz** (videodagi «hali bilmaganlarimiz»). Topish uchun `ming` ichida
   `if d9_test <= chiziq_test: print(seed)` qoʻshing.
4. `uv run --quiet homework/nuqtalar.py && go run ./homework/guru homework/nuqtalar.csv`:
   ```
   numpy: ball = 37.520865264908 + 5.338824873127 × soat
   go:    ball = 37.520865264908 + 5.338824873127 × soat
   ```
   12 xonagacha bir xil. Haqiqiy qonuniyat `40 + 5 × soat` edi: 10 ta shovqinli nuqtadan chiziq ham uni aniq topmaydi.
   Bu `go run` papka nomi bilan ishlaydi, shuning uchun buyruqni `lessons/ml-01-map` dan bering.
