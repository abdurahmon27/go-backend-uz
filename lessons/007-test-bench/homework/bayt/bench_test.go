package tozala

import "testing"

const matn = "Bugun havo yaxshi, darsdan keyin bozorga boramiz"

func BenchmarkContainsAny(b *testing.B) {
	for b.Loop() {
		tozalaCA(matn)
	}
}

func BenchmarkBayt(b *testing.B) {
	for b.Loop() {
		tozala(matn)
	}
}
