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
// TASK.md, 1-qadam: hozir hech narsa qilmaydi va doim nil qaytaradi.
func sellByID(products []Product, id int64, quantity int) error {
	return nil
}

// cheapest katalogdagi eng arzon mahsulotning manzilini qaytaradi; katalog boʻsh boʻlsa nil.
// TASK.md, 2-qadam: hozir doim nil qaytaradi.
func cheapest(products []Product) *Product {
	return nil
}
