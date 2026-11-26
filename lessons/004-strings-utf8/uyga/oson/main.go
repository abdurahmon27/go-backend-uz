package main

import (
	"fmt"
	"unicode/utf8"
)

// 🟢 Oson: ism oʻrniga oʻz ismingizni yozing.
// Ishga tushirishdan OLDIN ikkala sonni qogʻozga yozib qoʻying, keyin tekshiring.
// Keyin ʻ ni klaviaturadagi ' ga almashtiring: qaysi son oʻzgaradi?
func main() {
	ism := "Gʻayrat"
	fmt.Println(len(ism))
	fmt.Println(utf8.RuneCountInString(ism))
}
