package tozala

import "testing"

// Apostroflari allaqachon toʻgʻri matn: hech narsa oʻzgarmaydi.
const togri = "Oʻzbekiston togʻlari, maʼno va qoʻshiq"

func BenchmarkEski(b *testing.B) {
	for b.Loop() {
		tozalaEski(togri)
	}
}

func BenchmarkGuru(b *testing.B) {
	for b.Loop() {
		tozala(togri)
	}
}
