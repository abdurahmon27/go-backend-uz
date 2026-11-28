package main

import (
	"fmt"
	"sync"
	"time"
)

var offers = map[string]int{"chorsu": 40, "oloy": 120, "siyob": 5}

func fetch(name string) int {
	fmt.Println("boshlandi:", name)
	time.Sleep(time.Second)
	fmt.Println("tayyor:", name)
	return offers[name]
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	suppliers := []string{"chorsu", "oloy", "siyob"}
	counts := make([]int, len(suppliers))
	for i, name := range suppliers {
		wg.Go(func() { counts[i] = fetch(name) })
	}
	wg.Wait()
	fmt.Println("mahsulotlar:", counts)
	fmt.Println("jami:", time.Since(start))
}
