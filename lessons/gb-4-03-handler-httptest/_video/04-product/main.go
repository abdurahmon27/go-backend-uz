package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Product — katalogdagi bitta mahsulot (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

var olma = Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}

// healthHandler — server tirikligini tekshirish: har soʻrovga «ok».
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

// productHandler — GET'ga olma mahsulotini JSON qilib beradi, boshqa metodga 405.
func productHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(olma) // Encode xatosi: log'ga yozish GB-8.02'da
}

// main serverni ishga tushirmaydi (u GB-4.04'da): handler'ni oddiy funksiya kabi chaqirib, javobni chiqaradi.
func main() {
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rec := httptest.NewRecorder()
	productHandler(rec, req)
	fmt.Printf("%d %s", rec.Code, rec.Body.String())
}
