// savol: method value receiver'ni qachon oladi? Pauzadagi savol.
package main

import "fmt"

type Buyurtma struct {
	Summa int64
}

func (b Buyurtma) Jami() int64   { return b.Summa }
func (b *Buyurtma) JamiP() int64 { return b.Summa }

func main() {
	b := Buyurtma{Summa: 100}
	f := b.Jami  // value receiver: b ning nusxasi shu yerda olinadi
	g := b.JamiP // pointer receiver: &b shu yerda olinadi
	b.Summa = 200
	fmt.Println("f():", f())
	fmt.Println("g():", g())
}
