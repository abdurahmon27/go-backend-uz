package main

import "fmt"

// Vazifa: 55 va undan yuqori ball uchun "qoniqarli" degan yana bir pogʻona qoʻshing.
// ball := 60 boʻlganda dastur "qoniqarli" deb yozishi kerak.
func main() {
	ball := 60
	if ball >= 90 {
		fmt.Println("aʼlo")
	} else if ball >= 70 {
		fmt.Println("yaxshi")
	} else {
		fmt.Println("qayta topshiring")
	}
}
