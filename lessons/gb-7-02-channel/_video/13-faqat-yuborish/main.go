package main

// chan<- Product: faqat yuborish uchun; undan oʻqish kompilyatsiya xatosi (go build).
func parse(lines []string, out chan<- Product) {
	p := <-out
	_ = p
}

func main() {}
