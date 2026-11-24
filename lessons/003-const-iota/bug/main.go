package main

import "fmt"

type Status int

// v2
const (
	Unknown Status = iota
	New
	Cancelled
	Paid
	Shipped
)

func main() {
	saved := Status(2) // bazada saqlangan: v1'da bu Paid edi
	fmt.Println("to'langanmi?", saved == Paid)
	fmt.Println("bekor qilinganmi?", saved == Cancelled)
}
