package orta

import (
	"fmt"
	"slices"
	"testing"
)

func TestZanjir(t *testing.T) {
	summalar := []int64{150000, 42000, 310000, 99000}
	katta := Filtr(slices.Values(summalar), func(s int64) bool { return s >= 99000 })
	matn := Aylantir(katta, func(s int64) string { return fmt.Sprint(s, " soʻm") })

	got := slices.Collect(matn)
	want := []string{"150000 soʻm", "310000 soʻm", "99000 soʻm"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// Ichki iterator ham toʻxtashi kerak: break'dan keyin hech kim yield chaqirmaydi.
func TestBreak(t *testing.T) {
	n := 0
	for range Filtr(slices.Values([]int{1, 2, 3, 4, 5, 6}), func(i int) bool { return i%2 == 0 }) {
		n++
		break
	}
	if n != 1 {
		t.Fatalf("n = %d", n)
	}
}
