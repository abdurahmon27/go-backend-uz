package kesh

import (
	"fmt"
	"testing"
)

// Pauza savoli: har elementni yoki har 8-elementni (64 baytda bitta) qo'shamiz.
func BenchmarkQadam(b *testing.B) {
	for _, n := range []int{1 << 12, 1 << 24} {
		s := Sonlar(n)
		for _, k := range []int{1, 8} {
			b.Run(fmt.Sprintf("n=%d/k=%d", n, k), func(b *testing.B) {
				for b.Loop() {
					Qadam(s, k)
				}
			})
		}
	}
}

// Haqiqiy pointer'li ro'yxat: 2^20 tugun (16 MiB), yonma-yon va sochilgan.
func BenchmarkRoyxat(b *testing.B) {
	n := 1 << 20
	b.Run("yangi", func(b *testing.B) {
		t := RoyxatYangi(n)
		for b.Loop() {
			Royxat(t)
		}
		b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(n), "ns/el")
	})
	b.Run("sochilgan", func(b *testing.B) {
		t := RoyxatSochilgan(n, yangiRand())
		for b.Loop() {
			Royxat(t)
		}
		b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(n), "ns/el")
	})
	b.Run("slice", func(b *testing.B) {
		s := Sonlar(n)
		for b.Loop() {
			Ketma(s)
		}
		b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(n), "ns/el")
	})
}
