// savol: ikkalasida ham nol. Ular tengmi?
// go run ./savol
package main

import "fmt"

func main() {
	var a any = 0        // int
	var b any = int64(0) // int64
	fmt.Println("a == b:", a == b)
	fmt.Printf("a: %T, b: %T\n", a, b)
}
