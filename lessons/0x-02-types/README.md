# GO-X.02 · Turlar runtime'da ham bor: `"2" + 2` JS'da 22, Go'da kompilyatsiya xatosi

Kurs: **Go: JS/Python/PHP dasturchilari uchun** (GO-X), 2-dars · Lab: https://haywan.uz/lab/youtube/go#go-x-02 ·
Oldingi dars: [GO-X.01 · 20 qatorlik HTTP API](../0x-01-http-api/) · Keyingi dars: GO-X.03 · try/catch yoʻq: error qiymat va defer.

Kimlar uchun: boshqa tilda (JS/TS, Python, PHP) ishlaydigan, Go'ga yangi dasturchilar. Dasturlashni noldan
oʻrganayotgan boʻlsangiz, avval [Mohirdev](https://mohirdev.uz/kurslar/go-dasturlash-tili/) yoki
[dFarkhod](https://github.com/dFarkhod/golang-darslari) kurslarini koʻring.

Hamma natijalar shu kompyuterda olingan (2026-11-26): `go1.26.8 linux/amd64`, Node v24.10.0, Python 3.13.15,
TypeScript 5.9.3. Har bir Go dasturini https://go.dev/play da ham ishga tushirish mumkin.

**Bitta gʻoya:** Go'da har bir oʻzgaruvchining turi build paytida qatʼiy belgilanadi va binary ichida ham qoladi:
`"2" + 2` build boʻlmaydi, JSON'dagi `"2"` int maydonga tushmaydi, qiymat berilmagan oʻzgaruvchi esa `undefined`
emas, oʻz turining *zero value*'si.

## JS/TS · Python · Go yonma-yon

| | JavaScript / TypeScript | Python | Go |
|---|---|---|---|
| `"2" + 2` | `"22"` (`+` stringga aylantiradi) | `TypeError`, faqat qator ishlaganda | **build xatosi**: `mismatched types string and int` |
| turlar qachon tekshiriladi | `tsc` paytida, keyin **oʻchiriladi** | runtime'da (type hint'lar tekshirilmaydi) | build paytida **va** runtime'da (`%T`, `encoding/json`) |
| `2 + 1.5` (oʻzgaruvchilar) | `3.5` | `3.5` | build xatosi: `int` va `float64` → `float64(a) + b` |
| JSON `{"price":"2"}` → son maydon | `as Item` jim, natija `22`; tekshirish uchun zod kabi kutubxona | — | `0` + `json: cannot unmarshal string into Go struct field Item.price of type int` |
| qiymatsiz oʻzgaruvchi | `undefined` (`+ 1` → `NaN`) | `NameError` | zero value: `0`, `""`, `false`, `nil` |
| eʼlon / qiymat berish | `let x = 1` / `x = 2` | `x = 1` (eʼlon yoʻq) | `x := 1` / `x = 2` (`var x int` → 0) |
| roʻyxat ustidan | `for (const p of a)` → qiymat | `for p in a` → qiymat | `for p := range a` → **indeks**; qiymat: `for _, p := range a` |
| son → matn | `String(65)` → `"65"`, `String.fromCharCode(65)` → `"A"` | `str(65)` → `"65"` | `string(65)` → `"A"` (fromCharCode kabi!) · `strconv.Itoa(65)` → `"65"` |

PHP (oʻrnatilmagan, shuning uchun bu yerda ishga tushirilmadi; [PHP qoʻllanmasi](https://www.php.net/manual/en/language.operators.arithmetic.php)
boʻyicha): `"2" + 2` → `4` (raqamli string songa aylanadi), birlashtirish uchun alohida `.` operatori bor.

## Papkalar va natijalar

| Papka | Nima | Natija |
|---|---|---|
| `js/narx.ts`, `js/narx.js`, `js/narx.py` | hook: TS `price: number`, JSON'da `"2"` | `tsc --strict` jim → `22`; `narx.js`da tur yoʻq; Python `TypeError` |
| `hook/` | `price + count` (string + int) — ataylab build boʻlmaydi | `hook/main.go:8:14: invalid operation: price + count (mismatched types string and int)` |
| `aralash/` | bashorat 1: `a := 2; b := 1.5; a + b` — ataylab build boʻlmaydi | `mismatched types int and float64` |
| `aylantir/` | conversion `float64(a)`, `%T` | `int float64` / `3.5` |
| `json/` | bashorat 2: `{"price": "2"}` → `Price int` | `0 json: cannot unmarshal string into Go struct field Item.price of type int` |
| `nol/` | zero value | `0 "" false` |
| `yoq/` | JSON'da maydon yoʻq: `{}` | `0 <nil>` — xato yoʻq, jim nol |
| `deklar/` | `:=` ikki marta + ishlatilmagan oʻzgaruvchi — ataylab build boʻlmaydi | `no new variables on left side of :=` · `declared and not used: count` |
| `soya/` → `soya-tuzat/` | oʻrtadagi sinov: `if` ichida `discount := 10` (shadowing) → `=` | `chegirma: 10` / `150` → `140` |
| `range/` → `range-tuzat/` | bashorat 3: `for p := range prices` | `3` (indekslar 0+1+2) → `10` |
| `id/` → `id-tuzat/` | bashorat 4: `"user:" + string(id)` | `user:A` · `go vet` ushlaydi → `strconv.Itoa`: `user:65` |
| `chuqur/` | Chuqurroq: `2 + 1.5` (konstantalar; videoda faqat slayd) | `3.5` |
| `uyga/` | uy vazifasi: `oson/`, `orta/` va `javob/` | |

## Videodagi buyruqlar

```bash
npx -p typescript@5.9.3 tsc --strict js/narx.ts && node js/narx.js   # 22  (videoda: npx tsc --strict …)
cat js/narx.js                         # Item ham, number ham yoʻq
python3 js/narx.py                     # TypeError: can only concatenate str (not "int") to str
go run hook/main.go                    # hook/main.go:8:14: invalid operation: price + count (mismatched types string and int)
go run aralash/main.go                 # aralash/main.go:8:14: invalid operation: a + b (mismatched types int and float64)
go run json/main.go                    # 0 json: cannot unmarshal string into Go struct field Item.price of type int
node -e 'let total; console.log(total + 1)'   # NaN
python3 -c 'print(total + 1)'          # NameError: name 'total' is not defined
go run yoq/main.go                     # 0 <nil>
go run deklar/main.go                  # no new variables on left side of := / declared and not used: count
go run soya/main.go                    # chegirma: 10 / 150
go run soya-tuzat/main.go              # chegirma: 10 / 140
go run range/main.go                   # 3
go run range-tuzat/main.go             # 10
go run id/main.go                      # user:A
go vet id/main.go                      # conversion from int to string yields a string of one rune, not a string of digits
go run id-tuzat/main.go                # user:65
go run chuqur/main.go                  # 3.5
```

Node 24 `.ts` faylni toʻgʻridan-toʻgʻri ham ishga tushiradi (`node js/narx.ts` → `22`): u turlarni shunchaki
oʻchiradi, tekshirmaydi ([Node.js: Type stripping](https://nodejs.org/api/typescript.html#type-stripping)).

## Muhim nuqtalar (videodan qisqa)

- **Statik tur.** Har bir oʻzgaruvchining turi build paytida belgilanadi va oʻzgarmaydi. `:=` turni oʻng tomondan oladi
  (TS `let` kabi). Operatorning ikki tomoni bir xil turda boʻlishi shart ([spec: Operators](https://go.dev/ref/spec#Operators)).
  Aralash turlar faqat ochiq *conversion* bilan: `float64(a) + b`. Sabab ([Go FAQ](https://go.dev/doc/faq#conversions)):
  "The convenience of automatic conversion between numeric types in C is outweighed by the confusion it causes."
- **Build xatosi → binary yoʻq.** Eski binary esa oʻzgarmay qoladi (oʻlchandi): deploy skriptida `go build`ning exit
  code'ini tekshiring (`set -e`).
- **Runtime'da ham.** `encoding/json` maydon turini runtime'da biladi: `"2"` → `int` maydon = `UnmarshalTypeError`,
  maydon zero value'da qoladi ([Unmarshal](https://pkg.go.dev/encoding/json#Unmarshal)). TS'da buni zod kabi kutubxona qiladi.
  Ehtiyot boʻling: Unmarshal faqat **turni** tekshiradi, validatsiya qilmaydi (oʻlchandi, go1.26.8):
  `{"PRICE": 5}` → `5` (kalitlar katta-kichik harfga qaramay mos keladi), `{"price": null}` → `0`, xatosiz,
  `{"price": 2.0}` → xato (`cannot unmarshal number 2.0 … of type int`), notanish kalitlar jim tashlanadi
  (`Decoder.DisallowUnknownFields` bilan oʻzgaradi), yoʻq maydon esa jim zero value (`yoq/`; zod esa `Required` deydi).
  Tur maʼlumoti runtime'da faqat interface (`any`) va reflection orqali ishlatiladi (`%T`, `encoding/json`);
  asosiy tekshiruv build paytida.
- **Zero value.** `var x T` → `0`, `""`, `false`, `nil` ([spec](https://go.dev/ref/spec#The_zero_value)). Son/string/bool
  darhol ishlaydi; `nil` map'ga yozish esa panic: `assignment to entry in nil map`. "Narx 0" va
  "narx yuborilmagan" farqlanmaydi; yechimlar (pointer maydon, `omitempty`, alohida tekshiruv) GO-X.04 va keyin.
- **`:=` va `=`.** `:=` yangi oʻzgaruvchi, `=` mavjudiga qiymat. Ichki blokdagi `:=` tashqi oʻzgaruvchini *yashiradi*
  (shadowing). Chap tomonda bir nechta nom boʻlsa, kamida bittasi yangi boʻlishi kerak, qolganlari faqat qiymat oladi
  (`a, err := f(); b, err := g()` — ikkinchi `err` yangi emas). `go vet` shadowing'ni sukut boʻyicha ushlamaydi.
  Ishlatilmagan oʻzgaruvchi — build xatosi ([Go FAQ](https://go.dev/doc/faq#unused_variables_and_imports)).
  `:=` faqat funksiya ichida; paket darajasida `var`.
- **`for`/`range`.** `while` yoʻq: `for shart { }`. Bitta oʻzgaruvchili `range` slice'da **indeks** beradi
  ([spec](https://go.dev/ref/spec#For_range)). (JS `for…in` ham kalit beradi, lekin string: `"0"`, `"1"`, `"2"`; ularni `0` ga qoʻshsangiz `"0012"` chiqadi.)
- **`string(int)`.** "for historical reasons" Unicode belgi beradi: `string(65)` = `"A"`
  ([spec](https://go.dev/ref/spec#Conversions_to_and_from_a_string_type)). `go vet` ushlaydi, `go test` bu tekshiruvni
  oʻzi ishga tushiradi ([Go 1.15](https://go.dev/doc/go1.15#vet)). 1 114 111 dan katta har qanday id `"�"` boʻladi —
  hammasi bitta Redis kalitiga tushadi. Toʻgʻrisi: `strconv.Itoa` / `strconv.FormatInt`.

## Uy vazifasi

- 🟢 **Oson** — `uyga/oson/main.go`: savat summasi build boʻlmaydi. Uchta xatoni tuzating (kompilyator avval ikkitasini
  koʻrsatadi, birini tuzatsangiz uchinchisi chiqadi). Avval natijani taxmin qiling. Javob: `uyga/javob/oson/` → `jami: 40`.
- 🟡 **Oʻrta** — `uyga/orta/main.go`: GO-X.01 dagi `GET /users/{id}`. `r.PathValue("id")` string qaytaradi: uni `int`ga
  aylantiring, `{"id": 42, "next": 43}` qaytaring, `/users/abc` uchun `400 Bad Request`. Faylda `strconv.Atoi` qolipi bor
  (`if err != nil` — GO-X.03 mavzusi, hozircha nusxalab oling). Javob: `uyga/javob/orta/`.
  `net/http` paketi C kompilyatori yoʻq mashinada `cgo: C compiler "gcc" not found` berishi mumkin
  (agar `go env CGO_ENABLED` 1 boʻlsa): unda `CGO_ENABLED=0 go run uyga/orta/main.go`.
  ```
  curl -i localhost:3000/users/42    # HTTP/1.1 200 OK … {"id":42,"next":43}
  curl -i localhost:3000/users/abc   # HTTP/1.1 400 Bad Request … id son emas
  ```

## Manbalar

- Go spec (go1.26.8): [Operators](https://go.dev/ref/spec#Operators) · [Conversions](https://go.dev/ref/spec#Conversions) ·
  [The zero value](https://go.dev/ref/spec#The_zero_value) · [Short variable declarations](https://go.dev/ref/spec#Short_variable_declarations) ·
  [Declarations and scope](https://go.dev/ref/spec#Declarations_and_scope) · [For range](https://go.dev/ref/spec#For_range) ·
  [Constant expressions](https://go.dev/ref/spec#Constant_expressions)
- Go FAQ: [implicit numeric conversions](https://go.dev/doc/faq#conversions) · [unused variables](https://go.dev/doc/faq#unused_variables_and_imports)
- [encoding/json.Unmarshal](https://pkg.go.dev/encoding/json#Unmarshal) · [strconv](https://pkg.go.dev/strconv) · [fmt %T](https://pkg.go.dev/fmt) ·
  [Go 1.15 vet: string(x)](https://go.dev/doc/go1.15#vet) · [The Laws of Reflection](https://go.dev/blog/laws-of-reflection)
- [TypeScript Handbook: Erased Types](https://www.typescriptlang.org/docs/handbook/2/basic-types.html#erased-types) ·
  [Node.js: Type stripping](https://nodejs.org/api/typescript.html#type-stripping) ·
  [MDN: Addition (+)](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Addition) · [Zod](https://zod.dev/)

## Video matni

**hook.** TypeScript'da price maydonini number deb yozdingiz, backend esa narxni JSON'da string qilib yubordi. tsc strict rejimda ham jim. JavaScript esa string ikki plyus ikkini yigirma ikki qiladi. Go'da xuddi shu qoʻshish build boʻlmaydi. Nega shundayligini va bu backend uchun nimani oʻzgartirishini koʻramiz.

**eslash.** Avval oldingi darsdan savol. go build xato bilan tugadi, lekin papkada kechagi binary bor. Deploy skripti uni serverga koʻchirsa, nima boʻladi?

**eslash-javob.** Javob: eskisi oʻzgarmay qoladi. go build xato bersa, yangi fayl yozmaydi, va skript build xatosida toʻxtamasa, kechagi versiyani jimgina deploy qiladi. Bugun aynan shunday xatolarni koʻramiz: turlar mos kelmasa, binary umuman chiqmaydi.

**reja.** Bugungi gʻoya bitta: Go turni build paytida qat'iy tekshiradi va tur ma'lumotini runtime'gacha olib boradi, TypeScript esa uni oʻchiradi. Undan uchta amaliy farq chiqadi: statik tur, zero value, ikki nuqta teng va oddiy teng. Oxirida range'ning indeks tuzogʻi.

**tur-js.** Avval dalil. tsc chiqargan JavaScript faylda Item ham, number ham yoʻq, runtime ularni koʻrmaydi. Python'da shu kod TypeError beradi, lekin faqat shu qator ishlaganda, ya'ni balki production'da.

**tur-go.** Go'da har bir oʻzgaruvchining turi build paytida belgilanadi va oʻzgarmaydi, bu statik tur. Ikki nuqta teng turni oʻng tomondagi qiymatdan oladi, xuddi TypeScript'dagi let kabi: price string, count int. Go turli turdagi oʻzgaruvchilarni oʻzi aylantirmaydi, shuning uchun xato build paytida chiqdi va binary yaratilmadi.

**bashorat1.** Endi ikkala tomon ham son. a ikki, b bir butun oʻndan besh. Println nima chiqaradi?

**javob1.** Javob: build xatosi. JavaScript'dagi number aslida float64, butun son uchun alohida tur yoʻq. Go'da esa int va float64 ikki xil tur. JavaScript ham, Python ham uch butun oʻndan besh beradi, Go esa qoʻshmaydi.

**aylantir.** Tuzatish: turni ochiq aylantirasiz, float64, qavs ichida a. Bu conversion. Printf'dagi foiz T qiymatning turini chiqaradi, JavaScript'dagi typeof kabi, lekin aniqroq: number emas, int yoki float64. Go FAQ'da sababi yozilgan: avtomatik aylantirishning qulayligi u keltiradigan chalkashlikka arzimaydi. Qisqasi: aralash turlar faqat ochiq conversion bilan qoʻshiladi.

**json.** Endi runtime. Mana hook'dagi JSON, Go'da. Item turi TypeScript'dagi type Item kabi, Go'da bunday tur struct deyiladi, uni alohida darsda koʻramiz. Backtick ichidagi tag JSON'dagi price kaliti shu maydonga tushishini aytadi. Unmarshal JSON baytlarini item ichiga yozadi. Ampersand item'ning manzili, pointer'lar toʻrtinchi darsda.

**bashorat2.** err bu xato qiymati, keyingi darsning mavzusi. TypeScript'dagi Item turi runtime'da jim edi. Go'dagi Unmarshal nima qiladi?

**javob2.** Javob: xato qaytaradi. Unmarshal maydon turi int ekanini runtime'da biladi va stringni qabul qilmaydi. TypeScript'da bu ma'lumot build'dan keyin yoʻq, shuning uchun zod kabi validatsiya kutubxonasi kerak. Go'da tur ma'lumoti binary'da qoladi va encoding/json uni oʻqiydi. Nol esa qayerdan keldi?

**nol-js.** Qiymat berilmagan oʻzgaruvchi: JavaScript'da undefined va NaN, Python'da NameError.

**nol-go.** Go'da undefined yoʻq. var total int, ya'ni TypeScript'dagi let total: number, faqat avval nom, keyin tur. Qiymat berilmasa, oʻzgaruvchi oʻz turining nol qiymatini oladi, bu zero value deyiladi. int uchun nol, string uchun boʻsh string, bool uchun false. Printf'dagi foiz q stringni qoʻshtirnoq bilan koʻrsatadi.

**nol-nil.** Struct maydonlari ham shunday: var item Item yozsangiz, item.Price nol boʻladi. Keyingi darslardagi pointer, map va slice'ning zero value'si esa nil, ya'ni hech narsa. nil map'ga yozish dasturni yiqitadi, buni oʻsha darslarda koʻramiz.

**nol-savol.** JSON'da price umuman boʻlmasa-chi? zod bu holatda Required deb rad etadi. Go nima qiladi?

**nol-prod.** Javob: nol, xatosiz. err'ning zero value'si ham nil, ya'ni xato yoʻq. Unmarshal faqat kelgan kalitlarni yozadi, qolgan maydon zero value'da qoladi. Demak bepul mahsulot va unutilgan maydon bir xil koʻrinadi, toʻlov API'sida bu jiddiy bug. Hozircha muhim maydonni alohida tekshiring, toʻliq yechim pointer darsida. Qisqasi: notoʻgʻri tur xato beradi, yoʻq maydon esa jim zero value boʻladi.

**deklar.** Ikki nuqta teng TypeScript'dagi let kabi yangi oʻzgaruvchi e'lon qiladi. Oddiy teng mavjud oʻzgaruvchiga yangi qiymat beradi. Bir farqi bor: chap tomonda bir nechta nom boʻlsa, kamida bittasi yangi boʻlishi kerak, qolganlari qiymat oladi. Bu keyingi darsdagi err uchun muhim. Bu kodda ikkita build xatosi bor, toping.

**deklar-run.** Birinchisi: no new variables, chap tomonda yangi nom yoʻq. Ikkinchisi Go'ga xos: declared and not used. TypeScript'da bu linter ogohlantirishi, Go'da esa build xatosi. Go FAQ'ga koʻra, ishlatilmagan oʻzgaruvchi koʻpincha bug belgisi.

**soya.** Endi ikki nuqta teng bilan bogʻliq klassik tuzoq. total bir yuz ellik, discount nol. if ichida discount'ga oʻn beramiz, lekin ikki nuqta teng bilan.

**sinov.** Oxirgi Println nima chiqaradi?

**sinov-javob.** Javob: bir yuz ellik. if ichidagi ikki nuqta teng yangi discount yaratdi va u tashqaridagisini yashirdi, bu shadowing deyiladi. TypeScript'da yangi oʻzgaruvchini let soʻzidan koʻrasiz, Go'da esa farq bitta ikki nuqtada. Tuzatish: ichkarida oddiy teng yozing, shunda natija bir yuz qirq.

**range.** Endi for. Go'da while yoʻq: for'dan keyin faqat shart yozsangiz, u while boʻladi. Roʻyxat ustidan yurish uchun range, u har qadamda ikki qiymat beradi: indeks va element. Kvadrat qavsli int bu slice, int'lar roʻyxati, uni alohida darsda koʻramiz. Bu kod narxlarni qoʻshmoqchi.

**bashorat3.** Bu yerda oʻzgaruvchi bitta, p. Println nima chiqaradi?

**javob3.** Javob: uch. Bitta oʻzgaruvchi yozsangiz, u indeks boʻladi: nol, bir, ikki. TypeScript'dagi for of esa qiymat beradi, shuning uchun qoʻl oʻrganib qolgan.

**range-tuzat.** Tuzatish: indeks oʻrniga pastki chiziq, ya'ni blank identifier, u qiymatni tashlab yuboradi. U kerak, chunki ishlatilmagan indeks build xatosi boʻlardi. Endi natija oʻn.

**xato.** Oxirgisi, Go jamoasi alohida tekshiruv qoʻshgan xato. Redis kaliti yasaymiz: user, ikki nuqta va id. Eslatma: Unicode'da A harfining raqami oltmish besh. string, qavs ichida id: kalit nima boʻladi?

**xato-javob.** Javob: user, ikki nuqta, A. Go'dagi string, qavs ichida son, JavaScript'dagi String emas, String.fromCharCode kabi ishlaydi. Go spec'da yozilgan: bu imkoniyat tarixiy sabablar bilan qolgan.

**xato-vet.** Kompilyator jim, lekin go vet, ya'ni Go bilan keladigan statik tekshiruvchi, buni ushlaydi. Xabardagi rune bu bitta Unicode belgi. go test ham bu tekshiruvni oʻzi ishga tushiradi. Tuzatish: strconv paketidagi Itoa, u sonni oʻnlik raqamlarga aylantiradi. Teskarisi Atoi, uning qolipi uy vazifasi faylida bor.

**chuqur-runtime.** Chuqurroq. Tur ma'lumoti binary'da qanday yashaydi? Println istalgan turni, ya'ni any qabul qiladi. any qiymati ichida ikki narsa bor: qiymatning oʻzi va uning tur ma'lumoti. Foiz T shuni oʻqiydi. encoding/json esa reflect paketi bilan maydonlar va ularning turlarini runtime'da koʻradi. Tekshiruvlarning asosiysi esa baribir build paytida.

**chuqur-const.** Yana bir nuqta: ikki plyus bir butun oʻndan besh oʻzgaruvchisiz yozilsa, ishlaydi. Bular konstantalar, ularning turi hali tanlanmagan, shuning uchun Go ularni qoʻsha oladi. Konstantalar haqida alohida dars bor.

**takror.** Uchta savol bilan qaytaramiz. Birinchi: JSON'da price string boʻlib keldi, maydon esa int, Go nima qiladi? Ikkinchi: if ichida ikki nuqta teng yozsak, tashqaridagi oʻzgaruvchi oʻzgaradimi? Uchinchi: bitta oʻzgaruvchili range nima beradi?

**takror-javob.** Javoblar. Birinchisi: Unmarshal xato qaytaradi, maydon esa nol boʻlib qoladi. Ikkinchisi: yoʻq, ichkarida yangi oʻzgaruvchi paydo boʻladi, bu shadowing. Uchinchisi: indeks, qiymat kerak boʻlsa indeks oʻrniga pastki chiziq yozasiz.

**uyga.** Uyga vazifa. Oson: uyga oson papkasidagi savat summasi build boʻlmaydi, uchta xatoni tuzating va natijani taxmin qiling. Oʻrta: oldingi darsdagi users endpoint'ida id string. Uni Atoi bilan songa aylantiring, son boʻlmasa toʻrt yuzinchi status qaytaring, Atoi qolipi faylda tayyor. Javoblar havoladagi papkada.

**outro.** Bugun koʻrdik: Go turni build paytida tekshiradi va tur ma'lumotini runtime'gacha olib boradi, undefined oʻrniga esa zero value bor. Keyingi darsda try catch'siz yashaymiz: error qiymat va defer. Koʻrishguncha.
