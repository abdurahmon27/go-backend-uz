package main

import (
	"fmt"
	"sync"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ushlandi:", r)
		}
	}()
	var wg sync.WaitGroup
	wg.Go(func() { panic("chorsu: javob buzuq") })
	wg.Wait()
	fmt.Println("tugadi")
}
