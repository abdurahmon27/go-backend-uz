package main

import "fmt"

// Qiyin: []byte va []int oʻsishini yonma-yon oʻlchaymiz.
func main() {
	var b []byte
	var n []int
	fmt.Print("[]byte:")
	eski := -1
	for i := range 2000 {
		b = append(b, byte(i))
		if cap(b) != eski {
			eski = cap(b)
			fmt.Print(" ", eski)
		}
	}
	fmt.Print("\n[]int: ")
	eski = -1
	for i := range 2000 {
		n = append(n, i)
		if cap(n) != eski {
			eski = cap(n)
			fmt.Print(" ", eski)
		}
	}
	fmt.Println()
}
