package main

import "fmt"

func main() {
	const (
		Nomalum  = iota // 0: qiymat berilmagan (nol qiymat) holat
		Yangi           // 1
		Tolangan        // 2
	)
	fmt.Println(Yangi, Tolangan)
}
