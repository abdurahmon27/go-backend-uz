# GO-X.04 · Nusxa yoki havola? Qiymat semantikasi va pointer

Kurs: **Go: JS/Python/PHP dasturchilari uchun** (GO-X), 4-dars — kursning oxirgi darsi · Lab:
https://haywan.uz/lab/youtube/go#go-x-04 · Oldingi dars: [GO-X.03 · try/catch yoʻq: error qiymat va defer](../0x-03-errors/)
· Keyingi qadam: Go asoslari kursi, [GO-1.01 · Go nega bor? `go build` ichida nima bor](../001-why-go/).

Kimlar uchun: boshqa tilda (JS/TS, Python, PHP) ishlaydigan, Go'ga yangi dasturchilar. Dasturlashni noldan
oʻrganayotgan boʻlsangiz, avval [Mohirdev](https://mohirdev.uz/kurslar/go-dasturlash-tili/) yoki
[dFarkhod](https://github.com/dFarkhod/golang-darslari) kurslarini koʻring.

Hamma natijalar shu kompyuterda olingan (2026-11-26): `go1.26.8 linux/amd64`, Node v24.10.0, Python 3.13.15,
PHP 8.4.26, golangci-lint 2.12.2. Buyruqlar shu papkadan ishga tushiriladi.

> `mapxato/` va `uyga/orta/` **ataylab** build boʻlmaydi, shuning uchun `go vet ./...` shu papkada xato beradi.
> Alohida tekshiring: `go vet ./hook ./tuzat ./nusxa ./param ./pointer ./range ./manzil ./maptuzat ./slice ./uyga/oson ./uyga/javob/...`

> `net/http` paketi C kompilyatori yoʻq mashinada `cgo: C compiler "gcc" not found` berishi mumkin
> (agar `go env CGO_ENABLED` 1 boʻlsa): unda buyruqlar oldidan `export CGO_ENABLED=0`.

**Bitta gʻoya:** Go'da `=` va funksiyaga uzatish qiymatni nusxalaydi (struct'ni ham, massivni ham), JS/Python'dagidek
oʻsha obyektga havola bermaydi; asl qiymatni oʻzgartirish uchun uni pointer (`&user`, `*User`) yoki indeks orqali yozasiz.

## JS/TS · Python · PHP · Go yonma-yon

| | JavaScript / TypeScript | Python | PHP | Go |
|---|---|---|---|---|
| `b = a` (obyekt / struct) | oʻsha obyektga ikkinchi havola | oʻsha obyektga ikkinchi nom | massiv → **nusxa**; obyekt → oʻsha obyektning handle'i | struct → **nusxa** (hamma maydonlar) |
| funksiyaga uzatish | obyekt — havola, funksiya uni oʻzgartira oladi | "passed by assignment" — oʻsha obyekt | massiv — nusxa, obyekt — handle | **nusxa**; oʻzgartirish uchun `*User` parametri va `&user` |
| sikl elementi | `for (const user of users)` — oʻsha obyekt | `for user in users` — oʻsha dict | `foreach ($users as $u)` — massivda **nusxa**, `&$u` bilan havola | `for _, user := range users` — **nusxa**; `users[i]` — elementning oʻzi |
| `b = a` (roʻyxat) | oʻsha massiv (`a[0]` → 9) | oʻsha list (9) | **nusxa** (1) | massiv `[3]int` → **nusxa** (1); slice `[]int` → elementlar **umumiy** (9) |
| havola koʻrinadimi | yoʻq, doim yashirin | yoʻq | `&` bilan ochiq | ha: tur `*User`, chaqiruv `&user` |

## Papkalar va natijalar

| Papka | Nima | Natija |
|---|---|---|
| `hook/` | GO-X.01 xizmati + xotiradagi `users`; `POST /users/{id}/topup` sikldagi **nusxani** oshiradi | POST → `{"id":"42","balance":150}`, keyingi GET → `{"id":"42","balance":100}` |
| `js/server.js` | Node'dagi xuddi shu xizmat (`for…of`) | POST → 150, GET → 150 |
| `anchor/topup.py` | Python: `for user in users` | `150` |
| `anchor/topup.php` | PHP: massiv `foreach`, keyin obyekt `foreach` | `100` · `150` |
| `anchor/topup-ref.php` | PHP tuzatish: `foreach ($users as &$u)` | `150` |
| `nusxa/` | `c := user` — butun struct nusxalanadi | `100 150` |
| `param/` | bashorat 1: `topUp(user User, …)` — parametr ham nusxa | `100` |
| `pointer/` | `topUp(user *User, …)` va `topUp(&user, 50)` | `150` |
| `range/` | bashorat 2: A `for _, user := range` · B `users[i]` · C `p := &user` | `A: 100` · `B: 150` · `C: 100` |
| `manzil/` | oʻrtadagi sinov: `p := &users[i]` — elementning oʻz manzili | `150` |
| `tuzat/` | hook tuzatilgan: `for i := range users { … users[i] … }` | POST → 150, GET → 150 |
| `mapxato/` | keng tarqalgan xato: `users["42"].Balance += 50` — ataylab build boʻlmaydi | `mapxato/main.go:14:2: cannot assign to struct field users["42"].Balance in map` |
| `maptuzat/` | map tuzatish: oʻqi → oʻzgartir → qaytarib yoz; `map[string]*User` (kalit yoʻq → nil → panic) | `150` · `150` |
| `slice/` | bashorat 3: massiv `b := a` va slice `t := s` | `1 9` |
| `uyga/` | uy vazifasi: `oson/`, `orta/` va `javob/` | |

## Videodagi buyruqlar

```bash
export CGO_ENABLED=0                                # gcc yoʻq mashinada
go run ./hook &                                     # listening on :3000
curl -X POST localhost:3000/users/42/topup          # {"id":"42","balance":150}
curl localhost:3000/users/42                        # {"id":"42","balance":100}   ← pul yoʻqoldi
kill %1                                             # (go run'ning bolasi qolsa: kill $(lsof -t -i :3000))
node js/server.js &                                 # Node: xuddi shu POST va GET → 150 va 150
go run ../0x-03-errors/defer-oldin                  # panic: runtime error: invalid memory address or nil pointer dereference
python3 anchor/topup.py                             # 150
php anchor/topup.php                                # 100 / 150
php anchor/topup-ref.php                            # 150   (&$u bilan)
go run ./nusxa                                      # 100 150
go run ./param                                      # 100
go run ./pointer                                    # 150
go run ./range                                      # A: 100 / B: 150 / C: 100
go run ./manzil                                     # 150
go run ./tuzat &                                    # keyin POST → 150, GET → 150
go run ./mapxato                                    # cannot assign to struct field users["42"].Balance in map
go run ./maptuzat                                   # 150 / 150
go run ./slice                                      # 1 9
node -e 'const s=[1,2,3]; const t=s; t[0]=9; console.log(s[0])'   # 9   (JS massivi = Go slice kabi)
php -r '$a=[1,2,3]; $b=$a; $b[0]=9; echo $a[0];'                    # 1
```

PHP shu kompyuterda oʻrnatilmagan, u `nix shell nixpkgs#php -c php …` bilan ishga tushirildi.

## Muhim nuqtalar (videodan qisqa)

- **Go'da hamma narsa qiymat bilan uzatiladi.** [Go FAQ](https://go.dev/doc/faq#pass_by_value): "everything in Go is
  passed by value. That is, a function always gets a copy of the thing being passed, as if there were an assignment
  statement assigning the value to the parameter." `c := user` ham, `topUp(user, 50)` ham struct'ning hamma
  maydonlarini nusxalaydi (sayoz nusxa: maydon slice yoki pointer boʻlsa, faqat sarlavha yoki manzil nusxalanadi). Build ham, `go vet` ham jim — bu xato emas, qoida.
- **Pointer** — qiymatning manzili. `&user` manzilni beradi, `*User` — "User'ga pointer" turi. Pointer orqali maydonga
  oddiy nuqta bilan yoziladi: `user.Balance += amount` ([spec: Selectors](https://go.dev/ref/spec#Selectors)).
  `nil` pointer orqali oʻqish — panic ([spec: Address operators](https://go.dev/ref/spec#Address_operators)), GO-X.03'dagi
  `defer resp.Body.Close()` xatosi aynan shu edi. Pointer'ni nusxalash maʼlumotni nusxalamaydi: FAQ: "passing a pointer
  value makes a copy of the pointer, but not the data it points to."
- **Qachon pointer?** Chaqiruvchining qiymatini oʻzgartirish kerak boʻlsa yoki struct katta boʻlsa.
  [Code Review Comments: Pass Values](https://go.dev/wiki/CodeReviewComments#pass-values): "Don't pass pointers as
  function arguments just to save a few bytes." Kichik struct'ni qiymat bilan uzatish odatiy. (Metod receiver'lari —
  [GO-1.12](../012-methods-receivers/).)
- **`range` nusxa beradi.** `for _, user := range users { user.Balance += 50 }` slice'ga tegmaydi;
  `for _, user := range users { p := &user; … }` ham — bu **nusxaning** manzili. Tuzatish: `users[i].Balance += 50` yoki
  `p := &users[i]` ([100 Go Mistakes #30](https://100go.co/#ignoring-that-elements-are-copied-in-range-loops-30)).
  Go 1.22 dan beri har iteratsiyaning `user`i alohida oʻzgaruvchi, lekin u baribir elementning **nusxasi**.
- **Linter buni har doim ushlamaydi.** golangci-lint (sukut boʻyicha) hook'da faqat `Encode` natijasi tekshirilmaganini
  koʻrsatadi. govet'ning `unusedwrite` tekshiruvini yoqsangiz, hech qachon oʻqilmaydigan yozuvni topadi (uy vazifasidagi
  `u.Active = false`), lekin hook'dagi `user.Balance += 50; Encode(user)`ni topmaydi — nusxa javobga yozilgan, yaʼni oʻqilgan:

  ```yaml
  # .golangci.yml
  version: "2"
  linters:
    settings:
      govet:
        enable:
          - unusedwrite
  ```
- **Map'dagi struct maydoniga yozib boʻlmaydi:** `users["42"].Balance += 50` → `cannot assign to struct field … in map`.
  Map elementi addressable emas ([spec: Address operators](https://go.dev/ref/spec#Address_operators)): map oʻsganda
  elementlarni yangi jadvalga koʻchiradi ([go1.26.8 runtime/maps `grow`](https://github.com/golang/go/blob/go1.26.8/src/internal/runtime/maps/table.go)),
  R. Griesemer: "we don't want a surviving pointer to a map element" ([golang/go#3117](https://github.com/golang/go/issues/3117)).
  Tuzatish: oʻqing → oʻzgartiring → `users[id] = user`, yoki `map[string]*User` (unda kalit yoʻq boʻlsa pointer `nil` → yozish panic).
- **Massiv va slice.** `[3]int` — massiv, uzunlik turning bir qismi; `b := a` hamma elementlarni nusxalaydi
  ([Go blog: "Go's arrays are values"](https://go.dev/blog/slices-intro)). `[]int` — slice: kichik sarlavha
  (elementlarga pointer, uzunlik, sigʻim); `t := s` faqat sarlavhani nusxalaydi, elementlar umumiy. Shuning uchun
  slice'ni funksiyaga bersangiz, funksiya elementlarni oʻzgartira oladi — imzoda pointer boʻlmasa ham (FAQ: "Map and slice
  values behave like pointers"). `append` bilan bogʻliq tuzoqlar (`len`/`cap`,
  aliasing) — [GO-1.08](../008-slices/).
- **Chuqurroq:** JS massivida elementlar — obyektlarga havolalar; Go'dagi `[]User` esa `User`larning oʻzini ketma-ket
  saqlaydi (range shu joydan nusxa oladi). Pointer olsangiz, qiymat heap'ga koʻchishi mumkin — escape analysis
  ([GO-1.11](../011-pointers-escape/)). Pointer arifmetikasi yoʻq: [FAQ](https://go.dev/doc/faq#no_pointer_arithmetic)
  "Safety."

## Uy vazifasi

- 🟢 **Oson** — `uyga/oson/main.go`: `deactivate` (range nusxasi) va `rename` (qiymat parametri) foydalanuvchini
  oʻzgartirishi kerak, lekin hech narsa oʻzgarmaydi (`go run ./uyga/oson` → `true` va `[{42 Dilnoza true} {7 Jasur true}]`).
  Birini indeks bilan, ikkinchisini pointer bilan tuzating. Kutilgan: `true` va `[{42 Dilnoza false} {7 Bekzod true}]`.
  Javob: `uyga/javob/oson/`.
- 🟡 **Oʻrta** — `uyga/orta/main.go`: `PATCH /users/{id}` map'dagi foydalanuvchi nomini oʻzgartirishi kerak, lekin fayl
  build boʻlmaydi (`cannot assign to struct field users[id].Name in map`). Ikki usulda tuzating: a) oʻqing → oʻzgartiring →
  qaytarib yozing; b) `map[string]*User`. Tekshiruv:
  `curl -X PATCH -d '{"name":"Dilnoza A."}' localhost:3000/users/42` → `{"id":"42","name":"Dilnoza A.","balance":100}`,
  keyingi GET ham shuni beradi, `…/users/1` → `404 page not found`. ⚠️ `net/http` har soʻrovni alohida goroutine'da
  bajaradi ("creating a new service goroutine for each", `go doc net/http.Server.Serve`), shuning uchun haqiqiy xizmatda
  map'ga yozish `sync.Mutex` bilan himoyalanadi ([GO-1.09](../009-maps/)). Javob: `uyga/javob/orta/` (a usul; b usul izohda).

## Manbalar

- Go FAQ: When are function parameters passed by value? https://go.dev/doc/faq#pass_by_value · Why is there no pointer
  arithmetic? https://go.dev/doc/faq#no_pointer_arithmetic · stack or heap https://go.dev/doc/faq#stack_or_heap
- Go spec: Address operators, Selectors, Assignment statements, Calls, For range, Array types, Slice types https://go.dev/ref/spec
- A. Gerrand, Go Slices: usage and internals https://go.dev/blog/slices-intro
- Go Code Review Comments: Pass Values https://go.dev/wiki/CodeReviewComments#pass-values
- golang/go#3117 (map element field assignment) https://github.com/golang/go/issues/3117
- 100 Go Mistakes #30 https://100go.co/#ignoring-that-elements-are-copied-in-range-loops-30
- MDN: Working with objects ("objects are a reference type") https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Working_with_objects#comparing_objects
- Python FAQ: call by reference https://docs.python.org/3/faq/programming.html#how-do-i-write-a-function-with-output-parameters-call-by-reference
- PHP: Arrays https://www.php.net/manual/en/language.types.array.php · Objects and references https://www.php.net/manual/en/language.oop5.references.php · foreach https://www.php.net/manual/en/control-structures.foreach.php

## Video matni

**hook.** Balansni ellikka oshiradigan endpoint yozdingiz. POST javobida balans yuz ellik, hammasi joyida. Lekin keyingi GET yana yuz deydi, xato ham, log ham yo'q. Node'dagi xuddi shu for of sikli esa yuz ellikni saqlaydi. Qisqa javob: Go'dagi range elementning nusxasini beradi. Nega, va qanday tuzatiladi?

**eslash.** Avval oldingi darsdan savol. defer resp.Body.Close xato tekshiruvidan oldin yozilgan, va http.Get xato qaytardi, nima bo'ladi?

**eslash-javob.** Javob: panic, ya'ni dastur runtime xatosi bilan yiqiladi: nil pointer dereference, nil pointer orqali o'qish. JavaScript'dagi Cannot read properties of null xatosiga o'xshaydi. resp pointer edi, bugun pointer nima va Go'da nega kerak, shuni ko'ramiz.

**reja.** Bugungi g'oya: Go'da teng belgisi va funksiyaga uzatish qiymatni nusxalaydi, struct'ni ham. JavaScript va Python'da esa o'sha obyektga ikkinchi havola paydo bo'ladi. Uchta tushuncha: qiymat nusxasi, pointer, va slice bilan massiv farqi.

**kod.** Mana handler'ning o'zi: User struct, users esa User'lar slice'i, ma'lumotlar bazasi o'rnida. Jingalak qavsdagi qiymat struct literal, JavaScript'dagi obyekt literali kabi. Sikl JavaScript'dagi for of bilan bir xil ko'rinadi: topilgan user'ga ellik qo'shiladi va javobga yoziladi.

**anchor.** Python ham shunday: for bergan user ro'yxatdagi o'sha dict, chunki nom obyektga bog'lanadi, natija yuz ellik. PHP'da esa massiv nusxalanadi, foreach'dagi o'zgarish yo'qoladi va yuz qoladi. Obyekt bo'lsa, PHP uni handle, ya'ni obyektga ishora orqali o'zgartiradi.

**nusxa-go.** Go'da struct qiymat: birinchi qatorda struct literal'dan user yaratiladi. c ikki nuqta teng user esa havola emas: user'ning hamma maydonlari yangi joyga nusxalanadi. Bu value semantics, ya'ni qiymat semantikasi. c'ni o'zgartirsak, user'ga tegmaydi.

**param.** Endi ishdagi holat: balansni alohida funksiya oshiradi. JavaScript'da obyektni funksiyaga bersangiz, funksiya aynan o'sha obyektni o'zgartiradi. Go'da ham User'ni topUp'ga beramiz.

**bashorat1.** topUp'dan keyin Println nima chiqaradi?

**javob1.** Javob: yuz, chunki funksiyaga uzatish ham qiymat berish: parametr chaqiruvchining user'idan nusxa, funksiya tugaganda nusxa yo'qoladi. Go'ning rasmiy qoidasi shu: hamma narsa qiymat bilan uzatiladi. Build ham, go vet ham jim, chunki bu xato emas, qoida.

**nega.** Nega Go shunday tanlagan? JavaScript'da istalgan funksiya obyektingizni jim o'zgartira oladi va buni imzodan bilib bo'lmaydi. Go'da struct'ingizni o'zgartirmoqchi bo'lgan funksiya buni imzosida yozadi: pointer bilan, uni hozir ko'ramiz. Slice va map esa alohida holat, ular oxirida.

**pointer.** Ampersand user, user'ning xotiradagi manzilini beradi. Yulduzcha User esa turning nomi: User'ga pointer. Funksiya ichida maydonga oddiy nuqta bilan yozamiz, Go manzil bo'yicha o'zi boradi, alohida belgi kerak emas.

**pointer-mem.** Xotirada: p ikki nuqta teng ampersand user yangi User yaratmaydi. p faqat user qayerdaligini saqlaydi. p orqali yozsak, user'ning o'z qiymati o'zgaradi. topUp'ga ampersand user berganda esa manzil nusxalanadi. Ikkala nusxa bitta joyni ko'rsatadi, shuning uchun qoida o'zgarmaydi: hamma narsa qiymat bilan uzatiladi.

**qachon.** Qachon pointer kerak bo'ladi? Chaqiruvchining qiymatini o'zgartirish kerak bo'lsa yoki struct katta bo'lsa. Go jamoasining tavsiyasi: bir necha bayt tejash uchun pointer uzatmang. Oldingi darslardagi ikki va'da ham shu: Unmarshal item'ni to'ldirishi uchun ampersand item oladi, handler esa so'rovni yulduzcha http.Request, ya'ni pointer bilan oladi.

**range-kod.** Endi boshidagi handler'ga qaytamiz, undagi sikl A variant. Yana ikki variantni ko'ring, har biridan oldin users qaytadan yuzga tiklanadi. B variantda indeks, C variantda pointer bor.

**bashorat2.** Qaysi variant users'dagi balansni yuz ellik qiladi?

**javob2.** Javob: faqat B, chunki range bergan user elementning nusxasi, xuddi c ikki nuqta teng user kabi. C'dagi pointer esa shu nusxaning manzili. B'da esa indeks orqali olingan element slice ichidagi elementning o'zi.

**sinov-kod.** O'rtadagi savol: bu safar indeks bilan yuramiz, lekin indeks orqali olingan elementning manzilini ampersand bilan olamiz. Keyin p orqali yozamiz.

**sinov.** Endi ekranda nima chiqadi?

**sinov-javob.** Javob: yuz ellik, chunki indeks orqali olingan element slice'dagi asl element, ampersand esa uning o'z manzilini beradi. C'da manzil nusxaniki edi. Bu shakl struct katta bo'lsa yoki elementni boshqa funksiyaga berish kerak bo'lsa qulay.

**tuzat.** Handler'ni tuzatamiz: indeks bilan yuramiz va elementga indeks orqali yozamiz, javobga ham shu elementni beramiz. Endi POST ham, GET ham yuz ellik. Standart linterlar bu xatoni ko'rmaydi, chunki nusxa javobga yozilgan, ya'ni ishlatilgan. Bu demo: haqiqiy xizmatda parallel so'rovlar uchun mutex, ya'ni qulf kerak, u Go asoslari kursining to'qqizinchi darsida.

**map-qoida.** Ko'p xizmat foydalanuvchilarni slice'da emas, id bo'yicha map'da saqlaydi, map JavaScript'dagi Map kabi. Bir qoida: Go faqat manzili bor narsaning maydoniga yozadi, ya'ni addressable. O'zgaruvchi va slice elementi shunday, map elementi esa yo'q: map o'sganda elementlarni yangi xotiraga ko'chiradi. Butun elementni almashtirish esa mumkin.

**map-kod.** JavaScript'da Map'dan elementni olib, balansini oshirish ishlaydi. Go'da xuddi shuni yozamiz.

**map-savol.** Bu dastur nima qiladi?

**map-javob.** Javob: build xatosi, cannot assign to struct field. Map elementining maydoniga yozib bo'lmaydi, va Go jim nusxani o'zgartirish o'rniga buni build paytida to'xtatadi.

**map-tuzat.** Tuzatish ikki xil: nusxani o'qing, o'zgartiring va map'ga qaytarib yozing. Yoki map'da pointer saqlang: ampersand User literal yangi User'ning manzilini beradi, map pointer qaytaradi va yozish build bo'ladi. Faqat kalit yo'q bo'lsa pointer nil, va yozish panic beradi.

**slice-kod.** Oxirgi tushuncha: kvadrat qavs ichida uzunlik bo'lsa, bu massiv, uzunlik turning bir qismi. Bo'sh qavs esa slice, oldingi darslarda ko'rgan ro'yxatimiz. Amalda massiv kam uchraydi, deyarli hamma joyda slice, lekin farq nusxa qoidasini yaxshi ko'rsatadi.

**bashorat3.** Bu safar Println nima chiqaradi?

**javob3.** Javob: bir va to'qqiz. Massiv qiymat: b ikki nuqta teng a hamma elementlarni nusxalaydi, xuddi struct kabi. JavaScript'da buning uchun uch nuqtali spread yozardingiz. Slice esa JavaScript massivi kabi elementlarni bo'lishadi, chunki uning ichida elementlarga pointer bor.

**slice-rasm.** Rasmda: slice bu kichik sarlavha, inglizcha slice header, ichida elementlarga pointer, uzunlik va sig'im, ya'ni capacity. t ikki nuqta teng s faqat sarlavhani nusxalaydi, elementlar bitta qoladi. Diqqat: struct'da slice maydoni bo'lsa, struct nusxasi ham faqat sarlavhani oladi. append tuzoqlari Go asoslari kursining sakkizinchi darsida.

**chuqur.** Chuqurroq: JavaScript'dagi obyektlar massivi aslida havolalar ro'yxati. Go'da User slice'i User'larning o'zini ketma-ket saqlaydi, range nusxani aynan shu joydan oladi. Pointer funksiyadan tashqariga chiqib ketsa, qiymat heap'ga, ya'ni GC boshqaradigan xotiraga o'tishi mumkin: bu escape analysis, o'n birinchi darsda. Va unsafe paketisiz pointer arifmetikasi yo'q.

**takror.** Uchta savol bilan qaytaramiz. Birinchi: range bergan user'ni o'zgartirgan sikl nega users'ga tegmaydi, va ikki tuzatish qanday? Ikkinchi: topUp'ga ampersand user berilganda nima nusxalanadi va nega user o'zgaradi? Uchinchi, birinchi darsdan: go build qilingan binary'ni serverga ko'chirsangiz, u yerda Go kerakmi?

**takror-javob.** Javoblar: birinchisi user elementning nusxasi, tuzatish indeks orqali yozish yoki indeks orqali olingan elementning manzilini ampersand bilan olish. Ikkinchisi: manzil nusxalanadi, lekin u o'sha user'ni ko'rsatadi. Uchinchisi: kerak emas, binary ichida runtime ham bor, faqat bitta OS va bitta protsessor uchun yig'iladi.

**uyga.** Uyga ikkita vazifa bor. Oson: deactivate va rename foydalanuvchini o'zgartirishi kerak, lekin hech narsa o'zgarmaydi. Birini indeks bilan, ikkinchisini pointer bilan tuzating. O'rta: map'dagi foydalanuvchi nomini o'zgartiradigan PATCH handler build bo'lmaydi, uni ikki usulda tuzating. Javoblar GitHub'dagi dars papkasida.

**davom.** Bu bilan kurs tugadi: bitta binary, runtime'dagi turlar, error qiymat va endi nusxa bilan pointer. Keyingi qadam Go asoslari kursi, birinchi darsi: Go nega bor va go build ichida nima bor. Har darsning yozma versiyasi va kodi tavsifdagi lab sahifasida.

**outro.** Bugun ko'rdik: Go'da qiymat nusxalanadi, o'zgartirish kerak bo'lsa esa pointer'ni ochiq yozasiz. Go asoslari kursida ko'rishguncha.
