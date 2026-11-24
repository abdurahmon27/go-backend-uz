package main

import "fmt"

func hisoblagich() func() int {
	n := 0
	return func() int {
		n++ // tashqi n'ni oʻzgartiradi
		return n
	}
}

func main() {
	a := hisoblagich()
	b := hisoblagich()
	fmt.Println(a(), a(), a())
	fmt.Println(b())
}
