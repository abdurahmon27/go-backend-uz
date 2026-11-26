package main

import "fmt"

type Holat int // time.Duration kabi: asosi int boʻlgan alohida tur

const (
	Nomalum  Holat = iota // 0: zero value — holat berilmagan
	Yangi                 // 1: faqat nom → oldingi ifoda "Holat = iota" takrorlanadi
	Tolangan              // 2
)

func main() {
	var h Holat // zero value
	fmt.Println(Yangi, Tolangan, h == Nomalum)
}
