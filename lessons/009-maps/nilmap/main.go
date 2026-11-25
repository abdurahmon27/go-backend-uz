package main

import "fmt"

func main() {
	var narx map[string]int // nil map: make qilinmagan
	fmt.Println(narx["non"], len(narx))
	yoz(narx)
	fmt.Println("dastur davom etdi")
}

func yoz(m map[string]int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("qutqarildi:", r)
		}
	}()
	m["non"] = 5000 // ataylab xato: nil map'ga yozib boʻlmaydi
}
