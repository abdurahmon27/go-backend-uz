package main

import (
	"bytes"
	"fmt"
	"runtime"
	"slices"
)

// sarlavha katta "fayl"dan birinchi qatorni qaytaradi. Bo'lak butun massivga qaraydi:
// bo'lak tirik ekan, 64 MB ham tirik qoladi (go.dev/blog/slices-intro, FindDigits).
func sarlavha(fayl []byte) []byte {
	i := bytes.IndexByte(fayl, '\n')
	return fayl[:i]
}

// sarlavhaNusxa faqat kerakli baytlarni nusxalaydi.
func sarlavhaNusxa(fayl []byte) []byte {
	return slices.Clone(sarlavha(fayl))
}

func yangiFayl() []byte {
	f := make([]byte, 64<<20) // 64 MB
	copy(f, "id,nom,narx\n")
	return f
}

func heapMB() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc >> 20
}

func main() {
	a := sarlavha(yangiFayl())
	fmt.Printf("sarlavha:      %q, heap ≈ %d MB\n", a, heapMB())
	runtime.KeepAlive(a)
	a = nil

	b := sarlavhaNusxa(yangiFayl())
	fmt.Printf("sarlavhaNusxa: %q, heap ≈ %d MB\n", b, heapMB())
	runtime.KeepAlive(b)
}
