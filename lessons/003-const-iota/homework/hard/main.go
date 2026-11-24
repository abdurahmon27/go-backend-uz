package main

import (
	"fmt"
	"time"
)

func main() {
	// fmt.Println(1.5 * time.Second)
	// xato: 1.5 (untyped float constant) truncated to int64
	fmt.Println(1500 * time.Millisecond)
	fmt.Println(time.Second * 3 / 2)
}
