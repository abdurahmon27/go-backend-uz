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

// AddItem buyurtmaga qator qoʻshadi; narx mahsulotdan nusxalanadi.
func (o Order) AddItem(p Product, quantity int) {
	item := OrderItem{ProductID: p.ID, Price: p.Price, Quantity: quantity}
	o.Items = append(o.Items, item)
}

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
	non := Product{ID: 2, Name: "non", Price: 400_000, Stock: 100}
	var order Order
	order.AddItem(olma, 3)
	order.AddItem(non, 2)
	fmt.Println(len(order.Items), order.Total())
}
