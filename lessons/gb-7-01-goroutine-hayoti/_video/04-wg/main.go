package main

import (
	"fmt"
	"sync"
	"time"
)

func fetch(name string) {
	fmt.Println("boshlandi:", name)
	time.Sleep(time.Second)
	fmt.Println("tayyor:", name)
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	for _, name := range []string{"chorsu", "oloy", "siyob"} {
		wg.Go(func() { fetch(name) })
	}
	wg.Wait()
	fmt.Println("jami:", time.Since(start))
}
