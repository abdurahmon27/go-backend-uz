package tozala

import "testing"

var matnlar = []struct{ nom, s string }{
	{"toza", "Bugun havo yaxshi, darsdan keyin bozorga boramiz"},
	{"apostrof", "Oʻzbekiston togʻlari, maʼno va qoʻshiq"},
}

func BenchmarkEski(b *testing.B) {
	for _, m := range matnlar {
		b.Run(m.nom, func(b *testing.B) {
			for b.Loop() {
				tozalaEski(m.s)
			}
		})
	}
}

func BenchmarkTez(b *testing.B) {
	for _, m := range matnlar {
		b.Run(m.nom, func(b *testing.B) {
			for b.Loop() {
				tozala(m.s)
			}
		})
	}
}
