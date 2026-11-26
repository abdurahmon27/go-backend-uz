package main

import (
	"fmt"
	"time"
)

func main() {
	n := 2
	fmt.Println(time.Duration(n) * time.Second) // 2s — toʻgʻri
	fmt.Println(time.Duration(n))               // 2ns — notoʻgʻri tuzatish: birlik yoʻq
}
