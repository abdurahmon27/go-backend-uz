package main

import "fmt"

func main() {
	runes := []rune("paxta yogʻi")
	fmt.Println(len(runes), string(runes[:10]))
}
