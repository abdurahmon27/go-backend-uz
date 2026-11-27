package main

import (
	"fmt"
	"io"
	"os"
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

func writeTotal(w io.Writer, total int64) error {
	_, err := fmt.Fprintf(w, "jami: %d tiyin\n", total)
	return err
}

func main() {
	if err := writeTotal(os.Stdout, 4_400_000); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
