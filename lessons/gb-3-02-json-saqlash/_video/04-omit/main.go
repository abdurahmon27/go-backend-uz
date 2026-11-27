package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Product — katalogdagi bitta mahsulot. Tag'lar — JSON'dagi kalit nomlari.
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
	Note  string `json:"note,omitempty"`
}

func main() {
	products := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "oʻrik", Price: 3_500_000, Stock: 5, Note: "quritilgan"},
	}
	enc := json.NewEncoder(os.Stdout)
	for _, p := range products {
		if err := enc.Encode(p); err != nil {
			fmt.Fprintln(os.Stderr, "xato:", err)
			os.Exit(1)
		}
	}
}
