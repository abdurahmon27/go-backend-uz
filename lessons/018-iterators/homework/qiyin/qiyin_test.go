package qiyin

import (
	"slices"
	"testing"
)

func TestReyting(t *testing.T) {
	bb := []Buyurtma{
		{"Kamola", 100000}, {"Aziz", 150000}, {"Bobur", 99000},
		{"Kamola", 50000}, {"Dilnoz", 42000}, {"Bobur", 51000},
	}
	// Aziz 150000, Bobur 150000, Kamola 150000, Dilnoz 42000
	want := []string{"Aziz", "Bobur", "Kamola", "Dilnoz"}
	// Map tartibi har safar boshqa: natija baribir bir xil boʻlishi kerak.
	for range 50 {
		if got := Reyting(bb); !slices.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
