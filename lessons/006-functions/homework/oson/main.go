package main

import (
	"errors"
	"fmt"
)

// minMax: eng kichik va eng katta son; boʻsh slice uchun xato.
func minMax(sonlar []int) (int, int, error) {
	if len(sonlar) == 0 {
		return 0, 0, errors.New("boʻsh slice")
	}
	lo, hi := sonlar[0], sonlar[0]
	for _, n := range sonlar[1:] {
		lo = min(lo, n) // min/max: Go 1.21 dan beri ichki funksiyalar
		hi = max(hi, n)
	}
	return lo, hi, nil
}

func main() {
	fmt.Println(minMax([]int{7, -3, 12, 0}))
	fmt.Println(minMax(nil))
}
