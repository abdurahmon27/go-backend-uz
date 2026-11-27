package main

import (
	"errors"
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

// StockError — omborda soʻralgancha mahsulot yoʻq (GB-2.08).
type StockError struct {
	Name  string
	Stock int
	Want  int
}

func (e *StockError) Error() string {
	return fmt.Sprintf("%s: qoldiq %d, soʻralgan %d", e.Name, e.Stock, e.Want)
}

func main() {
	err := fmt.Errorf("buyurtma 7: %w", &StockError{Name: "olma", Stock: 40, Want: 41})
	if se, ok := errors.AsType(err); ok {
		fmt.Println(se)
	}
}
