# GO-X.03 · try/catch yoʻq: error qiymat va defer

Kurs: **Go: JS/Python/PHP dasturchilari uchun** (GO-X), 3-dars · Lab: https://haywan.uz/lab/youtube/go#go-x-03 ·
Oldingi dars: [GO-X.02 · Turlar runtime'da ham bor](../0x-02-types/) · Keyingi dars: GO-X.04 · Nusxa yoki havola?
Qiymat semantikasi va pointer.

Kimlar uchun: boshqa tilda (JS/TS, Python, PHP) ishlaydigan, Go'ga yangi dasturchilar. Dasturlashni noldan
oʻrganayotgan boʻlsangiz, avval [Mohirdev](https://mohirdev.uz/kurslar/go-dasturlash-tili/) yoki
[dFarkhod](https://github.com/dFarkhod/golang-darslari) kurslarini koʻring.

Hamma natijalar shu kompyuterda olingan (2026-11-26): `go1.26.8 linux/amd64`, Node v24.10.0, Python 3.13.15,
PHP 8.4.26, golangci-lint 2.12.2. Buyruqlar shu papkadan ishga tushiriladi (`config.json` shu yerda, unda
**ataylab** ortiqcha vergul bor).

> `net/http` paketi C kompilyatori yoʻq mashinada `cgo: C compiler "gcc" not found` berishi mumkin
> (agar `go env CGO_ENABLED` 1 boʻlsa): unda buyruqlar oldidan `export CGO_ENABLED=0`.

**Bitta gʻoya:** Go'da xato qaytariladigan oddiy qiymat, uni har chaqiruvdan keyin oʻzingiz tekshirasiz (tozalash ham oʻzingiz qoʻygan `defer` bilan), aks holda dastur notoʻgʻri qiymat bilan jim davom etadi.

## JS/TS · Python · PHP · Go yonma-yon

| | JavaScript / TypeScript | Python | PHP | Go |
|---|---|---|---|---|
| buzuq JSON | `JSON.parse` → `throw SyntaxError` | `json.load` → `raise JSONDecodeError` | `json_decode` → **jim `null`** (`JSON_THROW_ON_ERROR` bilan exception) | `json.Unmarshal` → **`error` qiymat qaytaradi** |
| xatoni kim koʻradi | eng yaqin `catch`; yoʻq boʻlsa jarayon yiqiladi | eng yaqin `except` | `json_last_error()`ni oʻzingiz soʻraysiz | chaqiruvchi: `if err != nil { … }`; tekshirmasa — **hech kim** |
| imzoda koʻrinadimi | yoʻq (TS'da `throws` yoʻq) | yoʻq | yoʻq | ha: `func ReadFile(name string) ([]byte, error)` |
| tozalash | `try { … } finally { fh.close() }` | `with open(…) as f:` | `finally` | `defer f.Close()` (funksiya tugaganda, har qanday `return`da) |
| HTTP 404 | `fetch` — xato emas (`res.ok`) | `requests` — xato emas (`raise_for_status()`) | — | `http.Get` — xato emas (`resp.StatusCode`) |
| stack trace | bor | bor | bor | standart `error`da **yoʻq**: faqat matn → kontekstni oʻzingiz qoʻshasiz (GO-1.14) |

## Papkalar va natijalar

| Papka | Nima | Natija |
|---|---|---|
| `config.json` | hook: oxirgi maydondan keyin ortiqcha vergul | — |
| `js/server.js` | Node'dagi xuddi shu xizmat | `SyntaxError: Expected double-quoted property name in JSON at position 58 (line 4 column 1)` … `at Object.<anonymous> (…/js/server.js:5:18)`, exit 1 |
| `hook/` | GO-X.01 xizmati + `json.Unmarshal(data, &cfg)` **natijasiz** | `listening on :0` — tasodifiy portda; `go build` va `go vet` jim |
| `eslash/` | takror (GO-X.02): `a, err := …; b, err := …` | `42 0 strconv.Atoi: parsing "x": invalid syntax` |
| `oqish/` | error qiymat: yoʻq fayl | `0 open yoq.json: no such file or directory` |
| `ikki/` | bashorat 1: `data := os.ReadFile(…)` — ataylab build boʻlmaydi | `assignment mismatch: 1 variable but os.ReadFile returns 2 values` |
| `config/` | `loadConfig`: xatoni yuqoriga qaytarish, `main`da `log.Fatal` | `invalid character '}' looking for beginning of object key string`, exit 1 |
| `tuzat/` | hook tuzatilgan: `if err := json.Unmarshal(…); err != nil` | `invalid character '}' …`, exit 1 (xizmat ishga tushmaydi) |
| `qamrov/` | oʻrtadagi sinov: `if cfg, err := loadConfig(…); …` keyin `cfg.Port` — ataylab build boʻlmaydi | `declared and not used: cfg` · `undefined: cfg` |
| `js/finally.mjs` | defer'ning JS tomoni: `try…finally` | `[ { id: 42, name: 'Dilnoza' } ]` |
| `http/` | `fetchUser`: `http.Get` + `defer resp.Body.Close()` (users API demo uchun `httptest`da) | `{42 Dilnoza} <nil>` · `Get "http://localhost:9/users/42": dial tcp [::1]:9: connect: connection refused` |
| `defer/` | bashorat 3: `defer` erta `return`dan **keyin** turibdi | `ochildi / o'qildi / yopildi / ochildi / xato: buzuq JSON` — `load(true)`da `yopildi` yoʻq |
| `defer-oldin/` | xato: `defer resp.Body.Close()` tekshiruvdan **oldin** (videoda faqat aytiladi) | `panic: runtime error: invalid memory address or nil pointer dereference` · `go vet`: `using resp before checking for errors` |
| `uyga/` | uy vazifasi: `oson/`, `orta/` va `javob/` | |

## Videodagi buyruqlar

```bash
export CGO_ENABLED=0                     # gcc yoʻq mashinada
cat config.json                          # "users_api": "http://localhost:9090",   ← ortiqcha vergul
node js/server.js                        # SyntaxError: Expected double-quoted property name in JSON at position 58 …
go run ./hook                            # listening on :0        (Ctrl+C bilan toʻxtating)
python3 -c 'import json; json.load(open("config.json"))'   # JSONDecodeError: Illegal trailing comma before end of object …
php -r 'var_dump(json_decode(file_get_contents("config.json")));'   # NULL
go run ./oqish                           # 0 open yoq.json: no such file or directory
go doc os.ReadFile | grep func           # func ReadFile(name string) ([]byte, error)
go run ikki/main.go                      # assignment mismatch: 1 variable but os.ReadFile returns 2 values
go build -o /dev/null ./hook && go vet ./hook; echo "exit $?"   # exit 0 — ikkalasi ham jim
ss -ltnp | grep hook                     # *:34329 … ("hook",…)  — port tasodifiy (go run ./hook ishlab turganda)
go run ./tuzat                           # invalid character '}' looking for beginning of object key string / exit status 1
golangci-lint run ./hook                 # hook/main.go:23:16: Error return value of `json.Unmarshal` is not checked (errcheck)
go run qamrov/main.go                    # declared and not used: cfg / undefined: cfg
node js/finally.mjs                      # [ { id: 42, name: 'Dilnoza' } ]
go run ./defer                           # ochildi o'qildi yopildi ochildi xato: buzuq JSON
go run ./http                            # {42 Dilnoza} <nil> / … connection refused
go run ./defer-oldin                     # panic: … nil pointer dereference   (defer tekshiruvdan oldin)
go vet ./defer-oldin                     # defer-oldin/main.go:12:8: using resp before checking for errors
```

PHP shu kompyuterda oʻrnatilmagan, u `nix shell nixpkgs#php -c php -r '…'` bilan ishga tushirildi.

## Muhim nuqtalar (videodan qisqa)

- **error — oddiy qiymat.** Xato beradigan funksiya uni oxirgi natija qilib qaytaradi; `nil` = xato yoʻq
  ([Effective Go: Errors](https://go.dev/doc/effective_go#errors)). Hech narsa "sakramaydi": keyingi qator baribir
  ishlaydi. Xato boʻlsa, boshqa natijalarga ishonmang (`data` boʻsh, `resp` nil boʻlishi mumkin).
  Nega exception emas? [Go FAQ](https://go.dev/doc/faq#exceptions): "coupling exceptions to a control structure, as in
  the try-catch-finally idiom, results in convoluted code. It also tends to encourage programmers to label too many
  ordinary errors, such as failing to open a file, as exceptional."
- **Natijani olsangiz — hammasini.** `data := os.ReadFile(…)` build boʻlmaydi. `_` bilan tashlash mumkin, lekin
  [Code Review Comments](https://go.dev/wiki/CodeReviewComments#handle-errors): "Do not discard errors using `_` variables."
- **Lekin chaqiruvni natijasiz yozsangiz — hammasi jim tashlanadi.** `json.Unmarshal(data, &cfg)` qatori build boʻladi,
  `go vet` ham jim ([spec: Expression statements](https://go.dev/ref/spec#Expression_statements)). Natija: `cfg.Port`
  zero value'da (0), `ListenAndServe(":0")` esa tasodifiy portni oladi
  ([net.Listen](https://pkg.go.dev/net#Listen): "If the port … is empty or "0" … a port number is automatically chosen").
  Buni [golangci-lint](https://golangci-lint.run/)'dagi **errcheck** ushlaydi (sukut boʻyicha yoqiq) — CI'ga qoʻshing.
  Bilib qoʻying: errcheck `defer resp.Body.Close()` va `w.Write(…)`ni ham koʻrsatadi (oʻqish uchun ochilgan body'da
  Close xatosini eʼtiborsiz qoldirish keng tarqalgan amaliyot), `_` bilan tashlanganlarni esa sukut boʻyicha koʻrmaydi
  (`check-blank: false`).
- **Xato matni:** kichik harf bilan, oxirida nuqtasiz ([Code Review: Error Strings](https://go.dev/wiki/CodeReviewComments#error-strings)) — chunki u boshqa matnga qoʻshiladi: `open yoq.json: no such file or directory`.
- **`if err != nil`.** Shart atrofida qavs yoʻq, `{ }` majburiy ([spec](https://go.dev/ref/spec#If_statements)). Xatoni
  hal qiling yoki yuqoriga qaytaring: ishga tushishda `log.Fatal`, handler'da `http.Error(w, …, 400/500)` + `return`,
  boshqa joyda `return …, err`. Asosiy yoʻl chap chetda, xato yoʻli `if` ichida
  ([Indent Error Flow](https://go.dev/wiki/CodeReviewComments#indent-error-flow)).
- **Qisqa shakl** `if err := f(); err != nil { … }`: `err` faqat shu `if` ichida yashaydi
  ([spec: Blocks](https://go.dev/ref/spec#Blocks)). Natija kerak boʻlsa (`cfg`), uni alohida qatorda oling.
- **`defer`** chaqiruvni oʻrab turgan funksiya tugaguncha kechiktiradi — `return` qaysi boʻlmasin. Lekin `finally`dan farqi: chaqiruv faqat `defer` qatori **bajarilganda** roʻyxatga yoziladi, undan oldingi `return`da hech narsa ishlamaydi (`defer/`)
  ([spec](https://go.dev/ref/spec#Defer_statements)). Yopish ochishning yonida turadi
  ([Go blog](https://go.dev/blog/defer-panic-and-recover)). `defer` xato tekshiruvidan **keyin**: xato boʻlsa `resp` nil,
  `defer resp.Body.Close()` esa `resp.Body`ni darhol oʻqiydi → panic (`defer-oldin/`; `go vet` buni ushlaydi). `resp.Body`ni oxirigacha oʻqib yopish chaqiruvchining vazifasi, aks holda keep-alive ulanish qayta ishlatilmasligi
  mumkin ([net/http.Response](https://pkg.go.dev/net/http#Response)). DB: `tx, err := db.BeginTx(…)` → `defer tx.Rollback()` (Commit'dan keyin u `sql.ErrTxDone` qaytaradi va hech narsa qilmaydi)
  ([Go docs](https://go.dev/doc/database/execute-transactions)).
- **Nega `if err != nil` qoladi?** Soʻrovnomalarda yillar davomida birinchi shikoyat; 2018 `check/handle`, 2019 `try`,
  2024 `?` takliflari koʻrildi va 2025-yilda Go jamoasi "will stop pursuing syntactic language changes for error
  handling" ([Go blog, 2025](https://go.dev/blog/error-syntax)).
- **Keyinroq (Go asoslari kursi):** xatoga kontekst qoʻshish `fmt.Errorf("…: %w", err)`, `errors.Is/As`
  ([GO-1.14](../014-errors-values/)); `defer` tartibi, sikldagi `defer`, `panic`/`recover`
  ([GO-1.15](../015-defer-panic-recover/)).

## Uy vazifasi

- 🟢 **Oson** — `uyga/oson/main.go`: config yuklovchi xatoni uch joyda tashlaydi (`data, _ :=`, natijasiz `Unmarshal`,
  `cfg, _ :=`). Avval taxmin qiling (`go run ./uyga/oson` → `port: 0`), keyin uchalasini `if err != nil` bilan
  tuzating. Tekshiruv: buzuq config → `invalid character '}' …`, `CONFIG=yoq.json` → `open yoq.json: no such file or
  directory`, `CONFIG=uyga/oson/config.json` → `port: 8080`. Javob: `uyga/javob/oson/`.
- 🟡 **Oʻrta** — `uyga/orta/main.go`: users API `404` qaytarsa, `fetchUser` xatosiz boʻsh foydalanuvchi beradi
  (`{0 } <nil>`): Go'da ham, `fetch`dagi kabi, non-2xx javob xato emas ([net/http.Get](https://pkg.go.dev/net/http#Get):
  "A non-2xx response doesn't cause an error."). `resp.StatusCode`ni tekshiring va
  `fmt.Errorf("users API: %s", resp.Status)` qaytaring. Kutilgan: `{0 } users API: 404 Not Found`.
  ⭐ GO-X.01 dagi `GET /users/{id}` handler'ini `fetchUser` bilan yozing: xato → `http.Error(w, err.Error(), http.StatusBadGateway)`.
  Javob: `uyga/javob/orta/`.

## Manbalar

- Go FAQ, Why does Go not have exceptions? https://go.dev/doc/faq#exceptions
- Go spec: If statements, Blocks, Expression statements, Defer statements https://go.dev/ref/spec
- Go Code Review Comments: Handle Errors, Indent Error Flow, In-Band Errors https://go.dev/wiki/CodeReviewComments
- Effective Go: Errors, Defer https://go.dev/doc/effective_go#errors
- A. Gerrand, Defer, Panic, and Recover (2010) https://go.dev/blog/defer-panic-and-recover
- R. Griesemer, [ On | No ] syntactic support for error handling (2025) https://go.dev/blog/error-syntax
- net/http: Response.Body, Get, Client.Do https://pkg.go.dev/net/http · net.Listen https://pkg.go.dev/net#Listen
- Executing transactions (`defer tx.Rollback()`) https://go.dev/doc/database/execute-transactions
- golangci-lint linters (errcheck) https://golangci-lint.run/docs/linters/
- PHP json_decode https://www.php.net/manual/en/function.json-decode.php · JSON_THROW_ON_ERROR https://www.php.net/manual/en/json.constants.php
- MDN: try...catch https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/try...catch · fetch() https://developer.mozilla.org/en-US/docs/Web/API/Window/fetch

## Video matni

**hook.** config.json'da bitta ortiqcha vergul qoldi. Node'dagi xizmat ishga tushmaydi va SyntaxError bilan yiqiladi. Go'dagi xuddi shu xizmat esa listening deb yozdi, lekin sakkiz ming sakson portida hech kim tinglamaydi. Xato qayerga yo'qoldi?

**eslash.** Avval oldingi darsdan savol. Birinchi qatorda a va err ikki nuqta teng bilan e'lon qilindi. Ikkinchi qatorda b va err yana ikki nuqta teng bilan. Bu build bo'ladimi?

**eslash-javob.** Javob: build bo'ladi, chunki chap tomonda b yangi, err esa faqat yangi qiymat oladi. Atoi oldingi uy vazifasidan, u stringni songa aylantiradi. Go kodida err aynan shunday qayta-qayta ishlatiladi, bugungi dars shu err haqida.

**reja.** Bugungi g'oya: Go'da xato throw qilinmaydi, u funksiya qaytaradigan oddiy qiymat. Uni tekshirmasangiz, hech narsa sakramaydi va dastur jim davom etadi. Uchta tushuncha: error qiymat throw o'rnida, if err nil'ga teng emas try catch o'rnida, va defer finally o'rnida.

**err-anchor.** Node'dagi throw'ni ko'rdik, Python'ning json.load'i ham raise qiladi. PHP'dagi json_decode esa sukut bo'yicha throw qilmaydi, jimgina null qaytaradi. Go'ning xatti-harakati shu ikkisining o'rtasida, hozir ko'ramiz.

**err-go.** Go'da xato qiymat. os.ReadFile ikki natija qaytaradi: baytlar slice'i va error. error bu matn beradigan Error metodi bor tur. Fayl yo'q bo'lsa, data nil, err esa xato matnini olib keladi, muvaffaqiyatda esa err nil. Keyingi qator baribir ishlaydi, JavaScript'dagidek catch'ga sakrash yo'q.

**bashorat1.** Endi bitta natija olamiz: data ikki nuqta teng os.ReadFile, err'siz. JavaScript'da massivdan birinchi elementni olish mumkin, Python esa tuple beradi. Go nima qiladi?

**javob1.** Javob: build xatosi, assignment mismatch. Natijani o'zgaruvchiga olsangiz, Go hamma natijani olishni talab qiladi. Pastki chiziq bilan ataylab tashlash mumkin, lekin Go jamoasining code review qo'llanmasi buni qilmang deydi.

**imzo.** go doc terminalda paket hujjatini ko'rsatadi. Imzoda ikki natija bor, error doim oxirida. JavaScript'da funksiya throw qilishini imzodan bilib bo'lmaydi, TypeScript'da ham throws yo'q. Go FAQ sababini aytadi: exception'lar fayl ochilmasligi kabi oddiy xatolarni favqulodda deb belgilashga undaydi.

**iferr.** Endi config'ni yuklaydigan funksiya. Parametrda ham avval nom, keyin tur, qavsdagi Config va error esa uning ikki natijasi. Keyin if err nil'ga teng emas: shart atrofida qavs yo'q, jingalak qavs esa majburiy. return cfg, err xatoni yuqoriga uzatadi, xuddi catch yozilmagan throw kabi, faqat ko'rinib turadi. Unmarshal faqat error qaytaradi, err allaqachon bor, shuning uchun oddiy teng.

**qayerda.** Xatoni qayerda hal qilasiz? Code review qo'llanmasi: hal qiling yoki yuqoriga qaytaring. Ishga tushishda log.Fatal bilan to'xtaysiz, handler'da mijozga to'rt yuz yoki besh yuz qaytarasiz, qolgan joyda yuqoriga qaytarasiz. Asosiy yo'l chap chetda qoladi, xato yo'li if ichida.

**xato-kod.** Endi boshidagi xizmatga qaytamiz. ReadFile tekshirilgan. Unmarshal esa natijasi olinmagan oddiy chaqiruv. Port shu cfg'dan olinadi, Sprintf uni ikki nuqtali manzilga aylantiradi. Bir ishora: fmt.Println ham ikki natija qaytaradi, n va err, lekin biz uni doim natijasiz chaqiramiz, va na build, na go vet bunga e'tiroz bildirgan.

**bashorat2.** Bu kodda go build va go vet nima deydi?

**javob2.** Javob: ikkalasi ham jim, chiqish kodi nol. Qisqasi, qoida ikki xil: natijani o'zgaruvchiga olsangiz, hammasini olasiz, chaqiruvni natijasiz yozsangiz, hamma natija jim tashlanadi. Unmarshal vergulga qoqilib xato qaytardi, uni hech kim o'qimadi, cfg.Port esa zero value'da, nolda qoldi.

**port0.** Nol port esa operatsion tizimga istalgan bo'sh portni tanla degani. ss buyrug'i tinglayotgan portlarni ko'rsatadi: xizmat tasodifiy portda, health check esa sakkiz ming sakson portiga boradi va javob olmaydi. Bu PHP'dagi jim null'ning Go'dagi ko'rinishi.

**tuzat.** Tuzatish: if ichida, nuqtali vergulgacha, qisqa e'lon yozamiz. Bu err faqat shu if ichida yashaydi. U tashqaridagi err'ni yashiradi, ya'ni shadowing, lekin bu yerda zararsiz.

**tuzat-run.** Endi xizmat buzuq config bilan ishga tushmaydi va sababini yozadi. golangci-lint Go'ning ESLint'i kabi linterlar to'plami, undagi errcheck tekshirilmagan error'larni topadi. Uni CI'ga qo'shing. Pastki chiziq bilan tashlanganini esa u sukut bo'yicha ko'rmaydi.

**sinov-kod.** O'rtadagi savol. Kimdir qisqa shaklni loadConfig bilan yozdi, natija ham if ichida e'lon qilindi.

**sinov.** Config to'g'ri bo'lsa, oxirgi qator nima chiqaradi?

**sinov-javob.** Javob: build xatosi, ikkita. if sarlavhasida e'lon qilingan cfg tashqarida yo'q, ichida esa ishlatilmagan. JavaScript'da blokdagi const'ni tashqarida ishlatsangiz, xato faqat ishlaganda chiqadi, Go'da esa build paytida. Natija keyin kerak bo'lsa, uni if'dan oldin alohida qatorda oling.

**defer-js.** Endi tozalash. JavaScript'da faylni yopish uchun try finally, Python'da with. Go'da finally bloki yo'q, uning o'rnida defer: bitta qator, u chaqiruvni funksiya tugaguncha kechiktiradi.

**defer-demo.** Fayl o'rniga ekranga yozamiz. errors.New matndan yangi error yasaydi. Diqqat: defer qatori erta return'dan keyin turibdi. main avval fail false, keyin fail true bilan chaqiradi va xatoni chiqaradi.

**bashorat3.** finally har doim ishlardi. load true chaqirilganda ekranda nima chiqadi?

**javob3.** Javob: ochildi va xato, yopildi yo'q. defer faqat o'sha qator bajarilganda ro'yxatga yoziladi, return undan oldin bo'ldi. fail false'da esa yopildi funksiya oxirida, o'qildi'dan keyin chiqdi. finally'dan farqi shu: u blokni emas, o'zidan keyingi yo'lni qamraydi. Qisqasi: defer funksiya oxirida ishlaydi, lekin faqat o'zi bajarilgan bo'lsa.

**defer-go.** Amalda: xizmatimiz foydalanuvchini boshqa xizmatdan HTTP orqali oladi. resp.Body'ni yopish chaqiruvchining vazifasi, aks holda keep-alive ulanish qayta ishlatilmasligi mumkin. defer xato tekshiruvidan keyin turadi: xato bo'lsa resp nil, defer qatori esa resp.Body'ni darhol o'qiydi va dastur yiqiladi. Bu tartibni go vet tekshiradi.

**defer-rasm.** Rasmda: defer qatori bajarilganda chaqiruv funksiyaning ro'yxatiga yoziladi, funksiya qaysi return bilan chiqmasin, avval ro'yxat bajariladi. Effective Go'dagi sabab: keyin yangi return qo'shsangiz ham, yopish unutilmaydi. DB'da ham shu qolip: defer tx.Rollback, ya'ni tranzaksiyani bekor qilish, commit'dan keyin u zararsiz.

**chuqur-nega.** Chuqurroq. if err nil'ga teng emas juda ko'p takrorlanadi, va Go so'rovnomalarida bu yillar davomida eng ko'p shikoyat. Go jamoasi yetti yil davomida check, try va so'roq belgisi kabi takliflarni ko'rib chiqdi. Ikki ming yigirma beshinchi yilda esa xato sintaksisini o'zgartirmaslikka qaror qildi.

**chuqur-stack.** Yana bir farq: standart error'da stack trace yo'q, u faqat qiymat va matn. log.Fatal invalid character dedi, lekin qaysi fayl ekanini aytmadi, Node esa fayl va qatorni ko'rsatgandi. Xatoga kontekst qo'shish va panic Go asoslari kursining o'n to'rtinchi va o'n beshinchi darslarida.

**takror.** Uchta savol bilan qaytaramiz. Birinchi: json.Unmarshal natijasiz yozildi, build bo'ladimi va xato qayerga ketadi? Ikkinchi: defer resp.Body.Close qachon ishlaydi va nega xato tekshiruvidan keyin turadi? Uchinchi, birinchi darsdan: ListenAndServe nega log.Fatal ichida?

**takror-javob.** Javoblar. Birinchisi: build ham, vet ham jim, xato tashlanadi va cfg zero value'da qoladi. Ikkinchisi: funksiya tugaganda, qaysi return bo'lmasin, xato bo'lsa esa resp nil va defer yiqiladi, go vet buni ushlaydi. Uchinchisi: ListenAndServe faqat xato bilan qaytadi, masalan port band bo'lsa, log.Fatal uni yozadi va dasturni to'xtatadi.

**uyga.** Uyga vazifa. Oson: config yuklovchi xatoni bir necha joyda tashlab yuboradi, natijasi port nol. Hammasini if err nil'ga teng emas bilan tuzating. O'rta: users API to'rt yuz to'rt qaytarsa, fetchUser xatosiz bo'sh foydalanuvchi beradi, chunki Go'da ham, fetch'dagi kabi, to'rt yuz to'rt xato emas. resp.StatusCode'ni tekshiring. Javoblar havoladagi papkada.

**outro.** Bugun ko'rdik: Go'da xato qaytariladigan qiymat, uni if err nil'ga teng emas bilan tekshirasiz, tozalashni esa defer bilan yozasiz. Keyingi darsda nusxa yoki havola: qiymat semantikasi va pointer. Ko'rishguncha.
