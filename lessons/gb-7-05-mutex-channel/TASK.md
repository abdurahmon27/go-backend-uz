# Vazifa · GB-7.05: egasi goroutine'ga `Stock` soʻrovi

Videoda `OwnedStore` qoldiq map'ini bitta goroutine'ga (`loop`) berdi: `Sell` map'ga tegmaydi, `sells` channel'iga
`sellReq` yuboradi va `reply`dan javob kutadi. GB-4.07'dagi `Store`da yana bir metod bor edi: `Stock(id)` — qoldiqni
qaytaradi. `OwnedStore`da u hozir hech narsa soʻramaydi va `-1` qaytaradi. Uni **qulfsiz**, egasi orqali yozing.

**Nima qilish kerak** (`store.go`, `VAZIFA` deb belgilangan joylar)

1. 🟢 **`Stock(id int64) int`.**
   - Yangi soʻrov turi: `stockReq` struct'i — `id int64` va javob uchun `reply chan int`.
   - `OwnedStore`ga ikkinchi channel maydoni (masalan `stocks chan stockReq`); `NewOwnedStore` uni ham `make` bilan yaratsin.
   - `Stock` `sellReq`dagidek javob channel'ini yaratadi, soʻrovni yuboradi va javobni qaytaradi. Yoʻq mahsulot: `0`.
   - `loop` endi ikki channel'ni kutadi: `for { select { case …: … case …: … } }` (`select`: GB-7.03).
   - `Close` `sells`ni yopadi. Yopiq channel'dan oʻqish kutmaydi, darhol nol qiymat beradi (GB-7.02), shuning uchun
     `case r, ok := <-s.sells:` va `ok` yolgʻon boʻlsa `return` — aks holda `loop` hech qachon tugamaydi.
   - Map faqat `loop`da yashaydi: `OwnedStore`ga map ham, `sync.Mutex` ham qoʻshmang (tekshiruv `OwnedStore` maydonlarida
     map, `sync.Mutex`, `sync.RWMutex`, `sync.Locker` va ularga koʻrsatkichni qidiradi).
2. 🟡 **Ixtiyoriy: `Requests() int64`.** `OwnedStore`dagi `calls atomic.Int64` maydoni tayyor. `Sell` har chaqiruvda
   uni bittaga oshirsin (qoldiq yetmagan chaqiruv ham sanaladi), `Requests` uni oʻqisin. `Sell`ni bir vaqtda koʻp
   goroutine chaqiradi, shuning uchun oddiy `int` va `++` emas, `Add` va `Load`. Yozmasangiz, bu tekshiruv `SKIP`.

`LockedStore`, `Sell`ning ishlash tartibi, `main.go` va `tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. Tekshiruv: `Stock(1)`, `Stock(3)`, `Stock(99)` qiymatlari; 8 goroutine oʻrikni (5 dona) sotayotgan
paytda yana 8 tasi `Stock` soʻraydi — hammasi 2 soniyada tugashi, aniq 5 ta sotuv oʻtishi va qoldiq 0 boʻlishi kerak;
`OwnedStore`da map va qulf yoʻqligi; 20 ta store `Close`dan keyin ularning `loop` goroutine'lari tugashi. Tekshiruv
faylidagi `reflect` (struct maydonlarini koʻradi) va `runtime.NumGoroutine` (goroutine'lar soni) — uning ichki qismi,
ularni oʻqish shart emas. Hozir shunday yiqiladi:

```
--- FAIL: TestVazifaStock (0.00s)
    tekshiruv_test.go:45: Stock(1) = -1: Stock hali yozilmagan. TASK.md 1-qadam: Stock soʻrovini loop'ga yuborib, javobni qaytaring
FAIL
```

`Close`dan keyin `loop` tugamasa, xabar shunday boʻladi (yechimdan `ok` tekshiruvini olib tashlab olingan haqiqiy natija):

```
--- FAIL: TestVazifaClose (2.01s)
    tekshiruv_test.go:130: Close'dan keyin loop goroutine'lari tugamadi (goroutine'lar: oldin 3, hozir 23).
        	Yopiq channel'dan oʻqish darhol nol qiymat beradi (GB-7.02): select'da `r, ok := <-s.sells` va `if !ok { return }`
```

**Maslahatlar**

- `Stock`ning shakli `Sell` bilan bir xil: `reply := make(chan int)`, `s.stocks <- stockReq{id, reply}`, `return <-reply`.
- `loop`dagi `select`ning ikkinchi `case`i bitta qator ish qiladi: `r.reply <- stock[r.id]`.
- `go test` "2 soniyada javob kelmadi" desa: `stocks` channel'i `make` bilan yaratilganini (`nil` channel'ga yuborish
  abadiy kutadi) va `loop`ning `select`ida shu channel borligini tekshiring.
- `-race` bilan ham tekshirib koʻring (Linux'da gcc kerak, README): `go test -race .` — ogohlantirish chiqmasligi kerak.

Javob: [`solution/store.go`](solution/store.go). Javobni tekshirish: `go test ./solution/`.
