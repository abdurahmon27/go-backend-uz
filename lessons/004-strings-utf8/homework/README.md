# 04 · Uy vazifasi javoblari

Avval oʻzingiz urinib koʻring. 🙂 Hammasi `go1.26.8` da ishga tushirilgan.

## 🟢 Oson
[`oson/main.go`](oson/main.go). `len` baytlarni, `utf8.RuneCountInString` rune'larni sanaydi:

```
Aziz       bayt=4 rune=4
Gʻayrat    bayt=8 rune=7
Oʻgʻiloy   bayt=10 rune=8
```

Har bir `ʻ` (U+02BB) bitta rune, lekin ikki bayt.

## 🟡 Oʻrta
[`orta/main.go`](orta/main.go). Rune'lar ustida yurib, `oʻ gʻ sh ch` juftligini bitta harf deb olamiz, tutuq belgisi
`ʼ` esa harf emas (CLDR `uz.xml` izohi: "the symbol ʼ [is] not considered [a] separate letter"):

```
oʻzbek  rune=6 harf=5 ["oʻ" "z" "b" "e" "k"]
shahar  rune=6 harf=5 ["sh" "a" "h" "a" "r"]
choʻgʻ  rune=6 harf=3 ["ch" "oʻ" "gʻ"]
maʼno   rune=5 harf=4 ["m" "a" "n" "o"]
Isʼhoq  rune=6 harf=5 ["I" "s" "h" "o" "q"]
```

`ng` haqida: 1995-yilgi imloda u alohida harf birikmasi sifatida qaralgan, CLDR esa uni alohida harf deb hisoblamaydi,
2026-yilgi qonun alifbosida esa u yoʻq. Shuning uchun bu yerda `ng` ni ikki harf deb oldik.

Eʼtibor bering: `Isʼhoq` da tutuq belgisi `s` va `h` ni ajratib turadi, shuning uchun bu yerda `sh` harfi yoʻq.
Bu ishlashi uchun matn avval darsdagi `tozala()` dan oʻtgan boʻlishi kerak.

## 🟠 Qiyin
[`qiyin/main.go`](qiyin/main.go):

```
kebzʻo
ʻgob
```

Rune'lar teskari boʻldi, lekin `ʻ` oʻz harfidan ajraldi: endi u `z` ga yopishgan. Sabab: `oʻ` bitta harf, lekin ikki
rune. "Harf" va "rune" bir xil narsa emas. Bu LC-01 (Valid Palindrome) masalasida yana chiqadi.

## 🔴 Guru
[`guru/main.go`](guru/main.go). Har bir soʻzni alifbodagi oʻrinlar roʻyxatiga aylantiramiz va `slices.Compare` bilan
solishtiramiz:

```
bayt:   [choy chumoli gilos gʻoz olma oʻrik pomidor sabzi shaftoli zarang]
alifbo: [gilos olma pomidor sabzi zarang oʻrik gʻoz shaftoli choy chumoli]
```

Alifbo: `a b d e f g h i j k l m n o p q r s t u v x y z oʻ gʻ sh ch` (CLDR `common/collation/uz.xml`).
`golang.org/x/text/collate` ni sinab koʻrdik: v0.42.0 (CLDR 23 jadvallari) `language.Uzbek` bilan ham
`[gʻoz gilos … zarang]` beradi, yaʼni bu qoidani qoʻllamaydi.

Murakkablashtirish: 2027-yildan alifbo `A B D E F G Ğ H … O Ö P … S Ş Ç T … Z` boʻladi (OʻRQ-1184). `alifbo`
oʻzgaruvchisini almashtirish yetarlimi? (Ishora: `ö` ning ikki xil yozilishi, NFC.)
