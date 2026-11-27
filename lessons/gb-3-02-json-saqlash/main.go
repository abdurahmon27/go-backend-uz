package main

import (
	"encoding/json"
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
// VAZIFA (1-qadam): json tag'larini qoʻying (TASK.md).
type OrderItem struct {
	ProductID int64
	Price     int64 // tiyin, buyurtma paytidagi narx
	Quantity  int
}

// Order — bitta buyurtma.
// VAZIFA (1-qadam): json tag'larini qoʻying (TASK.md).
type Order struct {
	ID     int64
	Items  []OrderItem
	Status OrderStatus
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
// VAZIFA (2-qadam): fayl yoʻq boʻlsa — boʻsh katalog; JSON xatosiga fayl nomi (TASK.md).
func loadProducts(path string) ([]Product, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var products []Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
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
