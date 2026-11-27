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

// Tiyin — pul miqdori, tiyinda.
type Tiyin int64

type Number interface {
	int | int64
}

func Sum[N Number](nums []N) N {
	var total N
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	fmt.Println(Sum([]Tiyin{1_200_000, 400_000}))
}
