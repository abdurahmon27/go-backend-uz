// Guru: join/ dagi tekshir'da summaXato ni *SummaXato turida eʼlon qilsak nima boʻladi?
// Toʻgʻri buyurtmada ham err != nil: Join "nil error"larni err != nil bilan tashlaydi,
// typed nil esa nil emas (GO-1.13). Println Error() ni chaqiradi, u nil pointer'dan
// oʻqiydi va panic qiladi; fmt uni ushlab %!v(PANIC=...) deb yozadi.
// go run ./homework/guru
package main

import (
	"errors"
	"fmt"
)

type Buyurtma struct {
	Summa  int64
	Manzil string
}

type SummaXato struct{ Summa int64 }

func (e *SummaXato) Error() string { return fmt.Sprintf("summa notoʻgʻri: %d", e.Summa) }

var ErrManzilYoq = errors.New("manzil boʻsh")

func tekshir(b Buyurtma) error {
	var summaXato *SummaXato // XATO: aniq tur, error emas
	var manzilXato error
	if b.Summa <= 0 {
		summaXato = &SummaXato{Summa: b.Summa}
	}
	if b.Manzil == "" {
		manzilXato = ErrManzilYoq
	}
	return errors.Join(summaXato, manzilXato)
}

func main() {
	err := tekshir(Buyurtma{Summa: 150000, Manzil: "Toshkent"}) // toʻgʻri buyurtma
	fmt.Println("err == nil:", err == nil)
	fmt.Println("err:", err)
}
