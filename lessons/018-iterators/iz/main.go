// Iz: kim qachon ishlaydi? Iterator va sikl tanasi navbatma-navbat gapiradi.
package main

import (
	"fmt"
	"iter"
)

func Summalar(ss []int64) iter.Seq[int64] {
	return func(yield func(int64) bool) {
		for _, s := range ss {
			fmt.Println("iterator:", s, "ni beraman")
			davom := yield(s)
			fmt.Println("iterator: yield", davom, "qaytardi")
			if !davom {
				return
			}
		}
		fmt.Println("iterator: hammasi berildi")
	}
}

func main() {
	for s := range Summalar([]int64{150000, 42000, 310000}) {
		fmt.Println("  sikl:", s)
		if s < 50000 {
			break
		}
	}
	fmt.Println("sikldan keyin")
}
