package main

import (
	"fmt"
	"time"
)

func main() {
	time.Sleep(2) // untyped 2 → time.Duration: 2 nanosekund soʻraydi
	var d time.Duration = 2
	fmt.Println(d)      // Duration birligi bilan chiqadi: 2ns
	d = 2 * time.Second // time.Second: typed constant, turi Duration (1 000 000 000 ns)
	fmt.Println(d)      // 2s
	time.Sleep(d)       // 2 soniya kutadi
}
