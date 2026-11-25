// katta: 4 KiB'lik struct. Qiymat bilan uzatilsa, har chaqiruvda 4096 bayt
// nusxalanadi; pointer bilan — faqat 8 baytlik manzil.
// go test -bench . ./katta
package katta

import "testing"

type Hisobot struct {
	Qatorlar [512]int64 // 4096 bayt
}

//go:noinline
func jamiQ(h Hisobot) int64 { return h.Qatorlar[0] + h.Qatorlar[511] }

//go:noinline
func jamiP(h *Hisobot) int64 { return h.Qatorlar[0] + h.Qatorlar[511] }

var h Hisobot

func BenchmarkQiymat(b *testing.B) {
	for b.Loop() {
		jamiQ(h)
	}
}

func BenchmarkPointer(b *testing.B) {
	for b.Loop() {
		jamiP(&h)
	}
}
