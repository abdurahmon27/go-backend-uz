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

// Total buyurtma summasini tiyinda qaytaradi.
func (o Order) Total() int64 {
	var total int64
	for _, it := range o.Items {
		total += it.Price * int64(it.Quantity)
	}
	return total
}

func main() {
	olma := Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}
	order := Order{
		ID: 1,
		Items: []OrderItem{
			{ProductID: olma.ID, Price: olma.Price, Quantity: 3},
		},
	}
	fmt.Println(order.Total())
}
