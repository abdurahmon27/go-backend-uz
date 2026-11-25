// fix: xato yoʻq boʻlsa, oddiy nil qaytaramiz.
// go run ./fix
package main

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

type SummaXato struct{ Summa int64 }

func (e *SummaXato) Error() string { return fmt.Sprintf("summa notoʻgʻri: %d", e.Summa) }

func tekshir(b Buyurtma) error {
	if b.Summa <= 0 {
		return &SummaXato{Summa: b.Summa}
	}
	return nil
}

func main() {
	for _, b := range []Buyurtma{{Summa: 150000}, {Summa: -5}} {
		err := tekshir(b)
		fmt.Printf("summa %d → err == nil: %v, err: %v\n", b.Summa, err == nil, err)
	}
}
