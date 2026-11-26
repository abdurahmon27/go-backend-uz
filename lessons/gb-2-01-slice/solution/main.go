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

// discounted aksiya narxlarini qaytaradi: har bir narx percent foizga arzonroq.
// Narxlar tiyinda (int64); prices oʻzgarmaydi.
func discounted(prices []int64, percent int64) []int64 {
	sale := make([]int64, len(prices))
	for i := range prices {
		sale[i] = prices[i] - prices[i]*percent/100
	}
	return sale
}

func main() {
	stock := []int{40, 12, 7} // olma, non, choy
	sold := []int{2, 1, 0}    // savatda
	fmt.Println("ombor oldin:", stock)
	left := afterOrder(stock, sold)
	fmt.Println("koʻrish:    ", left)
	fmt.Println("ombor keyin:", stock)

	prices := []int64{1_200_000, 500_000, 2_500_000} // tiyinda: 12 000, 5 000, 25 000 soʻm
	fmt.Println("narx oldin: ", prices)
	sale := discounted(prices, 10)
	fmt.Println("aksiya:     ", sale)
	fmt.Println("narx keyin: ", prices)
}
