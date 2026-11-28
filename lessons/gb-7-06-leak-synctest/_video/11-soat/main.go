package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// main 3 × 100 marta narx soʻraydi va har yuztadan keyin tirik goroutine'lar sonini chiqaradi.
// Leak boʻlsa, son oʻsadi (videoda: 201, 401, 601); tuzatilgan fastestPrice bilan doim 1.
func main() {
	for range 3 {
		for range 100 {
			fastestPrice(context.Background(), 1)
		}
		time.Sleep(time.Second)
		fmt.Println("goroutine'lar:", runtime.NumGoroutine())
	}
}
