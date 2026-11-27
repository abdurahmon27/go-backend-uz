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

// saveProducts katalogni path fayliga JSON qilib yozadi.
func saveProducts(path string, products []Product) error {
	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// loadProducts path faylidagi JSON katalogni oʻqiydi.
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
