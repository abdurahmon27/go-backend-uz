// asxato: errors.As ga & ni unutish. ATAYLAB XATO.
// Kompilyatsiya boʻladi, lekin ishlaganda panic. go vet buni topadi.
// go run ./asxato ; go vet ./asxato
package main

import (
	"errors"
	"fmt"
)

type SummaXato struct{ Summa int64 }

func (e *SummaXato) Error() string { return fmt.Sprintf("summa notoʻgʻri: %d", e.Summa) }

func main() {
	err := fmt.Errorf("buyurtma 9: %w", &SummaXato{Summa: -5})
	var sx *SummaXato
	if errors.As(err, sx) { // & yoʻq
		fmt.Println("summa =", sx.Summa)
	}
}
