package main

import "fmt"

// Boʻsh slice'ga bittadan beshta son qoʻshamiz va har safar cap'ni chiqaramiz.
func main() {
	var s []int
	for i := range 5 {
		s = append(s, i)
		fmt.Print(cap(s), " ")
	}
	fmt.Println()
}
