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
	olma := Product{Name: "olma", Stock: 40}
	p := &olma
	p.Stock = 37
	q := p
	q.Stock = 35
	fmt.Println(olma.Stock, p.Stock)
}
