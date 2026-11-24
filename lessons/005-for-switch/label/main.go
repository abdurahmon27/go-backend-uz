package main

import "fmt"

func main() {
	sonlar := []int{4, 7, -1, 8, -5}
qidiruv:
	for _, s := range sonlar {
		switch {
		case s < 0:
			fmt.Println("manfiy:", s, "→ toʻxtaymiz")
			break qidiruv // for'dan chiqadi
		case s%2 == 0:
			fmt.Println("juft:", s)
		}
	}
	fmt.Println("tugadi")
}
