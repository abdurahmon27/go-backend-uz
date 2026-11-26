// Guru: ikki iterator'ni yonma-yon yurgizish. range bunga yetmaydi — iter.Pull kerak.
package guru

import "iter"

// Zip a va b dan navbatma-navbat bittadan qiymat oladi; qisqasi tugaganda toʻxtaydi.
func Zip[A, B any](a iter.Seq[A], b iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		nextA, stopA := iter.Pull(a)
		defer stopA()
		nextB, stopB := iter.Pull(b)
		defer stopB()
		for {
			va, ok1 := nextA()
			vb, ok2 := nextB()
			if !ok1 || !ok2 || !yield(va, vb) {
				return
			}
		}
	}
}
