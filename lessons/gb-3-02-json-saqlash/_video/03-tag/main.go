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
}

func main() {
	p := Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
