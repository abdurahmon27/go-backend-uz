package main

import (
	"errors"
	"fmt"
	"math"
)

var ErrOverflow = errors.New("int64 overflow")

// safeAdd a+b ni qaytaradi; int64 ga sig'masa ErrOverflow.
// Go'da a+b panic bermaydi, jim aylanadi. Shuning uchun natijani hisoblab,
// keyin ishorasini tekshiramiz: musbat+musbat manfiy chiqsa yoki
// manfiy+manfiy musbat (yoki 0) chiqsa, overflow bo'lgan.
func safeAdd(a, b int64) (int64, error) {
	sum := a + b
	if (a > 0 && b > 0 && sum < 0) || (a < 0 && b < 0 && sum >= 0) {
		return 0, ErrOverflow
	}
	return sum, nil
}

func main() {
	fmt.Println(safeAdd(2, 3))                         // 5 <nil>
	fmt.Println(safeAdd(math.MaxInt64, 1))             // 0 int64 overflow
	fmt.Println(safeAdd(math.MinInt64, -1))            // 0 int64 overflow
	fmt.Println(safeAdd(math.MaxInt64, math.MinInt64)) // -1 <nil>
}
