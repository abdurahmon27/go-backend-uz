package main

import (
	"fmt"

	"example.com/tolov"
	"example.com/yetkazish"
)

// tolov v1.0.0 example.com/log v1.2.0 ni, yetkazish v1.0.0 esa v1.3.0 ni talab qiladi.
// Eng yangisi v1.5.0. Build'ga qaysi biri kiradi?
func main() {
	fmt.Println("tolov ichidagi log:    ", tolov.Log())
	fmt.Println("yetkazish ichidagi log:", yetkazish.Log())
}
