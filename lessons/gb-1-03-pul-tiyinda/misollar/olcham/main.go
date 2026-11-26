// olcham: butun son turlarining shu kompyuterdagi oʻlchami va chegarasi.
// math va unsafe paketlari darsda koʻrsatilmaydi: bu yerda faqat natija muhim.
package main

import (
	"fmt"
	"math"
	"unsafe"
)

func main() {
	var a int32
	var b int64
	var c int
	fmt.Println("int32:", unsafe.Sizeof(a), "bayt,", int64(math.MinInt32), "…", int64(math.MaxInt32))
	fmt.Println("int64:", unsafe.Sizeof(b), "bayt,", int64(math.MinInt64), "…", int64(math.MaxInt64))
	fmt.Println("int:  ", unsafe.Sizeof(c), "bayt (bu kompyuterda)")
}
