package main

import (
	"fmt"
	"math"
)

func main() {
	// faraz: hisoblagich int32 va har 10 millisekundda 1 ga oshadi
	const perSecond = 100
	ticks := float64(math.MaxInt32) + 1 // 2^31: shu qadamda to'ladi
	fmt.Printf("%.2f kun\n", ticks/perSecond/(24*60*60))

	var counter int32 = math.MaxInt32
	counter++
	fmt.Println(counter)
}
