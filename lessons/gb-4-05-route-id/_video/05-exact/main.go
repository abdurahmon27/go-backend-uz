package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Product — katalogdagi bitta mahsulot (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

var olma = Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}

// catalog — bozor katalogi (GB-4.03). Hozircha xotirada; bir vaqtdagi soʻrovlar xavfsizligi GB-4.07'da.
var catalog = []Product{
	olma,
	{ID: 2, Name: "non", Price: 400_000, Stock: 120},
	{ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
}

// healthHandler — server tirikligini tekshirish: har soʻrovga «ok» (GB-4.03).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

// homeHandler — bosh sahifa: API nomi.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "bozor API")
}

// catalogHandler — GET /products: catalog'ni JSON massiv qilib beradi. Metod tekshiruvi endi pattern'da.
func catalogHandler(w http.ResponseWriter, r *http.Request) {
	products := catalog
	if products == nil {
		products = []Product{} // nil slice JSON'da null boʻlardi
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products) // Encode xatosi: log'ga yozish GB-8.02'da
}

// findProduct katalogdan ID boʻyicha mahsulotni topadi: topilsa elementning manzili, topilmasa nil (GB-2.04).
func findProduct(products []Product, id int64) *Product {
	for i := range products {
		if products[i].ID == id {
			return &products[i]
		}
	}
	return nil
}

// productHandler — GET /products/{id}: id'ni yoʻldan oʻqiydi (r.PathValue), son boʻlmasa 400, topilmasa 404.
func productHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "id butun son boʻlishi kerak", http.StatusBadRequest)
		return
	}
	p := findProduct(catalog, id)
	if p == nil {
		http.Error(w, "mahsulot topilmadi", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p) // Encode xatosi: log'ga yozish GB-8.02'da
}

// newMux — hamma route'lar shu yerda: newServer ham, testlar ham bir xil mux'ni ishlatadi.
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /products", catalogHandler)
	mux.HandleFunc("GET /products/{id}", productHandler)
	return mux
}

// newServer — bozor API serveri (GB-4.04): endi mux'ni newMux yasaydi.
func newServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func main() {
	srv := newServer(":8765")
	fmt.Println("bozor API: http://localhost:8765")
	err := srv.ListenAndServe()
	fmt.Fprintln(os.Stderr, "xato:", err)
	os.Exit(1)
}
