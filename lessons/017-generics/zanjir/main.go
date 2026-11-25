// zanjir: Go 1.27 usuli. Aylantir endi metod — chaqiruvlar chapdan oʻngga oʻqiladi.
// go run ./zanjir
package main

import "fmt"

type Royxat[T any] []T

func (r Royxat[T]) Filtr(ok func(T) bool) Royxat[T] {
	var natija Royxat[T]
	for _, x := range r {
		if ok(x) {
			natija = append(natija, x)
		}
	}
	return natija
}

// U — metodning oʻz type parametri (Go 1.27).
func (r Royxat[T]) Aylantir[U any](f func(T) U) Royxat[U] {
	natija := make(Royxat[U], len(r))
	for i, x := range r {
		natija[i] = f(x)
	}
	return natija
}

func main() {
	r := Royxat[int64]{150000, 99000, 42000}
	fmt.Println(r.Filtr(katta).Aylantir(somda))
}

func katta(s int64) bool   { return s > 50000 }
func somda(s int64) string { return fmt.Sprintf("%d soʻm", s) }
