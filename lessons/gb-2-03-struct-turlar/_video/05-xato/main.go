package main

import "fmt"

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

func main() {
	p := Product{1, "olma", 1_200_000}
	fmt.Printf("%+v\n", p)
}
