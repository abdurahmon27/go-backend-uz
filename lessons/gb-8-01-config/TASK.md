# Vazifa · GB-8.01: `HOST` sozlamasi

Hozir bozor API'si har doim kompyuterning **hamma** tarmoq manzillarida tinglaydi (`addr()` → `":8765"`). Server oldida
boshqa server turganda (masalan, GB-8.06'dagi Caddy) API'ga faqat shu kompyuterdan kirish kerak: `127.0.0.1:8765`.
Buni kodni oʻzgartirmasdan, sozlama bilan tanlash kerak.

**Spetsifikatsiya** (`DATABASE_URL` qoʻyilgan, `PORT` yoʻq)

```
HOST yoʻq                      → addr() ":8765"            (hamma manzillar, hozirgidek)
HOST=127.0.0.1                 → addr() "127.0.0.1:8765"   (faqat shu kompyuter)
HOST=localhost                 → addr() "localhost:8765"
HOST=::1                       → addr() "[::1]:8765"       (IPv6 loopback: shu kompyuter)
HOST=0.0.0.0  va  -host 127.0.0.1 → addr() "127.0.0.1:8765"   (flag muhitdan ustun)
HOST=127.0.0.300               → xato: host "127.0.0.300": …   (bunday IP yoʻq)
HOST=bozor.uz                  → xato                       (domen emas: faqat IP yoki localhost)
```

Xato boshqa xatolar bilan **bitta** xabarda chiqadi: `HOST=bozor.uz` va `DATABASE_URL` yoʻq boʻlsa — ikki qator.

**Nima qilish kerak** (`cmd/api/config.go`)

1. `config`ga `host string` maydoni va `loadConfig`da `-host` flag'i: standart qiymati `envOr("HOST", "")` (`-port` qanday
   yozilgan boʻlsa, shunday).
2. Tekshiruv: `host` boʻsh emas, `"localhost"` emas va IP ham boʻlmasa — `errs`ga xato qoʻshing. Xato matnida `host`
   soʻzi va qiymat (`%q`) boʻlsin, masalan `host "bozor.uz": IP manzil yoki localhost kerak (HOST yoki -host)`.
3. `addr()` host'ni ham qoʻshsin: `""` → `":8765"`, `"127.0.0.1"` → `"127.0.0.1:8765"`, `"::1"` → `"[::1]:8765"`.

`run` allaqachon `newServer(cfg.addr(), …)`ni chaqiradi: boshqa joyni oʻzgartirish shart emas.
`cmd/api/tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (papkaning ildizida; Postgres kerak emas):

```bash
go test ./cmd/api -run Config
```

Hozir shunday yiqiladi (qisqartirilgan):

```
flag provided but not defined: -host
…
--- FAIL: TestConfigHost (0.00s)
    tekshiruv_test.go:85: HOST=127.0.0.1 []: addr() ":8765", kutilgan "127.0.0.1:8765" (TASK.md, 3-qadam: addr() host'ni ham qoʻshadi, net.JoinHostPort)
    …
    tekshiruv_test.go:80: HOST=0.0.0.0 [-host 127.0.0.1]: xato flag provided but not defined: -host
          kutilgan addr() "127.0.0.1:8765" (TASK.md, 1-qadam: -host flag'i, standarti envOr("HOST", ""))
    tekshiruv_test.go:94: HOST="127.0.0.300" va DATABASE_URL yoʻq: xato DATABASE_URL qoʻyilmagan: .env.example'ga qarang
          kutilgan ikkala muammo bitta xatoda: host "127.0.0.300" … va DATABASE_URL … (TASK.md, 2-qadam)
    …
FAIL	github.com/abdurahmon27/bozor/cmd/api
```

`ok` boʻlguncha ishlang. Keyin butun papka (API testlari haqiqiy Postgres'da):

```bash
docker compose up -d --wait
go test ./...
docker compose down -v
```

Qoʻlda ham tekshiring: `HOST=127.0.0.1 ./api` (`go build ./cmd/api`dan keyin) va ikkinchi terminalda
`ss -ltn | grep 8765` — `127.0.0.1:8765` koʻrinadi, `*:8765` emas.

**Maslahatlar**

- `net.ParseIP(s)` — `s` IP manzil boʻlmasa `nil` qaytaradi (IPv4 ham, IPv6 ham): `go doc net.ParseIP`.
- `net.JoinHostPort(host, port)` — `host` va `port`ni `host:port` qiladi, IPv6'ni `[` `]`ga oladi; `port` — matn:
  `strconv.Itoa(c.port)`.
- Flag'ni `fs.Parse`dan **oldin** roʻyxatga oling, aks holda `-host` notanish flag boʻladi.

Javob: [`solution/cmd/api/config.go`](solution/cmd/api/config.go). Tekshirish: `cd solution && go test ./cmd/api -run Config`.
