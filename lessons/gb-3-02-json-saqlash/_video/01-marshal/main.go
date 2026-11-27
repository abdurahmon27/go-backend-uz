package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
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
