// qiyin (javob): defer receiver'ni qachon oladi?
// defer b.Chiqar()            — method value: receiver (b ning nusxasi) defer qatorida olinadi → 100
// defer func() { b.Chiqar() }() — closure: b funksiya oxirida, chaqiruv paytida oʻqiladi → 200
// Tartib: defer LIFO, shuning uchun avval closure (200), keyin method value (100).
package main

import "fmt"

type Buyurtma struct{ Summa int64 }

func (b Buyurtma) Chiqar() { fmt.Println(b.Summa) }

func main() {
	b := Buyurtma{Summa: 100}
	defer b.Chiqar()
	defer func() { b.Chiqar() }()
	b.Summa = 200
}
