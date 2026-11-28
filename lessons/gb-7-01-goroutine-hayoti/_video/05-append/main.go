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
	var counts []int
	for _, name := range []string{"chorsu", "oloy", "siyob"} {
		wg.Go(func() { counts = append(counts, fetch(name)) })
	}
	wg.Wait()
	fmt.Println("mahsulotlar:", counts)
	fmt.Println("jami:", time.Since(start))
}
