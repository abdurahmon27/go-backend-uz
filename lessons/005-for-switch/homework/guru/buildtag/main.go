//go:build go1.22

// 3-usul: modul go 1.21 da qoladi, lekin shu fayl Go 1.22 semantikasi bilan
// kompilyatsiya qilinadi (fayl darajasidagi //go:build go1.22 qatori).
package main

import "fmt"

func main() {
	domenlar := []string{"a.uz", "b.uz", "c.uz"}
	var navbat []*string
	for _, d := range domenlar {
		navbat = append(navbat, &d)
	}
	for _, p := range navbat {
		fmt.Print(*p, " ")
	}
	fmt.Println("← //go:build go1.22")
}
