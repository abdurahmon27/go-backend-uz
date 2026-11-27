# GB-3.02 · JSON va struct tag'lar: mahsulotlarni `bozor.json`'ga saqlash

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #20"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2–3 darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang. Darsning «Chuqurroq» qismidagi
> `encoding/json/v2` Go 1.27'dan beri bor; Go 1.27'da `encoding/json` ham uning ustida ishlaydi (xato matnlari
> eski versiyalardagidan biroz farq qilishi mumkin, masalan `Go struct field .0.price`).

Oldingi dars: `gb-3-01` (fayllar, `bufio` va `io.Reader`) · keyingi dars: `gb-3-03` (jamoa asboblari: `go vet`,
`staticcheck`, `govulncheck`). Bu darsga kerak boʻlganlar: [`gb-2-03-struct-turlar`](../gb-2-03-struct-turlar)
(`Product`, `Order`, `OrderItem`), [`gb-2-04-pointer`](../gb-2-04-pointer) (`&x`: funksiya chaqiruvchining
oʻzgaruvchisiga faqat manzil orqali yozadi), [`gb-2-06-paketlar`](../gb-2-06-paketlar) (katta harf = eksport),
[`gb-2-07-interface-writer`](../gb-2-07-interface-writer) (`io.Writer`, `os.Stdout`),
[`gb-2-08-xatolar`](../gb-2-08-xatolar) (`errors.Is`, `errors.As`, `%w`), [`gb-2-09-defer-panic`](../gb-2-09-defer-panic)
(`defer f.Close()`).

**Bu darsdan keyin:** `encoding/json` bilan struct'ni JSON'ga (`json.Marshal`, `json.MarshalIndent`,
`json.NewEncoder`) va JSON'dan struct'ga (`json.Unmarshal(data, &v)`, `json.NewDecoder`) aylantirasiz;
`` `json:"name"` `` tag'lari bilan kalit nomlarini berasiz va `omitempty`'ni qayerga qoʻyish kerakligini bilasiz;
kichik harfli maydon nega jim yoʻqolishini, `Unmarshal`'ga `&` berilmasa nima boʻlishini (va uni `go vet`
topishini), buzuq va yoʻq fayl xatolarini koʻrasiz; katalogni `bozor.json`'ga saqlab, qayta yuklaysiz.

Bu papka 3-modulning boshqa darslari kabi **mustaqil**: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`,
`go 1.27`), oʻz `package main`'i, tashqi modul yoʻq. `Product` (`ID`, `Name`, `Price` tiyinda, `Stock`) va
`Order`/`OrderItem`/`OrderStatus` GB-2.03'dan; yangi maydon `Note` (izoh, `omitempty`). 3-modul oxiridagi bozor
CLI v0 (GB-3.06) modul darslarini birlashtiradi.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `Product` tag'lar bilan, `saveProducts`, `loadProducts`, `run` + `main` (GB-3.01 shakli: `bozor.json`'ni oʻqib chiqaradi, xatoni `os.Stderr`'ga); vazifa uchun `Order`, `OrderItem` (tag'siz) |
| `bozor.json` | videoda `saveProducts` yozgan fayl (`go run .` uni oʻqiydi) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/NN-…/` | videodagi har bir qadam (`main.go`, kerak boʻlsa `bozor.json`); `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 09-buzuq go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

```
$ go doc encoding/json.Marshal          (birinchi 4 qator)
package json // import "encoding/json"

func Marshal(v any) ([]byte, error)
    Marshal returns the JSON encoding of v.
```

`json.Marshal(p)`, tag'siz (`01-marshal`): kalitlar — maydon nomlarining oʻzi:

```
$ go run .
{"ID":1,"Name":"olma","Price":1200000,"Stock":40}
```

Tez-tez uchraydigan xato (`02-kichik`): maydonlar kichik harf bilan. `json` paketi ularni koʻrmaydi; xato ham, `go
vet` ogohlantirishi ham yoʻq:

```
$ go run .
{}
$ go vet .
$
```

Tag'lar bilan (`03-tag`); tag'da `json: "name"` (probel bilan, `03b-probel`) yozilsa, `json` uni tanimaydi va jim
`"Name"` yozadi, `go vet` esa koʻrsatadi:

```
$ go run .                               (03-tag)
{"id":1,"name":"olma","price":1200000,"stock":40}
$ go vet .                               (03b-probel)
main.go:12:2: struct field tag `json: "name"` not compatible with reflect.StructTag.Get: bad syntax for struct tag value
```

`Note` + `omitempty`, `json.NewEncoder(os.Stdout)` (`04-omit`); `Stock`'ga ham `omitempty` qoʻyilsa, qoldigʻi 0
boʻlgan mahsulotda `stock` kaliti yoʻqoladi (`05-quiz`):

```
$ go run .                               (04-omit)
{"id":1,"name":"olma","price":1200000,"stock":40}
{"id":2,"name":"oʻrik","price":3500000,"stock":5,"note":"quritilgan"}
$ go run .                               (05-quiz)
{"id":1,"name":"olma","price":1200000,"stock":40}
{"id":2,"name":"oʻrik","price":3500000,"note":"quritilgan"}
```

`saveProducts` (`06-save`): `json.MarshalIndent(products, "", "  ")` + `os.WriteFile(path, data, 0o644)`
(fayl huquqi `-rw-r--r--`, 189 bayt). Natija — papkadagi [`bozor.json`](bozor.json):

```
$ go run .
saqlandi: 2 ta mahsulot
```

`loadProducts` (`07-load`): `os.ReadFile` + `json.Unmarshal(data, &products)`. `&`'siz (`08-pointer`): kompilyator
jim, xato dastur ishlaganda; `go vet` esa oldindan topadi:

```
$ go run .                               (07-load)
1 olma: 40 dona
2 oʻrik: 5 dona
$ go run .                               (08-pointer)
xato: json: Unmarshal(non-pointer []main.Product)
exit status 1
$ go vet .                               (08-pointer)
main.go:34:26: call of Unmarshal passes non-pointer as second argument
```

Buzuq fayl (`09-buzuq`: olmaning `"price"` qiymati qoʻshtirnoqda) va yoʻq fayl (`10-yoq` = papkadagi `main.go`,
`bozor.json`'siz):

```
$ go run .                               (09-buzuq)
xato: json: cannot unmarshal string into Go struct field .0.price of type int64
exit status 1
$ rm bozor.json && go run .              (10-yoq: ./_video/run.sh final sh -c 'rm bozor.json && go run .')
xato: open bozor.json: no such file or directory
exit status 1
```

`json.NewDecoder(f).Decode(&products)` bilan `loadProducts` (`11-decoder`, `os.Open` + `defer f.Close()`): natija
`07-load` bilan bir xil.

Chuqurroq (`12-v2`): bitta kalit ikki marta kelganda `encoding/json` oxirgisini jim oladi, `encoding/json/v2` rad
etadi:

```
$ go run .
v1: <nil> non
v2: jsontext: duplicate object member name "name"
```

Hammasini qayta ishga tushirish:
`for s in 01-marshal 02-kichik 03-tag 03b-probel 04-omit 05-quiz 06-save 07-load 08-pointer 09-buzuq 10-yoq 11-decoder 12-v2 final; do ./_video/run.sh $s sh -c 'go vet . 2>&1; go run .'; done`

## Vazifa

[TASK.md](TASK.md): `Order` va `OrderItem`'ga jamoa formatidagi json tag'lari; `loadProducts` fayl yoʻq boʻlsa
boʻsh katalog qaytarsin, buzuq faylning xatosi fayl yoʻli bilan boshlansin va asl JSON xatosini `%w` bilan
saqlasin. Tekshirish: `go test .` (hozir FAIL, vazifa bajarilgach `ok`). Javobni tekshirish: `go test ./solution/`.

## Manbalar

- pkg.go.dev `encoding/json` (`Marshal`, `MarshalIndent`, `Unmarshal`, `NewEncoder`, `NewDecoder`, struct tag'lar,
  `omitempty`, "Migrating to v2") — https://pkg.go.dev/encoding/json
- pkg.go.dev `encoding/json/v2` — https://pkg.go.dev/encoding/json/v2 · Go 1.27 release notes — https://go.dev/doc/go1.27
- Go blog: JSON and Go (A. Gerrand) — https://go.dev/blog/json
- Go spec: Struct types (tags) — https://go.dev/ref/spec#Struct_types · String literals (raw) — https://go.dev/ref/spec#String_literals ·
  Exported identifiers — https://go.dev/ref/spec#Exported_identifiers
- pkg.go.dev `os` (`WriteFile`, `ReadFile`, `Open`, `ErrNotExist`) — https://pkg.go.dev/os
- `go vet` analyzers: structtag — https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/structtag ·
  unmarshal — https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unmarshal
- RFC 8259: The JSON Data Interchange Format — https://www.rfc-editor.org/rfc/rfc8259
