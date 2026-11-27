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
	products := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 3},
		{ID: 3, Name: "choy", Price: 2_500_000, Stock: 0},
	}
	orders := []Order{
		{ID: 7, Status: StatusPaid, Items: []OrderItem{{ProductID: 1, Price: 1_200_000, Quantity: 3}}},
		{ID: 8, Status: StatusNew, Items: []OrderItem{{ProductID: 2, Price: 400_000, Quantity: 5}}},
		{ID: 9, Status: StatusPaid, Items: []OrderItem{{ProductID: 2, Price: 400_000, Quantity: 2}}},
	}

	low := Filter(orders, func(p Product) bool { return p.Stock < 5 })
	fmt.Println(low, products)
}
