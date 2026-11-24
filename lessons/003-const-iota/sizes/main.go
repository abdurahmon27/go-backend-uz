package main

import "fmt"

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
	TB
)

func main() {
	fmt.Println(KB, MB, GB, TB)
}
