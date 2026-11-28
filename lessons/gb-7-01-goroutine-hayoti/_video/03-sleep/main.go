package main

import (
	"fmt"
	"time"
)

func fetch(name string) {
	fmt.Println("boshlandi:", name)
	time.Sleep(time.Second)
	fmt.Println("tayyor:", name)
}

func main() {
	start := time.Now()
	for _, name := range []string{"chorsu", "oloy", "siyob"} {
		go fetch(name)
	}
	time.Sleep(500 * time.Millisecond)
	fmt.Println("jami:", time.Since(start))
}
