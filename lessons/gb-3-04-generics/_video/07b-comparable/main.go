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

func Contains[T comparable](items []T, x T) bool {
	for _, v := range items {
		if v == x {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println(Contains([]string{"olma", "non"}, "non"), Contains([]int64{1, 2}, 3))
}
