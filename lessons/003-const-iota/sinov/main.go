package main

import (
	"fmt"
	"time"
)

func main() {
	const a = 3     // untyped
	const b int = 3 // typed: turi int
	x := a * time.Second
	y := b * time.Second // ataylab xato: int × time.Duration
	fmt.Println(x, y)
}
