# 06 · Uy vazifasi javoblari

Avval oʻzingiz urinib koʻring. 🙂 Hammasi `go1.26.8` da ishga tushirilgan.

## 🟢 Oson
[`oson/main.go`](oson/main.go). Eng kichik, eng katta va xato: `(int, int, error)`. `min`/`max` Go 1.21 dan beri ichki funksiyalar.

```
-3 12 <nil>
0 0 boʻsh slice
```

`fmt.Println(minMax(...))`: bir nechta qiymat qaytaradigan chaqiruvni toʻgʻridan-toʻgʻri boshqa funksiyaga argument qilib
berish mumkin (spec: *Calls*, "f(g(parameters_of_g))").

## 🟡 Oʻrta
[`orta/main.go`](orta/main.go). `birlashtir(ajratgich string, sozlar ...string)`: variadic parametr doim **oxirgi** boʻladi.

```
olma, nok, uzum
""
anor + behi
anor + behi
```

Argumentsiz chaqiruvda `sozlar` = `nil`, natija boʻsh string. Standart kutubxonada tayyori bor: `strings.Join`.

## 🟠 Qiyin
[`qiyin/main.go`](qiyin/main.go). `logla` funksiya oladi va uni oʻraydigan **yangi funksiya qaytaradi**: funksiya qiymat,
closure va qaytarilgan funksiya bir joyda. Closure `nom`, `next` va `n` ni ushlab qoladi, har `logla` chaqiruvida yangi `n`:

```
[salom #1] "Aziz" → "Salom, Aziz"
[salom #2] "Nodira" → "Salom, Nodira"
[ichki #1] "Gʻayrat" → "Salom, Gʻayrat"
[tashqi #1] "Gʻayrat" → "Salom, Gʻayrat"
```

HTTP middleware'lar (backend treki) aynan shunday quriladi: `func(next http.Handler) http.Handler`.

## 🔴 Guru
Vazifa: [`guru/xato/main.go`](guru/xato/main.go) (Temporal xatosining soddalashtirilgan shakli, 2021):

```
$ go run ./homework/guru/xato
saqla: <nil>
$ go vet -vettool=$(go env GOPATH)/bin/shadow ./homework/guru/xato
homework/guru/xato/main.go:17:9: declaration of "err" shadows declaration at line 15
```

`blob, err := kodla(...)` if ichida **yangi** `err` yaratadi, keyingi `err = yoz(blob)` esa aynan shu ichki `err` ga
yozadi. Tashqi `err` `nil` boʻlib qoladi.

Javob: [`guru/main.go`](guru/main.go). Ichki blokda `:=` ishlatmaymiz: `var blob string` + `blob, err = kodla(...)`.

```
saqla: baza: yozib boʻlmadi
```

`shadow` endi jim.
