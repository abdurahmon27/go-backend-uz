# CL-1.02 · Claude.ai'dan toʻgʻri foydalanish: model tanlash, Projects va fayllar

Oldingi dars: [CL-1.01 · LLM nima qiladi: keyingi token, context oynasi va ishonchli xato](../cl-01-next-token/) ·
Keyingi dars: CL-1.03 · Prompt — bu texnik topshiriq: kontekst, misol, cheklov

| Fayl | Nima |
|---|---|
| `main.go` | **skaner**: matn, log yoki papkada API kalit, parol, token, karta, telefon, email va `.env` fayllarini topadi; `-yashir` ularni `[YASHIRIN:tur]` bilan almashtiradi; `-qatiy` topilma boʻlsa chiqish kodi 1 |
| `main_test.go` | qoidalar uchun table test (topishi va **topmasligi** kerak boʻlgan qatorlar: buyurtma raqami, `0000 0000 0000 0000`, git hash, UUID, Go nomi) |
| `olchov_test.go` | oʻlchov: tasodifiy 24/32/40 belgili kalitlarning necha foizini «nomaʼlum kalit» qoidasi topadi (`go test -run Olchov -v`) |
| `misol.sh` | videodagi misollarni yaratadi: `misol/xato.log`, `misol/bozor/` (`.env`, `.env.example`, …). Kalitlar **soxta** |
| `loyiha-korsatma.md` | Claude.ai Project uchun koʻrsatma (*project instructions*) namunasi — shu kurs uchun |
| `tekshiruv-royxati.md` | chatga yuborishdan oldin 30 soniyalik tekshiruv va kalit ketib qolsa nima qilish |
| `homework/` | uy vazifasi javoblari: Telegram bot token qoidasi + test, `pre-commit` git hook |

Provayder prefiksli soxta kalitlar (`sk-ant-…`, `ghp_…`, `AKIA…`) repo'da toʻliq satr boʻlib turmaydi: `misol.sh` va
testlar ularni boʻlaklardan yigʻadi. Shunda GitHub secret scanning ularni haqiqiy kalit deb oʻylamaydi (GIT-1.02 dagi
`soxta_…` kalit bilan bir xil gʻoya). Qolgan soxta qiymatlar (`Qovun2026`, `Xq7L…`) hech qaysi provayder formatiga mos emas.

## Videodagi buyruqlar

Hammasi shu papkadan (`lessons/cl-02-claude-ai`) ishga tushiriladi.

```bash
./misol.sh && cat misol/xato.log         # 8 qatorli log: parol, Claude kaliti, email, telefon, karta, toʻlov kaliti
go run . misol/xato.log                  # jami 6: 3 ta sir, 3 ta shaxsiy maʼlumot
ls -A misol/bozor && go run . misol/bozor   # .env nomi bilan + ichidagi 2 sir; .env.example toza
go run . -yashir misol/xato.log          # log qaytadi, sirlar oʻrnida [YASHIRIN:…]
go run . $(go env GOROOT)/src/net/http | grep -v ^jami | awk '{print $3}' | sort | uniq -c | sort -rn
#   144 nomaʼlum-kalit · 15 url-parol · 7 parol · 4 email   (115 fayl, yolgʻon signallar)
go run . $(go env GOROOT)/src/net/http | grep nomaʼlum | cut -d: -f1 | xargs -n1 basename | sort | uniq -c | sort -rn
#   81 h2_bundle.go · 43 testcert.go · 8 serve_test.go · 7 transport_test.go · 3 request_test.go · 1 · 1
go test -run Olchov -v . | grep belgi
#   24 belgi: 8732 / 10000 · 32 belgi: 9732 / 10000 · 40 belgi: 9930 / 10000
go test ./...                            # qoidalar + uy vazifasi testlari
```
Natijalar: `go1.26.8 linux/amd64`, 2026-11-25.

Natija misoli:
```
misol/xato.log:2       sir     url-parol       post…6@
misol/xato.log:3       sir     anthropic-kalit sk-a…aD
misol/xato.log:4       shaxsiy email           aziz…om
misol/xato.log:5       shaxsiy telefon         +998…22
misol/xato.log:5       shaxsiy karta           4111…11
misol/xato.log:6       sir     nomaʼlum-kalit  Xq7L…Ge
jami 6: 3 ta sir, 3 ta shaxsiy maʼlumot. Chatga yuborishdan oldin olib tashlang.
```
Skaner qiymatni hech qachon toʻliq chiqarmaydi: 16 belgidan uzun boʻlsa faqat boshi va oxiri (`sk-a…aD`), qisqasi
toʻliq yulduzcha. Matnning oʻzini almashtirish uchun `-yashir` ishlating.

## Skaner qanday ishlaydi

**1. Maʼlum formatlar** (`qoidalar`, `main.go`):

| Qoida | Namuna | Manba |
|---|---|---|
| `anthropic-kalit` | `sk-ant-…` | Anthropic docs: "long-lived `sk-ant-...` API keys" ([Workload identity federation](https://platform.claude.com/docs/en/build-with-claude/workload-identity-federation)), Admin key `sk-ant-admin…` ([Administration API](https://platform.claude.com/docs/en/build-with-claude/administration-api)) |
| `github-token` | `ghp_`, `gho_`, `ghu_`, `ghs_`, `ghr_`, `github_pat_` | [GitHub: About authentication — token prefixes](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-authentication-to-github#githubs-token-formats) |
| `aws-kalit` | `AKIA…`, `ASIA…` (20 belgi) | [AWS IAM identifiers — unique ID prefixes](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html) |
| `maxfiy-kalit` | `-----BEGIN … PRIVATE KEY-----` | PEM sarlavhasi ([RFC 7468](https://www.rfc-editor.org/rfc/rfc7468)) |
| `url-parol` | `postgres://user:parol@…` | URL ichidagi `userinfo` ([RFC 3986 §3.2.1](https://www.rfc-editor.org/rfc/rfc3986#section-3.2.1)) |
| `parol` | `password=…`, `DB_PASSWORD=…`, `DB_PASS=…`, `parol: …`, `token=…` | umumiy kalit soʻzlar |
| `email`, `telefon` (`+998 90 …`, `+998 (90) …`) | | Oʻzbekiston kodi +998 ([ITU-T E.164 roʻyxati](https://www.itu.int/pub/T-SP-E.164D)). `90 123 45 67` kabi mahalliy shakl **ushlanmaydi** (har qanday 9 xonali son bilan adashmaslik uchun) |
| `karta` | 16 xonali son **va** Luhn tekshiruvi, hammasi bir xil raqam emas | Luhn formulasi (ISO/IEC 7812-1); shu sababli `id=2026112510140712` va `0000 0000 0000 0000` karta deb hisoblanmaydi |

**2. Nomaʼlum format** (`tasodifiy`): kamida 24 ta harf/raqamdan iborat soʻzda katta harf, kichik harf va kamida 2 ta
raqam boʻlsa, Shannon entropiyasi hisoblanadi: `H = −Σ p·log₂p`. 4 bit/belgidan oshsa — kalitga oʻxshaydi.
Git hash bu sinovga kirmaydi: unda katta harf yoʻq, 16 xil belgida (0–9, a–f) esa H ≤ log₂16 = 4.

### Chegaralar (oʻlchangan)

| Qoida varianti (raqam ≥ N, H > X) | 24 belgi topildi | 32 belgi | 40 belgi | `net/http` dagi «nomaʼlum kalit» |
|---|---|---|---|---|
| N = 3, X = 4.0 | 73 % | 91 % | 97 % | 128 |
| **N = 2, X = 4.0 (hozirgi)** | **87 %** | **97 %** | **99 %** | **144** |
| N = 1, X = 4.0 | 94 % | 99,6 % | 99,9 % | 196 |
| N = 1, X = 3.7 | 98 % | 99,7 % | 99,9 % | 488 |

(Tasodifiy base62 kalitlar, har biri 10 000 ta, `math/rand/v2` PCG seed 1,2; go1.26.8.) Bu — sozlash narxi: koʻproq
kalit topsangiz, koʻproq yolgʻon signal olasiz. Hozirgi variant 24 belgili kalitlarning taxminan har sakkiztadan
bittasini oʻtkazib yuboradi.

**Yolgʻon signallar.** Go'ning `net/http` paketida (115 fayl) skaner 144 ta «nomaʼlum kalit» topdi: 43 tasi
`internal/testcert/testcert.go` dagi test sertifikati va test RSA kaliti (Go uni `TESTING KEY` sarlavhasi bilan yozadi
va faqat ishga tushganda `PRIVATE KEY` ga almashtiradi — shuning uchun PEM qoidamiz uni ushlamadi, entropiya ushladi),
81 tasi `h2_bundle.go` dagi `http2defaultRFC9218Priority` kabi raqamli uzun nomlar, 20 tasi test nomlari. Yana 26 ta
`url-parol`/`parol`/`email` topilmasining 22 tasi test fayllarida. Bu skaner **ehtiyotkor**: ortiqcha ogohlantirish
arzon, sizgan kalit qimmat. Va aksincha: «topilma yoʻq» — «xavfsiz» degani emas.
Jiddiy loyihada tayyor asbobni ishlating (Anthropic [gitleaks](https://github.com/gitleaks/gitleaks) ni tavsiya qiladi)
va GitHub'da secret scanning / push protection'ni yoqing.

## Dars matni

**Hook.** Xatoni tushunish uchun Claude'ga server logini yopishtirmoqchisiz, unda esa baza paroli va Claude API kaliti
bor. Anthropic yozishicha, kalit sizishining eng koʻp sabablaridan biri — uni ochiq repository'ga yoki begona asbobga
tasodifan qoʻyish ([API Key Best Practices](https://support.claude.com/en/articles/9767949-api-key-best-practices-keep-your-keys-safe-and-secure)).
Chatga yopishtirilgan matn ham kompyuterdan chiqib, serverda saqlanadi (pastda: qancha).

**Modellar** ([Models overview](https://platform.claude.com/docs/en/about-claude/models/overview), 2026-11-25):

| Model | Anthropic taʼrifi | API narxi (kirish / chiqish, $ / 1M token) |
|---|---|---|
| Claude Fable 5.2 | demanding reasoning and long-horizon agentic work | 8 / 40 |
| Claude Opus 5.5 | long-running agentic coding and knowledge work | 4 / 20 |
| Claude Sonnet 5.5 | the best combination of speed and intelligence | 2 / 10 |
| Claude Haiku 5.5 | high-volume, latency-sensitive tasks such as classification, extraction, and routing | 0.10 dan / 0.50 dan |

Hammasida context window 1M token (API; pullik Claude.ai rejalarida ham chatda 1M —
[8606394](https://support.claude.com/en/articles/8606394-how-large-is-the-context-window-on-paid-claude-plans)).

**Tarif** ([claude.com/pricing](https://claude.com/pricing), [Fable on your plan](https://support.claude.com/en/articles/15424964-claude-fable-models-on-your-plan)):
Free — Sonnet, Haiku, Projects 5 tagacha; Pro — + Opus, Fable faqat *usage credits* bilan; Max va Ultra — Fable haftalik
limitning 50 % igacha. Limitga xabar uzunligi, fayllar, model va **effort** taʼsir qiladi
([Usage and length limits](https://support.claude.com/en/articles/11647753-how-do-usage-and-length-limits-work)).
Oʻzbekiston qoʻllab-quvvatlanadigan hududlar roʻyxatida bor ([Supported regions](https://www.anthropic.com/supported-countries), 2026-11-12 tahriri).

**Qaysi modelni tanlash.** Anthropic dasturchilar (API ilovalari) uchun ikki yoʻlni koʻrsatadi — mantiqi chatda ham
ishlaydi
([Choosing a model](https://platform.claude.com/docs/en/about-claude/models/choosing-a-model)): *efficiency-first* (Haiku
bilan boshlab, yetmasa kattasiga oʻtish) va *capability-first* (Opus bilan boshlab, keyin arzonlashtirish; Opus `xhigh`/`max`
effort'da ham yetmasa — Fable). Hujjat yana aytadi: effort'ni sozlash koʻpincha model almashtirishdan yaxshiroq richag.
**Kursdagi qoidamiz** (bizniki, hujjat emas): xato xabarini tushunish — Haiku; bitta funksiya yoki test — Sonnet; koʻp
faylli oʻzgarish, qiyin bug — Opus; Opus yuqori effort bilan ham uddalamasa — Fable.

**Model menyusi** (sxema, [8664678](https://support.claude.com/en/articles/8664678-change-the-model-effort-and-thinking-settings)):
```
[ xabar maydoni ............ ] [Sonnet 5.5 · High ▾] [↑]
                                  ├─ Haiku 5.5 / Sonnet 5.5 / Opus 5.5 / Fable 5.2
                                  ├─ More models
                                  └─ Effort: Low · Medium · High · Extra high · Max
```
Fable 5.2, Opus 5.5, Sonnet 5.5, Haiku 5.5 da thinking doim yoqiq; kamroq oʻylashi uchun effort'ni pasaytiring. Oʻzgarish keyingi javobdan boshlanadi.

**Projects** ([What are projects?](https://support.claude.com/en/articles/9517075-what-are-projects),
[Create and manage projects](https://support.claude.com/en/articles/9519177-how-can-i-create-and-manage-projects)):
```
claude.ai/projects → + New Project → nom + tavsif (Claude ularni KOʻRMAYDI)
                   → Set project instructions  (har bir chatga qoʻshiladi)
                   → + → fayllar               (Project bilim bazasi)
```
- Bilim bazasidagi fayllar Project'dagi hamma chatda ishlatiladi; chatlarning oʻzi esa bir-birini koʻrmaydi — maʼlumot
  boshqa chatga faqat bilim bazasi orqali oʻtadi. (Xotira — *memory* — yoqilgan boʻlsa, har Project'ning oʻz xotirasi bor;
  toʻliq ajratish uchun Project sahifasida "Use account memory" ni oʻchiring.)
- Pullik rejada bilim context window'ga yaqinlashsa, RAG yoqiladi: sigʻim 10× gacha, Claude faqat kerakli qismni qidiradi
  ([RAG for projects](https://support.claude.com/en/articles/11473015-retrieval-augmented-generation-rag-for-projects)).
- Project fayllari keshlanadi: qayta ishlatilganda limitdan kamroq sarflaydi
  ([Usage limit best practices](https://support.claude.com/en/articles/9797557-usage-limit-best-practices)).
- Koʻrsatmani qisqa tuting: umumiy kontekst, asosiy qoidalar, Claude'ning roli; vazifaga oid gapni chatning oʻziga yozing
  ([11647753](https://support.claude.com/en/articles/11647753-how-do-usage-and-length-limits-work)). Namuna: `loyiha-korsatma.md`.
- Diqqat: Anthropic Projects'ning yangi versiyasini (beta) bosqichma-bosqich chiqaryapti, avval Claude Code'da; ekranlar
  biroz farq qilishi mumkin. Videodagi sxemalar yordam markazidagi hozirgi versiya boʻyicha.

**Fayl chegaralari** ([Upload files to Claude](https://support.claude.com/en/articles/8241126-upload-files-to-claude)):
chat — 20 tagacha fayl, har biri 500 MB gacha; PDF ≤ 100 bet — matn + rasm/grafik, 101–1000 bet — faqat matn;
Project — fayl 30 MB gacha, faqat matn (multimodal PDF'dan tashqari). Roʻyxatdagi turlar: PDF, DOCX, CSV, TXT, HTML,
ODT, RTF, EPUB, JSON, XLSX, rasmlar — `.go` va `.md` bu roʻyxatda yoʻq. Ilova ularni qabul qilmasa, `.txt` qilib yuklang.

**Chatga nima yubormaslik.** Anthropic: moliyaviy maʼlumot, sogʻliq maʼlumoti, parol va login maʼlumotlari, maxfiy
hujjatlar bilan ehtiyot boʻling ([8325621](https://support.claude.com/en/articles/8325621-i-would-like-to-input-sensitive-data-into-my-chats-with-claude-who-can-view-my-conversations)).
API kalit — kredit karta raqami kabi: kim olsa, sizning hisobingizdan pul sarflaydi
([9767949](https://support.claude.com/en/articles/9767949-api-key-best-practices-keep-your-keys-safe-and-secure)).

**Maʼlumot qancha saqlanadi** (Free/Pro/Max, [How long do you store my data?](https://privacy.claude.com/en/articles/10023548-how-long-do-you-store-my-data),
[Is my data used for model training?](https://privacy.claude.com/en/articles/10023580-is-my-data-used-for-model-training)):
- Settings → Privacy → modelni yaxshilash sozlamasi. Ruxsat berilsa: anonimlashtirilgan holda 5 yilgacha.
- Oʻchirilgan suhbat: tarixdan darhol, serverdan 30 kun ichida.
- 👍/👎 bosilsa: butun suhbat 5 yilgacha.
- Incognito chat: tarixda yoʻq, oʻqitishga ishlatilmaydi, lekin 30 kun saqlanadi
  ([Use incognito chats](https://support.claude.com/en/articles/12260368-use-incognito-chats)). Project ichida incognito yoʻq.

**Pauza savoli.** Incognito chatga kalit yuborib, chatni yopdim — kalit xavfsizmi? Yoʻq: yuborilgan paytdayoq serverga ketdi, incognito
uni ham 30 kun saqlaydi;
yagona toʻgʻri yoʻl — kalitni bekor qilish (Console → API keys → `…` → Delete API key,
[8384961](https://support.claude.com/en/articles/8384961-what-should-i-do-if-i-suspect-my-api-key-has-been-compromised)).
Anthropic GitHub bilan hamkor: **ochiq GitHub repository**'da topilgan Claude kaliti avtomatik faolsizlantiriladi (*deactivates*)
([9767949](https://support.claude.com/en/articles/9767949-api-key-best-practices-keep-your-keys-safe-and-secure)) —
lekin hujjat bu himoyani faqat GitHub uchun aytadi.

## Uy vazifasi

1. **Oson.** Kursning bitta darsi uchun Project oching: dars `README.md`, kod va `go.mod` ni yuklang, `loyiha-korsatma.md`
   ni shu darsga moslang (masalan, «faqat shu darsdagi tushunchalar bilan tushuntir»).
2. **Oʻrta.** Bir xil savolni Haiku, Sonnet va Opus'ga bering, jadvalga yozing: javob toʻgʻrimi (`go run`/`go test`
   bilan tekshiring), qancha vaqt ketdi, qanchalik uzun.
3. **Qiyin.** Skanerga Telegram bot tokeni qoidasini qoʻshing va test yozing.
4. **Guru.** Skanerni `pre-commit` git hook qilib ulang: sir bor commit toʻxtasin.

### Javoblar
1. Namuna: `loyiha-korsatma.md` ning birinchi qatoriga «Hozir CL-1.01 darsini oʻrganyapman: token, context window,
   sampling. Boshqa mavzuga oʻtsam, eslat.» qoʻshing. Koʻrsatma qisqa qolsin — vazifaga oid gapni chatga yozing.
2. Jadval shabloni (natija sizniki boʻladi — biz Claude.ai'da oʻlchamadik):

   | Savol | Model | Effort | Toʻgʻrimi (`go test`) | Vaqt | Izoh |
   |---|---|---|---|---|---|
   | «`for range` ichida closure nega oxirgi qiymatni oladi?» | Haiku 5.5 | default | | | |
   | — | Sonnet 5.5 | default | | | |
   | — | Opus 5.5 | default | | | |

   Diqqat: Go 1.22 dan beri `for` sikli oʻzgaruvchisi har iteratsiyada yangi ([Go 1.22 release notes](https://go.dev/doc/go1.22#language)) —
   eskirgan javobni ushlash ham vazifaning bir qismi (CL-1.05).
3. `homework/telegram.go` + `homework/telegram_test.go` (`go test ./homework`). Format — Telegram rasmiy hujjatidagi
   namuna: `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11` ([Bot API: Authorizing your bot](https://core.telegram.org/bots/api#authorizing-your-bot)).
   `main.go` ga ulash: `qoidalar` ga `{"telegram-token", true, homework.TelegramToken, nil}` (yoki regex'ni koʻchiring).
4. `homework/pre-commit`: commit'ga ketayotgan fayllarning *index* dagi holatini vaqtinchalik papkaga chiqarib,
   `go run ./lessons/cl-02-claude-ai -qatiy` bilan tekshiradi. Oʻrnatish:
   `cp lessons/cl-02-claude-ai/homework/pre-commit .git/hooks/ && chmod +x .git/hooks/pre-commit`.
   Tekshirildi: parolli fayl bilan commit toʻxtadi, toza fayl bilan oʻtdi. Bu papkaning oʻzi soxta kalitlarga toʻla —
   shuning uchun hook uni ham ushlaydi (kerak boʻlsa `git commit --no-verify`, faqat yolgʻon signalda).
