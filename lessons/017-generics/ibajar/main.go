// ibajar: ATAYLAB kompilyatsiya boʻlmaydi. Generic metod interface metodini bajarmaydi.
// go build ./ibajar
package main

import "fmt"

type Bosh struct{}

func (Bosh) Nop[T any](x T) T { return x }

type MatnNop interface {
	Nop(string) string
}

func main() {
	var b Bosh
	fmt.Println(b.Nop("salom"), b.Nop(42)) // toʻgʻridan-toʻgʻri: ishlaydi
	var _ MatnNop = b                      // interface orqali: yoʻq
}
