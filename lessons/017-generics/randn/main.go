// randn: standart kutubxonadagi birinchi generic metod — (*rand.Rand).N (Go 1.27).
// go run ./randn
package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func main() {
	r := rand.New(rand.NewPCG(1, 2))  // doimiy seed: natija har safar bir xil
	fmt.Println(r.N(10))              // int
	fmt.Println(r.N(int64(150000)))   // int64
	fmt.Println(r.N(3 * time.Second)) // time.Duration
	fmt.Printf("%T\n", r.N[uint8])    // aniq tur bilan instantiate
}
