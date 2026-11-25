// atomfix: atomic.Int64 (Go 1.19+) har qanday tizimda 8 baytga tekislanadi.
package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type Hisoblagich struct {
	Yoniq bool
	Soni  atomic.Int64
}

func main() {
	h := &Hisoblagich{}
	fmt.Println("Soni offset:", unsafe.Offsetof(h.Soni))
	h.Soni.Add(1)
	fmt.Println("Soni:", h.Soni.Load())
}
