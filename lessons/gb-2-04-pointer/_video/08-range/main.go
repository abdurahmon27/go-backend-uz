package main

import "fmt"

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

func sell(p *Product, quantity int) {
	p.Stock -= quantity
}

// findProduct — XATO variant: range bergan p elementning nusxasi.
func findProduct(products []Product, id int64) *Product {
	for _, p := range products {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

func main() {
	catalog := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 100},
	}
	p := findProduct(catalog, 2)
	sell(p, 1)
	fmt.Println(p.Stock, catalog[1].Stock)
}
