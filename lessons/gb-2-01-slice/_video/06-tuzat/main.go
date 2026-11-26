package main

import "fmt"

// afterOrder buyurtmadan keyingi qoldiqni hisoblaydi; stock oʻzgarmaydi.
func afterOrder(stock []int, sold []int) []int {
	left := make([]int, len(stock))
	for i := range stock {
		left[i] = stock[i] - sold[i]
	}
	return left
}

func main() {
	stock := []int{40, 12, 7} // olma, non, choy
	sold := []int{2, 1, 0}    // savatda
	fmt.Println("ombor oldin:", stock)
	left := afterOrder(stock, sold)
	fmt.Println("koʻrish:    ", left)
	fmt.Println("ombor keyin:", stock)
}
