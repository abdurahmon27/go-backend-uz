package main

import "fmt"

func main() {
	r := 'ʻ'
	fmt.Printf("%U = %d = %b (%d bit)\n", r, r, r, len(fmt.Sprintf("%b", r)))
	b := []byte(string(r))
	// 2 baytli qolip: 110xxxxx 10xxxxxx, 11 ta x ga raqamning bitlari
	fmt.Printf("110|%05b  10|%06b  = % x\n", r>>6, r&0x3f, b)
}
