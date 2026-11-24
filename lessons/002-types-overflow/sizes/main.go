package main

import (
	"fmt"
	"math"
	"unsafe"
)

func main() {
	fmt.Println("int8 ", unsafe.Sizeof(int8(0)), "bayt:", math.MinInt8, "...", math.MaxInt8)
	fmt.Println("int16", unsafe.Sizeof(int16(0)), "bayt:", math.MinInt16, "...", math.MaxInt16)
	fmt.Println("int32", unsafe.Sizeof(int32(0)), "bayt:", math.MinInt32, "...", math.MaxInt32)
	fmt.Println("int64", unsafe.Sizeof(int64(0)), "bayt:", math.MinInt64, "...", math.MaxInt64)
	fmt.Println("uint8", unsafe.Sizeof(uint8(0)), "bayt:", 0, "...", math.MaxUint8)
	fmt.Println("int  ", unsafe.Sizeof(int(0)), "bayt (bu kompyuterda)")
}
