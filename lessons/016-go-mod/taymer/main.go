package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	t := time.NewTimer(time.Millisecond)
	time.Sleep(20 * time.Millisecond) // taymer allaqachon ishladi
	t.Reset(time.Hour)                // endi bir soat kutishi kerak

	select {
	case <-t.C:
		fmt.Println(runtime.Version(), "→ eski signal keldi")
	default:
		fmt.Println(runtime.Version(), "→ kanal boʻsh")
	}
}
