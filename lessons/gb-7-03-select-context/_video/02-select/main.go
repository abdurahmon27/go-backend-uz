package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	select {
	case n := <-supplier(3):
		fmt.Println("qoldiq:", n)
	case <-time.After(time.Second):
		fmt.Println("supplier javob bermadi")
	}
	fmt.Println(time.Since(start).Round(time.Second))
}
