// main'dagi recover boshqa goroutine'dagi panic'ni ushlay olmaydi.
package main

import (
	"fmt"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("qutqarildi:", r)
		}
	}()

	go func() {
		var narxlar map[string]int
		narxlar["olma"] = 1200 // nil map'ga yozish: panic
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("main tugadi")
}
