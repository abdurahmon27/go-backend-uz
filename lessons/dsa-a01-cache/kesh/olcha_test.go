package kesh

import (
	"fmt"
	"runtime"
	"testing"
)

// 2^10 dan 2^26 gacha element: 8 KiB dan 512 MiB gacha.
func olcha(b *testing.B, tayyorla func(n int) func() int64) {
	for k := 10; k <= 26; k++ {
		n := 1 << k
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			ish := tayyorla(n) // tayyorlash o'lchanmaydi: b.Loop birinchi chaqiruvda taymerni nolga qaytaradi
			runtime.GC()       // tayyorlashdan qolgan axlat o'lchov paytida yig'ilmasin
			for b.Loop() {
				ish() // b.Loop natijani o'zi "tirik" ushlaydi (GO-1.07)
			}
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(n), "ns/el")
		})
	}
}
