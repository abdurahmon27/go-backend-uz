package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	time.Sleep(2) // 2 sekund kutadi deb o'ylaysizmi?
	fmt.Println(time.Since(start) < time.Second)

	var d time.Duration = 2
	fmt.Println(d)
}
