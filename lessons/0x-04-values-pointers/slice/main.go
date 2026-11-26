// Bashorat 3 (GO-X.04): massiv qiymat — butunlay nusxalanadi; slice esa elementlarni bo'lishadi.
package main

import "fmt"

func main() {
	a := [3]int{1, 2, 3} // massiv: uzunligi turning bir qismi
	b := a
	b[0] = 9

	s := []int{1, 2, 3} // slice
	t := s
	t[0] = 9

	fmt.Println(a[0], s[0]) // 1 9
}
