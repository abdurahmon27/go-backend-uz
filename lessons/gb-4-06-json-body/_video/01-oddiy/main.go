package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Product — katalogdagi bitta mahsulot (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

var catalog = []Product{
	{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
	{ID: 2, Name: "non", Price: 400_000, Stock: 120},
	{ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
}

type productInput struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, e apiError) {
	writeJSON(w, status, map[string]apiError{"error": e})
}

func addProduct(in productInput) Product {
	var last int64
	for _, p := range catalog {
		last = max(last, p.ID)
	}
	p := Product{ID: last + 1, Name: in.Name, Price: in.Price, Stock: in.Stock}
	catalog = append(catalog, p)
	return p
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var in productInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, apiError{Code: "bad_json", Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, addProduct(in))
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", createProduct)
	return mux
}

func send(h http.Handler, method, path, body string) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	fmt.Printf("%d %s", rec.Code, rec.Body.String())
}

func main() {
	mux := newMux()
	send(mux, http.MethodPost, "/products", `{"name":"choy","price":2500000,"stock":30}`)
	send(mux, http.MethodPost, "/products", `{"name":`)
}
