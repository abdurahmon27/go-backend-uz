package main

import "fmt"

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

func sell(p Product, quantity int) {
	p.Stock -= quantity
	fmt.Println("sell ichida:", p.Stock)
}

func main() {
	olma := Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}
	sell(olma, 3)
	fmt.Println("main'da:", olma.Stock)
}
