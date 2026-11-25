package main

import (
	"math/rand/v2"
	"testing"
)

// TestOlchov: tasodifiy (base62) kalitlarning necha foizini «nomaʼlum-kalit» qoidasi topadi.
// go test -run Olchov -v   (seed qatʼiy: natija har safar bir xil)
func TestOlchov(t *testing.T) {
	const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	r := rand.New(rand.NewPCG(1, 2))
	for _, n := range []int{24, 32, 40} {
		top := 0
		for i := 0; i < 10000; i++ {
			b := make([]byte, n)
			for j := range b {
				b[j] = abc[r.IntN(len(abc))]
			}
			if len(qatorTekshir("x", "k="+string(b))) > 0 {
				top++
			}
		}
		t.Logf("%d belgi: %d / 10000 topildi", n, top)
	}
}
