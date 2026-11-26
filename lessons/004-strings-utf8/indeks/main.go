package main

import "fmt"

func main() {
	s := "oʻzbek"
	fmt.Println(s[1])                 // JS va Python'da s[1] → "ʻ"; Go'da: 1-oʻrindagi bayt (byte, son)
	fmt.Printf("%x\n", s[1])          // oʻsha bayt oʻn oltilikda: ca (ʻ ning birinchi bayti)
	fmt.Println(string([]rune(s)[1])) // rune'lar boʻyicha: 1-element — ʻ
}
