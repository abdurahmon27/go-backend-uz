package main

import (
	"fmt"
	"math"
)

func less(x int32) bool {
	return x < x+1
}

func main() {
	fmt.Println(less(5))
	fmt.Println(less(math.MaxInt32))
}
