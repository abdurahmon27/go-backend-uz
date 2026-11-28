package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {
	for range 3 {
		for range 100 {
			fastestPrice(1)
		}
		time.Sleep(time.Second)
		fmt.Println("goroutine'lar:", runtime.NumGoroutine())
	}
	pprof.Lookup("goroutineleak").WriteTo(os.Stdout, 1)
}
