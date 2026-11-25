package main

import "testing"

// Bir million buyurtmaning toʻlangan summasini hisoblaymiz.
// Ish bir xil, faqat xotirada oʻqiladigan baytlar soni har xil.

func BenchmarkBuyurtma(b *testing.B) {
	s := make([]Buyurtma, n)
	for i := range s {
		s[i] = Buyurtma{Tolandi: i%2 == 0, Summa: int64(i)}
	}
	for b.Loop() {
		var jami int64
		for i := range s {
			if s[i].Tolandi {
				jami += s[i].Summa
			}
		}
		_ = jami
	}
}

func BenchmarkIxcham(b *testing.B) {
	s := make([]Ixcham, n)
	for i := range s {
		s[i] = Ixcham{Tolandi: i%2 == 0, Summa: int64(i)}
	}
	for b.Loop() {
		var jami int64
		for i := range s {
			if s[i].Tolandi {
				jami += s[i].Summa
			}
		}
		_ = jami
	}
}
