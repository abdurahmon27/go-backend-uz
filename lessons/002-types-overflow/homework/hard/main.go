package main

import (
	"fmt"
	"math"
)

// addOverflows a+b int64'ga sig'masa true qaytaradi.
// Overflow faqat bir xil ishorali ikki son qo'shilganda bo'ladi:
// musbat + musbat manfiy chiqsa yoki manfiy + manfiy musbat chiqsa.
func addOverflows(a, b int64) bool {
	sum := a + b // Go'da bu jim aylanadi, panic bo'lmaydi
	return (a >= 0) == (b >= 0) && (sum >= 0) != (a >= 0)
}

func main() {
	fmt.Println(addOverflows(1, 2))                         // false
	fmt.Println(addOverflows(math.MaxInt64, 1))             // true
	fmt.Println(addOverflows(math.MinInt64, -1))            // true
	fmt.Println(addOverflows(math.MaxInt64, math.MinInt64)) // false: -1
}
