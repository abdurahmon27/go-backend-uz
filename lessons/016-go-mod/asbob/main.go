package main

import "fmt"

//go:generate go tool stringer -type=Status

// Status — 3-darsdagi enum: zero value maʼnoli holat emas (Unknown).
type Status int

const (
	Unknown Status = iota
	New
	Paid
	Shipped
)

func main() {
	fmt.Println(Unknown, New, Paid, Shipped)
}
