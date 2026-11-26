package main

import (
	"fmt"
	"time"
)

func main() {
	var kutish time.Duration = 2 // vaqt turi: nanosekundlarni sanaydi
	fmt.Println(kutish)          // Println birlik qoʻshadi: ns = nanosekund
	time.Sleep(kutish)           // = time.Sleep(2): deyarli kutmaydi
	kutish = 2 * time.Second     // time.Second = 1 soniya = 1 000 000 000 nanosekund
	fmt.Println(kutish)          // s = soniya
	time.Sleep(kutish)           // 2 soniya kutadi
}
