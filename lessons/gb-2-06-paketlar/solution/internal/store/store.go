// Package store bozor katalogini xotirada saqlaydi.
// U internal/ ichida: uni faqat bozor modulining oʻz kodi import qila oladi.
package store

import (
	"fmt"
	"slices"
)

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// Store — xotiradagi katalog. Nol qiymati tayyor: var st store.Store.
// products kichik harf bilan: paketdan tashqaridagi kod unga yoza olmaydi.
type Store struct {
	products []Product
}

// Add katalogga mahsulot qoʻshadi va unga ID beradi (1, 2, 3 …).
func (s *Store) Add(name string, price int64, stock int) Product {
	p := Product{ID: int64(len(s.products)) + 1, Name: name, Price: price, Stock: stock}
	s.products = append(s.products, p)
	return p
}

// List katalogdagi mahsulotlarning nusxasini qaytaradi: natijaga yozish
// katalogni oʻzgartirmaydi.
func (s *Store) List() []Product {
	return slices.Clone(s.products)
}

// Sell id'li mahsulotdan quantity dona sotadi va qoldiqni kamaytiradi.
// quantity musbat boʻlishi, mahsulot katalogda boʻlishi va qoldiq yetishi kerak;
// aks holda error qaytadi va qoldiq oʻzgarmaydi.
func (s *Store) Sell(id int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("dona soni musbat boʻlishi kerak, hozir %d", quantity)
	}
	for i := range s.products {
		if s.products[i].ID == id {
			if s.products[i].Stock < quantity {
				return fmt.Errorf("%s: qoldiq %d, soʻralgan %d", s.products[i].Name, s.products[i].Stock, quantity)
			}
			s.products[i].Stock -= quantity
			return nil
		}
	}
	return fmt.Errorf("mahsulot %d topilmadi", id)
}
