package main

import (
	"errors"
	"iter"
	"sync"
)

// Product — katalogdagi bitta mahsulot (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

// ErrStock — mahsulot yoʻq yoki omborda soʻralgancha dona yoʻq.
var ErrStock = errors.New("qoldiq yetmaydi")

// Store — xotiradagi katalog: ID → mahsulot.
type Store struct {
	mu       sync.Mutex
	products map[int64]Product
}

// Sell id'li mahsulotdan quantity dona sotadi (quantity musbat deb olinadi: bu laboratoriya soddalashtirilgan).
func (s *Store) Sell(id int64, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[id]
	if !ok || p.Stock < quantity {
		return ErrStock
	}
	p.Stock -= quantity
	s.products[id] = p
	return nil
}

// All katalogni birma-bir beradigan iterator (GB-3.05), qulf bilan.
func (s *Store) All() iter.Seq[Product] {
	return func(yield func(Product) bool) {
		s.mu.Lock()
		defer s.mu.Unlock()
		for _, p := range s.products {
			if !yield(p) {
				return
			}
		}
	}
}

func main() {
	st := &Store{products: map[int64]Product{
		3: {ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
	}}
	for p := range st.All() {
		st.Sell(p.ID, 1)
	}
}
