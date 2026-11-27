// Package store bozor katalogini xotirada saqlaydi (GB-2.06).
package store

import (
	"fmt"
	"slices"
)

// Product — katalogdagi bitta mahsulot (GB-2.03).
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// Store — xotiradagi katalog; products kichik harf bilan, tashqariga koʻrinmaydi.
type Store struct {
	products []Product
}

// Add katalogga mahsulot qoʻshadi va unga ID beradi (1, 2, 3 …).
func (s *Store) Add(name string, price int64, stock int) Product {
	p := Product{ID: int64(len(s.products)) + 1, Name: name, Price: price, Stock: stock}
	s.products = append(s.products, p)
	return p
}

// List katalogning nusxasini qaytaradi (GB-2.06).
func (s *Store) List() []Product {
	return slices.Clone(s.products)
}

// All mahsulotlarni birma-bir yield'ga beradi va har qadamni chiqaradi.
func (s *Store) All() func(yield func(Product) bool) {
	return func(yield func(Product) bool) {
		for _, p := range s.products {
			fmt.Println("iterator:", p.Name, "beraman")
			if !yield(p) {
				fmt.Println("iterator: yield false → toʻxtayman")
				return
			}
		}
	}
}
