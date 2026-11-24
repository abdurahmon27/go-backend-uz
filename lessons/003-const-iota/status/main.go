package main

import "fmt"

type Status int

const (
	Unknown Status = iota
	New
	Paid
	Shipped
)

func main() {
	fmt.Println(Unknown, New, Paid, Shipped)

	var s Status
	fmt.Println(s == Unknown)
}
