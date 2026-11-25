package main

import "fmt"

// Bittalab append qilamiz va cap oʻzgargan joylarni chiqaramiz.
func main() {
	var s []int
	eski := cap(s)
	for i := range 2560 {
		s = append(s, i)
		if cap(s) != eski {
			if eski >= 128 {
				fmt.Printf("%5d → %5d  ×%.2f\n", eski, cap(s), float64(cap(s))/float64(eski))
			}
			eski = cap(s)
		}
	}
}
