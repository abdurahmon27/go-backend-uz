// Oʻrta: iterator'ni oladigan va iterator qaytaradigan generic adapterlar (GO-1.17).
package orta

import "iter"

// Filtr faqat f true qaytargan qiymatlarni beradi.
func Filtr[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if f(v) && !yield(v) {
				return
			}
		}
	}
}

// Aylantir har bir qiymatni f orqali boshqa turga aylantiradi.
func Aylantir[V, W any](seq iter.Seq[V], f func(V) W) iter.Seq[W] {
	return func(yield func(W) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}
