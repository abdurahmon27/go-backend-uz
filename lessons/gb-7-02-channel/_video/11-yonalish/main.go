package main

// Yoʻnalishli channel turi: kompilyator xatosi (go vet / go build).
func inStock(in <-chan Product, out chan<- Product) {
	defer close(out)
	for p := range in {
		if p.Stock > 0 {
			out <- p
		}
	}
	close(in)
}

func main() {}
