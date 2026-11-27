package main

import (
	"errors"
	"fmt"
	"net/http"
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

// buy — POST /buy: oʻrikdan (ID 3) bitta dona sotadi. Soddalik uchun har xato 409 Conflict (haqiqiy API'da yoʻq
// mahsulot 404 boʻlardi).
func (s *Store) buy(w http.ResponseWriter, r *http.Request) {
	if err := s.Sell(3, 1); err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	fmt.Fprintln(w, "sotildi")
}

func main() {
	st := &Store{products: map[int64]Product{
		3: {ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
	}}
	http.HandleFunc("POST /buy", st.buy)
	http.ListenAndServe("localhost:8080", nil)
}
