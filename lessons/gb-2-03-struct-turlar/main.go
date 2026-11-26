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

func main() {
	p := Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}
	o := Order{
		ID: 1,
		Items: []OrderItem{
			{ProductID: p.ID, Price: p.Price, Quantity: 3},
		},
	}
	fmt.Printf("%+v\n", o)

	var empty Order
	fmt.Println(empty.Status == StatusNew, len(empty.Items))
}

// orderTotal buyurtmaning jami summasini tiyinda qaytaradi: har qatorda narx × dona.
// TASK.md, 1-qadam: hozir doim 0 qaytaradi.
func orderTotal(o Order) int64 {
	return 0
}

// canCancel buyurtmani bekor qilish mumkinmi: yangi va toʻlangan — ha,
// joʻnatilgan va bekor qilingan — yoʻq. TASK.md, 2-qadam: hozir doim false.
func canCancel(o Order) bool {
	return false
}
