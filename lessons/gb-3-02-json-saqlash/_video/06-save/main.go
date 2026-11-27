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

// run katalogni bozor.json'ga saqlaydi (GB-3.01'dagi shakl).
func run() error {
	products := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "oʻrik", Price: 3_500_000, Stock: 5, Note: "quritilgan"},
	}
	if err := saveProducts("bozor.json", products); err != nil {
		return err
	}
	fmt.Println("saqlandi:", len(products), "ta mahsulot")
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
