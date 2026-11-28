package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	n := <-supplier(3)
	fmt.Println("qoldiq:", n, time.Since(start).Round(time.Second))
}
