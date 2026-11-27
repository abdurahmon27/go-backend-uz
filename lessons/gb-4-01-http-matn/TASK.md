# Vazifa · GB-4.01: `parseRequestLine` va `statusClass`

Videoda `curl -v` bergan soʻrovning birinchi qatori `GET /products/1 HTTP/1.1` edi: **method**, **yoʻl** va
**versiya** boʻshliq bilan ajratilgan uch qism. Javobning birinchi qatoridagi status kodi esa birinchi raqami bilan
sinfga boʻlinadi (`2xx`, `4xx`, `5xx`…). `main.go`'da shu ikkitasi uchun funksiyalar bor, lekin ular hali boʻsh.

**Nima qilish kerak** (`main.go`'da, faqat shu ikki funksiyaning tanasini yozasiz):

1. `func parseRequestLine(line string) (Request, error)` qatorni `Request{Method, Path, Proto}`ga ajratsin:
   `"GET /products/1 HTTP/1.1"` → `Request{"GET", "/products/1", "HTTP/1.1"}`. Uch holatda `Request{}` va xato
   qaytarsin:
   - qator **aynan uch** boʻshliq bilan ajratilgan qismdan iborat emas (boʻsh qator, ikki qism, toʻrt qism,
     ketma-ket ikki boʻshliq);
   - yoʻl `/` bilan boshlanmaydi (`products/1`);
   - versiya `HTTP/` bilan boshlanmaydi (`FTP/1.1`).
   Soʻrovdagi `?limit=2` kabi qism yoʻlning ichida qoladi: uni ajratmaysiz. Qator oxiridagi `\r` va `\n` belgilari
   (videodagi `\r\n`) qismga kirmasin: `"GET /products/1 HTTP/1.1\r\n"` ham toʻgʻri qator (`strings.TrimRight(line, "\r\n")`).
2. `func statusClass(code int) string` status kodiga qarab sinf nomini qaytarsin:

   | kod | qaytariladi |
   |---|---|
   | 100–199 | `axborot` |
   | 200–299 | `muvaffaqiyat` |
   | 300–399 | `yoʻnaltirish` |
   | 400–499 | `mijoz xatosi` |
   | 500–599 | `server xatosi` |
   | qolgani | `nomaʼlum` |

   (`yoʻnaltirish`dagi `ʻ` — U+02BB, `nomaʼlum`dagi `ʼ` — U+02BC: ikkalasini `tekshiruv_test.go`'dan nusxa oling.)

Xato matni sizning ixtiyoringizda (`fmt.Errorf`, GB-1.07), test faqat xato bor-yoʻqligini tekshiradi.

**Tekshirish** (shu papkada):

```bash
go test
```

`ok` chiqquncha ishlang. Test fayli `tekshiruv_test.go`: uni oʻzgartirmang. Hozir u shunday yiqiladi
(qisqartirilgan):

```
--- FAIL: TestVazifaParseRequestLine (0.00s)
    --- FAIL: TestVazifaParseRequestLine/get_mahsulot (0.00s)
        tekshiruv_test.go:39: parseRequestLine("GET /products/1 HTTP/1.1"): xato kutilmagan, keldi: parseRequestLine hali yozilmagan (TASK.md, 1-qadam)
…
--- FAIL: TestVazifaStatusClass (0.00s)
    tekshiruv_test.go:70: statusClass(200) = ""; kutilgan "muvaffaqiyat" (TASK.md, 2-qadam)
…
FAIL
```

**Maslahatlar**

- `parseRequestLine`: `parts := strings.Split(line, " ")` qatorni boʻshliq boʻyicha boʻladi (GB-3.01); ketma-ket
  ikki boʻshliq boʻsh qism beradi, shuning uchun `len(parts) != 3` bu holni ham tutadi.
- `strings.HasPrefix(s, "/")` — `s` shu matndan boshlanadimi (`true`/`false`).
- `statusClass`: `if code >= 200 && code < 300 { return … }` zanjiri yoki shartli `switch`. Eng oxirida
  `nomaʼlum`.

Qoʻshimcha (testsiz): server bilan oʻzingiz sinab koʻring: `go run ./_video/server` bir terminalda, boshqasida
`curl -v http://localhost:8765/products/1`; `-v` chiqishidagi `>` qatorining birinchisini `parseRequestLine`ga bering.

Javob: [`solution/main.go`](solution/main.go).
