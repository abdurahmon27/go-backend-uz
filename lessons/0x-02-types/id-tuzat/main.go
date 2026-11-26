package main

import (
	"fmt"
	"strconv"
)

func main() {
	id := 65
	key := "user:" + strconv.Itoa(id)
	fmt.Println(key)
}
