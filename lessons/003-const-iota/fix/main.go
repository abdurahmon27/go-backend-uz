package main

import "fmt"

type Status int

// bazada saqlanadigan qiymatlar: har biri qo'lda yoziladi va hech qachon o'zgarmaydi
const (
	Unknown   Status = 0
	New       Status = 1
	Paid      Status = 2
	Shipped   Status = 3
	Cancelled Status = 4 // yangi qiymat oxiriga
)

func main() {
	saved := Status(2)
	fmt.Println("to'langanmi?", saved == Paid)
}
