// xarita: har bir maydon struct ichida qayerda turadi (Offsetof) va
// orada nechta boʻsh bayt (padding) bor.
package main

import (
	"fmt"
	"strings"
	"unsafe"
)

type Buyurtma struct {
	Tolandi    bool
	Summa      int64
	Yetkazildi bool
}

type Ixcham struct {
	Summa      int64
	Tolandi    bool
	Yetkazildi bool
}

type maydon struct {
	harf           string
	offset, olcham uintptr
}

// chiz: har bir bayt uchun bitta belgi; maydon harfi yoki "·" (padding).
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
	fmt.Printf("%-9s %2d bayt, align %d  [%s]\n", nom, jami, align, strings.Join(b, ""))
}

func main() {
	var x Buyurtma
	chiz("Buyurtma", unsafe.Sizeof(x), unsafe.Alignof(x),
		maydon{"T", unsafe.Offsetof(x.Tolandi), unsafe.Sizeof(x.Tolandi)},
		maydon{"S", unsafe.Offsetof(x.Summa), unsafe.Sizeof(x.Summa)},
		maydon{"Y", unsafe.Offsetof(x.Yetkazildi), unsafe.Sizeof(x.Yetkazildi)})

	var y Ixcham
	chiz("Ixcham", unsafe.Sizeof(y), unsafe.Alignof(y),
		maydon{"S", unsafe.Offsetof(y.Summa), unsafe.Sizeof(y.Summa)},
		maydon{"T", unsafe.Offsetof(y.Tolandi), unsafe.Sizeof(y.Tolandi)},
		maydon{"Y", unsafe.Offsetof(y.Yetkazildi), unsafe.Sizeof(y.Yetkazildi)})
}
