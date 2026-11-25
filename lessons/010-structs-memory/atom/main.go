// atom: 32-bitli tizimda int64 maydon 8 ga karrali manzilda turmasligi mumkin.
// GOARCH=386 go run ./atom — ataylab xato (panic).
package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type Hisoblagich struct {
	Yoniq bool
	Soni  int64 // bir nechta goroutine atomik oshiradi
}

func main() {
	h := &Hisoblagich{}
	fmt.Println("Soni offset:", unsafe.Offsetof(h.Soni))
	atomic.AddInt64(&h.Soni, 1)
	fmt.Println("Soni:", h.Soni)
}
