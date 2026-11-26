package main

import (
	"fmt"
	"math"
)

func main() {
	var id int64 = math.MaxInt64 // 9223372036854775807
	id++
	fmt.Println(id)
}
