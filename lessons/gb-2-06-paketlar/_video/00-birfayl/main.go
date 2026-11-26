package main

import "fmt"

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// Store — xotiradagi katalog.
type Store struct {
	products []Product
}

// Add katalogga mahsulot qoʻshadi va unga ID beradi.
func (s *Store) Add(name string, price int64, stock int) Product {
	p := Product{ID: int64(len(s.products)) + 1, Name: name, Price: price, Stock: stock}
	s.products = append(s.products, p)
	return p
}

// List katalogdagi mahsulotlarni qaytaradi.
func (s *Store) List() []Product {
	return s.products
}

func main() {
	var st Store
	st.Add("olma", 1_200_000, 40)
	st.products[0].Stock = -5
	fmt.Println(st.products[0].Stock)
}
