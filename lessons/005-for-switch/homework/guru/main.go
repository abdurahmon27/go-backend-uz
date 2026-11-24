package main

import "fmt"

// go.mod: go 1.21 (eski semantika), lekin natija toʻgʻri.
func main() {
	domenlar := []string{"a.uz", "b.uz", "c.uz"}

	// 1-usul: har qadamda yangi oʻzgaruvchi (Go 1.22 gacha shunday yozilardi)
	var navbat []*string
	for _, d := range domenlar {
		d := d
		navbat = append(navbat, &d)
	}
	for _, p := range navbat {
		fmt.Print(*p, " ")
	}
	fmt.Println("← d := d")

	// 2-usul: sikl oʻzgaruvchisi emas, slice elementining manzili
	navbat = nil
	for i := range domenlar {
		navbat = append(navbat, &domenlar[i])
	}
	for _, p := range navbat {
		fmt.Print(*p, " ")
	}
	fmt.Println("← &domenlar[i]")
}
