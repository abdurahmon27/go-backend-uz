# Vazifa · GB-4.07: `Restock` va `Total`ni qulf bilan yozing

Videoda `Sell`ni `sync.Mutex` bilan himoyaladik, `go test -race` toza chiqdi. `main.go`da yana ikki metod bor:
`Restock` allaqachon yozilgan, lekin **qulfsiz** (xuddi videodagi `Sell`ning oldingi holati), `Total` esa hali yozilmagan.
Server har ulanishni alohida goroutine'da ishlatadi, shuning uchun `Sell`, `Restock` va `Total` bir vaqtda chaqirilishi mumkin:
`products` mapiga **har** murojaat bir xil qulf ostida boʻlishi kerak.

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan joylar)

1. **`Restock`ni himoyalang.** Metodning boshiga `s.mu.Lock()` va keyingi qatorga `defer s.mu.Unlock()` qoʻying
   (videodagi `Sell`dagidek). Qolgan mantiq oʻzgarmaydi: mahsulot yoʻq boʻlsa `ErrNotFound`.
2. **`Total`ni yozing.** Hamma mahsulotlar `Stock`larining yigʻindisini qaytarsin (`map` ustida `for _, p := range s.products`).
   Qulf `Total` boshida **bir marta** olinadi. `Total` ichida `s.Stock(...)`ni chaqirmang: `Stock` ham `Lock` chaqiradi,
   qulfni ushlab turgan holda ikkinchi `Lock` oʻzini oʻzi kutib, dastur qotib qoladi (deadlock).

`Product`, `Store`, `Sell`, `Stock`, `buy`, `slowHandler`, `main`, `buy_test.go` va `tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. `tekshiruv_test.go` oʻzgartirilmaydi: u `Sell`, `Restock` va `Total`ni yuzlab goroutine'dan bir vaqtda
chaqiradi. Qulfsiz kod ba'zan oʻzi qulaydi (`fatal error: concurrent map writes`), shuning uchun tekshiruv ikkinchi jarayonda
ishlaydi va chiqishini oʻzbekcha xabarga aylantiradi. Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaTotalOddiy (0.00s)
    tekshiruv_test.go:33: Total() = 0, kutilgan 165 (40 + 120 + 5): Total hamma mahsulotlarning qoldigʻini qoʻshishi kerak
--- FAIL: TestVazifaParallel (0.01s)
    tekshiruv_test.go:47: bir vaqtda 900 ta chaqiruvda (Sell + Restock + Total) store buzildi yoki qotib qoldi (1-urinish).
        Tekshiring: Restock va Total ham s.mu.Lock() + defer s.mu.Unlock() ostida boʻlsin; …
        Tekshiruv jarayonining chiqishi (boshi):
        fatal error: concurrent map writes
…
FAIL
```

Mashinangizda C kompilyator (`gcc`) boʻlsa, hisobotni oʻzingiz oʻqing: `go test -race .` (qaysi qator `Restock`ni koʻrsatadi?).
`gcc` yoʻq boʻlsa, `-race` ishlamaydi (`go: -race requires cgo`): bu normal, yuqoridagi `go test .` yetadi.

**Maslahatlar**

- `Sell`ning boshiga qarang (`main.go`): ikki qator, `Lock` va `defer Unlock`, oʻsha ikki qatorni `Restock`ga koʻchiring.
- `Total`: `total := 0`, `for _, p := range s.products { total += p.Stock }`, `return total`. Qulf bu uch qatordan **oldin**.
- `Total` qotib qolsa, ichida `Stock(...)` yoki `Sell(...)` chaqirilmaganini tekshiring.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
