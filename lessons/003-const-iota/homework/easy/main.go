package main

import "fmt"

type Weekday int

const (
	Monday Weekday = iota + 1 // 0 = "kun tanlanmagan" (zero value)
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	fmt.Println(Monday, Friday, Sunday) // 1 5 7
	var d Weekday
	fmt.Println(d) // 0: hech qaysi kun emas
}
