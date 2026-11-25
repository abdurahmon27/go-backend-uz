// guru: oʻlchami kompilyatsiya paytida nomaʼlum slice. Go 1.25+ kompilyatori
// unga stack'da 32 baytlik joy ajratadi; sigʻmasa — heap.
// []int64 elementi 8 bayt: nechta element stack'ga sigʻadi?
// go test -bench . -benchmem ./homework/guru
package guru

import "testing"

//go:noinline
func ishlat(n int) int {
	s := make([]int64, n) // -gcflags=-m: make([]int64, n) does not escape
	s[0] = 1
	return len(s)
}

func BenchmarkMake4(b *testing.B) {
	for b.Loop() {
		ishlat(4)
	}
}

func BenchmarkMake5(b *testing.B) {
	for b.Loop() {
		ishlat(5)
	}
}
