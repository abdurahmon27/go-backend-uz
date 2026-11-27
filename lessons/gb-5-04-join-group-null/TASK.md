# Vazifa · GB-5.04: nima qancha sotilgan va qaysi buyurtmada izoh yoʻq

Videoda har **buyurtma**ning jamini hisobladik (`orders` + `order_items`). Bu yerda xuddi shu usulni boshqa tomondan qoʻllaysiz:
har **mahsulot** uchun nechta dona sotilgan. Ikkala `.sql` faylda hozir faqat izoh bor: `SELECT`ni oʻzingiz yozasiz
(odatdagidek `;` bilan tugating, izohni olib tashlang). Jadvallarga (`products.sql`, `orders.sql`) va `seed.sql`ga tegmang.

**Nima qilish kerak**

1. **`sotilgan.sql`.** Bitta `SELECT`: **har bir** mahsulot uchun `name` va `sotilgan` ustunlari, `sotilgan` = shu mahsulotning
   barcha buyurtmalardagi `quantity` yigʻindisi. Hech sotilmagan mahsulot ham chiqsin, `sotilgan` ustunida **0** bilan
   (`NULL` emas). Tartib: `sotilgan` **kamayish** boʻyicha, teng boʻlsa `name` **oʻsish** boʻyicha.
2. **`izohsiz.sql`.** Bitta `SELECT`: izohi **yoʻq** (`note` ustuni `NULL`) buyurtmalarning `id` va `customer` ustunlari, `id`
   oʻsish tartibida. Diqqat: boʻsh matn `''` — `NULL` emas, bunday buyurtma chiqmasligi kerak.

**Tekshirish** (shu papkada):

```bash
docker compose up -d --wait
go test .
```

Hozir shunday yiqiladi (qisqartirilgan; `.sql` fayllarda hali faqat izoh bor):

```
--- FAIL: TestSotilgan (0.00s)
    tekshiruv_test.go:132: sotilgan.sql'da SELECT yoʻq: har mahsulot uchun name va sotilgan (order_items.quantity yigʻindisi) ni tanlang; hech sotilmaganda 0 chiqsin (TASK.md)
--- FAIL: TestIzohsiz (0.00s)
    tekshiruv_test.go:143: izohsiz.sql'da SELECT yoʻq: izohi yoʻq (NULL) buyurtmalarning id va customer ustunlarini tanlang (TASK.md)
FAIL
```

Tekshiruv `tekshiruv_test.go`da: u `products.sql`, `orders.sql` va `seed.sql`ni alohida vaqtinchalik sxemada bajaradi
(sizning `public` sxemangizga tegmaydi), soʻngra sizning `.sql` faylingizni **ikki xil maʼlumotda** ishga tushiradi
(ikkinchisida boshqa mahsulotlar va buyurtmalar bor; har safar Postgres rejasini ham oʻzgartiradi), shuning uchun natijani qoʻlda
yozib qoʻyish va `ORDER BY`ni tashlab ketish oʻtmaydi. Natijani test `psql -At` shaklida solishtiradi: sarlavhasiz, ustunlar `|` bilan
(`non|3`) — bu `psql`ning odatdagi jadvalidan farq qiladi, xato emas. Oxirida sxemani
oʻchiradi. Go dasturidan Postgres'ga ulanishni GB-5.05'da koʻramiz. Docker ishlamasa test aniq xabar bilan yiqiladi,
oʻtkazib yuborilmaydi.

Tugagach: `docker compose down -v` (konteyner va uning maʼlumotlarini oʻchiradi).

**Maslahatlar**

- Sotilgan ustun mahsulotdan boshlanadi: `FROM products LEFT JOIN order_items ON order_items.product_id = products.id`.
  Chap jadval — hamma qatori qoladigan jadval.
- Mahsulot boʻyicha guruhlang: `GROUP BY products.id`; yigʻindi — `SUM(order_items.quantity)`.
- Sotilmagan mahsulotda `SUM` nima qaytaradi? Natijani `COALESCE(…, 0)` bilan almashtiring va ustunga `AS sotilgan` deb nom bering.
- Tartib: `ORDER BY sotilgan DESC, products.name` (`DESC` — kamayish; ikkinchi shart teng qatorlar uchun).
- Ikkita tuzoq: (1) `COALESCE`siz sotilmagan mahsulotning `sotilgan`i `NULL` boʻladi va `DESC` tartibda eng birinchi chiqadi;
  (2) `LEFT JOIN`dan keyin oʻng jadval ustuniga `WHERE` yozsangiz (masalan `WHERE order_items.quantity > 0`), juftsiz qatorlar yana
  tushib ketadi — shartni `ON`ga yozing yoki umuman yozmang.
- Ikkinchi vazifa: `NULL` bilan `=` ishlamaydi, `IS NULL` yozing.

Javob: [`solution/`](solution/). Tekshirish: `cd solution && docker compose up -d --wait && go test .`
