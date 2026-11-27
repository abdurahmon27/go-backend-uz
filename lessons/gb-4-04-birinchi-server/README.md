# GB-4.04 · Birinchi server, qatorma-qator: `ServeMux` va `ListenAndServe`

Kurs: **Go backend: noldan guru'gacha** · loyiha: **bozor** · video: YouTube'dagi "Golang backend #28"
(kanaldagi "Go backend" playlist'i) · yozma dars: [haywan.uz/lab](https://haywan.uz/lab)

> **Go versiyasi.** Modul 2 va undan keyingi darslar Go 1.27 bilan yozilgan (`go.mod`: `go 1.27`). Sizda eskiroq Go
> boʻlsa, `go run`/`go test` birinchi marta kerakli versiyani oʻzi yuklab oladi (`go: downloading go1.27…`). Bu oddiy
> hol; `GOTOOLCHAIN=local` qoʻygan boʻlsangiz, Go'ni go.dev/dl'dan yangilang.

Oldingi dars: [`gb-4-03-handler-httptest`](../gb-4-03-handler-httptest) (handler — oddiy funksiya, `httptest`) · keyingi dars:
[`gb-4-05-route-id`](../gb-4-05-route-id) (route'lar va `{id}`: `GET /products/{id}`, 404 va 405). Bu darsga kerak boʻlganlar:
[`gb-4-02-http-client`](../gb-4-02-http-client) (ikkinchi terminal, port, `Ctrl+C` → `signal: interrupt`, client `Timeout`,
`time.Second`), [`gb-4-01-http-matn`](../gb-4-01-http-matn) (`curl`, soʻrov va status qatori, `404`),
[`gb-1-07-funksiya-error`](../gb-1-07-funksiya-error) (`fmt.Fprintln(os.Stderr, …)`, `os.Exit(1)`, `echo $?`).

**Bu darsdan keyin:** `http.NewServeMux` va `mux.HandleFunc(pattern, handler)` bilan yoʻl → handler jadvalini yasaysiz,
`http.ListenAndServe(":8765", mux)` portni ochib nega qaytmasligini va qaytsa doim xato ekanini bilasiz, serverni ikkinchi
terminaldan `curl` bilan tekshirib `Ctrl+C` bilan toʻxtatasiz, `bind: address already in use`ni oʻqiysiz, `http.HandleFunc` +
`nil` (DefaultServeMux) nega ishda ishlatilmasligini va `http.Server` struct'ini (`Addr`, `Handler`, `ReadHeaderTimeout`) bilasiz.

Bu papka mustaqil: oʻz `go.mod`'i (`github.com/abdurahmon27/bozor`, `go 1.27`), oʻz `package main`'i, tashqi modul yoʻq.
`Product`, `olma`, `catalog`, `healthHandler`, `productHandler`, `catalogHandler` GB-4.03'dan (catalogHandler — oʻsha darsning
vazifasi, bu yerda yechilgan holda). Yangi: `newServer(addr string) *http.Server` (vazifa), server manzili `:8765` (GB-4.01'dagi port).

Oldingi darsdan nima oʻzgardi: `diff -ru ../gb-4-03-handler-httptest/solution/main.go solution/main.go` (handler'lar oʻsha-oʻsha;
`main` endi handler'ni chaqirmaydi, server ochadi).

## Papka

| Fayl | Nima |
|---|---|
| `main.go` | handler'lar (GB-4.03) va tayyor `main`; `newServer` boʻsh — vazifa (`nil` qaytaradi, `go run .` buni aytadi) |
| `tekshiruv_test.go` | vazifani tekshiradi (`go test .`); oʻzgartirmang |
| `TASK.md` | vazifa |
| `solution/` | javob: `main.go` (+ xuddi shu test fayli) |
| `_video/handlers.go` | GB-4.03 handler'lari: har bir video qadamiga qoʻshib koʻchiriladi |
| `_video/NN-…/main.go` | videodagi har bir qadamning `main` funksiyasi; `_` bilan boshlangan papkani `go test ./...`/`go vet ./...` chetlab oʻtadi |
| `_video/run.sh` | qadamni videodagidek yangi `bozor/` papkasida ishga tushiradi: `./_video/run.sh 02-server go run .` |

## Videodagi buyruqlar va natijalari (go1.27.4 linux/amd64, `CGO_ENABLED=0`, curl 8.22.0)

Server yoʻq — portni hech kim tinglamayapti (curl xato kodi 7, satrning bir qismi videoda `…` bilan yashirilgan):

```
$ curl http://localhost:8765/health
curl: (7) Failed to connect to localhost:8765 after 0 ms: Could not connect to server
```

Server (`02-server`, birinchi terminal): manzil chiqadi va buyruq satri qaytmaydi — dastur `ListenAndServe` ichida kutyapti:

```
$ go run .
bozor API: http://localhost:8765
```

Ikkinchi terminal:

```
$ curl http://localhost:8765/health
ok
$ curl http://localhost:8765/products
[{"id":1,"name":"olma","price":1200000,"stock":40},{"id":2,"name":"non","price":400000,"stock":120},{"id":3,"name":"oʻrik","price":3500000,"stock":5}]
$ curl -i http://localhost:8765/products/1
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 27 Nov 2026 10:32:55 GMT
Content-Length: 19

404 page not found
```

(`Date` sizda boshqa boʻladi. Videoda `404` javobining sarlavhalari `…` bilan yashirilgan.) `"/products"` patterni oxiri slesh
bilan tugamagan, shuning uchun faqat aynan `/products` yoʻliga mos; `/products/1` uchun mos pattern yoʻq va `ServeMux` oʻzi
`404 page not found` yozadi. Birinchi terminalda `Ctrl+C`:

```
bozor API: http://localhost:8765
^Csignal: interrupt
```

Shundan keyin ikkinchi terminaldagi `curl` yana `Failed to connect` beradi.

**Eng koʻp uchraydigan xato** — `ListenAndServe` qaytargan xatoni tashlab yuborish (`03-xato-yoq`). Birinchi terminalda server
ishlab turganda ikkinchisida shu nusxa:

```
$ go run .
bozor API: http://localhost:8765
$ echo $?
0
```

Dastur jim tugadi, chiqish kodi `0`, lekin server ishlamayapti. Xato tekshiriladigan nusxa (`04-xato-bor`, `02-server` bilan bir xil):

```
$ go run .
bozor API: http://localhost:8765
xato: listen tcp :8765: bind: address already in use
exit status 1
```

`bind` — portni dasturga biriktiradigan operatsion tizim chaqiruvi; `EADDRINUSE` — bu manzil boshqa dasturda
([bind(2)](https://man7.org/linux/man-pages/man2/bind.2.html)). Portni kim egallaganini koʻrish (videoda yoʻq):
`ss -ltnp 'sport = :8765'`.

`http.HandleFunc` va `nil` (`05-default`) ham ishlaydi (`curl …/health` → `ok`, `/products` → `404`, chunki u roʻyxatga olinmagan),
lekin yoʻllar `http.DefaultServeMux`ga — `net/http` ichidagi umumiy mux'ga — yoziladi.

`http.Server` (`06-http-server`): xuddi shu natijalar (`ok`, `/products` → `200`).

## `ReadHeaderTimeout` oʻlchovi (videoda raqamlari aytiladi)

Hech narsa yubormaydigan TCP ulanish (bash'ning `/dev/tcp`'i bilan; birinchi terminalda server ishlab turganda):

```bash
bash -c 'exec 3<>/dev/tcp/127.0.0.1/8765; s=$SECONDS; timeout 12 cat <&3; echo "kod $? · $((SECONDS-s)) s"'
```

| server | natija |
|---|---|
| `06-http-server` (`ReadHeaderTimeout: 5 * time.Second`) | `kod 0 · 5 s` — server ulanishni 5 soniyada yopdi, javob baytlari yoʻq |
| `02-server` (`http.ListenAndServe`, muddat yoʻq) | `kod 124 · 12 s` — 12 soniyadan keyin ham ochiq (`timeout` uzdi) |

`bind` xatosida `listen tcp`: TCP — HTTP ishlaydigan ulanish protokoli, `listen` — portni tinglash uchun ochish.

## Nima uchun shunday

- `http.NewServeMux()` boʻsh `*http.ServeMux` qaytaradi. `mux.HandleFunc(pattern, f)` jadvalga yozuv qoʻshadi; `f` qavssiz
  beriladi (funksiya qiymati), mux uni har mos soʻrovda `f(w, r)` qilib chaqiradi. `ServeMux`da `ServeHTTP` metodi bor, shuning
  uchun u `http.Handler` (GB-4.03) va `ListenAndServe`ga handler sifatida beriladi.
- Oxiri `/` bilan tugamagan pattern bitta yoʻlga mos; oxiri `/` bilan tugagani shu prefiksli hamma yoʻlga
  ([ServeMux, Patterns](https://pkg.go.dev/net/http#hdr-Patterns-ServeMux)). Go 1.22'dan beri patternda metod (`GET /products`)
  va `{id}` ham boʻladi — GB-4.05.
- `ListenAndServe` server ishlab turganda qaytmaydi: ulanishlarni qabul qilib, har birini alohida goroutine'da xizmat qiladi
  ([Server.Serve](https://pkg.go.dev/net/http#Server.Serve); bir vaqtda kelgan soʻrovlar — GB-4.07). U server toʻxtaganda
  qaytadi: ishga tusha olmasa (band port) yoki dastur uni `Shutdown`/`Close` bilan toʻxtatsa (`http.ErrServerClosed`, GB-8.03),
  va qaytgan qiymat doim `nil` boʻlmagan xato. Uni tashlab yuborish — band portda jim `exit 0` (yuqorida).
- `bozor API: http://localhost:8765` qatori `ListenAndServe`dan **oldin** chiqadi, shuning uchun u server ishlayotganini
  isbotlamaydi (band portda ham chiqdi). Isbot — `curl`.
- Patternda metod yozilmagan, shuning uchun u har qanday metodni oladi: `curl -i -X POST http://localhost:8765/products` →
  `405 Method Not Allowed`, `Allow: GET` (buni `catalogHandler`ning oʻzi qaytaradi, GB-4.03). Metodli pattern — GB-4.05.
- `":8765"` — host yoʻq, ya'ni shu kompyuterning hamma manzillari ([net.Listen](https://pkg.go.dev/net#Listen); `ss -ltn`
  `*:8765` koʻrsatadi): bir tarmoqdagi (masalan bitta Wi-Fi'dagi) boshqa kompyuter ham ulana oladi. Faqat shu kompyuterdan
  ulanish kerak boʻlsa: `"localhost:8765"` yoki `"127.0.0.1:8765"`. Konteynerda esa aynan hamma manzillar kerak (GB-8).
- `http.HandleFunc` / `ListenAndServe(addr, nil)` — `DefaultServeMux`. U global: import qilingan istalgan paket unga yoʻl
  qoʻshadi, masalan `net/http/pprof` `/debug/pprof/` yoʻllarini (profil ma'lumoti, buyruq qatori) ochadi
  ([net/http/pprof](https://pkg.go.dev/net/http/pprof)). Shuning uchun oʻz `ServeMux`imiz.
- `http.ListenAndServe(addr, h)` = `(&http.Server{Addr: addr, Handler: h}).ListenAndServe()`
  ([server.go](https://cs.opensource.google/go/go/+/refs/tags/go1.27.4:src/net/http/server.go)). `http.Server`da muddatlar
  bor; `ReadHeaderTimeout` 0 boʻlsa (va `ReadTimeout` ham 0 boʻlsa) sarlavhalarni kutish cheksiz. `ReadHeaderTimeout` faqat
  soʻrov sarlavhalarini kutishni cheklaydi: bitta soʻrov yuborib, keyin jim turgan keep-alive ulanish uni chetlab oʻtadi (keyingi
  soʻrovni kutish `IdleTimeout`, u ham 0 boʻlsa `ReadTimeout` bilan cheklanadi; ikkalasi ham 0 boʻlsa — cheksiz). Qolgan muddatlar
  (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`) va `Shutdown` — GB-8.03.
- Server har soʻrovni log'ga yozmaydi (ikkinchi terminalda curl ishlaganda birinchisida hech narsa chiqmaydi): log — GB-8.02.

## Manbalar

- [`net/http`: ListenAndServe](https://pkg.go.dev/net/http#ListenAndServe), [Server](https://pkg.go.dev/net/http#Server)
  (`Addr`, `Handler` — nil boʻlsa DefaultServeMux, `ReadHeaderTimeout`), [Server.ListenAndServe](https://pkg.go.dev/net/http#Server.ListenAndServe)
  ("always returns a non-nil error"), [ServeMux](https://pkg.go.dev/net/http#ServeMux), [NewServeMux](https://pkg.go.dev/net/http#NewServeMux),
  [ServeMux.HandleFunc](https://pkg.go.dev/net/http#ServeMux.HandleFunc), [HandleFunc](https://pkg.go.dev/net/http#HandleFunc),
  [DefaultServeMux](https://pkg.go.dev/net/http#DefaultServeMux).
- [`net`: Listen](https://pkg.go.dev/net#Listen) (boʻsh host), [bind(2)](https://man7.org/linux/man-pages/man2/bind.2.html) (`EADDRINUSE`).
- [`net/http/pprof`](https://pkg.go.dev/net/http/pprof) (DefaultServeMux'ga yoʻl qoʻshadi).
- [`net/http/httptest`: NewServer](https://pkg.go.dev/net/http/httptest#NewServer) (tekshiruv shu bilan).
