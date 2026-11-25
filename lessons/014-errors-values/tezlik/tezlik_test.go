// tezlik: errors.As va errors.AsType ni oʻlchaymiz.
// go test -bench . -benchmem ./tezlik
package tezlik

import (
	"errors"
	"fmt"
	"testing"
)

type SummaXato struct{ Summa int64 }

func (e *SummaXato) Error() string { return fmt.Sprintf("summa notoʻgʻri: %d", e.Summa) }

// ikki qavat oʻralgan xato: toʻlov → buyurtma → *SummaXato
var err = fmt.Errorf("toʻlov: %w", fmt.Errorf("buyurtma 9: %w", &SummaXato{Summa: -5}))

func BenchmarkAs(b *testing.B) {
	for b.Loop() {
		var sx *SummaXato
		if !errors.As(err, &sx) {
			b.Fatal("topilmadi")
		}
	}
}

func BenchmarkAsType(b *testing.B) {
	for b.Loop() {
		if _, ok := errors.AsType[*SummaXato](err); !ok {
			b.Fatal("topilmadi")
		}
	}
}
