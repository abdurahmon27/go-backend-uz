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

func main() {
	products := make(chan Product)
	go parse(chorsu, products)
	for p := range products {
		fmt.Println(p.Name, p.Stock)
	}
	fmt.Println("import tugadi")
}
