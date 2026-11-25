package main

import "fmt"

// boshi/ bilan bir xil, lekin slice'ning oʻzi ham Println'ga beriladi:
// endi u funksiyadan "qochadi" (heap'da yashaydi) va stekdagi zaxira joy ishlatilmaydi.
func main() {
	var s []int
	for i := range 5 {
		s = append(s, i)
		fmt.Print(cap(s), " ")
	}
	fmt.Println(s)
}
