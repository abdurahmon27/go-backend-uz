// million: bir million buyurtma xotirada qancha joy oladi?
package main

import (
	"fmt"
	"runtime"
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

const n = 1_000_000

// ajratildi: f ishlaganda heap'da necha bayt ajratildi.
func ajratildi(f func()) uint64 {
	var a, b runtime.MemStats
	runtime.ReadMemStats(&a)
	f()
	runtime.ReadMemStats(&b)
	return b.TotalAlloc - a.TotalAlloc
}

var (
	s1 []Buyurtma
	s2 []Ixcham
)

func main() {
	mb := func(b uint64) float64 { return float64(b) / 1e6 }
	fmt.Printf("[]Buyurtma: %.1f MB\n", mb(ajratildi(func() { s1 = make([]Buyurtma, n) })))
	fmt.Printf("[]Ixcham:   %.1f MB\n", mb(ajratildi(func() { s2 = make([]Ixcham, n) })))
}
