# Vazifa · GB-6.03: argon2id parametrlarini oshiring va eski hash'larni toping

Videoda `internal/password` OWASP'ning minimal sozlamasi bilan ishladi: `m=19456` KiB (19 MiB), `t=2`, `p=1`. Bazada
shunday satrlar allaqachon bor. Endi server xotirasi koʻp: OWASP jadvalining **birinchi qatoriga** oʻtamiz —
`m=47104` (46 MiB), `t=1`, `p=1` (himoya darajasi bir xil, faqat xotira koʻproq, vaqt kamroq).

Muammo: eski satrlar yoʻqolmaydi. Ular **tekshirilishi** kerak (Verify parametrlarni satrdan oʻqiydi, shuning uchun
bu oʻzi ishlaydi), lekin login paytida parol toʻgʻri chiqsa, ularni **yangi parametrlar bilan qayta saqlash** kerak.
Buning uchun bitta funksiya: `NeedsRehash`.

**Nima qilish kerak** (`internal/password/password.go`):

1. **Parametrlar.** `const` blokida `memory`ni `46 * 1024` (47104 KiB), `passes`ni `1` qiling. `threads`, `keyLen`,
   `saltLen` oʻzgarmaydi.
2. **`NeedsRehash(encoded string) (bool, error)`.** Satrni `decode` bilan oʻqing. Xato boʻlsa — `false` va oʻsha xato
   (u `ErrBadHash`ni oʻraydi). Aks holda `memory`, `passes`, `threads` yoki kalit uzunligi (`len(p.key)`) hozirgi
   konstantalardan **farq qilsa** — `true`, hammasi bir xil boʻlsa — `false`.

`tekshiruv_test.go`ni oʻzgartirmang. U uch narsani tekshiradi: yangi `Hash` `m=47104,t=1,p=1` bilan yasaydi; eski satr
(`m=19456,t=2`) hamon `Verify`dan oʻtadi; `NeedsRehash` eski satrga `true`, yangisiga `false`, buzuq satrga `ErrBadHash`
qaytaradi.

**Tekshirish** (papkaning ildizida):

```bash
go test ./...
```

Hammasi `ok` boʻlguncha ishlang. Hozir shunday yiqiladi (qisqartirilgan):

```
--- FAIL: TestVazifaParams (0.02s)
    tekshiruv_test.go:23: Hash hozir m=19456,t=2,p=1 bilan yasayapti, kutilgan m=47104,t=1,p=1 (OWASP'ning birinchi qatori, 46 MiB). password.go'dagi memory va passes konstantalarini oʻzgartiring (TASK.md, 1-qadam)
--- FAIL: TestVazifaNeedsRehash (0.02s)
    tekshiruv_test.go:38: NeedsRehash(eski satr m=19456,t=2) = false, <nil>; kutilgan true, nil. decode bilan oʻqing va parametrlarni konstantalar bilan solishtiring (TASK.md, 2-qadam)
    tekshiruv_test.go:45: NeedsRehash(sha256 hex) xatosi <nil>; kutilgan ErrBadHash: buzuq satrda decode'ning xatosini qaytaring
FAIL
FAIL	github.com/abdurahmon27/bozor/internal/password
```

**Maslahatlar**

- `decode` `params` qaytaradi: `p.memory`, `p.passes`, `p.threads`, `p.key`. Konstantalar shu paketda: `memory`, `passes`,
  `threads`, `keyLen`. Bir nechta shartni `||` bilan bitta `return`ga yigʻish mumkin (GB-1.05).
- Faqat 1-qadamni qilsangiz, `TestVazifaNeedsRehash` yiqilaveradi; faqat 2-qadamni qilsangiz, eski satr hozirgi
  parametrlar bilan bir xil boʻlib qoladi va `NeedsRehash` `false` qaytaradi.
- `go test ./internal/password -run TestVazifa -v` faqat vazifa testlarini koʻrsatadi.

Login'ning oʻzi (parol toʻgʻri + `NeedsRehash` → yangi `Hash`ni saqlash) — GB-6.04'da.

Javob: `solution/internal/password/password.go` (avval oʻzingiz urinib koʻring).
