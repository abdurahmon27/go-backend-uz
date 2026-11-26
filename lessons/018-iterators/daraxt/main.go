// Daraxt: rekursiv tuzilmani iterator bilan aylanish.
// Slice yaratilmaydi, channel ham, goroutine ham kerak emas.
package main

import (
	"cmp"
	"fmt"
	"iter"
)

type Daraxt[T cmp.Ordered] struct {
	chap, ong *Daraxt[T]
	qiymat    T
}

func (d *Daraxt[T]) Qosh(v T) *Daraxt[T] {
	if d == nil {
		return &Daraxt[T]{qiymat: v}
	}
	if v < d.qiymat {
		d.chap = d.chap.Qosh(v)
	} else {
		d.ong = d.ong.Qosh(v)
	}
	return d
}

// All qiymatlarni oʻsish tartibida beradi.
func (d *Daraxt[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		d.aylan(yield)
	}
}

func (d *Daraxt[T]) aylan(yield func(T) bool) bool {
	if d == nil {
		return true
	}
	return d.chap.aylan(yield) && yield(d.qiymat) && d.ong.aylan(yield)
}

func main() {
	var d *Daraxt[int64]
	for _, s := range []int64{150000, 42000, 310000, 99000, 7000} {
		d = d.Qosh(s)
	}
	for s := range d.All() {
		fmt.Println(s)
		if s >= 100000 {
			break
		}
	}
}
