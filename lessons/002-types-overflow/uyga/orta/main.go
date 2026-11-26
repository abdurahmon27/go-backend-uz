package main

import (
	"errors"
	"fmt"
	"math"
)

var ErrOverflow = errors.New("int64 overflow")

// Vazifa: a+b int64 ga sig'masa, natija o'rniga ErrOverflow qaytaring.
// Ishora: overflow faqat ishorasi bir xil ikki son qo'shilganda bo'ladi
// (musbat + musbat yoki manfiy + manfiy), va yig'indining ishorasi o'zgaradi.
func safeAdd(a, b int64) (int64, error) {
	return a + b, nil // hozir jim aylanadi
}

func main() {
	fmt.Println(safeAdd(2, 3))                         // kutilgan: 5 <nil>
	fmt.Println(safeAdd(math.MaxInt64, 1))             // kutilgan: 0 int64 overflow
	fmt.Println(safeAdd(math.MinInt64, -1))            // kutilgan: 0 int64 overflow
	fmt.Println(safeAdd(math.MaxInt64, math.MinInt64)) // kutilgan: -1 <nil>
}
