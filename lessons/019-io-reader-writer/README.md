# 19 · `io.Reader` va `io.Writer`: Go'ning eng chiroyli ikki qatori

Oldingi dars: GO-1.18 · Iterator'lar: `range` endi funksiya ustida ham ishlaydi · Keyingi dars: GO-1.20 · Fayllar,
`bufio` va JSON: v1 endi v2 ustida ishlaydi. Asos: GO-1.13 (interface = shartnoma) va GO-1.04 (matn = baytlar,
`oʻ` dagi `ʻ` = 2 bayt). Slice — GO-1.08.

**Bitta asosiy fikr:** katta maʼlumotni bir yoʻla emas, kichik boʻlaklab oʻqing va yozing. Hovuzdagi suvni butun
hovuzni koʻtarib emas, chelak bilan tashiysiz. Buni **oqim** (stream) deyishadi.

| Yangi soʻz | Oddiy tilda | Misol |
|---|---|---|
| oqim (stream) | maʼlumotni boʻlakma-boʻlak oʻqish/yozish | chelak bilan suv tashish |
| `io.Reader` (oʻquvchi) | "boʻsh chelak bering, toʻldiraman" — bitta metod: `Read(p []byte) (n int, err error)` | fayl, klaviatura, internet ulanishi |
| `io.Writer` (yozuvchi) | "toʻla chelak bering, toʻkaman" — bitta metod: `Write(p []byte) (n int, err error)` | fayl, ekran, sayt javobi |
| `io.EOF` | "maʼlumot tugadi" belgisi (end of file), nosozlik emas | — |
| `io.Copy` | nasos: Reader'dan oladi, Writer'ga beradi, oʻrtada bitta chelak | `io.Copy(&h, f)` |

> **Oʻlchov.** go1.26.8 linux/amd64, Linux 6.18, ext4, i5-12450H, 2026-11-26. "Xotira" = **dasturning oʻzi**
> egallagan eng koʻp operativ xotira (peak RSS): `/proc/self/status` dagi `VmHWM`, MiB (2²⁰ bayt), `ozim.RSS()`.
> Linux'ning fayl keshi (page cache) bunga kirmaydi — u ikkala dastur uchun ham faylni xotirada ushlab turishi mumkin.
> Videoda "megabayt" deyilgan. `getrusage` ning `Maxrss` i bu yerga yaramadi: `go run` dasturni vfork + exec bilan
> ishga tushiradi va Maxrss `go` buyrugʻining ~20 MiB xotirasini ham "meros" oladi.

| Papka | Nima | Videoda |
|---|---|---|
| `hammasi/` | `os.ReadFile` → butun fayl bitta slice'da → qatorlarni sanash | hook |
| `oqim/` | 32 KiB chelak, `f.Read` sikli, avval baytlar, keyin `io.EOF` | hook, chelak, savol |
| `bosh/` | hech narsa qilmaydigan Go dasturi — runtime'ning oʻzi qancha oladi | javob |
| `utf8/` | 4 baytlik chelak `ʻ` ni kesadi | — (qoʻshimcha) |
| `ozim/` | `Hisoblagich` (oʻz Writer), `Generator` (oʻz Reader), `RSS()` + testlar | oʻz Writer |
| `sanash/` | `io.Copy(&h, f)` — `oqim/` dagi sikl bitta qatorda | io.Copy |
| `siqish/` | `gzip.Writer` faylni oʻrab turadi: `io.Copy(gz, f)` | matryoshka |
| `yarat/` | 1 GiB test fayli: `io.Copy(f, ozim.NewGenerator(…))` | — (qoʻshimcha) |
| `nusxa/` | fayldan faylga `io.Copy` → Linux'da `copy_file_range` | — (qoʻshimcha) |
| `olch.sh` → `olch.txt` | diagramma: 16…1024 MiB faylda ikki usul | oʻlchov |
| `homework/` | uy vazifasi javoblari (testlar bilan) | vazifa |

Tekshirish (shu papkadan; bu mashinada C kompilyatori yoʻq, shuning uchun `CGO_ENABLED=0`):
```bash
CGO_ENABLED=0 go vet ./... && CGO_ENABLED=0 go test ./...
```

## Videodagi buyruqlar
Katta fayllar faqat `/tmp` da (gitda emas): 1 GiB fayl + 385 MiB `.gz` uchun ~1.5 GiB boʻsh joy kerak.
```bash
mkdir -p /tmp/io19
go run ./yarat 1024 /tmp/io19/buyurtma.csv     # yozildi: 1073741824 bayt · eng koʻp xotira: 2.7MiB
du -h --apparent-size /tmp/io19/buyurtma.csv   # 1.0G
go run ./hammasi /tmp/io19/buyurtma.csv        # buyurtmalar: 61003111 · eng koʻp xotira: 1028.1MiB
go run ./oqim /tmp/io19/buyurtma.csv           # buyurtmalar: 61003111 · eng koʻp xotira: 2.5MiB / Read chaqiruvi: 32769
go run ./bosh                                  # boʻsh dastur · eng koʻp xotira: 2.3–2.5MiB
sed -n 86,88p $(go env GOROOT)/src/io/io.go    # type Reader interface { Read(p []byte) (n int, err error) }
go run ./utf8                                  # "7000" true / " so\xca" false / "\xbbm" false
go run ./sanash /tmp/io19/buyurtma.csv         # buyurtmalar: 61003111 · bayt: 1073741824 · eng koʻp xotira: 2.5MiB
go run ./siqish /tmp/io19/buyurtma.csv /tmp/io19/buyurtma.csv.gz
#   oʻqildi: 1024.0MiB · yozildi: 384.6MiB · eng koʻp xotira: 3.1MiB   (~17 s, gzip.BestSpeed)
zcat /tmp/io19/buyurtma.csv.gz | cmp - /tmp/io19/buyurtma.csv && echo 'bir xil'   # bir xil
./olch.sh > olch.txt                           # diagramma (fayllarni oʻzi yaratadi va oʻchiradi)
rm -rf /tmp/io19                               # tozalash
```
VmHWM har ishga tushirishda ±0.3 MiB farq qiladi (oqim: 2.4–2.6, boʻsh dastur: 2.3–2.5).

`olch.txt` (MiB):

| fayl | butun fayl (`os.ReadFile`) | boʻlaklab (`oqim/`) |
|---:|---:|---:|
| 16 | 18.9 | 2.5 |
| 64 | 66.9 | 2.4 |
| 256 | 259.2 | 2.5 |
| 1024 | 1028.1 | 2.4 |

## Dars matni

**Hook.** 1 GiB faylda 61 003 111 ta buyurtma (har biri bitta qator `id,summa,tolandi`). Ikki dastur ularni sanaydi:
`hammasi` faylni xotiraga butunlay yuklaydi — 1028 MiB; `oqim` boʻlaklab oʻqiydi — qariyb 2.5 MiB.

**Xotira = ish stoli.** Dastur ishlatayotgan narsalar operativ xotirada (stolda) turadi, fayl esa diskda (shkafda).
10 GB faylni butunlay yuklash uchun 10 GB dan koʻp stol kerak; boʻlaklab oʻqish uchun baribir ~2.5 MiB.

**Gʻoya.** Hovuz va chelak: katta maʼlumotni kichik boʻlaklab oʻqish/yozish = oqim (stream). Go'da buni ikki kichik
shartnoma bajaradi.

**Reader.** GO-1.13 dagi kabi interface — shartnoma. `io/io.go` L86–L88:
```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```
`p` — boʻsh chelak (bayt slice, GO-1.08). Chelakni chaqiruvchi beradi, shuning uchun bitta chelak qayta-qayta
ishlatiladi. `n` — nechta bayt solindi, `err` — xato yoki "tugadi".

**Writer.** Teskarisi: toʻla chelak `p`, `n` — nechta bayt yozildi (io.go L90–L101). Ikkalasi ham faqat baytlarni
biladi. Reader: `*os.File`, `os.Stdin`, `net.Conn`; Writer: `*os.File`, `os.Stdout`, `http.ResponseWriter`.

**Chelak sikli** (`oqim/main.go`): `buf := make([]byte, 32*1024)`; `n, err := f.Read(buf)`; `bytes.Count(buf[:n], …)`.
Read chelakni toʻldirishi shart emas (io.go L58–L61: *"If some data is available but not len(p) bytes, Read
conventionally returns what is available instead of waiting for more"*) — doim `n` ga qarang; aniq `N` bayt kerak
boʻlsa `io.ReadFull` bor. Fayl tugasa `io.EOF` (io.go L37–L40). Tartib: avval `n` bayt, keyin `err` — *"Callers should
always process the n > 0 bytes returned before considering the error err"* (io.go L72–L75), chunki baʼzi Reader'lar
oxirgi baytlarni `EOF` bilan birga qaytaradi (`iotest.DataErrReader` aynan shunday qiladi).

**Savol.** 1 GiB ÷ 32 KiB = 1 048 576 ÷ 32 = 32 768 toʻla chelak + bitta `(0, io.EOF)` chaqiruvi = **32 769**.
(1 GiB fayl aynan 2³⁰ bayt chiqdi; `Generator` oxirgi qatorni butun qoldiradi, boshqa hajmlarda bir necha bayt ortiq.)
Xotira ≈ 2.5 MiB — `bosh/` (boʻsh dastur) ham 2.3–2.5 MiB: faylni oʻqish deyarli qoʻshimcha xotira olmadi.

**Oʻlchov.** `os.ReadFile` fayl hajmini `Stat` bilan oladi va bitta slice ochadi (`os/file.go` `readFileContents`):
xotira ≈ fayl. Boʻlaklab: hamma hajmda ≈ 2.5 MiB.

**Oʻz Writer'imiz** (`ozim.Hisoblagich`): baytlarni saqlamaydi, sanaydi; `return len(p), nil`. Writer qoidalari
(io.go L90–L98): *"Write must return a non-nil error if it returns n < len(p)"*; *"Implementations must not retain p"* —
`io.Copy` chelakni keyingi boʻlak uchun qayta ishlatadi.

**io.Copy — nasos** (`sanash/`): `io.Copy(&h, f)` — `oqim/` dagi sikl bitta qatorda; 61 003 111 qator, 1 073 741 824
bayt, ≈ 2.6 MiB. Fayl va Hisoblagich bir-birini tanimaydi, faqat shartnomani biladi.

**Matryoshka** (`siqish/`): `gzip.NewWriterLevel(out, gzip.BestSpeed)` — Writer ichida Writer. 1 GiB → 403 326 915
bayt (384.6 MiB, ≈ 37.6 %), ≈ 3.1 MiB xotira; `zcat … | cmp - asl_fayl` → bir xil.

### Qoʻshimcha (videoda yoʻq)
- **Chelak harfni kesadi (GO-1.04).** `ʻ` = `0xCA 0xBB`. 4 baytlik chelak bilan `"7000 soʻm"`: `" so\xca"`, `"\xbbm"` — ikkalasi ham
  `utf8.Valid` = false. Bizning sanoq toʻgʻri: `'\n'` = `0x0A`, koʻp baytli UTF-8 belgilarning hamma baytlari ≥ `0x80`
  (RFC 3629 §3), shuning uchun qator oxiri belgisi hech qachon kesilmaydi. Matnni qayta ishlash uchun — `bufio` (GO-1.20).
- **Oʻz Reader'imiz** (`ozim/generator.go`): buyurtmalarni oldindan saqlamaydi, `Read(p)` da `p` ga sigʻadigan
  qatorlarni yasaydi, qolgani `navbat` da kutadi, tugasa `0, io.EOF`. `len(p) == 0` da `0, nil` (io.go L77–L79) —
  `iotest.TestReader` bilan tekshirilgan (`ozim/ozim_test.go`). `yarat/`: 1 GiB fayl ≈ 2.7 MiB xotira bilan yoziladi.
- **Butun oʻqish qachon mayli:** kichik, hajmi maʼlum fayl (sozlamalar) — `os.ReadFile`. Hajmi nomaʼlum kirish
  (yuklangan fayl, HTTP body) — boʻlaklab; HTTP serverda `http.MaxBytesReader` (chegaradan oshsa xato qaytaradi;
  `io.LimitReader` esa jim EOF beradi — `LimitReader(r, n+1)` bilan oʻqib, `n` dan uzunini rad eting).
- **`io.ReadAll` ≈ 2× fayl:** go1.26 da u oʻsib boruvchi boʻlaklar yigʻadi va oxirida bitta toʻgʻri oʻlchamli
  slice'ga koʻchiradi (io.go L709–L750); 1 GiB faylda 2096 MiB oʻlchandi (`os.ReadFile`: 1028).
- **`io.Copy` ichida** (io.go L407–L450): avval `src.(WriterTo)`, keyin `dst.(ReaderFrom)` (GO-1.13 type assertion);
  boʻlmasa 32 KiB bufer. Shuning uchun "oʻzgarmas xotira" `io.Copy` ning emas, manba/manzilning xususiyati:
  `io.Copy(&bytes.Buffer{}, f)` `bytes.Buffer.ReadFrom` orqali butun faylni xotiraga yigʻadi.
- **Fayldan faylga** (`nusxa/`): `*os.File.WriteTo` (Go 1.22) faqat socket'ga zero-copy qila oladi, aks holda
  `genericWriteTo` → manzil `*os.File.ReadFrom` → Linux `copy_file_range` (`os/file.go` L265–L274,
  `os/zero_copy_linux.go` L42–L57). `strace -f -qq -e signal=none -e trace=copy_file_range /tmp/io19/nusxa …`:
  bitta chaqiruv `= 1073741824`, keyingisi `= 0` — baytlar dasturga kirmadi. Manzil `/dev/null` yoki pipe boʻlsa
  oddiy read/write ga qaytadi.

- `siqish/` da `gzip.NewWriterLevel` ning xatosi `_` ga tashlangan: u faqat notoʻgʻri daraja konstantasida xato
  qaytaradi, `gzip.BestSpeed` esa toʻgʻri. Oʻz kodingizda darajani tashqaridan olsangiz, xatoni tekshiring.

## Uy vazifasi
- 🟢 **Oson:** `io.Copy` bilan faylni ekranga chiqaring (oʻz `cat` ingiz) va 1 GiB faylda xotirani oʻlchang.
- 🟡 **Oʻrta:** `CRsiz{R io.Reader}` — Windows qator oxiridagi `\r` baytlarini olib tashlaydigan Reader.
- 🟠 **Qiyin:** `Chekli{W io.Writer; Qoldi int64}` — belgilangan baytdan keyin yozishni rad etadigan Writer
  (`n < len(p)` → xato).
- 🔴 **Guru:** `RuneSanagich` — chelak kesgan harflarni ham toʻgʻri sanaydigan Writer (`bufio` siz).

Avval oʻzingiz urinib koʻring. Javoblar [homework/](homework/) da (`CGO_ENABLED=0 go test ./homework/...`):
- Oson ([homework/oson](homework/oson/)): `io.Copy(os.Stdout, f)`;
  `go run ./homework/oson /tmp/io19/buyurtma.csv > /dev/null` → `eng koʻp xotira: 2.xMiB` (stderr'da).
- Oʻrta ([homework/orta](homework/orta/)): ichki Reader'dan oʻqing, `p[:n]` ichida `\r` dan boshqa baytlarni chelak
  boshiga suring, yangi uzunlikni qaytaring. Chelakda faqat `\r` boʻlsa, `(0, nil)` qaytarmang — qayta oʻqing
  (io.go L80–L83). `iotest.TestReader` + `iotest.OneByteReader` bilan test.
- Qiyin ([homework/qiyin](homework/qiyin/)): sigʻganini yozing va `ErrToldi` qaytaring; `io.Copy` toʻxtaydi: `n=5`, `"salom"`.
- Guru ([homework/guru](homework/guru/)): `utf8.FullRune` bilan tugallanmagan 1–3 baytni keyingi `Write` gacha saqlang;
  oqim oxirida qolganini `Close` sanaydi. Test `iotest.OneByteReader` bilan har koʻp baytli rune'ni kesadi.

## Manbalar
- `io` paketi: https://pkg.go.dev/io · manba: https://github.com/golang/go/blob/go1.26.8/src/io/io.go (EOF L37–L40,
  Reader L56–L88, Writer L90–L101, copyBuffer L407–L450, ReadAll L709–L750)
- `os.ReadFile`: https://pkg.go.dev/os#ReadFile · https://github.com/golang/go/blob/go1.26.8/src/os/file.go
- `compress/gzip`: https://pkg.go.dev/compress/gzip#NewWriterLevel
- `testing/iotest`: https://pkg.go.dev/testing/iotest (`DataErrReader`, `OneByteReader`, `TestReader`)
- UTF-8, RFC 3629 §3: https://www.rfc-editor.org/rfc/rfc3629#section-3
- `/proc/[pid]/status` (`VmHWM`): https://man7.org/linux/man-pages/man5/proc_pid_status.5.html
- Linux `copy_file_range(2)`: https://man7.org/linux/man-pages/man2/copy_file_range.2.html
- Go 1.26 release notes (`io.ReadAll`): https://go.dev/doc/go1.26
