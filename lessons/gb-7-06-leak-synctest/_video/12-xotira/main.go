// Videoda yoʻq (research uchun): leak boʻlgan goroutine'lar xotirani qancha egallaydi. Kechikishlar 1/2/3 ms.
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	var m runtime.MemStats
	for _, n := range []int{0, 1000, 10000} {
		for range n {
			fastestPrice(1)
		}
		time.Sleep(100 * time.Millisecond)
		runtime.GC()
		runtime.ReadMemStats(&m)
		fmt.Println(runtime.NumGoroutine(), "goroutine; StackInuse KB:", m.StackInuse/1024, "HeapInuse KB:", m.HeapInuse/1024)
	}
}
