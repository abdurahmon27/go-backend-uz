package main

import "fmt"

func kun(d int) string {
	s := ""
	switch d {
	case 6, 7: // bir nechta qiymat
		s = "dam olish"
	case 1:
		s = "hafta boshi, "
		fallthrough // keyingi case ham bajariladi
	case 2, 3, 4, 5:
		s += "ish kuni"
	default:
		s = "bunday kun yoʻq"
	}
	return s
}

func main() {
	for _, d := range []int{1, 3, 7, 9} {
		fmt.Println(d, kun(d))
	}
}
