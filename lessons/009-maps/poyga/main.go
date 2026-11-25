package main

import (
	"fmt"
	"sync"
)

func main() {
	m := map[int]int{}
	var wg sync.WaitGroup
	for g := range 2 {
		wg.Go(func() { yoz(m, g) }) // ataylab xato: ikki goroutine bitta map'ga yozadi
	}
	wg.Wait()
	fmt.Println("dastur davom etdi:", len(m))
}

func yoz(m map[int]int, g int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("qutqarildi:", r)
		}
	}()
	for i := range 1_000_000 {
		m[i] = g
	}
}
