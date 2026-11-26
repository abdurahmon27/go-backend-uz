# Vazifa · GB-1.01 · oʻz bozor modulingiz, versiya qatori bilan

Videoda koʻrdingiz: serverda eski fayl qolib ketishi mumkin. Shuning uchun dasturlar odatda ishga tushganda
oʻz versiyasini yozadi. Vazifa ikki qismli, ikkalasini `go test` tekshiradi.

## 1. Videodagi yoʻlni oʻzingiz bosib oʻting

Boʻsh joydan boshlang (bu kurs papkasidan tashqarida):

1. `mkdir bozor && cd bozor`
2. `go mod init github.com/<siz>/bozor`: `<siz>` oʻrniga oʻz GitHub nomingiz. `go.mod`'ni oching: ikki yozuv bor,
   `module …` va `go …` (sizdagi Go versiyasi).
3. `main.go`'ni videodagidek yozing (7 qator) va `go run .` bilan ishga tushiring.
4. Shu kurs papkasidagi `main_test.go` faylini oʻz `bozor` papkangizga nusxalang.

Test `go.mod`'ni ham tekshiradi: modul yoʻli `/bozor` bilan tugashi va `go` yozuvi boʻlishi kerak.
(Shu kurs papkasining oʻzida ham ishlashingiz mumkin: bu yerda `go.mod` tayyor.)

## 2. Versiya qatorini qoʻshing

`main.go`'ni oʻzgartiring: dastur ekranga aynan shu **ikki qatorni** chiqarsin.

```
bozor ishga tushdi
versiya: 0.1.0
```

Birinchi qator oʻzgarmaydi. Ikkinchi qatorda `versiya`, ikki nuqta, bitta boʻsh joy va `0.1.0` (bozorning birinchi
versiyasi). Fayl `go fmt` formatida boʻlishi kerak: test buni ham tekshiradi.

## Tekshirish

Papkangizda:

```bash
go test
```

`ok` chiqsa, vazifa bajarildi. `FAIL` chiqsa, `main_test.go:<qator>:` dan keyingi xabarni oʻqing: u nima kerakligini
va hozir nima ekanini aytadi. Versiya qatoridan oldin natija shunday:

```
--- FAIL: TestIkkiQator (0.00s)
    main_test.go:53: 2-qator «versiya: 0.1.0» boʻlishi kerak, hozir dastur faqat 1 qator chiqaryapti (TASK.md)
FAIL
exit status 1
FAIL	github.com/<siz>/bozor	0.002s
```

`main_test.go` faylini oʻzgartirmang: test yozishni 8-darsda (GB-1.08) oʻrganamiz. Hozir bilishingiz kerak
boʻlgani: `go test` shu faylni ishga tushiradi, u esa `main` funksiyangiz ekranga nima yozganini, `main.go`
formatini va `go.mod`'ni tekshiradi.

## Maslahatlar

1. `main` funksiyasi tanasidagi qatorlar tepadan pastga bajariladi: ikkinchi qator uchun ikkinchi buyruq kerak.
2. Har bir `fmt.Println` chiqargan matn oxiriga yangi qatorni oʻzi qoʻshadi.
3. Surilishni qoʻlda toʻgʻrilamang: `go fmt ./...` uni oʻzi qoʻyadi.
4. Xato chiqsa, avval `fayl:qator:ustun`'ni toping, keyin xabarni oʻqing; bir nechta boʻlsa, birinchisidan boshlang.

Qoʻshimcha (tekshirilmaydi): `go build`, keyin `./bozor`. Endi `main.go`'dagi versiyani `0.1.1` qiling va yana
`./bozor`'ni ishga tushiring, `go build`'siz. Qaysi versiya chiqdi va nega? Keyin `0.1.0`'ga qaytaring.

Javob: [`solution/main.go`](solution/main.go). Avval oʻzingiz urinib koʻring.
