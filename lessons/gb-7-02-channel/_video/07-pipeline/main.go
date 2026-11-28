package main

import "fmt"

func parse(lines []string, out chan<- Product) {
	defer close(out)
	for _, line := range lines {
		p, err := parseLine(line)
		if err != nil {
			continue // xatoni qaytarish: GB-7.04
		}
		out <- p
	}
}

func inStock(in <-chan Product, out chan<- Product) {
	defer close(out)
	for p := range in {
		if p.Stock > 0 {
			out <- p
		}
	}
}

func main() {
	raw := make(chan Product)
	ready := make(chan Product)
	go parse(chorsu, raw)
	go inStock(raw, ready)
	for p := range ready {
		fmt.Println(p.Name, p.Stock)
	}
}
