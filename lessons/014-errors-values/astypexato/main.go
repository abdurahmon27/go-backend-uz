// astypexato: AsType ga notoʻgʻri tur. ATAYLAB KOMPILYATSIYA BOʻLMAYDI.
// Error pointer receiver bilan yozilgan, shuning uchun SummaXato (pointer'siz) error emas.
// go build ./astypexato
package main

import (
	"errors"
	"fmt"
)

type SummaXato struct{ Summa int64 }

func (e *SummaXato) Error() string { return fmt.Sprintf("summa notoʻgʻri: %d", e.Summa) }

func main() {
	err := fmt.Errorf("buyurtma 9: %w", &SummaXato{Summa: -5})
	if sx, ok := errors.AsType[SummaXato](err); ok { // * yoʻq
		fmt.Println("summa =", sx.Summa)
	}
}
