// Videodagi "qora quti" server: bozor katalogi HTTP orqali. Uni GB-4.04'da qatorma-qator yozasiz;
// bu darsda faqat `curl -v` bilan uning soʻrov va javob matnini oʻqiymiz.
//
//	go run ./_video/server            (localhost:8765)
//	go run ./_video/server localhost:9000
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

type store struct {
	mu       sync.Mutex
	products []Product
	nextID   int64
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "ichki xato", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (s *store) find(id int64) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.products {
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}

func (s *store) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id butun son boʻlishi kerak"})
		return
	}
	p, ok := s.find(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "mahsulot topilmadi"})
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *store) create(w http.ResponseWriter, r *http.Request) {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mt != "application/json" {
		writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{"error": "Content-Type application/json boʻlishi kerak"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	var in Product
	if err := dec.Decode(&in); err != nil {
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "body juda katta"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON notoʻgʻri"})
		return
	}
	if dec.More() || in.Name == "" || in.Price <= 0 || in.Stock < 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name, price (>0) va stock (>=0) kerak, ortiqcha matn yoʻq"})
		return
	}
	s.mu.Lock()
	in.ID = s.nextID
	s.nextID++
	s.products = append(s.products, in)
	s.mu.Unlock()
	w.Header().Set("Location", fmt.Sprintf("/products/%d", in.ID))
	writeJSON(w, http.StatusCreated, in)
}

func main() {
	addr := "localhost:8765"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	s := &store{nextID: 3, products: []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 100},
	}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}", s.get)
	mux.HandleFunc("POST /products", s.create)
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Println("bozor serveri:", addr)
	log.Fatal(srv.ListenAndServe())
}
