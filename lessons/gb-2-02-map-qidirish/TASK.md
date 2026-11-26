# Vazifa · GB-2.02 · `reserve`: buyurtmani ombordan ayirish

Videoda bozor buyurtmadagi nomlarni `prices` map'idan comma ok bilan topdi, katalogda yoʻq mahsulotni rad etdi va
donalarni `count` map'ida sanadi. Endi omborni ham hisobga oling: mahsulot sotilsa, qoldigʻi kamaysin, yetmasa
buyurtma umuman oʻtmasin.

`main.go`'da `reserve` funksiyasi bor, lekin hozir hech narsa qilmaydi (`return nil`):

```go
func reserve(stock map[string]int, order []string) error
```

`stock` — ombor: mahsulot nomi → qoldiq (dona). `order` — buyurtmadagi nomlar, bir nom bir necha marta kelishi
mumkin (`olma non olma` = 2 olma, 1 non).

## 1-qadam: ayirish

Hammasi yetarli boʻlsa, `reserve` `nil` qaytaradi va har bir mahsulot qoldigʻi buyurtmadagi donasiga kamayadi.
Ombor `olma 5, non 2, choy 1` boʻlsa, `reserve(stock, []string{"olma", "non", "olma"})` dan keyin:
`olma 3, non 1, choy 1`.

## 2-qadam: omborda yoʻq mahsulot

Buyurtmadagi nom `stock`'da kalit sifatida yoʻq boʻlsa, xato qaytaring, matni aynan shunday:

```
anor omborda yoʻq
```

(`main`'dagi `katalogda yoʻq` xatosi `prices` katalogi haqida; `reserve` esa faqat omborni koʻradi.)

## 3-qadam: yetmaydi

Mahsulot bor, lekin buyurtmadagi donasi qoldiqdan koʻp boʻlsa, xato qaytaring:

```
choy: omborda 1 dona, buyurtmada 2
```

**Muhim:** xato boʻlsa, ombor **umuman oʻzgarmasin**, buyurtmadagi boshqa, yetarli mahsulotlar ham. `olma non non non`
buyurtmasida (`non` 2 dona) xato `non: omborda 2 dona, buyurtmada 3` boʻladi va `olma` 5 boʻlib qoladi. Bir nechta
muammo boʻlsa, buyurtmada birinchi kelgan muammoli nom haqida xabar bering: `anor choy choy` → `anor omborda yoʻq`.

## Tekshirish (shu papkada)

```bash
go test
```

`ok` chiqquncha ishlang. `tekshiruv_test.go` `reserve`'ni yuqoridagi buyurtmalar bilan chaqiradi; uni oʻzgartirmang.

## Maslahatlar

- `stock` — map, funksiyaga jadvalning oʻzi emas, faqat uning xotiradagi manzili nusxalanadi (videodagi `restock`): siz
  `stock["olma"] -= 2` deb yozsangiz, chaqiruvchining ombori oʻzgaradi. Shuning uchun avval **hammasini tekshiring**,
  keyin ayiring.
- Avval donalarni videodagi kabi yangi map'da sanang (`make`, `count[name]++`), keyin `order` ustida yurib har bir nomni
  `have, ok := stock[name]` bilan tekshiring.
- Xato matni uchun `fmt.Errorf("%s omborda yoʻq", name)` (GB-1.07). `yoʻq`dagi `ʻ` — oddiy apostrof emas, alohida
  belgi (U+02BB): matnni shu fayldan yoki `main.go`'dan nusxalang, aks holda test «boshqa matn» deydi.
- Ayirish uchun `stock[name] -= n`: `-=` — `+=`ning ayirish shakli (`x -= n` = `x = x - n`).

Javob: [`solution/main.go`](solution/main.go).
