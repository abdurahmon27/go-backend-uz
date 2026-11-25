// oson: avval qogʻozda hisoblang, keyin tekshiring.
//
//	A bool  -> 0
//	(7 bayt padding: B 8 ga karrali manzilda boshlanadi)
//	B int64 -> 8..15
//	C int32 -> 16..19
//	D bool  -> 20
//	jami 21, 8 ga karralisi -> 24
//
// Ixcham tartib: B, C, A, D = 8+4+1+1 = 14 -> 16.
package main

import (
	"fmt"
	"unsafe"
)

type Oson struct {
	A bool
	B int64
	C int32
	D bool
}

type OsonIxcham struct {
	B int64
	C int32
	A bool
	D bool
}

func main() {
	fmt.Println("Oson:      ", unsafe.Sizeof(Oson{}))       // 24
	fmt.Println("OsonIxcham:", unsafe.Sizeof(OsonIxcham{})) // 16
}
