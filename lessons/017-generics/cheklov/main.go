// cheklov: oʻz constraint'imiz. ~int64 — "asosiy turi int64 boʻlgan har qanday tur".
// go run ./cheklov
package main

import "fmt"

type Son interface {
	~int64 | ~float64
}

type Ball int64 // asosiy (underlying) turi int64

func Jami[T Son](xs []T) T {
	var s T
	for _, x := range xs {
		s += x
	}
	return s
}

func main() {
	fmt.Println(Jami([]Ball{90, 85, 70}))
	fmt.Println(Jami([]float64{1.5, 2.25}))
}
