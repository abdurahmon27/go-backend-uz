// Package qiyin: Guruhla — Go 1.27 generic metodi. K — metodning oʻz type parametri.
package qiyin

type Royxat[T any] []T

// Guruhla elementlarni kalit(x) boʻyicha guruhlaydi; har guruhda tartib saqlanadi.
func (r Royxat[T]) Guruhla[K comparable](kalit func(T) K) map[K]Royxat[T] {
	natija := make(map[K]Royxat[T])
	for _, x := range r {
		k := kalit(x)
		natija[k] = append(natija[k], x)
	}
	return natija
}
