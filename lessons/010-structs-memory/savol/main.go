// savol: pauzadagi savolning javobi. bool, int32, bool: necha bayt?
package main

import (
	"fmt"
	"strings"
	"unsafe"
)

type Savol struct {
	A bool
	B int32
	C bool
}

type Javob struct {
	B int32
	A bool
	C bool
}

type maydon struct {
	harf           string
	offset, olcham uintptr
}

func chiz(nom string, jami, align uintptr, mm ...maydon) {
	b := []string{}
	for range jami {
		b = append(b, "·")
	}
	for _, m := range mm {
		for i := m.offset; i < m.offset+m.olcham; i++ {
			b[i] = m.harf
		}
	}
	fmt.Printf("%-6s %2d bayt, align %d  [%s]\n", nom, jami, align, strings.Join(b, ""))
}

func main() {
	var x Savol
	chiz("Savol", unsafe.Sizeof(x), unsafe.Alignof(x),
		maydon{"A", unsafe.Offsetof(x.A), unsafe.Sizeof(x.A)},
		maydon{"B", unsafe.Offsetof(x.B), unsafe.Sizeof(x.B)},
		maydon{"C", unsafe.Offsetof(x.C), unsafe.Sizeof(x.C)})
	var y Javob
	chiz("Javob", unsafe.Sizeof(y), unsafe.Alignof(y),
		maydon{"B", unsafe.Offsetof(y.B), unsafe.Sizeof(y.B)},
		maydon{"A", unsafe.Offsetof(y.A), unsafe.Sizeof(y.A)},
		maydon{"C", unsafe.Offsetof(y.C), unsafe.Sizeof(y.C)})
}
