// Package store bozor katalogini xotirada saqlaydi.
// U internal/ ichida: uni faqat bozor modulining oʻz kodi import qila oladi.
package store

import (
	"errors"
	"fmt"
	"slices"
)

// ErrNotFound — katalogda bunday ID'li mahsulot yoʻq.
var ErrNotFound = errors.New("mahsulot topilmadi")

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// Store — xotiradagi katalog. Nol qiymati tayyor: var st store.Store.
type Store struct {
	products []Product
}

// Add katalogga mahsulot qoʻshadi va unga ID beradi (1, 2, 3 …).
func (s *Store) Add(name string, price int64, stock int) Product {
	p := Product{ID: int64(len(s.products)) + 1, Name: name, Price: price, Stock: stock}
	s.products = append(s.products, p)
	return p
}

// List katalogdagi mahsulotlarning nusxasini qaytaradi.
func (s *Store) List() []Product {
	return slices.Clone(s.products)
}

// Sell id'li mahsulotdan quantity dona sotadi va qoldiqni kamaytiradi.
// Xato boʻlsa qoldiq oʻzgarmaydi.
func (s *Store) Sell(id int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("dona soni musbat boʻlishi kerak, hozir %d", quantity)
	}
	for i := range s.products {
		p := &s.products[i]
		if p.ID == id {
			if p.Stock < quantity {
				return fmt.Errorf("%s: qoldiq %d, soʻralgan %d", p.Name, p.Stock, quantity)
			}
			p.Stock -= quantity
			return nil
		}
	}
	return ErrNotFound
}
