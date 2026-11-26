# Vazifa · GB-1.08 · `parseSold` uchun jadvalli test va yuqori chegara

Videoda `parsePrice` uchun jadvalli test yozdik va u ikki xatoni topdi: narxning yuqori chegarasi faqat
`main`da, koʻpaytirishdan keyin tekshirilardi, juda katta narx esa `int64`dan oshib manfiy boʻlardi.
`parseSold` (sotilayotgan dona soni) ham kirishni faqat pastdan tekshiradi: kamida 1 dona. Hozir
`./bozor add olma 12000 1000000` ni `main`dagi ombor tekshiruvi (`omborda 40 dona bor`) toʻxtatadi, lekin bu
tasodif: `parseSold` kirishni oʻqiydigan funksiya, uni boshqa buyruqlar ham chaqirishi mumkin, ombor tekshiruvi esa faqat
`add`da. Videodagi qoida shu: qoida kirish oʻqiladigan joyda turadi. Bitta buyurtmada **1000 donadan koʻp
boʻlmasin**, buni `parseSold`ning oʻzi rad etsin.

Ishni videodagi tartibda qiling: avval xatoni koʻrsatadigan test, keyin tuzatish.

## 1-qadam: test yozing (u yiqilsin)

Shu papkada `sold_test.go` faylini yarating. Unda `func TestParseSold(t *testing.T)` jadvalli test boʻlsin
(videodagi `price_test.go` kabi: `name`, `in`, `want`, `wantErr` maydonli `tests` jadvali,
`for _, tc := range tests`, har qator `t.Run` bilan). Kamida shu qatorlar:

| name | in | want | wantErr |
|---|---|---|---|
| `oddiy` | `"5"` | `5` | `false` |
| `chegara` | `"1000"` | `1000` | `false` |
| `nol` | `"0"` | `0` | `true` |
| `ortiq` | `"1001"` | `0` | `true` |
| `harf` | `"besh"` | `0` | `true` |

`go test -v -run TestParseSold` ni ishga tushiring: `ortiq` subtest'i **yiqilishi kerak**. Yiqilmasa, testingiz
bu xatoni koʻrmayapti: jadvalni tekshiring.

## 2-qadam: `parseSold`ni tuzating

`main.go`da `const maxSold = 1000` qoʻshing va `parseSold`da miqdor undan katta boʻlsa, xato qaytaring. Xato
matni aynan shunday boʻlsin (foydalanuvchi uni terminalda oʻqiydi):

```
miqdor 1000 donadan oshmasin, keldi 1001
```

## Tekshirish

Shu papkada:

```bash
go test
```

`ok` chiqsa, vazifa bajarildi. `go test` uchta test faylini birga ishga tushiradi: videodagi `price_test.go`,
sizning `sold_test.go` va `tekshiruv_test.go`. Oxirgisi vazifani tekshiradi, uni oʻzgartirmang. Uning
`TestVazifaParseSold` qismi videodagi jadval bilan bir xil tuzilgan; boshqa ikki funksiyasidagi
`os.ReadFile` va `strings.Contains` izohlarda tushuntirilgan. (`go test ./...` esa `solution/` papkani ham
tekshiradi; u yerda javob turgani uchun `ok` chiqadi.) Hozir (vazifadan oldin) natija shunday (`ortiq`, `million` — tekshiruvchining oʻz jadvalidagi qatorlar):

```
--- FAIL: TestVazifaJadvalBor (0.00s)
    tekshiruv_test.go:18: sold_test.go oʻqilmadi (open sold_test.go: no such file or directory): avval parseSold uchun jadvalli test yozing (TASK.md, 1-qadam)
--- FAIL: TestVazifaParseSold (0.00s)
    --- FAIL: TestVazifaParseSold/ortiq (0.00s)
        tekshiruv_test.go:48: parseSold("1001") = 1001, <nil>; xato kutilgan: true (TASK.md, 2-qadam)
    --- FAIL: TestVazifaParseSold/million (0.00s)
        tekshiruv_test.go:48: parseSold("1000000") = 1000000, <nil>; xato kutilgan: true (TASK.md, 2-qadam)
--- FAIL: TestVazifaXatoMatni (0.00s)
    tekshiruv_test.go:61: parseSold("1001") xato qaytarishi kerak, hozir xato yoʻq (TASK.md, 2-qadam)
FAIL
exit status 1
FAIL	github.com/abdurahmon27/bozor	0.002s
```

## Maslahatlar

1. Shu papkadagi `price_test.go`ni `sold_test.go` nomi bilan nusxalab boshlang: `parsePrice` → `parseSold`, `want`ning turi `int64` emas, `int`
   (`parseSold` `int` qaytaradi; turlar mos kelmasa, kompilyator GB-1.03'dagi `mismatched types` xatosini beradi).
2. Chegara qiymatlarini albatta qoʻshing: `1000` oʻtishi, `1001` yiqilishi kerak. Xatolar koʻpincha aynan chegarada.
3. Xato matnini `fmt.Errorf` bilan yasang, `parsePrice`dagi `maxPriceSom` tekshiruvi kabi: `maxSold` va kelgan son, ikkalasi `%d` bilan.

Javob: [`solution/`](solution/) (`sold_test.go` va `main.go`). Avval oʻzingiz urinib koʻring.
