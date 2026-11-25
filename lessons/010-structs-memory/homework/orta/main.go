// orta: string ham maydon. string = pointer + uzunlik: 16 bayt, align 8.
//
//	Faol  bool   -> 0      (7 bayt padding)
//	Ism   string -> 8..23
//	Admin bool   -> 24     (3 bayt padding)
//	Yosh  int32  -> 28..31
//	jami 32
//
// Ixcham: Ism, Yosh, Faol, Admin = 16+4+1+1 = 22 -> 24.
package main

import (
	"fmt"
	"unsafe"
)

type Foydalanuvchi struct {
	Faol  bool
	Ism   string
	Admin bool
	Yosh  int32
}

type FoydalanuvchiIxcham struct {
	Ism   string
	Yosh  int32
	Faol  bool
	Admin bool
}

func main() {
	var f Foydalanuvchi
	fmt.Println("string:", unsafe.Sizeof(f.Ism), "bayt, align", unsafe.Alignof(f.Ism))
	fmt.Println("Foydalanuvchi:      ", unsafe.Sizeof(f))                     // 32
	fmt.Println("FoydalanuvchiIxcham:", unsafe.Sizeof(FoydalanuvchiIxcham{})) // 24
}
