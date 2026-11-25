// join: bir nechta xatoni bitta error'ga yigʻish (Go 1.20).
// go run ./join
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
	var summaXato, manzilXato error // turi error: typed nil boʻlmaydi
	if b.Summa <= 0 {
		summaXato = &SummaXato{Summa: b.Summa}
	}
	if b.Manzil == "" {
		manzilXato = ErrManzilYoq
	}
	return errors.Join(summaXato, manzilXato)
}

func main() {
	err := tekshir(Buyurtma{Summa: -5})
	fmt.Println(err)
	fmt.Println("manzil xatosi bormi:", errors.Is(err, ErrManzilYoq))
	if sx, ok := errors.AsType[*SummaXato](err); ok {
		fmt.Println("summa xatosi:", sx.Summa)
	}
}
