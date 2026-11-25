// royxat: Go 1.18–1.26 usuli. Metod faqat receiver'dagi T ni ishlata oladi,
// yangi tur (U) kerak boʻlsa — paket darajasidagi funksiya.
// go run ./royxat
package main

import "fmt"

type Royxat[T any] []T

// Filtr: faqat T kerak — metod boʻla oladi (Go 1.18 dan beri).
func (r Royxat[T]) Filtr(ok func(T) bool) Royxat[T] {
	var natija Royxat[T]
	for _, x := range r {
		if ok(x) {
			natija = append(natija, x)
		}
	}
	return natija
}

// Aylantir: yangi U kerak — 1.27 gacha faqat funksiya.
func Aylantir[T, U any](r Royxat[T], f func(T) U) Royxat[U] {
	natija := make(Royxat[U], len(r))
	for i, x := range r {
		natija[i] = f(x)
	}
	return natija
}

func main() {
	r := Royxat[int64]{150000, 99000, 42000}
	fmt.Println(Aylantir(r.Filtr(katta), somda))
}

func katta(s int64) bool   { return s > 50000 }
func somda(s int64) string { return fmt.Sprintf("%d soʻm", s) }
