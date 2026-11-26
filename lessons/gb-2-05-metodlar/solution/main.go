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
	olma := Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}
	non := Product{ID: 2, Name: "non", Price: 400_000, Stock: 100}
	var order Order
	order.AddItem(olma, 3)
	order.AddItem(non, 2)
	order.Pay()
	fmt.Println(len(order.Items), order.Total(), order.Status)
}

// String holatning nomini qaytaradi; fmt uni oʻzi chaqiradi.
func (s OrderStatus) String() string {
	switch s {
	case StatusNew:
		return "yangi"
	case StatusPaid:
		return "toʻlangan"
	case StatusShipped:
		return "joʻnatilgan"
	case StatusCancelled:
		return "bekor qilingan"
	}
	return "nomaʼlum"
}

// Total buyurtma summasini tiyinda qaytaradi.
func (o *Order) Total() int64 {
	var total int64
	for _, it := range o.Items {
		total += it.Price * int64(it.Quantity)
	}
	return total
}

// Pay buyurtmani toʻlangan deb belgilaydi.
func (o *Order) Pay() {
	o.Status = StatusPaid
}

// AddItem buyurtmaga qator qoʻshadi; narx mahsulotdan nusxalanadi.
// Shu mahsulot buyurtmada bor boʻlsa, oʻsha qatorning Quantity'si oshadi.
func (o *Order) AddItem(p Product, quantity int) {
	for i := range o.Items {
		if o.Items[i].ProductID == p.ID {
			o.Items[i].Quantity += quantity
			return
		}
	}
	item := OrderItem{ProductID: p.ID, Price: p.Price, Quantity: quantity}
	o.Items = append(o.Items, item)
}

// Cancel buyurtmani bekor qiladi: yangi va toʻlangan buyurtma — ha;
// joʻnatilgan va bekor qilingan — yoʻq (error qaytadi, holat oʻzgarmaydi).
func (o *Order) Cancel() error {
	if o.Status != StatusNew && o.Status != StatusPaid {
		return fmt.Errorf("buyurtma %d: holati «%v», bekor qilib boʻlmaydi", o.ID, o.Status)
	}
	o.Status = StatusCancelled
	return nil
}
