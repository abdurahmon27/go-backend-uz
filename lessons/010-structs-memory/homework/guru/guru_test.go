// guru: false sharing. Ikki goroutine har biri oʻz hisoblagichini oshiradi.
// Yonma-yon turgan hisoblagichlar bitta cache line'da (64 bayt) boʻladi.
package guru

import (
	"sync"
	"sync/atomic"
	"testing"
)

type Yonma struct {
	A atomic.Int64
	B atomic.Int64 // A bilan bitta cache line'da
}

type Ajratilgan struct {
	A atomic.Int64
	_ [56]byte // A va B orasida 64 bayt boʻlsin
	B atomic.Int64
}

const n = 1_000_000

func ikkiGoroutine(a, b *atomic.Int64) {
	var wg sync.WaitGroup
	wg.Go(func() {
		for range n {
			a.Add(1)
		}
	})
	wg.Go(func() {
		for range n {
			b.Add(1)
		}
	})
	wg.Wait()
}

func BenchmarkYonma(b *testing.B) {
	var h Yonma
	for b.Loop() {
		ikkiGoroutine(&h.A, &h.B)
	}
}

func BenchmarkAjratilgan(b *testing.B) {
	var h Ajratilgan
	for b.Loop() {
		ikkiGoroutine(&h.A, &h.B)
	}
}
