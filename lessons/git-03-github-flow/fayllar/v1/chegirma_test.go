package main

import "testing"

func TestChegirma(t *testing.T) {
	holatlar := []struct{ narx, foiz, kutilgan int }{
		{3000, 0, 3000},
		{3000, 10, 2700},
	}
	for _, h := range holatlar {
		if got := chegirma(h.narx, h.foiz); got != h.kutilgan {
			t.Errorf("chegirma(%d, %d) = %d, kutilgan %d", h.narx, h.foiz, got, h.kutilgan)
		}
	}
}
