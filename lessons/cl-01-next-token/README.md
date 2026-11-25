# CL-1.01 · LLM nima qiladi: keyingi token, context oynasi va ishonchli xato

Claude bilan muhandislik trekining birinchi darsi · Keyingi dars: CL-1.02 · Claude.ai'dan toʻgʻri foydalanish: model tanlash, Projects va fayllar

| Fayl | Nima |
|---|---|
| `korpus.txt` | modelning butun "bilimi": 6 ta toʻgʻri gap (token = soʻz, nuqta ham token) |
| `main.go` | kichik til modeli: n-gram (3 tokengacha kontekst, backoff), `-oyna` (javob paytidagi context window, 1..3), `-t` (temperature), `-n`, `-seed`, `-korpus` |
| `tokenlar/tokenlar.py` | bir xil matn (Inson huquqlari umumjahon deklaratsiyasi) inglizcha va oʻzbekcha nechta token boʻlishi — **ochiq** tokenizer'larda (tiktoken `o200k_base`, `cl100k_base`) |
| `tokenlar/eski_claude.mjs` | xuddi shu, Anthropic'ning **eski** modellar uchun chiqargan ochiq tokenizer'ida (`@anthropic-ai/tokenizer` 0.0.4) |
| `tokenlar/udhr_*.txt` | deklaratsiya matni ([eric-muller/udhr](https://github.com/eric-muller/udhr), `data/udhr/udhr_eng.xml`, `udhr_uzn_latn.xml`; sarlavha va paragraflar, `‐` U+2010 → `-`) |
| `homework/` | uy vazifasi javoblari: `-chegara` ("bilmayman") va `-k` (top-k) bayroqlari, `korpus_rust.txt` |

## Videodagi buyruqlar
Hammasi shu papkadan (`lessons/cl-01-next-token`) ishga tushiriladi. `tokenlar.py` uchun [uv](https://docs.astral.sh/uv/) kerak
(u `tiktoken` ni oʻzi oʻrnatadi), `eski_claude.mjs` uchun avval `cd tokenlar && npm install`.
```bash
go run . python tilini                       # [tilini] → google (0.50) … python tilini google yaratgan .  ← notoʻgʻri
cat korpus.txt
uv run --quiet tokenlar/tokenlar.py soz      # we are learning: 3 token, oʻrganmoqdamiz: 8 (ʻ bilan), 7 (' bilan)
uv run --quiet tokenlar/tokenlar.py til && node tokenlar/eski_claude.mjs   # 1.15x belgi, 1.94x o200k, 2.54x Claude 3 gacha
go run . -oyna 1 python tilini               # oyna 1: python koʻrinmaydi → google
go run . -oyna 2 python tilini               # oyna 2: [python tilini] → guido  ✓
go run . rust                                # [rust] → ?  korpusda yoʻq
go run . -oyna 2 rust tilini                 # backoff: [tilini] → google … rust tilini google yaratgan .
for t in 0.5 1 2; do printf 't=%-4s %s\n' $t "$(go run . -t $t python tilini | head -1)"; done   # google 0.67 / 0.50 / 0.41
for t in 0 1 2; do echo "t=$t: $(go run . -t $t -n 10000 python tilini | grep -c guido) / 10000 guido"; done   # 0 / 2418 / 2891
for t in 0 1 2; do echo "t=$t: $(go run . -t $t -n 10000 go tilini | grep -c google) / 10000 google"; done     # 10000 / 5027 / 4105
```
Natijalar: `go1.26.8 linux/amd64`, tiktoken 0.14.0, `@anthropic-ai/tokenizer` 0.0.4, Node 24.10, 2026-11-25.
`-seed` (standart 1) bir xil boʻlsa, sampling natijasi ham bir xil (`math/rand/v2` PCG). Boshqa seed'da sonlar biroz
oʻzgaradi (10 000 namunada ±1 foiz atrofida); 100 namunada esa shovqin katta — shuning uchun videoda 10 000.
Qoʻshimcha: `tokenlar.py cl100k` (2011 / 5038, 2,51×), `tokenlar.py apostrof`.

## Dars matni

**Hook.** `korpus.txt` dagi oltita gapning hammasi toʻgʻri (qiziq tomoni: Guido van Rossum 2005–2012 yillarda Google'da ishlagan, shuning uchun bu xato ayniqsa ishonarli eshitiladi). Modelga `python tilini` deb boshlab beramiz, u davom ettiradi:
`python tilini google yaratgan .` — notoʻgʻri (Python'ni Guido van Rossum yaratgan). Claude'ning asosidagi model ham
keyingi soʻzni bashorat qilishga oʻrgatilgan (Anthropic Glossary, *Pretraining*: "autoregressive language models (such as
Claude's underlying model) are pretrained to predict the next word, given the previous context of text").

**Model qanday ishlaydi.** Bu n-gram model (Jurafsky & Martin, SLP3, 3-bob). `oqit` korpusdan **bir marta** oʻrganadi:
har bir token uchun undan oldingi 1..3 ta tokenni (kontekstni) yozib qoʻyadi va sanaydi. `-oyna` esa javob paytida model
matnning oxirgi nechta tokeniga qarashini belgilaydi — model oʻzgarmaydi, faqat u koʻradigan qism oʻzgaradi. `taqsimot`
oyna ichidagi eng uzun *tanish* kontekstni topadi (topilmasa, qisqaroq kontekstga oʻtadi — **backoff**) va sanoqlarni
ehtimollikka aylantiradi: `p ~ soni^(1/t)`. Silliqlash (smoothing) yoʻq: korpusda koʻrilmagan davomning ehtimolligi
aniq 0 (shuning uchun `graydon` hech qachon chiqmaydi; haqiqiy LLM esa lugʻatdagi har bir tokenga noldan katta
ehtimollik beradi). Teng sanoqda korpusda oldin uchragan
token oldinda turadi. Bu Claude emas: Claude — koʻp parametrli katta til modeli, butun kontekstga qaraydi va yordamchi boʻlishga
qoʻshimcha oʻqitilgan (Glossary: *"Claude is not a bare language model; it has already been fine-tuned to be a helpful
assistant"*, RLHF). Lekin chiqish shakli bir xil: har qadamda keyingi token uchun ehtimolliklar, keyin bitta tanlov.

**Token.** Anthropic Glossary, *Tokens*: token soʻz, soʻz boʻlagi, harf yoki bayt boʻlishi mumkin; bir token nechta
harfga toʻgʻri kelishi model va tilga bogʻliq; eski modellarda ≈3,5 inglizcha harf; Claude 4.7 va keyingi modellar yangi
tokenizer ishlatadi va bir xil matnga ≈30% koʻproq token beradi (Opus 5.5 migratsiya qoʻllanmasi: 1–1,35×).
Hozirgi Claude tokenizer'i ochiq emas; aniq son uchun Anthropic token counting API'ni yoki javobdagi `usage` ni tavsiya
qiladi. Shuning uchun biz **ochiq** tokenizer'larda oʻlchadik — bular Claude'ning aniq soni emas:

| Tokenizer | Inglizcha | Oʻzbekcha | Farq |
|---|---|---|---|
| tiktoken `o200k_base` (OpenAI, ochiq) | 2011 | 3911 | 1,94× |
| tiktoken `cl100k_base` (OpenAI, ochiq) | 2011 | 5038 | 2,51× |
| `@anthropic-ai/tokenizer` 0.0.4 (Anthropic, eski modellar; README: "As of the Claude 3 models, this algorithm is no longer accurate") | 2068 | 5249 | 2,54× |
| belgilar soni | 10638 | 12241 | 1,15× |

`o200k_base` da `we are learning` — 3 token, `oʻrganmoqdamiz` — 8 token (`o | ʻ | r | gan | mo | q | dam | iz`),
ASCII apostrof bilan 7 (`o | 'r | …`). BPE
(Sennrich va boshq., 2016) koʻp uchragan ketma-ketliklarni bitta tokenga birlashtiradi; oʻzbekcha soʻzning maydalanishi,
ehtimol, tokenizer oʻrgatilgan matnlarda oʻzbekcha kam boʻlganidan (bu xulosa, oʻlchov emas). Apostrof turi
(`ʻ` U+02BB / `'` / `‘`) butun matnda 3911 / 3892 / 3938 token beradi: farq kichik (`tokenlar.py apostrof`).
Chuqurroq: ML-4.01 (BPE oʻzbek matnini qanday kesadi).

**Context window.** Anthropic, *Context windows*: model javob yozayotganda koʻra oladigan butun matn, javobning oʻzi ham
kiradi; "working memory", oʻqitilgan maʼlumot emas; token koʻpaygan sari aniqlik va eslash pasayadi (*context rot*).
Hozirgi Claude modellari (Fable 5.2, Opus 5.5, Sonnet 5.5, Haiku 5.5), API'da: 1M token ≈ 555 ming soʻz (yangi
tokenizer'da); boshqa modellar 200k (Models overview). Chat ilovalari (claude.ai) uzun suhbatda oynani "first in, first
out" tarzda boshqarishi mumkin: eski xabarlar chiqib ketadi (Context windows, 1-izoh).
`-oyna 1` da model `google` ni tanlayotganda faqat `[tilini]` ni koʻradi, `python` oynadan chiqib ketgan. `-oyna 2` da
`[python tilini]` → `guido`. Claude'da bunday qisqa savol 1M oynadan chiqib ketmaydi: uning xatolari boshqa joydan
keladi. Korpusda yoʻq `rust` yolgʻiz berilsa, model jim qoladi (`→ ?`); `rust tilini` da esa backoff notanish `rust` ni
tashlab, tanish `[tilini]` ga tayanadi va yana `google` deydi. Jim turish oʻrniga taxmin — aynan backoff qarori.
Aslida Rust'ni Graydon Hoare boshlagan (2006, shaxsiy loyiha; 2009 dan Mozilla homiyligida).

**Sampling va temperature.** `t = 0`: har doim eng ehtimolli token (greedy). `t > 0`: ehtimollikka mos qurʼa
(`tanla`). `t = 0.5` sanoqlarni kvadratga koʻtaradi (google 0,67), `t = 2` kvadrat ildiz oladi (google 0,41).
10 000 namuna, `python tilini` (toʻgʻri: guido): t=0 → 0, t=1 → 2418, t=2 → 2891.
`go tilini` (toʻgʻri: google): t=0 → 10 000, t=1 → 5027, t=2 → 4105. Oyna 1 da ikkala savol uchun kontekst bir xil
(`[tilini]`), faqat toʻgʻri javob har xil: temperature haqiqatni bilmaydi, u faqat taqsimotni tekislaydi.
Claude'da: Glossary, *Temperature*: "Even with temperature set to 0, the results will not be fully deterministic".
Opus 5.5 da `temperature`, `top_p`, `top_k` olib tashlangan: "Use prompting to guide the model's behavior" (migratsiya qoʻllanmasi).

**Nega ishonchli, lekin notoʻgʻri?**
1. Model har doim keyingi tokenni berishi kerak. Anthropic: "language model training incentivizes hallucination: models
   are always supposed to give a guess for the next word."
2. Ehtimollik matnda koʻrinmaydi: bizning model `google` ni 0,50 bilan tanladi, gapda esa "balki" yoʻq.
3. Baholash taxminni ragʻbatlantiradi: Kalai va boshq. (2025): "language models hallucinate because the training and
   evaluation procedures reward guessing over acknowledging uncertainty" — xato javobga ball ayirilmaydigan imtihon kabi.
4. Claude ichida (Anthropic, *Tracing the thoughts…*, 2025, **Claude 3.5 Haiku** ustida): Claude "bilmayman" deyishga
   ham oʻrgatilgan, lekin "imperfect"; "bilmayman" sxemasi odatda yoqilgan; Michael Jordan kabi tanish nom haqida
   soʻralsa "known entity" belgisi uni oʻchiradi; notanish nom ("Michael Batkin") — javob yoʻq; nom tanish, lekin bilim
   yoʻq boʻlsa, sxema adashib oʻchadi va model "plausible — but unfortunately untrue" javob toʻqiydi.

**Nima qilamiz** (Anthropic, *Reduce hallucinations*): "bilmayman" deyishga ruxsat bering; kerakli hujjatni bering va
faqat undan javob soʻrang ("External knowledge restriction"); uzun hujjatda (>20k token) avval soʻzma-soʻz iqtibos;
har bir daʼvoga **berilgan hujjatdan** iqtibos, topilmasa daʼvo olib tashlansin (xotiradan "manba" soʻrash yangi
toʻqilgan manbaga olib kelishi mumkin); bir savolni bir necha marta berib solishtiring (best-of-N) — lekin bir xil
javob isbot emas (`t=0` da 10 000 marta bir xil xato). Hujjat: bu usullar xatoni kamaytiradi, butunlay yoʻqotmaydi.
Kodni ishga tushiring va test qiling — CL-1.05.

## Manbalar
- Anthropic Glossary (Tokens, Temperature, Pretraining, Context window): https://platform.claude.com/docs/en/about-claude/glossary
- Anthropic, Context windows: https://platform.claude.com/docs/en/build-with-claude/context-windows
- Anthropic, Models overview: https://platform.claude.com/docs/en/models/overview
- Anthropic, Claude Opus 5.5 migration guide: https://platform.claude.com/docs/en/models/opus-5-5/migration-guide
- Anthropic, Reduce hallucinations: https://platform.claude.com/docs/en/test-and-evaluate/strengthen-guardrails/reduce-hallucinations
- Anthropic, Tracing the thoughts of a large language model (2025-03-27): https://www.anthropic.com/research/tracing-thoughts-language-model
- Anthropic TypeScript tokenizer (eski modellar): https://github.com/anthropics/anthropic-tokenizer-typescript
- A. T. Kalai, O. Nachum, S. S. Vempala, E. Zhang, *Why Language Models Hallucinate* (2025): https://arxiv.org/abs/2509.04664
- R. Sennrich, B. Haddow, A. Birch, *Neural Machine Translation of Rare Words with Subword Units* (2016): https://arxiv.org/abs/1508.07909
- tiktoken: https://github.com/openai/tiktoken
- Korpusdagi faktlar: Go FAQ (tarix, 2009-11-10 ochiq kod) https://go.dev/doc/faq#history · Python History and License
  https://docs.python.org/3/license.html · D. Ritchie, *The Development of the C Language* https://www.bell-labs.com/usr/dmr/www/chist.html ·
  Linux README https://www.kernel.org/doc/html/latest/admin-guide/README.html · Google, "Google Gets the Message, Launches Gmail" (2004)
  http://googlepress.blogspot.com/2004/04/google-gets-message-launches-gmail.html · Rust: https://en.wikipedia.org/wiki/Rust_(programming_language) · Guido Google'da (2005–2012): https://en.wikipedia.org/wiki/Guido_van_Rossum
- D. Jurafsky, J. H. Martin, *Speech and Language Processing* (3-nashr loyihasi), 3-bob "N-gram Language Models" (backoff): https://web.stanford.edu/~jurafsky/slp3/

## Uy vazifasi
1. **Oson.** Korpusga `rust tilini graydon yaratgan .` gapini qoʻshing. `-oyna 2 rust tilini` va `-oyna 1 python tilini` nima deydi?
2. **Oʻrta.** Oʻzingiz yozgan bir paragrafni ikki tilda `tokenlar/` ga qoʻyib, `tokenlar.py` bilan oʻlchang.
3. **Qiyin.** Modelga "bilmayman" deyishni oʻrgating: eng yuqori ehtimollik chegaradan past boʻlsa, javob bermasin.
4. **Guru.** Top-k sampling qoʻshing (faqat eng ehtimolli `k` ta nomzoddan tanlash) va 100 namunada temperature bilan solishtiring.

### Javoblar
1. `go run . -korpus homework/korpus_rust.txt -oyna 2 rust tilini` → `[rust tilini] graydon 1.00` → `rust tilini graydon yaratgan .` ✓.
   `go run . -korpus homework/korpus_rust.txt python tilini` → `[tilini] → google  google 0.40  guido 0.20  dennis 0.20  graydon 0.20` →
   yana `python tilini google yaratgan .`: bilim qoʻshildi, lekin oyna 1 boʻlsa `python` baribir koʻrinmaydi. Bilim ham, kontekst ham kerak.
   Claude'ga yangi bilimni korpus orqali emas (Claude API fine-tuning bermaydi), kerakli hujjatni promptga qoʻyib berasiz.
2. Ikki faylni `udhr_eng.txt` / `udhr_uzn_latn.txt` oʻrniga qoʻying (yoki skriptdagi fayl nomini oʻzgartiring). Oʻzbekcha matn odatda
   koʻproq token oladi; deklaratsiyada `o200k_base` bilan 1,94×. Apostrof: `uv run --quiet tokenlar/tokenlar.py apostrof`.
3. `homework/main.go`, `-chegara`: `go run ./homework -chegara 0.6 python tilini` → `[tilini] → bilmayman eng yuqorisi google 0.50 < 0.60`;
   `-oyna 2` bilan → `guido 1.00`, javob beradi. Narxi: `go run ./homework -chegara 0.6 go tilini` ham "bilmayman" deydi,
   garchi toʻgʻri javob (`google`) eng ehtimolli boʻlsa ham. Jim turish xatoni kamaytiradi, lekin toʻgʻri javoblarni ham yoʻqotadi.
4. `homework/main.go`, `-k`: `go run ./homework -t 1 -k 2 -n 100 python tilini | sort | uniq -c | sort -rn` →
   `40 google yaratgan`, `33 guido yaratgan`, `27 google ochiq kodli qilgan` (`dennis` umuman yoʻq); `-k 1` = greedy (100 × google);
   `-k 0` = oddiy sampling (33 / 24 / 22 / 21). Temperature taqsimot shaklini oʻzgartiradi, top-k esa dumini kesadi.
