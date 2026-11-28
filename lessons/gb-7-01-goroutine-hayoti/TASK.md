# Vazifa · GB-7.01: `fetchAll` — goroutine'larni boshlaydi va ularni kutadi

Videoda uchta yetkazib beruvchining roʻyxatini `main` ichida bir vaqtda oʻqidik: `wg.Go` bilan har `fetch` alohida
goroutine'da, har goroutine oʻz `counts[i]`'iga yozadi, `wg.Wait()` hammasini kutadi. Endi shu ishni importer qayta
ishlata oladigan funksiyaga koʻchiramiz.

`main.go`dagi `fetchAll(suppliers, fetch)` yetkazuvchilar roʻyxatini va `fetch` funksiyasini oladi va
`counts[i] = fetch(suppliers[i])` natijalarini **xuddi `suppliers` tartibida** qaytaradi. Hozir u **ketma-ket** ishlaydi:
uchta yetkazuvchi uch soniya (`go run .` → `jami: 3.0…s`).

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan joylar)

1. **`fetchAll`ni bir vaqtda ishlaydigan qiling.** Har `fetch(name)` alohida goroutine'da boshlansin (`wg.Go`), har
   goroutine faqat oʻz elementiga yozsin (`counts[i] = …`), `fetchAll` esa qaytishdan **oldin** hammasini kutsin
   (`wg.Wait()`). `fetchAll` bu goroutine'larning egasi: u qaytganda ulardan birortasi ham ishlab turmasin.
   `go run .` → `jami: 1.0…s` va `mahsulotlar: [40 120 5]`.
2. 🟡 **Ixtiyoriy: `fetchAllAddDone`.** Xuddi shu funksiyani ishdagi koʻp kodda uchraydigan eski shaklda yozing:
   `wg.Add(1)` `go`dan oldin, ichida `defer wg.Done()`, oxirida `}()`. `wg.Go` ishlatmang. Yozmasangiz, bu tekshiruv `SKIP`.

`offers`, `fetch`, `main` va `tekshiruv_test.go`ni oʻzgartirmang. `wg` uchun fayl boshidagi `import`ga `"sync"`ni qoʻshing.

Laboratoriyada `fetch` xato qaytarmaydi. Haqiqiy yetkazuvchi xato berishi mumkin: xatoni qaytaradigan va birinchi
xatoda qolganlarini toʻxtatadigan variant — `errgroup`, GB-7.04.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. Test soxta `fetch` bilan beshta yetkazuvchini (150, 120, 90, 60, 30 ms) oʻqiydi va uchta narsani
tekshiradi: `fetchAll` qaytganda **hamma** `fetch` tugaganmi; natija `suppliers` tartibidami; ish bir vaqtda boʻldimi
(~150 ms, ketma-ket 450 ms) va `time.Sleep` bilan taxmin qilinmadimi. `tekshiruv_test.go` ichidagi `sync/atomic`
(tugagan `fetch`larni sanaydi) va `go/parser` (🟡 uchun `main.go` matnini oʻqiydi) — tekshiruvning ichki qismi, ularni
hozir oʻqish shart emas. Hozir shunday yiqiladi (vaqt sizda biroz boshqacha):

```
--- FAIL: TestVazifaFetchAll (0.45s)
    tekshiruv_test.go:97: fetchAll 5 ta yetkazuvchini 450ms'da oʻqidi: ular ketma-ket ishlayapti (450ms atrofida). Bir vaqtda boʻlsa, eng sekini, ya'ni ~150ms atrofida tugaydi: har fetch'ni wg.Go bilan alohida goroutine'da boshlang
FAIL
FAIL	github.com/abdurahmon27/bozor	0.457s
FAIL
```

Mashinangizda C kompilyator (`gcc`) boʻlsa, `go test -race .` bilan ham tekshiring: toza chiqishi kerak (README).

**Maslahatlar**

- Videodagi `main`ga qarang (README, "Videodagi qadamlar", `06-index`): `make`, `for i, name := range suppliers`,
  `wg.Go(func() { counts[i] = fetch(name) })`, `wg.Wait()`. Oʻsha toʻrt qator `fetchAll` ichiga koʻchadi.
- Test "N/5 ta fetch hali tugamagan" desa: `wg.Wait()` yoʻq yoki `return`dan keyin turibdi.
- Test tartib haqida yiqilsa: natijani `append` bilan yigʻyapsiz. Har goroutine `counts[i]`ga yozsin.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
