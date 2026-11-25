# Chatga yuborishdan oldin: 30 soniyalik tekshiruv

Claude.ai, ChatGPT, forum, Telegram guruh — qayerga boʻlmasin, matn, log yoki faylni yuborishdan oldin.

## Hech qachon yubormang
- [ ] **API kalit va token**: `sk-ant-…`, `ghp_…`, `github_pat_…`, `AKIA…`, bulut va toʻlov kalitlari
- [ ] **Parol**: `.env`, `config.yaml`, `DB_URL=postgres://user:PAROL@…`, `password=…`
- [ ] **Maxfiy kalit fayli**: `-----BEGIN … PRIVATE KEY-----`, `id_ed25519`, `*.pem`
- [ ] **Boshqa odamning shaxsiy maʼlumoti**: mijozlar telefoni, emaili, manzili, pasport maʼlumoti
- [ ] **Moliyaviy maʼlumot**: karta raqami, bank hisobi
- [ ] **Sogʻliq maʼlumoti**, maxfiy ish hujjatlari (ishxona ruxsat bermagan boʻlsa)

Anthropic oʻzi ham shuni tavsiya qiladi: moliyaviy maʼlumot, sogʻliq maʼlumoti, parollar va login maʼlumotlari,
maxfiy hujjatlar bilan ehtiyot boʻling
([support.claude.com/…/8325621](https://support.claude.com/en/articles/8325621-i-would-like-to-input-sensitive-data-into-my-chats-with-claude-who-can-view-my-conversations)).

## Yuborishdan oldin
1. `go run . fayl.log` — bu papkadagi skaner. Topilma boʻlsa, `go run . -yashir fayl.log` bilan almashtiring.
2. Project'ga papka yuklasangiz: `go run . papka/` — `.env` fayli roʻyxatda chiqsa, uni yuklamang.
3. Skaner «topilma yoʻq» desa ham, koʻz bilan bir marta qarang: u faqat oʻz qoidalarini biladi.
4. Kalit oʻrniga oʻzgaruvchi nomini yozing: `ANTHROPIC_API_KEY=<yashirin>`.

## Kalit baribir ketib qolsa
1. **Darhol bekor qiling** (revoke). Claude API kaliti: Claude Console → API keys → `…` → *Delete API key*
   ([8384961](https://support.claude.com/en/articles/8384961-what-should-i-do-if-i-suspect-my-api-key-has-been-compromised)).
2. Yangi kalit yarating va uni faqat muhit oʻzgaruvchisida yoki secret manager'da saqlang.
3. Chatni oʻchirish kalitni qaytarmaydi: oʻchirilgan suhbat Anthropic serverlaridan 30 kun ichida oʻchadi
   ([10023548](https://privacy.claude.com/en/articles/10023548-how-long-do-you-store-my-data)), kalit esa shu paytgacha ham ishlaydi.
4. Git tarixiga tushgan boʻlsa: GIT-1.02 darsi ([git-02-env-leak](../git-02-env-leak/)).
