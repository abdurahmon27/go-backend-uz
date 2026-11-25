// defer argumentlari defer qatorida hisoblanadi, closure ichidagisi esa oxirida.
package main

import (
	"fmt"
	"time"
)

func ish() {
	start := time.Now()
	defer fmt.Println("argument:", time.Since(start).Round(time.Millisecond))
	defer func() {
		fmt.Println("closure: ", time.Since(start).Round(time.Millisecond))
	}()

	time.Sleep(100 * time.Millisecond) // og'ir ish
}

func main() {
	ish()
}
