# GB-1.08 · Birinchi test: `go test`, `t.Errorf` va jadvalli testlar

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #8"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

Oldingi dars: [`gb-1-07-funksiya-error`](../gb-1-07-funksiya-error) · keyingi dars: `gb-2-01` (slice: `len`, `cap`,
`append`).

**Bu darsdan keyin:** Go funksiyasiga test yozasiz, `go test` bilan ishga tushirasiz, yiqilgan test xabarini
oʻqiysiz va jadvalli test (table-driven test) bilan xatoni topib tuzatasiz. Vositalar: `_test.go` fayli,
`func TestXxx(t *testing.T)`, `t.Errorf` / `t.Fatalf`, `go test -v`, `-run`, `t.Run` bilan subtest.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi bozor: `parsePrice` endi `maxPriceSom` (10 mln soʻm) chegarasini koʻpaytirishdan **oldin** tekshiradi; `main`dagi eski `price > 10_000_000*tiyinPerSom` tekshiruvi olib tashlandi |
| `price_test.go` | videodagi jadvalli test `TestParsePrice` |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa: `parseSold` uchun `sold_test.go` jadvali va 1000 donalik chegara |
| `solution/` | javob: `main.go`, `sold_test.go` (+ xuddi shu ikki test fayli) |
| `_video/NN-…/` | videodagi har bir qadamning `main.go` va `price_test.go` fayllari; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi (bir nechtasida ataylab xato bor) |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 03-jadval go test` |

Oldingi darsdan farq (`main.go` GB-1.07'ning `solution/main.go`'sidan boshlangan; GB-1.07'ning `main_test.go`
tekshiruvi oʻrniga bu darsning test fayllari keldi):

```bash
diff -u ../gb-1-07-funksiya-error/solution/main.go main.go
```

## Videodagi buyruqlar va natijalari (go1.26.8 linux/amd64, `CGO_ENABLED=0`)

GB-1.07'dagi bozor juda katta narxni qabul qiladi (`./_video/run.sh 00-gb107 bash -c 'go build && ./bozor add olma 100000000000000000 1'`):

```
bozor ishga tushdi
versiya: 0.1.0
olma, 1 dona   -84467440737095516.-16
QQS            -912720851596686.-13
yetkazish          15000.00
jami           -85380161588677202.-29
```

Sababi: `strconv.Atoi` uchun 100 000 000 000 000 000 toʻgʻri son (bu kompyuterda `int` 8 bayt), lekin
`* tiyinPerSom` dan keyin natija 10¹⁹ boʻlishi kerak, `int64`ning eng kattasi esa 9 223 372 036 854 775 807.
Ortigʻi aylanib, manfiy son chiqadi (GB-1.03). `main`dagi `price > 10_000_000*tiyinPerSom` tekshiruvi manfiy sonni
oʻtkazib yuboradi, chunki u koʻpaytirishdan keyin ishlaydi.

Birinchi test, ikki holat (`_video/01-test`): `go test` → `PASS` / `ok …`; `go test -v`:

```
=== RUN   TestParsePrice
--- PASS: TestParsePrice (0.00s)
PASS
ok  	github.com/abdurahmon27/bozor	0.002s
```

Testni bir marta yiqilgan holda koʻrish: `parsePrice`dan `* tiyinPerSom` olib tashlangan (`_video/02-buzuq`).
`t.Errorf` testni toʻxtatmaydi, shuning uchun ikkala xabar chiqadi:

```
--- FAIL: TestParsePrice (0.00s)
    price_test.go:8: parsePrice("12000") = 12000, <nil>; kutilgan 1200000, xatosiz
    price_test.go:12: parsePrice("500") = 500, <nil>; kutilgan 50000, xatosiz
FAIL
exit status 1
FAIL	github.com/abdurahmon27/bozor	0.002s
```

`exit status 1` faqat `go test` (paket nomisiz, papkaning ichida) shaklida chiqadi; `go test .` yoki `go test ./...`
uni yozmaydi, lekin oxirida yana bitta `FAIL` qatori qoʻshadi.

Jadvalli test GB-1.07'ning kodida (`_video/03-jadval`): chegara qatori `ortiq` va `ulkan` yiqiladi:

```
--- FAIL: TestParsePrice (0.00s)
    --- FAIL: TestParsePrice/ortiq (0.00s)
        price_test.go:23: parsePrice("10000001") = 1000000100, <nil>; xato kutilgan: true
    --- FAIL: TestParsePrice/ulkan (0.00s)
        price_test.go:23: parsePrice("100000000000000000") = -8446744073709551616, <nil>; xato kutilgan: true
FAIL
exit status 1
FAIL	github.com/abdurahmon27/bozor	0.002s
```

Birinchi tekshiruvda `t.Fatalf` oʻrniga `t.Errorf` boʻlsa (`_video/06-errorf`), har bir yiqilgan subtest ikkinchi,
foydasiz xabarni ham chiqaradi (`go test -run TestParsePrice/ulkan`):

```
--- FAIL: TestParsePrice (0.00s)
    --- FAIL: TestParsePrice/ulkan (0.00s)
        price_test.go:23: parsePrice("100000000000000000") = -8446744073709551616, <nil>; xato kutilgan: true
        price_test.go:26: parsePrice("100000000000000000") = -8446744073709551616; kutilgan 0
FAIL
exit status 1
FAIL	github.com/abdurahmon27/bozor	0.002s
```

Tuzatilgandan keyin (`_video/04-tuzat`, `main.go` bilan bir xil): `go test` → `PASS` / `ok`, va

```
$ go build && ./bozor add olma 100000000000000000 1
bozor ishga tushdi
versiya: 0.1.0
xato: narx 10000000 soʻmdan oshmasin, keldi 100000000000000000
```

Bitta subtest: `go test -v -run TestParsePrice/ulkan`

```
=== RUN   TestParsePrice
=== RUN   TestParsePrice/ulkan
--- PASS: TestParsePrice (0.00s)
    --- PASS: TestParsePrice/ulkan (0.00s)
PASS
ok  	github.com/abdurahmon27/bozor	0.002s
```

`-run` qiymati regexp va nomning istalgan joyiga mos kelishi yetadi: `-run TestParsePrice/ol` `nol`ni ishga
tushiradi, `-run Price` esa `TestParsePrice`ni. Aniq bitta nom uchun regexp'ni ikki tomondan bogʻlang:
`go test -v -run '^TestParsePrice$/^nol$'`.

Tez-tez uchraydigan xato: kichik harfli `func testParsePrice` test emas (`_video/05-kichik`):

```
testing: warning: no tests to run
PASS
ok  	github.com/abdurahmon27/bozor	0.002s
```

Papkada boshqa testlar boʻlsa, ogohlantirish ham chiqmaydi: test shunchaki ishlamaydi. Qaysi testlar borligini
`go test -list .` koʻrsatadi. `TestparsePrice` (Test'dan keyin kichik harf) esa `go test` ichidagi vet
tekshiruvida toʻxtaydi: `TestparsePrice has malformed name: first letter after 'Test' must not be lowercase`.

Chuqurroq: `go test -c` test dasturini ishga tushirmasdan `bozor.test` faylini qoldiradi (`.gitignore` uni
commit'dan chetlatadi):

```
$ go test -c && ls -lh bozor.test
-rwxr-xr-x … 4.2M … bozor.test
$ ./bozor.test
PASS
```

`go test` (paket nomisiz, papkaning ichida) natijani keshlamaydi; `go test .` yoki `go test ./...` kabi paket
koʻrsatilsa, oʻzgarmagan paket uchun vaqt oʻrniga `(cached)` chiqishi mumkin (`-count=1` keshni oʻchiradi).
Jadval siklidagi `tc` har aylanishda yangi oʻzgaruvchi (`go.mod`dagi `go 1.22` va undan yangi versiyalarda),
shuning uchun eski maqolalardagi `tc := tc` qatori kerak emas.

## Vazifa

[TASK.md](TASK.md): `sold_test.go`da `TestParseSold` jadvalini yozing (`ortiq` = `"1001"` qatori yiqilsin), keyin
`parseSold`ga 1000 donalik chegara qoʻshing. `go test` oʻtguncha.

Papka tekshiruvi:

```bash
gofmt -l . && go vet ./...          # hech narsa chiqmaydi
go test . ; echo $?                 # FAIL (TASK.md'dagi xabar), 1
go test ./solution/...              # ok
```

## Manbalar

- `testing` paketi: https://pkg.go.dev/testing (TestXxx, `*_test.go`, `T.Errorf`, `T.Fatalf`, `T.Run`)
- `go help test` (`-c`, local directory mode), `go help testflag` (`-v`, `-run`, `-list`):
  https://pkg.go.dev/cmd/go#hdr-Test_packages , https://pkg.go.dev/cmd/go#hdr-Testing_flags
- Go wiki, TableDrivenTests: https://go.dev/wiki/TableDrivenTests
- go.dev blog, Using Subtests and Sub-benchmarks: https://go.dev/blog/subtests
- Go Code Review Comments, Useful Test Failures: https://go.dev/wiki/CodeReviewComments#useful-test-failures
- Go 1.22: per-iteration loop variables: https://go.dev/blog/loopvar-preview
- Go spec, Integer overflow: https://go.dev/ref/spec#Integer_overflow
- `strconv.Atoi`: https://pkg.go.dev/strconv#Atoi
