package main

import "fmt"

// Vazifa: kiritilgan kod pin bilan teng boʻlmasa "notoʻgʻri kod", teng boʻlsa "xush kelibsiz" deb yozing.
// Faqat != (teng emasmi?) belgisidan foydalaning, == ishlatmang.
func main() {
	pin := 2580
	kiritilgan := 1111
	fmt.Println(pin, kiritilgan) // bu qatorni oʻchirib, oʻrniga if yozing
}
