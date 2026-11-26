package oson

import (
	"slices"
	"testing"
)

func TestOraliq(t *testing.T) {
	if got := slices.Collect(Oraliq(3, 6)); !slices.Equal(got, []int{3, 4, 5, 6}) {
		t.Fatalf("Oraliq(3, 6) = %v", got)
	}
	if got := slices.Collect(Oraliq(5, 4)); got != nil {
		t.Fatalf("Oraliq(5, 4) = %v, boʻsh boʻlishi kerak", got)
	}
}

// break dan keyin panic boʻlmasligi kerak: iterator yield natijasini tekshiradi.
func TestBreak(t *testing.T) {
	var olingan []int
	for i := range Oraliq(1, 1000) {
		olingan = append(olingan, i)
		if i == 3 {
			break
		}
	}
	if !slices.Equal(olingan, []int{1, 2, 3}) {
		t.Fatalf("olingan = %v", olingan)
	}
}
