package main

import "fmt"

func main() {
	narx := map[string]int{"non": 5000, "sut": 12000}
	narx["choy"] = 25000

	fmt.Println(narx["olma"]) // bunday kalit yoʻq

	n, bor := narx["olma"]
	fmt.Println(n, bor)

	delete(narx, "sut")
	fmt.Println(len(narx))
}
