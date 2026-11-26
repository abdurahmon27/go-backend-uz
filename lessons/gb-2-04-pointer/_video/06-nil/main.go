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

// findProduct katalogdan ID boʻyicha mahsulotni topadi: topilsa elementning manzili, topilmasa nil.
func findProduct(products []Product, id int64) *Product {
	for i := range products {
		if products[i].ID == id {
			return &products[i]
		}
	}
	return nil
}

func main() {
	catalog := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 100},
	}
	p := findProduct(catalog, 9)
	fmt.Println(p == nil)
	fmt.Println(p.Stock)
}
