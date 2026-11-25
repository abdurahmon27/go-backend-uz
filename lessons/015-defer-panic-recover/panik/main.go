// panic: oddiy ish to'xtaydi, lekin defer'lar baribir ishlaydi.
package main

import "fmt"

func main() {
	defer fmt.Println("defer: baza yopildi")

	narxlar := []int{1200, 3400, 560}
	i := 5
	fmt.Println("narx:", narxlar[i])
	fmt.Println("bu qator chiqmaydi")
}
