package main

import "fmt"

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

func main() {
	p := Product{ID: 1, Name: "olma", Price: 1_200_000}
	p.Stock = 40
	q := p
	q.Stock = 0
	fmt.Println(p.Stock, q.Stock)
}
