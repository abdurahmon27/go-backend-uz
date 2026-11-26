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

// ErrBadQuantity — soʻralgan dona soni musbat emas.
// VAZIFA (1-qadam): Sell bu xatoni hali qaytarmaydi (TASK.md).
var ErrBadQuantity = errors.New("dona soni musbat boʻlishi kerak")

// StockError — omborda soʻralgancha mahsulot yoʻq.
type StockError struct {
	Name  string
	Stock int // omborda bor
	Want  int // soʻralgan
}

func (e *StockError) Error() string {
	return fmt.Sprintf("%s: qoldiq %d, soʻralgan %d", e.Name, e.Stock, e.Want)
}

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
// Xato boʻlsa qoldiq oʻzgarmaydi. Chaqiruvchi xatoni tekshira oladi:
// errors.Is(err, ErrNotFound) — mahsulot yoʻq,
// errors.As(err, &se) (se *StockError) — qoldiq yetmaydi.
func (s *Store) Sell(id int64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("dona soni musbat boʻlishi kerak, hozir %d", quantity)
	}
	for i := range s.products {
		p := &s.products[i]
		if p.ID == id {
			if p.Stock < quantity {
				return &StockError{Name: p.Name, Stock: p.Stock, Want: quantity}
			}
			p.Stock -= quantity
			return nil
		}
	}
	return fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
}
