package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Product — katalogdagi bitta mahsulot (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

var (
	// ErrStock — mahsulot yoʻq yoki omborda soʻralgancha dona yoʻq.
	ErrStock = errors.New("qoldiq yetmaydi")
	// ErrNotFound — katalogda bunday ID'li mahsulot yoʻq.
	ErrNotFound = errors.New("mahsulot topilmadi")
)

// Store — xotiradagi katalog: ID → mahsulot. Server har ulanishni alohida goroutine'da ishlatadi, shuning uchun
// bir vaqtda kelgan soʻrovlar bir vaqtda ishlaydi: products'ga HAR murojaat (oʻqish ham) mu qulfi ostida boʻlishi kerak.
type Store struct {
	mu       sync.Mutex
	products map[int64]Product
}

// NewStore bozor katalogini (olma, non, oʻrik) yaratadi.
func NewStore() *Store {
	return &Store{products: map[int64]Product{
		1: {ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		2: {ID: 2, Name: "non", Price: 400_000, Stock: 120},
		3: {ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
	}}
}

// Sell id'li mahsulotdan quantity dona sotadi (quantity musbat deb olinadi: tekshiruv GB-2.08'dagi store paketida).
// Video: qulf shu yerda qoʻyiladi.
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

// Stock id'li mahsulotning qoldigʻini beradi.
func (s *Store) Stock(id int64) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.products[id].Stock
}

// Restock id'li mahsulotning qoldigʻiga quantity dona qoʻshadi; mahsulot yoʻq boʻlsa ErrNotFound.
//
// VAZIFA 1 (TASK.md): bu metod hozir qulfsiz — uni ham himoyalang.
func (s *Store) Restock(id int64, quantity int) error {
	p, ok := s.products[id]
	if !ok {
		return ErrNotFound
	}
	p.Stock += quantity
	s.products[id] = p
	return nil
}

// Total hamma mahsulotlarning qoldigʻi yigʻindisini beradi.
//
// VAZIFA 2 (TASK.md): hozir hech narsani sanamaydi.
func (s *Store) Total() int {
	return 0
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

// slowHandler — ogʻir ishni oʻrnini bosadi: bir soniya uxlaydi. Soʻrovlar bir vaqtda ishlashini koʻrsatadi.
func slowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("boshladi")
	time.Sleep(time.Second)
	fmt.Println("tugadi")
}

func main() {
	st := NewStore()
	http.HandleFunc("POST /buy", st.buy)
	http.HandleFunc("/slow", slowHandler)
	log.Fatal(http.ListenAndServe("localhost:8085", nil))
}
