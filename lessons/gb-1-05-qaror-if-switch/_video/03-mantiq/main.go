package main

import "fmt"

func main() {
	sold := 7
	stock := 40
	ok := sold > 0 && sold <= stock
	bad := sold <= 0 || sold > stock
	fmt.Println(ok, bad, !ok)
}
