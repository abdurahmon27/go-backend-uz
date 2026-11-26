package main

import "fmt"

func main() {
	var x int8 = 127
	x += 3 // 127 → −128 → −127 → −126: jim aylanadi (02-dars)
	fmt.Println(x)
}
