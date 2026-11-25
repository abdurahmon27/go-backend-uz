// hook: Go 1.27 dagi generic metod. Xuddi shu fayl ../eski da (go 1.26) kompilyatsiya boʻlmaydi.
// go run ./hook
package main

import "fmt"

type Royxat[T any] []T

// Aylantir har bir elementga f ni qoʻllaydi.
// U — metodning oʻz type parametri: Go 1.27 dan beri mumkin.
func (r Royxat[T]) Aylantir[U any](f func(T) U) Royxat[U] {
	natija := make(Royxat[U], len(r))
	for i, x := range r {
		natija[i] = f(x)
	}
	return natija
}

func main() {
	summalar := Royxat[int64]{150000, 99000, 42000}
	matn := summalar.Aylantir(func(s int64) string {
		return fmt.Sprintf("%d soʻm", s)
	})
	fmt.Println(matn)
}
