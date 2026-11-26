package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	time.Sleep(2) // 2 soniya? 2 millisekund?
	fmt.Println(time.Duration(2), time.Since(start) < time.Millisecond)
}
