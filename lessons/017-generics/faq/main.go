// faq: Go FAQ'dagi misol (generic_methods), Go 1.27 da ishga tushirilgan.
// Bosh qiymati any ichida: dastur ishlayotganda uning generic metodi koʻrinmaydi.
// go run ./faq
package main

import (
	"fmt"
	"reflect"
)

type Bosh struct{}

func (Bosh) Nop[T any](x T) T { return x }
func (Bosh) Oddiy()           {}

func main() {
	var x any = Bosh{}
	_, ok := x.(interface{ Nop(string) string })
	fmt.Println("Nop(string) string bormi?", ok)
	_, ok = x.(interface{ Nop(int) int })
	fmt.Println("Nop(int) int bormi?", ok)
	_, ok = x.(interface{ Oddiy() })
	fmt.Println("Oddiy() bormi?", ok)
	fmt.Println("reflect koʻrgan metodlar:", reflect.TypeOf(x).NumMethod())
}
