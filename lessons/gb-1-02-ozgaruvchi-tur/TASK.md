# Vazifa · GB-1.02 · sotuvdan keyin qoldiq

Bugun bozorda **3 qadoq olma sotildi**. `main.go`'ning oxiriga (oxirgi `fmt.Println` qatoridan keyin,
yopuvchi `}` dan oldin) shuni qoʻshing:

1. `sold` degan oʻzgaruvchi: nechta qadoq sotilgani, qiymati `3`. Qisqa shakl bilan eʼlon qiling.
2. `stock`'ni `sold` qadar kamaytiring. `stock` allaqachon bor: yangi oʻzgaruvchi emas, bor oʻzgaruvchiga
   yangi qiymat yoziladi.
3. Mahsulot nomi va yangi qoldiqni bitta `fmt.Println` bilan chiqaring.

Dastur ekranga aynan shu besh qatorni chiqarishi kerak (oxirgisi yangi):

```
bozor ishga tushdi
versiya: 0.1.0
 0 false
olma 12000 true 40 1.5
olma 37
```

## Tekshirish

```bash
go test
```

`ok` chiqsa, vazifa bajarildi. `FAIL` chiqsa, xabarni oʻqing: u qaysi qator nima boʻlishi kerakligini aytadi.
`main_test.go` faylini oʻzgartirmang: uni 8-darsda (GB-1.08) qatorma-qator ochamiz.

## Maslahatlar

- `stock := stock - sold` yozsangiz, kompilyator `no new variables on left side of :=` deydi. Videodagi oraliq savolni
  eslang: bor oʻzgaruvchiga bitta `=` bilan yoziladi.
- `sold`'ni eʼlon qilib, hech qayerda ishlatmasangiz, Go dasturni qurmaydi: `declared and not used: sold`.
  Eʼlon qilingan har bir oʻzgaruvchi ishlatilishi kerak.
- 4-qator oʻzgarib qolsa (`olma 12000 true 37 1.5`), `stock`'ni juda erta kamaytirgansiz: yangi qatorlarni
  mavjud `fmt.Println(name, price, active, stock, weight)` qatoridan **keyin** yozing.

Javob: [`solution/main.go`](solution/main.go). Avval oʻzingiz urinib koʻring.
