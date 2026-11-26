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
	n := Product{Name: "non", Price: 400_000}
	fmt.Printf("%+v\n", n)
}
