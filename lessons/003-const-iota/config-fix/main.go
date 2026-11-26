package main

import (
	"fmt"
	"time"
)

func main() {
	n := 2
	d := time.Duration(n) * time.Second // konversiya: int → time.Duration
	fmt.Println(d)
	time.Sleep(d) // 2 soniya kutadi
}
