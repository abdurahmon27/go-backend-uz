package main

// Boshqa goroutine tirik boʻlsa (masalan, server), runtime deadlock'ni aniqlamaydi: dastur abadiy osilib qoladi.
//   timeout 3 ./_video/run.sh 12-osilish go run .
import (
	"fmt"
	"time"
)

func main() {
	go func() {
		for {
			time.Sleep(time.Second) // tirik goroutine: masalan, ulanishlarni kutayotgan server
		}
	}()
	prices := make(chan int64)
	prices <- 12_000 // hech kim qabul qilmaydi
	fmt.Println("bu qator hech qachon chiqmaydi")
}
