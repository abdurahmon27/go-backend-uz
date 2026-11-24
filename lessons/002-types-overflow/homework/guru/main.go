package main

import (
	"fmt"
	"math"
)

// abs8 naive: -x. Lekin int8'da -(-128) = -128, chunki +128 sig'maydi.
func abs8(x int8) int8 {
	if x < 0 {
		return -x
	}
	return x
}

// abs8Safe natijani kattaroq turda qaytaradi: 128 int16'ga sig'adi.
func abs8Safe(x int8) int16 {
	v := int16(x)
	if v < 0 {
		return -v
	}
	return v
}

func main() {
	fmt.Println(abs8(-5), abs8(math.MinInt8))         // 5 -128
	fmt.Println(abs8Safe(-5), abs8Safe(math.MinInt8)) // 5 128
}
