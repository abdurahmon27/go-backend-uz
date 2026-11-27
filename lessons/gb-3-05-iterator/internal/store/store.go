// Package store bozor katalogini xotirada saqlaydi (GB-2.06).
// U internal/ ichida: uni faqat bozor modulining oʻz kodi import qila oladi.
package store

import (
	"iter"
	"slices"
)

// Product — katalogdagi bitta mahsulot (GB-2.03).
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// Store — xotiradagi katalog. Nol qiymati tayyor: var st store.Store.
// products kichik harf bilan: paketdan tashqaridagi kod unga tegolmaydi.
type Store struct {
	products []Product
}

// Add katalogga mahsulot qoʻshadi va unga ID beradi (1, 2, 3 …).
func (s *Store) Add(name string, price int64, stock int) Product {
	p := Product{ID: int64(len(s.products)) + 1, Name: name, Price: price, Stock: stock}
	s.products = append(s.products, p)
	return p
}

// List katalogning nusxasini qaytaradi (GB-2.06): har chaqiruvda butun slice koʻchiriladi.
func (s *Store) List() []Product {
	return slices.Clone(s.products)
}

// All katalogdagi mahsulotlarni birma-bir beradigan iterator qaytaradi:
//
//	for p := range st.All() { … }
//
// Slice nusxalanmaydi va tashqariga chiqmaydi; yield'ga har mahsulotning nusxasi beriladi.
// yield false qaytarsa (sikl break qildi), iterator darhol toʻxtaydi.
func (s *Store) All() iter.Seq[Product] {
	return func(yield func(Product) bool) {
		for _, p := range s.products {
			if !yield(p) {
				return
			}
		}
	}
}

// LowStock qoldigʻi limit'dan kam mahsulotlarni katalog tartibida beradi.
// VAZIFA (1-qadam, TASK.md): hozir u hech narsa bermaydi.
func (s *Store) LowStock(limit int) iter.Seq[Product] {
	return func(yield func(Product) bool) {}
}

// Names mahsulot nomlarini katalog tartibida beradi.
func (s *Store) Names() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, p := range s.products {
			if !yield(p.Name) {
				return
			}
		}
	}
}

// SortedNames mahsulot nomlarini slices.Sort tartibida (baytlar boʻyicha), yangi slice'da qaytaradi.
// VAZIFA (2-qadam, TASK.md): hozir u nil qaytaradi.
func (s *Store) SortedNames() []string {
	return nil
}
