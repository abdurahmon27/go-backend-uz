package main

import (
	"fmt"
)

// Product — katalogdagi bitta mahsulot (GB-2.03).
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// OrderItem — buyurtmadagi bitta qator (GB-2.03).
type OrderItem struct {
	ProductID int64
	Price     int64 // tiyin
	Quantity  int
}

// OrderStatus — buyurtma holati (GB-2.03).
type OrderStatus int

const (
	StatusNew OrderStatus = iota
	StatusPaid
)

// Order — bitta buyurtma (GB-2.03).
type Order struct {
	ID     int64
	Items  []OrderItem
	Status OrderStatus
}

func Filter[T any](items []T, keep func(T) bool) []T {
	var out []T
	for _, x := range items {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

func main() {
	fmt.Printf("%T\n", Filter[Product])
	fmt.Printf("%T\n", Filter[Order])
}
