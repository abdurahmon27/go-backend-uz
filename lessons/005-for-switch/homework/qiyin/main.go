package main

import "fmt"

func main() {
	jadval := [][]int{
		{3, 5, 7},
		{2, 0, 4},
		{0, 1, 9},
	}
	topildi := false
tashqi:
	for i, qator := range jadval {
		for j, v := range qator {
			if v == 0 {
				fmt.Printf("birinchi 0: qator %d, ustun %d\n", i, j)
				topildi = true
				break tashqi // ikkala sikldan chiqadi
			}
		}
	}
	if !topildi {
		fmt.Println("0 yoʻq")
	}
}
