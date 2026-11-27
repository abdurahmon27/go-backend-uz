# GB-3.01 · Fayllar, `bufio` va `io.Reader`: katalogni fayldan oʻqiymiz

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #19"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** 2–3-modul darslari Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go boʻlsa,
> `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy hol;
> `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-2-10-utf8-rune`](../gb-2-10-utf8-rune) (UTF-8: bayt, rune va oʻzbekcha matn) · keyingi dars:
`gb-3-02` (JSON va struct tag'lar: mahsulotlarni `bozor.json`ga saqlash). `io.Writer`, `writeReport` va `go doc`:
[`gb-2-07-interface-writer`](../gb-2-07-interface-writer); `run() error` va `defer f.Close()`:
[`gb-2-09-defer-panic`](../gb-2-09-defer-panic); `parsePrice`: [`gb-1-07-funksiya-error`](../gb-1-07-funksiya-error).

**Bu darsdan keyin:** faylni `os.Open` bilan oʻqish uchun ochib, `defer f.Close()` bilan yopasiz; `go doc io.Reader`'dagi
`Read(p []byte) (n int, err error)` shartnomasini oʻqiysiz (avval `n` bayt, keyin `err`; oxirida `io.EOF`, `==` bilan); `bufio.Scanner` (`Scan`, `Text`, `Err`) bilan qatorma-qator
oʻqiysiz; bitta `loadCatalog(r io.Reader)` funksiyasini fayl, `os.Stdin` (`cat katalog.txt | go run .`) va testdagi
`strings.NewReader(…)` bilan ishlatasiz (fayl nomini emas, `io.Reader`'ni qabul qilib); `sc.Err()` tekshirilmasa «token too long»da maʼlumot jim yoʻqolishini
koʻrasiz.

Bu papka 2–3-modul qoidasi boʻyicha **mustaqil lab**: oʻz `go.mod`'i bor va GB-2.07/2.09'dagidek bitta
`package main` (GB-2.06/2.08'dagi `cmd/bozor` + `internal/…` tuzilishiga hamma lablarni 3-modul oxiridagi bozor CLI
v0, GB-3.06, birlashtiradi). `Product` (GB-2.03), `parsePrice` (GB-1.07/1.08), `writeReport` (GB-2.07) va
`run`/`main` (GB-2.09) oʻsha darslardan; yangi: `parseLine` (nomni GB-2.10'dagidek `utf8.ValidString` bilan ham
tekshiradi), `loadCatalog`, `katalog.txt`. Oldingi dars (GB-2.10) labida `Product` yoʻq, shuning uchun farqni GB-2.09
bilan koʻring: `diff -ru ../gb-2-09-defer-panic/solution solution`.

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | videoning oxiridagi kod: `parsePrice`, `parseLine`, `loadCatalog(r io.Reader)`, `writeReport`, argument boʻlsa fayldan, boʻlmasa `os.Stdin`'dan oʻqiydigan `run` + `main`. `loadCatalog` — vazifaning boshlanishi |
| `katalog.txt` | videodagi katalog: `ID;nom;narx (soʻm);qoldiq`, 4 qator, 63 bayt |
| `katta.txt` | videodagi buzilgan eksport: 3-qator 70 014 bayt (Scanner chegarasi 64 KiB dan uzun) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli va `katalog.txt`); `cd solution && go run . katalog.txt` |
| `_video/NN-…/main.go` | videodagi har bir qadam; `_` bilan boshlangan papkani `go test ./...` va `go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi vaqtinchalik `bozor/` papkasida (`go.mod`, `katalog.txt`, `katta.txt` bilan) ishga tushiradi va oxirida `(chiqish kodi: N)` yozadi (quyida bu qator tushirilgan) |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`)

Katalog (`cat` faylni terminalga chiqaradi):

```
$ cat katalog.txt
1;olma;12000;40
2;oʻrik;35000;5
3;non;4000;100
4;choy;25000;8
```

Qoʻshimcha (videoda yoʻq): nisbiy yoʻl jarayonning **joriy papkasidan** qidiriladi (`01-open`; `go build` binary yasaydi, keyin bir papka
yuqoridan ishga tushiramiz):

```
$ ./_video/run.sh 01-open sh -c 'go build && cd .. && ./bozor/bozor'
xato: open katalog.txt: no such file or directory
```

`io.Reader` va `*os.File`'ning `Read`'i (hujjat matni videoda qisqartirilgan):

```
$ go doc io.Reader
package io // import "io"

type Reader interface {
	Read(p []byte) (n int, err error)
}
    Reader is the interface that wraps the basic Read method.

    Read reads up to len(p) bytes into p. It returns the number of bytes read (0
    <= n <= len(p)) and any error encountered. …
    Callers should always process the n > 0 bytes returned before considering
    the error err. …
$ go doc os.File.Read
package os // import "os"

func (f *File) Read(b []byte) (n int, err error)
    Read reads up to len(b) bytes from the File and stores them in b. It
    returns the number of bytes read and any error encountered. At end of file,
    Read returns 0, io.EOF.
```

`f.Read`'ni 20 baytlik bufer bilan oʻzimiz chaqiramiz (`02-read`). Sikl `err == io.EOF` boʻlganda tugaydi: `Read`
`io.EOF`'ni oʻramay qaytaradi, hujjatga koʻra u `==` bilan tekshiriladi; boshqa xato `return err` bilan qaytadi. Read
qatorni ham, belgini ham bilmaydi: birinchi boʻlak `ʻ`'ning (`ca bb`) birinchi baytida tugaydi; oxirgi chaqiruv `0`
bayt va `io.EOF`:

```
$ ./_video/run.sh 02-read go run .
20 "1;olma;12000;40\n2;o\xca"
20 "\xbbrik;35000;5\n3;non;4"
20 "000;100\n4;choy;25000"
3 ";8\n"
0 ""
```

`bufio.Scanner` (`03-scan`): qatorlar `\n`'siz, `oʻ` butun, oxirida `sc.Err()` = `nil`:

```
$ ./_video/run.sh 03-scan go run .
"1;olma;12000;40"
"2;oʻrik;35000;5"
"3;non;4000;100"
"4;choy;25000;8"
```

`loadCatalog(f)` + `writeReport(os.Stdout, …)` (`04-load`), keyin fayl nomi argumentdan yoki stdin'dan (`main.go`):

```
$ ./_video/run.sh 04-load go run .
olma    12000 soʻm   40 dona
oʻrik   35000 soʻm    5 dona
non      4000 soʻm  100 dona
choy    25000 soʻm    8 dona
$ ./_video/run.sh final go run . katalog.txt
(xuddi shu 4 qator)
$ ./_video/run.sh final sh -c 'cat katalog.txt | go run .'
(xuddi shu 4 qator)
```

`loadCatalog`'ga fayl nomini (`string`) bersangiz, kompilyator rad etadi: `string`'da `Read` metodi yoʻq
(`08-string`; birinchi qator `# github.com/abdurahmon27/bozor` videoda qisqartirilgan):

```
$ ./_video/run.sh 08-string go run .
./main.go:112:31: cannot use "katalog.txt" (constant of type string) as io.Reader value in argument to loadCatalog: string does not implement io.Reader (missing method Read)
```

Argument ham, pipe ham boʻlmasa (`go run .`), stdin — terminalning oʻzi: dastur siz yozadigan qatorlarni kutadi,
`Ctrl+D` (qator boshida) kirish tugaganini bildiradi.

Tez-tez uchraydigan xato: `sc.Err()` tekshirilmasa (`06-errsiz`), 64 KiB dan uzun qatorda Scanner toʻxtaydi va
katalogning yarmi **jim** yoʻqoladi; `return products, sc.Err()` bilan xato koʻrinadi:

```
$ ./_video/run.sh 00-cat wc -L katta.txt
70014 katta.txt
$ ./_video/run.sh 06-errsiz go run . katta.txt
olma    12000 soʻm   40 dona
oʻrik   35000 soʻm    5 dona
$ ./_video/run.sh final go run . katta.txt
xato: bufio.Scanner: token too long
exit status 1
```

(`wc -L` — eng uzun qatorning uzunligi. Qator haqiqatan uzun boʻlishi mumkin boʻlsa, `sc.Buffer(buf, max)` bilan
chegarani oshirish mumkin: https://pkg.go.dev/bufio#Scanner.Buffer.)

### Qoʻshimcha: Scanner `Read`'ni necha marta chaqiradi (`07-count`, videoda yoʻq)

`logReader` — boshqa Reader'ni oʻrab, har bir `Read`'ni chiqaradigan oʻz Reader'imiz (`Read` metodi bor, demak u
ham `io.Reader`). Scanner `Read`'ga 4096 baytlik buferini beradi: 63 baytli fayl uchun ikki chaqiruv (bizning
20 baytli siklda beshta edi):

```
$ ./_video/run.sh 07-count go run . katalog.txt
Read: len(p)=4096 n=63 err=<nil>
Read: len(p)=4033 n=0 err=EOF
4 ta mahsulot
```

## Vazifa

[TASK.md](TASK.md): `loadCatalog` boʻsh qatorlarni va `#` izohlarni tashlab oʻtsin; xato `«N-qator: »` bilan
boshlansin, asl xato (`sc.Err()` ham) `%w` bilan ichida qolsin. Tekshirish: `go test` (hozir FAIL, vazifa
bajarilgach `ok`).

## Manbalar

- pkg.go.dev `io.Reader` — https://pkg.go.dev/io#Reader
- pkg.go.dev `os.Open`, `os.File.Read`, `os.Stdin` — https://pkg.go.dev/os#Open
- pkg.go.dev `bufio.Scanner`, `ScanLines`, `MaxScanTokenSize`, `ErrTooLong` — https://pkg.go.dev/bufio#Scanner
- pkg.go.dev `strings.Split`, `strings.NewReader` — https://pkg.go.dev/strings
- Go blog: Strings, bytes, runes and characters in Go — https://go.dev/blog/strings
- Go Code Review Comments: Interfaces — https://go.dev/wiki/CodeReviewComments#interfaces
- read(2) (file offset) — https://man7.org/linux/man-pages/man2/read.2.html
- Bash Reference Manual: Pipelines — https://www.gnu.org/software/bash/manual/html_node/Pipelines.html
