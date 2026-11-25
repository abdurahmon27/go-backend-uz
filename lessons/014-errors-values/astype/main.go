// astype: oʻralgan xatodan oʻz turimizni olish. errors.As (Go 1.13) va errors.AsType (Go 1.26).
// go run ./astype
package main

import (
	"errors"
	"fmt"
)

type Buyurtma struct {
	ID    int
	Summa int64
}

type SummaXato struct{ Summa int64 }

func (e *SummaXato) Error() string { return fmt.Sprintf("summa notoʻgʻri: %d", e.Summa) }

func saqla(b Buyurtma) error {
	if b.Summa <= 0 {
		return fmt.Errorf("buyurtma %d: %w", b.ID, &SummaXato{Summa: b.Summa})
	}
	return nil
}

func main() {
	err := saqla(Buyurtma{ID: 9, Summa: -5})
	fmt.Println("err:", err)

	// eski usul (Go 1.13): avval oʻzgaruvchi, keyin uning manzili
	var sx *SummaXato
	if errors.As(err, &sx) {
		fmt.Println("As:     summa =", sx.Summa)
	}

	// yangi usul (Go 1.26): tur qavs ichida, natija qaytadi
	if sx, ok := errors.AsType[*SummaXato](err); ok {
		fmt.Println("AsType: summa =", sx.Summa)
	}
}
