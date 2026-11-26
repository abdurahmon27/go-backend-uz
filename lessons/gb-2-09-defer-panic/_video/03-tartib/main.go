package main

import "fmt"

func main() {
	fmt.Println("hisobot boshlandi")
	defer fmt.Println("defer 1: fayl yopildi")
	defer fmt.Println("defer 2: buyurtma yopildi")
	fmt.Println("hisobot tugadi")
}
