package main

import "testing"

func TestChegirma(t *testing.T) {
	holatlar := []struct {
		narx, foiz, kutilgan int
		xato                 bool
	}{
		{3000, 0, 3000, false},
		{3000, 10, 2700, false},
		{3000, 100, 0, false},
		{3000, 150, 0, true},
		{3000, -5, 0, true},
	}
	for _, h := range holatlar {
		got, err := chegirma(h.narx, h.foiz)
		if (err != nil) != h.xato || got != h.kutilgan {
			t.Errorf("chegirma(%d, %d) = %d, %v; kutilgan %d, xato=%v", h.narx, h.foiz, got, err, h.kutilgan, h.xato)
		}
	}
}
