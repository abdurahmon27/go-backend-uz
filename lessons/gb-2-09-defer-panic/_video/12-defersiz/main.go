package main

import "fmt"

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// findProduct katalogdan ID boʻyicha mahsulotni topadi: topilsa elementning manzili, topilmasa nil (GB-2.04).
func findProduct(products []Product, id int64) *Product {
	for i := range products {
		if products[i].ID == id {
			return &products[i]
		}
	}
	return nil
}

// sell p mahsulotdan quantity dona sotadi.
func sell(p *Product, quantity int) {
	defer fmt.Println("defer: sell tugadi")
	p.Stock -= quantity
	fmt.Println("sotildi")
}

// safely fn'ni chaqiradi; fn panic qilsa, panic'ni error qilib qaytaradi.
func safely(fn func()) (err error) {
	if r := recover(); r != nil {
		err = fmt.Errorf("panic: %v", r)
	}
	fn()
	return err
}

func main() {
	products := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 100},
		{ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
	}
	p := findProduct(products, 9) // topilmadi: nil
	err := safely(func() {
		sell(p, 1)
	})
	fmt.Println("xato:", err)
	fmt.Println("main oxiri")
}
