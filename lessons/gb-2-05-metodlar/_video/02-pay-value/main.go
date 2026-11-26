package main

import "fmt"

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// OrderItem — buyurtmadagi bitta qator.
type OrderItem struct {
	ProductID int64
	Price     int64 // tiyin, buyurtma paytidagi narx
	Quantity  int
}

// Order — bitta buyurtma.
type Order struct {
	ID     int64
	Items  []OrderItem
	Status OrderStatus
}

// OrderStatus — buyurtma holati.
type OrderStatus int

const (
	StatusNew OrderStatus = iota
	StatusPaid
	StatusShipped
	StatusCancelled
)

// Pay buyurtmani toʻlangan deb belgilaydi.
func (o Order) Pay() {
	o.Status = StatusPaid
}

func main() {
	var order Order
	order.Pay()
	fmt.Println(order.Status)
}
