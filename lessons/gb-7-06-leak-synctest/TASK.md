# Vazifa · GB-7.06: `cheapestPrice` — eng arzon taklif, leak'siz

Videoda `fastestPrice` uch yetkazib beruvchidan birdan soʻradi va **birinchi** javobni oldi; `synctest.Test` ortda qolgan
goroutine'larni ushladi, `context.WithCancel` + `defer cancel()` esa ularni chaqiruv bilan birga tugatdi. Endi bozor
**eng arzon** taklifni xohlaydi, lekin cheksiz kutolmaydi.

`offer.go`da hammasi tayyor, faqat bitta funksiya qolgan:

- `priceOffers` — uch yetkazib beruvchi (`priceOffer{delay, price}`): 100 ms'da 1 250 000 tiyin, 200 ms'da 1 150 000,
  3 soniyada 1 100 000 (eng arzoni eng sekini);
- `askOffer(ctx, o)` — bitta yetkazib beruvchidan soʻrash: `o.delay` oʻtguncha yoki `ctx` bekor qilinguncha kutadi
  (videodagi `quote` bilan bir xil);
- `quoteTimeout = time.Second`, `errNoOffer`.

**Nima qilish kerak:** `cheapestPrice(ctx context.Context) (int64, error)`ni yozing.

1. Hamma `priceOffers`dan **birdan** soʻrang: har taklif uchun bitta goroutine, ichida `askOffer`.
2. `quoteTimeout` (1 soniya) ichida kelgan takliflarning **eng arzonini** qaytaring (`best, nil`). Hamma taklif oldinroq
   kelsa, 1 soniyani kutmang: oxirgi javob kelishi bilan qayting.
3. 1 soniya ichida hech kim javob bermasa: `0, errNoOffer`.
4. **Chaqiruvchining** `ctx`'i oldinroq bekor qilinsa yoki muddati tugasa: `0, ctx.Err()` (masalan `context deadline exceeded`).
5. Qaytganda **birorta ham goroutine ortda qolmasin**: sekin yetkazib beruvchining soʻrovi toʻxtatilsin, natija channel'iga
   yuborish hech qachon osilib qolmasin.

`askOffer`, `priceOffers`, `price.go`, `store.go` va test fayllarini oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. `tekshiruv_test.go` har holatni synctest pufagida tekshiradi: soat soxta, shuning uchun
3 soniyalik taklif bir zumda "keladi" va vaqt **aniq** solishtiriladi (`1s`, `300ms`, `150ms`). Chaqiruvdan keyin u
`synctest.Wait()` qiladi va ortda qolgan goroutine'larni sanaydi. Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaEngArzon (0.00s)
    tekshiruv_test.go:57: cheapestPrice() = 0, VAZIFA: cheapestPrice hali yozilmagan (TASK.md) (0s oʻtdi); kutilgan 1150000, <nil>, aniq 1s.
        	1 soniya ichida ikki taklif keladi, arzoni 1 150 000; 3 soniyalik taklifni kutmang (quoteTimeout).
--- FAIL: TestVazifaHammasiTez (0.00s)
…
--- FAIL: TestVazifaChaqiruvchi (0.00s)
…
FAIL
```

Goroutine qoldirsangiz, avval shu xabar, keyin videodagi `deadlock` chiqadi:

```
    tekshiruv_test.go:45: goroutine leak: cheapestPrice qaytganda 1 ta goroutine hali kutib turibdi.
…
panic: deadlock: main bubble goroutine has exited but blocked goroutines remain [recovered, repanicked]
```

(`[recovered, repanicked]` — `testing` paketi bu panic'ni ushlab, test nomini yozib, qayta chiqargan; ma'nosi videodagi bilan bir xil.)

**Maslahatlar**

- Boshlanishi videodagi `fastestPrice`ga oʻxshaydi, faqat chaqiruvning oʻz muddati bor: `context.WithTimeout(ctx, quoteTimeout)`.
  `defer cancel()`ni unutmang (`go vet ./...` buni `lostcancel` deb aytadi). Natija channel'ining buffer'i: takliflar soni —
  muddat bilan bir paytda kelgan javob ham osilib qolmasin.
- Bitta `return <-ch` endi yetmaydi: javoblarni sikl ichida yigʻing, har aylanishda `select` bilan javob yoki muddat.
  "Hali taklif yoʻq" holatini alohida belgilang (`0` narxga tayanmang).
- Muddat tugaganda sababini ajrating: chaqiruvchining `ctx.Err()` nil boʻlmasa, oʻsha xato; aks holda bu sizning `quoteTimeout`ingiz.
