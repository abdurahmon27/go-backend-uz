package main

import (
	"fmt"
	"strings"
)

// Qiyin: kalitlar soni 8 dan oshsa, nechta tartib chiqadi?
func main() {
	for _, n := range []int{5, 8, 9, 16} {
		m := map[int]bool{}
		for i := range n {
			m[i] = true
		}
		tartiblar := map[string]bool{}
		for range 100_000 {
			var sb strings.Builder
			for k := range m {
				fmt.Fprint(&sb, k, " ")
			}
			tartiblar[sb.String()] = true
		}
		fmt.Printf("%2d ta kalit: %d xil tartib\n", n, len(tartiblar))
	}
}
