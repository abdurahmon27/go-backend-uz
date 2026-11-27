package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Product — katalogdagi bitta mahsulot. Tag'lar — bozor.json'dagi kalit nomlari.
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
	Note  string `json:"note,omitempty"`
}

// OrderStatus — buyurtma holati (GB-2.03).
type OrderStatus int

const (
	StatusNew OrderStatus = iota
	StatusPaid
	StatusShipped
	StatusCancelled
)

// OrderItem — buyurtmadagi bitta qator.
type OrderItem struct {
	ProductID int64 `json:"product_id"`
	Price     int64 `json:"price"` // tiyin, buyurtma paytidagi narx
	Quantity  int   `json:"quantity"`
}

// Order — bitta buyurtma. Status'da omitempty yoʻq: StatusNew = 0 ham yozilishi kerak.
type Order struct {
	ID     int64       `json:"id"`
	Items  []OrderItem `json:"items"`
	Status OrderStatus `json:"status"`
}

// saveProducts katalogni path fayliga JSON qilib yozadi.
func saveProducts(path string, products []Product) error {
	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// loadProducts path faylidagi JSON katalogni oʻqiydi.
// Fayl hali yoʻq boʻlsa (bozor birinchi marta ishga tushdi), katalog boʻsh: nil, nil.
// JSON buzuq boʻlsa, xato fayl nomi bilan %w orqali oʻraladi.
func loadProducts(path string) ([]Product, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var products []Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return products, nil
}

// run katalogni bozor.json'dan oʻqib, har mahsulotni chiqaradi (GB-3.01'dagi shakl).
func run() error {
	products, err := loadProducts("bozor.json")
	if err != nil {
		return err
	}
	for _, p := range products {
		fmt.Printf("%d %s: %d dona\n", p.ID, p.Name, p.Stock)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
