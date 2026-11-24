package main

import (
	"fmt"
	"strings"
)

// ataylab xato: funksiyalarni faqat nil bilan solishtirish mumkin
func main() {
	f := strings.ToUpper
	g := strings.ToUpper
	fmt.Println(f == g)
}
