// Nanosekund darajasidagi amallar: ularni bittalab soat bilan o'lchab bo'lmaydi,
// shuning uchun GO-1.07 dagi kabi benchmark sikli ishlatiladi (b.Loop).
//
//	taskset -c 2 go test -run '^$' -bench . -count 10 ./xotira
//
// Kesh va operativ xotiradan bitta o'qish (L1, L2, RAM) bu yerda qayta o'lchanmaydi:
// ular DSA-A.01 darsida o'lchangan (../dsa-a01-cache, Zanjir benchmark'i).
package xotira

import (
	"sync"
	"testing"
	"time"
)

// time.Now() ning o'zi qancha turadi? Shu sababli 1 ns lik amalni bittalab o'lchamaymiz.
func BenchmarkVaqt(b *testing.B) {
	var t time.Time
	for b.Loop() {
		t = time.Now()
	}
	_ = t
}

// Hech kim raqobat qilmayotgan mutex: Lock + Unlock juftligi.
func BenchmarkMutex(b *testing.B) {
	var mu sync.Mutex
	for b.Loop() {
		mu.Lock()
		mu.Unlock()
	}
}

const (
	mb    = 1 << 20    // 1 MiB
	katta = 512 << 20  // 512 MiB: L3 (12 MiB) dan ancha katta
	el    = mb / 8     // 1 MiB dagi int64 soni
	nKat  = katta / mb // nechta 1 MiB lik bo'lak
)

var yigindi int64

// 1 MiB ni ketma-ket o'qish (qo'shish), har safar keyingi bo'lakdan:
// 512 MiB aylanib chiqquncha bitta bo'lak keshda qolmaydi, demak o'qish operativ xotiradan.
func BenchmarkMegabayt(b *testing.B) {
	s := make([]int64, katta/8)
	for i := range s {
		s[i] = int64(i)
	}
	var y int64
	k := 0
	for b.Loop() {
		bolak := s[k*el : (k+1)*el]
		for _, v := range bolak {
			y += v
		}
		k = (k + 1) % nKat
	}
	yigindi = y
}
