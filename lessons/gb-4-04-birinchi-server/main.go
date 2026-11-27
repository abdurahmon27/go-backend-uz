package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Product — katalogdagi bitta mahsulot (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

var olma = Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}

// catalog — bozor katalogi (GB-4.03): /products shuni beradi.
var catalog = []Product{
	olma,
	{ID: 2, Name: "non", Price: 400_000, Stock: 120},
	{ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
}

// healthHandler — server tirikligini tekshirish: har soʻrovga «ok» (GB-4.03).
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

// productHandler — GET'ga olma mahsulotini JSON qilib beradi, boshqa metodga 405 (GB-4.03).
func productHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(olma) // Encode xatosi: log'ga yozish GB-8.02'da
}

// catalogHandler — GET'ga catalog'ni JSON massiv qilib beradi, boshqa metodga 405 (GB-4.03 vazifasi).
func catalogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	products := catalog
	if products == nil {
		products = []Product{} // nil slice JSON'da null boʻlardi
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products) // Encode xatosi: log'ga yozish GB-8.02'da
}

// newServer — VAZIFA (TASK.md): bozor API serverini yasab qaytaring.
// Kerak: mux (/health, /products, /products/1) va &http.Server{Addr, Handler, ReadHeaderTimeout}.
// main uni :8765 portida ishga tushiradi; tekshiruv_test.go esa uni portsiz, Handler va ReadHeaderTimeout orqali tekshiradi.
func newServer(addr string) *http.Server {
	// VAZIFA: shu yerga yozing
	return nil
}

func main() {
	srv := newServer(":8765")
	if srv == nil {
		fmt.Fprintln(os.Stderr, "xato: newServer hali yozilmagan (TASK.md)")
		os.Exit(1)
	}
	fmt.Println("bozor API: http://localhost:8765")
	err := srv.ListenAndServe()
	fmt.Fprintln(os.Stderr, "xato:", err)
	os.Exit(1)
}
