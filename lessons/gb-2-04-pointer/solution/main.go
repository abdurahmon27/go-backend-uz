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
	if p == nil {
		fmt.Println("mahsulot topilmadi: 9")
		return
	}
	sell(p, 1)
	fmt.Println(p.Stock)
}

// sellByID katalogdagi id raqamli mahsulotdan quantity dona sotadi va xatoni qaytaradi.
func sellByID(products []Product, id int64, quantity int) error {
	p := findProduct(products, id)
	if p == nil {
		return fmt.Errorf("mahsulot topilmadi: %d", id)
	}
	if quantity <= 0 {
		return fmt.Errorf("miqdor musbat boʻlishi kerak, keldi %d", quantity)
	}
	if quantity > p.Stock {
		return fmt.Errorf("%s: qoldiq %d, soʻralgan %d", p.Name, p.Stock, quantity)
	}
	sell(p, quantity)
	return nil
}

// cheapest katalogdagi eng arzon mahsulotning manzilini qaytaradi; katalog boʻsh boʻlsa nil.
func cheapest(products []Product) *Product {
	var best *Product
	for i := range products {
		if best == nil || products[i].Price < best.Price {
			best = &products[i]
		}
	}
	return best
}
