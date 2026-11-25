// qiyin: struct{} 0 bayt. Lekin oxirgi maydon boʻlsa, struct kattalashadi.
//
// Agar oxirdagi 0 baytli maydonga padding qoʻshilmasa, &x.B struct'dan
// keyingi baytga — boshqa obyektga — ishora qilardi. Shunda garbage
// collector u obyektni ham "tirik" deb oʻylashi mumkin edi. Shuning uchun
// kompilyator oxiridagi 0 baytli maydon uchun 1 bayt qoʻshadi, keyin
// struct oʻlchami 8 ga karrali qilib yaxlitlanadi: 8 + 1 -> 16.
// Manba: golang/go#9401, cmd/compile/internal/types/size.go.
package main

import (
	"fmt"
	"unsafe"
)

type Oxirida struct {
	A int64
	B struct{}
}

type Boshida struct {
	B struct{}
	A int64
}

func main() {
	fmt.Println("struct{}:", unsafe.Sizeof(struct{}{})) // 0
	fmt.Println("Oxirida: ", unsafe.Sizeof(Oxirida{}))  // 16
	fmt.Println("Boshida: ", unsafe.Sizeof(Boshida{}))  // 8
}
