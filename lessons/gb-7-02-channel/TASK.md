# Vazifa · GB-7.02: `unique` bosqichi

Videoning oxirida ikki yetkazib beruvchi (`chorsu`, `oloy`) bir vaqtda `raw` channel'iga yozdi, `inStock` qoldigʻi borlarni
oʻtkazdi va **non ikki marta keldi**: u ikkala roʻyxatda ham bor (ID 3). `main.go`da shu pipeline `importAll` funksiyasida
turibdi va natijani slice qilib qaytaradi. Endi takror ID'larni tashlab yuboradigan bosqich kerak.

**Nima qilish kerak** (`main.go`, `VAZIFA` deb belgilangan joylar)

1. **`unique`ni yozing.** `in`dan hamma mahsulotni oʻqisin (`for p := range in`), har ID'ni faqat **birinchi marta**
   `out`ga yuborsin, keyingilarini tashlab ketsin; tartib oʻzgarmaydi. `in` yopilgach, `out`ni yopsin
   (`defer close(out)`, videodagi `inStock`dek). Koʻrilgan ID'lar uchun `map[int64]bool` (GB-2.02).
2. **`importAll`ga qoʻying.** `inStock` va natija orasiga `unique`ni qoʻying: yangi channel yarating,
   `inStock` unga yozsin, `unique` undan oʻqib `ready`ga yozsin. Har bosqich oʻz goroutine'ida (`go …`).

`product.go`, `data.go`, `parse`, `inStock` va `tekshiruv_test.go`ni oʻzgartirmang.

**Tekshirish** (shu papkada):

```bash
go test .
```

`ok` boʻlguncha ishlang. Tekshiruv har bosqichni 2 soniyalik vaqt chegarasi bilan kutadi: bosqich `out`ni yopmasa,
`for p := range out` abadiy kutadi va test oʻzbekcha xabar bilan yiqiladi (vaqt chegarasi `select` va `time.After` bilan
yozilgan; ularni GB-7.03'da oʻrganamiz, testni oʻqish shart emas). Hozir shunday yiqiladi:

```
--- FAIL: TestVazifaUnique (2.00s)
    tekshiruv_test.go:54: unique 2 soniyada tugamadi: out yopilmadi, shuning uchun `for p := range out` abadiy kutyapti.
        	unique in'dagi hamma qiymatni oʻqishi (for p := range in) va oxirida out'ni yopishi kerak (defer close(out)).
--- FAIL: TestVazifaImportAll (0.00s)
    tekshiruv_test.go:101: importAll ichida unique bosqichi yoʻq: TASK.md 2-qadam, inStock → yangi channel → go unique(…, ready)
    tekshiruv_test.go:120: importAll natijasida ID 3 2 marta (1-urinish), kutilgan 1 marta. Natija: [{1 olma 1200000 40} {3 non 400000 100} {4 choy 2500000 8} {3 non 420000 60} {5 anor 1800000 25}]
        	TASK.md 2-qadam: inStock va natija orasiga unique bosqichini qoʻying.
FAIL
```

`Natija`dagi har `{1 olma 1200000 40}` — bitta `Product`, maydonlari tartibida: `ID`, `Name`, `Price` (tiyin), `Stock`
(`%v` struct'ni shunday chiqaradi). Tartib sizda boshqacha boʻladi: goroutine'lar har safar boshqa tartibda ishlaydi.
Tekshiruv `importAll` ichida `unique` chaqirilganini ham koʻradi (natijani oxirgi siklda tozalash vazifa emas).

**Maslahatlar**

- `unique`ning shakli `inStock` bilan bir xil: `defer close(out)`, `for p := range in`, ichida bitta `if`.
  Farqi — shart `p.Stock > 0` emas, "bu ID hali koʻrilmaganmi".
- `seen := make(map[int64]bool)`; ID koʻrilgan boʻlsa `continue`, aks holda `seen[p.ID] = true` va `out <- p`.
- `importAll` qotib qolsa (test "2 soniyada tugamadi" desa): `unique`ga `ready`ning oʻzini emas, yangi channel'ni
  `in` qilib berganingizni va `inStock` endi `ready`ga emas, yangi channel'ga yozayotganini tekshiring.

Javob: [`solution/main.go`](solution/main.go). Javobni tekshirish: `go test ./solution/`.
